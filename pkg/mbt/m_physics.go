package mbt

type MPhysics struct {
	broadphase             *BtDbvtBroadphase
	collisionConfiguration *BtDefaultCollisionConfiguration
	dispatcher             *BtCollisionDispatcher
	solver                 *BtSequentialImpulseConstraintSolver
	world                  *BtDiscreteDynamicsWorld
}

func NewMPhysics() *MPhysics {
	broadphase := NewBtDbvtBroadphase()
	collisionConfiguration := NewBtDefaultCollisionConfiguration()
	dispatcher := NewBtCollisionDispatcher(collisionConfiguration)
	solver := NewBtSequentialImpulseConstraintSolver()
	world := NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(NewBtVector3(float32(0), float32(-9.8 * 10.0), float32(0)))

	groundShape := NewBtStaticPlaneShape(NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := NewBtTransform()
	groundMotionState := NewBtDefaultMotionState(groundTransform)
	groundRigidBody := NewBtRigidBody(float32(0), groundMotionState, groundShape, NewBtVector3(float32(0), float32(0), float32(0)))

	world.AddRigidBody(groundRigidBody)

	p := &MPhysics{
		broadphase:             &broadphase,
		collisionConfiguration: &collisionConfiguration,
		dispatcher:             &dispatcher,
		solver:                 &solver,
		world:                  &world,
	}

	return p
}
