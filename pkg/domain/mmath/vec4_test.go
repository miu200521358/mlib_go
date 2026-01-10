package mmath

import "testing"

func TestVec4Ops(t *testing.T) {
	v := NewVec4()
	if !v.Equals(ZERO_VEC4) {
		t.Errorf("NewVec4")
	}
	v = Vec4{1, 2, 3, 1}
	if v.XY() != (Vec2{1, 2}) || v.XYZ() != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("XY/XYZ")
	}
	if v.String() == "" || v.MMD() != v {
		t.Errorf("String/MMD")
	}
	v2 := Vec4{3, 4, 5, 1}
	v.Add(v2)
	if v != (Vec4{4, 6, 8, 2}) {
		t.Errorf("Add")
	}
	if v2.AddedScalar(1) != (Vec4{4, 5, 6, 2}) {
		t.Errorf("AddedScalar")
	}
	v2.SubScalar(1)
	if v2 != (Vec4{2, 3, 4, 0}) {
		t.Errorf("SubScalar")
	}
	v2.Mul(Vec4{2, 2, 2, 2})
	if v2 != (Vec4{4, 6, 8, 0}) {
		t.Errorf("Mul")
	}
	if v2.MuledScalar(0.5) != (Vec4{2, 3, 4, 0}) {
		t.Errorf("MuledScalar")
	}
	v2.DivScalar(2)
	if v2 != (Vec4{2, 3, 4, 0}) {
		t.Errorf("DivScalar")
	}
	if v2.Dived(Vec4{1, 3, 4, 1}) != (Vec4{2, 1, 1, 0}) {
		t.Errorf("Dived")
	}
	if !v2.Equals(Vec4{2, 3, 4, 0}) || v2.NotEquals(Vec4{2, 3, 4, 0}) {
		t.Errorf("Equals")
	}
	if !v2.NearEquals(Vec4{2.001, 3.001, 4.001, 0.001}, 1e-2) || v2.NearEquals(Vec4{3, 4, 5, 1}, 1e-3) {
		t.Errorf("NearEquals")
	}
	if !v2.LessThan(Vec4{3, 4, 5, 1}) || v2.GreaterThan(Vec4{3, 4, 5, 1}) {
		t.Errorf("Compare")
	}
	v2.Negate()
	if v2 != (Vec4{-2, -3, -4, -0}) {
		t.Errorf("Negate")
	}
	if v2.Negated() != (Vec4{2, 3, 4, 0}) {
		t.Errorf("Negated")
	}
	v2.Abs()
	if v2 != (Vec4{2, 3, 4, 0}) {
		t.Errorf("Abs")
	}
	if v2.Absed() != (Vec4{2, 3, 4, 0}) {
		t.Errorf("Absed")
	}
	if v2.Hash() == 0 {
		t.Errorf("Hash")
	}
	if ZERO_VEC4.IsZero() != true || v2.IsZero() {
		t.Errorf("IsZero")
	}
	if v.Length() == 0 || v.LengthSqr() == 0 {
		t.Errorf("Length")
	}
	v3 := Vec4{1, 0, 0, 1}
	v3.Normalize()
	if v3 != (Vec4{1, 0, 0, 1}) {
		t.Errorf("Normalize")
	}
	if v3.Normalized() != (Vec4{1, 0, 0, 1}) {
		t.Errorf("Normalized")
	}
	if (Vec4{1, 2, 3, 1}).Dot(Vec4{4, 5, 6, 1}) == 0 {
		t.Errorf("Dot")
	}
	if Dot4(Vec4{1, 2, 3, 4}, Vec4{1, 1, 1, 1}) != 10 {
		t.Errorf("Dot4")
	}
	if (Vec4{1, 0, 0, 1}).Cross(Vec4{0, 1, 0, 1}) != (Vec4{0, 0, 1, 1}) {
		t.Errorf("Cross")
	}
	if (Vec4{2, 1, 3, 4}).Min() != (Vec4{1, 1, 1, 1}) || (Vec4{2, 1, 3, 4}).Max() != (Vec4{4, 4, 4, 4}) {
		t.Errorf("Min/Max")
	}
	clamp := (Vec4{-1, 2, 0, 0}).Clamped(ZERO_VEC4, ONE_VEC4)
	if clamp != (Vec4{0, 1, 0, 0}) {
		t.Errorf("Clamped")
	}
	clamp01 := Vec4{-1, 2, 0, 0}
	clamp01.Clamp01()
	if clamp01 != (Vec4{0, 1, 0, 0}) {
		t.Errorf("Clamp01")
	}
	if cp, err := (Vec4{1, 2, 3, 4}).Copy(); err != nil || *cp != (Vec4{1, 2, 3, 4}) {
		t.Errorf("Copy")
	}
	if len((Vec4{1, 2, 3, 4}).Vector()) != 4 {
		t.Errorf("Vector")
	}
	if (Vec4{1, 2, 3, 4}).Lerp(Vec4{3, 4, 5, 6}, 0.5) != (Vec4{2, 3, 4, 5}) {
		t.Errorf("Lerp")
	}
	if (Vec4{2, 2, 2, 2}).Vec3DividedByW() != (Vec3{r3Vec(1, 1, 1)}) {
		t.Errorf("Vec3DividedByW")
	}
	if (Vec4{2, 2, 2, 2}).DividedByW() != (Vec4{1, 1, 1, 1}) {
		t.Errorf("DividedByW")
	}
	v4 := Vec4{2, 2, 2, 2}
	v4.DivideByW()
	if v4 != (Vec4{1, 1, 1, 1}) {
		t.Errorf("DivideByW")
	}
	if (Vec4{0, 1e-15, 0, 0}).One() == ZERO_VEC4 {
		t.Errorf("One")
	}
	if (Vec4{1, 1, 1, 1}).Distance(Vec4{2, 1, 1, 1}) != 1 {
		t.Errorf("Distance")
	}
	cv := Vec4{1e-7, 1, 1, 1}
	cv.ClampIfVerySmall()
	if cv.X != 0 {
		t.Errorf("ClampIfVerySmall")
	}
	if (Vec4{1.234, 2.345, 3.456, 4.567}).Round(0.1) == ZERO_VEC4 {
		t.Errorf("Round")
	}
}

func TestVec4Extra(t *testing.T) {
	base := Vec4{1, 2, 3, 1}
	add := base
	add.AddScalar(1)
	if add != (Vec4{2, 3, 4, 2}) {
		t.Errorf("AddScalar")
	}
	sub := base
	sub.Sub(Vec4{1, 1, 1, 1})
	if sub != (Vec4{0, 1, 2, 0}) {
		t.Errorf("Sub")
	}
	if base.SubedScalar(1) != (Vec4{0, 1, 2, 0}) {
		t.Errorf("SubedScalar")
	}
	mul := base
	mul.MulScalar(2)
	if mul != (Vec4{2, 4, 6, 2}) {
		t.Errorf("MulScalar")
	}
	if base.Muled(Vec4{2, 3, 4, 1}) != (Vec4{2, 6, 12, 1}) {
		t.Errorf("Muled")
	}
	div := Vec4{4, 6, 8, 2}
	div.Div(Vec4{2, 3, 4, 1})
	if div != (Vec4{2, 2, 2, 2}) {
		t.Errorf("Div")
	}
	if base.DivedScalar(2) != (Vec4{0.5, 1, 1.5, 0.5}) {
		t.Errorf("DivedScalar")
	}
	if !base.LessThanOrEquals(Vec4{1, 2, 3, 1}) || base.LessThanOrEquals(Vec4{0, 2, 3, 1}) {
		t.Errorf("LessThanOrEquals")
	}
	if !base.GreaterThanOrEquals(Vec4{1, 2, 3, 1}) || base.GreaterThanOrEquals(Vec4{2, 2, 3, 1}) {
		t.Errorf("GreaterThanOrEquals")
	}
	if (Vec4{4, 3, 2, 1}).Min() != (Vec4{1, 1, 1, 1}) {
		t.Errorf("Min")
	}
	if (Vec4{1, 2, 3, 4}).Max() != (Vec4{4, 4, 4, 4}) {
		t.Errorf("Max")
	}
	if (Vec4{-1, 2, 0, 2}).Clamped01() != (Vec4{0, 1, 0, 1}) {
		t.Errorf("Clamped01")
	}
	if base.Lerp(Vec4{3, 4, 5, 6}, -0.1) != base {
		t.Errorf("Lerp t<=0")
	}
	if base.Lerp(Vec4{3, 4, 5, 6}, 1.1) != (Vec4{3, 4, 5, 6}) {
		t.Errorf("Lerp t>=1")
	}
	if base.Lerp(base, 0.5) != base {
		t.Errorf("Lerp equals")
	}
	tiny := Vec4{1e-7, 1e-7, 1e-7, 1e-7}
	tiny.ClampIfVerySmall()
	if tiny != (Vec4{}) {
		t.Errorf("ClampIfVerySmall all")
	}
}
