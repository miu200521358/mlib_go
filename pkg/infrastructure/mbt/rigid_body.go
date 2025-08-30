//go:build windows
// +build windows

package mbt

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
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
func (mp *MPhysics) initRigidBodies(modelIndex int, rigidBodies *pmx.RigidBodies) {
	// 剛体を順番にボーンと紐付けていく
	mp.rigidBodies[modelIndex] = make([]*rigidbodyValue, rigidBodies.Length())
	rigidBodies.ForEach(func(index int, rigidBody *pmx.RigidBody) bool {
		// 剛体の初期位置と回転
		btRigidBodyTransform := bt.NewBtTransform(newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBody.Position))

		// 物理設定の初期化
		mp.initRigidBody(modelIndex, rigidBody, btRigidBodyTransform, nil)

		return true
	})
}

// initRigidBodiesByBoneDeltas はボーンデルタ情報を使用して剛体を初期化します
func (mp *MPhysics) initRigidBodiesByBoneDeltas(
	modelIndex int, rigidBodies *pmx.RigidBodies,
	boneDeltas *delta.BoneDeltas, rigidBodyDeltas *delta.RigidBodyDeltas,
) {
	// 剛体を順番にボーンと紐付けていく
	mp.rigidBodies[modelIndex] = make([]*rigidbodyValue, rigidBodies.Length())
	rigidBodies.ForEach(func(index int, rigidBody *pmx.RigidBody) bool {
		// ボーンから見た剛体の初期位置
		var bone *pmx.Bone
		if rigidBody.Bone != nil {
			bone = rigidBody.Bone
		}

		// 剛体の初期位置と回転
		if bone == nil || !boneDeltas.Contains(bone.Index()) {
			return true
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

		var rigidBodyDelta *delta.RigidBodyDelta
		if rigidBodyDeltas != nil {
			rigidBodyDelta = rigidBodyDeltas.Get(rigidBody.Index())
		}

		// 物理設定の初期化
		mp.initRigidBody(modelIndex, rigidBody, btRigidBodyTransform, rigidBodyDelta)

		return true
	})
}

// initRigidBody は個別の剛体を初期化します
func (mp *MPhysics) initRigidBody(
	modelIndex int, rigidBody *pmx.RigidBody, btRigidBodyTransform bt.BtTransform,
	rigidBodyDelta *delta.RigidBodyDelta,
) {
	// 剛体の形状に基づいた衝突形状の生成
	btCollisionShape := mp.createCollisionShape(rigidBody, rigidBodyDelta)

	// 質量と慣性の計算
	mass, localInertia := mp.calculateMassAndInertia(rigidBody, btCollisionShape, rigidBodyDelta)

	// ボーンから見た剛体の初期位置
	bonePos := mp.getBonePosition(rigidBody)

	// 剛体のローカルトランスフォーム計算
	rigidBodyLocalPos := rigidBody.Position.Subed(bonePos)
	btRigidBodyLocalTransform := bt.NewBtTransform(
		newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBodyLocalPos))

	// 剛体のグローバル位置と回転
	motionState := bt.NewBtDefaultMotionState(btRigidBodyTransform)

	// 剛体の生成と物理パラメータの設定
	btRigidBody := bt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	mp.configureRigidBody(btRigidBody, rigidBody)

	// 剛体・剛体グループ・非衝突グループを追加
	group := 1 << rigidBody.CollisionGroup
	mp.world.AddRigidBody(btRigidBody, group, rigidBody.CollisionGroupMaskValue)
	mp.rigidBodies[modelIndex][rigidBody.Index()] = &rigidbodyValue{
		pmxRigidBody: rigidBody, btRigidBody: btRigidBody, btLocalTransform: &btRigidBodyLocalTransform,
		mask: rigidBody.CollisionGroupMaskValue, group: group}

	// 剛体の物理フラグを更新
	mp.updateFlag(modelIndex, rigidBody)
}

// createCollisionShape は剛体の形状に基づいた衝突形状を生成します
func (mp *MPhysics) createCollisionShape(
	rigidBody *pmx.RigidBody, rigidBodyDelta *delta.RigidBodyDelta,
) bt.BtCollisionShape {
	// マイナスサイズは許容しない
	size := rigidBody.Size
	if rigidBodyDelta != nil {
		size = rigidBody.Size.Muled(rigidBodyDelta.Size)
	}
	size.Clamp(mmath.MVec3Zero, mmath.MVec3MaxVal)

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
func (mp *MPhysics) calculateMassAndInertia(
	rigidBody *pmx.RigidBody, btCollisionShape bt.BtCollisionShape, rigidBodyDelta *delta.RigidBodyDelta,
) (float32, bt.BtVector3) {
	// 質量
	mass := float32(0.0)
	localInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))

	if rigidBody.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
		// ボーン追従ではない場合そのまま設定
		if rigidBodyDelta != nil {
			// 剛体デルタがある場合はその値を使用
			mass = float32(mmath.Clamped(rigidBodyDelta.Mass, 0, math.MaxFloat64))
		} else {
			mass = float32(mmath.Clamped(rigidBody.RigidBodyParam.Mass, 0, math.MaxFloat64))
		}
	}

	if mass != 0 {
		// 質量が設定されている場合、慣性を計算
		btCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	return mass, localInertia
}

// getBonePosition はボーンの位置を取得します
func (mp *MPhysics) getBonePosition(rigidBody *pmx.RigidBody) *mmath.MVec3 {
	if rigidBody.Bone != nil {
		return rigidBody.Bone.Position
	}
	return mmath.NewMVec3()
}

// configureRigidBody は剛体の物理パラメータを設定します
func (mp *MPhysics) configureRigidBody(btRigidBody bt.BtRigidBody, rigidBody *pmx.RigidBody) {
	btRigidBody.SetDamping(float32(rigidBody.RigidBodyParam.LinearDamping),
		float32(rigidBody.RigidBodyParam.AngularDamping))
	btRigidBody.SetRestitution(float32(rigidBody.RigidBodyParam.Restitution))
	btRigidBody.SetFriction(float32(rigidBody.RigidBodyParam.Friction))
	btRigidBody.SetUserIndex(rigidBody.Index())
}

// deleteRigidBodies はモデルの全剛体を削除します
func (mp *MPhysics) deleteRigidBodies(modelIndex int) {
	for _, r := range mp.rigidBodies[modelIndex] {
		if r == nil || r.btRigidBody == nil {
			continue
		}
		mp.world.RemoveRigidBody(r.btRigidBody)
		bt.DeleteBtRigidBody(r.btRigidBody)
	}
	mp.rigidBodies[modelIndex] = nil
}

// updateFlag は剛体の物理フラグを更新します
func (mp *MPhysics) updateFlag(modelIndex int, rigidBody *pmx.RigidBody) {
	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].btRigidBody

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
func (mp *MPhysics) UpdateTransform(
	modelIndex int,
	rigidBodyBone *pmx.Bone,
	boneGlobalMatrix *mmath.MMat4,
	rigidBody *pmx.RigidBody,
) {
	boneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneTransform)

	mat := mmath.NewGlMat4(boneGlobalMatrix)
	boneTransform.SetFromOpenGLMatrix(&mat[0])

	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].btRigidBody
	btRigidBodyLocalTransform := *mp.rigidBodies[modelIndex][rigidBody.Index()].btLocalTransform

	// 剛体のグローバル位置を確定
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	t := bt.NewBtTransform()
	defer bt.DeleteBtTransform(t)

	t.Mult(boneTransform, btRigidBodyLocalTransform)
	motionState.SetWorldTransform(t)
}

// GetRigidBodyBoneMatrix は剛体に基づいてボーン行列を取得します
func (mp *MPhysics) GetRigidBodyBoneMatrix(
	modelIndex int,
	rigidBody *pmx.RigidBody,
) *mmath.MMat4 {
	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].btRigidBody
	btRigidBodyLocalTransform := *mp.rigidBodies[modelIndex][rigidBody.Index()].btLocalTransform

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

// SaveRigidBodyStates は全剛体の物理状態を保存します
func (mp *MPhysics) SaveRigidBodyStates(modelIndex int) map[int]*physics.RigidbodyState {
	states := make(map[int]*physics.RigidbodyState)

	for rigidBodyIndex, r := range mp.rigidBodies[modelIndex] {
		if r == nil || r.btRigidBody == nil {
			continue
		}

		btRigidBody := r.btRigidBody

		// 線形速度取得
		linearVel := btRigidBody.GetLinearVelocity()
		linearVelocity := &mmath.MVec3{
			X: float64(linearVel.X()),
			Y: float64(linearVel.Y()),
			Z: float64(linearVel.Z()),
		}

		// 角速度取得
		angularVel := btRigidBody.GetAngularVelocity()
		angularVelocity := &mmath.MVec3{
			X: float64(angularVel.X()),
			Y: float64(angularVel.Y()),
			Z: float64(angularVel.Z()),
		}

		// ワールド変換行列取得
		motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
		worldTransform := bt.NewBtTransform()
		defer bt.DeleteBtTransform(worldTransform)
		motionState.GetWorldTransform(worldTransform)

		worldMatrixGL := mgl32.Mat4{}
		worldTransform.GetOpenGLMatrix(&worldMatrixGL[0])
		worldMatrix := newMMat4ByMgl(&worldMatrixGL)

		states[rigidBodyIndex] = &physics.RigidbodyState{
			LinearVelocity:  linearVelocity,
			AngularVelocity: angularVelocity,
			WorldTransform:  worldMatrix,
			IsActive:        btRigidBody.IsActive(),
		}
	}

	return states
}

// RestoreRigidBodyStates は保存された剛体の物理状態を復元します
func (mp *MPhysics) RestoreRigidBodyStates(modelIndex int, states map[int]*physics.RigidbodyState) {
	for rigidBodyIndex, state := range states {
		if rigidBodyIndex >= len(mp.rigidBodies[modelIndex]) {
			continue
		}

		r := mp.rigidBodies[modelIndex][rigidBodyIndex]
		if r == nil || r.btRigidBody == nil {
			continue
		}

		btRigidBody := r.btRigidBody
		pmxRigidBody := r.pmxRigidBody

		// サイズ変更後の慣性テンソル再計算（物理剛体のみ）
		if pmxRigidBody.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
			currentShape := btRigidBody.GetCollisionShape().(bt.BtCollisionShape)

			// 現在の質量を取得
			var currentMass float32
			if btRigidBody.GetInvMass() != 0 {
				currentMass = 1.0 / btRigidBody.GetInvMass()
			}

			if currentMass > 0 {
				// 新しい慣性テンソル計算
				currentShape.CalculateLocalInertia(currentMass,
					bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0)))
				// 質量と慣性を再設定
				btRigidBody.SetMassProps(currentMass,
					bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0)))
			}
		}

		// ボーン追従剛体は位置のみ復元、速度はゼロに
		if pmxRigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
			// 位置復元
			mat := mmath.NewGlMat4(state.WorldTransform)
			worldTransform := bt.NewBtTransform()
			defer bt.DeleteBtTransform(worldTransform)
			worldTransform.SetFromOpenGLMatrix(&mat[0])

			motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
			motionState.SetWorldTransform(worldTransform)

			// 速度をゼロに設定
			btRigidBody.SetLinearVelocity(bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0)))
			btRigidBody.SetAngularVelocity(bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0)))

			// 強制的に非アクティブにして物理演算の影響を受けないようにする
			btRigidBody.ForceActivationState(bt.DISABLE_SIMULATION)
		} else {
			// 物理剛体は全状態を復元
			// 位置復元
			mat := mmath.NewGlMat4(state.WorldTransform)
			worldTransform := bt.NewBtTransform()
			defer bt.DeleteBtTransform(worldTransform)
			worldTransform.SetFromOpenGLMatrix(&mat[0])

			motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
			motionState.SetWorldTransform(worldTransform)

			// 速度復元
			btRigidBody.SetLinearVelocity(bt.NewBtVector3(
				float32(state.LinearVelocity.X),
				float32(state.LinearVelocity.Y),
				float32(state.LinearVelocity.Z),
			))
			btRigidBody.SetAngularVelocity(bt.NewBtVector3(
				float32(state.AngularVelocity.X),
				float32(state.AngularVelocity.Y),
				float32(state.AngularVelocity.Z),
			))

			// アクティブ状態復元
			if state.IsActive {
				btRigidBody.SetActivationState(bt.ACTIVE_TAG)
			}
		}
	}
}
