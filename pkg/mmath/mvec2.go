package mmath

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/ungerik/go3d/float64/vec2"

)

type MVec2 vec2.T

var (
	MVec2Zero = MVec2{}

	// UnitX holds a vector with X set to one.
	MVec2UnitX = MVec2{1, 0}
	// UnitY holds a vector with Y set to one.
	MVec2UnitY = MVec2{0, 1}
	// UnitXY holds a vector with X and Y set to one.
	MVec2UnitXY = MVec2{1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MVec2MinVal = MVec2{-math.MaxFloat64, -math.MaxFloat64}
	// MaxVal holds a vector with the highest possible component values.
	MVec2MaxVal = MVec2{+math.MaxFloat64, +math.MaxFloat64}
)

// GetX returns the value of the X coordinate
func (v *MVec2) GetX() float64 {
	return v[0]
}

// SetX sets the value of the X coordinate
func (v *MVec2) SetX(x float64) {
	v[0] = x
}

// GetY returns the value of the Y coordinate
func (v *MVec2) GetY() float64 {
	return v[1]
}

// SetY sets the value of the Y coordinate
func (v *MVec2) SetY(y float64) {
	v[1] = y
}

// String T の文字列表現を返します。
func (v *MVec2) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f]", v.GetX(), v.GetY())
}

// GL OpenGL座標系に変換された2次元ベクトルを返します
func (v *MVec2) GL() MVec2 {
	return MVec2{-v.GetX(), v.GetY()}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *MVec2) MMD() MVec2 {
	return MVec2{v.GetX(), -v.GetY()}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec2) Add(other *MVec2) *MVec2 {
	return (*MVec2)((*vec2.T).Add((*vec2.T)(v), (*vec2.T)(other)))
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec2) AddScalar(s float64) *MVec2 {
	return (*MVec2)((*vec2.T).Add((*vec2.T)(v), &vec2.T{s, s}))
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec2) Added(other *MVec2) MVec2 {
	return MVec2((*vec2.T).Added((*vec2.T)(v), (*vec2.T)(other)))
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec2) Sub(other *MVec2) *MVec2 {
	return (*MVec2)((*vec2.T).Sub((*vec2.T)(v), (*vec2.T)(other)))
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec2) SubScalar(s float64) *MVec2 {
	return (*MVec2)((*vec2.T).Sub((*vec2.T)(v), &vec2.T{s, s}))
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec2) Subed(other *MVec2) MVec2 {
	return MVec2((*vec2.T).Subed((*vec2.T)(v), (*vec2.T)(other)))
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *MVec2) Mul(other *MVec2) *MVec2 {
	return (*MVec2)((*vec2.T).Mul((*vec2.T)(v), (*vec2.T)(other)))
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *MVec2) MulScalar(s float64) *MVec2 {
	return (*MVec2)((*vec2.T).Mul((*vec2.T)(v), &vec2.T{s, s}))
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *MVec2) Muled(other *MVec2) MVec2 {
	return MVec2((*vec2.T).Muled((*vec2.T)(v), (*vec2.T)(other)))
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec2) Div(other *MVec2) *MVec2 {
	return &MVec2{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
	}
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec2) DivScalar(s float64) *MVec2 {
	return &MVec2{
		v.GetX() / s,
		v.GetY() / s,
	}
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec2) Dived(other *MVec2) MVec2 {
	return MVec2{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
	}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *MVec2) Equals(other *MVec2) bool {
	return v.GetX() == other.GetX() && v.GetY() == other.GetY()
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *MVec2) NotEquals(other MVec2) bool {
	return v.GetX() != other.GetX() || v.GetY() != other.GetY()
}

// PracticallyEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *MVec2) PracticallyEquals(other *MVec2, epsilon float64) bool {
	return (*vec2.T).PracticallyEquals((*vec2.T)(v), (*vec2.T)(other), epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *MVec2) LessThan(other *MVec2) bool {
	return v.GetX() < other.GetX() && v.GetY() < other.GetY()
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *MVec2) LessThanOrEquals(other *MVec2) bool {
	return v.GetX() <= other.GetX() && v.GetY() <= other.GetY()
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *MVec2) GreaterThan(other *MVec2) bool {
	return v.GetX() > other.GetX() && v.GetY() > other.GetY()
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *MVec2) GreaterThanOrEquals(other *MVec2) bool {
	return v.GetX() >= other.GetX() && v.GetY() >= other.GetY()
}

// Invert ベクトルの各要素の符号を反転します (-v)
func (v *MVec2) Invert() *MVec2 {
	return (*MVec2)((*vec2.T).Invert((*vec2.T)(v)))
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *MVec2) Inverted() MVec2 {
	return MVec2((*vec2.T).Inverted((*vec2.T)(v)))
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec2) Abs() MVec2 {
	return MVec2{math.Abs(v.GetX()), math.Abs(v.GetY())}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec2) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f", v.GetX(), v.GetY())))
	return h.Sum64()
}

// Length ベクトルの長さを返します
func (v *MVec2) Length() float64 {
	return (*vec2.T).Length((*vec2.T)(v))
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec2) IsZero() bool {
	return (*vec2.T).IsZero((*vec2.T)(v))
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec2) LengthSqr() float64 {
	return (*vec2.T).LengthSqr((*vec2.T)(v))
}

// Normalize ベクトルを正規化します
func (v *MVec2) Normalize() *MVec2 {
	return (*MVec2)((*vec2.T).Normalize((*vec2.T)(v)))
}

// Normalized ベクトルを正規化した結果を返します
func (v *MVec2) Normalized() MVec2 {
	return MVec2((*vec2.T).Normalized((*vec2.T)(v)))
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *MVec2) Angle(other *MVec2) float64 {
	return vec2.Angle((*vec2.T)(v), (*vec2.T)(other))
}

// Degree ベクトルの角度(度数)を返します
func (v *MVec2) Degree(other *MVec2) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *MVec2) Dot(other *MVec2) float64 {
	return vec2.Dot((*vec2.T)(v), (*vec2.T)(other))
}

// Cross ベクトルの外積を返します
func (v *MVec2) Cross(other *MVec2) *MVec2 {
	result := MVec2(vec2.Cross((*vec2.T)(v), (*vec2.T)(other)))
	return &result
}

// IsLeftWinding ベクトルが他のベクトルより左回りかどうかをチェックします
func (v *MVec2) IsLeftWinding(other *MVec2) bool {
	return vec2.IsLeftWinding((*vec2.T)(v), (*vec2.T)(other))
}

// IsRightWinding ベクトルが他のベクトルより右回りかどうかをチェックします
func (v *MVec2) IsRightWinding(other *MVec2) bool {
	return vec2.IsRightWinding((*vec2.T)(v), (*vec2.T)(other))
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *MVec2) Min() *MVec2 {
	min := v.GetX()
	if v.GetY() < min {
		min = v.GetY()
	}
	return &MVec2{min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *MVec2) Max() *MVec2 {
	max := v.GetX()
	if v.GetY() > max {
		max = v.GetY()
	}
	return &MVec2{max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *MVec2) Interpolate(other *MVec2, t float64) MVec2 {
	return MVec2(vec2.Interpolate((*vec2.T)(v), (*vec2.T)(other), t))
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec2) Clamp(min, max *MVec2) *MVec2 {
	return (*MVec2)((*vec2.T).Clamp((*vec2.T)(v), (*vec2.T)(min), (*vec2.T)(max)))
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *MVec2) Clamped(min, max *MVec2) MVec2 {
	return MVec2((*vec2.T).Clamped((*vec2.T)(v), (*vec2.T)(min), (*vec2.T)(max)))
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec2) Clamp01() *MVec2 {
	return (*MVec2)((*vec2.T).Clamp01((*vec2.T)(v)))
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec2) Clamped01() MVec2 {
	return MVec2((*vec2.T).Clamped01((*vec2.T)(v)))
}

// Rotate ベクトルを回転します
func (v *MVec2) Rotate(angle float64) MVec2 {
	return MVec2((*vec2.T)(v).Rotated(angle))
}

// RotateAroundPoint ベクトルを指定された点を中心に回転します
func (v *MVec2) RotateAroundPoint(point *MVec2, angle float64) *MVec2 {
	return (*MVec2)((*vec2.T)(v).RotateAroundPoint((*vec2.T)(point), angle))
}

// Rotate90DegLeft ベクトルを90度左回転します
func (v *MVec2) Rotate90DegLeft() *MVec2 {
	return (*MVec2)((*vec2.T)(v).Rotate90DegLeft())
}

// Rotate90DegRight ベクトルを90度右回転します
func (v *MVec2) Rotate90DegRight() *MVec2 {
	return (*MVec2)((*vec2.T)(v).Rotate90DegRight())
}

// Copy
func (v *MVec2) Copy() *MVec2 {
	return &MVec2{v.GetX(), v.GetY()}
}

// Vector
func (v *MVec2) Vector() *[]float64 {
	return &[]float64{v.GetX(), v.GetY()}
}
