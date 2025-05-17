package vmd

type ShadowFrames struct {
	*BaseFrames[*ShadowFrame]
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		BaseFrames: NewBaseFrames[*ShadowFrame](NewShadowFrame, nilShadowFrame),
	}
}

func nilShadowFrame() *ShadowFrame {
	return nil
}

func (shadowFrames *ShadowFrames) Clean() {
	if shadowFrames.Length() > 1 {
		return
	} else {
		cf := shadowFrames.Get(shadowFrames.values.Min())
		if cf.Distance != 0 {
			return
		}
		shadowFrames.Delete(cf.Index())
	}
}
