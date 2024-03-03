package pmx

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
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
		Mass:           0,
		LinearDamping:  0,
		AngularDamping: 0,
		Restitution:    0,
		Friction:       0,
	}
}

// 剛体の形状
type Shape int

const (
	SHAPE_SPHERE  Shape = 0 // 球
	SHAPE_BOX     Shape = 1 // 箱
	SHAPE_CAPSULE Shape = 2 // カプセル
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
	*mcore.IndexNameModel
	BoneIndex          int                  // 関連ボーンIndex
	CollisionGroup     byte                 // グループ
	CollisionGroupMask CollisionGroup       // 非衝突グループフラグ
	ShapeType          Shape                // 形状
	Size               *mmath.MVec3         // サイズ(x,y,z)
	Position           *mmath.MVec3         // 位置(x,y,z)
	Rotation           *mmath.MRotation     // 回転(x,y,z) -> ラジアン角
	RigidBodyParam     *RigidBodyParam      // 剛体パラ
	PhysicsType        PhysicsType          // 剛体の物理演算
	XDirection         *mmath.MVec3         // X軸方向
	YDirection         *mmath.MVec3         // Y軸方向
	ZDirection         *mmath.MVec3         // Z軸方向
	IsSystem           bool                 // システムで追加した剛体か
	Matrix             *mmath.MMat4         // 剛体の行列
	BtTransform        mbt.BtTransform      // Bulletの変換行列
	BtCollisionShape   mbt.BtCollisionShape // 物理形状
	BtRigidBody        mbt.BtRigidBody      // 物理剛体
}

// NewRigidBody creates a new rigid body.
func NewRigidBody() *RigidBody {
	return &RigidBody{
		IndexNameModel:     &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		BoneIndex:          -1,
		CollisionGroup:     0,
		CollisionGroupMask: NewCollisionGroupFromSlice([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		ShapeType:          SHAPE_BOX,
		Size:               mmath.NewMVec3(),
		Position:           mmath.NewMVec3(),
		Rotation:           mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		RigidBodyParam:     NewRigidBodyParam(),
		PhysicsType:        PHYSICS_TYPE_STATIC,
		XDirection:         mmath.NewMVec3(),
		YDirection:         mmath.NewMVec3(),
		ZDirection:         mmath.NewMVec3(),
		IsSystem:           false,
	}
}

func (r *RigidBody) InitPhysics(modelPhysics *mphysics.MPhysics, bone *Bone) {
	switch r.ShapeType {
	case SHAPE_SPHERE:
		r.BtCollisionShape = mbt.NewBtSphereShape(float32(r.Size.GetX()))
	case SHAPE_BOX:
		r.BtCollisionShape = mbt.NewBtBoxShape(
			mbt.NewBtVector3(float32(r.Size.GetX()), float32(r.Size.GetY()), float32(r.Size.GetZ())))
	case SHAPE_CAPSULE:
		r.BtCollisionShape = mbt.NewBtCapsuleShape(float32(r.Size.GetX()), float32(r.Size.GetY()))
	}

	// 質量
	mass := float32(0.0)
	localInertia := mbt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
	if r.PhysicsType == PHYSICS_TYPE_STATIC {
		// ボーン追従の場合そのまま設定
		mass = float32(r.RigidBodyParam.Mass)
	}
	if mass != 0 {
		// 質量が設定されている場合、慣性を計算
		r.BtCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	// 剛体の回転
	rotationMat := r.Rotation.GetQuaternion().ToMat4()

	// 剛体の位置
	translationMat := mmath.NewMMat4()
	translationMat.Translate(r.Position)

	r.Matrix = mmath.NewMMat4()
	r.Matrix.Mul(translationMat)
	r.Matrix.Mul(rotationMat)

	// ボーンから見た剛体のローカル行列(ボーンがなければ初期行列)
	var rigidBodyOffsetGLMat *mgl32.Mat4
	if bone != nil {
		rigidBodyOffsetGLMat = bone.OffsetMatrix.Muled(r.Matrix).GL()
	} else {
		rigidBodyOffsetGLMat = mmath.NewMMat4().GL()
	}

	// OpenGL行列を設定
	r.BtTransform = mbt.NewBtTransform()
	r.BtTransform.SetFromOpenGLMatrix(&rigidBodyOffsetGLMat[0])

	var motionState mbt.BtDefaultMotionState
	if r.PhysicsType == PHYSICS_TYPE_STATIC {
		motionState = mphysics.NewStaticMotionState()
	} else if r.PhysicsType == PHYSICS_TYPE_DYNAMIC {
		motionState = mphysics.NewDynamicMotionState()
	} else if r.PhysicsType == PHYSICS_TYPE_DYNAMIC_BONE {
		motionState = mphysics.NewDynamicBoneMotionState()
	}
	motionState.SetWorldTransform(r.BtTransform)

	r.BtRigidBody = mbt.NewBtRigidBody(mass, motionState, r.BtCollisionShape, localInertia)
	r.BtRigidBody.SetDamping(float32(r.RigidBodyParam.LinearDamping), float32(r.RigidBodyParam.AngularDamping))
	r.BtRigidBody.SetRestitution(float32(r.RigidBodyParam.Restitution))
	r.BtRigidBody.SetFriction(float32(r.RigidBodyParam.Friction))
	// btRigidBody.SetUserIndex(mbt.)
	r.BtRigidBody.SetSleepingThresholds(0.01, (180.0 * 0.1 / math.Pi))
	r.BtRigidBody.SetActivationState(mbt.DISABLE_DEACTIVATION)

	modelPhysics.AddRigidBody(r.BtRigidBody)
}

// 剛体リスト
type RigidBodies struct {
	*mcore.IndexNameModelCorrection[*RigidBody]
}

func NewRigidBodies() *RigidBodies {
	return &RigidBodies{
		IndexNameModelCorrection: mcore.NewIndexNameModelCorrection[*RigidBody](),
	}
}
