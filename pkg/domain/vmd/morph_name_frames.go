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

func (morphNameFrames *MorphNameFrames) Copy() *MorphNameFrames {
	copied := NewMorphNameFrames(morphNameFrames.Name)
	for _, frame := range morphNameFrames.List() {
		copied.Append(frame.Copy().(*MorphFrame))
	}
	return copied
}
