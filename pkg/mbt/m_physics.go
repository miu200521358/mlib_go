package mbt

type MPhysics struct {
	World           BtDiscreteDynamicsWorld
	GroundTransform BtTransform
	FilterCallBack  MFilterCallback
	MotionState     BtDefaultMotionState
}

func NewMPhysics() *MPhysics {
	broadphase := NewBtDbvtBroadphase()
	collisionConfiguration := NewBtDefaultCollisionConfiguration()
	dispatcher := NewBtCollisionDispatcher(collisionConfiguration)
	solver := NewBtSequentialImpulseConstraintSolver()
	world := NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(NewBtVector3(float32(0), float32(-9.8*10.0), float32(0)))

	groundShape := NewBtStaticPlaneShape(NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := NewBtTransform()
	groundMotionState := NewBtDefaultMotionState(groundTransform)
	groundRigidBody := NewBtRigidBody(float32(0), groundMotionState, groundShape, NewBtVector3(float32(0), float32(0), float32(0)))

	world.AddRigidBody(groundRigidBody)

	filterCB := NewMFilterCallback(world.GetPairCache().GetOverlapFilterCallback())
	filterCB.m_nonFilterProxy = append(filterCB.m_nonFilterProxy, groundRigidBody.GetBroadphaseProxy())
	world.GetPairCache().SetOverlapFilterCallback(&filterCB)

	motionState := NewBtDefaultMotionState(groundTransform)

	p := &MPhysics{
		World:           world,
		GroundTransform: groundTransform,
		FilterCallBack:  filterCB,
		MotionState:     motionState,
	}

	return p
}
