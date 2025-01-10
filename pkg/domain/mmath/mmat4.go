package mmath

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type MMat4 mgl64.Mat4

var (
	// Zero holds a zero matrix.
	MMat4Zero = &MMat4{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}

	// Ident holds an ident matrix.
	MMat4Ident = &MMat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	// Ident holds an ident matrix.
	MMat4ScaleIdent = &MMat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
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

func NewMMat4ByValues(m11, m21, m31, m41, m12, m22, m32, m42, m13, m23, m33, m43, m14, m24, m34, m44 float64) *MMat4 {
	return &MMat4{
		m11, m21, m31, m41,
		m12, m22, m32, m42,
		m13, m23, m33, m43,
		m14, m24, m34, m44,
	}
}

func NewMMat4FromAxisAngle(axis *MVec3, angle float64) *MMat4 {
	m := MMat4(mgl64.HomogRotate3D(angle, mgl64.Vec3{axis.X, axis.Y, axis.Z}))
	return &m
}

func NewMMat4FromLookAt(eye, center, up *MVec3) *MMat4 {
	m := MMat4(mgl64.LookAtV(mgl64.Vec3{eye.X, eye.Y, eye.Z},
		mgl64.Vec3{center.X, center.Y, center.Z}, mgl64.Vec3{up.X, up.Y, up.Z}))
	return &m
}

// IsZero
func (mat *MMat4) IsZero() bool {
	return *mat == *MMat4Zero
}

// IsIdent
func (mat *MMat4) IsIdent() bool {
	return mat.NearEquals(MMat4Ident, 1e-10)
}

// String
func (mat *MMat4) String() string {
	return mgl64.Mat4(*mat).String()
}

func (mat *MMat4) Copy() *MMat4 {
	copied := NewMMat4ByValues(
		mat[0], mat[1], mat[2], mat[3], mat[4], mat[5], mat[6], mat[7], mat[8], mat[9], mat[10], mat[11], mat[12], mat[13], mat[14], mat[15])
	return copied
}

// NearEquals
func (mat *MMat4) NearEquals(other *MMat4, tolerance float64) bool {
	return mgl64.Mat4(*mat).ApproxEqualThreshold(mgl64.Mat4(*other), tolerance)
}

// Trace returns the trace value for the matrix.
func (mat *MMat4) Trace() float64 {
	return mgl64.Mat4(*mat).Trace()
}

// Trace3 returns the trace value for the 3x3 sub-matrix.
func (mat *MMat4) Trace3() float64 {
	return mgl64.Mat4(*mat).Mat3().Trace()
}

// MulVec3 はベクトルと行列の掛け算を行います
func (mat *MMat4) MulVec3(other *MVec3) *MVec3 {
	v := mgl64.Mat4(*mat).Mul4x1(mgl64.Vec4{other.X, other.Y, other.Z, 1})
	return &MVec3{v.X() / v.W(), v.Y() / v.W(), v.Z() / v.W()}
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
	return &MVec3{mat[12], mat[13], mat[14]}
}

func (mat *MMat4) Scale(s *MVec3) *MMat4 {
	*mat = *s.ToScaleMat4().Mul(mat)
	return mat
}

func (mat *MMat4) Scaled(s *MVec3) *MMat4 {
	return s.ToScaleMat4().Muled(mat)
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
	return &MQuaternion{q.X(), q.Y(), q.Z(), q.W}
}

// Transpose transposes the matrix.
func (mat *MMat4) Transpose() *MMat4 {
	tm := mgl64.Mat4(*mat).Transpose()
	return (*MMat4)(&tm)
}

// Mul は行列の掛け算を行います
func (mat1 *MMat4) Mul(mat2 *MMat4) *MMat4 {
	m := mgl64.Mat4(*mat1).Mul4(mgl64.Mat4(*mat2))
	*mat1 = MMat4(m)
	return mat1
}

func (mat1 *MMat4) Add(mat2 *MMat4) *MMat4 {
	m := mgl64.Mat4(*mat1).Add(mgl64.Mat4(*mat2))
	*mat1 = MMat4(m)
	return mat1
}

func (mat1 *MMat4) Muled(mat2 *MMat4) *MMat4 {
	m := MMat4(mgl64.Mat4(*mat1).Mul4(mgl64.Mat4(*mat2)))
	return &m
}

func (mat1 *MMat4) Added(mat2 *MMat4) *MMat4 {
	m := MMat4(mgl64.Mat4(*mat1).Add(mgl64.Mat4(*mat2)))
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

func (mat *MMat4) Det() float64 {
	return mgl64.Mat4(*mat).Det()
}

// 逆行列
func (mat *MMat4) Inverse() *MMat4 {
	im := mgl64.Mat4(*mat).Inv()
	*mat = MMat4(im)
	return mat
}

func (mat *MMat4) Inverted() *MMat4 {
	im := mgl64.Mat4(*mat).Inv()
	return (*MMat4)(&im)
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
	v := mgl64.Mat4(*mat).Col(0)
	return &MVec3{v.X(), v.Y(), v.Z()}
}

func (mat *MMat4) AxisY() *MVec3 {
	v := mgl64.Mat4(*mat).Col(1)
	return &MVec3{v.X(), v.Y(), v.Z()}
}

func (mat *MMat4) AxisZ() *MVec3 {
	v := mgl64.Mat4(*mat).Col(2)
	return &MVec3{v.X(), v.Y(), v.Z()}
}
