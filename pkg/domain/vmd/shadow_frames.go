package vmd

type ShadowFrames struct {
	*BaseFrames[*ShadowFrame]
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		BaseFrames: NewBaseFrames[*ShadowFrame](NewShadowFrame, NullShadowFrame),
	}
}
