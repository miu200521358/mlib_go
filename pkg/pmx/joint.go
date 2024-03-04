package pmx

import (
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
		RotationLimitMin:          mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		RotationLimitMax:          mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		SpringConstantTranslation: mmath.NewMVec3(),
		SpringConstantRotation:    mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
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
		Rotation:        mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		JointParam:      NewJointParam(),
		IsSystem:        false,
	}
}

func NewJointByName(name string) *Joint {
	j := NewJoint()
	j.Name = name
	return j
}

func (j *Joint) InitPhysics(modelPhysics *mphysics.MPhysics, rigidBodyA *RigidBody, rigidBodyB *RigidBody) {
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
	j.Constraint.SetLinearLowerLimit(j.JointParam.TranslationLimitMin.Bullet())
	j.Constraint.SetLinearUpperLimit(j.JointParam.TranslationLimitMax.Bullet())
	j.Constraint.SetAngularLowerLimit(j.JointParam.RotationLimitMin.GetRadians().Bullet())
	j.Constraint.SetAngularUpperLimit(j.JointParam.RotationLimitMax.GetRadians().Bullet())
	if j.JointParam.SpringConstantTranslation.GetX() != 0 {
		j.Constraint.EnableSpring(0, true)
		j.Constraint.SetStiffness(0, float32(j.JointParam.SpringConstantTranslation.GetX()))
	}
	if j.JointParam.SpringConstantTranslation.GetY() != 0 {
		j.Constraint.EnableSpring(1, true)
		j.Constraint.SetStiffness(1, float32(j.JointParam.SpringConstantTranslation.GetY()))
	}
	if j.JointParam.SpringConstantTranslation.GetZ() != 0 {
		j.Constraint.EnableSpring(2, true)
		j.Constraint.SetStiffness(2, float32(j.JointParam.SpringConstantTranslation.GetZ()))
	}
	if j.JointParam.SpringConstantRotation.GetRadians().GetX() != 0 {
		j.Constraint.EnableSpring(3, true)
		j.Constraint.SetStiffness(3, float32(j.JointParam.SpringConstantRotation.GetRadians().GetX()))
	}
	if j.JointParam.SpringConstantRotation.GetRadians().GetY() != 0 {
		j.Constraint.EnableSpring(4, true)
		j.Constraint.SetStiffness(4, float32(j.JointParam.SpringConstantRotation.GetRadians().GetY()))
	}
	if j.JointParam.SpringConstantRotation.GetRadians().GetZ() != 0 {
		j.Constraint.EnableSpring(5, true)
		j.Constraint.SetStiffness(5, float32(j.JointParam.SpringConstantRotation.GetRadians().GetZ()))
	}

	modelPhysics.AddJoint(j.Constraint)
}

// ジョイントリスト
type Joints struct {
	*mcore.IndexNameModelCorrection[*Joint]
}

func NewJoints() *Joints {
	return &Joints{
		IndexNameModelCorrection: mcore.NewIndexNameModelCorrection[*Joint](),
	}
}
