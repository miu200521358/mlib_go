package vmd

type BoneNameFrames struct {
	*BaseFrames[*BoneFrame]
	Name string // ボーン名
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		BaseFrames: NewBaseFrames[*BoneFrame](NewBoneFrame, NullBoneFrame),
		Name:       name,
	}
}

func (boneNameFrames *BoneNameFrames) Copy() *BoneNameFrames {
	copied := NewBoneNameFrames(boneNameFrames.Name)
	for _, frame := range boneNameFrames.List() {
		copied.Append(frame.Copy().(*BoneFrame))
	}
	return copied
}
