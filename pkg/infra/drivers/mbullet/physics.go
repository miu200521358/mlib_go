//go:build windows
// +build windows

// 指示: miu200521358
package mbullet

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet/bt"
)

// PhysicsConfig は物理エンジンの設定パラメータ。
type PhysicsConfig struct {
	MaxSubSteps   int
	FixedTimeStep float32
}

// WindConfig は風のパラメータ設定。
type WindConfig struct {
	Enabled          bool
	Direction        mmath.Vec3
	Speed            float32
	Randomness       float32
	TurbulenceFreqHz float32
	DragCoeff        float32
	LiftCoeff        float32
	MaxAcceleration  float32
}

// RigidBodyValue は剛体の物理エンジン内部表現。
type RigidBodyValue struct {
	RigidBody        *model.RigidBody
	BtRigidBody      bt.BtRigidBody
	BtLocalTransform *bt.BtTransform
	Mask             int
	Group            int
	PrevBoneMatrix   mmath.Mat4
	HasPrevBone      bool
}

// PhysicsEngine は Bullet 物理エンジンの実装本体。
type PhysicsEngine struct {
	world       bt.BtDiscreteDynamicsWorld
	drawer      bt.BtMDebugDraw
	liner       *mDebugDrawLiner
	config      PhysicsConfig
	DeformSpf   float32
	PhysicsSpf  float32
	joints      map[int][]*jointValue
	rigidBodies map[int][]*RigidBodyValue
	windCfg     WindConfig
	simTimeAcc  float32
}

// NewPhysicsEngine は物理エンジンのインスタンスを生成する。
func NewPhysicsEngine(gravity *mmath.Vec3) *PhysicsEngine {
	gravityVec := mmath.ZERO_VEC3
	if gravity != nil {
		gravityVec = *gravity
	}
	world := createWorld(gravityVec)

	engine := &PhysicsEngine{
		world: world,
		config: PhysicsConfig{
			FixedTimeStep: 1 / 60.0,
		},
		rigidBodies: make(map[int][]*RigidBodyValue),
		joints:      make(map[int][]*jointValue),
		windCfg: WindConfig{
			Enabled:          false,
			Direction:        mmath.UNIT_X_VEC3,
			Speed:            0,
			Randomness:       0,
			TurbulenceFreqHz: 0.5,
			DragCoeff:        0.8,
			LiftCoeff:        0.2,
			MaxAcceleration:  80.0,
		},
		simTimeAcc: 0,
	}

	engine.initDebugDrawer()

	return engine
}

// initDebugDrawer はデバッグ描画機能を初期化する。
func (mp *PhysicsEngine) initDebugDrawer() {
	liner := newMDebugDrawLiner()
	drawer := bt.NewBtMDebugDraw()
	drawer.SetLiner(liner)
	drawer.SetMDefaultColors(newConstBtMDefaultColors())
	mp.world.SetDebugDrawer(drawer)

	mp.drawer = drawer
	mp.liner = liner
}

// ResetWorld はワールドを再構築する。
func (mp *PhysicsEngine) ResetWorld(gravity *mmath.Vec3) {
	gravityVec := mmath.ZERO_VEC3
	if gravity != nil {
		gravityVec = *gravity
	}
	bt.DeleteBtDynamicsWorld(mp.world)
	world := createWorld(gravityVec)
	world.SetDebugDrawer(mp.drawer)
	mp.world = world
	mp.resetFollowBoneCache()
}

// StepSimulation は物理シミュレーションを1ステップ進める。
func (mp *PhysicsEngine) StepSimulation(timeStep float32, maxSubSteps int, fixedTimeStep float32) {
	mp.applyWindForces(timeStep)
	mp.world.StepSimulation(timeStep, maxSubSteps, fixedTimeStep)
}

// AddModel はモデルを物理エンジンに追加する。
func (mp *PhysicsEngine) AddModel(modelIndex int, model *model.PmxModel) {
	if model == nil || model.RigidBodies == nil || model.Joints == nil {
		return
	}
	mp.initRigidBodies(modelIndex, model)
	mp.initJoints(modelIndex, model)
}

// AddModelByDeltas は差分情報を使用してモデルを物理エンジンに追加する。
func (mp *PhysicsEngine) AddModelByDeltas(
	modelIndex int,
	model *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	physicsDeltas *delta.PhysicsDeltas,
) {
	if model == nil || model.RigidBodies == nil || model.Joints == nil || boneDeltas == nil {
		return
	}
	var rigidBodyDeltas *delta.RigidBodyDeltas
	if physicsDeltas != nil {
		rigidBodyDeltas = physicsDeltas.RigidBodies
	}

	mp.initRigidBodiesByBoneDeltas(modelIndex, model, boneDeltas, rigidBodyDeltas)
	mp.initJointsByBoneDeltas(modelIndex, model, boneDeltas, nil)
}

// DeleteModel はモデルを物理エンジンから削除する。
func (mp *PhysicsEngine) DeleteModel(modelIndex int) {
	mp.deleteJoints(modelIndex)
	mp.deleteRigidBodies(modelIndex)
}

// clearModelFollowBoneCache はモデル単位の追従行列キャッシュを破棄する。
func (mp *PhysicsEngine) clearModelFollowBoneCache(modelIndex int) {
	bodies, ok := mp.rigidBodies[modelIndex]
	if !ok || bodies == nil {
		return
	}
	for _, body := range bodies {
		if body == nil {
			continue
		}
		body.PrevBoneMatrix = mmath.Mat4{}
		body.HasPrevBone = false
	}
}

// resetFollowBoneCache は全モデルの追従行列キャッシュを破棄する。
func (mp *PhysicsEngine) resetFollowBoneCache() {
	for modelIndex := range mp.rigidBodies {
		mp.clearModelFollowBoneCache(modelIndex)
	}
}

// createWorld は Bullet ワールドを生成する。
func createWorld(gravity mmath.Vec3) bt.BtDiscreteDynamicsWorld {
	broadphase := bt.NewBtDbvtBroadphase()
	collisionConfiguration := bt.NewBtDefaultCollisionConfiguration()
	dispatcher := bt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := bt.NewBtSequentialImpulseConstraintSolver()
	world := bt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	world.SetGravity(bt.NewBtVector3(float32(gravity.X), float32(gravity.Y*10), float32(gravity.Z)))

	groundShape := bt.NewBtStaticPlaneShape(bt.NewBtVector3(float32(0), float32(1), float32(0)), float32(0))
	groundTransform := bt.NewBtTransform()
	groundTransform.SetIdentity()
	groundTransform.SetOrigin(bt.NewBtVector3(float32(0), float32(0), float32(0)))
	groundMotionState := bt.NewBtDefaultMotionState(groundTransform)
	groundRigidBody := bt.NewBtRigidBody(float32(0), groundMotionState, groundShape)

	groundRigidBody.SetUserIndex(-2)
	groundRigidBody.SetUserIndex2(-2)

	world.AddRigidBody(groundRigidBody, 1<<15, math.MaxUint16)

	return world
}
