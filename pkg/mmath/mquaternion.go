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

func NewMQuaternionByVec3(vec3 MVec3) *MQuaternion {
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

// String T の文字列表現を返します。
func (v *MQuaternion) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v.GetX(), v.GetY(), v.GetZ(), v.GetW())
}

// MMD MMD(MikuMikuDance)座標系に変換されたクォータニオンベクトルを返します
func (v *MQuaternion) MMD() *MQuaternion {
	return &MQuaternion{v.GetW(), mgl64.Vec3{v.GetX(), v.GetY(), -v.GetZ()}}
}

// NewMQuaternionFromAxisAngles は、軸周りの回転を表す四元数を返します。
func NewMQuaternionFromAxisAngles(axis *MVec3, angle float64) *MQuaternion {
	q := MQuaternion(mgl64.QuatRotate(angle, mgl64.Vec3(*axis)).Normalize())
	return &q
}

// NewMQuaternionFromRadiansは、オイラー角（ラジアン）回転を表す四元数を返します。
func NewMQuaternionFromRadians(xPitch, yHead, zRoll float64) *MQuaternion {
	q := mgl64.AnglesToQuat(xPitch, yHead, zRoll, mgl64.XYZ)
	return &MQuaternion{q.W, q.V}
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
	GIMBAL1_RAD    = 88.0 / 180.0 * math.Pi
	GIMBAL_MIN_RAD = GIMBAL1_RAD * 0.07
	GIMBAL0_5_RAD  = GIMBAL1_RAD * 0.5
	GIMBAL1_5_RAD  = GIMBAL1_RAD * 1.5
	GIMBAL2_RAD    = GIMBAL1_RAD * 2
	ONE_RAD        = math.Pi
)

// ToRadiansWithGimbalは、クォータニオンを三軸のオイラー角（ラジアン）回転を返します。
// ジンバルロックが発生しているか否かのフラグも返します
func (v *MQuaternion) ToRadiansWithGimbal(axisIndex int) (*MVec3, bool) {
	r := v.ToRadians()

	var other1Rad, other2Rad float64
	if axisIndex == 0 {
		other1Rad = math.Abs(r.GetY())
		other2Rad = math.Abs(r.GetZ())
	} else if axisIndex == 1 {
		other1Rad = math.Abs(r.GetX())
		other2Rad = math.Abs(r.GetZ())
	} else {
		other1Rad = math.Abs(r.GetX())
		other2Rad = math.Abs(r.GetY())
	}

	// ジンバルロックを判定する
	if other1Rad >= GIMBAL2_RAD && other2Rad >= GIMBAL2_RAD {
		return r, true
	}

	return r, false
}

// NewMQuaternionFromDegreesは、オイラー角（度）回転を表す四元数を返します。
func NewMQuaternionFromDegrees(xPitch, yHead, zRoll float64) *MQuaternion {
	xPitchRadian := DegToRad(xPitch)
	yHeadRadian := DegToRad(yHead)
	zRollRadian := DegToRad(zRoll)
	return NewMQuaternionFromRadians(xPitchRadian, yHeadRadian, zRollRadian)
}

// Utility functions to convert between degrees and radians
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

// ToDegreesは、クォータニオンのオイラー角（度）回転を返します。
func (quat *MQuaternion) ToDegrees() *MVec3 {
	vec := quat.ToRadians()
	return &MVec3{
		RadToDeg(vec.GetX()),
		RadToDeg(vec.GetY()),
		RadToDeg(vec.GetZ()),
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
	cos := quat.GetW()
	sin := math.Sqrt(1 - cos*cos)
	angle = math.Acos(cos) * 2

	var ooSin float64
	if math.Abs(sin) < 0.0005 {
		ooSin = 1
	} else {
		ooSin = 1 / sin
	}
	axis[0] = quat.GetX() * ooSin
	axis[1] = quat.GetY() * ooSin
	axis[2] = quat.GetZ() * ooSin

	return axis, angle
}

// Mulは、クォータニオンの積を返します。
func (q1 *MQuaternion) MulShort(q2 *MQuaternion) *MQuaternion {
	mat1 := q1.ToMat4()
	mat2 := q2.ToMat4()
	mat1.Mul(mat2)
	qq := mat1.Quaternion()

	return NewMQuaternionByValues(qq.GetX(), qq.GetY(), qq.GetZ(), qq.GetW())
}

func (q1 *MQuaternion) MuledShort(q2 *MQuaternion) *MQuaternion {
	copied := q1.Copy()
	copied.Mul(q2)
	return copied
}

// Mulは、クォータニオンの積を返します。
func (q1 *MQuaternion) Mul(q2 *MQuaternion) *MQuaternion {
	*q1 = MQuaternion(mgl64.Quat(*q1).Mul(mgl64.Quat(*q2)))
	return q1
}

func (q1 *MQuaternion) Muled(q2 *MQuaternion) *MQuaternion {
	copied := q1.Copy()
	copied.Mul(q2)
	return copied
}

// Normはクォータニオンのノルム値を返します。
func (quat *MQuaternion) Norm() float64 {
	return mgl64.Quat(*quat).Norm()
}

// Lengthはクォータニオンの長さ（ノルム）を返します。
func (quat *MQuaternion) Length() float64 {
	return mgl64.Quat(*quat).Len()
}

// Normalizeは、単位四位数に正規化する。
func (quat *MQuaternion) Normalize() *MQuaternion {
	q := mgl64.Quat(*quat).Normalize()
	return &MQuaternion{q.W, q.V}
}

// Normalizedは、単位を4進数に正規化したコピーを返す。
func (quat *MQuaternion) Normalized() *MQuaternion {
	copied := quat.Copy()
	return copied.Normalize()
}

// Negateはクォータニオンを反転する。
func (quat *MQuaternion) Negate() *MQuaternion {
	quat.SetX(-quat.GetX())
	quat.SetY(-quat.GetY())
	quat.SetZ(-quat.GetZ())
	quat.SetW(-quat.GetW())
	return quat
}

// Negatedは反転したクォータニオンを返します。
func (quat *MQuaternion) Negated() *MQuaternion {
	return NewMQuaternionByValues(-quat.GetX(), -quat.GetY(), -quat.GetZ(), -quat.GetW())
}

// Invertは、クォータニオンを反転させます。
func (quat *MQuaternion) Invert() *MQuaternion {
	quat.SetX(-quat.GetX())
	quat.SetY(-quat.GetY())
	quat.SetZ(-quat.GetZ())
	return quat
}

// Invertedは反転したクォータニオンを返します。
func (quat *MQuaternion) Inverted() *MQuaternion {
	return NewMQuaternionByValues(-quat.GetX(), -quat.GetY(), -quat.GetZ(), quat.GetW())
}

// SetShortestRotationは、クォータニオンが quat から other の方向への最短回転を表していない場合、そのクォータニオンを否定します。
// (quatの向きからotherの向きへの回転には2つの方向があります)
func (quat *MQuaternion) SetShortestRotation(other *MQuaternion) *MQuaternion {
	if !quat.IsShortestRotation(other) {
		quat.Negate()
	}
	return quat
}

// IsShortestRotationは、a から b への回転が可能な限り最短の回転かどうかを返します。
// (quatの向きから他の向きへの回転には2つの方向があります)
func (quat *MQuaternion) IsShortestRotation(other *MQuaternion) bool {
	return quat.Dot(other) >= 0
}

// IsUnitQuatは、クォータニオンが単位クォータニオンの許容範囲内にあるかどうかを返します。
func (quat *MQuaternion) IsUnitQuat(tolerance float64) bool {
	norm := quat.Norm()
	return norm >= (1.0-tolerance) && norm <= (1.0+tolerance)
}

// Shortenは、最短回転に変換します。
func (quat *MQuaternion) Shorten() *MQuaternion {
	if quat.GetW() < 0 {
		quat.Negate()
	}
	return quat
}

// RotateVec3は、四元数によって表される回転によって v を回転させます。
func (quat *MQuaternion) RotateVec3(v *MVec3) *MVec3 {
	r := mgl64.Quat(*quat).Rotate(mgl64.Vec3(*v))
	return &MVec3{r[0], r[1], r[2]}
}

// RotatedVec3は v の回転コピーを返します。
func (quat *MQuaternion) RotatedVec3(v *MVec3) *MVec3 {
	return quat.Copy().RotateVec3(v)
}

// Dotは2つのクォータニオンの内積を返します。
func (quat *MQuaternion) Dot(other *MQuaternion) float64 {
	return mgl64.Quat(*quat).Dot(mgl64.Quat(*other))
}

// MulScalarはクォータニオンにスカラーを掛け算します。
func (quat *MQuaternion) MulScalar(factor float64) *MQuaternion {
	if factor == 0.0 {
		return NewMQuaternion()
	}

	axis, angle := quat.AxisAngle()

	// factor をかけて角度を制限
	angle = math.Mod(angle*factor, math.Pi*2)

	return NewMQuaternionFromAxisAngles(&axis, angle)
}

// ToAxisAngleは、クォータニオンを軸と角度に変換します。
func (quat *MQuaternion) ToAxisAngle() (*MVec3, float64) {
	// クォータニオンを正規化
	quat.Normalize()

	// 角度を計算
	angle := 2 * math.Acos(quat.GetW())

	// 軸の成分を計算
	s := math.Sqrt(1 - quat.GetW()*quat.GetW())
	if s < 1e-9 {
		s = 1
	}
	axis := NewMVec3()
	axis.SetX(quat.GetX() / s)
	axis.SetY(quat.GetY() / s)
	axis.SetZ(quat.GetZ() / s)

	return axis, angle
}

// Slerpはt (0,1)におけるaとbの間の球面線形補間クォータニオンを返します。
// See http://en.wikipedia.org/wiki/Slerp
func (a *MQuaternion) Slerp(b *MQuaternion, t float64) *MQuaternion {
	q := mgl64.QuatSlerp(mgl64.Quat(*a), mgl64.Quat(*b), t)
	return (*MQuaternion)(&q)
}

func (q *MQuaternion) Lerp(other *MQuaternion, t float64) *MQuaternion {
	qq := mgl64.QuatLerp(mgl64.Quat(*q), mgl64.Quat(*other), t)
	return (*MQuaternion)(&qq)
}

// Vec3Diffは、2つのベクトル間の回転四元数を返します。
func (a *MVec3) Vec3Diff(b *MVec3) *MQuaternion {
	cr := a.Cross(b)
	sr := math.Sqrt(2 * (1 + a.Dot(b)))
	oosr := 1 / sr

	q := NewMQuaternionByValues(cr[0]*oosr, cr[1]*oosr, cr[2]*oosr, sr*0.5)
	return q.Normalize()
}

// ToDegreeは、クォータニオンを度に変換します。
func (quat *MQuaternion) ToDegree() float64 {
	return RadToDeg(quat.ToRadian())
}

// ToRadianは、クォータニオンをラジアンに変換します。
func (quat *MQuaternion) ToRadian() float64 {
	return 2 * math.Acos(math.Min(1, math.Max(-1, quat.Normalize().GetW())))
}

// ToSignedDegreeは、符号付き角度に変換します。
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

// ToSignedRadianは、符号付きラジアンに変換します。
func (quat *MQuaternion) ToSignedRadian(axisIndex int) float64 {
	// スカラー部分から基本的な角度を計算
	basicAngle := quat.ToRadian()

	// ベクトルの長さを使って、角度の正負を決定
	if quat.Vec3().Length() > 0 {
		if (axisIndex == 0 && quat.GetX() > 0) ||
			(axisIndex == 1 && quat.GetY() < 0) ||
			(axisIndex == 2 && quat.GetZ() < 0) {
			return -basicAngle
		}
	}

	// ベクトル部分がない場合は基本角度をそのまま使用
	return basicAngle
}

// ToThetaは、自分ともうひとつの値vとのtheta（変位量）を返します。
func (quat *MQuaternion) ToTheta(v *MQuaternion) float64 {
	return math.Acos(math.Min(1, math.Max(-1, quat.Normalize().Dot(v.Normalize()))))
}

// NewMQuaternionFromDirectionは、軸と角度からクォータニオンに変換します。
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

	result := NewMQuaternionFromAxes(xAxis, yAxis, zAxis)
	return result.Normalize()
}

// NewMQuaternionRotateはfromベクトルからtoベクトルまでの回転量を計算します。
func NewMQuaternionRotate(fromV, toV *MVec3) *MQuaternion {
	v0 := fromV.Normalize()
	v1 := toV.Normalize()
	d := v0.Dot(v1) + 1.0

	// dest vectorがsource vectorの逆方向に近い場合、任意の回転軸が有効です。
	if math.Abs(d) < 1e-6 {
		axis := MVec3UnitX.Cross(v0)
		if math.Abs(axis.LengthSqr()) < 1e-6 {
			axis = MVec3UnitY.Cross(v0)
		}
		axis.Normalize()
		// 同じくMQuaternion.fromAxisAndAngle(axis, 180.0)
		return NewMQuaternionByValues(axis.GetX(), axis.GetY(), axis.GetZ(), 0.0)
	}

	d = math.Sqrt(2.0 * d)
	axis := v0.Cross(v1).DivScalar(d)
	return NewMQuaternionByValues(axis.GetX(), axis.GetY(), axis.GetZ(), d*0.5)
}

// NewMQuaternionFromAxesは、3つの軸ベクトルからクォータニオンを作成します。
func NewMQuaternionFromAxes(xAxis, yAxis, zAxis *MVec3) *MQuaternion {
	mat := NewMMat4ByValues(
		xAxis.GetX(), xAxis.GetY(), xAxis.GetZ(), 0,
		yAxis.GetX(), yAxis.GetY(), yAxis.GetZ(), 0,
		zAxis.GetX(), zAxis.GetY(), zAxis.GetZ(), 0,
		0, 0, 0, 1,
	)
	qq := mat.Quaternion()
	return qq
}

// SeparateByAxisは、グローバル軸に基づいてクォータニオンを4つのクォータニオンに分割します。
func (quat *MQuaternion) SeparateByAxis(globalAxis *MVec3) (*MQuaternion, *MQuaternion, *MQuaternion, *MQuaternion) {
	localZAxis := MVec3UnitZ
	globalXAxis := globalAxis.Normalize()
	globalYAxis := localZAxis.Cross(globalXAxis)
	globalZAxis := globalXAxis.Cross(globalYAxis)

	if globalYAxis.Length() == 0 {
		localYAxis := MVec3UnitY
		globalZAxis = localYAxis.Cross(globalXAxis)
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

// Copyはクォータニオンのコピーを返します。
func (qq *MQuaternion) Copy() *MQuaternion {
	return NewMQuaternionByValues(qq.GetX(), qq.GetY(), qq.GetZ(), qq.GetW())
}

// Vectorはクォータニオンをベクトルに変換します。
func (v *MQuaternion) Vector() []float64 {
	return []float64{v.GetX(), v.GetY(), v.GetZ(), v.GetW()}
}

// ToMat4はクォータニオンを4x4行列に変換します。
func (v *MQuaternion) ToMat4() *MMat4 {
	m := mgl64.Quat(*v).Normalize().Mat4()
	return NewMMat4ByValues(
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15],
	)
}

// ToFixedAxisRotationは軸制限されたクォータニオンの回転を計算します。
func (quat *MQuaternion) ToFixedAxisRotation(fixedAxis *MVec3) *MQuaternion {
	normalizedFixedAxis := fixedAxis.Normalized()
	quatAxis := quat.GetXYZ().Normalized()
	rad := quat.ToRadian()
	if normalizedFixedAxis.Dot(quatAxis) < 0 {
		rad *= -1
	}
	result := NewMQuaternionFromAxisAngles(normalizedFixedAxis, rad)
	return result
}

func (quat *MQuaternion) IsIdent() bool {
	return quat.PracticallyEquals(&MQuaternionIdent, 1e-6)
}

// PracticallyEqualsは2つのクォータニオンがほぼ等しいかどうかを判定します。
func (quat *MQuaternion) PracticallyEquals(other *MQuaternion, epsilon float64) bool {
	return mgl64.Quat(*quat).ApproxEqualThreshold(mgl64.Quat(*other), epsilon)
}

// MulVec3は、ベクトルvをクォータニオンで回転させた結果の新しいベクトルを返します。
func (quat *MQuaternion) MulVec3(v *MVec3) *MVec3 {
	r := mgl64.Quat(*quat).Rotate(mgl64.Vec3(*v))
	return &MVec3{r[0], r[1], r[2]}
}

// VectorToDegreeは、与えられた2つのベクトルから角度に変換します。
func VectorToDegree(a *MVec3, b *MVec3) float64 {
	return RadToDeg(VectorToRadian(a, b))
}

// VectorToRadianは、与えられた2つのベクトルからラジアン角度に変換します。
func VectorToRadian(a *MVec3, b *MVec3) float64 {
	p := a.Dot(b)
	normA := a.Length()
	normB := b.Length()

	// 角度を計算
	cosAngle := p / (normA * normB)
	rad := math.Acos(math.Min(1, math.Max(-1, cosAngle)))

	return rad
}

// FindSlerpTは始点Q1、終点Q2、中間点Qtが与えられたとき、Slerp(Q1, Q2, t) ? Qtとなるtを見つけます。
func FindSlerpT(Q1, Q2, Qt *MQuaternion) float64 {
	tol := 1e-15
	return findSlerpTGoldenSection(Q1, Q2, Qt, tol)
}

// findSlerpTGoldenSectionは一貫したクォータニオンサインを確保した上でtを見つけます。
func findSlerpTGoldenSection(Q1, Q2, Qt *MQuaternion, tol float64) float64 {
	phi := (1 + math.Sqrt(5)) / 2
	maxIterations := 100

	// 初期範囲の設定
	a := 0.0
	b := 1.0
	c := b - (b-a)/phi
	d := a + (b-a)/phi

	q2 := Q2
	if Q1.Dot(Q2) < 0 {
		q2 = Q2.Negated()
	}
	Q2 = q2

	// 誤差の計算関数
	errorFunc := func(t float64) float64 {
		tQuat := Q1.Slerp(Q2, t)
		dot := math.Abs(tQuat.Dot(Qt))
		return 1 - dot
	}

	// 初期の誤差計算
	fc := errorFunc(c)
	fd := errorFunc(d)

	for i := 0; i < maxIterations; i++ {
		if math.Abs(b-a) < tol {
			return (a + b) / 2
		}
		if fc < fd {
			b = d
			d = c
			fd = fc
			c = b - (b-a)/phi
			fc = errorFunc(c)
		} else {
			a = c
			c = d
			fc = fd
			d = a + (b-a)/phi
			fd = errorFunc(d)
		}
	}

	// 終了条件に達したら範囲の中間点を返す
	return (a + b) / 2
}

// FindSlerpTBisectionは始点Q1、終点Q2、中間点Qtが与えられたとき、Slerp(Q1, Q2, t) ? Qtとなるtを二分法で見つけます。
func FindSlerpTBisection(Q1, Q2, Qt *MQuaternion, tol float64) float64 {
	low := 0.00001
	high := 0.99999
	mid := (low + high) / 2

	maxIterations := 50

	for i := 0; i < maxIterations; i++ {
		// Slerpで中間のクオータニオンを計算し、誤差を測定
		midQuat := Q1.Slerp(Q2, mid)
		dot := math.Abs(midQuat.Dot(Qt))
		err := 1 - dot

		// 誤差が許容範囲内であれば終了
		if err < tol {
			return mid
		}

		// 中間点の誤差に基づいて範囲を狭める
		lowQuat := Q1.Slerp(Q2, low)
		if lowQuat.Dot(Qt) < midQuat.Dot(Qt) {
			low = mid
		} else {
			high = mid
		}

		mid = (low + high) / 2
	}

	return mid
}
