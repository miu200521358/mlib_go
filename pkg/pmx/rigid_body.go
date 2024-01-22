package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type RigidBodyParam struct {
	// 質量
	Mass float64
	// 移動減衰
	LinearDamping float64
	// 回転減衰
	AngularDamping float64
	// 反発力
	Restitution float64
	// 摩擦力
	Friction float64
}

// 剛体の形状
type Shape int

const (
	// 球
	SHAPE_SPHERE Shape = 0
	// 箱
	SHAPE_BOX Shape = 1
	// カプセル
	SHAPE_CAPSULE Shape = 2
)

// 剛体物理の計算モード
type PhysicsType int

const (
	// ボーン追従(static)
	PHYSICS_TYPE_STATIC PhysicsType = 0
	// 物理演算(dynamic)
	PHYSICS_TYPE_DYNAMIC PhysicsType = 1
	// 物理演算 + Bone位置合わせ
	PHYSICS_TYPE_DYNAMIC_BONE PhysicsType = 2
)

type CollisionGroup struct {
	IsCollisions []uint16
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
	*mcore.IndexModel
	// 剛体名
	Name string
	// 剛体名英
	EnglishName string
	// 関連ボーンIndex
	BoneIndex int
	// グループ
	CollisionGroup byte
	// 非衝突グループフラグ
	CollisionGroupMask CollisionGroup
	// 形状
	ShapeType Shape
	// サイズ(x,y,z)
	Size mmath.MVec3
	// 位置(x,y,z)
	Position mmath.MVec3
	// 回転(x,y,z) -> ラジアン角
	Rotation mmath.MRotation
	// 剛体パラ
	RigidBodyParam RigidBodyParam
	// 剛体の物理演算
	PhysicsType PhysicsType
	// X軸方向
	XDirection mmath.MVec3
	// Y軸方向
	YDirection mmath.MVec3
	// Z軸方向
	ZDirection mmath.MVec3
	// システムで追加した剛体か
	IsSystem bool
}

// NewRigidBody creates a new rigid body.
func NewRigidBody() *RigidBody {
	return &RigidBody{
		IndexModel:         &mcore.IndexModel{Index: -1},
		Name:               "",
		EnglishName:        "",
		BoneIndex:          -1,
		CollisionGroup:     0,
		CollisionGroupMask: NewCollisionGroupFromSlice([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		ShapeType:          SHAPE_BOX,
		Size:               mmath.MVec3{},
		Position:           mmath.MVec3{},
		Rotation:           mmath.MRotation{},
		RigidBodyParam:     RigidBodyParam{},
		PhysicsType:        PHYSICS_TYPE_STATIC,
		XDirection:         mmath.MVec3{},
		YDirection:         mmath.MVec3{},
		ZDirection:         mmath.MVec3{},
		IsSystem:           false,
	}
}

// 剛体リスト
type RigidBodies struct {
	*mcore.IndexModelCorrection[*RigidBody]
}

func NewRigidBodies() *RigidBodies {
	return &RigidBodies{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*RigidBody](),
	}
}
