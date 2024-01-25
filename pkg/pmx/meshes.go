package pmx

import (
	"embed"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mutils"

)

type Meshes struct {
	meshes []*Mesh
	vao    *mgl.VAO
	vbo    *mgl.VBO
}

func NewMeshes(
	model *PmxModel,
	windowIndex int,
	resourceFiles embed.FS,
) *Meshes {
	// 頂点情報
	vertices := make([]float32, 0, len(model.Vertices.Indexes))
	for _, v := range model.Vertices.GetSortedData() {
		vertices = append(vertices, (*v).GL()...)
	}
	// println("vertices", mutils.JoinSlice(mutils.ConvertFloat32ToInterfaceSlice(vertices)))

	// 面情報
	faces := make([]uint32, 0, len(model.Faces.Indexes))
	for _, f := range model.Faces.GetSortedData() {
		vertices := f.VertexIndexes
		faces = append(faces, uint32(vertices[2]))
		faces = append(faces, uint32(vertices[1]))
		faces = append(faces, uint32(vertices[0]))
	}

	meshes := make([]*Mesh, 0, len(model.Materials.Indexes))
	prevVerticesCount := 0
	for _, m := range model.Materials.GetSortedData() {
		var texture *Texture
		if m.TextureIndex != -1 && mutils.ContainsInt(model.Textures.Indexes, m.TextureIndex) {
			texture = model.Textures.GetItem(m.TextureIndex)
		}

		var toonTexture *Texture
		// 個別Toon
		if m.ToonSharingFlag == TOON_SHARING_INDIVIDUAL &&
			m.ToonTextureIndex != -1 &&
			mutils.ContainsInt(model.Textures.Indexes, m.ToonTextureIndex) {
			toonTexture = model.Textures.GetItem(m.ToonTextureIndex)
		}
		// 共有Toon
		if m.ToonSharingFlag == TOON_SHARING_SHARING &&
			m.ToonTextureIndex != -1 &&
			mutils.ContainsInt(model.ToonTextures.Indexes, m.ToonTextureIndex) {
			toonTexture = model.ToonTextures.GetItem(m.ToonTextureIndex)
		}

		var sphereTexture *Texture
		if m.SphereMode != SPHERE_MODE_INVALID &&
			m.SphereTextureIndex != -1 &&
			mutils.ContainsInt(model.Textures.Indexes, m.SphereTextureIndex) {
			sphereTexture = model.Textures.GetItem(m.SphereTextureIndex)
		}

		materialGl := m.GL(
			model.Path,
			texture,
			toonTexture,
			sphereTexture,
			windowIndex,
			prevVerticesCount,
			resourceFiles,
		)
		mesh := NewMesh(
			faces,
			materialGl,
			prevVerticesCount,
		)
		meshes = append(meshes, mesh)

		prevVerticesCount += m.VerticesCount
	}

	vao := mgl.NewVAO()
	vao.Bind()
	vbo := mgl.NewVBO(gl.Ptr(vertices), len(vertices))
	vbo.Bind()
	vbo.Unbind()
	vao.Unbind()

	return &Meshes{
		meshes: meshes,
		vao:    vao,
		vbo:    vbo,
	}
}

func (m *Meshes) Delete() {
	for _, mesh := range m.meshes {
		mesh.Delete()
	}
	m.vao.Delete()
	m.vbo.Delete()
}

func (m *Meshes) Draw(shader *mgl.MShader, boneMatrixes []mgl32.Mat4, windowIndex int) {
	// 隠面消去
	// https://learnopengl.com/Advanced-OpenGL/Depth-testing
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	for _, mesh := range m.meshes {
		m.vao.Bind()
		m.vbo.Bind()

		// ブレンディングを有効にする
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

		shader.UseModelProgram()
		mesh.DrawModel(shader, windowIndex, boneMatrixes)
		shader.Unuse()

		m.vbo.Unbind()
		m.vao.Unbind()

		gl.Disable(gl.BLEND)
	}

	gl.Disable(gl.DEPTH_TEST)
}
