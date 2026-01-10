package mmath

import (
	"fmt"
	"hash/fnv"
	"math"

	"gonum.org/v1/gonum/spatial/r3"
)

type Vec4 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

var (
	ZERO_VEC4     = Vec4{}
	UNIT_XW_VEC4  = Vec4{X: 1, W: 1}
	UNIT_YW_VEC4  = Vec4{Y: 1, W: 1}
	UNIT_ZW_VEC4  = Vec4{Z: 1, W: 1}
	UNIT_W_VEC4   = Vec4{W: 1}
	ONE_VEC4      = Vec4{X: 1, Y: 1, Z: 1, W: 1}
	VEC4_MIN_VAL  = Vec4{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64, W: 1}
	VEC4_MAX_VAL  = Vec4{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64, W: 1}
)

func NewVec4() Vec4 {
	return Vec4{}
}

func (v Vec4) XY() Vec2 {
	return Vec2{v.X, v.Y}
}

func (v Vec4) XYZ() Vec3 {
	return Vec3{r3.Vec{X: v.X, Y: v.Y, Z: v.Z}}
}

func (v Vec4) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", v.X, v.Y, v.Z, v.W)
}

func (v Vec4) MMD() Vec4 {
	return v
}

func (v *Vec4) Add(other Vec4) *Vec4 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	v.W += other.W
	return v
}

func (v *Vec4) AddScalar(s float64) *Vec4 {
	v.X += s
	v.Y += s
	v.Z += s
	v.W += s
	return v
}

func (v Vec4) Added(other Vec4) Vec4 {
	return Vec4{v.X + other.X, v.Y + other.Y, v.Z + other.Z, v.W + other.W}
}

func (v Vec4) AddedScalar(s float64) Vec4 {
	return Vec4{v.X + s, v.Y + s, v.Z + s, v.W + s}
}

func (v *Vec4) Sub(other Vec4) *Vec4 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	v.W -= other.W
	return v
}

func (v *Vec4) SubScalar(s float64) *Vec4 {
	v.X -= s
	v.Y -= s
	v.Z -= s
	v.W -= s
	return v
}

func (v Vec4) Subed(other Vec4) Vec4 {
	return Vec4{v.X - other.X, v.Y - other.Y, v.Z - other.Z, v.W - other.W}
}

func (v Vec4) SubedScalar(s float64) Vec4 {
	return Vec4{v.X - s, v.Y - s, v.Z - s, v.W - s}
}

func (v *Vec4) Mul(other Vec4) *Vec4 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	v.W *= other.W
	return v
}

func (v *Vec4) MulScalar(s float64) *Vec4 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	v.W *= s
	return v
}

func (v Vec4) Muled(other Vec4) Vec4 {
	return Vec4{v.X * other.X, v.Y * other.Y, v.Z * other.Z, v.W * other.W}
}

func (v Vec4) MuledScalar(s float64) Vec4 {
	return Vec4{v.X * s, v.Y * s, v.Z * s, v.W * s}
}

func (v *Vec4) Div(other Vec4) *Vec4 {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	v.W /= other.W
	return v
}

func (v *Vec4) DivScalar(s float64) *Vec4 {
	v.X /= s
	v.Y /= s
	v.Z /= s
	v.W /= s
	return v
}

func (v Vec4) Dived(other Vec4) Vec4 {
	return Vec4{v.X / other.X, v.Y / other.Y, v.Z / other.Z, v.W / other.W}
}

func (v Vec4) DivedScalar(s float64) Vec4 {
	return Vec4{v.X / s, v.Y / s, v.Z / s, v.W / s}
}

func (v Vec4) Equals(other Vec4) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z && v.W == other.W
}

func (v Vec4) NotEquals(other Vec4) bool {
	return v.X != other.X || v.Y != other.Y || v.Z != other.Z || v.W != other.W
}

func (v Vec4) NearEquals(other Vec4, epsilon float64) bool {
	return math.Abs(v.X-other.X) <= epsilon && math.Abs(v.Y-other.Y) <= epsilon && math.Abs(v.Z-other.Z) <= epsilon && math.Abs(v.W-other.W) <= epsilon
}

func (v Vec4) LessThan(other Vec4) bool {
	return v.X < other.X && v.Y < other.Y && v.Z < other.Z && v.W < other.W
}

func (v Vec4) LessThanOrEquals(other Vec4) bool {
	return v.X <= other.X && v.Y <= other.Y && v.Z <= other.Z && v.W <= other.W
}

func (v Vec4) GreaterThan(other Vec4) bool {
	return v.X > other.X && v.Y > other.Y && v.Z > other.Z && v.W > other.W
}

func (v Vec4) GreaterThanOrEquals(other Vec4) bool {
	return v.X >= other.X && v.Y >= other.Y && v.Z >= other.Z && v.W >= other.W
}

func (v *Vec4) Negate() *Vec4 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	v.W = -v.W
	return v
}

func (v Vec4) Negated() Vec4 {
	return Vec4{-v.X, -v.Y, -v.Z, -v.W}
}

func (v *Vec4) Abs() *Vec4 {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	v.Z = math.Abs(v.Z)
	v.W = math.Abs(v.W)
	return v
}

func (v Vec4) Absed() Vec4 {
	return Vec4{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z), math.Abs(v.W)}
}

func (v Vec4) Hash() uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f,%.10f", v.X, v.Y, v.Z, v.W)))
	return h.Sum64()
}

func (v Vec4) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0 && v.W == 0
}

func (v Vec4) Length() float64 {
	return v.Vec3DividedByW().Length()
}

func (v Vec4) LengthSqr() float64 {
	return v.Vec3DividedByW().LengthSqr()
}

func (v *Vec4) Normalize() *Vec4 {
	v3 := v.Vec3DividedByW()
	v3.Normalize()
	v.X = v3.X
	v.Y = v3.Y
	v.Z = v3.Z
	v.W = 1
	return v
}

func (v Vec4) Normalized() Vec4 {
	vec := v
	vec.Normalize()
	return vec
}

func (v Vec4) Dot(other Vec4) float64 {
	a3 := v.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	return a3.Dot(b3)
}

func Dot4(a, b Vec4) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

func (v Vec4) Cross(other Vec4) Vec4 {
	a3 := v.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	c3 := a3.Cross(b3)
	return Vec4{c3.X, c3.Y, c3.Z, 1}
}

func (v Vec4) Min() Vec4 {
	min := v.X
	if v.Y < min {
		min = v.Y
	}
	if v.Z < min {
		min = v.Z
	}
	if v.W < min {
		min = v.W
	}
	return Vec4{min, min, min, min}
}

func (v Vec4) Max() Vec4 {
	max := v.X
	if v.Y > max {
		max = v.Y
	}
	if v.Z > max {
		max = v.Z
	}
	if v.W > max {
		max = v.W
	}
	return Vec4{max, max, max, max}
}

func (v *Vec4) Clamp(min, max Vec4) *Vec4 {
	v.X = Clamped(v.X, min.X, max.X)
	v.Y = Clamped(v.Y, min.Y, max.Y)
	v.Z = Clamped(v.Z, min.Z, max.Z)
	v.W = Clamped(v.W, min.W, max.W)
	return v
}

func (v Vec4) Clamped(min, max Vec4) Vec4 {
	result := v
	result.Clamp(min, max)
	return result
}

func (v *Vec4) Clamp01() *Vec4 {
	return v.Clamp(ZERO_VEC4, ONE_VEC4)
}

func (v Vec4) Clamped01() Vec4 {
	result := v
	result.Clamp01()
	return result
}

func (v Vec4) Copy() (*Vec4, error) {
	return deepCopy(v)
}

func (v Vec4) Vector() []float64 {
	return []float64{v.X, v.Y, v.Z, v.W}
}

func (v Vec4) Lerp(other Vec4, t float64) Vec4 {
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

func (v Vec4) Vec3DividedByW() Vec3 {
	oow := 1 / v.W
	return Vec3{r3.Vec{X: v.X * oow, Y: v.Y * oow, Z: v.Z * oow}}
}

func (v Vec4) DividedByW() Vec4 {
	oow := 1 / v.W
	return Vec4{v.X * oow, v.Y * oow, v.Z * oow, 1}
}

func (v *Vec4) DivideByW() *Vec4 {
	oow := 1 / v.W
	v.X *= oow
	v.Y *= oow
	v.Z *= oow
	v.W = 1
	return v
}

func (v Vec4) One() Vec4 {
	vec := v.Vector()
	epsilon := 1e-14
	for i := range vec {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return Vec4{vec[0], vec[1], vec[2], vec[3]}
}

func (v Vec4) Distance(other Vec4) float64 {
	return v.Vec3DividedByW().Distance(other.Vec3DividedByW())
}

func (v *Vec4) ClampIfVerySmall() *Vec4 {
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
	if math.Abs(v.W) < epsilon {
		v.W = 0
	}
	return v
}

func (v Vec4) Round(threshold float64) Vec4 {
	return Vec4{
		Round(v.X, threshold),
		Round(v.Y, threshold),
		Round(v.Z, threshold),
		Round(v.W, threshold),
	}
}
