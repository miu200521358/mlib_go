// 指示: miu200521358
package physics_api

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

const (
	PhysicsDefaultMaxSubSteps = 5 // 物理シミュレーションのデフォルト最大サブステップ数
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
	// UpdateRigidBodyShapeMass は剛体の形状・質量を更新する。
	UpdateRigidBodyShapeMass(modelIndex int, rigidBody *model.RigidBody, rigidBodyDelta *delta.RigidBodyDelta)
	// UpdateTransform は剛体の姿勢をボーン行列で更新する。
	UpdateTransform(modelIndex int, rigidBodyBone *model.Bone, boneGlobalMatrix *mmath.Mat4, rigidBody *model.RigidBody)
	// GetRigidBodyBoneMatrix は剛体の姿勢からボーン行列を取得する。
	GetRigidBodyBoneMatrix(modelIndex int, rigidBody *model.RigidBody) *mmath.Mat4

	// EnableWind は風の有効/無効を切り替える。
	EnableWind(enable bool)
	// SetWind は風向き/風速/ランダム性を設定する。
	SetWind(direction *mmath.Vec3, speed float32, randomness float32)
	// SetWindAdvanced は風の詳細パラメータを設定する。
	SetWindAdvanced(dragCoeff, liftCoeff, turbulenceFreqHz float32)
}

// IPhysicsDebug は物理デバッグ描画の契約を表す。
type IPhysicsDebug interface {
	// DrawDebugLines は物理デバッグ線を描画する。
	DrawDebugLines(shader graphics_api.IShader, visibleRigidBody, visibleJoint, isDrawRigidBodyFront bool)
}

// IPhysicsRaycast は物理レイキャストの契約を表す。
type IPhysicsRaycast interface {
	// RayTest はレイキャスト結果を取得する。
	RayTest(rayFrom, rayTo *mmath.Vec3, filter *RaycastFilter) *RaycastHit
}

// IPhysics は物理APIの統合インターフェースを表す。
type IPhysics interface {
	IPhysicsCore
	IPhysicsDebug
	IPhysicsRaycast
}

// RaycastFilter は衝突フィルタの条件を表す。
type RaycastFilter struct {
	Group int
	Mask  int
}

// RigidBodyRef は剛体を識別する参照を表す。
type RigidBodyRef struct {
	ModelIndex     int
	RigidBodyIndex int
}

// RaycastHit はレイキャスト結果を表す。
type RaycastHit struct {
	RigidBodyRef
	HitFraction float32
}

// PhysicsFactory は物理エンジン生成関数を表す。
type PhysicsFactory func(gravity *mmath.Vec3) IPhysics

var physicsFactory PhysicsFactory

// RegisterPhysicsFactory は物理エンジン生成関数を登録する。
func RegisterPhysicsFactory(factory PhysicsFactory) {
	physicsFactory = factory
}

// NewPhysics は登録済みの物理エンジンを生成する。
func NewPhysics(gravity *mmath.Vec3) IPhysics {
	if physicsFactory == nil {
		return nil
	}
	return physicsFactory(gravity)
}
