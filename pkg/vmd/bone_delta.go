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
	globalPosition      *mmath.MVec3       // グローバル位置
	framePosition       *mmath.MVec3       // キーフレ位置の変動量
	frameEffectPosition *mmath.MVec3       // キーフレ位置の変動量(付与親のみ)
	frameRotation       *mmath.MQuaternion // キーフレ回転の変動量
	frameEffectRotation *mmath.MQuaternion // キーフレ回転の変動量(付与親のみ)
	frameIkRotation     *mmath.MQuaternion // キーフレIK回転の変動量
	frameScale          *mmath.MVec3       // キーフレスケールの変動量
	unitMatrix          *mmath.MMat4
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

func (bd *BoneDelta) GlobalPosition() *mmath.MVec3 {
	if bd.globalPosition == nil {
		bd.globalPosition = mmath.NewMVec3()
	}
	return bd.globalPosition
}

func (bd *BoneDelta) LocalRotation() *mmath.MQuaternion {
	rot := bd.FrameRotation().Copy()
	if bd.frameIkRotation != nil && !bd.frameIkRotation.IsIdent() {
		// rot.Mul(bd.frameIkRotation)
		rot = bd.frameIkRotation.Muled(rot)
	}
	if bd.frameEffectRotation != nil && !bd.frameEffectRotation.IsIdent() {
		rot = bd.frameEffectRotation.Muled(rot)
		// rot.Mul(bd.frameEffectRotation)
	}

	if bd.Bone.HasFixedAxis() {
		return rot.ToFixedAxisRotation(bd.Bone.FixedAxis)
	}

	return rot
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

func (bd *BoneDelta) FrameIkRotation() *mmath.MQuaternion {
	if bd.frameIkRotation == nil {
		bd.frameIkRotation = mmath.NewMQuaternion()
	}
	return bd.frameIkRotation
}

func (bd *BoneDelta) FrameScale() *mmath.MVec3 {
	if bd.frameScale == nil {
		bd.frameScale = &mmath.MVec3{1, 1, 1}
	}
	return bd.frameScale
}

func (bd *BoneDelta) Copy() *BoneDelta {
	return &BoneDelta{
		Bone:                bd.Bone,
		Frame:               bd.Frame,
		globalMatrix:        bd.GlobalMatrix().Copy(),
		localMatrix:         bd.LocalMatrix().Copy(),
		globalPosition:      bd.GlobalPosition().Copy(),
		framePosition:       bd.FramePosition().Copy(),
		frameEffectPosition: bd.FrameEffectPosition().Copy(),
		frameRotation:       bd.FrameRotation().Copy(),
		frameEffectRotation: bd.FrameEffectRotation().Copy(),
		frameScale:          bd.FrameScale().Copy(),
		unitMatrix:          bd.unitMatrix.Copy(),
	}
}

func NewBoneDelta(
	bone *pmx.Bone,
	frame int,
	// globalMatrix, unitMatrix *mmath.MMat4,
	// framePosition, frameEffectPosition *mmath.MVec3,
	// frameRotation, frameEffectRotation *mmath.MQuaternion,
	// frameScale *mmath.MVec3,
) *BoneDelta {
	return &BoneDelta{
		Bone:  bone,
		Frame: frame,
		// globalMatrix: globalMatrix,
		//
		// localMatrix:         bone.OffsetMatrix.Muled(globalMatrix),
		// globalPosition:
		// framePosition:       framePosition,
		// frameEffectPosition: frameEffectPosition,
		// frameRotation:       frameRotation,
		// frameEffectRotation: frameEffectRotation,
		// frameScale:          frameScale,
		// unitMatrix:          unitMatrix,
	}
}

type BoneDeltas struct {
	Data  map[int]*BoneDelta
	Names map[string]int
}

func NewBoneDeltas() *BoneDeltas {
	return &BoneDeltas{
		Data:  make(map[int]*BoneDelta, 0),
		Names: make(map[string]int, 0),
	}
}

func (bts *BoneDeltas) Get(boneIndex int) *BoneDelta {
	if _, ok := bts.Data[boneIndex]; !ok {
		return nil
	}
	return bts.Data[boneIndex]
}

func (bts *BoneDeltas) GetByName(boneName string) *BoneDelta {
	if _, ok := bts.Names[boneName]; !ok {
		return nil
	}
	return bts.Get(bts.Names[boneName])
}

func (bts *BoneDeltas) Append(boneDelta *BoneDelta) {
	bts.Data[boneDelta.Bone.Index] = boneDelta
	bts.Names[boneDelta.Bone.Name] = boneDelta.Bone.Index
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

func (bds *BoneDeltas) SetGlobalMatrix(bone *pmx.Bone, globalMatrix *mmath.MMat4) {
	bd := bds.Get(bone.Index)
	bd.globalMatrix = globalMatrix

	var parentGlobalMatrix *mmath.MMat4
	if bd.Bone.ParentIndex >= 0 {
		parentGlobalMatrix = bds.Get(bd.Bone.ParentIndex).GlobalMatrix()
	} else {
		parentGlobalMatrix = mmath.NewMMat4()
	}
	unitMatrix := parentGlobalMatrix.Muled(globalMatrix.Inverted())

	bd.localMatrix = bone.OffsetMatrix.Muled(globalMatrix)
	bd.globalPosition = nil
	bd.frameRotation = unitMatrix.Quaternion()
	bd.frameEffectRotation = nil
	bd.framePosition = unitMatrix.Translation()
	bd.frameEffectPosition = nil
}
