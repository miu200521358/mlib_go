//go:build windows
// +build windows

package mphysics

import (
	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
	"github.com/miu200521358/mlib_go/pkg/mview"
)

type MPhysics struct {
	world            mbt.BtDiscreteDynamicsWorld
	drawer           mbt.BtMDebugDraw
	liner            *MDebugDrawLiner
	MaxSubSteps      int
	DeformFps        float32
	DeformSpf        float32
	PhysicsFps       float32
	PhysicsSpf       float32
	FixedTimeStep    float32
	joints           map[jointKey]mbt.BtTypedConstraint
	rigidBodies      map[rigidbodyKey]rigidbodyValue
	visibleRigidBody bool
	visibleJoint     bool
}

type rigidbodyKey struct {
	ModelIndex     int
	RigidBodyIndex int
}

type rigidbodyValue struct {
	RigidBody mbt.BtRigidBody
	Transform mbt.BtTransform
	Mask      int
	Group     int
}

type jointKey struct {
	ModelIndex int
	JointIndex int
}

func NewMPhysics(shader *mview.MShader) *MPhysics {
	world := createWorld()

	p := &MPhysics{
		world:       world,
		MaxSubSteps: 2,
		DeformFps:   30.0,
		PhysicsFps:  60.0,
		rigidBodies: make(map[rigidbodyKey]rigidbodyValue),
		joints:      make(map[jointKey]mbt.BtTypedConstraint),
	}

	// デバッグビューワー
	liner := NewMDebugDrawLiner(shader)
	drawer := mbt.NewBtMDebugDraw()
	drawer.SetLiner(liner)
	drawer.SetMDefaultColors(NewConstBtMDefaultColors())
	world.SetDebugDrawer(drawer)
	// mlog.D("world.GetDebugDrawer()=%+v\n", world.GetDebugDrawer())

	p.drawer = drawer
	p.liner = liner
	p.DeformSpf = 1.0 / p.DeformFps
	p.PhysicsSpf = 1.0 / p.PhysicsFps
	p.FixedTimeStep = 1 / 60.0

	p.VisibleRigidBody(false)
	p.VisibleJoint(false)

	return p
}

func createWorld() mbt.BtDiscreteDynamicsWorld {
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

	return world
}

func (p *MPhysics) ResetWorld() {
	world := createWorld()
	world.SetDebugDrawer(p.drawer)
	p.world = world
	for _, r := range p.rigidBodies {
		rigidBody := r.RigidBody
		group := r.Group
		mask := r.Mask
		p.world.AddRigidBody(rigidBody, group, mask)
	}
	for _, joint := range p.joints {
		p.world.AddConstraint(joint, true)
	}
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
	p.visibleRigidBody = enable
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
	p.visibleJoint = enable
}

func (p *MPhysics) DrawDebugLines() {
	if p.visibleRigidBody || p.visibleJoint {
		p.liner.DrawDebugLines()
		p.liner.vertices = []float32{}
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

func (p *MPhysics) GetRigidBody(modelIndex, rigidBodyIndex int) (mbt.BtRigidBody, mbt.BtTransform) {
	r := p.rigidBodies[rigidbodyKey{ModelIndex: modelIndex, RigidBodyIndex: rigidBodyIndex}]
	return r.RigidBody, r.Transform
}

func (p *MPhysics) AddRigidBody(rigidBody mbt.BtRigidBody, rigidBodyTransform mbt.BtTransform,
	modelIndex, rigidBodyIndex, group, mask int) {
	p.world.AddRigidBody(rigidBody, group, mask)
	rigidBodyKey := rigidbodyKey{ModelIndex: modelIndex, RigidBodyIndex: rigidBodyIndex}
	p.rigidBodies[rigidBodyKey] = rigidbodyValue{
		RigidBody: rigidBody, Transform: rigidBodyTransform, Mask: mask, Group: group}
}

func (p *MPhysics) DeleteRigidBodies(modelIndex int) {
	for k, r := range p.rigidBodies {
		rigidBody := r.RigidBody
		if k.ModelIndex == modelIndex {
			p.world.RemoveRigidBody(rigidBody)
		}
	}
}

func (p *MPhysics) AddJoint(modelIndex, jointIndex int, joint mbt.BtTypedConstraint) {
	p.world.AddConstraint(joint, true)
	key := jointKey{ModelIndex: modelIndex, JointIndex: jointIndex}
	p.joints[key] = joint
}

func (p *MPhysics) DeleteJoints(modelIndex int) {
	for k, joint := range p.joints {
		if k.ModelIndex == modelIndex {
			p.world.RemoveConstraint(joint)
		}
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
