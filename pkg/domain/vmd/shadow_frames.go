package vmd

type ShadowFrames struct {
	*BaseFrames[*ShadowFrame]
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		BaseFrames: NewBaseFrames[*ShadowFrame](NewShadowFrame, NullShadowFrame),
	}
}

func (shadowFrames *ShadowFrames) Copy() *ShadowFrames {
	copied := NewShadowFrames()
	for _, frame := range shadowFrames.List() {
		copied.Append(frame.Copy().(*ShadowFrame))
	}
	return copied
}

func (shadowFrames *ShadowFrames) Clean() {
	if shadowFrames.Len() > 1 {
		return
	} else {
		cf := shadowFrames.Get(shadowFrames.Indexes.Min())
		if cf.Distance != 0 {
			return
		}
		shadowFrames.Delete(cf.Index())
	}
}
