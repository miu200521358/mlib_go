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
