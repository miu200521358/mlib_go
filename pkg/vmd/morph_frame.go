package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
)

type MorphFrame struct {
	*BaseFrame         // キーフレ
	Ratio      float64 // モーフの割合
}

func NewMorphFrame(index float32) *MorphFrame {
	return &MorphFrame{
		BaseFrame: NewVmdBaseFrame(index),
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

func (mf *MorphFrame) Copy() mcore.IIndexFloatModel {
	return &MorphFrame{
		BaseFrame: NewVmdBaseFrame(mf.Index),
		Ratio:     mf.Ratio,
	}
}
