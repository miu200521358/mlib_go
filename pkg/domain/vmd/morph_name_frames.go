package vmd

type MorphNameFrames struct {
	*BaseFrames[*MorphFrame]
	Name string // モーフ名
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		BaseFrames: NewBaseFrames[*MorphFrame](),
		Name:       name,
	}
}
