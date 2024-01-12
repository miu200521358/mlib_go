package mquaternion

import (
	"fmt"
	"math"

	"github.com/ungerik/go3d/float64/mat3"
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"

	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"
)

type T quaternion.T

var (
	// Zero holds a zero quaternion.
	Zero = T{}

	// Ident holds an ident quaternion.
	Ident = T{0, 0, 0, 1}
)

// GetX returns the value of the X coordinate
func (v *T) GetX() float64 {
	return v[0]
}

// SetX sets the value of the X coordinate
func (v *T) SetX(x float64) {
	v[0] = x
}

// GetY returns the value of the Y coordinate
func (v *T) GetY() float64 {
	return v[1]
}

// SetY sets the value of the Y coordinate
func (v *T) SetY(y float64) {
	v[1] = y
}

// GetZ returns the value of the Z coordinate
func (v *T) GetZ() float64 {
	return v[2]
}

// SetZ sets the value of the Z coordinate
func (v *T) SetZ(z float64) {
	v[2] = z
}

// GetW returns the value of the W coordinate
func (v *T) GetW() float64 {
	return v[3]
}

// SetW sets the value of the W coordinate
func (v *T) SetW(w float64) {
	v[3] = w
}

// String T の文字列表現を返します。
func (v *T) String() string {
	return fmt.Sprintf("[x=%.5f, y=%.5f, z=%.5f, w=%.5f]", v.GetX(), v.GetY(), v.GetZ(), v.GetW())
}

// GL OpenGL座標系に変換されたクォータニオンベクトルを返します
func (v *T) GL() T {
	return T{-v.GetX(), v.GetY(), v.GetZ(), -v.GetW()}
}

// MMD MMD(MikuMikuDance)座標系に変換されたクォータニオンベクトルを返します
func (v *T) MMD() T {
	return T{v.GetX(), -v.GetY(), -v.GetZ(), v.GetW()}
}

// FromAxisAngle は、軸周りの回転を表す四元数を返します。
func FromAxisAngle(axis *mvec3.T, angle float64) T {
	return T(quaternion.FromAxisAngle((*vec3.T)(axis), angle))
}

// FromXAxisAngleは、X軸周りの回転を表す四元数を返します。
func FromXAxisAngle(angle float64) T {
	return T(quaternion.FromXAxisAngle(angle))
}

// FromYAxisAngleは、Y軸周りの回転を表す四元数を返します。
func FromYAxisAngle(angle float64) T {
	return T(quaternion.FromYAxisAngle(angle))
}

// FromZAxisAngleは、Z軸周りの回転を表す四元数を返します。
func FromZAxisAngle(angle float64) T {
	return T(quaternion.FromZAxisAngle(angle))
}

// FromEulerAnglesは、オイラー角（ラジアン）回転を表す四元数を返します。
func FromEulerAngles(xHead, yPitch, zRoll float64) T {
	return T(quaternion.FromEulerAngles(xHead, yPitch, zRoll))
}

// ToEulerAnglesは、クォータニオンのオイラー角（ラジアン）回転を返します。
func (quat *T) ToEulerAngles() mvec3.T {
	xHead, yPitch, zRoll := (*quaternion.T)(quat).ToEulerAngles()
	return mvec3.T{xHead, yPitch, zRoll}
}

// FromEulerAnglesDegreesは、オイラー角（度）回転を表す四元数を返します。
func FromEulerAnglesDegrees(xHead, yPitch, zRoll float64) T {
	xHeadRadian := math.Pi * xHead / 180.0
	yPitchRadian := math.Pi * yPitch / 180.0
	zRollRadian := math.Pi * zRoll / 180.0
	return T(quaternion.FromEulerAngles(yPitchRadian, xHeadRadian, zRollRadian))
}

// ToEulerAnglesDegreesは、クォータニオンのオイラー角（度）回転を返します。
func (quat *T) ToEulerAnglesDegrees() mvec3.T {
	xHead, yPitch, zRoll := (*quaternion.T)(quat).ToEulerAngles()
	return mvec3.T{
		180.0 * xHead / math.Pi,
		180.0 * yPitch / math.Pi,
		180.0 * zRoll / math.Pi,
	}
}

// FromVec4はvec4.Tをクォータニオンに変換する
func FromVec4(v *mvec4.T) T {
	return T(*v)
}

// Vec4は四元数をvec4.Tに変換する
func (quat *T) Vec4() *mvec4.T {
	return (*mvec4.T)(quat)
}

// Vec3は、クォータニオンのベクトル部分を返します。
func (quat *T) Vec3() *mvec3.T {
	vec3 := mvec3.T{quat.GetX(), quat.GetY(), quat.GetZ()}
	return &vec3
}

// AxisAngleは、正規化されたクォータニオンから、軸と回転角度の形で回転を取り出す。
func (quat *T) AxisAngle() (axis *mvec3.T, angle float64) {
	axisV3, angle := (*quaternion.T)(quat).AxisAngle()
	axis = &mvec3.T{axisV3[0], axisV3[1], axisV3[2]}
	return axis, angle
}

// Mul は、クォータニオンの積を返します。
func (quat *T) Mul(other *T) *T {
	mulQuat := quaternion.Mul((*quaternion.T)(quat), (*quaternion.T)(other))
	return (*T)(&mulQuat)
}

// Norm はクォータニオンのノルム値を返します。
func (quat *T) Norm() float64 {
	return (*quaternion.T)(quat).Norm()
}

// Normalizeは、単位四位数に正規化する。
func (quat *T) Normalize() *T {
	return (*T)((*quaternion.T)(quat).Normalize())
}

// Normalizedは、単位を4進数に正規化したコピーを返す。
func (quat *T) Normalized() T {
	return T((*quaternion.T)(quat).Normalized())
}

// Negate negates the quaternion.
func (quat *T) Negate() {
	(*quaternion.T)(quat).Negate()
}

// Negated returns a negated quaternion.
func (quat *T) Negated() T {
	return T((*quaternion.T)(quat).Negated())
}

// Invert inverts the quaternion.
func (quat *T) Invert() {
	(*quaternion.T)(quat).Invert()
}

// Inverted returns a inverted quaternion.
func (quat *T) Inverted() T {
	return T((*quaternion.T)(quat).Inverted())
}

// SetShortestRotation は、クォータニオンが quat から other の方向への最短回転を表していない場合、そのクォータニオンを否定します。
// (quatの向きからotherの向きへの回転には2つの方向があります)
func (quat *T) SetShortestRotation(other *T) *T {
	return (*T)((*quaternion.T)(quat).SetShortestRotation((*quaternion.T)(other)))
}

// IsShortestRotation は、a から b への回転が可能な限り最短の回転かどうかを返す。
// (quatの向きから他の向きへの回転には2つの方向がある)
func (quat *T) IsShortestRotation(other *T) bool {
	return quaternion.IsShortestRotation((*quaternion.T)(quat), (*quaternion.T)(other))
}

// IsUnitQuat は、クォータニオンが単位クォータニオンの許容範囲内にあるかどうかを返します。
func (quat *T) IsUnitQuat(tolerance float64) bool {
	return (*quaternion.T)(quat).IsUnitQuat(tolerance)
}

// RotateVec3 は、四元数によって表される回転によって v を回転させます。
// https://gamedev.stackexchange.com/questions/28395/rotating-vector3-by-a-quaternion
func (quat *T) RotateVec3(v *mvec3.T) {
	(*quaternion.T)(quat).RotateVec3((*vec3.T)(v))
}

// RotatedVec3 は v の回転コピーを返す。
// https://gamedev.stackexchange.com/questions/28395/rotating-vector3-by-a-quaternion
func (quat *T) RotatedVec3(v *mvec3.T) mvec3.T {
	return mvec3.T((*quaternion.T)(quat).RotatedVec3((*vec3.T)(v)))
}

// Dot は2つのクォータニオンの内積を返す。
func (quat *T) Dot(other *T) float64 {
	return quaternion.Dot((*quaternion.T)(quat), (*quaternion.T)(other))
}

// Mul
func Mul(a, b *T) T {
	return T(quaternion.Mul((*quaternion.T)(a), (*quaternion.T)(b)))
}

// Slerp は t (0,1) における a と b の間の球面線形補間クォータニオンを返す。
// See http://en.wikipedia.org/wiki/Slerp
func Slerp(a, b *T, t float64) T {
	return T(quaternion.Slerp((*quaternion.T)(a), (*quaternion.T)(b), t))
}

// Vec3Diff関数は、2つのベクトル間の回転四元数を返します。
func Vec3Diff(a, b *mvec3.T) T {
	return T(quaternion.Vec3Diff((*vec3.T)(a), (*vec3.T)(b)))
}

// ToDegree は、クォータニオンを度に変換します。
func (quat *T) ToDegree() float64 {
	w := quat.Normalize().GetW()
	radian := 2 * math.Acos(math.Min(1, math.Max(-1, w)))
	angle := radian * (180 / math.Pi)
	return angle
}

// ToRadian は、クォータニオンをラジアンに変換します。
func (quat *T) ToRadian() float64 {
	w := quat.Normalize().GetW()
	radian := 2 * math.Acos(math.Min(1, math.Max(-1, w)))
	return radian
}

// ToSignedDegree 符号付き角度に変換
func (quat *T) ToSignedDegree() float64 {
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
func (quat *T) ToSignedRadian() float64 {
	// スカラー部分から基本的な角度を計算
	basicAngle := quat.ToRadian()

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

// ToTheta 自分ともうひとつの値vとのtheta（変位量）を返す
func (quat *T) ToTheta(v *T) float64 {
	return math.Acos(math.Min(1, math.Max(-1, quat.Normalize().Dot(v.Normalize()))))
}

// 軸と角度からクォータニオンに変換する
func FromDirection(direction *mvec3.T, up *mvec3.T) *T {
	if direction.Length() == 0 {
		return &T{}
	}

	zAxis := direction.Normalize()
	xAxis := up.Cross(zAxis).Normalize()

	if xAxis.LengthSqr() == 0 {
		// collinear or invalid up vector derive shortest arc to new direction
		return Rotate(&mvec3.T{0.0, 0.0, 1.0}, zAxis)
	}

	yAxis := zAxis.Cross(xAxis)

	return FromAxes(xAxis, yAxis, zAxis).Normalize()
}

// Rotate fromベクトルからtoベクトルまでの回転量
func Rotate(fromV, toV *mvec3.T) *T {
	v0 := fromV.Normalize()
	v1 := toV.Normalize()
	d := v0.Dot(v1) + 1.0

	// if dest vector is close to the inverse of source vector, ANY axis of rotation is valid
	if math.Abs(d) < 1e-6 {
		axis := mvec3.UnitX.Cross(v0)
		if math.Abs(axis.LengthSqr()) < 1e-6 {
			axis = mvec3.UnitY.Cross(v0)
		}
		axis.Normalize()
		// same as MQuaternion.fromAxisAndAngle(axis, 180.0)
		return &T{axis.GetX(), axis.GetY(), axis.GetZ(), 0.0}
	}

	d = math.Sqrt(2.0 * d)
	axis := v0.Cross(v1).DivScalar(d)
	return &T{axis.GetX(), axis.GetY(), axis.GetZ(), d * 0.5}
}

// FromAxes
func FromAxes(xAxis, yAxis, zAxis *mvec3.T) *T {
	quat := mat3.Ident.AssignCoordinateSystem((*vec3.T)(xAxis), (*vec3.T)(yAxis), (*vec3.T)(zAxis)).Quaternion()
	return (*T)(&quat)
}

// SeparateByAxis separates the quaternion into four quaternions based on the global axis.
func (quat *T) SeparateByAxis(globalAxis *mvec3.T) (*T, *T, *T, *T) {
	localZAxis := &mvec3.UnitZ
	globalXAxis := globalAxis.Normalize()
	globalYAxis := localZAxis.Cross(globalXAxis)
	globalZAxis := globalXAxis.Cross(globalYAxis)

	if globalYAxis.Length() == 0 {
		localYAxis := &mvec3.UnitY
		globalZAxis := localYAxis.Cross(globalXAxis)
		globalYAxis = globalXAxis.Cross(globalZAxis)
	}

	// X成分を抽出する ------------

	// グローバル軸方向に伸ばす
	globalXVec := quat.RotatedVec3(globalXAxis)
	// YZの回転量（自身のねじれを無視する）
	yzQQ := Rotate(globalXAxis, globalXVec.Normalize())
	// 元々の回転量 から YZ回転 を除去して、除去されたX成分を求める
	invYzQQ := yzQQ.Inverted()
	xQQ := quat.Mul(&invYzQQ)

	// Y成分を抽出する ------------
	// グローバル軸方向に伸ばす
	globalYVec := quat.RotatedVec3(globalYAxis)
	// XZの回転量（自身のねじれを無視する）
	xzQQ := Rotate(globalYAxis, globalYVec.Normalize())
	// 元々の回転量 から XZ回転 を除去して、除去されたY成分を求める
	invXzQQ := xzQQ.Inverted()
	yQQ := quat.Mul(&invXzQQ)

	// Z成分を抽出する ------------
	// グローバル軸方向に伸ばす
	globalZVec := quat.RotatedVec3(globalZAxis)
	// XYの回転量（自身のねじれを無視する）
	xyQQ := Rotate(globalZAxis, globalZVec.Normalize())
	// 元々の回転量 から XY回転 を除去して、除去されたZ成分を求める
	invXyQQ := xyQQ.Inverted()
	zQQ := quat.Mul(&invXyQQ)

	return xQQ, yQQ, zQQ, yzQQ
}

// Copy
func (qq *T) Copy() *T {
	copied := *qq
	return &copied
}
