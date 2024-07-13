//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infra/bt"
)

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func MVec3Bullet(v *mmath.MVec3) bt.BtVector3 {
	return bt.NewBtVector3(float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()))
}

// Bullet+OpenGL座標系に変換されたクォータニオンベクトルを返します
func MRotationBullet(v *mmath.MRotation) bt.BtQuaternion {
	rx := float32(v.GetRadians().GetX())
	ry := float32(-v.GetRadians().GetY())
	rz := float32(-v.GetRadians().GetZ())
	return bt.NewBtQuaternion(ry, rx, rz)
}
