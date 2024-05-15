package mmath

import (
	"fmt"
	"math"
)

type MRotation struct {
	radians    MVec3
	degrees    MVec3
	quaternion MQuaternion
}

func NewRotation() *MRotation {
	model := &MRotation{
		radians:    NewMVec3(),
		degrees:    NewMVec3(),
		quaternion: NewMQuaternion(),
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
	return &m.quaternion
}

func (m *MRotation) SetQuaternion(v *MQuaternion) {
	m.quaternion = *v
	m.radians = *v.ToRadians()
	m.degrees = MVec3{
		180.0 * m.radians.GetX() / math.Pi,
		180.0 * m.radians.GetY() / math.Pi,
		180.0 * m.radians.GetZ() / math.Pi,
	}
}

func (m *MRotation) GetRadians() *MVec3 {
	return &m.radians
}

func (m *MRotation) SetRadians(v *MVec3) {
	m.radians = *v
	m.degrees = MVec3{
		180.0 * v.GetX() / math.Pi,
		180.0 * v.GetY() / math.Pi,
		180.0 * v.GetZ() / math.Pi,
	}
	m.quaternion = NewMQuaternionFromRadians(v.GetX(), v.GetY(), v.GetZ())
}

func (m *MRotation) GetDegrees() *MVec3 {
	return &m.degrees
}

func (m *MRotation) SetDegrees(v *MVec3) {
	m.degrees = *v
	m.radians = MVec3{
		math.Pi * v.GetX() / 180.0,
		math.Pi * v.GetY() / 180.0,
		math.Pi * v.GetZ() / 180.0,
	}
	m.quaternion = NewMQuaternionFromRadians(m.radians.GetX(), m.radians.GetY(), m.radians.GetZ())
}

// Copy
func (rot *MRotation) Copy() *MRotation {
	copied := &MRotation{
		radians:    rot.radians.Copy(),
		degrees:    rot.degrees.Copy(),
		quaternion: rot.quaternion.Copy(),
	}
	return copied
}

// Mul
func (rot *MRotation) Mul(v *MRotation) {
	qq := rot.quaternion.Mul(&v.quaternion)
	rot.SetQuaternion(qq)
}

func (rot *MRotation) String() string {
	return fmt.Sprintf("degrees: %s, radians: %s, quat: %s",
		rot.degrees.String(), rot.radians.String(), rot.quaternion.String())
}
