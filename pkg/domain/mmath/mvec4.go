package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
)

var (
	MVec4Zero = &MVec4{}

	// UnitXW holds a vector with X and W set to one.
	MVec4UnitXW = &MVec4{1, 0, 0, 1}
	// UnitYW holds a vector with Y and W set to one.
	MVec4UnitYW = &MVec4{0, 1, 0, 1}
	// UnitZW holds a vector with Z and W set to one.
	MVec4UnitZW = &MVec4{0, 0, 1, 1}
	// UnitW holds a vector with W set to one.
	MVec4UnitW = &MVec4{0, 0, 0, 1}
	// UnitXYZW holds a vector with X, Y, Z, W set to one.
	MVec4One = &MVec4{1, 1, 1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MVec4MinVal = &MVec4{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64, 1}
	// MaxVal holds a vector with the highest possible component values.
	MVec4MaxVal = &MVec4{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64, 1}
)

type MVec4 struct {
	X float64
	Y float64
	Z float64
	W float64
}

func NewMVec4() *MVec4 {
	return &MVec4{}
}

func (vec4 *MVec4) XY() *MVec2 {
	return &MVec2{vec4.X, vec4.Y}
}

func (vec4 *MVec4) XYZ() *MVec3 {
	return &MVec3{vec4.X, vec4.Y, vec4.Z}
}

// String T の文字列表現を返します。
func (vec4 *MVec4) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", vec4.X, vec4.Y, vec4.Z, vec4.W)
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (vec4 *MVec4) MMD() *MVec4 {
	return &MVec4{vec4.X, vec4.Y, vec4.Z, vec4.W}
}

// Add ベクトルに他のベクトルを加算します
func (vec4 *MVec4) Add(other *MVec4) *MVec4 {
	vec4.X += other.X
	vec4.Y += other.Y
	vec4.Z += other.Z
	vec4.W += other.W
	return vec4
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (vec4 *MVec4) AddScalar(s float64) *MVec4 {
	vec4.X += s
	vec4.Y += s
	vec4.Z += s
	vec4.W += s
	return vec4
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (vec4 *MVec4) Added(other *MVec4) *MVec4 {
	return &MVec4{vec4.X + other.X, vec4.Y + other.Y, vec4.Z + other.Z, vec4.W + other.W}
}

func (vec4 *MVec4) AddedScalar(s float64) *MVec4 {
	return &MVec4{vec4.X + s, vec4.Y + s, vec4.Z + s, vec4.W + s}
}

// Sub ベクトルから他のベクトルを減算します
func (vec4 *MVec4) Sub(other *MVec4) *MVec4 {
	vec4.X -= other.X
	vec4.Y -= other.Y
	vec4.Z -= other.Z
	vec4.W -= other.W
	return vec4
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (vec4 *MVec4) SubScalar(s float64) *MVec4 {
	vec4.X -= s
	vec4.Y -= s
	vec4.Z -= s
	vec4.W -= s
	return vec4
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (vec4 *MVec4) Subed(other *MVec4) *MVec4 {
	return &MVec4{vec4.X - other.X, vec4.Y - other.Y, vec4.Z - other.Z, vec4.W - other.W}
}

func (vec4 *MVec4) SubedScalar(s float64) *MVec4 {
	return &MVec4{vec4.X - s, vec4.Y - s, vec4.Z - s, vec4.W - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (vec4 *MVec4) Mul(other *MVec4) *MVec4 {
	vec4.X *= other.X
	vec4.Y *= other.Y
	vec4.Z *= other.Z
	vec4.W *= other.W
	return vec4
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (vec4 *MVec4) MulScalar(s float64) *MVec4 {
	vec4.X *= s
	vec4.Y *= s
	vec4.Z *= s
	vec4.W *= s
	return vec4
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (vec4 *MVec4) Muled(other *MVec4) *MVec4 {
	return &MVec4{vec4.X * other.X, vec4.Y * other.Y, vec4.Z * other.Z, vec4.W * other.W}
}

func (vec4 *MVec4) MuledScalar(s float64) *MVec4 {
	return &MVec4{vec4.X * s, vec4.Y * s, vec4.Z * s, vec4.W * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (vec4 *MVec4) Div(other *MVec4) *MVec4 {
	vec4.X /= other.X
	vec4.Y /= other.Y
	vec4.Z /= other.Z
	vec4.W /= other.W
	return vec4
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (vec4 *MVec4) DivScalar(s float64) *MVec4 {
	vec4.X /= s
	vec4.Y /= s
	vec4.Z /= s
	vec4.W /= s
	return vec4
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (vec4 *MVec4) Dived(other *MVec4) *MVec4 {
	return &MVec4{vec4.X / other.X, vec4.Y / other.Y, vec4.Z / other.Z, vec4.W / other.W}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (vec4 *MVec4) DivedScalar(s float64) *MVec4 {
	return &MVec4{vec4.X / s, vec4.Y / s, vec4.Z / s, vec4.W / s}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (vec4 *MVec4) Equals(other *MVec4) bool {
	return vec4.X == other.X && vec4.Y == other.Y && vec4.Z == other.Z && vec4.W == other.W
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (vec4 *MVec4) NotEquals(other MVec4) bool {
	return vec4.X != other.X || vec4.Y != other.Y || vec4.Z != other.Z || vec4.W != other.W
}

// NearEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (vec4 *MVec4) NearEquals(other *MVec4, epsilon float64) bool {
	return (math.Abs(vec4.X-other.X) <= epsilon) &&
		(math.Abs(vec4.Y-other.Y) <= epsilon) &&
		(math.Abs(vec4.Z-other.Z) <= epsilon) &&
		(math.Abs(vec4.W-other.W) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (vec4 *MVec4) LessThan(other *MVec4) bool {
	return vec4.X < other.X && vec4.Y < other.Y && vec4.Z < other.Z && vec4.W < other.W
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (vec4 *MVec4) LessThanOrEquals(other *MVec4) bool {
	return vec4.X <= other.X && vec4.Y <= other.Y && vec4.Z <= other.Z && vec4.W <= other.W
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (vec4 *MVec4) GreaterThan(other *MVec4) bool {
	return vec4.X > other.X && vec4.Y > other.Y && vec4.Z > other.Z && vec4.W > other.W
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (vec4 *MVec4) GreaterThanOrEquals(other *MVec4) bool {
	return vec4.X >= other.X && vec4.Y >= other.Y && vec4.Z >= other.Z && vec4.W >= other.W
}

// Inverse ベクトルの各要素の符号を反転します (-v)
func (vec4 *MVec4) Inverse() *MVec4 {
	vec4.X = -vec4.X
	vec4.Y = -vec4.Y
	vec4.Z = -vec4.Z
	vec4.W = -vec4.W
	return vec4
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (vec4 *MVec4) Inverted() *MVec4 {
	return &MVec4{-vec4.X, -vec4.Y, -vec4.Z, -vec4.W}
}

// Abs ベクトルの各要素の絶対値を返します
func (vec4 *MVec4) Abs() *MVec4 {
	vec4.X = math.Abs(vec4.X)
	vec4.Y = math.Abs(vec4.Y)
	vec4.Z = math.Abs(vec4.Z)
	vec4.W = math.Abs(vec4.W)
	return vec4
}

// Absed ベクトルの各要素の絶対値を返します
func (vec4 *MVec4) Absed() *MVec4 {
	return &MVec4{math.Abs(vec4.X), math.Abs(vec4.Y), math.Abs(vec4.Z), math.Abs(vec4.W)}
}

// Hash ベクトルのハッシュ値を計算します
func (vec4 *MVec4) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f,%.10f", vec4.X, vec4.Y, vec4.Z, vec4.W)))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (vec4 *MVec4) IsZero() bool {
	return vec4.X == 0 && vec4.Y == 0 && vec4.Z == 0 && vec4.W == 0
}

// Length ベクトルの長さを返します
func (vec4 *MVec4) Length() float64 {
	v3 := vec4.Vec3DividedByW()
	return v3.Length()
}

// LengthSqr ベクトルの長さの2乗を返します
func (vec4 *MVec4) LengthSqr() float64 {
	v3 := vec4.Vec3DividedByW()
	return v3.LengthSqr()
}

// Normalize ベクトルを正規化します
func (vec4 *MVec4) Normalize() *MVec4 {
	v3 := vec4.Vec3DividedByW()
	v3.Normalize()
	vec4.X = v3.X
	vec4.Y = v3.Y
	vec4.Z = v3.Z
	vec4.W = 1
	return vec4
}

// Normalized ベクトルを正規化した結果を返します
func (vec4 *MVec4) Normalized() *MVec4 {
	vec := *vec4
	vec.Normalize()
	return &vec
}

// Dot ベクトルの内積を返します
func (vec4 *MVec4) Dot(other *MVec4) float64 {
	a3 := vec4.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	return a3.Dot(b3)
}

// Dot4 returns the 4 element vdot product of two vectors.
func Dot4(vec1, vec2 *MVec4) float64 {
	return vec1.X*vec2.X + vec1.Y*vec2.Y + vec1.Z*vec2.Z + vec1.W*vec2.W
}

// Cross ベクトルの外積を返します
func (vec4 *MVec4) Cross(other *MVec4) *MVec4 {
	a3 := vec4.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	c3 := a3.Cross(b3)
	return &MVec4{c3.X, c3.Y, c3.Z, 1}
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (vec4 *MVec4) Min() *MVec4 {
	min := vec4.X
	if vec4.Y < min {
		min = vec4.Y
	}
	if vec4.Z < min {
		min = vec4.Z
	}
	if vec4.W < min {
		min = vec4.W
	}
	return &MVec4{min, min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (vec4 *MVec4) Max() *MVec4 {
	max := vec4.X
	if vec4.Y > max {
		max = vec4.Y
	}
	if vec4.Z > max {
		max = vec4.Z
	}
	if vec4.W > max {
		max = vec4.W
	}
	return &MVec4{max, max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (vec4 *MVec4) Interpolate(other *MVec4, t float64) *MVec4 {
	t1 := 1 - t
	return &MVec4{
		vec4.X*t1 + other.X*t,
		vec4.Y*t1 + other.Y*t,
		vec4.Z*t1 + other.Z*t,
		vec4.W*t1 + other.W*t,
	}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (vec4 *MVec4) Clamp(min, max *MVec4) *MVec4 {
	vec4.X = ClampedFloat(vec4.X, min.X, max.X)
	vec4.Y = ClampedFloat(vec4.Y, min.Y, max.Y)
	vec4.Z = ClampedFloat(vec4.Z, min.Z, max.Z)
	vec4.W = ClampedFloat(vec4.W, min.W, max.W)

	return vec4
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (vec4 *MVec4) Clamped(min, max *MVec4) *MVec4 {
	result := *vec4
	result.Clamp(min, max)
	return &result
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (vec4 *MVec4) Clamp01() *MVec4 {
	return vec4.Clamp(MVec4Zero, MVec4One)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (vec4 *MVec4) Clamped01() *MVec4 {
	result := *vec4
	result.Clamp01()
	return &result
}

// Copy
func (vec4 *MVec4) Copy() *MVec4 {
	copied := MVec4{vec4.X, vec4.Y, vec4.Z, vec4.W}
	return &copied
}

// Vector
func (vec4 *MVec4) Vector() []float64 {
	return []float64{vec4.X, vec4.Y, vec4.Z, vec4.W}
}

// 線形補間
func (vec4 *MVec4) Lerp(other *MVec4, t float64) *MVec4 {
	return (other.Subed(vec4)).MulScalar(t).Add(vec4)
}

// Vec3DividedByW returns a vec3.T version of the vector by dividing the first three vector components (XYZ) by the last one (W).
func (vec4 *MVec4) Vec3DividedByW() *MVec3 {
	oow := 1 / vec4.W
	return &MVec3{vec4.X * oow, vec4.Y * oow, vec4.Z * oow}
}

// DividedByW returns a copy of the vector with the first three components (XYZ) divided by the last one (W).
func (vec4 *MVec4) DividedByW() *MVec4 {
	oow := 1 / vec4.W
	return &MVec4{vec4.X * oow, vec4.Y * oow, vec4.Z * oow, 1}
}

// DivideByW divides the first three components (XYZ) by the last one (W).
func (vec4 *MVec4) DivideByW() *MVec4 {
	oow := 1 / vec4.W
	vec4.X *= oow
	vec4.Y *= oow
	vec4.Z *= oow
	vec4.W = 1
	return vec4
}

// One 0を1に変える
func (vec4 *MVec4) One() *MVec4 {
	vec := vec4.Vector()
	epsilon := 1e-14
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec4{vec[0], vec[1], vec[2], vec[3]}
}

func (vec4 *MVec4) Distance(other *MVec4) float64 {
	s := vec4.Subed(other)
	return s.Length()
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (vec4 *MVec4) ClampIfVerySmall() *MVec4 {
	epsilon := 1e-6
	if math.Abs(vec4.X) < epsilon {
		vec4.X = 0
	}
	if math.Abs(vec4.Y) < epsilon {
		vec4.Y = 0
	}
	if math.Abs(vec4.Z) < epsilon {
		vec4.Z = 0
	}
	if math.Abs(vec4.W) < epsilon {
		vec4.W = 0
	}
	return vec4
}

func (v *MVec4) Round(threshold float64) *MVec4 {
	return &MVec4{
		Round(v.X, threshold),
		Round(v.Y, threshold),
		Round(v.Z, threshold),
		Round(v.W, threshold),
	}
}
