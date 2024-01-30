package mmath

import (
	"math"
	"testing"

)

func TestMQuaternionFromAxisAngle(t *testing.T) {
	axis := MVec3{1, 0, 0}
	angle := math.Pi / 2
	expected := MQuaternion{0.7071067811865476, 0, 0, 0.7071067811865476}

	result := FromAxisAngle(&axis, angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionFromXAxisAngle(t *testing.T) {
	angle := math.Pi / 2
	expected := MQuaternion{0.7071067811865476, 0, 0, 0.7071067811865476}

	result := FromXAxisAngle(angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromXAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionFromYAxisAngle(t *testing.T) {
	angle := math.Pi / 2
	expected := MQuaternion{0, 0.7071067811865476, 0, 0.7071067811865476}

	result := FromYAxisAngle(angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromYAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionFromZAxisAngle(t *testing.T) {
	angle := math.Pi / 2
	expected := MQuaternion{0, 0, 0.7071067811865476, 0.7071067811865476}

	result := FromZAxisAngle(angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromZAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionFromEulerAngles(t *testing.T) {
	angles := MVec3{math.Pi / 2, math.Pi / 2, math.Pi / 2}
	expected := MQuaternion{0.7071067811865476, 0.0, 0.0, 0.7071067811865476}

	result := FromEulerAngles(angles.GetX(), angles.GetY(), angles.GetZ())

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromEulerAngles failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToEulerAngles(t *testing.T) {
	quat := MQuaternion{0.7071067811865476, 0.0, 0.0, 0.7071067811865476}
	expected := MVec3{1.5707963267948966, 0, 0}

	result := quat.ToEulerAngles()

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 {
		t.Errorf("ToEulerAngles failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToDegree(t *testing.T) {
	quat := MQuaternion{0.08715574274765817, 0.0, 0.0, 0.9961946980917455}
	expected := 10.0

	result := quat.ToDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToDegree2(t *testing.T) {
	quat := MQuaternion{0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885}
	expected := 35.81710117358426

	result := quat.ToDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToSignedDegree(t *testing.T) {
	quat := MQuaternion{0.08715574274765817, 0.0, 0.0, 0.9961946980917455}
	expected := 10.0

	result := quat.ToSignedDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToSignedDegree2(t *testing.T) {
	quat := MQuaternion{0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844}
	expected := 89.66927179998277

	result := quat.ToSignedDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionDot(t *testing.T) {
	quat1 := MQuaternion{0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844}
	quat2 := MQuaternion{0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885}
	expected := 0.6491836986795888

	result := quat1.Dot(&quat2)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionSlerp(t *testing.T) {
	quat1 := MQuaternion{0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844}
	quat2 := MQuaternion{0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885}
	tValue := 0.3
	expected := MQuaternion{0.3973722198386427, 0.19936467087655246, -0.27953105525419597, 0.851006131620254}

	result := Slerp(&quat1, &quat2, tValue)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("Slerp failed. Expected %v, got %v", expected, result)
	}
}

func TestMQuaternionToFixedAxisRotation(t *testing.T) {
	quat := MQuaternion{0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844}
	fixedAxis := MVec3{0.0, 1.0, 0.0}
	expected := &MQuaternion{0.00000, 0.70506, 0.00000, 0.70914}

	result := quat.ToFixedAxisRotation(&fixedAxis)

	if result.PracticallyEquals(expected, 1e-10) {
		t.Errorf("ToFixedAxisRotation failed. Expected %v, got %v", expected, result)
	}
}
