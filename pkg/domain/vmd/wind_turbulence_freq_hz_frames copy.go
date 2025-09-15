package vmd

import "github.com/tiendc/go-deepcopy"

type WindTurbulenceFreqHzFrame struct {
	*BaseFrame               // キーフレ
	TurbulenceFreqHz float32 // 風速
}

func NewWindTurbulenceFreqHzFrame(index float32) *WindTurbulenceFreqHzFrame {
	return &WindTurbulenceFreqHzFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		TurbulenceFreqHz: 0.0,
	}
}

func NewWindTurbulenceFreqHzFrameByValue(index float32, turbulenceFreqHz float32) *WindTurbulenceFreqHzFrame {
	return &WindTurbulenceFreqHzFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		TurbulenceFreqHz: turbulenceFreqHz,
	}
}

func (nextMf *WindTurbulenceFreqHzFrame) Copy() IBaseFrame {
	vv := &WindTurbulenceFreqHzFrame{
		BaseFrame:        nextMf.BaseFrame.Copy().(*BaseFrame),
		TurbulenceFreqHz: nextMf.TurbulenceFreqHz,
	}
	return vv
}

func (nextMf *WindTurbulenceFreqHzFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevWF := prevFrame.(*WindTurbulenceFreqHzFrame)

	prevIndex := prevFrame.Index()
	nextIndex := nextMf.Index()

	wf := NewWindTurbulenceFreqHzFrame(index)

	ry := float32(index-prevIndex) / float32(nextIndex-prevIndex)
	wf.TurbulenceFreqHz = prevWF.TurbulenceFreqHz + (nextMf.TurbulenceFreqHz-prevWF.TurbulenceFreqHz)*ry

	return wf
}

func (nextMf *WindTurbulenceFreqHzFrame) WindTurbulenceFreqHz() float32 {
	if nextMf.TurbulenceFreqHz <= 0 {
		return 0.0
	}
	return float32(1.0 / nextMf.TurbulenceFreqHz)
}

func (nextMf *WindTurbulenceFreqHzFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type WindTurbulenceFreqHzFrames struct {
	*BaseFrames[*WindTurbulenceFreqHzFrame]
}

func NewWindTurbulenceFreqHzFrames() *WindTurbulenceFreqHzFrames {
	return &WindTurbulenceFreqHzFrames{
		BaseFrames: NewBaseFrames(NewWindTurbulenceFreqHzFrame, nilWindTurbulenceFreqHzFrame),
	}
}

func nilWindTurbulenceFreqHzFrame() *WindTurbulenceFreqHzFrame {
	return nil
}

func (mf *WindTurbulenceFreqHzFrames) Copy() (*WindTurbulenceFreqHzFrames, error) {
	copied := new(WindTurbulenceFreqHzFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *WindTurbulenceFreqHzFrames) Get(frame float32) *WindTurbulenceFreqHzFrame {
	if mf.values.Len() == 0 {
		return &WindTurbulenceFreqHzFrame{
			BaseFrame:        NewFrame(frame).(*BaseFrame),
			TurbulenceFreqHz: 0.0, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*WindTurbulenceFreqHzFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*WindTurbulenceFreqHzFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*WindTurbulenceFreqHzFrame)
}
