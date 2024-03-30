package mmath

import (
	"math"
)

type MRotation struct {
	radians    *MVec3
	degrees    *MVec3
	quaternion *MQuaternion
}

func NewRotationModel() *MRotation {
	model := &MRotation{
		radians:    NewMVec3(),
		degrees:    NewMVec3(),
		quaternion: NewMQuaternion(),
	}
	return model
}

// NewRotationModelByRadians はラジアン角度からで回転を表すモデルを生成します。
func NewRotationModelByRadians(vRadians *MVec3) *MRotation {
	model := &MRotation{
		radians:    NewMVec3(),
		degrees:    NewMVec3(),
		quaternion: NewMQuaternion(),
	}
	if vRadians.Length() > 0 {
		model.SetRadians(vRadians)
	}
	return model
}

// NewRotationModelByDegrees は度数角度からで回転を表すモデルを生成します。
func NewRotationModelByDegrees(vDegrees *MVec3) *MRotation {
	model := &MRotation{
		radians:    NewMVec3(),
		degrees:    NewMVec3(),
		quaternion: NewMQuaternion(),
	}
	if vDegrees.Length() > 0 {
		model.SetDegrees(vDegrees)
	}
	return model
}

// NewRotationModelByQuaternion はクォータニオンからで回転を表すモデルを生成します。
func NewRotationModelByQuaternion(vQuaternion *MQuaternion) *MRotation {
	model := &MRotation{
		radians:    NewMVec3(),
		degrees:    NewMVec3(),
		quaternion: NewMQuaternion(),
	}
	if vQuaternion.GetXYZ().Length() > 0 {
		model.SetQuaternion(vQuaternion)
	}
	return model
}

func (m *MRotation) GetQuaternion() *MQuaternion {
	return m.quaternion
}

func (m *MRotation) SetQuaternion(v *MQuaternion) {
	m.quaternion = v
	m.radians = v.ToEulerAngles()
	m.degrees = &MVec3{
		180.0 * m.radians.GetX() / math.Pi,
		180.0 * m.radians.GetY() / math.Pi,
		180.0 * m.radians.GetZ() / math.Pi,
	}
}

func (m *MRotation) GetRadians() *MVec3 {
	return m.radians
}

func (m *MRotation) SetRadians(v *MVec3) {
	m.radians = v
	m.degrees = &MVec3{
		180.0 * v.GetX() / math.Pi,
		180.0 * v.GetY() / math.Pi,
		180.0 * v.GetZ() / math.Pi,
	}
	m.quaternion = NewMQuaternionFromEulerAngles(v.GetX(), v.GetY(), v.GetZ())
}

func (m *MRotation) GetDegrees() *MVec3 {
	return m.degrees
}

func (m *MRotation) SetDegrees(v *MVec3) {
	m.degrees = v
	m.radians = &MVec3{
		math.Pi * v.GetX() / 180.0,
		math.Pi * v.GetY() / 180.0,
		math.Pi * v.GetZ() / 180.0,
	}
	m.quaternion = NewMQuaternionFromEulerAngles(m.radians.GetX(), m.radians.GetY(), m.radians.GetZ())
}

// Copy
func (rot *MRotation) Copy() *MRotation {
	return &MRotation{
		degrees:    rot.degrees.Copy(),
		radians:    rot.radians.Copy(),
		quaternion: rot.quaternion.Copy(),
	}
}

// Add
func (rot *MRotation) Mul(v *MRotation) {
	rot.SetQuaternion(rot.quaternion.Mul(v.quaternion))
}
