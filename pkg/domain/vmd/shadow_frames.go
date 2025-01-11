package vmd

type ShadowFrames struct {
	*BaseFrames[*ShadowFrame]
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		BaseFrames: NewBaseFrames[*ShadowFrame](),
	}
}

func (shadowFrames *ShadowFrames) Clean() {
	if shadowFrames.Length() > 1 {
		return
	} else {
		cf := shadowFrames.Get(shadowFrames.Indexes.Min())
		if cf.Distance != 0 {
			return
		}
		shadowFrames.Delete(cf.Index())
	}
}
