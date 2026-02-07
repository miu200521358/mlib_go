// 指示: miu200521358
package motion

import "testing"

// TestCameraFrameLerp はカメラ補間を確認する。
func TestCameraFrameLerp(t *testing.T) {
	prev := NewCameraFrame(0)
	prev.Degrees = vec3Ptr(0, 0, 0)
	prev.Position = vec3Ptr(0, 0, 0)
	prev.Distance = 0
	prev.ViewOfAngle = 0
	prev.IsPerspectiveOff = false

	next := NewCameraFrame(10)
	next.Degrees = vec3Ptr(0, 90, 0)
	next.Position = vec3Ptr(10, 0, 0)
	next.Distance = 10
	next.ViewOfAngle = 60
	next.IsPerspectiveOff = true

	out := next.lerpFrame(prev, 5)
	if out == nil || out.Position == nil || out.Quaternion == nil {
		t.Fatalf("lerpFrame nil")
	}
	if out.Position.X < 4.9 || out.Position.X > 5.1 {
		t.Fatalf("Position lerp: got=%v", out.Position.X)
	}
	if out.IsPerspectiveOff != true {
		t.Fatalf("IsPerspectiveOff should use next")
	}
	if out.Degrees == nil {
		t.Fatalf("Degrees should not be nil")
	}
	if out.Degrees.Y < 44.9 || out.Degrees.Y > 45.1 {
		t.Fatalf("Degrees lerp mismatch: got=%v", out.Degrees.Y)
	}
	if out.Quaternion == nil {
		t.Fatalf("Quaternion should not be nil")
	}
}

// TestCameraFrameLerpKeepsMultiRotation は多回転の角度が保持されることを確認する。
func TestCameraFrameLerpKeepsMultiRotation(t *testing.T) {
	prev := NewCameraFrame(0)
	prev.Degrees = vec3Ptr(0, 0, 0)

	next := NewCameraFrame(10)
	next.Degrees = vec3Ptr(0, 720, 0)

	out := next.lerpFrame(prev, 5)
	if out == nil || out.Degrees == nil {
		t.Fatalf("lerpFrame nil")
	}
	if out.Degrees.Y < 359.9 || out.Degrees.Y > 360.1 {
		t.Fatalf("multi rotation mismatch: got=%v", out.Degrees.Y)
	}
}

// TestCameraFrameIsDefault は既定値判定を確認する。
func TestCameraFrameIsDefault(t *testing.T) {
	frame := NewCameraFrame(0)
	if !frame.IsDefault() {
		t.Fatalf("IsDefault should be true")
	}
	frame.Distance = 1
	if frame.IsDefault() {
		t.Fatalf("IsDefault should be false")
	}
}

// TestCameraFramesClean は既定値のみの削除を確認する。
func TestCameraFramesClean(t *testing.T) {
	frames := NewCameraFrames()
	frames.Append(NewCameraFrame(0))
	frames.Clean()
	if frames.Len() != 0 {
		t.Fatalf("Clean should delete default")
	}
}
