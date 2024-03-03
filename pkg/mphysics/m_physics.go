package mphysics

import "github.com/miu200521358/mlib_go/pkg/mbt"

type MPhysics struct {
	World           mbt.BtDiscreteDynamicsWorld
	GroundTransform mbt.BtTransform
	FilterCallBack  MFilterCallback
	MotionState     mbt.BtDefaultMotionState
	MaxSubSteps     int
	Fps             float64
}

func NewMPhysics() *MPhysics {
	broadphase := mbt.NewBtDbvtBroadphase()
	collisionConfiguration := mbt.NewBtDefaultCollisionConfiguration()
	dispatcher := mbt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := mbt.NewBtSequentialImpulseConstraintSolver()
	world := mbt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(mbt.NewBtVector3(float32(0), float32(-9.8*10.0), float32(0)))

	groundShape := mbt.NewBtStaticPlaneShape(mbt.NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := mbt.NewBtTransform()
	groundMotionState := mbt.NewBtDefaultMotionState(groundTransform)
	groundRigidBody := mbt.NewBtRigidBody(float32(0), groundMotionState, groundShape, mbt.NewBtVector3(float32(0), float32(0), float32(0)))

	world.AddRigidBody(groundRigidBody)

	filterCB := NewMFilterCallback(world.GetPairCache().GetOverlapFilterCallback())
	filterCB.m_nonFilterProxy = append(filterCB.m_nonFilterProxy, groundRigidBody.GetBroadphaseProxy())
	world.GetPairCache().SetOverlapFilterCallback(&filterCB)

	motionState := mbt.NewBtDefaultMotionState(groundTransform)

	p := &MPhysics{
		World:           world,
		GroundTransform: groundTransform,
		FilterCallBack:  filterCB,
		MotionState:     motionState,
		MaxSubSteps:     10,
		Fps:             60.0,
	}

	return p
}

func (p *MPhysics) AddRigidBody(rigidBody mbt.BtRigidBody) {
	p.World.AddRigidBody(rigidBody)
}

func (p *MPhysics) RemoveRigidBody(rigidBody mbt.BtRigidBody) {
	p.World.RemoveRigidBody(rigidBody)
}

func (p *MPhysics) AddJoint(joint mbt.BtTypedConstraint) {
	p.World.AddConstraint(joint, true)
}

func (p *MPhysics) RemoveJoint(joint mbt.BtTypedConstraint) {
	p.World.RemoveConstraint(joint)
}

func (p *MPhysics) Update() {
	p.World.StepSimulation(float32(1.0/60.0), p.MaxSubSteps, float32(1.0/p.Fps))
}
