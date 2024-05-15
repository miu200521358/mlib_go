package mmath

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type MMat3 [3]MVec3

func NewMMat3() *MMat3 {
	return &MMat3{
		MVec3{1, 0, 0},
		MVec3{0, 1, 0},
		MVec3{0, 0, 1},
	}
}

var (
	// Zero holds a zero matrix.
	MMat3Zero = MMat3{
		MVec3{0, 0, 0},
		MVec3{0, 0, 0},
		MVec3{0, 0, 0},
	}

	// Ident holds an ident matrix.
	MMat3Ident = MMat3{
		MVec3{1, 0, 0},
		MVec3{0, 1, 0},
		MVec3{0, 0, 1},
	}
)

// GL OpenGL座標系に変換されたベクトルを返します
func (m *MMat3) GL() *mgl64.Mat3 {
	mat := mgl64.Mat3([9]float64{
		m[0][0], -m[0][1], -m[0][2],
		-m[1][0], m[1][1], m[1][2],
		-m[2][0], m[2][1], m[2][2],
	})
	return &mat
}

// IsZero
func (m *MMat3) IsZero() bool {
	return *m == MMat3Zero
}

// IsIdent
func (m *MMat3) IsIdent() bool {
	return m.PracticallyEquals(&MMat3Ident, 1e-10)
}

// String
func (m *MMat3) String() string {
	return m[0].String() + "\n" +
		m[1].String() + "\n" +
		m[2].String()
}

func (m *MMat3) Copy() *MMat3 {
	copied := NewMMat3ByValues(m[0][0], m[0][1], m[0][2], m[1][0], m[1][1], m[1][2], m[2][0], m[2][1], m[2][2])
	return copied
}

func NewMMat3ByValues(m11, m12, m13, m21, m22, m23, m31, m32, m33 float64) *MMat3 {
	return &MMat3{
		MVec3{m11, m12, m13},
		MVec3{m21, m22, m23},
		MVec3{m31, m32, m33},
	}
}

// Slice returns the elements of the matrix as slice.
func (mat *MMat3) Slice() []float64 {
	return []float64{
		mat[0][0], mat[0][1], mat[0][2],
		mat[1][0], mat[1][1], mat[1][2],
		mat[2][0], mat[2][1], mat[2][2],
	}
}

// Scale multiplies the diagonal scale elements by f returns mat.
func (mat *MMat3) Scale(f float64) *MMat3 {
	mat[0][0] *= f
	mat[1][1] *= f
	mat[2][2] *= f
	return mat
}

// Scaled returns a copy of the matrix with the diagonal scale elements multiplied by f.
func (mat *MMat3) Scaled(f float64) *MMat3 {
	r := *mat
	return r.Scale(f)
}

// Scaling returns the scaling diagonal of the matrix.
func (mat *MMat3) Scaling() *MVec3 {
	return &MVec3{mat[0][0], mat[1][1], mat[2][2]}
}

// SetScaling sets the scaling diagonal of the matrix.
func (mat *MMat3) SetScaling(s *MVec3) *MMat3 {
	mat[0][0] = s[0]
	mat[1][1] = s[1]
	mat[2][2] = s[2]
	return mat
}

// ScaleVec2 multiplies the 2D scaling diagonal of the matrix by s.
func (mat *MMat3) ScaleVec2(s *MVec2) *MMat3 {
	mat[0][0] *= s[0]
	mat[1][1] *= s[1]
	return mat
}

// SetTranslation sets the 2D translation elements of the matrix.
func (mat *MMat3) SetTranslation(v *MVec2) *MMat3 {
	mat[2][0] = v[0]
	mat[2][1] = v[1]
	return mat
}

// Translate adds v to the 2D translation part of the matrix.
func (mat *MMat3) Translate(v *MVec2) *MMat3 {
	mat[0][2] += v[0]
	mat[1][2] += v[1]
	return mat
}

// TranslateX adds dx to the 2D X-translation element of the matrix.
func (mat *MMat3) TranslateX(dx float64) *MMat3 {
	mat[0][2] += dx
	return mat
}

// TranslateY adds dy to the 2D Y-translation element of the matrix.
func (mat *MMat3) TranslateY(dy float64) *MMat3 {
	mat[1][2] += dy
	return mat
}

// Trace returns the trace value for the matrix.
func (mat *MMat3) Trace() float64 {
	return mat[0][0] + mat[1][1] + mat[2][2]
}

// AssignMul multiplies a and b and assigns the result to mat.
func (mat *MMat3) AssignMul(a, b *MMat3) *MMat3 {
	mat[0] = *a.MulVec3(&b[0])
	mat[1] = *a.MulVec3(&b[1])
	mat[2] = *a.MulVec3(&b[2])
	return mat
}

// Mul multiplies every element by f and returns mat.
func (mat *MMat3) Mul(f float64) *MMat3 {
	mat[0][0] *= f
	mat[0][1] *= f
	mat[0][2] *= f

	mat[1][0] *= f
	mat[1][1] *= f
	mat[1][2] *= f

	mat[2][0] *= f
	mat[2][1] *= f
	mat[2][2] *= f

	return mat
}

// Muled returns a copy of the matrix with every element multiplied by f.
func (mat *MMat3) Muled(f float64) *MMat3 {
	result := *mat
	result.Mul(f)
	return &result
}

// MulVec3 multiplies v with mat and returns a new vector v' = M * v.
func (mat *MMat3) MulVec3(v *MVec3) *MVec3 {
	return &MVec3{
		mat[0][0]*v[0] + mat[1][0]*v[1] + mat[2][0]*v[2],
		mat[0][1]*v[0] + mat[1][1]*v[1] + mat[2][1]*v[2],
		mat[0][2]*v[0] + mat[1][2]*v[1] + mat[2][2]*v[2],
	}
}

// TransformVec3 multiplies v with mat and saves the result in v.
func (mat *MMat3) TransformVec3(v *MVec3) {
	// Use intermediate variables to not alter further computations.
	x := mat[0][0]*v[0] + mat[1][0]*v[1] + mat[2][0]*v[2]
	y := mat[0][1]*v[0] + mat[1][1]*v[1] + mat[2][1]*v[2]
	v[2] = mat[0][2]*v[0] + mat[1][2]*v[1] + mat[2][2]*v[2]
	v[0] = x
	v[1] = y
}

// Quaternion extracts a quaternion from the rotation part of the matrix.
func (mat *MMat3) Quaternion() *MQuaternion {
	tr := mat.Trace()

	s := math.Sqrt(tr + 1)
	w := s * 0.5
	s = 0.5 / s

	q := NewMQuaternionByValues(
		(mat[1][2]-mat[2][1])*s,
		(mat[2][0]-mat[0][2])*s,
		(mat[0][1]-mat[1][0])*s,
		w,
	)
	return q.Normalize()
}

// AssignQuaternion assigns a quaternion to the rotations part of the matrix and sets the other elements to their ident value.
func (mat *MMat3) AssignQuaternion(q *MQuaternion) *MMat3 {
	xx := q.V[0] * q.V[0] * 2
	yy := q.V[1] * q.V[1] * 2
	zz := q.V[2] * q.V[2] * 2
	xy := q.V[0] * q.V[1] * 2
	xz := q.V[0] * q.V[2] * 2
	yz := q.V[1] * q.V[2] * 2
	wx := q.W * q.V[0] * 2
	wy := q.W * q.V[1] * 2
	wz := q.W * q.V[2] * 2

	mat[0][0] = 1 - (yy + zz)
	mat[1][0] = xy - wz
	mat[2][0] = xz + wy

	mat[0][1] = xy + wz
	mat[1][1] = 1 - (xx + zz)
	mat[2][1] = yz - wx

	mat[0][2] = xz - wy
	mat[1][2] = yz + wx
	mat[2][2] = 1 - (xx + yy)

	return mat
}

// AssignXRotation assigns a rotation around the x axis to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat3) AssignXRotation(angle float64) *MMat3 {
	cosine := math.Cos(angle)
	sine := math.Sin(angle)

	mat[0][0] = 1
	mat[1][0] = 0
	mat[2][0] = 0

	mat[0][1] = 0
	mat[1][1] = cosine
	mat[2][1] = -sine

	mat[0][2] = 0
	mat[1][2] = sine
	mat[2][2] = cosine

	return mat
}

// AssignYRotation assigns a rotation around the y axis to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat3) AssignYRotation(angle float64) *MMat3 {
	cosine := math.Cos(angle)
	sine := math.Sin(angle)

	mat[0][0] = cosine
	mat[1][0] = 0
	mat[2][0] = sine

	mat[0][1] = 0
	mat[1][1] = 1
	mat[2][1] = 0

	mat[0][2] = -sine
	mat[1][2] = 0
	mat[2][2] = cosine

	return mat
}

// AssignZRotation assigns a rotation around the z axis to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat3) AssignZRotation(angle float64) *MMat3 {
	cosine := math.Cos(angle)
	sine := math.Sin(angle)

	mat[0][0] = cosine
	mat[1][0] = -sine
	mat[2][0] = 0

	mat[0][1] = sine
	mat[1][1] = cosine
	mat[2][1] = 0

	mat[0][2] = 0
	mat[1][2] = 0
	mat[2][2] = 1

	return mat
}

// AssignCoordinateSystem assigns the rotation of a orthogonal coordinates system to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat3) AssignCoordinateSystem(x, y, z *MVec3) *MMat3 {
	mat[0][0] = x[0]
	mat[1][0] = x[1]
	mat[2][0] = x[2]

	mat[0][1] = y[0]
	mat[1][1] = y[1]
	mat[2][1] = y[2]

	mat[0][2] = z[0]
	mat[1][2] = z[1]
	mat[2][2] = z[2]

	return mat
}

// AssignEulerRotation assigns Euler angle rotations to the rotation part of the matrix and sets the remaining elements to their ident value.
func (mat *MMat3) AssignEulerRotation(yHead, xPitch, zRoll float64) *MMat3 {
	sinH := math.Sin(yHead)
	cosH := math.Cos(yHead)
	sinP := math.Sin(xPitch)
	cosP := math.Cos(xPitch)
	sinR := math.Sin(zRoll)
	cosR := math.Cos(zRoll)

	mat[0][0] = cosR*cosH - sinR*sinP*sinH
	mat[1][0] = -sinR * cosP
	mat[2][0] = cosR*sinH + sinR*sinP*cosH

	mat[0][1] = sinR*cosH + cosR*sinP*sinH
	mat[1][1] = cosR * cosP
	mat[2][1] = sinR*sinH - cosR*sinP*cosH

	mat[0][2] = -cosP * sinH
	mat[1][2] = sinP
	mat[2][2] = cosP * cosH

	return mat
}

// ExtractEulerAngles extracts the rotation part of the matrix as Euler angle rotation values.
func (mat *MMat3) ExtractEulerAngles() (yHead, xPitch, zRoll float64) {
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

// Determinant returns the determinant of the matrix.
func (mat *MMat3) Determinant() float64 {
	// | a b c |
	// | d e f | = det A
	// | g h i |
	return mat[0][0]*mat[1][1]*mat[2][2] + // aei
		mat[1][0]*mat[2][1]*mat[0][2] + // dhc
		mat[2][0]*mat[0][1]*mat[1][2] - // gbf
		mat[2][0]*mat[1][1]*mat[0][2] - // gec
		mat[1][0]*mat[0][1]*mat[2][2] - // dbi
		mat[0][0]*mat[2][1]*mat[1][2] // ahf
}

// IsReflective returns true if the matrix can be reflected by a plane.
func (mat *MMat3) IsReflective() bool {
	return mat.Determinant() < 0
}

// PracticallyEquals compares two matrices if they are equal with each other within a delta tolerance.
func (mat *MMat3) PracticallyEquals(matrix *MMat3, allowedDelta float64) bool {
	return mat[0].PracticallyEquals(&matrix[0], allowedDelta) &&
		mat[1].PracticallyEquals(&matrix[1], allowedDelta) &&
		mat[2].PracticallyEquals(&matrix[2], allowedDelta)
}

// Transpose transposes the matrix.
func (mat *MMat3) Transpose() *MMat3 {
	mat[1][0], mat[0][1] = mat[0][1], mat[1][0]
	mat[2][0], mat[0][2] = mat[0][2], mat[2][0]
	mat[2][1], mat[1][2] = mat[1][2], mat[2][1]
	return mat
}

// Transposed returns a transposed copy the matrix.
func (mat *MMat3) Transposed() *MMat3 {
	result := *mat
	result.Transpose()
	return &result
}

// Adjugate computes the adjugate of this matrix and returns mat
func (matrix *MMat3) Adjugate() *MMat3 {
	mat := *matrix

	matrix[0][0] = +(mat[1][1]*mat[2][2] - mat[1][2]*mat[2][1])
	matrix[0][1] = -(mat[0][1]*mat[2][2] - mat[0][2]*mat[2][1])
	matrix[0][2] = +(mat[0][1]*mat[1][2] - mat[0][2]*mat[1][1])

	matrix[1][0] = -(mat[1][0]*mat[2][2] - mat[1][2]*mat[2][0])
	matrix[1][1] = +(mat[0][0]*mat[2][2] - mat[0][2]*mat[2][0])
	matrix[1][2] = -(mat[0][0]*mat[1][2] - mat[0][2]*mat[1][0])

	matrix[2][0] = +(mat[1][0]*mat[2][1] - mat[1][1]*mat[2][0])
	matrix[2][1] = -(mat[0][0]*mat[2][1] - mat[0][1]*mat[2][0])
	matrix[2][2] = +(mat[0][0]*mat[1][1] - mat[0][1]*mat[1][0])

	return matrix
}

// Adjugated returns an adjugated copy of the matrix.
func (mat *MMat3) Adjugated() *MMat3 {
	result := *mat
	result.Adjugate()
	return &result
}

// Invert inverts the given matrix. Destructive operation.
// Does not check if matrix is singular and may lead to strange results!
func (mat *MMat3) Invert() *MMat3 {
	initialDet := mat.Determinant()
	if initialDet == 0 {
		mat = NewMMat3()
	}

	mat.Adjugate()
	mat.Mul(1.0 / initialDet)
	return mat
}

// Inverted inverts a copy of the given matrix.
// Does not check if matrix is singular and may lead to strange results!
func (mat *MMat3) Inverted() *MMat3 {
	result := *mat
	return result.Invert()
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func (v *MMat3) ClampIfVerySmall() *MMat3 {
	epsilon := 1e-6
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(v[i][j]) < epsilon {
				v[i][j] = 0
			}
		}
	}
	return v
}
