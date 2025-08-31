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
		BaseFrames: NewBaseFrames(newRigidBodyFrame, nilRigidBodyFrame),
		Name:       name,
	}
}

func newRigidBodyFrame(index float32) *RigidBodyFrame {
	// デフォルト値とかを入れないよう、nilを返す
	return nil
}

func nilRigidBodyFrame() *RigidBodyFrame {
	return nil
}

func (rigidBodyNameFrames *RigidBodyNameFrames) Copy() (*RigidBodyNameFrames, error) {
	copied := new(RigidBodyNameFrames)
	err := deepcopy.Copy(copied, rigidBodyNameFrames)
	return copied, err
}
