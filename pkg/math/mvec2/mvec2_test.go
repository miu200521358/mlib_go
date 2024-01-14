package mvec2_test

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
)

func TestInterpolate(t *testing.T) {
	v1 := mvec2.T{1, 2}
	v2 := &mvec2.T{3, 4} // Pass the address of v2
	t1 := 0.5
	expected := mvec2.T{2, 3}

	result := v1.Interpolate(v2, t1) // Use v2 as a pointer

	if result != expected {
		t.Errorf("Interpolation failed. Expected %v, got %v", expected, result)
	}
}

func TestPracticallyEquals(t *testing.T) {
	v1 := mvec2.T{1, 2}
	v2 := mvec2.T{1.000001, 2.000001}
	epsilon := 0.00001

	if !v1.PracticallyEquals(&v2, epsilon) {
		t.Errorf("PracticallyEquals failed. Expected true, got false")
	}

	v3 := mvec2.T{1, 2}
	v4 := mvec2.T{1.0001, 2.0001}

	if v3.PracticallyEquals(&v4, epsilon) {
		t.Errorf("PracticallyEquals failed. Expected false, got true")
	}
}

func TestLessThan(t *testing.T) {
	v1 := mvec2.T{1, 2}
	v2 := mvec2.T{3, 4}

	if !v1.LessThan(&v2) {
		t.Errorf("LessThan failed. Expected true, got false")
	}

	v3 := mvec2.T{3, 4}
	v4 := mvec2.T{1, 2}

	if v3.LessThan(&v4) {
		t.Errorf("LessThan failed. Expected false, got true")
	}
}

func TestLessThanOrEquals(t *testing.T) {
	v1 := mvec2.T{1, 2}
	v2 := mvec2.T{3, 4}

	if !v1.LessThanOrEquals(&v2) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}

	v3 := mvec2.T{3, 4}
	v4 := mvec2.T{1, 2}

	if v3.LessThanOrEquals(&v4) {
		t.Errorf("LessThanOrEqual failed. Expected false, got true")
	}

	v5 := mvec2.T{1, 2}
	v6 := mvec2.T{1, 2}

	if !v5.LessThanOrEquals(&v6) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}
}

func TestGreaterThan(t *testing.T) {
	v1 := mvec2.T{1, 2}
	v2 := mvec2.T{3, 4}

	if v1.GreaterThan(&v2) {
		t.Errorf("GreaterThan failed. Expected false, got true")
	}

	v3 := mvec2.T{3, 4}
	v4 := mvec2.T{1, 2}

	if !v3.GreaterThan(&v4) {
		t.Errorf("GreaterThan failed. Expected true, got false")
	}
}

func TestGreaterThanOrEquals(t *testing.T) {
	v1 := mvec2.T{1, 2}
	v2 := mvec2.T{3, 4}

	if v1.GreaterThanOrEquals(&v2) {
		t.Errorf("GreaterThanOrEqual failed. Expected false, got true")
	}

	v3 := mvec2.T{3, 4}
	v4 := mvec2.T{1, 2}

	if !v3.GreaterThanOrEquals(&v4) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}

	v5 := mvec2.T{1, 2}
	v6 := mvec2.T{1, 2}

	if !v5.GreaterThanOrEquals(&v6) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}
}

func TestInverted(t *testing.T) {
	v1 := mvec2.T{1, 2}
	v2 := mvec2.T{3, 4}

	iv1 := v1.Inverted()
	if iv1.GetX() != -1 || iv1.GetY() != -2 {
		t.Errorf("Inverse failed. Expected (-1, -2), got (%v, %v)", iv1.GetX(), iv1.GetY())
	}

	iv2 := v2.Inverted()
	if iv2.GetX() != -3 || iv2.GetY() != -4 {
		t.Errorf("Inverse failed. Expected (-3, -4), got (%v, %v)", iv2.GetX(), iv2.GetY())
	}
}

func TestAbs(t *testing.T) {
	v1 := mvec2.T{-1, -2}
	expected1 := mvec2.T{1, 2}
	result1 := v1.Abs()
	if !result1.Equals(&expected1) {
		t.Errorf("Abs failed. Expected %v, got %v", expected1, result1)
	}

	v2 := mvec2.T{3, -4}
	expected2 := mvec2.T{3, 4}
	result2 := v2.Abs()
	if !result2.Equals(&expected2) {
		t.Errorf("Abs failed. Expected %v, got %v", expected2, result2)
	}

	v3 := mvec2.T{0, 0}
	expected3 := mvec2.T{0, 0}
	result3 := v3.Abs()
	if !result3.Equals(&expected3) {
		t.Errorf("Abs failed. Expected %v, got %v", expected3, result3)
	}
}

func TestHash(t *testing.T) {
	v := mvec2.T{1, 2}
	expected := uint64(4921663092573786862)
	result := v.Hash()
	if result != expected {
		t.Errorf("Hash failed. Expected %v, got %v", expected, result)
	}
}

func TestAngle(t *testing.T) {
	v1 := mvec2.T{1, 0}
	v2 := mvec2.T{0, 1}
	expected := math.Pi / 2
	result := v1.Angle(&v2)
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Angle failed. Expected %v, got %v", expected, result)
	}

	v3 := mvec2.T{1, 1}
	expected2 := math.Pi / 4
	result2 := v1.Angle(&v3)
	if math.Abs(result2-expected2) > 1e-10 {
		t.Errorf("Angle failed. Expected %v, got %v", expected2, result2)
	}
}

func TestDegree(t *testing.T) {
	v1 := mvec2.T{1, 0}
	v2 := mvec2.T{0, 1}
	expected := 90.0
	result := v1.Degree(&v2)
	if math.Abs(result-expected) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected, result)
	}

	v3 := mvec2.T{1, 1}
	expected2 := 45.0
	result2 := v1.Degree(&v3)
	if math.Abs(result2-expected2) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected2, result2)
	}
}
