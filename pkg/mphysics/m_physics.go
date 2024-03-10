package mphysics

import (
	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mgl"
)

type MPhysics struct {
	world       mbt.BtDiscreteDynamicsWorld
	MaxSubSteps int
	Fps         float32
	Spf         float32
	isDebug     bool
}

func NewMPhysics(shader *mgl.MShader) *MPhysics {
	broadphase := mbt.NewBtDbvtBroadphase()
	collisionConfiguration := mbt.NewBtDefaultCollisionConfiguration()
	dispatcher := mbt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := mbt.NewBtSequentialImpulseConstraintSolver()
	world := mbt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(mbt.NewBtVector3(float32(0), float32(-9.8*10), float32(0)))

	groundShape := mbt.NewBtStaticPlaneShape(mbt.NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := mbt.NewBtTransform()
	groundTransform.SetIdentity()
	groundTransform.SetOrigin(mbt.NewBtVector3(float32(0), float32(0), float32(0)))
	groundMotionState := mbt.NewBtDefaultMotionState(groundTransform)
	groundRigidBody := mbt.NewBtRigidBody(float32(0), groundMotionState, groundShape)

	world.AddRigidBody(groundRigidBody, 1<<15, 0xFFFF)

	// // 剛体衝突判定用のコールバック
	// world.GetPairCache().SetOverlapFilterCallback(NewMFilterCallback())

	// デバッグビューワー
	// world.SetDebugDrawer(NewMDebugView(shader))

	p := &MPhysics{
		world:       world,
		MaxSubSteps: 5,
		Fps:         60.0,
		Spf:         1.0 / 60.0,
		isDebug:     false,
	}

	return p
}

func (p *MPhysics) EnableDebug(enable bool) {
	p.isDebug = enable
	if enable {
		p.world.GetDebugDrawer().SetDebugMode(int(mbt.BtIDebugDrawDBG_DrawConstraints | mbt.BtIDebugDrawDBG_DrawConstraintLimits | mbt.BtIDebugDrawDBG_DrawWireframe | mbt.BtIDebugDrawDBG_DrawAabb))
	} else {
		p.world.GetDebugDrawer().SetDebugMode(0)
	}
}

func (p *MPhysics) DrawWorld() {
	if p.isDebug {
		p.world.DebugDrawWorld()
	}
}

// func (p *MPhysics) AddNonFilterProxy(proxy mbt.BtBroadphaseProxy) {
// 	pairCache := p.world.GetPairCache()

// }

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

func (p *MPhysics) Update() {
	p.world.StepSimulation(p.Spf, p.MaxSubSteps, p.Spf)
}
