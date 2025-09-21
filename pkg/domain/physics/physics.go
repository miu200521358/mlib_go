package physics

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// PhysicsConfig 物理エンジンの設定パラメータ
type PhysicsConfig struct {
	MaxSubSteps   int     // 最大ステップ数
	FixedTimeStep float32 // 固定タイムステップ
}

// WindConfig 風のパラメータ設定
type WindConfig struct {
	Enabled          bool         // 風の有効/無効
	Direction        *mmath.MVec3 // MMD座標系の風向（正規化前でOK）
	Speed            float32      // 基本風速 [unit/s]
	Randomness       float32      // 乱れの強さ(0..1)
	TurbulenceFreqHz float32      // 乱流の周波数[Hz]
	DragCoeff        float32      // 抵抗係数（0.5*rho*Cd*A を吸収）
	LiftCoeff        float32      // 揚力係数（0.5*rho*Cl*A を吸収）
	MaxAcceleration  float32      // 風力による最大加速度の上限 [unit/s^2]
}

// IPhysics 物理エンジンのインターフェース
type IPhysics interface {
	// シミュレーション関連
	StepSimulation(timeStep float32, maxSubSteps int, fixedTimeStep float32) // 物理シミュレーションを1ステップ進める
	ResetWorld(gravity *mmath.MVec3)
	GetWorld() bt.BtDiscreteDynamicsWorld

	// モデル管理
	AddModel(modelIndex int, model *pmx.PmxModel)
	AddModelByDeltas(modelIndex int, model *pmx.PmxModel, boneDeltas *delta.BoneDeltas, physicsDeltas *delta.PhysicsDeltas)
	DeleteModel(modelIndex int)

	// トランスフォーム関連
	UpdatePhysicsSelectively(modelIndex int, model *pmx.PmxModel, physicsDeltas *delta.PhysicsDeltas)
	UpdateTransform(modelIndex int, rigidBodyBone *pmx.Bone, boneGlobalMatrix *mmath.MMat4, rigidBody *pmx.RigidBody)
	UpdateRigidBodyShapeMass(modelIndex int, rigidBody *pmx.RigidBody, rigidBodyDelta *delta.RigidBodyDelta)
	UpdateRigidBodiesSelectively(modelIndex int, model *pmx.PmxModel, rigidBodyDeltas *delta.RigidBodyDeltas)
	GetRigidBodyBoneMatrix(modelIndex int, rigidBody *pmx.RigidBody) *mmath.MMat4

	// デバッグ表示
	DrawDebugLines(shader rendering.IShader, visibleRigidBody, visibleJoint, isDrawRigidBodyFront bool)
	DrawDebugHighlight(shader rendering.IShader, isDrawRigidBodyFront bool)
	UpdateDebugHoverByRigidBody(modelIndex int, rigidBody *pmx.RigidBody, enable bool)
	DebugHoverInfo() *DebugRigidBodyHover

	// 風関連
	EnableWind(enable bool)
	SetWind(direction *mmath.MVec3, speed float32, randomness float32)
	SetWindAdvanced(dragCoeff, liftCoeff, turbulenceFreqHz float32)

	// 剛体関連
	FindRigidBodyByCollisionHit(hitObj bt.BtCollisionObject, hasHit bool) (modelIndex int, rb *pmx.RigidBody, ok bool)
}

// DebugRigidBodyHover holds information about the rigid body under the debug cursor.
type DebugRigidBodyHover struct {
	RigidBody *pmx.RigidBody
	HitPoint  *mmath.MVec3
}
