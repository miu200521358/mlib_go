package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

// JointDelta は1つのボーンにおける変形（ポジション・回転・スケールなど）の差分を表す
type JointDelta struct {
	Joint                     *pmx.Joint
	Frame                     float32
	TranslationLimitMin       *mmath.MVec3 // 移動制限-下限(x,y,z)
	TranslationLimitMax       *mmath.MVec3 // 移動制限-上限(x,y,z)
	RotationLimitMin          *mmath.MVec3 // 回転制限-下限
	RotationLimitMax          *mmath.MVec3 // 回転制限-上限
	SpringConstantTranslation *mmath.MVec3 // バネ定数-移動(x,y,z)
	SpringConstantRotation    *mmath.MVec3 // バネ定数-回転(x,y,z)
}

// NewJointDelta は新規の JointDelta を生成するコンストラクタ
func NewJointDelta(bone *pmx.Joint, frame float32) *JointDelta {
	return &JointDelta{
		Joint: bone,
		Frame: frame,
	}
}

// NewJointDelta は新規の JointDelta を生成するコンストラクタ
func NewJointDeltaByValue(
	bone *pmx.Joint, frame float32,
	translationLimitMin, translationLimitMax *mmath.MVec3,
	rotationLimitMin, rotationLimitMax *mmath.MVec3,
	springConstantTranslation, springConstantRotation *mmath.MVec3,
) *JointDelta {
	return &JointDelta{
		Joint:                     bone,
		Frame:                     frame,
		TranslationLimitMin:       translationLimitMin,
		TranslationLimitMax:       translationLimitMax,
		RotationLimitMin:          rotationLimitMin,
		RotationLimitMax:          rotationLimitMax,
		SpringConstantTranslation: springConstantTranslation,
		SpringConstantRotation:    springConstantRotation,
	}
}
