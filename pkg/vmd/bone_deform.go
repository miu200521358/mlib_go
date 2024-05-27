package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type boneDeforms struct {
	deforms     map[int]*BoneDeform
	names       map[string]int
	boneIndexes []int
}

type BoneDeform struct {
	bone           *pmx.Bone
	position       *mmath.MVec3
	effectPosition *mmath.MVec3
	rotation       *mmath.MQuaternion
	effectRotation *mmath.MQuaternion
	scale          *mmath.MVec3
	unitMatrix     *mmath.MMat4
}

func getBoneDeform(boneDeformsMap map[bool]*boneDeforms, bone *pmx.Bone) *BoneDeform {
	if _, ok := boneDeformsMap[false].deforms[bone.Index]; ok {
		return boneDeformsMap[false].deforms[bone.Index]
	} else if _, ok := boneDeformsMap[true].deforms[bone.Index]; ok {
		return boneDeformsMap[true].deforms[bone.Index]
	}
	return nil
}
