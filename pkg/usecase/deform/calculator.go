package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

const maxEffectorRecursion = 10

// ----------------------------------------------------------------------------
// 公開メソッド (回転系)
// ----------------------------------------------------------------------------

// calculateTotalRotationMat ボーンの「回転」(モーフ含む) を再帰的に合成したマトリックスを求める
func calculateTotalRotationMat(
	deltas *delta.BoneDeltas, boneIndex int,
) *mmath.MMat4 {
	rotMat := accumulateTotalRotation(deltas, boneIndex, 0, 1.0).ToMat4()
	return applyCancelableRotation(deltas, boneIndex, rotMat)
}

// CalculateBoneRotation モーフを含まない「純粋なボーンの回転」だけを求める
func CalculateBoneRotation(
	deltas *delta.BoneDeltas, boneIndex int,
) *mmath.MQuaternion {
	rot := accumulateBoneRotation(deltas, boneIndex, 0, 1.0)
	return cancelBoneRotation(deltas, boneIndex, rot)
}

// ----------------------------------------------------------------------------
// 公開メソッド (位置系)
// ----------------------------------------------------------------------------

// calculateTotalPositionMat ボーンの「位置」を再帰的に合成したマトリックスを求める
func calculateTotalPositionMat(
	deltas *delta.BoneDeltas, boneIndex int,
) *mmath.MMat4 {
	posMat := accumulateTotalPosition(deltas, boneIndex, 0).ToMat4()
	return applyCancelablePosition(deltas, boneIndex, posMat)
}

// ----------------------------------------------------------------------------
// 公開メソッド (スケール系)
// ----------------------------------------------------------------------------

// calculateTotalScaleMat ボーンの「スケール」を再帰的に合成したマトリックスを求める
func calculateTotalScaleMat(
	deltas *delta.BoneDeltas, boneIndex int,
) *mmath.MMat4 {
	scaleMat := accumulateTotalScale(deltas, boneIndex, 0).ToScaleMat4()
	return applyCancelableScale(deltas, boneIndex, scaleMat)
}

// ----------------------------------------------------------------------------
// 公開メソッド (ローカル行列)
// ----------------------------------------------------------------------------

// calculateTotalLocalMat ボーンの「ローカル行列」を求める
func calculateTotalLocalMat(
	deltas *delta.BoneDeltas, boneIndex int,
) *mmath.MMat4 {
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return mmath.NewMMat4()
	}
	return bd.FilledTotalLocalMat()
}

// ----------------------------------------------------------------------------
// 以下、再帰的に合成するための非公開ヘルパー関数
// ----------------------------------------------------------------------------

// 回転の合成 (モーフ含む)
func accumulateTotalRotation(
	deltas *delta.BoneDeltas,
	boneIndex int,
	recursion int,
	factor float64,
) *mmath.MQuaternion {
	if recursion > maxEffectorRecursion {
		return mmath.NewMQuaternion()
	}
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return mmath.NewMQuaternion()
	}

	// すでに合成済みのトータル回転を前提として使う場合
	rot := bd.FilledTotalRotation()
	if rot == nil {
		// 未計算 or 取得不可なら単位回転
		rot = mmath.NewMQuaternion()
	}

	// ボーンが回転付与を持つ場合、エフェクタ先の回転を再帰合成
	if bd.Bone.IsEffectorRotation() {
		effectorRot := accumulateTotalRotation(deltas, bd.Bone.EffectIndex, recursion+1, bd.Bone.EffectFactor)
		rot.Mul(effectorRot)
	}

	return rot.MuledScalar(factor)
}

// 純粋なボーン回転だけ合成 (モーフを含めない)
func accumulateBoneRotation(
	deltas *delta.BoneDeltas,
	boneIndex int,
	recursion int,
	factor float64,
) *mmath.MQuaternion {
	if recursion > maxEffectorRecursion {
		return mmath.NewMQuaternion()
	}

	bd := deltas.Get(boneIndex)
	if bd == nil {
		return mmath.NewMQuaternion()
	}

	// フレーム上の回転
	rot := mmath.NewMQuaternion()
	if bd.FrameRotation != nil {
		rot = bd.FrameRotation.Copy()
	}

	// エフェクタ回転がある場合のみ再帰
	if bd.Bone.IsEffectorRotation() {
		effectorRot := accumulateBoneRotation(deltas, bd.Bone.EffectIndex, recursion+1, bd.Bone.EffectFactor)
		rot.Mul(effectorRot)
	}

	// 固定軸などがあればここで適用
	if bd.Bone.HasFixedAxis() && bd.Bone.NormalizedFixedAxis != nil {
		rot = rot.ToFixedAxisRotation(bd.Bone.NormalizedFixedAxis)
	}

	return rot.MuledScalar(factor)
}

// 位置合成
func accumulateTotalPosition(
	deltas *delta.BoneDeltas,
	boneIndex int,
	recursion int,
) *mmath.MVec3 {
	if recursion > maxEffectorRecursion {
		return mmath.NewMVec3()
	}
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return mmath.NewMVec3()
	}

	pos := bd.FilledTotalPosition().Copy()

	// 移動付与があれば再帰的に合成
	if bd.Bone.IsEffectorTranslation() {
		effectorPos := accumulateTotalPosition(deltas, bd.Bone.EffectIndex, recursion+1)
		pos.Mul(effectorPos.MuledScalar(bd.Bone.EffectFactor))
	}

	return pos
}

// スケール合成
func accumulateTotalScale(
	deltas *delta.BoneDeltas,
	boneIndex int,
	recursion int,
) *mmath.MVec3 {
	if recursion > maxEffectorRecursion {
		// デフォルトスケール = (1, 1, 1)
		return &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}

	return bd.FilledTotalScale()
}

// ----------------------------------------------------------------------------
// 以下、キャンセル行列を適用するヘルパー
// ----------------------------------------------------------------------------

// 回転キャンセル適用
func applyCancelableRotation(
	deltas *delta.BoneDeltas,
	boneIndex int,
	rotMat *mmath.MMat4,
) *mmath.MMat4 {
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return rotMat
	}

	parentMat := getParentCancelableRotationMat(deltas, bd.Bone.ParentIndex)

	// 自身のキャンセル成分が空なら、そのまま親をキャンセル
	hasSelfCancel := (bd.FrameCancelableRotation != nil && !bd.FrameCancelableRotation.IsIdent()) ||
		(bd.FrameMorphCancelableRotation != nil && !bd.FrameMorphCancelableRotation.IsIdent())

	if !hasSelfCancel {
		if parentMat == nil {
			return rotMat
		}
		return rotMat.Muled(parentMat.Inverted())
	}

	// 自身のキャンセル行列を適用
	if bd.FrameCancelableRotation != nil && !bd.FrameCancelableRotation.IsIdent() {
		rotMat = rotMat.Muled(bd.FrameCancelableRotation.ToMat4())
	}
	if bd.FrameMorphCancelableRotation != nil && !bd.FrameMorphCancelableRotation.IsIdent() {
		rotMat = rotMat.Muled(bd.FrameMorphCancelableRotation.ToMat4())
	}

	// 親のキャンセル
	if parentMat == nil {
		return rotMat
	}
	return rotMat.Muled(parentMat.Inverted())
}

// 純粋なボーン回転のキャンセル適用
func cancelBoneRotation(
	deltas *delta.BoneDeltas,
	boneIndex int,
	rot *mmath.MQuaternion,
) *mmath.MQuaternion {
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return rot
	}

	parentMat := getParentCancelableRotationMat(deltas, bd.Bone.ParentIndex)

	// 自身のキャンセルが無い場合
	if bd.FrameCancelableRotation == nil || bd.FrameCancelableRotation.IsIdent() {
		if parentMat == nil {
			return rot
		}
		return rot.ToMat4().Muled(parentMat.Inverted()).Quaternion()
	}

	// 自身のキャンセル適用
	newMat := rot.ToMat4().Muled(bd.FrameCancelableRotation.ToMat4())

	if parentMat == nil {
		return newMat.Quaternion()
	}
	return newMat.Muled(parentMat.Inverted()).Quaternion()
}

// 位置のキャンセル適用
func applyCancelablePosition(
	deltas *delta.BoneDeltas,
	boneIndex int,
	posMat *mmath.MMat4,
) *mmath.MMat4 {
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return posMat
	}

	parentMat := getParentCancelablePositionMat(deltas, bd.Bone.ParentIndex)

	hasSelfCancel := (bd.FrameCancelablePosition != nil && !bd.FrameCancelablePosition.IsZero()) ||
		(bd.FrameMorphCancelablePosition != nil && !bd.FrameMorphCancelablePosition.IsZero())

	if !hasSelfCancel {
		if parentMat == nil {
			return posMat
		}
		return posMat.Muled(parentMat.Inverted())
	}

	if bd.FrameCancelablePosition != nil && !bd.FrameCancelablePosition.IsZero() {
		posMat = posMat.Muled(bd.FrameCancelablePosition.ToMat4())
	}
	if bd.FrameMorphCancelablePosition != nil && !bd.FrameMorphCancelablePosition.IsZero() {
		posMat = posMat.Muled(bd.FrameMorphCancelablePosition.ToMat4())
	}

	if parentMat == nil {
		return posMat
	}
	return posMat.Muled(parentMat.Inverted())
}

// スケールのキャンセル適用
func applyCancelableScale(
	deltas *delta.BoneDeltas,
	boneIndex int,
	scaleMat *mmath.MMat4,
) *mmath.MMat4 {
	bd := deltas.Get(boneIndex)
	if bd == nil {
		return scaleMat
	}

	parentMat := getParentCancelableScaleMat(deltas, bd.Bone.ParentIndex)

	hasSelfCancel := (bd.FrameCancelableScale != nil && !bd.FrameCancelableScale.IsZero()) ||
		(bd.FrameMorphCancelableScale != nil && !bd.FrameMorphCancelableScale.IsZero())

	if !hasSelfCancel {
		if parentMat == nil {
			return scaleMat
		}
		return scaleMat.Muled(parentMat.Inverted())
	}

	// 自分のスケールキャンセル適用
	if bd.FrameCancelableScale != nil && !bd.FrameCancelableScale.IsZero() {
		scaleMat = scaleMat.Muled(bd.FrameCancelableScale.ToScaleMat4())
	}
	if bd.FrameMorphCancelableScale != nil && !bd.FrameMorphCancelableScale.IsZero() {
		scaleMat = scaleMat.Muled(bd.FrameMorphCancelableScale.ToScaleMat4())
	}

	if parentMat == nil {
		return scaleMat
	}
	return scaleMat.Muled(parentMat.Inverted())
}

// ----------------------------------------------------------------------------
// 以下、親のキャンセル行列を取得するためのプライベートヘルパー
// ----------------------------------------------------------------------------

func getParentCancelableRotationMat(
	deltas *delta.BoneDeltas,
	parentIndex int,
) *mmath.MMat4 {
	if !deltas.Contains(parentIndex) {
		return nil
	}
	pb := deltas.Get(parentIndex)
	var mat *mmath.MMat4

	if pb.FrameCancelableRotation != nil && !pb.FrameCancelableRotation.IsIdent() {
		mat = pb.FrameCancelableRotation.ToMat4()
	}
	if pb.FrameMorphCancelableRotation != nil && !pb.FrameMorphCancelableRotation.IsIdent() {
		if mat == nil {
			mat = pb.FrameMorphCancelableRotation.ToMat4()
		} else {
			mat = mat.Muled(pb.FrameMorphCancelableRotation.ToMat4())
		}
	}
	return mat
}

func getParentCancelablePositionMat(
	deltas *delta.BoneDeltas,
	parentIndex int,
) *mmath.MMat4 {
	if !deltas.Contains(parentIndex) {
		return nil
	}
	pb := deltas.Get(parentIndex)
	var mat *mmath.MMat4

	if pb.FrameCancelablePosition != nil && !pb.FrameCancelablePosition.IsZero() {
		mat = pb.FrameCancelablePosition.ToMat4()
	}
	if pb.FrameMorphCancelablePosition != nil && !pb.FrameMorphCancelablePosition.IsZero() {
		if mat == nil {
			mat = pb.FrameMorphCancelablePosition.ToMat4()
		} else {
			mat = mat.Muled(pb.FrameMorphCancelablePosition.ToMat4())
		}
	}
	return mat
}

func getParentCancelableScaleMat(
	deltas *delta.BoneDeltas,
	parentIndex int,
) *mmath.MMat4 {
	if !deltas.Contains(parentIndex) {
		return nil
	}
	pb := deltas.Get(parentIndex)
	var mat *mmath.MMat4

	if pb.FrameCancelableScale != nil && !pb.FrameCancelableScale.IsZero() {
		mat = pb.FrameCancelableScale.ToScaleMat4()
	}
	if pb.FrameMorphCancelableScale != nil && !pb.FrameMorphCancelableScale.IsZero() {
		if mat == nil {
			mat = pb.FrameMorphCancelableScale.ToScaleMat4()
		} else {
			mat = mat.Muled(pb.FrameMorphCancelableScale.ToScaleMat4())
		}
	}
	return mat
}
