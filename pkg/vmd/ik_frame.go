package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/vmd/delta"
)

type IkEnabledFrame struct {
	*delta.BaseFrame        // キーフレ
	BoneName         string // IKボーン名
	Enabled          bool   // IKON/OFF
}

func NewIkEnableFrame(index float32) *IkEnabledFrame {
	return &IkEnabledFrame{
		BaseFrame: delta.NewVmdBaseFrame(index),
		BoneName:  "",
		Enabled:   true,
	}
}

func (kf *IkEnabledFrame) Copy() *IkEnabledFrame {
	vv := &IkEnabledFrame{
		BoneName: kf.BoneName,
		Enabled:  kf.Enabled,
	}
	return vv
}

type IkFrame struct {
	*delta.BaseFrame                   // キーフレ
	Visible          bool              // 表示ON/OFF
	IkList           []*IkEnabledFrame // IKリスト
}

func NewIkFrame(index float32) *IkFrame {
	return &IkFrame{
		BaseFrame: delta.NewVmdBaseFrame(index),
		Visible:   true,
		IkList:    []*IkEnabledFrame{},
	}
}

func (ikf *IkFrame) Copy() mcore.IndexFloatModelInterface {
	vv := &IkFrame{
		Visible: ikf.Visible,
		IkList:  []*IkEnabledFrame{},
	}
	for _, v := range ikf.IkList {
		vv.IkList = append(vv.IkList, v.Copy())
	}
	return vv
}
