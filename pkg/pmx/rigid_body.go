package pmx

import (
	"fmt"
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

func (c *CollisionGroup) ToInt() int {
	// uint16の配列をintに変換
	result := 0
	for i, v := range c.IsCollisions {
		result |= int(v) << (16 * i)
	}
	// 非衝突を衝突グループに変換するため、intを反転させる
	return ^result
}

type RigidBody struct {
	*mcore.IndexNameModel
	BoneIndex                       int              // 関連ボーンIndex
	CollisionGroup                  byte             // グループ
	CollisionGroupMask              CollisionGroup   // 非衝突グループフラグ
	ShapeType                       Shape            // 形状
	Size                            *mmath.MVec3     // サイズ(x,y,z)
	Position                        *mmath.MVec3     // 位置(x,y,z)
	Rotation                        *mmath.MRotation // 回転(x,y,z) -> ラジアン角
	RigidBodyParam                  *RigidBodyParam  // 剛体パラ
	PhysicsType                     PhysicsType      // 剛体の物理演算
	XDirection                      *mmath.MVec3     // X軸方向
	YDirection                      *mmath.MVec3     // Y軸方向
	ZDirection                      *mmath.MVec3     // Z軸方向
	IsSystem                        bool             // システムで追加した剛体か
	Matrix                          *mmath.MMat4     // 剛体の行列
	BtRigidBody                     mbt.BtRigidBody  // 物理剛体
	BtInitialTransform              mbt.BtTransform  // 剛体のローカル変換情報
	BoneInvertedPositionBtTransform mbt.BtTransform  // ボーンの初期逆行列
	BonePositionBtTransform         mbt.BtTransform  // ボーンの初期位置
}

// NewRigidBody creates a new rigid body.
func NewRigidBody() *RigidBody {
	return &RigidBody{
		IndexNameModel:                  &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		BoneIndex:                       -1,
		CollisionGroup:                  0,
		CollisionGroupMask:              NewCollisionGroupFromSlice([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		ShapeType:                       SHAPE_BOX,
		Size:                            mmath.NewMVec3(),
		Position:                        mmath.NewMVec3(),
		Rotation:                        mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		RigidBodyParam:                  NewRigidBodyParam(),
		PhysicsType:                     PHYSICS_TYPE_STATIC,
		XDirection:                      mmath.NewMVec3(),
		YDirection:                      mmath.NewMVec3(),
		ZDirection:                      mmath.NewMVec3(),
		IsSystem:                        false,
		Matrix:                          mmath.NewMMat4(),
		BtRigidBody:                     nil,
		BtInitialTransform:              nil,
		BoneInvertedPositionBtTransform: nil,
	}
}

func (r *RigidBody) InitPhysics(modelPhysics *mphysics.MPhysics, bone *Bone) {
	var btCollisionShape mbt.BtCollisionShape
	switch r.ShapeType {
	case SHAPE_SPHERE:
		btCollisionShape = mbt.NewBtSphereShape(float32(r.Size.GetX()))
	case SHAPE_BOX:
		btCollisionShape = mbt.NewBtBoxShape(
			mbt.NewBtVector3(float32(r.Size.GetX()), float32(r.Size.GetY()), float32(r.Size.GetZ())))
	case SHAPE_CAPSULE:
		btCollisionShape = mbt.NewBtCapsuleShape(float32(r.Size.GetX()), float32(r.Size.GetY()))
	}

	// 質量
	mass := float32(0.0)
	localInertia := mbt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
	if r.PhysicsType != PHYSICS_TYPE_STATIC {
		// ボーン追従ではない場合そのまま設定
		mass = float32(r.RigidBodyParam.Mass)
	}
	if mass != 0 {
		// 質量が設定されている場合、慣性を計算
		btCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	// // ボーンのローカル位置
	// boneTransform := mbt.NewBtTransform()
	// boneTransform.SetIdentity()
	// boneTransform.SetOrigin(boneLocalPosition.Bullet())

	// 剛体の初期位置
	if bone == nil {
		r.BtInitialTransform = mbt.NewBtTransform(
			r.Rotation.GetQuaternion().Bullet(), mbt.NewBtVector3())
		r.BoneInvertedPositionBtTransform = mbt.NewBtTransform()
		r.BonePositionBtTransform = mbt.NewBtTransform()
	} else {
		r.BtInitialTransform = mbt.NewBtTransform(
			r.Rotation.GetQuaternion().Bullet(), r.Position.Bullet())

		r.BoneInvertedPositionBtTransform = mbt.NewBtTransform()
		r.BoneInvertedPositionBtTransform.SetIdentity()
		r.BoneInvertedPositionBtTransform.SetOrigin(r.Position.Inverted().Bullet())

		r.BonePositionBtTransform = mbt.NewBtTransform()
		r.BonePositionBtTransform.SetIdentity()
		r.BonePositionBtTransform.SetOrigin(r.Position.Bullet())
	}

	{
		fmt.Println("---------------------------------")
	}
	{
		mat := mgl32.Mat4{}
		r.BtInitialTransform.GetOpenGLMatrix(&mat[0])
		fmt.Printf("1. [%s] BtInitialTransform: \n%v\n", r.Name, mat)
	}

	motionState := mbt.NewBtDefaultMotionState(r.BtInitialTransform)

	r.BtRigidBody = mbt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	r.BtRigidBody.SetDamping(float32(r.RigidBodyParam.LinearDamping), float32(r.RigidBodyParam.AngularDamping))
	r.BtRigidBody.SetRestitution(float32(r.RigidBodyParam.Restitution))
	r.BtRigidBody.SetFriction(float32(r.RigidBodyParam.Friction))
	// btRigidBody.SetUserIndex(mbt.)
	r.BtRigidBody.SetSleepingThresholds(0.01, (180.0 * 0.1 / math.Pi))

	if r.PhysicsType == PHYSICS_TYPE_STATIC {
		// 剛体の位置更新に物理演算を使わない。
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		r.BtRigidBody.SetCollisionFlags(
			r.BtRigidBody.GetCollisionFlags() | int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を無効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが無効になる(剛体位置は判っているので不要)
		r.BtRigidBody.SetActivationState(mbt.DISABLE_SIMULATION)
	} else {
		// 物理演算・物理+ボーン位置合わせの場合
		// 剛体の位置更新に物理演算を使う。
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		r.BtRigidBody.SetCollisionFlags(
			r.BtRigidBody.GetCollisionFlags() & int(^mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を有効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが有効になる(剛体位置が変わるので必要)
		r.BtRigidBody.SetActivationState(mbt.ACTIVE_TAG)
	}

	modelPhysics.AddRigidBody(r.BtRigidBody, int(r.CollisionGroup), r.CollisionGroupMask.ToInt())
}

// func (r *RigidBody) SetActivation(activation bool) {
// 	if r.BtRigidBody == nil {
// 		return
// 	}

// 	if r.PhysicsType != PHYSICS_TYPE_STATIC {
// 		// 物理の場合
// 		if activation {
// 			r.BtRigidBody.SetCollisionFlags(
// 				r.BtRigidBody.GetCollisionFlags() & int(^mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
// 			r.BtRigidBody.SetMotionState(r.ActiveMotionState)
// 		} else {
// 			r.BtRigidBody.SetCollisionFlags(
// 				r.BtRigidBody.GetCollisionFlags() | int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
// 			r.BtRigidBody.SetMotionState(r.StaticMotionState)
// 		}
// 	} else {
// 		// ボーン追従の場合
// 		r.BtRigidBody.SetMotionState(r.StaticMotionState)
// 	}
// }

// func (r *RigidBody) ResetPhysics() {
// 	if r.BtRigidBody == nil {
// 		return
// 	}
// 	r.BtRigidBody.SetActivationState(mbt.DISABLE_SIMULATION)
// 	r.ActiveMotionState.Reset()
// }

func (r *RigidBody) UpdateTransform(
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
) {
	if r.BtRigidBody == nil || r.BtRigidBody.GetMotionState() == nil ||
		r.BoneIndex < 0 || r.BoneIndex >= len(boneTransforms) || r.PhysicsType != PHYSICS_TYPE_STATIC {
		return
	}

	{
		fmt.Println("----------")
	}

	// 剛体のグローバル位置と向き
	rigidBodyTransform := mbt.NewBtTransform()
	rigidBodyTransform.Mult(*boneTransforms[r.BoneIndex], r.BtInitialTransform)
	// rigidBodyTransform.SetFromOpenGLMatrix(&boneMatrixes[r.BoneIndex][0])

	{
		mat := mgl32.Mat4{}
		(*boneTransforms[r.BoneIndex]).GetOpenGLMatrix(&mat[0])
		fmt.Printf("2. [%d] boneTransform: \n%v\n", r.BoneIndex, mat)
	}
	{
		mat := mgl32.Mat4{}
		rigidBodyTransform.GetOpenGLMatrix(&mat[0])
		fmt.Printf("2. [%s] rigidBodyTransform: \n%v\n", r.Name, mat)
	}
	motionState := r.BtRigidBody.GetMotionState().(mbt.BtMotionState)
	motionState.SetWorldTransform(rigidBodyTransform)
}

func (r *RigidBody) UpdateMatrix(
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
) {
	if r.BtRigidBody == nil || r.BtRigidBody.GetMotionState() == nil ||
		r.BoneIndex < 0 || r.BoneIndex >= len(boneMatrixes) || r.PhysicsType == PHYSICS_TYPE_STATIC {
		return
	}

	{
		fmt.Println("----------")
	}

	motionState := r.BtRigidBody.GetMotionState().(mbt.BtMotionState)

	rigidBodyTransform := mbt.NewBtTransform()
	motionState.GetWorldTransform(rigidBodyTransform)

	if r.PhysicsType == PHYSICS_TYPE_DYNAMIC_BONE {
		{
			mat := mgl32.Mat4{}
			rigidBodyTransform.GetOpenGLMatrix(&mat[0])
			fmt.Printf("3. [%s] rigidBodyTransform Before: \n%v\n", r.Name, mat)
		}

		// 物理+ボーン追従はボーン移動成分を現在のボーン位置にする
		boneGlobalTransform := mbt.NewBtTransform()
		boneGlobalTransform.Mult(*boneTransforms[r.BoneIndex], r.BonePositionBtTransform)

		rigidBodyTransform.SetOrigin(boneGlobalTransform.GetOrigin().(mbt.BtVector3))
	}

	boneLocalTransform := mbt.NewBtTransform()
	boneLocalTransform.Mult(rigidBodyTransform, r.BoneInvertedPositionBtTransform)

	physicsBoneMatrix := mgl32.Mat4{}
	boneLocalTransform.GetOpenGLMatrix(&physicsBoneMatrix[0])

	{
		fmt.Printf("3. [%d] boneMatrixes: \n%v\n", r.BoneIndex, boneMatrixes[r.BoneIndex])
	}
	{
		mat := mgl32.Mat4{}
		(*boneTransforms[r.BoneIndex]).GetOpenGLMatrix(&mat[0])
		fmt.Printf("3. [%d] boneTransform: \n%v\n", r.BoneIndex, mat)
	}
	{
		mat := mgl32.Mat4{}
		rigidBodyTransform.GetOpenGLMatrix(&mat[0])
		fmt.Printf("3. [%s] rigidBodyTransform: \n%v\n", r.Name, mat)
	}
	{
		mat := mgl32.Mat4{}
		boneLocalTransform.GetOpenGLMatrix(&mat[0])
		fmt.Printf("3. [%s] boneLocalTransform: \n%v\n", r.Name, mat)
	}
	{
		fmt.Printf("3. [%s] physicsBoneMatrix: \n%v\n", r.Name, physicsBoneMatrix)
	}

	boneMatrixes[r.BoneIndex] = &physicsBoneMatrix
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
