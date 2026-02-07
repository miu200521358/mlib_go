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

// JointConstraintConfig はジョイント拘束の共通設定を表す。
type JointConstraintConfig struct {
	ERP                                float32
	StopERP                            float32
	CFM                                float32
	StopCFM                            float32
	DisableCollisionsBetweenLinkedBody bool
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
	AppliedSize      mmath.Vec3
	AppliedMass      float64
	HasAppliedParams bool
}

// worldResources は Bullet ワールドの所有リソースを保持する。
type worldResources struct {
	world                  bt.BtDiscreteDynamicsWorld
	broadphase             bt.BtBroadphaseInterface
	collisionConfiguration bt.BtCollisionConfiguration
	dispatcher             bt.BtDispatcher
	solver                 bt.BtConstraintSolver
	groundShape            bt.BtCollisionShape
	groundMotionState      bt.BtMotionState
	groundRigidBody        bt.BtRigidBody
}

// PhysicsEngine は Bullet 物理エンジンの実装本体。
type PhysicsEngine struct {
	world       bt.BtDiscreteDynamicsWorld
	worldRes    *worldResources
	drawer      bt.BtMDebugDraw
	liner       *mDebugDrawLiner
	config      PhysicsConfig
	DeformSpf   float32
	PhysicsSpf  float32
	joints      map[int][]*jointValue
	rigidBodies map[int][]*RigidBodyValue
	jointConfig JointConstraintConfig
	modelJoints map[int]JointConstraintConfig
	// FollowDeltaTransform で速度回転を許容する最大角度[rad]。
	// キーフレームの大ジャンプで速度まで回すとエネルギー注入が起きやすいため、しきい値で抑制する。
	followDeltaVelocityRotationMaxRad float64
	windCfg                           WindConfig
	simTimeAcc                        float32
}

// NewPhysicsEngine は物理エンジンのインスタンスを生成する。
func NewPhysicsEngine(gravity *mmath.Vec3) *PhysicsEngine {
	gravityVec := mmath.ZERO_VEC3
	if gravity != nil {
		gravityVec = *gravity
	}
	worldRes := createWorld(gravityVec)

	engine := &PhysicsEngine{
		world:    worldRes.world,
		worldRes: worldRes,
		config: PhysicsConfig{
			FixedTimeStep: 1 / 60.0,
		},
		rigidBodies: make(map[int][]*RigidBodyValue),
		joints:      make(map[int][]*jointValue),
		// 既定では拘束を硬めにし、必要時のみモデル単位で上書きする。
		jointConfig: JointConstraintConfig{
			ERP:                                float32(0.8),
			StopERP:                            float32(0.8),
			CFM:                                float32(0.01),
			StopCFM:                            float32(0.01),
			DisableCollisionsBetweenLinkedBody: false,
		},
		modelJoints:                       make(map[int]JointConstraintConfig),
		followDeltaVelocityRotationMaxRad: defaultFollowDeltaVelocityRotationMaxRadians,
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
	mp.disposeAllModels()
	mp.disposeWorldResources()

	worldRes := createWorld(gravityVec)
	worldRes.world.SetDebugDrawer(mp.drawer)
	mp.world = worldRes.world
	mp.worldRes = worldRes
	mp.simTimeAcc = 0
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
	var jointDeltas *delta.JointDeltas
	if physicsDeltas != nil {
		rigidBodyDeltas = physicsDeltas.RigidBodies
		jointDeltas = physicsDeltas.Joints
	}

	mp.initRigidBodiesByBoneDeltas(modelIndex, model, boneDeltas, rigidBodyDeltas)
	mp.initJointsByBoneDeltas(modelIndex, model, boneDeltas, jointDeltas)
}

// DeleteModel はモデルを物理エンジンから削除する。
func (mp *PhysicsEngine) DeleteModel(modelIndex int) {
	mp.deleteJoints(modelIndex)
	mp.deleteRigidBodies(modelIndex)
	delete(mp.modelJoints, modelIndex)
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

// SetFollowDeltaVelocityRotationMaxRadians はFollowDeltaTransformの速度回転許容角度[rad]を設定する。
func (mp *PhysicsEngine) SetFollowDeltaVelocityRotationMaxRadians(maxAngleRad float64) {
	mp.followDeltaVelocityRotationMaxRad = clampFollowDeltaVelocityRotationMaxRadians(maxAngleRad)
}

// createWorld は Bullet ワールドを生成する。
func createWorld(gravity mmath.Vec3) *worldResources {
	broadphase := bt.NewBtDbvtBroadphase()
	collisionConfiguration := bt.NewBtDefaultCollisionConfiguration()
	dispatcher := bt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := bt.NewBtSequentialImpulseConstraintSolver()
	world := bt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	// MMD互換の重力スケールに合わせるため、Y成分は10倍でBulletへ渡す。
	gravityVector := bt.NewBtVector3(float32(gravity.X), float32(gravity.Y*10), float32(gravity.Z))
	world.SetGravity(gravityVector)
	bt.DeleteBtVector3(gravityVector)

	groundNormal := bt.NewBtVector3(float32(0), float32(1), float32(0))
	groundShape := bt.NewBtStaticPlaneShape(groundNormal, float32(0))
	bt.DeleteBtVector3(groundNormal)
	groundTransform := bt.NewBtTransform()
	groundTransform.SetIdentity()
	groundOrigin := bt.NewBtVector3(float32(0), float32(0), float32(0))
	groundTransform.SetOrigin(groundOrigin)
	bt.DeleteBtVector3(groundOrigin)
	groundMotionState := bt.NewBtDefaultMotionState(groundTransform)
	bt.DeleteBtTransform(groundTransform)
	groundRigidBody := bt.NewBtRigidBody(float32(0), groundMotionState, groundShape)

	groundRigidBody.SetUserIndex(-2)
	groundRigidBody.SetUserIndex2(-2)

	world.AddRigidBody(groundRigidBody, 1<<15, math.MaxUint16)

	return &worldResources{
		world:                  world,
		broadphase:             broadphase,
		collisionConfiguration: collisionConfiguration,
		dispatcher:             dispatcher,
		solver:                 solver,
		groundShape:            groundShape,
		groundMotionState:      groundMotionState,
		groundRigidBody:        groundRigidBody,
	}
}

// disposeAllModels は登録済みモデルの剛体・ジョイントを全て解放する。
func (mp *PhysicsEngine) disposeAllModels() {
	processed := make(map[int]struct{})
	for modelIndex := range mp.rigidBodies {
		mp.deleteJoints(modelIndex)
		mp.deleteRigidBodies(modelIndex)
		delete(mp.modelJoints, modelIndex)
		processed[modelIndex] = struct{}{}
	}
	for modelIndex := range mp.joints {
		if _, exists := processed[modelIndex]; exists {
			continue
		}
		mp.deleteJoints(modelIndex)
		delete(mp.modelJoints, modelIndex)
	}
}

// disposeWorldResources はワールド構成要素を破棄する。
func (mp *PhysicsEngine) disposeWorldResources() {
	if mp.worldRes == nil {
		if mp.world != nil {
			bt.DeleteBtDiscreteDynamicsWorld(mp.world)
		}
		mp.world = nil
		return
	}
	resources := mp.worldRes
	if resources.world != nil && resources.groundRigidBody != nil {
		resources.world.RemoveRigidBody(resources.groundRigidBody)
	}
	if resources.groundRigidBody != nil {
		bt.DeleteBtRigidBody(resources.groundRigidBody)
		resources.groundRigidBody = nil
	}
	if resources.groundMotionState != nil {
		bt.DeleteBtMotionState(resources.groundMotionState)
		resources.groundMotionState = nil
	}
	if resources.groundShape != nil {
		bt.DeleteBtCollisionShape(resources.groundShape)
		resources.groundShape = nil
	}
	if resources.world != nil {
		bt.DeleteBtDiscreteDynamicsWorld(resources.world)
		resources.world = nil
	}
	if resources.solver != nil {
		bt.DeleteBtConstraintSolver(resources.solver)
		resources.solver = nil
	}
	if resources.dispatcher != nil {
		bt.DeleteBtDispatcher(resources.dispatcher)
		resources.dispatcher = nil
	}
	if resources.collisionConfiguration != nil {
		bt.DeleteBtCollisionConfiguration(resources.collisionConfiguration)
		resources.collisionConfiguration = nil
	}
	if resources.broadphase != nil {
		bt.DeleteBtBroadphaseInterface(resources.broadphase)
		resources.broadphase = nil
	}
	mp.world = nil
	mp.worldRes = nil
}
