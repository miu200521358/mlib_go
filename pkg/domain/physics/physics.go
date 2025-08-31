package physics

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// PhysicsConfig 物理エンジンの設定パラメータ
type PhysicsConfig struct {
	MaxSubSteps   int     // 最大ステップ数
	FixedTimeStep float32 // 固定タイムステップ
}

// IPhysics 物理エンジンのインターフェース
type IPhysics interface {
	// シミュレーション関連
	StepSimulation(timeStep float32, maxSubSteps int, fixedTimeStep float32) // 物理シミュレーションを1ステップ進める
	ResetWorld(gravity *mmath.MVec3)

	// モデル管理
	AddModel(modelIndex int, model *pmx.PmxModel)
	AddModelByDeltas(modelIndex int, model *pmx.PmxModel, boneDeltas *delta.BoneDeltas, physicsDeltas *delta.PhysicsDeltas)
	DeleteModel(modelIndex int)

	// トランスフォーム関連
	UpdateTransform(modelIndex int, rigidBodyBone *pmx.Bone, boneGlobalMatrix *mmath.MMat4, rigidBody *pmx.RigidBody)
	UpdateRigidBodyShape(modelIndex int, rigidBody *pmx.RigidBody, rigidBodyDelta *delta.RigidBodyDelta)
	UpdateRigidBodiesSelectively(modelIndex int, model *pmx.PmxModel, rigidBodyDeltas *delta.RigidBodyDeltas)
	GetRigidBodyBoneMatrix(modelIndex int, rigidBody *pmx.RigidBody) *mmath.MMat4

	// デバッグ表示
	DrawDebugLines(shader rendering.IShader, visibleRigidBody, visibleJoint, isDrawRigidBodyFront bool)
}
