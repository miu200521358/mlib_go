package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type LimitParam struct {
	LimitMin mmath.MVec3
	LimitMax mmath.MVec3
}

type JointParam struct {
	TranslationLimitMin       mmath.MVec3     // 移動制限-下限(x,y,z)
	TranslationLimitMax       mmath.MVec3     // 移動制限-上限(x,y,z)
	RotationLimitMin          mmath.MRotation // 回転制限-下限
	RotationLimitMax          mmath.MRotation // 回転制限-上限
	SpringConstantTranslation mmath.MVec3     // バネ定数-移動(x,y,z)
	SpringConstantRotation    mmath.MRotation // バネ定数-回転(x,y,z)
}

func NewJointParam() *JointParam {
	return &JointParam{
		TranslationLimitMin:       mmath.MVec3{},
		TranslationLimitMax:       mmath.MVec3{},
		RotationLimitMin:          mmath.MRotation{},
		RotationLimitMax:          mmath.MRotation{},
		SpringConstantTranslation: mmath.MVec3{},
		SpringConstantRotation:    mmath.MRotation{},
	}
}

type Joint struct {
	*mcore.IndexModel
	Name            string          // Joint名
	EnglishName     string          // Joint名英
	JointType       byte            // Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
	RigidbodyIndexA int             // 関連剛体AのIndex
	RigidbodyIndexB int             // 関連剛体BのIndex
	Position        mmath.MVec3     // 位置(x,y,z)
	Rotation        mmath.MRotation // 回転
	JointParam      JointParam      // ジョイントパラメーター
	IsSystem        bool
}

func NewJoint() *Joint {
	return &Joint{
		IndexModel:  &mcore.IndexModel{Index: -1},
		Name:        "",
		EnglishName: "",
		// Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		JointType:       0,
		RigidbodyIndexA: -1,
		RigidbodyIndexB: -1,
		Position:        mmath.MVec3{},
		Rotation:        mmath.MRotation{},
		JointParam:      *NewJointParam(),
		IsSystem:        false,
	}
}

// ジョイントリスト
type Joints struct {
	*mcore.IndexModelCorrection[*Joint]
}

func NewJoints() *Joints {
	return &Joints{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Joint](),
	}
}
