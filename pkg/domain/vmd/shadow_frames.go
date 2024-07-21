package vmd

type ShadowFrames struct {
	*BaseFrames[*ShadowFrame]
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		BaseFrames: NewBaseFrames[*ShadowFrame](NewShadowFrame, NullShadowFrame),
	}
}

func (shadowFrames *ShadowFrames) NewFrame(index int) *ShadowFrame {
	return NewShadowFrame(index)
}
