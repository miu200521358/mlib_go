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
		qq := NewMQuaternion()
		m.quaternion = &qq
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
		v := NewMVec3()
		m.radians = &v
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
	qq := NewMQuaternionFromRadians(v.GetX(), v.GetY(), v.GetZ())
	m.quaternion = &qq
}

func (m *MRotation) GetDegrees() *MVec3 {
	if m.degrees == nil {
		v := NewMVec3()
		m.degrees = &v
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
	qq := NewMQuaternionFromRadians(m.radians.GetX(), m.radians.GetY(), m.radians.GetZ())
	m.quaternion = &qq
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
		quat := NewMQuaternionByValues(rot.quaternion.GetX(), rot.quaternion.GetY(), rot.quaternion.GetZ(), rot.quaternion.GetW())
		copied.quaternion = &quat
	}
	return copied
}

// Mul
func (rot *MRotation) Mul(v *MRotation) {
	if rot.quaternion == nil {
		qq := NewMQuaternion()
		rot.quaternion = &qq
	}
	if v.quaternion == nil {
		qq := NewMQuaternion()
		v.quaternion = &qq
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
