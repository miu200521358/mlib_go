package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type BoneDelta struct {
	Bone                *pmx.Bone          // ボーン
	Frame               int                // キーフレーム
	globalMatrix        *mmath.MMat4       // グローバル行列
	localMatrix         *mmath.MMat4       // ローカル行列
	unitMatrix          *mmath.MMat4       // ボーン単体のデフォーム行列
	globalPosition      *mmath.MVec3       // グローバル位置
	framePosition       *mmath.MVec3       // キーフレ位置の変動量
	frameEffectPosition *mmath.MVec3       // キーフレ位置の変動量(付与親のみ)
	frameRotation       *mmath.MQuaternion // キーフレ回転の変動量
	frameEffectRotation *mmath.MQuaternion // キーフレ回転の変動量(付与親のみ)
	frameScale          *mmath.MVec3       // キーフレスケールの変動量
}

func (bd *BoneDelta) GlobalMatrix() *mmath.MMat4 {
	if bd.globalMatrix == nil {
		bd.globalMatrix = mmath.NewMMat4()
	}
	return bd.globalMatrix
}

func (bd *BoneDelta) LocalMatrix() *mmath.MMat4 {
	if bd.localMatrix == nil {
		bd.localMatrix = mmath.NewMMat4()
	}
	return bd.localMatrix
}

func (bd *BoneDelta) UnitMatrix() *mmath.MMat4 {
	if bd.unitMatrix == nil {
		bd.unitMatrix = mmath.NewMMat4()
	}
	return bd.unitMatrix
}

func (bd *BoneDelta) GlobalPosition() *mmath.MVec3 {
	if bd.globalPosition == nil {
		bd.globalPosition = mmath.NewMVec3()
	}
	return bd.globalPosition
}

func (bd *BoneDelta) GlobalRotation() *mmath.MQuaternion {
	return bd.GlobalMatrix().Quaternion()
}

func (bd *BoneDelta) FramePosition() *mmath.MVec3 {
	if bd.framePosition == nil {
		bd.framePosition = mmath.NewMVec3()
	}
	return bd.framePosition
}

func (bd *BoneDelta) FrameEffectPosition() *mmath.MVec3 {
	if bd.frameEffectPosition == nil {
		bd.frameEffectPosition = mmath.NewMVec3()
	}
	return bd.frameEffectPosition
}

func (bd *BoneDelta) FrameRotation() *mmath.MQuaternion {
	if bd.frameRotation == nil {
		bd.frameRotation = mmath.NewMQuaternion()
	}
	return bd.frameRotation
}

func (bd *BoneDelta) FrameEffectRotation() *mmath.MQuaternion {
	if bd.frameEffectRotation == nil {
		bd.frameEffectRotation = mmath.NewMQuaternion()
	}
	return bd.frameEffectRotation
}

func (bd *BoneDelta) FrameScale() *mmath.MVec3 {
	if bd.frameScale == nil {
		bd.frameScale = &mmath.MVec3{1, 1, 1}
	}
	return bd.frameScale
}

func NewBoneDelta(
	bone *pmx.Bone,
	frame int,
	globalMatrix, localMatrix, unitMatrix *mmath.MMat4,
	framePosition, frameEffectPosition *mmath.MVec3,
	frameRotation, frameEffectRotation *mmath.MQuaternion,
	frameScale *mmath.MVec3,
) *BoneDelta {
	return &BoneDelta{
		Bone:                bone,
		Frame:               frame,
		globalMatrix:        globalMatrix,
		localMatrix:         localMatrix,
		unitMatrix:          unitMatrix,
		globalPosition:      globalMatrix.Translation(),
		framePosition:       framePosition,
		frameEffectPosition: frameEffectPosition,
		frameRotation:       frameRotation,
		frameEffectRotation: frameEffectRotation,
		frameScale:          frameScale,
	}
}

type BoneDeltas struct {
	Data map[int]*BoneDelta
}

func NewBoneDeltas() *BoneDeltas {
	return &BoneDeltas{
		Data: make(map[int]*BoneDelta, 0),
	}
}

func (bts *BoneDeltas) Get(boneIndex int) *BoneDelta {
	return bts.Data[boneIndex]
}

func (bts *BoneDeltas) SetItem(boneIndex int, boneDelta *BoneDelta) {
	bts.Data[boneIndex] = boneDelta
}

func (bts *BoneDeltas) GetBoneIndexes() []int {
	boneIndexes := make([]int, 0)
	for key := range bts.Data {
		if !slices.Contains(boneIndexes, key) {
			boneIndexes = append(boneIndexes, key)
		}
	}
	return boneIndexes
}

func (bts *BoneDeltas) Contains(boneIndex int) bool {
	_, ok := bts.Data[boneIndex]
	return ok
}
