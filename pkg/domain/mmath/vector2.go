package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
)

// ----- 定数 -----

var (
	// VEC2_ZERO はゼロベクトルです
	VEC2_ZERO = &Vec2{}

	// VEC2_UNIT_X はX方向の単位ベクトルです
	VEC2_UNIT_X = &Vec2{X: 1, Y: 0}

	// VEC2_UNIT_Y はY方向の単位ベクトルです
	VEC2_UNIT_Y = &Vec2{X: 0, Y: 1}

	// VEC2_ONE は全要素が1のベクトルです
	VEC2_ONE = &Vec2{X: 1, Y: 1}

	// VEC2_MIN_VAL は最小値ベクトルです
	VEC2_MIN_VAL = &Vec2{X: -math.MaxFloat64, Y: -math.MaxFloat64}

	// VEC2_MAX_VAL は最大値ベクトルです
	VEC2_MAX_VAL = &Vec2{X: math.MaxFloat64, Y: math.MaxFloat64}
)

// ----- 型定義 -----

// Vec2 は2次元ベクトルを表します
type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// ----- コンストラクタ -----

// NewVec2 はゼロベクトルを作成します
func NewVec2() *Vec2 {
	return &Vec2{}
}

// NewVec2ByValues は指定した値でベクトルを作成します
func NewVec2ByValues(x, y float64) *Vec2 {
	return &Vec2{X: x, Y: y}
}

// ----- 文字列表現 -----

// String は文字列表現を返します
func (v *Vec2) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f]", v.X, v.Y)
}

// StringByDigits は指定桁数の文字列表現を返します
func (v *Vec2) StringByDigits(digits int) string {
	format := fmt.Sprintf("[x=%%.%df, y=%%.%df]", digits, digits)
	return fmt.Sprintf(format, v.X, v.Y)
}

// ----- 算術演算（破壊的） -----

// Add は他のベクトルを加算します（破壊的）
func (v *Vec2) Add(other *Vec2) *Vec2 {
	v.X += other.X
	v.Y += other.Y
	return v
}

// AddScalar はスカラーを加算します（破壊的）
func (v *Vec2) AddScalar(s float64) *Vec2 {
	v.X += s
	v.Y += s
	return v
}

// Sub は他のベクトルを減算します（破壊的）
func (v *Vec2) Sub(other *Vec2) *Vec2 {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

// SubScalar はスカラーを減算します（破壊的）
func (v *Vec2) SubScalar(s float64) *Vec2 {
	v.X -= s
	v.Y -= s
	return v
}

// Mul は他のベクトルと要素ごとに乗算します（破壊的）
func (v *Vec2) Mul(other *Vec2) *Vec2 {
	v.X *= other.X
	v.Y *= other.Y
	return v
}

// MulScalar はスカラーを乗算します（破壊的）
func (v *Vec2) MulScalar(s float64) *Vec2 {
	v.X *= s
	v.Y *= s
	return v
}

// Div は他のベクトルと要素ごとに除算します（破壊的）
func (v *Vec2) Div(other *Vec2) *Vec2 {
	v.X /= other.X
	v.Y /= other.Y
	return v
}

// DivScalar はスカラーで除算します（破壊的）
func (v *Vec2) DivScalar(s float64) *Vec2 {
	v.X /= s
	v.Y /= s
	return v
}

// ----- 算術演算（非破壊的） -----

// Added は加算結果を新しいベクトルで返します
func (v *Vec2) Added(other *Vec2) *Vec2 {
	return &Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

// AddedScalar はスカラー加算結果を新しいベクトルで返します
func (v *Vec2) AddedScalar(s float64) *Vec2 {
	return &Vec2{X: v.X + s, Y: v.Y + s}
}

// Subed は減算結果を新しいベクトルで返します
func (v *Vec2) Subed(other *Vec2) *Vec2 {
	return &Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

// SubedScalar はスカラー減算結果を新しいベクトルで返します
func (v *Vec2) SubedScalar(s float64) *Vec2 {
	return &Vec2{X: v.X - s, Y: v.Y - s}
}

// Muled は乗算結果を新しいベクトルで返します
func (v *Vec2) Muled(other *Vec2) *Vec2 {
	return &Vec2{X: v.X * other.X, Y: v.Y * other.Y}
}

// MuledScalar はスカラー乗算結果を新しいベクトルで返します
func (v *Vec2) MuledScalar(s float64) *Vec2 {
	return &Vec2{X: v.X * s, Y: v.Y * s}
}

// Dived は除算結果を新しいベクトルで返します
func (v *Vec2) Dived(other *Vec2) *Vec2 {
	return &Vec2{X: v.X / other.X, Y: v.Y / other.Y}
}

// DivedScalar はスカラー除算結果を新しいベクトルで返します
func (v *Vec2) DivedScalar(s float64) *Vec2 {
	return &Vec2{X: v.X / s, Y: v.Y / s}
}

// ----- 比較 -----

// Equals は他のベクトルと等しいかどうかを返します
func (v *Vec2) Equals(other *Vec2) bool {
	return v.X == other.X && v.Y == other.Y
}

// NotEquals は他のベクトルと等しくないかどうかを返します
func (v *Vec2) NotEquals(other *Vec2) bool {
	return v.X != other.X || v.Y != other.Y
}

// NearEquals は他のベクトルとほぼ等しいかどうかを返します
func (v *Vec2) NearEquals(other *Vec2, epsilon float64) bool {
	return math.Abs(v.X-other.X) <= epsilon && math.Abs(v.Y-other.Y) <= epsilon
}

// LessThan は他のベクトルより小さいかどうかを返します
func (v *Vec2) LessThan(other *Vec2) bool {
	return v.X < other.X && v.Y < other.Y
}

// LessThanOrEquals は他のベクトル以下かどうかを返します
func (v *Vec2) LessThanOrEquals(other *Vec2) bool {
	return v.X <= other.X && v.Y <= other.Y
}

// GreaterThan は他のベクトルより大きいかどうかを返します
func (v *Vec2) GreaterThan(other *Vec2) bool {
	return v.X > other.X && v.Y > other.Y
}

// GreaterThanOrEquals は他のベクトル以上かどうかを返します
func (v *Vec2) GreaterThanOrEquals(other *Vec2) bool {
	return v.X >= other.X && v.Y >= other.Y
}

// ----- 符号・絶対値 -----

// Negate は符号を反転します（破壊的）
func (v *Vec2) Negate() *Vec2 {
	v.X = -v.X
	v.Y = -v.Y
	return v
}

// Negated は符号を反転した新しいベクトルを返します
func (v *Vec2) Negated() *Vec2 {
	return &Vec2{X: -v.X, Y: -v.Y}
}

// Abs は絶対値にします（破壊的）
func (v *Vec2) Abs() *Vec2 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	return v
}

// Absed は絶対値の新しいベクトルを返します
func (v *Vec2) Absed() *Vec2 {
	return &Vec2{X: math.Abs(v.X), Y: math.Abs(v.Y)}
}

// ----- ベクトル演算 -----

// Length はベクトルの長さを返します
func (v *Vec2) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

// LengthSquared はベクトルの長さの2乗を返します
func (v *Vec2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Normalize は正規化します（破壊的）
func (v *Vec2) Normalize() *Vec2 {
	sl := v.LengthSquared()
	if sl == 0 || sl == 1 {
		return v
	}
	return v.MulScalar(1 / math.Sqrt(sl))
}

// Normalized は正規化した新しいベクトルを返します
func (v *Vec2) Normalized() *Vec2 {
	return v.Copy().Normalize()
}

// Dot は内積を返します
func (v *Vec2) Dot(other *Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Cross は外積（2Dでは擬似的な外積）を返します
func (v *Vec2) Cross(other *Vec2) float64 {
	return v.X*other.Y - v.Y*other.X
}

// Distance は他のベクトルとの距離を返します
func (v *Vec2) Distance(other *Vec2) float64 {
	return v.Subed(other).Length()
}

// Angle は他のベクトルとの角度（ラジアン）を返します
func (v *Vec2) Angle(other *Vec2) float64 {
	dot := v.Dot(other) / (v.Length() * other.Length())
	if dot > 1 {
		dot = 1
	} else if dot < -1 {
		dot = -1
	}
	return math.Acos(dot)
}

// Degree は他のベクトルとの角度（度）を返します
func (v *Vec2) Degree(other *Vec2) float64 {
	return RadToDeg(v.Angle(other))
}

// ----- クランプ -----

// Clamp は指定範囲内にクランプします（破壊的）
func (v *Vec2) Clamp(min, max *Vec2) *Vec2 {
	v.X = Clamped(v.X, min.X, max.X)
	v.Y = Clamped(v.Y, min.Y, max.Y)
	return v
}

// Clamped は指定範囲内にクランプした新しいベクトルを返します
func (v *Vec2) Clamped(min, max *Vec2) *Vec2 {
	return v.Copy().Clamp(min, max)
}

// Clamp01 は0～1の範囲内にクランプします（破壊的）
func (v *Vec2) Clamp01() *Vec2 {
	return v.Clamp(VEC2_ZERO, VEC2_ONE)
}

// Clamped01 は0～1の範囲内にクランプした新しいベクトルを返します
func (v *Vec2) Clamped01() *Vec2 {
	return v.Copy().Clamp01()
}

// Truncate は非常に小さい値をゼロにします（破壊的）
func (v *Vec2) Truncate(epsilon float64) *Vec2 {
	if math.Abs(v.X) < epsilon {
		v.X = 0
	}
	if math.Abs(v.Y) < epsilon {
		v.Y = 0
	}
	return v
}

// Truncated は非常に小さい値をゼロにした新しいベクトルを返します
func (v *Vec2) Truncated(epsilon float64) *Vec2 {
	return v.Copy().Truncate(epsilon)
}

// ----- 回転 -----

// Rotate は指定角度（ラジアン）だけ回転します（破壊的）
func (v *Vec2) Rotate(angle float64) *Vec2 {
	sin := math.Sin(angle)
	cos := math.Cos(angle)
	x := v.X*cos - v.Y*sin
	y := v.X*sin + v.Y*cos
	v.X = x
	v.Y = y
	return v
}

// Rotated は指定角度（ラジアン）だけ回転した新しいベクトルを返します
func (v *Vec2) Rotated(angle float64) *Vec2 {
	return v.Copy().Rotate(angle)
}

// RotateAroundPoint は指定点を中心に回転します（破壊的）
func (v *Vec2) RotateAroundPoint(point *Vec2, angle float64) *Vec2 {
	return v.Sub(point).Rotate(angle).Add(point)
}

// ----- 補間 -----

// Lerp は線形補間を行います
func (v *Vec2) Lerp(other *Vec2, t float64) *Vec2 {
	if t <= 0 {
		return v.Copy()
	}
	if t >= 1 {
		return other.Copy()
	}
	if v.Equals(other) {
		return v.Copy()
	}
	return v.Added(other.Subed(v).MuledScalar(t))
}

// ----- ユーティリティ -----

// Copy はコピーを返します
func (v *Vec2) Copy() *Vec2 {
	return &Vec2{X: v.X, Y: v.Y}
}

// Vector はスライス形式で返します
func (v *Vec2) Vector() []float64 {
	return []float64{v.X, v.Y}
}

// Hash はハッシュ値を返します
func (v *Vec2) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f", v.X, v.Y)))
	return h.Sum64()
}

// IsZero はゼロベクトルかどうかを返します
func (v *Vec2) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

// Round は四捨五入した新しいベクトルを返します
func (v *Vec2) Round() *Vec2 {
	return &Vec2{X: math.Round(v.X), Y: math.Round(v.Y)}
}

// One はゼロ要素を1に置き換えた新しいベクトルを返します（比率計算用）
func (v *Vec2) One() *Vec2 {
	epsilon := 1e-8
	x, y := v.X, v.Y
	if math.Abs(x) < epsilon {
		x = 1
	}
	if math.Abs(y) < epsilon {
		y = 1
	}
	return &Vec2{X: x, Y: y}
}

// MinElement は各要素の最小値を持つベクトルを返します
func (v *Vec2) MinElement() float64 {
	if v.X < v.Y {
		return v.X
	}
	return v.Y
}

// MaxElement は各要素の最大値を持つベクトルを返します
func (v *Vec2) MaxElement() float64 {
	if v.X > v.Y {
		return v.X
	}
	return v.Y
}
