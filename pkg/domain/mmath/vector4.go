package mmath

import (
	"fmt"
	"math"
)

// ----- 定数 -----

var (
	VEC4_ZERO = &Vec4{}
	VEC4_ONE  = &Vec4{X: 1, Y: 1, Z: 1, W: 1}
)

// ----- 型定義 -----

// Vec4 は4次元ベクトルを表します
type Vec4 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

// ----- コンストラクタ -----

func NewVec4() *Vec4 {
	return &Vec4{}
}

func NewVec4ByValues(x, y, z, w float64) *Vec4 {
	return &Vec4{X: x, Y: y, Z: z, W: w}
}

// ----- 文字列表現 -----

func (v *Vec4) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", v.X, v.Y, v.Z, v.W)
}

// ----- 算術演算（破壊的） -----

func (v *Vec4) Add(other *Vec4) *Vec4 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	v.W += other.W
	return v
}

func (v *Vec4) MulScalar(s float64) *Vec4 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	v.W *= s
	return v
}

// ----- 算術演算（非破壊的） -----

func (v *Vec4) Added(other *Vec4) *Vec4 {
	return &Vec4{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z, W: v.W + other.W}
}

func (v *Vec4) Subed(other *Vec4) *Vec4 {
	return &Vec4{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z, W: v.W - other.W}
}

func (v *Vec4) MuledScalar(s float64) *Vec4 {
	return &Vec4{X: v.X * s, Y: v.Y * s, Z: v.Z * s, W: v.W * s}
}

// ----- 比較 -----

func (v *Vec4) NearEquals(other *Vec4, epsilon float64) bool {
	return math.Abs(v.X-other.X) <= epsilon &&
		math.Abs(v.Y-other.Y) <= epsilon &&
		math.Abs(v.Z-other.Z) <= epsilon &&
		math.Abs(v.W-other.W) <= epsilon
}

// ----- ベクトル演算 -----

func (v *Vec4) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)
}

func (v *Vec4) Dot(other *Vec4) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

// ----- ユーティリティ -----

func (v *Vec4) Copy() *Vec4 {
	return &Vec4{X: v.X, Y: v.Y, Z: v.Z, W: v.W}
}

func (v *Vec4) Vector() []float64 {
	return []float64{v.X, v.Y, v.Z, v.W}
}

func (v *Vec4) GetXYZ() *Vec3 {
	return NewVec3ByValues(v.X, v.Y, v.Z)
}
