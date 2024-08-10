//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type jointValue struct {
	pmxJoint *pmx.Joint
	btJoint  bt.BtTypedConstraint
}

func (physics *MPhysics) initJoints(modelIndex int, rigidBodies *pmx.RigidBodies, j *pmx.Joints) {
	// ジョイントを順番に剛体と紐付けていく
	physics.joints[modelIndex] = make([]*jointValue, len(j.Data))
	for _, joint := range j.Data {
		if rigidBodies.Contains(joint.RigidbodyIndexA) && rigidBodies.Contains(joint.RigidbodyIndexB) {
			// ジョイントの位置と向き
			jointTransform := bt.NewBtTransform(MRotationBullet(joint.Rotation), MVec3Bullet(joint.Position))

			physics.initJoint(modelIndex, joint, jointTransform)
		}
	}
}

func (physics *MPhysics) initJointsByBoneDeltas(
	modelIndex int, rigidBodies *pmx.RigidBodies, j *pmx.Joints, boneDeltas *delta.BoneDeltas,
) {
	// ジョイントを順番に剛体と紐付けていく
	physics.joints[modelIndex] = make([]*jointValue, len(j.Data))
	for _, joint := range j.Data {
		if rigidBodies.Contains(joint.RigidbodyIndexA) && rigidBodies.Contains(joint.RigidbodyIndexB) {
			// ジョイントの位置と向き
			jointTransform := bt.NewBtTransform()

			var bone *pmx.Bone
			if rigidBodies.Get(joint.RigidbodyIndexA).Bone != nil {
				bone = rigidBodies.Get(joint.RigidbodyIndexA).Bone
			} else if rigidBodies.Get(joint.RigidbodyIndexB).Bone != nil {
				bone = rigidBodies.Get(joint.RigidbodyIndexB).Bone
			} else if rigidBodies.Get(joint.RigidbodyIndexA).JointedBone != nil {
				bone = rigidBodies.Get(joint.RigidbodyIndexA).JointedBone
			} else if rigidBodies.Get(joint.RigidbodyIndexB).JointedBone != nil {
				bone = rigidBodies.Get(joint.RigidbodyIndexB).JointedBone
			}

			if bone == nil {
				continue
			}

			boneTransform := bt.NewBtTransform()
			defer bt.DeleteBtTransform(boneTransform)

			mat := mgl.NewGlMat4(boneDeltas.Get(bone.Index()).FilledGlobalMatrix())
			boneTransform.SetFromOpenGLMatrix(&mat[0])

			jointLocalPos := joint.Position.Subed(bone.Position)
			btJointLocalTransform := bt.NewBtTransform(MRotationBullet(joint.Rotation),
				MVec3Bullet(jointLocalPos))

			jointTransform.Mult(boneTransform, btJointLocalTransform)

			physics.initJoint(modelIndex, joint, jointTransform)
		}
	}
}

func (physics *MPhysics) initJoint(
	modelIndex int, joint *pmx.Joint, jointTransform bt.BtTransform,
) {
	rigidBodyB := physics.rigidBodies[modelIndex][joint.RigidbodyIndexB].pmxRigidBody
	btRigidBodyA := physics.rigidBodies[modelIndex][joint.RigidbodyIndexA].btRigidBody
	btRigidBodyB := physics.rigidBodies[modelIndex][joint.RigidbodyIndexB].btRigidBody

	// 剛体Aの現在の位置と向きを取得
	worldTransformA := btRigidBodyA.GetWorldTransform().(bt.BtTransform)

	// 剛体Aのローカル座標系におけるジョイント
	jointLocalTransformA := bt.NewBtTransform()
	jointLocalTransformA.SetIdentity()
	jointLocalTransformA.Mult(worldTransformA.Inverse(), jointTransform)

	// 剛体Bの現在の位置と向きを取得
	worldTransformB := btRigidBodyB.GetWorldTransform().(bt.BtTransform)

	// 剛体Bのローカル座標系におけるジョイント
	jointLocalTransformB := bt.NewBtTransform()
	jointLocalTransformB.SetIdentity()
	jointLocalTransformB.Mult(worldTransformB.Inverse(), jointTransform)

	// ジョイント係数
	constraint := bt.NewBtGeneric6DofSpringConstraint(
		btRigidBodyA, btRigidBodyB, jointLocalTransformA, jointLocalTransformB, true)
	// 係数は符号を調整する必要がないため、そのまま設定
	constraint.SetLinearLowerLimit(bt.NewBtVector3(
		float32(joint.JointParam.TranslationLimitMin.X),
		float32(joint.JointParam.TranslationLimitMin.Y),
		float32(joint.JointParam.TranslationLimitMin.Z)))
	constraint.SetLinearUpperLimit(bt.NewBtVector3(
		float32(joint.JointParam.TranslationLimitMax.X),
		float32(joint.JointParam.TranslationLimitMax.Y),
		float32(joint.JointParam.TranslationLimitMax.Z)))
	constraint.SetAngularLowerLimit(bt.NewBtVector3(
		float32(joint.JointParam.RotationLimitMin.Radians().X),
		float32(joint.JointParam.RotationLimitMin.Radians().Y),
		float32(joint.JointParam.RotationLimitMin.Radians().Z)))
	constraint.SetAngularUpperLimit(bt.NewBtVector3(
		float32(joint.JointParam.RotationLimitMax.Radians().X),
		float32(joint.JointParam.RotationLimitMax.Radians().Y),
		float32(joint.JointParam.RotationLimitMax.Radians().Z)))

	if rigidBodyB.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
		// 剛体Bがボーン追従剛体の場合は、バネの値を設定しない
		constraint.EnableSpring(0, true)
		constraint.SetStiffness(0, float32(joint.JointParam.SpringConstantTranslation.X))
		constraint.EnableSpring(1, true)
		constraint.SetStiffness(1, float32(joint.JointParam.SpringConstantTranslation.Y))
		constraint.EnableSpring(2, true)
		constraint.SetStiffness(2, float32(joint.JointParam.SpringConstantTranslation.Z))
		constraint.EnableSpring(3, true)
		constraint.SetStiffness(3, float32(joint.JointParam.SpringConstantRotation.X))
		constraint.EnableSpring(4, true)
		constraint.SetStiffness(4, float32(joint.JointParam.SpringConstantRotation.Y))
		constraint.EnableSpring(5, true)
		constraint.SetStiffness(5, float32(joint.JointParam.SpringConstantRotation.Z))
	}

	constraint.SetParam(int(bt.BT_CONSTRAINT_ERP), float32(0.5), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_ERP), float32(0.5), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_CFM), float32(0.1), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_CFM), float32(0.1), 0)

	// デバッグ円の表示サイズ
	constraint.SetDbgDrawSize(float32(1.5))

	physics.world.AddConstraint(constraint)
	physics.joints[modelIndex][joint.Index()] = &jointValue{pmxJoint: joint, btJoint: constraint}
}

func (physics *MPhysics) deleteJoints(modelIndex int) {
	for _, j := range physics.joints[modelIndex] {
		physics.world.RemoveConstraint(j.btJoint)
		bt.DeleteBtTypedConstraint(j.btJoint)
	}
	physics.joints[modelIndex] = nil
}
