package pmx

import "github.com/miu200521358/mlib_go/pkg/mgl"

type Meshes struct {
	model     *PmxModel
	vertices  []*VertexGL
	faces     []*FaceGL
	materials []*MaterialGL
	vao       *mgl.VAO
	vbo       *mgl.VBO
}

func NewMeshes(model *PmxModel, windowIndex int) *Meshes {
	// 頂点情報
	vertices := make([]*VertexGL, 0, len(model.Vertices.Indexes))
	for _, v := range model.Vertices.Data {
		vertices = append(vertices, (*v).GL())
	}

	var faceDtype uint32

	if model.VertexCount == 1 {
		faceDtype = uint32(8)
	} else if model.VertexCount == 2 {
		faceDtype = uint32(16)
	} else {
		faceDtype = uint32(32)
	}

	// 面情報
	faces := make([]*FaceGL, 0, len(model.Faces.Indexes))
	for _, f := range model.Faces.Data {
		faces = append(faces, (*f).GL())
	}

	materials := make([]*MaterialGL, 0, len(model.Materials.Indexes))
	for _, m := range model.Materials.Data {
		materials = append(materials, (*m).GL(
			model.Path,
			nil,
			nil,
			nil,
			windowIndex,
		))
	}

	vao := mgl.NewVAO()
	vbo := mgl.NewVBO(faceDtype)

	meshes := &Meshes{
		model:     model,
		vertices:  vertices,
		faces:     faces,
		materials: materials,
		vao:       vao,
		vbo:       vbo,
	}

	return meshes
}
