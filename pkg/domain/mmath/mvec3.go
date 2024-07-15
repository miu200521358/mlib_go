package mmath

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

var (
	MVec3Zero = &MVec3{}

	MVec3UnitX = &MVec3{1, 0, 0}
	MVec3UnitY = &MVec3{0, 1, 0}
	MVec3UnitZ = &MVec3{0, 0, 1}
	MVec3One   = &MVec3{1, 1, 1}

	MVec3UnitXInv = &MVec3{-1, 0, 0}
	MVec3UnitYInv = &MVec3{0, -1, 0}
	MVec3UnitZInv = &MVec3{0, 0, -1}
	MVec3OneInv   = &MVec3{-1, -1, -1}

	MVec3MinVal = &MVec3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	MVec3MaxVal = &MVec3{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64}
)

type MVec3 struct {
	X float64
	Y float64
	Z float64
}

func NewMVec3() *MVec3 {
	return &MVec3{}
}

func (v *MVec3) GetXY() *MVec2 {
	return &MVec2{v.X, v.Y}
}

func (v *MVec3) IsOnlyX() bool {
	return !NearEquals(v.X, 0, 1e-10) &&
		NearEquals(v.Y, 0, 1e-10) &&
		NearEquals(v.Z, 0, 1e-10)
}

func (v *MVec3) IsOnlyY() bool {
	return NearEquals(v.X, 0, 1e-10) &&
		!NearEquals(v.Y, 0, 1e-10) &&
		NearEquals(v.Z, 0, 1e-10)
}

func (v *MVec3) IsOnlyZ() bool {
	return NearEquals(v.X, 0, 1e-10) &&
		NearEquals(v.Y, 0, 1e-10) &&
		!NearEquals(v.Z, 0, 1e-10)
}

// String T の文字列表現を返します。
func (v *MVec3) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f]", v.X, v.Y, v.Z)
}

// MMD MMD(MikuMikuDance)座標系に変換された3次元ベクトルを返します
func (v *MVec3) MMD() *MVec3 {
	return &MVec3{v.X, v.Y, v.Z}
}

// Add ベクトルに他のベクトルを加算します
func (v *MVec3) Add(other *MVec3) *MVec3 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (v *MVec3) AddScalar(s float64) *MVec3 {
	v.X += s
	v.Y += s
	v.Z += s
	return v
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (v *MVec3) Added(other *MVec3) *MVec3 {
	return &MVec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v *MVec3) AddedScalar(s float64) *MVec3 {
	return &MVec3{v.X + s, v.Y + s, v.Z + s}
}

// Sub ベクトルから他のベクトルを減算します
func (v *MVec3) Sub(other *MVec3) *MVec3 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (v *MVec3) SubScalar(s float64) *MVec3 {
	v.X -= s
	v.Y -= s
	v.Z -= s
	return v
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (v *MVec3) Subed(other *MVec3) *MVec3 {
	return &MVec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v *MVec3) SubedScalar(s float64) *MVec3 {
	return &MVec3{v.X - s, v.Y - s, v.Z - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (v *MVec3) Mul(other *MVec3) *MVec3 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (v *MVec3) MulScalar(s float64) *MVec3 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (v *MVec3) Muled(other *MVec3) *MVec3 {
	return &MVec3{v.X * other.X, v.Y * other.Y, v.Z * other.Z}
}

func (v *MVec3) MuledScalar(s float64) *MVec3 {
	return &MVec3{v.X * s, v.Y * s, v.Z * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (v *MVec3) Div(other *MVec3) *MVec3 {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	return v
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (v *MVec3) DivScalar(s float64) *MVec3 {
	v.X /= s
	v.Y /= s
	v.Z /= s
	return v
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (v *MVec3) Dived(other *MVec3) *MVec3 {
	return &MVec3{v.X / other.X, v.Y / other.Y, v.Z / other.Z}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (v *MVec3) DivedScalar(s float64) *MVec3 {
	return &MVec3{v.X / s, v.Y / s, v.Z / s}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (v *MVec3) Equals(other *MVec3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (v *MVec3) NotEquals(other MVec3) bool {
	return v.X != other.X || v.Y != other.Y || v.Z != other.Z
}

// NearEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (v *MVec3) NearEquals(other *MVec3, epsilon float64) bool {
	return (math.Abs(v.X-other.X) <= epsilon) &&
		(math.Abs(v.Y-other.Y) <= epsilon) &&
		(math.Abs(v.Z-other.Z) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (v *MVec3) LessThan(other *MVec3) bool {
	return v.X < other.X && v.Y < other.Y && v.Z < other.Z
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (v *MVec3) LessThanOrEquals(other *MVec3) bool {
	return v.X <= other.X && v.Y <= other.Y && v.Z <= other.Z
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (v *MVec3) GreaterThan(other *MVec3) bool {
	return v.X > other.X && v.Y > other.Y && v.Z > other.Z
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (v *MVec3) GreaterThanOrEquals(other *MVec3) bool {
	return v.X >= other.X && v.Y >= other.Y && v.Z >= other.Z
}

// Inverse ベクトルの各要素の符号を反転します (-v)
func (v *MVec3) Inverse() *MVec3 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

// Inverted ベクトルの各要素の符号を反転した結果を返します (-v)
func (v *MVec3) Inverted() *MVec3 {
	return &MVec3{-v.X, -v.Y, -v.Z}
}

// Abs ベクトルの各要素の絶対値を返します
func (v *MVec3) Abs() *MVec3 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	v.Z = math.Abs(v.Z)
	return v
}

// Absed ベクトルの各要素の絶対値を返します
func (v *MVec3) Absed() *MVec3 {
	return &MVec3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}
}

// Hash ベクトルのハッシュ値を計算します
func (v *MVec3) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f", v.X, v.Y, v.Z)))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (v *MVec3) IsZero() bool {
	return v.NearEquals(MVec3Zero, 1e-10)
}

// IsZero ベクトルが1ベクトルかどうかをチェックします
func (v *MVec3) IsOne() bool {
	return v.NearEquals(MVec3One, 1e-10)
}

// Length ベクトルの長さを返します
func (v *MVec3) Length() float64 {
	return mgl64.Vec3{v.X, v.Y, v.Z}.Len()
}

// LengthSqr ベクトルの長さの2乗を返します
func (v *MVec3) LengthSqr() float64 {
	return mgl64.Vec3{v.X, v.Y, v.Z}.LenSqr()
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
	vec := MVec3{v.X, v.Y, v.Z}
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
	return mgl64.Vec3{v.X, v.Y, v.Z}.Dot(mgl64.Vec3{other.X, other.Y, other.Z})
}

// Cross ベクトルの外積を返します
func (v1 *MVec3) Cross(v2 *MVec3) *MVec3 {
	v := mgl64.Vec3{v1.X, v1.Y, v1.Z}.Cross(mgl64.Vec3{v2.X, v2.Y, v2.Z})
	return &MVec3{v[0], v[1], v[2]}
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (v *MVec3) Min() *MVec3 {
	min := v.X
	if v.Y < min {
		min = v.Y
	}
	if v.Z < min {
		min = v.Z
	}
	return &MVec3{min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (v *MVec3) Max() *MVec3 {
	max := v.X
	if v.Y > max {
		max = v.Y
	}
	if v.Z > max {
		max = v.Z
	}
	return &MVec3{max, max, max}
}

// Interpolate ベクトルの線形補間を行います
func (v *MVec3) Interpolate(other *MVec3, t float64) *MVec3 {
	t1 := 1 - t
	return &MVec3{
		v.X*t1 + other.X*t,
		v.Y*t1 + other.Y*t,
		v.Z*t1 + other.Z*t,
	}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (v *MVec3) Clamp(min, max *MVec3) *MVec3 {
	v.X = ClampedFloat(v.X, min.X, max.X)
	v.Y = ClampedFloat(v.Y, min.Y, max.Y)
	v.Z = ClampedFloat(v.Z, min.Z, max.Z)
	return v
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (v *MVec3) Clamped(min, max *MVec3) *MVec3 {
	result := MVec3{v.X, v.Y, v.Z}
	result.Clamp(min, max)
	return &result
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (v *MVec3) Clamp01() *MVec3 {
	return v.Clamp(MVec3Zero, MVec3One)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (v *MVec3) Clamped01() *MVec3 {
	result := MVec3{v.X, v.Y, v.Z}
	result.Clamp01()
	return &result
}

// Copy
func (v *MVec3) Copy() *MVec3 {
	return &MVec3{v.X, v.Y, v.Z}
}

// Vector
func (v *MVec3) Vector() []float64 {
	return []float64{v.X, v.Y, v.Z}
}

func (v *MVec3) ToMat4() *MMat4 {
	mat := NewMMat4()
	mat[3] = v.X
	mat[7] = v.Y
	mat[11] = v.Z
	return mat
}

func (v *MVec3) ToScaleMat4() *MMat4 {
	mat := NewMMat4()
	mat[0] = v.X
	mat[5] = v.Y
	mat[10] = v.Z
	return mat
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (v *MVec3) ClampIfVerySmall() *MVec3 {
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
	return v
}

// 線形補間
func (v1 *MVec3) Lerp(v2 *MVec3, t float64) *MVec3 {
	return (v2.Subed(v1)).MulScalar(t).Add(v1)
}

func (v *MVec3) Round() *MVec3 {
	return &MVec3{
		math.Round(v.X),
		math.Round(v.Y),
		math.Round(v.Z),
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

	if math.IsNaN(xAxis.X) || math.IsNaN(xAxis.Y) || math.IsNaN(xAxis.Z) {
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

	if math.IsNaN(yAxis.X) || math.IsNaN(yAxis.Y) || math.IsNaN(yAxis.Z) {
		return NewMMat4()
	}

	zAxis = xAxis.Cross(yAxis)
	normZAxis := zAxis.Length()
	zAxis.DivScalar(normZAxis)

	// ローカル軸に合わせた回転行列を作成する
	rotationMatrix := NewMMat4ByValues(
		xAxis.X, yAxis.X, zAxis.X, 0,
		xAxis.Y, yAxis.Y, zAxis.Y, 0,
		xAxis.Z, yAxis.Z, zAxis.Z, 0,
		0, 0, 0, 1,
	)

	return rotationMatrix
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
	return v.Subed(other).Length()
}

func (v *MVec3) Distances(others []*MVec3) []float64 {
	distances := make([]float64, len(others))
	for i, other := range others {
		distances[i] = v.Distance(other)
	}
	return distances
}

// 2点間のベクトルと点Pの直交距離を計算
func DistanceFromPointToLine(a, b, p *MVec3) float64 {
	lineVec := b.Subed(a)               // 線分ABのベクトル
	pointVec := p.Subed(a)              // 点Pから点Aへのベクトル
	crossVec := lineVec.Cross(pointVec) // 外積ベクトル
	area := crossVec.Length()           // 平行四辺形の面積
	lineLength := lineVec.Length()      // 線分ABの長さ
	return area / lineLength            // 点Pから線分ABへの距離
}

// 2点間のベクトルと、点Pを含むカメラ平面と平行な面、との距離を計算
func DistanceFromPlaneToLine(near, far, forward, right, up, p *MVec3) float64 {
	// ステップ1: カメラ平面の法線ベクトルを計算
	normal := forward.Cross(right)

	// ステップ2: 点Pからカメラ平面へのベクトルを計算
	vectorToPlane := p.Subed(near)

	// ステップ3: 距離を計算
	distance := math.Abs(vectorToPlane.Dot(normal)) / normal.Length()

	return distance
}

// 2点間のベクトルと、点Pを含むカメラ平面と平行な面、との交点を計算
func IntersectLinePlane(near, far, forward, right, up, p *MVec3) *MVec3 {
	// ステップ1: カメラ平面の法線ベクトルを計算
	normal := forward.Cross(right)

	// ステップ2: nearからfarへのベクトルを計算
	direction := far.Subed(near)

	// ステップ3: 平面の方程式のD成分を計算
	D := -normal.Dot(p)

	// ステップ4: 方向ベクトルと法線ベクトルが平行かどうかを確認
	denom := normal.Dot(direction)
	if math.Abs(denom) < 1e-6 { // ほぼ0に近い場合、平行とみなす
		return nil // 平行ならば交点は存在しない
	}

	// ステップ5: 直線と平面の交点を計算
	t := -(normal.Dot(near) + D) / denom
	intersection := near.Added(direction.MuledScalar(t))
	return intersection
}

// DistanceLineToPoints 線分と点の距離を計算します
func DistanceLineToPoints(worldPos *MVec3, points []*MVec3) []float64 {
	distances := make([]float64, len(points))

	// worldPos の Z方向のベクトル
	worldDirection := worldPos.Added(MVec3UnitZInv)

	for i, p := range points {
		// 点PとworldPosのZ方向のベクトルとの距離を計算
		distances[i] = DistanceFromPointToLine(worldPos, worldDirection, p)
		mlog.D("DistanceLineToPoints[%d]: d: %.3f\n", i, distances[i])
	}

	return distances
}

func Distances(v *MVec3, others []*MVec3) []float64 {
	distances := make([]float64, len(others))
	for i, other := range others {
		distances[i] = v.Distance(other)
	}
	return distances
}

func (v *MVec3) Project(other *MVec3) *MVec3 {
	return other.MuledScalar(v.Dot(other) / other.LengthSqr())
}

// 標準偏差を加味したmean処理
func StdMeanVec3(values []MVec3, err float64) *MVec3 {
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
	return &MVec3{mean[0], mean[1], mean[2]}
}

// 点が直方体内にあるかどうかを判定する関数
func (point *MVec3) IsPointInsideBox(min, max *MVec3) bool {
	return point.X >= min.X && point.X <= max.X &&
		point.Y >= min.Y && point.Y <= max.Y &&
		point.Z >= min.Z && point.Z <= max.Z
}

// 直方体の境界を計算する関数
func CalculateBoundingBox(points ...*MVec3) (minPos, maxPos *MVec3) {
	minPos = &MVec3{math.Inf(1), math.Inf(1), math.Inf(1)}
	maxPos = &MVec3{math.Inf(-1), math.Inf(-1), math.Inf(-1)}

	for _, p := range points {
		if p.X < minPos.X {
			minPos.X = p.X
		}
		if p.Y < minPos.Y {
			minPos.Y = p.Y
		}
		if p.Z < minPos.Z {
			minPos.Z = p.Z
		}
		if p.X > maxPos.X {
			maxPos.X = p.X
		}
		if p.Y > maxPos.Y {
			maxPos.Y = p.Y
		}
		if p.Z > maxPos.Z {
			maxPos.Z = p.Z
		}
	}

	return minPos, maxPos
}
