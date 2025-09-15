package vmd

import "github.com/tiendc/go-deepcopy"

type WindSpeedFrame struct {
	*BaseFrame         // キーフレ
	Speed      float32 // 風速
}

func NewWindSpeedFrame(index float32) *WindSpeedFrame {
	return &WindSpeedFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Speed:     1.0,
	}
}

func NewWindSpeedFrameByValue(index float32, speed float32) *WindSpeedFrame {
	return &WindSpeedFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Speed:     speed,
	}
}

func (nextMf *WindSpeedFrame) Copy() IBaseFrame {
	vv := &WindSpeedFrame{
		BaseFrame: nextMf.BaseFrame.Copy().(*BaseFrame),
		Speed:     nextMf.Speed,
	}
	return vv
}

func (nextMf *WindSpeedFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevWF := prevFrame.(*WindSpeedFrame)

	prevIndex := prevFrame.Index()
	nextIndex := nextMf.Index()

	wf := NewWindSpeedFrame(index)

	ry := float32(index-prevIndex) / float32(nextIndex-prevIndex)
	wf.Speed = prevWF.Speed + (nextMf.Speed-prevWF.Speed)*ry

	return wf
}

func (nextMf *WindSpeedFrame) WindSpeed() float32 {
	if nextMf.Speed <= 0 {
		return 1.0
	}
	return float32(1.0 / nextMf.Speed)
}

func (nextMf *WindSpeedFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type WindSpeedFrames struct {
	*BaseFrames[*WindSpeedFrame]
}

func NewWindSpeedFrames() *WindSpeedFrames {
	return &WindSpeedFrames{
		BaseFrames: NewBaseFrames(NewWindSpeedFrame, nilWindSpeedFrame),
	}
}

func nilWindSpeedFrame() *WindSpeedFrame {
	return nil
}

func (mf *WindSpeedFrames) Copy() (*WindSpeedFrames, error) {
	copied := new(WindSpeedFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *WindSpeedFrames) Get(frame float32) *WindSpeedFrame {
	if mf.values.Len() == 0 {
		return &WindSpeedFrame{
			BaseFrame: NewFrame(frame).(*BaseFrame),
			Speed:     60, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*WindSpeedFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*WindSpeedFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*WindSpeedFrame)
}
