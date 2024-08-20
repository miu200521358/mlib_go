package vmd

type LightFrames struct {
	*BaseFrames[*LightFrame]
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		BaseFrames: NewBaseFrames[*LightFrame](NewLightFrame, NullLightFrame),
	}
}

func (lightFrames *LightFrames) Copy() *LightFrames {
	copied := NewLightFrames()
	for _, frame := range lightFrames.List() {
		copied.Append(frame.Copy().(*LightFrame))
	}
	return copied
}
