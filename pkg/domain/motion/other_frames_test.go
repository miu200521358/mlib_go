// 指示: miu200521358
package motion

import "testing"

// TestLightFrameLerp はライト補間を確認する。
func TestLightFrameLerp(t *testing.T) {
	prev := NewLightFrame(0)
	prev.Position = vec3(0, 0, 0)
	prev.Color = vec3(0, 0, 0)

	next := NewLightFrame(10)
	next.Position = vec3(10, 0, 0)
	next.Color = vec3(1, 0, 0)

	out := next.lerpFrame(prev, 5)
	if out.Position.X < 4.9 || out.Position.X > 5.1 {
		t.Fatalf("Position lerp: got=%v", out.Position.X)
	}
	if out.Color.X < 0.49 || out.Color.X > 0.51 {
		t.Fatalf("Color lerp: got=%v", out.Color.X)
	}
}

// TestShadowFrameLerp はシャドウ補間を確認する。
func TestShadowFrameLerp(t *testing.T) {
	prev := NewShadowFrame(0)
	prev.ShadowMode = 1
	prev.Distance = 0
	next := NewShadowFrame(10)
	next.ShadowMode = 2
	next.Distance = 10
	out := next.lerpFrame(prev, 5)
	if out.ShadowMode != 1 {
		t.Fatalf("ShadowMode should use prev")
	}
	if out.Distance < 4.9 || out.Distance > 5.1 {
		t.Fatalf("Distance lerp: got=%v", out.Distance)
	}
}

// TestLightFramesClean は既定値のみの削除を確認する。
func TestLightFramesClean(t *testing.T) {
	frames := NewLightFrames()
	frames.Append(NewLightFrame(0))
	frames.Clean()
	if frames.Len() != 0 {
		t.Fatalf("Clean should delete default light")
	}
}

// TestShadowFramesClean は既定値のみの削除を確認する。
func TestShadowFramesClean(t *testing.T) {
	frames := NewShadowFrames()
	frames.Append(NewShadowFrame(0))
	frames.Clean()
	if frames.Len() != 0 {
		t.Fatalf("Clean should delete default shadow")
	}
}

// TestIkFrameIsEnable はIKの有効判定を確認する。
func TestIkFrameIsEnable(t *testing.T) {
	frame := NewIkFrame(0)
	if !frame.IsEnable("bone") {
		t.Fatalf("IsEnable default should be true")
	}
	frame.IkList = []*IkEnabledFrame{{BaseFrame: NewBaseFrame(0), BoneName: "bone", Enabled: false}}
	if frame.IsEnable("bone") {
		t.Fatalf("IsEnable should be false")
	}
}

// TestIkFramesClean は既定値のみの削除を確認する。
func TestIkFramesClean(t *testing.T) {
	frames := NewIkFrames()
	frames.Append(NewIkFrame(0))
	frames.Clean()
	if frames.Len() != 0 {
		t.Fatalf("Clean should delete default")
	}
}
