package mmath

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/num/quat"
)

// ----- 定数 -----

var (
	// QUATERNION_ZERO はゼロクォータニオンです
	QUATERNION_ZERO = &Quaternion{}

	// QUATERNION_IDENTITY は単位クォータニオンです
	QUATERNION_IDENTITY = &Quaternion{Number: quat.Number{Real: 1}}
)

// ----- 型定義 -----

// Quaternion はクォータニオンを表します（gonum/num/quat.Numberを埋め込み）
type Quaternion struct {
	quat.Number
}

// X はX成分を返します
func (q *Quaternion) X() float64 { return q.Imag }

// Y はY成分を返します
func (q *Quaternion) Y() float64 { return q.Jmag }

// Z はZ成分を返します
func (q *Quaternion) Z() float64 { return q.Kmag }

// W はW成分を返します
func (q *Quaternion) W() float64 { return q.Real }

// SetX はX成分を設定します
func (q *Quaternion) SetX(v float64) { q.Imag = v }

// SetY はY成分を設定します
func (q *Quaternion) SetY(v float64) { q.Jmag = v }

// SetZ はZ成分を設定します
func (q *Quaternion) SetZ(v float64) { q.Kmag = v }

// SetW はW成分を設定します
func (q *Quaternion) SetW(v float64) { q.Real = v }

// ----- コンストラクタ -----

// NewQuaternion は単位クォータニオンを作成します
func NewQuaternion() *Quaternion {
	return &Quaternion{Number: quat.Number{Real: 1}}
}

// NewQuaternionByValues は指定した値でクォータニオンを作成します
func NewQuaternionByValues(x, y, z, w float64) *Quaternion {
	return &Quaternion{Number: quat.Number{Real: w, Imag: x, Jmag: y, Kmag: z}}
}

// NewQuaternionFromRadians はオイラー角（ラジアン）からクォータニオンを作成します
func NewQuaternionFromRadians(xPitch, yHead, zRoll float64) *Quaternion {
	cx := math.Cos(xPitch * 0.5)
	sx := math.Sin(xPitch * 0.5)
	cy := math.Cos(yHead * 0.5)
	sy := math.Sin(yHead * 0.5)
	cz := math.Cos(zRoll * 0.5)
	sz := math.Sin(zRoll * 0.5)

	return NewQuaternionByValues(
		sx*cy*cz+cx*sy*sz,
		cx*sy*cz-sx*cy*sz,
		cx*cy*sz+sx*sy*cz,
		cx*cy*cz-sx*sy*sz,
	)
}

// NewQuaternionFromDegrees はオイラー角（度）からクォータニオンを作成します
func NewQuaternionFromDegrees(xPitch, yHead, zRoll float64) *Quaternion {
	return NewQuaternionFromRadians(DegToRad(xPitch), DegToRad(yHead), DegToRad(zRoll))
}

// NewQuaternionFromAxisAngle は軸と角度からクォータニオンを作成します
func NewQuaternionFromAxisAngle(axis *Vec3, angle float64) *Quaternion {
	n := axis.Normalized()
	s := math.Sin(angle * 0.5)
	c := math.Cos(angle * 0.5)
	return NewQuaternionByValues(n.X*s, n.Y*s, n.Z*s, c)
}

// NewQuaternionRotate は2つのベクトル間の回転を表すクォータニオンを作成します
func NewQuaternionRotate(from, to *Vec3) *Quaternion {
	if from.NearEquals(to, 1e-6) || from.Length() == 0 || to.Length() == 0 {
		return NewQuaternion()
	}

	fn := from.Normalized()
	tn := to.Normalized()

	dot := fn.Dot(tn)
	if dot > 0.9999 {
		return NewQuaternion()
	}
	if dot < -0.9999 {
		// 180度回転
		axis := VEC3_UNIT_X.Cross(fn)
		if axis.Length() < 1e-6 {
			axis = VEC3_UNIT_Y.Cross(fn)
		}
		return NewQuaternionFromAxisAngle(axis.Normalized(), math.Pi)
	}

	axis := fn.Cross(tn)
	s := math.Sqrt((1 + dot) * 2)
	invS := 1 / s

	return NewQuaternionByValues(axis.X*invS, axis.Y*invS, axis.Z*invS, s*0.5)
}

// ----- 文字列表現 -----

func (q *Quaternion) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", q.X(), q.Y(), q.Z(), q.W())
}

// ----- 算術演算（破壊的） -----

// Mul はクォータニオンを乗算します（破壊的）
func (q *Quaternion) Mul(other *Quaternion) *Quaternion {
	q.Number = quat.Mul(q.Number, other.Number)
	return q
}

// Negate は符号を反転します（破壊的）
func (q *Quaternion) Negate() *Quaternion {
	q.Number = quat.Scale(-1, q.Number)
	return q
}

// Normalize は正規化します（破壊的）
func (q *Quaternion) Normalize() *Quaternion {
	n := quat.Abs(q.Number)
	if n == 0 {
		*q = *NewQuaternion()
		return q
	}
	q.Number = quat.Scale(1/n, q.Number)
	return q
}

// Inverse は逆クォータニオンにします（破壊的）
func (q *Quaternion) Inverse() *Quaternion {
	q.Number = quat.Inv(q.Number)
	return q
}

// ----- 算術演算（非破壊的） -----

// Muled はクォータニオンを乗算した結果を返します
func (q *Quaternion) Muled(other *Quaternion) *Quaternion {
	return &Quaternion{Number: quat.Mul(q.Number, other.Number)}
}

// Negated は符号を反転した結果を返します
func (q *Quaternion) Negated() *Quaternion {
	return &Quaternion{Number: quat.Scale(-1, q.Number)}
}

// Normalized は正規化した結果を返します
func (q *Quaternion) Normalized() *Quaternion {
	n := quat.Abs(q.Number)
	if n == 0 {
		return NewQuaternion()
	}
	return &Quaternion{Number: quat.Scale(1/n, q.Number)}
}

// Inverted は逆クォータニオンを返します
func (q *Quaternion) Inverted() *Quaternion {
	return &Quaternion{Number: quat.Inv(q.Number)}
}

// Conjugate は共役クォータニオンにします（破壊的）
func (q *Quaternion) Conjugate() *Quaternion {
	q.Number = quat.Conj(q.Number)
	return q
}

// Conjugated は共役クォータニオンを返します
func (q *Quaternion) Conjugated() *Quaternion {
	return &Quaternion{Number: quat.Conj(q.Number)}
}

// Log はクォータニオンの対数を返します
func (q *Quaternion) Log() *Quaternion {
	return &Quaternion{Number: quat.Log(q.Number)}
}

// Exp はクォータニオンの指数を返します
func (q *Quaternion) Exp() *Quaternion {
	return &Quaternion{Number: quat.Exp(q.Number)}
}

// ----- ベクトル演算 -----

// Length はクォータニオンの長さを返します
func (q *Quaternion) Length() float64 {
	return quat.Abs(q.Number)
}

// Dot は内積を返します
func (q *Quaternion) Dot(other *Quaternion) float64 {
	return q.X()*other.X() + q.Y()*other.Y() + q.Z()*other.Z() + q.W()*other.W()
}

// ----- 変換 -----

// ToRadians はオイラー角（ラジアン）に変換します
func (q *Quaternion) ToRadians() *Vec3 {
	x, y, z, w := q.X(), q.Y(), q.Z(), q.W()
	sx := -(2*y*z - 2*x*w)
	unlocked := math.Abs(sx) < 0.99999
	xPitch := math.Asin(math.Max(-1, math.Min(1, sx)))
	var yHead, zRoll float64
	if unlocked {
		yHead = math.Atan2(2*x*z+2*y*w, 2*w*w+2*z*z-1)
		zRoll = math.Atan2(2*x*y+2*z*w, 2*w*w+2*y*y-1)
	} else {
		yHead = math.Atan2(-(2*x*z - 2*y*w), 2*w*w+2*x*x-1)
		zRoll = 0
	}
	return NewVec3ByValues(xPitch, yHead, zRoll)
}

// ToDegrees はオイラー角（度）に変換します
func (q *Quaternion) ToDegrees() *Vec3 {
	rad := q.ToRadians()
	return NewVec3ByValues(RadToDeg(rad.X), RadToDeg(rad.Y), RadToDeg(rad.Z))
}

// ToAxisAngle は軸と角度に変換します
func (q *Quaternion) ToAxisAngle() (*Vec3, float64) {
	qn := q.Normalized()
	angle := 2.0 * math.Acos(math.Max(-1, math.Min(1, qn.W())))
	s := math.Sqrt(1 - qn.W()*qn.W())
	if s < 1e-9 {
		return VEC3_UNIT_X.Copy(), angle
	}
	return NewVec3ByValues(qn.X()/s, qn.Y()/s, qn.Z()/s), angle
}

// ToMat4 は4x4行列に変換します
func (q *Quaternion) ToMat4() *Mat4 {
	x, y, z, w := q.X(), q.Y(), q.Z(), q.W()
	xx, yy, zz := x*x, y*y, z*z
	xy, xz, yz := x*y, x*z, y*z
	wx, wy, wz := w*x, w*y, w*z

	return &Mat4{
		1 - 2*(yy+zz), 2 * (xy + wz), 2 * (xz - wy), 0,
		2 * (xy - wz), 1 - 2*(xx+zz), 2 * (yz + wx), 0,
		2 * (xz + wy), 2 * (yz - wx), 1 - 2*(xx+yy), 0,
		0, 0, 0, 1,
	}
}

// ToRadian は回転角度（ラジアン）を返します
func (q *Quaternion) ToRadian() float64 {
	return 2 * math.Acos(math.Min(1, math.Max(-1, q.Normalized().W())))
}

// ToDegree は回転角度（度）を返します
func (q *Quaternion) ToDegree() float64 {
	return RadToDeg(q.ToRadian())
}

// ----- 補間 -----

// Slerp は球面線形補間を行います
func (q *Quaternion) Slerp(other *Quaternion, t float64) *Quaternion {
	if t <= 0 {
		return q.Copy()
	}
	if t >= 1 {
		return other.Copy()
	}
	if q.NearEquals(other, 1e-8) {
		return q.Copy()
	}

	cosOmega := q.Dot(other)
	q2 := other.Copy()
	if cosOmega < 0 {
		cosOmega = -cosOmega
		q2.Negate()
	}

	var k1, k2 float64
	if cosOmega > 0.9999 {
		k1 = 1 - t
		k2 = t
	} else {
		sinOmega := math.Sqrt(1 - cosOmega*cosOmega)
		omega := math.Atan2(sinOmega, cosOmega)
		invSinOmega := 1 / sinOmega
		k1 = math.Sin((1-t)*omega) * invSinOmega
		k2 = math.Sin(t*omega) * invSinOmega
	}

	return NewQuaternionByValues(
		k1*q.X()+k2*q2.X(),
		k1*q.Y()+k2*q2.Y(),
		k1*q.Z()+k2*q2.Z(),
		k1*q.W()+k2*q2.W(),
	)
}

// Lerp は線形補間を行います
func (q *Quaternion) Lerp(other *Quaternion, t float64) *Quaternion {
	if t <= 0 {
		return q.Copy()
	}
	if t >= 1 {
		return other.Copy()
	}

	scale0 := 1 - t
	scale1 := t

	if q.Dot(other) < 0 {
		scale1 = -scale1
	}

	x := scale0*q.X() + scale1*other.X()
	y := scale0*q.Y() + scale1*other.Y()
	z := scale0*q.Z() + scale1*other.Z()
	w := scale0*q.W() + scale1*other.W()

	return NewQuaternionByValues(x, y, z, w).Normalize()
}

// ----- ベクトル回転 -----

// MulVec3 はベクトルをクォータニオンで回転します
func (q *Quaternion) MulVec3(v *Vec3) *Vec3 {
	qx, qy, qz, qw := q.X(), q.Y(), q.Z(), q.W()
	vx, vy, vz := v.X, v.Y, v.Z

	twoQx, twoQy, twoQz := 2*qx, 2*qy, 2*qz
	xx, xy, xz := qx*twoQx, qx*twoQy, qx*twoQz
	yy, yz, zz := qy*twoQy, qy*twoQz, qz*twoQz
	wx, wy, wz := qw*twoQx, qw*twoQy, qw*twoQz

	return NewVec3ByValues(
		vx*(1-(yy+zz))+vy*(xy-wz)+vz*(xz+wy),
		vx*(xy+wz)+vy*(1-(xx+zz))+vz*(yz-wx),
		vx*(xz-wy)+vy*(yz+wx)+vz*(1-(xx+yy)),
	)
}

// ----- ユーティリティ -----

// Copy はコピーを返します
func (q *Quaternion) Copy() *Quaternion {
	return &Quaternion{Number: q.Number}
}

// Vec3 はベクトル部分を返します
func (q *Quaternion) Vec3() *Vec3 {
	return NewVec3ByValues(q.X(), q.Y(), q.Z())
}

// Vec4 は4次元ベクトルとして返します
func (q *Quaternion) Vec4() *Vec4 {
	return NewVec4ByValues(q.X(), q.Y(), q.Z(), q.W())
}

// Vector はスライス形式で返します
func (q *Quaternion) Vector() []float64 {
	return []float64{q.X(), q.Y(), q.Z(), q.W()}
}

// IsIdent は単位クォータニオンかどうかを返します
func (q *Quaternion) IsIdent() bool {
	return q.NearEquals(QUATERNION_IDENTITY, 1e-6)
}

// NearEquals は他のクォータニオンとほぼ等しいかどうかを返します
func (q *Quaternion) NearEquals(other *Quaternion, epsilon float64) bool {
	return math.Abs(q.X()-other.X()) <= epsilon &&
		math.Abs(q.Y()-other.Y()) <= epsilon &&
		math.Abs(q.Z()-other.Z()) <= epsilon &&
		math.Abs(q.W()-other.W()) <= epsilon
}

// IsShortestRotation は最短回転かどうかを返します
func (q *Quaternion) IsShortestRotation(other *Quaternion) bool {
	return q.Dot(other) >= 0
}

// SetShortestRotation は最短回転になるよう調整します（破壊的）
func (q *Quaternion) SetShortestRotation(other *Quaternion) *Quaternion {
	if !q.IsShortestRotation(other) {
		q.Negate()
	}
	return q
}

// Shorten は最短回転に変換します（破壊的）
func (q *Quaternion) Shorten() *Quaternion {
	if q.W() < 0 {
		q.Negate()
	}
	return q
}
