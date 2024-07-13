package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type LightFrame struct {
	*BaseFrame              // キーフレ
	Position   *mmath.MVec3 // 位置
	Color      *mmath.MVec3 // 色
}

func NewLightFrame(index int) *LightFrame {
	return &LightFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Position:  mmath.NewMVec3(),
		Color:     mmath.NewMVec3(),
	}
}

func NullLightFrame() *LightFrame {
	return nil
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
	copied := NewLightFrame(lf.GetIndex())
	copied.Position = lf.Position
	copied.Color = lf.Color

	return copied
}

func (nextLf *LightFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevLf := prevFrame.(*LightFrame)
	// 線形補間
	t := float64(nextLf.GetIndex()-index) / float64(nextLf.GetIndex()-prevLf.GetIndex())
	vv := &LightFrame{
		Position: prevLf.Position.Lerp(nextLf.Position, t),
		Color:    prevLf.Color.Lerp(nextLf.Color, t),
	}
	return vv
}

func (lf *LightFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}
