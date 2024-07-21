//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

type IPhysics interface {
	ResetWorld()
	AddModel(modelIndex int, model *pmx.PmxModel)
	DeleteModel(modelIndex int)
	StepSimulation(timeStep float32)
	UpdateTransform(modelIndex int, rigidBodyBone *pmx.Bone, boneGlobalMatrix *mmath.MMat4, r *pmx.RigidBody)
	GetRigidBodyBoneMatrix(modelIndex int, rigidBody *pmx.RigidBody) *mmath.MMat4
	Exists(modelIndex int) bool
}

type MPhysics struct {
	world         bt.BtDiscreteDynamicsWorld // ワールド
	drawer        bt.BtMDebugDraw            // デバッグビューワー
	liner         *mDebugDrawLiner           // ライナー
	MaxSubSteps   int                        // 最大ステップ数
	DeformFps     float32                    // デフォームfps
	DeformSpf     float32                    // デフォームspf
	PhysicsFps    float32                    // 物理fps
	PhysicsSpf    float32                    // 物理spf
	FixedTimeStep float32                    // 固定タイムステップ
	joints        map[int][]*jointValue      // ジョイント
	rigidBodies   map[int][]*rigidbodyValue  // 剛体
}

func NewMPhysics() *MPhysics {
	world := createWorld()

	physics := &MPhysics{
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

	physics.drawer = drawer
	physics.liner = liner
	physics.DeformSpf = 1.0 / physics.DeformFps
	physics.PhysicsSpf = 1.0 / physics.PhysicsFps
	physics.FixedTimeStep = 1 / 60.0

	return physics
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

func (physics *MPhysics) ResetWorld() {
	world := createWorld()
	world.SetDebugDrawer(physics.drawer)
	physics.world = world
}

func (physics *MPhysics) AddModel(modelIndex int, model *pmx.PmxModel) {
	physics.initRigidBodies(modelIndex, model.RigidBodies)
	physics.initJoints(modelIndex, model.RigidBodies, model.Joints)
}

func (physics *MPhysics) DeleteModel(modelIndex int) {
	physics.deleteRigidBodies(modelIndex)
	physics.deleteJoints(modelIndex)
}

func (physics *MPhysics) StepSimulation(timeStep float32) {
	physics.world.StepSimulation(timeStep, physics.MaxSubSteps, physics.FixedTimeStep)
}

func (physics *MPhysics) Exists(modelIndex int) bool {
	_, ok := physics.rigidBodies[modelIndex]
	return ok
}
