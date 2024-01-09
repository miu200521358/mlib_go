package mvec4

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/ungerik/go3d/float64/vec4"

	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

type T vec4.T

var (
	Zero = T{}

	// UnitXW holds a vector with X and W set to one.
	UnitXW = T{1, 0, 0, 1}
	// UnitYW holds a vector with Y and W set to one.
	UnitYW = T{0, 1, 0, 1}
	// UnitZW holds a vector with Z and W set to one.
	UnitZW = T{0, 0, 1, 1}
	// UnitW holds a vector with W set to one.
	UnitW = T{0, 0, 0, 1}
	// UnitXYZW holds a vector with X, Y, Z, W set to one.
	UnitXYZW = T{1, 1, 1, 1}

	// Red holds the color red.
	Red = T{1, 0, 0, 1}
	// Green holds the color green.
	Green = T{0, 1, 0, 1}
	// Black holds the color black.
	Blue = T{0, 0, 1, 1}
	// Black holds the color black.
	Black = T{0, 0, 0, 1}
	// White holds the color white.
	White = T{1, 1, 1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MinVal = T{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64, 1}
	// MaxVal holds a vector with the highest possible component values.
	MaxVal = T{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64, 1}
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

// GetW returns the value of the W coordinate
func (v *T) GetW() float64 {
	return v[3]
}

// SetW sets the value of the W coordinate
func (v *T) SetW(w float64) {
	v[3] = w
}

// String T の文字列表現を返します。
func (v *T) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v.GetX(), v.GetY(), v.GetZ(), v.GetW())
}

// GL OpenGL座標系に変換された2次元ベクトルを返します
func (v *T) GL() T {
	return T{-v.GetX(), v.GetY(), v.GetZ(), v.GetW()}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *T) MMD() T {
	return T{v.GetX(), -v.GetY(), -v.GetZ(), v.GetW()}
}

// Add ベクトルに他のベクトルを加算します
func (v *T) Add(other *T) *T {
	return (*T)((*vec4.T).Add((*vec4.T)(v), (*vec4.T)(other)))
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *T) AddScalar(s float64) *T {
	return (*T)((*vec4.T).Add((*vec4.T)(v), &vec4.T{s, s}))
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *T) Added(other *T) T {
	return T{v.GetX() + other.GetX(), v.GetY() + other.GetY(), v.GetZ() + other.GetZ(), v.GetW() + other.GetW()}
}

// Sub ベクトルから他のベクトルを減算します
func (v *T) Sub(other *T) *T {
	return (*T)((*vec4.T).Sub((*vec4.T)(v), (*vec4.T)(other)))
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *T) SubScalar(s float64) *T {
	return (*T)((*vec4.T).Sub((*vec4.T)(v), &vec4.T{s, s}))
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *T) Subed(other *T) T {
	return T{v.GetX() - other.GetX(), v.GetY() - other.GetY(), v.GetZ() - other.GetZ(), v.GetW() - other.GetW()}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *T) Mul(other *T) *T {
	v[0] *= other[0]
	v[1] *= other[1]
	v[2] *= other[2]
	v[3] *= other[3]
	return v
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *T) MulScalar(s float64) *T {
	v[0] *= s
	v[1] *= s
	v[2] *= s
	v[3] *= s
	return v
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *T) Muled(other *T) T {
	return T{v.GetX() * other.GetX(), v.GetY() * other.GetY(), v.GetZ() * other.GetZ(), v.GetW() * other.GetW()}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *T) Div(other *T) *T {
	return &T{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
		v.GetW() / other.GetW(),
	}
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *T) DivScalar(s float64) *T {
	return &T{
		v.GetX() / s,
		v.GetY() / s,
		v.GetZ() / s,
		v.GetW() / s,
	}
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *T) Dived(other *T) T {
	return T{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
		v.GetW() / other.GetW(),
	}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *T) Equal(other *T) bool {
	return v.GetX() == other.GetX() && v.GetY() == other.GetY() && v.GetZ() == other.GetZ() && v.GetW() == other.GetW()
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *T) NotEqual(other T) bool {
	return v.GetX() != other.GetX() || v.GetY() != other.GetY() || v.GetZ() != other.GetZ() || v.GetW() != other.GetW()
}

// PracticallyEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *T) PracticallyEquals(other *T, epsilon float64) bool {
	return (math.Abs(v[0]-other[0]) <= epsilon) &&
		(math.Abs(v[1]-other[1]) <= epsilon) &&
		(math.Abs(v[2]-other[2]) <= epsilon) &&
		(math.Abs(v[3]-other[3]) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *T) LessThan(other *T) bool {
	return v.GetX() < other.GetX() && v.GetY() < other.GetY() && v.GetZ() < other.GetZ() && v.GetW() < other.GetW()
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *T) LessThanOrEqual(other *T) bool {
	return v.GetX() <= other.GetX() && v.GetY() <= other.GetY() && v.GetZ() <= other.GetZ() && v.GetW() <= other.GetW()
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *T) GreaterThan(other *T) bool {
	return v.GetX() > other.GetX() && v.GetY() > other.GetY() && v.GetZ() > other.GetZ() && v.GetW() > other.GetW()
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *T) GreaterThanOrEqual(other *T) bool {
	return v.GetX() >= other.GetX() && v.GetY() >= other.GetY() && v.GetZ() >= other.GetZ() && v.GetW() >= other.GetW()
}

// Invert ベクトルの各要素の符号を反転します (-v)
func (v *T) Invert() *T {
	v[0] = -v[0]
	v[1] = -v[1]
	v[2] = -v[2]
	v[3] = -v[3]
	return v
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *T) Inverted() T {
	return T((*vec4.T).Inverted((*vec4.T)(v)))
}

// Abs ベクトルの各要素の絶対値を返します
func (v *T) Abs() T {
	return T{math.Abs(v.GetX()), math.Abs(v.GetY()), math.Abs(v.GetZ()), math.Abs(v.GetW())}
}

// Hash ベクトルのハッシュ値を計算します
func (v *T) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f,%.10f", v.GetX(), v.GetY(), v.GetZ(), v.GetW())))
	return h.Sum64()
}

// Length ベクトルの長さを返します
func (v *T) Length() float64 {
	return (*vec4.T).Length((*vec4.T)(v))
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *T) IsZero() bool {
	return (*vec4.T).IsZero((*vec4.T)(v))
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *T) LengthSqr() float64 {
	return (*vec4.T).LengthSqr((*vec4.T)(v))
}

// Normalize ベクトルを正規化します
func (v *T) Normalize() *T {
	return (*T)((*vec4.T).Normalize((*vec4.T)(v)))
}

// Normalized ベクトルを正規化した結果を返します
func (v *T) Normalized() T {
	return T((*vec4.T).Normalized((*vec4.T)(v)))
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *T) Angle(other *T) float64 {
	return vec4.Angle((*vec4.T)(v.Normalize()), (*vec4.T)(other.Normalize()))
}

// Degree ベクトルの角度(度数)を返します
func (v *T) Degree(other *T) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *T) Dot(other *T) float64 {
	return vec4.Dot((*vec4.T)(v), (*vec4.T)(other))
}

// Cross ベクトルの外積を返します
func (v *T) Cross(other *T) *T {
	result := T(vec4.Cross((*vec4.T)(v), (*vec4.T)(other)))
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
	if v.GetW() < min {
		min = v.GetW()
	}
	return &T{min, min, min, min}
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
	if v.GetW() > max {
		max = v.GetW()
	}
	return &T{max, max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *T) Interpolate(other *T, t float64) T {
	return T(vec4.Interpolate((*vec4.T)(v), (*vec4.T)(other), t))
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *T) Clamp(min, max *T) *T {
	return (*T)((*vec4.T).Clamp((*vec4.T)(v), (*vec4.T)(min), (*vec4.T)(max)))
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *T) Clamped(min, max *T) T {
	return T((*vec4.T).Clamped((*vec4.T)(v), (*vec4.T)(min), (*vec4.T)(max)))
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *T) Clamp01() *T {
	return (*T)((*vec4.T).Clamp01((*vec4.T)(v)))
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *T) Clamped01() T {
	return T((*vec4.T).Clamped01((*vec4.T)(v)))
}

// DivByW ベクトルの各要素をWで除算します
func (v *T) DivByW() *T {
	return (*T)((*vec4.T).DivideByW((*vec4.T)(v)))
}

// DivedByW ベクトルの各要素をWで除算した結果を返します
func (v *T) DivedByW() T {
	return T((*vec4.T).DividedByW((*vec4.T)(v)))
}

// Vec3 ベクトルのX, Y, Zの要素を持つ3次元ベクトルを返します
func (v *T) Vec3() mvec3.T {
	return mvec3.T{v.GetX(), v.GetY(), v.GetZ()}
}

// Vec2 ベクトルのX, Yの要素を持つ2次元ベクトルを返します
func (v *T) Vec2() mvec2.T {
	return mvec2.T{v.GetX(), v.GetY()}
}
