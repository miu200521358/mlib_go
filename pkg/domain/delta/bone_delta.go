package delta

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type BoneDelta struct {
	Bone               *pmx.Bone          // ボーン
	Frame              int                // キーフレーム
	GlobalMatrix       *mmath.MMat4       // グローバル行列
	LocalMatrix        *mmath.MMat4       // ローカル行列
	GlobalPosition     *mmath.MVec3       // グローバル位置
	FramePosition      *mmath.MVec3       // キーフレ位置の変動量
	FrameMorphPosition *mmath.MVec3       // モーフ位置の変動量
	FrameRotation      *mmath.MQuaternion // キーフレ回転の変動量
	FrameMorphRotation *mmath.MQuaternion // モーフ回転の変動量
	FrameScale         *mmath.MVec3       // キーフレスケールの変動量
	FrameMorphScale    *mmath.MVec3       // モーフスケールの変動量
	UnitMatrix         *mmath.MMat4
	*MorphBoneDelta
}

func NewBoneDeltaByGlobalMatrix(
	bone *pmx.Bone, frame int, globalMatrix *mmath.MMat4, parentDelta *BoneDelta,
) *BoneDelta {
	var parentGlobalMatrix *mmath.MMat4
	if parentDelta != nil {
		parentGlobalMatrix = parentDelta.FilledGlobalMatrix()
	} else {
		parentGlobalMatrix = mmath.NewMMat4()
	}
	unitMatrix := parentGlobalMatrix.Muled(globalMatrix.Inverted())

	return &BoneDelta{
		Bone:           bone,
		Frame:          frame,
		MorphBoneDelta: NewMorphBoneDelta(),
		GlobalMatrix:   globalMatrix,
		LocalMatrix:    bone.OffsetMatrix.Muled(globalMatrix),
		UnitMatrix:     unitMatrix,
		// 物理演算後の移動を受け取ると逆オフセットかけても一部モデルで破綻するので一旦コメントアウト
		// framePosition:   unitMatrix.Translation(),
		FrameRotation: unitMatrix.Quaternion(),
	}
}

func (bd *BoneDelta) FilledGlobalMatrix() *mmath.MMat4 {
	if bd.GlobalMatrix == nil {
		bd.GlobalMatrix = mmath.NewMMat4()
	}
	return bd.GlobalMatrix
}

func (bd *BoneDelta) FilledLocalMatrix() *mmath.MMat4 {
	if bd.LocalMatrix == nil {
		// BOf行列: 自身のボーンのボーンオフセット行列をかけてローカル行列
		bd.LocalMatrix = bd.Bone.OffsetMatrix.Muled(bd.GlobalMatrix)
	}
	return bd.LocalMatrix
}

func (bd *BoneDelta) FilledGlobalPosition() *mmath.MVec3 {
	if bd.GlobalPosition == nil {
		if bd.GlobalMatrix != nil {
			bd.GlobalPosition = bd.GlobalMatrix.Translation()
		} else {
			bd.GlobalPosition = mmath.NewMVec3()
		}
	}
	return bd.GlobalPosition
}

func (bd *BoneDelta) FilledGlobalRotation() *mmath.MQuaternion {
	return bd.FilledGlobalMatrix().Quaternion()
}

func (bd *BoneDelta) FilledLocalPosition() *mmath.MVec3 {
	pos := bd.FilledFramePosition().Copy()

	if bd.FrameMorphPosition != nil && !bd.FrameMorphPosition.IsZero() {
		pos.Add(bd.FrameMorphPosition)
	}

	// if bd.FrameEffectPosition != nil && !bd.FrameEffectPosition.IsZero() {
	// 	pos.Add(bd.FrameEffectPosition)
	// }

	return pos
}

func (bd *BoneDelta) FilledLocalEffectorPosition(effectorFactor float64) *mmath.MVec3 {
	pos := bd.FilledFramePosition().Copy()

	if bd.FrameMorphPosition != nil && !bd.FrameMorphPosition.IsZero() {
		pos.Add(bd.FrameMorphPosition)
	}

	// if bd.FrameEffectPosition != nil && !bd.FrameEffectPosition.IsZero() {
	// 	pos.Add(bd.FrameEffectPosition)
	// }

	return pos.MuledScalar(effectorFactor)
}

func (bd *BoneDelta) FilledFramePosition() *mmath.MVec3 {
	if bd.FramePosition == nil {
		bd.FramePosition = mmath.NewMVec3()
	}
	return bd.FramePosition
}

func (bd *BoneDelta) FilledFrameMorphPosition() *mmath.MVec3 {
	if bd.FrameMorphPosition == nil {
		bd.FrameMorphPosition = mmath.NewMVec3()
	}
	return bd.FrameMorphPosition
}

func (bd *BoneDelta) FilledLocalRotation() *mmath.MQuaternion {
	rot := bd.FilledFrameRotation().Copy()

	if bd.FrameMorphRotation != nil && !bd.FrameMorphRotation.IsIdent() {
		// rot = bd.FrameMorphRotation.Muled(rot)
		rot.Mul(bd.FrameMorphRotation)
	}

	if bd.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bd.Bone.NormalizedFixedAxis)
	}

	return rot
}

func (bd *BoneDelta) FilledLocalEffectorRotation(effectorFactor float64) *mmath.MQuaternion {
	return bd.FilledLocalRotation().MuledScalar(effectorFactor)
}

func (bd *BoneDelta) FilledFrameRotation() *mmath.MQuaternion {
	if bd.FrameRotation == nil {
		bd.FrameRotation = mmath.NewMQuaternion()
	}
	return bd.FrameRotation
}

func (bd *BoneDelta) FilledFrameMorphRotation() *mmath.MQuaternion {
	if bd.FrameMorphRotation == nil {
		bd.FrameMorphRotation = mmath.NewMQuaternion()
	}
	return bd.FrameMorphRotation
}

func (bd *BoneDelta) FilledLocalScale() *mmath.MVec3 {
	pos := bd.FilledFrameScale().Copy()

	if bd.FrameMorphScale != nil && !bd.FrameMorphScale.IsZero() {
		pos.Add(bd.FrameMorphScale)
	}

	return pos
}

func (bd *BoneDelta) FilledFrameScale() *mmath.MVec3 {
	if bd.FrameScale == nil {
		bd.FrameScale = &mmath.MVec3{1, 1, 1}
	}
	return bd.FrameScale
}

func (bd *BoneDelta) FilledFrameMorphScale() *mmath.MVec3 {
	if bd.FrameMorphScale == nil {
		bd.FrameMorphScale = mmath.NewMVec3()
	}
	return bd.FrameMorphScale
}

func (bd *BoneDelta) Copy() *BoneDelta {
	return &BoneDelta{
		Bone:               bd.Bone,
		Frame:              bd.Frame,
		GlobalMatrix:       bd.FilledGlobalMatrix().Copy(),
		LocalMatrix:        bd.FilledLocalMatrix().Copy(),
		GlobalPosition:     bd.FilledGlobalPosition().Copy(),
		FramePosition:      bd.FilledFramePosition().Copy(),
		FrameMorphPosition: bd.FilledFrameMorphPosition().Copy(),
		FrameRotation:      bd.FilledFrameRotation().Copy(),
		FrameMorphRotation: bd.FilledFrameMorphRotation().Copy(),
		FrameScale:         bd.FilledFrameScale().Copy(),
		FrameMorphScale:    bd.FilledFrameMorphScale().Copy(),
		UnitMatrix:         bd.UnitMatrix.Copy(),
		MorphBoneDelta:     bd.MorphBoneDelta.Copy(),
	}
}

func NewBoneDelta(bone *pmx.Bone, frame int) *BoneDelta {
	return &BoneDelta{
		Bone:           bone,
		Frame:          frame,
		MorphBoneDelta: NewMorphBoneDelta(),
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
		distances[i] = worldPos.Distance(bd.FilledGlobalPosition())
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
			if !bd.FilledGlobalPosition().NearEquals(nearestBone.FilledGlobalPosition(), 0.01) {
				break
			}
		}
		boneIndexes = append(boneIndexes, sortedIndexes[i])
	}
	return boneIndexes
}

func (bds *BoneDeltas) LocalRotation(boneIndex int) *mmath.MQuaternion {
	return bds.localRotationLoop(boneIndex, 0)
}

func (bds *BoneDeltas) localRotationLoop(boneIndex int, loop int) *mmath.MQuaternion {
	bd := bds.Get(boneIndex)
	if bd == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}
	rot := bd.FilledLocalRotation()

	if bd.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := bds.localRotationLoop(bd.Bone.EffectIndex, loop+1)
		rot.Mul(effectorRot.MuledScalar(bd.Bone.EffectFactor))
	}

	return rot
}

func (bds *BoneDeltas) LocalPosition(boneIndex int) *mmath.MVec3 {
	return bds.localPositionLoop(boneIndex, 0)
}

func (bds *BoneDeltas) localPositionLoop(boneIndex int, loop int) *mmath.MVec3 {
	bd := bds.Get(boneIndex)
	if bd == nil || loop > 10 {
		return mmath.NewMVec3()
	}
	pos := bd.FilledLocalPosition()

	if bd.Bone.IsEffectorTranslation() {
		// 付与親移動がある場合、再帰で回転を取得する
		effectorPos := bds.localPositionLoop(bd.Bone.EffectIndex, loop+1)
		pos.Mul(effectorPos.MuledScalar(bd.Bone.EffectFactor))
	}

	return pos
}
