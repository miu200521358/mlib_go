// 指示: miu200521358
package physics

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// IPhysicsCore は物理エンジンのコア契約を表す。
type IPhysicsCore interface {
	// StepSimulation は物理シミュレーションを1ステップ進める。
	StepSimulation(timeStep float32, maxSubSteps int, fixedTimeStep float32)
	// ResetWorld はワールドを再構築する。
	ResetWorld(gravity *mmath.Vec3)

	// AddModel はモデルを物理エンジンに追加する。
	AddModel(modelIndex int, model *model.PmxModel)
	// AddModelByDeltas は差分情報を使用してモデルを追加する。
	AddModelByDeltas(modelIndex int, model *model.PmxModel, boneDeltas *delta.BoneDeltas, physicsDeltas *delta.PhysicsDeltas)
	// DeleteModel はモデルを物理エンジンから削除する。
	DeleteModel(modelIndex int)

	// UpdatePhysicsSelectively は差分に基づいて剛体/ジョイントを更新する。
	UpdatePhysicsSelectively(modelIndex int, model *model.PmxModel, physicsDeltas *delta.PhysicsDeltas)
	// UpdateRigidBodiesSelectively は剛体のみを差分更新する。
	UpdateRigidBodiesSelectively(modelIndex int, model *model.PmxModel, rigidBodyDeltas *delta.RigidBodyDeltas)
	// UpdateRigidBodyShapeMass は剛体の位置・形状・質量を更新する。
	UpdateRigidBodyShapeMass(modelIndex int, rigidBody *model.RigidBody, rigidBodyDelta *delta.RigidBodyDelta)
	// UpdateTransform は剛体の姿勢をボーン行列で更新する。
	UpdateTransform(modelIndex int, rigidBodyBone *model.Bone, boneGlobalMatrix *mmath.Mat4, rigidBody *model.RigidBody)
	// FollowDeltaTransform は前回ボーン姿勢との差分で剛体を追従更新する。
	FollowDeltaTransform(modelIndex int, rigidBodyBone *model.Bone, boneGlobalMatrix *mmath.Mat4, rigidBody *model.RigidBody)
	// GetRigidBodyBoneMatrix は剛体の姿勢からボーン行列を取得する。
	GetRigidBodyBoneMatrix(modelIndex int, rigidBody *model.RigidBody) *mmath.Mat4

	// EnableWind は風の有効/無効を切り替える。
	EnableWind(enable bool)
	// SetWind は風向き/風速/ランダム性を設定する。
	SetWind(direction *mmath.Vec3, speed float32, randomness float32)
	// SetWindAdvanced は風の詳細パラメータを設定する。
	SetWindAdvanced(dragCoeff, liftCoeff, turbulenceFreqHz float32)
}
