package vmd

import "github.com/tiendc/go-deepcopy"

type WindLiftCoeffFrame struct {
	*BaseFrame         // キーフレ
	LiftCoeff  float32 // 風速
}

func NewWindLiftCoeffFrame(index float32) *WindLiftCoeffFrame {
	return &WindLiftCoeffFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		LiftCoeff: 1.0,
	}
}

func NewWindLiftCoeffFrameByValue(index float32, liftCoeff float32) *WindLiftCoeffFrame {
	return &WindLiftCoeffFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		LiftCoeff: liftCoeff,
	}
}

func (nextMf *WindLiftCoeffFrame) Copy() IBaseFrame {
	vv := &WindLiftCoeffFrame{
		BaseFrame: nextMf.BaseFrame.Copy().(*BaseFrame),
		LiftCoeff: nextMf.LiftCoeff,
	}
	return vv
}

func (nextMf *WindLiftCoeffFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevWF := prevFrame.(*WindLiftCoeffFrame)

	prevIndex := prevFrame.Index()
	nextIndex := nextMf.Index()

	wf := NewWindLiftCoeffFrame(index)

	ry := float32(index-prevIndex) / float32(nextIndex-prevIndex)
	wf.LiftCoeff = prevWF.LiftCoeff + (nextMf.LiftCoeff-prevWF.LiftCoeff)*ry

	return wf
}

func (nextMf *WindLiftCoeffFrame) WindLiftCoeff() float32 {
	if nextMf.LiftCoeff <= 0 {
		return 1.0
	}
	return float32(1.0 / nextMf.LiftCoeff)
}

func (nextMf *WindLiftCoeffFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type WindLiftCoeffFrames struct {
	*BaseFrames[*WindLiftCoeffFrame]
}

func NewWindLiftCoeffFrames() *WindLiftCoeffFrames {
	return &WindLiftCoeffFrames{
		BaseFrames: NewBaseFrames(NewWindLiftCoeffFrame, nilWindLiftCoeffFrame),
	}
}

func nilWindLiftCoeffFrame() *WindLiftCoeffFrame {
	return nil
}

func (mf *WindLiftCoeffFrames) Copy() (*WindLiftCoeffFrames, error) {
	copied := new(WindLiftCoeffFrames)
	err := deepcopy.Copy(copied, mf)
	return copied, err
}

func (mf *WindLiftCoeffFrames) Get(frame float32) *WindLiftCoeffFrame {
	if mf.values.Len() == 0 {
		return &WindLiftCoeffFrame{
			BaseFrame: NewFrame(frame).(*BaseFrame),
			LiftCoeff: 60, // デフォルト値
		}
	}

	if mf.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return mf.values.Get(NewFrame(frame)).(*WindLiftCoeffFrame)
	}

	nextFrame := mf.NextFrame(frame)
	if nextFrame > frame {
		// ドンピシャルフレームがなく、次のキーフレがある場合、次のキーフレを取得
		return mf.values.Get(NewFrame(nextFrame)).(*WindLiftCoeffFrame)
	}

	// 次のキーフレがない場合、前のキーフレを取得
	prevFrame := mf.PrevFrame(frame)
	return mf.values.Get(NewFrame(prevFrame)).(*WindLiftCoeffFrame)
}
