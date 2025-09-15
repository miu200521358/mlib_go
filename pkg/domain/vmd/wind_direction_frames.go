package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type WindDirectionFrame struct {
	*BaseFrame // キーフレ
	Direction  *mmath.MVec3
}

func NewWindDirectionFrame(index float32) *WindDirectionFrame {
	return &WindDirectionFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Direction: &mmath.MVec3{X: 0, Y: 0, Z: 0},
	}
}

func NewWindDirectionFrameByValue(index float32, direction *mmath.MVec3) *WindDirectionFrame {
	return &WindDirectionFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Direction: direction,
	}
}

func (nextWf *WindDirectionFrame) Copy() IBaseFrame {
	vv := &WindDirectionFrame{
		BaseFrame: nextWf.BaseFrame.Copy().(*BaseFrame),
		Direction: nextWf.Direction.Copy(),
	}
	return vv
}

func (nextWf *WindDirectionFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevWF := prevFrame.(*WindDirectionFrame)

	prevIndex := prevFrame.Index()
	nextIndex := nextWf.Index()

	wf := NewWindDirectionFrame(index)

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	wf.Direction = prevWF.Direction.Added((nextWf.Direction.Subed(prevWF.Direction)).MuledScalar(ry)).Effective()

	return wf
}

func (nextWf *WindDirectionFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type WindDirectionFrames struct {
	*BaseFrames[*WindDirectionFrame]
}

func NewWindDirectionFrames() *WindDirectionFrames {
	return &WindDirectionFrames{
		BaseFrames: NewBaseFrames(NewWindDirectionFrame, nilWindDirectionFrame),
	}
}

func nilWindDirectionFrame() *WindDirectionFrame {
	return nil
}

func (mf *WindDirectionFrames) Copy() (*WindDirectionFrames, error) {
	copied := new(WindDirectionFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *WindDirectionFrames) Get(frame float32) *WindDirectionFrame {
	if mf.values.Len() == 0 {
		return &WindDirectionFrame{
			BaseFrame: NewFrame(frame).(*BaseFrame),
			Direction: &mmath.MVec3{X: 0, Y: -9.8, Z: 0}, // デフォルト重力値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*WindDirectionFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*WindDirectionFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*WindDirectionFrame)
}
