package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/mmath"
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
		RotationLimitMin:          mmath.NewRotation(),
		RotationLimitMax:          mmath.NewRotation(),
		SpringConstantTranslation: mmath.NewMVec3(),
		SpringConstantRotation:    mmath.NewRotation(),
	}
}

type Joint struct {
	*core.IndexNameModel
	JointType       byte             // Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
	RigidbodyIndexA int              // 関連剛体AのIndex
	RigidbodyIndexB int              // 関連剛体BのIndex
	Position        *mmath.MVec3     // 位置(x,y,z)
	Rotation        *mmath.MRotation // 回転
	JointParam      *JointParam      // ジョイントパラメーター
	IsSystem        bool
}

func NewJoint() *Joint {
	return &Joint{
		IndexNameModel: &core.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		// Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		JointType:       0,
		RigidbodyIndexA: -1,
		RigidbodyIndexB: -1,
		Position:        mmath.NewMVec3(),
		Rotation:        mmath.NewRotation(),
		JointParam:      NewJointParam(),
		IsSystem:        false,
	}
}

func (j *Joint) Copy() core.IIndexNameModel {
	copied := NewJoint()
	copier.CopyWithOption(copied, j, copier.Option{DeepCopy: true})
	return copied
}

func NewJointByName(name string) *Joint {
	j := NewJoint()
	j.Name = name
	return j
}

// ジョイントリスト
type Joints struct {
	*core.IndexNameModels[*Joint]
}

func NewJoints(count int) *Joints {
	return &Joints{
		IndexNameModels: core.NewIndexNameModels[*Joint](count, func() *Joint { return nil }),
	}
}
