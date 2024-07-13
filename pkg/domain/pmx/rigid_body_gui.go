//go:build windows
// +build windows

package pmx

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
)

func (r *RigidBody) UpdateFlags(
	modelIndex int, modelPhysics *mbt.MPhysics, enablePhysics, resetPhysics bool,
) bool {
	btRigidBody, _ := modelPhysics.GetRigidBody(modelIndex, r.Index)
	if btRigidBody == nil {
		return false
	}

	if r.PhysicsType == PHYSICS_TYPE_STATIC || resetPhysics {
		// 剛体の位置更新に物理演算を使わない。もしくは物理演算OFF時
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		btRigidBody.SetCollisionFlags(
			btRigidBody.GetCollisionFlags() | int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を無効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが無効になる(剛体位置は判っているので不要)
		btRigidBody.SetActivationState(bt.DISABLE_SIMULATION)

		if resetPhysics {
			return true
		}
	} else {
		prevActivationState := btRigidBody.GetActivationState()

		if enablePhysics {
			// 物理演算・物理+ボーン位置合わせの場合
			if prevActivationState != bt.ACTIVE_TAG {
				btRigidBody.SetCollisionFlags(0 & ^int(bt.BtCollisionObjectCF_NO_CONTACT_RESPONSE))
				btRigidBody.ForceActivationState(bt.ACTIVE_TAG)

				return true
			} else {
				// 剛体の位置更新に物理演算を使う。
				// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
				btRigidBody.SetCollisionFlags(
					btRigidBody.GetCollisionFlags() & ^int(bt.BtCollisionObjectCF_NO_CONTACT_RESPONSE))
				// 毎ステップの剛体位置通知を有効にする
				// MotionState::setWorldTransformの毎ステップ呼び出しが有効になる(剛体位置が変わるので必要)
				btRigidBody.SetActivationState(bt.ACTIVE_TAG)
			}
		} else {
			// 物理OFF時
			btRigidBody.SetCollisionFlags(
				btRigidBody.GetCollisionFlags() | int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT))
			btRigidBody.SetActivationState(bt.ISLAND_SLEEPING)

			if prevActivationState != bt.ISLAND_SLEEPING {
				return true
			}
		}
	}

	return false
}

func (r *RigidBody) initPhysics(modelIndex int, modelPhysics *mbt.MPhysics) {
	var btCollisionShape bt.BtCollisionShape

	// マイナスサイズは許容しない
	size := r.Size.Absed()

	switch r.ShapeType {
	case SHAPE_SPHERE:
		// 球剛体
		btCollisionShape = bt.NewBtSphereShape(float32(size.GetX()))
	case SHAPE_BOX:
		// 箱剛体
		btCollisionShape = bt.NewBtBoxShape(
			bt.NewBtVector3(float32(size.GetX()), float32(size.GetY()), float32(size.GetZ())))
	case SHAPE_CAPSULE:
		// カプセル剛体
		btCollisionShape = bt.NewBtCapsuleShape(float32(size.GetX()), float32(size.GetY()))
	}
	// btCollisionShape.SetMargin(0.0001)

	// r.CorrectPhysicsType = r.PhysicsType
	// if r.PhysicsType == PHYSICS_TYPE_DYNAMIC_BONE && r.BoneIndex < 0 {
	// 	// ボーン追従 + 物理剛体の場合、ボーンIndexが設定されていない場合は物理剛体に変更
	// 	r.CorrectPhysicsType = PHYSICS_TYPE_DYNAMIC
	// }

	// 質量
	mass := float32(0.0)
	localInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
	if r.PhysicsType != PHYSICS_TYPE_STATIC {
		// ボーン追従ではない場合そのまま設定
		mass = float32(r.RigidBodyParam.Mass)
	}
	if mass != 0 {
		// 質量が設定されている場合、慣性を計算
		btCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	// // ボーンのローカル位置
	// boneTransform := bt.NewBtTransform()
	// boneTransform.SetIdentity()
	// boneTransform.SetOrigin(boneLocalPosition.Bullet())

	// 剛体の初期位置と回転
	btRigidBodyTransform := bt.NewBtTransform(mbt.MRotationBullet(r.Rotation), mbt.MVec3Bullet(r.Position))

	// ボーンから見た剛体の初期位置
	var bPos *mmath.MVec3
	if r.Bone != nil {
		bPos = r.Bone.Position
	} else if r.JointedBone != nil {
		bPos = r.JointedBone.Position
	} else {
		bPos = mmath.NewMVec3()
	}
	rbLocalPos := r.Position.Subed(bPos)
	btRigidBodyLocalTransform := bt.NewBtTransform(mbt.MRotationBullet(r.Rotation), mbt.MVec3Bullet(rbLocalPos))

	// {
	// 	mlog.V("---------------------------------")
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	r.BtRigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	mlog.V("1. [%s] BtRigidBodyTransform: \n%v\n", r.Name, mat)
	// }

	// 剛体のグローバル位置と回転
	motionState := bt.NewBtDefaultMotionState(btRigidBodyTransform)

	btRigidBody := bt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	btRigidBody.SetDamping(float32(r.RigidBodyParam.LinearDamping), float32(r.RigidBodyParam.AngularDamping))
	btRigidBody.SetRestitution(float32(r.RigidBodyParam.Restitution))
	btRigidBody.SetFriction(float32(r.RigidBodyParam.Friction))
	btRigidBody.SetUserIndex(r.Index)
	// btRigidBody.SetSleepingThresholds(0.1, (180.0 * 0.1 / math.Pi))

	// mlog.V("name: %s, group: %d, mask: %d\n", r.Name, r.CollisionGroup, r.CollisionGroupMaskValue)

	// modelPhysics.AddNonFilterProxy(btRigidBody.GetBroadphaseProxy())
	// 剛体・剛体グループ・非衝突グループを追加
	modelPhysics.AddRigidBody(btRigidBody, btRigidBodyLocalTransform, modelIndex, r.Index,
		1<<r.CollisionGroup, r.CollisionGroupMaskValue)

	r.UpdateFlags(modelIndex, modelPhysics, true, false)
}

func (r *RigidBody) UpdateTransform(
	modelIndex int,
	modelPhysics *mbt.MPhysics,
	rigidBodyBone *Bone,
	boneTransform bt.BtTransform,
) {
	btRigidBody, btRigidBodyLocalTransform := modelPhysics.GetRigidBody(modelIndex, r.Index)

	if btRigidBody == nil || btRigidBody.GetMotionState() == nil {
		return
	}

	// {
	// 	mlog.V("----------")
	// }

	// 剛体のグローバル位置を確定
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	t := bt.NewBtTransform()
	defer bt.DeleteBtTransform(t)

	t.Mult(boneTransform, btRigidBodyLocalTransform)
	motionState.SetWorldTransform(t)

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

}

func (r *RigidBody) GetRigidBodyBoneMatrix(
	modelIndex int,
	modelPhysics *mbt.MPhysics,
) *mmath.MMat4 {
	btRigidBody, btRigidBodyLocalTransform := modelPhysics.GetRigidBody(modelIndex, r.Index)

	if btRigidBody == nil || btRigidBody.GetMotionState() == nil {
		return nil
	}

	// if r.Name == "前髪" {
	// 	state := btRigidBody.GetActivationState()
	// 	isStatic := btRigidBody.IsStaticOrKinematicObject()
	// 	flags := btRigidBody.GetCollisionFlags()
	// 	mlog.V("name: %s, state: %v, static: %v, flags: %v, isKinematic: %v, isStatic: %v\n", r.Name, state, isStatic, flags, flags&int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT) != 0, flags&int(bt.BtCollisionObjectCF_STATIC_OBJECT) != 0)
	// }

	// {
	// 	mlog.V("----------")
	// }

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	rigidBodyGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(rigidBodyGlobalTransform)

	motionState.GetWorldTransform(rigidBodyGlobalTransform)

	// if r.CorrectPhysicsType == PHYSICS_TYPE_DYNAMIC_BONE {
	// 	{
	// 		mat := mgl32.Mat4{}
	// 		rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 		mlog.V("3. [%s] rigidBodyTransform Before: \n%v\n", r.Name, mat)
	// 	}

	// 	if r.BoneIndex >= 0 && r.BoneIndex < len(boneTransforms) {
	// 		// 物理+ボーン追従はボーン移動成分を現在のボーン位置にする
	// 		boneGlobalTransform := bt.NewBtTransform()
	// 		boneGlobalTransform.Mult(*boneTransforms[r.BoneIndex], r.BtRigidBodyTransform)

	// 		rigidBodyTransform.SetOrigin(boneGlobalTransform.GetOrigin().(bt.BtVector3))
	// 	} else if r.JointedBoneIndex >= 0 && r.JointedBoneIndex < len(boneTransforms) {
	// 		// ジョイントで繋がっているボーンがある場合はそのボーン位置にする
	// 		boneGlobalTransform := bt.NewBtTransform()
	// 		boneGlobalTransform.Mult(*boneTransforms[r.JointedBoneIndex], r.BtRigidBodyTransform)

	// 		rigidBodyTransform.SetOrigin(boneGlobalTransform.GetOrigin().(bt.BtVector3))
	// 	}

	// 	{
	// 		mat := mgl32.Mat4{}
	// 		rigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 		mlog.V("3. [%s] rigidBodyTransform After: \n%v\n", r.Name, mat)
	// 	}
	// }

	// ボーンのグローバル位置を剛体の現在グローバル行列に初期位置ローカル行列を掛けて求める
	boneGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneGlobalTransform)

	invRigidBodyLocalTransform := btRigidBodyLocalTransform.Inverse()
	defer bt.DeleteBtTransform(invRigidBodyLocalTransform)

	boneGlobalTransform.Mult(rigidBodyGlobalTransform, invRigidBodyLocalTransform)

	boneGlobalMatrixGL := mgl32.Mat4{}
	boneGlobalTransform.GetOpenGLMatrix(&boneGlobalMatrixGL[0])

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

	return mmath.NewMMat4ByMgl(&boneGlobalMatrixGL)
}

func (r *RigidBodies) initPhysics(modelIndex int, physics *mbt.MPhysics) {
	// 剛体を順番にボーンと紐付けていく
	for _, rigidBody := range r.Data {
		// 物理設定の初期化
		rigidBody.initPhysics(modelIndex, physics)
	}
}
