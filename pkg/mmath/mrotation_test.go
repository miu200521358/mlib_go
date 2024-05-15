package mmath

import (
	"testing"
)

func TestNewRotationByRadians(t *testing.T) {
	radians := MVec3{1, 2, 3}
	rotation := NewRotationByRadians(&radians)

	if *rotation.GetRadians() != radians {
		t.Errorf("Expected GetRadians() to return %v, but got %v", radians, rotation.radians.String())
	}
}

func TestNewRotationByDegrees(t *testing.T) {
	degrees := MVec3{90, 180, 270}
	rotation := NewRotationByDegrees(&degrees)

	if *rotation.GetDegrees() != degrees {
		t.Errorf("Expected GetDegrees() to return %v, but got %v", degrees, rotation.degrees.String())
	}
}

func TestNewRotationByQuaternion(t *testing.T) {
	quaternion := NewMQuaternionByValues(1, 0, 0, 0)
	rotation := NewRotationByQuaternion(&quaternion)

	if *rotation.GetQuaternion() != quaternion {
		t.Errorf("Expected GetQuaternion() to return %v, but got %v", quaternion, rotation.quaternion.String())
	}
}

func TestT_Copy(t *testing.T) {
	qq := NewMQuaternionByValues(1, 0, 0, 0)
	rot := &MRotation{
		radians:    &MVec3{1, 2, 3},
		degrees:    &MVec3{90, 180, 270},
		quaternion: &qq,
	}

	copied := rot.Copy()

	if &copied == &rot {
		t.Error("Expected Copy() to return a different instance")
	}

	if !copied.GetRadians().PracticallyEquals(rot.GetRadians(), 1e-10) {
		t.Errorf("Copied instance does not match the original (radians) %s %s", copied.radians.String(), rot.radians.String())
	}
	if !copied.GetDegrees().PracticallyEquals(rot.GetDegrees(), 1e-10) {
		t.Errorf("Copied instance does not match the original (degrees) %s %s", copied.degrees.String(), rot.degrees.String())
	}
	if 1-copied.GetQuaternion().Dot(rot.GetQuaternion()) > 1e-10 {
		t.Errorf("Copied instance does not match the original (quaternion) %s %s", copied.quaternion.String(), rot.quaternion.String())
	}
}
