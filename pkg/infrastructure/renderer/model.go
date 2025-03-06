//go:build windows
// +build windows

package render

import (
	"sync"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type RenderModel struct {
	windowIndex              int                      // ウィンドウインデックス
	invisibleMaterialIndexes map[int]struct{}         // 非表示材質インデックス
	selectedVertexes         map[int]struct{}         // 選択頂点インデックス
	noSelectedVertexes       map[int]struct{}         // 非選択頂点インデックス
	model                    *pmx.PmxModel            // 元モデル
	hash                     string                   // ハッシュ
	meshes                   []*Mesh                  // メッシュ
	textures                 []*textureGl             // テクスチャ
	toonTextures             []*textureGl             // トゥーンテクスチャ
	vertices                 []float32                // 頂点情報
	vao                      rendering.IVertexArray   // 頂点VAO
	vbo                      rendering.IVertexBuffer  // 頂点VBO
	normalVertices           []float32                // 法線情報
	normalVao                rendering.IVertexArray   // 法線VAO
	normalVbo                rendering.IVertexBuffer  // 法線VBO
	normalIbo                rendering.IElementBuffer // 法線IBO
	selectedVertexVao        rendering.IVertexArray   // 選択頂点VAO
	selectedVertexVbo        rendering.IVertexBuffer  // 選択頂点VBO
	selectedVertexIbo        rendering.IElementBuffer // 選択頂点IBO
	boneLineCount            int                      // ボーン情報
	boneLineVao              rendering.IVertexArray   // ボーンVAO
	boneLineVbo              rendering.IVertexBuffer  // ボーンVBO
	boneLineIbo              rendering.IElementBuffer // ボーンIBO
	boneLineIndexes          []int                    // ボーンインデックス
	bonePointCount           int                      // ボーン情報
	bonePointVao             rendering.IVertexArray   // ボーンVAO
	bonePointVbo             rendering.IVertexBuffer  // ボーンVBO
	bonePointIbo             rendering.IElementBuffer // ボーンIBO
	bonePointIndexes         []int                    // ボーンインデックス
	cursorPositionVao        rendering.IVertexArray   // カーソルラインVAO
	cursorPositionVbo        rendering.IVertexBuffer  // カーソルラインVBO
	ssbo                     uint32                   // SSBO
	vertexCount              int                      // 頂点数
}

func NewRenderModel(windowIndex int, model *pmx.PmxModel) *RenderModel {
	renderModel := &RenderModel{
		windowIndex: windowIndex,
		model:       model,
		vertexCount: model.Vertices.Length(),
	}
	renderModel.initToonTexturesGl(windowIndex)
	renderModel.initTexturesGl(windowIndex, model.Textures, model.Path())

	renderModel.initializeBuffer(model)
	renderModel.hash = model.Hash()
	renderModel.selectedVertexes = make(map[int]struct{})

	return renderModel
}

func (renderModel *RenderModel) Model() *pmx.PmxModel {
	return renderModel.model
}

func (renderModel *RenderModel) Hash() string {
	return renderModel.hash
}

func (renderModel *RenderModel) Delete() {
	if renderModel.vao != nil {
		renderModel.vao.Delete()
	}
	if renderModel.vbo != nil {
		renderModel.vbo.Delete()
	}
	if renderModel.normalVao != nil {
		renderModel.normalVao.Delete()
	}
	if renderModel.normalVbo != nil {
		renderModel.normalVbo.Delete()
	}
	if renderModel.normalIbo != nil {
		renderModel.normalIbo.Delete()
	}
	if renderModel.selectedVertexVao != nil {
		renderModel.selectedVertexVao.Delete()
	}
	if renderModel.selectedVertexVbo != nil {
		renderModel.selectedVertexVbo.Delete()
	}
	if renderModel.selectedVertexIbo != nil {
		renderModel.selectedVertexIbo.Delete()
	}
	if renderModel.boneLineVao != nil {
		renderModel.boneLineVao.Delete()
	}
	if renderModel.boneLineVbo != nil {
		renderModel.boneLineVbo.Delete()
	}
	if renderModel.boneLineIbo != nil {
		renderModel.boneLineIbo.Delete()
	}
	if renderModel.ssbo != 0 {
		gl.DeleteBuffers(1, &renderModel.ssbo)
	}
	for _, mesh := range renderModel.meshes {
		if mesh != nil {
			mesh.delete()
		}
	}
	for _, texture := range renderModel.textures {
		if texture != nil {
			texture.delete()
		}
	}
	for _, texture := range renderModel.toonTextures {
		if texture != nil {
			texture.delete()
		}
	}
}

func (renderModel *RenderModel) InvisibleMaterials() []int {
	if renderModel.invisibleMaterialIndexes == nil {
		return nil
	}

	indexes := make([]int, 0, len(renderModel.invisibleMaterialIndexes))
	for i := range renderModel.invisibleMaterialIndexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (renderModel *RenderModel) SetInvisibleMaterials(indexes []int) {
	if len(indexes) == 0 {
		renderModel.invisibleMaterialIndexes = nil
		return
	}
	renderModel.invisibleMaterialIndexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		renderModel.invisibleMaterialIndexes[i] = struct{}{}
	}
}

func (renderModel *RenderModel) ExistInvisibleMaterial(index int) bool {
	_, ok := renderModel.invisibleMaterialIndexes[index]
	return ok
}

func (renderModel *RenderModel) SelectedVertexes() []int {
	indexes := make([]int, 0, len(renderModel.selectedVertexes))
	for i := range renderModel.selectedVertexes {
		if _, ok := renderModel.noSelectedVertexes[i]; !ok {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (renderModel *RenderModel) NoSelectedVertexes() []int {
	indexes := make([]int, 0, len(renderModel.noSelectedVertexes))
	for i := range renderModel.noSelectedVertexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (renderModel *RenderModel) SetSelectedVertexes(indexes []int) {
	renderModel.selectedVertexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		renderModel.selectedVertexes[i] = struct{}{}
	}
}

func (renderModel *RenderModel) SetNoSelectedVertexes(indexes []int) {
	renderModel.noSelectedVertexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		renderModel.noSelectedVertexes[i] = struct{}{}
	}
}

func (renderModel *RenderModel) ClearSelectedVertexes() {
	renderModel.UpdateNoSelectedVertexes(renderModel.SelectedVertexes())
	renderModel.selectedVertexes = make(map[int]struct{})
}

func (renderModel *RenderModel) UpdateSelectedVertexes(indexes []int) {
	for _, index := range indexes {
		renderModel.selectedVertexes[index] = struct{}{}
		delete(renderModel.noSelectedVertexes, index)
	}
}

func (renderModel *RenderModel) UpdateNoSelectedVertexes(indexes []int) {
	if renderModel.noSelectedVertexes == nil {
		renderModel.noSelectedVertexes = make(map[int]struct{}, len(indexes))
	}
	for _, index := range indexes {
		renderModel.noSelectedVertexes[index] = struct{}{}
		delete(renderModel.selectedVertexes, index)
	}
}

// bufferData は描画に必要なバッファデータを保持する構造体
type bufferData struct {
	Vertices            []float32
	NormalVertices      []float32
	NormalFaces         []uint32
	SelectedVertices    []float32
	SelectedVertexFaces []uint32
	Faces               []uint32
	BoneLines           []float32
	BoneLineFaces       []uint32
	BoneLineIndexes     []int
	BonePoints          []float32
	BonePointFaces      []uint32
	BonePointIndexes    []int
}

func (renderModel *RenderModel) initializeBuffer(model *pmx.PmxModel) {
	// データの準備
	bufferData := renderModel.prepareBufferData(model)

	// OpenGL バッファの生成
	renderModel.createVertexBuffers(bufferData)
	renderModel.createNormalBuffers(bufferData)
	renderModel.createBoneBuffers(bufferData)
	renderModel.createSelectedVertexBuffers(bufferData)
	renderModel.createCursorPositionBuffers()

	// メッシュの初期化
	renderModel.initializeMeshes(model, bufferData.Faces)

	// SSBOの作成
	renderModel.createSSBO(model)
}

// prepareBufferData は描画に必要なすべてのバッファデータを準備します
func (renderModel *RenderModel) prepareBufferData(model *pmx.PmxModel) *bufferData {
	data := &bufferData{
		Vertices:            make([]float32, 0, model.Vertices.Length()*3),
		NormalVertices:      make([]float32, 0, model.Vertices.Length()*6),
		NormalFaces:         make([]uint32, 0, model.Vertices.Length()*2),
		SelectedVertices:    make([]float32, 0, model.Vertices.Length()*3),
		SelectedVertexFaces: make([]uint32, 0, model.Vertices.Length()),
		Faces:               make([]uint32, 0, model.Faces.Length()*3),
		BoneLines:           make([]float32, 0),
		BoneLineFaces:       make([]uint32, 0, model.Bones.Length()*2),
		BoneLineIndexes:     make([]int, model.Bones.Length()*2),
		BonePoints:          make([]float32, 0),
		BonePointFaces:      make([]uint32, model.Bones.Length()),
		BonePointIndexes:    make([]int, model.Bones.Length()),
	}

	// 並列処理でデータを準備
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 頂点と法線情報の準備
	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for v := range model.Vertices.Iterator() {
			i := v.Index
			vertex := v.Value
			vgl := newVertexGl(vertex)

			mu.Lock()
			// 頂点情報
			data.Vertices = append(data.Vertices, vgl...)

			// 法線情報
			data.NormalVertices = append(data.NormalVertices, vgl...)
			data.NormalVertices = append(data.NormalVertices, newVertexNormalGl(vertex)...)
			data.NormalFaces = append(data.NormalFaces, uint32(n), uint32(n+1))

			// 選択頂点情報
			data.SelectedVertices = append(data.SelectedVertices, newSelectedVertexGl(vertex)...)
			data.SelectedVertexFaces = append(data.SelectedVertexFaces, uint32(i))
			mu.Unlock()

			n += 2
		}
	}()

	// 面情報の準備
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range model.Faces.Iterator() {
			face := v.Value
			vertices := face.VertexIndexes
			mu.Lock()
			data.Faces = append(data.Faces, uint32(vertices[2]), uint32(vertices[1]), uint32(vertices[0]))
			mu.Unlock()
		}
	}()

	// ボーン情報の準備
	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for v := range model.Bones.Iterator() {
			bone := v.Value
			boneStart := newBoneGl(bone)
			boneEnd := newTailBoneGl(bone)

			mu.Lock()
			// ボーンライン情報
			data.BoneLines = append(data.BoneLines, boneStart...)
			data.BoneLines = append(data.BoneLines, boneEnd...)
			data.BoneLineFaces = append(data.BoneLineFaces, uint32(n), uint32(n+1))
			data.BoneLineIndexes[n] = bone.Index()
			data.BoneLineIndexes[n+1] = bone.Index()

			// ボーンポイント情報
			data.BonePoints = append(data.BonePoints, boneStart...)
			data.BonePointFaces[bone.Index()] = uint32(bone.Index())
			data.BonePointIndexes[bone.Index()] = bone.Index()
			mu.Unlock()

			n += 2
		}
	}()

	wg.Wait() // すべての処理が完了するまで待機
	return data
}

// createVertexBuffers は頂点バッファを生成します
func (renderModel *RenderModel) createVertexBuffers(data *bufferData) {
	renderModel.vao = mgl.NewVertexArray()
	renderModel.vao.Bind()

	renderModel.vbo = renderModel.createVertexBuffer(data.Vertices, rendering.BufferUsageStatic)
	renderModel.vbo.Bind()
	renderModel.vbo.SetAttribute(rendering.VertexAttribute{
		Index:     0,
		Size:      3,
		Type:      gl.FLOAT,
		Normalize: false,
		Stride:    0,
		Offset:    0,
	})
	renderModel.vbo.Unbind()

	renderModel.vao.Unbind()
}

// createNormalBuffers は法線表示用のバッファを生成します
func (renderModel *RenderModel) createNormalBuffers(data *bufferData) {
	renderModel.normalVao = mgl.NewVertexArray()
	renderModel.normalVao.Bind()

	renderModel.normalVbo = renderModel.createVertexBuffer(data.NormalVertices, rendering.BufferUsageStatic)
	renderModel.normalVbo.Bind()
	renderModel.normalVbo.SetAttribute(rendering.VertexAttribute{
		Index:     0,
		Size:      3,
		Type:      gl.FLOAT,
		Normalize: false,
		Stride:    0,
		Offset:    0,
	})

	renderModel.normalIbo = renderModel.createElementBuffer(data.NormalFaces, rendering.BufferUsageStatic)
	renderModel.normalIbo.Bind()
	renderModel.normalIbo.Unbind()
	renderModel.normalVbo.Unbind()

	renderModel.normalVao.Unbind()
}

// createBoneBuffers はボーン表示用のバッファを生成します
func (renderModel *RenderModel) createBoneBuffers(data *bufferData) {
	// ボーンライン用バッファ
	renderModel.boneLineCount = len(data.BoneLines) / 3
	renderModel.boneLineIndexes = data.BoneLineIndexes
	renderModel.boneLineVao = mgl.NewVertexArray()
	renderModel.boneLineVao.Bind()

	renderModel.boneLineVbo = renderModel.createVertexBuffer(data.BoneLines, rendering.BufferUsageStatic)
	renderModel.boneLineVbo.Bind()
	renderModel.boneLineVbo.SetAttribute(rendering.VertexAttribute{
		Index:     0,
		Size:      3,
		Type:      gl.FLOAT,
		Normalize: false,
		Stride:    0,
		Offset:    0,
	})

	renderModel.boneLineIbo = renderModel.createElementBuffer(data.BoneLineFaces, rendering.BufferUsageStatic)
	renderModel.boneLineIbo.Bind()
	renderModel.boneLineIbo.Unbind()
	renderModel.boneLineVbo.Unbind()

	renderModel.boneLineVao.Unbind()

	// ボーンポイント用バッファ
	renderModel.bonePointCount = len(data.BonePoints) / 3
	renderModel.bonePointIndexes = data.BonePointIndexes
	renderModel.bonePointVao = mgl.NewVertexArray()
	renderModel.bonePointVao.Bind()

	renderModel.bonePointVbo = renderModel.createVertexBuffer(data.BonePoints, rendering.BufferUsageStatic)
	renderModel.bonePointVbo.Bind()
	renderModel.bonePointVbo.SetAttribute(rendering.VertexAttribute{
		Index:     0,
		Size:      3,
		Type:      gl.FLOAT,
		Normalize: false,
		Stride:    0,
		Offset:    0,
	})

	renderModel.bonePointIbo = renderModel.createElementBuffer(data.BonePointFaces, rendering.BufferUsageStatic)
	renderModel.bonePointIbo.Bind()
	renderModel.bonePointIbo.Unbind()
	renderModel.bonePointVbo.Unbind()

	renderModel.bonePointVao.Unbind()
}

// createSelectedVertexBuffers は選択頂点表示用のバッファを生成します
func (renderModel *RenderModel) createSelectedVertexBuffers(data *bufferData) {
	renderModel.selectedVertexVao = mgl.NewVertexArray()
	renderModel.selectedVertexVao.Bind()

	renderModel.selectedVertexVbo = renderModel.createVertexBuffer(data.SelectedVertices, rendering.BufferUsageStatic)
	renderModel.selectedVertexVbo.Bind()
	renderModel.selectedVertexVbo.SetAttribute(rendering.VertexAttribute{
		Index:     0,
		Size:      3,
		Type:      gl.FLOAT,
		Normalize: false,
		Stride:    0,
		Offset:    0,
	})

	renderModel.selectedVertexIbo = renderModel.createElementBuffer(data.SelectedVertexFaces, rendering.BufferUsageStatic)
	renderModel.selectedVertexIbo.Bind()
	renderModel.selectedVertexIbo.Unbind()
	renderModel.selectedVertexVbo.Unbind()

	renderModel.selectedVertexVao.Unbind()
}

// createCursorPositionBuffers はカーソル位置表示用のバッファを生成します
func (renderModel *RenderModel) createCursorPositionBuffers() {
	cursorPositions := []float32{0, 0, 0}
	renderModel.cursorPositionVao = mgl.NewVertexArray()
	renderModel.cursorPositionVao.Bind()

	renderModel.cursorPositionVbo = rendering.IVertexBuffer(mgl.NewVertexBuffer())
	renderModel.cursorPositionVbo.Bind()
	renderModel.cursorPositionVbo.BufferData(len(cursorPositions)*4, gl.Ptr(cursorPositions), rendering.BufferUsageDynamic)
	renderModel.cursorPositionVbo.SetAttribute(rendering.VertexAttribute{
		Index:     0,
		Size:      3,
		Type:      gl.FLOAT,
		Normalize: false,
		Stride:    0,
		Offset:    0,
	})
	renderModel.cursorPositionVbo.Unbind()

	renderModel.cursorPositionVao.Unbind()
}

// createSSBO はシェーダストレージバッファオブジェクトを作成します
func (renderModel *RenderModel) createSSBO(model *pmx.PmxModel) {
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, model.Vertices.Length()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	renderModel.ssbo = ssbo
}

// initializeMeshes はメッシュ情報を初期化します
func (renderModel *RenderModel) initializeMeshes(model *pmx.PmxModel, faces []uint32) {
	renderModel.meshes = make([]*Mesh, model.Materials.Length())
	prevVerticesCount := 0

	for v := range model.Materials.Iterator() {
		i := v.Index
		material := v.Value
		// テクスチャの選択
		var texGl *textureGl
		if material.TextureIndex != -1 && model.Textures.Contains(material.TextureIndex) {
			texGl = renderModel.textures[material.TextureIndex]
		}

		var toonTexGl *textureGl
		if material.ToonSharingFlag == pmx.TOON_SHARING_INDIVIDUAL &&
			material.ToonTextureIndex != -1 &&
			model.Textures.Contains(material.ToonTextureIndex) {
			// 個別Toon
			toonTexGl = renderModel.textures[material.ToonTextureIndex]
		} else if material.ToonSharingFlag == pmx.TOON_SHARING_SHARING &&
			material.ToonTextureIndex != -1 {
			// 共有Toon
			toonTexGl = renderModel.toonTextures[material.ToonTextureIndex]
		}

		var sphereTexGl *textureGl
		if material.SphereMode != pmx.SPHERE_MODE_INVALID &&
			material.SphereTextureIndex != -1 &&
			model.Textures.Contains(material.SphereTextureIndex) {
			sphereTexGl = renderModel.textures[material.SphereTextureIndex]
		}

		materialGl := &materialGL{
			Material:          material,
			texture:           texGl,
			sphereTexture:     sphereTexGl,
			toonTexture:       toonTexGl,
			prevVerticesCount: prevVerticesCount,
		}

		renderModel.meshes[i] = newMesh(faces, materialGl, prevVerticesCount)
		prevVerticesCount += material.VerticesCount
	}
}

// createVertexBuffer はバッファデータから頂点バッファを作成するヘルパーメソッド
func (renderModel *RenderModel) createVertexBuffer(data []float32, usage rendering.BufferUsage) rendering.IVertexBuffer {
	buffer := rendering.IVertexBuffer(mgl.NewVertexBuffer())
	if len(data) > 0 {
		buffer.BufferData(len(data)*4, gl.Ptr(data), usage)
	}
	return buffer
}

// createElementBuffer はインデックスデータから要素バッファを作成するヘルパーメソッド
func (renderModel *RenderModel) createElementBuffer(data []uint32, usage rendering.BufferUsage) rendering.IElementBuffer {
	buffer := rendering.IElementBuffer(mgl.NewElementBuffer())
	if len(data) > 0 {
		buffer.BufferData(len(data)*4, gl.Ptr(data), usage)
	}
	return buffer
}
