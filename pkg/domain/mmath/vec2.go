// 指示: miu200521358
package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
)

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	ZERO_VEC2    = Vec2{}
	UNIT_X_VEC2  = Vec2{X: 1}
	UNIT_Y_VEC2  = Vec2{Y: 1}
	UNIT_XY_VEC2 = Vec2{X: 1, Y: 1}
	VEC2_MIN_VAL = Vec2{X: -math.MaxFloat64, Y: -math.MaxFloat64}
	VEC2_MAX_VAL = Vec2{X: math.MaxFloat64, Y: math.MaxFloat64}
)

// NewVec2 はVec2を生成する。
func NewVec2() Vec2 {
	return Vec2{}
}

// String は文字列表現を返す。
func (v Vec2) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f]", v.X, v.Y)
}

// Add は加算する。
func (v *Vec2) Add(other Vec2) *Vec2 {
	v.X += other.X
	v.Y += other.Y
	return v
}

// AddScalar はスカラーを加算する。
func (v *Vec2) AddScalar(s float64) *Vec2 {
	v.X += s
	v.Y += s
	return v
}

// Added は加算結果を返す。
func (v Vec2) Added(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

// AddedScalar はスカラー加算結果を返す。
func (v Vec2) AddedScalar(s float64) Vec2 {
	return Vec2{v.X + s, v.Y + s}
}

// Sub は減算する。
func (v *Vec2) Sub(other Vec2) *Vec2 {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

// SubScalar はスカラーを減算する。
func (v *Vec2) SubScalar(s float64) *Vec2 {
	v.X -= s
	v.Y -= s
	return v
}

// Subed は減算結果を返す。
func (v Vec2) Subed(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

// SubedScalar はスカラー減算結果を返す。
func (v Vec2) SubedScalar(s float64) Vec2 {
	return Vec2{v.X - s, v.Y - s}
}

// Mul は乗算する。
func (v *Vec2) Mul(other Vec2) *Vec2 {
	v.X *= other.X
	v.Y *= other.Y
	return v
}

// MulScalar はスカラーを乗算する。
func (v *Vec2) MulScalar(s float64) *Vec2 {
	v.X *= s
	v.Y *= s
	return v
}

// Muled は乗算結果を返す。
func (v Vec2) Muled(other Vec2) Vec2 {
	return Vec2{v.X * other.X, v.Y * other.Y}
}

// MuledScalar はスカラー乗算結果を返す。
func (v Vec2) MuledScalar(s float64) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

// Div は除算する。
func (v *Vec2) Div(other Vec2) *Vec2 {
	v.X /= other.X
	v.Y /= other.Y
	return v
}

// DivScalar はスカラーで除算する。
func (v *Vec2) DivScalar(s float64) *Vec2 {
	v.X /= s
	v.Y /= s
	return v
}

// Dived は除算結果を返す。
func (v Vec2) Dived(other Vec2) Vec2 {
	return Vec2{v.X / other.X, v.Y / other.Y}
}

// DivedScalar はスカラー除算結果を返す。
func (v Vec2) DivedScalar(s float64) Vec2 {
	return Vec2{v.X / s, v.Y / s}
}

// Equals は等しいか判定する。
func (v Vec2) Equals(other Vec2) bool {
	return v.X == other.X && v.Y == other.Y
}

// NotEquals は等しくないか判定する。
func (v Vec2) NotEquals(other Vec2) bool {
	return v.X != other.X || v.Y != other.Y
}

// NearEquals は近似的に等しいか判定する。
func (v Vec2) NearEquals(other Vec2, epsilon float64) bool {
	return math.Abs(v.X-other.X) <= epsilon && math.Abs(v.Y-other.Y) <= epsilon
}

// LessThan は小さいか判定する。
func (v Vec2) LessThan(other Vec2) bool {
	return v.X < other.X && v.Y < other.Y
}

// LessThanOrEquals は以下か判定する。
func (v Vec2) LessThanOrEquals(other Vec2) bool {
	return v.X <= other.X && v.Y <= other.Y
}

// GreaterThan は大きいか判定する。
func (v Vec2) GreaterThan(other Vec2) bool {
	return v.X > other.X && v.Y > other.Y
}

// GreaterThanOrEquals は以上か判定する。
func (v Vec2) GreaterThanOrEquals(other Vec2) bool {
	return v.X >= other.X && v.Y >= other.Y
}

// Negate は符号を反転する。
func (v *Vec2) Negate() *Vec2 {
	v.X = -v.X
	v.Y = -v.Y
	return v
}

// Negated は符号反転結果を返す。
func (v Vec2) Negated() Vec2 {
	return Vec2{-v.X, -v.Y}
}

// Abs は絶対値化する。
func (v *Vec2) Abs() *Vec2 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	return v
}

// Absed は絶対値化した結果を返す。
func (v Vec2) Absed() Vec2 {
	return Vec2{math.Abs(v.X), math.Abs(v.Y)}
}

// Hash はハッシュ値を返す。
func (v Vec2) Hash() uint64 {
	h := fnv.New64a()
	_, _ = fmt.Fprintf(h, "%.10f,%.10f", v.X, v.Y)
	return h.Sum64()
}

// IsZero はゼロか判定する。
func (v Vec2) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

// Length は長さを返す。
func (v Vec2) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

// LengthSqr は長さの二乗を返す。
func (v Vec2) LengthSqr() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Normalize は正規化する。
func (v *Vec2) Normalize() *Vec2 {
	sl := v.LengthSqr()
	if sl == 0 || sl == 1 {
		return v
	}
	return v.MulScalar(1 / math.Sqrt(sl))
}

// Normalized は正規化結果を返す。
func (v Vec2) Normalized() Vec2 {
	vec := v
	vec.Normalize()
	return vec
}

// Angle は角度を返す。
func (v Vec2) Angle(other Vec2) float64 {
	denom := v.Length() * other.Length()
	if denom == 0 {
		return 0
	}
	return angleFromCosVec2(v.Dot(other) / denom)
}

// Degree は度数法の角度を返す。
func (v Vec2) Degree(other Vec2) float64 {
	return RadToDeg(v.Angle(other))
}

// Dot は内積を返す。
func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Cross は外積を返す。
func (v Vec2) Cross(other Vec2) Vec2 {
	return Vec2{
		v.Y*other.X - v.X*other.Y,
		v.X*other.Y - v.Y*other.X,
	}
}

// Min は最小値を返す。
func (v Vec2) Min() Vec2 {
	min := v.X
	if v.Y < min {
		min = v.Y
	}
	return Vec2{min, min}
}

// Max は最大値を返す。
func (v Vec2) Max() Vec2 {
	max := v.X
	if v.Y > max {
		max = v.Y
	}
	return Vec2{max, max}
}

// Clamp は範囲内に収める。
func (v *Vec2) Clamp(min, max Vec2) *Vec2 {
	v.X = Clamped(v.X, min.X, max.X)
	v.Y = Clamped(v.Y, min.Y, max.Y)
	return v
}

// Clamped は範囲内に収めた結果を返す。
func (v Vec2) Clamped(min, max Vec2) Vec2 {
	result := v
	result.Clamp(min, max)
	return result
}

// Clamp01 は0〜1に収める。
func (v *Vec2) Clamp01() *Vec2 {
	return v.Clamp(ZERO_VEC2, UNIT_XY_VEC2)
}

// Clamped01 は0〜1に収めた結果を返す。
func (v Vec2) Clamped01() Vec2 {
	result := v
	result.Clamp01()
	return result
}

// Rotate は回転する。
func (v *Vec2) Rotate(angle float64) *Vec2 {
	sinus := math.Sin(angle)
	cosinus := math.Cos(angle)
	v.X = v.X*cosinus - v.Y*sinus
	v.Y = v.X*sinus + v.Y*cosinus
	return v
}

// Rotated は回転結果を返す。
func (v Vec2) Rotated(angle float64) Vec2 {
	result := v
	result.Rotate(angle)
	return result
}

// RotateAroundPoint は指定点の周りで回転する。
func (v *Vec2) RotateAroundPoint(point Vec2, angle float64) *Vec2 {
	return v.Sub(point).Rotate(angle).Add(point)
}

// Copy はコピーを返す。
func (v Vec2) Copy() (Vec2, error) {
	return deepCopy(v)
}

// Vector はスライス表現を返す。
func (v Vec2) Vector() []float64 {
	return []float64{v.X, v.Y}
}

// Lerp は線形補間する。
func (v Vec2) Lerp(other Vec2, t float64) Vec2 {
	if t <= 0 {
		return v
	}
	if t >= 1 {
		return other
	}
	if v.Equals(other) {
		return v
	}
	return other.Subed(v).MuledScalar(t).Added(v)
}

// Round は丸める。
func (v Vec2) Round() Vec2 {
	return Vec2{math.Round(v.X), math.Round(v.Y)}
}

// One は微小値を1に補正した結果を返す。
func (v Vec2) One() Vec2 {
	vec := v.Vector()
	epsilon := 1e-8
	for i := range vec {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return Vec2{vec[0], vec[1]}
}

// Distance は距離を返す。
func (v Vec2) Distance(other Vec2) float64 {
	return v.Subed(other).Length()
}

// ClampIfVerySmall は微小値を0に丸める。
func (v *Vec2) ClampIfVerySmall() *Vec2 {
	epsilon := 1e-6
	if math.Abs(v.X) < epsilon {
		v.X = 0
	}
	if math.Abs(v.Y) < epsilon {
		v.Y = 0
	}
	return v
}

// angleFromCosVec2 はcos値から角度を返す。
func angleFromCosVec2(val float64) float64 {
	if val > 1 {
		val = val - 2
	} else if val < -1 {
		val = val + 2
	}
	return math.Acos(val)
}

