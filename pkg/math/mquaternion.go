package math

import (
	"fmt"

	"github.com/ungerik/go3d/float64/quaternion"
)

type MQuaternion quaternion.T

// String は MQuaternion の文字列表現を返します。
func (v MQuaternion) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v[0], v[1], v[2], v[3])
}

// GL OpenGL座標系に変換された4次元ベクトルを返します
func (v MQuaternion) GL() MQuaternion {
	return MQuaternion{-v[0], v[1], v[2], -v[3]}
}

// MMD MMD(MikuMikuDance)座標系に変換された4次元ベクトルを返します
func (v MQuaternion) MMD() MQuaternion {
	return MQuaternion{v[0], -v[1], -v[2], v[3]}
}
