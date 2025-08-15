package vmd

import (
	"github.com/tiendc/go-deepcopy"
)

type RigidBodyNameFrames struct {
	*BaseFrames[*RigidBodyFrame]
	Name string // 剛体名
}

func NewRigidBodyNameFrames(name string) *RigidBodyNameFrames {
	return &RigidBodyNameFrames{
		BaseFrames: NewBaseFrames(NewRigidBodyFrame, nilRigidBodyFrame),
		Name:       name,
	}
}

func nilRigidBodyFrame() *RigidBodyFrame {
	return nil
}

func (rigidBodyNameFrames *RigidBodyNameFrames) Copy() (*RigidBodyNameFrames, error) {
	copied := new(RigidBodyNameFrames)
	err := deepcopy.Copy(copied, rigidBodyNameFrames)
	return copied, err
}
