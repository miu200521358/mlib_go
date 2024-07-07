package vmd

type VmdDeltas struct {
	Vertices             *VertexDeltas
	Bones                *BoneDeltas
	Morphs               *MorphDeltas
	SelectedVertexDeltas *SelectedVertexMorphDeltas
}

func NewVmdDeltas() *VmdDeltas {
	return &VmdDeltas{
		Vertices:             NewVertexDeltas(),
		SelectedVertexDeltas: NewSelectedVertexMorphDeltas(),
	}
}
