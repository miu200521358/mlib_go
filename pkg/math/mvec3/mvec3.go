package mvec3

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

type T vec3.T

var (
	Zero = T{}

	// UnitX holds a vector with X set to one.
	UnitX = T{1, 0, 0}
	// UnitY holds a vector with Y set to one.
	UnitY = T{0, 1, 0}
	// UnitZ holds a vector with Z set to one.
	UnitZ = T{0, 0, 1}
	// UnitXYZ holds a vector with X, Y, Z set to one.
	UnitXYZ = T{1, 1, 1}

	// Red holds the color red.
	Red = T{1, 0, 0}
	// Green holds the color green.
	Green = T{0, 1, 0}
	// Blue holds the color black.
	Blue = T{0, 0, 1}
	// Black holds the color black.
	Black = T{0, 0, 0}
	// White holds the color white.
	White = T{1, 1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MinVal = T{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	// MaxVal holds a vector with the highest possible component values.
	MaxVal = T{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64}
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

// GetZ returns the value of the Z coordinate
func (v *T) GetZ() float64 {
	return v[2]
}

// SetZ sets the value of the Z coordinate
func (v *T) SetZ(z float64) {
	v[2] = z
}

// String T の文字列表現を返します。
func (v *T) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f]", v.GetX(), v.GetY(), v.GetZ())
}

// GL OpenGL座標系に変換された2次元ベクトルを返します
func (v *T) GL() T {
	return T{-v.GetX(), v.GetY(), v.GetZ()}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *T) MMD() T {
	return T{v.GetX(), -v.GetY(), -v.GetZ()}
}

// Add ベクトルに他のベクトルを加算します
func (v *T) Add(other *T) *T {
	return (*T)((*vec3.T).Add((*vec3.T)(v), (*vec3.T)(other)))
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *T) AddScalar(s float64) *T {
	return (*T)((*vec3.T).Add((*vec3.T)(v), &vec3.T{s, s}))
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *T) Added(other *T) T {
	return T((*vec3.T).Added((*vec3.T)(v), (*vec3.T)(other)))
}

// Sub ベクトルから他のベクトルを減算します
func (v *T) Sub(other *T) *T {
	return (*T)((*vec3.T).Sub((*vec3.T)(v), (*vec3.T)(other)))
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *T) SubScalar(s float64) *T {
	return (*T)((*vec3.T).Sub((*vec3.T)(v), &vec3.T{s, s}))
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *T) Subed(other *T) T {
	return T((*vec3.T).Subed((*vec3.T)(v), (*vec3.T)(other)))
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *T) Mul(other *T) *T {
	return (*T)((*vec3.T).Mul((*vec3.T)(v), (*vec3.T)(other)))
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *T) MulScalar(s float64) *T {
	return (*T)((*vec3.T).Mul((*vec3.T)(v), &vec3.T{s, s}))
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *T) Muled(other *T) T {
	return T((*vec3.T).Muled((*vec3.T)(v), (*vec3.T)(other)))
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *T) Div(other *T) *T {
	return &T{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
	}
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *T) DivScalar(s float64) *T {
	return &T{
		v.GetX() / s,
		v.GetY() / s,
		v.GetZ() / s,
	}
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *T) Dived(other *T) T {
	return T{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
	}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *T) Equal(other *T) bool {
	return v.GetX() == other.GetX() && v.GetY() == other.GetY() && v.GetZ() == other.GetZ()
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *T) NotEqual(other T) bool {
	return v.GetX() != other.GetX() || v.GetY() != other.GetY() || v.GetZ() != other.GetZ()
}

// PracticallyEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *T) PracticallyEquals(other *T, epsilon float64) bool {
	return (*vec3.T).PracticallyEquals((*vec3.T)(v), (*vec3.T)(other), epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *T) LessThan(other *T) bool {
	return v.GetX() < other.GetX() && v.GetY() < other.GetY() && v.GetZ() < other.GetZ()
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *T) LessThanOrEqual(other *T) bool {
	return v.GetX() <= other.GetX() && v.GetY() <= other.GetY() && v.GetZ() <= other.GetZ()
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *T) GreaterThan(other *T) bool {
	return v.GetX() > other.GetX() && v.GetY() > other.GetY() && v.GetZ() > other.GetZ()
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *T) GreaterThanOrEqual(other *T) bool {
	return v.GetX() >= other.GetX() && v.GetY() >= other.GetY() && v.GetZ() >= other.GetZ()
}

// Invert ベクトルの各要素の符号を反転します (-v)
func (v *T) Invert() *T {
	return (*T)((*vec3.T).Invert((*vec3.T)(v)))
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *T) Inverted() T {
	return T((*vec3.T).Inverted((*vec3.T)(v)))
}

// Abs ベクトルの各要素の絶対値を返します
func (v *T) Abs() T {
	return T{math.Abs(v.GetX()), math.Abs(v.GetY()), math.Abs(v.GetZ())}
}

// Hash ベクトルのハッシュ値を計算します
func (v *T) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f", v.GetX(), v.GetY(), v.GetZ())))
	return h.Sum64()
}

// Length ベクトルの長さを返します
func (v *T) Length() float64 {
	return (*vec3.T).Length((*vec3.T)(v))
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *T) IsZero() bool {
	return (*vec3.T).IsZero((*vec3.T)(v))
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *T) LengthSqr() float64 {
	return (*vec3.T).LengthSqr((*vec3.T)(v))
}

// Normalize ベクトルを正規化します
func (v *T) Normalize() *T {
	return (*T)((*vec3.T).Normalize((*vec3.T)(v)))
}

// Normalized ベクトルを正規化した結果を返します
func (v *T) Normalized() T {
	return T((*vec3.T).Normalized((*vec3.T)(v)))
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *T) Angle(other *T) float64 {
	return vec3.Angle((*vec3.T)(v), (*vec3.T)(other))
}

// Degree ベクトルの角度(度数)を返します
func (v *T) Degree(other *T) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *T) Dot(other *T) float64 {
	return vec3.Dot((*vec3.T)(v), (*vec3.T)(other))
}

// Cross ベクトルの外積を返します
func (v *T) Cross(other *T) *T {
	result := T(vec3.Cross((*vec3.T)(v), (*vec3.T)(other)))
	return &result
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *T) Min() *T {
	min := v.GetX()
	if v.GetY() < min {
		min = v.GetY()
	}
	if v.GetZ() < min {
		min = v.GetZ()
	}
	return &T{min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *T) Max() *T {
	max := v.GetX()
	if v.GetY() > max {
		max = v.GetY()
	}
	if v.GetZ() > max {
		max = v.GetZ()
	}
	return &T{max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *T) Interpolate(other *T, t float64) T {
	return T(vec3.Interpolate((*vec3.T)(v), (*vec3.T)(other), t))
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *T) Clamp(min, max *T) *T {
	return (*T)((*vec3.T).Clamp((*vec3.T)(v), (*vec3.T)(min), (*vec3.T)(max)))
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *T) Clamped(min, max *T) T {
	return T((*vec3.T).Clamped((*vec3.T)(v), (*vec3.T)(min), (*vec3.T)(max)))
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *T) Clamp01() *T {
	return (*T)((*vec3.T).Clamp01((*vec3.T)(v)))
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *T) Clamped01() T {
	return T((*vec3.T).Clamped01((*vec3.T)(v)))
}

// Copy
func (v *T) Copy() *T {
	copied := *v
	return &copied
}
