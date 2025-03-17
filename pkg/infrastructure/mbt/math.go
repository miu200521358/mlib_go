//go:build windows
// +build windows

package mbt

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func newBulletFromVec(vec *mmath.MVec3) bt.BtVector3 {
	return bt.NewBtVector3(float32(-vec.X), float32(vec.Y), float32(vec.Z))
}

// Bullet+OpenGL座標系に変換されたラジアン→クォータニオンベクトルを返します
func newBulletFromRad(rad *mmath.MVec3) bt.BtQuaternion {
	return bt.NewBtQuaternion(float32(-rad.Y), float32(rad.X), float32(-rad.Z))
}

// NewMMat4ByMgl OpenGL座標系からMMD座標系に変換された行列を返します
func newMMat4ByMgl(mat *mgl32.Mat4) *mmath.MMat4 {
	mm := mmath.NewMMat4ByValues(
		float64(mat[0]), float64(-mat[1]), float64(-mat[2]), float64(mat[3]),
		float64(-mat[4]), float64(mat[5]), float64(mat[6]), float64(mat[7]),
		float64(-mat[8]), float64(mat[9]), float64(mat[10]), float64(mat[11]),
		float64(-mat[12]), float64(mat[13]), float64(mat[14]), float64(mat[15]),
	)
	return mm
}
