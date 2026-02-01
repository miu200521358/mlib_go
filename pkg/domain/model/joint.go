// 指示: miu200521358
package model

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// JointParam はジョイントのパラメータを表す。
type JointParam struct {
	Position                  mmath.Vec3
	Rotation                  mmath.Vec3
	TranslationLimitMin       mmath.Vec3
	TranslationLimitMax       mmath.Vec3
	RotationLimitMin          mmath.Vec3
	RotationLimitMax          mmath.Vec3
	SpringConstantTranslation mmath.Vec3
	SpringConstantRotation    mmath.Vec3
}

// Joint はジョイント要素を表す。
type Joint struct {
	index           int
	name            string
	EnglishName     string
	RigidBodyIndexA int
	RigidBodyIndexB int
	Param           JointParam
}

// Index はジョイント index を返す。
func (j *Joint) Index() int {
	return j.index
}

// SetIndex はジョイント index を設定する。
func (j *Joint) SetIndex(index int) {
	j.index = index
}

// Name はジョイント名を返す。
func (j *Joint) Name() string {
	return j.name
}

// SetName はジョイント名を設定する。
func (j *Joint) SetName(name string) {
	j.name = name
}

// IsValid はジョイントが有効か判定する。
func (j *Joint) IsValid() bool {
	return j != nil && j.index >= 0
}
