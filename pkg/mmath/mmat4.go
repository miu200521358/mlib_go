package mmath

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type MMat4 mgl64.Mat4

var (
	// Zero holds a zero matrix.
	MMat4Zero = MMat4{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}

	// Ident holds an ident matrix.
	MMat4Ident = MMat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
)

func NewMMat4() *MMat4 {
	return &MMat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func NewMMat4ByValues(m11, m12, m13, m14, m21, m22, m23, m24, m31, m32, m33, m34, m41, m42, m43, m44 float64) *MMat4 {
	return &MMat4{
		m11, m12, m13, m14,
		m21, m22, m23, m24,
		m31, m32, m33, m34,
		m41, m42, m43, m44,
	}
}

func NewMMat4FromAxisAngle(axis *MVec3, angle float64) *MMat4 {
	m := MMat4(mgl64.HomogRotate3D(angle, mgl64.Vec3(*axis)))
	return &m
}

func NewMMat4FromLookAt(eye, center, up *MVec3) *MMat4 {
	m := MMat4(mgl64.LookAtV(mgl64.Vec3(*eye), mgl64.Vec3(*center), mgl64.Vec3(*up)))
	return &m
}

// IsZero
func (m *MMat4) IsZero() bool {
	return *m == MMat4Zero
}

// IsIdent
func (m *MMat4) IsIdent() bool {
	return m.NearEquals(&MMat4Ident, 1e-10)
}

// String
func (m *MMat4) String() string {
	return mgl64.Mat4(*m).String()
}

func (m *MMat4) Copy() *MMat4 {
	copied := NewMMat4ByValues(
		m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15])
	return copied
}

// NearEquals
func (m *MMat4) NearEquals(other *MMat4, tolerance float64) bool {
	return mgl64.Mat4(*m).ApproxEqualThreshold(mgl64.Mat4(*other), tolerance)
}

// Trace returns the trace value for the matrix.
func (mat *MMat4) Trace() float64 {
	return mgl64.Mat4(*mat).Trace()
}

// Trace3 returns the trace value for the 3x3 sub-matrix.
func (mat *MMat4) Trace3() float64 {
	return mgl64.Mat4(*mat).Mat3().Trace()
}

// MulVec3 multiplies v (converted to a vec4 as (v_1, v_2, v_3, 1))
// with mat and divides the result by w. Returns a new vec3.
func (mat *MMat4) MulVec3(other *MVec3) *MVec3 {
	return other.ToMat4().Mul(mat.Copy()).Translation()
}

// Translate adds v to the translation part of the matrix.
func (mat *MMat4) Translate(v *MVec3) *MMat4 {
	*mat = *v.ToMat4().Mul(mat)
	return mat
}

func (mat *MMat4) Translated(v *MVec3) *MMat4 {
	return v.ToMat4().Mul(mat)
}

// 行列の移動情報
func (mat *MMat4) Translation() *MVec3 {
	return &MVec3{mat[3], mat[7], mat[11]}
}

func (mat *MMat4) Scale(s *MVec3) *MMat4 {
	*mat = *s.ToScaleMat4().Mul(mat)
	return mat
}

func (mat *MMat4) Scaled(s *MVec3) *MMat4 {
	return s.ToScaleMat4().Mul(mat)
}

func (mat *MMat4) Scaling() *MVec3 {
	return &MVec3{mat[0], mat[5], mat[10]}
}

func (mat *MMat4) Rotate(quat *MQuaternion) *MMat4 {
	*mat = *quat.ToMat4().Mul(mat)
	return mat
}

func (mat *MMat4) Rotated(quat *MQuaternion) *MMat4 {
	return quat.ToMat4().Mul(mat)
}

func (mat *MMat4) Quaternion() *MQuaternion {
	q := mgl64.Mat4ToQuat(mgl64.Mat4(*mat))
	return &MQuaternion{q.W, q.V}
}

func (mat *MMat4) Mat3() *MMat3 {
	return &MMat3{
		mat[0], mat[1], mat[2],
		mat[4], mat[5], mat[6],
		mat[8], mat[9], mat[10],
	}
}

// Transpose transposes the matrix.
func (mat *MMat4) Transpose() *MMat4 {
	tm := mgl64.Mat4(*mat).Transpose()
	return (*MMat4)(&tm)
}

// Mul は行列の掛け算を行います
func (m1 *MMat4) Mul(m2 *MMat4) *MMat4 {
	m := mgl64.Mat4(*m1).Mul4(mgl64.Mat4(*m2))
	*m1 = MMat4(m)
	return m1
}

func (m1 *MMat4) Muled(m2 *MMat4) *MMat4 {
	m := MMat4(mgl64.Mat4(*m1).Mul4(mgl64.Mat4(*m2)))
	return &m
}

func (mat *MMat4) MulScalar(v float64) *MMat4 {
	m := mgl64.Mat4(*mat.Copy()).Mul(v)
	*mat = MMat4(m)
	return mat
}

func (mat *MMat4) MuledScalar(v float64) *MMat4 {
	copied := mat.Copy()
	copied.MulScalar(v)
	return copied
}

func (m *MMat4) Det() float64 {
	return mgl64.Mat4(*m).Det()
}

// 逆行列
func (m *MMat4) Inverse() *MMat4 {
	im := mgl64.Mat4(*m).Inv()
	return (*MMat4)(&im)
}

func (mat *MMat4) Inverted() *MMat4 {
	copied := mat.Copy()
	return copied.Inverse()
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (mat *MMat4) ClampIfVerySmall() *MMat4 {
	epsilon := 1e-6
	for i := range mat {
		if math.Abs(mat[i]) < epsilon {
			mat[i] = 0
		}
	}
	return mat
}

func (mat *MMat4) AxisX() *MVec3 {
	v := MVec3(mgl64.Mat4(*mat).Row(0).Vec3())
	return &v
}

func (mat *MMat4) AxisY() *MVec3 {
	v := MVec3(mgl64.Mat4(*mat).Row(1).Vec3())
	return &v
}

func (mat *MMat4) AxisZ() *MVec3 {
	v := MVec3(mgl64.Mat4(*mat).Row(2).Vec3())
	return &v
}
