//go:build windows
// +build windows

package mmath

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Gl OpenGL座標系に変換された3次元ベクトルを返します
func NewGlVec3(v *MVec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(-v.X), float32(v.Y), float32(v.Z)}
}

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func NewGlMat4(m *MMat4) mgl32.Mat4 {
	mat := mgl32.Mat4{
		float32(m[0]), float32(-m[1]), float32(-m[2]), float32(m[3]),
		float32(-m[4]), float32(m[5]), float32(m[6]), float32(m[7]),
		float32(-m[8]), float32(m[9]), float32(m[10]), float32(m[11]),
		float32(-m[12]), float32(m[13]), float32(m[14]), float32(m[15]),
	}
	return mat
}
