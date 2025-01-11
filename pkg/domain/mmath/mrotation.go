package mmath

import (
	"fmt"
	"math"
	"reflect"
)

type MRotation struct {
	radians    *MVec3
	degrees    *MVec3
	quaternion *MQuaternion
}

func NewMRotation() *MRotation {
	model := &MRotation{
		radians:    nil,
		degrees:    nil,
		quaternion: nil,
	}
	return model
}

// NewMRotationFromRadians はラジアン角度からで回転を表すモデルを生成します。
func NewMRotationFromRadians(vRadians *MVec3) *MRotation {
	model := NewMRotation()
	model.SetRadians(vRadians)
	return model
}

// NewMRotationFromDegrees は度数角度からで回転を表すモデルを生成します。
func NewMRotationFromDegrees(vDegrees *MVec3) *MRotation {
	model := NewMRotation()
	model.SetDegrees(vDegrees)
	return model
}

// NewMRotationFromQuaternion はクォータニオンからで回転を表すモデルを生成します。
func NewMRotationFromQuaternion(vQuaternion *MQuaternion) *MRotation {
	model := NewMRotation()
	model.SetQuaternion(vQuaternion)
	return model
}

func (rot *MRotation) Quaternion() *MQuaternion {
	if rot.quaternion == nil {
		rot.quaternion = NewMQuaternion()
	}
	return rot.quaternion
}

func (rot *MRotation) SetQuaternion(v *MQuaternion) {
	rot.quaternion = v
	rot.radians = v.ToRadians()
	rot.degrees = &MVec3{
		180.0 * rot.radians.X / math.Pi,
		180.0 * rot.radians.Y / math.Pi,
		180.0 * rot.radians.Z / math.Pi,
	}
}

func (rot *MRotation) Radians() *MVec3 {
	if rot.radians == nil {
		rot.radians = NewMVec3()
	}
	return rot.radians
}

func (rot *MRotation) GetRadiansMMD() *MVec3 {
	if rot.radians == nil {
		rot.radians = NewMVec3()
	}
	return rot.radians.MMD()
}

func (rot *MRotation) SetRadians(v *MVec3) {
	rot.radians = v
	rot.degrees = &MVec3{
		180.0 * v.X / math.Pi,
		180.0 * v.Y / math.Pi,
		180.0 * v.Z / math.Pi,
	}
	rot.quaternion = NewMQuaternionFromRadians(v.X, v.Y, v.Z)
}

func (rot *MRotation) Degrees() *MVec3 {
	if rot.degrees == nil {
		rot.degrees = NewMVec3()
	}
	return rot.degrees
}

func (rot *MRotation) DegreesMMD() *MVec3 {
	if rot.degrees == nil {
		rot.degrees = NewMVec3()
	}
	return rot.degrees.MMD()
}

func (rot *MRotation) SetDegrees(v *MVec3) {
	rot.degrees = v
	rot.radians = &MVec3{
		math.Pi * v.X / 180.0,
		math.Pi * v.Y / 180.0,
		math.Pi * v.Z / 180.0,
	}
	rot.quaternion = NewMQuaternionFromRadians(rot.radians.X, rot.radians.Y, rot.radians.Z)
}

// Copy
func (rot *MRotation) Copy() *MRotation {
	copied := NewMRotation()
	if rot.radians != nil {
		copied.radians = &MVec3{rot.radians.X, rot.radians.Y, rot.radians.Z}
	}
	if rot.degrees != nil {
		copied.degrees = &MVec3{rot.degrees.X, rot.degrees.Y, rot.degrees.Z}
	}
	if rot.quaternion != nil {
		copied.quaternion = NewMQuaternionByValues(
			rot.quaternion.X, rot.quaternion.Y, rot.quaternion.Z, rot.quaternion.W)
	}
	return copied
}

// Mul
func (rot *MRotation) Mul(v *MRotation) {
	if rot.quaternion == nil {
		rot.quaternion = NewMQuaternion()
	}
	if v.quaternion == nil {
		v.quaternion = NewMQuaternion()
	}
	qq := rot.quaternion.Mul(v.quaternion)
	rot.SetQuaternion(qq)
}

func (rot *MRotation) String() string {
	stringOrNil := func(v fmt.Stringer) string {
		if v == nil || reflect.ValueOf(v).IsNil() {
			return "nil"
		}
		return v.String()
	}

	return fmt.Sprintf("degrees: %s, radians: %s, quat: %s",
		stringOrNil(rot.degrees), stringOrNil(rot.radians), stringOrNil(rot.quaternion))
}
