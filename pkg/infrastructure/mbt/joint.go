//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// jointValue はジョイントの物理エンジン内部表現を格納する構造体です
type jointValue struct {
	pmxJoint *pmx.Joint
	btJoint  bt.BtTypedConstraint
}

// initJoints はモデルのジョイントを初期化します
func (physics *MPhysics) initJoints(modelIndex int, rigidBodies *pmx.RigidBodies, j *pmx.Joints) {
	// ジョイントを順番に剛体と紐付けていく
	physics.joints[modelIndex] = make([]*jointValue, j.Length())
	for v := range j.Iterator() {
		joint := v.Value
		if rigidBodies.Contains(joint.RigidbodyIndexA) && rigidBodies.Contains(joint.RigidbodyIndexB) {
			// ジョイントの位置と向き
			jointTransform := bt.NewBtTransform(newBulletFromRad(joint.Rotation), newBulletFromVec(joint.Position))

			physics.initJoint(modelIndex, joint, jointTransform)
		}
	}
}

// initJointsByBoneDeltas はボーンデルタ情報を使用してジョイントを初期化します
func (physics *MPhysics) initJointsByBoneDeltas(
	modelIndex int, rigidBodies *pmx.RigidBodies, j *pmx.Joints, boneDeltas *delta.BoneDeltas,
) {
	// ジョイントを順番に剛体と紐付けていく
	physics.joints[modelIndex] = make([]*jointValue, j.Length())
	for v := range j.Iterator() {
		joint := v.Value
		if !physics.canCreateJoint(joint, rigidBodies) {
			continue
		}

		// ボーン情報の取得
		bone := physics.findReferenceBone(joint, rigidBodies)
		if bone == nil || !boneDeltas.Contains(bone.Index()) {
			continue
		}

		// ジョイントのグローバル変換を計算
		jointTransform := physics.calculateJointTransform(joint, bone, boneDeltas)

		// ジョイント初期化
		physics.initJoint(modelIndex, joint, jointTransform)
	}
}

// canCreateJoint はジョイントが作成可能かチェックします
func (physics *MPhysics) canCreateJoint(joint *pmx.Joint, rigidBodies *pmx.RigidBodies) bool {
	return rigidBodies.Contains(joint.RigidbodyIndexA) && rigidBodies.Contains(joint.RigidbodyIndexB)
}

// findReferenceBone はジョイントの参照ボーンを検索します
func (physics *MPhysics) findReferenceBone(joint *pmx.Joint, rigidBodies *pmx.RigidBodies) *pmx.Bone {
	// 剛体AまたはBに関連するボーンを検索
	if rb, err := rigidBodies.Get(joint.RigidbodyIndexA); err == nil {
		if rb.Bone != nil {
			return rb.Bone
		} else if rb.JointedBone != nil {
			return rb.JointedBone
		}
	}

	if rb, err := rigidBodies.Get(joint.RigidbodyIndexB); err == nil {
		if rb.Bone != nil {
			return rb.Bone
		} else if rb.JointedBone != nil {
			return rb.JointedBone
		}
	}

	return nil
}

// calculateJointTransform はジョイントの変換行列を計算します
func (physics *MPhysics) calculateJointTransform(
	joint *pmx.Joint, bone *pmx.Bone, boneDeltas *delta.BoneDeltas,
) bt.BtTransform {
	jointTransform := bt.NewBtTransform()
	boneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneTransform)

	// ボーンのグローバル変換行列を設定
	mat := mmath.NewGlMat4(boneDeltas.Get(bone.Index()).FilledGlobalMatrix())
	boneTransform.SetFromOpenGLMatrix(&mat[0])

	// ジョイントのローカル位置と回転を計算
	jointLocalPos := joint.Position.Subed(bone.Position)
	btJointLocalTransform := bt.NewBtTransform(
		newBulletFromRad(joint.Rotation),
		newBulletFromVec(jointLocalPos),
	)

	// ボーングローバル変換とジョイントローカル変換を乗算
	jointTransform.Mult(boneTransform, btJointLocalTransform)

	return jointTransform
}

// initJoint は個別のジョイントを初期化します
func (physics *MPhysics) initJoint(
	modelIndex int, joint *pmx.Joint, jointTransform bt.BtTransform,
) {
	// 関連する剛体が存在するか確認
	if !physics.validateJointRigidBodies(modelIndex, joint) {
		return
	}

	// 剛体の取得
	rigidBodyB := physics.rigidBodies[modelIndex][joint.RigidbodyIndexB].pmxRigidBody
	btRigidBodyA := physics.rigidBodies[modelIndex][joint.RigidbodyIndexA].btRigidBody
	btRigidBodyB := physics.rigidBodies[modelIndex][joint.RigidbodyIndexB].btRigidBody

	// 剛体ローカル座標系におけるジョイント変換を計算
	jointLocalTransforms := physics.calculateJointLocalTransforms(btRigidBodyA, btRigidBodyB, jointTransform)
	jointLocalTransformA := jointLocalTransforms[0]
	jointLocalTransformB := jointLocalTransforms[1]

	// ジョイント拘束の作成
	constraint := physics.createJointConstraint(btRigidBodyA, btRigidBodyB, jointLocalTransformA, jointLocalTransformB)

	// ジョイントのパラメータ設定
	physics.configureJointConstraint(constraint, joint, rigidBodyB)

	// 物理ワールドにジョイントを追加
	physics.world.AddConstraint(constraint)
	physics.joints[modelIndex][joint.Index()] = &jointValue{pmxJoint: joint, btJoint: constraint}
}

// validateJointRigidBodies はジョイントに関連する剛体が有効か検証します
func (physics *MPhysics) validateJointRigidBodies(modelIndex int, joint *pmx.Joint) bool {
	return physics.rigidBodies[modelIndex][joint.RigidbodyIndexB] != nil &&
		physics.rigidBodies[modelIndex][joint.RigidbodyIndexA] != nil &&
		physics.rigidBodies[modelIndex][joint.RigidbodyIndexB].pmxRigidBody != nil &&
		physics.rigidBodies[modelIndex][joint.RigidbodyIndexA].pmxRigidBody != nil
}

// calculateJointLocalTransforms は剛体のローカル座標系におけるジョイント変換を計算します
func (physics *MPhysics) calculateJointLocalTransforms(
	btRigidBodyA bt.BtRigidBody,
	btRigidBodyB bt.BtRigidBody,
	jointTransform bt.BtTransform,
) [2]bt.BtTransform {
	// 剛体Aのローカル座標系におけるジョイント
	worldTransformA := btRigidBodyA.GetWorldTransform().(bt.BtTransform)
	jointLocalTransformA := bt.NewBtTransform()
	jointLocalTransformA.SetIdentity()
	jointLocalTransformA.Mult(worldTransformA.Inverse(), jointTransform)

	// 剛体Bのローカル座標系におけるジョイント
	worldTransformB := btRigidBodyB.GetWorldTransform().(bt.BtTransform)
	jointLocalTransformB := bt.NewBtTransform()
	jointLocalTransformB.SetIdentity()
	jointLocalTransformB.Mult(worldTransformB.Inverse(), jointTransform)

	return [2]bt.BtTransform{jointLocalTransformA, jointLocalTransformB}
}

// createJointConstraint はジョイント拘束を作成します
func (physics *MPhysics) createJointConstraint(
	btRigidBodyA bt.BtRigidBody,
	btRigidBodyB bt.BtRigidBody,
	jointLocalTransformA bt.BtTransform,
	jointLocalTransformB bt.BtTransform,
) bt.BtGeneric6DofSpringConstraint {
	return bt.NewBtGeneric6DofSpringConstraint(
		btRigidBodyA, btRigidBodyB, jointLocalTransformA, jointLocalTransformB, true)
}

// configureJointConstraint はジョイント拘束のパラメータを設定します
func (physics *MPhysics) configureJointConstraint(
	constraint bt.BtGeneric6DofSpringConstraint,
	joint *pmx.Joint,
	rigidBodyB *pmx.RigidBody,
) {
	// 平行移動制限設定
	constraint.SetLinearLowerLimit(bt.NewBtVector3(
		float32(joint.JointParam.TranslationLimitMin.X),
		float32(joint.JointParam.TranslationLimitMin.Y),
		float32(joint.JointParam.TranslationLimitMin.Z)))
	constraint.SetLinearUpperLimit(bt.NewBtVector3(
		float32(joint.JointParam.TranslationLimitMax.X),
		float32(joint.JointParam.TranslationLimitMax.Y),
		float32(joint.JointParam.TranslationLimitMax.Z)))

	// 回転制限設定
	constraint.SetAngularLowerLimit(bt.NewBtVector3(
		float32(joint.JointParam.RotationLimitMin.X),
		float32(joint.JointParam.RotationLimitMin.Y),
		float32(joint.JointParam.RotationLimitMin.Z)))
	constraint.SetAngularUpperLimit(bt.NewBtVector3(
		float32(joint.JointParam.RotationLimitMax.X),
		float32(joint.JointParam.RotationLimitMax.Y),
		float32(joint.JointParam.RotationLimitMax.Z)))

	// バネパラメータ設定
	physics.configureJointSprings(constraint, joint, rigidBodyB)

	// 基本パラメータ設定
	physics.configureBasicJointParams(constraint)
}

// configureJointSprings はジョイントのバネパラメータを設定します
func (physics *MPhysics) configureJointSprings(
	constraint bt.BtGeneric6DofSpringConstraint,
	joint *pmx.Joint,
	rigidBodyB *pmx.RigidBody,
) {
	if rigidBodyB.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
		// 剛体Bがボーン追従剛体の場合は、バネの値を設定しない
		// 平行移動バネ設定
		constraint.EnableSpring(0, true)
		constraint.SetStiffness(0, float32(joint.JointParam.SpringConstantTranslation.X))
		constraint.EnableSpring(1, true)
		constraint.SetStiffness(1, float32(joint.JointParam.SpringConstantTranslation.Y))
		constraint.EnableSpring(2, true)
		constraint.SetStiffness(2, float32(joint.JointParam.SpringConstantTranslation.Z))

		// 回転バネ設定
		constraint.EnableSpring(3, true)
		constraint.SetStiffness(3, float32(joint.JointParam.SpringConstantRotation.X))
		constraint.EnableSpring(4, true)
		constraint.SetStiffness(4, float32(joint.JointParam.SpringConstantRotation.Y))
		constraint.EnableSpring(5, true)
		constraint.SetStiffness(5, float32(joint.JointParam.SpringConstantRotation.Z))
	}
}

// configureBasicJointParams はジョイントの基本パラメータを設定します
func (physics *MPhysics) configureBasicJointParams(constraint bt.BtTypedConstraint) {
	constraint.SetParam(int(bt.BT_CONSTRAINT_ERP), float32(0.5), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_ERP), float32(0.5), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_CFM), float32(0.1), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_CFM), float32(0.1), 0)

	// デバッグ円の表示サイズ
	g6dof, ok := constraint.(bt.BtGeneric6DofSpringConstraint)
	if ok {
		g6dof.SetDbgDrawSize(float32(1.5))
	}
}

// deleteJoints はモデルの全ジョイントを削除します
func (physics *MPhysics) deleteJoints(modelIndex int) {
	for _, j := range physics.joints[modelIndex] {
		if j == nil || j.btJoint == nil {
			continue
		}
		physics.world.RemoveConstraint(j.btJoint)
		bt.DeleteBtTypedConstraint(j.btJoint)
	}
	physics.joints[modelIndex] = nil
}
