package deform

import "github.com/miu200521358/mlib_go/pkg/mcore"

type BaseFrame struct {
	*mcore.IndexFloatModel
	Registered bool // 登録対象のキーフレであるか
	Read       bool // VMDファイルから読み込んだキーフレであるか
}

func NewVmdBaseFrame(index float32) *BaseFrame {
	return &BaseFrame{
		IndexFloatModel: &mcore.IndexFloatModel{Index: index},
		Registered:      false,
		Read:            false,
	}
}