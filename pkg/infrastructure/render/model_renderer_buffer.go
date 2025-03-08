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

	// ボーン関連バッファを一括で初期化
	var boneLines []float32
	var boneLineFaces []uint32
	var boneLineIndexes []int
	var bonePoints []float32
	var bonePointFaces []uint32
	var bonePointIndexes []int
	wg.Add(1)
	go func() {
		defer wg.Done()
		boneLines, boneLineFaces, boneLineIndexes,
			bonePoints, bonePointFaces, bonePointIndexes = createAllBoneData(model)
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
	md.boneLineBufferHandle = factory.CreateVertexBuffer(gl.Ptr(boneLines), len(boneLines))
	md.boneLineIbo = factory.CreateElementBuffer(gl.Ptr(boneLineFaces), len(boneLineFaces))

	// ボーンポイントバッファの設定
	md.bonePointIndexes = bonePointIndexes
	md.bonePointCount = len(bonePointIndexes)
	md.bonePointBufferHandle = factory.CreateBoneBuffer(gl.Ptr(bonePoints), len(bonePoints))
	md.bonePointIbo = factory.CreateElementBuffer(gl.Ptr(bonePointFaces), len(bonePointFaces))

	// 選択頂点バッファの設定
	md.selectedVertexBufferHandle = factory.CreateVertexBuffer(gl.Ptr(selectedVertexVertices), len(selectedVertexVertices))

	// 選択頂点用のインデックスバッファ（すべての頂点のインデックス）
	indexes := createAllVertexIndexesData(model)
	md.selectedVertexIbo = factory.CreateElementBuffer(gl.Ptr(indexes), len(indexes))

	// カーソル位置表示用バッファの初期化
	md.cursorPositionBufferHandle = factory.CreateDebugBuffer()

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
	vertices := make([]float32, 0, vertexCount*vertexDataSize)
	normalVertices := make([]float32, 0, vertexCount*2*vertexDataSize)
	selectedVertices := make([]float32, 0, vertexCount*vertexDataSize)

	var wg sync.WaitGroup
	var muVert, muNormal, muSelected sync.Mutex

	// 並列処理のためのバッチサイズ
	batchSize := 1000
	batches := (vertexCount + batchSize - 1) / batchSize

	for b := range batches {
		wg.Add(1)
		go func(batchIndex int) {
			defer wg.Done()
			start := batchIndex * batchSize
			end := min(start+batchSize, vertexCount)

			batchVertices := make([]float32, 0, (end-start)*vertexDataSize)
			batchNormalVertices := make([]float32, 0, (end-start)*2*vertexDataSize)
			batchSelectedVertices := make([]float32, 0, (end-start)*vertexDataSize)

			for i := start; i < end; i++ {
				if vertex, err := model.Vertices.Get(i); err == nil {
					// 通常の頂点データ
					vgl := newVertexGl(vertex)
					batchVertices = append(batchVertices, vgl...)

					// 法線データ
					normalVgl := newVertexNormalGl(vertex)
					batchNormalVertices = append(batchNormalVertices, vgl...)       // 頂点位置
					batchNormalVertices = append(batchNormalVertices, normalVgl...) // 法線方向の終点

					// 選択頂点データ（頂点データと同じ）
					batchSelectedVertices = append(batchSelectedVertices, vgl...)
				} else {
					emptyVgl := make([]float32, vertexDataSize)
					batchVertices = append(batchVertices, emptyVgl...)

					batchNormalVertices = append(batchNormalVertices, emptyVgl...)
					batchNormalVertices = append(batchNormalVertices, emptyVgl...)

					batchSelectedVertices = append(batchSelectedVertices, emptyVgl...)
				}
			}

			// スレッドセーフにスライスを更新
			muVert.Lock()
			vertices = append(vertices, batchVertices...)
			muVert.Unlock()

			muNormal.Lock()
			normalVertices = append(normalVertices, batchNormalVertices...)
			muNormal.Unlock()

			muSelected.Lock()
			selectedVertices = append(selectedVertices, batchSelectedVertices...)
			muSelected.Unlock()
		}(b)
	}

	wg.Wait()
	return vertices, normalVertices, selectedVertices
}

// createAllBoneData はボーンライン・ポイントデータを一括生成します
func createAllBoneData(model *pmx.PmxModel) ([]float32, []uint32, []int, []float32, []uint32, []int) {
	boneCount := model.Bones.Length()
	boneLines := make([]float32, 0, boneCount*2*7) // 線の始点と終点
	bonePoints := make([]float32, 0, boneCount*7)  // ボーン位置のみ
	boneLineFaces := make([]uint32, 0, boneCount*2)
	bonePointFaces := make([]uint32, boneCount)
	boneLineIndexes := make([]int, boneCount*2)
	bonePointIndexes := make([]int, boneCount)

	var wg sync.WaitGroup
	var muLine, muPoint sync.Mutex

	// 並列処理のためのバッチサイズ
	batchSize := 500
	batches := (boneCount + batchSize - 1) / batchSize

	for b := range batches {
		wg.Add(1)
		go func(batchIndex int) {
			defer wg.Done()
			start := batchIndex * batchSize
			end := min(start+batchSize, boneCount)

			batchLines := make([]float32, 0, (end-start)*2*7)
			batchPoints := make([]float32, 0, (end-start)*7)
			batchLineFaces := make([]uint32, 0, (end-start)*2)
			batchPointFaces := make([]uint32, end-start)
			batchLineIndexes := make(map[int]int, (end-start)*2)
			batchPointIndexes := make([]int, end-start)

			for i := start; i < end; i++ {
				var bone *pmx.Bone
				if b, err := model.Bones.Get(i); err == nil {
					bone = b
				} else {
					bone = pmx.NewBone()
				}
				n := i * 2

				// ボーンラインデータ
				boneStartGL := newBoneGl(bone)
				batchLines = append(batchLines, boneStartGL...)

				boneEndGL := newTailBoneGl(bone)
				batchLines = append(batchLines, boneEndGL...)

				batchLineFaces = append(batchLineFaces, uint32(n), uint32(n+1))

				batchLineIndexes[n] = bone.Index()
				batchLineIndexes[n+1] = bone.Index()

				// ボーンポイントデータ (始点のみ)
				batchPoints = append(batchPoints, boneStartGL...)

				localIdx := i - start
				batchPointFaces[localIdx] = uint32(bone.Index())
				batchPointIndexes[localIdx] = bone.Index()
			}

			// スレッドセーフにデータを統合
			muLine.Lock()
			boneLines = append(boneLines, batchLines...)
			boneLineFaces = append(boneLineFaces, batchLineFaces...)

			// ラインインデックスの更新
			for idx, boneIdx := range batchLineIndexes {
				boneLineIndexes[idx] = boneIdx
			}
			muLine.Unlock()

			muPoint.Lock()
			bonePoints = append(bonePoints, batchPoints...)

			// ポイントインデックスの更新
			for i := start; i < end; i++ {
				localIdx := i - start
				bonePointFaces[i] = batchPointFaces[localIdx]
				bonePointIndexes[i] = batchPointIndexes[localIdx]
			}
			muPoint.Unlock()
		}(b)
	}

	wg.Wait()
	return boneLines, boneLineFaces, boneLineIndexes, bonePoints, bonePointFaces, bonePointIndexes
}

// createIndexesData はインデックスデータを生成します
func createIndexesData(model *pmx.PmxModel) []uint32 {
	faces := make([]uint32, 0, model.Faces.Length()*3)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 面情報の並列処理
	batchSize := 1000 // 一度に処理する面数
	faceCount := model.Faces.Length()
	batches := (faceCount + batchSize - 1) / batchSize

	for b := range batches {
		wg.Add(1)
		go func(batchIndex int) {
			defer wg.Done()
			start := batchIndex * batchSize
			end := min(start+batchSize, faceCount)

			batchFaces := make([]uint32, 0, batchSize*3)

			for i := start; i < end; i++ {
				var vertices [3]int
				if face, err := model.Faces.Get(i); err == nil {
					vertices = face.VertexIndexes
				} else {
					vertices = [3]int{0, 0, 0}
				}
				// 頂点の順序を反転（OpenGL用）
				batchFaces = append(batchFaces, uint32(vertices[2]), uint32(vertices[1]), uint32(vertices[0]))
			}

			mu.Lock()
			faces = append(faces, batchFaces...)
			mu.Unlock()
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
