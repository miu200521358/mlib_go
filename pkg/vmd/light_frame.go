package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type LightFrame struct {
	*BaseFrame              // キーフレ
	Position   *mmath.MVec3 // 位置
	Color      *mmath.MVec3 // 色
}

func NewLightFrame(index float32) *LightFrame {
	return &LightFrame{
		BaseFrame: NewVmdBaseFrame(index),
		Position:  mmath.NewMVec3(),
		Color:     mmath.NewMVec3(),
	}
}

func (lf *LightFrame) Add(v *LightFrame) {
	lf.Position.Add(v.Position)
	lf.Color.Add(v.Color)
}

func (lf *LightFrame) Added(v *LightFrame) *LightFrame {
	copied := lf.Copy().(*LightFrame)

	copied.Position.Add(v.Position)
	copied.Color.Add(v.Color)

	return copied
}
