package vmd

type MorphNameFrames struct {
	*BaseFrames[*MorphFrame]
	Name string // ボーン名
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		BaseFrames: NewBaseFrames[*MorphFrame](NewMorphFrame, NullMorphFrame),
		Name:       name,
	}
}

func (i *MorphNameFrames) NewFrame(index int) *MorphFrame {
	return NewMorphFrame(index)
}
