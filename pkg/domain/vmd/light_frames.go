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

func (lightFrames *LightFrames) Clean() {
	if lightFrames.Len() > 1 {
		return
	} else {
		cf := lightFrames.Get(lightFrames.Indexes.Min())
		if !(cf.Position == nil || cf.Position.Length() == 0 ||
			cf.Color == nil || cf.Color.Length() == 0) {
			return
		}
		lightFrames.Delete(cf.Index())
	}
}
