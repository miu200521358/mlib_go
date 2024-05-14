package vmd

type LightFrames struct {
	*BaseFrames[*LightFrame]
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		BaseFrames: NewBaseFrames[*LightFrame](NewLightFrame),
	}
}

func (i *LightFrames) NewFrame(index int) *LightFrame {
	return NewLightFrame(index)
}
