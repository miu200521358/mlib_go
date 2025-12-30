// 指示: miu200521358
package mmath

import (
	"testing"
)

func TestMVec4Lerp(t *testing.T) {
	v1 := MVec4{1, 2, 3, 4}
	v2 := &MVec4{4, 5, 6, 8} // Pass the address of v2
	t1 := 0.5
	expected := MVec4{2.5, 3.5, 4.5, 6}

	result := v1.Lerp(v2, t1) // Use v2 as a pointer

	if !result.NearEquals(&expected, 1e-8) {
		t.Errorf("Lerp failed. Expected %v, got %v", expected, result)
	}
}

func TestMVec4NearEquals(t *testing.T) {
	v1 := MVec4{1, 2, 3, 4}
	v2 := MVec4{1.000001, 2.000001, 3.000001, 4.000001}
	epsilon := 0.00001

	if !v1.NearEquals(&v2, epsilon) {
		t.Errorf("NearEquals failed. Expected true, got false")
	}

	v3 := MVec4{1, 2, 3, 4}
	v4 := MVec4{1.0001, 2.0001, 3.0001, 4.0001}

	if v3.NearEquals(&v4, epsilon) {
		t.Errorf("NearEquals failed. Expected false, got true")
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
	v1 := MVec4{1, 2, 3, 4}
	v2 := MVec4{3, 4, 5, 6}

	if v1.GreaterThanOrEquals(&v2) {
		t.Errorf("GreaterThanOrEqual failed. Expected false, got true")
	}

	v3 := MVec4{3, 4, 5, 6}
	v4 := MVec4{1, 2, 3, 4}

	if !v3.GreaterThanOrEquals(&v4) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}

	v5 := MVec4{1, 2, 3, 4}
	v6 := MVec4{1, 2, 3, 4}

	if !v5.GreaterThanOrEquals(&v6) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}
}

func TestMVec4Negated(t *testing.T) {
	v1 := MVec4{1, 2, 3, 4}
	v2 := MVec4{3, 4, 5, 6}

	iv1 := v1.Negated()
	if iv1.X != -1 || iv1.Y != -2 || iv1.Z != -3 {
		t.Errorf("Inverse failed. Expected (-1, -2, -3), got (%v, %v, %v)", iv1.X, iv1.Y, iv1.Z)
	}

	iv2 := v2.Negated()
	if iv2.X != -3 || iv2.Y != -4 || iv2.Z != -5 {
		t.Errorf("Inverse failed. Expected (-3, -4, -5), got (%v, %v, %v)", iv2.X, iv2.Y, iv2.Z)
	}
}

func TestMVec4Abs(t *testing.T) {
	v1 := MVec4{-1, -2, -3, -4}
	expected1 := MVec4{1, 2, 3, 4}
	result1 := v1.Abs()
	if !result1.Equals(&expected1) {
		t.Errorf("Abs failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec4{3, -4, 5, -6}
	expected2 := MVec4{3, 4, 5, 6}
	result2 := v2.Abs()
	if !result2.Equals(&expected2) {
		t.Errorf("Abs failed. Expected %v, got %v", expected2, result2)
	}

	v3 := MVec4{0, 0, 0, 0}
	expected3 := MVec4{0, 0, 0, 0}
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

