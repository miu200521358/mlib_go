//go:build windows
// +build windows

package animation

import (
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl/buffer"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

type RenderModel struct {
	hash              string       // ハッシュ
	meshes            []*Mesh      // メッシュ
	textures          []*textureGl // テクスチャ
	toonTextures      []*textureGl // トゥーンテクスチャ
	vertices          []float32    // 頂点情報
	vao               *buffer.VAO  // 頂点VAO
	vbo               *buffer.VBO  // 頂点VBO
	normalVertices    []float32    // 法線情報
	normalVao         *buffer.VAO  // 法線VAO
	normalVbo         *buffer.VBO  // 法線VBO
	normalIbo         *buffer.IBO  // 法線IBO
	selectedVertexVao *buffer.VAO  // 選択頂点VAO
	selectedVertexVbo *buffer.VBO  // 選択頂点VBO
	selectedVertexIbo *buffer.IBO  // 選択頂点IBO
	boneLineCount     int          // ボーン情報
	boneLineVao       *buffer.VAO  // ボーンVAO
	boneLineVbo       *buffer.VBO  // ボーンVBO
	boneLineIbo       *buffer.IBO  // ボーンIBO
	boneLineIndexes   []int        // ボーンインデックス
	bonePointCount    int          // ボーン情報
	bonePointVao      *buffer.VAO  // ボーンVAO
	bonePointVbo      *buffer.VBO  // ボーンVBO
	bonePointIbo      *buffer.IBO  // ボーンIBO
	bonePointIndexes  []int        // ボーンインデックス
	ssbo              uint32       // SSBO
	vertexCount       int          // 頂点数
}

func NewRenderModel(windowIndex int, model *pmx.PmxModel) *RenderModel {
	m := &RenderModel{
		vertexCount: model.Vertices.Len(),
	}
	m.initToonTexturesGl(windowIndex)
	m.initTexturesGl(windowIndex, model.Textures, model.Path())

	m.initializeBuffer(model)
	m.hash = model.Hash()

	return m
}

func (renderModel *RenderModel) Hash() string {
	return renderModel.hash
}

func (renderModel *RenderModel) Delete() {
	renderModel.vao.Delete()
	renderModel.vbo.Delete()
	renderModel.normalVao.Delete()
	renderModel.normalVbo.Delete()
	renderModel.normalIbo.Delete()
	renderModel.selectedVertexVao.Delete()
	renderModel.selectedVertexVbo.Delete()
	renderModel.selectedVertexIbo.Delete()
	renderModel.boneLineVao.Delete()
	renderModel.boneLineVbo.Delete()
	renderModel.boneLineIbo.Delete()
	gl.DeleteBuffers(1, &renderModel.ssbo)
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

func (renderModel *RenderModel) initializeBuffer(
	model *pmx.PmxModel,
) {
	// 頂点情報
	renderModel.vertices = make([]float32, 0, len(model.Vertices.Data))
	renderModel.normalVertices = make([]float32, 0, len(model.Vertices.Data)*2)
	normalFaces := make([]uint32, 0, len(model.Vertices.Data)*2)
	selectedVertices := make([]float32, 0, len(model.Vertices.Data))
	selectedVertexFaces := make([]uint32, 0, len(model.Vertices.Data))

	// WaitGroupを用いて並列処理を管理
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 頂点情報の並列処理
	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for i, vertex := range model.Vertices.Data {
			vgl := newVertexGl(vertex)

			mu.Lock()
			renderModel.vertices = append(renderModel.vertices, vgl...)

			// 法線
			renderModel.normalVertices = append(renderModel.normalVertices, vgl...)
			renderModel.normalVertices = append(renderModel.normalVertices, newVertexNormalGl(vertex)...)
			normalFaces = append(normalFaces, uint32(n), uint32(n+1))

			// 選択頂点
			selectedVertices = append(selectedVertices, newSelectedVertexGl(vertex)...)
			selectedVertexFaces = append(selectedVertexFaces, uint32(i))

			mu.Unlock()
			n += 2
		}
	}()

	// 面情報の並列処理
	faces := make([]uint32, 0, len(model.Faces.Data)*3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for _, face := range model.Faces.Data {
			vertices := face.VertexIndexes
			mu.Lock()
			faces = append(faces, uint32(vertices[2]), uint32(vertices[1]), uint32(vertices[0]))
			mu.Unlock()

			n += 3
		}
	}()

	// メッシュ情報の並列処理
	renderModel.meshes = make([]*Mesh, len(model.Materials.Data))
	prevVerticesCount := 0

	for i, m := range model.Materials.Data {
		// テクスチャ
		var texGl *textureGl
		if m.TextureIndex != -1 && model.Textures.Contains(m.TextureIndex) {
			texGl = renderModel.textures[m.TextureIndex]
		}

		var toonTexGl *textureGl
		if m.ToonSharingFlag == pmx.TOON_SHARING_INDIVIDUAL &&
			m.ToonTextureIndex != -1 &&
			model.Textures.Contains(m.ToonTextureIndex) {
			// 個別Toon
			toonTexGl = renderModel.textures[m.ToonTextureIndex]
		} else if m.ToonSharingFlag == pmx.TOON_SHARING_SHARING &&
			m.ToonTextureIndex != -1 {
			// 共有Toon
			toonTexGl = renderModel.toonTextures[m.ToonTextureIndex]
		}

		var sphereTexGl *textureGl
		if m.SphereMode != pmx.SPHERE_MODE_INVALID &&
			m.SphereTextureIndex != -1 &&
			model.Textures.Contains(m.SphereTextureIndex) {
			sphereTexGl = renderModel.textures[m.SphereTextureIndex]
		}

		materialGl := &materialGL{
			Material:          m,
			texture:           texGl,
			sphereTexture:     sphereTexGl,
			toonTexture:       toonTexGl,
			prevVerticesCount: prevVerticesCount,
		}

		mesh := newMesh(
			faces,
			materialGl,
			prevVerticesCount,
		)

		renderModel.meshes[i] = mesh
		prevVerticesCount += m.VerticesCount
	}

	// ボーン情報の並列処理
	bones := make([]float32, 0)
	boneFaces := make([]uint32, 0, len(model.Bones.Data)*2)
	boneIndexes := make([]int, len(model.Bones.Data)*2)

	bonePoints := make([]float32, 0)
	bonePointFaces := make([]uint32, len(model.Bones.Data))
	bonePointIndexes := make([]int, len(model.Bones.Data))

	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for _, bone := range model.Bones.Data {
			bones = append(bones, newBoneGl(bone)...)
			bones = append(bones, newTailBoneGl(bone)...)
			boneFaces = append(boneFaces, uint32(n), uint32(n+1))
			boneIndexes[n] = bone.Index()
			boneIndexes[n+1] = bone.Index()

			bonePoints = append(bonePoints, newBoneGl(bone)...)
			bonePointFaces[bone.Index()] = uint32(bone.Index())
			bonePointIndexes[bone.Index()] = bone.Index()

			n += 2

			// if bone.ParentIndex >= 0 && model.Bones.Contains(bone.ParentIndex) &&
			// 	!model.Bones.Get(bone.ParentIndex).Position.IsZero() {
			// 	mu.Lock()
			// 	renderModel.bones = append(renderModel.bones, newBoneGl(bone)...)
			// 	renderModel.bones = append(renderModel.bones, newParentBoneGl(bone)...)
			// 	boneFaces = append(boneFaces, uint32(n), uint32(n+1))
			// 	renderModel.boneIndexes = append(renderModel.boneIndexes, bone.Index, bone.ParentIndex)
			// 	mu.Unlock()

			// 	n += 2
			// }
		}
	}()

	// WaitGroupの完了を待つ
	wg.Wait()

	// GLオブジェクトの生成は並列化しない
	renderModel.vao = buffer.NewVAO()
	renderModel.vao.Bind()

	renderModel.vbo = buffer.NewVBOForVertex(gl.Ptr(renderModel.vertices), len(renderModel.vertices))
	renderModel.vbo.BindVertex(nil, nil)
	renderModel.vbo.Unbind()
	renderModel.vao.Unbind()

	renderModel.normalVao = buffer.NewVAO()
	renderModel.normalVao.Bind()
	renderModel.normalVbo = buffer.NewVBOForVertex(gl.Ptr(renderModel.normalVertices), len(renderModel.normalVertices))
	renderModel.normalVbo.BindVertex(nil, nil)
	renderModel.normalIbo = buffer.NewIBO(gl.Ptr(normalFaces), len(normalFaces))
	renderModel.normalIbo.Bind()
	renderModel.normalIbo.Unbind()
	renderModel.normalVbo.Unbind()
	renderModel.normalVao.Unbind()

	renderModel.boneLineCount = len(bones)
	renderModel.boneLineIndexes = boneIndexes
	renderModel.boneLineVao = buffer.NewVAO()
	renderModel.boneLineVao.Bind()
	renderModel.boneLineVbo = buffer.NewVBOForBone(gl.Ptr(bones), len(bones))
	renderModel.boneLineVbo.BindBone(nil, nil)
	renderModel.boneLineIbo = buffer.NewIBO(gl.Ptr(boneFaces), len(boneFaces))
	renderModel.boneLineIbo.Bind()
	renderModel.boneLineIbo.Unbind()
	renderModel.boneLineVbo.Unbind()
	renderModel.boneLineVao.Unbind()

	renderModel.bonePointCount = len(bonePoints)
	renderModel.bonePointIndexes = bonePointIndexes
	renderModel.bonePointVao = buffer.NewVAO()
	renderModel.bonePointVao.Bind()
	renderModel.bonePointVbo = buffer.NewVBOForBone(gl.Ptr(bonePoints), len(bonePoints))
	renderModel.bonePointVbo.BindBone(nil, nil)
	renderModel.bonePointIbo = buffer.NewIBO(gl.Ptr(bonePointFaces), len(bonePointFaces))
	renderModel.bonePointIbo.Bind()
	renderModel.bonePointIbo.Unbind()
	renderModel.bonePointVbo.Unbind()
	renderModel.bonePointVao.Unbind()

	renderModel.selectedVertexVao = buffer.NewVAO()
	renderModel.selectedVertexVao.Bind()
	renderModel.selectedVertexVbo = buffer.NewVBOForVertex(gl.Ptr(selectedVertices), len(selectedVertices))
	renderModel.selectedVertexVbo.BindVertex(nil, nil)
	renderModel.selectedVertexIbo = buffer.NewIBO(gl.Ptr(selectedVertexFaces), len(selectedVertexFaces))
	renderModel.selectedVertexIbo.Bind()
	renderModel.selectedVertexIbo.Unbind()
	renderModel.selectedVertexVbo.Unbind()
	renderModel.selectedVertexVao.Unbind()

	// SSBOの作成
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, model.Vertices.Len()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)

	renderModel.ssbo = ssbo
}

func (renderModel *RenderModel) Render(
	shader mgl.IShader, appState state.IAppState, animationState state.IAnimationState,
) {
	deltas := animationState.VmdDeltas()

	renderModel.vao.Bind()
	defer renderModel.vao.Unbind()

	renderModel.vbo.BindVertex(
		animationState.RenderDeltas().VertexMorphDeltaIndexes, animationState.RenderDeltas().VertexMorphDeltas)
	defer renderModel.vbo.Unbind()

	paddedMatrixes, matrixWidth, matrixHeight, err := createBoneMatrixes(deltas.Bones)
	if err != nil {
		return
	}

	for i, mesh := range renderModel.meshes {
		if mesh == nil || animationState.RenderDeltas() == nil || len(animationState.RenderDeltas().MeshDeltas) <= i {
			continue
		}

		mesh.ibo.Bind()

		mesh.drawModel(animationState.WindowIndex(), shader, paddedMatrixes, matrixWidth, matrixHeight,
			animationState.RenderDeltas().MeshDeltas[i])

		if mesh.material.DrawFlag.IsDrawingEdge() {
			// エッジ描画
			mesh.drawEdge(animationState.WindowIndex(), shader, paddedMatrixes,
				matrixWidth, matrixHeight, animationState.RenderDeltas().MeshDeltas[i])
		}

		if appState.IsShowWire() {
			mesh.drawWire(animationState.WindowIndex(), shader, paddedMatrixes, matrixWidth, matrixHeight, false)
		}

		mesh.ibo.Unbind()
	}

	if appState.IsShowNormal() {
		renderModel.drawNormal(animationState.WindowIndex(), shader, paddedMatrixes, matrixWidth, matrixHeight)
	}

	if appState.IsShowBoneAll() || appState.IsShowBoneEffector() || appState.IsShowBoneIk() ||
		appState.IsShowBoneFixed() || appState.IsShowBoneRotate() ||
		appState.IsShowBoneTranslate() || appState.IsShowBoneVisible() {
		renderModel.drawBone(animationState.WindowIndex(),
			shader, animationState.Model().Bones, appState,
			paddedMatrixes, matrixWidth, matrixHeight)
	}

}

func (renderModel *RenderModel) drawNormal(
	windowIndex int,
	shader mgl.IShader,
	paddedMatrixes []float32,
	width, height int,
) {
	program := shader.Program(mgl.PROGRAM_TYPE_NORMAL)
	gl.UseProgram(program)

	renderModel.normalVao.Bind()
	renderModel.normalVbo.BindVertex(nil, nil)
	renderModel.normalIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)

	windowOpacityUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_WINDOW_OPACITY))
	gl.Uniform1f(windowOpacityUniform, shader.WindowOpacity())

	normalColor := mgl32.Vec4{0.3, 0.3, 0.7, 0.5}
	specularUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &normalColor[0])

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(len(renderModel.normalVertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.normalIbo.Unbind()
	renderModel.normalVbo.Unbind()
	renderModel.normalVao.Unbind()

	gl.UseProgram(0)
}

func (renderModel *RenderModel) drawSelectedVertex(
	windowIndex int,
	selectedVertexMorphIndexes []int,
	selectedVertexDeltas [][]float32,
	shader mgl.IShader,
	paddedMatrixes []float32,
	width, height int,
) [][]float32 {
	// モデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(mgl.PROGRAM_TYPE_SELECTED_VERTEX)
	gl.UseProgram(program)

	renderModel.selectedVertexVao.Bind()
	renderModel.selectedVertexVbo.BindVertex(selectedVertexMorphIndexes, selectedVertexDeltas)
	renderModel.selectedVertexIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)

	windowOpacityUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_WINDOW_OPACITY))
	gl.Uniform1f(windowOpacityUniform, shader.WindowOpacity())

	vertexColor := mgl32.Vec4{1.0, 0.4, 0.0, 0.7}
	specularUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &vertexColor[0])
	gl.PointSize(5.0) // 選択頂点のサイズ

	// 点描画
	gl.DrawElements(
		gl.POINTS,
		int32(len(renderModel.vertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	// SSBOからデータを読み込む
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, renderModel.ssbo)
	ptr := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_ONLY)
	results := make([][]float32, renderModel.vertexCount)

	if ptr != nil {
		// 頂点数を取得
		for i := range results {
			results[i] = make([]float32, 4)
		}

		// SSBOから読み取り
		for i := 0; i < renderModel.vertexCount; i++ {
			for j := 0; j < 4; j++ {
				results[i][j] = *(*float32)(unsafe.Pointer(uintptr(ptr) + uintptr((i*4+j)*4)))
			}
		}
	}

	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	renderModel.selectedVertexIbo.Unbind()
	renderModel.selectedVertexVbo.Unbind()
	renderModel.selectedVertexVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	return results
}

func (renderModel *RenderModel) drawBone(
	windowIndex int,
	shader mgl.IShader,
	bones *pmx.Bones,
	appState state.IAppState,
	paddedMatrixes []float32,
	width, height int,
) {
	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(mgl.PROGRAM_TYPE_BONE)
	gl.UseProgram(program)

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)

	windowOpacityUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_WINDOW_OPACITY))
	gl.Uniform1f(windowOpacityUniform, shader.WindowOpacity())

	renderModel.boneLineVao.Bind()
	renderModel.boneLineVbo.BindBone(renderModel.fetchBoneLineDeltas(bones, appState))
	renderModel.boneLineIbo.Bind()

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(renderModel.boneLineCount),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.boneLineIbo.Unbind()
	renderModel.boneLineVbo.Unbind()
	renderModel.boneLineVao.Unbind()

	renderModel.bonePointVao.Bind()
	renderModel.bonePointVbo.BindBone(renderModel.fetchBonePointDeltas(bones, appState))
	renderModel.bonePointIbo.Bind()

	gl.PointSize(5.0) // 選択頂点のサイズ

	// 点描画
	gl.DrawElements(
		gl.POINTS,
		int32(len(renderModel.vertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.bonePointIbo.Unbind()
	renderModel.bonePointVbo.Unbind()
	renderModel.bonePointVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

}

func (renderModel *RenderModel) fetchBoneLineDeltas(
	bones *pmx.Bones, appState state.IAppState,
) ([]int, [][]float32) {
	indexes := make([]int, len(renderModel.boneLineIndexes))
	deltas := make([][]float32, len(renderModel.boneLineIndexes))

	for i, boneIndex := range renderModel.boneLineIndexes {
		indexes[i] = i
		deltas[i] = getBoneDebugColor(bones.Get(boneIndex), appState)
	}

	return indexes, deltas
}

func (renderModel *RenderModel) fetchBonePointDeltas(
	bones *pmx.Bones, appState state.IAppState,
) ([]int, [][]float32) {
	indexes := make([]int, len(renderModel.bonePointIndexes))
	deltas := make([][]float32, len(renderModel.bonePointIndexes))

	for i, boneIndex := range renderModel.bonePointIndexes {
		indexes[i] = i
		deltas[i] = getBoneDebugColor(bones.Get(boneIndex), appState)
	}

	return indexes, deltas
}
