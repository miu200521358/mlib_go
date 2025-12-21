package mmodel

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

// JointParam はジョイントパラメータを表します。
type JointParam struct {
	TranslationLimitMin       *mmath.Vec3 // 移動制限-下限
	TranslationLimitMax       *mmath.Vec3 // 移動制限-上限
	RotationLimitMin          *mmath.Vec3 // 回転制限-下限
	RotationLimitMax          *mmath.Vec3 // 回転制限-上限
	SpringConstantTranslation *mmath.Vec3 // バネ定数-移動
	SpringConstantRotation    *mmath.Vec3 // バネ定数-回転
}

// NewJointParam は新しいジョイントパラメータを生成します。
func NewJointParam() *JointParam {
	return &JointParam{
		TranslationLimitMin:       mmath.NewVec3(),
		TranslationLimitMax:       mmath.NewVec3(),
		RotationLimitMin:          mmath.NewVec3(),
		RotationLimitMax:          mmath.NewVec3(),
		SpringConstantTranslation: mmath.NewVec3(),
		SpringConstantRotation:    mmath.NewVec3(),
	}
}

// String は文字列表現を返します。
func (p *JointParam) String() string {
	return fmt.Sprintf("TranslationLimitMin: %v, TranslationLimitMax: %v, RotationLimitMin: %v, RotationLimitMax: %v, SpringConstantTranslation: %v, SpringConstantRotation: %v",
		p.TranslationLimitMin, p.TranslationLimitMax, p.RotationLimitMin, p.RotationLimitMax, p.SpringConstantTranslation, p.SpringConstantRotation)
}

// Joint はジョイントを表します。
type Joint struct {
	mcore.IndexNameModel
	JointType       byte        // Joint種類（0:スプリング6DOF）
	RigidBodyIndexA int         // 関連剛体AのIndex
	RigidBodyIndexB int         // 関連剛体BのIndex
	Position        *mmath.Vec3 // 位置
	Rotation        *mmath.Vec3 // 回転
	JointParam      *JointParam // ジョイントパラメータ
	IsSystem        bool        // システム追加ジョイント
}

// NewJoint は新しいジョイントを生成します。
func NewJoint() *Joint {
	return &Joint{
		IndexNameModel:  *mcore.NewIndexNameModel(-1, "", ""),
		JointType:       0,
		RigidBodyIndexA: -1,
		RigidBodyIndexB: -1,
		Position:        mmath.NewVec3(),
		Rotation:        mmath.NewVec3(),
		JointParam:      NewJointParam(),
		IsSystem:        false,
	}
}

// NewJointByName は名前を指定して新しいジョイントを生成します。
func NewJointByName(name string) *Joint {
	j := NewJoint()
	j.SetName(name)
	return j
}

// IsValid はジョイントが有効かどうかを返します。
func (j *Joint) IsValid() bool {
	return j != nil && j.Index() >= 0
}

// Copy は深いコピーを作成します。
func (j *Joint) Copy() (*Joint, error) {
	cp := &Joint{}
	if err := deepcopy.Copy(cp, j); err != nil {
		return nil, err
	}
	return cp, nil
}
