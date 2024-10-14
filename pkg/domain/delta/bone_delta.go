package delta

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type BoneDelta struct {
	Bone                         *pmx.Bone          // ボーン
	Frame                        float32            // キーフレーム
	GlobalIkOffMatrix            *mmath.MMat4       // IKオフ時のグローバル行列
	GlobalMatrix                 *mmath.MMat4       // グローバル行列
	LocalMatrix                  *mmath.MMat4       // ローカル行列
	GlobalPosition               *mmath.MVec3       // グローバル位置
	UnitMatrix                   *mmath.MMat4       // 親ボーンからの変位行列
	FramePosition                *mmath.MVec3       // キーフレ位置の変動量
	FrameMorphPosition           *mmath.MVec3       // モーフ位置の変動量
	FrameCancelablePosition      *mmath.MVec3       // キャンセル位置の変動量
	FrameMorphCancelablePosition *mmath.MVec3       // モーフキャンセル位置の変動量
	FrameRotation                *mmath.MQuaternion // キーフレ回転の変動量
	FrameMorphRotation           *mmath.MQuaternion // モーフ回転の変動量
	FrameCancelableRotation      *mmath.MQuaternion // キャンセル回転の変動量
	FrameMorphCancelableRotation *mmath.MQuaternion // モーフキャンセル回転の変動量
	FrameScale                   *mmath.MVec3       // キーフレスケールの変動量
	FrameMorphScale              *mmath.MVec3       // モーフスケールの変動量
	FrameCancelableScale         *mmath.MVec3       // キャンセルスケールの変動量
	FrameMorphCancelableScale    *mmath.MVec3       // モーフキャンセルスケールの変動量
	FrameLocalMat                *mmath.MMat4       // キーフレのローカル変動行列
	FrameLocalMorphMat           *mmath.MMat4       // モーフのローカル変動行列
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
	unitMatrix := globalMatrix.Inverted().Muled(parentGlobalMatrix)

	return &BoneDelta{
		Bone:         bone,
		Frame:        frame,
		GlobalMatrix: globalMatrix,
		LocalMatrix:  globalMatrix.Muled(bone.Extend.OffsetMatrix),
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
		boneDelta.LocalMatrix = boneDelta.FilledGlobalMatrix().Muled(boneDelta.Bone.Extend.OffsetMatrix)
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

func (boneDelta *BoneDelta) FilledGlobalBoneRotation() *mmath.MQuaternion {
	return boneDelta.FilledGlobalMatrix().Quaternion()
}

func (boneDelta *BoneDelta) FilledGlobalRotation() *mmath.MQuaternion {
	return boneDelta.FilledGlobalMatrix().Quaternion()
}

func (boneDelta *BoneDelta) FilledTotalPosition() *mmath.MVec3 {
	pos := boneDelta.FilledFramePosition().Copy()

	if boneDelta.FrameMorphPosition != nil && !boneDelta.FrameMorphPosition.IsZero() {
		pos.Add(boneDelta.FrameMorphPosition)
	}

	return pos
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

func (boneDelta *BoneDelta) FilledFrameCancelablePosition() *mmath.MVec3 {
	if boneDelta.FrameCancelablePosition == nil {
		boneDelta.FrameCancelablePosition = mmath.NewMVec3()
	}
	return boneDelta.FrameCancelablePosition
}

func (boneDelta *BoneDelta) FilledFrameMorphCancelablePosition() *mmath.MVec3 {
	if boneDelta.FrameMorphCancelablePosition == nil {
		boneDelta.FrameMorphCancelablePosition = mmath.NewMVec3()
	}
	return boneDelta.FrameMorphCancelablePosition
}

func (boneDelta *BoneDelta) FilledTotalRotation() *mmath.MQuaternion {
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

func (boneDelta *BoneDelta) FilledFrameCancelableRotation() *mmath.MQuaternion {
	if boneDelta.FrameCancelableRotation == nil {
		boneDelta.FrameCancelableRotation = mmath.NewMQuaternion()
	}
	return boneDelta.FrameCancelableRotation
}

func (boneDelta *BoneDelta) FilledFrameMorphCancelableRotation() *mmath.MQuaternion {
	if boneDelta.FrameMorphCancelableRotation == nil {
		boneDelta.FrameMorphCancelableRotation = mmath.NewMQuaternion()
	}
	return boneDelta.FrameMorphCancelableRotation
}

func (boneDelta *BoneDelta) FilledTotalScale() *mmath.MVec3 {
	scale := boneDelta.FilledFrameScale().Copy()

	if boneDelta.FrameMorphScale != nil && !boneDelta.FrameMorphScale.IsOne() {
		scale.Mul(boneDelta.FrameMorphScale)
	}

	return scale
}

func (boneDelta *BoneDelta) FilledFrameScale() *mmath.MVec3 {
	if boneDelta.FrameScale == nil {
		boneDelta.FrameScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return boneDelta.FrameScale
}

func (boneDelta *BoneDelta) FilledFrameMorphScale() *mmath.MVec3 {
	if boneDelta.FrameMorphScale == nil {
		boneDelta.FrameMorphScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return boneDelta.FrameMorphScale
}

func (boneDelta *BoneDelta) FilledFrameCancelableScale() *mmath.MVec3 {
	if boneDelta.FrameCancelableScale == nil {
		boneDelta.FrameCancelableScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return boneDelta.FrameCancelableScale
}

func (boneDelta *BoneDelta) FilledFrameMorphCancelableScale() *mmath.MVec3 {
	if boneDelta.FrameMorphCancelableScale == nil {
		boneDelta.FrameMorphCancelableScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return boneDelta.FrameMorphCancelableScale
}

func (boneDelta *BoneDelta) FilledTotalLocalMat() *mmath.MMat4 {
	mat := boneDelta.FilledFrameLocalMat().Copy()

	if boneDelta.FrameLocalMorphMat != nil && !boneDelta.FrameLocalMorphMat.IsIdent() {
		mat.Mul(boneDelta.FrameLocalMorphMat)
	}

	return mat
}

func (boneDelta *BoneDelta) FilledFrameLocalMat() *mmath.MMat4 {
	if boneDelta.FrameLocalMat == nil {
		boneDelta.FrameLocalMat = mmath.NewMMat4()
	}
	return boneDelta.FrameLocalMat
}

func (boneDelta *BoneDelta) FilledFrameLocalMorphMat() *mmath.MMat4 {
	if boneDelta.FrameLocalMorphMat == nil {
		boneDelta.FrameLocalMorphMat = mmath.NewMMat4()
	}
	return boneDelta.FrameLocalMorphMat
}

func (boneDelta *BoneDelta) Copy() *BoneDelta {
	return &BoneDelta{
		Bone:                         boneDelta.Bone,
		Frame:                        boneDelta.Frame,
		GlobalMatrix:                 boneDelta.FilledGlobalMatrix().Copy(),
		LocalMatrix:                  boneDelta.FilledLocalMatrix().Copy(),
		GlobalPosition:               boneDelta.FilledGlobalPosition().Copy(),
		FramePosition:                boneDelta.FilledFramePosition().Copy(),
		FrameMorphPosition:           boneDelta.FilledFrameMorphPosition().Copy(),
		FrameCancelablePosition:      boneDelta.FilledFrameCancelablePosition().Copy(),
		FrameMorphCancelablePosition: boneDelta.FilledFrameMorphCancelablePosition().Copy(),
		FrameRotation:                boneDelta.FilledFrameRotation().Copy(),
		FrameMorphRotation:           boneDelta.FilledFrameMorphRotation().Copy(),
		FrameCancelableRotation:      boneDelta.FilledFrameCancelableRotation().Copy(),
		FrameMorphCancelableRotation: boneDelta.FilledFrameMorphCancelableRotation().Copy(),
		FrameScale:                   boneDelta.FilledFrameScale().Copy(),
		FrameMorphScale:              boneDelta.FilledFrameMorphScale().Copy(),
		FrameCancelableScale:         boneDelta.FilledFrameCancelableScale().Copy(),
		FrameMorphCancelableScale:    boneDelta.FilledFrameMorphCancelableScale().Copy(),
		FrameLocalMat:                boneDelta.FilledFrameLocalMat().Copy(),
		FrameLocalMorphMat:           boneDelta.FilledFrameLocalMorphMat().Copy(),
		UnitMatrix:                   boneDelta.UnitMatrix.Copy(),
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
	if boneIndex < 0 || boneIndex >= len(boneDeltas.Data) {
		return false
	}
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

func (boneDeltas *BoneDeltas) TotalRotationMat(boneIndex int) *mmath.MMat4 {
	rotMat := boneDeltas.totalRotationLoop(boneIndex, 0, 1.0).ToMat4()
	return boneDeltas.totalCancelRotation(boneIndex, rotMat)
}

func (boneDeltas *BoneDeltas) totalCancelRotation(boneIndex int, rotMat *mmath.MMat4) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)

	// 親のキャンセル付き回転行列
	var parentCancelableRotMat *mmath.MMat4
	if boneDeltas.Contains(boneDelta.Bone.ParentIndex) {
		parentBoneDelta := boneDeltas.Get(boneDelta.Bone.ParentIndex)
		if parentBoneDelta.FrameCancelableRotation != nil && !parentBoneDelta.FrameCancelableRotation.IsIdent() {
			parentCancelableRotMat = parentBoneDelta.FrameCancelableRotation.ToMat4()
		}
		if parentBoneDelta.FrameMorphCancelableRotation != nil && !parentBoneDelta.FrameMorphCancelableRotation.IsIdent() {
			if parentCancelableRotMat == nil {
				parentCancelableRotMat = parentBoneDelta.FrameMorphCancelableRotation.ToMat4()
			} else {
				parentCancelableRotMat = parentCancelableRotMat.Muled(parentBoneDelta.FrameMorphCancelableRotation.ToMat4())
			}
		}
	}

	// キャンセル付き回転
	if (boneDelta.FrameCancelableRotation == nil || boneDelta.FrameCancelableRotation.IsIdent()) &&
		(boneDelta.FrameMorphCancelableRotation == nil || boneDelta.FrameMorphCancelableRotation.IsIdent()) {
		// 親の回転をキャンセルする
		if parentCancelableRotMat == nil {
			return rotMat
		}
		return rotMat.Muled(parentCancelableRotMat.Inverted())
	}

	if parentCancelableRotMat == nil {
		if boneDelta.FrameCancelableRotation != nil && !boneDelta.FrameCancelableRotation.IsIdent() {
			rotMat = rotMat.Muled(boneDelta.FrameCancelableRotation.ToMat4())
		}
		if boneDelta.FrameMorphCancelableRotation != nil && !boneDelta.FrameMorphCancelableRotation.IsIdent() {
			rotMat = rotMat.Muled(boneDelta.FrameMorphCancelableRotation.ToMat4())
		}
		return rotMat
	}

	if boneDelta.FrameCancelableRotation != nil && !boneDelta.FrameCancelableRotation.IsIdent() {
		rotMat = rotMat.Muled(boneDelta.FrameCancelableRotation.ToMat4())
	}
	if boneDelta.FrameMorphCancelableRotation != nil && !boneDelta.FrameMorphCancelableRotation.IsIdent() {
		rotMat = rotMat.Muled(boneDelta.FrameMorphCancelableRotation.ToMat4())
	}

	// 親の回転をキャンセルする
	return rotMat.Muled(parentCancelableRotMat.Inverted())
}

func (boneDeltas *BoneDeltas) totalRotationLoop(boneIndex int, loop int, factor float64) *mmath.MQuaternion {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}

	rot := boneDelta.FilledTotalRotation()

	if boneDelta.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := boneDeltas.totalRotationLoop(boneDelta.Bone.EffectIndex, loop+1, boneDelta.Bone.EffectFactor)
		rot.Mul(effectorRot)
	}

	return rot.MuledScalar(factor)
}

// 該当ボーンまでの付与親を加味した全ての回転（モーフは含まない）
func (boneDeltas *BoneDeltas) TotalBoneRotation(boneIndex int) *mmath.MQuaternion {
	rot := boneDeltas.totalBoneRotationLoop(boneIndex, 0, 1.0)
	return boneDeltas.totalBoneCancelRotation(boneIndex, rot)
}

func (boneDeltas *BoneDeltas) totalBoneCancelRotation(boneIndex int, rot *mmath.MQuaternion) *mmath.MQuaternion {
	boneDelta := boneDeltas.Get(boneIndex)

	// 親のキャンセル付き回転行列
	var parentCancelableRotMat *mmath.MMat4
	if boneDeltas.Contains(boneDelta.Bone.ParentIndex) {
		parentBoneDelta := boneDeltas.Get(boneDelta.Bone.ParentIndex)
		if parentBoneDelta.FrameCancelableRotation != nil && !parentBoneDelta.FrameCancelableRotation.IsIdent() {
			parentCancelableRotMat = parentBoneDelta.FrameCancelableRotation.ToMat4()
		}
	}

	// キャンセル付き回転
	if boneDelta.FrameCancelableRotation == nil || boneDelta.FrameCancelableRotation.IsIdent() {
		// 親の回転をキャンセルする
		if parentCancelableRotMat == nil {
			return rot
		}
		return rot.ToMat4().Muled(parentCancelableRotMat.Inverted()).Quaternion()
	}

	if parentCancelableRotMat == nil {
		if boneDelta.FrameCancelableRotation != nil && !boneDelta.FrameCancelableRotation.IsIdent() {
			rot = rot.ToMat4().Muled(boneDelta.FrameCancelableRotation.ToMat4()).Quaternion()
		}
		return rot
	}

	if boneDelta.FrameCancelableRotation != nil && !boneDelta.FrameCancelableRotation.IsIdent() {
		rot = rot.ToMat4().Muled(boneDelta.FrameCancelableRotation.ToMat4()).Quaternion()
	}

	// 親の回転をキャンセルする
	return rot.ToMat4().Muled(parentCancelableRotMat.Inverted()).Quaternion()
}

func (boneDeltas *BoneDeltas) totalBoneRotationLoop(boneIndex int, loop int, factor float64) *mmath.MQuaternion {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}

	rot := boneDelta.FilledFrameRotation().Copy()

	if boneDelta.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := boneDeltas.totalBoneRotationLoop(boneDelta.Bone.EffectIndex, loop+1, boneDelta.Bone.EffectFactor)
		rot.Mul(effectorRot)
	}

	rot = rot.MuledScalar(factor)

	if boneDelta.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(boneDelta.Bone.Extend.NormalizedFixedAxis)
	}

	return rot
}

func (boneDeltas *BoneDeltas) TotalPositionMat(boneIndex int) *mmath.MMat4 {
	posMat := boneDeltas.totalPositionLoop(boneIndex, 0).ToMat4()
	return boneDeltas.totalCancelPosition(boneIndex, posMat)
}

func (boneDeltas *BoneDeltas) totalCancelPosition(boneIndex int, posMat *mmath.MMat4) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)

	// 親のキャンセル付き移動行列
	var parentCancelablePosMat *mmath.MMat4
	if boneDeltas.Contains(boneDelta.Bone.ParentIndex) {
		parentBoneDelta := boneDeltas.Get(boneDelta.Bone.ParentIndex)
		if parentBoneDelta.FrameCancelablePosition != nil && !parentBoneDelta.FrameCancelablePosition.IsZero() {
			parentCancelablePosMat = parentBoneDelta.FrameCancelablePosition.ToMat4()
		}
		if parentBoneDelta.FrameMorphCancelablePosition != nil && !parentBoneDelta.FrameMorphCancelablePosition.IsZero() {
			if parentCancelablePosMat == nil {
				parentCancelablePosMat = parentBoneDelta.FrameMorphCancelablePosition.ToMat4()
			} else {
				parentCancelablePosMat = parentCancelablePosMat.Muled(parentBoneDelta.FrameMorphCancelablePosition.ToMat4())
			}
		}
	}

	// キャンセル付き移動
	if (boneDelta.FrameCancelablePosition == nil || boneDelta.FrameCancelablePosition.IsZero()) &&
		(boneDelta.FrameMorphCancelablePosition == nil || boneDelta.FrameMorphCancelablePosition.IsZero()) {
		// 親の移動をキャンセルする
		if parentCancelablePosMat == nil {
			return posMat
		}
		return posMat.Muled(parentCancelablePosMat.Inverted())
	}

	if parentCancelablePosMat == nil {
		if boneDelta.FrameCancelablePosition != nil && !boneDelta.FrameCancelablePosition.IsZero() {
			posMat = posMat.Muled(boneDelta.FrameCancelablePosition.ToMat4())
		}
		if boneDelta.FrameMorphCancelablePosition != nil && !boneDelta.FrameMorphCancelablePosition.IsZero() {
			posMat = posMat.Muled(boneDelta.FrameMorphCancelablePosition.ToMat4())
		}
		return posMat
	}

	if boneDelta.FrameCancelablePosition != nil && !boneDelta.FrameCancelablePosition.IsZero() {
		posMat = posMat.Muled(boneDelta.FrameCancelablePosition.ToMat4())
	}
	if boneDelta.FrameMorphCancelablePosition != nil && !boneDelta.FrameMorphCancelablePosition.IsZero() {
		posMat = posMat.Muled(boneDelta.FrameMorphCancelablePosition.ToMat4())
	}

	// 親の移動をキャンセルする
	return posMat.Muled(parentCancelablePosMat.Inverted())
}

func (boneDeltas *BoneDeltas) totalPositionLoop(boneIndex int, loop int) *mmath.MVec3 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMVec3()
	}
	pos := boneDelta.FilledTotalPosition()

	if boneDelta.Bone.IsEffectorTranslation() {
		// 付与親移動がある場合、再帰で回転を取得する
		effectorPos := boneDeltas.totalPositionLoop(boneDelta.Bone.EffectIndex, loop+1)
		pos.Mul(effectorPos.MuledScalar(boneDelta.Bone.EffectFactor))
	}

	return pos
}

func (boneDeltas *BoneDeltas) TotalScaleMat(boneIndex int) *mmath.MMat4 {
	scaleMat := boneDeltas.totalScaleMatLoop(boneIndex, 0).ToScaleMat4()
	return boneDeltas.totalCancelScale(boneIndex, scaleMat)
}

func (boneDeltas *BoneDeltas) totalScaleMatLoop(boneIndex int, loop int) *mmath.MVec3 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	scale := boneDelta.FilledTotalScale()

	return scale
}

func (boneDeltas *BoneDeltas) totalCancelScale(boneIndex int, scaleMat *mmath.MMat4) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)

	// 親のキャンセル付きスケール行列
	var parentCancelableScaleMat *mmath.MMat4
	if boneDeltas.Contains(boneDelta.Bone.ParentIndex) {
		parentBoneDelta := boneDeltas.Get(boneDelta.Bone.ParentIndex)
		if parentBoneDelta.FrameCancelableScale != nil && !parentBoneDelta.FrameCancelableScale.IsZero() {
			parentCancelableScaleMat = parentBoneDelta.FrameCancelableScale.ToScaleMat4()
		}
		if parentBoneDelta.FrameMorphCancelableScale != nil && !parentBoneDelta.FrameMorphCancelableScale.IsZero() {
			if parentCancelableScaleMat == nil {
				parentCancelableScaleMat = parentBoneDelta.FrameMorphCancelableScale.ToScaleMat4()
			} else {
				parentCancelableScaleMat = parentCancelableScaleMat.Muled(parentBoneDelta.FrameMorphCancelableScale.ToScaleMat4())
			}
		}
	}

	// キャンセル付きスケール
	if (boneDelta.FrameCancelableScale == nil || boneDelta.FrameCancelableScale.IsZero()) ||
		(boneDelta.FrameMorphCancelableScale == nil || boneDelta.FrameMorphCancelableScale.IsZero()) {
		// 親のスケールをキャンセルする
		if parentCancelableScaleMat == nil {
			return scaleMat
		}
		return scaleMat.Muled(parentCancelableScaleMat.Inverted())
	}

	if parentCancelableScaleMat == nil {
		if boneDelta.FrameCancelableScale != nil && !boneDelta.FrameCancelableScale.IsZero() {
			scaleMat = scaleMat.Muled(boneDelta.FrameCancelableScale.ToScaleMat4())
		}
		if boneDelta.FrameMorphCancelableScale != nil && !boneDelta.FrameMorphCancelableScale.IsZero() {
			scaleMat = scaleMat.Muled(boneDelta.FrameMorphCancelableScale.ToScaleMat4())
		}
		return scaleMat
	}

	if boneDelta.FrameCancelableScale != nil && !boneDelta.FrameCancelableScale.IsZero() {
		scaleMat = scaleMat.Muled(boneDelta.FrameCancelableScale.ToScaleMat4())
	}
	if boneDelta.FrameMorphCancelableScale != nil && !boneDelta.FrameMorphCancelableScale.IsZero() {
		scaleMat = scaleMat.Muled(boneDelta.FrameMorphCancelableScale.ToScaleMat4())
	}

	// 親のスケールをキャンセルする
	return scaleMat.Muled(parentCancelableScaleMat.Inverted())
}

func (boneDeltas *BoneDeltas) TotalLocalMat(boneIndex int) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil {
		return mmath.NewMMat4()
	}

	// ローカル変換行列
	return boneDelta.FilledTotalLocalMat()
}
