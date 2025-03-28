package mmath

import (
	"math"
	"testing"
)

func TestMQuaternionToEulerAngles(t *testing.T) {
	quat := NewMQuaternionByValues(0.7071067811865476, 0.0, 0.0, 0.7071067811865476)
	expected := MVec3{1.5707963267948966, 0, 0}

	result := quat.ToRadians()

	if !result.NearEquals(&expected, 1e-10) {
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

	if !result.NearEquals(expected, 1e-10) {
		t.Errorf("Slerp failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToFixedAxisRotation(t *testing.T) {
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{1, 0, 0}
		expected := NewMQuaternionByValues(0.866025403784439, 0, 0, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{0, 1, 0}
		expected := NewMQuaternionByValues(0, 0.866025403784439, 0, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{0, 0, 1}
		expected := NewMQuaternionByValues(0, 0, 0.866025403784439, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
		fixedAxis := MVec3{0.5, 0.7, 0.2}
		expected := NewMQuaternionByValues(0.49029033784546, 0.686406472983644, 0.196116135138184, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
	{
		quat := NewMQuaternionByValues(0.5, 0.5, 0.5, -0.5)
		fixedAxis := MVec3{0.5, -0.7, 0.2}
		expected := NewMQuaternionByValues(-0.49029033784546, 0.686406472983644, -0.196116135138184, 0.5)

		result := quat.ToFixedAxisRotation(&fixedAxis)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
		}
	}
}

func TestMQuaternionNormalized(t *testing.T) {
	quat1 := NewMQuaternionByValues(2, 3, 4, 1)
	expected1 := NewMQuaternionByValues(0.36514837, 0.54772256, 0.73029674, 0.18257419)

	result1 := quat1.Normalized()

	if !result1.NearEquals(expected1, 1e-7) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected1, result1)
	}

	quat2 := NewMQuaternionByValues(0, 0, 0, 1)
	expected2 := NewMQuaternionByValues(0, 0, 0, 1)

	result2 := quat2.Normalized()

	if !result2.NearEquals(expected2, 1e-10) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected2, result2)
	}

	quat3 := NewMQuaternion()
	expected3 := NewMQuaternionByValues(0, 0, 0, 1)

	result3 := quat3.Normalized()

	if !result3.NearEquals(expected3, 1e-10) {
		t.Errorf("Normalized failed. Expected %v, got %v", expected3, result3)
	}
}

func TestFromEulerAnglesDegrees(t *testing.T) {
	expected1 := NewMQuaternionByValues(0, 0, 0, 1)

	result1 := NewMQuaternionFromDegrees(0, 0, 0)

	if !result1.NearEquals(expected1, 1e-8) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(0.08715574, 0.0, 0.0, 0.9961947)

	result2 := NewMQuaternionFromDegrees(10, 0, 0)

	if !result2.NearEquals(expected2, 1e-6) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected2, result2)
	}

	expected3 := NewMQuaternionByValues(0.12767944069578063, 0.14487812541736914, 0.2685358227515692, 0.943714364147489)

	result3 := NewMQuaternionFromDegrees(10, 20, 30)

	if !result3.NearEquals(expected3, 1e-5) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected3, result3)
	}

	expected4 := NewMQuaternionByValues(0.47386805375453483, 0.20131048764138493, -0.6147244358103234, 0.5975257510887351)

	result4 := NewMQuaternionFromDegrees(60, -20, -80)

	if !result4.NearEquals(expected4, 1e-5) {
		t.Errorf("FromEulerAnglesDegrees failed. Expected %v, got %v", expected4, result4)
	}
}

func TestMQuaternionToEulerAnglesDegrees(t *testing.T) {
	expected1 := &MVec3{0, 0, 0}

	qq1 := NewMQuaternionByValues(0, 0, 0, 1)
	result1 := qq1.ToDegrees()

	if !result1.NearEquals(expected1, 1e-8) {
		t.Errorf("ToEulerAnglesDegrees failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := &MVec3{10, 0, 0}

	qq2 := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)
	result2 := qq2.ToDegrees()

	if !result2.NearEquals(expected2, 1e-5) {
		t.Errorf("ToEulerAnglesDegrees failed. Expected %v, got %v", expected2, result2)
	}

	expected3 := &MVec3{10, 20, 30}

	qq3 := NewMQuaternionByValues(0.12767944, 0.14487813, 0.23929834, 0.95154852)
	result3 := qq3.ToDegrees()

	if !result3.NearEquals(expected3, 1e-5) {
		t.Errorf("ToEulerAnglesDegrees failed. Expected %v, got %v", expected3, result3)
	}

	expected4 := &MVec3{60, -20, -80}

	qq4 := NewMQuaternionByValues(0.47386805, 0.20131049, -0.48170221, 0.70914465)
	result4 := qq4.ToDegrees()

	if !result4.NearEquals(expected4, 1e-5) {
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

	if !result1.NearEquals(expected1, 1e-8) {
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

	if !result2.NearEquals(expected2, 1e-8) {
		t.Errorf("MQuaternionMultiply failed. Expected %v, got %v", expected2, result2)
	}
}

func TestNewMQuaternionFromAxisAngles(t *testing.T) {
	expected1 := NewMQuaternionByValues(
		0.0691722994246875, 0.138344598849375, 0.207516898274062, 0.965925826289068)
	result1 := NewMQuaternionFromAxisAnglesRotate(&MVec3{1, 2, 3}, DegToRad(30))

	if !result1.NearEquals(expected1, 1e-5) {
		t.Errorf("NewMQuaternionFromAxisAngles failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(
		-0.116858651016609, 0.779057673444061, -0.389528836722031, 0.477158760259608)
	result2 := NewMQuaternionFromAxisAnglesRotate(&MVec3{-3, 20, -10}, DegToRad(123))

	if !result2.NearEquals(expected2, 1e-5) {
		t.Errorf("NewMQuaternionFromAxisAngles failed. Expected %v, got %v", expected2, result2)
	}

	axis := MVec3{1, 0, 0}
	angle := math.Pi / 2
	expected := NewMQuaternionByValues(0.707106781186548, 0, 0, 0.7071067811865476)

	result := NewMQuaternionFromAxisAnglesRotate(&axis, angle)

	if !result.NearEquals(expected, 1e-10) {
		t.Errorf("NewMQuaternionFromAxisAngles failed. Expected %v, got %v", expected, result)
	}
}
func TestMQuaternionFromDirection(t *testing.T) {
	expected1 := NewMQuaternionByValues(
		-0.3115472173245163, -0.045237910083403, -0.5420603160713341, 0.7791421414666787)
	result1 := NewMQuaternionFromDirection(&MVec3{1, 2, 3}, &MVec3{4, 5, 6})

	if !result1.NearEquals(expected1, 1e-7) {
		t.Errorf("MQuaternionFromDirection failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(
		-0.543212292317204, -0.6953153333136457, -0.20212324833235548, 0.42497433477564167)
	result2 := NewMQuaternionFromDirection(&MVec3{-10, 20, -15}, &MVec3{40, -5, 6})

	if !result2.NearEquals(expected2, 1e-7) {
		t.Errorf("MQuaternionFromDirection failed. Expected %v, got %v", expected2, result2)
	}

}

func TestMQuaternionRotate(t *testing.T) {
	expected1 := NewMQuaternionByValues(
		-0.04597839511020707, 0.0919567902204141, -0.04597839511020706, 0.9936377222602503)
	result1 := NewMQuaternionRotate(&MVec3{1, 2, 3}, &MVec3{4, 5, 6})

	if !result1.NearEquals(expected1, 1e-5) {
		t.Errorf("MQuaternionRotate failed. Expected %v, got %v", expected1, result1)
	}

	expected2 := NewMQuaternionByValues(
		0.042643949239185255, -0.511727390870223, -0.7107324873197542, 0.48080755245182594)
	result2 := NewMQuaternionRotate(&MVec3{-10, 20, -15}, &MVec3{40, -5, 6})

	if !result2.NearEquals(expected2, 1e-5) {
		t.Errorf("MQuaternionRotate failed. Expected %v, got %v", expected2, result2)
	}
}

func TestMQuaternionToMatrix4x4(t *testing.T) {
	expected := NewMMat4ByValues(
		0.45487413, -0.49240388, -0.74204309, 0.0,
		0.87398231, 0.08682409, 0.47813857, 0.0,
		-0.17101007, -0.8660254, 0.46984631, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)

	qq1 := NewMQuaternionByValues(
		0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844)
	result1 := qq1.ToMat4()

	if !result1.NearEquals(expected, 1e-5) {
		t.Errorf("ToMatrix4x4 failed. Expected %v, got %v", expected, result1)
	}

	expected2 := NewMMat4ByValues(
		-0.28213944, 0.69636424, -0.65990468, 0.0,
		0.48809647, 0.69636424, 0.52615461, 0.0,
		0.82592928, -0.17364818, -0.53636474, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)

	// np.array([10, 123, 45])
	qq2 := NewMQuaternionByValues(
		0.3734504874442106, 0.7929168339527322, 0.11114231087966482, 0.4684709324967611)
	result2 := qq2.ToMat4()

	if !result2.NearEquals(expected2, 1e-5) {
		t.Errorf("ToMatrix4x4 failed. Expected %v, got %v", expected, result2)
	}
}

func TestMQuaternionMulVec3(t *testing.T) {
	expected := &MVec3{16.89808539, -29.1683191, 16.23772986}
	//  np.array([60, -20, -80]),
	qq := NewMQuaternionByValues(
		0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844)
	result := qq.MulVec3(&MVec3{10, 20, 30})

	if !result.NearEquals(expected, 1e-5) {
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

func TestMQuaternionMulScalar(t *testing.T) {
	{
		quat := NewMQuaternionFromDegrees(90, 0, 0)
		factor := 0.5
		expected := NewMQuaternionFromDegrees(45, 0, 0)

		result := quat.MuledScalar(factor)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("MulScalar failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionFromDegrees(-24.53194, 180, 180)
		factor := 0.5
		expected := NewMQuaternionFromDegrees(77.73403, 0, 0)

		result := quat.MuledScalar(factor)

		if !result.NearEquals(expected, 1e-6) {
			t.Errorf("MulScalar failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionFromDegrees(-24.53194, 180, 180)
		factor := -0.5
		expected := NewMQuaternionFromDegrees(-77.73403, 0, 0)

		result := quat.MuledScalar(factor)

		if !result.NearEquals(expected, 1e-6) {
			t.Errorf("MulScalar failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionByValues(0, 0, 0, 1)
		factor := 0.5
		expected := NewMQuaternionByValues(0, 0, 0, 1)

		result := quat.MuledScalar(factor)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("MulScalar failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
		}
	}

	{
		quat := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)
		factor := 1.0
		expected := NewMQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455)

		result := quat.MuledScalar(factor)

		if !result.NearEquals(expected, 1e-10) {
			t.Errorf("MulScalar failed. Expected %v, got %v(%v)", expected, result, result.ToDegrees())
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

	if !actualMat.NearEquals(expectedMat, 1e-10) {
		t.Errorf("Test case 1 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

	// Test case 2: Non-identity quaternion
	quat = NewMQuaternionByValues(0.5, 0.5, 0.5, 0.5)
	expectedMat = NewMMat4ByValues(
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		1.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)
	actualMat = quat.ToMat4()

	if !actualMat.NearEquals(expectedMat, 1e-10) {
		t.Errorf("Test case 2 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

	// Test case 3: Random quaternion
	quat = NewMQuaternionByValues(0.1, 0.2, 0.3, 0.4)
	expectedMat = NewMMat4ByValues(
		0.13333333, 0.93333333, -0.33333333, 0.0,
		-0.66666667, 0.33333333, 0.66666667, 0.0,
		0.73333333, 0.13333333, 0.66666667, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)
	actualMat = quat.Normalize().ToMat4()

	if !actualMat.NearEquals(expectedMat, 1e-5) {
		t.Errorf("Test case 3 failed: Expected %v, but got %v", expectedMat, actualMat)
	}

}

func TestFindSlerpT(t *testing.T) {
	tests := []struct {
		q0, q1 *MQuaternion
		t      float64
	}{
		{
			q0: NewMQuaternionFromDegrees(10, 20, 30),
			q1: NewMQuaternionFromDegrees(60, -20, -80),
			t:  0.3,
		},
		{
			q0: NewMQuaternionFromDegrees(10, 20, 30),
			q1: NewMQuaternionFromDegrees(60, -20, -80),
			t:  0.6,
		},
		{
			q0: NewMQuaternionFromDegrees(-10, 20, 30),
			q1: NewMQuaternionFromDegrees(60, 20, -80),
			t:  0.9,
		},
		{
			q0: &MQuaternion{X: -0.008417034521698952, Y: -0.001099314889870584, Z: -0.33594009280204773, W: 0.9418451189994812},
			q1: &MQuaternion{X: -0.00845526841302267, Y: -0.0014676861230259702, Z: -0.3248561346517127, W: 0.9457245085714213},
			t:  0.5,
		},
		{
			q0: &MQuaternion{X: 0.01753191463649273, Y: 0.023771926760673523, Z: 0, W: 0.9995636940002441},
			q1: &MQuaternion{X: 0.3140714764595032, Y: 0.42585673928260803, Z: 0, W: 0.848531186580658},
			t:  0.5,
		},
	}

	for _, tt := range tests {
		q := tt.q0.Slerp(tt.q1, tt.t)
		result := FindSlerpT(tt.q0, tt.q1, q, tt.t)
		if math.Abs(tt.t-result) > 1e-3 {
			t.Errorf("FindSlerpT failed. Expected %v, got %v", tt.t, result)
		}
	}
}
