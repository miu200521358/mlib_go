package pmx

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mgl"
)

type Meshes struct {
	meshes []*Mesh
	vao    *mgl.VAO
	vbo    *mgl.VBO
	ibo    *mgl.IBO
}

func NewMeshes(model *PmxModel, windowIndex int) *Meshes {
	// 頂点情報
	vertices := make([]float32, 0, len(model.Vertices.Indexes))
	for _, v := range model.Vertices.Data {
		vertices = append(vertices, (*v).GL()...)
	}

	var faceDtype uint32

	if model.VertexCountType == 1 {
		faceDtype = uint32(8)
	} else if model.VertexCountType == 2 {
		faceDtype = uint32(16)
	} else {
		faceDtype = uint32(32)
	}

	// 面情報
	faces := make([]uint32, 0, len(model.Faces.Indexes))
	for _, f := range model.Faces.Data {
		vertices := (*f).VertexIndexes
		faces = append(faces, uint32(vertices[2]))
		faces = append(faces, uint32(vertices[1]))
		faces = append(faces, uint32(vertices[0]))
	}

	meshes := make([]*Mesh, 0, len(model.Materials.Indexes))
	prevVerticesCount := 0
	for _, m := range model.Materials.Data {
		materialGl := (*m).GL(
			model.Path,
			nil,
			nil,
			nil,
			windowIndex,
		)
		mesh := NewMesh(
			materialGl,
			prevVerticesCount,
		)
		meshes = append(meshes, mesh)
	}

	vao := mgl.NewVAO()
	vao.Bind()
	vbo := mgl.NewVBO(gl.Ptr(vertices), len(vertices))
	vbo.Bind()
	ibo := mgl.NewIBO(gl.Ptr(faces), len(faces), faceDtype)
	ibo.Bind()
	ibo.Unbind()
	vbo.Unbind()
	vao.Unbind()

	return &Meshes{
		meshes: meshes,
		vao:    vao,
		vbo:    vbo,
		ibo:    ibo,
	}
}

func (m *Meshes) Delete() {
	m.ibo.Delete()
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
		m.ibo.Bind()

		shader.UseModelProgram()
		mesh.DrawModel(shader, m.ibo.FaceCount, m.ibo.Dtype, windowIndex, boneMatrixes)
		shader.Unuse()

		m.ibo.Unbind()
		m.vbo.Unbind()
		m.vao.Unbind()
	}
}
