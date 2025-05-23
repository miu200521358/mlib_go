package vmd

type LightFrames struct {
	*BaseFrames[*LightFrame]
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		BaseFrames: NewBaseFrames[*LightFrame](NewLightFrame, nilLightFrame),
	}
}

func nilLightFrame() *LightFrame {
	return nil
}

func (lightFrames *LightFrames) Clean() {
	if lightFrames.Length() > 1 {
		return
	} else {
		cf := lightFrames.Get(lightFrames.values.Min())
		if !(cf.Position == nil || cf.Position.Length() == 0 ||
			cf.Color == nil || cf.Color.Length() == 0) {
			return
		}
		lightFrames.Delete(cf.Index())
	}
}
