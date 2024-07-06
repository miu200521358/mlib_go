package vmd

type VmdDeltas struct {
	Vertices *VertexDeltas
	Bones    *BoneDeltas
	Morphs   *MorphDeltas
}

func NewVmdDeltas() *VmdDeltas {
	return &VmdDeltas{
		Vertices: NewVertexDeltas(),
	}
}
