package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type BoneDelta struct {
	Bone               *pmx.Bone          // ボーン
	Frame              int                // キーフレーム
	globalMatrix       *mmath.MMat4       // グローバル行列
	localMatrix        *mmath.MMat4       // ローカル行列
	globalPosition     *mmath.MVec3       // グローバル位置
	framePosition      *mmath.MVec3       // キーフレ位置の変動量
	frameMorphPosition *mmath.MVec3       // モーフ位置の変動量
	frameRotation      *mmath.MQuaternion // キーフレ回転の変動量
	frameMorphRotation *mmath.MQuaternion // モーフ回転の変動量
	frameScale         *mmath.MVec3       // キーフレスケールの変動量
	frameMorphScale    *mmath.MVec3       // モーフスケールの変動量
	unitMatrix         *mmath.MMat4
	*MorphFrameDelta
}

func NewBoneDeltaByGlobalMatrix(
	bone *pmx.Bone, frame int, globalMatrix *mmath.MMat4, parentDelta *BoneDelta,
) *BoneDelta {
	var parentGlobalMatrix *mmath.MMat4
	if parentDelta != nil {
		parentGlobalMatrix = parentDelta.GlobalMatrix()
	} else {
		parentGlobalMatrix = mmath.NewMMat4()
	}
	unitMatrix := parentGlobalMatrix.Muled(globalMatrix.Inverted())

	return &BoneDelta{
		Bone:            bone,
		Frame:           frame,
		MorphFrameDelta: NewMorphFrameDelta(),
		globalMatrix:    globalMatrix,
		localMatrix:     bone.OffsetMatrix.Muled(globalMatrix),
		unitMatrix:      unitMatrix,
		// 物理演算後の移動を受け取ると逆オフセットかけても一部モデルで破綻するので一旦コメントアウト
		// framePosition:   unitMatrix.Translation(),
		frameRotation: unitMatrix.Quaternion(),
	}
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

	// if bd.frameEffectPosition != nil && !bd.frameEffectPosition.IsZero() {
	// 	pos.Add(bd.frameEffectPosition)
	// }

	return pos
}

func (bd *BoneDelta) LocalEffectorPosition(effectorFactor float64) *mmath.MVec3 {
	pos := bd.FramePosition().Copy()

	if bd.frameMorphPosition != nil && !bd.frameMorphPosition.IsZero() {
		pos.Add(bd.frameMorphPosition)
	}

	// if bd.frameEffectPosition != nil && !bd.frameEffectPosition.IsZero() {
	// 	pos.Add(bd.frameEffectPosition)
	// }

	return pos.MuledScalar(effectorFactor)
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

func (bd *BoneDelta) LocalRotation() *mmath.MQuaternion {
	rot := bd.FrameRotation().Copy()

	if bd.frameMorphRotation != nil && !bd.frameMorphRotation.IsIdent() {
		// rot = bd.frameMorphRotation.Muled(rot)
		rot.Mul(bd.frameMorphRotation)
	}

	if bd.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bd.Bone.NormalizedFixedAxis)
	}

	return rot
}

func (bd *BoneDelta) LocalEffectorRotation(effectorFactor float64) *mmath.MQuaternion {
	return bd.LocalRotation().MuledScalar(effectorFactor)
}

func (bd *BoneDelta) FrameRotation() *mmath.MQuaternion {
	if bd.frameRotation == nil {
		bd.frameRotation = mmath.NewMQuaternion()
	}
	return bd.frameRotation
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
		Bone:               bd.Bone,
		Frame:              bd.Frame,
		globalMatrix:       bd.GlobalMatrix().Copy(),
		localMatrix:        bd.LocalMatrix().Copy(),
		globalPosition:     bd.GlobalPosition().Copy(),
		framePosition:      bd.FramePosition().Copy(),
		frameMorphPosition: bd.FrameMorphPosition().Copy(),
		frameRotation:      bd.FrameRotation().Copy(),
		frameMorphRotation: bd.FrameMorphRotation().Copy(),
		frameScale:         bd.FrameScale().Copy(),
		frameMorphScale:    bd.FrameMorphScale().Copy(),
		unitMatrix:         bd.unitMatrix.Copy(),
		MorphFrameDelta:    bd.MorphFrameDelta.Copy(),
	}
}

func NewBoneDelta(bone *pmx.Bone, frame int) *BoneDelta {
	return &BoneDelta{
		Bone:            bone,
		Frame:           frame,
		MorphFrameDelta: NewMorphFrameDelta(),
	}
}

type BoneDeltas struct {
	Data  []*BoneDelta
	Bones *pmx.Bones
}

func NewBoneDeltas(bones *pmx.Bones) *BoneDeltas {
	return &BoneDeltas{
		Data:  make([]*BoneDelta, bones.Len()),
		Bones: bones,
	}
}

func (bds *BoneDeltas) Get(boneIndex int) *BoneDelta {
	if boneIndex < 0 || boneIndex >= len(bds.Data) {
		return nil
	}

	return bds.Data[boneIndex]
}

func (bds *BoneDeltas) GetByName(boneName string) *BoneDelta {
	if _, ok := bds.Bones.NameIndexes[boneName]; ok {
		return bds.Data[bds.Bones.NameIndexes[boneName]]
	}
	return nil
}

func (bds *BoneDeltas) Update(boneDelta *BoneDelta) {
	bds.Data[boneDelta.Bone.Index] = boneDelta
}

func (bds *BoneDeltas) GetBoneIndexes() []int {
	boneIndexes := make([]int, 0)
	for key := range bds.Data {
		if !slices.Contains(boneIndexes, key) {
			boneIndexes = append(boneIndexes, key)
		}
	}
	return boneIndexes
}

func (bds *BoneDeltas) Contains(boneIndex int) bool {
	return bds.Data[boneIndex] != nil
}

func (bds *BoneDeltas) GetNearestBoneIndexes(worldPos *mmath.MVec3) []int {
	boneIndexes := make([]int, 0)
	distances := make([]float64, len(bds.Data))
	for i := range len(bds.Data) {
		bd := bds.Get(i)
		if bd == nil {
			continue
		}
		distances[i] = worldPos.Distance(bd.GlobalPosition())
	}
	if len(distances) == 0 {
		return boneIndexes
	}
	sortedDistances := mmath.Float64Slice(distances)
	sortedIndexes := mmath.ArgSort(sortedDistances)
	nearestBone := bds.Get(sortedIndexes[0])
	for i := range sortedIndexes {
		bd := bds.Get(sortedIndexes[i])
		if bd == nil {
			continue
		}
		if len(boneIndexes) > 0 {
			if !bd.GlobalPosition().NearEquals(nearestBone.GlobalPosition(), 0.01) {
				break
			}
		}
		boneIndexes = append(boneIndexes, sortedIndexes[i])
	}
	return boneIndexes
}

func (bds *BoneDeltas) LocalRotation(boneIndex int, loop int) *mmath.MQuaternion {
	bd := bds.Get(boneIndex)
	if bd == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}
	rot := bd.LocalRotation()

	if bd.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := bds.LocalRotation(bd.Bone.EffectIndex, loop+1)
		rot.Mul(effectorRot.MuledScalar(bd.Bone.EffectFactor))
	}

	return rot
}

func (bds *BoneDeltas) LocalPosition(boneIndex int, loop int) *mmath.MVec3 {
	bd := bds.Get(boneIndex)
	if bd == nil || loop > 10 {
		return mmath.NewMVec3()
	}
	pos := bd.LocalPosition()

	if bd.Bone.IsEffectorTranslation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorPos := bds.LocalPosition(bd.Bone.EffectIndex, loop+1)
		pos.Mul(effectorPos.MuledScalar(bd.Bone.EffectFactor))
	}

	return pos
}
