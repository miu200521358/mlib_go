//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

func initJointPhysics(
	modelIndex int, modelPhysics *MPhysics, rigidBodyA *pmx.RigidBody, rigidBodyB *pmx.RigidBody, j *pmx.Joint,
) {
	// ジョイントの位置と向き
	jointTransform := bt.NewBtTransform(MRotationBullet(j.Rotation), MVec3Bullet(j.Position))

	btRigidBodyA, _ := modelPhysics.GetRigidBody(modelIndex, rigidBodyA.Index)
	btRigidBodyB, _ := modelPhysics.GetRigidBody(modelIndex, rigidBodyB.Index)
	if btRigidBodyA == nil || btRigidBodyB == nil {
		return
	}

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
		float32(j.JointParam.TranslationLimitMin.GetX()),
		float32(j.JointParam.TranslationLimitMin.GetY()),
		float32(j.JointParam.TranslationLimitMin.GetZ())))
	constraint.SetLinearUpperLimit(bt.NewBtVector3(
		float32(j.JointParam.TranslationLimitMax.GetX()),
		float32(j.JointParam.TranslationLimitMax.GetY()),
		float32(j.JointParam.TranslationLimitMax.GetZ())))
	constraint.SetAngularLowerLimit(bt.NewBtVector3(
		float32(j.JointParam.RotationLimitMin.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetZ())))
	constraint.SetAngularUpperLimit(bt.NewBtVector3(
		float32(j.JointParam.RotationLimitMax.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetZ())))

	if rigidBodyB.PhysicsType != pmx.PHYSICS_TYPE_STATIC {
		// 剛体Bがボーン追従剛体の場合は、バネの値を設定しない
		constraint.EnableSpring(0, true)
		constraint.SetStiffness(0, float32(j.JointParam.SpringConstantTranslation.GetX()))
		constraint.EnableSpring(1, true)
		constraint.SetStiffness(1, float32(j.JointParam.SpringConstantTranslation.GetY()))
		constraint.EnableSpring(2, true)
		constraint.SetStiffness(2, float32(j.JointParam.SpringConstantTranslation.GetZ()))
		constraint.EnableSpring(3, true)
		constraint.SetStiffness(3, float32(j.JointParam.SpringConstantRotation.GetX()))
		constraint.EnableSpring(4, true)
		constraint.SetStiffness(4, float32(j.JointParam.SpringConstantRotation.GetY()))
		constraint.EnableSpring(5, true)
		constraint.SetStiffness(5, float32(j.JointParam.SpringConstantRotation.GetZ()))
	}

	constraint.SetParam(int(bt.BT_CONSTRAINT_ERP), float32(0.5), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_ERP), float32(0.5), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_CFM), float32(0.1), 0)
	constraint.SetParam(int(bt.BT_CONSTRAINT_STOP_CFM), float32(0.1), 0)

	// デバッグ円の表示サイズ
	constraint.SetDbgDrawSize(float32(1.5))

	modelPhysics.AddJoint(modelIndex, j.Index, constraint)
}

func initJointsPhysics(modelIndex int, modelPhysics *MPhysics, rigidBodies *pmx.RigidBodies, j *pmx.Joints) {
	// ジョイントを順番に剛体と紐付けていく
	for _, joint := range j.Data {
		if joint.RigidbodyIndexA >= 0 && rigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && rigidBodies.Contains(joint.RigidbodyIndexB) {
			initJointPhysics(
				modelIndex, modelPhysics, rigidBodies.Get(joint.RigidbodyIndexA),
				rigidBodies.Get(joint.RigidbodyIndexB), joint)
		}
	}
}
