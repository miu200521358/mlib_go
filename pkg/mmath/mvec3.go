package mmath

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

type MVec3 vec3.T

var (
	MVec3Zero = MVec3{}

	// UnitX holds a vector with X set to one.
	MVec3UnitX = MVec3{1, 0, 0}
	// UnitY holds a vector with Y set to one.
	MVec3UnitY = MVec3{0, 1, 0}
	// UnitZ holds a vector with Z set to one.
	MVec3UnitZ = MVec3{0, 0, 1}
	// UnitXYZ holds a vector with X, Y, Z set to one.
	MVec3UnitXYZ = MVec3{1, 1, 1}

	// Red holds the color red.
	MVec3Red = MVec3{1, 0, 0}
	// Green holds the color green.
	MVec3Green = MVec3{0, 1, 0}
	// Blue holds the color black.
	MVec3Blue = MVec3{0, 0, 1}
	// Black holds the color black.
	MVec3Black = MVec3{0, 0, 0}
	// White holds the color white.
	MVec3White = MVec3{1, 1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MVec3MinVal = MVec3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	// MaxVal holds a vector with the highest possible component values.
	MVec3MaxVal = MVec3{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64}
)

// String T の文字列表現を返します。
func (v *MVec3) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f]", v.GetX(), v.GetY(), v.GetZ())
}

// GetX returns the value of the X coordinate
func (v *MVec3) GetX() float64 {
	return v[0]
}

// SetX sets the value of the X coordinate
func (v *MVec3) SetX(x float64) {
	v[0] = x
}

// GetY returns the value of the Y coordinate
func (v *MVec3) GetY() float64 {
	return v[1]
}

// SetY sets the value of the Y coordinate
func (v *MVec3) SetY(y float64) {
	v[1] = y
}

// GetZ returns the value of the Z coordinate
func (v *MVec3) GetZ() float64 {
	return v[2]
}

// SetZ sets the value of the Z coordinate
func (v *MVec3) SetZ(z float64) {
	v[2] = z
}

// Gl OpenGL座標系に変換された2次元ベクトルを返します
func (v *MVec3) GL() [3]float32 {
	return [3]float32{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ())}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *MVec3) MMD() *MVec3 {
	return &MVec3{v.GetX(), -v.GetY(), -v.GetZ()}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec3) Add(other *MVec3) *MVec3 {
	return (*MVec3)((*vec3.T).Add((*vec3.T)(v), (*vec3.T)(other)))
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec3) AddScalar(s float64) *MVec3 {
	return (*MVec3)((*vec3.T).Add((*vec3.T)(v), &vec3.T{s, s}))
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec3) Added(other *MVec3) MVec3 {
	return MVec3((*vec3.T).Added((*vec3.T)(v), (*vec3.T)(other)))
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec3) Sub(other *MVec3) *MVec3 {
	return (*MVec3)((*vec3.T).Sub((*vec3.T)(v), (*vec3.T)(other)))
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec3) SubScalar(s float64) *MVec3 {
	return (*MVec3)((*vec3.T).Sub((*vec3.T)(v), &vec3.T{s, s}))
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec3) Subed(other *MVec3) MVec3 {
	return MVec3((*vec3.T).Subed((*vec3.T)(v), (*vec3.T)(other)))
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *MVec3) Mul(other *MVec3) *MVec3 {
	return (*MVec3)((*vec3.T).Mul((*vec3.T)(v), (*vec3.T)(other)))
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *MVec3) MulScalar(s float64) *MVec3 {
	return (*MVec3)((*vec3.T).Mul((*vec3.T)(v), &vec3.T{s, s}))
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *MVec3) Muled(other *MVec3) MVec3 {
	return MVec3((*vec3.T).Muled((*vec3.T)(v), (*vec3.T)(other)))
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec3) Div(other *MVec3) *MVec3 {
	return &MVec3{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
	}
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec3) DivScalar(s float64) *MVec3 {
	return &MVec3{
		v.GetX() / s,
		v.GetY() / s,
		v.GetZ() / s,
	}
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec3) Dived(other *MVec3) MVec3 {
	return MVec3{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
	}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *MVec3) Equals(other *MVec3) bool {
	return v.GetX() == other.GetX() && v.GetY() == other.GetY() && v.GetZ() == other.GetZ()
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *MVec3) NotEquals(other MVec3) bool {
	return v.GetX() != other.GetX() || v.GetY() != other.GetY() || v.GetZ() != other.GetZ()
}

// PracticallyEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *MVec3) PracticallyEquals(other *MVec3, epsilon float64) bool {
	return (*vec3.T).PracticallyEquals((*vec3.T)(v), (*vec3.T)(other), epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *MVec3) LessThan(other *MVec3) bool {
	return v.GetX() < other.GetX() && v.GetY() < other.GetY() && v.GetZ() < other.GetZ()
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *MVec3) LessThanOrEquals(other *MVec3) bool {
	return v.GetX() <= other.GetX() && v.GetY() <= other.GetY() && v.GetZ() <= other.GetZ()
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *MVec3) GreaterThan(other *MVec3) bool {
	return v.GetX() > other.GetX() && v.GetY() > other.GetY() && v.GetZ() > other.GetZ()
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *MVec3) GreaterThanOrEquals(other *MVec3) bool {
	return v.GetX() >= other.GetX() && v.GetY() >= other.GetY() && v.GetZ() >= other.GetZ()
}

// Invert ベクトルの各要素の符号を反転します (-v)
func (v *MVec3) Invert() *MVec3 {
	return (*MVec3)((*vec3.T).Invert((*vec3.T)(v)))
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *MVec3) Inverted() MVec3 {
	return MVec3((*vec3.T).Inverted((*vec3.T)(v)))
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec3) Abs() MVec3 {
	return MVec3{math.Abs(v.GetX()), math.Abs(v.GetY()), math.Abs(v.GetZ())}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec3) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f", v.GetX(), v.GetY(), v.GetZ())))
	return h.Sum64()
}

// Length ベクトルの長さを返します
func (v *MVec3) Length() float64 {
	return (*vec3.T).Length((*vec3.T)(v))
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec3) IsZero() bool {
	return (*vec3.T).IsZero((*vec3.T)(v))
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec3) LengthSqr() float64 {
	return (*vec3.T).LengthSqr((*vec3.T)(v))
}

// Normalize ベクトルを正規化します
func (v *MVec3) Normalize() *MVec3 {
	return (*MVec3)((*vec3.T).Normalize((*vec3.T)(v)))
}

// Normalized ベクトルを正規化した結果を返します
func (v *MVec3) Normalized() MVec3 {
	return MVec3((*vec3.T).Normalized((*vec3.T)(v)))
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *MVec3) Angle(other *MVec3) float64 {
	return vec3.Angle((*vec3.T)(v), (*vec3.T)(other))
}

// Degree ベクトルの角度(度数)を返します
func (v *MVec3) Degree(other *MVec3) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *MVec3) Dot(other *MVec3) float64 {
	return vec3.Dot((*vec3.T)(v), (*vec3.T)(other))
}

// Cross ベクトルの外積を返します
func (v *MVec3) Cross(other *MVec3) *MVec3 {
	result := MVec3(vec3.Cross((*vec3.T)(v), (*vec3.T)(other)))
	return &result
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *MVec3) Min() *MVec3 {
	min := v.GetX()
	if v.GetY() < min {
		min = v.GetY()
	}
	if v.GetZ() < min {
		min = v.GetZ()
	}
	return &MVec3{min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *MVec3) Max() *MVec3 {
	max := v.GetX()
	if v.GetY() > max {
		max = v.GetY()
	}
	if v.GetZ() > max {
		max = v.GetZ()
	}
	return &MVec3{max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *MVec3) Interpolate(other *MVec3, t float64) MVec3 {
	return MVec3(vec3.Interpolate((*vec3.T)(v), (*vec3.T)(other), t))
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec3) Clamp(min, max *MVec3) *MVec3 {
	return (*MVec3)((*vec3.T).Clamp((*vec3.T)(v), (*vec3.T)(min), (*vec3.T)(max)))
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *MVec3) Clamped(min, max *MVec3) MVec3 {
	return MVec3((*vec3.T).Clamped((*vec3.T)(v), (*vec3.T)(min), (*vec3.T)(max)))
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec3) Clamp01() *MVec3 {
	return (*MVec3)((*vec3.T).Clamp01((*vec3.T)(v)))
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec3) Clamped01() MVec3 {
	return MVec3((*vec3.T).Clamped01((*vec3.T)(v)))
}

// Copy
func (v *MVec3) Copy() *MVec3 {
	copied := *v
	return &copied
}

// Vector
func (v *MVec3) Vector() *[]float64 {
	return &[]float64{v.GetX(), v.GetY(), v.GetZ()}
}
