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

// initRigidBodies はモデルの剛体を初期化します
func (mp *MPhysics) initRigidBodies(modelIndex int, rigidBodies *pmx.RigidBodies) {
	// 剛体を順番にボーンと紐付けていく
	mp.rigidBodies[modelIndex] = make([]*physics.RigidBodyValue, rigidBodies.Length())
	rigidBodies.ForEach(func(index int, rigidBody *pmx.RigidBody) bool {
		// 剛体の初期位置と回転
		btRigidBodyTransform := bt.NewBtTransform(NewBulletFromRad(rigidBody.Rotation), NewBulletFromVec(rigidBody.Position))

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
	mp.rigidBodies[modelIndex] = make([]*physics.RigidBodyValue, rigidBodies.Length())
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
		btRigidBodyLocalTransform := bt.NewBtTransform(NewBulletFromRad(rigidBody.Rotation),
			NewBulletFromVec(rigidBodyLocalPos))
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
		NewBulletFromRad(rigidBody.Rotation), NewBulletFromVec(rigidBodyLocalPos))

	// 剛体のグローバル位置と回転
	motionState := bt.NewBtDefaultMotionState(btRigidBodyTransform)

	// 剛体の生成と物理パラメータの設定
	btRigidBody := bt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	mp.configureRigidBody(btRigidBody, modelIndex, rigidBody)

	// 剛体・剛体グループ・非衝突グループを追加
	group := 1 << rigidBody.CollisionGroup
	mp.world.AddRigidBody(btRigidBody, group, rigidBody.CollisionGroupMaskValue)
	mp.rigidBodies[modelIndex][rigidBody.Index()] = &physics.RigidBodyValue{
		PmxRigidBody: rigidBody, BtRigidBody: btRigidBody, BtLocalTransform: &btRigidBodyLocalTransform,
		Mask: rigidBody.CollisionGroupMaskValue, Group: group}

	// 剛体の物理フラグを更新
	mp.updateFlag(modelIndex, rigidBody)
}

// createCollisionShape は剛体の形状に基づいた衝突形状を生成します
func (mp *MPhysics) createCollisionShape(
	rigidBody *pmx.RigidBody, rigidBodyDelta *delta.RigidBodyDelta,
) bt.BtCollisionShape {
	// マイナスサイズは許容しない
	size := rigidBody.Size
	if rigidBodyDelta != nil && rigidBodyDelta.Size != nil {
		size = rigidBodyDelta.Size.Copy()
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
		if rigidBodyDelta != nil && rigidBodyDelta.Mass != 0.0 {
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
func (mp *MPhysics) configureRigidBody(btRigidBody bt.BtRigidBody, modelIndex int, rigidBody *pmx.RigidBody) {
	btRigidBody.SetDamping(float32(rigidBody.RigidBodyParam.LinearDamping),
		float32(rigidBody.RigidBodyParam.AngularDamping))
	btRigidBody.SetRestitution(float32(rigidBody.RigidBodyParam.Restitution))
	btRigidBody.SetFriction(float32(rigidBody.RigidBodyParam.Friction))
	btRigidBody.SetUserIndex(modelIndex)
	btRigidBody.SetUserIndex2(rigidBody.Index())
}

// deleteRigidBodies はモデルの全剛体を削除します
func (mp *MPhysics) deleteRigidBodies(modelIndex int) {
	for _, r := range mp.rigidBodies[modelIndex] {
		if r == nil || r.BtRigidBody == nil {
			continue
		}
		mp.world.RemoveRigidBody(r.BtRigidBody)
		bt.DeleteBtRigidBody(r.BtRigidBody)
	}
	mp.rigidBodies[modelIndex] = nil
}

// updateFlag は剛体の物理フラグを更新します
func (mp *MPhysics) updateFlag(modelIndex int, rigidBody *pmx.RigidBody) {
	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody

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

	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody
	btRigidBodyLocalTransform := *mp.rigidBodies[modelIndex][rigidBody.Index()].BtLocalTransform

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
	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody
	btRigidBodyLocalTransform := *mp.rigidBodies[modelIndex][rigidBody.Index()].BtLocalTransform

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

	return NewMMat4ByMgl(&boneGlobalMatrixGL)
}

// UpdateRigidBodiesSelectively は変更が必要な剛体のみを選択的に更新します
func (mp *MPhysics) UpdateRigidBodiesSelectively(
	modelIndex int,
	model *pmx.PmxModel,
	rigidBodyDeltas *delta.RigidBodyDeltas,
) {
	if rigidBodyDeltas == nil {
		return
	}

	// 変更がある剛体のみ更新
	rigidBodyDeltas.ForEach(func(index int, rigidBodyDelta *delta.RigidBodyDelta) bool {
		if rigidBodyDelta == nil {
			return true
		}

		rigidBody, err := model.RigidBodies.Get(index)
		if err != nil || rigidBody == nil {
			return true
		}

		// 個別剛体の形状を更新
		mp.UpdateRigidBodyShapeMass(modelIndex, rigidBody, rigidBodyDelta)

		return true
	})
}

// UpdateRigidBodyShapeMass はサイズ・質量変更時に剛体の形状を更新します
func (mp *MPhysics) UpdateRigidBodyShapeMass(
	modelIndex int,
	rigidBody *pmx.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) {
	if rigidBodyDelta == nil || (rigidBodyDelta.Size == nil && rigidBodyDelta.Mass == 0.0) {
		return
	}

	r := mp.rigidBodies[modelIndex][rigidBody.Index()]
	if r == nil || r.BtRigidBody == nil {
		return
	}

	btRigidBody := r.BtRigidBody

	// 現在の状態を保存
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)
	motionState.GetWorldTransform(currentTransform)

	// 現在の物理パラメータを保存
	currentLinearVel := btRigidBody.GetLinearVelocity()
	currentAngularVel := btRigidBody.GetAngularVelocity()
	currentMass := float32(0.0)
	if btRigidBody.GetInvMass() != 0 {
		currentMass = 1.0 / btRigidBody.GetInvMass()
	}

	// 新しい質量を決定
	newMass := currentMass
	if rigidBodyDelta.Mass != 0.0 && rigidBody.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
		newMass = float32(mmath.Clamped(rigidBodyDelta.Mass, 0, math.MaxFloat64))
	}

	// サイズが変更される場合は形状を再作成
	needShapeUpdate := rigidBodyDelta.Size != nil
	var newShape bt.BtCollisionShape
	if needShapeUpdate {
		// 一旦物理世界から削除
		mp.world.RemoveRigidBody(btRigidBody)

		// 古い形状を取得して削除
		oldShape := btRigidBody.GetCollisionShape()
		if oldShape != nil {
			if collisionShape, ok := oldShape.(bt.BtCollisionShape); ok {
				bt.DeleteBtCollisionShape(collisionShape)
			}
		}

		// 新しい形状を作成
		newShape = mp.createCollisionShape(rigidBody, rigidBodyDelta)
		btRigidBody.SetCollisionShape(newShape)
	}

	// 新しい慣性テンソルを計算
	newInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
	if newMass > 0 && rigidBody.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
		currentShape := btRigidBody.GetCollisionShape()
		if currentShape != nil {
			if collisionShape, ok := currentShape.(bt.BtCollisionShape); ok {
				collisionShape.CalculateLocalInertia(newMass, newInertia)
			}
		}
	}

	// 質量と慣性を更新
	btRigidBody.SetMassProps(newMass, newInertia)

	// 内部状態を更新
	btRigidBody.UpdateInertiaTensor()

	// 位置と速度を復元
	motionState.SetWorldTransform(currentTransform)
	btRigidBody.SetLinearVelocity(currentLinearVel)
	btRigidBody.SetAngularVelocity(currentAngularVel)

	// サイズが変更された場合は物理世界に再追加
	if needShapeUpdate {
		mp.world.AddRigidBody(btRigidBody, r.Group, r.Mask)
	}

	// 物理フラグを更新
	mp.updateFlag(modelIndex, rigidBody)

	// 剛体をアクティブにしてワールドの再計算を促す
	btRigidBody.Activate(true)
}

// FindRigidBodyByCollisionHit はレイキャストで得た btCollisionObject から
// (modelIndex, rigidBodyIndex) を逆引きする。
func (mp *MPhysics) FindRigidBodyByCollisionHit(hitObj bt.BtCollisionObject, hasHit bool) (modelIndex int, rb *physics.RigidBodyValue, ok bool) {
	if hitObj == nil || !hasHit {
		return -1, nil, false
	}
	defer func() {
		if r := recover(); r != nil {
			// SWIG ラッパが 0 ポインタを内包していた場合などを握り潰す
			modelIndex = -1
			rb = nil
			ok = false
		}
	}()

	// モデルIndexと剛体Indexを取得
	modelIndex = hitObj.GetUserIndex()
	rigidBodyIndex := hitObj.GetUserIndex2()

	if _, ok := mp.rigidBodies[modelIndex]; ok {
		// 該当モデルが存在する場合、その剛体Indexを使う
		if len(mp.rigidBodies[modelIndex]) > rigidBodyIndex {
			v := mp.rigidBodies[modelIndex][rigidBodyIndex]
			if v != nil && v.BtRigidBody != nil {
				return modelIndex, v, true
			}
		}
	}

	return -1, nil, false
}
