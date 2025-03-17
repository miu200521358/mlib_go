package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type VmdDeltas struct {
	frame      float32
	modelHash  string
	motionHash string
	Bones      *BoneDeltas
	Morphs     *MorphDeltas
}

func NewVmdDeltas(frame float32, bones *pmx.Bones, modelHash, motionHash string) *VmdDeltas {
	return &VmdDeltas{
		frame:      frame,
		modelHash:  modelHash,
		motionHash: motionHash,
		Bones:      NewBoneDeltas(bones),
		Morphs:     nil,
	}
}

func (vmdDeltas *VmdDeltas) Frame() float32 {
	return vmdDeltas.frame
}

func (vmdDeltas *VmdDeltas) SetFrame(frame float32) {
	vmdDeltas.frame = frame
}

func (vmdDeltas *VmdDeltas) ModelHash() string {
	return vmdDeltas.modelHash
}

func (vmdDeltas *VmdDeltas) SetModelHash(hash string) {
	vmdDeltas.modelHash = hash
}

func (vmdDeltas *VmdDeltas) MotionHash() string {
	return vmdDeltas.motionHash
}

func (vmdDeltas *VmdDeltas) SetMotionHash(hash string) {
	vmdDeltas.motionHash = hash
}
