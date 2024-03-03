package mphysics

import "github.com/miu200521358/mlib_go/pkg/mbt"

type MPhysics struct {
	World           mbt.BtDiscreteDynamicsWorld
	GroundTransform mbt.BtTransform
	FilterCallBack  MFilterCallback
	MotionState     mbt.BtDefaultMotionState
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
	}

	return p
}
