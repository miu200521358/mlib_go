package mbt

type MPhysics struct {
	world    *BtDiscreteDynamicsWorld
	filterCB *MFilterCallback
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
	world.GetPairCache().SetOverlapFilterCallback(filterCB)

	p := &MPhysics{
		world:    &world,
		filterCB: filterCB,
	}

	return p
}
