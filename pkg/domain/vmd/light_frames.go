package vmd

type LightFrames struct {
	*BaseFrames[*LightFrame]
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		BaseFrames: NewBaseFrames[*LightFrame](NewLightFrame, NullLightFrame),
	}
}

func (i *LightFrames) NewFrame(index int) *LightFrame {
	return NewLightFrame(index)
}
