// 指示: miu200521358
package mmath

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

type Mat4 [16]float64

var (
	ZERO_MAT4 = Mat4{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
	IDENT_MAT4 = Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	IDENT_SCALE_MAT4 = Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
	}
)

// NewMat4 はMat4を生成する。
func NewMat4() Mat4 {
	return IDENT_MAT4
}

// NewMat4ByValues はMat4を生成する。
func NewMat4ByValues(m11, m21, m31, m41, m12, m22, m32, m42, m13, m23, m33, m43, m14, m24, m34, m44 float64) Mat4 {
	return Mat4{
		m11, m21, m31, m41,
		m12, m22, m32, m42,
		m13, m23, m33, m43,
		m14, m24, m34, m44,
	}
}

// NewMat4FromAxisAngle はMat4を生成する。
func NewMat4FromAxisAngle(axis Vec3, angle float64) Mat4 {
	rot := r3.NewRotation(angle, axis.Vec)
	qq := Quaternion{quat.Number(rot)}
	return qq.ToMat4()
}

// NewMat4FromLookAt はMat4を生成する。
func NewMat4FromLookAt(eye, center, up Vec3) Mat4 {
	f := center.Subed(eye).Normalized()
	s := f.Cross(up).Normalized()
	u := s.Cross(f)

	m := NewMat4ByValues(
		s.X, s.Y, s.Z, 0,
		u.X, u.Y, u.Z, 0,
		-f.X, -f.Y, -f.Z, 0,
		0, 0, 0, 1,
	)

	m[12] = -s.Dot(eye)
	m[13] = -u.Dot(eye)
	m[14] = f.Dot(eye)
	return m
}

// IsZero はゼロか判定する。
func (m Mat4) IsZero() bool {
	return m == ZERO_MAT4
}

// IsIdent は単位行列か判定する。
func (m Mat4) IsIdent() bool {
	return m.NearEquals(IDENT_MAT4, 1e-10)
}

// String は文字列表現を返す。
func (m Mat4) String() string {
	return fmt.Sprintf("[%g %g %g %g; %g %g %g %g; %g %g %g %g; %g %g %g %g]",
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	)
}

// Copy はコピーを返す。
func (m Mat4) Copy() (Mat4, error) {
	return deepCopy(m)
}

// NearEquals は近似的に等しいか判定する。
func (m Mat4) NearEquals(other Mat4, tolerance float64) bool {
	for i := range m {
		if math.Abs(m[i]-other[i]) > tolerance {
			return false
		}
	}
	return true
}

// Trace はトレースを返す。
func (m Mat4) Trace() float64 {
	return m[0] + m[5] + m[10] + m[15]
}

// Trace3 は3x3のトレースを返す。
func (m Mat4) Trace3() float64 {
	return m[0] + m[5] + m[10]
}

// MulVec3 はベクトルを変換する。
func (m Mat4) MulVec3(other Vec3) Vec3 {
	x := other.X*m[0] + other.Y*m[4] + other.Z*m[8] + m[12]
	y := other.X*m[1] + other.Y*m[5] + other.Z*m[9] + m[13]
	z := other.X*m[2] + other.Y*m[6] + other.Z*m[10] + m[14]
	w := other.X*m[3] + other.Y*m[7] + other.Z*m[11] + m[15]

	if w != 0 && w != 1 {
		invW := 1.0 / w
		return Vec3{r3.Vec{X: x * invW, Y: y * invW, Z: z * invW}}
	}
	return Vec3{r3.Vec{X: x, Y: y, Z: z}}
}

// Translate は平行移動する。
func (m *Mat4) Translate(v Vec3) *Mat4 {
	result := v.ToMat4().Muled(*m)
	*m = result
	return m
}

// Translated は平行移動結果を返す。
func (m Mat4) Translated(v Vec3) Mat4 {
	return v.ToMat4().Muled(m)
}

// Translation は平行移動成分を返す。
func (m Mat4) Translation() Vec3 {
	return Vec3{r3.Vec{X: m[12], Y: m[13], Z: m[14]}}
}

// Scale は拡大縮小する。
func (m *Mat4) Scale(s Vec3) *Mat4 {
	result := s.ToScaleMat4().Muled(*m)
	*m = result
	return m
}

// Scaled は拡大縮小結果を返す。
func (m Mat4) Scaled(s Vec3) Mat4 {
	return s.ToScaleMat4().Muled(m)
}

// Scaling は拡大縮小成分を返す。
func (m Mat4) Scaling() Vec3 {
	return Vec3{r3.Vec{X: m[0], Y: m[5], Z: m[10]}}
}

// Rotate は回転する。
func (m *Mat4) Rotate(q Quaternion) *Mat4 {
	result := q.ToMat4().Muled(*m)
	*m = result
	return m
}

// Rotated は回転結果を返す。
func (m Mat4) Rotated(q Quaternion) Mat4 {
	return q.ToMat4().Muled(m)
}

// Quaternion はクォータニオンを返す。
func (m Mat4) Quaternion() Quaternion {
	trace := m[0] + m[5] + m[10] + 1.0

	var x, y, z, w float64
	if trace > 1e-5 {
		s := 0.5 / math.Sqrt(trace)
		w = 0.25 / s
		x = (m[9] - m[6]) * s
		y = (m[2] - m[8]) * s
		z = (m[4] - m[1]) * s
	} else if m[0] > m[5] && m[0] > m[10] {
		s := 2.0 * math.Sqrt(1.0+m[0]-m[5]-m[10])
		x = 0.25 * s
		y = (m[1] + m[4]) / s
		z = (m[2] + m[8]) / s
		w = (m[9] - m[6]) / s
	} else if m[5] > m[10] {
		s := 2.0 * math.Sqrt(1.0+m[5]-m[0]-m[10])
		x = (m[1] + m[4]) / s
		y = 0.25 * s
		z = (m[6] + m[9]) / s
		w = (m[2] - m[8]) / s
	} else {
		s := 2.0 * math.Sqrt(1.0+m[10]-m[0]-m[5])
		x = (m[2] + m[8]) / s
		y = (m[6] + m[9]) / s
		z = 0.25 * s
		w = (m[4] - m[1]) / s
	}

	return NewQuaternionByValues(-x, -y, -z, w)
}

// Transpose は転置する。
func (m *Mat4) Transpose() *Mat4 {
	*m = Mat4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
	return m
}

// Mul は乗算する。
func (m *Mat4) Mul(other Mat4) *Mat4 {
	*m = m.Muled(other)
	return m
}

// Muled は乗算結果を返す。
func (m Mat4) Muled(other Mat4) Mat4 {
	return Mat4{
		m[0]*other[0] + m[4]*other[1] + m[8]*other[2] + m[12]*other[3],
		m[1]*other[0] + m[5]*other[1] + m[9]*other[2] + m[13]*other[3],
		m[2]*other[0] + m[6]*other[1] + m[10]*other[2] + m[14]*other[3],
		m[3]*other[0] + m[7]*other[1] + m[11]*other[2] + m[15]*other[3],

		m[0]*other[4] + m[4]*other[5] + m[8]*other[6] + m[12]*other[7],
		m[1]*other[4] + m[5]*other[5] + m[9]*other[6] + m[13]*other[7],
		m[2]*other[4] + m[6]*other[5] + m[10]*other[6] + m[14]*other[7],
		m[3]*other[4] + m[7]*other[5] + m[11]*other[6] + m[15]*other[7],

		m[0]*other[8] + m[4]*other[9] + m[8]*other[10] + m[12]*other[11],
		m[1]*other[8] + m[5]*other[9] + m[9]*other[10] + m[13]*other[11],
		m[2]*other[8] + m[6]*other[9] + m[10]*other[10] + m[14]*other[11],
		m[3]*other[8] + m[7]*other[9] + m[11]*other[10] + m[15]*other[11],

		m[0]*other[12] + m[4]*other[13] + m[8]*other[14] + m[12]*other[15],
		m[1]*other[12] + m[5]*other[13] + m[9]*other[14] + m[13]*other[15],
		m[2]*other[12] + m[6]*other[13] + m[10]*other[14] + m[14]*other[15],
		m[3]*other[12] + m[7]*other[13] + m[11]*other[14] + m[15]*other[15],
	}
}

// Add は加算する。
func (m *Mat4) Add(other Mat4) *Mat4 {
	for i := range m {
		m[i] += other[i]
	}
	return m
}

// Added は加算結果を返す。
func (m Mat4) Added(other Mat4) Mat4 {
	result := m
	for i := range result {
		result[i] += other[i]
	}
	return result
}

// MulScalar はスカラーを乗算する。
func (m *Mat4) MulScalar(v float64) *Mat4 {
	for i := range m {
		m[i] *= v
	}
	return m
}

// MuledScalar はスカラー乗算結果を返す。
func (m Mat4) MuledScalar(v float64) Mat4 {
	result := m
	for i := range result {
		result[i] *= v
	}
	return result
}

// Det は行列式を返す。
func (m Mat4) Det() float64 {
	return mat.Det(m.toDense())
}

// Inverse は逆行列にする。
func (m *Mat4) Inverse() *Mat4 {
	inv := m.Inverted()
	*m = inv
	return m
}

// Inverted は逆行列を返す。
func (m Mat4) Inverted() Mat4 {
	det := m.Det()
	if math.Abs(det) < 1e-10 {
		return NewMat4()
	}
	var inv mat.Dense
	if err := inv.Inverse(m.toDense()); err != nil || math.IsNaN(inv.At(0, 0)) || math.IsInf(inv.At(0, 0), 0) {
		return NewMat4()
	}
	return mat4FromDense(&inv)
}

// ClampIfVerySmall は微小値を0に丸める。
func (m *Mat4) ClampIfVerySmall() *Mat4 {
	epsilon := 1e-6
	for i := range m {
		if math.Abs(m[i]) < epsilon {
			m[i] = 0
		}
	}
	return m
}

// AxisX はX軸ベクトルを返す。
func (m Mat4) AxisX() Vec3 {
	return Vec3{r3.Vec{X: m[0], Y: m[1], Z: m[2]}}
}

// AxisY はY軸ベクトルを返す。
func (m Mat4) AxisY() Vec3 {
	return Vec3{r3.Vec{X: m[4], Y: m[5], Z: m[6]}}
}

// AxisZ はZ軸ベクトルを返す。
func (m Mat4) AxisZ() Vec3 {
	return Vec3{r3.Vec{X: m[8], Y: m[9], Z: m[10]}}
}

// toDense はmat.Denseへ変換する。
func (m Mat4) toDense() *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	})
}

// mat4FromDense はmat.DenseからMat4へ変換する。
func mat4FromDense(d *mat.Dense) Mat4 {
	return Mat4{
		d.At(0, 0), d.At(1, 0), d.At(2, 0), d.At(3, 0),
		d.At(0, 1), d.At(1, 1), d.At(2, 1), d.At(3, 1),
		d.At(0, 2), d.At(1, 2), d.At(2, 2), d.At(3, 2),
		d.At(0, 3), d.At(1, 3), d.At(2, 3), d.At(3, 3),
	}
}

