package renderer

import (
	"sync"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl/buffer"
)

type RenderModel struct {
	Model             *pmx.PmxModel
	Initialized       bool
	meshes            []*Mesh
	vertices          []float32
	vao               *buffer.VAO
	vbo               *buffer.VBO
	normalVertices    []float32
	normalVao         *buffer.VAO
	normalVbo         *buffer.VBO
	normalIbo         *buffer.IBO
	selectedVertexVao *buffer.VAO
	selectedVertexVbo *buffer.VBO
	selectedVertexIbo *buffer.IBO
	bones             []float32
	boneVao           *buffer.VAO
	boneVbo           *buffer.VBO
	boneIbo           *buffer.IBO
	boneIndexes       []int
	ssbo              uint32
	vertexCount       int
}

func NewRenderModel(windowIndex int, model *pmx.PmxModel) *RenderModel {
	m := &RenderModel{
		Model:       model,
		Initialized: false,
	}
	model.ToonTextures = pmx.NewToonTextures()
	initToonTexturesGl(windowIndex, model.ToonTextures)

	m.initialize(model, windowIndex)

	return m
}

func (renderModel *RenderModel) initialize(
	model *pmx.PmxModel,
	windowIndex int,
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

	// テクスチャの gl.GenTextures はスレッドセーフではないので、並列化しない
	for i, m := range model.Materials.Data {
		// テクスチャ
		var texture *pmx.Texture
		if m.TextureIndex != -1 && model.Textures.Contains(m.TextureIndex) {
			texture = model.Textures.Get(m.TextureIndex)
		}

		var toonTexture *pmx.Texture
		// 個別Toon
		if m.ToonSharingFlag == pmx.TOON_SHARING_INDIVIDUAL &&
			m.ToonTextureIndex != -1 &&
			model.Textures.Contains(m.ToonTextureIndex) {
			toonTexture = model.Textures.Get(m.ToonTextureIndex)
		}
		// 共有Toon
		if m.ToonSharingFlag == pmx.TOON_SHARING_SHARING &&
			m.ToonTextureIndex != -1 &&
			model.ToonTextures.Contains(m.ToonTextureIndex) {
			toonTexture = model.ToonTextures.Get(m.ToonTextureIndex)
		}

		var sphereTexture *pmx.Texture
		if m.SphereMode != pmx.SPHERE_MODE_INVALID &&
			m.SphereTextureIndex != -1 &&
			model.Textures.Contains(m.SphereTextureIndex) {
			sphereTexture = model.Textures.Get(m.SphereTextureIndex)
		}

		materialGl := newMaterialGl(
			m,
			model.GetPath(),
			texture,
			toonTexture,
			sphereTexture,
			windowIndex,
			prevVerticesCount,
		)
		mesh := NewMesh(
			faces,
			materialGl,
			prevVerticesCount,
		)
		mu.Lock()
		renderModel.meshes[i] = mesh
		mu.Unlock()

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
	renderModel.vertexCount = model.Vertices.Len()
}
