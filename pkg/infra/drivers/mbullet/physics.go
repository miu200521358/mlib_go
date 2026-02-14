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
	jointConfig JointConstraintConfig
	modelJoints map[int]JointConstraintConfig
	// FollowDeltaTransform で速度回転を許容する最大角度[rad]。
	// キーフレームの大ジャンプで速度まで回すとエネルギー注入が起きやすいため、しきい値で抑制する。
	followDeltaVelocityRotationMaxRad float64
	windCfg                           WindConfig
	simTimeAcc                        float32
}

// PhysicsStepDiagnostics は1ステップ分の物理診断サマリを保持する。
type PhysicsStepDiagnostics struct {
	ManifoldCount                 int
	ActiveManifoldCount           int
	ContactPointCount             int
	PenetrationPointCount         int
	MaxPenetrationDepth           float32
	MaxPenetrationModelIndexA     int
	MaxPenetrationRigidBodyIndexA int
	MaxPenetrationModelIndexB     int
	MaxPenetrationRigidBodyIndexB int
	MaxAppliedImpulse             float32
	MaxImpulseModelIndexA         int
	MaxImpulseRigidBodyIndexA     int
	MaxImpulseModelIndexB         int
	MaxImpulseRigidBodyIndexB     int
	DynamicBodyCount              int
	MaxLinearSpeed                float32
	MaxLinearSpeedModelIndex      int
	MaxLinearSpeedRigidBodyIndex  int
	MaxAngularSpeed               float32
	MaxAngularSpeedModelIndex     int
	MaxAngularSpeedRigidBodyIndex int
}

// TrackedRigidBodyDiagnostics は追跡対象剛体の属性と状態を保持する。
type TrackedRigidBodyDiagnostics struct {
	Reason                     string
	ModelIndex                 int
	RigidBodyIndex             int
	Resolved                   bool
	HasAppliedParams           bool
	PhysicsType                int
	Shape                      int
	Group                      int
	Mask                       int
	BaseMass                   float64
	AppliedMass                float64
	AppliedSize                mmath.Vec3
	AppliedPosition            mmath.Vec3
	LinearDamping              float32
	AngularDamping             float32
	Restitution                float32
	Friction                   float32
	CollisionFlags             int
	ActivationState            int
	CcdMotionThreshold         float32
	CcdSweptSphereRadius       float32
	ContactProcessingThreshold float32
	LinearSpeed                float32
	AngularSpeed               float32
}

// trackedRigidBodyTarget は剛体追跡対象を表す。
type trackedRigidBodyTarget struct {
	reason         string
	modelIndex     int
	rigidBodyIndex int
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
}

// StepSimulation は物理シミュレーションを1ステップ進める。
func (mp *PhysicsEngine) StepSimulation(timeStep float32, maxSubSteps int, fixedTimeStep float32) {
	mp.applyWindForces(timeStep)
	mp.world.StepSimulation(timeStep, maxSubSteps, fixedTimeStep)
}

// CollectStepDiagnostics は現在ステップの接触/速度サマリを返す。
func (mp *PhysicsEngine) CollectStepDiagnostics() PhysicsStepDiagnostics {
	diagnostics := newPhysicsStepDiagnostics()
	if mp == nil {
		return diagnostics
	}
	mp.collectContactDiagnostics(&diagnostics)
	mp.collectRigidBodyDiagnostics(&diagnostics)
	return diagnostics
}

// CollectTrackedRigidBodyDiagnostics はステップ診断で注目すべき剛体属性を返す。
func (mp *PhysicsEngine) CollectTrackedRigidBodyDiagnostics(
	stepDiagnostics PhysicsStepDiagnostics,
) []TrackedRigidBodyDiagnostics {
	targets := resolveTrackedRigidBodyTargets(stepDiagnostics)
	diagnostics := make([]TrackedRigidBodyDiagnostics, 0, len(targets))
	for _, target := range targets {
		diagnostics = append(diagnostics, mp.collectTrackedRigidBodyDiagnostic(target))
	}
	return diagnostics
}

// newPhysicsStepDiagnostics は診断サマリの初期値を生成する。
func newPhysicsStepDiagnostics() PhysicsStepDiagnostics {
	return PhysicsStepDiagnostics{
		MaxPenetrationModelIndexA:     -1,
		MaxPenetrationRigidBodyIndexA: -1,
		MaxPenetrationModelIndexB:     -1,
		MaxPenetrationRigidBodyIndexB: -1,
		MaxImpulseModelIndexA:         -1,
		MaxImpulseRigidBodyIndexA:     -1,
		MaxImpulseModelIndexB:         -1,
		MaxImpulseRigidBodyIndexB:     -1,
		MaxLinearSpeedModelIndex:      -1,
		MaxLinearSpeedRigidBodyIndex:  -1,
		MaxAngularSpeedModelIndex:     -1,
		MaxAngularSpeedRigidBodyIndex: -1,
	}
}

// newTrackedRigidBodyDiagnostics は追跡剛体の初期値を生成する。
func newTrackedRigidBodyDiagnostics(
	reason string,
	modelIndex int,
	rigidBodyIndex int,
) TrackedRigidBodyDiagnostics {
	return TrackedRigidBodyDiagnostics{
		Reason:                     reason,
		ModelIndex:                 modelIndex,
		RigidBodyIndex:             rigidBodyIndex,
		Resolved:                   false,
		HasAppliedParams:           false,
		PhysicsType:                -1,
		Shape:                      -1,
		Group:                      -1,
		Mask:                       -1,
		CollisionFlags:             -1,
		ActivationState:            -1,
		CcdMotionThreshold:         float32(-1),
		CcdSweptSphereRadius:       float32(-1),
		ContactProcessingThreshold: float32(-1),
	}
}

// resolveTrackedRigidBodyTargets は診断対象の剛体キーを列挙する。
func resolveTrackedRigidBodyTargets(stepDiagnostics PhysicsStepDiagnostics) []trackedRigidBodyTarget {
	targets := make([]trackedRigidBodyTarget, 0, 6)

	appendTarget := func(reason string, modelIndex int, rigidBodyIndex int) {
		if !shouldTrackRigidBody(modelIndex, rigidBodyIndex) {
			return
		}
		targets = append(targets, newTrackedRigidBodyTarget(reason, modelIndex, rigidBodyIndex))
	}

	appendTarget(
		"maxPenetrationA",
		stepDiagnostics.MaxPenetrationModelIndexA,
		stepDiagnostics.MaxPenetrationRigidBodyIndexA,
	)
	appendTarget(
		"maxPenetrationB",
		stepDiagnostics.MaxPenetrationModelIndexB,
		stepDiagnostics.MaxPenetrationRigidBodyIndexB,
	)
	appendTarget(
		"maxImpulseA",
		stepDiagnostics.MaxImpulseModelIndexA,
		stepDiagnostics.MaxImpulseRigidBodyIndexA,
	)
	appendTarget(
		"maxImpulseB",
		stepDiagnostics.MaxImpulseModelIndexB,
		stepDiagnostics.MaxImpulseRigidBodyIndexB,
	)
	appendTarget(
		"maxLinearSpeed",
		stepDiagnostics.MaxLinearSpeedModelIndex,
		stepDiagnostics.MaxLinearSpeedRigidBodyIndex,
	)
	appendTarget(
		"maxAngularSpeed",
		stepDiagnostics.MaxAngularSpeedModelIndex,
		stepDiagnostics.MaxAngularSpeedRigidBodyIndex,
	)

	return targets
}

// shouldTrackRigidBody は追跡対象として有効な剛体キーか判定する。
func shouldTrackRigidBody(modelIndex int, rigidBodyIndex int) bool {
	return modelIndex != -1 || rigidBodyIndex != -1
}

// newTrackedRigidBodyTarget は剛体追跡キーを生成する。
func newTrackedRigidBodyTarget(reason string, modelIndex int, rigidBodyIndex int) trackedRigidBodyTarget {
	return trackedRigidBodyTarget{
		reason:         reason,
		modelIndex:     modelIndex,
		rigidBodyIndex: rigidBodyIndex,
	}
}

// collectTrackedRigidBodyDiagnostic は追跡対象剛体の属性を収集する。
func (mp *PhysicsEngine) collectTrackedRigidBodyDiagnostic(
	target trackedRigidBodyTarget,
) TrackedRigidBodyDiagnostics {
	diagnostics := newTrackedRigidBodyDiagnostics(target.reason, target.modelIndex, target.rigidBodyIndex)
	if mp == nil {
		return diagnostics
	}
	body := mp.getRigidBodyValueByIndex(target.modelIndex, target.rigidBodyIndex)
	if body == nil || body.RigidBody == nil || body.BtRigidBody == nil {
		return diagnostics
	}

	diagnostics.Resolved = true
	diagnostics.HasAppliedParams = body.HasAppliedParams
	diagnostics.PhysicsType = int(body.RigidBody.PhysicsType)
	diagnostics.Shape = int(body.RigidBody.Shape)
	diagnostics.Group = body.Group
	diagnostics.Mask = body.Mask
	diagnostics.BaseMass = body.RigidBody.Param.Mass
	diagnostics.AppliedMass = body.AppliedMass
	diagnostics.AppliedSize = body.AppliedSize
	diagnostics.AppliedPosition = body.AppliedPosition
	diagnostics.LinearDamping = float32(body.RigidBody.Param.LinearDamping)
	diagnostics.AngularDamping = float32(body.RigidBody.Param.AngularDamping)
	diagnostics.Restitution = float32(body.RigidBody.Param.Restitution)
	diagnostics.Friction = float32(body.RigidBody.Param.Friction)
	diagnostics.CollisionFlags = body.BtRigidBody.GetCollisionFlags()
	diagnostics.ActivationState = body.BtRigidBody.GetActivationState()
	diagnostics.CcdMotionThreshold = body.BtRigidBody.GetCcdMotionThreshold()
	diagnostics.CcdSweptSphereRadius = body.BtRigidBody.GetCcdSweptSphereRadius()
	diagnostics.ContactProcessingThreshold = body.BtRigidBody.GetContactProcessingThreshold()

	linearVelocity := body.BtRigidBody.GetLinearVelocity()
	if linearVelocity != nil {
		diagnostics.LinearSpeed = linearVelocity.Length()
	}
	angularVelocity := body.BtRigidBody.GetAngularVelocity()
	if angularVelocity != nil {
		diagnostics.AngularSpeed = angularVelocity.Length()
	}

	return diagnostics
}

// collectContactDiagnostics は接触情報のサマリを集計する。
func (mp *PhysicsEngine) collectContactDiagnostics(diagnostics *PhysicsStepDiagnostics) {
	if mp == nil || diagnostics == nil || mp.worldRes == nil || mp.worldRes.dispatcher == nil {
		return
	}
	manifoldCount := mp.worldRes.dispatcher.GetNumManifolds()
	if manifoldCount <= 0 {
		return
	}
	diagnostics.ManifoldCount = manifoldCount
	for manifoldIndex := 0; manifoldIndex < manifoldCount; manifoldIndex++ {
		manifold := mp.worldRes.dispatcher.GetManifoldByIndexInternal(manifoldIndex)
		if manifold == nil {
			continue
		}
		contactCount := manifold.GetNumContacts()
		if contactCount <= 0 {
			continue
		}
		diagnostics.ActiveManifoldCount++
		bodyA := manifold.GetBody0()
		bodyB := manifold.GetBody1()
		modelIndexA, rigidBodyIndexA := resolveCollisionObjectIndexes(bodyA)
		modelIndexB, rigidBodyIndexB := resolveCollisionObjectIndexes(bodyB)
		for contactIndex := 0; contactIndex < contactCount; contactIndex++ {
			contactPointRaw := manifold.GetContactPoint(contactIndex)
			contactPoint, ok := contactPointRaw.(bt.BtManifoldPoint)
			if !ok || contactPoint == nil {
				continue
			}
			diagnostics.ContactPointCount++
			distance := contactPoint.GetDistance()
			if distance < 0 {
				penetrationDepth := -distance
				diagnostics.PenetrationPointCount++
				if penetrationDepth > diagnostics.MaxPenetrationDepth {
					diagnostics.MaxPenetrationDepth = penetrationDepth
					diagnostics.MaxPenetrationModelIndexA = modelIndexA
					diagnostics.MaxPenetrationRigidBodyIndexA = rigidBodyIndexA
					diagnostics.MaxPenetrationModelIndexB = modelIndexB
					diagnostics.MaxPenetrationRigidBodyIndexB = rigidBodyIndexB
				}
			}
			appliedImpulse := contactPoint.GetAppliedImpulse()
			if appliedImpulse > diagnostics.MaxAppliedImpulse {
				diagnostics.MaxAppliedImpulse = appliedImpulse
				diagnostics.MaxImpulseModelIndexA = modelIndexA
				diagnostics.MaxImpulseRigidBodyIndexA = rigidBodyIndexA
				diagnostics.MaxImpulseModelIndexB = modelIndexB
				diagnostics.MaxImpulseRigidBodyIndexB = rigidBodyIndexB
			}
		}
	}
}

// collectRigidBodyDiagnostics は剛体速度サマリを集計する。
func (mp *PhysicsEngine) collectRigidBodyDiagnostics(diagnostics *PhysicsStepDiagnostics) {
	if mp == nil || diagnostics == nil {
		return
	}
	for modelIndex, bodies := range mp.rigidBodies {
		for _, rigidBodyValue := range bodies {
			if rigidBodyValue == nil || rigidBodyValue.RigidBody == nil || rigidBodyValue.BtRigidBody == nil {
				continue
			}
			if rigidBodyValue.RigidBody.PhysicsType == model.PHYSICS_TYPE_STATIC {
				continue
			}
			diagnostics.DynamicBodyCount++
			linearVelocity := rigidBodyValue.BtRigidBody.GetLinearVelocity()
			if linearVelocity != nil {
				linearSpeed := linearVelocity.Length()
				if linearSpeed > diagnostics.MaxLinearSpeed {
					diagnostics.MaxLinearSpeed = linearSpeed
					diagnostics.MaxLinearSpeedModelIndex = modelIndex
					diagnostics.MaxLinearSpeedRigidBodyIndex = rigidBodyValue.RigidBody.Index()
				}
			}
			angularVelocity := rigidBodyValue.BtRigidBody.GetAngularVelocity()
			if angularVelocity != nil {
				angularSpeed := angularVelocity.Length()
				if angularSpeed > diagnostics.MaxAngularSpeed {
					diagnostics.MaxAngularSpeed = angularSpeed
					diagnostics.MaxAngularSpeedModelIndex = modelIndex
					diagnostics.MaxAngularSpeedRigidBodyIndex = rigidBodyValue.RigidBody.Index()
				}
			}
		}
	}
}

// resolveCollisionObjectIndexes は衝突オブジェクトからモデル/剛体インデックスを取得する。
func resolveCollisionObjectIndexes(collisionObject bt.BtCollisionObject) (int, int) {
	if collisionObject == nil {
		return -1, -1
	}
	return collisionObject.GetUserIndex(), collisionObject.GetUserIndex2()
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

// createWorld は Bullet ワールドを生成する。
func createWorld(gravity mmath.Vec3) bt.BtDiscreteDynamicsWorld {
	broadphase := bt.NewBtDbvtBroadphase()
	collisionConfiguration := bt.NewBtDefaultCollisionConfiguration()
	dispatcher := bt.NewBtCollisionDispatcher(collisionConfiguration)
	solver := bt.NewBtSequentialImpulseConstraintSolver()
	world := bt.NewBtDiscreteDynamicsWorld(dispatcher, broadphase, solver, collisionConfiguration)
	// MMD互換の重力スケールに合わせるため、Y成分は10倍でBulletへ渡す。
	gravityVector := bt.NewBtVector3(float32(gravity.X), float32(gravity.Y*10), float32(gravity.Z))
	world.SetGravity(gravityVector)
	bt.DeleteBtVector3(gravityVector)

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
