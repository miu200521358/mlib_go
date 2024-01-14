package mrotation

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/math/mquaternion"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"

)

func TestNewBaseRotationModelByRadians(t *testing.T) {
	radians := mvec3.T{1, 2, 3}
	rotation := NewBaseRotationModelByRadians(&radians)

	if *rotation.GetRadians() != radians {
		t.Errorf("Expected GetRadians() to return %v, but got %v", radians, rotation.radians.String())
	}
}

func TestNewBaseRotationModelByDegrees(t *testing.T) {
	degrees := mvec3.T{90, 180, 270}
	rotation := NewBaseRotationModelByDegrees(&degrees)

	if *rotation.GetDegrees() != degrees {
		t.Errorf("Expected GetDegrees() to return %v, but got %v", degrees, rotation.degrees.String())
	}
}

func TestNewBaseRotationModelByQuaternion(t *testing.T) {
	quaternion := mquaternion.T{1, 0, 0, 0}
	rotation := NewBaseRotationModelByQuaternion(&quaternion)

	if *rotation.GetQuaternion() != quaternion {
		t.Errorf("Expected GetQuaternion() to return %v, but got %v", quaternion, rotation.quaternion.String())
	}
}

func TestT_Copy(t *testing.T) {
	rot := &T{
		radians:    mvec3.T{1, 2, 3},
		degrees:    mvec3.T{90, 180, 270},
		quaternion: mquaternion.T{1, 0, 0, 0},
	}

	copied := rot.Copy()

	if copied == rot {
		t.Error("Expected Copy() to return a different instance")
	}

	if copied.GetRadians() != rot.GetRadians() {
		t.Errorf("Copied instance does not match the original (radians) %s %s", copied.radians.String(), rot.radians.String())
	}
	if copied.GetDegrees() != rot.GetDegrees() {
		t.Errorf("Copied instance does not match the original (degrees) %s %s", copied.degrees.String(), rot.degrees.String())
	}
	if copied.GetQuaternion() != rot.GetQuaternion() {
		t.Errorf("Copied instance does not match the original (quaternion) %s %s", copied.quaternion.String(), rot.quaternion.String())
	}
}
