package mmath

import (
	"math"

)

type MRotation struct {
	radians    MVec3
	degrees    MVec3
	quaternion MQuaternion
}

// NewBaseRotationModelByRadians はラジアン角度からで回転を表すモデルを生成します。
func NewBaseRotationModelByRadians(vRadians *MVec3) *MRotation {
	model := &MRotation{
		radians:    MVec3{},
		degrees:    MVec3{},
		quaternion: MQuaternion{},
	}
	model.SetRadians(*vRadians)
	return model
}

// NewBaseRotationModelByDegrees は度数角度からで回転を表すモデルを生成します。
func NewBaseRotationModelByDegrees(vDegrees *MVec3) *MRotation {
	model := &MRotation{
		radians:    MVec3{},
		degrees:    MVec3{},
		quaternion: MQuaternion{},
	}
	model.SetDegrees(*vDegrees)
	return model
}

// NewBaseRotationModelByQuaternion はクォータニオンからで回転を表すモデルを生成します。
func NewBaseRotationModelByQuaternion(vQuaternion *MQuaternion) *MRotation {
	model := &MRotation{
		radians:    MVec3{},
		degrees:    MVec3{},
		quaternion: MQuaternion{},
	}
	model.SetQuaternion(*vQuaternion)
	return model
}

func (m *MRotation) GetQuaternion() *MQuaternion {
	return &m.quaternion
}

func (m *MRotation) SetQuaternion(v MQuaternion) {
	m.quaternion = v
	m.radians = v.ToEulerAngles()
	m.degrees = MVec3{
		180.0 * m.radians.GetX() / math.Pi,
		180.0 * m.radians.GetY() / math.Pi,
		180.0 * m.radians.GetZ() / math.Pi,
	}
}

func (m *MRotation) GetRadians() *MVec3 {
	return &m.radians
}

func (m *MRotation) SetRadians(v MVec3) {
	m.radians = v
	m.degrees = MVec3{
		180.0 * v.GetX() / math.Pi,
		180.0 * v.GetY() / math.Pi,
		180.0 * v.GetZ() / math.Pi,
	}
	m.quaternion = FromEulerAngles(v.GetX(), v.GetY(), v.GetZ())
}

func (m *MRotation) GetDegrees() *MVec3 {
	return &m.degrees
}

func (m *MRotation) SetDegrees(v MVec3) {
	m.degrees = v
	m.radians = MVec3{
		math.Pi * v.GetX() / 180.0,
		math.Pi * v.GetY() / 180.0,
		math.Pi * v.GetZ() / 180.0,
	}
	m.quaternion = FromEulerAnglesDegrees(v.GetX(), v.GetY(), v.GetZ())
}

// Copy
func (rot *MRotation) Copy() *MRotation {
	return &MRotation{*rot.radians.Copy(), *rot.degrees.Copy(), *rot.quaternion.Copy()}
}
