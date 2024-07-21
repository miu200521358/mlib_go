package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type VmdDeltas struct {
	Bones  *BoneDeltas
	Morphs *MorphDeltas
}

func NewVmdDeltas(materials *pmx.Materials, bones *pmx.Bones) *VmdDeltas {
	return &VmdDeltas{
		Bones: NewBoneDeltas(bones),
	}
}
