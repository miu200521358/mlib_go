package mmath

import (
	"testing"
)

func TestMMat4_NearEquals(t *testing.T) {
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

	if !mat1.NearEquals(mat2, 1e-5) {
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

	if mat3.NearEquals(mat4, 1e-5) {
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
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		1, 2, 3, 1,
	}

	mat.Translate(v)

	// Verify the matrix values
	if !mat.NearEquals(expectedMat, 1e-10) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat, mat)
	}

	{
		m := NewMMat4ByValues(
			-0.28213944, 0.48809647, 0.82592928, 0.0,
			0.69636424, 0.69636424, -0.17364818, 0.0,
			-0.65990468, 0.52615461, -0.53636474, 0.0,
			0.0, 0.0, 0.0, 1.0,
		)
		m.Translate(&MVec3{10, 20, 30})

		expectedMat := NewMMat4ByValues(
			-0.28213944, 0.48809647, 0.82592928, 0.0,
			0.69636424, 0.69636424, -0.17364818, 0.0,
			-0.65990468, 0.52615461, -0.53636474, 0.0,
			10.0, 20.0, 30.0, 1.0,
		)

		if !m.NearEquals(expectedMat, 1e-5) {
			t.Errorf("Expected mat to be %v, got %v", expectedMat, m)
		}

		m.Translate(&MVec3{-8, -12, 3})

		expectedMat2 := NewMMat4ByValues(
			-0.28213944, 0.48809647, 0.82592928, 0.0,
			0.69636424, 0.69636424, -0.17364818, 0.0,
			-0.65990468, 0.52615461, -0.53636474, 0.0,
			2.0, 8.0, 33.0, 1.0,
		)

		if !m.NearEquals(expectedMat2, 1e-5) {
			t.Errorf("Expected mat to be %v, got %v", expectedMat2, m)
		}
	}
}

func TestMMat4_Rotate(t *testing.T) {
	m := &MMat4{
		-0.28213944, 0.48809647, 0.82592928, 0.0,
		0.69636424, 0.69636424, -0.17364818, 0.0,
		-0.65990468, 0.52615461, -0.53636474, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
	q := NewMQuaternionFromDegrees(10, 20, 30)
	expectedQ := NewMQuaternionByValues(0.127679440695781, 0.144878125417369, 0.2685358227515692, 0.943714364147489)

	if !q.NearEquals(expectedQ, 1e-10) {
		t.Errorf("Expected q to be %v, got %v", expectedQ, q)
	}

	m.Rotate(q)

	expectedMat := &MMat4{
		-0.1764503, 0.11357786, 0.97773481, 0.0,
		0.18012426, 0.98027284, -0.08136594, 0.0,
		-0.96768825, 0.16175671, -0.19342756, 0.0,
		0., 0., 0., 1.,
	}

	if !m.NearEquals(expectedMat, 1e-5) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat, m)
	}

	q2 := NewMQuaternionFromDegrees(-40, 20, -32)

	expectedMat2 := &MMat4{
		0.250348, 0.755653, 0.605239, 0.0,
		0.603851, 0.366775, -0.707701, 0.0,
		-0.756763, 0.542646, -0.364480, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	m.Rotate(q2)

	if !m.NearEquals(expectedMat2, 1e-5) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat2, m)
	}
}

func TestMMat4_Scale(t *testing.T) {
	m := NewMMat4()
	m.Rotate(NewMQuaternionFromDegrees(10, 20, 30))

	expectedMat := &MMat4{
		0.81379768, 0.54383814, -0.20487413, 0.0,
		-0.46984631, 0.82317294, 0.31879578, 0.0,
		0.34202014, -0.16317591, 0.92541658, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	if !m.NearEquals(expectedMat, 1e-5) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat, m)
	}

	m.Translate(&MVec3{1, 2, 3})

	expectedMat2 := NewMMat4ByValues(
		0.81379768, 0.54383814, -0.20487413, 0.0,
		-0.46984631, 0.82317294, 0.31879578, 0.0,
		0.34202014, -0.16317591, 0.92541658, 0.0,
		1.0, 2.0, 3.0, 1.0,
	)

	if !m.NearEquals(expectedMat2, 1e-5) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat2, m)
	}

	m.Scale(&MVec3{4, 5, 6})

	expectedMat3 := NewMMat4ByValues(
		3.25519072, 2.7191907, -1.22924478, 0.0,
		-1.87938524, 4.1158647, 1.91277468, 0.0,
		1.36808056, -0.81587955, 5.55249948, 0.0,
		4.0, 10.0, 18.0, 1.0,
	)

	if !m.NearEquals(expectedMat3, 1e-5) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat3, m)
	}

	m.Scale(&MVec3{-0.8, -0.12, 0.3})

	expectedMat4 := NewMMat4ByValues(
		-2.60415258, -0.32630288, -0.36877343, 0.,
		1.50350819, -0.49390376, 0.5738324, 0.,
		-1.09446445, 0.09790555, 1.66574984, 0.,
		-3.2, -1.2, 5.4, 1.,
	)

	if !m.NearEquals(expectedMat4, 1e-5) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat4, m)
	}
}

func TestMMat4_Quaternion(t *testing.T) {
	{
		mat := &MMat4{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		}
		expectedQ := NewMQuaternionByValues(0, 0, 0, 1)

		q := mat.Quaternion().Shorten()

		// Verify the quaternion values
		if !q.NearEquals(expectedQ, 1e-10) {
			t.Errorf("Expected q to be %v, got %v", expectedQ, q)
		}
	}

	{
		expectedQ := NewMQuaternionByValues(
			-0.37345048680206483, -0.7929168335960588, -0.11114230870887896, 0.46847093448041827)
		m := NewMMat4ByValues(
			-0.28213944, 0.48809647, 0.82592928, 0.0,
			0.69636424, 0.69636424, -0.17364818, 0.0,
			-0.65990468, 0.52615461, -0.53636474, 0.0,
			0.0, 0.0, 0.0, 1.0,
		)
		q := m.Quaternion().Shorten()

		if !q.NearEquals(expectedQ, 1e-5) {
			t.Errorf("Expected q to be %v, got %v", expectedQ, q)
		}
	}

	{
		expectedQ := NewMQuaternionByValues(
			-0.12767944069578063, -0.14487812541736916, -0.2392983377447303, 0.9515485246437885)
		m := NewMMat4ByValues(
			0.84349327, -0.41841204, 0.33682409, 0.0,
			0.49240388, 0.85286853, -0.17364818, 0.0,
			-0.21461018, 0.31232456, 0.92541658, 0.0,
			0.0, 0.0, 0.0, 1.0)
		q := m.Quaternion().Shorten()

		if !q.NearEquals(expectedQ, 1e-5) {
			t.Errorf("Expected q to be %v, got %v", expectedQ, q)
		}
	}
}

func TestMMat4_Mul(t *testing.T) {
	a := &MMat4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	mat := &MMat4{
		17, 18, 19, 20,
		21, 22, 23, 24,
		25, 26, 27, 28,
		29, 30, 31, 32,
	}

	expectedMat := &MMat4{
		250, 260, 270, 280,
		618, 644, 670, 696,
		986, 1028, 1070, 1112,
		1354, 1412, 1470, 1528,
	}
	mat.Mul(a)

	// Verify the matrix values
	if !mat.NearEquals(expectedMat, 1e-10) {
		t.Errorf("Expected mat to be %v, got %v", expectedMat, mat)
	}
}

func TestMMat4_Translation(t *testing.T) {
	mat := &MMat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		1, 2, 3, 1,
	}

	expectedVec := MVec3{1, 2, 3}

	result := mat.Translation()

	// Verify the vector values
	if !result.NearEquals(&expectedVec, 1e-5) {
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
	if !result1.NearEquals(&expected1, 1e-5) {
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
	if !result2.NearEquals(&expected2, 1e-5) {
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

	if !mat1.NearEquals(&expected1, 1e-5) {
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

	if !mat1.NearEquals(&expected3, 1e-5) {
		t.Errorf("Expected matrix to be %v, got %v", expected3, mat1)
	}
}

func TestMMat4_MulVec3(t *testing.T) {
	mat := &MMat4{
		-0.28213944, 0.69636424, -0.65990468, 0,
		0.48809647, 0.69636424, 0.52615461, 0,
		0.82592928, -0.17364818, -0.53636474, 0,
		0., 0., 0., 1.,
	}

	v := &MVec3{10, 20, 30}

	expectedVec := MVec3{31.7184134, 15.6814818, -12.166896800000002}

	result := mat.MulVec3(v)

	// Verify the vector values
	if !result.NearEquals(&expectedVec, 1e-5) {
		t.Errorf("Expected vector to be %v, got %v", expectedVec, result)
	}
}
