//go:build windows
// +build windows

package mbt

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func MVec3Bullet(vec *mmath.MVec3) bt.BtVector3 {
	return bt.NewBtVector3(float32(-vec.X), float32(vec.Y), float32(vec.Z))
}

// Bullet+OpenGL座標系に変換されたクォータニオンベクトルを返します
func MRotationBullet(vec *mmath.MRotation) bt.BtQuaternion {
	r := vec.Radians()
	return bt.NewBtQuaternion(float32(-r.Y), float32(r.X), float32(-r.Z))
}
