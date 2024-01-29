package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type BoneTree struct {
	BoneName      string
	FrameNo       int
	GlobalMatrix  *mmath.MMat4
	LocalMatrix   *mmath.MMat4
	Position      *mmath.MVec3
	FramePosition *mmath.MVec3
	FrameRotation *mmath.MQuaternion
	FrameScale    *mmath.MVec3
}

func NewBoneTree(
	boneName string,
	frameNo int,
	globalMatrix, localMatrix *mmath.MMat4,
	framePosition *mmath.MVec3,
	frameRotation *mmath.MQuaternion,
	frameScale *mmath.MVec3,
) *BoneTree {
	p := globalMatrix.Translation()
	return &BoneTree{
		BoneName:      boneName,
		FrameNo:       frameNo,
		GlobalMatrix:  globalMatrix,
		LocalMatrix:   localMatrix,
		Position:      &p,
		FramePosition: framePosition,
		FrameRotation: frameRotation,
		FrameScale:    frameScale,
	}
}

type BoneNameFrameNo struct {
	BoneName string
	FrameNo  int
}

type BoneTrees struct {
	Data map[BoneNameFrameNo]*BoneTree
}

func NewBoneTrees() *BoneTrees {
	return &BoneTrees{
		Data: make(map[BoneNameFrameNo]*BoneTree, 0),
	}
}

func (bts *BoneTrees) GetItem(boneName string, frameNo int) *BoneTree {
	return bts.Data[BoneNameFrameNo{boneName, frameNo}]
}

func (bts *BoneTrees) SetItem(boneName string, frameNo int, boneTree *BoneTree) {
	bts.Data[BoneNameFrameNo{boneName, frameNo}] = boneTree
}

func (bts *BoneTrees) GetBoneNames() []string {
	boneNames := make([]string, 0)
	for key := range bts.Data {
		boneNames = append(boneNames, key.BoneName)
	}
	return boneNames
}

func (bts *BoneTrees) GetFrameNos() []int {
	frameNos := make([]int, 0)
	for key := range bts.Data {
		frameNos = append(frameNos, key.FrameNo)
	}
	return frameNos
}

func (bts *BoneTrees) Contains(boneName string, frameNo int) bool {
	_, ok := bts.Data[BoneNameFrameNo{boneName, frameNo}]
	return ok
}
