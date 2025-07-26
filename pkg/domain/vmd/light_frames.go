package vmd

import "github.com/tiendc/go-deepcopy"

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

func (lightFrames *LightFrames) Copy() (*LightFrames, error) {
	copied := new(LightFrames)
	err := deepcopy.Copy(copied, lightFrames)
	return copied, err
}
