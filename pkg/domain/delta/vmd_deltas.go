package delta

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type VmdDeltas struct {
	Vertices               *VertexDeltas
	Bones                  *BoneDeltas
	Morphs                 *MorphDeltas
	SelectedVertexDeltas   *SelectedVertexMorphDeltas
	BoneGlDeltas           []mgl32.Mat4
	MeshGlDeltas           []*MeshDelta
	VertexMorphIndexes     []int
	VertexMorphGlDeltas    [][]float32
	SelectedVertexIndexes  []int
	SelectedVertexGlDeltas [][]float32
}

func NewVmdDeltas(vertices *pmx.Vertices) *VmdDeltas {
	return &VmdDeltas{
		Vertices:             NewVertexDeltas(vertices),
		SelectedVertexDeltas: NewSelectedVertexMorphDeltas(),
	}
}
