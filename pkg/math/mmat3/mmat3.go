package mmat3

import (
	"fmt"

	"github.com/ungerik/go3d/float64/mat3"
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"

	"github.com/miu200521358/mlib_go/pkg/math/mquaternion"
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

type T mat3.T

var (
	// Zero holds a zero matrix.
	Zero = T{}

	// Ident holds an ident matrix.
	Ident = T{}
)

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (m T) GL() T {
	vec := m
	vec[0][1], vec[0][2] = -vec[0][1], -vec[0][2]
	vec[1][0], vec[2][0] = -vec[1][0], -vec[2][0]
	return vec
}

// IsZero
func (m T) IsZero() bool {
	return m == Zero
}

// Scale
func (mat *T) Scale(f float64) *T {
	return (*T)((*mat3.T)(mat).Scale(f))
}

// Scaled
func (mat *T) Scaled(f float64) T {
	return (T)((*mat3.T)(mat).Scaled(f))
}

// Trace
func (mat *T) Trace() float64 {
	return (*mat3.T)(mat).Trace()
}

// Scaling
func (mat *T) Scaling() mvec3.T {
	return (mvec3.T)((*mat3.T)(mat).Scaling())
}

// SetScaling
func (mat *T) SetScaling(s *mvec3.T) *T {
	return (*T)((*mat3.T)(mat).SetScaling((*vec3.T)(s)))
}

func (mat *T) ScaleVec2(s *mvec2.T) *T {
	return (*T)((*mat3.T)(mat).ScaleVec2((*vec2.T)(s)))
}

func (mat *T) SetTranslation(v *mvec2.T) *T {
	return (*T)((*mat3.T)(mat).SetTranslation((*vec2.T)(v)))
}

func (mat *T) Translate(v *vec2.T) *T {
	return (*T)((*mat3.T)(mat).Translate(v))
}

func (mat *T) TranslateX(dx float64) *T {
	return (*T)((*mat3.T)(mat).TranslateX(dx))
}

func (mat *T) TranslateY(dy float64) *T {
	return (*T)((*mat3.T)(mat).TranslateY(dy))
}

func (mat *T) AssignMul(a, b *T) *T {
	return (*T)((*mat3.T)(mat).AssignMul((*mat3.T)(a), (*mat3.T)(b)))
}

func (mat *T) Mul(f float64) *T {
	return (*T)((*mat3.T)(mat).Mul(f))
}

func (mat *T) Muled(f float64) T {
	return (T)((*mat3.T)(mat).Muled(f))
}

func (mat *T) MulVec3(v *mvec3.T) mvec3.T {
	return (mvec3.T)((*mat3.T)(mat).MulVec3((*vec3.T)(v)))
}

func (mat *T) TransformVec3(v *mvec3.T) {
	(*mat3.T)(mat).TransformVec3((*vec3.T)(v))
}

func (mat *T) Quaternion() mquaternion.T {
	return (mquaternion.T)((*mat3.T)(mat).Quaternion())
}

func (mat *T) AssignQuaternion(q *mquaternion.T) *T {
	return (*T)((*mat3.T)(mat).AssignQuaternion((*quaternion.T)(q)))
}

func (mat *T) AssignXRotation(angle float64) *T {
	return (*T)((*mat3.T)(mat).AssignXRotation(angle))
}

func (mat *T) AssignYRotation(angle float64) *T {
	return (*T)((*mat3.T)(mat).AssignYRotation(angle))
}

func (mat *T) AssignZRotation(angle float64) *T {
	return (*T)((*mat3.T)(mat).AssignZRotation(angle))
}

func (mat *T) AssignCoordinateSystem(x, y, z *vec3.T) *T {
	return (*T)((*mat3.T)(mat).AssignCoordinateSystem(x, y, z))
}

func (mat *T) AssignEulerRotation(yHead, xPitch, zRoll float64) *T {
	return (*T)((*mat3.T)(mat).AssignEulerRotation(yHead, xPitch, zRoll))
}

func (mat *T) ExtractEulerAngles() (yHead, xPitch, zRoll float64) {
	return (*mat3.T)(mat).ExtractEulerAngles()
}

func (mat *T) Determinant() float64 {
	return (*mat3.T)(mat).Determinant()
}

func (mat *T) IsReflective() bool {
	return (*mat3.T)(mat).IsReflective()
}

func (mat *T) PracticallyEquals(matrix *T, allowedDelta float64) bool {
	return (*mat3.T)(mat).PracticallyEquals((*mat3.T)(matrix), allowedDelta)
}

func (mat *T) Transpose() *T {
	return (*T)((*mat3.T)(mat).Transpose())
}

func (mat *T) Transposed() T {
	return (T)((*mat3.T)(mat).Transposed())
}

func (matrix *T) Adjugate() *T {
	return (*T)((*mat3.T)(matrix).Adjugate())
}

func (mat *T) Adjugated() T {
	return (T)((*mat3.T)(mat).Adjugated())
}

func (mat *T) Invert() (*T, error) {
	m, err := (*mat3.T)(mat).Invert()
	return (*T)(m), err
}

func (mat *T) Inverted() (T, error) {
	m, err := (*mat3.T)(mat).Inverted()
	return (T)(m), err
}

// Copy
func (mat *T) Copy() *T {
	return &T{
		vec3.T{mat[0].Slice()[0],
			mat[0].Slice()[1],
			mat[0].Slice()[2]},
		vec3.T{mat[1].Slice()[0],
			mat[1].Slice()[1],
			mat[1].Slice()[2]},
		vec3.T{mat[2].Slice()[0],
			mat[2].Slice()[1],
			mat[2].Slice()[2]},
	}
}

func (mat *T) String() string {
	return fmt.Sprintf("%v\n%v\n%v", mat[0], mat[1], mat[2])
}
