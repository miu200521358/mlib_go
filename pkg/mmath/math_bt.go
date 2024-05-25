//go:build windows
// +build windows

package mmath

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
)

// Gl OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) GL() mgl32.Vec3 {
	return mgl32.Vec3{float32(v.GetX()), float32(v.GetY()), float32(-v.GetZ())}
}

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) Bullet() mbt.BtVector3 {
	return mbt.NewBtVector3(float32(v.GetX()), float32(v.GetY()), float32(-v.GetZ()))
}

// Bullet+OpenGL座標系に変換されたクォータニオンベクトルを返します
func (v *MQuaternion) Bullet() mbt.BtQuaternion {
	return mbt.NewBtQuaternion(float32(v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(v.GetW()))
}

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (m *MMat4) GL() mgl32.Mat4 {
	return mgl32.Mat4{
		float32(m[0]), float32(-m[1]), float32(m[8]), float32(m[12]),
		float32(-m[4]), float32(m[5]), float32(m[6]), float32(m[13]),
		float32(m[2]), float32(m[9]), float32(m[10]), float32(m[14]),
		float32(-m[3]), float32(m[7]), float32(m[11]), float32(m[15]),
	}
}

// NewMMat4ByMgl OpenGL座標系からMMD座標系に変換された行列を返します
func NewMMat4ByMgl(m *mgl32.Mat4) *MMat4 {
	return NewMMat4ByValues(
		float64(m[0]), float64(m[1]), float64(m[2]), float64(m[3]),
		float64(m[4]), float64(m[5]), float64(m[6]), float64(m[7]),
		float64(m[8]), float64(m[9]), float64(m[10]), float64(m[11]),
		float64(m[12]), float64(m[13]), float64(m[14]), float64(m[15]),
	)
}
