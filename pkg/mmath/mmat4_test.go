package mmath

import (
	"testing"
)

func TestMMat4_PracticallyEquals(t *testing.T) {
	mat1 := &MMat4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	mat2 := &MMat4{
		1.00000001, 2.00000001, 3.00000001, 4.00000001,
		5.00000001, 6.00000001, 7.00000001, 8.00000001,
		9.00000001, 10.00000001, 11.00000001, 12.00000001,
		13.00000001, 14.00000001, 15.00000001, 16.00000001,
	}

	if !mat1.PracticallyEquals(mat2, 1e-8) {
		t.Errorf("Expected mat1 to be practically equal to mat2")
	}

	mat3 := &MMat4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	mat4 := &MMat4{
		1.0001, 2.0001, 3.0001, 4.0001,
		5.0001, 6.0001, 7.0001, 8.0001,
		9.0001, 10.0001, 11.0001, 12.0001,
		13.0001, 14.0001, 15.0001, 16.0001,
	}

	if mat3.PracticallyEquals(mat4, 1e-8) {
		t.Errorf("Expected mat3 to not be practically equal to mat4")
	}
}

func TestMMat4_Translate(t *testing.T) {
	mat := &MMat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	v := &MVec3{1, 2, 3}
	expectedMat := &MMat4{
		1, 0, 0, 1,
		0, 1, 0, 2,
		0, 0, 1, 3,
		0, 0, 0, 1,
	}

	mat.Translate(v)

	// Verify the matrix values
	if !mat.PracticallyEquals(expectedMat, 1e-10) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat, mat)
	}
}

func TestMMat4_Quaternion(t *testing.T) {
	mat := &MMat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	expectedQ := NewMQuaternionByValues(0, 0, 0, 1)

	q := mat.Quaternion()

	// Verify the quaternion values
	if !q.PracticallyEquals(expectedQ, 1e-10) {
		t.Errorf("Expected q to be %v, got %v", expectedQ, q)
	}
}

func TestMMat4_Mul(t *testing.T) {
	mat := &MMat4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	a := &MMat4{
		17, 18, 19, 20,
		21, 22, 23, 24,
		25, 26, 27, 28,
		29, 30, 31, 32,
	}

	expectedMat := &MMat4{
		394, 564, 734, 520,
		458, 644, 830, 720,
		522, 724, 926, 920,
		1222, 1406, 1590, 1528,
	}
	mat.Mul(a)

	// Verify the matrix values
	if !mat.PracticallyEquals(expectedMat, 1e-10) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat, mat)
	}
}

func TestMMat4_Translation(t *testing.T) {
	mat := &MMat4{
		1, 0, 0, 1,
		0, 1, 0, 2,
		0, 0, 1, 3,
		0, 0, 0, 1,
	}

	expectedVec := MVec3{1, 2, 3}

	result := mat.Translation()

	// Verify the vector values
	if !result.PracticallyEquals(&expectedVec, 1e-8) {
		t.Errorf("Expected translation to be %v, got %v", expectedVec, result)
	}
}

func TestMMat4_Inverse(t *testing.T) {
	mat1 := &MMat4{
		-0.28213944, 0.48809647, 0.82592928, 0.0,
		0.69636424, 0.69636424, -0.17364818, 0.0,
		-0.65990468, 0.52615461, -0.53636474, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	expected1 := MMat4{
		-0.28213944, 0.69636424, -0.65990468, 0.0,
		0.48809647, 0.69636424, 0.52615461, 0.0,
		0.82592928, -0.17364818, -0.53636474, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	result1 := mat1.Inverse()

	// Verify the matrix values
	if !result1.PracticallyEquals(&expected1, 1e-8) {
		t.Errorf("Expected inverse matrix to be %v, got %v", expected1, result1)
	}

	mat2 := &MMat4{
		0.45487413, 0.87398231, -0.17101007, 0.0,
		-0.49240388, 0.08682409, -0.8660254, 0.0,
		-0.74204309, 0.47813857, 0.46984631, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	expected2 := MMat4{
		0.45487413, -0.49240388, -0.74204309, 0.0,
		0.87398231, 0.08682409, 0.47813857, 0.0,
		-0.17101007, -0.8660254, 0.46984631, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	result2 := mat2.Inverse()

	// Verify the matrix values
	if !result2.PracticallyEquals(&expected2, 1e-8) {
		t.Errorf("Expected inverse matrix to be %v, got %v", expected2, result2)
	}
}

func TestMMat4Mul(t *testing.T) {
	mat1 := &MMat4{
		-0.28213944, 0.48809647, 0.82592928, 0.0,
		0.69636424, 0.69636424, -0.17364818, 0.0,
		-0.65990468, 0.52615461, -0.53636474, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	mat2 := &MMat4{
		0.81379768, -0.46984631, 0.34202014, 0.,
		0.54383814, 0.82317294, -0.16317591, 0.,
		-0.20487413, 0.31879578, 0.92541658, 0.,
		0., 0., 0., 1.,
	}

	mat1.Mul(mat2)

	expected1 := MMat4{
		-0.7824892813287088, 0.24998307969608058, 0.5702797450534226, 0,
		0.5274705571536827, 0.7528179178495863, 0.3937511650919034, 0,
		-0.330885678728, 0.6089118411450198, -0.720930672993596, 0,
		0, 0, 0, 1,
	}

	if !mat1.PracticallyEquals(&expected1, 1e-8) {
		t.Errorf("Expected matrix to be %v, got %v", expected1, mat1)
	}

	mat3 := &MMat4{
		0.79690454, 0.49796122, 0.34202014, 0.,
		-0.59238195, 0.53314174, 0.60402277, 0.,
		0.11843471, -0.68395505, 0.71984631, 0.,
		0., 0., 0., 1.,
	}

	mat1.Mul(mat3)

	expected3 := MMat4{
		-0.47407894480040325, 0.7823468930993056, 0.40395851874113675, 0,
		0.5448866127486663, 0.6210698073823508, -0.5633567882156808, 0,
		-0.6916268772680453, -0.046964001131448774, -0.7207264663039717, 0,
		0, 0, 0, 1,
	}

	if !mat1.PracticallyEquals(&expected3, 1e-8) {
		t.Errorf("Expected matrix to be %v, got %v", expected3, mat1)
	}
}
