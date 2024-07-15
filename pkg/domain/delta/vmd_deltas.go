package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type VmdDeltas struct {
	Bones  *BoneDeltas
	Morphs *MorphDeltas
}

func NewVmdDeltas(vertices *pmx.Vertices) *VmdDeltas {
	return &VmdDeltas{}
}
