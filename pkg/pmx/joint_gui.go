//go:build windows
// +build windows

package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
)

func (j *Joint) initPhysics(modelPhysics *mphysics.MPhysics, rigidBodyA *RigidBody, rigidBodyB *RigidBody) {
	// ジョイントの位置と向き
	jointTransform := mbt.NewBtTransform(j.Rotation.Bullet(), j.Position.Bullet())

	btRigidBodyA, _ := modelPhysics.GetRigidBody(rigidBodyA.Index)
	btRigidBodyB, _ := modelPhysics.GetRigidBody(rigidBodyB.Index)
	if btRigidBodyA == nil || btRigidBodyB == nil {
		return
	}

	// 剛体Aの現在の位置と向きを取得
	worldTransformA := btRigidBodyA.GetWorldTransform().(mbt.BtTransform)

	// 剛体Aのローカル座標系におけるジョイント
	jointLocalTransformA := mbt.NewBtTransform()
	jointLocalTransformA.SetIdentity()
	jointLocalTransformA.Mult(worldTransformA.Inverse(), jointTransform)

	// 剛体Bの現在の位置と向きを取得
	worldTransformB := btRigidBodyB.GetWorldTransform().(mbt.BtTransform)

	// 剛体Bのローカル座標系におけるジョイント
	jointLocalTransformB := mbt.NewBtTransform()
	jointLocalTransformB.SetIdentity()
	jointLocalTransformB.Mult(worldTransformB.Inverse(), jointTransform)

	// ジョイント係数
	constraint := mbt.NewBtGeneric6DofSpringConstraint(
		btRigidBodyA, btRigidBodyB, jointLocalTransformA, jointLocalTransformB, true)
	// 係数は符号を調整する必要がないため、そのまま設定
	constraint.SetLinearLowerLimit(mbt.NewBtVector3(
		float32(j.JointParam.TranslationLimitMin.GetX()),
		float32(j.JointParam.TranslationLimitMin.GetY()),
		float32(j.JointParam.TranslationLimitMin.GetZ())))
	constraint.SetLinearUpperLimit(mbt.NewBtVector3(
		float32(j.JointParam.TranslationLimitMax.GetX()),
		float32(j.JointParam.TranslationLimitMax.GetY()),
		float32(j.JointParam.TranslationLimitMax.GetZ())))
	constraint.SetAngularLowerLimit(mbt.NewBtVector3(
		float32(j.JointParam.RotationLimitMin.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetZ())))
	constraint.SetAngularUpperLimit(mbt.NewBtVector3(
		float32(j.JointParam.RotationLimitMax.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetZ())))
	constraint.EnableSpring(0, true)
	constraint.SetStiffness(0, float32(j.JointParam.SpringConstantTranslation.GetX()))
	constraint.EnableSpring(1, true)
	constraint.SetStiffness(1, float32(j.JointParam.SpringConstantTranslation.GetY()))
	constraint.EnableSpring(2, true)
	constraint.SetStiffness(2, float32(j.JointParam.SpringConstantTranslation.GetZ()))
	constraint.EnableSpring(3, true)
	constraint.SetStiffness(3, float32(j.JointParam.SpringConstantRotation.GetRadians().GetX()))
	constraint.EnableSpring(4, true)
	constraint.SetStiffness(4, float32(j.JointParam.SpringConstantRotation.GetRadians().GetY()))
	constraint.EnableSpring(5, true)
	constraint.SetStiffness(5, float32(j.JointParam.SpringConstantRotation.GetRadians().GetZ()))

	constraint.SetParam(int(mbt.BT_CONSTRAINT_ERP), float32(0.5), 0)
	constraint.SetParam(int(mbt.BT_CONSTRAINT_STOP_ERP), float32(0.5), 0)
	constraint.SetParam(int(mbt.BT_CONSTRAINT_CFM), float32(0.1), 0)
	constraint.SetParam(int(mbt.BT_CONSTRAINT_STOP_CFM), float32(0.1), 0)

	// 円の表示サイズ
	constraint.SetDbgDrawSize(float32(1.5))

	modelPhysics.AddJoint(constraint)
}

func (j *Joints) initPhysics(modelPhysics *mphysics.MPhysics, rigidBodies *RigidBodies) {
	// ジョイントを順番に剛体と紐付けていく
	for _, joint := range j.GetSortedData() {
		if joint.RigidbodyIndexA >= 0 && rigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && rigidBodies.Contains(joint.RigidbodyIndexB) {
			joint.initPhysics(
				modelPhysics, rigidBodies.Get(joint.RigidbodyIndexA),
				rigidBodies.Get(joint.RigidbodyIndexB))
		}
	}
}
