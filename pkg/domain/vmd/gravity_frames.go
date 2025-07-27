package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type GravityFrame struct {
	*BaseFrame // キーフレ
	Gravity    *mmath.MVec3
}

func NewGravityFrame(index float32) *GravityFrame {
	return &GravityFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Gravity:   &mmath.MVec3{X: 0, Y: -9.8, Z: 0}, // デフォルト重力値
	}
}

func NewGravityFrameByValue(index float32, gravity *mmath.MVec3) *GravityFrame {
	return &GravityFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Gravity:   gravity,
	}
}

func (nextMf *GravityFrame) Copy() IBaseFrame {
	vv := &GravityFrame{
		BaseFrame: nextMf.BaseFrame.Copy().(*BaseFrame),
		Gravity:   nextMf.Gravity.Copy(),
	}
	return vv
}

func (nextMf *GravityFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*GravityFrame)
	// 補間なしで前のキーフレを引き継ぐ
	vv := &GravityFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Gravity:   prevMf.Gravity.Copy(),
	}
	return vv
}

func (nextMf *GravityFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type GravityFrames struct {
	*BaseFrames[*GravityFrame]
}

func NewGravityFrames() *GravityFrames {
	return &GravityFrames{
		BaseFrames: NewBaseFrames(NewGravityFrame, nilGravityFrame),
	}
}

func nilGravityFrame() *GravityFrame {
	return nil
}

func (mf *GravityFrames) Copy() (*GravityFrames, error) {
	copied := new(GravityFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *GravityFrames) Get(frame float32) *GravityFrame {
	if mf.values.Len() == 0 {
		return &GravityFrame{
			BaseFrame: NewFrame(frame).(*BaseFrame),
			Gravity:   &mmath.MVec3{X: 0, Y: -9.8, Z: 0}, // デフォルト重力値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*GravityFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*GravityFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*GravityFrame)
}
