//go:build windows
// +build windows

package render

import (
	"sync"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

// ModelRenderer (旧 RenderModel) のバッファ初期化処理をまとめたファイルです。
// このファイルでは、頂点情報、法線、ボーンライン・ボーンポイント、選択頂点、カーソル位置、SSBO など
// モデル描画に必要な OpenGL バッファの生成・データ転送処理を行います。

// initializeBuffers は、モデルの頂点バッファおよび関連バッファを初期化します。
func (mr *ModelRenderer) initializeBuffers(factory *mgl.BufferFactory, model *pmx.PmxModel) {
	// ModelDrawer を初期化
	md := mr.ModelDrawer

	// WaitGroupを用いて並列処理を管理
	var wg sync.WaitGroup

	// 頂点データ、法線データ、選択頂点データを一括生成
	var selectedVertexVertices []float32
	wg.Add(1)
	go func() {
		defer wg.Done()
		md.vertices, md.normalVertices, selectedVertexVertices = createAllVertexData(model)
	}()

	// 面データ一括生成
	wg.Add(1)
	go func() {
		defer wg.Done()
		md.faces = createIndexesData(model)
	}()

	// ボーンデバッグ関連バッファを一括で初期化
	var boneLines, bonePoints []float32
	var boneLineFaces, bonePointFaces []uint32
	var boneLineIndexes, bonePointIndexes []int
	wg.Add(1)
	go func() {
		defer wg.Done()
		boneLines, boneLineFaces, boneLineIndexes,
			bonePoints, bonePointFaces, bonePointIndexes = createAllBoneDebugData(model)
	}()

	// WaitGroupの完了を待つ
	wg.Wait()

	// GLオブジェクトの生成は並列化しない ------------

	// メインのモデル頂点バッファ
	mr.bufferHandle = factory.CreateVertexBuffer(gl.Ptr(md.vertices), len(md.vertices))

	// 法線バッファの初期化
	md.normalBufferHandle = factory.CreateVertexBuffer(gl.Ptr(md.normalVertices), len(md.normalVertices))

	// 法線表示用のインデックスバッファ
	normalIndexes := createNormalIndexesData(model)

	md.normalIbo = factory.CreateElementBuffer(gl.Ptr(normalIndexes), len(normalIndexes))

	// ボーンラインバッファの設定
	md.boneLineIndexes = boneLineIndexes
	md.boneLineCount = len(boneLineIndexes)
	if len(boneLines) > 0 {
		md.boneLineBufferHandle = factory.CreateBoneBuffer(gl.Ptr(boneLines), len(boneLines))
	}
	if len(boneLineFaces) > 0 {
		md.boneLineIbo = factory.CreateElementBuffer(gl.Ptr(boneLineFaces), len(boneLineFaces))
	}

	// ボーンポイントバッファの設定
	md.bonePointIndexes = bonePointIndexes
	md.bonePointCount = len(bonePointIndexes)
	if len(bonePoints) > 0 {
		md.bonePointBufferHandle = factory.CreateBoneBuffer(gl.Ptr(bonePoints), len(bonePoints))
	}
	if len(bonePointFaces) > 0 {
		md.bonePointIbo = factory.CreateElementBuffer(gl.Ptr(bonePointFaces), len(bonePointFaces))
	}

	// 選択頂点バッファの設定
	md.selectedVertexBufferHandle = factory.CreateVertexBuffer(gl.Ptr(selectedVertexVertices), len(selectedVertexVertices))

	// 選択頂点用のインデックスバッファ（すべての頂点のインデックス）
	indexes := createAllVertexIndexesData(model)
	md.selectedVertexIbo = factory.CreateElementBuffer(gl.Ptr(indexes), len(indexes))

	// // カーソル位置表示用バッファの初期化
	// md.cursorPositionBufferHandle = factory.CreateDebugBuffer()

	// SSBOの作成
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, model.Vertices.Length()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)

	mr.ssbo = ssbo
}

// createAllVertexData は頂点データ、法線データ、選択頂点データを一括で生成します
func createAllVertexData(model *pmx.PmxModel) ([]float32, []float32, []float32) {
	vertexCount := model.Vertices.Length()
	vertices := make([]float32, vertexCount*vertexDataSize)
	normalVertices := make([]float32, vertexCount*2*vertexDataSize)
	selectedVertices := make([]float32, vertexCount*vertexDataSize)

	var wg sync.WaitGroup

	// 並列処理のためのバッチサイズ
	batchSize := 1000
	batches := (vertexCount + batchSize - 1) / batchSize

	for b := range batches {
		wg.Add(1)
		go func(batchIndex int) {
			defer wg.Done()
			start := batchIndex * batchSize
			end := min(start+batchSize, vertexCount)

			// 各バッチが担当する位置を計算
			vertexOffset := start * vertexDataSize
			normalOffset := start * 2 * vertexDataSize
			selectedOffset := start * vertexDataSize

			for i := start; i < end; i++ {
				localIdx := i - start
				vertexLocalOffset := localIdx * vertexDataSize
				normalLocalOffset := localIdx * 2 * vertexDataSize
				selectedLocalOffset := localIdx * vertexDataSize

				if vertex, err := model.Vertices.Get(i); err == nil {
					// 通常の頂点データ
					vgl := newVertexGl(vertex)
					copy(vertices[vertexOffset+vertexLocalOffset:], vgl)

					// 法線データ
					normalVgl := newVertexNormalGl(vertex)
					copy(normalVertices[normalOffset+normalLocalOffset:], vgl)                      // 頂点位置
					copy(normalVertices[normalOffset+normalLocalOffset+vertexDataSize:], normalVgl) // 法線方向の終点

					// 選択頂点データ（頂点データと同じ）
					copy(selectedVertices[selectedOffset+selectedLocalOffset:], vgl)
				} else {
					// 空データの場合
					emptyVgl := make([]float32, vertexDataSize)
					copy(vertices[vertexOffset+vertexLocalOffset:], emptyVgl)

					copy(normalVertices[normalOffset+normalLocalOffset:], emptyVgl)
					copy(normalVertices[normalOffset+normalLocalOffset+vertexDataSize:], emptyVgl)

					copy(selectedVertices[selectedOffset+selectedLocalOffset:], emptyVgl)
				}
			}
		}(b)
	}

	wg.Wait()
	return vertices, normalVertices, selectedVertices
}

// createAllBoneDebugData はボーンライン・ポイントデータを一括生成します
func createAllBoneDebugData(model *pmx.PmxModel) ([]float32, []uint32, []int, []float32, []uint32, []int) {
	// ボーン情報の並列処理
	boneLines := make([]float32, 0)
	boneLineFaces := make([]uint32, 0, model.Bones.Length()*2)
	boneLineIndexes := make([]int, model.Bones.Length()*2)

	bonePoints := make([]float32, 0)
	bonePointFaces := make([]uint32, model.Bones.Length())
	bonePointIndexes := make([]int, model.Bones.Length())

	n := 0
	model.Bones.ForEach(func(index int, bone *pmx.Bone) bool {
		boneLines = append(boneLines, newBoneGl(bone)...)
		boneLines = append(boneLines, newTailBoneGl(bone)...)
		boneLineFaces = append(boneLineFaces, uint32(n), uint32(n+1))
		boneLineIndexes[n] = bone.Index()
		boneLineIndexes[n+1] = bone.Index()

		bonePoints = append(bonePoints, newBoneGl(bone)...)
		bonePointFaces[bone.Index()] = uint32(bone.Index())
		bonePointIndexes[bone.Index()] = bone.Index()

		n += 2
		return true
	})

	return boneLines, boneLineFaces, boneLineIndexes, bonePoints, bonePointFaces, bonePointIndexes
}

// createIndexesData はインデックスデータを生成します
func createIndexesData(model *pmx.PmxModel) []uint32 {
	faceCount := model.Faces.Length()
	faces := make([]uint32, faceCount*3)
	var wg sync.WaitGroup

	// 面情報の並列処理
	batchSize := 1000 // 一度に処理する面数
	batches := (faceCount + batchSize - 1) / batchSize

	for b := range batches {
		wg.Add(1)
		go func(batchIndex int) {
			defer wg.Done()
			start := batchIndex * batchSize
			end := min(start+batchSize, faceCount)
			offset := start * 3

			for i := start; i < end; i++ {
				var vertices [3]int
				if face, err := model.Faces.Get(i); err == nil {
					vertices = face.VertexIndexes
				} else {
					vertices = [3]int{0, 0, 0}
				}

				localIdx := (i - start) * 3
				// 頂点の順序を反転（OpenGL用）
				faces[offset+localIdx] = uint32(vertices[2])
				faces[offset+localIdx+1] = uint32(vertices[1])
				faces[offset+localIdx+2] = uint32(vertices[0])
			}
		}(b)
	}

	wg.Wait()
	return faces
}

// createNormalIndexesData は法線インデックスデータを生成します
func createNormalIndexesData(model *pmx.PmxModel) []uint32 {
	normalFaces := make([]uint32, 0, model.Vertices.Length()*2)

	// 各頂点の開始位置と対応する法線位置を結ぶインデックスを生成
	for i := range model.Vertices.Length() {
		n := i * 2
		normalFaces = append(normalFaces, uint32(n), uint32(n+1))
	}

	return normalFaces
}

// createAllVertexIndexesData はすべての頂点インデックスデータを生成します
func createAllVertexIndexesData(model *pmx.PmxModel) []uint32 {
	indexes := make([]uint32, model.Vertices.Length())
	for i := range model.Vertices.Length() {
		indexes[i] = uint32(i)
	}
	return indexes
}
