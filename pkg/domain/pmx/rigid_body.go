package pmx

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type RigidBodyParam struct {
	Mass           float64 // 質量
	LinearDamping  float64 // 移動減衰
	AngularDamping float64 // 回転減衰
	Restitution    float64 // 反発力
	Friction       float64 // 摩擦力
}

func NewRigidBodyParam() *RigidBodyParam {
	return &RigidBodyParam{
		Mass:           1,
		LinearDamping:  0.5,
		AngularDamping: 0.5,
		Restitution:    0,
		Friction:       0.5,
	}
}

func (param *RigidBodyParam) String() string {
	return fmt.Sprintf("Mass: %.5f, LinearDamping: %.5f, AngularDamping: %.5f, Restitution: %.5f, Friction: %.5f",
		param.Mass, param.LinearDamping, param.AngularDamping, param.Restitution, param.Friction)
}

// 剛体の形状
type Shape int

const (
	SHAPE_NONE    Shape = -1 // なし
	SHAPE_SPHERE  Shape = 0  // 球
	SHAPE_BOX     Shape = 1  // 箱
	SHAPE_CAPSULE Shape = 2  // カプセル
)

// 剛体物理の計算モード
type PhysicsType int

const (
	PHYSICS_TYPE_STATIC       PhysicsType = 0 // ボーン追従(static)
	PHYSICS_TYPE_DYNAMIC      PhysicsType = 1 // 物理演算(dynamic)
	PHYSICS_TYPE_DYNAMIC_BONE PhysicsType = 2 // 物理演算 + Bone位置合わせ
)

type CollisionGroup struct {
	IsCollisions []uint16
}

func (cg CollisionGroup) Value() int {
	var value uint16 = 0
	for i, v := range cg.IsCollisions {
		if v == 1 {
			value |= CollisionGroupFlags[i]
		}
	}
	return int(value)
}

var CollisionGroupFlags = []uint16{
	0x0001, // 0:グループ1
	0x0002, // 1:グループ2
	0x0004, // 2:グループ3
	0x0008, // 3:グループ4
	0x0010, // 4:グループ5
	0x0020, // 5:グループ6
	0x0040, // 6:グループ7
	0x0080, // 7:グループ8
	0x0100, // 8:グループ9
	0x0200, // 9:グループ10
	0x0400, // 10:グループ11
	0x0800, // 11:グループ12
	0x1000, // 12:グループ13
	0x2000, // 13:グループ14
	0x4000, // 14:グループ15
	0x8000, // 15:グループ16
}

func NewCollisionGroupFromSlice(collisionGroup []uint16) CollisionGroup {
	groups := CollisionGroup{}
	collisionGroupMask := uint16(0)
	for i, v := range collisionGroup {
		if v == 1 {
			collisionGroupMask |= CollisionGroupFlags[i]
		}
	}
	groups.IsCollisions = NewCollisionGroup(collisionGroupMask)

	return groups
}

func NewCollisionGroupAll() CollisionGroup {
	groups := CollisionGroup{}
	collisionGroupMask := uint16(0)
	for i, v := range []uint16{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} {
		if v == 1 {
			collisionGroupMask |= CollisionGroupFlags[i]
		}
	}
	groups.IsCollisions = NewCollisionGroup(collisionGroupMask)

	return groups
}

func NewCollisionGroup(collisionGroupMask uint16) []uint16 {
	collisionGroup := []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i, v := range CollisionGroupFlags {
		if collisionGroupMask&v == v {
			collisionGroup[i] = 0
		} else {
			collisionGroup[i] = 1
		}
	}
	return collisionGroup
}

type RigidBody struct {
	index                   int             // 剛体INDEX
	name                    string          // 剛体名
	englishName             string          // 剛体英名
	BoneIndex               int             // 関連ボーンIndex
	CollisionGroup          byte            // グループ
	CollisionGroupMask      CollisionGroup  // 非衝突グループフラグ
	CollisionGroupMaskValue int             // 非衝突グループフラグ値
	ShapeType               Shape           // 形状
	Size                    *mmath.MVec3    // サイズ(x,y,z)
	Position                *mmath.MVec3    // 位置(x,y,z)
	Rotation                *mmath.MVec3    // 回転(x,y,z) -> ラジアン角
	RigidBodyParam          *RigidBodyParam // 剛体パラ
	PhysicsType             PhysicsType     // 剛体の物理演算
	XDirection              *mmath.MVec3    // X軸方向
	YDirection              *mmath.MVec3    // Y軸方向
	ZDirection              *mmath.MVec3    // Z軸方向
	IsSystem                bool            // システムで追加した剛体か
	Matrix                  *mmath.MMat4    // 剛体の行列
	Bone                    *Bone           // 繋がっているボーン
}

// NewRigidBody creates a new rigid body.
func NewRigidBody() *RigidBody {
	return &RigidBody{
		index:                   -1,
		name:                    "",
		englishName:             "",
		BoneIndex:               -1,
		CollisionGroup:          0,
		CollisionGroupMask:      NewCollisionGroupAll(),
		CollisionGroupMaskValue: 0,
		ShapeType:               SHAPE_BOX,
		Size:                    mmath.NewMVec3(),
		Position:                mmath.NewMVec3(),
		Rotation:                mmath.NewMVec3(),
		RigidBodyParam:          NewRigidBodyParam(),
		PhysicsType:             PHYSICS_TYPE_STATIC,
		XDirection:              mmath.NewMVec3(),
		YDirection:              mmath.NewMVec3(),
		ZDirection:              mmath.NewMVec3(),
		IsSystem:                false,
		Matrix:                  mmath.NewMMat4(),
		Bone:                    nil,
	}
}

func (rigidBody *RigidBody) Index() int {
	return rigidBody.index
}

func (rigidBody *RigidBody) SetIndex(index int) {
	rigidBody.index = index
}

func (rigidBody *RigidBody) Name() string {
	return rigidBody.name
}

func (rigidBody *RigidBody) SetName(name string) {
	rigidBody.name = name
}

func (rigidBody *RigidBody) EnglishName() string {
	return rigidBody.englishName
}

func (rigidBody *RigidBody) SetEnglishName(englishName string) {
	rigidBody.englishName = englishName
}

func (rigidBody *RigidBody) IsValid() bool {
	return rigidBody != nil && rigidBody.index >= 0
}

func (rigidBody *RigidBody) AsDynamic() bool {
	return rigidBody.PhysicsType == PHYSICS_TYPE_DYNAMIC || rigidBody.PhysicsType == PHYSICS_TYPE_DYNAMIC_BONE
}

func (rigidBody *RigidBody) Copy() core.IIndexNameModel {
	return &RigidBody{
		index:                   rigidBody.index,
		name:                    rigidBody.name,
		englishName:             rigidBody.englishName,
		BoneIndex:               rigidBody.BoneIndex,
		CollisionGroup:          rigidBody.CollisionGroup,
		CollisionGroupMask:      rigidBody.CollisionGroupMask,
		CollisionGroupMaskValue: rigidBody.CollisionGroupMaskValue,
		ShapeType:               rigidBody.ShapeType,
		Size:                    rigidBody.Size.Copy(),
		Position:                rigidBody.Position.Copy(),
		Rotation:                rigidBody.Rotation.Copy(),
		RigidBodyParam:          rigidBody.RigidBodyParam,
		PhysicsType:             rigidBody.PhysicsType,
		XDirection:              rigidBody.XDirection.Copy(),
		YDirection:              rigidBody.YDirection.Copy(),
		ZDirection:              rigidBody.ZDirection.Copy(),
		IsSystem:                rigidBody.IsSystem,
		Matrix:                  rigidBody.Matrix.Copy(),
		Bone:                    nil,
	}
}
