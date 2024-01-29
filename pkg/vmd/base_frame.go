package vmd

import "github.com/miu200521358/mlib_go/pkg/mcore"

type BaseFrame struct {
	*mcore.IndexModel
	Registered bool // 登録対象のキーフレであるか
	Read       bool // VMDファイルから読み込んだキーフレであるか
}

func NewVmdBaseFrame(index int) *BaseFrame {
	return &BaseFrame{
		IndexModel: &mcore.IndexModel{Index: index},
		Registered: false,
		Read:       false,
	}
}
