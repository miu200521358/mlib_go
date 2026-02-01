// 指示: miu200521358
package motion

import "testing"

// TestMorphFrameLerp はモーフ補間を確認する。
func TestMorphFrameLerp(t *testing.T) {
	prev := NewMorphFrame(0)
	prev.Ratio = 0
	next := NewMorphFrame(10)
	next.Ratio = 1
	out := next.lerpFrame(prev, 5)
	if out == nil {
		t.Fatalf("lerpFrame nil")
	}
	if out.Ratio < 0.49 || out.Ratio > 0.51 {
		t.Fatalf("Ratio lerp: got=%v", out.Ratio)
	}
}

// TestMorphNameFramesContainsActive は有効判定を確認する。
func TestMorphNameFramesContainsActive(t *testing.T) {
	frames := NewMorphNameFrames("m")
	frames.Append(NewMorphFrame(0))
	if frames.ContainsActive() {
		t.Fatalf("ContainsActive should be false")
	}

	mf := NewMorphFrame(1)
	mf.Ratio = 0.5
	frames.Append(mf)
	if !frames.ContainsActive() {
		t.Fatalf("ContainsActive should be true")
	}
}
