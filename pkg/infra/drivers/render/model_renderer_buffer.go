//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"sync"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
)

// ModelRenderer (旧 RenderModel) のバッファ初期化処理をまとめたファイルです。
// このファイルでは、頂点情報、法線、ボーンライン・ボーンポイント、選択頂点、カーソル位置、SSBO など
// モデル描画に必要な OpenGL バッファの生成・データ転送処理を行います。

// initializeBuffers は、モデルの頂点バッファおよび関連バッファを初期化します。
func (mr *ModelRenderer) initializeBuffers(factory *mgl.BufferFactory, modelData *model.PmxModel) {
	// ModelDrawer を初期化
	md := mr.ModelDrawer

	// WaitGroupを用いて並列処理を管理
	var wg sync.WaitGroup

	// 頂点データ、法線データ、選択頂点データを一括生成
	var selectedVertexVertices []float32
	wg.Add(1)
	go func() {
		defer wg.Done()
		md.vertices, md.normalVertices, selectedVertexVertices = createAllVertexData(modelData)
	}()

	// 面データ一括生成
	wg.Add(1)
	go func() {
		defer wg.Done()
		md.faces = createIndexesData(modelData)
	}()

	// ボーンデバッグ関連バッファを一括で初期化
	var boneLines, bonePoints []float32
	var boneLineFaces, bonePointFaces []uint32
	var boneLineIndexes, bonePointIndexes []int
	wg.Add(1)
	go func() {
		defer wg.Done()
		boneLines, boneLineFaces, boneLineIndexes,
			bonePoints, bonePointFaces, bonePointIndexes = createAllBoneDebugData(modelData)
	}()

	// WaitGroupの完了を待つ
	wg.Wait()

	// GLオブジェクトの生成は並列化しない ------------

	// メインのモデル頂点バッファ
	mr.bufferHandle = factory.NewVertexBuffer(gl.Ptr(md.vertices), len(md.vertices))

	// 法線バッファの初期化
	md.normalBufferHandle = factory.NewVertexBuffer(gl.Ptr(md.normalVertices), len(md.normalVertices))

	// 法線表示用のインデックスバッファ
	normalIndexes := createNormalIndexesData(modelData)
	md.normalIndexCount = len(normalIndexes)
	md.normalIbo = factory.NewIndexBuffer(gl.Ptr(normalIndexes), len(normalIndexes))

	// ボーンラインバッファの設定
	md.boneLineIndexes = boneLineIndexes
	md.boneLineCount = len(boneLineIndexes)
	if len(boneLines) > 0 {
		md.boneLineBufferHandle = factory.NewBoneBuffer(gl.Ptr(boneLines), len(boneLines))
	}
	if len(boneLineFaces) > 0 {
		md.boneLineIbo = factory.NewIndexBuffer(gl.Ptr(boneLineFaces), len(boneLineFaces))
	}

	// ボーンポイントバッファの設定
	md.bonePointIndexes = bonePointIndexes
	md.bonePointCount = len(bonePointIndexes)
	if len(bonePoints) > 0 {
		md.bonePointBufferHandle = factory.NewBoneBuffer(gl.Ptr(bonePoints), len(bonePoints))
	}
	if len(bonePointFaces) > 0 {
		md.bonePointIbo = factory.NewIndexBuffer(gl.Ptr(bonePointFaces), len(bonePointFaces))
	}

	// 選択頂点バッファの設定
	md.selectedVertexBufferHandle = factory.NewVertexBuffer(gl.Ptr(selectedVertexVertices), len(selectedVertexVertices))
	// 選択頂点用の元データ参照を保持して、VBOのベースコピーが破棄されないようにする。
	md.selectedVertexVertices = selectedVertexVertices

	// 選択頂点用のインデックスバッファ（すべての頂点のインデックス）
	indexes := createAllVertexIndexesData(modelData)
	md.selectedVertexCount = len(indexes)
	md.selectedVertexIbo = factory.NewIndexBuffer(gl.Ptr(indexes), len(indexes))

	// カーソル位置表示用バッファの初期化
	md.cursorPositionBufferHandle = factory.NewDebugBuffer(nil, 0)

	// SSBOの作成
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, modelData.Vertices.Len()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)

	mr.ssbo = ssbo
}

// createAllVertexData は頂点データ、法線データ、選択頂点データを一括で生成します
func createAllVertexData(modelData *model.PmxModel) ([]float32, []float32, []float32) {
	vertexCount := modelData.Vertices.Len()
	vertices := make([]float32, vertexCount*vertexDataSize)
	normalVertices := make([]float32, vertexCount*2*vertexDataSize)
	selectedVertices := make([]float32, vertexCount*vertexDataSize)

	var wg sync.WaitGroup

	// 並列処理のためのバッチサイズ
	batchSize := 1000
	batches := (vertexCount + batchSize - 1) / batchSize

	for b := 0; b < batches; b++ {
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

				if vertex, err := modelData.Vertices.Get(i); err == nil {
					// 通常の頂点データ
					vgl := newVertexGl(vertex)
					copy(vertices[vertexOffset+vertexLocalOffset:], vgl)

					// 法線データ
					normalVgl := newVertexNormalGl(vertex)
					copy(normalVertices[normalOffset+normalLocalOffset:], vgl)                      // 頂点位置
					copy(normalVertices[normalOffset+normalLocalOffset+vertexDataSize:], normalVgl) // 法線方向の終点

					// 選択頂点データ（非選択時は描画されないUVにしておく）
					svgl := newSelectedVertexGl(vertex)
					copy(selectedVertices[selectedOffset+selectedLocalOffset:], svgl)
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
func createAllBoneDebugData(modelData *model.PmxModel) ([]float32, []uint32, []int, []float32, []uint32, []int) {
	// ボーン情報の並列処理
	boneLines := make([]float32, 0)
	boneLineFaces := make([]uint32, 0, modelData.Bones.Len()*2)
	boneLineIndexes := make([]int, modelData.Bones.Len()*2)

	bonePoints := make([]float32, 0)
	bonePointFaces := make([]uint32, modelData.Bones.Len())
	bonePointIndexes := make([]int, modelData.Bones.Len())

	n := 0
	for _, bone := range modelData.Bones.Values() {
		if bone == nil {
			continue
		}
		tailPos := calcTailPosition(bone, modelData.Bones)
		boneLines = append(boneLines, newBoneGl(bone)...)
		boneLines = append(boneLines, newTailBoneGl(bone, tailPos)...)
		boneLineFaces = append(boneLineFaces, uint32(n), uint32(n+1))
		boneLineIndexes[n] = bone.Index()
		boneLineIndexes[n+1] = bone.Index()

		bonePoints = append(bonePoints, newBoneGl(bone)...)
		bonePointFaces[bone.Index()] = uint32(bone.Index())
		bonePointIndexes[bone.Index()] = bone.Index()

		n += 2
	}

	return boneLines, boneLineFaces, boneLineIndexes, bonePoints, bonePointFaces, bonePointIndexes
}

// createIndexesData はインデックスデータを生成します
func createIndexesData(modelData *model.PmxModel) []uint32 {
	faceCount := modelData.Faces.Len()
	faces := make([]uint32, faceCount*3)
	var wg sync.WaitGroup

	// 面情報の並列処理
	batchSize := 1000 // 一度に処理する面数
	batches := (faceCount + batchSize - 1) / batchSize

	for b := 0; b < batches; b++ {
		wg.Add(1)
		go func(batchIndex int) {
			defer wg.Done()
			start := batchIndex * batchSize
			end := min(start+batchSize, faceCount)
			offset := start * 3

			for i := start; i < end; i++ {
				var vertices [3]int
				if face, err := modelData.Faces.Get(i); err == nil {
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
func createNormalIndexesData(modelData *model.PmxModel) []uint32 {
	normalFaces := make([]uint32, 0, modelData.Vertices.Len()*2)

	// 各頂点の開始位置と対応する法線位置を結ぶインデックスを生成
	for i := 0; i < modelData.Vertices.Len(); i++ {
		n := i * 2
		normalFaces = append(normalFaces, uint32(n), uint32(n+1))
	}

	return normalFaces
}

// createAllVertexIndexesData はすべての頂点インデックスデータを生成します
func createAllVertexIndexesData(modelData *model.PmxModel) []uint32 {
	indexes := make([]uint32, modelData.Vertices.Len())
	for i := 0; i < modelData.Vertices.Len(); i++ {
		indexes[i] = uint32(i)
	}
	return indexes
}
