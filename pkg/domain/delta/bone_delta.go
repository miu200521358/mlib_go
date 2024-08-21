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
	FrameLocalPosition      *mmath.MVec3       // キーフレ位置のローカル変動量
	FrameMorphPosition      *mmath.MVec3       // モーフ位置の変動量
	FrameLocalMorphPosition *mmath.MVec3       // モーフ位置のローカル変動量
	FrameRotation           *mmath.MQuaternion // キーフレ回転の変動量
	FrameLocalRotation      *mmath.MQuaternion // キーフレ回転のローカル変動量
	FrameMorphRotation      *mmath.MQuaternion // モーフ回転の変動量
	FrameLocalMorphRotation *mmath.MQuaternion // モーフ回転のローカル変動量
	FrameScale              *mmath.MVec3       // キーフレスケールの変動量
	FrameLocalScale         *mmath.MVec3       // キーフレスケールのローカル変動量
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

func (boneDelta *BoneDelta) FilledTotalLocalPosition() *mmath.MVec3 {
	pos := boneDelta.FilledFrameLocalPosition().Copy()

	if boneDelta.FrameLocalMorphPosition != nil && !boneDelta.FrameLocalMorphPosition.IsZero() {
		pos.Add(boneDelta.FrameLocalMorphPosition)
	}

	return pos
}

func (boneDelta *BoneDelta) FilledFramePosition() *mmath.MVec3 {
	if boneDelta.FramePosition == nil {
		boneDelta.FramePosition = mmath.NewMVec3()
	}
	return boneDelta.FramePosition
}

func (boneDelta *BoneDelta) FilledFrameLocalPosition() *mmath.MVec3 {
	if boneDelta.FrameLocalPosition == nil {
		boneDelta.FrameLocalPosition = mmath.NewMVec3()
	}
	return boneDelta.FrameLocalPosition
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

func (boneDelta *BoneDelta) FilledTotalLocalRot() *mmath.MQuaternion {
	rot := boneDelta.FilledFrameLocalRotation().Copy()

	if boneDelta.FrameLocalMorphRotation != nil && !boneDelta.FrameLocalMorphRotation.IsIdent() {
		rot.Mul(boneDelta.FrameLocalMorphRotation)
	}

	// if boneDelta.Bone.HasFixedAxis() {
	// 	rot = rot.ToFixedAxisRotation(boneDelta.Bone.Extend.NormalizedFixedAxis)
	// }

	return rot
}

func (boneDelta *BoneDelta) FilledFrameRotation() *mmath.MQuaternion {
	if boneDelta.FrameRotation == nil {
		boneDelta.FrameRotation = mmath.NewMQuaternion()
	}
	return boneDelta.FrameRotation
}

func (boneDelta *BoneDelta) FilledFrameLocalRotation() *mmath.MQuaternion {
	if boneDelta.FrameLocalRotation == nil {
		boneDelta.FrameLocalRotation = mmath.NewMQuaternion()
	}
	return boneDelta.FrameLocalRotation
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

func (boneDelta *BoneDelta) FilledTotalScale() *mmath.MVec3 {
	scale := boneDelta.FilledFrameScale().Copy()

	if boneDelta.FrameMorphScale != nil && !boneDelta.FrameMorphScale.IsZero() {
		scale.Mul(boneDelta.FrameMorphScale)
	}

	return scale
}

func (boneDelta *BoneDelta) FilledTotalLocalScale() *mmath.MVec3 {
	scale := boneDelta.FilledFrameLocalScale().Copy()

	if boneDelta.FrameLocalMorphScale != nil && !boneDelta.FrameLocalMorphScale.IsZero() {
		scale.Mul(boneDelta.FrameLocalMorphScale)
	}

	return scale
}

func (boneDelta *BoneDelta) FilledFrameScale() *mmath.MVec3 {
	if boneDelta.FrameScale == nil {
		boneDelta.FrameScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return boneDelta.FrameScale
}

func (boneDelta *BoneDelta) FilledFrameLocalScale() *mmath.MVec3 {
	if boneDelta.FrameLocalScale == nil {
		boneDelta.FrameLocalScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return boneDelta.FrameLocalScale
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

func (boneDeltas *BoneDeltas) TotalRotationMat(boneIndex int) *mmath.MMat4 {
	return boneDeltas.totalRotationLoop(boneIndex, 0, 1.0).ToMat4()
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

func (boneDeltas *BoneDeltas) TotalLocalRotationMat(boneIndex int) *mmath.MMat4 {
	return boneDeltas.totalLocalRotationLoop(boneIndex, 0)
}

func (boneDeltas *BoneDeltas) totalLocalRotationLoop(boneIndex int, loop int) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMMat4()
	}

	var parentRotMat *mmath.MMat4
	if boneDelta.Bone.ParentIndex >= 0 {
		parentBoneDelta := boneDeltas.Get(boneDelta.Bone.ParentIndex)
		parentRotMat = parentBoneDelta.FilledTotalLocalRot().ToMat4()
	}

	rot := boneDelta.FilledTotalLocalRot()
	if rot.IsIdent() {
		if parentRotMat == nil {
			// 親の回転が定義されていない場合、もしくは親を引き継ぐ場合、単位行列を返す
			return mmath.NewMMat4()
		}
		// 親の回転が指定されている場合、キャンセルする
		return parentRotMat.Inverted()
	}

	// 回転行列
	rotMat := rot.ToMat4()

	if parentRotMat == nil {
		return rotMat
	}

	// 親の回転をキャンセルする
	return rotMat.Muled(parentRotMat.Inverted())
}

func (boneDeltas *BoneDeltas) TotalPositionMat(boneIndex int) *mmath.MMat4 {
	return boneDeltas.totalPositionLoop(boneIndex, 0).ToMat4()
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

func (boneDeltas *BoneDeltas) TotalLocalPositionMat(boneIndex int) *mmath.MMat4 {
	return boneDeltas.totalLocalPositionLoop(boneIndex, 0)
}

func (boneDeltas *BoneDeltas) totalLocalPositionLoop(boneIndex int, loop int) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMMat4()
	}

	var parentPosMat *mmath.MMat4
	if boneDelta.Bone.ParentIndex >= 0 {
		parentBoneDelta := boneDeltas.Get(boneDelta.Bone.ParentIndex)
		parentPos := parentBoneDelta.FilledTotalLocalPosition()

		if !parentPos.IsZero() {
			// 親の位置行列
			parentPosMat = parentPos.ToMat4()
		}
	}

	pos := boneDelta.FilledTotalLocalPosition()
	if pos.IsZero() {
		if parentPosMat == nil {
			// 親の位置が定義されていない場合、単位行列を返す
			return mmath.NewMMat4()
		}

		// 親の位置が指定されている場合、キャンセルする
		return parentPosMat.Inverted()
	}

	// ローカル軸に沿った位置行列
	posMat := pos.ToMat4()

	if parentPosMat == nil {
		return posMat
	}

	// 親の位置をキャンセルする
	return posMat.Muled(parentPosMat.Inverted())
}

func (boneDeltas *BoneDeltas) TotalScaleMat(boneIndex int) *mmath.MMat4 {
	return boneDeltas.totalScaleMatLoop(boneIndex, 0).ToScaleMat4()
}

func (boneDeltas *BoneDeltas) totalScaleMatLoop(boneIndex int, loop int) *mmath.MVec3 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	scale := boneDelta.FilledTotalScale()

	return scale
}

func (boneDeltas *BoneDeltas) TotalLocalScaleMat(boneIndex int) *mmath.MMat4 {
	return boneDeltas.totalLocalScaleMatLoop(boneIndex, 0)
}

func (boneDeltas *BoneDeltas) totalLocalScaleMatLoop(boneIndex int, loop int) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMMat4()
	}

	var parentScaleMat *mmath.MMat4
	if boneDelta.Bone.ParentIndex >= 0 {
		parentBoneDelta := boneDeltas.Get(boneDelta.Bone.ParentIndex)
		parentScale := parentBoneDelta.FilledTotalLocalScale()

		if !parentScale.IsOne() {
			// 親のローカル軸に沿ったスケール行列
			parentScaleMat = parentBoneDelta.Bone.Extend.LocalAxis.ToScaleLocalMat(parentScale)
		}
	}

	scale := boneDelta.FilledTotalLocalScale()
	if scale.IsOne() {
		if parentScaleMat == nil {
			// 親のスケールが定義されていない場合、単位行列を返す
			return mmath.NewMMat4()
		}
		// 親のスケールが指定されている場合、キャンセルする
		return parentScaleMat.Inverted()
	}

	// ローカル軸に沿ったスケール行列
	scaleMat := boneDelta.Bone.Extend.LocalAxis.ToScaleLocalMat(scale)

	if parentScaleMat == nil {
		return scaleMat
	}

	// 親のスケールをキャンセルする
	return scaleMat.Muled(parentScaleMat.Inverted())
}
