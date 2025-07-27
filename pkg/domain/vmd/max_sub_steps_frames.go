package vmd

import "github.com/tiendc/go-deepcopy"

type MaxSubStepsFrame struct {
	*BaseFrame      // キーフレ
	MaxSubSteps int // 最大演算回数
}

func NewMaxSubStepsFrame(index float32) *MaxSubStepsFrame {
	return &MaxSubStepsFrame{
		BaseFrame:   NewFrame(index).(*BaseFrame),
		MaxSubSteps: 2, // デフォルト値
	}
}

func NewMaxSubStepsFrameByValue(index float32, maxSubStep int) *MaxSubStepsFrame {
	return &MaxSubStepsFrame{
		BaseFrame:   NewFrame(index).(*BaseFrame),
		MaxSubSteps: maxSubStep,
	}
}

func (nextMf *MaxSubStepsFrame) Copy() IBaseFrame {
	vv := &MaxSubStepsFrame{
		BaseFrame:   nextMf.BaseFrame.Copy().(*BaseFrame),
		MaxSubSteps: nextMf.MaxSubSteps,
	}
	return vv
}

func (nextMf *MaxSubStepsFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*MaxSubStepsFrame)
	// 補間なしで前のキーフレを引き継ぐ
	vv := &MaxSubStepsFrame{
		BaseFrame:   NewFrame(index).(*BaseFrame),
		MaxSubSteps: prevMf.MaxSubSteps,
	}
	return vv
}

func (nextMf *MaxSubStepsFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type MaxSubStepsFrames struct {
	*BaseFrames[*MaxSubStepsFrame]
}

func NewMaxSubStepsFrames() *MaxSubStepsFrames {
	return &MaxSubStepsFrames{
		BaseFrames: NewBaseFrames(NewMaxSubStepsFrame, nilMaxSubStepsFrame),
	}
}

func nilMaxSubStepsFrame() *MaxSubStepsFrame {
	return nil
}

func (mf *MaxSubStepsFrames) Copy() (*MaxSubStepsFrames, error) {
	copied := new(MaxSubStepsFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *MaxSubStepsFrames) Get(frame float32) *MaxSubStepsFrame {
	if mf.values.Len() == 0 {
		return &MaxSubStepsFrame{
			BaseFrame:   NewFrame(frame).(*BaseFrame),
			MaxSubSteps: 2, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*MaxSubStepsFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*MaxSubStepsFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*MaxSubStepsFrame)
}
