package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type BoneDelta struct {
	BoneName      string
	Frame         float32
	GlobalMatrix  *mmath.MMat4
	LocalMatrix   *mmath.MMat4
	Position      *mmath.MVec3
	FramePosition *mmath.MVec3
	FrameRotation *mmath.MQuaternion
	FrameScale    *mmath.MVec3
}

func NewBoneDelta(
	boneName string,
	frame float32,
	globalMatrix, localMatrix *mmath.MMat4,
	framePosition *mmath.MVec3,
	frameRotation *mmath.MQuaternion,
	frameScale *mmath.MVec3,
) *BoneDelta {
	p := globalMatrix.Translation()
	return &BoneDelta{
		BoneName:      boneName,
		Frame:         frame,
		GlobalMatrix:  globalMatrix,
		LocalMatrix:   localMatrix,
		Position:      p,
		FramePosition: framePosition,
		FrameRotation: frameRotation,
		FrameScale:    frameScale,
	}
}

type BoneNameFrameNo struct {
	BoneName string
	Frame    float32
}

type BoneDeltas struct {
	Data map[BoneNameFrameNo]*BoneDelta
}

func NewBoneDeltas() *BoneDeltas {
	return &BoneDeltas{
		Data: make(map[BoneNameFrameNo]*BoneDelta, 0),
	}
}

func (bts *BoneDeltas) GetItem(boneName string, frame float32) *BoneDelta {
	return bts.Data[BoneNameFrameNo{boneName, frame}]
}

func (bts *BoneDeltas) SetItem(boneName string, frame float32, boneDelta *BoneDelta) {
	bts.Data[BoneNameFrameNo{boneName, frame}] = boneDelta
}

func (bts *BoneDeltas) GetBoneNames() []string {
	boneNames := make([]string, 0)
	for key := range bts.Data {
		if slices.Contains(boneNames, key.BoneName) {
			boneNames = append(boneNames, key.BoneName)
		}
	}
	return boneNames
}

func (bts *BoneDeltas) GetFrameNos() []float32 {
	frames := make([]float32, 0)
	for key := range bts.Data {
		if slices.Contains(frames, key.Frame) {
			frames = append(frames, key.Frame)
		}
	}
	return frames
}

func (bts *BoneDeltas) Contains(boneName string, frame float32) bool {
	_, ok := bts.Data[BoneNameFrameNo{boneName, frame}]
	return ok
}
