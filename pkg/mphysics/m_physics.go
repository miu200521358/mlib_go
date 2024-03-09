package mphysics

import "github.com/miu200521358/mlib_go/pkg/mbt"

type MPhysics struct {
	world           mbt.BtDiscreteDynamicsWorld
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
	groundTransform.SetIdentity()
	groundTransform.SetOrigin(mbt.NewBtVector3(float32(0), float32(0), float32(0)))
	groundMotionState := mbt.NewBtDefaultMotionState(groundTransform)
	groundRigidBody := mbt.NewBtRigidBody(float32(0), groundMotionState, groundShape)

	world.AddRigidBody(groundRigidBody, 1<<15, 0xFFFF)

	callback := NewMFilterCallback()
	world.GetPairCache().SetOverlapFilterCallback(&callback)

	motionState := mbt.NewBtDefaultMotionState(groundTransform)

	p := &MPhysics{
		world:           world,
		GroundTransform: groundTransform,
		MotionState:     motionState,
		MaxSubSteps:     5,
		Fps:             60.0,
	}

	return p
}

func (p *MPhysics) AddNonFilterProxy(proxy mbt.BtBroadphaseProxy) {
	p.FilterCallBack.nonFilterProxy = append(p.FilterCallBack.nonFilterProxy, proxy)
}

func (p *MPhysics) AddRigidBody(rigidBody mbt.BtRigidBody, group int, mask int) {
	p.world.AddRigidBody(rigidBody, group, mask)
}

func (p *MPhysics) RemoveRigidBody(rigidBody mbt.BtRigidBody) {
	p.world.RemoveRigidBody(rigidBody)
}

func (p *MPhysics) AddJoint(joint mbt.BtTypedConstraint) {
	p.world.AddConstraint(joint, true)
}

func (p *MPhysics) RemoveJoint(joint mbt.BtTypedConstraint) {
	p.world.RemoveConstraint(joint)
}

func (p *MPhysics) Update(elapsed float32) {
	p.world.StepSimulation(elapsed, p.MaxSubSteps, float32(1.0/p.Fps))
}
