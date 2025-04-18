package vmd

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

type MorphFrame struct {
	*BaseFrame         // キーフレ
	Ratio      float64 // モーフの割合
}

func NewMorphFrame(index float32) *MorphFrame {
	return &MorphFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Ratio:     0.0,
	}
}

func (mf *MorphFrame) Add(v *MorphFrame) {
	mf.Ratio += v.Ratio
}

func (mf *MorphFrame) Added(v *MorphFrame) *MorphFrame {
	copied := mf.Copy().(*MorphFrame)
	copied.Ratio += v.Ratio
	return copied
}

func (mf *MorphFrame) Copy() IBaseFrame {
	copied := NewMorphFrame(mf.Index())
	copied.Ratio = mf.Ratio
	return copied
}

func (nextMf *MorphFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*MorphFrame)

	prevIndex := prevMf.Index()
	nextIndex := nextMf.Index()

	mf := NewMorphFrame(index)

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	mf.Ratio = mmath.Effective(prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*ry)

	return mf
}

func (mf *MorphFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}
