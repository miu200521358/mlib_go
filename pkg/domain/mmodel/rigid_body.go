package mmodel

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

// Shape は剛体の形状を表します。
type Shape int

const (
	SHAPE_NONE    Shape = -1 // なし
	SHAPE_SPHERE  Shape = 0  // 球
	SHAPE_BOX     Shape = 1  // 箱
	SHAPE_CAPSULE Shape = 2  // カプセル
)

// PhysicsType は剛体物理の計算モードを表します。
type PhysicsType int

const (
	PHYSICS_TYPE_STATIC       PhysicsType = 0 // ボーン追従(static)
	PHYSICS_TYPE_DYNAMIC      PhysicsType = 1 // 物理演算(dynamic)
	PHYSICS_TYPE_DYNAMIC_BONE PhysicsType = 2 // 物理演算 + Bone位置合わせ
)

// CollisionGroupFlags は衝突グループフラグです。
var CollisionGroupFlags = []uint16{
	0x0001, 0x0002, 0x0004, 0x0008,
	0x0010, 0x0020, 0x0040, 0x0080,
	0x0100, 0x0200, 0x0400, 0x0800,
	0x1000, 0x2000, 0x4000, 0x8000,
}

// CollisionGroup は衝突グループを表します。
type CollisionGroup struct {
	IsCollisions []uint16 // 衝突判定フラグ（16グループ）
}

// Value は衝突グループのビットマスク値を返します。
func (cg CollisionGroup) Value() int {
	var value uint16 = 0
	for i, v := range cg.IsCollisions {
		if v == 1 {
			value |= CollisionGroupFlags[i]
		}
	}
	return int(value)
}

// NewCollisionGroupFromSlice はスライスから衝突グループを生成します。
func NewCollisionGroupFromSlice(collisionGroup []uint16) CollisionGroup {
	groups := CollisionGroup{}
	collisionGroupMask := uint16(0)
	for i, v := range collisionGroup {
		if v == 1 {
			collisionGroupMask |= CollisionGroupFlags[i]
		}
	}
	groups.IsCollisions = NewCollisionGroupMask(collisionGroupMask)
	return groups
}

// NewCollisionGroupAll は全グループ衝突を生成します。
func NewCollisionGroupAll() CollisionGroup {
	groups := CollisionGroup{}
	var collisionGroupMask uint16 = 0xFFFF
	groups.IsCollisions = NewCollisionGroupMask(collisionGroupMask)
	return groups
}

// NewCollisionGroupMask はマスク値から衝突グループ配列を生成します。
func NewCollisionGroupMask(collisionGroupMask uint16) []uint16 {
	collisionGroup := make([]uint16, 16)
	for i, v := range CollisionGroupFlags {
		if collisionGroupMask&v == v {
			collisionGroup[i] = 0
		} else {
			collisionGroup[i] = 1
		}
	}
	return collisionGroup
}

// RigidBodyParam は剛体パラメータを表します。
type RigidBodyParam struct {
	Mass           float64 // 質量
	LinearDamping  float64 // 移動減衰
	AngularDamping float64 // 回転減衰
	Restitution    float64 // 反発力
	Friction       float64 // 摩擦力
}

// NewRigidBodyParam は新しい剛体パラメータを生成します。
func NewRigidBodyParam() *RigidBodyParam {
	return &RigidBodyParam{
		Mass:           1,
		LinearDamping:  0.5,
		AngularDamping: 0.5,
		Restitution:    0,
		Friction:       0.5,
	}
}

// String は文字列表現を返します。
func (p *RigidBodyParam) String() string {
	return fmt.Sprintf("Mass: %.5f, LinearDamping: %.5f, AngularDamping: %.5f, Restitution: %.5f, Friction: %.5f",
		p.Mass, p.LinearDamping, p.AngularDamping, p.Restitution, p.Friction)
}

// RigidBody は剛体を表します。
type RigidBody struct {
	mcore.IndexNameModel
	BoneIndex               int             // 関連ボーンIndex
	CollisionGroup          byte            // グループ
	CollisionGroupMask      CollisionGroup  // 非衝突グループフラグ
	CollisionGroupMaskValue int             // 非衝突グループフラグ値
	ShapeType               Shape           // 形状
	Size                    *mmath.Vec3     // サイズ
	Position                *mmath.Vec3     // 位置
	Rotation                *mmath.Vec3     // 回転（ラジアン）
	RigidBodyParam          *RigidBodyParam // 剛体パラメータ
	PhysicsType             PhysicsType     // 物理演算モード
	XDirection              *mmath.Vec3     // X軸方向
	YDirection              *mmath.Vec3     // Y軸方向
	ZDirection              *mmath.Vec3     // Z軸方向
	IsSystem                bool            // システム追加剛体
	Matrix                  *mmath.Mat4     // 剛体行列
	Bone                    *Bone           // 関連ボーン参照（Setup後）
}

// NewRigidBody は新しい剛体を生成します。
func NewRigidBody() *RigidBody {
	return &RigidBody{
		IndexNameModel:          *mcore.NewIndexNameModel(-1, "", ""),
		BoneIndex:               -1,
		CollisionGroup:          0,
		CollisionGroupMask:      NewCollisionGroupAll(),
		CollisionGroupMaskValue: 0,
		ShapeType:               SHAPE_BOX,
		Size:                    mmath.NewVec3(),
		Position:                mmath.NewVec3(),
		Rotation:                mmath.NewVec3(),
		RigidBodyParam:          NewRigidBodyParam(),
		PhysicsType:             PHYSICS_TYPE_STATIC,
		XDirection:              mmath.NewVec3(),
		YDirection:              mmath.NewVec3(),
		ZDirection:              mmath.NewVec3(),
		IsSystem:                false,
		Matrix:                  mmath.NewMat4(),
		Bone:                    nil,
	}
}

// IsValid は剛体が有効かどうかを返します。
func (r *RigidBody) IsValid() bool {
	return r != nil && r.Index() >= 0
}

// AsDynamic は物理演算対象かどうかを返します。
func (r *RigidBody) AsDynamic() bool {
	return r.PhysicsType == PHYSICS_TYPE_DYNAMIC || r.PhysicsType == PHYSICS_TYPE_DYNAMIC_BONE
}

// Copy は深いコピーを作成します。
func (r *RigidBody) Copy() (*RigidBody, error) {
	cp := &RigidBody{}
	if err := deepcopy.Copy(cp, r); err != nil {
		return nil, err
	}
	// Boneは参照のためnilにする
	cp.Bone = nil
	return cp, nil
}
