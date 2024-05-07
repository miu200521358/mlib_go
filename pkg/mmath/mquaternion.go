package mmath

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type MQuaternion mgl64.Quat

func NewMQuaternion() *MQuaternion {
	return &MQuaternion{1., mgl64.Vec3{0, 0, 0}}
}

func NewMQuaternionByVec3(vec3 *MVec3) *MQuaternion {
	return NewMQuaternionByValues(vec3.GetX(), vec3.GetY(), vec3.GetZ(), 0)
}

// 指定された値でクォータニオンを作成します。
// ただし必ず最短距離クォータニオンにします
func NewMQuaternionByValuesShort(x, y, z, w float64) *MQuaternion {
	qq := &MQuaternion{w, mgl64.Vec3{x, y, z}}
	if !MQuaternionIdent.IsShortestRotation(qq) {
		qq.Negate()
	}
	return qq
}

// NewMQuaternionByValuesOriginal は、指定された値でクォータニオンを作成します。
// ただし、強制的に最短距離クォータニオンにはしません
func NewMQuaternionByValues(x, y, z, w float64) *MQuaternion {
	return &MQuaternion{w, mgl64.Vec3{x, y, z}}
}

var (
	// Zero holds a zero quaternion.
	MQuaternionZero = MQuaternion{}

	// Ident holds an ident quaternion.
	MQuaternionIdent = MQuaternion{1., mgl64.Vec3{0, 0, 0}}
)

// GetX returns the value of the X coordinate
func (v *MQuaternion) GetX() float64 {
	return v.V[0]
}

// SetX sets the value of the X coordinate
func (v *MQuaternion) SetX(x float64) {
	v.V[0] = x
}

// GetY returns the value of the Y coordinate
func (v *MQuaternion) GetY() float64 {
	return v.V[1]
}

// SetY sets the value of the Y coordinate
func (v *MQuaternion) SetY(y float64) {
	v.V[1] = y
}

// GetZ returns the value of the Z coordinate
func (v *MQuaternion) GetZ() float64 {
	return v.V[2]
}

// SetZ sets the value of the Z coordinate
func (v *MQuaternion) SetZ(z float64) {
	v.V[2] = z
}

// GetW returns the value of the W coordinate
func (v *MQuaternion) GetW() float64 {
	return v.W
}

// SetW sets the value of the W coordinate
func (v *MQuaternion) SetW(w float64) {
	v.W = w
}

func (v *MQuaternion) GetXYZ() *MVec3 {
	return &MVec3{v.GetX(), v.GetY(), v.GetZ()}
}

func (v *MQuaternion) SetXYZ(vec3 *MVec3) {
	v.SetX(vec3.GetX())
	v.SetY(vec3.GetY())
	v.SetZ(vec3.GetZ())
}

func (v *MQuaternion) AssignVec3(vec3 *MVec3) {
	v.SetX(vec3.GetX())
	v.SetY(vec3.GetY())
	v.SetZ(vec3.GetZ())
}

// String T の文字列表現を返します。
func (v *MQuaternion) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v.GetX(), v.GetY(), v.GetZ(), v.GetW())
}

// MMD MMD(MikuMikuDance)座標系に変換されたクォータニオンベクトルを返します
func (v *MQuaternion) MMD() *MQuaternion {
	return &MQuaternion{v.GetW(), mgl64.Vec3{-v.GetX(), -v.GetY(), v.GetZ()}}
}

// NewMQuaternionFromAxisAngles は、軸周りの回転を表す四元数を返します。
func NewMQuaternionFromAxisAngles(axis *MVec3, angle float64) *MQuaternion {
	angle *= 0.5
	sin := math.Sin(angle)
	q := NewMQuaternionByValues(axis[0]*sin, axis[1]*sin, axis[2]*sin, math.Cos(angle))
	return q.Normalize()
}

// NewMQuaternionFromXAxisAngleは、X軸周りの回転を表す四元数を返します。
func NewMQuaternionFromXAxisAngle(angle float64) *MQuaternion {
	angle *= 0.5
	q := NewMQuaternionByValues(math.Sin(angle), 0, 0, math.Cos(angle))
	return q.Normalize()
}

// NewMQuaternionFromYAxisAngleは、Y軸周りの回転を表す四元数を返します。
func NewMQuaternionFromYAxisAngle(angle float64) *MQuaternion {
	angle *= 0.5
	q := NewMQuaternionByValues(0, math.Sin(angle), 0, math.Cos(angle))
	return q.Normalize()
}

// NewMQuaternionFromZAxisAngleは、Z軸周りの回転を表す四元数を返します。
func NewMQuaternionFromZAxisAngle(angle float64) *MQuaternion {
	angle *= 0.5
	q := NewMQuaternionByValues(0, 0, math.Sin(angle), math.Cos(angle))
	return q.Normalize()
}

// NewMQuaternionFromRadiansは、オイラー角（ラジアン）回転を表す四元数を返します。
func NewMQuaternionFromRadians(xPitch, yHead, zRoll float64) *MQuaternion {
	qy := NewMQuaternionFromYAxisAngle(yHead)
	qx := NewMQuaternionFromXAxisAngle(xPitch)
	qz := NewMQuaternionFromZAxisAngle(zRoll)
	q := qy.Mul(qx)
	return q.Mul(qz).Normalize()
}

// 参考URL:
// https://qiita.com/aa_debdeb/items/abe90a9bd0b4809813da
// https://site.nicovideo.jp/ch/userblomaga_thanks/archive/ar805999

// ToRadiansは、クォータニオンを三軸のオイラー角（ラジアン）回転を返します。
func (v *MQuaternion) ToRadians() *MVec3 {
	sx := -(2*v.GetY()*v.GetZ() - 2*v.GetX()*v.GetW())
	unlocked := math.Abs(sx) < 0.99999
	xPitch := math.Asin(math.Max(-1, math.Min(1, sx)))
	var yHead, zRoll float64
	if unlocked {
		yHead = math.Atan2(2*v.GetX()*v.GetZ()+2*v.GetY()*v.GetW(), 2*v.GetW()*v.GetW()+2*v.GetZ()*v.GetZ()-1)
		zRoll = math.Atan2(2*v.GetX()*v.GetY()+2*v.GetZ()*v.GetW(), 2*v.GetW()*v.GetW()+2*v.GetY()*v.GetY()-1)
	} else {
		yHead = math.Atan2(-(2*v.GetX()*v.GetZ() - 2*v.GetY()*v.GetW()), 2*v.GetW()*v.GetW()+2*v.GetX()*v.GetX()-1)
		zRoll = 0
	}

	return &MVec3{xPitch, yHead, zRoll}
}

const (
	GIMBAL2_RAD = math.Pi * 88.0 * 2 / 180
	ONE_RAD     = math.Pi
)

// ToRadiansWithGimbalは、クォータニオンを三軸のオイラー角（ラジアン）回転を返します。
// ジンバルロックが発生しているか否かのフラグも返します
func (v *MQuaternion) ToRadiansWithGimbal(axisIndex int) (*MVec3, bool) {
	r := v.ToRadians()

	var other1Rad, other2Rad float64
	if axisIndex == 0 {
		other1Rad = math.Abs(r.Vector()[1])
		other2Rad = math.Abs(r.Vector()[2])
	} else if axisIndex == 1 {
		other1Rad = math.Abs(r.Vector()[0])
		other2Rad = math.Abs(r.Vector()[2])
	} else {
		other1Rad = math.Abs(r.Vector()[0])
		other2Rad = math.Abs(r.Vector()[1])
	}

	// ジンバルロックを判定する
	if other1Rad >= GIMBAL2_RAD && other2Rad >= GIMBAL2_RAD {
		return r, true
	}

	return r, false
}

// NewMQuaternionFromDegreesは、オイラー角（度）回転を表す四元数を返します。
func NewMQuaternionFromDegrees(xPitch, yHead, zRoll float64) *MQuaternion {
	xPitchRadian := math.Pi * xPitch / 180.0
	yHeadRadian := math.Pi * yHead / 180.0
	zRollRadian := math.Pi * zRoll / 180.0
	return NewMQuaternionFromRadians(xPitchRadian, yHeadRadian, zRollRadian)
}

// ToDegreesは、クォータニオンのオイラー角（度）回転を返します。
func (quat *MQuaternion) ToDegrees() *MVec3 {
	vec := quat.ToRadians()
	return &MVec3{
		180.0 * vec.GetX() / math.Pi,
		180.0 * vec.GetY() / math.Pi,
		180.0 * vec.GetZ() / math.Pi,
	}
}

// NewMQuaternionFromVec4はvec4.Tをクォータニオンに変換する
func NewMQuaternionFromVec4(v *MVec4) *MQuaternion {
	return NewMQuaternionByValues(v.GetX(), v.GetY(), v.GetZ(), v.GetW())
}

// Vec4は四元数をvec4.Tに変換する
func (quat *MQuaternion) Vec4() *MVec4 {
	return &MVec4{quat.GetX(), quat.GetY(), quat.GetZ(), quat.GetW()}
}

// Vec3は、クォータニオンのベクトル部分を返します。
func (quat *MQuaternion) Vec3() *MVec3 {
	vec3 := MVec3{quat.GetX(), quat.GetY(), quat.GetZ()}
	return &vec3
}

// AxisAngleは、正規化されたクォータニオンから、軸と回転角度の形で回転を取り出す。
func (quat *MQuaternion) AxisAngle() (axis MVec3, angle float64) {
	cos := quat.W
	sin := math.Sqrt(1 - cos*cos)
	angle = math.Acos(cos) * 2

	var ooSin float64
	if math.Abs(sin) < 0.0005 {
		ooSin = 1
	} else {
		ooSin = 1 / sin
	}
	axis[0] = quat.V[0] * ooSin
	axis[1] = quat.V[1] * ooSin
	axis[2] = quat.V[2] * ooSin

	return axis, angle
}

// Mul は、クォータニオンの積を返します。
func (q1 *MQuaternion) MulShort(q2 *MQuaternion) *MQuaternion {
	mat1 := q1.ToMat4()
	mat2 := q2.ToMat4()
	mat1.Mul(mat2)
	qq := mat1.Quaternion()

	return NewMQuaternionByValues(qq.V[0], qq.V[1], qq.V[2], qq.W)
}

func (q1 *MQuaternion) MuledShort(q2 *MQuaternion) *MQuaternion {
	copied := q1.Copy()
	copied.Mul(q2)
	return copied
}

// Mul は、クォータニオンの積を返します。
func (q1 *MQuaternion) Mul(q2 *MQuaternion) *MQuaternion {
	mat1 := q1.ToMat4()
	mat2 := q2.ToMat4()
	mat1.Mul(mat2)
	qq := mat1.Quaternion()

	q1.SetX(qq.GetX())
	q1.SetY(qq.GetY())
	q1.SetZ(qq.GetZ())
	q1.SetW(qq.GetW())
	return q1
}

func (q1 *MQuaternion) Muled(q2 *MQuaternion) *MQuaternion {
	copied := q1.Copy()
	copied.Mul(q2)
	return copied
}

// Norm はクォータニオンのノルム値を返します。
func (quat *MQuaternion) Norm() float64 {
	return quat.V[0]*quat.V[0] + quat.V[1]*quat.V[1] + quat.V[2]*quat.V[2] + quat.W*quat.W
}

// Len gives the Length of the quaternion, also known as its Norm. This is the
// same thing as the Len of a Vec4.
func (q1 MQuaternion) Length() float64 {
	return math.Sqrt(float64(q1.W*q1.W + q1.V[0]*q1.V[0] + q1.V[1]*q1.V[1] + q1.V[2]*q1.V[2]))
}

// Normalizeは、単位四位数に正規化する。
func (quat *MQuaternion) Normalize() *MQuaternion {
	norm := quat.Norm()
	if norm != 1 && norm != 0 {
		ool := 1 / math.Sqrt(norm)
		quat.V[0] *= ool
		quat.V[1] *= ool
		quat.V[2] *= ool
		quat.W *= ool
	}
	return quat
}

// Normalizedは、単位を4進数に正規化したコピーを返す。
func (quat *MQuaternion) Normalized() *MQuaternion {
	norm := quat.Norm()
	if norm != 1 && norm != 0 {
		ool := 1 / math.Sqrt(norm)
		q := NewMQuaternionByValues(
			quat.V[0]*ool,
			quat.V[1]*ool,
			quat.V[2]*ool,
			quat.W*ool,
		)
		return q
	} else {
		return quat
	}
}

// Negate negates the quaternion.
func (quat *MQuaternion) Negate() *MQuaternion {
	quat.V[0] = -quat.V[0]
	quat.V[1] = -quat.V[1]
	quat.V[2] = -quat.V[2]
	quat.W = -quat.W
	return quat
}

// Negated returns a negated quaternion.
func (quat *MQuaternion) Negated() *MQuaternion {
	return NewMQuaternionByValues(-quat.V[0], -quat.V[1], -quat.V[2], -quat.W)
}

// Invert inverts the quaternion.
func (quat *MQuaternion) Invert() *MQuaternion {
	quat.V[0] = -quat.V[0]
	quat.V[1] = -quat.V[1]
	quat.V[2] = -quat.V[2]
	return quat
}

// Inverted returns a inverted quaternion.
func (quat *MQuaternion) Inverted() *MQuaternion {
	return NewMQuaternionByValues(-quat.V[0], -quat.V[1], -quat.V[2], quat.W)
}

// SetShortestRotation は、クォータニオンが quat から other の方向への最短回転を表していない場合、そのクォータニオンを否定します。
// (quatの向きからotherの向きへの回転には2つの方向があります)
func (quat *MQuaternion) SetShortestRotation(other *MQuaternion) *MQuaternion {
	if !quat.IsShortestRotation(other) {
		quat.Negate()
	}
	return quat
}

// IsShortestRotation は、a から b への回転が可能な限り最短の回転かどうかを返す。
// (quatの向きから他の向きへの回転には2つの方向がある)
func (quat *MQuaternion) IsShortestRotation(other *MQuaternion) bool {
	return quat.Dot(other) >= 0
}

// IsUnitQuat は、クォータニオンが単位クォータニオンの許容範囲内にあるかどうかを返します。
func (quat *MQuaternion) IsUnitQuat(tolerance float64) bool {
	norm := quat.Norm()
	return norm >= (1.0-tolerance) && norm <= (1.0+tolerance)
}

// 最短回転に変換します
func (quat *MQuaternion) Shorten() *MQuaternion {
	if quat.W < 0 {
		quat.Negate()
	}
	return quat
}

// RotateVec3 は、四元数によって表される回転によって v を回転させます。
// https://gamedev.stackexchange.com/questions/28395/rotating-vector3-by-a-quaternion
func (quat *MQuaternion) RotateVec3(v *MVec3) {
	u := MVec3{quat.V[0], quat.V[1], quat.V[2]}
	s := quat.W
	vt1 := u.MuledScalar(2 * u.Dot(v))
	vt2 := v.MuledScalar(s*s - u.Dot(&u))
	vt3 := u.Cross(v)
	vt3 = vt3.MulScalar(2 * s)
	v[0] = vt1[0] + vt2[0] + vt3[0]
	v[1] = vt1[1] + vt2[1] + vt3[1]
	v[2] = vt1[2] + vt2[2] + vt3[2]
}

// RotatedVec3 は v の回転コピーを返す。
// https://gamedev.stackexchange.com/questions/28395/rotating-vector3-by-a-quaternion
func (quat *MQuaternion) RotatedVec3(v *MVec3) *MVec3 {
	u := MVec3{quat.V[0], quat.V[1], quat.V[2]}
	s := quat.W
	vt1 := u.MuledScalar(2 * u.Dot(v))
	vt2 := v.MuledScalar(s*s - u.Dot(&u))
	vt3 := u.Cross(v)
	vt3 = vt3.MulScalar(2 * s)
	return &MVec3{vt1[0] + vt2[0] + vt3[0], vt1[1] + vt2[1] + vt3[1], vt1[2] + vt2[2] + vt3[2]}
}

// Dot は2つのクォータニオンの内積を返す。
func (quat *MQuaternion) Dot(other *MQuaternion) float64 {
	return quat.V[0]*other.V[0] + quat.V[1]*other.V[1] + quat.V[2]*other.V[2] + quat.W*other.W
}

// MulScalar
func (quat *MQuaternion) MulScalar(factor float64) *MQuaternion {
	if factor == 0.0 {
		return NewMQuaternion()
	}

	// qq := NewMQuaternion().Lerp(quat, math.Abs(factor))
	// if factor < 0 {
	// 	qq.Invert()
	// }

	// return qq.Normalize()

	axis, angle := quat.ToAxisAngle()

	// factor をかけて角度を制限
	angle = math.Mod(angle*factor, math.Pi*2)

	return NewMQuaternionFromAxisAngles(axis, angle)
}

// ToAxisAngle converts the quaternion to an axis and an angle.
func (quat *MQuaternion) ToAxisAngle() (*MVec3, float64) {
	// Normalize the quaternion
	quat.Normalize()

	// Calculate the angle
	angle := 2 * math.Acos(quat.W)

	// Calculate the components of the axis
	s := math.Sqrt(1 - quat.W*quat.W)
	if s < 1e-9 {
		s = 1
	}
	axis := NewMVec3()
	axis.SetX(quat.V[0] / s)
	axis.SetY(quat.V[1] / s)
	axis.SetZ(quat.V[2] / s)

	return axis, angle
}

// Slerp は t (0,1) における a と b の間の球面線形補間クォータニオンを返す。
// See http://en.wikipedia.org/wiki/Slerp
func (a *MQuaternion) Slerp(b *MQuaternion, t float64) *MQuaternion {
	q := mgl64.QuatSlerp(mgl64.Quat(*a), mgl64.Quat(*b), t)
	return (*MQuaternion)(&q)
}

func (q *MQuaternion) Lerp(other *MQuaternion, t float64) *MQuaternion {
	qq := mgl64.QuatLerp(mgl64.Quat(*q), mgl64.Quat(*other), t)
	return (*MQuaternion)(&qq)
}

// Vec3Diff関数は、2つのベクトル間の回転四元数を返します。
func (a *MVec3) Vec3Diff(b *MVec3) *MQuaternion {
	cr := a.Cross(b)
	sr := math.Sqrt(2 * (1 + a.Dot(b)))
	oosr := 1 / sr

	q := NewMQuaternionByValues(cr[0]*oosr, cr[1]*oosr, cr[2]*oosr, sr*0.5)
	return q.Normalize()
}

// ToDegree は、クォータニオンを度に変換します。
func (quat *MQuaternion) ToDegree() float64 {
	w := quat.Normalize().GetW()
	radian := 2 * math.Acos(math.Min(1, math.Max(-1, w)))
	angle := radian * (180 / math.Pi)
	return angle
}

// ToRadian は、クォータニオンをラジアンに変換します。
func (quat *MQuaternion) ToRadian() float64 {
	w := quat.Normalize().GetW()
	radian := 2 * math.Acos(math.Min(1, math.Max(-1, w)))
	return radian
}

// ToSignedDegree 符号付き角度に変換
func (quat *MQuaternion) ToSignedDegree() float64 {
	// スカラー部分から基本的な角度を計算
	basicAngle := quat.ToDegree()

	// ベクトルの長さを使って、角度の正負を決定
	if quat.Vec3().Length() > 0 {
		// ベクトルの向きに基づいて角度を調整
		if quat.GetW() >= 0 {
			return basicAngle
		} else {
			return -basicAngle
		}
	}

	// ベクトル部分がない場合は基本角度をそのまま使用
	return basicAngle
}

// ToSignedRadian 符号付きラジアンに変換
func (quat *MQuaternion) ToSignedRadian(axisIndex int) float64 {
	// スカラー部分から基本的な角度を計算
	basicAngle := quat.ToRadian()

	// ベクトルの長さを使って、角度の正負を決定
	if quat.Vec3().Length() > 0 {
		if (axisIndex == 0 && quat.GetX() < 0) ||
			(axisIndex == 1 && quat.GetY() < 0) ||
			(axisIndex == 2 && quat.GetZ() < 0) {
			return -basicAngle
		}
	}

	// ベクトル部分がない場合は基本角度をそのまま使用
	return basicAngle
}

// ToTheta 自分ともうひとつの値vとのtheta（変位量）を返す
func (quat *MQuaternion) ToTheta(v *MQuaternion) float64 {
	return math.Acos(math.Min(1, math.Max(-1, quat.Normalize().Dot(v.Normalize()))))
}

// 軸と角度からクォータニオンに変換する
func NewMQuaternionFromDirection(direction *MVec3, up *MVec3) *MQuaternion {
	if direction.Length() == 0 {
		return NewMQuaternion()
	}

	zAxis := direction.Normalize()
	xAxis := up.Cross(zAxis).Normalize()

	if xAxis.LengthSqr() == 0 {
		// collinear or invalid up vector derive shortest arc to new direction
		return NewMQuaternionRotate(&MVec3{0.0, 0.0, 1.0}, zAxis)
	}

	yAxis := zAxis.Cross(xAxis)

	return NewMQuaternionFromAxes(xAxis, yAxis, zAxis).Normalize()
}

// NewMQuaternionRotate fromベクトルからtoベクトルまでの回転量
func NewMQuaternionRotate(fromV, toV *MVec3) *MQuaternion {
	v0 := fromV.Normalize()
	v1 := toV.Normalize()
	d := v0.Dot(v1) + 1.0

	// if dest vector is close to the inverse of source vector, ANY axis of rotation is valid
	if math.Abs(d) < 1e-6 {
		axis := MVec3UnitX.Cross(v0)
		if math.Abs(axis.LengthSqr()) < 1e-6 {
			axis = MVec3UnitY.Cross(v0)
		}
		axis.Normalize()
		// same as MQuaternion.fromAxisAndAngle(axis, 180.0)
		return NewMQuaternionByValues(axis.GetX(), axis.GetY(), axis.GetZ(), 0.0)
	}

	d = math.Sqrt(2.0 * d)
	axis := v0.Cross(v1).DivScalar(d)
	return NewMQuaternionByValues(axis.GetX(), axis.GetY(), axis.GetZ(), d*0.5)
}

// NewMQuaternionFromAxes
func NewMQuaternionFromAxes(xAxis, yAxis, zAxis *MVec3) *MQuaternion {
	x := MVec3{xAxis.GetX(), yAxis.GetX(), zAxis.GetX()}
	y := MVec3{xAxis.GetY(), yAxis.GetY(), zAxis.GetY()}
	z := MVec3{xAxis.GetZ(), yAxis.GetZ(), zAxis.GetZ()}
	mat := NewMMat3().AssignCoordinateSystem(&x, &y, &z)
	qq := mat.Quaternion()
	return qq
}

// SeparateByAxis separates the quaternion into four quaternions based on the global axis.
func (quat *MQuaternion) SeparateByAxis(globalAxis *MVec3) (*MQuaternion, *MQuaternion, *MQuaternion, *MQuaternion) {
	localZAxis := MVec3UnitZ
	globalXAxis := globalAxis.Normalize()
	globalYAxis := localZAxis.Cross(globalXAxis)
	globalZAxis := globalXAxis.Cross(globalYAxis)

	if globalYAxis.Length() == 0 {
		localYAxis := MVec3UnitY
		globalZAxis := localYAxis.Cross(globalXAxis)
		globalYAxis = globalXAxis.Cross(globalZAxis)
	}

	// X成分を抽出する ------------

	// グローバル軸方向に伸ばす
	globalXVec := quat.RotatedVec3(globalXAxis)
	// YZの回転量（自身のねじれを無視する）
	yzQQ := NewMQuaternionRotate(globalXAxis, globalXVec.Normalize())
	// 元々の回転量 から YZ回転 を除去して、除去されたX成分を求める
	invYzQQ := yzQQ.Inverted()
	xQQ := quat.Mul(invYzQQ)

	// Y成分を抽出する ------------
	// グローバル軸方向に伸ばす
	globalYVec := quat.RotatedVec3(globalYAxis)
	// XZの回転量（自身のねじれを無視する）
	xzQQ := NewMQuaternionRotate(globalYAxis, globalYVec.Normalize())
	// 元々の回転量 から XZ回転 を除去して、除去されたY成分を求める
	invXzQQ := xzQQ.Inverted()
	yQQ := quat.Mul(invXzQQ)

	// Z成分を抽出する ------------
	// グローバル軸方向に伸ばす
	globalZVec := quat.RotatedVec3(globalZAxis)
	// XYの回転量（自身のねじれを無視する）
	xyQQ := NewMQuaternionRotate(globalZAxis, globalZVec.Normalize())
	// 元々の回転量 から XY回転 を除去して、除去されたZ成分を求める
	invXyQQ := xyQQ.Inverted()
	zQQ := quat.Mul(invXyQQ)

	return xQQ, yQQ, zQQ, yzQQ
}

// Copy
func (qq *MQuaternion) Copy() *MQuaternion {
	return NewMQuaternionByValues(qq.GetX(), qq.GetY(), qq.GetZ(), qq.GetW())
}

// Vector
func (v *MQuaternion) Vector() []float64 {
	return []float64{v.GetX(), v.GetY(), v.GetZ(), v.GetW()}
}

func (v *MQuaternion) ToMat4() *MMat4 {
	mat := NewMMat4()
	mat.AssignQuaternion(v)
	return mat
}

// ToFixedAxisRotation 軸制限されたクォータニオンの回転
// fixedAxis: 軸制限を表す3次元ベクトル
func (quat *MQuaternion) ToFixedAxisRotation(fixedAxis *MVec3) *MQuaternion {
	normalizedFixedAxis := fixedAxis.Normalized()
	quatAxis := quat.GetXYZ().Normalized()
	rad := quat.ToRadian()
	if normalizedFixedAxis.Dot(quatAxis) < 0 {
		rad *= -1
	}
	return NewMQuaternionFromAxisAngles(normalizedFixedAxis, rad)
}

// PracticallyEquals
func (quat *MQuaternion) PracticallyEquals(other *MQuaternion, epsilon float64) bool {
	return (math.Abs(quat.V[0]-other.V[0]) <= epsilon) &&
		(math.Abs(quat.V[1]-other.V[1]) <= epsilon) &&
		(math.Abs(quat.V[2]-other.V[2]) <= epsilon) &&
		(math.Abs(quat.W-other.W) <= epsilon)
}

// MulVec3 multiplies v (converted to a vec4 as (v_1, v_2, v_3, 1))
// with mat and divides the result by w. Returns a new vec3.
func (quat *MQuaternion) MulVec3(v *MVec3) *MVec3 {
	return quat.ToMat4().MulVec3(v)
}

// ラジアン角度をオイラー角度に変換
func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

// オイラー角度をラジアン角度に変換
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// VectorToDegree は、与えられた2つのベクトルから角度に変換します。
func VectorToDegree(a *MVec3, b *MVec3) float64 {
	return RadToDeg(VectorToRadian(a, b))
}

// VectorToRadian は、与えられた2つのベクトルからラジアン角度に変換します。
func VectorToRadian(a *MVec3, b *MVec3) float64 {
	p := a.Dot(b)
	normA := a.Length()
	normB := b.Length()

	// 角度を計算
	cosAngle := p / (normA * normB)
	rad := math.Acos(math.Min(1, math.Max(-1, cosAngle)))

	return rad
}

// FindSlerpT 始点Q1、終点Q2、中間点Qtが与えられたとき、Slerp(Q1, Q2, t) ≈ Qtとなるtを見つけます。
func FindSlerpT(Q1, Q2, Qt *MQuaternion) float64 {
	t0 := 0.5
	eps := 1e-15
	err := 1e-20

	return findSlerpT(Q1, Q2, Qt, t0, eps, err)
}

// findSlerpT uses Newton's method to approximate the t value for which Slerp(Q1, Q2, t) approximates Qt.
func findSlerpT(Q1, Q2, Qt *MQuaternion, t0, eps, err float64) float64 {
	maxIterations := 20
	for i := 0; i < maxIterations; i++ {
		f := quatError(Q1, Q2, Qt, t0)
		df := quatErrorDerivative(Q1, Q2, Qt, t0, eps)

		// Update t0 using Newton's method
		t1 := t0 - f/df

		if math.Abs(t1-t0) < err {
			return t1
		}

		t0 = t1
	}
	return t0
}

// quatError calculates the error between Slerp(Q1, Q2, t) and the target quaternion Qt.
func quatError(Q1, Q2, Qt *MQuaternion, t float64) float64 {
	t2 := Q1.Slerp(Q2, t)
	dot := math.Abs(t2.Dot(Qt))
	return 1 - dot
}

// quatErrorDerivative calculates the derivative of the error function with respect to t using central difference.
func quatErrorDerivative(Q1, Q2, Qt *MQuaternion, t, eps float64) float64 {
	fPlus := quatError(Q1, Q2, Qt, t+eps)
	fMinus := quatError(Q1, Q2, Qt, t-eps)
	return (fPlus - fMinus) / (2 * eps)
}
