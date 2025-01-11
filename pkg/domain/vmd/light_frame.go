package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type LightFrame struct {
	*BaseFrame              // キーフレ
	Position   *mmath.MVec3 // 位置
	Color      *mmath.MVec3 // 色
}

func NewLightFrame(index float32) *LightFrame {
	return &LightFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
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

func (lf *LightFrame) Copy() IBaseFrame {
	copied := NewLightFrame(lf.Index())
	copied.Position = lf.Position
	copied.Color = lf.Color

	return copied
}

func (nextLf *LightFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevLf := prevFrame.(*LightFrame)
	// 線形補間
	t := float64(nextLf.Index()-index) / float64(nextLf.Index()-prevLf.Index())
	vv := &LightFrame{
		Position: prevLf.Position.Lerp(nextLf.Position, t),
		Color:    prevLf.Color.Lerp(nextLf.Color, t),
	}
	return vv
}

func (lf *LightFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}
