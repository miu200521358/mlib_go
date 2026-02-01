// 指示: miu200521358
package motion

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

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
	q1 := mmath.NewQuaternionFromDegrees(0, 0, 0)
	q2 := mmath.NewQuaternionFromDegrees(0, 90, 0)
	if !out.Quaternion.NearEquals(q1.Slerp(q2, 0.5), 1e-6) {
		t.Fatalf("Quaternion lerp mismatch")
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
