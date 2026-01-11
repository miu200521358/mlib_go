package mmath

import (
	"math"
	"testing"
)

func TestVec2Ops(t *testing.T) {
	v := NewVec2()
	if !v.Equals(ZERO_VEC2) {
		t.Errorf("NewVec2")
	}
	v = Vec2{1, 2}
	if v.String() == "" {
		t.Errorf("String")
	}
	v2 := Vec2{3, 4}
	v.Add(v2)
	if v.X != 4 || v.Y != 6 {
		t.Errorf("Add")
	}
	v.AddScalar(1)
	if v.X != 5 || v.Y != 7 {
		t.Errorf("AddScalar")
	}
	if v2.Added(Vec2{1, 1}) != (Vec2{4, 5}) {
		t.Errorf("Added")
	}
	if v2.Subed(Vec2{1, 1}) != (Vec2{2, 3}) {
		t.Errorf("Subed")
	}
	v2.SubScalar(1)
	if v2.X != 2 || v2.Y != 3 {
		t.Errorf("SubScalar")
	}
	v2.Mul(Vec2{2, 2})
	if v2.X != 4 || v2.Y != 6 {
		t.Errorf("Mul")
	}
	if v2.MuledScalar(0.5) != (Vec2{2, 3}) {
		t.Errorf("MuledScalar")
	}
	v2.DivScalar(2)
	if v2.X != 2 || v2.Y != 3 {
		t.Errorf("DivScalar")
	}
	if v2.Dived(Vec2{1, 3}) != (Vec2{2, 1}) {
		t.Errorf("Dived")
	}
	if !v2.Equals(Vec2{2, 3}) || v2.NotEquals(Vec2{2, 3}) {
		t.Errorf("Equals")
	}
	if !v2.NearEquals(Vec2{2.00001, 3.00001}, 1e-2) || v2.NearEquals(Vec2{3, 4}, 1e-3) {
		t.Errorf("NearEquals")
	}
	if !v2.LessThan(Vec2{3, 4}) || v2.GreaterThan(Vec2{3, 4}) {
		t.Errorf("Compare")
	}
	v2.Negate()
	if v2.X != -2 || v2.Y != -3 {
		t.Errorf("Negate")
	}
	if v2.Negated() != (Vec2{2, 3}) {
		t.Errorf("Negated")
	}
	v2.Abs()
	if v2.X != 2 || v2.Y != 3 {
		t.Errorf("Abs")
	}
	if v2.Absed() != (Vec2{2, 3}) {
		t.Errorf("Absed")
	}
	if v2.Hash() == 0 {
		t.Errorf("Hash")
	}
	if ZERO_VEC2.IsZero() != true || v2.IsZero() {
		t.Errorf("IsZero")
	}
	lenv := Vec2{3, 4}
	if lenv.Length() != 5 || lenv.LengthSqr() != 25 {
		t.Errorf("Length")
	}
	zero := Vec2{}
	zero.Normalize()
	if zero != (Vec2{}) {
		t.Errorf("Normalize zero")
	}
	unit := Vec2{1, 0}
	unit.Normalize()
	if unit != (Vec2{1, 0}) {
		t.Errorf("Normalize unit")
	}
	normSrc := Vec2{3, 4}
	norm := normSrc.Normalized()
	if math.Abs(norm.X-0.6) > 1e-6 || math.Abs(norm.Y-0.8) > 1e-6 {
		t.Errorf("Normalized")
	}
	if zero.Angle(unit) != 0 {
		t.Errorf("Angle denom")
	}
	if angleFromCosVec2(1.5) == 0 || angleFromCosVec2(-1.5) == 0 {
		t.Errorf("angleFromCosVec2")
	}
	if math.Abs(unit.Degree(unit)) > 1e-9 {
		t.Errorf("Degree")
	}
	dotSrc := Vec2{1, 2}
	if dotSrc.Dot(Vec2{3, 4}) != 11 {
		t.Errorf("Dot")
	}
	if dotSrc.Cross(Vec2{3, 4}) == (Vec2{}) {
		t.Errorf("Cross")
	}
	minmax := Vec2{2, 1}
	if minmax.Min() != (Vec2{1, 1}) || minmax.Max() != (Vec2{2, 2}) {
		t.Errorf("Min/Max")
	}
	clampSrc := Vec2{-1, 2}
	clamp := clampSrc.Clamped(ZERO_VEC2, UNIT_XY_VEC2)
	if clamp != (Vec2{0, 1}) {
		t.Errorf("Clamped")
	}
	clamp01 := Vec2{-1, 2}
	clamp01.Clamp01()
	if clamp01 != (Vec2{0, 1}) {
		t.Errorf("Clamp01")
	}
	rot := Vec2{1, 1}
	rot.Rotate(math.Pi / 4)
	if rot == (Vec2{}) {
		t.Errorf("Rotate")
	}
	rot2Src := Vec2{1, 1}
	rot2 := rot2Src.Rotated(0)
	if rot2 != (Vec2{1, 1}) {
		t.Errorf("Rotated")
	}
	rot3 := Vec2{1, 0}
	rot3.RotateAroundPoint(Vec2{1, 0}, math.Pi/2)
	if rot3 != (Vec2{1, 0}) {
		t.Errorf("RotateAroundPoint")
	}
	copySrc := Vec2{1, 2}
	if cp, err := copySrc.Copy(); err != nil || cp != (Vec2{1, 2}) {
		t.Errorf("Copy")
	}
	if len(copySrc.Vector()) != 2 {
		t.Errorf("Vector")
	}
	lerpSrc := Vec2{1, 2}
	if lerpSrc.Lerp(Vec2{3, 4}, 0.5) != (Vec2{2, 3}) {
		t.Errorf("Lerp")
	}
	roundSrc := Vec2{1.2, 2.6}
	if roundSrc.Round() != (Vec2{1, 3}) {
		t.Errorf("Round")
	}
	oneSrc := Vec2{0, 1e-9}
	if oneSrc.One() == (Vec2{}) {
		t.Errorf("One")
	}
	distSrc := Vec2{1, 1}
	if distSrc.Distance(Vec2{2, 1}) != 1 {
		t.Errorf("Distance")
	}
	cvs := Vec2{1e-7, 1}
	cvs.ClampIfVerySmall()
	if cvs.X != 0 {
		t.Errorf("ClampIfVerySmall")
	}
}

func TestVec2Extra(t *testing.T) {
	base := Vec2{1, 2}
	if base.AddedScalar(1) != (Vec2{2, 3}) {
		t.Errorf("AddedScalar")
	}
	if base.SubedScalar(1) != (Vec2{0, 1}) {
		t.Errorf("SubedScalar")
	}
	if base.Muled(Vec2{2, 3}) != (Vec2{2, 6}) {
		t.Errorf("Muled")
	}
	div := Vec2{4, 6}
	div.Div(Vec2{2, 3})
	if div != (Vec2{2, 2}) {
		t.Errorf("Div")
	}
	if !base.LessThanOrEquals(Vec2{1, 2}) || base.LessThanOrEquals(Vec2{0, 2}) {
		t.Errorf("LessThanOrEquals")
	}
	if !base.GreaterThanOrEquals(Vec2{1, 2}) || base.GreaterThanOrEquals(Vec2{2, 2}) {
		t.Errorf("GreaterThanOrEquals")
	}
	if (Vec2{1, 2}).Max() != (Vec2{2, 2}) {
		t.Errorf("Max")
	}
	if (Vec2{-1, 2}).Clamped01() != (Vec2{0, 1}) {
		t.Errorf("Clamped01")
	}
	if base.Lerp(Vec2{3, 4}, -0.1) != base {
		t.Errorf("Lerp t<=0")
	}
	if base.Lerp(Vec2{3, 4}, 1.1) != (Vec2{3, 4}) {
		t.Errorf("Lerp t>=1")
	}
	if base.Lerp(base, 0.5) != base {
		t.Errorf("Lerp equals")
	}
	tiny := Vec2{1e-7, 1e-7}
	tiny.ClampIfVerySmall()
	if tiny != (Vec2{}) {
		t.Errorf("ClampIfVerySmall all")
	}
}
