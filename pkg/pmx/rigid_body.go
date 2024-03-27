package pmx

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/jinzhu/copier"

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
	BoneIndex                    int              // 関連ボーンIndex
	CollisionGroup               byte             // グループ
	CollisionGroupMask           CollisionGroup   // 非衝突グループフラグ
	CollisionGroupMaskValue      int              // 非衝突グループフラグ値
	ShapeType                    Shape            // 形状
	Size                         *mmath.MVec3     // サイズ(x,y,z)
	Position                     *mmath.MVec3     // 位置(x,y,z)
	Rotation                     *mmath.MRotation // 回転(x,y,z) -> ラジアン角
	RigidBodyParam               *RigidBodyParam  // 剛体パラ
	PhysicsType                  PhysicsType      // 剛体の物理演算
	CorrectPhysicsType           PhysicsType      // 剛体の物理演算(補正後)
	XDirection                   *mmath.MVec3     // X軸方向
	YDirection                   *mmath.MVec3     // Y軸方向
	ZDirection                   *mmath.MVec3     // Z軸方向
	IsSystem                     bool             // システムで追加した剛体か
	Matrix                       *mmath.MMat4     // 剛体の行列
	BtRigidBody                  mbt.BtRigidBody  // 物理剛体
	BtRigidBodyTransform         mbt.BtTransform  // 剛体の初期位置・回転情報
	BtRigidBodyPositionTransform mbt.BtTransform  // 剛体の初期位置情報
	JointedBoneIndex             int              // ジョイントで繋がってるボーンIndex
}

// NewRigidBody creates a new rigid body.
func NewRigidBody() *RigidBody {
	return &RigidBody{
		IndexNameModel:               &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		BoneIndex:                    -1,
		CollisionGroup:               0,
		CollisionGroupMask:           NewCollisionGroupFromSlice([]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		CollisionGroupMaskValue:      0,
		ShapeType:                    SHAPE_BOX,
		Size:                         mmath.NewMVec3(),
		Position:                     mmath.NewMVec3(),
		Rotation:                     mmath.NewRotationModel(),
		RigidBodyParam:               NewRigidBodyParam(),
		PhysicsType:                  PHYSICS_TYPE_STATIC,
		XDirection:                   mmath.NewMVec3(),
		YDirection:                   mmath.NewMVec3(),
		ZDirection:                   mmath.NewMVec3(),
		IsSystem:                     false,
		Matrix:                       mmath.NewMMat4(),
		BtRigidBody:                  nil,
		BtRigidBodyTransform:         nil,
		BtRigidBodyPositionTransform: nil,
		JointedBoneIndex:             -1,
	}
}

func (r *RigidBody) Copy() mcore.IIndexNameModel {
	copied := NewMorph()
	copier.CopyWithOption(copied, r, copier.Option{DeepCopy: true})
	return copied
}

// func (r *RigidBody) resetPhysics(enablePhysics bool) bool {
// 	// 物理ON＆自身が物理剛体＆現在のステートがActiveではない場合、True
// 	return enablePhysics &&
// 		r.CorrectPhysicsType != PHYSICS_TYPE_STATIC &&
// 		r.BtRigidBody.GetActivationState() == mbt.ISLAND_SLEEPING
// }

func (r *RigidBody) updateFlags(enablePhysics bool) bool {
	if r.CorrectPhysicsType == PHYSICS_TYPE_STATIC {
		// 剛体の位置更新に物理演算を使わない。もしくは物理演算OFF時
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		r.BtRigidBody.SetCollisionFlags(
			r.BtRigidBody.GetCollisionFlags() | int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を無効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが無効になる(剛体位置は判っているので不要)
		r.BtRigidBody.SetActivationState(mbt.DISABLE_SIMULATION)
		// if prevActivationState != mbt.DISABLE_SIMULATION {
		// 	return true
		// }
	} else {
		prevActivationState := r.BtRigidBody.GetActivationState()

		if enablePhysics {
			// 物理演算・物理+ボーン位置合わせの場合
			if prevActivationState != mbt.ACTIVE_TAG {
				r.BtRigidBody.SetCollisionFlags(0 & ^int(mbt.BtCollisionObjectCF_NO_CONTACT_RESPONSE))
				r.BtRigidBody.ForceActivationState(mbt.ACTIVE_TAG)

				// localInertia := mbt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
				// r.BtRigidBody.GetCollisionShape().(mbt.BtCollisionShape).CalculateLocalInertia(
				// 	float32(r.RigidBodyParam.Mass), localInertia)

				// r.BtRigidBody.GetMotionState().(mbt.BtMotionState).SetWorldTransform(r.BtRigidBodyTransform)

				return true
			} else {
				// 剛体の位置更新に物理演算を使う。
				// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
				r.BtRigidBody.SetCollisionFlags(
					r.BtRigidBody.GetCollisionFlags() & ^int(mbt.BtCollisionObjectCF_NO_CONTACT_RESPONSE))
				// 毎ステップの剛体位置通知を有効にする
				// MotionState::setWorldTransformの毎ステップ呼び出しが有効になる(剛体位置が変わるので必要)
				r.BtRigidBody.SetActivationState(mbt.ACTIVE_TAG)
			}
		} else {
			// 物理OFF時
			r.BtRigidBody.SetCollisionFlags(
				r.BtRigidBody.GetCollisionFlags() | int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
			r.BtRigidBody.SetActivationState(mbt.ISLAND_SLEEPING)

			if prevActivationState != mbt.ISLAND_SLEEPING {
				return true
			}
		}
		// if prevActivationState != mbt.ACTIVE_TAG {
		// 	// r.BtRigidBody.ForceActivationState(mbt.ACTIVE_TAG)
		// 	// r.BtRigidBody.Activate(true)
		// 	modelPhysics.RemoveRigidBody(r.BtRigidBody)
		// 	r.InitPhysics(modelPhysics)
		// 	return true
		// } else {
		// 	r.BtRigidBody.SetActivationState(mbt.ACTIVE_TAG)
		// }
	}

	return false
}

func (r *RigidBody) initPhysics(modelPhysics *mphysics.MPhysics) {
	var btCollisionShape mbt.BtCollisionShape
	switch r.ShapeType {
	case SHAPE_SPHERE:
		// 球剛体
		btCollisionShape = mbt.NewBtSphereShape(float32(r.Size.GetX()))
	case SHAPE_BOX:
		// 箱剛体
		btCollisionShape = mbt.NewBtBoxShape(
			mbt.NewBtVector3(float32(r.Size.GetX()), float32(r.Size.GetY()), float32(r.Size.GetZ())))
	case SHAPE_CAPSULE:
		// カプセル剛体
		btCollisionShape = mbt.NewBtCapsuleShape(float32(r.Size.GetX()), float32(r.Size.GetY()))
	}

	r.CorrectPhysicsType = r.PhysicsType
	// if r.PhysicsType == PHYSICS_TYPE_DYNAMIC_BONE && r.BoneIndex < 0 {
	// 	// ボーン追従 + 物理剛体の場合、ボーンIndexが設定されていない場合は物理剛体に変更
	// 	r.CorrectPhysicsType = PHYSICS_TYPE_DYNAMIC
	// }

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

	// 剛体の初期位置と回転
	r.BtRigidBodyTransform = mbt.NewBtTransform(r.Rotation.GetQuaternion().Bullet(), r.Position.Bullet())

	// 剛体の初期位置
	r.BtRigidBodyPositionTransform = mbt.NewBtTransform()
	r.BtRigidBodyPositionTransform.SetIdentity()
	r.BtRigidBodyPositionTransform.SetOrigin(r.Position.Bullet())

	// {
	// 	mlog.V("---------------------------------")
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	r.BtRigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	mlog.V("1. [%s] BtRigidBodyTransform: \n%v\n", r.Name, mat)
	// }

	// 剛体のグローバル位置と回転
	motionState := mbt.NewBtDefaultMotionState(r.BtRigidBodyTransform)

	r.BtRigidBody = mbt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	r.BtRigidBody.SetDamping(float32(r.RigidBodyParam.LinearDamping), float32(r.RigidBodyParam.AngularDamping))
	r.BtRigidBody.SetRestitution(float32(r.RigidBodyParam.Restitution))
	r.BtRigidBody.SetFriction(float32(r.RigidBodyParam.Friction))
	r.BtRigidBody.SetUserIndex(r.Index)
	// r.BtRigidBody.SetSleepingThresholds(0.1, (180.0 * 0.1 / math.Pi))

	r.updateFlags(true)

	// mlog.V("name: %s, group: %d, mask: %d\n", r.Name, r.CollisionGroup, r.CollisionGroupMaskValue)

	// modelPhysics.AddNonFilterProxy(r.BtRigidBody.GetBroadphaseProxy())
	// 剛体・剛体グループ・非衝突グループを追加
	modelPhysics.AddRigidBody(r.BtRigidBody, 1<<r.CollisionGroup, r.CollisionGroupMaskValue)
}

func (r *RigidBody) updateTransform(
	boneTransforms []*mbt.BtTransform,
	isForce bool,
) {
	if r.BtRigidBody == nil || r.BtRigidBody.GetMotionState() == nil ||
		(r.CorrectPhysicsType == PHYSICS_TYPE_DYNAMIC && !isForce) {
		return
	}

	// {
	// 	mlog.V("----------")
	// }

	// 剛体のグローバル位置を確定
	rigidBodyTransform := mbt.NewBtTransform()
	if r.BoneIndex >= 0 && r.BoneIndex < len(boneTransforms) {
		rigidBodyTransform.Mult(*boneTransforms[r.BoneIndex], r.BtRigidBodyTransform)
	} else if r.JointedBoneIndex >= 0 && r.JointedBoneIndex < len(boneTransforms) {
		rigidBodyTransform.Mult(*boneTransforms[r.JointedBoneIndex], r.BtRigidBodyTransform)
	}

	// if r.BoneIndex >= 0 && r.BoneIndex < len(boneTransforms) {
	// 	mat := mgl32.Mat4{}
	// 	(*boneTransforms[r.BoneIndex]).GetOpenGLMatrix(&mat[0])
	// 	mlog.V("2. [%d] boneTransform: \n%v\n", r.BoneIndex, mat)
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	mlog.V("2. [%s] rigidBodyTransform: \n%v\n", r.Name, mat)
	// }

	motionState := r.BtRigidBody.GetMotionState().(mbt.BtMotionState)
	motionState.SetWorldTransform(rigidBodyTransform)
}

func (r *RigidBody) updateMatrix(
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
) {
	if r.BtRigidBody == nil || r.BtRigidBody.GetMotionState() == nil || r.CorrectPhysicsType == PHYSICS_TYPE_STATIC {
		return
	}

	// if r.Name == "前髪" {
	// 	state := r.BtRigidBody.GetActivationState()
	// 	isStatic := r.BtRigidBody.IsStaticOrKinematicObject()
	// 	flags := r.BtRigidBody.GetCollisionFlags()
	// 	mlog.V("name: %s, state: %v, static: %v, flags: %v, isKinematic: %v, isStatic: %v\n", r.Name, state, isStatic, flags, flags&int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT) != 0, flags&int(mbt.BtCollisionObjectCF_STATIC_OBJECT) != 0)
	// }

	// {
	// 	mlog.V("----------")
	// }

	motionState := r.BtRigidBody.GetMotionState().(mbt.BtMotionState)

	rigidBodyTransform := mbt.NewBtTransform()
	motionState.GetWorldTransform(rigidBodyTransform)

	// if r.CorrectPhysicsType == PHYSICS_TYPE_DYNAMIC_BONE {
	// 	{
	// 		mat := mgl32.Mat4{}
	// 		rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 		mlog.V("3. [%s] rigidBodyTransform Before: \n%v\n", r.Name, mat)
	// 	}

	// 	if r.BoneIndex >= 0 && r.BoneIndex < len(boneTransforms) {
	// 		// 物理+ボーン追従はボーン移動成分を現在のボーン位置にする
	// 		boneGlobalTransform := mbt.NewBtTransform()
	// 		boneGlobalTransform.Mult(*boneTransforms[r.BoneIndex], r.BtRigidBodyTransform)

	// 		rigidBodyTransform.SetOrigin(boneGlobalTransform.GetOrigin().(mbt.BtVector3))
	// 	} else if r.JointedBoneIndex >= 0 && r.JointedBoneIndex < len(boneTransforms) {
	// 		// ジョイントで繋がっているボーンがある場合はそのボーン位置にする
	// 		boneGlobalTransform := mbt.NewBtTransform()
	// 		boneGlobalTransform.Mult(*boneTransforms[r.JointedBoneIndex], r.BtRigidBodyTransform)

	// 		rigidBodyTransform.SetOrigin(boneGlobalTransform.GetOrigin().(mbt.BtVector3))
	// 	}

	// 	{
	// 		mat := mgl32.Mat4{}
	// 		rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 		mlog.V("3. [%s] rigidBodyTransform After: \n%v\n", r.Name, mat)
	// 	}
	// }

	boneLocalTransform := mbt.NewBtTransform()
	boneLocalTransform.Mult(rigidBodyTransform, r.BtRigidBodyTransform.Inverse())

	physicsBoneMatrix := mgl32.Mat4{}
	boneLocalTransform.GetOpenGLMatrix(&physicsBoneMatrix[0])

	// if r.BoneIndex >= 0 && r.BoneIndex < len(boneTransforms) {
	// 	mat := mgl32.Mat4{}
	// 	(*boneTransforms[r.BoneIndex]).GetOpenGLMatrix(&mat[0])
	// 	mlog.V("3. [%d] boneTransform: \n%v\n", r.BoneIndex, mat)
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	mlog.V("3. [%s] rigidBodyTransform: \n%v\n", r.Name, mat)
	// }
	// {
	// 	mlog.V("3. [%s] physicsBoneMatrix: \n%v\n", r.Name, physicsBoneMatrix)
	// }

	if r.BoneIndex >= 0 && r.BoneIndex < len(boneTransforms) {
		boneMatrixes[r.BoneIndex] = &physicsBoneMatrix
	} else if r.JointedBoneIndex >= 0 && r.JointedBoneIndex < len(boneTransforms) {
		boneMatrixes[r.JointedBoneIndex] = &physicsBoneMatrix
	}
}

func (r *RigidBody) deletePhysics() {
	if r.BtRigidBody != nil {
		r.BtRigidBody.SetUserIndex(-1)
		r.BtRigidBody.SetMotionState(nil)
		r.BtRigidBody.SetCollisionShape(nil)
		r.BtRigidBody = nil
	}
}

// 剛体リスト
type RigidBodies struct {
	*mcore.IndexNameModels[*RigidBody]
}

func NewRigidBodies() *RigidBodies {
	return &RigidBodies{
		IndexNameModels: mcore.NewIndexNameModelCorrection[*RigidBody](),
	}
}

func (r *RigidBodies) initPhysics(physics *mphysics.MPhysics) {
	// 剛体を順番にボーンと紐付けていく
	for _, rigidBody := range r.GetSortedData() {
		// 物理設定の初期化
		rigidBody.initPhysics(physics)
	}
}

func (r *RigidBodies) deletePhysics(modelPhysics *mphysics.MPhysics) {
	for _, rigidBody := range r.Data {
		modelPhysics.DeleteRigidBody(rigidBody.BtRigidBody)
		rigidBody.deletePhysics()
	}
}
