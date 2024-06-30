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
	frameMorphPosition  *mmath.MVec3       // モーフ位置の変動量
	frameEffectPosition *mmath.MVec3       // キーフレ位置の変動量(付与親のみ)
	frameRotation       *mmath.MQuaternion // キーフレ回転の変動量
	frameMorphRotation  *mmath.MQuaternion // モーフ回転の変動量
	frameEffectRotation *mmath.MQuaternion // キーフレ回転の変動量(付与親のみ)
	frameScale          *mmath.MVec3       // キーフレスケールの変動量
	frameMorphScale     *mmath.MVec3       // モーフスケールの変動量
	unitMatrix          *mmath.MMat4
	*MorphFrameDelta
}

func (bd *BoneDelta) GlobalMatrix() *mmath.MMat4 {
	if bd.globalMatrix == nil {
		bd.globalMatrix = mmath.NewMMat4()
	}
	return bd.globalMatrix
}

func (bd *BoneDelta) LocalMatrix() *mmath.MMat4 {
	if bd.localMatrix == nil {
		// BOf行列: 自身のボーンのボーンオフセット行列をかけてローカル行列
		bd.localMatrix = bd.Bone.OffsetMatrix.Muled(bd.globalMatrix)
	}
	return bd.localMatrix
}

func (bd *BoneDelta) GlobalPosition() *mmath.MVec3 {
	if bd.globalPosition == nil {
		bd.globalPosition = bd.globalMatrix.Translation()
	}
	return bd.globalPosition
}

func (bd *BoneDelta) GlobalRotation() *mmath.MQuaternion {
	return bd.GlobalMatrix().Quaternion()
}

func (bd *BoneDelta) LocalPosition() *mmath.MVec3 {
	pos := bd.FramePosition().Copy()

	if bd.frameMorphPosition != nil && !bd.frameMorphPosition.IsZero() {
		pos.Add(bd.frameMorphPosition)
	}

	if bd.frameEffectPosition != nil && !bd.frameEffectPosition.IsZero() {
		pos.Add(bd.frameEffectPosition)
	}

	return pos
}

func (bd *BoneDelta) FramePosition() *mmath.MVec3 {
	if bd.framePosition == nil {
		bd.framePosition = mmath.NewMVec3()
	}
	return bd.framePosition
}

func (bd *BoneDelta) FrameMorphPosition() *mmath.MVec3 {
	if bd.frameMorphPosition == nil {
		bd.frameMorphPosition = mmath.NewMVec3()
	}
	return bd.frameMorphPosition
}

func (bd *BoneDelta) FrameEffectPosition() *mmath.MVec3 {
	if bd.frameEffectPosition == nil {
		bd.frameEffectPosition = mmath.NewMVec3()
	}
	return bd.frameEffectPosition
}

func (bd *BoneDelta) LocalRotation() *mmath.MQuaternion {
	rot := bd.FrameRotation().Copy()

	if bd.frameMorphRotation != nil && !bd.frameMorphRotation.IsIdent() {
		// rot = bd.frameMorphRotation.Muled(rot)
		rot.Mul(bd.frameMorphRotation)
	}

	if bd.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bd.Bone.NormalizedFixedAxis)
	}

	if bd.frameEffectRotation != nil && !bd.frameEffectRotation.IsIdent() {
		// rot = bd.frameEffectRotation.Muled(rot)
		rot.Mul(bd.frameEffectRotation)
	}

	return rot
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

func (bd *BoneDelta) FrameMorphRotation() *mmath.MQuaternion {
	if bd.frameMorphRotation == nil {
		bd.frameMorphRotation = mmath.NewMQuaternion()
	}
	return bd.frameMorphRotation
}

func (bd *BoneDelta) LocalScale() *mmath.MVec3 {
	pos := bd.FrameScale().Copy()

	if bd.frameMorphScale != nil && !bd.frameMorphScale.IsZero() {
		pos.Add(bd.frameMorphScale)
	}

	return pos
}

func (bd *BoneDelta) FrameScale() *mmath.MVec3 {
	if bd.frameScale == nil {
		bd.frameScale = &mmath.MVec3{1, 1, 1}
	}
	return bd.frameScale
}

func (bd *BoneDelta) FrameMorphScale() *mmath.MVec3 {
	if bd.frameMorphScale == nil {
		bd.frameMorphScale = mmath.NewMVec3()
	}
	return bd.frameMorphScale
}

func (bd *BoneDelta) Copy() *BoneDelta {
	return &BoneDelta{
		Bone:                bd.Bone,
		Frame:               bd.Frame,
		globalMatrix:        bd.GlobalMatrix().Copy(),
		localMatrix:         bd.LocalMatrix().Copy(),
		globalPosition:      bd.GlobalPosition().Copy(),
		framePosition:       bd.FramePosition().Copy(),
		frameMorphPosition:  bd.FrameMorphPosition().Copy(),
		frameEffectPosition: bd.FrameEffectPosition().Copy(),
		frameRotation:       bd.FrameRotation().Copy(),
		frameMorphRotation:  bd.FrameMorphRotation().Copy(),
		frameEffectRotation: bd.FrameEffectRotation().Copy(),
		frameScale:          bd.FrameScale().Copy(),
		frameMorphScale:     bd.FrameMorphScale().Copy(),
		unitMatrix:          bd.unitMatrix.Copy(),
		MorphFrameDelta:     bd.MorphFrameDelta.Copy(),
	}
}

func NewBoneDelta(
	bone *pmx.Bone,
	frame int,
) *BoneDelta {
	return &BoneDelta{
		Bone:            bone,
		Frame:           frame,
		MorphFrameDelta: NewMorphFrameDelta(),
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
	if bd == nil {
		bd = NewBoneDelta(bone, 0)
	}
	bd.globalMatrix = globalMatrix
	bds.Append(bd)
}

// FillLocalMatrix 物理演算後にグローバル行列を埋め終わった後に呼び出して、ローカル行列を計算する
func (bds *BoneDeltas) FillLocalMatrix() {
	for _, bd := range bds.Data {
		var parentGlobalMatrix *mmath.MMat4
		if bd.Bone.ParentIndex >= 0 {
			parentGlobalMatrix = bds.Get(bd.Bone.ParentIndex).GlobalMatrix()
		} else {
			parentGlobalMatrix = mmath.NewMMat4()
		}
		unitMatrix := parentGlobalMatrix.Muled(bd.globalMatrix.Inverted())

		bd.localMatrix = bd.Bone.OffsetMatrix.Muled(bd.globalMatrix)
		bd.globalPosition = nil
		bd.framePosition = unitMatrix.Translation()
		bd.frameMorphPosition = nil
		bd.frameEffectPosition = nil
		bd.frameRotation = unitMatrix.Quaternion()
		bd.frameMorphRotation = nil
		bd.frameEffectRotation = nil
	}
}
