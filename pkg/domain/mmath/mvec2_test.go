package mmath

import (
	"math"
	"testing"
)

func TestMVec2LerpVec2(t *testing.T) {
	v1 := &MVec2{1, 2}
	v2 := &MVec2{3, 4} // Pass the address of v2
	t1 := 0.5
	expected := MVec2{2, 3}

	result := LerpVec2(v1, v2, t1) // Use v2 as a pointer

	if !result.NearEquals(&expected, 1e-8) {
		t.Errorf("TestMVec2LerpVec2 failed. Expected %v, got %v", expected, result)
	}
}

func TestMVec2NearEquals(t *testing.T) {
	v1 := MVec2{1, 2}
	v2 := MVec2{1.000001, 2.000001}
	epsilon := 0.00001

	if !v1.NearEquals(&v2, epsilon) {
		t.Errorf("NearEquals failed. Expected true, got false")
	}

	v3 := MVec2{1, 2}
	v4 := MVec2{1.0001, 2.0001}

	if v3.NearEquals(&v4, epsilon) {
		t.Errorf("NearEquals failed. Expected false, got true")
	}
}

func TestMVec2LessThan(t *testing.T) {
	v1 := MVec2{1, 2}
	v2 := MVec2{3, 4}

	if !v1.LessThan(&v2) {
		t.Errorf("LessThan failed. Expected true, got false")
	}

	v3 := MVec2{3, 4}
	v4 := MVec2{1, 2}

	if v3.LessThan(&v4) {
		t.Errorf("LessThan failed. Expected false, got true")
	}
}

func TestMVec2LessThanOrEquals(t *testing.T) {
	v1 := MVec2{1, 2}
	v2 := MVec2{3, 4}

	if !v1.LessThanOrEquals(&v2) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}

	v3 := MVec2{3, 4}
	v4 := MVec2{1, 2}

	if v3.LessThanOrEquals(&v4) {
		t.Errorf("LessThanOrEqual failed. Expected false, got true")
	}

	v5 := MVec2{1, 2}
	v6 := MVec2{1, 2}

	if !v5.LessThanOrEquals(&v6) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}
}

func TestMVec2GreaterThan(t *testing.T) {
	v1 := MVec2{1, 2}
	v2 := MVec2{3, 4}

	if v1.GreaterThan(&v2) {
		t.Errorf("GreaterThan failed. Expected false, got true")
	}

	v3 := MVec2{3, 4}
	v4 := MVec2{1, 2}

	if !v3.GreaterThan(&v4) {
		t.Errorf("GreaterThan failed. Expected true, got false")
	}
}

func TestMVec2GreaterThanOrEquals(t *testing.T) {
	v1 := MVec2{1, 2}
	v2 := MVec2{3, 4}

	if v1.GreaterThanOrEquals(&v2) {
		t.Errorf("GreaterThanOrEqual failed. Expected false, got true")
	}

	v3 := MVec2{3, 4}
	v4 := MVec2{1, 2}

	if !v3.GreaterThanOrEquals(&v4) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}

	v5 := MVec2{1, 2}
	v6 := MVec2{1, 2}

	if !v5.GreaterThanOrEquals(&v6) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}
}

func TestMVec2Inverted(t *testing.T) {
	v1 := MVec2{1, 2}
	v2 := MVec2{3, 4}

	iv1 := v1.Inverted()
	if iv1.X != -1 || iv1.Y != -2 {
		t.Errorf("Inverse failed. Expected (-1, -2), got (%v, %v)", iv1.X, iv1.Y)
	}

	iv2 := v2.Inverted()
	if iv2.X != -3 || iv2.Y != -4 {
		t.Errorf("Inverse failed. Expected (-3, -4), got (%v, %v)", iv2.X, iv2.Y)
	}
}

func TestMVec2Abs(t *testing.T) {
	v1 := MVec2{-1, -2}
	expected1 := MVec2{1, 2}
	result1 := v1.Abs()
	if !result1.Equals(&expected1) {
		t.Errorf("Abs failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec2{3, -4}
	expected2 := MVec2{3, 4}
	result2 := v2.Abs()
	if !result2.Equals(&expected2) {
		t.Errorf("Abs failed. Expected %v, got %v", expected2, result2)
	}

	v3 := MVec2{0, 0}
	expected3 := MVec2{0, 0}
	result3 := v3.Abs()
	if !result3.Equals(&expected3) {
		t.Errorf("Abs failed. Expected %v, got %v", expected3, result3)
	}
}

func TestMVec2Hash(t *testing.T) {
	v := MVec2{1, 2}
	expected := uint64(4921663092573786862)
	result := v.Hash()
	if result != expected {
		t.Errorf("Hash failed. Expected %v, got %v", expected, result)
	}
}

func TestMVec2Angle(t *testing.T) {
	v1 := MVec2{1, 0}
	v2 := MVec2{0, 1}
	expected := math.Pi / 2
	result := v1.Angle(&v2)
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Angle failed. Expected %v, got %v", expected, result)
	}

	v3 := MVec2{1, 1}
	expected2 := math.Pi / 4
	result2 := v1.Angle(&v3)
	if math.Abs(result2-expected2) > 1e-10 {
		t.Errorf("Angle failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMVec2Degree(t *testing.T) {
	v1 := MVec2{1, 0}
	v2 := MVec2{0, 1}
	expected := 90.0
	result := v1.Degree(&v2)
	if math.Abs(result-expected) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected, result)
	}

	v3 := MVec2{1, 1}
	expected2 := 45.0
	result2 := v1.Degree(&v3)
	if math.Abs(result2-expected2) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected2, result2)
	}
}
