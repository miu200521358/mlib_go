package mmath

import (
	"errors"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

var (
	// Zero holds a zero quaternion.
	MQuaternionZero = &MQuaternion{}

	// Ident holds an ident quaternion.
	MQuaternionIdent = &MQuaternion{0, 0, 0, 1}

	MQuaternionUnitX = &MQuaternion{1, 0, 0, 0}
	MQuaternionUnitY = &MQuaternion{0, 1, 0, 0}
	MQuaternionUnitZ = &MQuaternion{0, 0, 1, 0}
)

type MQuaternion struct {
	X float64
	Y float64
	Z float64
	W float64
}

func NewMQuaternion() *MQuaternion {
	return &MQuaternion{X: 0, Y: 0, Z: 0, W: 1}
}

// 指定された値でクォータニオンを作成します。
// ただし必ず最短距離クォータニオンにします
func NewMQuaternionByValuesShort(x, y, z, w float64) *MQuaternion {
	qq := &MQuaternion{X: x, Y: y, Z: z, W: w}
	if !MQuaternionIdent.IsShortestRotation(qq) {
		qq.Negate()
	}
	return qq
}

// NewMQuaternionByValuesOriginal は、指定された値でクォータニオンを作成します。
// ただし、強制的に最短距離クォータニオンにはしません
func NewMQuaternionByValues(x, y, z, w float64) *MQuaternion {
	return &MQuaternion{X: x, Y: y, Z: z, W: w}
}

func (quat *MQuaternion) XYZ() *MVec3 {
	return &MVec3{quat.X, quat.Y, quat.Z}
}

func (quat *MQuaternion) SetXYZ(v3 *MVec3) {
	quat.X = v3.X
	quat.Y = v3.Y
	quat.Z = v3.Z
}

// String T の文字列表現を返します。
func (quat *MQuaternion) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", quat.X, quat.Y, quat.Z, quat.W)
}

// MMD MMD(MikuMikuDance)座標系に変換されたクォータニオンベクトルを返します
func (quat *MQuaternion) MMD() *MQuaternion {
	return &MQuaternion{quat.X, quat.Y, quat.Z, quat.W}
}

// NewMQuaternionFromAxisAngles は、軸周りの回転を表す四元数を返します。
func NewMQuaternionFromAxisAngles(axis *MVec3, angle float64) *MQuaternion {
	axis.Normalize()
	m := MMat4(mgl64.HomogRotate3D(angle, mgl64.Vec3{axis.X, axis.Y, axis.Z}))
	return m.Quaternion()
}

// NewMQuaternionFromAxisAnglesRotate は、軸周りの回転を表す四元数を返します。
func NewMQuaternionFromAxisAnglesRotate(axis *MVec3, angle float64) *MQuaternion {
	x := axis.Normalized()
	m := mgl64.QuatRotate(angle, mgl64.Vec3{x.X, x.Y, x.Z}).Normalize()
	return &MQuaternion{m.X(), m.Y(), m.Z(), m.W}
}

// NewMQuaternionFromRadiansは、オイラー角（ラジアン）回転を表す四元数を返します。
func NewMQuaternionFromRadians(xPitch, yHead, zRoll float64) *MQuaternion {
	q := mgl64.AnglesToQuat(xPitch, yHead, zRoll, mgl64.XYZ).Normalize()
	return &MQuaternion{q.X(), q.Y(), q.Z(), q.W}
}

// 参考URL:
// https://qiita.com/aa_debdeb/items/abe90a9bd0b4809813da
// https://site.nicovideo.jp/ch/userblomaga_thanks/archive/ar805999

// ToRadiansは、クォータニオンを三軸のオイラー角（ラジアン）回転を返します。
func (quat *MQuaternion) ToRadians() *MVec3 {
	sx := -(2*quat.Y*quat.Z - 2*quat.X*quat.W)
	unlocked := math.Abs(sx) < 0.99999
	xPitch := math.Asin(math.Max(-1, math.Min(1, sx)))
	var yHead, zRoll float64
	if unlocked {
		yHead = math.Atan2(2*quat.X*quat.Z+2*quat.Y*quat.W, 2*quat.W*quat.W+2*quat.Z*quat.Z-1)
		zRoll = math.Atan2(2*quat.X*quat.Y+2*quat.Z*quat.W, 2*quat.W*quat.W+2*quat.Y*quat.Y-1)
	} else {
		yHead = math.Atan2(-(2*quat.X*quat.Z - 2*quat.Y*quat.W), 2*quat.W*quat.W+2*quat.X*quat.X-1)
		zRoll = 0
	}

	return &MVec3{xPitch, yHead, zRoll}
}

const (
	GIMBAL1_RAD = 88.0 / 180.0 * math.Pi
	GIMBAL2_RAD = GIMBAL1_RAD * 2
	ONE_RAD     = math.Pi
	HALF_RAD    = math.Pi / 2
)

// ToRadiansWithGimbalは、クォータニオンを三軸のオイラー角（ラジアン）回転を返します。
// ジンバルロックが発生しているか否かのフラグも返します
func (quat *MQuaternion) ToRadiansWithGimbal(axisIndex int) (*MVec3, bool) {
	r := quat.ToRadians()

	var other1Rad, other2Rad float64
	if axisIndex == 0 {
		other1Rad = math.Abs(r.Y)
		other2Rad = math.Abs(r.Z)
	} else if axisIndex == 1 {
		other1Rad = math.Abs(r.X)
		other2Rad = math.Abs(r.Z)
	} else {
		other1Rad = math.Abs(r.X)
		other2Rad = math.Abs(r.Y)
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

// ToDegreesは、クォータニオンのオイラー角（度）回転を返します。
func (quat *MQuaternion) ToDegrees() *MVec3 {
	vec := quat.ToRadians()
	return &MVec3{
		RadToDeg(vec.X),
		RadToDeg(vec.Y),
		RadToDeg(vec.Z),
	}
}

// ToDegreesは、クォータニオンのオイラー角（度）回転を返します。
func (quat *MQuaternion) ToMMDDegrees() *MVec3 {
	vec := quat.MMD().ToRadians()
	return &MVec3{
		RadToDeg(vec.X),
		RadToDeg(-vec.Y),
		RadToDeg(-vec.Z),
	}
}

// Vec4は四元数をvec4.Tに変換する
func (quat *MQuaternion) Vec4() *MVec4 {
	return &MVec4{quat.X, quat.Y, quat.Z, quat.W}
}

// Vec3は、クォータニオンのベクトル部分を返します。
func (quat *MQuaternion) Vec3() *MVec3 {
	vec3 := MVec3{quat.X, quat.Y, quat.Z}
	return &vec3
}

// Mulは、クォータニオンの積を返します。
func (quat1 *MQuaternion) MulShort(quat2 *MQuaternion) *MQuaternion {
	// 直接クォータニオン乗算を実装して行列変換を回避
	x := quat1.W*quat2.X + quat1.X*quat2.W + quat1.Y*quat2.Z - quat1.Z*quat2.Y
	y := quat1.W*quat2.Y + quat1.Y*quat2.W + quat1.Z*quat2.X - quat1.X*quat2.Z
	z := quat1.W*quat2.Z + quat1.Z*quat2.W + quat1.X*quat2.Y - quat1.Y*quat2.X
	w := quat1.W*quat2.W - quat1.X*quat2.X - quat1.Y*quat2.Y - quat1.Z*quat2.Z

	return NewMQuaternionByValues(x, y, z, w)
}

func (q1 *MQuaternion) MuledShort(q2 *MQuaternion) *MQuaternion {
	copied := q1.Copy()
	copied.Mul(q2)
	return copied
}

// Mulは、クォータニオンの積を返します。
func (quat1 *MQuaternion) Mul(quat2 *MQuaternion) *MQuaternion {
	// 直接クォータニオン乗算を計算
	x := quat1.W*quat2.X + quat1.X*quat2.W + quat1.Y*quat2.Z - quat1.Z*quat2.Y
	y := quat1.W*quat2.Y + quat1.Y*quat2.W + quat1.Z*quat2.X - quat1.X*quat2.Z
	z := quat1.W*quat2.Z + quat1.Z*quat2.W + quat1.X*quat2.Y - quat1.Y*quat2.X
	w := quat1.W*quat2.W - quat1.X*quat2.X - quat1.Y*quat2.Y - quat1.Z*quat2.Z

	quat1.X = x
	quat1.Y = y
	quat1.Z = z
	quat1.W = w
	return quat1
}

func (quat1 *MQuaternion) Muled(quat2 *MQuaternion) *MQuaternion {
	// 直接クォータニオン乗算を計算して新しいインスタンスを返す
	x := quat1.W*quat2.X + quat1.X*quat2.W + quat1.Y*quat2.Z - quat1.Z*quat2.Y
	y := quat1.W*quat2.Y + quat1.Y*quat2.W + quat1.Z*quat2.X - quat1.X*quat2.Z
	z := quat1.W*quat2.Z + quat1.Z*quat2.W + quat1.X*quat2.Y - quat1.Y*quat2.X
	w := quat1.W*quat2.W - quat1.X*quat2.X - quat1.Y*quat2.Y - quat1.Z*quat2.Z

	return &MQuaternion{x, y, z, w}
}

// Normはクォータニオンのノルム値を返します。
func (quat *MQuaternion) Norm() float64 {
	return mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Norm()
}

// Lengthはクォータニオンの長さ（ノルム）を返します。
func (quat *MQuaternion) Length() float64 {
	return mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Len()
}

// Normalizeは、単位四位数に正規化する。
func (quat *MQuaternion) Normalize() *MQuaternion {
	qq := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Normalize()
	*quat = MQuaternion{qq.V[0], qq.V[1], qq.V[2], qq.W}
	return quat
}

// Normalizedは、単位を4進数に正規化したコピーを返す。
func (quat *MQuaternion) Normalized() *MQuaternion {
	qq := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Normalize()
	return &MQuaternion{qq.V[0], qq.V[1], qq.V[2], qq.W}
}

// Negateはクォータニオンを反転する。
func (quat *MQuaternion) Negate() *MQuaternion {
	quat.X *= -1
	quat.Y *= -1
	quat.Z *= -1
	quat.W *= -1
	return quat
}

// Negatedは反転したクォータニオンを返します。
func (quat *MQuaternion) Negated() *MQuaternion {
	return NewMQuaternionByValues(-quat.X, -quat.Y, -quat.Z, -quat.W)
}

// Inverseは、クォータニオンを反転させます。
func (quat *MQuaternion) Inverse() *MQuaternion {
	// 長さの二乗を計算
	lenSq := quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W

	// 長さが0に近い場合は単位クォータニオンを返す
	if lenSq < 1e-10 {
		*quat = MQuaternion{0, 0, 0, 1}
		return quat
	}

	// 単位クォータニオンの場合は共役を取るだけで良い
	if math.Abs(lenSq-1.0) < 1e-10 {
		quat.X = -quat.X
		quat.Y = -quat.Y
		quat.Z = -quat.Z
		return quat
	}

	// 一般的なケース
	invLenSq := 1.0 / lenSq
	quat.X = -quat.X * invLenSq
	quat.Y = -quat.Y * invLenSq
	quat.Z = -quat.Z * invLenSq
	quat.W = quat.W * invLenSq
	return quat
}

// Invertedは反転したクォータニオンを返します。
func (quat *MQuaternion) Inverted() *MQuaternion {
	// 長さの二乗を計算
	lenSq := quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W

	// 長さが0に近い場合は単位クォータニオンを返す
	if lenSq < 1e-10 {
		return &MQuaternion{0, 0, 0, 1}
	}

	// 単位クォータニオンの場合は共役を取るだけで良い
	if math.Abs(lenSq-1.0) < 1e-10 {
		return &MQuaternion{-quat.X, -quat.Y, -quat.Z, quat.W}
	}

	// 一般的なケース
	invLenSq := 1.0 / lenSq
	return &MQuaternion{
		-quat.X * invLenSq,
		-quat.Y * invLenSq,
		-quat.Z * invLenSq,
		quat.W * invLenSq,
	}
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
	if quat.W < 0 {
		quat.Negate()
	}
	return quat
}

// Dotは2つのクォータニオンの内積を返します。
func (quat *MQuaternion) Dot(other *MQuaternion) float64 {
	q1 := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}
	q2 := mgl64.Quat{V: mgl64.Vec3{other.X, other.Y, other.Z}, W: other.W}
	return q1.Dot(q2)
}

// MuledScalarはクォータニオンにスカラーを掛け算します。
func (quat *MQuaternion) MuledScalar(factor float64) *MQuaternion {
	if factor == 0.0 {
		return NewMQuaternion()
	} else if factor == 1.0 {
		return quat.Copy()
	} else if factor == -1.0 {
		return quat.Inverted()
	}

	// 特殊なケースでなければ、軸角度の表現を使って効率的に計算
	axis, angle := quat.ToAxisAngle()
	return NewMQuaternionFromAxisAngles(axis, angle*factor)
}

// ToAxisAngleは、クォータニオンを軸と角度に変換します。
func (quat *MQuaternion) ToAxisAngle() (*MVec3, float64) {
	// 長さを計算してクォータニオンが正規化されているか確認
	lenSq := quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W

	// 正規化が必要な場合
	var normW float64
	if math.Abs(lenSq-1.0) > 1e-10 {
		invLen := 1.0 / math.Sqrt(lenSq)
		normW = quat.W * invLen
	} else {
		normW = quat.W
	}

	// 角度を計算
	angle := 2.0 * math.Acos(math.Max(-1.0, math.Min(1.0, normW)))

	// 軸の計算
	s := math.Sqrt(1.0 - normW*normW)

	// 角度が非常に小さい場合、任意の軸を選択
	if s < 1e-9 {
		return &MVec3{1, 0, 0}, angle
	}

	invS := 1.0 / s
	return &MVec3{
		quat.X * invS,
		quat.Y * invS,
		quat.Z * invS,
	}, angle
}

// Slerpはt (0,1)におけるaとbの間の球面線形補間クォータニオンを返します。
func (quat *MQuaternion) Slerp(other *MQuaternion, t float64) *MQuaternion {
	// 早期リターン条件をチェック
	if t <= 0.0 {
		return quat.Copy()
	}
	if t >= 1.0 {
		return other.Copy()
	}
	if quat.NearEquals(other, 1e-8) {
		return quat.Copy()
	}

	// コサイン値を計算
	cosOmega := quat.X*other.X + quat.Y*other.Y + quat.Z*other.Z + quat.W*other.W

	// 最短経路補間のため、必要に応じて符号を反転
	q2x, q2y, q2z, q2w := other.X, other.Y, other.Z, other.W
	if cosOmega < 0.0 {
		cosOmega = -cosOmega
		q2x = -q2x
		q2y = -q2y
		q2z = -q2z
		q2w = -q2w
	}

	// 角度の補間方法を決定
	var k1, k2 float64
	if cosOmega > 0.9999 {
		// 角度が非常に小さい場合、線形補間を使用
		k1 = 1.0 - t
		k2 = t
	} else {
		// 通常のslerp計算
		sinOmega := math.Sqrt(1.0 - cosOmega*cosOmega)
		omega := math.Atan2(sinOmega, cosOmega)
		invSinOmega := 1.0 / sinOmega

		k1 = math.Sin((1.0-t)*omega) * invSinOmega
		k2 = math.Sin(t*omega) * invSinOmega
	}

	// 補間結果を計算
	return &MQuaternion{
		k1*quat.X + k2*q2x,
		k1*quat.Y + k2*q2y,
		k1*quat.Z + k2*q2z,
		k1*quat.W + k2*q2w,
	}
}

func (quat *MQuaternion) Lerp(other *MQuaternion, t float64) *MQuaternion {
	// 早期リターン条件をチェック
	if t <= 0.0 {
		return quat.Copy()
	}
	if t >= 1.0 {
		return other.Copy()
	}
	if quat.NearEquals(other, 1e-8) {
		return quat.Copy()
	}

	// クォータニオンの線形補間
	scale0 := 1.0 - t
	scale1 := t

	// 最短経路補間のため、必要に応じて符号を反転
	dot := quat.X*other.X + quat.Y*other.Y + quat.Z*other.Z + quat.W*other.W
	if dot < 0 {
		scale1 = -scale1
	}

	// 補間結果を計算
	x := scale0*quat.X + scale1*other.X
	y := scale0*quat.Y + scale1*other.Y
	z := scale0*quat.Z + scale1*other.Z
	w := scale0*quat.W + scale1*other.W

	// 正規化
	len := math.Sqrt(x*x + y*y + z*z + w*w)
	if len > 0 {
		invLen := 1.0 / len
		x *= invLen
		y *= invLen
		z *= invLen
		w *= invLen
	}

	return &MQuaternion{x, y, z, w}
}

// ToDegreeは、クォータニオンを度に変換します。
func (quat *MQuaternion) ToDegree() float64 {
	return RadToDeg(quat.ToRadian())
}

// ToRadianは、クォータニオンをラジアンに変換します。
func (quat *MQuaternion) ToRadian() float64 {
	return 2 * math.Acos(math.Min(1, math.Max(-1, quat.Normalize().W)))
}

// ToSignedDegreeは、符号付き角度に変換します。
func (quat *MQuaternion) ToSignedDegree() float64 {
	// スカラー部分から基本的な角度を計算
	basicAngle := quat.ToDegree()

	// ベクトルの長さを使って、角度の正負を決定
	if quat.Vec3().Length() > 0 {
		// ベクトルの向きに基づいて角度を調整
		if quat.W >= 0 {
			return basicAngle
		} else {
			return -basicAngle
		}
	}

	// ベクトル部分がない場合は基本角度をそのまま使用
	return basicAngle
}

// ToSignedRadianは、符号付きラジアンに変換します。
func (quat *MQuaternion) ToSignedRadian() float64 {
	// スカラー部分から基本的な角度を計算
	basicAngle := quat.ToRadian()

	// ベクトルの長さを使って、角度の正負を決定
	if quat.Vec3().Length() > 0 {
		// ベクトルの向きに基づいて角度を調整
		if quat.W >= 0 {
			return basicAngle
		} else {
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

	zAxis := direction.Normalized()
	xAxis := up.Cross(zAxis).Normalized()

	if xAxis.LengthSqr() == 0 {
		// collinear or invalid up vector derive shortest arc to new direction
		return NewMQuaternionRotate(&MVec3{0.0, 0.0, 1.0}, zAxis)
	}

	yAxis := zAxis.Cross(xAxis)

	return NewMQuaternionFromAxes(xAxis, yAxis, zAxis).Normalize()
}

// NewMQuaternionRotateはfromベクトルからtoベクトルまでの回転量を計算します。
func NewMQuaternionRotate(fromV, toV *MVec3) *MQuaternion {
	if fromV.NearEquals(toV, 1e-6) || fromV.Length() == 0 || toV.Length() == 0 {
		return NewMQuaternion()
	}
	v := mgl64.QuatBetweenVectors(mgl64.Vec3{fromV.X, fromV.Y, fromV.Z}, mgl64.Vec3{toV.X, toV.Y, toV.Z})
	return NewMQuaternionByValues(v.V[0], v.V[1], v.V[2], v.W)
}

// NewMQuaternionFromAxesは、3つの軸ベクトルからクォータニオンを作成します。
func NewMQuaternionFromAxes(xAxis, yAxis, zAxis *MVec3) *MQuaternion {
	mat := NewMMat4ByValues(
		xAxis.X, xAxis.Y, xAxis.Z, 0,
		yAxis.X, yAxis.Y, yAxis.Z, 0,
		zAxis.X, zAxis.Y, zAxis.Z, 0,
		0, 0, 0, 1,
	)
	qq := mat.Quaternion()
	return qq
}

// SeparateByAxisは、グローバル軸に基づいてクォータニオンを2つのクォータニオン(捩りとそれ以外)に分割します。
func (quat *MQuaternion) SeparateTwistByAxis(globalAxis *MVec3) (twistQQ *MQuaternion, yzQQ *MQuaternion) {
	globalXAxis := globalAxis.Normalized()

	// X成分を抽出する ------------

	// グローバル軸方向に伸ばす
	globalXVec := quat.MulVec3(globalXAxis)
	// YZの回転量（自身のねじれを無視する）
	yzQQ = NewMQuaternionRotate(globalXAxis, globalXVec.Normalize())
	// 元々の回転量 から YZ回転 を除去して、除去されたX成分を求める
	twistQQ = yzQQ.Inverted().Mul(quat)

	return twistQQ, yzQQ
}

// SeparateByAxisは、グローバル軸に基づいてクォータニオンを3つのクォータニオン(x, y, z)に分割します。
func (quat *MQuaternion) SeparateByAxis(globalAxis *MVec3) (*MQuaternion, *MQuaternion, *MQuaternion) {
	localZAxis := MVec3{0, 0, -1}
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
	globalXVec := quat.MulVec3(globalXAxis)
	// YZの回転量（自身のねじれを無視する）
	yzQQ := NewMQuaternionRotate(globalXAxis, globalXVec.Normalize())
	// 元々の回転量 から YZ回転 を除去して、除去されたX成分を求める
	xQQ := yzQQ.Inverse().Mul(quat)

	// Y成分を抽出する ------------
	// グローバル軸方向に伸ばす
	globalYVec := quat.MulVec3(globalYAxis)
	// XZの回転量（自身のねじれを無視する）
	xzQQ := NewMQuaternionRotate(globalYAxis, globalYVec.Normalize())
	// 元々の回転量 から XZ回転 を除去して、除去されたY成分を求める
	yQQ := xzQQ.Inverse().Mul(quat)

	// Z成分を抽出する ------------
	// グローバル軸方向に伸ばす
	globalZVec := quat.MulVec3(globalZAxis)
	// XYの回転量（自身のねじれを無視する）
	xyQQ := NewMQuaternionRotate(globalZAxis, globalZVec.Normalize())
	// 元々の回転量 から XY回転 を除去して、除去されたZ成分を求める
	zQQ := xyQQ.Inverse().Mul(quat)

	return xQQ, yQQ, zQQ
}

// Copyはクォータニオンのコピーを返します。
func (quat *MQuaternion) Copy() *MQuaternion {
	return NewMQuaternionByValues(quat.X, quat.Y, quat.Z, quat.W)
}

// Vectorはクォータニオンをベクトルに変換します。
func (quat *MQuaternion) Vector() []float64 {
	return []float64{quat.X, quat.Y, quat.Z, quat.W}
}

// ToMat4はクォータニオンを4x4行列に変換します。
func (quat *MQuaternion) ToMat4() *MMat4 {
	m := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Mat4()
	return (*MMat4)(&m)
}

// ToFixedAxisRotationは軸制限されたクォータニオンの回転を計算します。
func (quat *MQuaternion) ToFixedAxisRotation(fixedAxis *MVec3) *MQuaternion {
	normalizedFixedAxis := fixedAxis.Normalized()
	quatAxis := quat.XYZ().Normalized()
	rad := quat.ToRadian()
	if normalizedFixedAxis.Dot(quatAxis) < 0 {
		rad *= -1
	}
	return NewMQuaternionFromAxisAngles(normalizedFixedAxis, rad)
}

func (quat *MQuaternion) IsIdent() bool {
	return quat.NearEquals(MQuaternionIdent, 1e-6)
}

// NearEqualsは2つのクォータニオンがほぼ等しいかどうかを判定します。
func (quat *MQuaternion) NearEquals(other *MQuaternion, epsilon float64) bool {
	q1 := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}
	q2 := mgl64.Quat{V: mgl64.Vec3{other.X, other.Y, other.Z}, W: other.W}
	return q1.ApproxEqualThreshold(q2, epsilon)
}

// MulVec3は、クォータニオン分ベクトルを回した結果を返します
func (quat *MQuaternion) MulVec3(v *MVec3) *MVec3 {
	// クォータニオンを使った直接的なベクトル回転の最適化実装
	// 行列変換を避けて計算速度を向上

	// クォータニオンの各要素を取得
	qx, qy, qz, qw := quat.X, quat.Y, quat.Z, quat.W

	// ベクトルの各要素を取得
	vx, vy, vz := v.X, v.Y, v.Z

	// 計算の一時変数
	twoQx, twoQy, twoQz := 2.0*qx, 2.0*qy, 2.0*qz

	xx, xy, xz := qx*twoQx, qx*twoQy, qx*twoQz
	yy, yz, zz := qy*twoQy, qy*twoQz, qz*twoQz
	wx, wy, wz := qw*twoQx, qw*twoQy, qw*twoQz

	// 回転後のベクトルを計算
	x := vx*(1.0-(yy+zz)) + vy*(xy-wz) + vz*(xz+wy)
	y := vx*(xy+wz) + vy*(1.0-(xx+zz)) + vz*(yz-wx)
	z := vx*(xz-wy) + vy*(yz+wx) + vz*(1.0-(xx+yy))

	return &MVec3{x, y, z}
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

func (q MQuaternion) Log() (MQuaternion, error) {
	if math.Abs(q.W) > 1.0 {
		return MQuaternion{}, errors.New("invalid quaternion scalar part: must be within [-1, 1]")
	}

	vNorm := q.Norm()
	if vNorm == 0 {
		return MQuaternion{W: 1, X: 0, Y: 0, Z: 0}, nil // Logarithm of a pure scalar quaternion
	}

	angle := math.Acos(q.W)
	scale := angle / vNorm

	return MQuaternion{
		W: 0,
		X: scale * q.X,
		Y: scale * q.Y,
		Z: scale * q.Z,
	}, nil
}

// FindSlerpTは始点Q1、終点Q2、中間点Qtが与えられたとき、Slerp(Q0, Q1, t) = Qtとなるtを見つけます。
func FindSlerpT(Q1, Q2, Qt *MQuaternion, initialT float64) float64 {
	tol := 1e-10
	phi := (1 + math.Sqrt(5)) / 2
	maxIterations := 100

	if math.Abs(Q1.Dot(Q2)) > 0.9999 {
		return initialT
	}

	a := 0.0
	b := 1.0
	c := b - (b-a)/phi
	d := a + (b-a)/phi

	q2 := Q2
	if Q1.Dot(Q2) < 0 {
		q2 = Q2.Negated()
	}

	errorFunc := func(t float64) float64 {
		tQuat := Q1.Slerp(q2, t)
		theta := math.Acos(tQuat.Dot(Qt))
		return theta
	}

	fc := errorFunc(c)
	fd := errorFunc(d)

	for i := 0; i < maxIterations; i++ {
		if math.Abs(b-a) < tol || math.Min(fc, fd) < tol {
			break
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

	return (a + b) / 2
}

func FindLerpT(Q1, Q2, Qt *MQuaternion) float64 {
	tol := 1e-8
	phi := (1 + math.Sqrt(5)) / 2
	maxIterations := 100

	a := 0.0
	b := 1.0
	c := b - (b-a)/phi
	d := a + (b-a)/phi

	q2 := Q2
	if Q1.Dot(Q2) < 0 {
		q2 = Q2.Negated()
	}

	errorFunc := func(t float64) float64 {
		tQuat := Q1.Lerp(q2, t)
		theta := math.Acos(tQuat.Dot(Qt))
		return theta
	}

	fc := errorFunc(c)
	fd := errorFunc(d)

	for i := 0; i < maxIterations; i++ {
		if math.Abs(b-a) < tol || math.Min(fc, fd) < tol {
			break
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

	return (a + b) / 2
}
