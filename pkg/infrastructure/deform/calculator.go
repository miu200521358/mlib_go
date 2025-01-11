package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TotalRotationMat(boneDeltas *delta.BoneDeltas, boneIndex int) *mmath.MMat4 {
	rotMat := totalRotationLoop(boneDeltas, boneIndex, 0, 1.0).ToMat4()
	return totalCancelRotation(boneDeltas, boneIndex, rotMat)
}

func totalCancelRotation(boneDeltas *delta.BoneDeltas, boneIndex int, rotMat *mmath.MMat4) *mmath.MMat4 {
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

func totalRotationLoop(boneDeltas *delta.BoneDeltas, boneIndex int, loop int, factor float64) *mmath.MQuaternion {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}

	rot := boneDelta.FilledTotalRotation()

	if boneDelta.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := totalRotationLoop(boneDeltas, boneDelta.Bone.EffectIndex, loop+1, boneDelta.Bone.EffectFactor)
		rot.Mul(effectorRot)
	}

	return rot.MuledScalar(factor)
}

// 該当ボーンまでの付与親を加味した全ての回転（モーフは含まない）
func TotalBoneRotation(boneDeltas *delta.BoneDeltas, boneIndex int) *mmath.MQuaternion {
	rot := totalBoneRotationLoop(boneDeltas, boneIndex, 0, 1.0)
	return totalBoneCancelRotation(boneDeltas, boneIndex, rot)
}

func totalBoneCancelRotation(boneDeltas *delta.BoneDeltas, boneIndex int, rot *mmath.MQuaternion) *mmath.MQuaternion {
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

func totalBoneRotationLoop(boneDeltas *delta.BoneDeltas, boneIndex int, loop int, factor float64) *mmath.MQuaternion {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMQuaternion()
	}

	rot := boneDelta.FilledFrameRotation().Copy()

	if boneDelta.Bone.IsEffectorRotation() {
		// 付与親回転がある場合、再帰で回転を取得する
		effectorRot := totalBoneRotationLoop(boneDeltas, boneDelta.Bone.EffectIndex, loop+1, boneDelta.Bone.EffectFactor)
		rot.Mul(effectorRot)
	}

	rot = rot.MuledScalar(factor)

	if boneDelta.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(boneDelta.Bone.NormalizedFixedAxis)
	}

	return rot
}

func TotalPositionMat(boneDeltas *delta.BoneDeltas, boneIndex int) *mmath.MMat4 {
	posMat := totalPositionLoop(boneDeltas, boneIndex, 0).ToMat4()
	return totalCancelPosition(boneDeltas, boneIndex, posMat)
}

func totalCancelPosition(boneDeltas *delta.BoneDeltas, boneIndex int, posMat *mmath.MMat4) *mmath.MMat4 {
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

func totalPositionLoop(boneDeltas *delta.BoneDeltas, boneIndex int, loop int) *mmath.MVec3 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return mmath.NewMVec3()
	}
	pos := boneDelta.FilledTotalPosition()

	if boneDelta.Bone.IsEffectorTranslation() {
		// 付与親移動がある場合、再帰で回転を取得する
		effectorPos := totalPositionLoop(boneDeltas, boneDelta.Bone.EffectIndex, loop+1)
		pos.Mul(effectorPos.MuledScalar(boneDelta.Bone.EffectFactor))
	}

	return pos
}

func TotalScaleMat(boneDeltas *delta.BoneDeltas, boneIndex int) *mmath.MMat4 {
	scaleMat := totalScaleMatLoop(boneDeltas, boneIndex, 0).ToScaleMat4()
	return totalCancelScale(boneDeltas, boneIndex, scaleMat)
}

func totalScaleMatLoop(boneDeltas *delta.BoneDeltas, boneIndex int, loop int) *mmath.MVec3 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil || loop > 10 {
		return &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	scale := boneDelta.FilledTotalScale()

	return scale
}

func totalCancelScale(boneDeltas *delta.BoneDeltas, boneIndex int, scaleMat *mmath.MMat4) *mmath.MMat4 {
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

func TotalLocalMat(boneDeltas *delta.BoneDeltas, boneIndex int) *mmath.MMat4 {
	boneDelta := boneDeltas.Get(boneIndex)
	if boneDelta == nil {
		return mmath.NewMMat4()
	}

	// ローカル変換行列
	return boneDelta.FilledTotalLocalMat()
}
