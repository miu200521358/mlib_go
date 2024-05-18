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

func NewRotation() *MRotation {
	model := &MRotation{
		radians:    nil,
		degrees:    nil,
		quaternion: nil,
	}
	return model
}

// NewRotationByRadians はラジアン角度からで回転を表すモデルを生成します。
func NewRotationByRadians(vRadians *MVec3) *MRotation {
	model := NewRotation()
	model.SetRadians(vRadians)
	return model
}

// NewRotationByDegrees は度数角度からで回転を表すモデルを生成します。
func NewRotationByDegrees(vDegrees *MVec3) *MRotation {
	model := NewRotation()
	model.SetDegrees(vDegrees)
	return model
}

// NewRotationByQuaternion はクォータニオンからで回転を表すモデルを生成します。
func NewRotationByQuaternion(vQuaternion *MQuaternion) *MRotation {
	model := NewRotation()
	model.SetQuaternion(vQuaternion)
	return model
}

func (m *MRotation) GetQuaternion() *MQuaternion {
	if m.quaternion == nil {
		m.quaternion = NewMQuaternion()
	}
	return m.quaternion
}

func (m *MRotation) SetQuaternion(v *MQuaternion) {
	m.quaternion = v
	m.radians = v.ToRadians()
	m.degrees = &MVec3{
		180.0 * m.radians.GetX() / math.Pi,
		180.0 * m.radians.GetY() / math.Pi,
		180.0 * m.radians.GetZ() / math.Pi,
	}
}

func (m *MRotation) GetRadians() *MVec3 {
	if m.radians == nil {
		m.radians = NewMVec3()
	}
	return m.radians
}

func (m *MRotation) GetRadiansMMD() *MVec3 {
	if m.radians == nil {
		m.radians = NewMVec3()
	}
	if m.quaternion != nil {
		return m.quaternion.MMD().ToRadians()
	}
	return m.radians
}

func (m *MRotation) SetRadians(v *MVec3) {
	m.radians = v
	m.degrees = &MVec3{
		180.0 * v.GetX() / math.Pi,
		180.0 * v.GetY() / math.Pi,
		180.0 * v.GetZ() / math.Pi,
	}
	m.quaternion = NewMQuaternionFromRadians(v.GetX(), v.GetY(), v.GetZ())
}

func (m *MRotation) GetDegrees() *MVec3 {
	if m.degrees == nil {
		m.degrees = NewMVec3()
	}
	return m.degrees
}

func (m *MRotation) GetDegreesMMD() *MVec3 {
	if m.degrees == nil {
		m.degrees = NewMVec3()
	}
	if m.quaternion != nil {
		return m.quaternion.MMD().ToDegrees()
	}
	return m.degrees
}

func (m *MRotation) SetDegrees(v *MVec3) {
	m.degrees = v
	m.radians = &MVec3{
		math.Pi * v.GetX() / 180.0,
		math.Pi * v.GetY() / 180.0,
		math.Pi * v.GetZ() / 180.0,
	}
	m.quaternion = NewMQuaternionFromRadians(m.radians.GetX(), m.radians.GetY(), m.radians.GetZ())
}

// Copy
func (rot *MRotation) Copy() *MRotation {
	copied := NewRotation()
	if rot.radians != nil {
		copied.radians = &MVec3{rot.radians.GetX(), rot.radians.GetY(), rot.radians.GetZ()}
	}
	if rot.degrees != nil {
		copied.degrees = &MVec3{rot.degrees.GetX(), rot.degrees.GetY(), rot.degrees.GetZ()}
	}
	if rot.quaternion != nil {
		copied.quaternion = NewMQuaternionByValues(
			rot.quaternion.GetX(), rot.quaternion.GetY(), rot.quaternion.GetZ(), rot.quaternion.GetW())
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
