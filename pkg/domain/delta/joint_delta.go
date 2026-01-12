// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	sharedtime "github.com/miu200521358/mlib_go/pkg/shared/contracts/time"
)

// JointDelta はジョイント差分を表す。
type JointDelta struct {
	Joint                     *model.Joint
	Frame                     sharedtime.Frame
	TranslationLimitMin       mmath.Vec3
	TranslationLimitMax       mmath.Vec3
	RotationLimitMin          mmath.Vec3
	RotationLimitMax          mmath.Vec3
	SpringConstantTranslation mmath.Vec3
	SpringConstantRotation    mmath.Vec3
}

// NewJointDelta はJointDeltaを生成する。
func NewJointDelta(joint *model.Joint, frame sharedtime.Frame) *JointDelta {
	if joint == nil {
		return nil
	}
	param := joint.Param
	return &JointDelta{
		Joint:                     joint,
		Frame:                     frame,
		TranslationLimitMin:       param.TranslationLimitMin,
		TranslationLimitMax:       param.TranslationLimitMax,
		RotationLimitMin:          param.RotationLimitMin,
		RotationLimitMax:          param.RotationLimitMax,
		SpringConstantTranslation: param.SpringConstantTranslation,
		SpringConstantRotation:    param.SpringConstantRotation,
	}
}

// NewJointDeltaByValue は値を指定してJointDeltaを生成する。
func NewJointDeltaByValue(
	joint *model.Joint,
	frame sharedtime.Frame,
	translationMin, translationMax mmath.Vec3,
	rotationMin, rotationMax mmath.Vec3,
	springTrans, springRot mmath.Vec3,
) *JointDelta {
	if joint == nil {
		return nil
	}
	return &JointDelta{
		Joint:                     joint,
		Frame:                     frame,
		TranslationLimitMin:       translationMin,
		TranslationLimitMax:       translationMax,
		RotationLimitMin:          rotationMin,
		RotationLimitMax:          rotationMax,
		SpringConstantTranslation: springTrans,
		SpringConstantRotation:    springRot,
	}
}
