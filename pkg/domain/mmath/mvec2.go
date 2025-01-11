package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
)

var (
	MVec2Zero = &MVec2{}

	// UnitX holds a vector with X set to one.
	MVec2UnitX = &MVec2{1, 0}
	// UnitY holds a vector with Y set to one.
	MVec2UnitY = &MVec2{0, 1}
	// UnitXY holds a vector with X and Y set to one.
	MVec2UnitXY = &MVec2{1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MVec2MinVal = &MVec2{-math.MaxFloat64, -math.MaxFloat64}
	// MaxVal holds a vector with the highest possible component values.
	MVec2MaxVal = &MVec2{+math.MaxFloat64, +math.MaxFloat64}
)

type MVec2 struct {
	X float64
	Y float64
}

func NewMVec2() *MVec2 {
	return &MVec2{}
}

// String 文字列表現を返します。
func (vec2 *MVec2) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f]", vec2.X, vec2.Y)
}

// Add ベクトルに他のベクトルを加算します
func (vec2 *MVec2) Add(other *MVec2) *MVec2 {
	vec2.X += other.X
	vec2.Y += other.Y
	return vec2
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (vec2 *MVec2) AddScalar(s float64) *MVec2 {
	vec2.X += s
	vec2.Y += s
	return vec2
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (vec2 *MVec2) Added(other *MVec2) *MVec2 {
	return &MVec2{vec2.X + other.X, vec2.Y + other.Y}
}

func (vec2 *MVec2) AddedScalar(s float64) *MVec2 {
	return &MVec2{vec2.X + s, vec2.Y + s}
}

// Sub ベクトルから他のベクトルを減算します
func (vec2 *MVec2) Sub(other *MVec2) *MVec2 {
	vec2.X -= other.X
	vec2.Y -= other.Y
	return vec2
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (vec2 *MVec2) SubScalar(s float64) *MVec2 {
	vec2.X -= s
	vec2.Y -= s
	return vec2
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (vec2 *MVec2) Subed(other *MVec2) *MVec2 {
	return &MVec2{vec2.X - other.X, vec2.Y - other.Y}
}

func (vec2 *MVec2) SubedScalar(s float64) *MVec2 {
	return &MVec2{vec2.X - s, vec2.Y - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (vec2 *MVec2) Mul(other *MVec2) *MVec2 {
	vec2.X *= other.X
	vec2.Y *= other.Y
	return vec2
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (vec2 *MVec2) MulScalar(s float64) *MVec2 {
	vec2.X *= s
	vec2.Y *= s
	return vec2
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (vec2 *MVec2) Muled(other *MVec2) *MVec2 {
	return &MVec2{vec2.X * other.X, vec2.Y * other.Y}
}

func (vec2 *MVec2) MuledScalar(s float64) *MVec2 {
	return &MVec2{vec2.X * s, vec2.Y * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (vec2 *MVec2) Div(other *MVec2) *MVec2 {
	vec2.X /= other.X
	vec2.Y /= other.Y
	return vec2
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (vec2 *MVec2) DivScalar(s float64) *MVec2 {
	vec2.X /= s
	vec2.Y /= s
	return vec2
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (vec2 *MVec2) Dived(other *MVec2) *MVec2 {
	return &MVec2{vec2.X / other.X, vec2.Y / other.Y}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (vec2 *MVec2) DivedScalar(s float64) *MVec2 {
	return &MVec2{vec2.X / s, vec2.Y / s}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (vec2 *MVec2) Equals(other *MVec2) bool {
	return vec2.X == other.X && vec2.Y == other.Y
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (vec2 *MVec2) NotEquals(other MVec2) bool {
	return vec2.X != other.X || vec2.Y != other.Y
}

// NearEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (vec2 *MVec2) NearEquals(other *MVec2, epsilon float64) bool {
	return (math.Abs(vec2.X-other.X) <= epsilon) &&
		(math.Abs(vec2.Y-other.Y) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (vec2 *MVec2) LessThan(other *MVec2) bool {
	return vec2.X < other.X && vec2.Y < other.Y
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (vec2 *MVec2) LessThanOrEquals(other *MVec2) bool {
	return vec2.X <= other.X && vec2.Y <= other.Y
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (vec2 *MVec2) GreaterThan(other *MVec2) bool {
	return vec2.X > other.X && vec2.Y > other.Y
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (vec2 *MVec2) GreaterThanOrEquals(other *MVec2) bool {
	return vec2.X >= other.X && vec2.Y >= other.Y
}

// Negate ベクトルの各要素の符号を反転します (-v)
func (vec2 *MVec2) Negate() *MVec2 {
	vec2.X = -vec2.X
	vec2.Y = -vec2.Y
	return vec2
}

// Negated ベクトルの各要素の符号を反転した結果を返します (-v)
func (vec2 *MVec2) Negated() *MVec2 {
	return &MVec2{-vec2.X, -vec2.Y}
}

// Abs ベクトルの各要素の絶対値を返します
func (vec2 *MVec2) Abs() *MVec2 {
	vec2.X = math.Abs(vec2.X)
	vec2.Y = math.Abs(vec2.Y)
	return vec2
}

// Absed ベクトルの各要素の絶対値を返します
func (vec2 *MVec2) Absed() *MVec2 {
	return &MVec2{math.Abs(vec2.X), math.Abs(vec2.Y)}
}

// Hash ベクトルのハッシュ値を計算します
func (vec2 *MVec2) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f", vec2.X, vec2.Y)))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (vec2 *MVec2) IsZero() bool {
	return vec2.X == 0 && vec2.Y == 0
}

// Length ベクトルの長さを返します
func (vec2 *MVec2) Length() float64 {
	return math.Hypot(vec2.X, vec2.Y)
}

// LengthSqr ベクトルの長さの2乗を返します
func (vec2 *MVec2) LengthSqr() float64 {
	return vec2.X*vec2.X + vec2.Y*vec2.Y
}

// Normalize ベクトルを正規化します
func (vec2 *MVec2) Normalize() *MVec2 {
	sl := vec2.LengthSqr()
	if sl == 0 || sl == 1 {
		return vec2
	}
	return vec2.MulScalar(1 / math.Sqrt(sl))
}

// Normalized ベクトルを正規化した結果を返します
func (vec2 *MVec2) Normalized() *MVec2 {
	vec := vec2.Copy()
	vec.Normalize()
	return vec
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (vec2 *MVec2) Angle(other *MVec2) float64 {
	v := vec2.Dot(other) / (vec2.Length() * other.Length())
	// prevent NaN
	if v > 1. {
		v = v - 2
	} else if v < -1. {
		v = v + 2
	}
	return math.Acos(v)
}

// Degree ベクトルの角度(度数)を返します
func (vec2 *MVec2) Degree(other *MVec2) float64 {
	radian := vec2.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (vec2 *MVec2) Dot(other *MVec2) float64 {
	return vec2.X*other.X + vec2.Y*other.Y
}

// Cross ベクトルの外積を返します
func (vec2 *MVec2) Cross(other *MVec2) *MVec2 {
	return &MVec2{
		vec2.Y*other.X - vec2.X*other.Y,
		vec2.X*other.Y - vec2.Y*other.X,
	}
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (vec2 *MVec2) Min() *MVec2 {
	min := vec2.X
	if vec2.Y < min {
		min = vec2.Y
	}
	return &MVec2{min, min}
}

// Max ベクトルの各要素の最大値を返します
func (vec2 *MVec2) Max() *MVec2 {
	max := vec2.X
	if vec2.Y > max {
		max = vec2.Y
	}
	return &MVec2{max, max}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (vec2 *MVec2) Clamp(min, max *MVec2) *MVec2 {
	vec2.X = Clamped(vec2.X, min.X, max.X)
	vec2.Y = Clamped(vec2.Y, min.Y, max.Y)
	return vec2
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (vec2 *MVec2) Clamped(min, max *MVec2) *MVec2 {
	result := *vec2
	result.Clamp(min, max)
	return &result
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (vec2 *MVec2) Clamp01() *MVec2 {
	return vec2.Clamp(MVec2Zero, MVec2UnitXY)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (vec2 *MVec2) Clamped01() *MVec2 {
	result := vec2.Copy()
	result.Clamp01()
	return result
}

func (vec2 *MVec2) Rotate(angle float64) *MVec2 {
	sinus := math.Sin(angle)
	cosinus := math.Cos(angle)
	vec2.X = vec2.X*cosinus - vec2.Y*sinus
	vec2.Y = vec2.X*sinus + vec2.Y*cosinus
	return vec2
}

// Rotated ベクトルを回転します
func (vec2 *MVec2) Rotated(angle float64) *MVec2 {
	copied := vec2.Copy()
	return copied.Rotate(angle)
}

// RotateAroundPoint ベクトルを指定された点を中心に回転します
func (vec2 *MVec2) RotateAroundPoint(point *MVec2, angle float64) *MVec2 {
	return vec2.Sub(point).Rotate(angle).Add(point)
}

// Copy
func (vec2 *MVec2) Copy() *MVec2 {
	return &MVec2{vec2.X, vec2.Y}
}

// Vector
func (vec2 *MVec2) Vector() []float64 {
	return []float64{vec2.X, vec2.Y}
}

// 線形補間
func (v1 *MVec2) Lerp(v2 *MVec2, t float64) *MVec2 {
	if t <= 0 {
		return v1
	} else if t >= 1 {
		return v2
	}
	return (v2.Sub(v1)).MulScalar(t).Added(v1)
}

func (vec2 *MVec2) Round() *MVec2 {
	return &MVec2{
		math.Round(vec2.X),
		math.Round(vec2.Y),
	}
}

// One 0を1に変える
func (vec2 *MVec2) One() *MVec2 {
	vec := vec2.Vector()
	epsilon := 1e-8
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec2{vec[0], vec[1]}
}

func (vec2 *MVec2) Distance(other *MVec2) float64 {
	s := vec2.Subed(other)
	return s.Length()
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (vec2 *MVec2) ClampIfVerySmall() *MVec2 {
	epsilon := 1e-6
	if math.Abs(vec2.X) < epsilon {
		vec2.X = 0.0
	}
	if math.Abs(vec2.Y) < epsilon {
		vec2.Y = 0.0
	}
	return vec2
}
