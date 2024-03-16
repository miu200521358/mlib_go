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
	world.SetGravity(mbt.NewBtVector3(float32(0), float32(-9.8), float32(0)))

	groundShape := mbt.NewBtStaticPlaneShape(mbt.NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := mbt.NewBtTransform()
	groundTransform.SetIdentity()
	groundTransform.SetOrigin(mbt.NewBtVector3(float32(0), float32(0), float32(0)))
	groundMotionState := mbt.NewBtDefaultMotionState(groundTransform)
	groundRigidBody := mbt.NewBtRigidBody(float32(0), groundMotionState, groundShape)

	world.AddRigidBody(groundRigidBody, 1<<15, 0xFFFF)

	// デバッグビューワー
	drawer := mbt.NewBtMDebugDraw()
	drawer.SetLiner(NewMDebugDrawLiner(shader))
	drawer.SetMDefaultColors(NewConstBtMDefaultColors())
	world.SetDebugDrawer(drawer)
	// fmt.Printf("world.GetDebugDrawer()=%+v\n", world.GetDebugDrawer())

	p := &MPhysics{
		world:       world,
		MaxSubSteps: 5,
		Fps:         30.0,
		isDebug:     false,
	}
	p.Spf = 1.0 / p.Fps

	p.EnableDebug(false)

	return p
}

func (p *MPhysics) EnableDebug(enable bool) {
	p.isDebug = enable
	if enable {
		p.world.GetDebugDrawer().SetDebugMode(int(
			mbt.BtIDebugDrawDBG_DrawConstraints |
				mbt.BtIDebugDrawDBG_DrawConstraintLimits |
				mbt.BtIDebugDrawDBG_DrawWireframe |
				mbt.BtIDebugDrawDBG_DrawContactPoints,
		))
	} else {
		p.world.GetDebugDrawer().SetDebugMode(0)
	}
}

func (p *MPhysics) DrawWorld() {
	if p.isDebug {
		// fmt.Printf("DrawWorld p.world=%+v\n", p.world)

		// // 標準出力を一時的にリダイレクトする
		// old := os.Stdout // keep backup of the real stdout
		// r, w, _ := os.Pipe()
		// os.Stdout = w

		p.world.DebugDrawWorld()

		// // 標準出力を元に戻す
		// w.Close()
		// os.Stdout = old // restoring the real stdout

		// var buf bytes.Buffer
		// buf.ReadFrom(r)
		// fmt.Print(buf.String())
	}
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

func (p *MPhysics) Update(timeStep float32) {
	// // 標準出力を一時的にリダイレクトする
	// old := os.Stdout // keep backup of the real stdout
	// r, w, _ := os.Pipe()
	// os.Stdout = w

	p.world.StepSimulation(timeStep, p.MaxSubSteps, timeStep/float32(p.MaxSubSteps))

	// // p.frame += float32(elapsed)
	// fmt.Printf("timeStep: %.8f [p.world.StepSimulation]\n", timeStep)

	// // 標準出力を元に戻す
	// w.Close()
	// os.Stdout = old // restoring the real stdout

	// var buf bytes.Buffer
	// buf.ReadFrom(r)
	// fmt.Print(buf.String())
}
