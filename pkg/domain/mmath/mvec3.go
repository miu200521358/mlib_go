// 指示: miu200521358
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
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
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

func (vec3 *MVec3) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f]", vec3.X, vec3.Y, vec3.Z)
}

func (vec3 *MVec3) StringByDigits(digits int) string {
	format := fmt.Sprintf("[x=%%.%df, y=%%.%df, z=%%.%df]", digits, digits, digits)
	return fmt.Sprintf(format, vec3.X, vec3.Y, vec3.Z)
}

func (vec3 *MVec3) MMD() *MVec3 {
	return &MVec3{vec3.X, vec3.Y, vec3.Z}
}

func (vec3 *MVec3) Add(other *MVec3) *MVec3 {
	vec3.X += other.X
	vec3.Y += other.Y
	vec3.Z += other.Z
	return vec3
}

func (vec3 *MVec3) AddScalar(s float64) *MVec3 {
	vec3.X += s
	vec3.Y += s
	vec3.Z += s
	return vec3
}

func (vec3 *MVec3) Added(other *MVec3) *MVec3 {
	return &MVec3{vec3.X + other.X, vec3.Y + other.Y, vec3.Z + other.Z}
}

func (vec3 *MVec3) AddedScalar(s float64) *MVec3 {
	return &MVec3{vec3.X + s, vec3.Y + s, vec3.Z + s}
}

func (vec3 *MVec3) Sub(other *MVec3) *MVec3 {
	vec3.X -= other.X
	vec3.Y -= other.Y
	vec3.Z -= other.Z
	return vec3
}

func (vec3 *MVec3) SubScalar(s float64) *MVec3 {
	vec3.X -= s
	vec3.Y -= s
	vec3.Z -= s
	return vec3
}

func (vec3 *MVec3) Subed(other *MVec3) *MVec3 {
	return &MVec3{vec3.X - other.X, vec3.Y - other.Y, vec3.Z - other.Z}
}

func (vec3 *MVec3) SubedScalar(s float64) *MVec3 {
	return &MVec3{vec3.X - s, vec3.Y - s, vec3.Z - s}
}

func (vec3 *MVec3) Mul(other *MVec3) *MVec3 {
	vec3.X *= other.X
	vec3.Y *= other.Y
	vec3.Z *= other.Z
	return vec3
}

func (vec3 *MVec3) MulScalar(s float64) *MVec3 {
	vec3.X *= s
	vec3.Y *= s
	vec3.Z *= s
	return vec3
}

func (vec3 *MVec3) Muled(other *MVec3) *MVec3 {
	return &MVec3{vec3.X * other.X, vec3.Y * other.Y, vec3.Z * other.Z}
}

func (vec3 *MVec3) MuledScalar(s float64) *MVec3 {
	return &MVec3{vec3.X * s, vec3.Y * s, vec3.Z * s}
}

func (vec3 *MVec3) Div(other *MVec3) *MVec3 {
	vec3.X /= other.X
	vec3.Y /= other.Y
	vec3.Z /= other.Z
	return vec3
}

func (vec3 *MVec3) DivScalar(s float64) *MVec3 {
	vec3.X /= s
	vec3.Y /= s
	vec3.Z /= s
	return vec3
}

func (vec3 *MVec3) Dived(other *MVec3) *MVec3 {
	return &MVec3{vec3.X / other.X, vec3.Y / other.Y, vec3.Z / other.Z}
}

func (vec3 *MVec3) DivedScalar(s float64) *MVec3 {
	return &MVec3{vec3.X / s, vec3.Y / s, vec3.Z / s}
}

func (vec3 *MVec3) Equals(other *MVec3) bool {
	return vec3.X == other.X && vec3.Y == other.Y && vec3.Z == other.Z
}

func (vec3 *MVec3) NotEquals(other MVec3) bool {
	return vec3.X != other.X || vec3.Y != other.Y || vec3.Z != other.Z
}

func (vec3 *MVec3) NearEquals(other *MVec3, epsilon float64) bool {
	return (math.Abs(vec3.X-other.X) <= epsilon) &&
		(math.Abs(vec3.Y-other.Y) <= epsilon) &&
		(math.Abs(vec3.Z-other.Z) <= epsilon)
}

func (vec3 *MVec3) LessThan(other *MVec3) bool {
	return vec3.X < other.X && vec3.Y < other.Y && vec3.Z < other.Z
}

func (vec3 *MVec3) LessThanOrEquals(other *MVec3) bool {
	return vec3.X <= other.X && vec3.Y <= other.Y && vec3.Z <= other.Z
}

func (vec3 *MVec3) GreaterThan(other *MVec3) bool {
	return vec3.X > other.X && vec3.Y > other.Y && vec3.Z > other.Z
}

func (vec3 *MVec3) GreaterThanOrEquals(other *MVec3) bool {
	return vec3.X >= other.X && vec3.Y >= other.Y && vec3.Z >= other.Z
}

func (vec3 *MVec3) Negate() *MVec3 {
	vec3.X = -vec3.X
	vec3.Y = -vec3.Y
	vec3.Z = -vec3.Z
	return vec3
}

func (vec3 *MVec3) Negated() *MVec3 {
	return &MVec3{-vec3.X, -vec3.Y, -vec3.Z}
}

func (vec3 *MVec3) Abs() *MVec3 {
	vec3.X = math.Abs(vec3.X)
	vec3.Y = math.Abs(vec3.Y)
	vec3.Z = math.Abs(vec3.Z)
	return vec3
}

func (vec3 *MVec3) Absed() *MVec3 {
	return &MVec3{math.Abs(vec3.X), math.Abs(vec3.Y), math.Abs(vec3.Z)}
}

func (vec3 *MVec3) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f", vec3.X, vec3.Y, vec3.Z)))
	return h.Sum64()
}

func (vec3 *MVec3) Truncate(epsilon float64) *MVec3 {
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

func (vec3 *MVec3) Truncated(epsilon float64) *MVec3 {
	vec := vec3.Copy()
	vec.Truncate(epsilon)
	return vec
}

func (vec3 *MVec3) MergeIfZero(v float64) *MVec3 {
	if vec3.X == 0 {
		vec3.X = v
	}
	if vec3.Y == 0 {
		vec3.Y = v
	}
	if vec3.Z == 0 {
		vec3.Z = v
	}
	return vec3
}

func (vec3 *MVec3) MergeIfZeros(v *MVec3) *MVec3 {
	if vec3.X == 0 {
		vec3.X = v.X
	}
	if vec3.Y == 0 {
		vec3.Y = v.Y
	}
	if vec3.Z == 0 {
		vec3.Z = v.Z
	}
	return vec3
}

func (vec3 *MVec3) IsZero() bool {
	return vec3 == nil || vec3.NearEquals(MVec3Zero, 1e-10)
}

func (vec3 *MVec3) IsOne() bool {
	return vec3.NearEquals(MVec3One, 1e-10)
}

func (vec3 *MVec3) Length() float64 {
	return mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.Len()
}

func (vec3 *MVec3) LengthSqr() float64 {
	return mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.LenSqr()
}

func (vec3 *MVec3) Normalize() *MVec3 {
	sl := vec3.LengthSqr()
	if sl == 0 || sl == 1 {
		return vec3
	}
	return vec3.MulScalar(1 / math.Sqrt(sl))
}

func (vec3 *MVec3) Normalized() *MVec3 {
	vec := MVec3{vec3.X, vec3.Y, vec3.Z}
	vec.Normalize()
	return &vec
}

func (vec3 *MVec3) Angle(other *MVec3) float64 {
	vec := vec3.Dot(other) / (vec3.Length() * other.Length())
	if vec > 1. {
		return 0
	} else if vec < -1. {
		return math.Pi
	}
	return math.Acos(vec)
}

func (vec3 *MVec3) Degree(other *MVec3) float64 {
	radian := vec3.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

func (vec3 *MVec3) Dot(other *MVec3) float64 {
	return mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.Dot(mgl64.Vec3{other.X, other.Y, other.Z})
}

func (vec3 *MVec3) Cross(other *MVec3) *MVec3 {
	v := mgl64.Vec3{vec3.X, vec3.Y, vec3.Z}.Cross(mgl64.Vec3{other.X, other.Y, other.Z})
	return &MVec3{v[0], v[1], v[2]}
}

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

func (vec3 *MVec3) Clamp(min, max *MVec3) *MVec3 {
	vec3.X = Clamped(vec3.X, min.X, max.X)
	vec3.Y = Clamped(vec3.Y, min.Y, max.Y)
	vec3.Z = Clamped(vec3.Z, min.Z, max.Z)
	return vec3
}

func (vec3 *MVec3) Clamped(min, max *MVec3) *MVec3 {
	result := MVec3{vec3.X, vec3.Y, vec3.Z}
	result.Clamp(min, max)
	return &result
}

func (vec3 *MVec3) Clamp01() *MVec3 {
	return vec3.Clamp(MVec3Zero, MVec3One)
}

func (vec3 *MVec3) Clamped01() *MVec3 {
	result := MVec3{vec3.X, vec3.Y, vec3.Z}
	result.Clamp01()
	return &result
}

func (vec3 *MVec3) Copy() *MVec3 {
	return &MVec3{vec3.X, vec3.Y, vec3.Z}
}

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

		return vec3.Added((other.Subed(vec3)).MuledScalar(t))
	}
}

func (vec3 *MVec3) Slerp(other *MVec3, t float64) *MVec3 {
	switch {
	case t <= 0:
		return vec3
	case t >= 1:
		return other
	default:
		if vec3.NearEquals(other, 1e-8) {
			return vec3.Copy()
		}

		v0 := vec3.Normalized()
		v1 := other.Normalized()

		dot := v0.Dot(v1)

		if dot > 1.0 {
			dot = 1.0
		} else if dot < -1.0 {
			dot = -1.0
		}

		theta := math.Acos(dot)
		sinTheta := math.Sin(theta)

		s0 := math.Sin((1-t)*theta) / sinTheta
		s1 := math.Sin(t*theta) / sinTheta

		result := v0.MuledScalar(s0).Add(v1.MuledScalar(s1))

		return result.MuledScalar(vec3.Length())
	}
}

func (vec3 *MVec3) ToLocalMat() *MMat4 {
	if vec3.IsZero() {
		return NewMMat4()
	}

	v := vec3.Normalized()

	var up *MVec3
	if math.Abs(v.Y) < 1-1e-6 {
		up = MVec3UnitY
	} else {
		up = MVec3UnitZ
	}

	u := up.Cross(v).Normalized()

	w := v.Cross(u).Normalized()

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

	rotationMatrix := vec3.ToLocalMat()

	return rotationMatrix.Muled(scales.ToScaleMat4()).Muled(rotationMatrix.Inverted())
}

func (vec3 *MVec3) One() *MVec3 {
	vec := vec3.Vector()
	epsilon := 1e-3
	for i := range vec {
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

func DistanceFromPointToLine(vec1, vec2, point *MVec3) float64 {
	lineVec := vec2.Subed(vec1)
	pointVec := point.Subed(vec1)
	crossVec := lineVec.Cross(pointVec)
	area := crossVec.Length()
	lineLength := lineVec.Length()
	return area / lineLength
}

func DistanceFromPlaneToLine(near, far, forward, right, up, point *MVec3) float64 {
	normal := forward.Cross(right)

	vectorToPlane := point.Subed(near)

	distance := math.Abs(vectorToPlane.Dot(normal)) / normal.Length()

	return distance
}

func IntersectLinePlane(near, far, forward, right, up, point *MVec3) *MVec3 {
	normal := forward.Cross(right)

	direction := far.Subed(near)

	D := -normal.Dot(point)

	denom := normal.Dot(direction)
	if math.Abs(denom) < 1e-6 {
		return nil
	}

	t := -(normal.Dot(near) + D) / denom
	intersection := near.Added(direction.MuledScalar(t))
	return intersection
}

func IntersectLinePoint(near, far, point *MVec3) *MVec3 {
	direction := far.Subed(near)

	t := (point.X - near.X) / direction.X
	intersection := near.Added(direction.MuledScalar(t))
	return intersection
}

func DistanceLineToPoints(worldPos *MVec3, points []*MVec3) []float64 {
	distances := make([]float64, len(points))

	worldDirection := worldPos.Added(MVec3UnitZNeg)

	for i, p := range points {
		distances[i] = DistanceFromPointToLine(worldPos, worldDirection, p)
	}

	return distances
}

func (vec3 *MVec3) Project(other *MVec3) *MVec3 {
	return other.MuledScalar(vec3.Dot(other) / other.LengthSqr())
}

func (vec3 *MVec3) IsPointInsideBox(min, max *MVec3) bool {
	return vec3.X >= min.X && vec3.X <= max.X &&
		vec3.Y >= min.Y && vec3.Y <= max.Y &&
		vec3.Z >= min.Z && vec3.Z <= max.Z
}

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

func MinVec3(vectors []*MVec3) *MVec3 {
	if len(vectors) == 0 {
		return &MVec3{0, 0, 0}
	}

	min := vectors[0].Copy()
	for _, v := range vectors[1:] {
		min.X = math.Min(min.X, v.X)
		min.Y = math.Min(min.Y, v.Y)
		min.Z = math.Min(min.Z, v.Z)
	}

	return min
}

func MaxVec3(vectors []*MVec3) *MVec3 {
	if len(vectors) == 0 {
		return &MVec3{0, 0, 0}
	}

	max := vectors[0].Copy()
	for _, v := range vectors[1:] {
		max.X = math.Max(max.X, v.X)
		max.Y = math.Max(max.Y, v.Y)
		max.Z = math.Max(max.Z, v.Z)
	}

	return max
}

func MedianVec3(vectors []*MVec3) *MVec3 {
	if len(vectors) == 0 {
		return &MVec3{0, 0, 0}
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

	return &MVec3{
		X: xValues[len(xValues)/2],
		Y: yValues[len(yValues)/2],
		Z: zValues[len(zValues)/2],
	}
}

