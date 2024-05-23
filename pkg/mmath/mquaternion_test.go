package mmath

import (
	"math"
	"testing"
)

func TestMQuaternionToEulerAngles(t *testing.T) {
	quat := NewMQuaternionByValues(0.7071067811865476, 0.0, 0.0, 0.7071067811865476)
	expected := MVec3{1.5707963267948966, 0, 0}

	result := quat.ToRadians()

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 {
		t.Errorf("ToEulerAngles failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToDegree(t *testing.T) {
	quat := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)
	expected := 10.0

	result := quat.ToDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToDegree2(t *testing.T) {
	quat := NewMQuaternionByValues(0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885)
	expected := 35.81710117358426

	result := quat.ToDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToSignedDegree(t *testing.T) {
	quat := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)
	expected := 10.0

	result := quat.ToSignedDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToSignedDegree2(t *testing.T) {
	quat := NewMQuaternionByValues(0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844)
	expected := 89.66927179998277

	result := quat.ToSignedDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionDot(t *testing.T) {
	// np.array([60, -20, -80]),
	quat1 := NewMQuaternionByValues(0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844)
	// np.array([10, 20, 30]),
	quat2 := NewMQuaternionByValues(0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885)
	expected := 0.6491836986795888

	result := quat1.Dot(quat2)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}

	// np.array([10, 23, 45]),
	quat3 := NewMQuaternionByValues(0.1549093965157679, 0.15080756177478563, 0.3575205710320892, 0.908536845412201)
	// np.array([12, 20, 42]),
	quat4 := NewMQuaternionByValues(0.15799222008931638, 0.1243359045760714, 0.33404459937562386, 0.9208654879256133)

	expected2 := 0.9992933154462645

	result2 := quat3.Dot(quat4)

	if math.Abs(result2-expected2) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMQuaternionSlerp(t *testing.T) {
	// np.array([60, -20, -80])
	quat1 := NewMQuaternionByValues(0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844)
	// np.array([10, 20, 30]),
	quat2 := NewMQuaternionByValues(0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885)
	tValue := 0.3
	expected := NewMQuaternionByValues(0.3973722198386427, 0.19936467087655246, -0.27953105525419597, 0.851006131620254)

	result := quat1.Slerp(quat2, tValue)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("Slerp failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToFixedAxisRotation(t *testing.T) {
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{1, 0, 0}
		expected := NewMQuaternionByValues(0.866025403784439, 0, 0, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.PracticallyEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{0, 1, 0}
		expected := NewMQuaternionByValues(0, 0.866025403784439, 0, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.PracticallyEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{0, 0, 1}
		expected := NewMQuaternionByValues(0, 0, 0.866025403784439, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.PracticallyEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{0.5, 0.7, 0.2}
		expected := NewMQuaternionByValues(0.49029033784546, 0.686406472983644, 0.196116135138184, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.PracticallyEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
}

func TestMQuaternionNormalized(t *testing.T) {
	quat1 := NewMQuaternionByValues(2, 3, 4, 1)
	expected1 := NewMQuaternionByValues(0.36514837, 0.54772256, 0.73029674, 0.18257419)

	result1 := quat1.Normalized()

	if !result1.PracticallyEquals(expected1, 1e-7) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected1, result1)
	}

	quat2 := NewMQuaternionByValues(0, 0, 0, 1)
	expected2 := NewMQuaternionByValues(0, 0, 0, 1)

	result2 := quat2.Normalized()

	if !result2.PracticallyEquals(expected2, 1e-10) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected2, result2)
	}

	quat3 := NewMQuaternion()
	expected3 := NewMQuaternionByValues(0, 0, 0, 1)

	result3 := quat3.Normalized()

	if !result3.PracticallyEquals(expected3, 1e-10) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected3, result3)
	}
}

func TestFromEulerAnglesDegrees(t *testing.T) {
	expected1 := NewMQuaternionByValues(0, 0, 0, 1)

	result1 := NewMQuaternionFromDegrees(0, 0, 0)

	if !result1.PracticallyEquals(expected1, 1e-8) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(0.08715574, 0.0, 0.0, 0.9961947)

	result2 := NewMQuaternionFromDegrees(10, 0, 0)

	if !result2.PracticallyEquals(expected2, 1e-6) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected2, result2)
	}

	expected3 := NewMQuaternionByValues(0.12767944069578063, 0.14487812541736914, 0.2685358227515692, 0.943714364147489)

	result3 := NewMQuaternionFromDegrees(10, 20, 30)

	if !result3.PracticallyEquals(expected3, 1e-5) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected3, result3)
	}

	expected4 := NewMQuaternionByValues(0.47386805375453483, 0.20131048764138493, -0.6147244358103234, 0.5975257510887351)

	result4 := NewMQuaternionFromDegrees(60, -20, -80)

	if !result4.PracticallyEquals(expected4, 1e-5) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected4, result4)
	}
}

func TestMQuaternionToEulerAnglesDegrees(t *testing.T) {
	expected1 := &MVec3{0, 0, 0}

	qq1 := NewMQuaternionByValues(0, 0, 0, 1)
	result1 := qq1.ToDegrees()

	if !result1.PracticallyEquals(expected1, 1e-8) {
		t.Errorf("ToEulerAnglesDegrees failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := &MVec3{10, 0, 0}

	qq2 := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)
	result2 := qq2.ToDegrees()

	if !result2.PracticallyEquals(expected2, 1e-5) {
		t.Errorf("ToEulerAnglesDegrees failed. Expected %v, got %v", expected2, result2)
	}

	expected3 := &MVec3{10, 20, 30}

	qq3 := NewMQuaternionByValues(0.12767944, 0.14487813, 0.23929834, 0.95154852)
	result3 := qq3.ToDegrees()

	if !result3.PracticallyEquals(expected3, 1e-5) {
		t.Errorf("ToEulerAnglesDegrees failed. Expected %v, got %v", expected3, result3)
	}

	expected4 := &MVec3{60, -20, -80}

	qq4 := NewMQuaternionByValues(0.47386805, 0.20131049, -0.48170221, 0.70914465)
	result4 := qq4.ToDegrees()

	if !result4.PracticallyEquals(expected4, 1e-5) {
		t.Errorf("ToEulerAnglesDegrees failed. Expected %v, got %v", expected4, result4)
	}
}

func TestMQuaternionMultiply(t *testing.T) {
	expected1 := NewMQuaternionByValues(
		0.6594130183457979, 0.11939693791117263, -0.24571599091322077, 0.7003873887093154)
	q11 := NewMQuaternionByValues(
		0.4738680537545347,
		0.20131048764138487,
		-0.48170221425083437,
		0.7091446481376844,
	)
	q12 := NewMQuaternionByValues(
		0.12767944069578063,
		0.14487812541736916,
		0.2392983377447303,
		0.9515485246437885,
	)
	result1 := q11.Mul(q12)

	if !result1.PracticallyEquals(expected1, 1e-8) {
		t.Errorf("MQuaternionMultiply failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(
		0.4234902605993554, 0.46919555165368526, -0.3316158006229952, 0.7003873887093154)
	q21 := NewMQuaternionByValues(
		0.12767944069578063,
		0.14487812541736916,
		0.2392983377447303,
		0.9515485246437885,
	)
	q22 := NewMQuaternionByValues(
		0.4738680537545347,
		0.20131048764138487,
		-0.48170221425083437,
		0.7091446481376844,
	)
	result2 := q21.Mul(q22)

	if !result2.PracticallyEquals(expected2, 1e-8) {
		t.Errorf("MQuaternionMultiply failed. Expected %v, got %v", expected2, result2)
	}
}

func TestNewMQuaternionFromAxisAngles(t *testing.T) {
	expected1 := NewMQuaternionByValues(
		0.25511557978461696, 0.5102311595692339, 0.7653467393538509, -0.29803451698791944)
	result1 := NewMQuaternionFromAxisAngles(&MVec3{1, 2, 3}, 30)

	if !result1.PracticallyEquals(expected1, 1e-5) {
		t.Errorf("NewMQuaternionFromAxisAngles failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(
		0.1329649118205802, -0.8864327454705346, 0.4432163727352673, 0.01079661620640226)
	result2 := NewMQuaternionFromAxisAngles(&MVec3{-3, 20, -10}, 123)

	if !result2.PracticallyEquals(expected2, 1e-5) {
		t.Errorf("NewMQuaternionFromAxisAngles failed. Expected %v, got %v", expected2, result2)
	}

	axis := MVec3{1, 0, 0}
	angle := math.Pi / 2
	expected := NewMQuaternionByValues(0.7071067811865476, 0, 0, 0.7071067811865476)

	result := NewMQuaternionFromAxisAngles(&axis, angle)

	if !result.PracticallyEquals(expected, 1e-10) {
		t.Errorf("NewMQuaternionFromAxisAngles failed. Expected %v, got %v", expected, result)
	}
}
func TestMQuaternionFromDirection(t *testing.T) {
	expected1 := NewMQuaternionByValues(
		0.3115472173245163, 0.04523791008340301, 0.5420603160713341, 0.7791421414666787)
	result1 := NewMQuaternionFromDirection(&MVec3{1, 2, 3}, &MVec3{4, 5, 6})

	if !result1.PracticallyEquals(expected1, 1e-7) {
		t.Errorf("MQuaternionFromDirection failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(
		0.543212292317204, 0.6953153333136457, 0.20212324833235548, 0.42497433477564167)
	result2 := NewMQuaternionFromDirection(&MVec3{-10, 20, -15}, &MVec3{40, -5, 6})

	if !result2.PracticallyEquals(expected2, 1e-7) {
		t.Errorf("MQuaternionFromDirection failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMQuaternionRotate(t *testing.T) {
	expected1 := NewMQuaternionByValues(
		-0.04597839511020707, 0.0919567902204141, -0.04597839511020706, 0.9936377222602503)
	result1 := NewMQuaternionRotate(&MVec3{1, 2, 3}, &MVec3{4, 5, 6})

	if !result1.PracticallyEquals(expected1, 1e-5) {
		t.Errorf("MQuaternionRotate failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(
		0.042643949239185255, -0.511727390870223, -0.7107324873197542, 0.48080755245182594)
	result2 := NewMQuaternionRotate(&MVec3{-10, 20, -15}, &MVec3{40, -5, 6})

	if !result2.PracticallyEquals(expected2, 1e-5) {
		t.Errorf("MQuaternionRotate failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMQuaternionToMatrix4x4(t *testing.T) {
	expected := &MMat4{
		0.45487412870286215, -0.49240387650610407, -0.7420430913623366, 0,
		0.8739823124221016, 0.08682408883346526, 0.47813857318935704, 0,
		-0.1710100716628345, -0.8660254037844384, 0.46984631039295444, 0,
		0, 0, 0, 1,
	}

	qq1 := NewMQuaternionByValues(
		0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844)
	result1 := qq1.ToMat4()

	if !result1.PracticallyEquals(expected, 1e-5) {
		t.Errorf("ToMatrix4x4 failed. Expected %v, got %v", expected, qq1)
	}

	expected2 := &MMat4{
		-0.28213943766659333, 0.6963642403200191, -0.6599046768410582, 0,
		0.4880964722494531, 0.6963642403200189, 0.526154614710507, 0,
		0.8259292775358055, -0.17364817766693036, -0.5363647442758863, 0,
		0, 0, 0, 1,
	}

	// np.array([10, 123, 45])
	qq2 := NewMQuaternionByValues(
		0.3734504874442106, 0.7929168339527322, 0.11114231087966482, 0.4684709324967611)
	result2 := qq2.ToMat4()

	if !result2.PracticallyEquals(expected2, 1e-5) {
		t.Errorf("ToMatrix4x4 failed. Expected %v, got %v", expected, qq1)
	}
}

func TestMQuaternionMulVec3(t *testing.T) {
	expected := &MVec3{16.89808539, -29.1683191, 16.23772986}
	//  np.array([60, -20, -80]),
	qq := NewMQuaternionByValues(
		0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844)
	result := qq.MulVec3(&MVec3{10, 20, 30})

	if !result.PracticallyEquals(expected, 1e-5) {
		t.Errorf("MulVec3 failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionVectorToDegree(t *testing.T) {
	expected := 81.78678929826181

	result := VectorToDegree(&MVec3{10, 20, 30}, &MVec3{30, -20, 10})

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("VectorToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionMulFactor(t *testing.T) {
	{
		quat := NewMQuaternionFromDegrees(90, 0, 0)
		factor := 0.5
		expected := NewMQuaternionFromDegrees(45, 0, 0)

		result := quat.MulScalar(factor)

		if !result.PracticallyEquals(expected, 1e-10) {
			t.Errorf("MulFactor failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionFromDegrees(-24.53194, 180, 180)
		factor := 0.5
		expected := NewMQuaternionFromDegrees(77.73403, 0, 0)

		result := quat.MulScalar(factor)

		if !result.PracticallyEquals(expected, 1e-6) {
			t.Errorf("MulFactor failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionFromDegrees(-24.53194, 180, 180)
		factor := -0.5
		expected := NewMQuaternionFromDegrees(-77.73403, 0, 0)

		result := quat.MulScalar(factor)

		if !result.PracticallyEquals(expected, 1e-6) {
			t.Errorf("MulFactor failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionByValues(0, 0, 0, 1)
		factor := 0.5
		expected := NewMQuaternionByValues(0, 0, 0, 1)

		result := quat.MulScalar(factor)

		if !result.PracticallyEquals(expected, 1e-10) {
			t.Errorf("MulFactor failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)
		factor := 1.0
		expected := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)

		result := quat.MulScalar(factor)

		if !result.PracticallyEquals(expected, 1e-10) {
			t.Errorf("MulFactor failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}
}

func TestMQuaternion_ToMat4(t *testing.T) {
	// Test case 1: Identity quaternion
	quat := NewMQuaternionByValues(0, 0, 0, 1).Normalize()
	expectedMat := NewMMat4ByValues(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	actualMat := quat.ToMat4()

	if !actualMat.PracticallyEquals(expectedMat, 1e-10) {
		t.Errorf("Test case 1 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

	// Test case 2: Non-identity quaternion
	quat = NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
	expectedMat = NewMMat4ByValues(0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1)
	actualMat = quat.ToMat4()

	if !actualMat.PracticallyEquals(expectedMat, 1e-10) {
		t.Errorf("Test case 2 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

	if !actualMat.Quaternion().PracticallyEquals(quat.Normalized(), 1e-10) {
		t.Errorf("Test case 2 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

	// Test case 3: Random quaternion
	quat = NewMQuaternionByValues(0.1, 0.2, 0.3, 0.4)
	expectedMat = NewMMat4ByValues(
		0.13333333333333353, 0.9333333333333332, -0.33333333333333326, 0,
		-0.6666666666666666, 0.3333333333333335, 0.6666666666666665, 0,
		0.7333333333333332, 0.13333333333333336, 0.6666666666666667, 0,
		0, 0, 0, 1,
	)
	actualMat = quat.ToMat4()

	if !actualMat.PracticallyEquals(expectedMat, 1e-10) {
		t.Errorf("Test case 3 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

	if !actualMat.Quaternion().PracticallyEquals(quat.Normalized(), 1e-10) {
		t.Errorf("Test case 2 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

}
