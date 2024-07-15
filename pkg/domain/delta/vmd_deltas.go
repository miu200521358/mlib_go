package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type VmdDeltas struct {
	Vertices             *VertexDeltas
	Bones                *BoneDeltas
	Morphs               *MorphDeltas
	SelectedVertexDeltas *SelectedVertexMorphDeltas
}

func NewVmdDeltas(vertices *pmx.Vertices) *VmdDeltas {
	return &VmdDeltas{
		Vertices:             NewVertexDeltas(vertices),
		SelectedVertexDeltas: NewSelectedVertexMorphDeltas(),
	}
}
