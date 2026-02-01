// 指示: miu200521358
package motion

import "github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"

// Frame はモーションのフレーム番号を表す。
type Frame = mtime.Frame

// IBaseFrame はフレームの共通インターフェース。
type IBaseFrame interface {
	Index() Frame
}

// BaseFrame は共通フレーム情報を表す。
type BaseFrame struct {
	index Frame
	Read  bool
}

// NewBaseFrame はBaseFrameを生成する。
func NewBaseFrame(index Frame) *BaseFrame {
	return &BaseFrame{index: index}
}

// Index はフレーム番号を返す。
func (f *BaseFrame) Index() Frame {
	if f == nil {
		return 0
	}
	return f.index
}

// Copy はフレームを複製する。
func (f *BaseFrame) Copy() (BaseFrame, error) {
	if f == nil {
		return BaseFrame{}, nil
	}
	return BaseFrame{index: f.index, Read: f.Read}, nil
}
