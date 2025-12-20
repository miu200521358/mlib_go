package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
	"sort"

	"gonum.org/v1/gonum/spatial/r3"
)

// ----- 定数 -----

var (
	// VEC3_ZERO はゼロベクトルです
	VEC3_ZERO = &Vec3{}

	// VEC3_UNIT_X はX方向の単位ベクトルです
	VEC3_UNIT_X = &Vec3{Vec: r3.Vec{X: 1, Y: 0, Z: 0}}

	// VEC3_UNIT_Y はY方向の単位ベクトルです
	VEC3_UNIT_Y = &Vec3{Vec: r3.Vec{X: 0, Y: 1, Z: 0}}

	// VEC3_UNIT_Z はZ方向の単位ベクトルです
	VEC3_UNIT_Z = &Vec3{Vec: r3.Vec{X: 0, Y: 0, Z: 1}}

	// VEC3_ONE は全要素が1のベクトルです
	VEC3_ONE = &Vec3{Vec: r3.Vec{X: 1, Y: 1, Z: 1}}

	// VEC3_UNIT_X_NEG は負のX方向の単位ベクトルです
	VEC3_UNIT_X_NEG = &Vec3{Vec: r3.Vec{X: -1, Y: 0, Z: 0}}

	// VEC3_UNIT_Y_NEG は負のY方向の単位ベクトルです
	VEC3_UNIT_Y_NEG = &Vec3{Vec: r3.Vec{X: 0, Y: -1, Z: 0}}

	// VEC3_UNIT_Z_NEG は負のZ方向の単位ベクトルです
	VEC3_UNIT_Z_NEG = &Vec3{Vec: r3.Vec{X: 0, Y: 0, Z: -1}}

	// VEC3_MIN_VAL は最小値ベクトルです
	VEC3_MIN_VAL = &Vec3{Vec: r3.Vec{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64}}

	// VEC3_MAX_VAL は最大値ベクトルです
	VEC3_MAX_VAL = &Vec3{Vec: r3.Vec{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}}
)

// ----- 型定義 -----

// Vec3 は3次元ベクトルを表します（gonum/r3.Vecを埋め込み）
type Vec3 struct {
	r3.Vec
}

// ----- コンストラクタ -----

// NewVec3 はゼロベクトルを作成します
func NewVec3() *Vec3 {
	return &Vec3{}
}

// NewVec3ByValues は指定した値でベクトルを作成します
func NewVec3ByValues(x, y, z float64) *Vec3 {
	return &Vec3{Vec: r3.Vec{X: x, Y: y, Z: z}}
}

// ----- 文字列表現 -----

// String は文字列表現を返します
func (v *Vec3) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f]", v.X, v.Y, v.Z)
}

// StringByDigits は指定桁数の文字列表現を返します
func (v *Vec3) StringByDigits(digits int) string {
	format := fmt.Sprintf("[x=%%.%df, y=%%.%df, z=%%.%df]", digits, digits, digits)
	return fmt.Sprintf(format, v.X, v.Y, v.Z)
}

// ----- アクセサ -----

// GetXY はXYコンポーネントをVec2として返します
func (v *Vec3) GetXY() *Vec2 {
	return &Vec2{X: v.X, Y: v.Y}
}

// IsOnlyX はXのみ非ゼロかどうかを返します
func (v *Vec3) IsOnlyX() bool {
	return !NearEquals(v.X, 0.0, 1e-10) &&
		NearEquals(v.Y, 0.0, 1e-10) &&
		NearEquals(v.Z, 0.0, 1e-10)
}

// IsOnlyY はYのみ非ゼロかどうかを返します
func (v *Vec3) IsOnlyY() bool {
	return NearEquals(v.X, 0.0, 1e-10) &&
		!NearEquals(v.Y, 0.0, 1e-10) &&
		NearEquals(v.Z, 0.0, 1e-10)
}

// IsOnlyZ はZのみ非ゼロかどうかを返します
func (v *Vec3) IsOnlyZ() bool {
	return NearEquals(v.X, 0.0, 1e-10) &&
		NearEquals(v.Y, 0.0, 1e-10) &&
		!NearEquals(v.Z, 0.0, 1e-10)
}

// ----- 算術演算（破壊的） -----

// Add は他のベクトルを加算します（破壊的）
func (v *Vec3) Add(other *Vec3) *Vec3 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// AddScalar はスカラーを加算します（破壊的）
func (v *Vec3) AddScalar(s float64) *Vec3 {
	v.X += s
	v.Y += s
	v.Z += s
	return v
}

// Sub は他のベクトルを減算します（破壊的）
func (v *Vec3) Sub(other *Vec3) *Vec3 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// SubScalar はスカラーを減算します（破壊的）
func (v *Vec3) SubScalar(s float64) *Vec3 {
	v.X -= s
	v.Y -= s
	v.Z -= s
	return v
}

// Mul は他のベクトルと要素ごとに乗算します（破壊的）
func (v *Vec3) Mul(other *Vec3) *Vec3 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

// MulScalar はスカラーを乗算します（破壊的）
func (v *Vec3) MulScalar(s float64) *Vec3 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

// Div は他のベクトルと要素ごとに除算します（破壊的）
func (v *Vec3) Div(other *Vec3) *Vec3 {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	return v
}

// DivScalar はスカラーで除算します（破壊的）
func (v *Vec3) DivScalar(s float64) *Vec3 {
	v.X /= s
	v.Y /= s
	v.Z /= s
	return v
}

// ----- 算術演算（非破壊的） -----

// Added は加算結果を新しいベクトルで返します
func (v *Vec3) Added(other *Vec3) *Vec3 {
	return &Vec3{Vec: r3.Add(v.Vec, other.Vec)}
}

// AddedScalar はスカラー加算結果を新しいベクトルで返します
func (v *Vec3) AddedScalar(s float64) *Vec3 {
	return &Vec3{Vec: r3.Vec{X: v.X + s, Y: v.Y + s, Z: v.Z + s}}
}

// Subed は減算結果を新しいベクトルで返します
func (v *Vec3) Subed(other *Vec3) *Vec3 {
	return &Vec3{Vec: r3.Sub(v.Vec, other.Vec)}
}

// SubedScalar はスカラー減算結果を新しいベクトルで返します
func (v *Vec3) SubedScalar(s float64) *Vec3 {
	return &Vec3{Vec: r3.Vec{X: v.X - s, Y: v.Y - s, Z: v.Z - s}}
}

// Muled は乗算結果を新しいベクトルで返します
func (v *Vec3) Muled(other *Vec3) *Vec3 {
	return &Vec3{Vec: r3.Vec{X: v.X * other.X, Y: v.Y * other.Y, Z: v.Z * other.Z}}
}

// MuledScalar はスカラー乗算結果を新しいベクトルで返します
func (v *Vec3) MuledScalar(s float64) *Vec3 {
	return &Vec3{Vec: r3.Scale(s, v.Vec)}
}

// Dived は除算結果を新しいベクトルで返します
func (v *Vec3) Dived(other *Vec3) *Vec3 {
	return &Vec3{Vec: r3.Vec{X: v.X / other.X, Y: v.Y / other.Y, Z: v.Z / other.Z}}
}

// DivedScalar はスカラー除算結果を新しいベクトルで返します
func (v *Vec3) DivedScalar(s float64) *Vec3 {
	return &Vec3{Vec: r3.Scale(1/s, v.Vec)}
}

// ----- 比較 -----

// Equals は他のベクトルと等しいかどうかを返します
func (v *Vec3) Equals(other *Vec3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// NotEquals は他のベクトルと等しくないかどうかを返します
func (v *Vec3) NotEquals(other *Vec3) bool {
	return v.X != other.X || v.Y != other.Y || v.Z != other.Z
}

// NearEquals は他のベクトルとほぼ等しいかどうかを返します
func (v *Vec3) NearEquals(other *Vec3, epsilon float64) bool {
	return math.Abs(v.X-other.X) <= epsilon &&
		math.Abs(v.Y-other.Y) <= epsilon &&
		math.Abs(v.Z-other.Z) <= epsilon
}

// LessThan は他のベクトルより小さいかどうかを返します
func (v *Vec3) LessThan(other *Vec3) bool {
	return v.X < other.X && v.Y < other.Y && v.Z < other.Z
}

// LessThanOrEquals は他のベクトル以下かどうかを返します
func (v *Vec3) LessThanOrEquals(other *Vec3) bool {
	return v.X <= other.X && v.Y <= other.Y && v.Z <= other.Z
}

// GreaterThan は他のベクトルより大きいかどうかを返します
func (v *Vec3) GreaterThan(other *Vec3) bool {
	return v.X > other.X && v.Y > other.Y && v.Z > other.Z
}

// GreaterThanOrEquals は他のベクトル以上かどうかを返します
func (v *Vec3) GreaterThanOrEquals(other *Vec3) bool {
	return v.X >= other.X && v.Y >= other.Y && v.Z >= other.Z
}

// ----- 符号・絶対値 -----

// Negate は符号を反転します（破壊的）
func (v *Vec3) Negate() *Vec3 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

// Negated は符号を反転した新しいベクトルを返します
func (v *Vec3) Negated() *Vec3 {
	return &Vec3{Vec: r3.Scale(-1, v.Vec)}
}

// Abs は絶対値にします（破壊的）
func (v *Vec3) Abs() *Vec3 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	v.Z = math.Abs(v.Z)
	return v
}

// Absed は絶対値の新しいベクトルを返します
func (v *Vec3) Absed() *Vec3 {
	return &Vec3{Vec: r3.Vec{X: math.Abs(v.X), Y: math.Abs(v.Y), Z: math.Abs(v.Z)}}
}

// ----- ベクトル演算 -----

// Length はベクトルの長さを返します
func (v *Vec3) Length() float64 {
	return r3.Norm(v.Vec)
}

// LengthSquared はベクトルの長さの2乗を返します
func (v *Vec3) LengthSquared() float64 {
	return r3.Norm2(v.Vec)
}

// Normalize は正規化します（破壊的）
func (v *Vec3) Normalize() *Vec3 {
	sl := v.LengthSquared()
	if sl == 0 || sl == 1 {
		return v
	}
	return v.MulScalar(1 / math.Sqrt(sl))
}

// Normalized は正規化した新しいベクトルを返します
func (v *Vec3) Normalized() *Vec3 {
	n := r3.Unit(v.Vec)
	return &Vec3{Vec: n}
}

// Dot は内積を返します
func (v *Vec3) Dot(other *Vec3) float64 {
	return r3.Dot(v.Vec, other.Vec)
}

// Cross は外積を返します
func (v *Vec3) Cross(other *Vec3) *Vec3 {
	return &Vec3{Vec: r3.Cross(v.Vec, other.Vec)}
}

// Distance は他のベクトルとの距離を返します
func (v *Vec3) Distance(other *Vec3) float64 {
	return v.Subed(other).Length()
}

// Angle は他のベクトルとの角度（ラジアン）を返します
func (v *Vec3) Angle(other *Vec3) float64 {
	dot := v.Dot(other) / (v.Length() * other.Length())
	if dot > 1 {
		return 0
	} else if dot < -1 {
		return math.Pi
	}
	return math.Acos(dot)
}

// Degree は他のベクトルとの角度（度）を返します
func (v *Vec3) Degree(other *Vec3) float64 {
	return RadToDeg(v.Angle(other))
}

// Project は他のベクトルへの射影を返します
func (v *Vec3) Project(other *Vec3) *Vec3 {
	return other.MuledScalar(v.Dot(other) / other.LengthSquared())
}

// Cos は他のベクトルとのコサイン（角度のコサイン）を返します
func (v *Vec3) Cos(other *Vec3) float64 {
	return r3.Cos(v.Vec, other.Vec)
}

// ----- クランプ -----

// Clamp は指定範囲内にクランプします（破壊的）
func (v *Vec3) Clamp(min, max *Vec3) *Vec3 {
	v.X = Clamped(v.X, min.X, max.X)
	v.Y = Clamped(v.Y, min.Y, max.Y)
	v.Z = Clamped(v.Z, min.Z, max.Z)
	return v
}

// Clamped は指定範囲内にクランプした新しいベクトルを返します
func (v *Vec3) Clamped(min, max *Vec3) *Vec3 {
	return v.Copy().Clamp(min, max)
}

// Clamp01 は0～1の範囲内にクランプします（破壊的）
func (v *Vec3) Clamp01() *Vec3 {
	return v.Clamp(VEC3_ZERO, VEC3_ONE)
}

// Clamped01 は0～1の範囲内にクランプした新しいベクトルを返します
func (v *Vec3) Clamped01() *Vec3 {
	return v.Copy().Clamp01()
}

// Truncate は非常に小さい値をゼロにします（破壊的）
func (v *Vec3) Truncate(epsilon float64) *Vec3 {
	if math.Abs(v.X) < epsilon {
		v.X = 0
	}
	if math.Abs(v.Y) < epsilon {
		v.Y = 0
	}
	if math.Abs(v.Z) < epsilon {
		v.Z = 0
	}
	return v
}

// Truncated は非常に小さい値をゼロにした新しいベクトルを返します
func (v *Vec3) Truncated(epsilon float64) *Vec3 {
	return v.Copy().Truncate(epsilon)
}

// ----- 補間 -----

// Lerp は線形補間を行います
func (v *Vec3) Lerp(other *Vec3, t float64) *Vec3 {
	if t <= 0 {
		return v.Copy()
	}
	if t >= 1 {
		return other.Copy()
	}
	if v.NearEquals(other, 1e-8) {
		return v.Copy()
	}
	return v.Added(other.Subed(v).MuledScalar(t))
}

// Slerp は球面線形補間を行います
func (v *Vec3) Slerp(other *Vec3, t float64) *Vec3 {
	if t <= 0 {
		return v.Copy()
	}
	if t >= 1 {
		return other.Copy()
	}
	if v.NearEquals(other, 1e-8) {
		return v.Copy()
	}

	v0 := v.Normalized()
	v1 := other.Normalized()

	dot := v0.Dot(v1)
	if dot > 1.0 {
		dot = 1.0
	} else if dot < -1.0 {
		dot = -1.0
	}

	theta := math.Acos(dot)
	sinTheta := math.Sin(theta)

	if sinTheta < 1e-10 {
		return v.Lerp(other, t)
	}

	s0 := math.Sin((1-t)*theta) / sinTheta
	s1 := math.Sin(t*theta) / sinTheta

	result := v0.MuledScalar(s0).Add(v1.MuledScalar(s1))
	return result.MuledScalar(v.Length())
}

// ----- ユーティリティ -----

// Copy はコピーを返します
func (v *Vec3) Copy() *Vec3 {
	return &Vec3{Vec: v.Vec}
}

// Vector はスライス形式で返します
func (v *Vec3) Vector() []float64 {
	return []float64{v.X, v.Y, v.Z}
}

// Hash はハッシュ値を返します
func (v *Vec3) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f", v.X, v.Y, v.Z)))
	return h.Sum64()
}

// IsZero はゼロベクトルかどうかを返します
func (v *Vec3) IsZero() bool {
	return v == nil || v.NearEquals(VEC3_ZERO, 1e-10)
}

// IsOne は1ベクトルかどうかを返します
func (v *Vec3) IsOne() bool {
	return v.NearEquals(VEC3_ONE, 1e-10)
}

// Round は指定閾値で四捨五入した新しいベクトルを返します
func (v *Vec3) Round(threshold float64) *Vec3 {
	return &Vec3{Vec: r3.Vec{
		X: Round(v.X, threshold),
		Y: Round(v.Y, threshold),
		Z: Round(v.Z, threshold),
	}}
}

// One はゼロ要素を1に置き換えた新しいベクトルを返します（比率計算用）
func (v *Vec3) One() *Vec3 {
	epsilon := 1e-3
	x, y, z := v.X, v.Y, v.Z
	if math.Abs(x) < epsilon {
		x = 1
	}
	if math.Abs(y) < epsilon {
		y = 1
	}
	if math.Abs(z) < epsilon {
		z = 1
	}
	return &Vec3{Vec: r3.Vec{X: x, Y: y, Z: z}}
}

// Effective は有効な値を返します（NaN/Infの場合は0）
func (v *Vec3) Effective() *Vec3 {
	v.X = Effective(v.X)
	v.Y = Effective(v.Y)
	v.Z = Effective(v.Z)
	return v
}

// RadToDeg はラジアンを度に変換した新しいベクトルを返します
func (v *Vec3) RadToDeg() *Vec3 {
	return &Vec3{Vec: r3.Vec{X: RadToDeg(v.X), Y: RadToDeg(v.Y), Z: RadToDeg(v.Z)}}
}

// DegToRad は度をラジアンに変換した新しいベクトルを返します
func (v *Vec3) DegToRad() *Vec3 {
	return &Vec3{Vec: r3.Vec{X: DegToRad(v.X), Y: DegToRad(v.Y), Z: DegToRad(v.Z)}}
}

// IsPointInsideBox は点が直方体内にあるかどうかを判定します
func (v *Vec3) IsPointInsideBox(min, max *Vec3) bool {
	return v.X >= min.X && v.X <= max.X &&
		v.Y >= min.Y && v.Y <= max.Y &&
		v.Z >= min.Z && v.Z <= max.Z
}

// Distances は他のベクトル群との距離を返します
func (v *Vec3) Distances(others []*Vec3) []float64 {
	distances := make([]float64, len(others))
	for i, other := range others {
		distances[i] = v.Distance(other)
	}
	return distances
}

// ----- 集計関数 -----

// MeanVec3 はベクトルの平均を返します
func MeanVec3(vectors []*Vec3) *Vec3 {
	if len(vectors) == 0 {
		return NewVec3()
	}
	sum := NewVec3()
	for _, v := range vectors {
		sum.Add(v)
	}
	return sum.MuledScalar(1.0 / float64(len(vectors)))
}

// MinVec3 は各要素の最小値を持つベクトルを返します
func MinVec3(vectors []*Vec3) *Vec3 {
	if len(vectors) == 0 {
		return NewVec3()
	}
	min := vectors[0].Copy()
	for _, v := range vectors[1:] {
		min.X = math.Min(min.X, v.X)
		min.Y = math.Min(min.Y, v.Y)
		min.Z = math.Min(min.Z, v.Z)
	}
	return min
}

// MaxVec3 は各要素の最大値を持つベクトルを返します
func MaxVec3(vectors []*Vec3) *Vec3 {
	if len(vectors) == 0 {
		return NewVec3()
	}
	max := vectors[0].Copy()
	for _, v := range vectors[1:] {
		max.X = math.Max(max.X, v.X)
		max.Y = math.Max(max.Y, v.Y)
		max.Z = math.Max(max.Z, v.Z)
	}
	return max
}

// MedianVec3 は各要素の中央値を持つベクトルを返します
func MedianVec3(vectors []*Vec3) *Vec3 {
	if len(vectors) == 0 {
		return NewVec3()
	}
	xValues := make([]float64, len(vectors))
	yValues := make([]float64, len(vectors))
	zValues := make([]float64, len(vectors))
	for i, v := range vectors {
		xValues[i] = v.X
		yValues[i] = v.Y
		zValues[i] = v.Z
	}
	sort.Float64s(xValues)
	sort.Float64s(yValues)
	sort.Float64s(zValues)

	return &Vec3{Vec: r3.Vec{
		X: xValues[len(xValues)/2],
		Y: yValues[len(yValues)/2],
		Z: zValues[len(zValues)/2],
	}}
}

// SortVec3 はベクトルをX, Y, Zの順でソートします
func SortVec3(vectors []Vec3) []Vec3 {
	sort.Slice(vectors, func(i, j int) bool {
		if vectors[i].X == vectors[j].X {
			if vectors[i].Y == vectors[j].Y {
				return vectors[i].Z < vectors[j].Z
			}
			return vectors[i].Y < vectors[j].Y
		}
		return vectors[i].X < vectors[j].X
	})
	return vectors
}
