//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func MVec3Bullet(v *mmath.MVec3) bt.BtVector3 {
	return bt.NewBtVector3(float32(-v.X), float32(v.Y), float32(v.Z))
}

// Bullet+OpenGL座標系に変換されたクォータニオンベクトルを返します
func MRotationBullet(v *mmath.MRotation) bt.BtQuaternion {
	r := v.GetRadians()
	return bt.NewBtQuaternion(float32(-r.Y), float32(r.X), float32(-r.Z))
}
