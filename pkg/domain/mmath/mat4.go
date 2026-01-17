// 指示: miu200521358
package mmath

import (
	"fmt"
	"math"

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
	var m Mat4
	m.SetByValues(m11, m21, m31, m41, m12, m22, m32, m42, m13, m23, m33, m43, m14, m24, m34, m44)
	return m
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
	a00, a01, a02, a03 := m[0], m[1], m[2], m[3]
	a10, a11, a12, a13 := m[4], m[5], m[6], m[7]
	a20, a21, a22, a23 := m[8], m[9], m[10], m[11]
	a30, a31, a32, a33 := m[12], m[13], m[14], m[15]
	m[0], m[1], m[2], m[3] = a00, a10, a20, a30
	m[4], m[5], m[6], m[7] = a01, a11, a21, a31
	m[8], m[9], m[10], m[11] = a02, a12, a22, a32
	m[12], m[13], m[14], m[15] = a03, a13, a23, a33
	return m
}

// Mul は乗算する。
func (m *Mat4) Mul(other Mat4) *Mat4 {
	m.MulTo(other, m)
	return m
}

// MulTo は乗算結果をoutへ書き込む。
func (m Mat4) MulTo(other Mat4, out *Mat4) {
	if out == nil {
		return
	}
	a00, a01, a02, a03 := m[0], m[1], m[2], m[3]
	a10, a11, a12, a13 := m[4], m[5], m[6], m[7]
	a20, a21, a22, a23 := m[8], m[9], m[10], m[11]
	a30, a31, a32, a33 := m[12], m[13], m[14], m[15]

	b00, b01, b02, b03 := other[0], other[1], other[2], other[3]
	b10, b11, b12, b13 := other[4], other[5], other[6], other[7]
	b20, b21, b22, b23 := other[8], other[9], other[10], other[11]
	b30, b31, b32, b33 := other[12], other[13], other[14], other[15]

	out[0] = a00*b00 + a10*b01 + a20*b02 + a30*b03
	out[1] = a01*b00 + a11*b01 + a21*b02 + a31*b03
	out[2] = a02*b00 + a12*b01 + a22*b02 + a32*b03
	out[3] = a03*b00 + a13*b01 + a23*b02 + a33*b03
	out[4] = a00*b10 + a10*b11 + a20*b12 + a30*b13
	out[5] = a01*b10 + a11*b11 + a21*b12 + a31*b13
	out[6] = a02*b10 + a12*b11 + a22*b12 + a32*b13
	out[7] = a03*b10 + a13*b11 + a23*b12 + a33*b13
	out[8] = a00*b20 + a10*b21 + a20*b22 + a30*b23
	out[9] = a01*b20 + a11*b21 + a21*b22 + a31*b23
	out[10] = a02*b20 + a12*b21 + a22*b22 + a32*b23
	out[11] = a03*b20 + a13*b21 + a23*b22 + a33*b23
	out[12] = a00*b30 + a10*b31 + a20*b32 + a30*b33
	out[13] = a01*b30 + a11*b31 + a21*b32 + a31*b33
	out[14] = a02*b30 + a12*b31 + a22*b32 + a32*b33
	out[15] = a03*b30 + a13*b31 + a23*b32 + a33*b33
}

// MulToPtr は乗算結果をoutへ書き込む。
func (m *Mat4) MulToPtr(other *Mat4, out *Mat4) {
	if m == nil || other == nil || out == nil {
		return
	}
	a00, a01, a02, a03 := m[0], m[1], m[2], m[3]
	a10, a11, a12, a13 := m[4], m[5], m[6], m[7]
	a20, a21, a22, a23 := m[8], m[9], m[10], m[11]
	a30, a31, a32, a33 := m[12], m[13], m[14], m[15]

	b00, b01, b02, b03 := other[0], other[1], other[2], other[3]
	b10, b11, b12, b13 := other[4], other[5], other[6], other[7]
	b20, b21, b22, b23 := other[8], other[9], other[10], other[11]
	b30, b31, b32, b33 := other[12], other[13], other[14], other[15]

	out[0] = a00*b00 + a10*b01 + a20*b02 + a30*b03
	out[1] = a01*b00 + a11*b01 + a21*b02 + a31*b03
	out[2] = a02*b00 + a12*b01 + a22*b02 + a32*b03
	out[3] = a03*b00 + a13*b01 + a23*b02 + a33*b03
	out[4] = a00*b10 + a10*b11 + a20*b12 + a30*b13
	out[5] = a01*b10 + a11*b11 + a21*b12 + a31*b13
	out[6] = a02*b10 + a12*b11 + a22*b12 + a32*b13
	out[7] = a03*b10 + a13*b11 + a23*b12 + a33*b13
	out[8] = a00*b20 + a10*b21 + a20*b22 + a30*b23
	out[9] = a01*b20 + a11*b21 + a21*b22 + a31*b23
	out[10] = a02*b20 + a12*b21 + a22*b22 + a32*b23
	out[11] = a03*b20 + a13*b21 + a23*b22 + a33*b23
	out[12] = a00*b30 + a10*b31 + a20*b32 + a30*b33
	out[13] = a01*b30 + a11*b31 + a21*b32 + a31*b33
	out[14] = a02*b30 + a12*b31 + a22*b32 + a32*b33
	out[15] = a03*b30 + a13*b31 + a23*b32 + a33*b33
}

// MulTranslateTo は平行移動の乗算結果をoutへ書き込む。
func (m *Mat4) MulTranslateTo(v *Vec3, out *Mat4) {
	if m == nil || v == nil || out == nil {
		return
	}
	a00, a01, a02, a03 := m[0], m[1], m[2], m[3]
	a10, a11, a12, a13 := m[4], m[5], m[6], m[7]
	a20, a21, a22, a23 := m[8], m[9], m[10], m[11]
	a30, a31, a32, a33 := m[12], m[13], m[14], m[15]
	x, y, z := v.X, v.Y, v.Z

	out[0], out[1], out[2], out[3] = a00, a01, a02, a03
	out[4], out[5], out[6], out[7] = a10, a11, a12, a13
	out[8], out[9], out[10], out[11] = a20, a21, a22, a23
	out[12] = a00*x + a10*y + a20*z + a30
	out[13] = a01*x + a11*y + a21*z + a31
	out[14] = a02*x + a12*y + a22*z + a32
	out[15] = a03*x + a13*y + a23*z + a33
}

// Muled は乗算結果を返す。
func (m Mat4) Muled(other Mat4) Mat4 {
	var out Mat4
	m.MulTo(other, &out)
	return out
}

// SetByValues は値を設定する。
func (m *Mat4) SetByValues(m11, m21, m31, m41, m12, m22, m32, m42, m13, m23, m33, m43, m14, m24, m34, m44 float64) *Mat4 {
	if m == nil {
		return nil
	}
	m[0], m[1], m[2], m[3] = m11, m21, m31, m41
	m[4], m[5], m[6], m[7] = m12, m22, m32, m42
	m[8], m[9], m[10], m[11] = m13, m23, m33, m43
	m[12], m[13], m[14], m[15] = m14, m24, m34, m44
	return m
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
	return mat4Det(m)
}

// Inverse は逆行列にする。
func (m *Mat4) Inverse() *Mat4 {
	inv := m.Inverted()
	*m = inv
	return m
}

// Inverted は逆行列を返す。
func (m Mat4) Inverted() Mat4 {
	det := mat4Det(m)
	if math.IsNaN(det) || math.IsInf(det, 0) || math.Abs(det) < 1e-10 {
		return NewMat4()
	}
	inv, ok := mat4Inverse(m, det)
	if !ok {
		return NewMat4()
	}
	return inv
}

// InvertedTo は逆行列をoutへ書き込む。
func (m *Mat4) InvertedTo(out *Mat4) {
	if m == nil || out == nil {
		return
	}
	det := mat4Det(*m)
	if math.IsNaN(det) || math.IsInf(det, 0) || math.Abs(det) < 1e-10 {
		*out = NewMat4()
		return
	}
	inv, ok := mat4Inverse(*m, det)
	if !ok {
		*out = NewMat4()
		return
	}
	*out = inv
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

func mat4Det(m Mat4) float64 {
	a00, a01, a02, a03 := m[0], m[1], m[2], m[3]
	a10, a11, a12, a13 := m[4], m[5], m[6], m[7]
	a20, a21, a22, a23 := m[8], m[9], m[10], m[11]
	a30, a31, a32, a33 := m[12], m[13], m[14], m[15]

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	return b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06
}

func mat4Inverse(m Mat4, det float64) (Mat4, bool) {
	a00, a01, a02, a03 := m[0], m[1], m[2], m[3]
	a10, a11, a12, a13 := m[4], m[5], m[6], m[7]
	a20, a21, a22, a23 := m[8], m[9], m[10], m[11]
	a30, a31, a32, a33 := m[12], m[13], m[14], m[15]

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	if det == 0 {
		return Mat4{}, false
	}
	invDet := 1.0 / det

	inv := Mat4{
		(a11*b11 - a12*b10 + a13*b09) * invDet,
		(a02*b10 - a01*b11 - a03*b09) * invDet,
		(a31*b05 - a32*b04 + a33*b03) * invDet,
		(a22*b04 - a21*b05 - a23*b03) * invDet,
		(a12*b08 - a10*b11 - a13*b07) * invDet,
		(a00*b11 - a02*b08 + a03*b07) * invDet,
		(a32*b02 - a30*b05 - a33*b01) * invDet,
		(a20*b05 - a22*b02 + a23*b01) * invDet,
		(a10*b10 - a11*b08 + a13*b06) * invDet,
		(a01*b08 - a00*b10 - a03*b06) * invDet,
		(a30*b04 - a31*b02 + a33*b00) * invDet,
		(a21*b02 - a20*b04 - a23*b00) * invDet,
		(a11*b07 - a10*b09 - a12*b06) * invDet,
		(a00*b09 - a01*b07 + a02*b06) * invDet,
		(a31*b01 - a30*b03 - a32*b00) * invDet,
		(a20*b03 - a21*b01 + a22*b00) * invDet,
	}

	if math.IsNaN(inv[0]) || math.IsInf(inv[0], 0) {
		return Mat4{}, false
	}
	return inv, true
}
