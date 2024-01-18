package mquaternion_test

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/math/mquaternion"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

func TestFromAxisAngle(t *testing.T) {
	axis := mvec3.T{1, 0, 0}
	angle := math.Pi / 2
	expected := mquaternion.T{0.7071067811865476, 0, 0, 0.7071067811865476}

	result := mquaternion.FromAxisAngle(&axis, angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestFromXAxisAngle(t *testing.T) {
	angle := math.Pi / 2
	expected := mquaternion.T{0.7071067811865476, 0, 0, 0.7071067811865476}

	result := mquaternion.FromXAxisAngle(angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromXAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestFromYAxisAngle(t *testing.T) {
	angle := math.Pi / 2
	expected := mquaternion.T{0, 0.7071067811865476, 0, 0.7071067811865476}

	result := mquaternion.FromYAxisAngle(angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromYAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestFromZAxisAngle(t *testing.T) {
	angle := math.Pi / 2
	expected := mquaternion.T{0, 0, 0.7071067811865476, 0.7071067811865476}

	result := mquaternion.FromZAxisAngle(angle)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromZAxisAngle failed. Expected %v, got %v", expected, result)
	}
}

func TestFromEulerAngles(t *testing.T) {
	angles := mvec3.T{math.Pi / 2, math.Pi / 2, math.Pi / 2}
	expected := mquaternion.T{0.7071067811865476, 0.0, 0.0, 0.7071067811865476}

	result := mquaternion.FromEulerAngles(angles.GetX(), angles.GetY(), angles.GetZ())

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("FromEulerAngles failed. Expected %v, got %v", expected, result)
	}
}

func TestToEulerAngles(t *testing.T) {
	quat := mquaternion.T{0.7071067811865476, 0.0, 0.0, 0.7071067811865476}
	expected := mvec3.T{1.5707963267948966, 0, 0}

	result := quat.ToEulerAngles()

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 {
		t.Errorf("ToEulerAngles failed. Expected %v, got %v", expected, result)
	}
}

func TestToDegree(t *testing.T) {
	quat := mquaternion.T{0.08715574274765817, 0.0, 0.0, 0.9961946980917455}
	expected := 10.0

	result := quat.ToDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestToDegree2(t *testing.T) {
	quat := mquaternion.T{0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885}
	expected := 35.81710117358426

	result := quat.ToDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestToSignedDegree(t *testing.T) {
	quat := mquaternion.T{0.08715574274765817, 0.0, 0.0, 0.9961946980917455}
	expected := 10.0

	result := quat.ToSignedDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestToSignedDegree2(t *testing.T) {
	quat := mquaternion.T{0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844}
	expected := 89.66927179998277

	result := quat.ToSignedDegree()

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestDot(t *testing.T) {
	quat1 := mquaternion.T{0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844}
	quat2 := mquaternion.T{0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885}
	expected := 0.6491836986795888

	result := quat1.Dot(&quat2)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("ToDegree failed. Expected %v, got %v", expected, result)
	}
}

func TestSlerp(t *testing.T) {
	quat1 := mquaternion.T{0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844}
	quat2 := mquaternion.T{0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885}
	tValue := 0.3
	expected := mquaternion.T{0.3973722198386427, 0.19936467087655246, -0.27953105525419597, 0.851006131620254}

	result := mquaternion.Slerp(&quat1, &quat2, tValue)

	if math.Abs(result.GetX()-expected.GetX()) > 1e-10 || math.Abs(result.GetY()-expected.GetY()) > 1e-10 || math.Abs(result.GetZ()-expected.GetZ()) > 1e-10 || math.Abs(result.GetW()-expected.GetW()) > 1e-10 {
		t.Errorf("Slerp failed. Expected %v, got %v", expected, result)
	}
}
