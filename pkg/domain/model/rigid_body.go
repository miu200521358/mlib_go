// 指示: miu200521358
package model

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// CollisionGroup は剛体の衝突設定を表す。
type CollisionGroup struct {
	Group byte
	Mask  uint16
}

// Shape は剛体形状を表す。
type Shape int

const (
	// SHAPE_NONE は形状なし。
	SHAPE_NONE Shape = -1
	// SHAPE_SPHERE は球。
	SHAPE_SPHERE Shape = 0
	// SHAPE_BOX は箱。
	SHAPE_BOX Shape = 1
	// SHAPE_CAPSULE はカプセル。
	SHAPE_CAPSULE Shape = 2
)

// PhysicsType は剛体物理の種類を表す。
type PhysicsType int

const (
	// PHYSICS_TYPE_STATIC は静的物理。
	PHYSICS_TYPE_STATIC PhysicsType = iota
	// PHYSICS_TYPE_DYNAMIC は動的物理。
	PHYSICS_TYPE_DYNAMIC
	// PHYSICS_TYPE_DYNAMIC_BONE は動的物理+ボーン追従。
	PHYSICS_TYPE_DYNAMIC_BONE
)

// RigidBodyParam は剛体パラメータを表す。
type RigidBodyParam struct {
	Mass           float64
	LinearDamping  float64
	AngularDamping float64
	Restitution    float64
	Friction       float64
}

// RigidBody は剛体要素を表す。
type RigidBody struct {
	index          int
	name           string
	EnglishName    string
	BoneIndex      int
	CollisionGroup CollisionGroup
	Shape          Shape
	Size           mmath.Vec3
	Position       mmath.Vec3
	Rotation       mmath.Vec3
	Param          RigidBodyParam
	PhysicsType    PhysicsType
	IsSystem       bool // システム追加剛体の場合はtrue
}

// Index は剛体 index を返す。
func (r *RigidBody) Index() int {
	return r.index
}

// SetIndex は剛体 index を設定する。
func (r *RigidBody) SetIndex(index int) {
	r.index = index
}

// Name は剛体名を返す。
func (r *RigidBody) Name() string {
	return r.name
}

// SetName は剛体名を設定する。
func (r *RigidBody) SetName(name string) {
	r.name = name
}

// IsValid は剛体が有効か判定する。
func (r *RigidBody) IsValid() bool {
	return r != nil && r.index >= 0
}
