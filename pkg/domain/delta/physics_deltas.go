package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type PhysicsDeltas struct {
	frame       float32
	modelHash   string
	motionHash  string
	RigidBodies *RigidBodyDeltas
	Joints      *JointDeltas
}

func NewPhysicsDeltas(
	frame float32, rigidBodies *pmx.RigidBodies, joints *pmx.Joints, modelHash, motionHash string,
) *PhysicsDeltas {
	return &PhysicsDeltas{
		frame:       frame,
		modelHash:   modelHash,
		motionHash:  motionHash,
		RigidBodies: NewRigidBodyDeltas(rigidBodies),
		Joints:      NewJointDeltas(joints),
	}
}

func (vmdDeltas *PhysicsDeltas) Frame() float32 {
	return vmdDeltas.frame
}

func (vmdDeltas *PhysicsDeltas) SetFrame(frame float32) {
	vmdDeltas.frame = frame
}

func (vmdDeltas *PhysicsDeltas) ModelHash() string {
	return vmdDeltas.modelHash
}

func (vmdDeltas *PhysicsDeltas) SetModelHash(hash string) {
	vmdDeltas.modelHash = hash
}

func (vmdDeltas *PhysicsDeltas) MotionHash() string {
	return vmdDeltas.motionHash
}

func (vmdDeltas *PhysicsDeltas) SetMotionHash(hash string) {
	vmdDeltas.motionHash = hash
}
