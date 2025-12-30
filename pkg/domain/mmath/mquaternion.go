// 指示: miu200521358
package mmath

import (
	"errors"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

var (
	MQuaternionZero = &MQuaternion{}

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

func NewMQuaternionByValuesShort(x, y, z, w float64) *MQuaternion {
	qq := &MQuaternion{X: x, Y: y, Z: z, W: w}
	if !MQuaternionIdent.IsShortestRotation(qq) {
		qq.Negate()
	}
	return qq
}

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

func (quat *MQuaternion) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", quat.X, quat.Y, quat.Z, quat.W)
}

func (quat *MQuaternion) MMD() *MQuaternion {
	return &MQuaternion{quat.X, quat.Y, quat.Z, quat.W}
}

func NewMQuaternionFromAxisAngles(axis *MVec3, angle float64) *MQuaternion {
	axis.Normalize()
	m := MMat4(mgl64.HomogRotate3D(angle, mgl64.Vec3{axis.X, axis.Y, axis.Z}))
	return m.Quaternion()
}

func NewMQuaternionFromAxisAnglesRotate(axis *MVec3, angle float64) *MQuaternion {
	x := axis.Normalized()
	m := mgl64.QuatRotate(angle, mgl64.Vec3{x.X, x.Y, x.Z}).Normalize()
	return &MQuaternion{m.X(), m.Y(), m.Z(), m.W}
}

func NewMQuaternionFromRadians(xPitch, yHead, zRoll float64) *MQuaternion {
	q := mgl64.AnglesToQuat(xPitch, yHead, zRoll, mgl64.XYZ).Normalize()
	return &MQuaternion{q.X(), q.Y(), q.Z(), q.W}
}

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

	if other1Rad >= GIMBAL2_RAD && other2Rad >= GIMBAL2_RAD {
		return r, true
	}

	return r, false
}

func NewMQuaternionFromDegrees(xPitch, yHead, zRoll float64) *MQuaternion {
	xPitchRadian := DegToRad(xPitch)
	yHeadRadian := DegToRad(yHead)
	zRollRadian := DegToRad(zRoll)
	return NewMQuaternionFromRadians(xPitchRadian, yHeadRadian, zRollRadian)
}

func (quat *MQuaternion) ToDegrees() *MVec3 {
	vec := quat.ToRadians()
	return &MVec3{
		RadToDeg(vec.X),
		RadToDeg(vec.Y),
		RadToDeg(vec.Z),
	}
}

func (quat *MQuaternion) ToMMDDegrees() *MVec3 {
	vec := quat.MMD().ToRadians()
	return &MVec3{
		RadToDeg(vec.X),
		RadToDeg(-vec.Y),
		RadToDeg(-vec.Z),
	}
}

func (quat *MQuaternion) Vec4() *MVec4 {
	return &MVec4{quat.X, quat.Y, quat.Z, quat.W}
}

func (quat *MQuaternion) Vec3() *MVec3 {
	vec3 := MVec3{quat.X, quat.Y, quat.Z}
	return &vec3
}

func (quat1 *MQuaternion) MulShort(quat2 *MQuaternion) *MQuaternion {
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

func (quat1 *MQuaternion) Mul(quat2 *MQuaternion) *MQuaternion {
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
	x := quat1.W*quat2.X + quat1.X*quat2.W + quat1.Y*quat2.Z - quat1.Z*quat2.Y
	y := quat1.W*quat2.Y + quat1.Y*quat2.W + quat1.Z*quat2.X - quat1.X*quat2.Z
	z := quat1.W*quat2.Z + quat1.Z*quat2.W + quat1.X*quat2.Y - quat1.Y*quat2.X
	w := quat1.W*quat2.W - quat1.X*quat2.X - quat1.Y*quat2.Y - quat1.Z*quat2.Z

	return &MQuaternion{x, y, z, w}
}

func (quat *MQuaternion) Norm() float64 {
	return mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Norm()
}

func (quat *MQuaternion) Length() float64 {
	return mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Len()
}

func (quat *MQuaternion) Normalize() *MQuaternion {
	qq := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Normalize()
	*quat = MQuaternion{qq.V[0], qq.V[1], qq.V[2], qq.W}
	return quat
}

func (quat *MQuaternion) Normalized() *MQuaternion {
	qq := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Normalize()
	return &MQuaternion{qq.V[0], qq.V[1], qq.V[2], qq.W}
}

func (quat *MQuaternion) Negate() *MQuaternion {
	quat.X *= -1
	quat.Y *= -1
	quat.Z *= -1
	quat.W *= -1
	return quat
}

func (quat *MQuaternion) Negated() *MQuaternion {
	return NewMQuaternionByValues(-quat.X, -quat.Y, -quat.Z, -quat.W)
}

func (quat *MQuaternion) Inverse() *MQuaternion {
	lenSq := quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W

	if lenSq < 1e-10 {
		*quat = MQuaternion{0, 0, 0, 1}
		return quat
	}

	if math.Abs(lenSq-1.0) < 1e-10 {
		quat.X = -quat.X
		quat.Y = -quat.Y
		quat.Z = -quat.Z
		return quat
	}

	invLenSq := 1.0 / lenSq
	quat.X = -quat.X * invLenSq
	quat.Y = -quat.Y * invLenSq
	quat.Z = -quat.Z * invLenSq
	quat.W = quat.W * invLenSq
	return quat
}

func (quat *MQuaternion) Inverted() *MQuaternion {
	lenSq := quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W

	if lenSq < 1e-10 {
		return &MQuaternion{0, 0, 0, 1}
	}

	if math.Abs(lenSq-1.0) < 1e-10 {
		return &MQuaternion{-quat.X, -quat.Y, -quat.Z, quat.W}
	}

	invLenSq := 1.0 / lenSq
	return &MQuaternion{
		-quat.X * invLenSq,
		-quat.Y * invLenSq,
		-quat.Z * invLenSq,
		quat.W * invLenSq,
	}
}

func (quat *MQuaternion) SetShortestRotation(other *MQuaternion) *MQuaternion {
	if !quat.IsShortestRotation(other) {
		quat.Negate()
	}
	return quat
}

func (quat *MQuaternion) IsShortestRotation(other *MQuaternion) bool {
	return quat.Dot(other) >= 0
}

func (quat *MQuaternion) IsUnitQuat(tolerance float64) bool {
	norm := quat.Norm()
	return norm >= (1.0-tolerance) && norm <= (1.0+tolerance)
}

func (quat *MQuaternion) Shorten() *MQuaternion {
	if quat.W < 0 {
		quat.Negate()
	}
	return quat
}

func (quat *MQuaternion) Dot(other *MQuaternion) float64 {
	q1 := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}
	q2 := mgl64.Quat{V: mgl64.Vec3{other.X, other.Y, other.Z}, W: other.W}
	return q1.Dot(q2)
}

func (quat *MQuaternion) MuledScalar(factor float64) *MQuaternion {
	if factor == 0.0 {
		return NewMQuaternion()
	} else if factor == 1.0 {
		return quat.Copy()
	} else if factor == -1.0 {
		return quat.Inverted()
	}

	return MQuaternionIdent.SlerpExtended(quat, factor)
}

func (quat *MQuaternion) SlerpExtended(other *MQuaternion, t float64) *MQuaternion {
	if quat.NearEquals(other, 1e-8) {
		return quat.Copy()
	}

	cosOmega := quat.X*other.X + quat.Y*other.Y + quat.Z*other.Z + quat.W*other.W

	q2x, q2y, q2z, q2w := other.X, other.Y, other.Z, other.W
	if cosOmega < 0.0 {
		cosOmega = -cosOmega
		q2x = -q2x
		q2y = -q2y
		q2z = -q2z
		q2w = -q2w
	}

	var result *MQuaternion
	if cosOmega > 0.9999 {
		result = &MQuaternion{
			quat.X*(1-t) + q2x*t,
			quat.Y*(1-t) + q2y*t,
			quat.Z*(1-t) + q2z*t,
			quat.W*(1-t) + q2w*t,
		}
	} else {
		omega := math.Acos(cosOmega)
		sinOmega := math.Sin(omega)

		angle := t * omega

		s1 := math.Sin(omega-angle) / sinOmega
		s2 := math.Sin(angle) / sinOmega

		result = &MQuaternion{
			s1*quat.X + s2*q2x,
			s1*quat.Y + s2*q2y,
			s1*quat.Z + s2*q2z,
			s1*quat.W + s2*q2w,
		}
	}

	return result.Normalize()
}

func (quat *MQuaternion) ToAxisAngle() (*MVec3, float64) {
	lenSq := quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W

	var normW float64
	if math.Abs(lenSq-1.0) > 1e-10 {
		invLen := 1.0 / math.Sqrt(lenSq)
		normW = quat.W * invLen
	} else {
		normW = quat.W
	}

	angle := 2.0 * math.Acos(math.Max(-1.0, math.Min(1.0, normW)))

	s := math.Sqrt(1.0 - normW*normW)

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

func (quat *MQuaternion) Slerp(other *MQuaternion, t float64) *MQuaternion {
	if t <= 0.0 {
		return quat.Copy()
	}
	if t >= 1.0 {
		return other.Copy()
	}
	if quat.NearEquals(other, 1e-8) {
		return quat.Copy()
	}

	cosOmega := quat.X*other.X + quat.Y*other.Y + quat.Z*other.Z + quat.W*other.W

	q2x, q2y, q2z, q2w := other.X, other.Y, other.Z, other.W
	if cosOmega < 0.0 {
		cosOmega = -cosOmega
		q2x = -q2x
		q2y = -q2y
		q2z = -q2z
		q2w = -q2w
	}

	var k1, k2 float64
	if cosOmega > 0.9999 {
		k1 = 1.0 - t
		k2 = t
	} else {
		sinOmega := math.Sqrt(1.0 - cosOmega*cosOmega)
		omega := math.Atan2(sinOmega, cosOmega)
		invSinOmega := 1.0 / sinOmega

		k1 = math.Sin((1.0-t)*omega) * invSinOmega
		k2 = math.Sin(t*omega) * invSinOmega
	}

	return &MQuaternion{
		k1*quat.X + k2*q2x,
		k1*quat.Y + k2*q2y,
		k1*quat.Z + k2*q2z,
		k1*quat.W + k2*q2w,
	}
}

func (quat *MQuaternion) Lerp(other *MQuaternion, t float64) *MQuaternion {
	if t <= 0.0 {
		return quat.Copy()
	}
	if t >= 1.0 {
		return other.Copy()
	}
	if quat.NearEquals(other, 1e-8) {
		return quat.Copy()
	}

	scale0 := 1.0 - t
	scale1 := t

	dot := quat.X*other.X + quat.Y*other.Y + quat.Z*other.Z + quat.W*other.W
	if dot < 0 {
		scale1 = -scale1
	}

	x := scale0*quat.X + scale1*other.X
	y := scale0*quat.Y + scale1*other.Y
	z := scale0*quat.Z + scale1*other.Z
	w := scale0*quat.W + scale1*other.W

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

func (quat *MQuaternion) ToDegree() float64 {
	return RadToDeg(quat.ToRadian())
}

func (quat *MQuaternion) ToRadian() float64 {
	return 2 * math.Acos(math.Min(1, math.Max(-1, quat.Normalize().W)))
}

func (quat *MQuaternion) ToSignedDegree() float64 {
	basicAngle := quat.ToDegree()

	if quat.Vec3().Length() > 0 {
		if quat.W >= 0 {
			return basicAngle
		} else {
			return -basicAngle
		}
	}

	return basicAngle
}

func (quat *MQuaternion) ToSignedRadian() float64 {
	basicAngle := quat.ToRadian()

	if quat.Vec3().Length() > 0 {
		if quat.W >= 0 {
			return basicAngle
		} else {
			return -basicAngle
		}
	}

	return basicAngle
}

func (quat *MQuaternion) ToTheta(v *MQuaternion) float64 {
	return math.Acos(math.Min(1, math.Max(-1, quat.Normalize().Dot(v.Normalize()))))
}

func NewMQuaternionFromDirection(direction *MVec3, up *MVec3) *MQuaternion {
	if direction.Length() == 0 {
		return NewMQuaternion()
	}

	zAxis := direction.Normalized()
	xAxis := up.Cross(zAxis).Normalized()

	if xAxis.LengthSqr() == 0 {
		return NewMQuaternionRotate(&MVec3{0.0, 0.0, 1.0}, zAxis)
	}

	yAxis := zAxis.Cross(xAxis)

	return NewMQuaternionFromAxes(xAxis, yAxis, zAxis).Normalize()
}

func NewMQuaternionRotate(fromV, toV *MVec3) *MQuaternion {
	if fromV.NearEquals(toV, 1e-6) || fromV.Length() == 0 || toV.Length() == 0 {
		return NewMQuaternion()
	}
	v := mgl64.QuatBetweenVectors(mgl64.Vec3{fromV.X, fromV.Y, fromV.Z}, mgl64.Vec3{toV.X, toV.Y, toV.Z})
	return NewMQuaternionByValues(v.V[0], v.V[1], v.V[2], v.W)
}

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

func (quat *MQuaternion) SeparateTwistByAxis(globalAxis *MVec3) (twistQQ *MQuaternion, yzQQ *MQuaternion) {
	globalXAxis := globalAxis.Normalized()

	globalXVec := quat.MulVec3(globalXAxis)
	yzQQ = NewMQuaternionRotate(globalXAxis, globalXVec.Normalize())
	twistQQ = yzQQ.Inverted().Mul(quat)

	return twistQQ, yzQQ
}

func (quat *MQuaternion) SeparateByAxis(globalAxis *MVec3) (xQQ, yQQ, zQQ *MQuaternion) {
	localZAxis := MVec3{0, 0, -1}
	globalXAxis := globalAxis.Normalize()
	globalYAxis := localZAxis.Cross(globalXAxis)
	globalZAxis := globalXAxis.Cross(globalYAxis)

	if globalYAxis.Length() == 0 {
		localYAxis := MVec3UnitY
		globalZAxis = localYAxis.Cross(globalXAxis)
		globalYAxis = globalXAxis.Cross(globalZAxis)
	}

	globalXVec := quat.MulVec3(globalXAxis)
	yzQQ := NewMQuaternionRotate(globalXAxis, globalXVec.Normalize())
	xQQ = yzQQ.Inverse().Mul(quat)

	globalYVec := quat.MulVec3(globalYAxis)
	xzQQ := NewMQuaternionRotate(globalYAxis, globalYVec.Normalize())
	yQQ = xzQQ.Inverse().Mul(quat)

	globalZVec := quat.MulVec3(globalZAxis)
	xyQQ := NewMQuaternionRotate(globalZAxis, globalZVec.Normalize())
	zQQ = xyQQ.Inverse().Mul(quat)

	return xQQ, yQQ, zQQ
}

func (quat *MQuaternion) Copy() *MQuaternion {
	return NewMQuaternionByValues(quat.X, quat.Y, quat.Z, quat.W)
}

func (quat *MQuaternion) Vector() []float64 {
	return []float64{quat.X, quat.Y, quat.Z, quat.W}
}

func (quat *MQuaternion) ToMat4() *MMat4 {
	m := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}.Mat4()
	return (*MMat4)(&m)
}

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

func (quat *MQuaternion) NearEquals(other *MQuaternion, epsilon float64) bool {
	q1 := mgl64.Quat{V: mgl64.Vec3{quat.X, quat.Y, quat.Z}, W: quat.W}
	q2 := mgl64.Quat{V: mgl64.Vec3{other.X, other.Y, other.Z}, W: other.W}
	return q1.ApproxEqualThreshold(q2, epsilon)
}

func (quat *MQuaternion) MulVec3(v *MVec3) *MVec3 {

	qx, qy, qz, qw := quat.X, quat.Y, quat.Z, quat.W

	vx, vy, vz := v.X, v.Y, v.Z

	twoQx, twoQy, twoQz := 2.0*qx, 2.0*qy, 2.0*qz

	xx, xy, xz := qx*twoQx, qx*twoQy, qx*twoQz
	yy, yz, zz := qy*twoQy, qy*twoQz, qz*twoQz
	wx, wy, wz := qw*twoQx, qw*twoQy, qw*twoQz

	x := vx*(1.0-(yy+zz)) + vy*(xy-wz) + vz*(xz+wy)
	y := vx*(xy+wz) + vy*(1.0-(xx+zz)) + vz*(yz-wx)
	z := vx*(xz-wy) + vy*(yz+wx) + vz*(1.0-(xx+yy))

	return &MVec3{x, y, z}
}

func VectorToDegree(a *MVec3, b *MVec3) float64 {
	return RadToDeg(VectorToRadian(a, b))
}

func VectorToRadian(a *MVec3, b *MVec3) float64 {
	p := a.Dot(b)
	normA := a.Length()
	normB := b.Length()

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

