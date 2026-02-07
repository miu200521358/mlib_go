// 指示: miu200521358
package graphics_api

import (
	"math"
	"testing"
)

func TestCameraResetPresetPositionRecenterHorizontalLookAt(t *testing.T) {
	cam := NewDefaultCamera(1280, 720)
	cam.LookAtCenter.X = 3.25
	cam.LookAtCenter.Y = 14.5
	cam.LookAtCenter.Z = -2.75

	beforeDistance := cam.OrbitDistance()

	cam.ResetPresetPosition(0, 0)

	const eps = 1e-9
	if math.Abs(cam.LookAtCenter.X-InitialLookAtCenterX) > eps {
		t.Fatalf("LookAtCenter.X: got=%f want=%f", cam.LookAtCenter.X, InitialLookAtCenterX)
	}
	if math.Abs(cam.LookAtCenter.Z-InitialLookAtCenterZ) > eps {
		t.Fatalf("LookAtCenter.Z: got=%f want=%f", cam.LookAtCenter.Z, InitialLookAtCenterZ)
	}
	if math.Abs(cam.LookAtCenter.Y-14.5) > eps {
		t.Fatalf("LookAtCenter.Y: got=%f want=%f", cam.LookAtCenter.Y, 14.5)
	}
	if math.Abs(cam.OrbitDistance()-beforeDistance) > eps {
		t.Fatalf("OrbitDistance: got=%f want=%f", cam.OrbitDistance(), beforeDistance)
	}
	if math.Abs(cam.Position.X-InitialLookAtCenterX) > eps {
		t.Fatalf("Position.X: got=%f want=%f", cam.Position.X, InitialLookAtCenterX)
	}
}

func TestCameraResetPresetPositionNilReceiver(t *testing.T) {
	var cam *Camera
	cam.ResetPresetPosition(0, 0)
}
