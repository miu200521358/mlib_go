package mrotation

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/math/mquaternion"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

type T struct {
	radians    mvec3.T
	degrees    mvec3.T
	quaternion mquaternion.T
}

// NewBaseRotationModelByRadians はラジアン角度からで回転を表すモデルを生成します。
func NewBaseRotationModelByRadians(vRadians *mvec3.T) *T {
	model := &T{
		radians:    mvec3.T{},
		degrees:    mvec3.T{},
		quaternion: mquaternion.T{},
	}
	model.SetRadians(*vRadians)
	return model
}

// NewBaseRotationModelByDegrees は度数角度からで回転を表すモデルを生成します。
func NewBaseRotationModelByDegrees(vDegrees *mvec3.T) *T {
	model := &T{
		radians:    mvec3.T{},
		degrees:    mvec3.T{},
		quaternion: mquaternion.T{},
	}
	model.SetDegrees(*vDegrees)
	return model
}

// NewBaseRotationModelByQuaternion はクォータニオンからで回転を表すモデルを生成します。
func NewBaseRotationModelByQuaternion(vQuaternion *mquaternion.T) *T {
	model := &T{
		radians:    mvec3.T{},
		degrees:    mvec3.T{},
		quaternion: mquaternion.T{},
	}
	model.SetQuaternion(*vQuaternion)
	return model
}

func (m *T) GetQuaternion() *mquaternion.T {
	return &m.quaternion
}

func (m *T) SetQuaternion(v mquaternion.T) {
	m.quaternion = v
	m.radians = v.ToEulerAngles()
	m.degrees = mvec3.T{
		180.0 * m.radians.GetX() / math.Pi,
		180.0 * m.radians.GetY() / math.Pi,
		180.0 * m.radians.GetZ() / math.Pi,
	}
}

func (m *T) GetRadians() *mvec3.T {
	return &m.radians
}

func (m *T) SetRadians(v mvec3.T) {
	m.radians = v
	m.degrees = mvec3.T{
		180.0 * v.GetX() / math.Pi,
		180.0 * v.GetY() / math.Pi,
		180.0 * v.GetZ() / math.Pi,
	}
	m.quaternion = mquaternion.FromEulerAngles(v.GetX(), v.GetY(), v.GetZ())
}

func (m *T) GetDegrees() *mvec3.T {
	return &m.degrees
}

func (m *T) SetDegrees(v mvec3.T) {
	m.degrees = v
	m.radians = mvec3.T{
		math.Pi * v.GetX() / 180.0,
		math.Pi * v.GetY() / 180.0,
		math.Pi * v.GetZ() / 180.0,
	}
	m.quaternion = mquaternion.FromEulerAnglesDegrees(v.GetX(), v.GetY(), v.GetZ())
}

// Copy
func (rot *T) Copy() *T {
	return &T{*rot.radians.Copy(), *rot.degrees.Copy(), *rot.quaternion.Copy()}
}
