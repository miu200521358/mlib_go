package delta

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type BoneDelta struct {
	Bone                    *pmx.Bone          // ボーン
	Frame                   float32            // キーフレーム
	GlobalMatrix            *mmath.MMat4       // グローバル行列
	LocalMatrix             *mmath.MMat4       // ローカル行列
	GlobalPosition          *mmath.MVec3       // グローバル位置
	UnitMatrix              *mmath.MMat4       // 親ボーンからの変位行列
	FramePosition           *mmath.MVec3       // キーフレ位置の変動量
	FrameMorphPosition      *mmath.MVec3       // モーフ位置の変動量
	FrameLocalMorphPosition *mmath.MVec3       // モーフ位置のローカル変動量
	FrameRotation           *mmath.MQuaternion // キーフレ回転の変動量
	FrameMorphRotation      *mmath.MQuaternion // モーフ回転の変動量
	FrameLocalMorphRotation *mmath.MQuaternion // モーフ回転のローカル変動量
	FrameScale              *mmath.MVec3       // キーフレスケールの変動量
	FrameMorphScale         *mmath.MVec3       // モーフスケールの変動量
	FrameLocalMorphScale    *mmath.MVec3       // モーフスケールのローカル変動量
}

func NewBoneDeltaByGlobalMatrix(
	bone *pmx.Bone, frame float32, globalMatrix *mmath.MMat4, parentDelta *BoneDelta,
) *BoneDelta {
	var parentGlobalMatrix *mmath.MMat4
	if parentDelta != nil {
		parentGlobalMatrix = parentDelta.FilledGlobalMatrix()
	} else {
		parentGlobalMatrix = mmath.NewMMat4()
	}
	unitMatrix := parentGlobalMatrix.Muled(globalMatrix.Inverted())

	return &BoneDelta{
		Bone:         bone,
		Frame:        frame,
		GlobalMatrix: globalMatrix,
		LocalMatrix:  bone.Extend.OffsetMatrix.Muled(globalMatrix),
		UnitMatrix:   unitMatrix,
		// 物理演算後の移動を受け取ると逆オフセットかけても一部モデルで破綻するので一旦コメントアウト
		// framePosition:   unitMatrix.Translation(),
		FrameRotation: unitMatrix.Quaternion(),
	}
}

func (boneDelta *BoneDelta) FilledGlobalMatrix() *mmath.MMat4 {
	if boneDelta.GlobalMatrix == nil {
		boneDelta.GlobalMatrix = mmath.NewMMat4()
	}
	return boneDelta.GlobalMatrix
}

func (boneDelta *BoneDelta) FilledLocalMatrix() *mmath.MMat4 {
	if boneDelta.LocalMatrix == nil {
		// BOf行列: 自身のボーンのボーンオフセット行列をかけてローカル行列
		boneDelta.LocalMatrix = boneDelta.Bone.Extend.OffsetMatrix.Muled(boneDelta.FilledGlobalMatrix())
	}
	return boneDelta.LocalMatrix
}

func (boneDelta *BoneDelta) FilledGlobalPosition() *mmath.MVec3 {
	if boneDelta.GlobalPosition == nil {
		if boneDelta.GlobalMatrix != nil {
			boneDelta.GlobalPosition = boneDelta.GlobalMatrix.Translation()
		} else {
			boneDelta.GlobalPosition = mmath.NewMVec3()
		}
	}
	return boneDelta.GlobalPosition
}

func (boneDelta *BoneDelta) FilledGlobalRotation() *mmath.MQuaternion {
	return boneDelta.FilledGlobalMatrix().Quaternion()
}

func (boneDelta *BoneDelta) FilledLocalPosition() *mmath.MVec3 {
	pos := boneDelta.FilledFramePosition().Copy()

	if boneDelta.FrameMorphPosition != nil && !boneDelta.FrameMorphPosition.IsZero() {
		pos.Add(boneDelta.FrameMorphPosition)
	}

	// if boneDelta.FrameEffectPosition != nil && !boneDelta.FrameEffectPosition.IsZero() {
	// 	pos.Add(boneDelta.FrameEffectPosition)
	// }

	return pos
}

func (boneDelta *BoneDelta) FilledLocalEffectorPosition(effectorFactor float64) *mmath.MVec3 {
	pos := boneDelta.FilledFramePosition().Copy()

	if boneDelta.FrameMorphPosition != nil && !boneDelta.FrameMorphPosition.IsZero() {
		pos.Add(boneDelta.FrameMorphPosition)
	}

	// if boneDelta.FrameEffectPosition != nil && !boneDelta.FrameEffectPosition.IsZero() {
	// 	pos.Add(boneDelta.FrameEffectPosition)
	// }

	return pos.MuledScalar(effectorFactor)
}

func (boneDelta *BoneDelta) FilledFramePosition() *mmath.MVec3 {
	if boneDelta.FramePosition == nil {
		boneDelta.FramePosition = mmath.NewMVec3()
	}
	return boneDelta.FramePosition
}

func (boneDelta *BoneDelta) FilledFrameMorphPosition() *mmath.MVec3 {
	if boneDelta.FrameMorphPosition == nil {
		boneDelta.FrameMorphPosition = mmath.NewMVec3()
	}
	return boneDelta.FrameMorphPosition
}

func (boneDelta *BoneDelta) FilledFrameLocalMorphPosition() *mmath.MVec3 {
	if boneDelta.FrameLocalMorphPosition == nil {
		boneDelta.FrameLocalMorphPosition = mmath.NewMVec3()
	}
	return boneDelta.FrameLocalMorphPosition
}

func (boneDelta *BoneDelta) FilledLocalRotation() *mmath.MQuaternion {
	rot := boneDelta.FilledFrameRotation().Copy()

	if boneDelta.FrameMorphRotation != nil && !boneDelta.FrameMorphRotation.IsIdent() {
		// rot = boneDelta.FrameMorphRotation.Muled(rot)
		rot.Mul(boneDelta.FrameMorphRotation)
	}

	if boneDelta.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(boneDelta.Bone.Extend.NormalizedFixedAxis)
	}

	return rot
}

func (boneDelta *BoneDelta) FilledLocalEffectorRotation(effectorFactor float64) *mmath.MQuaternion {
	return boneDelta.FilledLocalRotation().MuledScalar(effectorFactor)
}

func (boneDelta *BoneDelta) FilledFrameRotation() *mmath.MQuaternion {
	if boneDelta.FrameRotation == nil {
		boneDelta.FrameRotation = mmath.NewMQuaternion()
	}
	return boneDelta.FrameRotation
}

func (boneDelta *BoneDelta) FilledFrameMorphRotation() *mmath.MQuaternion {
	if boneDelta.FrameMorphRotation == nil {
		boneDelta.FrameMorphRotation = mmath.NewMQuaternion()
	}
	return boneDelta.FrameMorphRotation
}

func (boneDelta *BoneDelta) FilledFrameLocalMorphRotation() *mmath.MQuaternion {
	if boneDelta.FrameLocalMorphRotation == nil {
		boneDelta.FrameLocalMorphRotation = mmath.NewMQuaternion()
	}
	return boneDelta.FrameLocalMorphRotation
}

func (boneDelta *BoneDelta) FilledLocalScale() *mmath.MVec3 {
	pos := boneDelta.FilledFrameScale().Copy()

	if boneDelta.FrameMorphScale != nil && !boneDelta.FrameMorphScale.IsZero() {
		pos.Add(boneDelta.FrameMorphScale)
	}

	return pos
}

func (boneDelta *BoneDelta) FilledFrameScale() *mmath.MVec3 {
	if boneDelta.FrameScale == nil {
		boneDelta.FrameScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return boneDelta.FrameScale
}

func (boneDelta *BoneDelta) FilledFrameMorphScale() *mmath.MVec3 {
	if boneDelta.FrameMorphScale == nil {
		boneDelta.FrameMorphScale = mmath.NewMVec3()
	}
	return boneDelta.FrameMorphScale
}

func (boneDelta *BoneDelta) FilledFrameLocalMorphScale() *mmath.MVec3 {
	if boneDelta.FrameLocalMorphScale == nil {
		boneDelta.FrameLocalMorphScale = mmath.NewMVec3()
	}
	return boneDelta.FrameLocalMorphScale
}

func (boneDelta *BoneDelta) Copy() *BoneDelta {
	return &BoneDelta{
		Bone:                    boneDelta.Bone,
		Frame:                   boneDelta.Frame,
		GlobalMatrix:            boneDelta.FilledGlobalMatrix().Copy(),
		LocalMatrix:             boneDelta.FilledLocalMatrix().Copy(),
		GlobalPosition:          boneDelta.FilledGlobalPosition().Copy(),
		FramePosition:           boneDelta.FilledFramePosition().Copy(),
		FrameMorphPosition:      boneDelta.FilledFrameMorphPosition().Copy(),
		FrameLocalMorphPosition: boneDelta.FilledFrameLocalMorphPosition().Copy(),
		FrameRotation:           boneDelta.FilledFrameRotation().Copy(),
		FrameMorphRotation:      boneDelta.FilledFrameMorphRotation().Copy(),
		FrameLocalMorphRotation: boneDelta.FilledFrameLocalMorphRotation().Copy(),
		FrameScale:              boneDelta.FilledFrameScale().Copy(),
		FrameMorphScale:         boneDelta.FilledFrameMorphScale().Copy(),
		FrameLocalMorphScale:    boneDelta.FilledFrameLocalMorphScale().Copy(),
		UnitMatrix:              boneDelta.UnitMatrix.Copy(),
	}
}

func NewBoneDelta(bone *pmx.Bone, frame float32) *BoneDelta {
	return &BoneDelta{
		Bone:  bone,
		Frame: frame,
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

func (boneDeltas *BoneDeltas) Get(boneIndex int) *BoneDelta {
	if boneIndex < 0 || boneIndex >= len(boneDeltas.Data) {
		return nil
	}

	return boneDeltas.Data[boneIndex]
}

func (boneDeltas *BoneDeltas) GetByName(boneName string) *BoneDelta {
	if _, ok := boneDeltas.Bones.NameIndexes[boneName]; ok {
		return boneDeltas.Data[boneDeltas.Bones.NameIndexes[boneName]]
	}
	return nil
}

func (boneDeltas *BoneDeltas) Update(boneDelta *BoneDelta) {
	boneDeltas.Data[boneDelta.Bone.Index()] = boneDelta
}

func (boneDeltas *BoneDeltas) GetBoneIndexes() []int {
	boneIndexes := make([]int, 0)
	for key := range boneDeltas.Data {
		if !slices.Contains(boneIndexes, key) {
			boneIndexes = append(boneIndexes, key)
		}
	}
	return boneIndexes
}

func (boneDeltas *BoneDeltas) Contains(boneIndex int) bool {
	return boneDeltas.Data[boneIndex] != nil
}

func (boneDeltas *BoneDeltas) GetNearestBoneIndexes(worldPos *mmath.MVec3) []int {
	boneIndexes := make([]int, 0)
	distances := make([]float64, len(boneDeltas.Data))
	for i := range len(boneDeltas.Data) {
		boneDelta := boneDeltas.Get(i)
		if boneDelta == nil {
			continue
		}
		distances[i] = worldPos.Distance(boneDelta.FilledGlobalPosition())
	}
	if len(distances) == 0 {
		return boneIndexes
	}
	sortedDistances := mmath.Float64Slice(distances)
	sortedIndexes := mmath.ArgSort(sortedDistances)
	nearestBone := boneDeltas.Get(sortedIndexes[0])
	for i := range sortedIndexes {
		boneDelta := boneDeltas.Get(sortedIndexes[i])
		if boneDelta == nil {
			continue
		}
		if len(boneIndexes) > 0 {
			if !boneDelta.FilledGlobalPosition().NearEquals(nearestBone.FilledGlobalPosition(), 0.01) {
				break
			}
		}
		boneIndexes = append(boneIndexes, sortedIndexes[i])
	}
	return boneIndexes
}

func (boneDeltas *BoneDeltas) LocalRotation(boneIndex int) *mmath.MQuaternion {
	return boneDeltas.localRotationLoop(boneIndex, 0)
}

func (boneDeltas *BoneDeltas) localRotationLoop(boneIndex int, loop int) *mmath.MQuaternion {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}
	rot := boneDelta.FilledLocalRotation()

	if boneDelta.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := boneDeltas.localRotationLoop(boneDelta.Bone.EffectIndex, loop+1)
		rot.Mul(effectorRot.MuledScalar(boneDelta.Bone.EffectFactor))
	}

	return rot
}

func (boneDeltas *BoneDeltas) LocalPosition(boneIndex int) *mmath.MVec3 {
	return boneDeltas.localPositionLoop(boneIndex, 0)
}

func (boneDeltas *BoneDeltas) localPositionLoop(boneIndex int, loop int) *mmath.MVec3 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMVec3()
	}
	pos := boneDelta.FilledLocalPosition()

	if boneDelta.Bone.IsEffectorTranslation() {
		// 付与親移動がある場合、再帰で回転を取得する
		effectorPos := boneDeltas.localPositionLoop(boneDelta.Bone.EffectIndex, loop+1)
		pos.Mul(effectorPos.MuledScalar(boneDelta.Bone.EffectFactor))
	}

	return pos
}
