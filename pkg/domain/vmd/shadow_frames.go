package vmd

type ShadowFrames struct {
	*BaseFrames[*ShadowFrame]
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		BaseFrames: NewBaseFrames[*ShadowFrame](NewShadowFrame, NullShadowFrame),
	}
}

func (shadowFrames *ShadowFrames) Copy() *ShadowFrames {
	copied := NewShadowFrames()
	for _, frame := range shadowFrames.List() {
		copied.Append(frame.Copy().(*ShadowFrame))
	}
	return copied
}
