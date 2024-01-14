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
	TranslationLimitMin       mvec3.T
	TranslationLimitMax       mvec3.T
	RotationLimitMin          mrotation.T
	RotationLimitMax          mrotation.T
	SpringConstantTranslation mvec3.T
	SpringConstantRotation    mvec3.T
}

type Joint struct {
	index_model.IndexModel
	Index           int
	Name            string
	EnglishName     string
	JointType       int
	RigidbodyIndexA int
	RigidbodyIndexB int
	Position        mvec3.T
	Rotation        mrotation.T
	Param           Param
	IsSystem        bool
}

func NewJoint(
	index int,
	name string,
	englishName string,
	jointType int,
	rigidbodyIndexA int,
	rigidbodyIndexB int,
	position mvec3.T,
	rotation mrotation.T,
	param Param,
	isSystem bool,
) *Joint {
	return &Joint{
		Index:           index,
		Name:            name,
		EnglishName:     englishName,
		JointType:       jointType,
		RigidbodyIndexA: rigidbodyIndexA,
		RigidbodyIndexB: rigidbodyIndexB,
		Position:        position,
		Rotation:        rotation,
		Param:           param,
		IsSystem:        isSystem,
	}
}

// ジョイントリスト
type Joints struct {
	index_model.IndexModelCorrection[*Joint]
}

func NewJoints(name string) *Joints {
	return &Joints{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Joint](),
	}
}
