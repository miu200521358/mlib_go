package mmath

import (
	"math"

	"github.com/jinzhu/copier"
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

func (m *MRotation) SetRadians(v *MVec3) {
	m.radians = v
	m.degrees = &MVec3{
		180.0 * v.GetX() / math.Pi,
		180.0 * v.GetY() / math.Pi,
		180.0 * v.GetZ() / math.Pi,
	}
	qq := NewMQuaternionFromRadians(v.GetX(), v.GetY(), v.GetZ())
	m.quaternion = qq
}

func (m *MRotation) GetDegrees() *MVec3 {
	if m.degrees == nil {
		m.degrees = NewMVec3()
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
	m.quaternion = qq
}

// Copy
func (rot *MRotation) Copy() *MRotation {
	copied := NewRotation()
	copier.CopyWithOption(copied, rot, copier.Option{DeepCopy: true})
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
