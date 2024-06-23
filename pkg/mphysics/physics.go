//go:build windows
// +build windows

package mphysics

import (
	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
	"github.com/miu200521358/mlib_go/pkg/mview"
)

type MPhysics struct {
	world               mbt.BtDiscreteDynamicsWorld
	MaxSubSteps         int
	Fps                 float32
	Spf                 float32
	PhysicsSpf          float64
	FixedTimeStep       float32
	joints              []mbt.BtTypedConstraint
	rigidBodies         map[int]mbt.BtRigidBody
	rigidBodyTransforms map[int]mbt.BtTransform // 剛体の初期位置・回転情報
}

func NewMPhysics(shader *mview.MShader) *MPhysics {
	broadphase := mbt.NewBtDbvtBroadphase()
	collisionConfiguration := mbt.NewBtDefaultCollisionConfiguration()
	dispatcher := mbt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := mbt.NewBtSequentialImpulseConstraintSolver()
	// solver.GetM_analyticsData().SetM_numIterationsUsed(200)
	world := mbt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(mbt.NewBtVector3(float32(0), float32(-9.8*10), float32(0)))
	// world.GetSolverInfo().(mbt.BtContactSolverInfo).SetM_numIterations(100)
	// world.GetSolverInfo().(mbt.BtContactSolverInfo).SetM_splitImpulse(1)

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
	// mlog.D("world.GetDebugDrawer()=%+v\n", world.GetDebugDrawer())

	p := &MPhysics{
		world:               world,
		MaxSubSteps:         5,
		Fps:                 30.0,
		PhysicsSpf:          1.0 / 60.0,
		rigidBodies:         make(map[int]mbt.BtRigidBody),
		rigidBodyTransforms: make(map[int]mbt.BtTransform),
	}
	p.Spf = 1.0 / p.Fps
	p.FixedTimeStep = 1 / 60.0

	p.VisibleRigidBody(false)
	p.VisibleJoint(false)

	return p
}

func (p *MPhysics) VisibleRigidBody(enable bool) {
	if enable {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() | int(
				mbt.BtIDebugDrawDBG_DrawWireframe|
					mbt.BtIDebugDrawDBG_DrawContactPoints,
			))
	} else {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() & ^int(
				mbt.BtIDebugDrawDBG_DrawWireframe|
					mbt.BtIDebugDrawDBG_DrawContactPoints,
			))
	}
}

func (p *MPhysics) VisibleJoint(enable bool) {
	if enable {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() | int(
				mbt.BtIDebugDrawDBG_DrawConstraints|
					mbt.BtIDebugDrawDBG_DrawConstraintLimits,
			))
	} else {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() & ^int(
				mbt.BtIDebugDrawDBG_DrawConstraints|
					mbt.BtIDebugDrawDBG_DrawConstraintLimits,
			))
	}
}

func (p *MPhysics) DebugDrawWorld() {
	// mlog.D("DrawWorld p.world=%+v\n", p.world)

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

func (p *MPhysics) GetRigidBody(index int) (mbt.BtRigidBody, mbt.BtTransform) {
	return p.rigidBodies[index], p.rigidBodyTransforms[index]
}

func (p *MPhysics) AddRigidBody(rigidBody mbt.BtRigidBody, rigidBodyTransform mbt.BtTransform,
	rigidBodyIndex int, group int, mask int) {
	p.world.AddRigidBody(rigidBody, group, mask)
	p.rigidBodies[rigidBodyIndex] = rigidBody
	p.rigidBodyTransforms[rigidBodyIndex] = rigidBodyTransform
}

func (p *MPhysics) DeleteRigidBodies() {
	for _, rigidBody := range p.rigidBodies {
		p.world.RemoveRigidBody(rigidBody)
	}
}

func (p *MPhysics) AddJoint(joint mbt.BtTypedConstraint) {
	p.world.AddConstraint(joint, true)
	p.joints = append(p.joints, joint)
}

func (p *MPhysics) DeleteJoints() {
	for _, joint := range p.joints {
		p.world.RemoveConstraint(joint)
	}
}

func (p *MPhysics) Update(timeStep float32) {
	// // 標準出力を一時的にリダイレクトする
	// old := os.Stdout // keep backup of the real stdout
	// r, w, _ := os.Pipe()
	// os.Stdout = w
	// for range int(timeStep / p.Spf) {
	// mlog.I("timeStep=%.8f, spf: %.8f", timeStep, p.Spf)
	p.world.StepSimulation(timeStep, p.MaxSubSteps, p.FixedTimeStep)
	// }

	// // p.frame += float32(elapsed)
	// mlog.D("timeStep: %.8f [p.world.StepSimulation]\n", timeStep)

	// // 標準出力を元に戻す
	// w.Close()
	// os.Stdout = old // restoring the real stdout

	// var buf bytes.Buffer
	// buf.ReadFrom(r)
	// fmt.Print(buf.String())
}
