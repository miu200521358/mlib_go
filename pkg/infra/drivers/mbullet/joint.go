//go:build windows
// +build windows

// 指示: miu200521358
package mbullet

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet/bt"
)

// jointValue はジョイントの物理エンジン内部表現を保持します。
type jointValue struct {
	joint   *model.Joint
	btJoint bt.BtTypedConstraint
}

// SetJointConstraintConfig は全モデル共通の拘束設定を更新します。
func (mp *PhysicsEngine) SetJointConstraintConfig(config JointConstraintConfig) {
	mp.jointConfig = config
}

// SetModelJointConstraintConfig はモデルごとの拘束設定を更新します。
// 連結剛体衝突設定は拘束生成時に適用されるため、既存拘束へは次回再生成時に反映されます。
func (mp *PhysicsEngine) SetModelJointConstraintConfig(modelIndex int, config *JointConstraintConfig) {
	if config == nil {
		delete(mp.modelJoints, modelIndex)
		return
	}
	mp.modelJoints[modelIndex] = *config
}

// SetModelDisableCollisionsBetweenLinkedBody は連結剛体の衝突無効フラグをモデル単位で設定します。
// 衝突無効フラグは拘束生成時に適用されるため、既存拘束へは次回再生成時に反映されます。
func (mp *PhysicsEngine) SetModelDisableCollisionsBetweenLinkedBody(modelIndex int, disable bool) {
	config := mp.resolveJointConstraintConfig(modelIndex)
	config.DisableCollisionsBetweenLinkedBody = disable
	mp.modelJoints[modelIndex] = config
}

// resolveJointConstraintConfig はモデル単位の拘束設定を取得します。
func (mp *PhysicsEngine) resolveJointConstraintConfig(modelIndex int) JointConstraintConfig {
	if config, ok := mp.modelJoints[modelIndex]; ok {
		return config
	}
	return mp.jointConfig
}

// initJoints はモデルのジョイントを初期化します。
func (mp *PhysicsEngine) initJoints(modelIndex int, pmxModel *model.PmxModel) {
	if pmxModel == nil || pmxModel.Joints == nil || pmxModel.RigidBodies == nil {
		return
	}

	joints := pmxModel.Joints.Values()
	mp.joints[modelIndex] = make([]*jointValue, len(joints))
	for _, joint := range joints {
		if joint == nil {
			continue
		}
		if !mp.canCreateJoint(joint, pmxModel.RigidBodies) {
			continue
		}
		jointTransform := bt.NewBtTransform(newBulletFromRad(joint.Param.Rotation), newBulletFromVec(joint.Param.Position))
		mp.initJoint(modelIndex, joint, jointTransform, nil)
	}
}

// initJointsByBoneDeltas はボーンデルタ情報を使用してジョイントを初期化します。
func (mp *PhysicsEngine) initJointsByBoneDeltas(
	modelIndex int,
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	jointDeltas *delta.JointDeltas,
) {
	if pmxModel == nil || pmxModel.Joints == nil || pmxModel.RigidBodies == nil || boneDeltas == nil {
		return
	}

	joints := pmxModel.Joints.Values()
	mp.joints[modelIndex] = make([]*jointValue, len(joints))
	for _, joint := range joints {
		if joint == nil {
			continue
		}
		if !mp.canCreateJoint(joint, pmxModel.RigidBodies) {
			continue
		}

		bone := mp.findReferenceBone(joint, pmxModel.Bones, pmxModel.RigidBodies)
		if bone == nil || !boneDeltas.Contains(bone.Index()) {
			continue
		}

		jointTransform := mp.calculateJointTransform(joint, bone, boneDeltas)

		var jointDelta *delta.JointDelta
		if jointDeltas != nil {
			jointDelta = jointDeltas.Get(joint.Index())
		}

		mp.initJoint(modelIndex, joint, jointTransform, jointDelta)
	}
}

// canCreateJoint はジョイントが作成可能か判定します。
func (mp *PhysicsEngine) canCreateJoint(
	joint *model.Joint,
	rigidBodies *collection.NamedCollection[*model.RigidBody],
) bool {
	if joint == nil || rigidBodies == nil {
		return false
	}
	return rigidBodies.Contains(joint.RigidBodyIndexA) && rigidBodies.Contains(joint.RigidBodyIndexB)
}

// findReferenceBone はジョイントの参照ボーンを取得します。
func (mp *PhysicsEngine) findReferenceBone(
	joint *model.Joint,
	bones *model.BoneCollection,
	rigidBodies *collection.NamedCollection[*model.RigidBody],
) *model.Bone {
	if joint == nil || bones == nil || rigidBodies == nil {
		return nil
	}

	if rb, err := rigidBodies.Get(joint.RigidBodyIndexA); err == nil {
		bone := mp.getRigidBodyBone(bones, rb)
		if bone != nil {
			return bone
		}
	}
	if rb, err := rigidBodies.Get(joint.RigidBodyIndexB); err == nil {
		bone := mp.getRigidBodyBone(bones, rb)
		if bone != nil {
			return bone
		}
	}
	return nil
}

// calculateJointTransform はジョイントの変換行列を計算します。
func (mp *PhysicsEngine) calculateJointTransform(
	joint *model.Joint,
	bone *model.Bone,
	boneDeltas *delta.BoneDeltas,
) bt.BtTransform {
	jointTransform := bt.NewBtTransform()
	boneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneTransform)

	mat := newMglMat4FromMat4(boneDeltas.Get(bone.Index()).FilledGlobalMatrix())
	boneTransform.SetFromOpenGLMatrix(&mat[0])

	jointLocalPos := joint.Param.Position.Subed(bone.Position)
	btJointLocalTransform := bt.NewBtTransform(
		newBulletFromRad(joint.Param.Rotation),
		newBulletFromVec(jointLocalPos),
	)

	jointTransform.Mult(boneTransform, btJointLocalTransform)
	return jointTransform
}

// initJoint は個別のジョイントを初期化します。
func (mp *PhysicsEngine) initJoint(
	modelIndex int,
	joint *model.Joint,
	jointTransform bt.BtTransform,
	jointDelta *delta.JointDelta,
) {
	if !mp.validateJointRigidBodies(modelIndex, joint) {
		return
	}

	rigidBodyB := mp.rigidBodies[modelIndex][joint.RigidBodyIndexB].RigidBody
	btRigidBodyA := mp.rigidBodies[modelIndex][joint.RigidBodyIndexA].BtRigidBody
	btRigidBodyB := mp.rigidBodies[modelIndex][joint.RigidBodyIndexB].BtRigidBody

	jointLocalTransformA, jointLocalTransformB := mp.calculateJointLocalTransforms(btRigidBodyA, btRigidBodyB, jointTransform)
	constraint := mp.createJointConstraint(btRigidBodyA, btRigidBodyB, jointLocalTransformA, jointLocalTransformB)
	constraintConfig := mp.resolveJointConstraintConfig(modelIndex)

	mp.configureJointConstraint(constraint, joint, rigidBodyB, jointDelta, constraintConfig)

	// 連結剛体同士の自己衝突有無はモデル単位設定で切り替える。
	mp.world.AddConstraint(constraint, constraintConfig.DisableCollisionsBetweenLinkedBody)
	mp.joints[modelIndex][joint.Index()] = &jointValue{joint: joint, btJoint: constraint}
}

// validateJointRigidBodies はジョイントに関連する剛体が有効か検証します。
func (mp *PhysicsEngine) validateJointRigidBodies(modelIndex int, joint *model.Joint) bool {
	if joint == nil {
		return false
	}
	return mp.rigidBodies[modelIndex][joint.RigidBodyIndexB] != nil &&
		mp.rigidBodies[modelIndex][joint.RigidBodyIndexA] != nil &&
		mp.rigidBodies[modelIndex][joint.RigidBodyIndexB].RigidBody != nil &&
		mp.rigidBodies[modelIndex][joint.RigidBodyIndexA].RigidBody != nil
}

// calculateJointLocalTransforms は剛体のローカル座標系におけるジョイント変換を計算します。
func (mp *PhysicsEngine) calculateJointLocalTransforms(
	btRigidBodyA bt.BtRigidBody,
	btRigidBodyB bt.BtRigidBody,
	jointTransform bt.BtTransform,
) (bt.BtTransform, bt.BtTransform) {
	worldTransformA := btRigidBodyA.GetWorldTransform().(bt.BtTransform)
	jointLocalTransformA := bt.NewBtTransform()
	jointLocalTransformA.SetIdentity()
	jointLocalTransformA.Mult(worldTransformA.Inverse(), jointTransform)

	worldTransformB := btRigidBodyB.GetWorldTransform().(bt.BtTransform)
	jointLocalTransformB := bt.NewBtTransform()
	jointLocalTransformB.SetIdentity()
	jointLocalTransformB.Mult(worldTransformB.Inverse(), jointTransform)

	return jointLocalTransformA, jointLocalTransformB
}

// createJointConstraint はジョイント拘束を作成します。
func (mp *PhysicsEngine) createJointConstraint(
	btRigidBodyA bt.BtRigidBody,
	btRigidBodyB bt.BtRigidBody,
	jointLocalTransformA bt.BtTransform,
	jointLocalTransformB bt.BtTransform,
) bt.BtGeneric6DofSpringConstraint {
	return bt.NewBtGeneric6DofSpringConstraint(
		btRigidBodyA, btRigidBodyB, jointLocalTransformA, jointLocalTransformB, true,
	)
}

// configureJointConstraint はジョイント拘束のパラメータを設定します。
func (mp *PhysicsEngine) configureJointConstraint(
	constraint bt.BtGeneric6DofSpringConstraint,
	joint *model.Joint,
	rigidBodyB *model.RigidBody,
	jointDelta *delta.JointDelta,
	constraintConfig JointConstraintConfig,
) {
	translationLimitMin := joint.Param.TranslationLimitMin
	translationLimitMax := joint.Param.TranslationLimitMax
	rotationLimitMin := joint.Param.RotationLimitMin
	rotationLimitMax := joint.Param.RotationLimitMax
	if jointDelta != nil {
		translationLimitMin = jointDelta.TranslationLimitMin
		translationLimitMax = jointDelta.TranslationLimitMax
		rotationLimitMin = jointDelta.RotationLimitMin
		rotationLimitMax = jointDelta.RotationLimitMax
	}

	constraint.SetLinearLowerLimit(bt.NewBtVector3(
		float32(translationLimitMin.X),
		float32(translationLimitMin.Y),
		float32(translationLimitMin.Z),
	))
	constraint.SetLinearUpperLimit(bt.NewBtVector3(
		float32(translationLimitMax.X),
		float32(translationLimitMax.Y),
		float32(translationLimitMax.Z),
	))

	constraint.SetAngularLowerLimit(bt.NewBtVector3(
		float32(rotationLimitMin.X),
		float32(rotationLimitMin.Y),
		float32(rotationLimitMin.Z),
	))
	constraint.SetAngularUpperLimit(bt.NewBtVector3(
		float32(rotationLimitMax.X),
		float32(rotationLimitMax.Y),
		float32(rotationLimitMax.Z),
	))

	mp.configureJointSprings(constraint, joint, rigidBodyB, jointDelta)
	mp.configureBasicJointParams(constraint, constraintConfig)
}

// configureJointSprings はジョイントのバネパラメータを設定します。
func (mp *PhysicsEngine) configureJointSprings(
	constraint bt.BtGeneric6DofSpringConstraint,
	joint *model.Joint,
	rigidBodyB *model.RigidBody,
	jointDelta *delta.JointDelta,
) {
	springConstantTranslation := joint.Param.SpringConstantTranslation
	springConstantRotation := joint.Param.SpringConstantRotation
	if jointDelta != nil {
		springConstantTranslation = jointDelta.SpringConstantTranslation
		springConstantRotation = jointDelta.SpringConstantRotation
	}

	if rigidBodyB.PhysicsType != model.PHYSICS_TYPE_STATIC {
		constraint.EnableSpring(0, true)
		constraint.SetStiffness(0, float32(springConstantTranslation.X))
		constraint.EnableSpring(1, true)
		constraint.SetStiffness(1, float32(springConstantTranslation.Y))
		constraint.EnableSpring(2, true)
		constraint.SetStiffness(2, float32(springConstantTranslation.Z))

		constraint.EnableSpring(3, true)
		constraint.SetStiffness(3, float32(springConstantRotation.X))
		constraint.EnableSpring(4, true)
		constraint.SetStiffness(4, float32(springConstantRotation.Y))
		constraint.EnableSpring(5, true)
		constraint.SetStiffness(5, float32(springConstantRotation.Z))
	}
}

// configureBasicJointParams はジョイントの基本パラメータを設定します。
func (mp *PhysicsEngine) configureBasicJointParams(constraint bt.BtTypedConstraint, constraintConfig JointConstraintConfig) {
	for axis := 0; axis < 6; axis++ {
		constraint.SetParam(int(bt.BT_CONSTRAINT_ERP), constraintConfig.ERP, axis)
		constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_ERP), constraintConfig.StopERP, axis)
		constraint.SetParam(int(bt.BT_CONSTRAINT_CFM), constraintConfig.CFM, axis)
		constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_CFM), constraintConfig.StopCFM, axis)
	}

	g6dof, ok := constraint.(bt.BtGeneric6DofSpringConstraint)
	if ok {
		g6dof.SetDbgDrawSize(float32(1.5))
	}
}

// UpdateJointParameters はジョイントパラメータを更新します。
func (mp *PhysicsEngine) UpdateJointParameters(
	modelIndex int,
	joint *model.Joint,
	jointDelta *delta.JointDelta,
) {
	if jointDelta == nil || joint == nil {
		return
	}

	j := mp.joints[modelIndex][joint.Index()]
	if j == nil || j.btJoint == nil {
		return
	}

	constraint, ok := j.btJoint.(bt.BtGeneric6DofSpringConstraint)
	if !ok {
		return
	}

	rigidBodyB := mp.rigidBodies[modelIndex][joint.RigidBodyIndexB].RigidBody
	if rigidBodyB == nil {
		return
	}

	constraint.SetLinearLowerLimit(bt.NewBtVector3(
		float32(jointDelta.TranslationLimitMin.X),
		float32(jointDelta.TranslationLimitMin.Y),
		float32(jointDelta.TranslationLimitMin.Z),
	))
	constraint.SetLinearUpperLimit(bt.NewBtVector3(
		float32(jointDelta.TranslationLimitMax.X),
		float32(jointDelta.TranslationLimitMax.Y),
		float32(jointDelta.TranslationLimitMax.Z),
	))
	constraint.SetAngularLowerLimit(bt.NewBtVector3(
		float32(jointDelta.RotationLimitMin.X),
		float32(jointDelta.RotationLimitMin.Y),
		float32(jointDelta.RotationLimitMin.Z),
	))
	constraint.SetAngularUpperLimit(bt.NewBtVector3(
		float32(jointDelta.RotationLimitMax.X),
		float32(jointDelta.RotationLimitMax.Y),
		float32(jointDelta.RotationLimitMax.Z),
	))

	if rigidBodyB.PhysicsType != model.PHYSICS_TYPE_STATIC {
		constraint.EnableSpring(0, true)
		constraint.SetStiffness(0, float32(jointDelta.SpringConstantTranslation.X))
		constraint.EnableSpring(1, true)
		constraint.SetStiffness(1, float32(jointDelta.SpringConstantTranslation.Y))
		constraint.EnableSpring(2, true)
		constraint.SetStiffness(2, float32(jointDelta.SpringConstantTranslation.Z))

		constraint.EnableSpring(3, true)
		constraint.SetStiffness(3, float32(jointDelta.SpringConstantRotation.X))
		constraint.EnableSpring(4, true)
		constraint.SetStiffness(4, float32(jointDelta.SpringConstantRotation.Y))
		constraint.EnableSpring(5, true)
		constraint.SetStiffness(5, float32(jointDelta.SpringConstantRotation.Z))
	}
	mp.configureBasicJointParams(constraint, mp.resolveJointConstraintConfig(modelIndex))
}

// UpdateJointsSelectively は変更が必要なジョイントのみを選択的に更新します。
func (mp *PhysicsEngine) UpdateJointsSelectively(
	modelIndex int,
	pmxModel *model.PmxModel,
	jointDeltas *delta.JointDeltas,
) {
	if pmxModel == nil || jointDeltas == nil {
		return
	}

	jointDeltas.ForEach(func(index int, jointDelta *delta.JointDelta) bool {
		if jointDelta == nil || jointDelta.Joint == nil {
			return true
		}

		mp.UpdateJointParameters(modelIndex, jointDelta.Joint, jointDelta)
		return true
	})
}

// deleteJoints はモデルの全ジョイントを削除します。
func (mp *PhysicsEngine) deleteJoints(modelIndex int) {
	for _, j := range mp.joints[modelIndex] {
		if j == nil || j.btJoint == nil {
			continue
		}
		mp.world.RemoveConstraint(j.btJoint)
		bt.DeleteBtTypedConstraint(j.btJoint)
	}
	mp.joints[modelIndex] = nil
}
