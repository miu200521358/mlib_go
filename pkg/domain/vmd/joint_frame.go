package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type JointFrame struct {
	*BaseFrame                             // キーフレ
	TranslationLimitMin       *mmath.MVec3 // 移動制限-下限(x,y,z)
	TranslationLimitMax       *mmath.MVec3 // 移動制限-上限(x,y,z)
	RotationLimitMin          *mmath.MVec3 // 回転制限-下限
	RotationLimitMax          *mmath.MVec3 // 回転制限-上限
	SpringConstantTranslation *mmath.MVec3 // バネ定数-移動(x,y,z)
	SpringConstantRotation    *mmath.MVec3 // バネ定数-回転(x,y,z)
}

func NewJointFrameByValues(index float32, translationLimitMin, translationLimitMax, rotationLimitMin, rotationLimitMax, springConstantTranslation, springConstantRotation *mmath.MVec3) *JointFrame {
	return &JointFrame{
		BaseFrame:                 NewFrame(index).(*BaseFrame),
		TranslationLimitMin:       translationLimitMin,
		TranslationLimitMax:       translationLimitMax,
		RotationLimitMin:          rotationLimitMin,
		RotationLimitMax:          rotationLimitMax,
		SpringConstantTranslation: springConstantTranslation,
		SpringConstantRotation:    springConstantRotation,
	}
}

func (mf *JointFrame) Copy() IBaseFrame {
	copied := new(JointFrame)
	deepcopy.Copy(mf, copied)
	return copied
}

func (nextMf *JointFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*JointFrame)

	prevIndex := prevMf.Index()
	nextIndex := nextMf.Index()

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	translationLimitMin := prevMf.TranslationLimitMin.Lerp(nextMf.TranslationLimitMin, ry)
	translationLimitMax := prevMf.TranslationLimitMax.Lerp(nextMf.TranslationLimitMax, ry)
	rotationLimitMin := prevMf.RotationLimitMin.Lerp(nextMf.RotationLimitMin, ry)
	rotationLimitMax := prevMf.RotationLimitMax.Lerp(nextMf.RotationLimitMax, ry)
	springConstantTranslation := prevMf.SpringConstantTranslation.Lerp(nextMf.SpringConstantTranslation, ry)
	springConstantRotation := prevMf.SpringConstantRotation.Lerp(nextMf.SpringConstantRotation, ry)

	return NewJointFrameByValues(index, translationLimitMin, translationLimitMax, rotationLimitMin, rotationLimitMax, springConstantTranslation, springConstantRotation)
}

func (mf *JointFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}
