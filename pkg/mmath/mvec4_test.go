package mmath

import (
	"math"
	"testing"
)

func TestMVec4Interpolate(t *testing.T) {
	v1 := MVec4{1, 2, 3}
	v2 := &MVec4{4, 5, 6} // Pass the address of v2
	t1 := 0.5
	expected := MVec4{2.5, 3.5, 4.5}

	result := v1.Interpolate(v2, t1) // Use v2 as a pointer

	if !result.PracticallyEquals(&expected, 1e-8) {
		t.Errorf("Interpolation failed. Expected %v, got %v", expected, result)
	}
}

func TestMVec4PracticallyEquals(t *testing.T) {
	v1 := MVec4{1, 2, 3}
	v2 := MVec4{1.000001, 2.000001, 3.000001}
	epsilon := 0.00001

	if !v1.PracticallyEquals(&v2, epsilon) {
		t.Errorf("PracticallyEquals failed. Expected true, got false")
	}

	v3 := MVec4{1, 2, 3}
	v4 := MVec4{1.0001, 2.0001, 3.0001}

	if v3.PracticallyEquals(&v4, epsilon) {
		t.Errorf("PracticallyEquals failed. Expected false, got true")
	}
}

func TestMVec4LessThan(t *testing.T) {
	v1 := MVec4{1, 2, 3, 4}
	v2 := MVec4{4, 5, 6, 7}

	if !v1.LessThan(&v2) {
		t.Errorf("LessThan failed. Expected true, got false")
	}

	v3 := MVec4{3, 4, 5, 6}
	v4 := MVec4{1, 2, 3, 4}

	if v3.LessThan(&v4) {
		t.Errorf("LessThan failed. Expected false, got true")
	}
}

func TestMVec4LessThanOrEquals(t *testing.T) {
	v1 := MVec4{1, 2, 3, 4}
	v2 := MVec4{3, 4, 5, 6}

	if !v1.LessThanOrEquals(&v2) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}

	v3 := MVec4{3, 4, 5, 6}
	v4 := MVec4{1, 2, 3, 4}

	if v3.LessThanOrEquals(&v4) {
		t.Errorf("LessThanOrEqual failed. Expected false, got true")
	}

	v5 := MVec4{1, 2, 3, 4}
	v6 := MVec4{1, 2, 3, 4}

	if !v5.LessThanOrEquals(&v6) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}
}

func TestMVec4GreaterThan(t *testing.T) {
	v1 := MVec4{1, 2, 3, 4}
	v2 := MVec4{3, 4, 5, 6}

	if v1.GreaterThan(&v2) {
		t.Errorf("GreaterThan failed. Expected false, got true")
	}

	v3 := MVec4{3, 4, 5, 6}
	v4 := MVec4{1, 2, 3, 4}

	if !v3.GreaterThan(&v4) {
		t.Errorf("GreaterThan failed. Expected true, got false")
	}
}

func TestMVec4GreaterThanOrEquals(t *testing.T) {
	v1 := MVec4{1, 2, 3}
	v2 := MVec4{3, 4, 5}

	if v1.GreaterThanOrEquals(&v2) {
		t.Errorf("GreaterThanOrEqual failed. Expected false, got true")
	}

	v3 := MVec4{3, 4, 5}
	v4 := MVec4{1, 2, 3}

	if !v3.GreaterThanOrEquals(&v4) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}

	v5 := MVec4{1, 2, 3}
	v6 := MVec4{1, 2, 3}

	if !v5.GreaterThanOrEquals(&v6) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}
}

func TestMVec4Inverted(t *testing.T) {
	v1 := MVec4{1, 2, 3}
	v2 := MVec4{3, 4, 5}

	iv1 := v1.Inverted()
	if iv1.GetX() != -1 || iv1.GetY() != -2 || iv1.GetZ() != -3 {
		t.Errorf("Inverse failed. Expected (-1, -2, -3), got (%v, %v, %v)", iv1.GetX(), iv1.GetY(), iv1.GetZ())
	}

	iv2 := v2.Inverted()
	if iv2.GetX() != -3 || iv2.GetY() != -4 || iv2.GetZ() != -5 {
		t.Errorf("Inverse failed. Expected (-3, -4, -5), got (%v, %v, %v)", iv2.GetX(), iv2.GetY(), iv2.GetZ())
	}
}

func TestMVec4Abs(t *testing.T) {
	v1 := MVec4{-1, -2, -3}
	expected1 := MVec4{1, 2, 3}
	result1 := v1.Abs()
	if !result1.Equals(&expected1) {
		t.Errorf("Abs failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec4{3, -4, 5}
	expected2 := MVec4{3, 4, 5}
	result2 := v2.Abs()
	if !result2.Equals(&expected2) {
		t.Errorf("Abs failed. Expected %v, got %v", expected2, result2)
	}

	v3 := MVec4{0, 0}
	expected3 := MVec4{0, 0}
	result3 := v3.Abs()
	if !result3.Equals(&expected3) {
		t.Errorf("Abs failed. Expected %v, got %v", expected3, result3)
	}
}

func TestMVec4Hash(t *testing.T) {
	v := MVec4{1, 2, 3, 4}
	expected := uint64(13473159861922604751)
	result := v.Hash()
	if result != expected {
		t.Errorf("Hash failed. Expected %v, got %v", expected, result)
	}
}

func TestMVec4Angle(t *testing.T) {
	v1 := MVec4{1, 0, 0, 1}
	v2 := MVec4{0, 1, 0, 1}
	expected := math.Pi / 2
	result := v1.Angle(&v2)
	if result != expected {
		t.Errorf("Angle failed. Expected %v, got %v", expected, result)
	}

	v3 := MVec4{1, 1, 1, 1}
	expected2 := 0.9553166181245092
	result2 := v1.Angle(&v3)
	if result2 != expected2 {
		t.Errorf("Angle failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMVec4Degree(t *testing.T) {
	v1 := MVec4{1, 0, 0}
	v2 := MVec4{0, 1, 0}
	expected := 90.0
	result := v1.Degree(&v2)
	if math.Abs(result-expected) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected, result)
	}

	v3 := MVec4{1, 1, 1}
	expected2 := 54.735610317245346
	result2 := v1.Degree(&v3)
	if math.Abs(result2-expected2) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected2, result2)
	}
}
