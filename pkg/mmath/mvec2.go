package mmath

import (
	"fmt"

	"github.com/ungerik/go3d/float64/vec2"

)

type MVec2 vec2.T

// String は MVec2 の文字列表現を返します。
func (v MVec2) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f]", v[0], v[1])
}

// GL OpenGL座標系に変換された2次元ベクトルを返します
func (v MVec2) GL() MVec2 {
	return MVec2{-v[0], v[1]}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v MVec2) MMD() MVec2 {
	return MVec2{v[0], -v[1]}
}

// CalcByRatio ベクトルの線形補間を行います
func (v MVec2) CalcByRatio(next MVec2, x, y float64) MVec2 {
	return MVec2{v[0] + (next[0]-v[0])*x, v[1] + (next[1]-v[1])*y}
}
