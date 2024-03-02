package mbt

type MPhysics struct {
	broadphase             *BtDbvtBroadphase
	collisionConfiguration *BtCollisionConfiguration
}

func NewMPhysics() *MPhysics {
	broadphase := NewBtDbvtBroadphase()
	// collisionConfiguration := NewBtCollisionConfiguration()

	p := &MPhysics{
		broadphase: &broadphase,
		// collisionConfiguration: &collisionConfiguration,
	}

	return p
}
