package mmath

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"sort"

	"gonum.org/v1/gonum/spatial/r3"
)

type Vec3 struct {
	r3.Vec
}

var (
	ZERO_VEC3       = Vec3{}
	UNIT_X_VEC3     = Vec3{r3.Vec{X: 1}}
	UNIT_Y_VEC3     = Vec3{r3.Vec{Y: 1}}
	UNIT_Z_VEC3     = Vec3{r3.Vec{Z: 1}}
	ONE_VEC3        = Vec3{r3.Vec{X: 1, Y: 1, Z: 1}}
	UNIT_X_NEG_VEC3 = Vec3{r3.Vec{X: -1}}
	UNIT_Y_NEG_VEC3 = Vec3{r3.Vec{Y: -1}}
	UNIT_Z_NEG_VEC3 = Vec3{r3.Vec{Z: -1}}
	ONE_NEG_VEC3    = Vec3{r3.Vec{X: -1, Y: -1, Z: -1}}
	VEC3_MIN_VAL    = Vec3{r3.Vec{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64}}
	VEC3_MAX_VAL    = Vec3{r3.Vec{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}}
)

func NewVec3() Vec3 {
	return Vec3{}
}

func (v Vec3) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}{v.X, v.Y, v.Z})
}

func (v *Vec3) UnmarshalJSON(data []byte) error {
	var tmp struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.X = tmp.X
	v.Y = tmp.Y
	v.Z = tmp.Z
	return nil
}

func (v Vec3) XY() Vec2 {
	return Vec2{v.X, v.Y}
}

func (v Vec3) IsOnlyX() bool {
	return !NearEquals(v.X, 0, 1e-10) && NearEquals(v.Y, 0, 1e-10) && NearEquals(v.Z, 0, 1e-10)
}

func (v Vec3) IsOnlyY() bool {
	return NearEquals(v.X, 0, 1e-10) && !NearEquals(v.Y, 0, 1e-10) && NearEquals(v.Z, 0, 1e-10)
}

func (v Vec3) IsOnlyZ() bool {
	return NearEquals(v.X, 0, 1e-10) && NearEquals(v.Y, 0, 1e-10) && !NearEquals(v.Z, 0, 1e-10)
}

func (v Vec3) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f]", v.X, v.Y, v.Z)
}

func (v Vec3) StringByDigits(digits int) string {
	format := fmt.Sprintf("[x=%%.%df, y=%%.%df, z=%%.%df]", digits, digits, digits)
	return fmt.Sprintf(format, v.X, v.Y, v.Z)
}

func (v Vec3) MMD() Vec3 {
	return v
}

func (v *Vec3) Add(other Vec3) *Vec3 {
	v.Vec = r3.Add(v.Vec, other.Vec)
	return v
}

func (v *Vec3) AddScalar(s float64) *Vec3 {
	v.X += s
	v.Y += s
	v.Z += s
	return v
}

func (v Vec3) Added(other Vec3) Vec3 {
	return Vec3{r3.Add(v.Vec, other.Vec)}
}

func (v Vec3) AddedScalar(s float64) Vec3 {
	return Vec3{r3.Vec{X: v.X + s, Y: v.Y + s, Z: v.Z + s}}
}

func (v *Vec3) Sub(other Vec3) *Vec3 {
	v.Vec = r3.Sub(v.Vec, other.Vec)
	return v
}

func (v *Vec3) SubScalar(s float64) *Vec3 {
	v.X -= s
	v.Y -= s
	v.Z -= s
	return v
}

func (v Vec3) Subed(other Vec3) Vec3 {
	return Vec3{r3.Sub(v.Vec, other.Vec)}
}

func (v Vec3) SubedScalar(s float64) Vec3 {
	return Vec3{r3.Vec{X: v.X - s, Y: v.Y - s, Z: v.Z - s}}
}

func (v *Vec3) Mul(other Vec3) *Vec3 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

func (v *Vec3) MulScalar(s float64) *Vec3 {
	v.Vec = r3.Scale(s, v.Vec)
	return v
}

func (v Vec3) Muled(other Vec3) Vec3 {
	return Vec3{r3.Vec{X: v.X * other.X, Y: v.Y * other.Y, Z: v.Z * other.Z}}
}

func (v Vec3) MuledScalar(s float64) Vec3 {
	return Vec3{r3.Scale(s, v.Vec)}
}

func (v *Vec3) Div(other Vec3) *Vec3 {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	return v
}

func (v *Vec3) DivScalar(s float64) *Vec3 {
	v.X /= s
	v.Y /= s
	v.Z /= s
	return v
}

func (v Vec3) Dived(other Vec3) Vec3 {
	return Vec3{r3.Vec{X: v.X / other.X, Y: v.Y / other.Y, Z: v.Z / other.Z}}
}

func (v Vec3) DivedScalar(s float64) Vec3 {
	return Vec3{r3.Vec{X: v.X / s, Y: v.Y / s, Z: v.Z / s}}
}

func (v Vec3) Equals(other Vec3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

func (v Vec3) NotEquals(other Vec3) bool {
	return v.X != other.X || v.Y != other.Y || v.Z != other.Z
}

func (v Vec3) NearEquals(other Vec3, epsilon float64) bool {
	return math.Abs(v.X-other.X) <= epsilon && math.Abs(v.Y-other.Y) <= epsilon && math.Abs(v.Z-other.Z) <= epsilon
}

func (v Vec3) LessThan(other Vec3) bool {
	return v.X < other.X && v.Y < other.Y && v.Z < other.Z
}

func (v Vec3) LessThanOrEquals(other Vec3) bool {
	return v.X <= other.X && v.Y <= other.Y && v.Z <= other.Z
}

func (v Vec3) GreaterThan(other Vec3) bool {
	return v.X > other.X && v.Y > other.Y && v.Z > other.Z
}

func (v Vec3) GreaterThanOrEquals(other Vec3) bool {
	return v.X >= other.X && v.Y >= other.Y && v.Z >= other.Z
}

func (v *Vec3) Negate() *Vec3 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

func (v Vec3) Negated() Vec3 {
	return Vec3{r3.Vec{X: -v.X, Y: -v.Y, Z: -v.Z}}
}

func (v *Vec3) Abs() *Vec3 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	v.Z = math.Abs(v.Z)
	return v
}

func (v Vec3) Absed() Vec3 {
	return Vec3{r3.Vec{X: math.Abs(v.X), Y: math.Abs(v.Y), Z: math.Abs(v.Z)}}
}

func (v Vec3) Hash() uint64 {
	h := fnv.New64a()
	_, _ = fmt.Fprintf(h, "%.10f,%.10f,%.10f", v.X, v.Y, v.Z)
	return h.Sum64()
}

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

func (v Vec3) Truncated(epsilon float64) Vec3 {
	vec := v
	vec.Truncate(epsilon)
	return vec
}

func (v *Vec3) MergeIfZero(val float64) *Vec3 {
	if v.X == 0 {
		v.X = val
	}
	if v.Y == 0 {
		v.Y = val
	}
	if v.Z == 0 {
		v.Z = val
	}
	return v
}

func (v *Vec3) MergeIfZeros(other Vec3) *Vec3 {
	if v.X == 0 {
		v.X = other.X
	}
	if v.Y == 0 {
		v.Y = other.Y
	}
	if v.Z == 0 {
		v.Z = other.Z
	}
	return v
}

func (v Vec3) IsZero() bool {
	return v.NearEquals(ZERO_VEC3, 1e-10)
}

func (v Vec3) IsOne() bool {
	return v.NearEquals(ONE_VEC3, 1e-10)
}

func (v Vec3) Length() float64 {
	return r3.Norm(v.Vec)
}

func (v Vec3) LengthSqr() float64 {
	return r3.Norm2(v.Vec)
}

func (v *Vec3) Normalize() *Vec3 {
	sl := v.LengthSqr()
	if sl == 0 || sl == 1 {
		return v
	}
	v.Vec = r3.Scale(1/math.Sqrt(sl), v.Vec)
	return v
}

func (v Vec3) Normalized() Vec3 {
	vec := v
	vec.Normalize()
	return vec
}

func (v Vec3) Angle(other Vec3) float64 {
	denom := v.Length() * other.Length()
	if denom == 0 {
		return 0
	}
	return angleFromCosVec3(v.Dot(other) / denom)
}

func (v Vec3) Degree(other Vec3) float64 {
	return RadToDeg(v.Angle(other))
}

func (v Vec3) Dot(other Vec3) float64 {
	return r3.Dot(v.Vec, other.Vec)
}

func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{r3.Cross(v.Vec, other.Vec)}
}

func (v Vec3) Min() Vec3 {
	min := v.X
	if v.Y < min {
		min = v.Y
	}
	if v.Z < min {
		min = v.Z
	}
	return Vec3{r3.Vec{X: min, Y: min, Z: min}}
}

func (v Vec3) Max() Vec3 {
	max := v.X
	if v.Y > max {
		max = v.Y
	}
	if v.Z > max {
		max = v.Z
	}
	return Vec3{r3.Vec{X: max, Y: max, Z: max}}
}

func (v *Vec3) Clamp(min, max Vec3) *Vec3 {
	v.X = Clamped(v.X, min.X, max.X)
	v.Y = Clamped(v.Y, min.Y, max.Y)
	v.Z = Clamped(v.Z, min.Z, max.Z)
	return v
}

func (v Vec3) Clamped(min, max Vec3) Vec3 {
	result := v
	result.Clamp(min, max)
	return result
}

func (v *Vec3) Clamp01() *Vec3 {
	return v.Clamp(ZERO_VEC3, ONE_VEC3)
}

func (v Vec3) Clamped01() Vec3 {
	result := v
	result.Clamp01()
	return result
}

func (v Vec3) Copy() (*Vec3, error) {
	return deepCopy(v)
}

func (v Vec3) Vector() []float64 {
	return []float64{v.X, v.Y, v.Z}
}

func (v Vec3) ToMat4() Mat4 {
	mat := NewMat4()
	mat[12] = v.X
	mat[13] = v.Y
	mat[14] = v.Z
	return mat
}

func (v Vec3) ToScaleMat4() Mat4 {
	mat := NewMat4()
	mat[0] = v.X
	mat[5] = v.Y
	mat[10] = v.Z
	return mat
}

func (v *Vec3) ClampIfVerySmall() *Vec3 {
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

func (v Vec3) RadToDeg() Vec3 {
	return Vec3{r3.Vec{X: RadToDeg(v.X), Y: RadToDeg(v.Y), Z: RadToDeg(v.Z)}}
}

func (v Vec3) DegToRad() Vec3 {
	return Vec3{r3.Vec{X: DegToRad(v.X), Y: DegToRad(v.Y), Z: DegToRad(v.Z)}}
}

func (v Vec3) RadToQuaternion() Quaternion {
	return NewQuaternionFromRadians(v.X, v.Y, v.Z)
}

func (v Vec3) DegToQuaternion() Quaternion {
	return NewQuaternionFromDegrees(v.X, v.Y, v.Z)
}

func (v Vec3) Lerp(other Vec3, t float64) Vec3 {
	if t <= 0 {
		return v
	}
	if t >= 1 {
		return other
	}
	if v.NearEquals(other, 1e-8) {
		return v
	}
	return v.Added(other.Subed(v).MuledScalar(t))
}

func (v Vec3) Slerp(other Vec3, t float64) Vec3 {
	if t <= 0 {
		return v
	}
	if t >= 1 {
		return other
	}
	if v.NearEquals(other, 1e-8) {
		return v
	}

	v0 := v.Normalized()
	v1 := other.Normalized()
	dot := v0.Dot(v1)
	dot = math.Max(-1.0, math.Min(1.0, dot))
	theta := math.Acos(dot)
	sinTheta := math.Sin(theta)
	s0 := math.Sin((1-t)*theta) / sinTheta
	s1 := math.Sin(t*theta) / sinTheta
	result := v0.MuledScalar(s0).Added(v1.MuledScalar(s1))
	return result.MuledScalar(v.Length())
}

func (v Vec3) ToLocalMat() Mat4 {
	if v.IsZero() {
		return NewMat4()
	}

	vv := v.Normalized()
	var up Vec3
	if math.Abs(vv.Y) < 1-1e-6 {
		up = UNIT_Y_VEC3
	} else {
		up = UNIT_Z_VEC3
	}

	u := up.Cross(vv).Normalized()
	w := vv.Cross(u).Normalized()

	return NewMat4ByValues(
		vv.X, vv.Y, vv.Z, 0,
		w.X, w.Y, w.Z, 0,
		u.X, u.Y, u.Z, 0,
		0, 0, 0, 1,
	)
}

func (v Vec3) ToScaleLocalMat(scales Vec3) Mat4 {
	if v.IsZero() || v.IsOne() {
		return NewMat4()
	}

	rotationMatrix := v.ToLocalMat()
	return rotationMatrix.Muled(scales.ToScaleMat4()).Muled(rotationMatrix.Inverted())
}

func (v Vec3) One() Vec3 {
	vec := v.Vector()
	epsilon := 1e-3
	for i := range vec {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return Vec3{r3.Vec{X: vec[0], Y: vec[1], Z: vec[2]}}
}

func (v Vec3) Distance(other Vec3) float64 {
	return v.Subed(other).Length()
}

func (v Vec3) Distances(others []Vec3) []float64 {
	distances := make([]float64, len(others))
	for i, other := range others {
		distances[i] = v.Distance(other)
	}
	return distances
}

func (v *Vec3) Effective() *Vec3 {
	v.X = Effective(v.X)
	v.Y = Effective(v.Y)
	v.Z = Effective(v.Z)
	return v
}

func DistanceFromPointToLine(vec1, vec2, point Vec3) float64 {
	lineVec := vec2.Subed(vec1)
	pointVec := point.Subed(vec1)
	crossVec := lineVec.Cross(pointVec)
	area := crossVec.Length()
	lineLength := lineVec.Length()
	return area / lineLength
}

func DistanceFromPlaneToLine(near, far, forward, right, up, point Vec3) float64 {
	_ = up
	normal := forward.Cross(right)
	vectorToPlane := point.Subed(near)
	distance := math.Abs(vectorToPlane.Dot(normal)) / normal.Length()
	return distance
}

func IntersectLinePlane(near, far, forward, right, up, point Vec3) (Vec3, error) {
	_ = up
	normal := forward.Cross(right)
	direction := far.Subed(near)
	D := -normal.Dot(point)
	denom := normal.Dot(direction)
	if math.Abs(denom) < 1e-6 {
		return Vec3{}, fmt.Errorf("line and plane are parallel")
	}
	t := -(normal.Dot(near) + D) / denom
	intersection := near.Added(direction.MuledScalar(t))
	return intersection, nil
}

func IntersectLinePoint(near, far, point Vec3) Vec3 {
	direction := far.Subed(near)
	t := (point.X - near.X) / direction.X
	intersection := near.Added(direction.MuledScalar(t))
	return intersection
}

func DistanceLineToPoints(worldPos Vec3, points []Vec3) []float64 {
	distances := make([]float64, len(points))
	worldDirection := worldPos.Added(UNIT_Z_NEG_VEC3)
	for i, p := range points {
		distances[i] = DistanceFromPointToLine(worldPos, worldDirection, p)
	}
	return distances
}

func (v Vec3) Project(other Vec3) Vec3 {
	return other.MuledScalar(v.Dot(other) / other.LengthSqr())
}

func (v Vec3) IsPointInsideBox(min, max Vec3) bool {
	return v.X >= min.X && v.X <= max.X && v.Y >= min.Y && v.Y <= max.Y && v.Z >= min.Z && v.Z <= max.Z
}

func (v Vec3) Vec3Diff(other Vec3) Quaternion {
	cr := v.Cross(other)
	sr := math.Sqrt(2 * (1 + v.Dot(other)))
	oosr := 1 / sr
	q := NewQuaternionByValues(cr.X*oosr, cr.Y*oosr, cr.Z*oosr, sr*0.5)
	return q.Normalized()
}

func (v Vec3) Round(threshold float64) Vec3 {
	return Vec3{r3.Vec{X: Round(v.X, threshold), Y: Round(v.Y, threshold), Z: Round(v.Z, threshold)}}
}

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

func MeanVec3(vectors []Vec3) Vec3 {
	if len(vectors) == 0 {
		return Vec3{}
	}
	sum := Vec3{}
	for _, v := range vectors {
		sum.Add(v)
	}
	return sum.MuledScalar(1.0 / float64(len(vectors)))
}

func MinVec3(vectors []Vec3) Vec3 {
	if len(vectors) == 0 {
		return Vec3{}
	}
	min := vectors[0]
	for _, v := range vectors[1:] {
		min.X = math.Min(min.X, v.X)
		min.Y = math.Min(min.Y, v.Y)
		min.Z = math.Min(min.Z, v.Z)
	}
	return min
}

func MaxVec3(vectors []Vec3) Vec3 {
	if len(vectors) == 0 {
		return Vec3{}
	}
	max := vectors[0]
	for _, v := range vectors[1:] {
		max.X = math.Max(max.X, v.X)
		max.Y = math.Max(max.Y, v.Y)
		max.Z = math.Max(max.Z, v.Z)
	}
	return max
}

func MedianVec3(vectors []Vec3) Vec3 {
	if len(vectors) == 0 {
		return Vec3{}
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
	return Vec3{r3.Vec{X: xValues[len(xValues)/2], Y: yValues[len(yValues)/2], Z: zValues[len(zValues)/2]}}
}

func angleFromCosVec3(val float64) float64 {
	if val > 1 {
		return 0
	}
	if val < -1 {
		return math.Pi
	}
	return math.Acos(val)
}
