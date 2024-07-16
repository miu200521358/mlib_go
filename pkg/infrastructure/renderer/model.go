package renderer

import (
	"slices"
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl/buffer"
)

type RenderModel struct {
	Initialized       bool         // 描画初期化済みフラグ
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
	bones             []float32    // ボーン情報
	boneVao           *buffer.VAO  // ボーンVAO
	boneVbo           *buffer.VBO  // ボーンVBO
	boneIbo           *buffer.IBO  // ボーンIBO
	boneIndexes       []int        // ボーンインデックス
	ssbo              uint32       // SSBO
	vertexCount       int          // 頂点数
}

func NewRenderModel(windowIndex int, model *pmx.PmxModel) *RenderModel {
	m := &RenderModel{
		Initialized: false,
		vertexCount: model.Vertices.Len(),
	}
	m.initToonTexturesGl(windowIndex)
	m.initTexturesGl(windowIndex, model.Textures, model.GetPath())

	m.initializeBuffer(model)

	return m
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
	renderModel.boneVao.Delete()
	renderModel.boneVbo.Delete()
	renderModel.boneIbo.Delete()
	gl.DeleteBuffers(1, &renderModel.ssbo)
	for _, mesh := range renderModel.meshes {
		mesh.delete()
	}
	for _, texture := range renderModel.textures {
		texture.delete()
	}
	for _, texture := range renderModel.toonTextures {
		texture.delete()
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
	renderModel.bones = make([]float32, 0, len(model.Bones.Data)*4)
	boneFaces := make([]uint32, 0, len(model.Bones.Data)*4)
	renderModel.boneIndexes = make([]int, 0, len(model.Bones.Data)*4)

	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for _, bone := range model.Bones.Data {
			mu.Lock()
			renderModel.bones = append(renderModel.bones, newBoneGl(bone)...)
			renderModel.bones = append(renderModel.bones, newTailBoneGl(bone)...)
			boneFaces = append(boneFaces, uint32(n), uint32(n+1))
			renderModel.boneIndexes = append(renderModel.boneIndexes, bone.Index, bone.Index)
			mu.Unlock()

			n += 2

			if bone.ParentIndex >= 0 && model.Bones.Contains(bone.ParentIndex) &&
				!model.Bones.Get(bone.ParentIndex).Position.IsZero() {
				mu.Lock()
				renderModel.bones = append(renderModel.bones, newBoneGl(bone)...)
				renderModel.bones = append(renderModel.bones, newParentBoneGl(bone)...)
				boneFaces = append(boneFaces, uint32(n), uint32(n+1))
				renderModel.boneIndexes = append(renderModel.boneIndexes, bone.Index, bone.ParentIndex)
				mu.Unlock()

				n += 2
			}
		}
	}()

	// WaitGroupの完了を待つ
	wg.Wait()

	// 以下の部分は並列化する必要がないのでそのままにする
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

	renderModel.boneVao = buffer.NewVAO()
	renderModel.boneVao.Bind()
	renderModel.boneVbo = buffer.NewVBOForVertex(gl.Ptr(renderModel.bones), len(renderModel.bones))
	renderModel.boneVbo.BindVertex(nil, nil)
	renderModel.boneIbo = buffer.NewIBO(gl.Ptr(boneFaces), len(boneFaces))
	renderModel.boneIbo.Bind()
	renderModel.boneIbo.Unbind()
	renderModel.boneVbo.Unbind()
	renderModel.boneVao.Unbind()

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
	shader *mgl.MShader, animationStates *AnimationStates, windowIndex int,
	isShowNormal, isShowWire, isShowSelectedVertex bool, isShowBones map[pmx.BoneFlag]bool,
) {
	deltas := animationStates.Now.VmdDeltas

	renderModel.vao.Bind()
	defer renderModel.vao.Unbind()

	renderModel.vbo.BindVertex(animationStates.Now.vertexMorphDeltaIndexes, animationStates.Now.vertexMorphDeltas)
	defer renderModel.vbo.Unbind()

	paddedMatrixes, matrixWidth, matrixHeight := createBoneMatrixes(deltas.Bones)

	for i, mesh := range renderModel.meshes {
		mesh.ibo.Bind()

		mesh.drawModel(shader, paddedMatrixes, matrixWidth, matrixHeight, animationStates.Now.meshDeltas[i])

		if mesh.material.DrawFlag.IsDrawingEdge() {
			// エッジ描画
			mesh.drawEdge(shader, paddedMatrixes, matrixWidth, matrixHeight, animationStates.Now.meshDeltas[i])
		}

		if isShowWire {
			mesh.drawWire(shader, paddedMatrixes, matrixWidth, matrixHeight, ((len(animationStates.Now.InvisibleMaterialIndexes) > 0 && len(animationStates.Next.InvisibleMaterialIndexes) == 0 &&
				slices.Contains(animationStates.Now.InvisibleMaterialIndexes, mesh.material.Index)) ||
				slices.Contains(animationStates.Next.InvisibleMaterialIndexes, mesh.material.Index)))
		}

		mesh.ibo.Unbind()
	}

	if isShowNormal {
		renderModel.drawNormal(shader, paddedMatrixes, matrixWidth, matrixHeight)
	}

	isShowBone := false
	for _, drawBone := range isShowBones {
		if drawBone {
			isShowBone = true
			break
		}
	}

	if isShowBone {
		renderModel.drawBone(shader, animationStates.Now.Model.Bones, isShowBones,
			paddedMatrixes, matrixWidth, matrixHeight)
	}

}

func (renderModel *RenderModel) drawNormal(
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
) {
	program := shader.GetProgram(mgl.PROGRAM_TYPE_NORMAL)
	gl.UseProgram(program)

	renderModel.normalVao.Bind()
	renderModel.normalVbo.BindVertex(nil, nil)
	renderModel.normalIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)

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
	selectedVertexMorphIndexes []int,
	selectedVertexDeltas [][]float32,
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
) [][]float32 {
	// モデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.GetProgram(mgl.PROGRAM_TYPE_SELECTED_VERTEX)
	gl.UseProgram(program)

	renderModel.selectedVertexVao.Bind()
	renderModel.selectedVertexVbo.BindVertex(selectedVertexMorphIndexes, selectedVertexDeltas)
	renderModel.selectedVertexIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)

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
	shader *mgl.MShader,
	bones *pmx.Bones,
	isShowBones map[pmx.BoneFlag]bool,
	paddedMatrixes []float32,
	width, height int,
) {
	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	program := shader.GetProgram(mgl.PROGRAM_TYPE_BONE)
	gl.UseProgram(program)

	renderModel.boneVao.Bind()
	renderModel.boneVbo.BindVertex(renderModel.fetchBoneDebugDeltas(bones, isShowBones))
	renderModel.boneIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(len(renderModel.bones)),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.boneIbo.Unbind()
	renderModel.boneVbo.Unbind()
	renderModel.boneVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

}

func (renderModel *RenderModel) fetchBoneDebugDeltas(
	bones *pmx.Bones, isShowBones map[pmx.BoneFlag]bool,
) ([]int, [][]float32) {
	indexes := make([]int, 0)
	deltas := make([][]float32, 0)

	for _, boneIndex := range renderModel.boneIndexes {
		bone := bones.Get(boneIndex)
		indexes = append(indexes, boneIndex)
		deltas = append(deltas, newBoneDebugAlphaGl(bone, isShowBones))
	}

	return indexes, deltas
}
