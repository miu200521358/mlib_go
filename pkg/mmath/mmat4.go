package mmath

import (
	"fmt"
	"math"

	"github.com/ungerik/go3d/float64/mat3"
	"github.com/ungerik/go3d/float64/mat4"
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
	"github.com/ungerik/go3d/float64/vec4"
)

type MMat4 mat4.T

var (
	// Zero holds a zero matrix.
	MMat4Zero = MMat4{}

	// Ident holds an ident matrix.
	MMat4Ident = MMat4{
		vec4.T{1, 0, 0, 0},
		vec4.T{0, 1, 0, 0},
		vec4.T{0, 0, 1, 0},
		vec4.T{0, 0, 0, 1},
	}
)

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (m MMat4) GL() *[16]float32 {
	mat := [16]float32{
		float32(m[0][0]), float32(-m[0][1]), float32(-m[0][2]), float32(m[0][3]),
		float32(-m[1][0]), float32(m[1][1]), float32(m[1][2]), float32(m[1][3]),
		float32(-m[2][0]), float32(m[2][1]), float32(m[2][2]), float32(m[2][3]),
		float32(-m[3][0]), float32(m[3][1]), float32(m[3][2]), float32(m[3][3]),
	}
	return &mat
}

// IsZero
func (m MMat4) IsZero() bool {
	return m == MMat4Zero
}

// Scale
func (mat *MMat4) Scale(f float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).Scale(f))
}

// Scaled
func (mat *MMat4) Scaled(f float64) MMat4 {
	return (MMat4)((*mat4.T)(mat).Scaled(f))
}

// Trace
func (mat *MMat4) Trace() float64 {
	return (*mat4.T)(mat).Trace()
}

// Trace3
func (mat *MMat4) Trace3() float64 {
	return (*mat4.T)(mat).Trace3()
}

// AssignMat3x3
func (mat *MMat4) AssignMat3x3(m *MMat3) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignMat3x3((*mat3.T)(m)))
}

func (mat *MMat4) AssignMul(a, b *MMat4) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignMul((*mat4.T)(a), (*mat4.T)(b)))
}

func (mat *MMat4) MulVec4(v *MVec4) MVec4 {
	return (MVec4)((*mat4.T)(mat).MulVec4((*vec4.T)(v)))
}

func (mat *MMat4) TransformVec4(v *MVec4) {
	(*mat4.T)(mat).TransformVec4((*vec4.T)(v))
}

func (mat *MMat4) MulVec3(v *MVec3) MVec3 {
	return (MVec3)((*mat4.T)(mat).MulVec3((*vec3.T)(v)))
}

func (mat *MMat4) TransformVec3(v *MVec3) {
	(*mat4.T)(mat).TransformVec3((*vec3.T)(v))
}

func (mat *MMat4) MulVec3W(v *MVec3, w float64) MVec3 {
	return (MVec3)((*mat4.T)(mat).MulVec3W((*vec3.T)(v), w))
}

func (mat *MMat4) TransformVec3W(v *MVec3, w float64) {
	(*mat4.T)(mat).TransformVec3W((*vec3.T)(v), w)
}

func (mat *MMat4) SetTranslation(v *MVec3) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).SetTranslation((*vec3.T)(v)))
}

func (mat *MMat4) Translate(v *MVec3) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).Translate((*vec3.T)(v)))
}

func (mat *MMat4) TranslateX(dx float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).TranslateX(dx))
}

func (mat *MMat4) TranslateY(dy float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).TranslateY(dy))
}

func (mat *MMat4) TranslateZ(dz float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).TranslateZ(dz))
}

func (mat *MMat4) Scaling() MVec4 {
	return (MVec4)((*mat4.T)(mat).Scaling())
}

func (mat *MMat4) SetScaling(s *MVec4) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).SetScaling((*vec4.T)(s)))
}

func (mat *MMat4) ScaleVec3(s *MVec3) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).ScaleVec3((*vec3.T)(s)))
}

func (mat *MMat4) Quaternion() MQuaternion {
	return (MQuaternion)((*mat4.T)(mat).Quaternion())
}

func (mat *MMat4) AssignQuaternion(q *MQuaternion) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignQuaternion((*quaternion.T)(q)))
}

func (mat *MMat4) AssignXRotation(angle float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignXRotation(angle))
}

func (mat *MMat4) AssignYRotation(angle float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignYRotation(angle))
}

func (mat *MMat4) AssignZRotation(angle float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignZRotation(angle))
}

func (mat *MMat4) AssignCoordinateSystem(x, y, z *MVec3) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignCoordinateSystem((*vec3.T)(x), (*vec3.T)(y), (*vec3.T)(z)))
}

func (mat *MMat4) AssignEulerRotation(yHead, xPitch, zRoll float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignEulerRotation(yHead, xPitch, zRoll))
}

func (mat *MMat4) ExtractEulerAngles() (yHead, xPitch, zRoll float64) {
	return (*mat4.T)(mat).ExtractEulerAngles()
}

func (mat *MMat4) AssignPerspectiveProjection(left, right, bottom, top, znear, zfar float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignPerspectiveProjection(left, right, bottom, top, znear, zfar))
}

func (mat *MMat4) AssignOrthogonalProjection(left, right, bottom, top, znear, zfar float64) *MMat4 {
	return (*MMat4)((*mat4.T)(mat).AssignOrthogonalProjection(left, right, bottom, top, znear, zfar))
}

func (mat *MMat4) Determinant3x3() float64 {
	return (*mat4.T)(mat).Determinant3x3()
}

func (mat *MMat4) IsReflective() bool {
	return (*mat4.T)(mat).IsReflective()
}

func (mat *MMat4) Transpose() *MMat4 {
	return (*MMat4)((*mat4.T)(mat).Transpose())
}

func (mat *MMat4) Transpose3x3() *MMat4 {
	return (*MMat4)((*mat4.T)(mat).Transpose3x3())
}

// Copy
func (mat *MMat4) Copy() *MMat4 {
	return &MMat4{
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

func (mat *MMat4) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v]", mat[0], mat[1], mat[2], mat[3])
}

func (mat *MMat4) PracticallyEquals(compareVector *MMat4, allowedDelta float64) bool {
	return (math.Abs(mat[0][0]-compareVector[0][0]) <= allowedDelta) &&
		(math.Abs(mat[0][1]-compareVector[0][1]) <= allowedDelta) &&
		(math.Abs(mat[0][2]-compareVector[0][2]) <= allowedDelta) &&
		(math.Abs(mat[0][3]-compareVector[0][3]) <= allowedDelta) &&
		(math.Abs(mat[1][0]-compareVector[1][0]) <= allowedDelta) &&
		(math.Abs(mat[1][1]-compareVector[1][1]) <= allowedDelta) &&
		(math.Abs(mat[1][2]-compareVector[1][2]) <= allowedDelta) &&
		(math.Abs(mat[1][3]-compareVector[1][3]) <= allowedDelta) &&
		(math.Abs(mat[2][0]-compareVector[2][0]) <= allowedDelta) &&
		(math.Abs(mat[2][1]-compareVector[2][1]) <= allowedDelta) &&
		(math.Abs(mat[2][2]-compareVector[2][2]) <= allowedDelta) &&
		(math.Abs(mat[2][3]-compareVector[2][3]) <= allowedDelta) &&
		(math.Abs(mat[3][0]-compareVector[3][0]) <= allowedDelta) &&
		(math.Abs(mat[3][1]-compareVector[3][1]) <= allowedDelta) &&
		(math.Abs(mat[3][2]-compareVector[3][2]) <= allowedDelta) &&
		(math.Abs(mat[3][3]-compareVector[3][3]) <= allowedDelta)
}
