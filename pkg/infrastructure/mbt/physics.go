//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// PhysicsConfig 物理エンジンの設定パラメータ
type PhysicsConfig struct {
	MaxSubSteps   int     // 最大ステップ数
	FixedTimeStep float32 // 固定タイムステップ
}

// MPhysics 物理エンジンの実装
type MPhysics struct {
	world       bt.BtDiscreteDynamicsWorld // ワールド
	drawer      bt.BtMDebugDraw            // デバッグビューワー
	liner       *mDebugDrawLiner           // ライナー
	config      PhysicsConfig              // 設定パラメータ
	DeformSpf   float32                    // デフォームspf
	PhysicsSpf  float32                    // 物理spf
	joints      map[int][]*jointValue      // ジョイント
	rigidBodies map[int][]*rigidbodyValue  // 剛体
}

// NewMPhysics は物理エンジンのインスタンスを生成します
func NewMPhysics() *MPhysics {
	world := createWorld()

	// デフォルト設定
	config := PhysicsConfig{
		MaxSubSteps:   2,
		FixedTimeStep: 1 / 60.0,
	}

	physics := &MPhysics{
		world:       world,
		config:      config,
		rigidBodies: make(map[int][]*rigidbodyValue),
		joints:      make(map[int][]*jointValue),
	}

	// デバッグビューワーの初期化
	physics.initDebugDrawer()

	return physics
}

// initDebugDrawer はデバッグ描画機能を初期化します
func (physics *MPhysics) initDebugDrawer() {
	liner := newMDebugDrawLiner()
	drawer := bt.NewBtMDebugDraw()
	drawer.SetLiner(liner)
	drawer.SetMDefaultColors(newConstBtMDefaultColors())
	physics.world.SetDebugDrawer(drawer)

	physics.drawer = drawer
	physics.liner = liner
}

// ResetWorld はワールドをリセットします
func (physics *MPhysics) ResetWorld() {
	// ワールド削除
	bt.DeleteBtDynamicsWorld(physics.world)
	// ワールド作成
	world := createWorld()
	world.SetDebugDrawer(physics.drawer)
	physics.world = world
}

// AddModel はモデルを物理エンジンに追加します
func (physics *MPhysics) AddModel(modelIndex int, model *pmx.PmxModel) {
	// 根元から追加していく
	physics.initRigidBodies(modelIndex, model.RigidBodies)
	physics.initJoints(modelIndex, model.RigidBodies, model.Joints)
}

// AddModelByBoneDeltas はボーンデルタ情報を使用してモデルを物理エンジンに追加します
func (physics *MPhysics) AddModelByBoneDeltas(modelIndex int, model *pmx.PmxModel, boneDeltas *delta.BoneDeltas) {
	// 根元から追加していく
	physics.initRigidBodiesByBoneDeltas(modelIndex, model.RigidBodies, boneDeltas)
	physics.initJointsByBoneDeltas(modelIndex, model.RigidBodies, model.Joints, boneDeltas)
}

// DeleteModel はモデルを物理エンジンから削除します
func (physics *MPhysics) DeleteModel(modelIndex int) {
	// 末端から削除していく
	physics.deleteJoints(modelIndex)
	physics.deleteRigidBodies(modelIndex)
}

// StepSimulation は物理シミュレーションを1ステップ進めます
func (physics *MPhysics) StepSimulation(timeStep float32) {
	physics.world.StepSimulation(timeStep, physics.config.MaxSubSteps, physics.config.FixedTimeStep)
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
