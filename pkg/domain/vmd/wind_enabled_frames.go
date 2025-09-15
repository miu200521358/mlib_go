package vmd

import "github.com/tiendc/go-deepcopy"

type WindEnabledFrame struct {
	*BaseFrame      // キーフレ
	Enabled    bool // 物理リセット種別
}

func NewWindEnabledFrame(index float32) *WindEnabledFrame {
	return &WindEnabledFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Enabled:   false, // デフォルト値
	}
}

func NewWindEnabledFrameByValue(index float32, physicsEnabled bool) *WindEnabledFrame {
	return &WindEnabledFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Enabled:   physicsEnabled,
	}
}

func (nextWf *WindEnabledFrame) Copy() IBaseFrame {
	vv := &WindEnabledFrame{
		BaseFrame: nextWf.BaseFrame.Copy().(*BaseFrame),
		Enabled:   nextWf.Enabled,
	}
	return vv
}

func (nextWf *WindEnabledFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevWf := prevFrame.(*WindEnabledFrame)
	// 補間なしで前のキーフレを引き継ぐ
	vv := &WindEnabledFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Enabled:   prevWf.Enabled,
	}
	return vv
}

func (nextWf *WindEnabledFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type WindEnabledFrames struct {
	*BaseFrames[*WindEnabledFrame]
}

func NewWindEnabledFrames() *WindEnabledFrames {
	return &WindEnabledFrames{
		BaseFrames: NewBaseFrames(NewWindEnabledFrame, nilWindEnabledFrame),
	}
}

func nilWindEnabledFrame() *WindEnabledFrame {
	return nil
}

func (mf *WindEnabledFrames) Copy() (*WindEnabledFrames, error) {
	copied := new(WindEnabledFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *WindEnabledFrames) Get(frame float32) *WindEnabledFrame {
	if mf.values.Len() == 0 {
		return &WindEnabledFrame{
			BaseFrame: NewFrame(frame).(*BaseFrame),
			Enabled:   false, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*WindEnabledFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*WindEnabledFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*WindEnabledFrame)
}
