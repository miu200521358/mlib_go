//go:build windows
// +build windows

// 指示: miu200521358
package mbullet

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet/bt"
)

// initRigidBodies はモデルの剛体を初期化します。
func (mp *PhysicsEngine) initRigidBodies(modelIndex int, pmxModel *model.PmxModel) {
	if pmxModel == nil || pmxModel.RigidBodies == nil {
		return
	}

	rigidBodies := pmxModel.RigidBodies.Values()
	mp.rigidBodies[modelIndex] = make([]*RigidBodyValue, len(rigidBodies))
	for _, rigidBody := range rigidBodies {
		if rigidBody == nil {
			continue
		}
		btRigidBodyTransform := bt.NewBtTransform(newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBody.Position))
		mp.initRigidBody(modelIndex, pmxModel.Bones, rigidBody, btRigidBodyTransform, nil)
	}
}

// initRigidBodiesByBoneDeltas はボーンデルタ情報を使用して剛体を初期化します。
func (mp *PhysicsEngine) initRigidBodiesByBoneDeltas(
	modelIndex int,
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBodyDeltas *delta.RigidBodyDeltas,
) {
	if pmxModel == nil || pmxModel.RigidBodies == nil || pmxModel.Bones == nil || boneDeltas == nil {
		return
	}

	rigidBodies := pmxModel.RigidBodies.Values()
	mp.rigidBodies[modelIndex] = make([]*RigidBodyValue, len(rigidBodies))
	for _, rigidBody := range rigidBodies {
		if rigidBody == nil {
			continue
		}
		bone := mp.getRigidBodyBone(pmxModel.Bones, rigidBody)
		if bone == nil || !boneDeltas.Contains(bone.Index()) {
			continue
		}

		btRigidBodyTransform := bt.NewBtTransform()
		boneTransform := bt.NewBtTransform()
		defer bt.DeleteBtTransform(boneTransform)

		mat := newMglMat4FromMat4(boneDeltas.Get(bone.Index()).FilledGlobalMatrix())
		boneTransform.SetFromOpenGLMatrix(&mat[0])

		rigidBodyLocalPos := rigidBody.Position.Subed(bone.Position)
		btRigidBodyLocalTransform := bt.NewBtTransform(newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBodyLocalPos))
		defer bt.DeleteBtTransform(btRigidBodyLocalTransform)

		btRigidBodyTransform.Mult(boneTransform, btRigidBodyLocalTransform)

		var rigidBodyDelta *delta.RigidBodyDelta
		if rigidBodyDeltas != nil {
			rigidBodyDelta = rigidBodyDeltas.Get(rigidBody.Index())
		}

		mp.initRigidBody(modelIndex, pmxModel.Bones, rigidBody, btRigidBodyTransform, rigidBodyDelta)
	}
}

// initRigidBody は個別の剛体を初期化します。
func (mp *PhysicsEngine) initRigidBody(
	modelIndex int,
	bones *model.BoneCollection,
	rigidBody *model.RigidBody,
	btRigidBodyTransform bt.BtTransform,
	rigidBodyDelta *delta.RigidBodyDelta,
) {
	btCollisionShape := mp.createCollisionShape(rigidBody, rigidBodyDelta)
	mass, localInertia := mp.calculateMassAndInertia(rigidBody, btCollisionShape, rigidBodyDelta)

	bonePos := mp.getBonePosition(bones, rigidBody)
	btRigidBodyLocalTransform := bt.NewBtTransform(
		newBulletFromRad(rigidBody.Rotation),
		newBulletFromVec(rigidBody.Position.Subed(bonePos)),
	)

	motionState := bt.NewBtDefaultMotionState(btRigidBodyTransform)
	btRigidBody := bt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	mp.configureRigidBody(btRigidBody, modelIndex, rigidBody)

	group := 1 << int(rigidBody.CollisionGroup.Group)
	mask := int(rigidBody.CollisionGroup.Mask)
	mp.world.AddRigidBody(btRigidBody, group, mask)
	mp.rigidBodies[modelIndex][rigidBody.Index()] = &RigidBodyValue{
		RigidBody:        rigidBody,
		BtRigidBody:      btRigidBody,
		BtLocalTransform: &btRigidBodyLocalTransform,
		Mask:             mask,
		Group:            group,
	}

	mp.updateFlag(modelIndex, rigidBody)
}

// createCollisionShape は剛体の形状に基づいた衝突形状を生成します。
func (mp *PhysicsEngine) createCollisionShape(
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) bt.BtCollisionShape {
	size := rigidBody.Size
	if rigidBodyDelta != nil {
		size = rigidBodyDelta.Size
	}
	size.Clamp(mmath.ZERO_VEC3, mmath.VEC3_MAX_VAL)

	switch rigidBody.Shape {
	case model.SHAPE_SPHERE:
		return bt.NewBtSphereShape(float32(size.X))
	case model.SHAPE_BOX:
		return bt.NewBtBoxShape(bt.NewBtVector3(float32(size.X), float32(size.Y), float32(size.Z)))
	case model.SHAPE_CAPSULE:
		return bt.NewBtCapsuleShape(float32(size.X), float32(size.Y))
	default:
		return bt.NewBtSphereShape(float32(size.X))
	}
}

// calculateMassAndInertia は剛体の質量と慣性を計算します。
func (mp *PhysicsEngine) calculateMassAndInertia(
	rigidBody *model.RigidBody,
	btCollisionShape bt.BtCollisionShape,
	rigidBodyDelta *delta.RigidBodyDelta,
) (float32, bt.BtVector3) {
	mass := float32(0.0)
	localInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))

	if rigidBody.PhysicsType != model.PHYSICS_TYPE_STATIC {
		baseMass := rigidBody.Param.Mass
		if rigidBodyDelta != nil {
			baseMass = rigidBodyDelta.Mass
		}
		mass = float32(mmath.Clamped(baseMass, 0, math.MaxFloat64))
	}

	if mass != 0 {
		btCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	return mass, localInertia
}

// getRigidBodyBone は剛体に紐づくボーンを取得します。
func (mp *PhysicsEngine) getRigidBodyBone(bones *model.BoneCollection, rigidBody *model.RigidBody) *model.Bone {
	if bones == nil || rigidBody == nil || rigidBody.BoneIndex < 0 {
		return nil
	}
	bone, err := bones.Get(rigidBody.BoneIndex)
	if err != nil {
		return nil
	}
	return bone
}

// getBonePosition は剛体に紐づくボーン位置を取得します。
func (mp *PhysicsEngine) getBonePosition(bones *model.BoneCollection, rigidBody *model.RigidBody) mmath.Vec3 {
	bone := mp.getRigidBodyBone(bones, rigidBody)
	if bone == nil {
		return mmath.NewVec3()
	}
	return bone.Position
}

// configureRigidBody は剛体の物理パラメータを設定します。
func (mp *PhysicsEngine) configureRigidBody(btRigidBody bt.BtRigidBody, modelIndex int, rigidBody *model.RigidBody) {
	btRigidBody.SetDamping(float32(rigidBody.Param.LinearDamping), float32(rigidBody.Param.AngularDamping))
	btRigidBody.SetRestitution(float32(rigidBody.Param.Restitution))
	btRigidBody.SetFriction(float32(rigidBody.Param.Friction))
	btRigidBody.SetUserIndex(modelIndex)
	btRigidBody.SetUserIndex2(rigidBody.Index())
}

// deleteRigidBodies はモデルの全剛体を削除します。
func (mp *PhysicsEngine) deleteRigidBodies(modelIndex int) {
	for _, r := range mp.rigidBodies[modelIndex] {
		if r == nil || r.BtRigidBody == nil {
			continue
		}
		mp.world.RemoveRigidBody(r.BtRigidBody)
		bt.DeleteBtRigidBody(r.BtRigidBody)
	}
	mp.rigidBodies[modelIndex] = nil
}

// updateFlag は剛体の物理フラグを更新します。
func (mp *PhysicsEngine) updateFlag(modelIndex int, rigidBody *model.RigidBody) {
	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody

	if rigidBody.PhysicsType == model.PHYSICS_TYPE_STATIC {
		btRigidBody.SetCollisionFlags(
			btRigidBody.GetCollisionFlags() | int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT),
		)
		btRigidBody.SetActivationState(bt.DISABLE_SIMULATION)
		return
	}

	btRigidBody.SetCollisionFlags(
		btRigidBody.GetCollisionFlags() &
			^int(bt.BtCollisionObjectCF_NO_CONTACT_RESPONSE) &
			^int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT),
	)
	btRigidBody.SetActivationState(bt.DISABLE_DEACTIVATION)
}

// UpdateTransform はボーン行列に基づいて剛体の位置を更新します。
func (mp *PhysicsEngine) UpdateTransform(
	modelIndex int,
	rigidBodyBone *model.Bone,
	boneGlobalMatrix *mmath.Mat4,
	rigidBody *model.RigidBody,
) {
	if rigidBodyBone == nil || boneGlobalMatrix == nil || rigidBody == nil {
		return
	}

	boneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneTransform)

	mat := newMglMat4FromMat4(*boneGlobalMatrix)
	boneTransform.SetFromOpenGLMatrix(&mat[0])

	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody
	btRigidBodyLocalTransform := *mp.rigidBodies[modelIndex][rigidBody.Index()].BtLocalTransform

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)

	currentTransform.Mult(boneTransform, btRigidBodyLocalTransform)
	motionState.SetWorldTransform(currentTransform)
}

// GetRigidBodyBoneMatrix は剛体に基づいてボーン行列を取得します。
func (mp *PhysicsEngine) GetRigidBodyBoneMatrix(
	modelIndex int,
	rigidBody *model.RigidBody,
) *mmath.Mat4 {
	if rigidBody == nil {
		return nil
	}

	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody
	btRigidBodyLocalTransform := *mp.rigidBodies[modelIndex][rigidBody.Index()].BtLocalTransform

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	rigidBodyGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(rigidBodyGlobalTransform)

	motionState.GetWorldTransform(rigidBodyGlobalTransform)

	boneGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneGlobalTransform)

	invRigidBodyLocalTransform := btRigidBodyLocalTransform.Inverse()
	defer bt.DeleteBtTransform(invRigidBodyLocalTransform)

	boneGlobalTransform.Mult(rigidBodyGlobalTransform, invRigidBodyLocalTransform)

	boneGlobalMatrixGL := mgl32.Mat4{}
	boneGlobalTransform.GetOpenGLMatrix(&boneGlobalMatrixGL[0])

	boneGlobalMatrix := newMat4FromMgl(&boneGlobalMatrixGL)
	return &boneGlobalMatrix
}

// UpdateRigidBodiesSelectively は変更が必要な剛体のみを選択的に更新します。
func (mp *PhysicsEngine) UpdateRigidBodiesSelectively(
	modelIndex int,
	pmxModel *model.PmxModel,
	rigidBodyDeltas *delta.RigidBodyDeltas,
) {
	if pmxModel == nil || rigidBodyDeltas == nil {
		return
	}

	rigidBodyDeltas.ForEach(func(index int, rigidBodyDelta *delta.RigidBodyDelta) bool {
		if rigidBodyDelta == nil || rigidBodyDelta.RigidBody == nil {
			return true
		}

		mp.UpdateRigidBodyShapeMass(modelIndex, rigidBodyDelta.RigidBody, rigidBodyDelta)

		return true
	})
}

// UpdateRigidBodyShapeMass はサイズ・質量変更時に剛体の形状を更新します。
func (mp *PhysicsEngine) UpdateRigidBodyShapeMass(
	modelIndex int,
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) {
	if rigidBody == nil || rigidBodyDelta == nil {
		return
	}

	massChanged := !mmath.NearEquals(rigidBodyDelta.Mass, rigidBody.Param.Mass, 1e-10)
	sizeChanged := !rigidBodyDelta.Size.NearEquals(rigidBody.Size, 1e-10)
	if !massChanged && !sizeChanged {
		return
	}

	r := mp.rigidBodies[modelIndex][rigidBody.Index()]
	if r == nil || r.BtRigidBody == nil {
		return
	}

	btRigidBody := r.BtRigidBody
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)
	motionState.GetWorldTransform(currentTransform)

	currentLinearVel := btRigidBody.GetLinearVelocity()
	currentAngularVel := btRigidBody.GetAngularVelocity()
	currentMass := float32(0.0)
	if btRigidBody.GetInvMass() != 0 {
		currentMass = 1.0 / btRigidBody.GetInvMass()
	}

	newMass := currentMass
	if rigidBody.PhysicsType != model.PHYSICS_TYPE_STATIC {
		newMass = float32(mmath.Clamped(rigidBodyDelta.Mass, 0, math.MaxFloat64))
	}

	needShapeUpdate := sizeChanged
	if needShapeUpdate {
		mp.world.RemoveRigidBody(btRigidBody)

		oldShape := btRigidBody.GetCollisionShape()
		if oldShape != nil {
			if collisionShape, ok := oldShape.(bt.BtCollisionShape); ok {
				bt.DeleteBtCollisionShape(collisionShape)
			}
		}

		newShape := mp.createCollisionShape(rigidBody, rigidBodyDelta)
		btRigidBody.SetCollisionShape(newShape)
	}

	newInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
	if newMass > 0 && rigidBody.PhysicsType != model.PHYSICS_TYPE_STATIC {
		currentShape := btRigidBody.GetCollisionShape()
		if currentShape != nil {
			if collisionShape, ok := currentShape.(bt.BtCollisionShape); ok {
				collisionShape.CalculateLocalInertia(newMass, newInertia)
			}
		}
	}

	btRigidBody.SetMassProps(newMass, newInertia)
	btRigidBody.UpdateInertiaTensor()

	motionState.SetWorldTransform(currentTransform)
	btRigidBody.SetLinearVelocity(currentLinearVel)
	btRigidBody.SetAngularVelocity(currentAngularVel)

	if needShapeUpdate {
		mp.world.AddRigidBody(btRigidBody, r.Group, r.Mask)
	}

	mp.updateFlag(modelIndex, rigidBody)
	btRigidBody.Activate(true)
}

// findRigidBodyByCollisionObject は衝突オブジェクトから剛体参照を取得します。
func (mp *PhysicsEngine) findRigidBodyByCollisionObject(
	hitObj bt.BtCollisionObject,
	hasHit bool,
) (modelIndex int, rigidBodyIndex int, ok bool) {
	if hitObj == nil || !hasHit {
		return -1, -1, false
	}
	defer func() {
		if r := recover(); r != nil {
			modelIndex = -1
			rigidBodyIndex = -1
			ok = false
		}
	}()

	modelIndex = hitObj.GetUserIndex()
	rigidBodyIndex = hitObj.GetUserIndex2()

	bodies, exists := mp.rigidBodies[modelIndex]
	if !exists {
		return -1, -1, false
	}
	if rigidBodyIndex < 0 || rigidBodyIndex >= len(bodies) {
		return -1, -1, false
	}
	if bodies[rigidBodyIndex] == nil || bodies[rigidBodyIndex].BtRigidBody == nil {
		return -1, -1, false
	}

	return modelIndex, rigidBodyIndex, true
}
