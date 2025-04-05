package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
	"sort"

	"github.com/go-gl/mathgl/mgl64"
)

var (
	MVec3Zero = &MVec3{}

	MVec3UnitX = &MVec3{1, 0, 0}
	MVec3UnitY = &MVec3{0, 1, 0}
	MVec3UnitZ = &MVec3{0, 0, 1}
	MVec3One   = &MVec3{1, 1, 1}

	MVec3UnitXNeg = &MVec3{-1, 0, 0}
	MVec3UnitYNeg = &MVec3{0, -1, 0}
	MVec3UnitZNeg = &MVec3{0, 0, -1}
	MVec3OneNeg   = &MVec3{-1, -1, -1}

	MVec3MinVal = &MVec3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	MVec3MaxVal = &MVec3{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64}
)

type MVec3 struct {
	X float64 `json:"x"` // X座標
	Y float64 `json:"y"` // Y座標
	Z float64 `json:"z"` // Z座標
}

func NewMVec3() *MVec3 {
	return &MVec3{}
}

func (vec3 *MVec3) GetXY() *MVec2 {
	return &MVec2{vec3.X, vec3.Y}
}

func (vec3 *MVec3) IsOnlyX() bool {
	return !NearEquals(vec3.X, 0, 1e-10) &&
		NearEquals(vec3.Y, 0, 1e-10) &&
		NearEquals(vec3.Z, 0, 1e-10)
}

func (vec3 *MVec3) IsOnlyY() bool {
	return NearEquals(vec3.X, 0, 1e-10) &&
		!NearEquals(vec3.Y, 0, 1e-10) &&
		NearEquals(vec3.Z, 0, 1e-10)
}

func (vec3 *MVec3) IsOnlyZ() bool {
	return NearEquals(vec3.X, 0, 1e-10) &&
		NearEquals(vec3.Y, 0, 1e-10) &&
		!NearEquals(vec3.Z, 0, 1e-10)
}

// String T の文字列表現を返します。
func (vec3 *MVec3) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f]", vec3.X, vec3.Y, vec3.Z)
}

// MMD MMD(MikuMikuDance)座標系に変換された3次元ベクトルを返します
func (vec3 *MVec3) MMD() *MVec3 {
	return &MVec3{vec3.X, vec3.Y, vec3.Z}
}

// Add ベクトルに他のベクトルを加算します
func (vec3 *MVec3) Add(other *MVec3) *MVec3 {
	vec3.X += other.X
	vec3.Y += other.Y
	vec3.Z += other.Z
	return vec3
}

// AddScalar ベクトルの各要素にスカラーを加算します
func (vec3 *MVec3) AddScalar(s float64) *MVec3 {
	vec3.X += s
	vec3.Y += s
	vec3.Z += s
	return vec3
}

// Added ベクトルに他のベクトルを加算した結果を返します
func (vec3 *MVec3) Added(other *MVec3) *MVec3 {
	return &MVec3{vec3.X + other.X, vec3.Y + other.Y, vec3.Z + other.Z}
}

func (vec3 *MVec3) AddedScalar(s float64) *MVec3 {
	return &MVec3{vec3.X + s, vec3.Y + s, vec3.Z + s}
}

// Sub ベクトルから他のベクトルを減算します
func (vec3 *MVec3) Sub(other *MVec3) *MVec3 {
	vec3.X -= other.X
	vec3.Y -= other.Y
	vec3.Z -= other.Z
	return vec3
}

// SubScalar ベクトルの各要素からスカラーを減算します
func (vec3 *MVec3) SubScalar(s float64) *MVec3 {
	vec3.X -= s
	vec3.Y -= s
	vec3.Z -= s
	return vec3
}

// Subed ベクトルから他のベクトルを減算した結果を返します
func (vec3 *MVec3) Subed(other *MVec3) *MVec3 {
	return &MVec3{vec3.X - other.X, vec3.Y - other.Y, vec3.Z - other.Z}
}

func (vec3 *MVec3) SubedScalar(s float64) *MVec3 {
	return &MVec3{vec3.X - s, vec3.Y - s, vec3.Z - s}
}

// Mul ベクトルの各要素に他のベクトルの各要素を乗算します
func (vec3 *MVec3) Mul(other *MVec3) *MVec3 {
	vec3.X *= other.X
	vec3.Y *= other.Y
	vec3.Z *= other.Z
	return vec3
}

// MulScalar ベクトルの各要素にスカラーを乗算します
func (vec3 *MVec3) MulScalar(s float64) *MVec3 {
	vec3.X *= s
	vec3.Y *= s
	vec3.Z *= s
	return vec3
}

// Muled ベクトルの各要素に他のベクトルの各要素を乗算した結果を返します
func (vec3 *MVec3) Muled(other *MVec3) *MVec3 {
	return &MVec3{vec3.X * other.X, vec3.Y * other.Y, vec3.Z * other.Z}
}

func (vec3 *MVec3) MuledScalar(s float64) *MVec3 {
	return &MVec3{vec3.X * s, vec3.Y * s, vec3.Z * s}
}

// Div ベクトルの各要素を他のベクトルの各要素で除算します
func (vec3 *MVec3) Div(other *MVec3) *MVec3 {
	vec3.X /= other.X
	vec3.Y /= other.Y
	vec3.Z /= other.Z
	return vec3
}

// DivScalar ベクトルの各要素をスカラーで除算します
func (vec3 *MVec3) DivScalar(s float64) *MVec3 {
	vec3.X /= s
	vec3.Y /= s
	vec3.Z /= s
	return vec3
}

// Dived ベクトルの各要素を他のベクトルの各要素で除算した結果を返します
func (vec3 *MVec3) Dived(other *MVec3) *MVec3 {
	return &MVec3{vec3.X / other.X, vec3.Y / other.Y, vec3.Z / other.Z}
}

// DivedScalar ベクトルの各要素をスカラーで除算した結果を返します
func (vec3 *MVec3) DivedScalar(s float64) *MVec3 {
	return &MVec3{vec3.X / s, vec3.Y / s, vec3.Z / s}
}

// Equal ベクトルが他のベクトルと等しいかどうかをチェックします
func (vec3 *MVec3) Equals(other *MVec3) bool {
	return vec3.X == other.X && vec3.Y == other.Y && vec3.Z == other.Z
}

// NotEqual ベクトルが他のベクトルと等しくないかどうかをチェックします
func (vec3 *MVec3) NotEquals(other MVec3) bool {
	return vec3.X != other.X || vec3.Y != other.Y || vec3.Z != other.Z
}

// NearEquals ベクトルが他のベクトルとほぼ等しいかどうかをチェックします
func (vec3 *MVec3) NearEquals(other *MVec3, epsilon float64) bool {
	return (math.Abs(vec3.X-other.X) <= epsilon) &&
		(math.Abs(vec3.Y-other.Y) <= epsilon) &&
		(math.Abs(vec3.Z-other.Z) <= epsilon)
}

// LessThan ベクトルが他のベクトルより小さいかどうかをチェックします (<)
func (vec3 *MVec3) LessThan(other *MVec3) bool {
	return vec3.X < other.X && vec3.Y < other.Y && vec3.Z < other.Z
}

// LessThanOrEqual ベクトルが他のベクトル以下かどうかをチェックします (<=)
func (vec3 *MVec3) LessThanOrEquals(other *MVec3) bool {
	return vec3.X <= other.X && vec3.Y <= other.Y && vec3.Z <= other.Z
}

// GreaterThan ベクトルが他のベクトルより大きいかどうかをチェックします (>)
func (vec3 *MVec3) GreaterThan(other *MVec3) bool {
	return vec3.X > other.X && vec3.Y > other.Y && vec3.Z > other.Z
}

// GreaterThanOrEqual ベクトルが他のベクトル以上かどうかをチェックします (>=)
func (vec3 *MVec3) GreaterThanOrEquals(other *MVec3) bool {
	return vec3.X >= other.X && vec3.Y >= other.Y && vec3.Z >= other.Z
}

// Negate ベクトルの各要素の符号を反転します (-v)
func (vec3 *MVec3) Negate() *MVec3 {
	vec3.X = -vec3.X
	vec3.Y = -vec3.Y
	vec3.Z = -vec3.Z
	return vec3
}

// Negated ベクトルの各要素の符号を反転した結果を返します (-v)
func (vec3 *MVec3) Negated() *MVec3 {
	return &MVec3{-vec3.X, -vec3.Y, -vec3.Z}
}

// Abs ベクトルの各要素の絶対値を返します
func (vec3 *MVec3) Abs() *MVec3 {
	vec3.X = math.Abs(vec3.X)
	vec3.Y = math.Abs(vec3.Y)
	vec3.Z = math.Abs(vec3.Z)
	return vec3
}

// Absed ベクトルの各要素の絶対値を返します
func (vec3 *MVec3) Absed() *MVec3 {
	return &MVec3{math.Abs(vec3.X), math.Abs(vec3.Y), math.Abs(vec3.Z)}
}

// Hash ベクトルのハッシュ値を計算します
func (vec3 *MVec3) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f", vec3.X, vec3.Y, vec3.Z)))
	return h.Sum64()
}

// IsZero ベクトルがゼロベクトルかどうかをチェックします
func (vec3 *MVec3) IsZero() bool {
	return vec3.NearEquals(MVec3Zero, 1e-10)
}

// IsZero ベクトルが1ベクトルかどうかをチェックします
func (vec3 *MVec3) IsOne() bool {
	return vec3.NearEquals(MVec3One, 1e-10)
}

// Length ベクトルの長さを返します
func (vec3 *MVec3) Length() float64 {
	return mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.Len()
}

// LengthSqr ベクトルの長さの2乗を返します
func (vec3 *MVec3) LengthSqr() float64 {
	return mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.LenSqr()
}

// Normalize ベクトルを正規化します
func (vec3 *MVec3) Normalize() *MVec3 {
	sl := vec3.LengthSqr()
	if sl == 0 || sl == 1 {
		return vec3
	}
	return vec3.MulScalar(1 / math.Sqrt(sl))
}

// Normalized ベクトルを正規化した結果を返します
func (vec3 *MVec3) Normalized() *MVec3 {
	vec := MVec3{vec3.X, vec3.Y, vec3.Z}
	vec.Normalize()
	return &vec
}

// Angle ベクトルの角度(ラジアン角度)を返します
func (vec3 *MVec3) Angle(other *MVec3) float64 {
	vec := vec3.Dot(other) / (vec3.Length() * other.Length())
	// prevent NaN
	if vec > 1. {
		return 0
	} else if vec < -1. {
		return math.Pi
	}
	return math.Acos(vec)
}

// Degree ベクトルの角度(度数)を返します
func (vec3 *MVec3) Degree(other *MVec3) float64 {
	radian := vec3.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

// Dot ベクトルの内積を返します
func (vec3 *MVec3) Dot(other *MVec3) float64 {
	return mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.Dot(mgl64.Vec3{other.X, other.Y, other.Z})
}

// Cross ベクトルの外積を返します
func (vec3 *MVec3) Cross(other *MVec3) *MVec3 {
	v := mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.Cross(mgl64.Vec3{other.X, other.Y, other.Z})
	return &MVec3{v[0], v[1], v[2]}
}

// Min ベクトルの各要素の最小値をTの各要素に設定して返します
func (vec3 *MVec3) Min() *MVec3 {
	min := vec3.X
	if vec3.Y < min {
		min = vec3.Y
	}
	if vec3.Z < min {
		min = vec3.Z
	}
	return &MVec3{min, min, min}
}

// Max ベクトルの各要素の最大値を返します
func (vec3 *MVec3) Max() *MVec3 {
	max := vec3.X
	if vec3.Y > max {
		max = vec3.Y
	}
	if vec3.Z > max {
		max = vec3.Z
	}
	return &MVec3{max, max, max}
}

// Clamp ベクトルの各要素を指定された範囲内にクランプします
func (vec3 *MVec3) Clamp(min, max *MVec3) *MVec3 {
	vec3.X = Clamped(vec3.X, min.X, max.X)
	vec3.Y = Clamped(vec3.Y, min.Y, max.Y)
	vec3.Z = Clamped(vec3.Z, min.Z, max.Z)
	return vec3
}

// Clamped ベクトルの各要素を指定された範囲内にクランプした結果を返します
func (vec3 *MVec3) Clamped(min, max *MVec3) *MVec3 {
	result := MVec3{vec3.X, vec3.Y, vec3.Z}
	result.Clamp(min, max)
	return &result
}

// Clamp01 ベクトルの各要素を0.0～1.0の範囲内にクランプします
func (vec3 *MVec3) Clamp01() *MVec3 {
	return vec3.Clamp(MVec3Zero, MVec3One)
}

// Clamped01 ベクトルの各要素を0.0～1.0の範囲内にクランプした結果を返します
func (vec3 *MVec3) Clamped01() *MVec3 {
	result := MVec3{vec3.X, vec3.Y, vec3.Z}
	result.Clamp01()
	return &result
}

// Copy
func (vec3 *MVec3) Copy() *MVec3 {
	return &MVec3{vec3.X, vec3.Y, vec3.Z}
}

// Vector
func (vec3 *MVec3) Vector() []float64 {
	return []float64{vec3.X, vec3.Y, vec3.Z}
}

func (vec3 *MVec3) ToMat4() *MMat4 {
	mat := NewMMat4()
	mat[12] = vec3.X
	mat[13] = vec3.Y
	mat[14] = vec3.Z
	return mat
}

func (vec3 *MVec3) ToScaleMat4() *MMat4 {
	mat := NewMMat4()
	mat[0] = vec3.X
	mat[5] = vec3.Y
	mat[10] = vec3.Z
	return mat
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (vec3 *MVec3) ClampIfVerySmall() *MVec3 {
	epsilon := 1e-6
	if math.Abs(vec3.X) < epsilon {
		vec3.X = 0
	}
	if math.Abs(vec3.Y) < epsilon {
		vec3.Y = 0
	}
	if math.Abs(vec3.Z) < epsilon {
		vec3.Z = 0
	}
	return vec3
}

func (vec3 *MVec3) RadToDeg() *MVec3 {
	return &MVec3{RadToDeg(vec3.X), RadToDeg(vec3.Y), RadToDeg(vec3.Z)}
}

func (vec3 *MVec3) DegToRad() *MVec3 {
	return &MVec3{DegToRad(vec3.X), DegToRad(vec3.Y), DegToRad(vec3.Z)}
}

func (vec3 *MVec3) RadToQuaternion() *MQuaternion {
	return NewMQuaternionFromRadians(vec3.X, vec3.Y, vec3.Z)
}

func (vec3 *MVec3) DegToQuaternion() *MQuaternion {
	return NewMQuaternionFromDegrees(vec3.X, vec3.Y, vec3.Z)
}

// 線形補間
func (vec3 *MVec3) Lerp(other *MVec3, t float64) *MVec3 {
	switch {
	case t <= 0:
		return vec3
	case t >= 1:
		return other
	default:
		if vec3.NearEquals(other, 1e-8) {
			return vec3.Copy()
		}

		return vec3.Add((other.Subed(vec3)).MulScalar(t))
	}
}

// ToLocalMat 自身をローカル軸とした場合の回転行列を取得します
func (vec3 *MVec3) ToLocalMat() *MMat4 {
	if vec3.IsZero() {
		return NewMMat4()
	}

	// 正規化されたローカルX軸ベクトルを取得
	v := vec3.Normalized()

	// ローカルY軸を計算
	var up *MVec3
	if math.Abs(v.Y) < 1-1e-6 {
		up = MVec3UnitY
	} else {
		up = MVec3UnitZ
	}

	// ローカルZ軸を計算
	u := up.Cross(v).Normalized()

	// ローカルY軸を計算
	w := v.Cross(u).Normalized()

	// ローカル座標系の回転行列を構築
	rotationMatrix := NewMMat4ByValues(
		v.X, v.Y, v.Z, 0,
		w.X, w.Y, w.Z, 0,
		u.X, u.Y, u.Z, 0,
		0, 0, 0, 1,
	)

	return rotationMatrix
}

func (vec3 *MVec3) ToScaleLocalMat(scales *MVec3) *MMat4 {
	if vec3.IsZero() || vec3.IsOne() {
		return NewMMat4()
	}

	// 軸方向の回転行列
	rotationMatrix := vec3.ToLocalMat()

	return rotationMatrix.Muled(scales.ToScaleMat4()).Muled(rotationMatrix.Inverted())
}

// One 0を1に変える
func (vec3 *MVec3) One() *MVec3 {
	vec := vec3.Vector()
	epsilon := 1e-14
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec3{vec[0], vec[1], vec[2]}
}

func (vec3 *MVec3) Distance(other *MVec3) float64 {
	return vec3.Subed(other).Length()
}

func (vec3 *MVec3) Distances(others []*MVec3) []float64 {
	distances := make([]float64, len(others))
	for i, other := range others {
		distances[i] = vec3.Distance(other)
	}
	return distances
}

func (vec3 *MVec3) Effective() *MVec3 {
	vec3.X = Effective(vec3.X)
	vec3.Y = Effective(vec3.Y)
	vec3.Z = Effective(vec3.Z)
	return vec3
}

// 2点間のベクトルと点Pの直交距離を計算
func DistanceFromPointToLine(vec1, vec2, point *MVec3) float64 {
	lineVec := vec2.Subed(vec1)         // 線分ABのベクトル
	pointVec := point.Subed(vec1)       // 点Pから点Aへのベクトル
	crossVec := lineVec.Cross(pointVec) // 外積ベクトル
	area := crossVec.Length()           // 平行四辺形の面積
	lineLength := lineVec.Length()      // 線分ABの長さ
	return area / lineLength            // 点Pから線分ABへの距離
}

// 2点間のベクトルと、点Pを含むカメラ平面と平行な面、との距離を計算
func DistanceFromPlaneToLine(near, far, forward, right, up, point *MVec3) float64 {
	// ステップ1: カメラ平面の法線ベクトルを計算
	normal := forward.Cross(right)

	// ステップ2: 点Pからカメラ平面へのベクトルを計算
	vectorToPlane := point.Subed(near)

	// ステップ3: 距離を計算
	distance := math.Abs(vectorToPlane.Dot(normal)) / normal.Length()

	return distance
}

// 2点間のベクトルと、点Pを含むカメラ平面と平行な面、との交点を計算
func IntersectLinePlane(near, far, forward, right, up, point *MVec3) *MVec3 {
	// ステップ1: カメラ平面の法線ベクトルを計算
	normal := forward.Cross(right)

	// ステップ2: nearからfarへのベクトルを計算
	direction := far.Subed(near)

	// ステップ3: 平面の方程式のD成分を計算
	D := -normal.Dot(point)

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

// 2点間のベクトルと、点Pとの交点を計算
func IntersectLinePoint(near, far, point *MVec3) *MVec3 {
	// ステップ1: nearからfarへのベクトルを計算
	direction := far.Subed(near)

	// ステップ2: 直線と点の交点を計算
	t := (point.X - near.X) / direction.X
	intersection := near.Added(direction.MuledScalar(t))
	return intersection
}

// DistanceLineToPoints 線分と点の距離を計算します
func DistanceLineToPoints(worldPos *MVec3, points []*MVec3) []float64 {
	distances := make([]float64, len(points))

	// worldPos の Z方向のベクトル
	worldDirection := worldPos.Added(MVec3UnitZNeg)

	for i, p := range points {
		// 点PとworldPosのZ方向のベクトルとの距離を計算
		distances[i] = DistanceFromPointToLine(worldPos, worldDirection, p)
	}

	return distances
}

func (vec3 *MVec3) Project(other *MVec3) *MVec3 {
	return other.MuledScalar(vec3.Dot(other) / other.LengthSqr())
}

// 点が直方体内にあるかどうかを判定する関数
func (vec3 *MVec3) IsPointInsideBox(min, max *MVec3) bool {
	return vec3.X >= min.X && vec3.X <= max.X &&
		vec3.Y >= min.Y && vec3.Y <= max.Y &&
		vec3.Z >= min.Z && vec3.Z <= max.Z
}

// Vec3Diffは、2つのベクトル間の回転四元数を返します。
func (vec1 *MVec3) Vec3Diff(vec2 *MVec3) *MQuaternion {
	cr := vec1.Cross(vec2)
	sr := math.Sqrt(2 * (1 + vec1.Dot(vec2)))
	oosr := 1 / sr

	q := NewMQuaternionByValues(cr.X*oosr, cr.Y*oosr, cr.Z*oosr, sr*0.5)
	return q.Normalize()
}

func (v *MVec3) Round(threshold float64) *MVec3 {
	return &MVec3{
		Round(v.X, threshold),
		Round(v.Y, threshold),
		Round(v.Z, threshold),
	}
}

func SortVec3(vectors []MVec3) []MVec3 {
	// X, Y, Zの順にソート
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

func MeanVec3(vectors []*MVec3) *MVec3 {
	if len(vectors) == 0 {
		return &MVec3{0, 0, 0}
	}

	sum := &MVec3{0, 0, 0}
	for _, v := range vectors {
		sum.Add(v)
	}

	return sum.MuledScalar(1.0 / float64(len(vectors)))
}
