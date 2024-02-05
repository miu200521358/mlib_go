package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type BoneTree struct {
	BoneName      string
	Frame         float32
	GlobalMatrix  *mmath.MMat4
	LocalMatrix   *mmath.MMat4
	Position      *mmath.MVec3
	FramePosition *mmath.MVec3
	FrameRotation *mmath.MQuaternion
	FrameScale    *mmath.MVec3
}

func NewBoneTree(
	boneName string,
	frame float32,
	globalMatrix, localMatrix *mmath.MMat4,
	framePosition *mmath.MVec3,
	frameRotation *mmath.MQuaternion,
	frameScale *mmath.MVec3,
) *BoneTree {
	p := globalMatrix.Translation()
	return &BoneTree{
		BoneName:      boneName,
		Frame:         frame,
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
	Frame    float32
}

type BoneTrees struct {
	Data map[BoneNameFrameNo]*BoneTree
}

func NewBoneTrees() *BoneTrees {
	return &BoneTrees{
		Data: make(map[BoneNameFrameNo]*BoneTree, 0),
	}
}

func (bts *BoneTrees) GetItem(boneName string, frame float32) *BoneTree {
	return bts.Data[BoneNameFrameNo{boneName, frame}]
}

func (bts *BoneTrees) SetItem(boneName string, frame float32, boneTree *BoneTree) {
	bts.Data[BoneNameFrameNo{boneName, frame}] = boneTree
}

func (bts *BoneTrees) GetBoneNames() []string {
	boneNames := make([]string, 0)
	for key := range bts.Data {
		boneNames = append(boneNames, key.BoneName)
	}
	return boneNames
}

func (bts *BoneTrees) GetFrameNos() []float32 {
	frames := make([]float32, 0)
	for key := range bts.Data {
		frames = append(frames, key.Frame)
	}
	return frames
}

func (bts *BoneTrees) Contains(boneName string, frame float32) bool {
	_, ok := bts.Data[BoneNameFrameNo{boneName, frame}]
	return ok
}
