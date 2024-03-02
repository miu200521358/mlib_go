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

	p := &MPhysics{
		broadphase:             &broadphase,
		collisionConfiguration: &collisionConfiguration,
		dispatcher:             &dispatcher,
		solver:                 &solver,
		world:                  &world,
	}

	return p
}
