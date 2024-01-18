package joint

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mrotation"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"

)

type LimitParam struct {
	LimitMin mvec3.T
	LimitMax mvec3.T
}

type Param struct {
	// 移動制限-下限(x,y,z)
	TranslationLimitMin mvec3.T
	// 移動制限-上限(x,y,z)
	TranslationLimitMax mvec3.T
	// 回転制限-下限
	RotationLimitMin mrotation.T
	// 回転制限-上限
	RotationLimitMax mrotation.T
	// バネ定数-移動(x,y,z)
	SpringConstantTranslation mvec3.T
	// バネ定数-回転(x,y,z)
	SpringConstantRotation mrotation.T
}

func NewParam() *Param {
	return &Param{
		TranslationLimitMin:       mvec3.T{},
		TranslationLimitMax:       mvec3.T{},
		RotationLimitMin:          mrotation.T{},
		RotationLimitMax:          mrotation.T{},
		SpringConstantTranslation: mvec3.T{},
		SpringConstantRotation:    mrotation.T{},
	}
}

type Joint struct {
	*index_model.IndexModel
	// Joint名
	Name string
	// Joint名英
	EnglishName string
	// Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
	JointType byte
	// 関連剛体AのIndex
	RigidbodyIndexA int
	// 関連剛体BのIndex
	RigidbodyIndexB int
	// 位置(x,y,z)
	Position mvec3.T
	// 回転
	Rotation mrotation.T
	// ジョイントパラメーター
	Param    Param
	IsSystem bool
}

func NewJoint() *Joint {
	return &Joint{
		IndexModel:  &index_model.IndexModel{Index: -1},
		Name:        "",
		EnglishName: "",
		// Joint種類 - 0:スプリング6DOF   | PMX2.0では 0 のみ(拡張用)
		JointType:       0,
		RigidbodyIndexA: -1,
		RigidbodyIndexB: -1,
		Position:        mvec3.T{},
		Rotation:        mrotation.T{},
		Param:           *NewParam(),
		IsSystem:        false,
	}
}

// ジョイントリスト
type Joints struct {
	*index_model.IndexModelCorrection[*Joint]
}

func NewJoints() *Joints {
	return &Joints{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Joint](),
	}
}
