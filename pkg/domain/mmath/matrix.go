package mmath

import (
	"fmt"
	"math"
)

// ----- 定数 -----

// MAT4_IDENTITY は単位行列です
var MAT4_IDENTITY = NewMat4()

// ----- 型定義 -----

// Mat4 は4x4行列を表します（列優先順序）
// OpenGL互換のため、インデックスは以下の通り:
//
//	[0]  [4]  [8]  [12]
//	[1]  [5]  [9]  [13]
//	[2]  [6]  [10] [14]
//	[3]  [7]  [11] [15]
type Mat4 [16]float64

// ----- コンストラクタ -----

// NewMat4 は単位行列を作成します
func NewMat4() *Mat4 {
	return &Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// NewMat4ByValues は指定した値で行列を作成します
func NewMat4ByValues(
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float64,
) *Mat4 {
	return &Mat4{
		m00, m10, m20, m30,
		m01, m11, m21, m31,
		m02, m12, m22, m32,
		m03, m13, m23, m33,
	}
}

// ----- 文字列表現 -----

func (m *Mat4) String() string {
	return fmt.Sprintf(
		"[[%.4f, %.4f, %.4f, %.4f], [%.4f, %.4f, %.4f, %.4f], [%.4f, %.4f, %.4f, %.4f], [%.4f, %.4f, %.4f, %.4f]]",
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	)
}

// ----- 要素アクセス -----

// At は指定した行・列の要素を返します
func (m *Mat4) At(row, col int) float64 {
	return m[col*4+row]
}

// Set は指定した行・列の要素を設定します
func (m *Mat4) Set(row, col int, val float64) {
	m[col*4+row] = val
}

// ----- 算術演算（破壊的） -----

// Mul は行列を乗算します（破壊的）
func (m *Mat4) Mul(other *Mat4) *Mat4 {
	result := Mat4{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := 0.0
			for k := 0; k < 4; k++ {
				sum += m[k*4+i] * other[j*4+k]
			}
			result[j*4+i] = sum
		}
	}
	*m = result
	return m
}

// ----- 算術演算（非破壊的） -----

// Muled は行列を乗算した結果を返します
func (m *Mat4) Muled(other *Mat4) *Mat4 {
	result := &Mat4{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := 0.0
			for k := 0; k < 4; k++ {
				sum += m[k*4+i] * other[j*4+k]
			}
			result[j*4+i] = sum
		}
	}
	return result
}

// ----- 逆行列 -----

// Inverse は逆行列を計算します（破壊的）
func (m *Mat4) Inverse() *Mat4 {
	inv := m.Inverted()
	*m = *inv
	return m
}

// Inverted は逆行列を返します
func (m *Mat4) Inverted() *Mat4 {
	inv := &Mat4{}

	inv[0] = m[5]*m[10]*m[15] - m[5]*m[11]*m[14] - m[9]*m[6]*m[15] + m[9]*m[7]*m[14] + m[13]*m[6]*m[11] - m[13]*m[7]*m[10]
	inv[4] = -m[4]*m[10]*m[15] + m[4]*m[11]*m[14] + m[8]*m[6]*m[15] - m[8]*m[7]*m[14] - m[12]*m[6]*m[11] + m[12]*m[7]*m[10]
	inv[8] = m[4]*m[9]*m[15] - m[4]*m[11]*m[13] - m[8]*m[5]*m[15] + m[8]*m[7]*m[13] + m[12]*m[5]*m[11] - m[12]*m[7]*m[9]
	inv[12] = -m[4]*m[9]*m[14] + m[4]*m[10]*m[13] + m[8]*m[5]*m[14] - m[8]*m[6]*m[13] - m[12]*m[5]*m[10] + m[12]*m[6]*m[9]
	inv[1] = -m[1]*m[10]*m[15] + m[1]*m[11]*m[14] + m[9]*m[2]*m[15] - m[9]*m[3]*m[14] - m[13]*m[2]*m[11] + m[13]*m[3]*m[10]
	inv[5] = m[0]*m[10]*m[15] - m[0]*m[11]*m[14] - m[8]*m[2]*m[15] + m[8]*m[3]*m[14] + m[12]*m[2]*m[11] - m[12]*m[3]*m[10]
	inv[9] = -m[0]*m[9]*m[15] + m[0]*m[11]*m[13] + m[8]*m[1]*m[15] - m[8]*m[3]*m[13] - m[12]*m[1]*m[11] + m[12]*m[3]*m[9]
	inv[13] = m[0]*m[9]*m[14] - m[0]*m[10]*m[13] - m[8]*m[1]*m[14] + m[8]*m[2]*m[13] + m[12]*m[1]*m[10] - m[12]*m[2]*m[9]
	inv[2] = m[1]*m[6]*m[15] - m[1]*m[7]*m[14] - m[5]*m[2]*m[15] + m[5]*m[3]*m[14] + m[13]*m[2]*m[7] - m[13]*m[3]*m[6]
	inv[6] = -m[0]*m[6]*m[15] + m[0]*m[7]*m[14] + m[4]*m[2]*m[15] - m[4]*m[3]*m[14] - m[12]*m[2]*m[7] + m[12]*m[3]*m[6]
	inv[10] = m[0]*m[5]*m[15] - m[0]*m[7]*m[13] - m[4]*m[1]*m[15] + m[4]*m[3]*m[13] + m[12]*m[1]*m[7] - m[12]*m[3]*m[5]
	inv[14] = -m[0]*m[5]*m[14] + m[0]*m[6]*m[13] + m[4]*m[1]*m[14] - m[4]*m[2]*m[13] - m[12]*m[1]*m[6] + m[12]*m[2]*m[5]
	inv[3] = -m[1]*m[6]*m[11] + m[1]*m[7]*m[10] + m[5]*m[2]*m[11] - m[5]*m[3]*m[10] - m[9]*m[2]*m[7] + m[9]*m[3]*m[6]
	inv[7] = m[0]*m[6]*m[11] - m[0]*m[7]*m[10] - m[4]*m[2]*m[11] + m[4]*m[3]*m[10] + m[8]*m[2]*m[7] - m[8]*m[3]*m[6]
	inv[11] = -m[0]*m[5]*m[11] + m[0]*m[7]*m[9] + m[4]*m[1]*m[11] - m[4]*m[3]*m[9] - m[8]*m[1]*m[7] + m[8]*m[3]*m[5]
	inv[15] = m[0]*m[5]*m[10] - m[0]*m[6]*m[9] - m[4]*m[1]*m[10] + m[4]*m[2]*m[9] + m[8]*m[1]*m[6] - m[8]*m[2]*m[5]

	det := m[0]*inv[0] + m[1]*inv[4] + m[2]*inv[8] + m[3]*inv[12]
	if det == 0 {
		return NewMat4()
	}

	det = 1.0 / det
	for i := 0; i < 16; i++ {
		inv[i] *= det
	}

	return inv
}

// ----- 平行移動・回転・スケール -----

// Translation は平行移動成分を返します
func (m *Mat4) Translation() *Vec3 {
	return NewVec3ByValues(m[12], m[13], m[14])
}

// SetTranslation は平行移動成分を設定します
func (m *Mat4) SetTranslation(v *Vec3) {
	m[12] = v.X
	m[13] = v.Y
	m[14] = v.Z
}

// Translate は平行移動を適用します（破壊的）
func (m *Mat4) Translate(v *Vec3) *Mat4 {
	t := NewMat4()
	t[12] = v.X
	t[13] = v.Y
	t[14] = v.Z
	return m.Mul(t)
}

// Scale はスケールを適用します（破壊的）
func (m *Mat4) Scale(v *Vec3) *Mat4 {
	s := NewMat4()
	s[0] = v.X
	s[5] = v.Y
	s[10] = v.Z
	return m.Mul(s)
}

// RotateX はX軸周りの回転を適用します（破壊的）
func (m *Mat4) RotateX(angle float64) *Mat4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	r := NewMat4()
	r[5] = c
	r[6] = s
	r[9] = -s
	r[10] = c
	return m.Mul(r)
}

// RotateY はY軸周りの回転を適用します（破壊的）
func (m *Mat4) RotateY(angle float64) *Mat4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	r := NewMat4()
	r[0] = c
	r[2] = -s
	r[8] = s
	r[10] = c
	return m.Mul(r)
}

// RotateZ はZ軸周りの回転を適用します（破壊的）
func (m *Mat4) RotateZ(angle float64) *Mat4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	r := NewMat4()
	r[0] = c
	r[1] = s
	r[4] = -s
	r[5] = c
	return m.Mul(r)
}

// ----- ベクトル変換 -----

// MulVec3 はベクトルを行列で変換します（w=1として）
func (m *Mat4) MulVec3(v *Vec3) *Vec3 {
	return NewVec3ByValues(
		m[0]*v.X+m[4]*v.Y+m[8]*v.Z+m[12],
		m[1]*v.X+m[5]*v.Y+m[9]*v.Z+m[13],
		m[2]*v.X+m[6]*v.Y+m[10]*v.Z+m[14],
	)
}

// MulVec4 はベクトルを行列で変換します
func (m *Mat4) MulVec4(v *Vec4) *Vec4 {
	return NewVec4ByValues(
		m[0]*v.X+m[4]*v.Y+m[8]*v.Z+m[12]*v.W,
		m[1]*v.X+m[5]*v.Y+m[9]*v.Z+m[13]*v.W,
		m[2]*v.X+m[6]*v.Y+m[10]*v.Z+m[14]*v.W,
		m[3]*v.X+m[7]*v.Y+m[11]*v.Z+m[15]*v.W,
	)
}

// ----- ユーティリティ -----

// Copy はコピーを返します
func (m *Mat4) Copy() *Mat4 {
	result := *m
	return &result
}

// IsIdent は単位行列かどうかを返します
func (m *Mat4) IsIdent() bool {
	return m.NearEquals(MAT4_IDENTITY, 1e-10)
}

// NearEquals は他の行列とほぼ等しいかどうかを返します
func (m *Mat4) NearEquals(other *Mat4, epsilon float64) bool {
	for i := 0; i < 16; i++ {
		if math.Abs(m[i]-other[i]) > epsilon {
			return false
		}
	}
	return true
}

// GL は行列をOpenGL用のfloat32スライスとして返します
func (m *Mat4) GL() []float32 {
	return []float32{
		float32(m[0]), float32(m[1]), float32(m[2]), float32(m[3]),
		float32(m[4]), float32(m[5]), float32(m[6]), float32(m[7]),
		float32(m[8]), float32(m[9]), float32(m[10]), float32(m[11]),
		float32(m[12]), float32(m[13]), float32(m[14]), float32(m[15]),
	}
}
