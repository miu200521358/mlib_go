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

// initializeBuffers は、与えられた pmx.PmxModel の情報から各種バッファを初期化します。
func (mr *ModelRenderer) initializeBuffers(model *pmx.PmxModel) {
	// --- 頂点・法線・選択頂点用データの作成 ---
	mr.vertices = make([]float32, 0, model.Vertices.Length())
	mr.normalVertices = make([]float32, 0, model.Vertices.Length()*2)
	normalFaces := make([]uint32, 0, model.Vertices.Length()*2)
	selectedVertices := make([]float32, 0, model.Vertices.Length())
	selectedVertexFaces := make([]uint32, 0, model.Vertices.Length())
	mr.faces = make([]uint32, 0, model.Faces.Length()*3)

	// --- ボーン情報用データの作成 ---
	var boneLineFaces []uint32
	boneLineFaces = make([]uint32, 0)
	mr.boneLineIndexes = make([]int, 0)
	var bonePointFaces []uint32
	bonePointFaces = make([]uint32, 0)
	mr.bonePointIndexes = make([]int, 0)

	// 並列処理で頂点情報を生成
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for v := range model.Vertices.Iterator() {
			i := v.Index
			vertex := v.Value
			vgl := newVertexGl(vertex)
			// 選択頂点用データ (新SelectedVertexGl)
			selectedV := newSelectedVertexGl(vertex)
			// 法線 (newVertexNormalGl)
			normalV := newVertexNormalGl(vertex)
			mu.Lock()
			mr.vertices = append(mr.vertices, vgl...)
			// 法線は、元の頂点とその法線の2頂点を並べる想定
			mr.normalVertices = append(mr.normalVertices, vgl...)
			mr.normalVertices = append(mr.normalVertices, normalV...)
			normalFaces = append(normalFaces, uint32(n), uint32(n+1))
			// 選択頂点情報
			selectedVertices = append(selectedVertices, selectedV...)
			selectedVertexFaces = append(selectedVertexFaces, uint32(i))
			mu.Unlock()
			n += 2
		}
	}()

	// 面情報 (通常のメッシュ面) を生成
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range model.Faces.Iterator() {
			face := v.Value
			vertices := face.VertexIndexes
			mu.Lock()
			// 頂点の順序を反転して (頂点順が逆なら) 格納
			mr.faces = append(mr.faces, uint32(vertices[0]), uint32(vertices[1]), uint32(vertices[2]))
			mu.Unlock()
		}
	}()

	// ボーン情報を生成
	var boneVertices []float32
	wg.Add(1)
	go func() {
		defer wg.Done()
		nBone := 0
		for v := range model.Bones.Iterator() {
			bone := v.Value
			// 各ボーンに対して、始点と末端の頂点データを生成
			boneStart := newBoneGl(bone)
			boneEnd := newTailBoneGl(bone)
			mu.Lock()
			boneVertices = append(boneVertices, boneStart...)
			boneVertices = append(boneVertices, boneEnd...)
			// インデックスとして、始点と末端を登録
			boneLineFaces = append(boneLineFaces, uint32(nBone), uint32(nBone+1))
			// ボーンのインデックスを記録（重複しても可）
			mr.boneLineIndexes = append(mr.boneLineIndexes, bone.Index(), bone.Index())
			// ボーンポイントは、先頭の位置を使う (簡易的な実装)
			mr.bonePointIndexes = append(mr.bonePointIndexes, bone.Index())
			mu.Unlock()
			nBone += 2
		}
	}()

	wg.Wait()

	// --- VAO/VBO/ElementBuffer の生成 ---
	// 頂点情報の VAO/VBO
	mr.vao = mgl.NewVertexArray()
	mr.vao.Bind()
	mr.vbo = mgl.NewVertexBuffer()
	if len(mr.vertices) > 0 {
		mr.vbo.BufferData(len(mr.vertices)*4, gl.Ptr(mr.vertices), gl.STATIC_DRAW)
	}
	// BindVertex の代替として、単に VBO の BufferData を呼ぶ（初期化時）
	mr.vbo.Bind()
	mr.vbo.Unbind()
	mr.vao.Unbind()

	// 法線情報の VAO/VBO/ElementBuffer
	mr.normalVao = mgl.NewVertexArray()
	mr.normalVao.Bind()
	mr.vbo = mgl.NewVertexBuffer()
	if len(mr.normalVertices) > 0 {
		mr.normalVbo.BufferData(len(mr.normalVertices)*4, gl.Ptr(mr.normalVertices), gl.STATIC_DRAW)
	}
	mr.normalVbo.Bind()
	if len(normalFaces) == 0 {
		mr.normalIbo = mgl.NewElementBuffer()
		mr.normalIbo.BufferData(0, nil, gl.STATIC_DRAW)
	} else {
		mr.normalIbo = mgl.NewElementBuffer()
		mr.normalIbo.Bind()
		mr.normalIbo.BufferData(len(normalFaces)*4, gl.Ptr(normalFaces), gl.STATIC_DRAW)
		mr.normalIbo.Unbind()
	}
	mr.normalVbo.Unbind()
	mr.normalVao.Unbind()

	// ボーンライン用 VAO/VBO/ElementBuffer
	mr.boneLineVao = mgl.NewVertexArray()
	mr.boneLineVao.Bind()
	mr.vbo = mgl.NewVertexBuffer()
	if len(boneVertices) > 0 {
		mr.boneLineVbo.BufferData(len(boneVertices)*4, gl.Ptr(boneVertices), gl.STATIC_DRAW)
	}
	// BindBone の代替として、単に Bind を行う
	mr.boneLineVbo.Bind()
	if len(boneLineFaces) == 0 {
		mr.boneLineIbo = mgl.NewElementBuffer()
		mr.boneLineIbo.BufferData(0, nil, gl.STATIC_DRAW)
	} else {
		mr.boneLineIbo = mgl.NewElementBuffer()
		mr.boneLineIbo.Bind()
		mr.boneLineIbo.BufferData(len(boneLineFaces)*4, gl.Ptr(boneLineFaces), gl.STATIC_DRAW)
		mr.boneLineIbo.Unbind()
	}
	mr.boneLineVbo.Unbind()
	mr.boneLineVao.Unbind()
	// 保存：ボーンラインの描画点数（インデックス数）
	mr.boneLineCount = len(boneLineFaces)

	// ボーンポイント用 VAO/VBO/ElementBuffer
	mr.bonePointVao = mgl.NewVertexArray()
	mr.bonePointVao.Bind()
	mr.bonePointVbo = mgl.NewVertexBuffer()
	mr.bonePointVbo.BufferData(len(boneVertices)*4, gl.Ptr(boneVertices), gl.STATIC_DRAW)
	mr.bonePointVbo.Bind()
	if len(bonePointFaces) == 0 {
		mr.bonePointIbo = mgl.NewElementBuffer()
		mr.bonePointIbo.BufferData(0, nil, gl.STATIC_DRAW)
	} else {
		mr.bonePointIbo = mgl.NewElementBuffer()
		mr.bonePointIbo.Bind()
		mr.bonePointIbo.BufferData(len(bonePointFaces)*4, gl.Ptr(bonePointFaces), gl.STATIC_DRAW)
		mr.bonePointIbo.Unbind()
	}
	mr.bonePointVbo.Unbind()
	mr.bonePointVao.Unbind()
	// 保存：ボーンポイントの数（シンプルに bonePointFaces の長さ）
	mr.bonePointCount = len(bonePointFaces)

	// 選択頂点用 VAO/VBO/ElementBuffer
	mr.selectedVertexVao = mgl.NewVertexArray()
	mr.selectedVertexVao.Bind()
	mr.selectedVertexVbo = mgl.NewVertexBuffer()
	if len(selectedVertices) > 0 {
		mr.selectedVertexVbo.BufferData(len(selectedVertices)*4, gl.Ptr(selectedVertices), gl.STATIC_DRAW)
	}
	mr.selectedVertexVbo.Bind()
	if len(selectedVertexFaces) == 0 {
		mr.selectedVertexIbo = mgl.NewElementBuffer()
		mr.selectedVertexIbo.BufferData(0, nil, gl.STATIC_DRAW)
	} else {
		mr.selectedVertexIbo = mgl.NewElementBuffer()
		mr.selectedVertexIbo.Bind()
		mr.selectedVertexIbo.BufferData(len(selectedVertexFaces)*4, gl.Ptr(selectedVertexFaces), gl.STATIC_DRAW)
		mr.selectedVertexIbo.Unbind()
	}
	mr.selectedVertexVbo.Unbind()
	mr.selectedVertexVao.Unbind()

	// カーソル位置用 VAO/VBO
	mr.cursorPositionVao = mgl.NewVertexArray()
	mr.cursorPositionVao.Bind()
	mr.cursorPositionVbo = mgl.NewVertexBuffer()
	mr.cursorPositionVbo.BufferData(3*4, gl.Ptr([]float32{0, 0, 0}), gl.STATIC_DRAW)
	mr.cursorPositionVbo.Unbind()
	mr.cursorPositionVao.Unbind()

	// SSBO の作成
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, model.Vertices.Length()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}
