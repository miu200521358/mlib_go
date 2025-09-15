package vmd

import "github.com/tiendc/go-deepcopy"

type WindRandomnessFrame struct {
	*BaseFrame         // キーフレ
	Randomness float32 // 風速
}

func NewWindRandomnessFrame(index float32) *WindRandomnessFrame {
	return &WindRandomnessFrame{
		BaseFrame:  NewFrame(index).(*BaseFrame),
		Randomness: 0.0,
	}
}

func NewWindRandomnessFrameByValue(index float32, randomness float32) *WindRandomnessFrame {
	return &WindRandomnessFrame{
		BaseFrame:  NewFrame(index).(*BaseFrame),
		Randomness: randomness,
	}
}

func (nextMf *WindRandomnessFrame) Copy() IBaseFrame {
	vv := &WindRandomnessFrame{
		BaseFrame:  nextMf.BaseFrame.Copy().(*BaseFrame),
		Randomness: nextMf.Randomness,
	}
	return vv
}

func (nextMf *WindRandomnessFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevWF := prevFrame.(*WindRandomnessFrame)

	prevIndex := prevFrame.Index()
	nextIndex := nextMf.Index()

	wf := NewWindRandomnessFrame(index)

	ry := float32(index-prevIndex) / float32(nextIndex-prevIndex)
	wf.Randomness = prevWF.Randomness + (nextMf.Randomness-prevWF.Randomness)*ry

	return wf
}

func (nextMf *WindRandomnessFrame) WindRandomness() float32 {
	if nextMf.Randomness <= 0 {
		return 0.0
	}
	return float32(1.0 / nextMf.Randomness)
}

func (nextMf *WindRandomnessFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type WindRandomnessFrames struct {
	*BaseFrames[*WindRandomnessFrame]
}

func NewWindRandomnessFrames() *WindRandomnessFrames {
	return &WindRandomnessFrames{
		BaseFrames: NewBaseFrames(NewWindRandomnessFrame, nilWindRandomnessFrame),
	}
}

func nilWindRandomnessFrame() *WindRandomnessFrame {
	return nil
}

func (mf *WindRandomnessFrames) Copy() (*WindRandomnessFrames, error) {
	copied := new(WindRandomnessFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *WindRandomnessFrames) Get(frame float32) *WindRandomnessFrame {
	if mf.values.Len() == 0 {
		return &WindRandomnessFrame{
			BaseFrame:  NewFrame(frame).(*BaseFrame),
			Randomness: 0.0, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*WindRandomnessFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*WindRandomnessFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*WindRandomnessFrame)
}
