package mvec2

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/ungerik/go3d/float64/vec2"

)

type T vec2.T

var (
	Zero = T{}

	// UnitX holds a vector with X set to one.
	UnitX = T{1, 0}
	// UnitY holds a vector with Y set to one.
	UnitY = T{0, 1}
	// UnitXY holds a vector with X and Y set to one.
	UnitXY = T{1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MinVal = T{-math.MaxFloat64, -math.MaxFloat64}
	// MaxVal holds a vector with the highest possible component values.
	MaxVal = T{+math.MaxFloat64, +math.MaxFloat64}
)

// GetX returns the value of the X coordinate
func (v *T) GetX() float64 {
	return v[0]
}

// SetX sets the value of the X coordinate
func (v *T) SetX(x float64) {
	v[0] = x
}

// GetY returns the value of the Y coordinate
func (v *T) GetY() float64 {
	return v[1]
}

// SetY sets the value of the Y coordinate
func (v *T) SetY(y float64) {
	v[1] = y
}

// String T の文字列表現を返します。
func (v *T) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f]", v.GetX(), v.GetY())
}

// GL OpenGL座標系に変換された2次元ベクトルを返します
func (v *T) GL() T {
	return T{-v.GetX(), v.GetY()}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *T) MMD() T {
	return T{v.GetX(), -v.GetY()}
}

// Add ベクトルに他のベクトルを加算します
func (v *T) Add(other *T) *T {
	return (*T)((*vec2.T).Add((*vec2.T)(v), (*vec2.T)(other)))
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *T) AddScalar(s float64) *T {
	return (*T)((*vec2.T).Add((*vec2.T)(v), &vec2.T{s, s}))
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *T) Added(other *T) T {
	return T((*vec2.T).Added((*vec2.T)(v), (*vec2.T)(other)))
}

// Sub ベクトルから他のベクトルを減算します
func (v *T) Sub(other *T) *T {
	return (*T)((*vec2.T).Sub((*vec2.T)(v), (*vec2.T)(other)))
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *T) SubScalar(s float64) *T {
	return (*T)((*vec2.T).Sub((*vec2.T)(v), &vec2.T{s, s}))
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *T) Subed(other *T) T {
	return T((*vec2.T).Subed((*vec2.T)(v), (*vec2.T)(other)))
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *T) Mul(other *T) *T {
	return (*T)((*vec2.T).Mul((*vec2.T)(v), (*vec2.T)(other)))
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *T) MulScalar(s float64) *T {
	return (*T)((*vec2.T).Mul((*vec2.T)(v), &vec2.T{s, s}))
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *T) Muled(other *T) T {
	return T((*vec2.T).Muled((*vec2.T)(v), (*vec2.T)(other)))
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *T) Div(other *T) *T {
	return &T{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
	}
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *T) DivScalar(s float64) *T {
	return &T{
		v.GetX() / s,
		v.GetY() / s,
	}
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *T) Dived(other *T) T {
	return T{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
	}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *T) Equals(other *T) bool {
	return v.GetX() == other.GetX() && v.GetY() == other.GetY()
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *T) NotEquals(other T) bool {
	return v.GetX() != other.GetX() || v.GetY() != other.GetY()
}

// PracticallyEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *T) PracticallyEquals(other *T, epsilon float64) bool {
	return (*vec2.T).PracticallyEquals((*vec2.T)(v), (*vec2.T)(other), epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *T) LessThan(other *T) bool {
	return v.GetX() < other.GetX() && v.GetY() < other.GetY()
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *T) LessThanOrEquals(other *T) bool {
	return v.GetX() <= other.GetX() && v.GetY() <= other.GetY()
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *T) GreaterThan(other *T) bool {
	return v.GetX() > other.GetX() && v.GetY() > other.GetY()
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *T) GreaterThanOrEquals(other *T) bool {
	return v.GetX() >= other.GetX() && v.GetY() >= other.GetY()
}

// Invert ベクトルの各要素の符号を反転します (-v)
func (v *T) Invert() *T {
	return (*T)((*vec2.T).Invert((*vec2.T)(v)))
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *T) Inverted() T {
	return T((*vec2.T).Inverted((*vec2.T)(v)))
}

// Abs ベクトルの各要素の絶対値を返します
func (v *T) Abs() T {
	return T{math.Abs(v.GetX()), math.Abs(v.GetY())}
}

// Hash ベクトルのハッシュ値を計算します
func (v *T) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f", v.GetX(), v.GetY())))
	return h.Sum64()
}

// Length ベクトルの長さを返します
func (v *T) Length() float64 {
	return (*vec2.T).Length((*vec2.T)(v))
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *T) IsZero() bool {
	return (*vec2.T).IsZero((*vec2.T)(v))
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *T) LengthSqr() float64 {
	return (*vec2.T).LengthSqr((*vec2.T)(v))
}

// Normalize ベクトルを正規化します
func (v *T) Normalize() *T {
	return (*T)((*vec2.T).Normalize((*vec2.T)(v)))
}

// Normalized ベクトルを正規化した結果を返します
func (v *T) Normalized() T {
	return T((*vec2.T).Normalized((*vec2.T)(v)))
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *T) Angle(other *T) float64 {
	return vec2.Angle((*vec2.T)(v), (*vec2.T)(other))
}

// Degree ベクトルの角度(度数)を返します
func (v *T) Degree(other *T) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *T) Dot(other *T) float64 {
	return vec2.Dot((*vec2.T)(v), (*vec2.T)(other))
}

// Cross ベクトルの外積を返します
func (v *T) Cross(other *T) *T {
	result := T(vec2.Cross((*vec2.T)(v), (*vec2.T)(other)))
	return &result
}

// IsLeftWinding ベクトルが他のベクトルより左回りかどうかをチェックします
func (v *T) IsLeftWinding(other *T) bool {
	return vec2.IsLeftWinding((*vec2.T)(v), (*vec2.T)(other))
}

// IsRightWinding ベクトルが他のベクトルより右回りかどうかをチェックします
func (v *T) IsRightWinding(other *T) bool {
	return vec2.IsRightWinding((*vec2.T)(v), (*vec2.T)(other))
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *T) Min() *T {
	min := v.GetX()
	if v.GetY() < min {
		min = v.GetY()
	}
	return &T{min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *T) Max() *T {
	max := v.GetX()
	if v.GetY() > max {
		max = v.GetY()
	}
	return &T{max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *T) Interpolate(other *T, t float64) T {
	return T(vec2.Interpolate((*vec2.T)(v), (*vec2.T)(other), t))
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *T) Clamp(min, max *T) *T {
	return (*T)((*vec2.T).Clamp((*vec2.T)(v), (*vec2.T)(min), (*vec2.T)(max)))
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *T) Clamped(min, max *T) T {
	return T((*vec2.T).Clamped((*vec2.T)(v), (*vec2.T)(min), (*vec2.T)(max)))
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *T) Clamp01() *T {
	return (*T)((*vec2.T).Clamp01((*vec2.T)(v)))
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *T) Clamped01() T {
	return T((*vec2.T).Clamped01((*vec2.T)(v)))
}

// Rotate ベクトルを回転します
func (v *T) Rotate(angle float64) T {
	return T((*vec2.T)(v).Rotated(angle))
}

// RotateAroundPoint ベクトルを指定された点を中心に回転します
func (v *T) RotateAroundPoint(point *T, angle float64) *T {
	return (*T)((*vec2.T)(v).RotateAroundPoint((*vec2.T)(point), angle))
}

// Rotate90DegLeft ベクトルを90度左回転します
func (v *T) Rotate90DegLeft() *T {
	return (*T)((*vec2.T)(v).Rotate90DegLeft())
}

// Rotate90DegRight ベクトルを90度右回転します
func (v *T) Rotate90DegRight() *T {
	return (*T)((*vec2.T)(v).Rotate90DegRight())
}

// Copy
func (v *T) Copy() *T {
	return &T{v.GetX(), v.GetY()}
}

// Vector
func (v *T) Vector() *[]float64 {
	return &[]float64{v.GetX(), v.GetY()}
}
