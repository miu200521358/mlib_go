package vmd

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/mcore"
)

// VMDリーダー
type VmdVpdMotionReader struct {
	mcore.BaseReader[*VmdMotion]
	vmdReader *VmdMotionReader
	vpdReader *VpdMotionReader
}

func NewVmdVpdMotionReader() *VmdVpdMotionReader {
	reader := new(VmdVpdMotionReader)
	reader.vmdReader = &VmdMotionReader{}
	reader.vpdReader = &VpdMotionReader{}
	return reader
}

// 指定されたパスのファイルからデータを読み込む
func (r *VmdVpdMotionReader) ReadByFilepath(path string) (mcore.IHashModel, error) {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return r.vpdReader.ReadByFilepath(path)
	} else {
		return r.vmdReader.ReadByFilepath(path)
	}
}

func (r *VmdVpdMotionReader) ReadNameByFilepath(path string) (string, error) {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return r.vpdReader.ReadNameByFilepath(path)
	} else {
		return r.vmdReader.ReadNameByFilepath(path)
	}
}
