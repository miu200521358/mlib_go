package mmath

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

var (
	MVec3Zero = MVec3{}

	// UnitX holds a vector with X set to one.
	MVec3UnitX = MVec3{1, 0, 0}
	// UnitY holds a vector with Y set to one.
	MVec3UnitY = MVec3{0, 1, 0}
	// UnitZ holds a vector with Z set to one.
	MVec3UnitZ = MVec3{0, 0, 1}
	// UnitXYZ holds a vector with X, Y, Z set to one.
	MVec3UnitXYZ = MVec3{1, 1, 1}

	// Red holds the color red.
	MVec3Red = MVec3{1, 0, 0}
	// Green holds the color green.
	MVec3Green = MVec3{0, 1, 0}
	// Blue holds the color black.
	MVec3Blue = MVec3{0, 0, 1}
	// Black holds the color black.
	MVec3Black = MVec3{0, 0, 0}
	// White holds the color white.
	MVec3White = MVec3{1, 1, 1}

	// MinVal holds a vector with the smallest possible component values.
	MVec3MinVal = MVec3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	// MaxVal holds a vector with the highest possible component values.
	MVec3MaxVal = MVec3{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64}
)

type MVec3 mgl64.Vec3

func NewMVec3() *MVec3 {
	return &MVec3{0, 0, 0}
}

// GetX returns the value of the X coordinate
func (v *MVec3) GetX() float64 {
	return v[0]
}

// SetX sets the value of the X coordinate
func (v *MVec3) SetX(x float64) {
	v[0] = x
}

// GetY returns the value of the Y coordinate
func (v *MVec3) GetY() float64 {
	return v[1]
}

// SetY sets the value of the Y coordinate
func (v *MVec3) SetY(y float64) {
	v[1] = y
}

// GetZ returns the value of the Z coordinate
func (v *MVec3) GetZ() float64 {
	return v[2]
}

// SetZ sets the value of the Z coordinate
func (v *MVec3) SetZ(z float64) {
	v[2] = z
}

// String T の文字列表現を返します。
func (v *MVec3) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f]", v.GetX(), v.GetY(), v.GetZ())
}

// Gl OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) GL() mgl32.Vec3 {
	return mgl32.Vec3{float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ())}
}

// Bullet+OpenGL座標系に変換された3次元ベクトルを返します
func (v *MVec3) Bullet() mbt.BtVector3 {
	return mbt.NewBtVector3(float32(-v.GetX()), float32(v.GetY()), float32(v.GetZ()))
}

// MMD MMD(MikuMikuDance)座標系に変換された2次元ベクトルを返します
func (v *MVec3) MMD() *MVec3 {
	return &MVec3{v.GetX(), -v.GetY(), -v.GetZ()}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec3) Add(other *MVec3) *MVec3 {
	v[0] += other[0]
	v[1] += other[1]
	v[2] += other[2]
	return v
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec3) AddScalar(s float64) *MVec3 {
	v[0] += s
	v[1] += s
	v[2] += s
	return v
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec3) Added(other *MVec3) *MVec3 {
	return &MVec3{v[0] + other[0], v[1] + other[1], v[2] + other[2]}
}

func (v *MVec3) AddedScalar(s float64) *MVec3 {
	return &MVec3{v[0] + s, v[1] + s, v[2] + s}
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec3) Sub(other *MVec3) *MVec3 {
	v[0] -= other[0]
	v[1] -= other[1]
	v[2] -= other[2]
	return v
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec3) SubScalar(s float64) *MVec3 {
	v[0] -= s
	v[1] -= s
	v[2] -= s
	return v
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec3) Subed(other *MVec3) *MVec3 {
	return &MVec3{v[0] - other[0], v[1] - other[1], v[2] - other[2]}
}

func (v *MVec3) SubedScalar(s float64) *MVec3 {
	return &MVec3{v[0] - s, v[1] - s, v[2] - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *MVec3) Mul(other *MVec3) *MVec3 {
	v[0] *= other[0]
	v[1] *= other[1]
	v[2] *= other[2]
	return v
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *MVec3) MulScalar(s float64) *MVec3 {
	v[0] *= s
	v[1] *= s
	v[2] *= s
	return v
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *MVec3) Muled(other *MVec3) *MVec3 {
	return &MVec3{v[0] * other[0], v[1] * other[1], v[2] * other[2]}
}

func (v *MVec3) MuledScalar(s float64) *MVec3 {
	return &MVec3{v[0] * s, v[1] * s, v[2] * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec3) Div(other *MVec3) *MVec3 {
	v[0] /= other[0]
	v[1] /= other[1]
	v[2] /= other[2]
	return v
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec3) DivScalar(s float64) *MVec3 {
	v[0] /= s
	v[1] /= s
	v[2] /= s
	return v
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec3) Dived(other *MVec3) *MVec3 {
	return &MVec3{v[0] / other[0], v[1] / other[1], v[2] / other[2]}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (v *MVec3) DivedScalar(s float64) *MVec3 {
	return &MVec3{v[0] / s, v[1] / s, v[2] / s}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *MVec3) Equals(other *MVec3) bool {
	return v.GetX() == other.GetX() && v.GetY() == other.GetY() && v.GetZ() == other.GetZ()
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *MVec3) NotEquals(other MVec3) bool {
	return v.GetX() != other.GetX() || v.GetY() != other.GetY() || v.GetZ() != other.GetZ()
}

// PracticallyEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *MVec3) PracticallyEquals(other *MVec3, epsilon float64) bool {
	return (math.Abs(v[0]-other[0]) <= epsilon) &&
		(math.Abs(v[1]-other[1]) <= epsilon) &&
		(math.Abs(v[2]-other[2]) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *MVec3) LessThan(other *MVec3) bool {
	return v.GetX() < other.GetX() && v.GetY() < other.GetY() && v.GetZ() < other.GetZ()
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *MVec3) LessThanOrEquals(other *MVec3) bool {
	return v.GetX() <= other.GetX() && v.GetY() <= other.GetY() && v.GetZ() <= other.GetZ()
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *MVec3) GreaterThan(other *MVec3) bool {
	return v.GetX() > other.GetX() && v.GetY() > other.GetY() && v.GetZ() > other.GetZ()
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *MVec3) GreaterThanOrEquals(other *MVec3) bool {
	return v.GetX() >= other.GetX() && v.GetY() >= other.GetY() && v.GetZ() >= other.GetZ()
}

// Invert ベクトルの各要素の符号を反転します (-v)
func (v *MVec3) Invert() *MVec3 {
	v[0] = -v[0]
	v[1] = -v[1]
	v[2] = -v[2]
	return v
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *MVec3) Inverted() *MVec3 {
	return &MVec3{-v[0], -v[1], -v[2]}
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec3) Abs() *MVec3 {
	v[0] = math.Abs(v[0])
	v[1] = math.Abs(v[1])
	v[2] = math.Abs(v[2])
	return v
}

// Absed ベクトルの各要素の絶対値を返します
func (v *MVec3) Absed() *MVec3 {
	return &MVec3{math.Abs(v[0]), math.Abs(v[1]), math.Abs(v[2])}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec3) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f", v.GetX(), v.GetY(), v.GetZ())))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec3) IsZero() bool {
	return v[0] == 0 && v[1] == 0 && v[2] == 0
}

// Length ベクトルの長さを返します
func (v *MVec3) Length() float64 {
	return math.Sqrt(v.LengthSqr())
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec3) LengthSqr() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

// Normalize ベクトルを正規化します
func (v *MVec3) Normalize() *MVec3 {
	sl := v.LengthSqr()
	if sl == 0 || sl == 1 {
		return v
	}
	return v.MulScalar(1 / math.Sqrt(sl))
}

// Normalized ベクトルを正規化した結果を返します
func (v *MVec3) Normalized() *MVec3 {
	vec := *v
	vec.Normalize()
	return &vec
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (v *MVec3) Angle(other *MVec3) float64 {
	vec := v.Dot(other) / (v.Length() * other.Length())
	// prevent NaN
	if vec > 1. {
		return 0
	} else if vec < -1. {
		return math.Pi
	}
	return math.Acos(vec)
}

// Degree ベクトルの角度(度数)を返します
func (v *MVec3) Degree(other *MVec3) float64 {
	radian := v.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (v *MVec3) Dot(other *MVec3) float64 {
	return v[0]*other[0] + v[1]*other[1] + v[2]*other[2]
}

// Cross ベクトルの外積を返します
func (v1 *MVec3) Cross(v2 *MVec3) *MVec3 {
	return &MVec3{v1[1]*v2[2] - v1[2]*v2[1], v1[2]*v2[0] - v1[0]*v2[2], v1[0]*v2[1] - v1[1]*v2[0]}
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *MVec3) Min() *MVec3 {
	min := v.GetX()
	if v.GetY() < min {
		min = v.GetY()
	}
	if v.GetZ() < min {
		min = v.GetZ()
	}
	return &MVec3{min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *MVec3) Max() *MVec3 {
	max := v.GetX()
	if v.GetY() > max {
		max = v.GetY()
	}
	if v.GetZ() > max {
		max = v.GetZ()
	}
	return &MVec3{max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *MVec3) Interpolate(other *MVec3, t float64) *MVec3 {
	t1 := 1 - t
	return &MVec3{
		v[0]*t1 + other[0]*t,
		v[1]*t1 + other[1]*t,
		v[2]*t1 + other[2]*t,
	}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec3) Clamp(min, max *MVec3) *MVec3 {
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
func (v *MVec3) Clamped(min, max *MVec3) *MVec3 {
	result := *v
	result.Clamp(min, max)
	return &result
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec3) Clamp01() *MVec3 {
	return v.Clamp(&MVec3Zero, &MVec3UnitXYZ)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec3) Clamped01() *MVec3 {
	result := *v
	result.Clamp01()
	return &result
}

// Copy
func (v *MVec3) Copy() *MVec3 {
	return &MVec3{v.GetX(), v.GetY(), v.GetZ()}
}

// Vector
func (v *MVec3) Vector() []float64 {
	return []float64{v.GetX(), v.GetY(), v.GetZ()}
}

func (v *MVec3) ToMat4() *MMat4 {
	mat := NewMMat4()
	mat[0][3] = v.GetX()
	mat[1][3] = v.GetY()
	mat[2][3] = v.GetZ()
	return mat
}

// 線形補間
func LerpFloat(v1, v2 float64, t float64) float64 {
	return v1 + ((v2 - v1) * t)
}

// Clamp01 ベクトルの各要素をmin～maxの範囲内にクランプします
func ClampFloat(v float64, min float64, max float64) float64 {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

// 線形補間
func LerpVec3(v1, v2 *MVec3, t float64) *MVec3 {
	return (v2.Sub(v1)).MulScalar(t).Added(v1)
}

func (v *MVec3) Round() *MVec3 {
	return &MVec3{
		math.Round(v.GetX()),
		math.Round(v.GetY()),
		math.Round(v.GetZ()),
	}
}

// ToLocalMatrix4x4 自身をローカル軸とした場合の回転行列を取得します
func (v *MVec3) ToLocalMatrix4x4() *MMat4 {
	if v.IsZero() {
		return NewMMat4()
	}

	// ローカルX軸の方向ベクトル
	xAxis := v.Copy()
	normXAxis := xAxis.Length()
	if normXAxis == 0 {
		return NewMMat4()
	}
	xAxis.DivScalar(normXAxis)

	if math.IsNaN(xAxis.GetX()) || math.IsNaN(xAxis.GetY()) || math.IsNaN(xAxis.GetZ()) {
		return NewMMat4()
	}

	// ローカルZ軸の方向ベクトル
	zAxis := &MVec3{0.0, 0.0, -1.0}
	if zAxis.Equals(v) {
		// 自身がほぼZ軸ベクトルの場合、別ベクトルを与える
		zAxis = &MVec3{0.0, 1.0, 0.0}
	}

	// ローカルY軸の方向ベクトル
	yAxis := zAxis.Cross(xAxis)
	normYAxis := yAxis.Length()
	if normYAxis == 0 {
		return NewMMat4()
	}
	yAxis.DivScalar(normYAxis)

	if math.IsNaN(yAxis.GetX()) || math.IsNaN(yAxis.GetY()) || math.IsNaN(yAxis.GetZ()) {
		return NewMMat4()
	}

	zAxis = xAxis.Cross(yAxis)
	normZAxis := zAxis.Length()
	zAxis.DivScalar(normZAxis)

	// ローカル軸に合わせた回転行列を作成する
	rotationMatrix := NewMMat4()
	rotationMatrix[0][0] = xAxis.GetX()
	rotationMatrix[1][0] = xAxis.GetY()
	rotationMatrix[2][0] = xAxis.GetZ()
	rotationMatrix[0][1] = yAxis.GetX()
	rotationMatrix[1][1] = yAxis.GetY()
	rotationMatrix[2][1] = yAxis.GetZ()
	rotationMatrix[0][2] = zAxis.GetX()
	rotationMatrix[1][2] = zAxis.GetY()
	rotationMatrix[2][2] = zAxis.GetZ()

	return rotationMatrix
}

// 標準偏差を加味したmean処理
func StdMeanVec3(values []MVec3, err float64) *MVec3 {
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
	return &MVec3{mean[0], mean[1], mean[2]}
}

// One 0を1に変える
func (v *MVec3) One() *MVec3 {
	vec := v.Vector()
	epsilon := 1e-14
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec3{vec[0], vec[1], vec[2]}
}

func (v *MVec3) Distance(other *MVec3) float64 {
	s := v.Subed(other)
	return s.Length()
}

func (v *MVec3) Project(other *MVec3) *MVec3 {
	return other.MuledScalar(v.Dot(other) / other.LengthSqr())
}

// ボーンから見た頂点ローカル位置を求める
// vertexPositions: グローバル頂点位置
// startBonePosition: 親ボーン位置
// endBonePosition: 子ボーン位置
func GetVertexLocalPositions(vertexPositions []*MVec3, startBonePosition *MVec3, endBonePosition *MVec3) []*MVec3 {
	vertexSize := len(vertexPositions)
	boneVector := endBonePosition.Sub(startBonePosition)
	boneDirection := boneVector.Normalized()

	localPositions := make([]*MVec3, vertexSize)
	for i := 0; i < vertexSize; i++ {
		vertexPosition := vertexPositions[i]
		subedVertexPosition := vertexPosition.Subed(startBonePosition)
		projection := subedVertexPosition.Project(boneDirection)
		localPosition := endBonePosition.Added(projection)
		localPositions[i] = localPosition
	}

	return localPositions
}
