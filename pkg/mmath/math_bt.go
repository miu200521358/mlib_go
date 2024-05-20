//go:build windows
// +build windows

package mmath

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
)

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (v *MQuaternion) GL() [4]float32 {
	return [4]float32{float32(v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(v.GetW())}
}

// Bullet+OpenGL座標系に変換されたクォータニオンベクトルを返します
func (v *MQuaternion) Bullet() mbt.BtQuaternion {
	return mbt.NewBtQuaternion(float32(-v.GetX()), float32(-v.GetY()), float32(v.GetZ()), float32(v.GetW()))
}

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (m *MMat4) GL() mgl32.Mat4 {
	return mgl32.Mat4{
		float32(m[0]), float32(m[4]), float32(m[8]), float32(m[12]),
		float32(m[1]), float32(m[5]), float32(m[9]), float32(m[13]),
		float32(m[2]), float32(m[6]), float32(m[10]), float32(m[14]),
		float32(m[3]), float32(m[7]), float32(m[11]), float32(m[15]),
	}
}

// NewMMat4ByMgl OpenGL座標系からMMD座標系に変換された行列を返します
func NewMMat4ByMgl(m *mgl32.Mat4) *MMat4 {
	return NewMMat4ByValues(
		float64(m.Col(0).X()), float64(m.Col(1).X()), float64(m.Col(2).X()), float64(m.Col(3).X()),
		float64(m.Col(0).Y()), float64(m.Col(1).Y()), float64(m.Col(2).Y()), float64(m.Col(3).Y()),
		float64(m.Col(0).Z()), float64(m.Col(1).Z()), float64(m.Col(2).Z()), float64(m.Col(3).Z()),
		float64(m.Col(0).W()), float64(m.Col(1).W()), float64(m.Col(2).W()), float64(m.Col(3).W()),
	)
}

// Bullet+OpenGL座標系に変換された行列ベクトルを返します
func (v *MMat4) Bullet() mbt.BtMatrix3x3 {
	glMat := v.GL()
	return mbt.NewBtMatrix3x3(
		float32(glMat[0]), float32(glMat[1]), float32(glMat[2]),
		float32(glMat[4]), float32(glMat[5]), float32(glMat[6]),
		float32(glMat[8]), float32(glMat[9]), float32(glMat[10]))
}

// Gl OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) GL() mgl32.Vec3 {
	return mgl32.Vec3{float32(v.GetX()), float32(v.GetY()), float32(v.GetZ())}
}

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) Bullet() mbt.BtVector3 {
	return mbt.NewBtVector3(float32(v.GetX()), float32(v.GetY()), float32(v.GetZ()))
}

// GL OpenGL座標系に変換された2次元ベクトルを返します
func (v *MVec2) GL() mgl32.Vec2 {
	return mgl32.Vec2{float32(v.GetX()), float32(v.GetY())}
}

// GL OpenGL座標系に変換された4次元ベクトルを返します
func (v *MVec4) GL() mgl32.Vec4 {
	return mgl32.Vec4{float32(v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(v.GetW())}
}

// GL OpenGL座標系に変換されたベクトルを返します
func (m *MMat3) GL() mgl32.Mat3 {
	return mgl32.Mat3{
		float32(m[0]), float32(m[1]), float32(m[2]),
		float32(m[3]), float32(m[4]), float32(m[5]),
		float32(m[6]), float32(m[7]), float32(m[8]),
	}
}
