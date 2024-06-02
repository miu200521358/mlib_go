package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type ShadowFrame struct {
	*BaseFrame         // キーフレ
	ShadowMode int     // セルフ影モード
	Distance   float64 // 影範囲距離
}

func NewShadowFrame(index int) *ShadowFrame {
	return &ShadowFrame{
		BaseFrame:  NewFrame(index).(*BaseFrame),
		ShadowMode: 0,
		Distance:   0.0,
	}
}

func NullShadowFrame() *ShadowFrame {
	return nil
}

func (sf *ShadowFrame) Copy() IBaseFrame {
	vv := &ShadowFrame{
		ShadowMode: sf.ShadowMode,
		Distance:   sf.Distance,
	}
	return vv
}

func (nextSf *ShadowFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevSf := prevFrame.(*ShadowFrame)

	prevIndex := prevSf.GetIndex()
	nextIndex := nextSf.GetIndex()

	sf := NewShadowFrame(index)

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	sf.ShadowMode = prevSf.ShadowMode
	sf.Distance = mmath.LerpFloat(prevSf.Distance, nextSf.Distance, ry)

	return sf
}

func (sf *ShadowFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}
