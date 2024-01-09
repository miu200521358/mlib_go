package mvec3_test

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

func TestInterpolate(t *testing.T) {
	v1 := mvec3.T{1, 2, 3}
	v2 := &mvec3.T{4, 5, 6} // Pass the address of v2
	t1 := 0.5
	expected := mvec3.T{2.5, 3.5, 4.5}

	result := v1.Interpolate(v2, t1) // Use v2 as a pointer

	if result != expected {
		t.Errorf("Interpolation failed. Expected %v, got %v", expected, result)
	}
}

func TestPracticallyEquals(t *testing.T) {
	v1 := mvec3.T{1, 2, 3}
	v2 := mvec3.T{1.000001, 2.000001, 3.000001}
	epsilon := 0.00001

	if !v1.PracticallyEquals(&v2, epsilon) {
		t.Errorf("PracticallyEquals failed. Expected true, got false")
	}

	v3 := mvec3.T{1, 2, 3}
	v4 := mvec3.T{1.0001, 2.0001, 3.0001}

	if v3.PracticallyEquals(&v4, epsilon) {
		t.Errorf("PracticallyEquals failed. Expected false, got true")
	}
}

func TestLessThan(t *testing.T) {
	v1 := mvec3.T{1, 2, 3}
	v2 := mvec3.T{4, 5, 6}

	if !v1.LessThan(&v2) {
		t.Errorf("LessThan failed. Expected true, got false")
	}

	v3 := mvec3.T{3, 4, 5}
	v4 := mvec3.T{1, 2, 3}

	if v3.LessThan(&v4) {
		t.Errorf("LessThan failed. Expected false, got true")
	}
}

func TestLessThanOrEqual(t *testing.T) {
	v1 := mvec3.T{1, 2, 3}
	v2 := mvec3.T{3, 4, 5}

	if !v1.LessThanOrEqual(&v2) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}

	v3 := mvec3.T{3, 4, 5}
	v4 := mvec3.T{1, 2, 3}

	if v3.LessThanOrEqual(&v4) {
		t.Errorf("LessThanOrEqual failed. Expected false, got true")
	}

	v5 := mvec3.T{1, 2}
	v6 := mvec3.T{1, 2}

	if !v5.LessThanOrEqual(&v6) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}
}

func TestGreaterThan(t *testing.T) {
	v1 := mvec3.T{1, 2, 3}
	v2 := mvec3.T{3, 4, 5}

	if v1.GreaterThan(&v2) {
		t.Errorf("GreaterThan failed. Expected false, got true")
	}

	v3 := mvec3.T{3, 4, 5}
	v4 := mvec3.T{1, 2, 3}

	if !v3.GreaterThan(&v4) {
		t.Errorf("GreaterThan failed. Expected true, got false")
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	v1 := mvec3.T{1, 2, 3}
	v2 := mvec3.T{3, 4, 5}

	if v1.GreaterThanOrEqual(&v2) {
		t.Errorf("GreaterThanOrEqual failed. Expected false, got true")
	}

	v3 := mvec3.T{3, 4, 5}
	v4 := mvec3.T{1, 2, 3}

	if !v3.GreaterThanOrEqual(&v4) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}

	v5 := mvec3.T{1, 2, 3}
	v6 := mvec3.T{1, 2, 3}

	if !v5.GreaterThanOrEqual(&v6) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}
}

func TestInverted(t *testing.T) {
	v1 := mvec3.T{1, 2, 3}
	v2 := mvec3.T{3, 4, 5}

	iv1 := v1.Inverted()
	if iv1.GetX() != -1 || iv1.GetY() != -2 || iv1.GetZ() != -3 {
		t.Errorf("Inverse failed. Expected (-1, -2, -3), got (%v, %v, %v)", iv1.GetX(), iv1.GetY(), iv1.GetZ())
	}

	iv2 := v2.Inverted()
	if iv2.GetX() != -3 || iv2.GetY() != -4 || iv2.GetZ() != -5 {
		t.Errorf("Inverse failed. Expected (-3, -4, -5), got (%v, %v, %v)", iv2.GetX(), iv2.GetY(), iv2.GetZ())
	}
}

func TestAbs(t *testing.T) {
	v1 := mvec3.T{-1, -2, -3}
	expected1 := mvec3.T{1, 2, 3}
	result1 := v1.Abs()
	if !result1.Equal(&expected1) {
		t.Errorf("Abs failed. Expected %v, got %v", expected1, result1)
	}

	v2 := mvec3.T{3, -4, 5}
	expected2 := mvec3.T{3, 4, 5}
	result2 := v2.Abs()
	if !result2.Equal(&expected2) {
		t.Errorf("Abs failed. Expected %v, got %v", expected2, result2)
	}

	v3 := mvec3.T{0, 0}
	expected3 := mvec3.T{0, 0}
	result3 := v3.Abs()
	if !result3.Equal(&expected3) {
		t.Errorf("Abs failed. Expected %v, got %v", expected3, result3)
	}
}

func TestHash(t *testing.T) {
	v := mvec3.T{1, 2, 3}
	expected := uint64(17648364615301650315)
	result := v.Hash()
	if result != expected {
		t.Errorf("Hash failed. Expected %v, got %v", expected, result)
	}
}

func TestAngle(t *testing.T) {
	v1 := mvec3.T{1, 0, 0}
	v2 := mvec3.T{0, 1, 0}
	expected := math.Pi / 2
	result := v1.Angle(&v2)
	if result != expected {
		t.Errorf("Angle failed. Expected %v, got %v", expected, result)
	}

	v3 := mvec3.T{1, 1, 1}
	expected2 := 0.9553166181245092
	result2 := v1.Angle(&v3)
	if result2 != expected2 {
		t.Errorf("Angle failed. Expected %v, got %v", expected2, result2)
	}
}

func TestDegree(t *testing.T) {
	v1 := mvec3.T{1, 0, 0}
	v2 := mvec3.T{0, 1, 0}
	expected := 90.0
	result := v1.Degree(&v2)
	if math.Abs(result-expected) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected, result)
	}

	v3 := mvec3.T{1, 1, 1}
	expected2 := 54.735610317245346
	result2 := v1.Degree(&v3)
	if math.Abs(result2-expected2) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected2, result2)
	}
}
