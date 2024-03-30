package mmath

import (
	"testing"
)

func TestMMat4_Translate(t *testing.T) {
	mat := &MMat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	v := &MVec3{1, 2, 3}
	expectedMat := &MMat4{
		{1, 0, 0, 1},
		{0, 1, 0, 2},
		{0, 0, 1, 3},
		{0, 0, 0, 1},
	}

	mat.Translate(v)

	// Verify the matrix values
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if mat[i][j] != expectedMat[i][j] {
				t.Errorf("Expected mat[%d][%d] to be %f, got %f", i, j, expectedMat[i][j], mat[i][j])
			}
		}
	}
}

func TestMMat4_Quaternion(t *testing.T) {
	mat := &MMat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	expectedQ := NewMQuaternionByValues(0, 0, 0, 1)

	q := mat.Quaternion()

	// Verify the quaternion values
	if !q.PracticallyEquals(expectedQ, 1e-10) {
		t.Errorf("Expected q to be %v, got %v", expectedQ, q)
	}
}

func TestMMat4_AssignEulerRotation(t *testing.T) {
	mat := &MMat4{}
	xPitch := 0.5
	yHead := 0.3
	zRoll := 0.2

	expectedMat := &MMat4{
		{0.90815, -0.17435, 0.38062, 0.00000},
		{0.32865, 0.86009, -0.39017, 0.00000},
		{-0.25934, 0.47943, 0.83839, 0.00000},
		{0.00000, 0.00000, 0.00000, 1.00000},
	}

	mat.AssignEulerRotation(xPitch, yHead, zRoll)

	if !mat.PracticallyEquals(expectedMat, 1e-5) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat, mat)
	}
}

func TestMMat4_Mul(t *testing.T) {
	mat := &MMat4{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}

	a := &MMat4{
		{17, 18, 19, 20},
		{21, 22, 23, 24},
		{25, 26, 27, 28},
		{29, 30, 31, 32},
	}

	expectedMat := &MMat4{
		{250, 260, 270, 280},
		{618, 644, 670, 696},
		{986, 1028, 1070, 1112},
		{1354, 1412, 1470, 1528},
	}

	mat.Mul(a)

	// Verify the matrix values
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if mat[i][j] != expectedMat[i][j] {
				t.Errorf("Expected mat[%d][%d] to be %f, got %f", i, j, expectedMat[i][j], mat[i][j])
			}
		}
	}
}

func TestMMat4_Translation(t *testing.T) {
	mat := &MMat4{
		{1, 0, 0, 1},
		{0, 1, 0, 2},
		{0, 0, 1, 3},
		{0, 0, 0, 1},
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
		{-0.28213944, 0.48809647, 0.82592928, 0.0},
		{0.69636424, 0.69636424, -0.17364818, 0.0},
		{-0.65990468, 0.52615461, -0.53636474, 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}

	expected1 := MMat4{
		{-0.28213944, 0.69636424, -0.65990468, 0.0},
		{0.48809647, 0.69636424, 0.52615461, 0.0},
		{0.82592928, -0.17364818, -0.53636474, 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}

	result1 := mat1.Inverse()

	// Verify the matrix values
	if !result1.PracticallyEquals(&expected1, 1e-8) {
		t.Errorf("Expected inverse matrix to be %v, got %v", expected1, result1)
	}

	mat2 := &MMat4{
		{0.45487413, 0.87398231, -0.17101007, 0.0},
		{-0.49240388, 0.08682409, -0.8660254, 0.0},
		{-0.74204309, 0.47813857, 0.46984631, 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}

	expected2 := MMat4{
		{0.45487413, -0.49240388, -0.74204309, 0.0},
		{0.87398231, 0.08682409, 0.47813857, 0.0},
		{-0.17101007, -0.8660254, 0.46984631, 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}

	result2 := mat2.Inverse()

	// Verify the matrix values
	if !result2.PracticallyEquals(&expected2, 1e-8) {
		t.Errorf("Expected inverse matrix to be %v, got %v", expected2, result2)
	}
}

func TestMMat4Mul(t *testing.T) {
	mat1 := &MMat4{
		{-0.28213944, 0.48809647, 0.82592928, 0.0},
		{0.69636424, 0.69636424, -0.17364818, 0.0},
		{-0.65990468, 0.52615461, -0.53636474, 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}

	expected1 := MMat4{
		{-0.13337049, 0.79765275, 0.58818569, 0.0},
		{0.98098506, 0.19068573, -0.03615618, 0.0},
		{-0.14099869, 0.5721792, -0.80791727, 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}

	mat2 := &MMat4{
		{0.81379768, -0.46984631, 0.34202014, 0.},
		{0.54383814, 0.82317294, -0.16317591, 0.},
		{-0.20487413, 0.31879578, 0.92541658, 0.},
		{0., 0., 0., 1.},
	}

	mat1.Mul(mat2)

	if !mat1.PracticallyEquals(&expected1, 1e-8) {
		t.Errorf("Expected matrix to be %v, got %v", expected1, mat1)
	}

	expected3 := MMat4{
		{-0.50913704, -0.04344393, 0.85958833, 0.0},
		{0.66451052, 0.61488424, 0.42466828, 0.0},
		{-0.54699658, 0.78741983, -0.28419139, 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}

	mat3 := &MMat4{
		{0.79690454, 0.49796122, 0.34202014, 0.},
		{-0.59238195, 0.53314174, 0.60402277, 0.},
		{0.11843471, -0.68395505, 0.71984631, 0.},
		{0., 0., 0., 1.},
	}

	mat1.Mul(mat3)

	if !mat1.PracticallyEquals(&expected3, 1e-8) {
		t.Errorf("Expected matrix to be %v, got %v", expected3, mat1)
	}
}
