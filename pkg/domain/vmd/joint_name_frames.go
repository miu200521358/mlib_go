package vmd

import (
	"github.com/tiendc/go-deepcopy"
)

type JointNameFrames struct {
	*BaseFrames[*JointFrame]
	Name string // ジョイント名
}

func NewJointNameFrames(name string) *JointNameFrames {
	return &JointNameFrames{
		BaseFrames: NewBaseFrames(newJointFrame, nilJointFrame),
		Name:       name,
	}
}

func newJointFrame(index float32) *JointFrame {
	return nil
}

func nilJointFrame() *JointFrame {
	return nil
}

func (jointNameFrames *JointNameFrames) Copy() (*JointNameFrames, error) {
	copied := new(JointNameFrames)
	err := deepcopy.Copy(copied, jointNameFrames)
	return copied, err
}
