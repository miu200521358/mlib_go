// 指示: miu200521358
package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
)

var (
	MVec4Zero = &MVec4{}

	MVec4UnitXW = &MVec4{1, 0, 0, 1}
	MVec4UnitYW = &MVec4{0, 1, 0, 1}
	MVec4UnitZW = &MVec4{0, 0, 1, 1}
	MVec4UnitW  = &MVec4{0, 0, 0, 1}
	MVec4One    = &MVec4{1, 1, 1, 1}

	MVec4MinVal = &MVec4{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64, 1}
	MVec4MaxVal = &MVec4{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64, 1}
)

type MVec4 struct {
	X float64
	Y float64
	Z float64
	W float64
}

func NewMVec4() *MVec4 {
	return &MVec4{}
}

func (vec4 *MVec4) XY() *MVec2 {
	return &MVec2{vec4.X, vec4.Y}
}

func (vec4 *MVec4) XYZ() *MVec3 {
	return &MVec3{vec4.X, vec4.Y, vec4.Z}
}

func (vec4 *MVec4) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", vec4.X, vec4.Y, vec4.Z, vec4.W)
}

func (vec4 *MVec4) MMD() *MVec4 {
	return &MVec4{vec4.X, vec4.Y, vec4.Z, vec4.W}
}

func (vec4 *MVec4) Add(other *MVec4) *MVec4 {
	vec4.X += other.X
	vec4.Y += other.Y
	vec4.Z += other.Z
	vec4.W += other.W
	return vec4
}

func (vec4 *MVec4) AddScalar(s float64) *MVec4 {
	vec4.X += s
	vec4.Y += s
	vec4.Z += s
	vec4.W += s
	return vec4
}

func (vec4 *MVec4) Added(other *MVec4) *MVec4 {
	return &MVec4{vec4.X + other.X, vec4.Y + other.Y, vec4.Z + other.Z, vec4.W + other.W}
}

func (vec4 *MVec4) AddedScalar(s float64) *MVec4 {
	return &MVec4{vec4.X + s, vec4.Y + s, vec4.Z + s, vec4.W + s}
}

func (vec4 *MVec4) Sub(other *MVec4) *MVec4 {
	vec4.X -= other.X
	vec4.Y -= other.Y
	vec4.Z -= other.Z
	vec4.W -= other.W
	return vec4
}

func (vec4 *MVec4) SubScalar(s float64) *MVec4 {
	vec4.X -= s
	vec4.Y -= s
	vec4.Z -= s
	vec4.W -= s
	return vec4
}

func (vec4 *MVec4) Subed(other *MVec4) *MVec4 {
	return &MVec4{vec4.X - other.X, vec4.Y - other.Y, vec4.Z - other.Z, vec4.W - other.W}
}

func (vec4 *MVec4) SubedScalar(s float64) *MVec4 {
	return &MVec4{vec4.X - s, vec4.Y - s, vec4.Z - s, vec4.W - s}
}

func (vec4 *MVec4) Mul(other *MVec4) *MVec4 {
	vec4.X *= other.X
	vec4.Y *= other.Y
	vec4.Z *= other.Z
	vec4.W *= other.W
	return vec4
}

func (vec4 *MVec4) MulScalar(s float64) *MVec4 {
	vec4.X *= s
	vec4.Y *= s
	vec4.Z *= s
	vec4.W *= s
	return vec4
}

func (vec4 *MVec4) Muled(other *MVec4) *MVec4 {
	return &MVec4{vec4.X * other.X, vec4.Y * other.Y, vec4.Z * other.Z, vec4.W * other.W}
}

func (vec4 *MVec4) MuledScalar(s float64) *MVec4 {
	return &MVec4{vec4.X * s, vec4.Y * s, vec4.Z * s, vec4.W * s}
}

func (vec4 *MVec4) Div(other *MVec4) *MVec4 {
	vec4.X /= other.X
	vec4.Y /= other.Y
	vec4.Z /= other.Z
	vec4.W /= other.W
	return vec4
}

func (vec4 *MVec4) DivScalar(s float64) *MVec4 {
	vec4.X /= s
	vec4.Y /= s
	vec4.Z /= s
	vec4.W /= s
	return vec4
}

func (vec4 *MVec4) Dived(other *MVec4) *MVec4 {
	return &MVec4{vec4.X / other.X, vec4.Y / other.Y, vec4.Z / other.Z, vec4.W / other.W}
}

func (vec4 *MVec4) DivedScalar(s float64) *MVec4 {
	return &MVec4{vec4.X / s, vec4.Y / s, vec4.Z / s, vec4.W / s}
}

func (vec4 *MVec4) Equals(other *MVec4) bool {
	return vec4.X == other.X && vec4.Y == other.Y && vec4.Z == other.Z && vec4.W == other.W
}

func (vec4 *MVec4) NotEquals(other MVec4) bool {
	return vec4.X != other.X || vec4.Y != other.Y || vec4.Z != other.Z || vec4.W != other.W
}

func (vec4 *MVec4) NearEquals(other *MVec4, epsilon float64) bool {
	return (math.Abs(vec4.X-other.X) <= epsilon) &&
		(math.Abs(vec4.Y-other.Y) <= epsilon) &&
		(math.Abs(vec4.Z-other.Z) <= epsilon) &&
		(math.Abs(vec4.W-other.W) <= epsilon)
}

func (vec4 *MVec4) LessThan(other *MVec4) bool {
	return vec4.X < other.X && vec4.Y < other.Y && vec4.Z < other.Z && vec4.W < other.W
}

func (vec4 *MVec4) LessThanOrEquals(other *MVec4) bool {
	return vec4.X <= other.X && vec4.Y <= other.Y && vec4.Z <= other.Z && vec4.W <= other.W
}

func (vec4 *MVec4) GreaterThan(other *MVec4) bool {
	return vec4.X > other.X && vec4.Y > other.Y && vec4.Z > other.Z && vec4.W > other.W
}

func (vec4 *MVec4) GreaterThanOrEquals(other *MVec4) bool {
	return vec4.X >= other.X && vec4.Y >= other.Y && vec4.Z >= other.Z && vec4.W >= other.W
}

func (vec4 *MVec4) Negate() *MVec4 {
	vec4.X = -vec4.X
	vec4.Y = -vec4.Y
	vec4.Z = -vec4.Z
	vec4.W = -vec4.W
	return vec4
}

func (vec4 *MVec4) Negated() *MVec4 {
	return &MVec4{-vec4.X, -vec4.Y, -vec4.Z, -vec4.W}
}

func (vec4 *MVec4) Abs() *MVec4 {
	vec4.X = math.Abs(vec4.X)
	vec4.Y = math.Abs(vec4.Y)
	vec4.Z = math.Abs(vec4.Z)
	vec4.W = math.Abs(vec4.W)
	return vec4
}

func (vec4 *MVec4) Absed() *MVec4 {
	return &MVec4{math.Abs(vec4.X), math.Abs(vec4.Y), math.Abs(vec4.Z), math.Abs(vec4.W)}
}

func (vec4 *MVec4) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f,%.10f,%.10f", vec4.X, vec4.Y, vec4.Z, vec4.W)))
	return h.Sum64()
}

func (vec4 *MVec4) IsZero() bool {
	return vec4.X == 0 && vec4.Y == 0 && vec4.Z == 0 && vec4.W == 0
}

func (vec4 *MVec4) Length() float64 {
	v3 := vec4.Vec3DividedByW()
	return v3.Length()
}

func (vec4 *MVec4) LengthSqr() float64 {
	v3 := vec4.Vec3DividedByW()
	return v3.LengthSqr()
}

func (vec4 *MVec4) Normalize() *MVec4 {
	v3 := vec4.Vec3DividedByW()
	v3.Normalize()
	vec4.X = v3.X
	vec4.Y = v3.Y
	vec4.Z = v3.Z
	vec4.W = 1
	return vec4
}

func (vec4 *MVec4) Normalized() *MVec4 {
	vec := *vec4
	vec.Normalize()
	return &vec
}

func (vec4 *MVec4) Dot(other *MVec4) float64 {
	a3 := vec4.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	return a3.Dot(b3)
}

func Dot4(vec1, vec2 *MVec4) float64 {
	return vec1.X*vec2.X + vec1.Y*vec2.Y + vec1.Z*vec2.Z + vec1.W*vec2.W
}

func (vec4 *MVec4) Cross(other *MVec4) *MVec4 {
	a3 := vec4.Vec3DividedByW()
	b3 := other.Vec3DividedByW()
	c3 := a3.Cross(b3)
	return &MVec4{c3.X, c3.Y, c3.Z, 1}
}

func (vec4 *MVec4) Min() *MVec4 {
	min := vec4.X
	if vec4.Y < min {
		min = vec4.Y
	}
	if vec4.Z < min {
		min = vec4.Z
	}
	if vec4.W < min {
		min = vec4.W
	}
	return &MVec4{min, min, min, min}
}

func (vec4 *MVec4) Max() *MVec4 {
	max := vec4.X
	if vec4.Y > max {
		max = vec4.Y
	}
	if vec4.Z > max {
		max = vec4.Z
	}
	if vec4.W > max {
		max = vec4.W
	}
	return &MVec4{max, max, max, max}
}

func (vec4 *MVec4) Clamp(min, max *MVec4) *MVec4 {
	vec4.X = Clamped(vec4.X, min.X, max.X)
	vec4.Y = Clamped(vec4.Y, min.Y, max.Y)
	vec4.Z = Clamped(vec4.Z, min.Z, max.Z)
	vec4.W = Clamped(vec4.W, min.W, max.W)

	return vec4
}

func (vec4 *MVec4) Clamped(min, max *MVec4) *MVec4 {
	result := *vec4
	result.Clamp(min, max)
	return &result
}

func (vec4 *MVec4) Clamp01() *MVec4 {
	return vec4.Clamp(MVec4Zero, MVec4One)
}

func (vec4 *MVec4) Clamped01() *MVec4 {
	result := *vec4
	result.Clamp01()
	return &result
}

func (vec4 *MVec4) Copy() *MVec4 {
	copied := MVec4{vec4.X, vec4.Y, vec4.Z, vec4.W}
	return &copied
}

func (vec4 *MVec4) Vector() []float64 {
	return []float64{vec4.X, vec4.Y, vec4.Z, vec4.W}
}

func (vec4 *MVec4) Lerp(other *MVec4, t float64) *MVec4 {
	if t <= 0 {
		return vec4.Copy()
	} else if t >= 1 {
		return other.Copy()
	}

	if vec4.Equals(other) {
		return vec4.Copy()
	}

	return (other.Subed(vec4)).MuledScalar(t).Added(vec4)
}

func (vec4 *MVec4) Vec3DividedByW() *MVec3 {
	oow := 1 / vec4.W
	return &MVec3{vec4.X * oow, vec4.Y * oow, vec4.Z * oow}
}

func (vec4 *MVec4) DividedByW() *MVec4 {
	oow := 1 / vec4.W
	return &MVec4{vec4.X * oow, vec4.Y * oow, vec4.Z * oow, 1}
}

func (vec4 *MVec4) DivideByW() *MVec4 {
	oow := 1 / vec4.W
	vec4.X *= oow
	vec4.Y *= oow
	vec4.Z *= oow
	vec4.W = 1
	return vec4
}

func (vec4 *MVec4) One() *MVec4 {
	vec := vec4.Vector()
	epsilon := 1e-14
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec4{vec[0], vec[1], vec[2], vec[3]}
}

func (vec4 *MVec4) Distance(other *MVec4) float64 {
	s := vec4.Subed(other)
	return s.Length()
}

func (vec4 *MVec4) ClampIfVerySmall() *MVec4 {
	epsilon := 1e-6
	if math.Abs(vec4.X) < epsilon {
		vec4.X = 0
	}
	if math.Abs(vec4.Y) < epsilon {
		vec4.Y = 0
	}
	if math.Abs(vec4.Z) < epsilon {
		vec4.Z = 0
	}
	if math.Abs(vec4.W) < epsilon {
		vec4.W = 0
	}
	return vec4
}

func (v *MVec4) Round(threshold float64) *MVec4 {
	return &MVec4{
		Round(v.X, threshold),
		Round(v.Y, threshold),
		Round(v.Z, threshold),
		Round(v.W, threshold),
	}
}

