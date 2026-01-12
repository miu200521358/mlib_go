// 指示: miu200521358
package motion

import sharedtime "github.com/miu200521358/mlib_go/pkg/shared/contracts/time"

// Frame はモーションのフレーム番号を表す。
type Frame = sharedtime.Frame

// IBaseFrame はフレームの共通インターフェース。
type IBaseFrame interface {
	Index() Frame
	Copy() (IBaseFrame, error)
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
func (f *BaseFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*BaseFrame)(nil), nil
	}
	return &BaseFrame{index: f.index, Read: f.Read}, nil
}
