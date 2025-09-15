package vmd

import "github.com/tiendc/go-deepcopy"

type WindDragCoeffFrame struct {
	*BaseFrame         // キーフレ
	DragCoeff  float32 // 風速
}

func NewWindDragCoeffFrame(index float32) *WindDragCoeffFrame {
	return &WindDragCoeffFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		DragCoeff: 1.0,
	}
}

func NewWindDragCoeffFrameByValue(index float32, dragCoeff float32) *WindDragCoeffFrame {
	return &WindDragCoeffFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		DragCoeff: dragCoeff,
	}
}

func (nextMf *WindDragCoeffFrame) Copy() IBaseFrame {
	vv := &WindDragCoeffFrame{
		BaseFrame: nextMf.BaseFrame.Copy().(*BaseFrame),
		DragCoeff: nextMf.DragCoeff,
	}
	return vv
}

func (nextMf *WindDragCoeffFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevWF := prevFrame.(*WindDragCoeffFrame)

	prevIndex := prevFrame.Index()
	nextIndex := nextMf.Index()

	wf := NewWindDragCoeffFrame(index)

	ry := float32(index-prevIndex) / float32(nextIndex-prevIndex)
	wf.DragCoeff = prevWF.DragCoeff + (nextMf.DragCoeff-prevWF.DragCoeff)*ry

	return wf
}

func (nextMf *WindDragCoeffFrame) WindDragCoeff() float32 {
	if nextMf.DragCoeff <= 0 {
		return 1.0
	}
	return float32(1.0 / nextMf.DragCoeff)
}

func (nextMf *WindDragCoeffFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type WindDragCoeffFrames struct {
	*BaseFrames[*WindDragCoeffFrame]
}

func NewWindDragCoeffFrames() *WindDragCoeffFrames {
	return &WindDragCoeffFrames{
		BaseFrames: NewBaseFrames(NewWindDragCoeffFrame, nilWindDragCoeffFrame),
	}
}

func nilWindDragCoeffFrame() *WindDragCoeffFrame {
	return nil
}

func (mf *WindDragCoeffFrames) Copy() (*WindDragCoeffFrames, error) {
	copied := new(WindDragCoeffFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *WindDragCoeffFrames) Get(frame float32) *WindDragCoeffFrame {
	if mf.values.Len() == 0 {
		return &WindDragCoeffFrame{
			BaseFrame: NewFrame(frame).(*BaseFrame),
			DragCoeff: 60, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*WindDragCoeffFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*WindDragCoeffFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*WindDragCoeffFrame)
}
