package rigidbody

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mrotation"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"

)

type Param struct {
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
type PhysicsMode int

const (
	// ボーン追従(static)
	PHYSICS_MODE_STATIC PhysicsMode = 0
	// 物理演算(dynamic)
	PHYSICS_MODE_DYNAMIC PhysicsMode = 1
	// 物理演算 + Bone位置合わせ
	PHYSICS_MODE_DYNAMIC_BONE PhysicsMode = 2
)

type CollisionGroup int

const (
	// 0:グループなし
	COLLISION_NONE    CollisionGroup = 0x0000
	COLLISION_GROUP01 CollisionGroup = 0x0001
	COLLISION_GROUP02 CollisionGroup = 0x0002
	COLLISION_GROUP03 CollisionGroup = 0x0004
	COLLISION_GROUP04 CollisionGroup = 0x0008
	COLLISION_GROUP05 CollisionGroup = 0x0010
	COLLISION_GROUP06 CollisionGroup = 0x0020
	COLLISION_GROUP07 CollisionGroup = 0x0040
	COLLISION_GROUP08 CollisionGroup = 0x0080
	COLLISION_GROUP09 CollisionGroup = 0x0100
	COLLISION_GROUP10 CollisionGroup = 0x0200
	COLLISION_GROUP11 CollisionGroup = 0x0400
	COLLISION_GROUP12 CollisionGroup = 0x0800
	COLLISION_GROUP13 CollisionGroup = 0x1000
	COLLISION_GROUP14 CollisionGroup = 0x2000
	COLLISION_GROUP15 CollisionGroup = 0x4000
	// 床との衝突
	COLLISION_GROUP16 CollisionGroup = 0x8000
)

type RigidBody struct {
	*index_model.IndexModel
	// 剛体名
	Name string
	// 剛体名英
	EnglishName string
	// 関連ボーンIndex
	BoneIndex int
	// グループ
	CollisionGroup int
	// 非衝突グループフラグ
	NoCollisionGroup CollisionGroup
	// 形状
	ShapeType Shape
	// サイズ(x,y,z)
	ShapeSize mvec3.T
	// 位置(x,y,z)
	ShapePosition mvec3.T
	// 回転(x,y,z) -> ラジアン角
	ShapeRotation mrotation.T
	// 剛体パラ
	Param Param
	// 剛体の物理演算
	Mode PhysicsMode
	// X軸方向
	XDirection mvec3.T
	// Y軸方向
	YDirection mvec3.T
	// Z軸方向
	ZDirection mvec3.T
	// システムで追加した剛体か
	IsSystem bool
}

// NewRigidBody creates a new rigid body.
func NewRigidBody(
	name string,
	englishName string,
	boneIndex int,
	collisionGroup int,
	noCollisionGroup CollisionGroup,
	shapeType Shape,
	shapeSize mvec3.T,
	shapePosition mvec3.T,
	shapeRotation mrotation.T,
	param Param,
	mode PhysicsMode,
	xDirection mvec3.T,
	yDirection mvec3.T,
	zDirection mvec3.T,
	isSystem bool,
) *RigidBody {
	return &RigidBody{
		IndexModel:       &index_model.IndexModel{Index: -1},
		Name:             name,
		EnglishName:      englishName,
		BoneIndex:        boneIndex,
		CollisionGroup:   collisionGroup,
		NoCollisionGroup: noCollisionGroup,
		ShapeType:        shapeType,
		ShapeSize:        shapeSize,
		ShapePosition:    shapePosition,
		ShapeRotation:    shapeRotation,
		Param:            param,
		Mode:             mode,
		XDirection:       xDirection,
		YDirection:       yDirection,
		ZDirection:       zDirection,
		IsSystem:         isSystem,
	}
}

// 剛体リスト
type RigidBodies struct {
	*index_model.IndexModelCorrection[*RigidBody]
}

func NewRigidBodies(name string) *RigidBodies {
	return &RigidBodies{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*RigidBody](),
	}
}
