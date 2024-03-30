package mmath

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

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

type MVec4 mgl64.Vec4

func NewMVec4() *MVec4 {
	return &MVec4{0, 0, 0, 0}
}

// GetX returns the value of the X coordinate
func (v *MVec4) GetX() float64 {
	return v[0]
}

// SetX sets the value of the X coordinate
func (v *MVec4) SetX(x float64) {
	v[0] = x
}

func (v *MVec4) AddX(x float64) {
	v[0] += x
}

// GetY returns the value of the Y coordinate
func (v *MVec4) GetY() float64 {
	return v[1]
}

// SetY sets the value of the Y coordinate
func (v *MVec4) SetY(y float64) {
	v[1] = y
}

func (v *MVec4) AddY(y float64) {
	v[1] += y
}

// GetZ returns the value of the Z coordinate
func (v *MVec4) GetZ() float64 {
	return v[2]
}

// SetZ sets the value of the Z coordinate
func (v *MVec4) SetZ(z float64) {
	v[2] = z
}

func (v *MVec4) AddZ(z float64) {
	v[2] += z
}

// GetW returns the value of the W coordinate
func (v *MVec4) GetW() float64 {
	return v[3]
}

// SetW sets the value of the W coordinate
func (v *MVec4) SetW(w float64) {
	v[3] = w
}

func (v *MVec4) AddW(w float64) {
	v[3] += w
}

func (v *MVec4) GetXY() *MVec2 {
	return &MVec2{v[0], v[1]}
}

func (v *MVec4) GetXYZ() *MVec3 {
	return &MVec3{v[0], v[1], v[2]}
}

// String T の文字列表現を返します。
func (v *MVec4) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v.GetX(), v.GetY(), v.GetZ(), v.GetW())
}

// GL OpenGL座標系に変換された4次元ベクトルを返します
func (v *MVec4) GL() mgl32.Vec4 {
	return mgl32.Vec4{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()), float32(-v.GetW())}
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *MVec4) MMD() *MVec4 {
	return &MVec4{v.GetX(), -v.GetY(), -v.GetZ(), v.GetW()}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec4) Add(other *MVec4) *MVec4 {
	v[0] += other[0]
	v[1] += other[1]
	v[2] += other[2]
	v[3] += other[3]
	return v
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec4) AddScalar(s float64) *MVec4 {
	v[0] += s
	v[1] += s
	v[2] += s
	v[3] += s
	return v
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec4) Added(other *MVec4) *MVec4 {
	return &MVec4{v[0] + other[0], v[1] + other[1], v[2] + other[2], v[3] + other[3]}
}

func (v *MVec4) AddedScalar(s float64) *MVec4 {
	return &MVec4{v[0] + s, v[1] + s, v[2] + s, v[3] + s}
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec4) Sub(other *MVec4) *MVec4 {
	v[0] -= other[0]
	v[1] -= other[1]
	v[2] -= other[2]
	v[3] -= other[3]
	return v
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec4) SubScalar(s float64) *MVec4 {
	v[0] -= s
	v[1] -= s
	v[2] -= s
	v[3] -= s
	return v
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec4) Subed(other *MVec4) *MVec4 {
	return &MVec4{v[0] - other[0], v[1] - other[1], v[2] - other[2], v[3] - other[3]}
}

func (v *MVec4) SubedScalar(s float64) *MVec4 {
	return &MVec4{v[0] - s, v[1] - s, v[2] - s, v[3] - s}
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
func (v *MVec4) Muled(other *MVec4) *MVec4 {
	return &MVec4{v[0] * other[0], v[1] * other[1], v[2] * other[2], v[3] * other[3]}
}

func (v *MVec4) MuledScalar(s float64) *MVec4 {
	return &MVec4{v[0] * s, v[1] * s, v[2] * s, v[3] * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec4) Div(other *MVec4) *MVec4 {
	v[0] /= other[0]
	v[1] /= other[1]
	v[2] /= other[2]
	v[3] /= other[3]
	return v
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec4) DivScalar(s float64) *MVec4 {
	v[0] /= s
	v[1] /= s
	v[2] /= s
	v[3] /= s
	return v
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec4) Dived(other *MVec4) *MVec4 {
	return &MVec4{v[0] / other[0], v[1] / other[1], v[2] / other[2], v[3] / other[3]}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (v *MVec4) DivedScalar(s float64) *MVec4 {
	return &MVec4{v[0] / s, v[1] / s, v[2] / s, v[3] / s}
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
func (v *MVec4) Inverted() *MVec4 {
	return &MVec4{-v[0], -v[1], -v[2], -v[3]}
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec4) Abs() *MVec4 {
	v[0] = math.Abs(v[0])
	v[1] = math.Abs(v[1])
	v[2] = math.Abs(v[2])
	v[3] = math.Abs(v[3])
	return v
}

// Absed ベクトルの各要素の絶対値を返します
func (v *MVec4) Absed() *MVec4 {
	return &MVec4{math.Abs(v[0]), math.Abs(v[1]), math.Abs(v[2]), math.Abs(v[3])}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec4) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f,%.10f", v.GetX(), v.GetY(), v.GetZ(), v.GetW())))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec4) IsZero() bool {
	return v[0] == 0 && v[1] == 0 && v[2] == 0 && v[3] == 0
}

// Length ベクトルの長さを返します
func (v *MVec4) Length() float64 {
	return v.Vec3DividedByW().Length()
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec4) LengthSqr() float64 {
	return v.Vec3DividedByW().LengthSqr()
}

// Normalize ベクトルを正規化します
func (v *MVec4) Normalize() *MVec4 {
	v3 := v.Vec3DividedByW().Normalize()
	v[0] = v3[0]
	v[1] = v3[1]
	v[2] = v3[2]
	v[3] = 1
	return v
}

// Normalized ベクトルを正規化した結果を返します
func (v *MVec4) Normalized() *MVec4 {
	return v.Copy().Normalize()
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *MVec4) Angle(other *MVec4) float64 {
	vec := v.Dot(other) / (v.Length() * other.Length())
	// prevent NaN
	if vec > 1. {
		vec = vec - 2
	} else if vec < -1. {
		vec = vec + 2
	}
	return math.Acos(vec)
}

// Degree ベクトルの角度(度数)を返します
func (v *MVec4) Degree(other *MVec4) float64 {
	return v.Angle(other) * (180 / math.Pi)
}

// Dot ベクトルの内積を返します
func (v *MVec4) Dot(other *MVec4) float64 {
	return v.Vec3DividedByW().Dot(other.Vec3DividedByW())
}

// Dot4 returns the 4 element vdot product of two vectors.
func Dot4(a, b *MVec4) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
}

// Cross ベクトルの外積を返します
func (v *MVec4) Cross(other *MVec4) *MVec4 {
	c := v.Vec3DividedByW().Cross(other.Vec3DividedByW())
	return &MVec4{c[0], c[1], c[2], 1}
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
func (v *MVec4) Interpolate(other *MVec4, t float64) *MVec4 {
	t1 := 1 - t
	return &MVec4{
		v[0]*t1 + other[0]*t,
		v[1]*t1 + other[1]*t,
		v[2]*t1 + other[2]*t,
		v[3]*t1 + other[3]*t,
	}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec4) Clamp(min, max *MVec4) *MVec4 {
	for i := range v {
		if v[i] < min[i] {
			v[i] = min[i]
		} else if v[i] > max[i] {
			v[i] = max[i]
		}
	}
	return v
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *MVec4) Clamped(min, max *MVec4) *MVec4 {
	return v.Copy().Clamp(min, max)
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec4) Clamp01() *MVec4 {
	return v.Clamp(&MVec4Zero, &MVec4UnitXYZW)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec4) Clamped01() *MVec4 {
	return v.Copy().Clamp01()
}

// Copy
func (v *MVec4) Copy() *MVec4 {
	copied := NewMVec4()
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return copied
}

// Vector
func (v *MVec4) Vector() []float64 {
	return []float64{v.GetX(), v.GetY(), v.GetZ(), v.GetW()}
}

// 線形補間
func LerpVec4(v1, v2 *MVec4, t float64) *MVec4 {
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

// Vec3DividedByW returns a vec3.T version of the vector by dividing the first three vector components (XYZ) by the last one (W).
func (vec *MVec4) Vec3DividedByW() *MVec3 {
	oow := 1 / vec[3]
	return &MVec3{vec[0] * oow, vec[1] * oow, vec[2] * oow}
}

// DividedByW returns a copy of the vector with the first three components (XYZ) divided by the last one (W).
func (vec *MVec4) DividedByW() *MVec4 {
	oow := 1 / vec[3]
	return &MVec4{vec[0] * oow, vec[1] * oow, vec[2] * oow, 1}
}

// DivideByW divides the first three components (XYZ) by the last one (W).
func (vec *MVec4) DivideByW() *MVec4 {
	oow := 1 / vec[3]
	vec[0] *= oow
	vec[1] *= oow
	vec[2] *= oow
	vec[3] = 1
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

	medianStandardValues := mutils.Median(npStandardLengths)
	stdStandardValues := mutils.Std(npStandardLengths)

	// 中央値から標準偏差の一定範囲までの値を取得
	var filteredStandardValues [][]float64
	for i := 0; i < len(npStandardVectors); i++ {
		if npStandardLengths[i] >= medianStandardValues-err*stdStandardValues &&
			npStandardLengths[i] <= medianStandardValues+err*stdStandardValues {
			filteredStandardValues = append(filteredStandardValues, npStandardVectors[i])
		}
	}

	mean := mutils.Mean2DVertical(filteredStandardValues)
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
	return v.Subed(other).Length()
}
