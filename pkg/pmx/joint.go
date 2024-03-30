package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
)

type JointParam struct {
	TranslationLimitMin       *mmath.MVec3     // 移動制限-下限(x,y,z)
	TranslationLimitMax       *mmath.MVec3     // 移動制限-上限(x,y,z)
	RotationLimitMin          *mmath.MRotation // 回転制限-下限
	RotationLimitMax          *mmath.MRotation // 回転制限-上限
	SpringConstantTranslation *mmath.MVec3     // バネ定数-移動(x,y,z)
	SpringConstantRotation    *mmath.MRotation // バネ定数-回転(x,y,z)
}

func NewJointParam() *JointParam {
	return &JointParam{
		TranslationLimitMin:       mmath.NewMVec3(),
		TranslationLimitMax:       mmath.NewMVec3(),
		RotationLimitMin:          mmath.NewRotationModel(),
		RotationLimitMax:          mmath.NewRotationModel(),
		SpringConstantTranslation: mmath.NewMVec3(),
		SpringConstantRotation:    mmath.NewRotationModel(),
	}
}

type Joint struct {
	*mcore.IndexNameModel
	JointType       byte             // Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
	RigidbodyIndexA int              // 関連剛体AのIndex
	RigidbodyIndexB int              // 関連剛体BのIndex
	Position        *mmath.MVec3     // 位置(x,y,z)
	Rotation        *mmath.MRotation // 回転
	JointParam      *JointParam      // ジョイントパラメーター
	IsSystem        bool
	Constraint      mbt.BtGeneric6DofSpringConstraint // Bulletのジョイント
}

func NewJoint() *Joint {
	return &Joint{
		IndexNameModel: &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		// Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		JointType:       0,
		RigidbodyIndexA: -1,
		RigidbodyIndexB: -1,
		Position:        mmath.NewMVec3(),
		Rotation:        mmath.NewRotationModel(),
		JointParam:      NewJointParam(),
		IsSystem:        false,
	}
}

func (j *Joint) Copy() mcore.IIndexNameModel {
	copied := NewJoint()
	copier.CopyWithOption(copied, j, copier.Option{DeepCopy: true})
	return copied
}

func NewJointByName(name string) *Joint {
	j := NewJoint()
	j.Name = name
	return j
}

func (j *Joint) initPhysics(modelPhysics *mphysics.MPhysics, rigidBodyA *RigidBody, rigidBodyB *RigidBody) {
	// ジョイントの位置と向き
	jointTransform := mbt.NewBtTransform(j.Rotation.GetQuaternion().Bullet(), j.Position.Bullet())

	// 剛体Aの現在の位置と向きを取得
	worldTransformA := rigidBodyA.BtRigidBody.GetWorldTransform().(mbt.BtTransform)

	// 剛体Aのローカル座標系におけるジョイント
	jointLocalTransformA := mbt.NewBtTransform()
	jointLocalTransformA.SetIdentity()
	jointLocalTransformA.Mult(worldTransformA.Inverse(), jointTransform)

	// 剛体Bの現在の位置と向きを取得
	worldTransformB := rigidBodyB.BtRigidBody.GetWorldTransform().(mbt.BtTransform)

	// 剛体Bのローカル座標系におけるジョイント
	jointLocalTransformB := mbt.NewBtTransform()
	jointLocalTransformB.SetIdentity()
	jointLocalTransformB.Mult(worldTransformB.Inverse(), jointTransform)

	// ジョイント係数
	j.Constraint = mbt.NewBtGeneric6DofSpringConstraint(
		rigidBodyA.BtRigidBody, rigidBodyB.BtRigidBody, jointLocalTransformA, jointLocalTransformB, true)
	// 係数は符号を調整する必要がないため、そのまま設定
	j.Constraint.SetLinearLowerLimit(mbt.NewBtVector3(
		float32(j.JointParam.TranslationLimitMin.GetX()),
		float32(j.JointParam.TranslationLimitMin.GetY()),
		float32(j.JointParam.TranslationLimitMin.GetZ())))
	j.Constraint.SetLinearUpperLimit(mbt.NewBtVector3(
		float32(j.JointParam.TranslationLimitMax.GetX()),
		float32(j.JointParam.TranslationLimitMax.GetY()),
		float32(j.JointParam.TranslationLimitMax.GetZ())))
	j.Constraint.SetAngularLowerLimit(mbt.NewBtVector3(
		float32(j.JointParam.RotationLimitMin.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMin.GetRadians().GetZ())))
	j.Constraint.SetAngularUpperLimit(mbt.NewBtVector3(
		float32(j.JointParam.RotationLimitMax.GetRadians().GetX()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetY()),
		float32(j.JointParam.RotationLimitMax.GetRadians().GetZ())))
	j.Constraint.EnableSpring(0, true)
	j.Constraint.SetStiffness(0, float32(j.JointParam.SpringConstantTranslation.GetX()))
	j.Constraint.EnableSpring(1, true)
	j.Constraint.SetStiffness(1, float32(j.JointParam.SpringConstantTranslation.GetY()))
	j.Constraint.EnableSpring(2, true)
	j.Constraint.SetStiffness(2, float32(j.JointParam.SpringConstantTranslation.GetZ()))
	j.Constraint.EnableSpring(3, true)
	j.Constraint.SetStiffness(3, float32(j.JointParam.SpringConstantRotation.GetRadians().GetX()))
	j.Constraint.EnableSpring(4, true)
	j.Constraint.SetStiffness(4, float32(j.JointParam.SpringConstantRotation.GetRadians().GetY()))
	j.Constraint.EnableSpring(5, true)
	j.Constraint.SetStiffness(5, float32(j.JointParam.SpringConstantRotation.GetRadians().GetZ()))

	// j.Constraint.SetParam(int(mbt.BT_CONSTRAINT_ERP), float32(0.8), 0)
	// j.Constraint.SetParam(int(mbt.BT_CONSTRAINT_STOP_ERP), float32(0.8), 0)
	// j.Constraint.SetParam(int(mbt.BT_CONSTRAINT_CFM), float32(0.2), 0)
	// j.Constraint.SetParam(int(mbt.BT_CONSTRAINT_STOP_CFM), float32(0.2), 0)

	// 円の表示サイズ
	j.Constraint.SetDbgDrawSize(float32(1.5))

	modelPhysics.AddJoint(j.Constraint)
}

func (j *Joint) deletePhysics() {
	j.Constraint = nil
}

// ジョイントリスト
type Joints struct {
	*mcore.IndexNameModels[*Joint]
}

func NewJoints() *Joints {
	return &Joints{
		IndexNameModels: mcore.NewIndexNameModels[*Joint](),
	}
}

func (j *Joints) initPhysics(modelPhysics *mphysics.MPhysics, rigidBodies *RigidBodies) {
	// ジョイントを順番に剛体と紐付けていく
	for _, joint := range j.GetSortedData() {
		if joint.RigidbodyIndexA >= 0 && rigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && rigidBodies.Contains(joint.RigidbodyIndexB) {
			joint.initPhysics(
				modelPhysics, rigidBodies.GetItem(joint.RigidbodyIndexA),
				rigidBodies.GetItem(joint.RigidbodyIndexB))
		}
	}
}

func (j *Joints) deletePhysics(modelPhysics *mphysics.MPhysics) {
	for _, joint := range j.Data {
		modelPhysics.DeleteJoint(joint.Constraint)
		joint.deletePhysics()
	}
}
