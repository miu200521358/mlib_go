package vmd

import "github.com/miu200521358/mlib_go/pkg/pmx"

type VmdMotion struct {
	BoneFrames *BoneFrames
}

func NewVmdMotion() *VmdMotion {
	return &VmdMotion{
		BoneFrames: NewBoneFrames(),
	}
}

func (m *VmdMotion) Animate(fno int, model pmx.PmxModel) BoneTrees {
	return m.AnimateBone([]int{fno}, model, nil, false, false, "")
}

func (m *VmdMotion) AnimateBone(
	fnos []int,
	model pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) BoneTrees {
	// ボーン変形行列操作
	return m.BoneFrames.Animate(fnos, model, boneNames, isCalcIk, isOutLog, description)
}
