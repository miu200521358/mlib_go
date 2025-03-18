//go:build windows
// +build windows

package mbt

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// rigidbodyValue は剛体の物理エンジン内部表現を格納する構造体です
type rigidbodyValue struct {
	pmxRigidBody     *pmx.RigidBody  // PMXモデルの剛体定義
	btRigidBody      bt.BtRigidBody  // Bullet物理エンジンの剛体
	btLocalTransform *bt.BtTransform // 剛体のローカルトランスフォーム
	mask             int             // 衝突マスク
	group            int             // 衝突グループ
}

// initRigidBodies はモデルの剛体を初期化します
func (physics *MPhysics) initRigidBodies(modelIndex int, rigidBodies *pmx.RigidBodies) {
	// 剛体を順番にボーンと紐付けていく
	physics.rigidBodies[modelIndex] = make([]*rigidbodyValue, rigidBodies.Length())
	rigidBodies.ForEach(func(index int, rigidBody *pmx.RigidBody) {
		// 剛体の初期位置と回転
		btRigidBodyTransform := bt.NewBtTransform(newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBody.Position))

		// 物理設定の初期化
		physics.initRigidBody(modelIndex, rigidBody, btRigidBodyTransform)
	})
}

// initRigidBodiesByBoneDeltas はボーンデルタ情報を使用して剛体を初期化します
func (physics *MPhysics) initRigidBodiesByBoneDeltas(
	modelIndex int, rigidBodies *pmx.RigidBodies, boneDeltas *delta.BoneDeltas,
) {
	// 剛体を順番にボーンと紐付けていく
	physics.rigidBodies[modelIndex] = make([]*rigidbodyValue, rigidBodies.Length())
	rigidBodies.ForEach(func(index int, rigidBody *pmx.RigidBody) {
		// ボーンから見た剛体の初期位置
		var bone *pmx.Bone
		if rigidBody.Bone != nil {
			bone = rigidBody.Bone
		}

		// 剛体の初期位置と回転
		if bone == nil || !boneDeltas.Contains(bone.Index()) {
			return
		}

		btRigidBodyTransform := bt.NewBtTransform()
		boneTransform := bt.NewBtTransform()
		defer bt.DeleteBtTransform(boneTransform)

		mat := mmath.NewGlMat4(boneDeltas.Get(bone.Index()).FilledGlobalMatrix())
		boneTransform.SetFromOpenGLMatrix(&mat[0])

		rigidBodyLocalPos := rigidBody.Position.Subed(bone.Position)
		btRigidBodyLocalTransform := bt.NewBtTransform(newBulletFromRad(rigidBody.Rotation),
			newBulletFromVec(rigidBodyLocalPos))
		defer bt.DeleteBtTransform(btRigidBodyLocalTransform)

		btRigidBodyTransform.Mult(boneTransform, btRigidBodyLocalTransform)

		// 物理設定の初期化
		physics.initRigidBody(modelIndex, rigidBody, btRigidBodyTransform)
	})
}

// initRigidBody は個別の剛体を初期化します
func (physics *MPhysics) initRigidBody(
	modelIndex int, rigidBody *pmx.RigidBody, btRigidBodyTransform bt.BtTransform,
) {
	// 剛体の形状に基づいた衝突形状の生成
	btCollisionShape := physics.createCollisionShape(rigidBody)

	// 質量と慣性の計算
	mass, localInertia := physics.calculateMassAndInertia(rigidBody, btCollisionShape)

	// ボーンから見た剛体の初期位置
	bonePos := physics.getBonePosition(rigidBody)

	// 剛体のローカルトランスフォーム計算
	rigidBodyLocalPos := rigidBody.Position.Subed(bonePos)
	btRigidBodyLocalTransform := bt.NewBtTransform(
		newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBodyLocalPos))

	// 剛体のグローバル位置と回転
	motionState := bt.NewBtDefaultMotionState(btRigidBodyTransform)

	// 剛体の生成と物理パラメータの設定
	btRigidBody := bt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	physics.configureRigidBody(btRigidBody, rigidBody)

	// 剛体・剛体グループ・非衝突グループを追加
	group := 1 << rigidBody.CollisionGroup
	physics.world.AddRigidBody(btRigidBody, group, rigidBody.CollisionGroupMaskValue)
	physics.rigidBodies[modelIndex][rigidBody.Index()] = &rigidbodyValue{
		pmxRigidBody: rigidBody, btRigidBody: btRigidBody, btLocalTransform: &btRigidBodyLocalTransform,
		mask: rigidBody.CollisionGroupMaskValue, group: group}

	// 剛体の物理フラグを更新
	physics.updateFlag(modelIndex, rigidBody)
}

// createCollisionShape は剛体の形状に基づいた衝突形状を生成します
func (physics *MPhysics) createCollisionShape(rigidBody *pmx.RigidBody) bt.BtCollisionShape {
	// マイナスサイズは許容しない
	size := rigidBody.Size.Clamped(mmath.MVec3Zero, mmath.MVec3MaxVal)

	switch rigidBody.ShapeType {
	case pmx.SHAPE_SPHERE:
		// 球剛体
		return bt.NewBtSphereShape(float32(size.X))
	case pmx.SHAPE_BOX:
		// 箱剛体
		return bt.NewBtBoxShape(
			bt.NewBtVector3(float32(size.X), float32(size.Y), float32(size.Z)))
	case pmx.SHAPE_CAPSULE:
		// カプセル剛体
		return bt.NewBtCapsuleShape(float32(size.X), float32(size.Y))
	default:
		// デフォルトは球
		return bt.NewBtSphereShape(float32(size.X))
	}
}

// calculateMassAndInertia は剛体の質量と慣性を計算します
func (physics *MPhysics) calculateMassAndInertia(rigidBody *pmx.RigidBody, btCollisionShape bt.BtCollisionShape) (float32, bt.BtVector3) {
	// 質量
	mass := float32(0.0)
	localInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))

	if rigidBody.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
		// ボーン追従ではない場合そのまま設定
		mass = float32(mmath.Clamped(rigidBody.RigidBodyParam.Mass, 0, math.MaxFloat64))
	}

	if mass != 0 {
		// 質量が設定されている場合、慣性を計算
		btCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	return mass, localInertia
}

// getBonePosition はボーンの位置を取得します
func (physics *MPhysics) getBonePosition(rigidBody *pmx.RigidBody) *mmath.MVec3 {
	if rigidBody.Bone != nil {
		return rigidBody.Bone.Position
	}
	return mmath.NewMVec3()
}

// configureRigidBody は剛体の物理パラメータを設定します
func (physics *MPhysics) configureRigidBody(btRigidBody bt.BtRigidBody, rigidBody *pmx.RigidBody) {
	btRigidBody.SetDamping(float32(rigidBody.RigidBodyParam.LinearDamping),
		float32(rigidBody.RigidBodyParam.AngularDamping))
	btRigidBody.SetRestitution(float32(rigidBody.RigidBodyParam.Restitution))
	btRigidBody.SetFriction(float32(rigidBody.RigidBodyParam.Friction))
	btRigidBody.SetUserIndex(rigidBody.Index())
}

// deleteRigidBodies はモデルの全剛体を削除します
func (physics *MPhysics) deleteRigidBodies(modelIndex int) {
	for _, r := range physics.rigidBodies[modelIndex] {
		if r == nil || r.btRigidBody == nil {
			continue
		}
		physics.world.RemoveRigidBody(r.btRigidBody)
		bt.DeleteBtRigidBody(r.btRigidBody)
	}
	physics.rigidBodies[modelIndex] = nil
}

// updateFlag は剛体の物理フラグを更新します
func (physics *MPhysics) updateFlag(modelIndex int, rigidBody *pmx.RigidBody) {
	btRigidBody := physics.rigidBodies[modelIndex][rigidBody.Index()].btRigidBody

	if rigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
		// 剛体の位置更新に物理演算を使わない。もしくは物理演算OFF時
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		btRigidBody.SetCollisionFlags(
			btRigidBody.GetCollisionFlags() | int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を無効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが無効になる(剛体位置は判っているので不要)
		btRigidBody.SetActivationState(bt.DISABLE_SIMULATION)
	} else {
		// 物理演算・物理+ボーン位置合わせの場合
		// 剛体の位置更新に物理演算を使う。
		// MotionState::getWorldTransformが毎ステップコールされるようになるのでここで剛体位置を更新する。
		btRigidBody.SetCollisionFlags(btRigidBody.GetCollisionFlags() &
			^int(bt.BtCollisionObjectCF_NO_CONTACT_RESPONSE) & ^int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT))
		// 毎ステップの剛体位置通知を有効にする
		// MotionState::setWorldTransformの毎ステップ呼び出しが有効になる(剛体位置が変わるので必要)
		btRigidBody.SetActivationState(bt.DISABLE_DEACTIVATION)
	}
}

// UpdateTransform はボーン行列に基づいて剛体の位置を更新します
func (physics *MPhysics) UpdateTransform(
	modelIndex int,
	rigidBodyBone *pmx.Bone,
	boneGlobalMatrix *mmath.MMat4,
	rigidBody *pmx.RigidBody,
) {
	boneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneTransform)

	mat := mmath.NewGlMat4(boneGlobalMatrix)
	boneTransform.SetFromOpenGLMatrix(&mat[0])

	btRigidBody := physics.rigidBodies[modelIndex][rigidBody.Index()].btRigidBody
	btRigidBodyLocalTransform := *physics.rigidBodies[modelIndex][rigidBody.Index()].btLocalTransform

	// 剛体のグローバル位置を確定
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	t := bt.NewBtTransform()
	defer bt.DeleteBtTransform(t)

	t.Mult(boneTransform, btRigidBodyLocalTransform)
	motionState.SetWorldTransform(t)
}

// GetRigidBodyBoneMatrix は剛体に基づいてボーン行列を取得します
func (physics *MPhysics) GetRigidBodyBoneMatrix(
	modelIndex int,
	rigidBody *pmx.RigidBody,
) *mmath.MMat4 {
	btRigidBody := physics.rigidBodies[modelIndex][rigidBody.Index()].btRigidBody
	btRigidBodyLocalTransform := *physics.rigidBodies[modelIndex][rigidBody.Index()].btLocalTransform

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	rigidBodyGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(rigidBodyGlobalTransform)

	motionState.GetWorldTransform(rigidBodyGlobalTransform)

	// ボーンのグローバル位置を剛体の現在グローバル行列に初期位置ローカル行列を掛けて求める
	boneGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneGlobalTransform)

	invRigidBodyLocalTransform := btRigidBodyLocalTransform.Inverse()
	defer bt.DeleteBtTransform(invRigidBodyLocalTransform)

	boneGlobalTransform.Mult(rigidBodyGlobalTransform, invRigidBodyLocalTransform)

	boneGlobalMatrixGL := mgl32.Mat4{}
	boneGlobalTransform.GetOpenGLMatrix(&boneGlobalMatrixGL[0])

	return newMMat4ByMgl(&boneGlobalMatrixGL)
}
