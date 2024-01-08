package mmath

import (
	"fmt"

	"github.com/ungerik/go3d/float64/vec3"

)

type MVec3 vec3.T

// String は MVec3 の文字列表現を返します。
func (v MVec3) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f]", v[0], v[1], v[2])
}

// GL OpenGL座標系に変換された3次元ベクトルを返します
func (v MVec3) GL() MVec3 {
	return MVec3{-v[0], v[1], v[2]}
}

// MMD MMD(MikuMikuDance)座標系に変換された3次元ベクトルを返します
func (v MVec3) MMD() MVec3 {
	return MVec3{v[0], -v[1], -v[2]}
}

// CalcByRatio ベクトルの線形補間を行います
func (v MVec3) CalcByRatio(next MVec3, x, y, z float64) MVec3 {
	return MVec3{v[0] + (next[0]-v[0])*x, v[1] + (next[1]-v[1])*y, v[2] + (next[2]-v[2])*z}
}
