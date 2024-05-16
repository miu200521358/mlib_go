//go:build windows
// +build windows

package pmx

import (
	"embed"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mview"
)

type Meshes struct {
	meshes   []*Mesh
	vertices []float32
	vao      *mview.VAO
	vbo      *mview.VBO
}

func NewMeshes(
	model *PmxModel,
	windowIndex int,
	resourceFiles embed.FS,
) *Meshes {
	// 頂点情報
	vertices := make([]float32, 0, len(model.Vertices.Data))
	for i := range len(model.Vertices.Data) {
		vertices = append(vertices, model.Vertices.GetItem(i).GL()...)
	}
	// println("vertices", mutils.JoinSlice(mutils.ConvertFloat32ToInterfaceSlice(vertices)))

	// 面情報
	faces := make([]uint32, 0, len(model.Faces.Data)*3)
	for i := range len(model.Faces.Data) {
		vertices := model.Faces.GetItem(i).VertexIndexes
		faces = append(faces, uint32(vertices[2]))
		faces = append(faces, uint32(vertices[1]))
		faces = append(faces, uint32(vertices[0]))
	}

	meshes := make([]*Mesh, len(model.Materials.GetIndexes()))
	prevVerticesCount := 0
	for i, m := range model.Materials.GetSortedData() {
		var texture *Texture
		if m.TextureIndex != -1 && model.Textures.Contains(m.TextureIndex) {
			texture = model.Textures.GetItem(m.TextureIndex)
		}

		var toonTexture *Texture
		// 個別Toon
		if m.ToonSharingFlag == TOON_SHARING_INDIVIDUAL &&
			m.ToonTextureIndex != -1 &&
			model.Textures.Contains(m.ToonTextureIndex) {
			toonTexture = model.Textures.GetItem(m.ToonTextureIndex)
		}
		// 共有Toon
		if m.ToonSharingFlag == TOON_SHARING_SHARING &&
			m.ToonTextureIndex != -1 &&
			model.ToonTextures.Contains(m.ToonTextureIndex) {
			toonTexture = model.ToonTextures.GetItem(m.ToonTextureIndex)
		}

		var sphereTexture *Texture
		if m.SphereMode != SPHERE_MODE_INVALID &&
			m.SphereTextureIndex != -1 &&
			model.Textures.Contains(m.SphereTextureIndex) {
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
		meshes[i] = mesh

		prevVerticesCount += m.VerticesCount
	}

	vao := mview.NewVAO()
	vao.Bind()
	vbo := mview.NewVBOForVertex(gl.Ptr(vertices), len(vertices))
	vbo.BindVertex(nil, nil)
	vbo.Unbind()
	vao.Unbind()

	return &Meshes{
		meshes:   meshes,
		vertices: vertices,
		vao:      vao,
		vbo:      vbo,
	}
}

func (m *Meshes) delete() {
	for _, mesh := range m.meshes {
		mesh.delete()
	}
	m.vao.Delete()
	m.vbo.Delete()
}

func (m *Meshes) Draw(
	shader *mview.MShader,
	boneMatrixes []*mgl32.Mat4,
	vertexDeltas [][]float32,
	materialDeltas []*Material,
	windowIndex int,
) {
	// 隠面消去
	// https://learnopengl.com/Advanced-OpenGL/Depth-testing
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	m.vao.Bind()
	m.vbo.BindVertex(m.vertices, vertexDeltas)

	for i, mesh := range m.meshes {
		shader.UseModelProgram()
		mesh.DrawModel(shader, windowIndex, boneMatrixes, materialDeltas[i])
		shader.Unuse()
	}

	m.vbo.Unbind()
	m.vao.Unbind()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.DEPTH_TEST)
}
