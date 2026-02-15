//go:build windows
// +build windows

// 指示: miu200521358
package mbullet

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet/bt"
)

// newBulletFromVec は MMD 座標系のベクトルを Bullet+OpenGL 座標へ変換する。
func newBulletFromVec(vec mmath.Vec3) bt.BtVector3 {
	return bt.NewBtVector3(float32(-vec.X), float32(vec.Y), float32(vec.Z))
}

// newBulletFromRad はラジアン表現を Bullet+OpenGL 座標のクォータニオンへ変換する。
func newBulletFromRad(rad mmath.Vec3) bt.BtQuaternion {
	return bt.NewBtQuaternion(float32(-rad.Y), float32(rad.X), float32(-rad.Z))
}

// newMglMat4FromMat4 は MMD 座標系の行列を OpenGL 座標の行列へ変換する。
func newMglMat4FromMat4(mat mmath.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(mat[0]), float32(-mat[1]), float32(-mat[2]), float32(mat[3]),
		float32(-mat[4]), float32(mat[5]), float32(mat[6]), float32(mat[7]),
		float32(-mat[8]), float32(mat[9]), float32(mat[10]), float32(mat[11]),
		float32(-mat[12]), float32(mat[13]), float32(mat[14]), float32(mat[15]),
	}
}

// newMat4FromMgl は OpenGL 座標系の行列を MMD 座標の行列へ変換する。
func newMat4FromMgl(mat *mgl32.Mat4) mmath.Mat4 {
	return mmath.NewMat4ByValues(
		float64(mat[0]), float64(-mat[1]), float64(-mat[2]), float64(mat[3]),
		float64(-mat[4]), float64(mat[5]), float64(mat[6]), float64(mat[7]),
		float64(-mat[8]), float64(mat[9]), float64(mat[10]), float64(mat[11]),
		float64(-mat[12]), float64(mat[13]), float64(mat[14]), float64(mat[15]),
	)
}
