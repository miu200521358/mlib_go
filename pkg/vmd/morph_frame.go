package vmd

type MorphFrame struct {
	*BaseFrame         // キーフレ
	Ratio      float64 // モーフの割合
}

func NewMorphFrame(index int) *MorphFrame {
	return &MorphFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Ratio:     0.0,
	}
}

func (mf *MorphFrame) Add(v *MorphFrame) {
	mf.Ratio += v.Ratio
}

func (mf *MorphFrame) Added(v *MorphFrame) *MorphFrame {
	copied := mf.Copy().(*MorphFrame)
	copied.Ratio += v.Ratio
	return copied
}

func (mf *MorphFrame) Copy() IBaseFrame {
	copied := NewMorphFrame(mf.GetIndex())
	copied.Ratio = mf.Ratio
	return copied
}

func (nextMf *MorphFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevMf := prevFrame.(*MorphFrame)

	prevIndex := prevMf.GetIndex()
	nextIndex := nextMf.GetIndex()

	mf := NewMorphFrame(index)

	ry := (index - prevIndex) / (nextIndex - prevIndex)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*float64(ry)

	return mf
}

func (mf *MorphFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}
