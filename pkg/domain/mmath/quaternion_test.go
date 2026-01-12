// 指示: miu200521358
package mmath

import (
	"math"
	"testing"
)

func TestQuaternionBasics(t *testing.T) {
	q := NewQuaternion()
	if !q.IsIdent() {
		t.Errorf("NewQuaternion")
	}
	qv := NewQuaternionByValues(1, 2, 3, 4)
	if qv.X() != 1 || qv.Y() != 2 || qv.Z() != 3 || qv.W() != 4 {
		t.Errorf("NewQuaternionByValues")
	}
	qs := NewQuaternionByValuesShort(1, 2, 3, -1)
	if qs.W() <= 0 {
		t.Errorf("NewQuaternionByValuesShort")
	}
	vec := qv.XYZ()
	if vec != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("XYZ")
	}
	qv.SetXYZ(Vec3{r3Vec(4, 5, 6)})
	if qv.X() != 4 || qv.Y() != 5 || qv.Z() != 6 {
		t.Errorf("SetXYZ")
	}
	if qv.String() == "" || qv.MMD() != qv {
		t.Errorf("String/MMD")
	}
}

func TestQuaternionOps(t *testing.T) {
	q := NewQuaternionByValues(1, 0, 0, 0)
	q2 := NewQuaternionByValues(0, 1, 0, 0)
	if q.MulShort(q2) == quaternionZero {
		t.Errorf("MulShort")
	}
	if q.MuledShort(q2) == quaternionZero {
		t.Errorf("MuledShort")
	}
	q3 := q
	q3.Mul(q2)
	if q3 == quaternionZero {
		t.Errorf("Mul")
	}
	if q.Muled(q2) == quaternionZero {
		t.Errorf("Muled")
	}
	if q.Norm() == 0 || q.Length() == 0 {
		t.Errorf("Norm/Length")
	}
	qn := Quaternion{}
	qn.Normalize()
	if !qn.IsIdent() {
		t.Errorf("Normalize")
	}
	if qn.Normalized().IsIdent() == false {
		t.Errorf("Normalized")
	}
	q4 := NewQuaternionByValues(1, 2, 3, 4)
	q4.Negate()
	if q4.W() != -4 {
		t.Errorf("Negate")
	}
	if q4.Negated().W() != 4 {
		t.Errorf("Negated")
	}
	q5 := Quaternion{}
	q5.Inverse()
	if !q5.IsIdent() {
		t.Errorf("Inverse")
	}
	if NewQuaternionByValues(1, 2, 3, 4).Inverted() == quaternionZero {
		t.Errorf("Inverted")
	}
	q6 := NewQuaternionByValues(1, 0, 0, 0)
	q6.SetShortestRotation(NewQuaternionByValues(-1, 0, 0, 0))
	if !q6.IsShortestRotation(NewQuaternionByValues(-1, 0, 0, 0)) {
		t.Errorf("SetShortestRotation")
	}
	if !NewQuaternion().IsUnitQuat(1e-6) {
		t.Errorf("IsUnitQuat")
	}
	if q.Dot(q2) != 0 {
		t.Errorf("Dot")
	}
}

func TestQuaternionInterpolation(t *testing.T) {
	q := NewQuaternion()
	if q.MuledScalar(0) != NewQuaternion() {
		t.Errorf("MuledScalar zero")
	}
	if q.MuledScalar(1) != q {
		t.Errorf("MuledScalar one")
	}
	if !q.MuledScalar(-1).IsIdent() {
		t.Errorf("MuledScalar minus")
	}
	if q.MuledScalar(0.5) == quaternionZero {
		t.Errorf("MuledScalar mid")
	}
	q2 := NewQuaternionByValues(0, 0, 0, -1)
	qNear := NewQuaternionFromAxisAngles(Vec3{r3Vec(1, 0, 0)}, 1e-6)
	if q.SlerpExtended(qNear, 0.3) == quaternionZero {
		t.Errorf("SlerpExtended near")
	}
	if q.SlerpExtended(q2, 0.3) == quaternionZero {
		t.Errorf("SlerpExtended")
	}
	if q.Slerp(q, 0.5) != q {
		t.Errorf("Slerp same")
	}
	if q.Slerp(q2, 0) != q || q.Slerp(q2, 1) != q2 {
		t.Errorf("Slerp bounds")
	}
	if q.Slerp(q2, 0.5) == quaternionZero {
		t.Errorf("Slerp mid")
	}
	if q.Slerp(qNear, 0.5) == quaternionZero {
		t.Errorf("Slerp near")
	}
	if q.Lerp(q, 0.5) != q {
		t.Errorf("Lerp same")
	}
	if q.Lerp(q2, 0) != q || q.Lerp(q2, 1) != q2 {
		t.Errorf("Lerp bounds")
	}
	if q.Lerp(q2, 0.5) == quaternionZero {
		t.Errorf("Lerp")
	}
}

func TestQuaternionConversions(t *testing.T) {
	axis := Vec3{r3Vec(1, 0, 0)}
	qa := NewQuaternionFromAxisAngles(axis, math.Pi/2)
	qr := NewQuaternionFromAxisAnglesRotate(axis, math.Pi/2)
	if qa == quaternionZero || qr == quaternionZero {
		t.Errorf("AxisAngles")
	}
	q := NewQuaternionFromRadians(0, 0, 0)
	if !q.IsIdent() {
		t.Errorf("FromRadians")
	}
	qd := NewQuaternionFromDegrees(0, 0, 0)
	if !qd.IsIdent() {
		t.Errorf("FromDegrees")
	}
	if q.ToRadians() != (Vec3{r3Vec(0, 0, 0)}) {
		t.Errorf("ToRadians")
	}
	if _, gimbal := q.ToRadiansWithGimbal(0); gimbal {
		t.Errorf("ToRadiansWithGimbal")
	}
	if q.ToDegrees() != (Vec3{r3Vec(0, 0, 0)}) {
		t.Errorf("ToDegrees")
	}
	if q.ToMMDDegrees() != (Vec3{r3Vec(0, 0, 0)}) {
		t.Errorf("ToMMDDegrees")
	}
	if q.Vec4() != (Vec4{0, 0, 0, 1}) {
		t.Errorf("Vec4")
	}
	if q.Vec3() != (Vec3{r3Vec(0, 0, 0)}) {
		t.Errorf("Vec3")
	}
	if _, angle := q.ToAxisAngle(); angle != 0 {
		t.Errorf("ToAxisAngle")
	}
	if q.ToMat4().IsIdent() == false {
		t.Errorf("ToMat4")
	}
	fixed := q.ToFixedAxisRotation(Vec3{r3Vec(1, 0, 0)})
	if fixed.IsIdent() == false {
		t.Errorf("ToFixedAxisRotation")
	}
	qdeg := NewQuaternionByValues(0, 0, 0, 1)
	if qdeg.ToDegree() != 0 || qdeg.ToRadian() != 0 {
		t.Errorf("ToDegree/ToRadian")
	}
	qs := NewQuaternionByValues(1, 0, 0, -1)
	if qs.ToSignedDegree() >= 0 || qs.ToSignedRadian() >= 0 {
		t.Errorf("ToSigned")
	}
	qt := NewQuaternionByValues(0, 0, 0, 1)
	if qt.ToTheta(NewQuaternion()) != 0 {
		t.Errorf("ToTheta")
	}
	gimbalQ := NewQuaternionFromDegrees(0, 180, 180)
	if _, gimbal := gimbalQ.ToRadiansWithGimbal(0); !gimbal {
		t.Errorf("ToRadiansWithGimbal true")
	}
	if _, angle := NewQuaternionByValues(1, 1, 1, 1).ToAxisAngle(); angle == 0 {
		t.Errorf("ToAxisAngle norm")
	}
}

func TestQuaternionAdvanced(t *testing.T) {
	if NewQuaternionFromDirection(Vec3{}, Vec3{}).IsIdent() == false {
		t.Errorf("FromDirection zero")
	}
	collinear := NewQuaternionFromDirection(Vec3{r3Vec(0, 0, 1)}, Vec3{r3Vec(0, 0, 1)})
	if collinear.IsIdent() == false {
		t.Errorf("FromDirection collinear")
	}
	rot := NewQuaternionRotate(Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(-1, 0, 0)})
	if rot.IsIdent() {
		t.Errorf("NewQuaternionRotate opposite")
	}
	if NewQuaternionRotate(Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(1, 0, 0)}).IsIdent() == false {
		t.Errorf("NewQuaternionRotate same")
	}
	if NewQuaternionRotate(Vec3{}, Vec3{r3Vec(1, 0, 0)}).IsIdent() == false {
		t.Errorf("NewQuaternionRotate zero")
	}
	axes := NewQuaternionFromAxes(Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(0, 1, 0)}, Vec3{r3Vec(0, 0, 1)})
	if axes.IsIdent() == false {
		t.Errorf("NewQuaternionFromAxes")
	}
	q := NewQuaternionFromDegrees(10, 20, 30)
	twist, yz := q.SeparateTwistByAxis(Vec3{r3Vec(1, 0, 0)})
	if twist == quaternionZero || yz == quaternionZero {
		t.Errorf("SeparateTwistByAxis")
	}
	xq, yq, zq := q.SeparateByAxis(Vec3{r3Vec(0, 0, -1)})
	if xq == quaternionZero || yq == quaternionZero || zq == quaternionZero {
		t.Errorf("SeparateByAxis")
	}
	if cp, err := q.Copy(); err != nil || cp == quaternionZero {
		t.Errorf("Copy")
	}
	if len(q.Vector()) != 4 {
		t.Errorf("Vector")
	}
	if NewQuaternion().MulVec3(Vec3{r3Vec(1, 2, 3)}) != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("MulVec3")
	}
	if VectorToDegree(Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(1, 0, 0)}) != 0 {
		t.Errorf("VectorToDegree")
	}
	if _, err := NewQuaternionByValues(0, 0, 0, 2).Log(); err == nil {
		t.Errorf("Log error")
	}
	qZero := Quaternion{}
	if logq, err := qZero.Log(); err != nil || logq.W() != 1 {
		t.Errorf("Log zero")
	}
	if _, err := NewQuaternionByValues(1, 0, 0, 1).Log(); err != nil {
		t.Errorf("Log")
	}
	st := FindSlerpT(NewQuaternion(), NewQuaternionByValues(0, 0, 0, -1), NewQuaternion(), 0.5)
	if math.IsNaN(st) || st < 0 || st > 1 {
		t.Errorf("FindSlerpT")
	}
	lt := FindLerpT(NewQuaternion(), NewQuaternionByValues(0, 0, 0, -1), NewQuaternion())
	if math.IsNaN(lt) || lt < 0 || lt > 1 {
		t.Errorf("FindLerpT")
	}
}

func TestQuaternionExtra(t *testing.T) {
	qp := NewQuaternionByValuesShort(1, 2, 3, 1)
	if qp.W() <= 0 {
		t.Errorf("NewQuaternionByValuesShort positive")
	}
	qs := Quaternion{}
	qs.set(1, 2, 3, 4)
	if qs.X() != 1 || qs.Y() != 2 || qs.Z() != 3 || qs.W() != 4 {
		t.Errorf("set")
	}

	locked := NewQuaternionFromAxisAngles(Vec3{r3Vec(1, 0, 0)}, math.Pi/2)
	lockedRad := locked.ToRadians()
	if math.Abs(lockedRad.Z) > 1e-9 {
		t.Errorf("ToRadians locked")
	}
	if _, g := NewQuaternion().ToRadiansWithGimbal(0); g {
		t.Errorf("ToRadiansWithGimbal axis0")
	}
	if _, g := NewQuaternion().ToRadiansWithGimbal(1); g {
		t.Errorf("ToRadiansWithGimbal axis1")
	}
	if _, g := NewQuaternion().ToRadiansWithGimbal(2); g {
		t.Errorf("ToRadiansWithGimbal axis2")
	}

	qInv := NewQuaternionByValues(1, 0, 0, 1)
	qInv.Inverse()
	if qInv == quaternionZero {
		t.Errorf("Inverse nonzero")
	}
	zeroQ := Quaternion{}
	if zeroQ.Inverted() != quaternionIdent {
		t.Errorf("Inverted zero")
	}

	axis := Vec3{r3Vec(0, 0, 1)}
	base := NewQuaternion()
	if base.SlerpExtended(base, 0.3) != base {
		t.Errorf("SlerpExtended near")
	}
	near := NewQuaternionFromAxisAngles(axis, 1e-3)
	if base.SlerpExtended(near, 0.5) == quaternionZero {
		t.Errorf("SlerpExtended linear")
	}
	far := NewQuaternionFromAxisAngles(axis, math.Pi/2)
	if base.SlerpExtended(far, 0.5) == quaternionZero {
		t.Errorf("SlerpExtended spherical")
	}
	neg := NewQuaternionFromAxisAngles(axis, 3*math.Pi/2)
	if base.SlerpExtended(neg, 0.5) == quaternionZero {
		t.Errorf("SlerpExtended neg")
	}
	if base.Slerp(neg, 0.5) == quaternionZero {
		t.Errorf("Slerp neg")
	}
	if base.Slerp(near, 0.5) == quaternionZero {
		t.Errorf("Slerp linear")
	}

	pos := NewQuaternionByValues(1, 0, 0, 1)
	if pos.ToSignedDegree() < 0 || pos.ToSignedRadian() < 0 {
		t.Errorf("ToSigned positive")
	}
	zeroAxis := NewQuaternionByValues(0, 0, 0, 1)
	if zeroAxis.ToSignedDegree() != 0 || zeroAxis.ToSignedRadian() != 0 {
		t.Errorf("ToSigned zero")
	}

	dir := Vec3{r3Vec(1, 0, 0)}
	up := Vec3{r3Vec(0, 1, 0)}
	if NewQuaternionFromDirection(dir, up).IsIdent() {
		t.Errorf("NewQuaternionFromDirection")
	}
	rot := NewQuaternionRotate(Vec3{r3Vec(0, 1, 0)}, Vec3{r3Vec(0, -1, 0)})
	if rot.IsIdent() {
		t.Errorf("NewQuaternionRotate axis")
	}
	parallel := NewQuaternionRotate(Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(2, 0, 0)})
	if !parallel.IsIdent() {
		t.Errorf("NewQuaternionRotate parallel")
	}

	fixed := NewQuaternionFromAxisAngles(Vec3{r3Vec(1, 0, 0)}, math.Pi/3).ToFixedAxisRotation(Vec3{r3Vec(-1, 0, 0)})
	if fixed == quaternionZero {
		t.Errorf("ToFixedAxisRotation")
	}

	q1 := NewQuaternion()
	q2 := NewQuaternionFromAxisAngles(Vec3{r3Vec(0, 1, 0)}, math.Pi/2)
	qt := q1.Slerp(q2, 0.3)
	if found := FindSlerpT(q1, q2, qt, 0.2); math.IsNaN(found) || found < 0 || found > 1 {
		t.Errorf("FindSlerpT search")
	}
	q2neg := NewQuaternionFromAxisAngles(Vec3{r3Vec(0, 1, 0)}, 3*math.Pi/2)
	if found := FindSlerpT(q1, q2neg, qt, 0.4); math.IsNaN(found) || found < 0 || found > 1 {
		t.Errorf("FindSlerpT neg")
	}
	if found := FindLerpT(q1, q2neg, qt); math.IsNaN(found) || found < 0 || found > 1 {
		t.Errorf("FindLerpT neg")
	}
	if found := FindLerpT(q1, q2, qt); math.IsNaN(found) || found < 0 || found > 1 {
		t.Errorf("FindLerpT pos")
	}
}
