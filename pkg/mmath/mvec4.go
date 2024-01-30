package mmath

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ungerik/go3d/float64/vec4"
)

type MVec4 vec4.T

func NewMVec4() *MVec4 {
	return &MVec4{0, 0, 0, 0}
}

var (
	MVec4Zero = MVec4{}

	// UnitXW holds a vector with X and W set to one.
	MVec4UnitXW = MVec4{1, 0, 0, 1}
	// UnitYW holds a vector with Y and W set to one.
	MVec4UnitYW = MVec4{0, 1, 0, 1}
	// UnitZW holds a vector with Z and W set to one.
	MVec4UnitZW = MVec4{0, 0, 1, 1}
	// UnitW holds a vector with W set to one.
	MVec4UnitW = MVec4{0, 0, 0, 1}
	// UnitXYZW holds a vector with X, Y, Z, W set to one.
	MVec4UnitXYZW = MVec4{1, 1, 1, 1}

	// Red holds the color red.
	MVec4Red = MVec4{1, 0, 0, 1}
	// Green holds the color green.
	MVec4Green = MVec4{0, 1, 0, 1}
	// Black holds the color black.
	MVec4Blue = MVec4{0, 0, 1, 1}
	// Black holds the color black.
	MVec4Black = MVec4{0, 0, 0, 1}
	// White holds the color white.
	MVec4White = MVec4{1, 1, 1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MVec4MinVal = MVec4{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64, 1}
	// MaxVal holds a vector with the highest possible component values.
	MVec4MaxVal = MVec4{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64, 1}
)

// String T の文字列表現を返します。
func (v *MVec4) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v.GetX(), v.GetY(), v.GetZ(), v.GetW())
}

// GetX returns the value of the X coordinate
func (v *MVec4) GetX() float64 {
	return v[0]
}

// SetX sets the value of the X coordinate
func (v *MVec4) SetX(x float64) {
	v[0] = x
}

// GetY returns the value of the Y coordinate
func (v *MVec4) GetY() float64 {
	return v[1]
}

// SetY sets the value of the Y coordinate
func (v *MVec4) SetY(y float64) {
	v[1] = y
}

// GetZ returns the value of the Z coordinate
func (v *MVec4) GetZ() float64 {
	return v[2]
}

// SetZ sets the value of the Z coordinate
func (v *MVec4) SetZ(z float64) {
	v[2] = z
}

// GetW returns the value of the W coordinate
func (v *MVec4) GetW() float64 {
	return v[3]
}

// SetW sets the value of the W coordinate
func (v *MVec4) SetW(w float64) {
	v[3] = w
}

// GL OpenGL座標系に変換された4次元ベクトルを返します
func (v *MVec4) GL() [4]float32 {
	return [4]float32{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(-v.GetW())}
}

func (v *MVec4) Mgl() mgl32.Vec4 {
	return mgl32.Vec4{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(-v.GetW())}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *MVec4) MMD() *MVec4 {
	return &MVec4{v.GetX(), -v.GetY(), -v.GetZ(), v.GetW()}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec4) Add(other *MVec4) *MVec4 {
	return (*MVec4)((*vec4.T).Add((*vec4.T)(v), (*vec4.T)(other)))
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec4) AddScalar(s float64) *MVec4 {
	return (*MVec4)((*vec4.T).Add((*vec4.T)(v), &vec4.T{s, s}))
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec4) Added(other *MVec4) MVec4 {
	return MVec4{v.GetX() + other.GetX(), v.GetY() + other.GetY(), v.GetZ() + other.GetZ(), v.GetW() + other.GetW()}
}

func (v *MVec4) AddedScalar(s float64) MVec4 {
	return MVec4{v.GetX() + s, v.GetY() + s, v.GetZ() + s, v.GetW() + s}
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec4) Sub(other *MVec4) *MVec4 {
	return (*MVec4)((*vec4.T).Sub((*vec4.T)(v), (*vec4.T)(other)))
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec4) SubScalar(s float64) *MVec4 {
	return (*MVec4)((*vec4.T).Sub((*vec4.T)(v), &vec4.T{s, s}))
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec4) Subed(other *MVec4) MVec4 {
	return MVec4{v.GetX() - other.GetX(), v.GetY() - other.GetY(), v.GetZ() - other.GetZ(), v.GetW() - other.GetW()}
}

func (v *MVec4) SubedScalar(s float64) MVec4 {
	return MVec4{v.GetX() - s, v.GetY() - s, v.GetZ() - s, v.GetW() - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *MVec4) Mul(other *MVec4) *MVec4 {
	v[0] *= other[0]
	v[1] *= other[1]
	v[2] *= other[2]
	v[3] *= other[3]
	return v
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *MVec4) MulScalar(s float64) *MVec4 {
	v[0] *= s
	v[1] *= s
	v[2] *= s
	v[3] *= s
	return v
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *MVec4) Muled(other *MVec4) MVec4 {
	return MVec4{v.GetX() * other.GetX(), v.GetY() * other.GetY(), v.GetZ() * other.GetZ(), v.GetW() * other.GetW()}
}

func (v *MVec4) MuledScalar(s float64) MVec4 {
	return MVec4{v.GetX() * s, v.GetY() * s, v.GetZ() * s, v.GetW() * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec4) Div(other *MVec4) *MVec4 {
	return &MVec4{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
		v.GetW() / other.GetW(),
	}
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec4) DivScalar(s float64) *MVec4 {
	return &MVec4{
		v.GetX() / s,
		v.GetY() / s,
		v.GetZ() / s,
		v.GetW() / s,
	}
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec4) Dived(other *MVec4) MVec4 {
	return MVec4{
		v.GetX() / other.GetX(),
		v.GetY() / other.GetY(),
		v.GetZ() / other.GetZ(),
		v.GetW() / other.GetW(),
	}
}

func (v *MVec4) DivedScalar(s float64) MVec4 {
	return MVec4{
		v.GetX() / s,
		v.GetY() / s,
		v.GetZ() / s,
		v.GetW() / s,
	}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *MVec4) Equals(other *MVec4) bool {
	return v.GetX() == other.GetX() && v.GetY() == other.GetY() && v.GetZ() == other.GetZ() && v.GetW() == other.GetW()
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *MVec4) NotEquals(other MVec4) bool {
	return v.GetX() != other.GetX() || v.GetY() != other.GetY() || v.GetZ() != other.GetZ() || v.GetW() != other.GetW()
}

// PracticallyEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *MVec4) PracticallyEquals(other *MVec4, epsilon float64) bool {
	return (math.Abs(v[0]-other[0]) <= epsilon) &&
		(math.Abs(v[1]-other[1]) <= epsilon) &&
		(math.Abs(v[2]-other[2]) <= epsilon) &&
		(math.Abs(v[3]-other[3]) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *MVec4) LessThan(other *MVec4) bool {
	return v.GetX() < other.GetX() && v.GetY() < other.GetY() && v.GetZ() < other.GetZ() && v.GetW() < other.GetW()
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *MVec4) LessThanOrEquals(other *MVec4) bool {
	return v.GetX() <= other.GetX() && v.GetY() <= other.GetY() && v.GetZ() <= other.GetZ() && v.GetW() <= other.GetW()
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *MVec4) GreaterThan(other *MVec4) bool {
	return v.GetX() > other.GetX() && v.GetY() > other.GetY() && v.GetZ() > other.GetZ() && v.GetW() > other.GetW()
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *MVec4) GreaterThanOrEquals(other *MVec4) bool {
	return v.GetX() >= other.GetX() && v.GetY() >= other.GetY() && v.GetZ() >= other.GetZ() && v.GetW() >= other.GetW()
}

// Invert ベクトルの各要素の符号を反転します (-v)
func (v *MVec4) Invert() *MVec4 {
	v[0] = -v[0]
	v[1] = -v[1]
	v[2] = -v[2]
	v[3] = -v[3]
	return v
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *MVec4) Inverted() MVec4 {
	return MVec4((*vec4.T).Inverted((*vec4.T)(v)))
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec4) Abs() MVec4 {
	return MVec4{math.Abs(v.GetX()), math.Abs(v.GetY()), math.Abs(v.GetZ()), math.Abs(v.GetW())}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec4) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f,%.10f", v.GetX(), v.GetY(), v.GetZ(), v.GetW())))
	return h.Sum64()
}

// Length ベクトルの長さを返します
func (v *MVec4) Length() float64 {
	return (*vec4.T).Length((*vec4.T)(v))
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec4) IsZero() bool {
	return (*vec4.T).IsZero((*vec4.T)(v))
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec4) LengthSqr() float64 {
	return (*vec4.T).LengthSqr((*vec4.T)(v))
}

// Normalize ベクトルを正規化します
func (v *MVec4) Normalize() *MVec4 {
	return (*MVec4)((*vec4.T).Normalize((*vec4.T)(v)))
}

// Normalized ベクトルを正規化した結果を返します
func (v *MVec4) Normalized() MVec4 {
	return MVec4((*vec4.T).Normalized((*vec4.T)(v)))
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *MVec4) Angle(other *MVec4) float64 {
	return vec4.Angle((*vec4.T)(v.Normalize()), (*vec4.T)(other.Normalize()))
}

// Degree ベクトルの角度(度数)を返します
func (v *MVec4) Degree(other *MVec4) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *MVec4) Dot(other *MVec4) float64 {
	return vec4.Dot((*vec4.T)(v), (*vec4.T)(other))
}

// Cross ベクトルの外積を返します
func (v *MVec4) Cross(other *MVec4) *MVec4 {
	result := MVec4(vec4.Cross((*vec4.T)(v), (*vec4.T)(other)))
	return &result
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *MVec4) Min() *MVec4 {
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
	return &MVec4{min, min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *MVec4) Max() *MVec4 {
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
	return &MVec4{max, max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *MVec4) Interpolate(other *MVec4, t float64) MVec4 {
	return MVec4(vec4.Interpolate((*vec4.T)(v), (*vec4.T)(other), t))
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec4) Clamp(min, max *MVec4) *MVec4 {
	return (*MVec4)((*vec4.T).Clamp((*vec4.T)(v), (*vec4.T)(min), (*vec4.T)(max)))
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *MVec4) Clamped(min, max *MVec4) MVec4 {
	return MVec4((*vec4.T).Clamped((*vec4.T)(v), (*vec4.T)(min), (*vec4.T)(max)))
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec4) Clamp01() *MVec4 {
	return (*MVec4)((*vec4.T).Clamp01((*vec4.T)(v)))
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec4) Clamped01() MVec4 {
	return MVec4((*vec4.T).Clamped01((*vec4.T)(v)))
}

// DivByW ベクトルの各要素をWで除算します
func (v *MVec4) DivByW() *MVec4 {
	return (*MVec4)((*vec4.T).DivideByW((*vec4.T)(v)))
}

// DivedByW ベクトルの各要素をWで除算した結果を返します
func (v *MVec4) DivedByW() MVec4 {
	return MVec4((*vec4.T).DividedByW((*vec4.T)(v)))
}

// Vec3 ベクトルのX, Y, Zの要素を持つ3次元ベクトルを返します
func (v *MVec4) Vec3() MVec3 {
	return MVec3{v.GetX(), v.GetY(), v.GetZ()}
}

// Vec2 ベクトルのX, Yの要素を持つ2次元ベクトルを返します
func (v *MVec4) Vec2() MVec2 {
	return MVec2{v.GetX(), v.GetY()}
}

// Copy
func (v *MVec4) Copy() *MVec4 {
	return &MVec4{v.GetX(), v.GetY(), v.GetZ(), v.GetW()}
}

// Vector
func (v *MVec4) Vector() *[]float64 {
	return &[]float64{v.GetX(), v.GetY(), v.GetZ(), v.GetW()}
}

func LerpInt(v1, v2, t int) int {
	return v1 + (v2-v1)*t
}

func LerpFloat(v1, v2, t float64) float64 {
	return v1 + (v2-v1)*t
}

// 線形補間
func LerpVec4(v1, v2 *MVec4, t float64) MVec4 {
	return (v2.Sub(v1)).MulScalar(t).Added(v1)
}

func (v *MVec4) Round() *MVec4 {
	return &MVec4{
		math.Round(v.GetX()),
		math.Round(v.GetY()),
		math.Round(v.GetZ()),
		math.Round(v.GetW()),
	}
}
