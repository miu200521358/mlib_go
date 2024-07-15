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

func (v *MVec4) GetXY() *MVec2 {
	return &MVec2{v.X, v.Y}
}

func (v *MVec4) GetXYZ() *MVec3 {
	return &MVec3{v.X, v.Y, v.Z}
}

// String T の文字列表現を返します。
func (v *MVec4) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", v.X, v.Y, v.Z, v.W)
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *MVec4) MMD() *MVec4 {
	return &MVec4{v.X, v.Y, v.Z, v.W}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec4) Add(other *MVec4) *MVec4 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	v.W += other.W
	return v
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec4) AddScalar(s float64) *MVec4 {
	v.X += s
	v.Y += s
	v.Z += s
	v.W += s
	return v
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec4) Added(other *MVec4) *MVec4 {
	return &MVec4{v.X + other.X, v.Y + other.Y, v.Z + other.Z, v.W + other.W}
}

func (v *MVec4) AddedScalar(s float64) *MVec4 {
	return &MVec4{v.X + s, v.Y + s, v.Z + s, v.W + s}
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec4) Sub(other *MVec4) *MVec4 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	v.W -= other.W
	return v
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec4) SubScalar(s float64) *MVec4 {
	v.X -= s
	v.Y -= s
	v.Z -= s
	v.W -= s
	return v
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec4) Subed(other *MVec4) *MVec4 {
	return &MVec4{v.X - other.X, v.Y - other.Y, v.Z - other.Z, v.W - other.W}
}

func (v *MVec4) SubedScalar(s float64) *MVec4 {
	return &MVec4{v.X - s, v.Y - s, v.Z - s, v.W - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *MVec4) Mul(other *MVec4) *MVec4 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	v.W *= other.W
	return v
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *MVec4) MulScalar(s float64) *MVec4 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	v.W *= s
	return v
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *MVec4) Muled(other *MVec4) *MVec4 {
	return &MVec4{v.X * other.X, v.Y * other.Y, v.Z * other.Z, v.W * other.W}
}

func (v *MVec4) MuledScalar(s float64) *MVec4 {
	return &MVec4{v.X * s, v.Y * s, v.Z * s, v.W * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec4) Div(other *MVec4) *MVec4 {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	v.W /= other.W
	return v
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec4) DivScalar(s float64) *MVec4 {
	v.X /= s
	v.Y /= s
	v.Z /= s
	v.W /= s
	return v
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec4) Dived(other *MVec4) *MVec4 {
	return &MVec4{v.X / other.X, v.Y / other.Y, v.Z / other.Z, v.W / other.W}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (v *MVec4) DivedScalar(s float64) *MVec4 {
	return &MVec4{v.X / s, v.Y / s, v.Z / s, v.W / s}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *MVec4) Equals(other *MVec4) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z && v.W == other.W
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *MVec4) NotEquals(other MVec4) bool {
	return v.X != other.X || v.Y != other.Y || v.Z != other.Z || v.W != other.W
}

// NearEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *MVec4) NearEquals(other *MVec4, epsilon float64) bool {
	return (math.Abs(v.X-other.X) <= epsilon) &&
		(math.Abs(v.Y-other.Y) <= epsilon) &&
		(math.Abs(v.Z-other.Z) <= epsilon) &&
		(math.Abs(v.W-other.W) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *MVec4) LessThan(other *MVec4) bool {
	return v.X < other.X && v.Y < other.Y && v.Z < other.Z && v.W < other.W
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *MVec4) LessThanOrEquals(other *MVec4) bool {
	return v.X <= other.X && v.Y <= other.Y && v.Z <= other.Z && v.W <= other.W
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *MVec4) GreaterThan(other *MVec4) bool {
	return v.X > other.X && v.Y > other.Y && v.Z > other.Z && v.W > other.W
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *MVec4) GreaterThanOrEquals(other *MVec4) bool {
	return v.X >= other.X && v.Y >= other.Y && v.Z >= other.Z && v.W >= other.W
}

// Inverse ベクトルの各要素の符号を反転します (-v)
func (v *MVec4) Inverse() *MVec4 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	v.W = -v.W
	return v
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *MVec4) Inverted() *MVec4 {
	return &MVec4{-v.X, -v.Y, -v.Z, -v.W}
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec4) Abs() *MVec4 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	v.Z = math.Abs(v.Z)
	v.W = math.Abs(v.W)
	return v
}

// Absed ベクトルの各要素の絶対値を返します
func (v *MVec4) Absed() *MVec4 {
	return &MVec4{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z), math.Abs(v.W)}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec4) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f,%.10f", v.X, v.Y, v.Z, v.W)))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec4) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0 && v.W == 0
}

// Length ベクトルの長さを返します
func (v *MVec4) Length() float64 {
	v3 := v.Vec3DividedByW()
	return v3.Length()
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec4) LengthSqr() float64 {
	v3 := v.Vec3DividedByW()
	return v3.LengthSqr()
}

// Normalize ベクトルを正規化します
func (v *MVec4) Normalize() *MVec4 {
	v3 := v.Vec3DividedByW()
	v3.Normalize()
	v.X = v3.X
	v.Y = v3.Y
	v.Z = v3.Z
	v.W = 1
	return v
}

// Normalized ベクトルを正規化した結果を返します
func (v *MVec4) Normalized() *MVec4 {
	vec := *v
	vec.Normalize()
	return &vec
}

// Dot ベクトルの内積を返します
func (v *MVec4) Dot(other *MVec4) float64 {
	a3 := v.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	return a3.Dot(b3)
}

// Dot4 returns the 4 element vdot product of two vectors.
func Dot4(a, b *MVec4) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

// Cross ベクトルの外積を返します
func (v *MVec4) Cross(other *MVec4) *MVec4 {
	a3 := v.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	c3 := a3.Cross(b3)
	return &MVec4{c3.X, c3.Y, c3.Z, 1}
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *MVec4) Min() *MVec4 {
	min := v.X
	if v.Y < min {
		min = v.Y
	}
	if v.Z < min {
		min = v.Z
	}
	if v.W < min {
		min = v.W
	}
	return &MVec4{min, min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *MVec4) Max() *MVec4 {
	max := v.X
	if v.Y > max {
		max = v.Y
	}
	if v.Z > max {
		max = v.Z
	}
	if v.W > max {
		max = v.W
	}
	return &MVec4{max, max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *MVec4) Interpolate(other *MVec4, t float64) *MVec4 {
	t1 := 1 - t
	return &MVec4{
		v.X*t1 + other.X*t,
		v.Y*t1 + other.Y*t,
		v.Z*t1 + other.Z*t,
		v.W*t1 + other.W*t,
	}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec4) Clamp(min, max *MVec4) *MVec4 {
	v.X = ClampedFloat(v.X, min.X, max.X)
	v.Y = ClampedFloat(v.Y, min.Y, max.Y)
	v.Z = ClampedFloat(v.Z, min.Z, max.Z)
	v.W = ClampedFloat(v.W, min.W, max.W)

	return v
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *MVec4) Clamped(min, max *MVec4) *MVec4 {
	result := *v
	result.Clamp(min, max)
	return &result
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec4) Clamp01() *MVec4 {
	return v.Clamp(MVec4Zero, MVec4One)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec4) Clamped01() *MVec4 {
	result := *v
	result.Clamp01()
	return &result
}

// Copy
func (v *MVec4) Copy() *MVec4 {
	copied := MVec4{v.X, v.Y, v.Z, v.W}
	return &copied
}

// Vector
func (v *MVec4) Vector() []float64 {
	return []float64{v.X, v.Y, v.Z, v.W}
}

// 線形補間
func (v1 *MVec4) Lerp(v2 *MVec4, t float64) *MVec4 {
	return (v2.Subed(v1)).MulScalar(t).Add(v1)
}

func (v *MVec4) Round() *MVec4 {
	return &MVec4{
		math.Round(v.X),
		math.Round(v.Y),
		math.Round(v.Z),
		math.Round(v.W),
	}
}

// Vec3DividedByW returns a vec3.T version of the vector by dividing the first three vector components (XYZ) by the last one (W).
func (vec *MVec4) Vec3DividedByW() *MVec3 {
	oow := 1 / vec.W
	return &MVec3{vec.X * oow, vec.Y * oow, vec.Z * oow}
}

// DividedByW returns a copy of the vector with the first three components (XYZ) divided by the last one (W).
func (vec *MVec4) DividedByW() *MVec4 {
	oow := 1 / vec.W
	return &MVec4{vec.X * oow, vec.Y * oow, vec.Z * oow, 1}
}

// DivideByW divides the first three components (XYZ) by the last one (W).
func (vec *MVec4) DivideByW() *MVec4 {
	oow := 1 / vec.W
	vec.X *= oow
	vec.Y *= oow
	vec.Z *= oow
	vec.W = 1
	return vec
}

// 標準偏差を加味したmean処理
func StdMeanVec4(values []MVec4, err float64) *MVec4 {
	npStandardVectors := make([][]float64, len(values))
	npStandardLengths := make([]float64, len(values))

	for i, v := range values {
		npStandardVectors[i] = v.Vector()
		npStandardLengths[i] = v.Length()
	}

	medianStandardValues := Median(npStandardLengths)
	stdStandardValues := Std(npStandardLengths)

	// 中央値から標準偏差の一定範囲までの値を取得
	var filteredStandardValues [][]float64
	for i := 0; i < len(npStandardVectors); i++ {
		if npStandardLengths[i] >= medianStandardValues-err*stdStandardValues &&
			npStandardLengths[i] <= medianStandardValues+err*stdStandardValues {
			filteredStandardValues = append(filteredStandardValues, npStandardVectors[i])
		}
	}

	mean := Mean2DVertical(filteredStandardValues)
	return &MVec4{mean[0], mean[1], mean[2], mean[3]}
}

// One 0を1に変える
func (v *MVec4) One() *MVec4 {
	vec := v.Vector()
	epsilon := 1e-14
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec4{vec[0], vec[1], vec[2], vec[3]}
}

func (v *MVec4) Distance(other *MVec4) float64 {
	s := v.Subed(other)
	return s.Length()
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (v *MVec4) ClampIfVerySmall() *MVec4 {
	epsilon := 1e-6
	if math.Abs(v.X) < epsilon {
		v.X = 0
	}
	if math.Abs(v.Y) < epsilon {
		v.Y = 0
	}
	if math.Abs(v.Z) < epsilon {
		v.Z = 0
	}
	if math.Abs(v.W) < epsilon {
		v.W = 0
	}
	return v
}
