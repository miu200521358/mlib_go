package mmat4

import (
	"fmt"

	"github.com/ungerik/go3d/float64/mat3"
	"github.com/ungerik/go3d/float64/mat4"
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
	"github.com/ungerik/go3d/float64/vec4"

	"github.com/miu200521358/mlib_go/pkg/math/mmat3"
	"github.com/miu200521358/mlib_go/pkg/math/mquaternion"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"

)

type T mat4.T

var (
	// Zero holds a zero matrix.
	Zero = T{}

	// Ident holds an ident matrix.
	Ident = T{
		vec4.T{1, 0, 0, 0},
		vec4.T{0, 1, 0, 0},
		vec4.T{0, 0, 1, 0},
		vec4.T{0, 0, 0, 1},
	}
)

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (m T) GL() T {
	vec := m
	vec[0][1], vec[0][2] = -vec[0][1], -vec[0][2]
	vec[1][0], vec[2][0] = -vec[1][0], -vec[2][0]
	vec[3][0] = -vec[3][0]
	return vec
}

// IsZero
func (m T) IsZero() bool {
	return m == Zero
}

// Scale
func (mat *T) Scale(f float64) *T {
	return (*T)((*mat4.T)(mat).Scale(f))
}

// Scaled
func (mat *T) Scaled(f float64) T {
	return (T)((*mat4.T)(mat).Scaled(f))
}

// Trace
func (mat *T) Trace() float64 {
	return (*mat4.T)(mat).Trace()
}

// Trace3
func (mat *T) Trace3() float64 {
	return (*mat4.T)(mat).Trace3()
}

// AssignMat3x3
func (mat *T) AssignMat3x3(m *mmat3.T) *T {
	return (*T)((*mat4.T)(mat).AssignMat3x3((*mat3.T)(m)))
}

func (mat *T) AssignMul(a, b *T) *T {
	return (*T)((*mat4.T)(mat).AssignMul((*mat4.T)(a), (*mat4.T)(b)))
}

func (mat *T) MulVec4(v *mvec4.T) mvec4.T {
	return (mvec4.T)((*mat4.T)(mat).MulVec4((*vec4.T)(v)))
}

func (mat *T) TransformVec4(v *mvec4.T) {
	(*mat4.T)(mat).TransformVec4((*vec4.T)(v))
}

func (mat *T) MulVec3(v *mvec3.T) mvec3.T {
	return (mvec3.T)((*mat4.T)(mat).MulVec3((*vec3.T)(v)))
}

func (mat *T) TransformVec3(v *mvec3.T) {
	(*mat4.T)(mat).TransformVec3((*vec3.T)(v))
}

func (mat *T) MulVec3W(v *mvec3.T, w float64) mvec3.T {
	return (mvec3.T)((*mat4.T)(mat).MulVec3W((*vec3.T)(v), w))
}

func (mat *T) TransformVec3W(v *mvec3.T, w float64) {
	(*mat4.T)(mat).TransformVec3W((*vec3.T)(v), w)
}

func (mat *T) SetTranslation(v *mvec3.T) *T {
	return (*T)((*mat4.T)(mat).SetTranslation((*vec3.T)(v)))
}

func (mat *T) Translate(v *mvec3.T) *T {
	return (*T)((*mat4.T)(mat).Translate((*vec3.T)(v)))
}

func (mat *T) TranslateX(dx float64) *T {
	return (*T)((*mat4.T)(mat).TranslateX(dx))
}

func (mat *T) TranslateY(dy float64) *T {
	return (*T)((*mat4.T)(mat).TranslateY(dy))
}

func (mat *T) TranslateZ(dz float64) *T {
	return (*T)((*mat4.T)(mat).TranslateZ(dz))
}

func (mat *T) Scaling() mvec4.T {
	return (mvec4.T)((*mat4.T)(mat).Scaling())
}

func (mat *T) SetScaling(s *mvec4.T) *T {
	return (*T)((*mat4.T)(mat).SetScaling((*vec4.T)(s)))
}

func (mat *T) ScaleVec3(s *mvec3.T) *T {
	return (*T)((*mat4.T)(mat).ScaleVec3((*vec3.T)(s)))
}

func (mat *T) Quaternion() mquaternion.T {
	return (mquaternion.T)((*mat4.T)(mat).Quaternion())
}

func (mat *T) AssignQuaternion(q *mquaternion.T) *T {
	return (*T)((*mat4.T)(mat).AssignQuaternion((*quaternion.T)(q)))
}

func (mat *T) AssignXRotation(angle float64) *T {
	return (*T)((*mat4.T)(mat).AssignXRotation(angle))
}

func (mat *T) AssignYRotation(angle float64) *T {
	return (*T)((*mat4.T)(mat).AssignYRotation(angle))
}

func (mat *T) AssignZRotation(angle float64) *T {
	return (*T)((*mat4.T)(mat).AssignZRotation(angle))
}

func (mat *T) AssignCoordinateSystem(x, y, z *mvec3.T) *T {
	return (*T)((*mat4.T)(mat).AssignCoordinateSystem((*vec3.T)(x), (*vec3.T)(y), (*vec3.T)(z)))
}

func (mat *T) AssignEulerRotation(yHead, xPitch, zRoll float64) *T {
	return (*T)((*mat4.T)(mat).AssignEulerRotation(yHead, xPitch, zRoll))
}

func (mat *T) ExtractEulerAngles() (yHead, xPitch, zRoll float64) {
	return (*mat4.T)(mat).ExtractEulerAngles()
}

func (mat *T) AssignPerspectiveProjection(left, right, bottom, top, znear, zfar float64) *T {
	return (*T)((*mat4.T)(mat).AssignPerspectiveProjection(left, right, bottom, top, znear, zfar))
}

func (mat *T) AssignOrthogonalProjection(left, right, bottom, top, znear, zfar float64) *T {
	return (*T)((*mat4.T)(mat).AssignOrthogonalProjection(left, right, bottom, top, znear, zfar))
}

func (mat *T) Determinant3x3() float64 {
	return (*mat4.T)(mat).Determinant3x3()
}

func (mat *T) IsReflective() bool {
	return (*mat4.T)(mat).IsReflective()
}

func (mat *T) Transpose() *T {
	return (*T)((*mat4.T)(mat).Transpose())
}

func (mat *T) Transpose3x3() *T {
	return (*T)((*mat4.T)(mat).Transpose3x3())
}

// Copy
func (mat *T) Copy() *T {
	return &T{
		vec4.T{mat[0].Slice()[0],
			mat[0].Slice()[1],
			mat[0].Slice()[2],
			mat[0].Slice()[3]},
		vec4.T{mat[1].Slice()[0],
			mat[1].Slice()[1],
			mat[1].Slice()[2],
			mat[1].Slice()[3]},
		vec4.T{mat[2].Slice()[0],
			mat[2].Slice()[1],
			mat[2].Slice()[2],
			mat[2].Slice()[3]},
	}
}

func (mat *T) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v]", mat[0], mat[1], mat[2], mat[3])
}

