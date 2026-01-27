// 指示: miu200521358
package physics_api

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/physics"
)

// IPhysicsCore は物理エンジンのコア契約を表す。
type IPhysicsCore = physics.IPhysicsCore

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
