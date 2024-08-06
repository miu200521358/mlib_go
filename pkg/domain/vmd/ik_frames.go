package vmd

type IkFrames struct {
	*BaseFrames[*IkFrame]
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		BaseFrames: NewBaseFrames[*IkFrame](NewIkFrame, NullNewIkFrame),
	}
}
