//go:build windows
// +build windows

package mbt

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type MPhysics struct {
	world         bt.BtDiscreteDynamicsWorld
	drawer        bt.BtMDebugDraw
	liner         *MDebugDrawLiner
	MaxSubSteps   int
	DeformFps     float32
	DeformSpf     float32
	PhysicsFps    float32
	PhysicsSpf    float32
	FixedTimeStep float32
	joints        map[jointKey]bt.BtTypedConstraint
	rigidBodies   map[rigidbodyKey]rigidbodyValue
}

type rigidbodyKey struct {
	ModelIndex     int
	RigidBodyIndex int
}

type rigidbodyValue struct {
	RigidBody bt.BtRigidBody
	Transform bt.BtTransform
	Mask      int
	Group     int
}

type jointKey struct {
	ModelIndex int
	JointIndex int
}

func NewMPhysics(shader *mgl.MShader) *MPhysics {
	world := createWorld()

	p := &MPhysics{
		world:       world,
		MaxSubSteps: 2,
		DeformFps:   30.0,
		PhysicsFps:  60.0,
		rigidBodies: make(map[rigidbodyKey]rigidbodyValue),
		joints:      make(map[jointKey]bt.BtTypedConstraint),
	}

	// デバッグビューワー
	liner := NewMDebugDrawLiner(shader)
	drawer := bt.NewBtMDebugDraw()
	drawer.SetLiner(liner)
	drawer.SetMDefaultColors(NewConstBtMDefaultColors())
	world.SetDebugDrawer(drawer)
	// mlog.D("world.GetDebugDrawer()=%+v\n", world.GetDebugDrawer())

	p.drawer = drawer
	p.liner = liner
	p.DeformSpf = 1.0 / p.DeformFps
	p.PhysicsSpf = 1.0 / p.PhysicsFps
	p.FixedTimeStep = 1 / 60.0

	return p
}

func createWorld() bt.BtDiscreteDynamicsWorld {
	broadphase := bt.NewBtDbvtBroadphase()
	collisionConfiguration := bt.NewBtDefaultCollisionConfiguration()
	dispatcher := bt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := bt.NewBtSequentialImpulseConstraintSolver()
	// solver.GetM_analyticsData().SetM_numIterationsUsed(200)
	world := bt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(bt.NewBtVector3(float32(0), float32(-9.8*10), float32(0)))
	// world.GetSolverInfo().(bt.BtContactSolverInfo).SetM_numIterations(100)
	// world.GetSolverInfo().(bt.BtContactSolverInfo).SetM_splitImpulse(1)

	groundShape := bt.NewBtStaticPlaneShape(bt.NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := bt.NewBtTransform()
	groundTransform.SetIdentity()
	groundTransform.SetOrigin(bt.NewBtVector3(float32(0), float32(0), float32(0)))
	groundMotionState := bt.NewBtDefaultMotionState(groundTransform)
	groundRigidBody := bt.NewBtRigidBody(float32(0), groundMotionState, groundShape)

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

func (p *MPhysics) DrawDebugLines(isDrawRigidBodyFront bool) {
	// // 標準出力を一時的にリダイレクトする
	// old := os.Stdout // keep backup of the real stdout
	// r, w, _ := os.Pipe()
	// os.Stdout = w
	if isDrawRigidBodyFront {
		// モデルメッシュの前面に描画するために深度テストを無効化
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.ALWAYS)
	}

	p.liner.DrawDebugLines()
	p.liner.vertices = []float32{}

	// 深度テストを有効に戻す
	if isDrawRigidBodyFront {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LEQUAL)
	}
}

func (p *MPhysics) DebugDrawWorld(visibleRigidBody, visibleJoint bool) {
	// mlog.D("DrawWorld p.world=%+v\n", p.world)

	if visibleRigidBody {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() | int(
				bt.BtIDebugDrawDBG_DrawWireframe|
					bt.BtIDebugDrawDBG_DrawContactPoints,
			))
	} else {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() & ^int(
				bt.BtIDebugDrawDBG_DrawWireframe|
					bt.BtIDebugDrawDBG_DrawContactPoints,
			))
	}

	if visibleJoint {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() | int(
				bt.BtIDebugDrawDBG_DrawConstraints|
					bt.BtIDebugDrawDBG_DrawConstraintLimits,
			))
	} else {
		p.world.GetDebugDrawer().SetDebugMode(
			p.world.GetDebugDrawer().GetDebugMode() & ^int(
				bt.BtIDebugDrawDBG_DrawConstraints|
					bt.BtIDebugDrawDBG_DrawConstraintLimits,
			))
	}

	p.world.DebugDrawWorld()

	// // 標準出力を元に戻す
	// w.Close()
	// os.Stdout = old // restoring the real stdout

	// var buf bytes.Buffer
	// buf.ReadFrom(r)
	// fmt.Print(buf.String())
}

func (p *MPhysics) GetRigidBody(modelIndex, rigidBodyIndex int) (bt.BtRigidBody, bt.BtTransform) {
	r := p.rigidBodies[rigidbodyKey{ModelIndex: modelIndex, RigidBodyIndex: rigidBodyIndex}]
	return r.RigidBody, r.Transform
}

func (p *MPhysics) AddRigidBody(rigidBody bt.BtRigidBody, rigidBodyTransform bt.BtTransform,
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
			delete(p.rigidBodies, k)
		}
	}
}

func (p *MPhysics) AddJoint(modelIndex, jointIndex int, joint bt.BtTypedConstraint) {
	p.world.AddConstraint(joint, true)
	key := jointKey{ModelIndex: modelIndex, JointIndex: jointIndex}
	p.joints[key] = joint
}

func (p *MPhysics) DeleteJoints(modelIndex int) {
	for k, joint := range p.joints {
		if k.ModelIndex == modelIndex {
			p.world.RemoveConstraint(joint)
			delete(p.joints, k)
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
