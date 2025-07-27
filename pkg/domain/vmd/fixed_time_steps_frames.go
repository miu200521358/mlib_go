package vmd

import "github.com/tiendc/go-deepcopy"

type FixedTimeStepFrame struct {
	*BaseFrame               // キーフレ
	FixedTimeStepNum float64 // 演算頻度
}

func NewFixedTimeStepFrame(index float32) *FixedTimeStepFrame {
	return &FixedTimeStepFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		FixedTimeStepNum: 60, // デフォルト値
	}
}

func NewFixedTimeStepFrameByValue(index float32, fixedTimeStep float64) *FixedTimeStepFrame {
	return &FixedTimeStepFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		FixedTimeStepNum: fixedTimeStep,
	}
}

func (nextMf *FixedTimeStepFrame) Copy() IBaseFrame {
	vv := &FixedTimeStepFrame{
		BaseFrame:        nextMf.BaseFrame.Copy().(*BaseFrame),
		FixedTimeStepNum: nextMf.FixedTimeStepNum,
	}
	return vv
}

func (nextMf *FixedTimeStepFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*FixedTimeStepFrame)
	// 補間なしで前のキーフレを引き継ぐ
	vv := &FixedTimeStepFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		FixedTimeStepNum: prevMf.FixedTimeStepNum,
	}
	return vv
}

func (nextMf *FixedTimeStepFrame) FixedTimeStep() float32 {
	if nextMf.FixedTimeStepNum <= 0 {
		return 1.0 / 60.0 // デフォルト値
	}
	return float32(1.0 / nextMf.FixedTimeStepNum)
}

func (nextMf *FixedTimeStepFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type FixedTimeStepFrames struct {
	*BaseFrames[*FixedTimeStepFrame]
}

func NewFixedTimeStepFrames() *FixedTimeStepFrames {
	return &FixedTimeStepFrames{
		BaseFrames: NewBaseFrames(NewFixedTimeStepFrame, nilFixedTimeStepFrame),
	}
}

func nilFixedTimeStepFrame() *FixedTimeStepFrame {
	return nil
}

func (mf *FixedTimeStepFrames) Copy() (*FixedTimeStepFrames, error) {
	copied := new(FixedTimeStepFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *FixedTimeStepFrames) Get(frame float32) *FixedTimeStepFrame {
	if mf.values.Len() == 0 {
		return &FixedTimeStepFrame{
			BaseFrame:        NewFrame(frame).(*BaseFrame),
			FixedTimeStepNum: 60, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*FixedTimeStepFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*FixedTimeStepFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*FixedTimeStepFrame)
}
