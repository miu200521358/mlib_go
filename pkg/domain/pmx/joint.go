package pmx

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type JointParam struct {
	TranslationLimitMin       *mmath.MVec3 // 移動制限-下限(x,y,z)
	TranslationLimitMax       *mmath.MVec3 // 移動制限-上限(x,y,z)
	RotationLimitMin          *mmath.MVec3 // 回転制限-下限
	RotationLimitMax          *mmath.MVec3 // 回転制限-上限
	SpringConstantTranslation *mmath.MVec3 // バネ定数-移動(x,y,z)
	SpringConstantRotation    *mmath.MVec3 // バネ定数-回転(x,y,z)
}

func NewJointParam() *JointParam {
	return &JointParam{
		TranslationLimitMin:       mmath.NewMVec3(),
		TranslationLimitMax:       mmath.NewMVec3(),
		RotationLimitMin:          mmath.NewMVec3(),
		RotationLimitMax:          mmath.NewMVec3(),
		SpringConstantTranslation: mmath.NewMVec3(),
		SpringConstantRotation:    mmath.NewMVec3(),
	}
}

func (param *JointParam) String() string {
	return fmt.Sprintf("TranslationLimitMin: %v, TranslationLimitMax: %v, RotationLimitMin: %v, RotationLimitMax: %v, SpringConstantTranslation: %v, SpringConstantRotation: %v",
		param.TranslationLimitMin, param.TranslationLimitMax, param.RotationLimitMin, param.RotationLimitMax, param.SpringConstantTranslation, param.SpringConstantRotation)

}

type Joint struct {
	index           int          // ジョイントINDEX
	name            string       // ジョイント名
	englishName     string       // ジョイント英名
	JointType       byte         // Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
	RigidBodyIndexA int          // 関連剛体AのIndex
	RigidBodyIndexB int          // 関連剛体BのIndex
	Position        *mmath.MVec3 // 位置(x,y,z)
	Rotation        *mmath.MVec3 // 回転
	JointParam      *JointParam  // ジョイントパラメーター
	IsSystem        bool
}

func NewJoint() *Joint {
	return &Joint{
		index:       -1,
		name:        "",
		englishName: "",
		// Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		JointType:       0,
		RigidBodyIndexA: -1,
		RigidBodyIndexB: -1,
		Position:        mmath.NewMVec3(),
		Rotation:        mmath.NewMVec3(),
		JointParam:      NewJointParam(),
		IsSystem:        false,
	}
}

func (joint *Joint) Index() int {
	return joint.index
}

func (joint *Joint) SetIndex(index int) {
	joint.index = index
}

func (joint *Joint) Name() string {
	return joint.name
}

func (joint *Joint) SetName(name string) {
	joint.name = name
}

func (joint *Joint) EnglishName() string {
	return joint.englishName
}

func (joint *Joint) SetEnglishName(englishName string) {
	joint.englishName = englishName
}

func (joint *Joint) IsValid() bool {
	return joint != nil && joint.index >= 0
}

func (joint *Joint) Copy() core.IIndexNameModel {
	return &Joint{
		index:           joint.index,
		name:            joint.name,
		englishName:     joint.englishName,
		JointType:       joint.JointType,
		RigidBodyIndexA: joint.RigidBodyIndexA,
		RigidBodyIndexB: joint.RigidBodyIndexB,
		Position:        joint.Position.Copy(),
		Rotation:        joint.Rotation.Copy(),
		JointParam:      joint.JointParam,
		IsSystem:        joint.IsSystem,
	}
}

func NewJointByName(name string) *Joint {
	joint := NewJoint()
	joint.SetName(name)
	return joint
}
