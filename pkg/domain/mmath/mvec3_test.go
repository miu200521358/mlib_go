package mmath

import (
	"math"
	"testing"
)

func TestMVec3Interpolate(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := &MVec3{4, 5, 6} // Pass the address of v2
	t1 := 0.5
	expected := MVec3{2.5, 3.5, 4.5}

	result := v1.Interpolate(v2, t1) // Use v2 as a pointer

	if !result.NearEquals(&expected, 1e-8) {
		t.Errorf("Interpolation failed. Expected %v, got %v", expected, result)
	}
}

func TestMVec3NearEquals(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := MVec3{1.000001, 2.000001, 3.000001}
	epsilon := 0.00001

	if !v1.NearEquals(&v2, epsilon) {
		t.Errorf("NearEquals failed. Expected true, got false")
	}

	v3 := MVec3{1, 2, 3}
	v4 := MVec3{1.0001, 2.0001, 3.0001}

	if v3.NearEquals(&v4, epsilon) {
		t.Errorf("NearEquals failed. Expected false, got true")
	}
}

func TestMVec3LessThan(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := MVec3{4, 5, 6}

	if !v1.LessThan(&v2) {
		t.Errorf("LessThan failed. Expected true, got false")
	}

	v3 := MVec3{3, 4, 5}
	v4 := MVec3{1, 2, 3}

	if v3.LessThan(&v4) {
		t.Errorf("LessThan failed. Expected false, got true")
	}
}

func TestMVec3LessThanOrEquals(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := MVec3{3, 4, 5}

	if !v1.LessThanOrEquals(&v2) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}

	v3 := MVec3{3, 4, 5}
	v4 := MVec3{1, 2, 3}

	if v3.LessThanOrEquals(&v4) {
		t.Errorf("LessThanOrEqual failed. Expected false, got true")
	}

	v5 := MVec3{1, 2, 3}
	v6 := MVec3{1, 2, 3}

	if !v5.LessThanOrEquals(&v6) {
		t.Errorf("LessThanOrEqual failed. Expected true, got false")
	}
}

func TestMVec3GreaterThan(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := MVec3{3, 4, 5}

	if v1.GreaterThan(&v2) {
		t.Errorf("GreaterThan failed. Expected false, got true")
	}

	v3 := MVec3{3, 4, 5}
	v4 := MVec3{1, 2, 3}

	if !v3.GreaterThan(&v4) {
		t.Errorf("GreaterThan failed. Expected true, got false")
	}
}

func TestMVec3GreaterThanOrEquals(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := MVec3{3, 4, 5}

	if v1.GreaterThanOrEquals(&v2) {
		t.Errorf("GreaterThanOrEqual failed. Expected false, got true")
	}

	v3 := MVec3{3, 4, 5}
	v4 := MVec3{1, 2, 3}

	if !v3.GreaterThanOrEquals(&v4) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}

	v5 := MVec3{1, 2, 3}
	v6 := MVec3{1, 2, 3}

	if !v5.GreaterThanOrEquals(&v6) {
		t.Errorf("GreaterThanOrEqual failed. Expected true, got false")
	}
}

func TestMVec3Inverted(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := MVec3{3, 4, 5}

	iv1 := v1.Inverted()
	if iv1.X != -1 || iv1.Y != -2 || iv1.Z != -3 {
		t.Errorf("Inverse failed. Expected (-1, -2, -3), got (%v, %v, %v)", iv1.X, iv1.Y, iv1.Z)
	}

	iv2 := v2.Inverted()
	if iv2.X != -3 || iv2.Y != -4 || iv2.Z != -5 {
		t.Errorf("Inverse failed. Expected (-3, -4, -5), got (%v, %v, %v)", iv2.X, iv2.Y, iv2.Z)
	}
}

func TestMVec3Abs(t *testing.T) {
	v1 := MVec3{-1, -2, -3}
	expected1 := MVec3{1, 2, 3}
	result1 := v1.Abs()
	if !result1.Equals(&expected1) {
		t.Errorf("Abs failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec3{3, -4, 5}
	expected2 := MVec3{3, 4, 5}
	result2 := v2.Abs()
	if !result2.Equals(&expected2) {
		t.Errorf("Abs failed. Expected %v, got %v", expected2, result2)
	}

	v3 := MVec3{0, 0, 0}
	expected3 := MVec3{0, 0, 0}
	result3 := v3.Abs()
	if !result3.Equals(&expected3) {
		t.Errorf("Abs failed. Expected %v, got %v", expected3, result3)
	}
}

func TestMVec3Hash(t *testing.T) {
	v := MVec3{1, 2, 3}
	expected := uint64(17648364615301650315)
	result := v.Hash()
	if result != expected {
		t.Errorf("Hash failed. Expected %v, got %v", expected, result)
	}
}

func TestMVec3Angle(t *testing.T) {
	v1 := MVec3{1, 0, 0}
	v2 := MVec3{0, 1, 0}
	expected := math.Pi / 2
	result := v1.Angle(&v2)
	if result != expected {
		t.Errorf("Angle failed. Expected %v, got %v", expected, result)
	}

	v3 := MVec3{1, 1, 1}
	expected2 := 0.9553166181245092
	result2 := v1.Angle(&v3)
	if result2 != expected2 {
		t.Errorf("Angle failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMVec3Degree(t *testing.T) {
	v1 := MVec3{1, 0, 0}
	v2 := MVec3{0, 1, 0}
	expected := 90.0
	result := v1.Degree(&v2)
	if math.Abs(result-expected) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected, result)
	}

	v3 := MVec3{1, 1, 1}
	expected2 := 54.735610317245346
	result2 := v1.Degree(&v3)
	if math.Abs(result2-expected2) > 0.00001 {
		t.Errorf("Degree failed. Expected %v, got %v", expected2, result2)
	}
}
func TestStdMean(t *testing.T) {
	values := []MVec3{
		{1, 2, 3},
		{1.5, 1.2, 20.3},
		{1.8, 0.3, 1.3},
		{15, 0.2, 1.3},
		{1.3, 2.2, 2.3},
	}

	err := 1.5
	expected := MVec3{1.36666667, 1.5, 2.2}

	result := StdMeanVec3(values, err)

	if !result.NearEquals(&expected, 1e-8) {
		t.Errorf("StdMean failed. Expected %v, got %v", expected, result)
	}
}
func TestMVec3One(t *testing.T) {
	v1 := MVec3{1, 2, 3.2}
	expected1 := MVec3{1, 2, 3.2}
	result1 := v1.One()
	if !result1.NearEquals(&expected1, 1e-8) {
		t.Errorf("One failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec3{0, 2, 3.2}
	expected2 := MVec3{1, 2, 3.2}
	result2 := v2.One()
	if !result2.NearEquals(&expected2, 1e-8) {
		t.Errorf("One failed. Expected %v, got %v", expected2, result2)
	}

	v3 := MVec3{1, 0, 3.2}
	expected3 := MVec3{1, 1, 3.2}
	result3 := v3.One()
	if !result3.NearEquals(&expected3, 1e-8) {
		t.Errorf("One failed. Expected %v, got %v", expected3, result3)
	}

	v4 := MVec3{2, 0, 0}
	expected4 := MVec3{2, 1, 1}
	result4 := v4.One()
	if !result4.NearEquals(&expected4, 1e-8) {
		t.Errorf("One failed. Expected %v, got %v", expected4, result4)
	}
}

func TestMVec3Length(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	expected1 := 3.7416573867739413
	result1 := v1.Length()
	if math.Abs(result1-expected1) > 1e-10 {
		t.Errorf("Length failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec3{2.3, 0.2, 9}
	expected2 := 9.291393867445294
	result2 := v2.Length()
	if math.Abs(result2-expected2) > 1e-10 {
		t.Errorf("Length failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMVec3LengthSqr(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	expected1 := 14.0
	result1 := v1.LengthSqr()
	if math.Abs(result1-expected1) > 1e-10 {
		t.Errorf("LengthSqr failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec3{2.3, 0.2, 9}
	expected2 := 86.33000000000001
	result2 := v2.LengthSqr()
	if math.Abs(result2-expected2) > 1e-10 {
		t.Errorf("LengthSqr failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMVec3Normalized(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	expected1 := MVec3{0.2672612419124244, 0.5345224838248488, 0.8017837257372732}
	result1 := v1.Normalized()
	if !result1.NearEquals(&expected1, 1e-8) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected1, result1)
	}

	v2 := MVec3{2.3, 0.2, 9}
	expected2 := MVec3{0.24754089997827142, 0.021525295650284472, 0.9686383042628013}
	result2 := v2.Normalized()
	if !result2.NearEquals(&expected2, 1e-8) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected2, result2)
	}
}
func TestMVec3Distance(t *testing.T) {
	v1 := MVec3{1, 2, 3}
	v2 := MVec3{2.3, 0.2, 9}
	expected1 := 6.397655820689325
	result1 := v1.Distance(&v2)
	if math.Abs(result1-expected1) > 1e-10 {
		t.Errorf("Distance failed. Expected %v, got %v", expected1, result1)
	}

	v3 := MVec3{-1, -0.3, 3}
	v4 := MVec3{-2.3, 0.2, 9.33333333}
	expected2 := 6.484682804030502
	result2 := v3.Distance(&v4)
	if math.Abs(result2-expected2) > 1e-10 {
		t.Errorf("Distance failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMVector3DGetLocalMatrix(t *testing.T) {
	v1 := MVec3{0.8, 0.6, 1}
	localMatrix := v1.ToLocalMatrix4x4()

	expected1 := MMat4{
		0.56568542, 0.6, 0.56568542, 0.,
		0.42426407, -0.8, 0.42426407, 0.,
		0.70710678, 0., -0.70710678, 0.,
		0., 0., 0., 1.,
	}

	if !localMatrix.NearEquals(&expected1, 1e-8) {
		t.Errorf("Local matrix calculation failed. Expected %v, got %v", expected1, localMatrix)
	}

	v2 := MVec3{1, 0, 0}
	localVector1 := localMatrix.MulVec3(&v2)

	expected2 := &MVec3{0.56568542, 0.42426407, 0.70710678}
	if !localVector1.NearEquals(expected2, 1e-8) {
		t.Errorf("Local vector calculation failed. Expected %v, got %v", expected2, localVector1)
	}

	v3 := MVec3{1, 0, 1}
	localVector2 := localMatrix.MulVec3(&v3)

	expected3 := MVec3{1.13137085, 0.848528137, -1.11022302e-16}
	if !localVector2.NearEquals(&expected3, 1e-8) {
		t.Errorf("Local vector calculation failed. Expected %v, got %v", expected3, localVector2)
	}

	v4 := MVec3{0, 0, -0.5}
	localMatrix2 := v4.ToLocalMatrix4x4()

	expected4 := MMat4{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	if !localMatrix2.NearEquals(&expected4, 1e-8) {
		t.Errorf("Local matrix calculation failed. Expected %v, got %v", expected4, localMatrix2)
	}
}

func TestGetVertexLocalPositions(t *testing.T) {
	vertexLocalPositions := GetVertexLocalPositions(
		[]*MVec3{{1, 0, 0}, {0.5, 3, 2}, {-1, -2, 3}},
		&MVec3{0.5, 0.5, 1},
		&MVec3{0.7, 2, 1.5},
	)

	expected := []*MVec3{
		{0.10944881889763777, 0.8208661417322836, 0.2736220472440945},
		{0.5346456692913384, 4.009842519685039, 1.3366141732283463},
		{-0.04015748031496058, -0.3011811023622044, -0.10039370078740151},
	}

	for i := 0; i < 3; i++ {
		if vertexLocalPositions[i].NearEquals(expected[i], 1e-8) {
			t.Errorf("GetVertexLocalPositions failed. Expected %v, got %v", expected, vertexLocalPositions)
			break
		}
	}
}
