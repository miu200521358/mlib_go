package mmath

import (
	"fmt"

	"github.com/ungerik/go3d/float64/vec4"

)

type MVec4 vec4.T

// String は MVec4 の文字列表現を返します。
func (v MVec4) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v[0], v[1], v[2], v[3])
}

// GL OpenGL座標系に変換された4次元ベクトルを返します
func (v MVec4) GL() MVec4 {
	return MVec4{-v[0], v[1], v[2], v[3]}
}

// MMD MMD(MikuMikuDance)座標系に変換された4次元ベクトルを返します
func (v MVec4) MMD() MVec4 {
	return MVec4{v[0], -v[1], -v[2], v[3]}
}

// CalcByRatio ベクトルの線形補間を行います
func (v MVec4) CalcByRatio(next MVec4, x, y, z, w float64) MVec4 {
	return MVec4{v[0] + (next[0]-v[0])*x, v[1] + (next[1]-v[1])*y, v[2] + (next[2]-v[2])*z, v[3] + (next[3]-v[3])*w}
}
