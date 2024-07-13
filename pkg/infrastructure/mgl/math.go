//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// Gl OpenGL座標系に変換された3次元ベクトルを返します
func NewGlVec3FromMVec3(v *mmath.MVec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ())}
}

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func NewGlMat4FromMMat4(m *mmath.MMat4) mgl32.Mat4 {
	mat := mgl32.Mat4{
		float32(m[0]), float32(-m[4]), float32(-m[8]), float32(m[12]),
		float32(-m[1]), float32(m[5]), float32(m[9]), float32(m[13]),
		float32(-m[2]), float32(m[6]), float32(m[10]), float32(m[14]),
		float32(-m[3]), float32(m[7]), float32(m[11]), float32(m[15]),
	}
	return mat
}

// NewMMat4ByMgl OpenGL座標系からMMD座標系に変換された行列を返します
func NewMMat4ByMgl(m *mgl32.Mat4) *mmath.MMat4 {
	mm := mmath.NewMMat4ByValues(
		float64(m.Col(0).X()), float64(-m.Col(1).X()), float64(-m.Col(2).X()), float64(-m.Col(3).X()),
		float64(-m.Col(0).Y()), float64(m.Col(1).Y()), float64(m.Col(2).Y()), float64(m.Col(3).Y()),
		float64(-m.Col(0).Z()), float64(m.Col(1).Z()), float64(m.Col(2).Z()), float64(m.Col(3).Z()),
		float64(m.Col(0).W()), float64(m.Col(1).W()), float64(m.Col(2).W()), float64(m.Col(3).W()),
	)
	m = nil
	return mm
}
