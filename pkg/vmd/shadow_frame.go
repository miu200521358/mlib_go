package vmd

import "github.com/miu200521358/mlib_go/pkg/mcore"

type ShadowFrame struct {
	*BaseFrame         // キーフレ
	ShadowMode int     // セルフ影モード
	Distance   float64 // 影範囲距離
}

func NewShadowFrame(index float32) *ShadowFrame {
	return &ShadowFrame{
		BaseFrame:  NewVmdBaseFrame(index),
		ShadowMode: 0,
		Distance:   0.0,
	}
}

func (sf *ShadowFrame) Copy() mcore.IndexFloatModelInterface {
	vv := &ShadowFrame{
		ShadowMode: sf.ShadowMode,
		Distance:   sf.Distance,
	}
	return vv
}
