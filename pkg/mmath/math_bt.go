//go:build windows
// +build windows

package mmath

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
)

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (v *MQuaternion) GL() [4]float32 {
	return [4]float32{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(-v.GetW())}
}

// Bullet+OpenGL座標系に変換されたクォータニオンベクトルを返します
func (v *MQuaternion) Bullet() mbt.BtQuaternion {
	return mbt.NewBtQuaternion(float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(-v.GetW()))
}

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (m *MMat4) GL() *mgl32.Mat4 {
	tm := m.Transpose()
	mat := mgl32.Mat4{
		float32(tm[0][0]), float32(-tm[0][1]), float32(-tm[0][2]), float32(tm[0][3]),
		float32(-tm[1][0]), float32(tm[1][1]), float32(tm[1][2]), float32(tm[1][3]),
		float32(-tm[2][0]), float32(tm[2][1]), float32(tm[2][2]), float32(tm[2][3]),
		float32(-tm[3][0]), float32(tm[3][1]), float32(tm[3][2]), float32(tm[3][3]),
	}
	return &mat
}

// Bullet+OpenGL座標系に変換された行列ベクトルを返します
func (v *MMat4) Bullet() mbt.BtMatrix3x3 {
	glMat := v.GL()
	return mbt.NewBtMatrix3x3(
		float32(glMat[0]), float32(glMat[1]), float32(glMat[2]),
		float32(glMat[4]), float32(glMat[5]), float32(glMat[6]),
		float32(glMat[8]), float32(glMat[9]), float32(glMat[10]))
}

// Bullet+OpenGL座標系に変換 + Z軸の逆行列を加味した行列ベクトルを返します
func (v *MMat4) BulletInvZ() mbt.BtMatrix3x3 {
	invZ := NewMMat4()
	invZ.ScaleVec3(&MVec3{1, 1, -1})

	mat := NewMMat4()
	mat.Mul(invZ)
	mat.Mul(v)
	mat.Mul(invZ)

	return mbt.NewBtMatrix3x3(
		float32(mat[0][0]), float32(mat[0][1]), float32(mat[0][2]),
		float32(mat[1][0]), float32(mat[1][1]), float32(mat[1][2]),
		float32(mat[2][0]), float32(mat[2][1]), float32(mat[2][2]))
}

// Gl OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) GL() *mgl32.Vec3 {
	return &mgl32.Vec3{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ())}
}

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) Bullet() mbt.BtVector3 {
	return mbt.NewBtVector3(float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()))
}
