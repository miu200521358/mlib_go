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

func (m *MRotation) Quaternion() *MQuaternion {
	if m.quaternion == nil {
		m.quaternion = NewMQuaternion()
	}
	return m.quaternion
}

func (m *MRotation) SetQuaternion(v *MQuaternion) {
	m.quaternion = v
	m.radians = v.ToRadians()
	m.degrees = &MVec3{
		180.0 * m.radians.X / math.Pi,
		180.0 * m.radians.Y / math.Pi,
		180.0 * m.radians.Z / math.Pi,
	}
}

func (m *MRotation) Radians() *MVec3 {
	if m.radians == nil {
		m.radians = NewMVec3()
	}
	return m.radians
}

func (m *MRotation) GetRadiansMMD() *MVec3 {
	if m.radians == nil {
		m.radians = NewMVec3()
	}
	return m.radians.MMD()
}

func (m *MRotation) SetRadians(v *MVec3) {
	m.radians = v
	m.degrees = &MVec3{
		180.0 * v.X / math.Pi,
		180.0 * v.Y / math.Pi,
		180.0 * v.Z / math.Pi,
	}
	m.quaternion = NewMQuaternionFromRadians(v.X, v.Y, v.Z)
}

func (m *MRotation) Degrees() *MVec3 {
	if m.degrees == nil {
		m.degrees = NewMVec3()
	}
	return m.degrees
}

func (m *MRotation) DegreesMMD() *MVec3 {
	if m.degrees == nil {
		m.degrees = NewMVec3()
	}
	return m.degrees.MMD()
}

func (m *MRotation) SetDegrees(v *MVec3) {
	m.degrees = v
	m.radians = &MVec3{
		math.Pi * v.X / 180.0,
		math.Pi * v.Y / 180.0,
		math.Pi * v.Z / 180.0,
	}
	m.quaternion = NewMQuaternionFromRadians(m.radians.X, m.radians.Y, m.radians.Z)
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
