//go:build windows
// +build windows

package pmx

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
)

func (r *RigidBody) UpdateFlags(modelPhysics *mphysics.MPhysics, enablePhysics bool) bool {
	btRigidBody, _ := modelPhysics.GetRigidBody(r.Index)

	if r.CorrectPhysicsType == PHYSICS_TYPE_STATIC {
		// 剛体の位置更新に物理演算を使わない。もしくは物理演算OFF時
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		btRigidBody.SetCollisionFlags(
			btRigidBody.GetCollisionFlags() | int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を無効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが無効になる(剛体位置は判っているので不要)
		btRigidBody.SetActivationState(mbt.DISABLE_SIMULATION)
		// if prevActivationState != mbt.DISABLE_SIMULATION {
		// 	return true
		// }
	} else {
		prevActivationState := btRigidBody.GetActivationState()

		if enablePhysics {
			// 物理演算・物理+ボーン位置合わせの場合
			if prevActivationState != mbt.ACTIVE_TAG {
				btRigidBody.SetCollisionFlags(0 & ^int(mbt.BtCollisionObjectCF_NO_CONTACT_RESPONSE))
				btRigidBody.ForceActivationState(mbt.ACTIVE_TAG)

				// localInertia := mbt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
				// btRigidBody.GetCollisionShape().(mbt.BtCollisionShape).CalculateLocalInertia(
				// 	float32(r.RigidBodyParam.Mass), localInertia)

				// btRigidBody.GetMotionState().(mbt.BtMotionState).SetWorldTransform(r.BtRigidBodyTransform)

				return true
			} else {
				// 剛体の位置更新に物理演算を使う。
				// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
				btRigidBody.SetCollisionFlags(
					btRigidBody.GetCollisionFlags() & ^int(mbt.BtCollisionObjectCF_NO_CONTACT_RESPONSE))
				// 毎ステップの剛体位置通知を有効にする
				// MotionState::setWorldTransformの毎ステップ呼び出しが有効になる(剛体位置が変わるので必要)
				btRigidBody.SetActivationState(mbt.ACTIVE_TAG)
			}
		} else {
			// 物理OFF時
			btRigidBody.SetCollisionFlags(
				btRigidBody.GetCollisionFlags() | int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT))
			btRigidBody.SetActivationState(mbt.ISLAND_SLEEPING)

			if prevActivationState != mbt.ISLAND_SLEEPING {
				return true
			}
		}
		// if prevActivationState != mbt.ACTIVE_TAG {
		// 	// btRigidBody.ForceActivationState(mbt.ACTIVE_TAG)
		// 	// btRigidBody.Activate(true)
		// 	modelPhysics.RemoveRigidBody(r.BtRigidBody)
		// 	r.InitPhysics(modelPhysics)
		// 	return true
		// } else {
		// 	btRigidBody.SetActivationState(mbt.ACTIVE_TAG)
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
	// btCollisionShape.SetMargin(0.0001)

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
	btRigidBodyTransform := mbt.NewBtTransform(r.Rotation.Bullet(), r.Position.Bullet())

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
	btRigidBodyLocalTransform := mbt.NewBtTransform(r.Rotation.Bullet(), rbLocalPos.Bullet())

	// {
	// 	mlog.V("---------------------------------")
	// }
	// {
	// 	mat := mgl32.Mat4{}
	// 	r.BtRigidBodyTransform.GetOpenGLMatrix(&mat[0])
	// 	mlog.V("1. [%s] BtRigidBodyTransform: \n%v\n", r.Name, mat)
	// }

	// 剛体のグローバル位置と回転
	motionState := mbt.NewBtDefaultMotionState(btRigidBodyTransform)

	btRigidBody := mbt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	btRigidBody.SetDamping(float32(r.RigidBodyParam.LinearDamping), float32(r.RigidBodyParam.AngularDamping))
	btRigidBody.SetRestitution(float32(r.RigidBodyParam.Restitution))
	btRigidBody.SetFriction(float32(r.RigidBodyParam.Friction))
	btRigidBody.SetUserIndex(r.Index)
	// btRigidBody.SetSleepingThresholds(0.1, (180.0 * 0.1 / math.Pi))

	// mlog.V("name: %s, group: %d, mask: %d\n", r.Name, r.CollisionGroup, r.CollisionGroupMaskValue)

	// modelPhysics.AddNonFilterProxy(btRigidBody.GetBroadphaseProxy())
	// 剛体・剛体グループ・非衝突グループを追加
	modelPhysics.AddRigidBody(btRigidBody, btRigidBodyLocalTransform, r.Index,
		1<<r.CollisionGroup, r.CollisionGroupMaskValue)

	r.UpdateFlags(modelPhysics, true)
}

func (r *RigidBody) UpdateTransform(
	modelPhysics *mphysics.MPhysics,
	rigidBodyBone *Bone,
	boneTransform mbt.BtTransform,
	isForce bool,
) {
	btRigidBody, btRigidBodyLocalTransform := modelPhysics.GetRigidBody(r.Index)

	if btRigidBody == nil || btRigidBody.GetMotionState() == nil || boneTransform == nil ||
		(r.CorrectPhysicsType == PHYSICS_TYPE_DYNAMIC && !isForce) {
		return
	}

	// {
	// 	mlog.V("----------")
	// }

	// 剛体のグローバル位置を確定
	motionState := btRigidBody.GetMotionState().(mbt.BtMotionState)

	t := mbt.NewBtTransform()
	// defer mbt.DeleteBtTransform(t)

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
	modelPhysics *mphysics.MPhysics,
) *mmath.MMat4 {
	btRigidBody, btRigidBodyLocalTransform := modelPhysics.GetRigidBody(r.Index)

	if btRigidBody == nil || btRigidBody.GetMotionState() == nil || r.CorrectPhysicsType == PHYSICS_TYPE_STATIC {
		return nil
	}

	// if r.Name == "前髪" {
	// 	state := btRigidBody.GetActivationState()
	// 	isStatic := btRigidBody.IsStaticOrKinematicObject()
	// 	flags := btRigidBody.GetCollisionFlags()
	// 	mlog.V("name: %s, state: %v, static: %v, flags: %v, isKinematic: %v, isStatic: %v\n", r.Name, state, isStatic, flags, flags&int(mbt.BtCollisionObjectCF_KINEMATIC_OBJECT) != 0, flags&int(mbt.BtCollisionObjectCF_STATIC_OBJECT) != 0)
	// }

	// {
	// 	mlog.V("----------")
	// }

	motionState := btRigidBody.GetMotionState().(mbt.BtMotionState)

	rigidBodyGlobalTransform := mbt.NewBtTransform()
	// defer mbt.DeleteBtTransform(rigidBodyGlobalTransform)

	motionState.GetWorldTransform(rigidBodyGlobalTransform)

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

	// ボーンのグローバル位置を剛体の現在グローバル行列に初期位置ローカル行列を掛けて求める
	boneGlobalTransform := mbt.NewBtTransform()
	// defer mbt.DeleteBtTransform(boneGlobalTransform)

	boneGlobalTransform.Mult(rigidBodyGlobalTransform, btRigidBodyLocalTransform.Inverse())

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

func (r *RigidBodies) initPhysics(physics *mphysics.MPhysics) {
	// 剛体を順番にボーンと紐付けていく
	for _, rigidBody := range r.GetSortedData() {
		// 物理設定の初期化
		rigidBody.initPhysics(physics)
	}
}
