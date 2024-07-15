//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
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
	joints        map[int][]*jointValue
	rigidBodies   map[int][]*rigidbodyValue
}

func NewMPhysics(shader *mgl.MShader) *MPhysics {
	world := createWorld()

	p := &MPhysics{
		world:       world,
		MaxSubSteps: 2,
		DeformFps:   30.0,
		PhysicsFps:  60.0,
		rigidBodies: make(map[int][]*rigidbodyValue),
		joints:      make(map[int][]*jointValue),
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
	// 一旦削除
	for modelIndex := range p.rigidBodies {
		p.DeleteRigidBodies(modelIndex)
	}
	for modelIndex := range p.joints {
		p.DeleteJoints(modelIndex)
	}
	// 登録し直し
	for modelIndex := range p.rigidBodies {
		for _, r := range p.rigidBodies[modelIndex] {
			if r != nil {
				p.initRigidBody(modelIndex, r.pmxRigidBody)
			}
		}
	}
	for modelIndex := range p.joints {
		for _, j := range p.joints[modelIndex] {
			if j != nil {
				p.initJoint(modelIndex, j.pmxJoint)
			}
		}
	}
}

func (physics *MPhysics) AddModel(modelIndex int, model *pmx.PmxModel) {
	// pm.physics = physics
	physics.initRigidBodies(modelIndex, model.RigidBodies)
	physics.initJoints(modelIndex, model.RigidBodies, model.Joints)
}

func (p *MPhysics) DrawDebugLines(visibleRigidBody, visibleJoint, isDrawRigidBodyFront bool) {

	// 物理デバッグ表示
	p.DebugDrawWorld(visibleRigidBody, visibleJoint)

	// // 標準出力を一時的にリダイレクトする
	// old := os.Stdout // keep backup of the real stdout
	// r, w, _ := os.Pipe()
	// os.Stdout = w
	p.liner.DrawDebugLines(isDrawRigidBodyFront)
	p.liner.vertices = []float32{}

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
	r := p.rigidBodies[modelIndex][rigidBodyIndex]
	return *r.btRigidBody, *r.btLocalTransform
}

func (p *MPhysics) DeleteRigidBodies(modelIndex int) {
	for _, r := range p.rigidBodies[modelIndex] {
		p.world.RemoveRigidBody(*r.btRigidBody)
	}
}

func (p *MPhysics) DeleteJoints(modelIndex int) {
	for _, j := range p.joints[modelIndex] {
		p.world.RemoveConstraint(j.btJoint)
	}
}

func (p *MPhysics) StepSimulation(timeStep float32) {
	p.world.StepSimulation(timeStep, p.MaxSubSteps, p.FixedTimeStep)
}
