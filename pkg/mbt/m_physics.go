package mbt

type MPhysics struct {
	broadphase             *BtDbvtBroadphase
	collisionConfiguration *BtDefaultCollisionConfiguration
}

func NewMPhysics() *MPhysics {
	broadphase := NewBtDbvtBroadphase()
	collisionConfiguration := NewBtDefaultCollisionConfiguration()

	p := &MPhysics{
		broadphase:             &broadphase,
		collisionConfiguration: &collisionConfiguration,
	}

	return p
}
