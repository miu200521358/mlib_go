package mmath

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ungerik/go3d/float64/mat3"
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"

)

type MMat3 mat3.T

var (
	// Zero holds a zero matrix.
	MMat3Zero = MMat3{}

	// Ident holds an ident matrix.
	MMat3Ident = MMat3{}
)

// GL OpenGL座標系に変換されたベクトルを返します
func (m MMat3) GL() mgl32.Mat3 {
	return mgl32.Mat3([9]float32{
		float32(m[0][0]), float32(-m[0][1]), float32(-m[0][2]),
		float32(-m[1][0]), float32(m[1][1]), float32(m[1][2]),
		float32(-m[2][0]), float32(m[2][1]), float32(m[2][2]),
	})
}

// IsZero
func (m MMat3) IsZero() bool {
	return m == MMat3Zero
}

// Scale
func (mat *MMat3) Scale(f float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).Scale(f))
}

// Scaled
func (mat *MMat3) Scaled(f float64) MMat3 {
	return (MMat3)((*mat3.T)(mat).Scaled(f))
}

// Trace
func (mat *MMat3) Trace() float64 {
	return (*mat3.T)(mat).Trace()
}

// Scaling
func (mat *MMat3) Scaling() MVec3 {
	return (MVec3)((*mat3.T)(mat).Scaling())
}

// SetScaling
func (mat *MMat3) SetScaling(s *MVec3) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).SetScaling((*vec3.T)(s)))
}

func (mat *MMat3) ScaleVec2(s *MVec2) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).ScaleVec2((*vec2.T)(s)))
}

func (mat *MMat3) SetTranslation(v *MVec2) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).SetTranslation((*vec2.T)(v)))
}

func (mat *MMat3) Translate(v *vec2.T) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).Translate(v))
}

func (mat *MMat3) TranslateX(dx float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).TranslateX(dx))
}

func (mat *MMat3) TranslateY(dy float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).TranslateY(dy))
}

func (mat *MMat3) AssignMul(a, b *MMat3) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).AssignMul((*mat3.T)(a), (*mat3.T)(b)))
}

func (mat *MMat3) Mul(f float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).Mul(f))
}

func (mat *MMat3) Muled(f float64) MMat3 {
	return (MMat3)((*mat3.T)(mat).Muled(f))
}

func (mat *MMat3) MulVec3(v *MVec3) MVec3 {
	return (MVec3)((*mat3.T)(mat).MulVec3((*vec3.T)(v)))
}

func (mat *MMat3) TransformVec3(v *MVec3) {
	(*mat3.T)(mat).TransformVec3((*vec3.T)(v))
}

func (mat *MMat3) Quaternion() MQuaternion {
	return (MQuaternion)((*mat3.T)(mat).Quaternion())
}

func (mat *MMat3) AssignQuaternion(q *MQuaternion) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).AssignQuaternion((*quaternion.T)(q)))
}

func (mat *MMat3) AssignXRotation(angle float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).AssignXRotation(angle))
}

func (mat *MMat3) AssignYRotation(angle float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).AssignYRotation(angle))
}

func (mat *MMat3) AssignZRotation(angle float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).AssignZRotation(angle))
}

func (mat *MMat3) AssignCoordinateSystem(x, y, z *vec3.T) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).AssignCoordinateSystem(x, y, z))
}

func (mat *MMat3) AssignEulerRotation(yHead, xPitch, zRoll float64) *MMat3 {
	return (*MMat3)((*mat3.T)(mat).AssignEulerRotation(yHead, xPitch, zRoll))
}

func (mat *MMat3) ExtractEulerAngles() (yHead, xPitch, zRoll float64) {
	return (*mat3.T)(mat).ExtractEulerAngles()
}

func (mat *MMat3) Determinant() float64 {
	return (*mat3.T)(mat).Determinant()
}

func (mat *MMat3) IsReflective() bool {
	return (*mat3.T)(mat).IsReflective()
}

func (mat *MMat3) PracticallyEquals(matrix *MMat3, allowedDelta float64) bool {
	return (*mat3.T)(mat).PracticallyEquals((*mat3.T)(matrix), allowedDelta)
}

func (mat *MMat3) Transpose() *MMat3 {
	return (*MMat3)((*mat3.T)(mat).Transpose())
}

func (mat *MMat3) Transposed() MMat3 {
	return (MMat3)((*mat3.T)(mat).Transposed())
}

func (matrix *MMat3) Adjugate() *MMat3 {
	return (*MMat3)((*mat3.T)(matrix).Adjugate())
}

func (mat *MMat3) Adjugated() MMat3 {
	return (MMat3)((*mat3.T)(mat).Adjugated())
}

func (mat *MMat3) Invert() (*MMat3, error) {
	m, err := (*mat3.T)(mat).Invert()
	return (*MMat3)(m), err
}

func (mat *MMat3) Inverted() (MMat3, error) {
	m, err := (*mat3.T)(mat).Inverted()
	return (MMat3)(m), err
}

// Copy
func (mat *MMat3) Copy() *MMat3 {
	return &MMat3{
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

func (mat *MMat3) String() string {
	return fmt.Sprintf("%v\n%v\n%v", mat[0], mat[1], mat[2])
}
