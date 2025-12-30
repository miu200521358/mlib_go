// 指示: miu200521358
package mmath

import (
	"fmt"
	"hash/fnv"
	"math"
)

var (
	MVec2Zero = &MVec2{}

	MVec2UnitX  = &MVec2{1, 0}
	MVec2UnitY  = &MVec2{0, 1}
	MVec2UnitXY = &MVec2{1, 1}

	MVec2MinVal = &MVec2{-math.MaxFloat64, -math.MaxFloat64}
	MVec2MaxVal = &MVec2{+math.MaxFloat64, +math.MaxFloat64}
)

type MVec2 struct {
	X float64
	Y float64
}

func NewMVec2() *MVec2 {
	return &MVec2{}
}

func (vec2 *MVec2) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f]", vec2.X, vec2.Y)
}

func (vec2 *MVec2) Add(other *MVec2) *MVec2 {
	vec2.X += other.X
	vec2.Y += other.Y
	return vec2
}

func (vec2 *MVec2) AddScalar(s float64) *MVec2 {
	vec2.X += s
	vec2.Y += s
	return vec2
}

func (vec2 *MVec2) Added(other *MVec2) *MVec2 {
	return &MVec2{vec2.X + other.X, vec2.Y + other.Y}
}

func (vec2 *MVec2) AddedScalar(s float64) *MVec2 {
	return &MVec2{vec2.X + s, vec2.Y + s}
}

func (vec2 *MVec2) Sub(other *MVec2) *MVec2 {
	vec2.X -= other.X
	vec2.Y -= other.Y
	return vec2
}

func (vec2 *MVec2) SubScalar(s float64) *MVec2 {
	vec2.X -= s
	vec2.Y -= s
	return vec2
}

func (vec2 *MVec2) Subed(other *MVec2) *MVec2 {
	return &MVec2{vec2.X - other.X, vec2.Y - other.Y}
}

func (vec2 *MVec2) SubedScalar(s float64) *MVec2 {
	return &MVec2{vec2.X - s, vec2.Y - s}
}

func (vec2 *MVec2) Mul(other *MVec2) *MVec2 {
	vec2.X *= other.X
	vec2.Y *= other.Y
	return vec2
}

func (vec2 *MVec2) MulScalar(s float64) *MVec2 {
	vec2.X *= s
	vec2.Y *= s
	return vec2
}

func (vec2 *MVec2) Muled(other *MVec2) *MVec2 {
	return &MVec2{vec2.X * other.X, vec2.Y * other.Y}
}

func (vec2 *MVec2) MuledScalar(s float64) *MVec2 {
	return &MVec2{vec2.X * s, vec2.Y * s}
}

func (vec2 *MVec2) Div(other *MVec2) *MVec2 {
	vec2.X /= other.X
	vec2.Y /= other.Y
	return vec2
}

func (vec2 *MVec2) DivScalar(s float64) *MVec2 {
	vec2.X /= s
	vec2.Y /= s
	return vec2
}

func (vec2 *MVec2) Dived(other *MVec2) *MVec2 {
	return &MVec2{vec2.X / other.X, vec2.Y / other.Y}
}

func (vec2 *MVec2) DivedScalar(s float64) *MVec2 {
	return &MVec2{vec2.X / s, vec2.Y / s}
}

func (vec2 *MVec2) Equals(other *MVec2) bool {
	return vec2.X == other.X && vec2.Y == other.Y
}

func (vec2 *MVec2) NotEquals(other MVec2) bool {
	return vec2.X != other.X || vec2.Y != other.Y
}

func (vec2 *MVec2) NearEquals(other *MVec2, epsilon float64) bool {
	return (math.Abs(vec2.X-other.X) <= epsilon) &&
		(math.Abs(vec2.Y-other.Y) <= epsilon)
}

func (vec2 *MVec2) LessThan(other *MVec2) bool {
	return vec2.X < other.X && vec2.Y < other.Y
}

func (vec2 *MVec2) LessThanOrEquals(other *MVec2) bool {
	return vec2.X <= other.X && vec2.Y <= other.Y
}

func (vec2 *MVec2) GreaterThan(other *MVec2) bool {
	return vec2.X > other.X && vec2.Y > other.Y
}

func (vec2 *MVec2) GreaterThanOrEquals(other *MVec2) bool {
	return vec2.X >= other.X && vec2.Y >= other.Y
}

func (vec2 *MVec2) Negate() *MVec2 {
	vec2.X = -vec2.X
	vec2.Y = -vec2.Y
	return vec2
}

func (vec2 *MVec2) Negated() *MVec2 {
	return &MVec2{-vec2.X, -vec2.Y}
}

func (vec2 *MVec2) Abs() *MVec2 {
	vec2.X = math.Abs(vec2.X)
	vec2.Y = math.Abs(vec2.Y)
	return vec2
}

func (vec2 *MVec2) Absed() *MVec2 {
	return &MVec2{math.Abs(vec2.X), math.Abs(vec2.Y)}
}

func (vec2 *MVec2) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.10f,%.10f", vec2.X, vec2.Y)))
	return h.Sum64()
}

func (vec2 *MVec2) IsZero() bool {
	return vec2.X == 0 && vec2.Y == 0
}

func (vec2 *MVec2) Length() float64 {
	return math.Hypot(vec2.X, vec2.Y)
}

func (vec2 *MVec2) LengthSqr() float64 {
	return vec2.X*vec2.X + vec2.Y*vec2.Y
}

func (vec2 *MVec2) Normalize() *MVec2 {
	sl := vec2.LengthSqr()
	if sl == 0 || sl == 1 {
		return vec2
	}
	return vec2.MulScalar(1 / math.Sqrt(sl))
}

func (vec2 *MVec2) Normalized() *MVec2 {
	vec := vec2.Copy()
	vec.Normalize()
	return vec
}

func (vec2 *MVec2) Angle(other *MVec2) float64 {
	v := vec2.Dot(other) / (vec2.Length() * other.Length())
	if v > 1. {
		v = v - 2
	} else if v < -1. {
		v = v + 2
	}
	return math.Acos(v)
}

func (vec2 *MVec2) Degree(other *MVec2) float64 {
	radian := vec2.Angle(other)
	degree := radian * (180 / math.Pi)
	return degree
}

func (vec2 *MVec2) Dot(other *MVec2) float64 {
	return vec2.X*other.X + vec2.Y*other.Y
}

func (vec2 *MVec2) Cross(other *MVec2) *MVec2 {
	return &MVec2{
		vec2.Y*other.X - vec2.X*other.Y,
		vec2.X*other.Y - vec2.Y*other.X,
	}
}

func (vec2 *MVec2) Min() *MVec2 {
	min := vec2.X
	if vec2.Y < min {
		min = vec2.Y
	}
	return &MVec2{min, min}
}

func (vec2 *MVec2) Max() *MVec2 {
	max := vec2.X
	if vec2.Y > max {
		max = vec2.Y
	}
	return &MVec2{max, max}
}

func (vec2 *MVec2) Clamp(min, max *MVec2) *MVec2 {
	vec2.X = Clamped(vec2.X, min.X, max.X)
	vec2.Y = Clamped(vec2.Y, min.Y, max.Y)
	return vec2
}

func (vec2 *MVec2) Clamped(min, max *MVec2) *MVec2 {
	result := *vec2
	result.Clamp(min, max)
	return &result
}

func (vec2 *MVec2) Clamp01() *MVec2 {
	return vec2.Clamp(MVec2Zero, MVec2UnitXY)
}

func (vec2 *MVec2) Clamped01() *MVec2 {
	result := vec2.Copy()
	result.Clamp01()
	return result
}

func (vec2 *MVec2) Rotate(angle float64) *MVec2 {
	sinus := math.Sin(angle)
	cosinus := math.Cos(angle)
	vec2.X = vec2.X*cosinus - vec2.Y*sinus
	vec2.Y = vec2.X*sinus + vec2.Y*cosinus
	return vec2
}

func (vec2 *MVec2) Rotated(angle float64) *MVec2 {
	copied := vec2.Copy()
	return copied.Rotate(angle)
}

func (vec2 *MVec2) RotateAroundPoint(point *MVec2, angle float64) *MVec2 {
	return vec2.Sub(point).Rotate(angle).Add(point)
}

func (vec2 *MVec2) Copy() *MVec2 {
	return &MVec2{vec2.X, vec2.Y}
}

func (vec2 *MVec2) Vector() []float64 {
	return []float64{vec2.X, vec2.Y}
}

func (v1 *MVec2) Lerp(v2 *MVec2, t float64) *MVec2 {
	if t <= 0 {
		return v1.Copy()
	} else if t >= 1 {
		return v2.Copy()
	}

	if v1.Equals(v2) {
		return v1.Copy()
	}

	return (v2.Subed(v1)).MuledScalar(t).Added(v1)
}

func (vec2 *MVec2) Round() *MVec2 {
	return &MVec2{
		math.Round(vec2.X),
		math.Round(vec2.Y),
	}
}

func (vec2 *MVec2) One() *MVec2 {
	vec := vec2.Vector()
	epsilon := 1e-8
	for i := 0; i < len(vec); i++ {
		if math.Abs(vec[i]) < epsilon {
			vec[i] = 1
		}
	}
	return &MVec2{vec[0], vec[1]}
}

func (vec2 *MVec2) Distance(other *MVec2) float64 {
	s := vec2.Subed(other)
	return s.Length()
}

func (vec2 *MVec2) ClampIfVerySmall() *MVec2 {
	epsilon := 1e-6
	if math.Abs(vec2.X) < epsilon {
		vec2.X = 0.0
	}
	if math.Abs(vec2.Y) < epsilon {
		vec2.Y = 0.0
	}
	return vec2
}

