package mmath

import (
	"errors"
	"fmt"
	"math"

	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

type Quaternion struct {
	quat.Number
}

var (
	quaternionZero  = Quaternion{}
	quaternionIdent = Quaternion{quat.Number{Real: 1}}
)

const (
	GIMBAL1_RAD = 88.0 / 180.0 * math.Pi
	GIMBAL2_RAD = GIMBAL1_RAD * 2
	ONE_RAD     = math.Pi
	HALF_RAD    = math.Pi / 2
)

func NewQuaternion() Quaternion {
	return quaternionIdent
}

func NewQuaternionByValues(x, y, z, w float64) Quaternion {
	return Quaternion{quat.Number{Real: w, Imag: x, Jmag: y, Kmag: z}}
}

func NewQuaternionByValuesShort(x, y, z, w float64) Quaternion {
	q := NewQuaternionByValues(x, y, z, w)
	if q.W() < 0 {
		return q.Negated()
	}
	return q
}

func (q Quaternion) X() float64 { return q.Imag }
func (q Quaternion) Y() float64 { return q.Jmag }
func (q Quaternion) Z() float64 { return q.Kmag }
func (q Quaternion) W() float64 { return q.Real }

func (q *Quaternion) set(x, y, z, w float64) {
	q.Imag = x
	q.Jmag = y
	q.Kmag = z
	q.Real = w
}

func (q Quaternion) XYZ() Vec3 {
	return Vec3{r3.Vec{X: q.Imag, Y: q.Jmag, Z: q.Kmag}}
}

func (q *Quaternion) SetXYZ(v Vec3) {
	q.Imag = v.X
	q.Jmag = v.Y
	q.Kmag = v.Z
}

func (q Quaternion) String() string {
	return fmt.Sprintf("[x=%.7f, y=%.7f, z=%.7f, w=%.7f]", q.Imag, q.Jmag, q.Kmag, q.Real)
}

func (q Quaternion) MMD() Quaternion {
	return q
}

func NewQuaternionFromAxisAngles(axis Vec3, angle float64) Quaternion {
	rot := r3.NewRotation(angle, axis.Vec)
	qq := quat.Number(rot)
	return Quaternion{qq}
}

func NewQuaternionFromAxisAnglesRotate(axis Vec3, angle float64) Quaternion {
	rot := r3.NewRotation(angle, axis.Vec)
	qq := quat.Number(rot)
	return Quaternion{qq}
}

func NewQuaternionFromRadians(xPitch, yHead, zRoll float64) Quaternion {
	cx, sx := math.Cos(xPitch*0.5), math.Sin(xPitch*0.5)
	cy, sy := math.Cos(yHead*0.5), math.Sin(yHead*0.5)
	cz, sz := math.Cos(zRoll*0.5), math.Sin(zRoll*0.5)

	w := cx*cy*cz - sx*sy*sz
	x := sx*cy*cz + cx*sy*sz
	y := cx*sy*cz - sx*cy*sz
	z := cx*cy*sz + sx*sy*cz

	return NewQuaternionByValues(x, y, z, w).Normalized()
}

func (q Quaternion) ToRadians() Vec3 {
	sx := -(2*q.Jmag*q.Kmag - 2*q.Imag*q.Real)
	unlocked := math.Abs(sx) < 0.99999
	xPitch := math.Asin(math.Max(-1, math.Min(1, sx)))
	var yHead, zRoll float64
	if unlocked {
		yHead = math.Atan2(2*q.Imag*q.Kmag+2*q.Jmag*q.Real, 2*q.Real*q.Real+2*q.Kmag*q.Kmag-1)
		zRoll = math.Atan2(2*q.Imag*q.Jmag+2*q.Kmag*q.Real, 2*q.Real*q.Real+2*q.Jmag*q.Jmag-1)
	} else {
		yHead = math.Atan2(-(2*q.Imag*q.Kmag-2*q.Jmag*q.Real), 2*q.Real*q.Real+2*q.Imag*q.Imag-1)
		zRoll = 0
	}
	return Vec3{r3.Vec{X: xPitch, Y: yHead, Z: zRoll}}
}

func (q Quaternion) ToRadiansWithGimbal(axisIndex int) (Vec3, bool) {
	r := q.ToRadians()
	var other1Rad, other2Rad float64
	switch axisIndex {
	case 0:
		other1Rad = math.Abs(r.Y)
		other2Rad = math.Abs(r.Z)
	case 1:
		other1Rad = math.Abs(r.X)
		other2Rad = math.Abs(r.Z)
	default:
		other1Rad = math.Abs(r.X)
		other2Rad = math.Abs(r.Y)
	}
	if other1Rad >= GIMBAL2_RAD && other2Rad >= GIMBAL2_RAD {
		return r, true
	}
	return r, false
}

func NewQuaternionFromDegrees(xPitch, yHead, zRoll float64) Quaternion {
	return NewQuaternionFromRadians(DegToRad(xPitch), DegToRad(yHead), DegToRad(zRoll))
}

func (q Quaternion) ToDegrees() Vec3 {
	vec := q.ToRadians()
	return Vec3{r3.Vec{X: RadToDeg(vec.X), Y: RadToDeg(vec.Y), Z: RadToDeg(vec.Z)}}
}

func (q Quaternion) ToMMDDegrees() Vec3 {
	vec := q.MMD().ToRadians()
	return Vec3{r3.Vec{X: RadToDeg(vec.X), Y: RadToDeg(-vec.Y), Z: RadToDeg(-vec.Z)}}
}

func (q Quaternion) Vec4() Vec4 {
	return Vec4{q.Imag, q.Jmag, q.Kmag, q.Real}
}

func (q Quaternion) Vec3() Vec3 {
	return Vec3{r3.Vec{X: q.Imag, Y: q.Jmag, Z: q.Kmag}}
}

func (q Quaternion) MulShort(other Quaternion) Quaternion {
	x := q.Real*other.Imag + q.Imag*other.Real + q.Jmag*other.Kmag - q.Kmag*other.Jmag
	y := q.Real*other.Jmag + q.Jmag*other.Real + q.Kmag*other.Imag - q.Imag*other.Kmag
	z := q.Real*other.Kmag + q.Kmag*other.Real + q.Imag*other.Jmag - q.Jmag*other.Imag
	w := q.Real*other.Real - q.Imag*other.Imag - q.Jmag*other.Jmag - q.Kmag*other.Kmag
	return NewQuaternionByValues(x, y, z, w)
}

func (q Quaternion) MuledShort(other Quaternion) Quaternion {
	copied := q
	copied.Mul(other)
	return copied
}

func (q *Quaternion) Mul(other Quaternion) *Quaternion {
	qq := quat.Mul(q.Number, other.Number)
	q.Number = qq
	return q
}

func (q Quaternion) Muled(other Quaternion) Quaternion {
	return Quaternion{quat.Mul(q.Number, other.Number)}
}

func (q Quaternion) Norm() float64 {
	return quat.Abs(q.Number)
}

func (q Quaternion) Length() float64 {
	return quat.Abs(q.Number)
}

func (q *Quaternion) Normalize() *Quaternion {
	lenSq := q.Imag*q.Imag + q.Jmag*q.Jmag + q.Kmag*q.Kmag + q.Real*q.Real
	if lenSq < 1e-10 {
		q.Number = quaternionIdent.Number
		return q
	}
	invLen := 1.0 / math.Sqrt(lenSq)
	q.Imag *= invLen
	q.Jmag *= invLen
	q.Kmag *= invLen
	q.Real *= invLen
	return q
}

func (q Quaternion) Normalized() Quaternion {
	vec := q
	vec.Normalize()
	return vec
}

func (q *Quaternion) Negate() *Quaternion {
	q.Imag = -q.Imag
	q.Jmag = -q.Jmag
	q.Kmag = -q.Kmag
	q.Real = -q.Real
	return q
}

func (q Quaternion) Negated() Quaternion {
	return Quaternion{quat.Scale(-1, q.Number)}
}

func (q *Quaternion) Inverse() *Quaternion {
	lenSq := q.Imag*q.Imag + q.Jmag*q.Jmag + q.Kmag*q.Kmag + q.Real*q.Real
	if lenSq < 1e-10 {
		q.Number = quaternionIdent.Number
		return q
	}
	q.Number = quat.Inv(q.Number)
	return q
}

func (q Quaternion) Inverted() Quaternion {
	lenSq := q.Imag*q.Imag + q.Jmag*q.Jmag + q.Kmag*q.Kmag + q.Real*q.Real
	if lenSq < 1e-10 {
		return quaternionIdent
	}
	return Quaternion{quat.Inv(q.Number)}
}

func (q *Quaternion) SetShortestRotation(other Quaternion) *Quaternion {
	if !q.IsShortestRotation(other) {
		q.Negate()
	}
	return q
}

func (q Quaternion) IsShortestRotation(other Quaternion) bool {
	return q.Dot(other) >= 0
}

func (q Quaternion) IsUnitQuat(tolerance float64) bool {
	norm := q.Norm()
	return norm >= (1.0-tolerance) && norm <= (1.0+tolerance)
}

func (q Quaternion) Dot(other Quaternion) float64 {
	return q.Imag*other.Imag + q.Jmag*other.Jmag + q.Kmag*other.Kmag + q.Real*other.Real
}

func (q Quaternion) MuledScalar(factor float64) Quaternion {
	if factor == 0.0 {
		return NewQuaternion()
	}
	if factor == 1.0 {
		return q
	}
	if factor == -1.0 {
		return q.Inverted()
	}
	return quaternionIdent.SlerpExtended(q, factor)
}

func (q Quaternion) SlerpExtended(other Quaternion, t float64) Quaternion {
	if q.NearEquals(other, 1e-8) {
		return q
	}

	cosOmega := q.Dot(other)
	q2x, q2y, q2z, q2w := other.Imag, other.Jmag, other.Kmag, other.Real
	if cosOmega < 0 {
		cosOmega = -cosOmega
		q2x = -q2x
		q2y = -q2y
		q2z = -q2z
		q2w = -q2w
	}

	var result Quaternion
	if cosOmega > 0.9999 {
		result = NewQuaternionByValues(
			q.Imag*(1-t)+q2x*t,
			q.Jmag*(1-t)+q2y*t,
			q.Kmag*(1-t)+q2z*t,
			q.Real*(1-t)+q2w*t,
		)
	} else {
		omega := math.Acos(cosOmega)
		sinOmega := math.Sin(omega)
		angle := t * omega
		s1 := math.Sin(omega-angle) / sinOmega
		s2 := math.Sin(angle) / sinOmega
		result = NewQuaternionByValues(
			s1*q.Imag+s2*q2x,
			s1*q.Jmag+s2*q2y,
			s1*q.Kmag+s2*q2z,
			s1*q.Real+s2*q2w,
		)
	}
	return result.Normalized()
}

func (q Quaternion) ToAxisAngle() (Vec3, float64) {
	lenSq := q.Imag*q.Imag + q.Jmag*q.Jmag + q.Kmag*q.Kmag + q.Real*q.Real
	var normW float64
	if math.Abs(lenSq-1.0) > 1e-10 {
		invLen := 1.0 / math.Sqrt(lenSq)
		normW = q.Real * invLen
	} else {
		normW = q.Real
	}

	angle := 2.0 * math.Acos(math.Max(-1.0, math.Min(1.0, normW)))
	s := math.Sqrt(1.0 - normW*normW)
	if s < 1e-9 {
		return Vec3{r3.Vec{X: 1}}, angle
	}
	invS := 1.0 / s
	return Vec3{r3.Vec{X: q.Imag * invS, Y: q.Jmag * invS, Z: q.Kmag * invS}}, angle
}

func (q Quaternion) Slerp(other Quaternion, t float64) Quaternion {
	if t <= 0 {
		return q
	}
	if t >= 1 {
		return other
	}
	if q.NearEquals(other, 1e-8) {
		return q
	}

	cosOmega := q.Dot(other)
	q2x, q2y, q2z, q2w := other.Imag, other.Jmag, other.Kmag, other.Real
	if cosOmega < 0 {
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

	return NewQuaternionByValues(
		k1*q.Imag+k2*q2x,
		k1*q.Jmag+k2*q2y,
		k1*q.Kmag+k2*q2z,
		k1*q.Real+k2*q2w,
	)
}

func (q Quaternion) Lerp(other Quaternion, t float64) Quaternion {
	if t <= 0 {
		return q
	}
	if t >= 1 {
		return other
	}
	if q.NearEquals(other, 1e-8) {
		return q
	}

	scale0 := 1.0 - t
	scale1 := t
	dot := q.Dot(other)
	if dot < 0 {
		scale1 = -scale1
	}

	x := scale0*q.Imag + scale1*other.Imag
	y := scale0*q.Jmag + scale1*other.Jmag
	z := scale0*q.Kmag + scale1*other.Kmag
	w := scale0*q.Real + scale1*other.Real

	len := math.Sqrt(x*x + y*y + z*z + w*w)
	if len > 0 {
		invLen := 1.0 / len
		x *= invLen
		y *= invLen
		z *= invLen
		w *= invLen
	}

	return NewQuaternionByValues(x, y, z, w)
}

func (q *Quaternion) ToDegree() float64 {
	return RadToDeg(q.ToRadian())
}

func (q *Quaternion) ToRadian() float64 {
	q.Normalize()
	return 2 * math.Acos(math.Min(1, math.Max(-1, q.Real)))
}

func (q *Quaternion) ToSignedDegree() float64 {
	basicAngle := q.ToDegree()
	if q.Vec3().Length() > 0 {
		if q.Real >= 0 {
			return basicAngle
		}
		return -basicAngle
	}
	return basicAngle
}

func (q *Quaternion) ToSignedRadian() float64 {
	basicAngle := q.ToRadian()
	if q.Vec3().Length() > 0 {
		if q.Real >= 0 {
			return basicAngle
		}
		return -basicAngle
	}
	return basicAngle
}

func (q *Quaternion) ToTheta(other Quaternion) float64 {
	q.Normalize()
	v := other.Normalized()
	return math.Acos(math.Min(1, math.Max(-1, q.Dot(v))))
}

func NewQuaternionFromDirection(direction, up Vec3) Quaternion {
	if direction.Length() == 0 {
		return NewQuaternion()
	}

	zAxis := direction.Normalized()
	xAxis := up.Cross(zAxis).Normalized()
	if xAxis.LengthSqr() == 0 {
		return NewQuaternionRotate(Vec3{r3.Vec{Z: 1}}, zAxis)
	}
	yAxis := zAxis.Cross(xAxis)
	return NewQuaternionFromAxes(xAxis, yAxis, zAxis).Normalized()
}

func NewQuaternionRotate(fromV, toV Vec3) Quaternion {
	if fromV.NearEquals(toV, 1e-6) || fromV.Length() == 0 || toV.Length() == 0 {
		return NewQuaternion()
	}

	v0 := fromV.Normalized()
	v1 := toV.Normalized()
	dot := v0.Dot(v1)
	if dot >= 1.0 {
		return NewQuaternion()
	}
	if dot <= -1.0 {
		axis := Vec3{r3.Vec{X: 1}}
		if math.Abs(v0.X) > 0.9 {
			axis = Vec3{r3.Vec{Y: 1}}
		}
		axis = v0.Cross(axis).Normalized()
		return NewQuaternionFromAxisAngles(axis, math.Pi)
	}

	cross := v0.Cross(v1)
	s := math.Sqrt((1 + dot) * 2)
	oos := 1 / s
	return NewQuaternionByValues(cross.X*oos, cross.Y*oos, cross.Z*oos, s*0.5).Normalized()
}

func NewQuaternionFromAxes(xAxis, yAxis, zAxis Vec3) Quaternion {
	mat := NewMat4ByValues(
		xAxis.X, xAxis.Y, xAxis.Z, 0,
		yAxis.X, yAxis.Y, yAxis.Z, 0,
		zAxis.X, zAxis.Y, zAxis.Z, 0,
		0, 0, 0, 1,
	)
	return mat.Quaternion()
}

func (q Quaternion) SeparateTwistByAxis(globalAxis Vec3) (Quaternion, Quaternion) {
	globalXAxis := globalAxis.Normalized()
	globalXVec := q.MulVec3(globalXAxis)
	globalXVec.Normalize()
	yzQQ := NewQuaternionRotate(globalXAxis, globalXVec)
	twistQQ := yzQQ.Inverted().Muled(q)
	return twistQQ, yzQQ
}

func (q Quaternion) SeparateByAxis(globalAxis Vec3) (Quaternion, Quaternion, Quaternion) {
	localZAxis := Vec3{r3.Vec{Z: -1}}
	globalXAxis := globalAxis.Normalized()
	globalYAxis := localZAxis.Cross(globalXAxis)
	globalZAxis := globalXAxis.Cross(globalYAxis)

	if globalYAxis.Length() == 0 {
		localYAxis := UNIT_Y_VEC3
		globalZAxis = localYAxis.Cross(globalXAxis)
		globalYAxis = globalXAxis.Cross(globalZAxis)
	}

	globalXVec := q.MulVec3(globalXAxis)
	globalXVec.Normalize()
	yzQQ := NewQuaternionRotate(globalXAxis, globalXVec)
	xQQ := yzQQ.Inverted().Muled(q)

	globalYVec := q.MulVec3(globalYAxis)
	globalYVec.Normalize()
	xzQQ := NewQuaternionRotate(globalYAxis, globalYVec)
	yQQ := xzQQ.Inverted().Muled(q)

	globalZVec := q.MulVec3(globalZAxis)
	globalZVec.Normalize()
	xyQQ := NewQuaternionRotate(globalZAxis, globalZVec)
	zQQ := xyQQ.Inverted().Muled(q)

	return xQQ, yQQ, zQQ
}

func (q Quaternion) Copy() (*Quaternion, error) {
	return deepCopy(q)
}

func (q Quaternion) Vector() []float64 {
	return []float64{q.Imag, q.Jmag, q.Kmag, q.Real}
}

func (q Quaternion) ToMat4() Mat4 {
	x, y, z, w := q.Imag, q.Jmag, q.Kmag, q.Real
	xx, yy, zz := x*x, y*y, z*z
	xy, xz, yz := x*y, x*z, y*z
	wx, wy, wz := w*x, w*y, w*z

	m00 := 1 - 2*(yy+zz)
	m01 := 2 * (xy - wz)
	m02 := 2 * (xz + wy)

	m10 := 2 * (xy + wz)
	m11 := 1 - 2*(xx+zz)
	m12 := 2 * (yz - wx)

	m20 := 2 * (xz - wy)
	m21 := 2 * (yz + wx)
	m22 := 1 - 2*(xx+yy)

	return NewMat4ByValues(
		m00, m10, m20, 0,
		m01, m11, m21, 0,
		m02, m12, m22, 0,
		0, 0, 0, 1,
	)
}

func (q Quaternion) ToFixedAxisRotation(fixedAxis Vec3) Quaternion {
	normalizedFixedAxis := fixedAxis.Normalized()
	quatAxis := q.XYZ().Normalized()
	rad := q.ToRadian()
	if normalizedFixedAxis.Dot(quatAxis) < 0 {
		rad *= -1
	}
	return NewQuaternionFromAxisAngles(normalizedFixedAxis, rad)
}

func (q Quaternion) IsIdent() bool {
	return q.NearEquals(quaternionIdent, 1e-6)
}

func (q Quaternion) NearEquals(other Quaternion, epsilon float64) bool {
	return math.Abs(q.Imag-other.Imag) <= epsilon && math.Abs(q.Jmag-other.Jmag) <= epsilon &&
		math.Abs(q.Kmag-other.Kmag) <= epsilon && math.Abs(q.Real-other.Real) <= epsilon
}

func (q Quaternion) MulVec3(v Vec3) Vec3 {
	rot := r3.Rotation(q.Normalized().Number)
	return Vec3{rot.Rotate(v.Vec)}
}

func VectorToDegree(a, b Vec3) float64 {
	return RadToDeg(VectorToRadian(a, b))
}

func VectorToRadian(a, b Vec3) float64 {
	p := a.Dot(b)
	normA := a.Length()
	normB := b.Length()
	cosAngle := p / (normA * normB)
	return math.Acos(math.Min(1, math.Max(-1, cosAngle)))
}

func (q Quaternion) Log() (Quaternion, error) {
	if math.Abs(q.Real) > 1.0 {
		return Quaternion{}, errors.New("invalid quaternion scalar part")
	}
	vNorm := q.Norm()
	if vNorm == 0 {
		return Quaternion{quat.Number{Real: 1}}, nil
	}
	angle := math.Acos(q.Real)
	scale := angle / vNorm
	return NewQuaternionByValues(scale*q.Imag, scale*q.Jmag, scale*q.Kmag, 0), nil
}

func FindSlerpT(q1, q2, qt Quaternion, initialT float64) float64 {
	tol := 1e-10
	phi := (1 + math.Sqrt(5)) / 2
	maxIterations := 100

	if math.Abs(q1.Dot(q2)) > 0.9999 {
		return initialT
	}

	a := 0.0
	b := 1.0
	c := b - (b-a)/phi
	d := a + (b-a)/phi

	q2c := q2
	if q1.Dot(q2) < 0 {
		q2c = q2.Negated()
	}

	errorFunc := func(t float64) float64 {
		tQuat := q1.Slerp(q2c, t)
		theta := math.Acos(tQuat.Dot(qt))
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

func FindLerpT(q1, q2, qt Quaternion) float64 {
	tol := 1e-8
	phi := (1 + math.Sqrt(5)) / 2
	maxIterations := 100

	a := 0.0
	b := 1.0
	c := b - (b-a)/phi
	d := a + (b-a)/phi

	q2c := q2
	if q1.Dot(q2) < 0 {
		q2c = q2.Negated()
	}

	errorFunc := func(t float64) float64 {
		tQuat := q1.Lerp(q2c, t)
		theta := math.Acos(tQuat.Dot(qt))
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
