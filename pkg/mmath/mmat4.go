package mmath

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mbt"
)

type MMat4 [4]MVec4

func NewMMat4FromQuaternion(rot *MQuaternion) *MMat4 {
	mat := NewMMat4()
	mat.AssignQuaternion(rot)
	return mat
}

var (
	// Zero holds a zero matrix.
	MMat4Zero = MMat4{
		MVec4{0, 0, 0, 0},
		MVec4{0, 0, 0, 0},
		MVec4{0, 0, 0, 0},
		MVec4{0, 0, 0, 0},
	}

	// Ident holds an ident matrix.
	MMat4Ident = MMat4{
		MVec4{1, 0, 0, 0},
		MVec4{0, 1, 0, 0},
		MVec4{0, 0, 1, 0},
		MVec4{0, 0, 0, 1},
	}
)

func NewMMat4() *MMat4 {
	return &MMat4{
		MVec4{1, 0, 0, 0},
		MVec4{0, 1, 0, 0},
		MVec4{0, 0, 1, 0},
		MVec4{0, 0, 0, 1},
	}
}

func NewMMat4ByValues(m11, m12, m13, m14, m21, m22, m23, m24, m31, m32, m33, m34, m41, m42, m43, m44 float64) *MMat4 {
	return &MMat4{
		MVec4{m11, m12, m13, m14},
		MVec4{m21, m22, m23, m24},
		MVec4{m31, m32, m33, m34},
		MVec4{m41, m42, m43, m44},
	}
}

func NewMMat4ByVec4(v1, v2, v3, v4 *MVec4) *MMat4 {
	return &MMat4{
		*v1,
		*v2,
		*v3,
		*v4,
	}
}

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (m *MMat4) GL() *mgl32.Mat4 {
	tm := m.Transpose()
	mat := mgl32.Mat4{
		float32(tm[0][0]), float32(-tm[0][1]), float32(-tm[0][2]), float32(tm[0][3]),
		float32(-tm[1][0]), float32(tm[1][1]), float32(tm[1][2]), float32(tm[1][3]),
		float32(-tm[2][0]), float32(tm[2][1]), float32(tm[2][2]), float32(tm[2][3]),
		float32(-tm[3][0]), float32(tm[3][1]), float32(tm[3][2]), float32(tm[3][3]),
	}
	return &mat
}

// Bullet+OpenGL座標系に変換された行列ベクトルを返します
func (v *MMat4) Bullet() mbt.BtMatrix3x3 {
	glMat := v.GL()
	return mbt.NewBtMatrix3x3(
		float32(glMat[0]), float32(glMat[1]), float32(glMat[2]),
		float32(glMat[4]), float32(glMat[5]), float32(glMat[6]),
		float32(glMat[8]), float32(glMat[9]), float32(glMat[10]))
}

// Bullet+OpenGL座標系に変換 + Z軸の逆行列を加味した行列ベクトルを返します
func (v *MMat4) BulletInvZ() mbt.BtMatrix3x3 {
	invZ := NewMMat4()
	invZ.ScaleVec3(&MVec3{1, 1, -1})

	mat := NewMMat4()
	mat.Mul(invZ)
	mat.Mul(v)
	mat.Mul(invZ)

	return mbt.NewBtMatrix3x3(
		float32(mat[0][0]), float32(mat[0][1]), float32(mat[0][2]),
		float32(mat[1][0]), float32(mat[1][1]), float32(mat[1][2]),
		float32(mat[2][0]), float32(mat[2][1]), float32(mat[2][2]))
}

// IsZero
func (m *MMat4) IsZero() bool {
	return *m == MMat4Zero
}

// IsIdent
func (m *MMat4) IsIdent() bool {
	return m.PracticallyEquals(&MMat4Ident, 1e-10)
}

// String
func (m *MMat4) String() string {
	return m[0].String() + "\n" +
		m[1].String() + "\n" +
		m[2].String() + "\n" +
		m[3].String()
}

func (m *MMat4) Copy() *MMat4 {
	copied := NewMMat4()
	copier.CopyWithOption(copied, m, copier.Option{DeepCopy: true})
	return copied
}

// PracticallyEquals
func (m *MMat4) PracticallyEquals(other *MMat4, tolerance float64) bool {
	return m[0].PracticallyEquals(&other[0], tolerance) &&
		m[1].PracticallyEquals(&other[1], tolerance) &&
		m[2].PracticallyEquals(&other[2], tolerance) &&
		m[3].PracticallyEquals(&other[3], tolerance)
}

// Scale multiplies the diagonal scale elements by f returns mat.
func (mat *MMat4) Scale(f float64) *MMat4 {
	mat[0][0] *= f
	mat[1][1] *= f
	mat[2][2] *= f
	return mat
}

// Scaled returns a copy of the matrix with the diagonal scale elements multiplied by f.
func (mat *MMat4) Scaled(f float64) *MMat4 {
	r := *mat
	return r.Scale(f)
}

// Trace returns the trace value for the matrix.
func (mat *MMat4) Trace() float64 {
	return mat[0][0] + mat[1][1] + mat[2][2] + mat[3][3]
}

// Trace3 returns the trace value for the 3x3 sub-matrix.
func (mat *MMat4) Trace3() float64 {
	return mat[0][0] + mat[1][1] + mat[2][2]
}

// AssignMat3x3 assigns a 3x3 sub-matrix and sets the rest of the matrix to the ident value.
func (mat *MMat4) AssignMat3x3(m *MMat3) *MMat4 {
	*mat = MMat4{
		MVec4{m[0][0], m[1][0], m[2][0], 0},
		MVec4{m[0][1], m[1][1], m[2][1], 0},
		MVec4{m[0][2], m[1][2], m[2][2], 0},
		MVec4{0, 0, 0, 1},
	}
	return mat
}

// AssignMul multiplies a and b and assigns the result to T.
func (mat *MMat4) AssignMul(a, b *MMat4) *MMat4 {
	mat[0] = *a.MulVec4(&b[0])
	mat[1] = *a.MulVec4(&b[1])
	mat[2] = *a.MulVec4(&b[2])
	mat[3] = *a.MulVec4(&b[3])
	return mat
}

// MulVec4 multiplies v with mat and returns a new vector v' = M * v.
func (mat *MMat4) MulVec4(v *MVec4) *MVec4 {
	return &MVec4{
		mat[0][0]*v[0] + mat[1][0]*v[1] + mat[2][0]*v[2] + mat[3][0]*v[3],
		mat[0][1]*v[0] + mat[1][1]*v[1] + mat[2][1]*v[2] + mat[3][1]*v[3],
		mat[0][2]*v[0] + mat[1][2]*v[1] + mat[2][2]*v[2] + mat[3][2]*v[3],
		mat[0][3]*v[0] + mat[1][3]*v[1] + mat[2][3]*v[2] + mat[3][3]*v[3],
	}
}

// TransformVec4 multiplies v with mat and saves the result in v.
func (mat *MMat4) TransformVec4(v *MVec4) {
	// Use intermediate variables to not alter further computations.
	x := mat[0][0]*v[0] + mat[1][0]*v[1] + mat[2][0]*v[2] + mat[3][0]*v[3]
	y := mat[0][1]*v[0] + mat[1][1]*v[1] + mat[2][1]*v[2] + mat[3][1]*v[3]
	z := mat[0][2]*v[0] + mat[1][2]*v[1] + mat[2][2]*v[2] + mat[3][2]*v[3]
	v[3] = mat[0][3]*v[0] + mat[1][3]*v[1] + mat[2][3]*v[2] + mat[3][3]*v[3]
	v[0] = x
	v[1] = y
	v[2] = z
}

// MulVec3 multiplies v (converted to a vec4 as (v_1, v_2, v_3, 1))
// with mat and divides the result by w. Returns a new vec3.
func (mat *MMat4) MulVec3(other *MVec3) *MVec3 {
	s := [4]float64{
		mat[0][0]*other[0] + mat[0][1]*other[1] + mat[0][2]*other[2] + mat[0][3],
		mat[1][0]*other[0] + mat[1][1]*other[1] + mat[1][2]*other[2] + mat[1][3],
		mat[2][0]*other[0] + mat[2][1]*other[1] + mat[2][2]*other[2] + mat[2][3],
		mat[3][0]*other[0] + mat[3][1]*other[1] + mat[3][2]*other[2] + mat[3][3],
	}

	if s[3] == 1.0 {
		return &MVec3{s[0], s[1], s[2]}
	} else if s[3] == 0.0 {
		return NewMVec3()
	} else {
		return &MVec3{s[0] / s[3], s[1] / s[3], s[2] / s[3]}
	}
}

// TransformVec3 multiplies v (converted to a vec4 as (v_1, v_2, v_3, 1))
// with mat, divides the result by w and saves the result in v.
func (mat *MMat4) TransformVec3(v *MVec3) {
	x := mat[0][0]*v[0] + mat[1][0]*v[1] + mat[2][0]*v[2] + mat[3][0]
	y := mat[0][1]*v[0] + mat[1][1]*v[1] + mat[2][1]*v[2] + mat[3][1]
	z := mat[0][2]*v[0] + mat[1][2]*v[1] + mat[2][2]*v[2] + mat[3][2]
	w := mat[0][3]*v[0] + mat[1][3]*v[1] + mat[2][3]*v[2] + mat[3][3]
	oow := 1 / w
	v[0] = x * oow
	v[1] = y * oow
	v[2] = z * oow
}

// MulVec3W multiplies v with mat with w as fourth component of the vector.
// Useful to differentiate between vectors (w = 0) and points (w = 1)
// without transforming them to vec4.
func (mat *MMat4) MulVec3W(v *MVec3, w float64) *MVec3 {
	result := *v
	mat.TransformVec3W(&result, w)
	return &result
}

// TransformVec3W multiplies v with mat with w as fourth component of the vector and
// saves the result in v.
// Useful to differentiate between vectors (w = 0) and points (w = 1)
// without transforming them to vec4.
func (mat *MMat4) TransformVec3W(v *MVec3, w float64) {
	// use intermediate variables to not alter further computations
	x := mat[0][0]*v[0] + mat[1][0]*v[1] + mat[2][0]*v[2] + mat[3][0]*w
	y := mat[0][1]*v[0] + mat[1][1]*v[1] + mat[2][1]*v[2] + mat[3][1]*w
	v[2] = mat[0][2]*v[0] + mat[1][2]*v[1] + mat[2][2]*v[2] + mat[3][2]*w
	v[0] = x
	v[1] = y
}

// SetTranslation sets the translation elements of the matrix.
func (mat *MMat4) SetTranslation(v *MVec3) *MMat4 {
	mat[0][3] = v[0]
	mat[1][3] = v[1]
	mat[2][3] = v[2]
	return mat
}

// Translate adds v to the translation part of the matrix.
func (mat *MMat4) Translate(v *MVec3) *MMat4 {
	mat[0][3] += v[0]
	mat[1][3] += v[1]
	mat[2][3] += v[2]
	return mat
}

// TranslateX adds dx to the X-translation element of the matrix.
func (mat *MMat4) TranslateX(dx float64) *MMat4 {
	mat[0][3] += dx
	return mat
}

// TranslateY adds dy to the Y-translation element of the matrix.
func (mat *MMat4) TranslateY(dy float64) *MMat4 {
	mat[1][3] += dy
	return mat
}

// TranslateZ adds dz to the Z-translation element of the matrix.
func (mat *MMat4) TranslateZ(dz float64) *MMat4 {
	mat[2][3] += dz
	return mat
}

// Scaling returns the scaling diagonal of the matrix.
func (mat *MMat4) Scaling() *MVec3 {
	return &MVec3{mat[0][0], mat[1][1], mat[2][2]}
}

// SetScaling sets the scaling diagonal of the matrix.
func (mat *MMat4) SetScaling(s *MVec4) *MMat4 {
	mat[0][0] = s[0]
	mat[1][1] = s[1]
	mat[2][2] = s[2]
	mat[3][3] = s[3]
	return mat
}

// ScaleVec3 multiplies the scaling diagonal of the matrix by s.
func (mat *MMat4) ScaleVec3(s *MVec3) *MMat4 {
	mat[0][0] *= s[0]
	mat[1][1] *= s[1]
	mat[2][2] *= s[2]
	return mat
}

// Quaternion extracts a quaternion from the rotation part of the matrix.
func (mat *MMat4) Quaternion() *MQuaternion {
	trace := mat[0][0] + mat[1][1] + mat[2][2]

	q := MQuaternion{}
	if 0 < trace {
		s := 0.5 / math.Sqrt(trace+1)
		q.W = 0.25 / s
		q.V[0] = (mat[2][1] - mat[1][2]) * s
		q.V[1] = (mat[0][2] - mat[2][0]) * s
		q.V[2] = (mat[1][0] - mat[0][1]) * s
	} else {
		if mat[0][0] > mat[1][1] && mat[0][0] > mat[2][2] {
			s := 2 * math.Sqrt(1+mat[0][0]-mat[1][1]-mat[2][2])
			q.W = (mat[2][1] - mat[1][2]) / s
			q.V[0] = 0.25 * s
			q.V[1] = (mat[0][1] + mat[1][0]) / s
			q.V[2] = (mat[0][2] + mat[2][0]) / s
		} else if mat[1][1] > mat[2][2] {
			s := 2 * math.Sqrt(1+mat[1][1]-mat[0][0]-mat[2][2])
			q.W = (mat[0][2] - mat[2][0]) / s
			q.V[0] = (mat[0][1] + mat[1][0]) / s
			q.V[1] = 0.25 * s
			q.V[2] = (mat[1][2] + mat[2][1]) / s
		} else {
			s := 2 * math.Sqrt(1+mat[2][2]-mat[0][0]-mat[1][1])
			q.W = (mat[1][0] - mat[0][1]) / s
			q.V[0] = (mat[0][2] + mat[2][0]) / s
			q.V[1] = (mat[1][2] + mat[2][1]) / s
			q.V[2] = 0.25 * s
		}
	}

	return q.Normalize()
}

func (mat *MMat4) Mat3() *MMat3 {
	return &MMat3{
		MVec3{mat[0][0], mat[0][1], mat[0][2]},
		MVec3{mat[1][0], mat[1][1], mat[1][2]},
		MVec3{mat[2][0], mat[2][1], mat[2][2]},
	}
}

// AssignQuaternion assigns a quaternion to the rotations part of the matrix and sets the other elements to their ident value.
func (mat *MMat4) AssignQuaternion(q *MQuaternion) *MMat4 {
	xx := q.GetX() * q.GetX() * 2
	yy := q.GetY() * q.GetY() * 2
	zz := q.GetZ() * q.GetZ() * 2
	xy := q.GetX() * q.GetY() * 2
	xz := q.GetX() * q.GetZ() * 2
	yz := q.GetY() * q.GetZ() * 2
	wx := q.GetW() * q.GetX() * 2
	wy := q.GetW() * q.GetY() * 2
	wz := q.GetW() * q.GetZ() * 2

	mat[0][0] = 1 - (yy + zz)
	mat[0][1] = xy - wz
	mat[0][2] = xz + wy
	mat[0][3] = 0

	mat[1][0] = xy + wz
	mat[1][1] = 1 - (xx + zz)
	mat[1][2] = yz - wx
	mat[1][3] = 0

	mat[2][0] = xz - wy
	mat[2][1] = yz + wx
	mat[2][2] = 1 - (xx + yy)
	mat[2][3] = 0

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = 0
	mat[3][3] = 1

	return mat
}

// AssignXRotation assigns a rotation around the x axis to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat4) AssignXRotation(angle float64) *MMat4 {
	cosine := math.Cos(angle)
	sine := math.Sin(angle)

	mat[0][0] = 1
	mat[0][1] = 0
	mat[0][2] = 0
	mat[0][3] = 0

	mat[1][0] = 0
	mat[1][1] = cosine
	mat[1][2] = -sine
	mat[1][3] = 0

	mat[2][0] = 0
	mat[2][1] = sine
	mat[2][2] = cosine
	mat[2][3] = 0

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = 0
	mat[3][3] = 1

	return mat
}

// AssignYRotation assigns a rotation around the y axis to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat4) AssignYRotation(angle float64) *MMat4 {
	cosine := math.Cos(angle)
	sine := math.Sin(angle)

	mat[0][0] = cosine
	mat[0][1] = 0
	mat[0][2] = sine
	mat[0][3] = 0

	mat[1][0] = 0
	mat[1][1] = 1
	mat[1][2] = 0
	mat[1][3] = 0

	mat[2][0] = -sine
	mat[2][1] = 0
	mat[2][2] = cosine
	mat[2][3] = 0

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = 0
	mat[3][3] = 1

	return mat
}

// AssignZRotation assigns a rotation around the z axis to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat4) AssignZRotation(angle float64) *MMat4 {
	cosine := math.Cos(angle)
	sine := math.Sin(angle)

	mat[0][0] = cosine
	mat[0][1] = -sine
	mat[0][2] = 0
	mat[0][3] = 0

	mat[1][0] = sine
	mat[1][1] = cosine
	mat[1][2] = 0
	mat[1][3] = 0

	mat[2][0] = 0
	mat[2][1] = 0
	mat[2][2] = 1
	mat[2][3] = 0

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = 0
	mat[3][3] = 1

	return mat
}

// AssignCoordinateSystem assigns the rotation of a orthogonal coordinates system to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat4) AssignCoordinateSystem(x, y, z *MVec3) *MMat4 {
	mat[0][0] = x[0]
	mat[0][1] = x[1]
	mat[0][2] = x[2]
	mat[0][3] = 0

	mat[1][0] = y[0]
	mat[1][1] = y[1]
	mat[1][2] = y[2]
	mat[1][3] = 0

	mat[2][0] = z[0]
	mat[2][1] = z[1]
	mat[2][2] = z[2]
	mat[2][3] = 0

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = 0
	mat[3][3] = 1

	return mat
}

// AssignEulerRotation assigns Euler angle rotations to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat4) AssignEulerRotation(xPitch, yHead, zRoll float64) *MMat4 {
	sinH := math.Sin(yHead)
	cosH := math.Cos(yHead)
	sinP := math.Sin(xPitch)
	cosP := math.Cos(xPitch)
	sinR := math.Sin(zRoll)
	cosR := math.Cos(zRoll)

	mat[0][0] = cosR*cosH - sinR*sinP*sinH
	mat[0][1] = -sinR * cosP
	mat[0][2] = cosR*sinH + sinR*sinP*cosH
	mat[0][3] = 0

	mat[1][0] = sinR*cosH + cosR*sinP*sinH
	mat[1][1] = cosR * cosP
	mat[1][2] = sinR*sinH - cosR*sinP*cosH
	mat[1][3] = 0

	mat[2][0] = -cosP * sinH
	mat[2][1] = sinP
	mat[2][2] = cosP * cosH
	mat[2][3] = 0

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = 0
	mat[3][3] = 1

	return mat
}

// ExtractEulerAngles extracts the rotation part of the matrix as Euler angle rotation values.
func (mat *MMat4) ExtractEulerAngles() (xPitch, yHead, zRoll float64) {
	xPitch = math.Asin(mat[1][2])
	f12 := math.Abs(mat[1][2])
	if f12 > (1.0-0.0001) && f12 < (1.0+0.0001) { // f12 == 1.0
		yHead = 0.0
		zRoll = math.Atan2(mat[0][1], mat[0][0])
	} else {
		yHead = math.Atan2(-mat[0][2], mat[2][2])
		zRoll = math.Atan2(-mat[1][0], mat[1][1])
	}
	return yHead, xPitch, zRoll
}

// AssignPerspectiveProjection assigns a perspective projection transformation.
func (mat *MMat4) AssignPerspectiveProjection(left, right, bottom, top, znear, zfar float64) *MMat4 {
	near2 := znear + znear
	ooFarNear := 1 / (zfar - znear)

	mat[0][0] = near2 / (right - left)
	mat[0][1] = 0
	mat[0][2] = (right + left) / (right - left)
	mat[0][3] = 0

	mat[1][0] = 0
	mat[1][1] = near2 / (top - bottom)
	mat[1][2] = (top + bottom) / (top - bottom)
	mat[1][3] = 0

	mat[2][0] = 0
	mat[2][1] = 0
	mat[2][2] = -(zfar + znear) * ooFarNear
	mat[2][3] = -2 * zfar * znear * ooFarNear

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = -1
	mat[3][3] = 0

	return mat
}

// AssignOrthogonalProjection assigns an orthogonal projection transformation.
func (mat *MMat4) AssignOrthogonalProjection(left, right, bottom, top, znear, zfar float64) *MMat4 {
	ooRightLeft := 1 / (right - left)
	ooTopBottom := 1 / (top - bottom)
	ooFarNear := 1 / (zfar - znear)

	mat[0][0] = 2 * ooRightLeft
	mat[0][1] = 0
	mat[0][2] = 0
	mat[0][3] = -(right + left) * ooRightLeft

	mat[1][0] = 0
	mat[1][1] = 2 * ooTopBottom
	mat[1][2] = 0
	mat[1][3] = -(top + bottom) * ooTopBottom

	mat[2][0] = 0
	mat[2][1] = 0
	mat[2][2] = -2 * ooFarNear
	mat[2][3] = -(zfar + znear) * ooFarNear

	mat[3][0] = 0
	mat[3][1] = 0
	mat[3][2] = 0
	mat[3][3] = 1

	return mat
}

// Determinant3x3 returns the determinant of the 3x3 sub-matrix.
func (mat *MMat4) Determinant3x3() float64 {
	return mat[0][0]*mat[1][1]*mat[2][2] +
		mat[1][0]*mat[2][1]*mat[0][2] +
		mat[2][0]*mat[0][1]*mat[1][2] -
		mat[2][0]*mat[1][1]*mat[0][2] -
		mat[1][0]*mat[0][1]*mat[2][2] -
		mat[0][0]*mat[2][1]*mat[1][2]
}

// IsReflective returns true if the matrix can be reflected by a plane.
func (mat *MMat4) IsReflective() bool {
	return mat.Determinant3x3() < 0
}

func swap(a, b *float64) {
	*a, *b = *b, *a
}

// Transpose transposes the matrix.
func (mat *MMat4) Transpose() *MMat4 {
	swap(&mat[3][0], &mat[0][3])
	swap(&mat[3][1], &mat[1][3])
	swap(&mat[3][2], &mat[2][3])
	return mat.Transpose3x3()
}

// Transpose3x3 transposes the 3x3 sub-matrix.
func (mat *MMat4) Transpose3x3() *MMat4 {
	swap(&mat[1][0], &mat[0][1])
	swap(&mat[2][0], &mat[0][2])
	swap(&mat[2][1], &mat[1][2])
	return mat
}

// Mul は行列の掛け算を行います
func (m1 *MMat4) Mul(m2 *MMat4) {
	var result MMat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := 0.0
			for k := 0; k < 4; k++ {
				sum += m1[i][k] * m2[k][j]
			}
			result[i][j] = sum
		}
	}
	*m1 = result
}

func (mat *MMat4) Muled(a *MMat4) *MMat4 {
	copied := mat.Copy()
	copied.Mul(a)
	return copied
}

func (mat *MMat4) MuledFactor(v float64) *MMat4 {
	copied := mat.Copy()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			copied[i][j] *= v
		}
	}
	return copied
}

// 行列の移動情報
func (mat *MMat4) Translation() *MVec3 {
	return &MVec3{mat[0][3], mat[1][3], mat[2][3]}
}

func (m *MMat4) Det() float64 {
	return m[0][0]*m[1][1]*m[2][2]*m[3][3] - m[0][0]*m[1][1]*m[2][3]*m[3][2] - m[0][0]*m[1][2]*m[2][1]*m[3][3] + m[0][0]*m[1][2]*m[2][3]*m[3][1] + m[0][0]*m[1][3]*m[2][1]*m[3][2] - m[0][0]*m[1][3]*m[2][2]*m[3][1] - m[0][1]*m[1][0]*m[2][2]*m[3][3] + m[0][1]*m[1][0]*m[2][3]*m[3][2] + m[0][1]*m[1][2]*m[2][0]*m[3][3] - m[0][1]*m[1][2]*m[2][3]*m[3][0] - m[0][1]*m[1][3]*m[2][0]*m[3][2] + m[0][1]*m[1][3]*m[2][2]*m[3][0] + m[0][2]*m[1][0]*m[2][1]*m[3][3] - m[0][2]*m[1][0]*m[2][3]*m[3][1] - m[0][2]*m[1][1]*m[2][0]*m[3][3] + m[0][2]*m[1][1]*m[2][3]*m[3][0] + m[0][2]*m[1][3]*m[2][0]*m[3][1] - m[0][2]*m[1][3]*m[2][1]*m[3][0] - m[0][3]*m[1][0]*m[2][1]*m[3][2] + m[0][3]*m[1][0]*m[2][2]*m[3][1] + m[0][3]*m[1][1]*m[2][0]*m[3][2] - m[0][3]*m[1][1]*m[2][2]*m[3][0] - m[0][3]*m[1][2]*m[2][0]*m[3][1] + m[0][3]*m[1][2]*m[2][1]*m[3][0]
}

// 逆行列
func (m *MMat4) Inverse() *MMat4 {
	det := m.Det()
	if mgl64.FloatEqual(det, float64(0.0)) {
		return NewMMat4()
	}

	retMat := MMat4{
		MVec4{
			-m[1][3]*m[2][2]*m[3][1] + m[1][2]*m[2][3]*m[3][1] + m[1][3]*m[2][1]*m[3][2] - m[1][1]*m[2][3]*m[3][2] - m[1][2]*m[2][1]*m[3][3] + m[1][1]*m[2][2]*m[3][3],
			m[0][3]*m[2][2]*m[3][1] - m[0][2]*m[2][3]*m[3][1] - m[0][3]*m[2][1]*m[3][2] + m[0][1]*m[2][3]*m[3][2] + m[0][2]*m[2][1]*m[3][3] - m[0][1]*m[2][2]*m[3][3],
			-m[0][3]*m[1][2]*m[3][1] + m[0][2]*m[1][3]*m[3][1] + m[0][3]*m[1][1]*m[3][2] - m[0][1]*m[1][3]*m[3][2] - m[0][2]*m[1][1]*m[3][3] + m[0][1]*m[1][2]*m[3][3],
			m[0][3]*m[1][2]*m[2][1] - m[0][2]*m[1][3]*m[2][1] - m[0][3]*m[1][1]*m[2][2] + m[0][1]*m[1][3]*m[2][2] + m[0][2]*m[1][1]*m[2][3] - m[0][1]*m[1][2]*m[2][3],
		},
		MVec4{
			m[1][3]*m[2][2]*m[3][0] - m[1][2]*m[2][3]*m[3][0] - m[1][3]*m[2][0]*m[3][2] + m[1][0]*m[2][3]*m[3][2] + m[1][2]*m[2][0]*m[3][3] - m[1][0]*m[2][2]*m[3][3],
			-m[0][3]*m[2][2]*m[3][0] + m[0][2]*m[2][3]*m[3][0] + m[0][3]*m[2][0]*m[3][2] - m[0][0]*m[2][3]*m[3][2] - m[0][2]*m[2][0]*m[3][3] + m[0][0]*m[2][2]*m[3][3],
			m[0][3]*m[1][2]*m[3][0] - m[0][2]*m[1][3]*m[3][0] - m[0][3]*m[1][0]*m[3][2] + m[0][0]*m[1][3]*m[3][2] + m[0][2]*m[1][0]*m[3][3] - m[0][0]*m[1][2]*m[3][3],
			-m[0][3]*m[1][2]*m[2][0] + m[0][2]*m[1][3]*m[2][0] + m[0][3]*m[1][0]*m[2][2] - m[0][0]*m[1][3]*m[2][2] - m[0][2]*m[1][0]*m[2][3] + m[0][0]*m[1][2]*m[2][3],
		},
		MVec4{
			-m[1][3]*m[2][1]*m[3][0] + m[1][1]*m[2][3]*m[3][0] + m[1][3]*m[2][0]*m[3][1] - m[1][0]*m[2][3]*m[3][1] - m[1][1]*m[2][0]*m[3][3] + m[1][0]*m[2][1]*m[3][3],
			m[0][3]*m[2][1]*m[3][0] - m[0][1]*m[2][3]*m[3][0] - m[0][3]*m[2][0]*m[3][1] + m[0][0]*m[2][3]*m[3][1] + m[0][1]*m[2][0]*m[3][3] - m[0][0]*m[2][1]*m[3][3],
			-m[0][3]*m[1][1]*m[3][0] + m[0][1]*m[1][3]*m[3][0] + m[0][3]*m[1][0]*m[3][1] - m[0][0]*m[1][3]*m[3][1] - m[0][1]*m[1][0]*m[3][3] + m[0][0]*m[1][1]*m[3][3],
			m[0][3]*m[1][1]*m[2][0] - m[0][1]*m[1][3]*m[2][0] - m[0][3]*m[1][0]*m[2][1] + m[0][0]*m[1][3]*m[2][1] + m[0][1]*m[1][0]*m[2][3] - m[0][0]*m[1][1]*m[2][3],
		},
		MVec4{
			m[1][2]*m[2][1]*m[3][0] - m[1][1]*m[2][2]*m[3][0] - m[1][2]*m[2][0]*m[3][1] + m[1][0]*m[2][2]*m[3][1] + m[1][1]*m[2][0]*m[3][2] - m[1][0]*m[2][1]*m[3][2],
			-m[0][2]*m[2][1]*m[3][0] + m[0][1]*m[2][2]*m[3][0] + m[0][2]*m[2][0]*m[3][1] - m[0][0]*m[2][2]*m[3][1] - m[0][1]*m[2][0]*m[3][2] + m[0][0]*m[2][1]*m[3][2],
			m[0][2]*m[1][1]*m[3][0] - m[0][1]*m[1][2]*m[3][0] - m[0][2]*m[1][0]*m[3][1] + m[0][0]*m[1][2]*m[3][1] + m[0][1]*m[1][0]*m[3][2] - m[0][0]*m[1][1]*m[3][2],
			-m[0][2]*m[1][1]*m[2][0] + m[0][1]*m[1][2]*m[2][0] + m[0][2]*m[1][0]*m[2][1] - m[0][0]*m[1][2]*m[2][1] - m[0][1]*m[1][0]*m[2][2] + m[0][0]*m[1][1]*m[2][2],
		},
	}

	invDet := 1 / det

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			retMat[i][j] *= invDet
		}
	}

	return &retMat
}

func (m *MMat4) Inverted() *MMat4 {
	copied := m.Copy()
	return copied.Inverse()
}
