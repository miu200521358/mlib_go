package rigidbody

import (
	"github.com/miu200521358/mlib_go/pkg/math/mrotation"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
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
type RigidBodyShape int

const (
	// 球
	Sphere RigidBodyShape = 0
	// 箱
	Box RigidBodyShape = 1
	// カプセル
	Capsule RigidBodyShape = 2
)

// 剛体物理の計算モード
type RigidBodyMode int

const (
	// ボーン追従(static)
	Static RigidBodyMode = 0
	// 物理演算(dynamic)
	Dynamic RigidBodyMode = 1
	// 物理演算 + Bone位置合わせ
	DynamicBone RigidBodyMode = 2
)

type RigidBodyCollisionGroup int

const (
	// 0:グループなし
	None    RigidBodyCollisionGroup = 0x0000
	Group01 RigidBodyCollisionGroup = 0x0001
	Group02 RigidBodyCollisionGroup = 0x0002
	Group03 RigidBodyCollisionGroup = 0x0004
	Group04 RigidBodyCollisionGroup = 0x0008
	Group05 RigidBodyCollisionGroup = 0x0010
	Group06 RigidBodyCollisionGroup = 0x0020
	Group07 RigidBodyCollisionGroup = 0x0040
	Group08 RigidBodyCollisionGroup = 0x0080
	Group09 RigidBodyCollisionGroup = 0x0100
	Group10 RigidBodyCollisionGroup = 0x0200
	Group11 RigidBodyCollisionGroup = 0x0400
	Group12 RigidBodyCollisionGroup = 0x0800
	Group13 RigidBodyCollisionGroup = 0x1000
	Group14 RigidBodyCollisionGroup = 0x2000
	Group15 RigidBodyCollisionGroup = 0x4000
	Group16 RigidBodyCollisionGroup = 0x8000
)

type T struct {
	// 剛体INDEX
	Index int
	// 剛体名
	Name string
	// 剛体名英
	EnglishName string
	// 関連ボーンIndex
	BoneIndex int
	// グループ
	CollisionGroup int
	// 非衝突グループフラグ
	NoCollisionGroup RigidBodyCollisionGroup
	// 形状
	ShapeType RigidBodyShape
	// サイズ(x,y,z)
	ShapeSize mvec3.T
	// 位置(x,y,z)
	ShapePosition mvec3.T
	// 回転(x,y,z) -> ラジアン角
	ShapeRotation mrotation.T
	// 剛体パラ
	Param RigidBodyParam
	// 剛体の物理演算
	Mode RigidBodyMode
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
	index int,
	name string,
	englishName string,
	boneIndex int,
	collisionGroup int,
	noCollisionGroup RigidBodyCollisionGroup,
	shapeType RigidBodyShape,
	shapeSize mvec3.T,
	shapePosition mvec3.T,
	shapeRotation mrotation.T,
	param RigidBodyParam,
	mode RigidBodyMode,
	xDirection mvec3.T,
	yDirection mvec3.T,
	zDirection mvec3.T,
	isSystem bool,
) *T {
	return &T{
		Index:            index,
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

// Copy
func (v *T) Copy() *T {
	copied := *v
	return &copied
}
