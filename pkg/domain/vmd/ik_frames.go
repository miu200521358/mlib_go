package vmd

type IkFrames struct {
	*BaseFrames[*IkFrame]
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		BaseFrames: NewBaseFrames[*IkFrame](NewIkFrame, NullNewIkFrame),
	}
}

func (ikFrames *IkFrames) Copy() *IkFrames {
	copied := NewIkFrames()
	for _, frame := range ikFrames.List() {
		copied.Append(frame.Copy().(*IkFrame))
	}
	return copied
}

func (ikFrames *IkFrames) Clean() {
	if ikFrames.Len() > 1 {
		return
	} else {
		cf := ikFrames.Get(ikFrames.Indexes.Min())
		if !(cf.IkList == nil || len(cf.IkList) == 0) {
			return
		}
		ikFrames.Delete(cf.Index())
	}
}
