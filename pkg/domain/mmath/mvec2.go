package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
)

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

type MVec2 struct {
	X float64
	Y float64
}

func NewMVec2() *MVec2 {
	return &MVec2{}
}

// String 文字列表現を返します。
func (v *MVec2) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f]", v.X, v.Y)
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *MVec2) MMD() *MVec2 {
	return &MVec2{v.X, v.Y}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec2) Add(other *MVec2) *MVec2 {
	v.X += other.X
	v.Y += other.Y
	return v
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec2) AddScalar(s float64) *MVec2 {
	v.X += s
	v.Y += s
	return v
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec2) Added(other *MVec2) *MVec2 {
	return &MVec2{v.X + other.X, v.Y + other.Y}
}

func (v *MVec2) AddedScalar(s float64) *MVec2 {
	return &MVec2{v.X + s, v.Y + s}
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec2) Sub(other *MVec2) *MVec2 {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec2) SubScalar(s float64) *MVec2 {
	v.X -= s
	v.Y -= s
	return v
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec2) Subed(other *MVec2) *MVec2 {
	return &MVec2{v.X - other.X, v.Y - other.Y}
}

func (v *MVec2) SubedScalar(s float64) *MVec2 {
	return &MVec2{v.X - s, v.Y - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *MVec2) Mul(other *MVec2) *MVec2 {
	v.X *= other.X
	v.Y *= other.Y
	return v
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *MVec2) MulScalar(s float64) *MVec2 {
	v.X *= s
	v.Y *= s
	return v
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *MVec2) Muled(other *MVec2) *MVec2 {
	return &MVec2{v.X * other.X, v.Y * other.Y}
}

func (v *MVec2) MuledScalar(s float64) *MVec2 {
	return &MVec2{v.X * s, v.Y * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec2) Div(other *MVec2) *MVec2 {
	v.X /= other.X
	v.Y /= other.Y
	return v
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec2) DivScalar(s float64) *MVec2 {
	v.X /= s
	v.Y /= s
	return v
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec2) Dived(other *MVec2) *MVec2 {
	return &MVec2{v.X / other.X, v.Y / other.Y}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (v *MVec2) DivedScalar(s float64) *MVec2 {
	return &MVec2{v.X / s, v.Y / s}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *MVec2) Equals(other *MVec2) bool {
	return v.X == other.X && v.Y == other.Y
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *MVec2) NotEquals(other MVec2) bool {
	return v.X != other.X || v.Y != other.Y
}

// NearEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *MVec2) NearEquals(other *MVec2, epsilon float64) bool {
	return (math.Abs(v.X-other.X) <= epsilon) &&
		(math.Abs(v.Y-other.Y) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *MVec2) LessThan(other *MVec2) bool {
	return v.X < other.X && v.Y < other.Y
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *MVec2) LessThanOrEquals(other *MVec2) bool {
	return v.X <= other.X && v.Y <= other.Y
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *MVec2) GreaterThan(other *MVec2) bool {
	return v.X > other.X && v.Y > other.Y
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *MVec2) GreaterThanOrEquals(other *MVec2) bool {
	return v.X >= other.X && v.Y >= other.Y
}

// Inverse ベクトルの各要素の符号を反転します (-v)
func (v *MVec2) Inverse() *MVec2 {
	v.X = -v.X
	v.Y = -v.Y
	return v
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *MVec2) Inverted() *MVec2 {
	return &MVec2{-v.X, -v.Y}
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec2) Abs() *MVec2 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	return v
}

// Absed ベクトルの各要素の絶対値を返します
func (v *MVec2) Absed() *MVec2 {
	return &MVec2{math.Abs(v.X), math.Abs(v.Y)}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec2) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f", v.X, v.Y)))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec2) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

// Length ベクトルの長さを返します
func (v *MVec2) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec2) LengthSqr() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Normalize ベクトルを正規化します
func (v *MVec2) Normalize() *MVec2 {
	sl := v.LengthSqr()
	if sl == 0 || sl == 1 {
		return v
	}
	return v.MulScalar(1 / math.Sqrt(sl))
}

// Normalized ベクトルを正規化した結果を返します
func (v *MVec2) Normalized() *MVec2 {
	vec := *v
	vec.Normalize()
	return &vec
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (a *MVec2) Angle(b *MVec2) float64 {
	v := a.Dot(b) / (a.Length() * b.Length())
	// prevent NaN
	if v > 1. {
		v = v - 2
	} else if v < -1. {
		v = v + 2
	}
	return math.Acos(v)
}

// Degree ベクトルの角度(度数)を返します
func (v *MVec2) Degree(other *MVec2) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *MVec2) Dot(other *MVec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Cross ベクトルの外積を返します
func (v *MVec2) Cross(other *MVec2) *MVec2 {
	return &MVec2{
		v.Y*other.X - v.X*other.Y,
		v.X*other.Y - v.Y*other.X,
	}
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *MVec2) Min() *MVec2 {
	min := v.X
	if v.Y < min {
		min = v.Y
	}
	return &MVec2{min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *MVec2) Max() *MVec2 {
	max := v.X
	if v.Y > max {
		max = v.Y
	}
	return &MVec2{max, max}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec2) Clamp(min, max *MVec2) *MVec2 {
	v.X = ClampedFloat(v.X, min.X, max.X)
	v.Y = ClampedFloat(v.Y, min.Y, max.Y)
	return v
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *MVec2) Clamped(min, max *MVec2) *MVec2 {
	result := *v
	result.Clamp(min, max)
	return &result
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec2) Clamp01() *MVec2 {
	return v.Clamp(&MVec2Zero, &MVec2UnitXY)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec2) Clamped01() *MVec2 {
	result := *v
	result.Clamp01()
	return &result
}

func (v *MVec2) Rotate(angle float64) *MVec2 {
	sinus := math.Sin(angle)
	cosinus := math.Cos(angle)
	v.X = v.X*cosinus - v.Y*sinus
	v.Y = v.X*sinus + v.Y*cosinus
	return v
}

// Rotated ベクトルを回転します
func (v *MVec2) Rotated(angle float64) *MVec2 {
	copied := v.Copy()
	return copied.Rotate(angle)
}

// RotateAroundPoint ベクトルを指定された点を中心に回転します
func (v *MVec2) RotateAroundPoint(point *MVec2, angle float64) *MVec2 {
	return v.Sub(point).Rotate(angle).Add(point)
}

// Rotate90DegLeft ベクトルを90度左回転します
func (v *MVec2) Rotate90DegLeft() *MVec2 {
	temp := v.X
	v.X = -v.Y
	v.Y = temp
	return v
}

// Rotate90DegRight ベクトルを90度右回転します
func (v *MVec2) Rotate90DegRight() *MVec2 {
	temp := v.X
	v.X = v.Y
	v.Y = -temp
	return v
}

// Copy
func (v *MVec2) Copy() *MVec2 {
	return &MVec2{v.X, v.Y}
}

// Vector
func (v *MVec2) Vector() []float64 {
	return []float64{v.X, v.Y}
}

// 線形補間
func LerpVec2(v1, v2 *MVec2, t float64) *MVec2 {
	return (v2.Sub(v1)).MulScalar(t).Added(v1)
}

func (v *MVec2) Round() *MVec2 {
	return &MVec2{
		math.Round(v.X),
		math.Round(v.Y),
	}
}

// 標準偏差を加味したmean処理
func StdMeanVec2(values []MVec2, err float64) *MVec2 {
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
	return &MVec2{mean[0], mean[1]}
}

// One 0を1に変える
func (v *MVec2) One() *MVec2 {
	vec := v.Vector()
	epsilon := 1e-14
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec2{vec[0], vec[1]}
}

func (v *MVec2) Distance(other *MVec2) float64 {
	s := v.Subed(other)
	return s.Length()
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (v *MVec2) ClampIfVerySmall() *MVec2 {
	epsilon := 1e-6
	if math.Abs(v.X) < epsilon {
		v.X = 0.0
	}
	if math.Abs(v.Y) < epsilon {
		v.Y = 0.0
	}
	return v
}
