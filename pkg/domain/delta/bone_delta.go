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
	*MorphFrameDelta
}

func NewBoneDeltaByGlobalMatrix(
	bone *pmx.Bone, frame int, globalMatrix *mmath.MMat4, parentDelta *BoneDelta,
) *BoneDelta {
	var parentGlobalMatrix *mmath.MMat4
	if parentDelta != nil {
		parentGlobalMatrix = parentDelta.GetGlobalMatrix()
	} else {
		parentGlobalMatrix = mmath.NewMMat4()
	}
	unitMatrix := parentGlobalMatrix.Muled(globalMatrix.Inverted())

	return &BoneDelta{
		Bone:            bone,
		Frame:           frame,
		MorphFrameDelta: NewMorphFrameDelta(),
		GlobalMatrix:    globalMatrix,
		LocalMatrix:     bone.OffsetMatrix.Muled(globalMatrix),
		UnitMatrix:      unitMatrix,
		// 物理演算後の移動を受け取ると逆オフセットかけても一部モデルで破綻するので一旦コメントアウト
		// framePosition:   unitMatrix.Translation(),
		FrameRotation: unitMatrix.Quaternion(),
	}
}

func (bd *BoneDelta) GetGlobalMatrix() *mmath.MMat4 {
	if bd.GlobalMatrix == nil {
		bd.GlobalMatrix = mmath.NewMMat4()
	}
	return bd.GlobalMatrix
}

func (bd *BoneDelta) GetLocalMatrix() *mmath.MMat4 {
	if bd.LocalMatrix == nil {
		// BOf行列: 自身のボーンのボーンオフセット行列をかけてローカル行列
		bd.LocalMatrix = bd.Bone.OffsetMatrix.Muled(bd.GlobalMatrix)
	}
	return bd.LocalMatrix
}

func (bd *BoneDelta) GetGlobalPosition() *mmath.MVec3 {
	if bd.GlobalPosition == nil {
		if bd.GlobalMatrix != nil {
			bd.GlobalPosition = bd.GlobalMatrix.Translation()
		} else {
			bd.GlobalPosition = mmath.NewMVec3()
		}
	}
	return bd.GlobalPosition
}

func (bd *BoneDelta) GetGlobalRotation() *mmath.MQuaternion {
	return bd.GetGlobalMatrix().Quaternion()
}

func (bd *BoneDelta) GetLocalPosition() *mmath.MVec3 {
	pos := bd.GetFramePosition().Copy()

	if bd.FrameMorphPosition != nil && !bd.FrameMorphPosition.IsZero() {
		pos.Add(bd.FrameMorphPosition)
	}

	// if bd.FrameEffectPosition != nil && !bd.FrameEffectPosition.IsZero() {
	// 	pos.Add(bd.FrameEffectPosition)
	// }

	return pos
}

func (bd *BoneDelta) GetLocalEffectorPosition(effectorFactor float64) *mmath.MVec3 {
	pos := bd.GetFramePosition().Copy()

	if bd.FrameMorphPosition != nil && !bd.FrameMorphPosition.IsZero() {
		pos.Add(bd.FrameMorphPosition)
	}

	// if bd.FrameEffectPosition != nil && !bd.FrameEffectPosition.IsZero() {
	// 	pos.Add(bd.FrameEffectPosition)
	// }

	return pos.MuledScalar(effectorFactor)
}

func (bd *BoneDelta) GetFramePosition() *mmath.MVec3 {
	if bd.FramePosition == nil {
		bd.FramePosition = mmath.NewMVec3()
	}
	return bd.FramePosition
}

func (bd *BoneDelta) GetFrameMorphPosition() *mmath.MVec3 {
	if bd.FrameMorphPosition == nil {
		bd.FrameMorphPosition = mmath.NewMVec3()
	}
	return bd.FrameMorphPosition
}

func (bd *BoneDelta) GetLocalRotation() *mmath.MQuaternion {
	rot := bd.GetFrameRotation().Copy()

	if bd.FrameMorphRotation != nil && !bd.FrameMorphRotation.IsIdent() {
		// rot = bd.FrameMorphRotation.Muled(rot)
		rot.Mul(bd.FrameMorphRotation)
	}

	if bd.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bd.Bone.NormalizedFixedAxis)
	}

	return rot
}

func (bd *BoneDelta) GetLocalEffectorRotation(effectorFactor float64) *mmath.MQuaternion {
	return bd.GetLocalRotation().MuledScalar(effectorFactor)
}

func (bd *BoneDelta) GetFrameRotation() *mmath.MQuaternion {
	if bd.FrameRotation == nil {
		bd.FrameRotation = mmath.NewMQuaternion()
	}
	return bd.FrameRotation
}

func (bd *BoneDelta) GetFrameMorphRotation() *mmath.MQuaternion {
	if bd.FrameMorphRotation == nil {
		bd.FrameMorphRotation = mmath.NewMQuaternion()
	}
	return bd.FrameMorphRotation
}

func (bd *BoneDelta) GetLocalScale() *mmath.MVec3 {
	pos := bd.GetFrameScale().Copy()

	if bd.FrameMorphScale != nil && !bd.FrameMorphScale.IsZero() {
		pos.Add(bd.FrameMorphScale)
	}

	return pos
}

func (bd *BoneDelta) GetFrameScale() *mmath.MVec3 {
	if bd.FrameScale == nil {
		bd.FrameScale = &mmath.MVec3{1, 1, 1}
	}
	return bd.FrameScale
}

func (bd *BoneDelta) GetFrameMorphScale() *mmath.MVec3 {
	if bd.FrameMorphScale == nil {
		bd.FrameMorphScale = mmath.NewMVec3()
	}
	return bd.FrameMorphScale
}

func (bd *BoneDelta) Copy() *BoneDelta {
	return &BoneDelta{
		Bone:               bd.Bone,
		Frame:              bd.Frame,
		GlobalMatrix:       bd.GetGlobalMatrix().Copy(),
		LocalMatrix:        bd.GetLocalMatrix().Copy(),
		GlobalPosition:     bd.GetGlobalPosition().Copy(),
		FramePosition:      bd.GetFramePosition().Copy(),
		FrameMorphPosition: bd.GetFrameMorphPosition().Copy(),
		FrameRotation:      bd.GetFrameRotation().Copy(),
		FrameMorphRotation: bd.GetFrameMorphRotation().Copy(),
		FrameScale:         bd.GetFrameScale().Copy(),
		FrameMorphScale:    bd.GetFrameMorphScale().Copy(),
		UnitMatrix:         bd.UnitMatrix.Copy(),
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
		distances[i] = worldPos.Distance(bd.GetGlobalPosition())
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
			if !bd.GetGlobalPosition().NearEquals(nearestBone.GetGlobalPosition(), 0.01) {
				break
			}
		}
		boneIndexes = append(boneIndexes, sortedIndexes[i])
	}
	return boneIndexes
}

func (bds *BoneDeltas) GetLocalRotation(boneIndex int, loop int) *mmath.MQuaternion {
	bd := bds.Get(boneIndex)
	if bd == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}
	rot := bd.GetLocalRotation()

	if bd.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := bds.GetLocalRotation(bd.Bone.EffectIndex, loop+1)
		rot.Mul(effectorRot.MuledScalar(bd.Bone.EffectFactor))
	}

	return rot
}

func (bds *BoneDeltas) LocalPosition(boneIndex int, loop int) *mmath.MVec3 {
	bd := bds.Get(boneIndex)
	if bd == nil || loop > 10 {
		return mmath.NewMVec3()
	}
	pos := bd.GetLocalPosition()

	if bd.Bone.IsEffectorTranslation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorPos := bds.LocalPosition(bd.Bone.EffectIndex, loop+1)
		pos.Mul(effectorPos.MuledScalar(bd.Bone.EffectFactor))
	}

	return pos
}
