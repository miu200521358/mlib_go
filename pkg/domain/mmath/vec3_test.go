package mmath

import (
	"encoding/json"
	"math"
	"testing"

	"gonum.org/v1/gonum/spatial/r3"
)

func TestVec3Ops(t *testing.T) {
	v := NewVec3()
	if !v.Equals(ZERO_VEC3) {
		t.Errorf("NewVec3")
	}
	data, err := json.Marshal(Vec3{r3Vec(1, 2, 3)})
	if err != nil || len(data) == 0 {
		t.Errorf("MarshalJSON")
	}
	var vjson Vec3
	if err := json.Unmarshal([]byte(`{"x":1,"y":2,"z":3}`), &vjson); err != nil {
		t.Errorf("UnmarshalJSON")
	}
	if vjson != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("UnmarshalJSON value")
	}
	v = Vec3{r3Vec(1, 2, 3)}
	if v.XY() != (Vec2{1, 2}) {
		t.Errorf("XY")
	}
	if !(Vec3{r3Vec(1, 0, 0)}).IsOnlyX() || !(Vec3{r3Vec(0, 1, 0)}).IsOnlyY() || !(Vec3{r3Vec(0, 0, 1)}).IsOnlyZ() {
		t.Errorf("IsOnly")
	}
	if v.String() == "" || v.StringByDigits(3) == "" {
		t.Errorf("String")
	}
	if v.MMD() != v {
		t.Errorf("MMD")
	}
	v2 := Vec3{r3Vec(4, 5, 6)}
	v.Add(v2)
	if v != (Vec3{r3Vec(5, 7, 9)}) {
		t.Errorf("Add")
	}
	if v2.Added(Vec3{r3Vec(1, 1, 1)}) != (Vec3{r3Vec(5, 6, 7)}) {
		t.Errorf("Added")
	}
	if v2.Subed(Vec3{r3Vec(1, 1, 1)}) != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("Subed")
	}
	v2.SubScalar(1)
	if v2 != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("SubScalar")
	}
	v2.Mul(Vec3{r3Vec(2, 2, 2)})
	if v2 != (Vec3{r3Vec(6, 8, 10)}) {
		t.Errorf("Mul")
	}
	if v2.MuledScalar(0.5) != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("MuledScalar")
	}
	v2.DivScalar(2)
	if v2 != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("DivScalar")
	}
	if v2.Dived(Vec3{r3Vec(1, 2, 5)}) != (Vec3{r3Vec(3, 2, 1)}) {
		t.Errorf("Dived")
	}
	if !v2.Equals(Vec3{r3Vec(3, 4, 5)}) || v2.NotEquals(Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("Equals")
	}
	if !v2.NearEquals(Vec3{r3Vec(3.001, 4.001, 5.001)}, 1e-2) || v2.NearEquals(Vec3{r3Vec(4, 5, 6)}, 1e-3) {
		t.Errorf("NearEquals")
	}
	if !v2.LessThan(Vec3{r3Vec(4, 5, 6)}) || v2.GreaterThan(Vec3{r3Vec(4, 5, 6)}) {
		t.Errorf("Compare")
	}
	v2.Negate()
	if v2 != (Vec3{r3Vec(-3, -4, -5)}) {
		t.Errorf("Negate")
	}
	if v2.Negated() != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("Negated")
	}
	v2.Abs()
	if v2 != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("Abs")
	}
	if v2.Absed() != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("Absed")
	}
	if v2.Hash() == 0 {
		t.Errorf("Hash")
	}
	v2.Truncate(4.5)
	if v2 != (Vec3{r3Vec(0, 0, 5)}) {
		t.Errorf("Truncate")
	}
	v3 := Vec3{r3Vec(0, 0, 1)}
	v3.MergeIfZero(2)
	if v3 != (Vec3{r3Vec(2, 2, 1)}) {
		t.Errorf("MergeIfZero")
	}
	v3.MergeIfZeros(Vec3{r3Vec(3, 4, 5)})
	if v3 != (Vec3{r3Vec(2, 2, 1)}) {
		t.Errorf("MergeIfZeros")
	}
	if !ZERO_VEC3.IsZero() || !ONE_VEC3.IsOne() {
		t.Errorf("IsZero/IsOne")
	}
	lenv := Vec3{r3Vec(3, 4, 0)}
	if lenv.Length() != 5 || lenv.LengthSqr() != 25 {
		t.Errorf("Length")
	}
	zero := Vec3{}
	zero.Normalize()
	unit := Vec3{r3Vec(1, 0, 0)}
	unit.Normalize()
	if zero != (Vec3{}) || unit != (Vec3{r3Vec(1, 0, 0)}) {
		t.Errorf("Normalize")
	}
	norm := (Vec3{r3Vec(3, 4, 0)}).Normalized()
	if math.Abs(norm.X-0.6) > 1e-6 || math.Abs(norm.Y-0.8) > 1e-6 {
		t.Errorf("Normalized")
	}
	if angleFromCosVec3(1.5) != 0 || angleFromCosVec3(-1.5) != math.Pi {
		t.Errorf("angleFromCosVec3")
	}
	if math.Abs((Vec3{r3Vec(1, 0, 0)}).Degree(Vec3{r3Vec(1, 0, 0)})) > 1e-9 {
		t.Errorf("Degree")
	}
	if (Vec3{r3Vec(1, 2, 3)}).Dot(Vec3{r3Vec(4, 5, 6)}) != 32 {
		t.Errorf("Dot")
	}
	if (Vec3{r3Vec(1, 0, 0)}).Cross(Vec3{r3Vec(0, 1, 0)}) != (Vec3{r3Vec(0, 0, 1)}) {
		t.Errorf("Cross")
	}
	if (Vec3{r3Vec(2, 1, 3)}).Min() != (Vec3{r3Vec(1, 1, 1)}) || (Vec3{r3Vec(2, 1, 3)}).Max() != (Vec3{r3Vec(3, 3, 3)}) {
		t.Errorf("Min/Max")
	}
	clamp := (Vec3{r3Vec(-1, 2, 0)}).Clamped(ZERO_VEC3, ONE_VEC3)
	if clamp != (Vec3{r3Vec(0, 1, 0)}) {
		t.Errorf("Clamped")
	}
	clamp01 := Vec3{r3Vec(-1, 2, 0)}
	clamp01.Clamp01()
	if clamp01 != (Vec3{r3Vec(0, 1, 0)}) {
		t.Errorf("Clamp01")
	}
	if cp, err := (Vec3{r3Vec(1, 2, 3)}).Copy(); err != nil || cp != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("Copy")
	}
	if len((Vec3{r3Vec(1, 2, 3)}).Vector()) != 3 {
		t.Errorf("Vector")
	}
	if (Vec3{r3Vec(1, 2, 3)}).ToMat4().MulVec3(Vec3{r3Vec(0, 0, 0)}) != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("ToMat4")
	}
	if (Vec3{r3Vec(2, 3, 4)}).ToScaleMat4()[0] != 2 {
		t.Errorf("ToScaleMat4")
	}
	cv := Vec3{r3Vec(1e-7, 1, 1)}
	cv.ClampIfVerySmall()
	if cv.X != 0 {
		t.Errorf("ClampIfVerySmall")
	}
	if (Vec3{r3Vec(math.Pi, 0, 0)}).RadToDeg().X != 180 {
		t.Errorf("RadToDeg")
	}
	if (Vec3{r3Vec(180, 0, 0)}).DegToRad().X != math.Pi {
		t.Errorf("DegToRad")
	}
	if (Vec3{r3Vec(0, 0, 0)}).RadToQuaternion().IsIdent() == false {
		t.Errorf("RadToQuaternion")
	}
	if (Vec3{r3Vec(0, 0, 0)}).DegToQuaternion().IsIdent() == false {
		t.Errorf("DegToQuaternion")
	}
	if (Vec3{r3Vec(1, 2, 3)}).Lerp(Vec3{r3Vec(4, 5, 6)}, 0.5) != (Vec3{r3Vec(2.5, 3.5, 4.5)}) {
		t.Errorf("Lerp")
	}
	if (Vec3{r3Vec(1, 0, 0)}).Slerp(Vec3{r3Vec(1, 0, 0)}, 0.3) != (Vec3{r3Vec(1, 0, 0)}) {
		t.Errorf("Slerp same")
	}
	if (Vec3{}).ToLocalMat() != NewMat4() {
		t.Errorf("ToLocalMat zero")
	}
	if (Vec3{r3Vec(0, 1, 0)}).ToLocalMat() == NewMat4() {
		t.Errorf("ToLocalMat")
	}
	if (Vec3{}).ToScaleLocalMat(Vec3{r3Vec(2, 2, 2)}) != NewMat4() {
		t.Errorf("ToScaleLocalMat zero")
	}
	if (Vec3{r3Vec(1, 0, 0)}).One() == ZERO_VEC3 {
		t.Errorf("One")
	}
	if (Vec3{r3Vec(1, 1, 1)}).Distance(Vec3{r3Vec(2, 1, 1)}) != 1 {
		t.Errorf("Distance")
	}
	if len((Vec3{r3Vec(0, 0, 0)}).Distances([]Vec3{{r3Vec(1, 0, 0)}})) != 1 {
		t.Errorf("Distances")
	}
	vEff := Vec3{r3Vec(math.NaN(), 1, math.Inf(1))}
	vEff.Effective()
	if vEff.X != 0 || vEff.Z != 0 {
		t.Errorf("Effective")
	}
	if DistanceFromPointToLine(Vec3{r3Vec(0, 0, 0)}, Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(0, 1, 0)}) == 0 {
		t.Errorf("DistanceFromPointToLine")
	}
	if DistanceFromPlaneToLine(Vec3{r3Vec(0, 0, 0)}, Vec3{r3Vec(0, 1, 0)}, Vec3{r3Vec(0, 0, 1)}, Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(0, 1, 0)}, Vec3{r3Vec(0, 0, 0)}) != 0 {
		t.Errorf("DistanceFromPlaneToLine")
	}
	if _, err := IntersectLinePlane(Vec3{r3Vec(0, 0, 0)}, Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(0, 0, 1)}, Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(0, 1, 0)}, Vec3{r3Vec(0, 0, 0)}); err == nil {
		t.Errorf("IntersectLinePlane error")
	}
	if hit, err := IntersectLinePlane(Vec3{r3Vec(0, 0, 0)}, Vec3{r3Vec(0, 1, 0)}, Vec3{r3Vec(0, 0, 1)}, Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(0, 1, 0)}, Vec3{r3Vec(0, 1, 0)}); err != nil || hit.Y == 0 {
		t.Errorf("IntersectLinePlane")
	}
	if IntersectLinePoint(Vec3{r3Vec(0, 0, 0)}, Vec3{r3Vec(1, 0, 0)}, Vec3{r3Vec(2, 0, 0)}).X != 2 {
		t.Errorf("IntersectLinePoint")
	}
	if len(DistanceLineToPoints(Vec3{r3Vec(0, 0, 0)}, []Vec3{{r3Vec(0, 0, 1)}})) != 1 {
		t.Errorf("DistanceLineToPoints")
	}
	if (Vec3{r3Vec(1, 0, 0)}).Project(Vec3{r3Vec(1, 0, 0)}) != (Vec3{r3Vec(1, 0, 0)}) {
		t.Errorf("Project")
	}
	if !(Vec3{r3Vec(1, 1, 1)}).IsPointInsideBox(Vec3{r3Vec(0, 0, 0)}, Vec3{r3Vec(2, 2, 2)}) {
		t.Errorf("IsPointInsideBox")
	}
	if (Vec3{r3Vec(1, 0, 0)}).Vec3Diff(Vec3{r3Vec(1, 0, 0)}).IsIdent() == false {
		t.Errorf("Vec3Diff")
	}
	if (Vec3{r3Vec(1.234, 2.345, 3.456)}).Round(0.1) == (Vec3{}) {
		t.Errorf("Round")
	}
	sorted := SortVec3([]Vec3{{r3Vec(1, 2, 3)}, {r3Vec(0, 2, 3)}})
	if sorted[0].X != 0 {
		t.Errorf("SortVec3")
	}
	if MeanVec3([]Vec3{{r3Vec(1, 2, 3)}, {r3Vec(3, 4, 5)}}).X != 2 {
		t.Errorf("MeanVec3")
	}
	if MinVec3([]Vec3{{r3Vec(1, 2, 3)}, {r3Vec(0, 4, 5)}}).X != 0 {
		t.Errorf("MinVec3")
	}
	if MaxVec3([]Vec3{{r3Vec(1, 2, 3)}, {r3Vec(4, 0, 5)}}).X != 4 {
		t.Errorf("MaxVec3")
	}
	if MedianVec3([]Vec3{{r3Vec(1, 9, 3)}, {r3Vec(4, 0, 5)}, {r3Vec(2, 7, 6)}}).X != 2 {
		t.Errorf("MedianVec3")
	}
}

func TestVec3Extra(t *testing.T) {
	var bad Vec3
	if err := bad.UnmarshalJSON([]byte(`{`)); err == nil {
		t.Errorf("UnmarshalJSON error")
	}

	base := Vec3{r3Vec(1, 2, 3)}
	add := base
	add.AddScalar(1)
	if add != (Vec3{r3Vec(2, 3, 4)}) {
		t.Errorf("AddScalar")
	}
	if base.AddedScalar(1) != (Vec3{r3Vec(2, 3, 4)}) {
		t.Errorf("AddedScalar")
	}
	sub := base
	sub.Sub(Vec3{r3Vec(1, 1, 1)})
	if sub != (Vec3{r3Vec(0, 1, 2)}) {
		t.Errorf("Sub")
	}
	if base.SubedScalar(1) != (Vec3{r3Vec(0, 1, 2)}) {
		t.Errorf("SubedScalar")
	}
	mul := base
	mul.MulScalar(2)
	if mul != (Vec3{r3Vec(2, 4, 6)}) {
		t.Errorf("MulScalar")
	}
	if base.Muled(Vec3{r3Vec(2, 3, 4)}) != (Vec3{r3Vec(2, 6, 12)}) {
		t.Errorf("Muled")
	}
	div := Vec3{r3Vec(4, 6, 8)}
	div.Div(Vec3{r3Vec(2, 3, 4)})
	if div != (Vec3{r3Vec(2, 2, 2)}) {
		t.Errorf("Div")
	}
	if base.DivedScalar(2) != (Vec3{r3Vec(0.5, 1, 1.5)}) {
		t.Errorf("DivedScalar")
	}
	if !base.LessThanOrEquals(Vec3{r3Vec(1, 2, 3)}) || base.LessThanOrEquals(Vec3{r3Vec(0, 2, 3)}) {
		t.Errorf("LessThanOrEquals")
	}
	if !base.GreaterThanOrEquals(Vec3{r3Vec(1, 2, 3)}) || base.GreaterThanOrEquals(Vec3{r3Vec(2, 2, 3)}) {
		t.Errorf("GreaterThanOrEquals")
	}

	trunc := Vec3{r3Vec(1e-7, 1e-7, 1e-7)}
	trunc.Truncate(1e-6)
	if trunc != (Vec3{}) {
		t.Errorf("Truncate")
	}
	if (Vec3{r3Vec(1e-7, 1, 1e-7)}).Truncated(1e-6) != (Vec3{r3Vec(0, 1, 0)}) {
		t.Errorf("Truncated")
	}

	merge := Vec3{r3Vec(0, 1, 0)}
	merge.MergeIfZero(2)
	if merge != (Vec3{r3Vec(2, 1, 2)}) {
		t.Errorf("MergeIfZero")
	}
	merge2 := Vec3{}
	merge2.MergeIfZeros(Vec3{r3Vec(1, 2, 3)})
	if merge2 != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("MergeIfZeros")
	}

	zeroVec := Vec3{}
	if zeroVec.Angle(zeroVec) != 0 {
		t.Errorf("Angle denom")
	}
	if (Vec3{r3Vec(2, 3, 1)}).Min() != (Vec3{r3Vec(1, 1, 1)}) {
		t.Errorf("Min Z")
	}
	if (Vec3{r3Vec(1, 2, 3)}).Max() != (Vec3{r3Vec(3, 3, 3)}) {
		t.Errorf("Max Z")
	}
	if (Vec3{r3Vec(-1, 2, -0.5)}).Clamped01() != (Vec3{r3Vec(0, 1, 0)}) {
		t.Errorf("Clamped01")
	}

	tiny := Vec3{r3Vec(1e-7, 1e-7, 1e-7)}
	tiny.ClampIfVerySmall()
	if tiny != (Vec3{}) {
		t.Errorf("ClampIfVerySmall all")
	}

	if base.Lerp(Vec3{r3Vec(3, 4, 5)}, -0.1) != base {
		t.Errorf("Lerp t<=0")
	}
	if base.Lerp(Vec3{r3Vec(3, 4, 5)}, 1.1) != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("Lerp t>=1")
	}
	if base.Lerp(base.AddedScalar(1e-9), 0.5) != base {
		t.Errorf("Lerp near")
	}

	if base.Slerp(Vec3{r3Vec(3, 4, 5)}, -0.1) != base {
		t.Errorf("Slerp t<=0")
	}
	if base.Slerp(Vec3{r3Vec(3, 4, 5)}, 1.1) != (Vec3{r3Vec(3, 4, 5)}) {
		t.Errorf("Slerp t>=1")
	}
	if base.Slerp(base.AddedScalar(1e-9), 0.5) != base {
		t.Errorf("Slerp near")
	}
	if (Vec3{r3Vec(1, 0, 0)}).Slerp(Vec3{r3Vec(0, 1, 0)}, 0.5) == (Vec3{}) {
		t.Errorf("Slerp")
	}

	if (Vec3{r3Vec(0, 1, 0)}).ToLocalMat() == NewMat4() {
		t.Errorf("ToLocalMat up")
	}
	if (Vec3{r3Vec(1, 0, 0)}).ToLocalMat() == NewMat4() {
		t.Errorf("ToLocalMat flat")
	}
	if (Vec3{r3Vec(0, 1, 0)}).ToScaleLocalMat(Vec3{r3Vec(2, 2, 2)}) == NewMat4() {
		t.Errorf("ToScaleLocalMat")
	}

	sorted := SortVec3([]Vec3{{r3Vec(1, 1, 2)}, {r3Vec(1, 1, 1)}, {r3Vec(1, 0, 5)}, {r3Vec(0, 2, 3)}})
	if sorted[0].X != 0 || sorted[len(sorted)-1].Z != 2 {
		t.Errorf("SortVec3 branches")
	}

	if MeanVec3(nil) != (Vec3{}) {
		t.Errorf("MeanVec3 empty")
	}
	if MinVec3(nil) != (Vec3{}) {
		t.Errorf("MinVec3 empty")
	}
	if MaxVec3(nil) != (Vec3{}) {
		t.Errorf("MaxVec3 empty")
	}
	if MedianVec3(nil) != (Vec3{}) {
		t.Errorf("MedianVec3 empty")
	}
}

func r3Vec(x, y, z float64) r3.Vec {
	return r3.Vec{X: x, Y: y, Z: z}
}
