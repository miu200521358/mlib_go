package mmath

import (
	"github.com/go-gl/mathgl/mgl64"
)

type MMat3 mgl64.Mat3

func NewMMat3() *MMat3 {
	return &MMat3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

var (
	// Zero holds a zero matrix.
	MMat3Zero = MMat3{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}

	// Ident holds an ident matrix.
	MMat3Ident = MMat3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
)

// IsZero
func (m *MMat3) IsZero() bool {
	return *m == MMat3Zero
}

// IsIdent
func (m *MMat3) IsIdent() bool {
	return m.NearEquals(&MMat3Ident, 1e-10)
}

// String
func (m *MMat3) String() string {
	return mgl64.Mat3(*m).String()
}

func (m *MMat3) Copy() *MMat3 {
	copied := NewMMat3ByValues(m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8])
	return copied
}

func NewMMat3ByValues(m11, m12, m13, m21, m22, m23, m31, m32, m33 float64) *MMat3 {
	return (*MMat3)(&mgl64.Mat3{m11, m12, m13, m21, m22, m23, m31, m32, m33})
}

func (m *MMat3) NearEquals(other *MMat3, tolerance float64) bool {
	return mgl64.Mat3(*m).ApproxEqualThreshold(mgl64.Mat3(*other), tolerance)
}

// Mul は行列の掛け算を行います
func (m1 *MMat3) Mul(m2 *MMat3) {
	m := mgl64.Mat3(*m1).Mul3(mgl64.Mat3(*m2))
	*m1 = MMat3(m)
}
