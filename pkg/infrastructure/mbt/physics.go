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
	liner         *mDebugDrawLiner
	MaxSubSteps   int
	DeformFps     float32
	DeformSpf     float32
	PhysicsFps    float32
	PhysicsSpf    float32
	FixedTimeStep float32
	joints        map[int][]*jointValue
	rigidBodies   map[int][]*rigidbodyValue
}

func NewMPhysics() *MPhysics {
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
	liner := newMDebugDrawLiner()
	drawer := bt.NewBtMDebugDraw()
	drawer.SetLiner(liner)
	drawer.SetMDefaultColors(newConstBtMDefaultColors())
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
		p.DeleteModel(modelIndex)
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

func (physics *MPhysics) DeleteModel(modelIndex int) {
	// pm.physics = physics
	physics.deleteRigidBodies(modelIndex)
	physics.deleteJoints(modelIndex)
}

func (p *MPhysics) DrawDebugLines(shader *mgl.MShader, visibleRigidBody, visibleJoint, isDrawRigidBodyFront bool) {
	// 物理デバッグ取得
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

	// デバッグ情報取得
	p.world.DebugDrawWorld()

	// デバッグ描画
	p.liner.drawDebugLines(shader, isDrawRigidBodyFront)
}

func (p *MPhysics) StepSimulation(timeStep float32) {
	p.world.StepSimulation(timeStep, p.MaxSubSteps, p.FixedTimeStep)
}
