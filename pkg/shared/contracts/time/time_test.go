// 指示: miu200521358
package time

import "testing"

// TestTimeConstants はFPS/SPFの定数値を確認する。
func TestTimeConstants(t *testing.T) {
	if DefaultFps != 30 {
		t.Errorf("DefaultFps: got=%v", DefaultFps)
	}
}

// TestTimeConversions は時間変換の基本動作を確認する。
func TestTimeConversions(t *testing.T) {
	sec := FramesToSeconds(60, 30)
	if sec != 2 {
		t.Errorf("FramesToSeconds: got=%v", sec)
	}
	frames := SecondsToFrames(2, 30)
	if frames != 60 {
		t.Errorf("SecondsToFrames: got=%v", frames)
	}
	spf := FpsToSpf(30)
	expectedSpf := Seconds(1.0 / 30.0)
	diffSpf := float32(spf - expectedSpf)
	if diffSpf < 0 {
		diffSpf = -diffSpf
	}
	if diffSpf > 1e-5 {
		t.Errorf("FpsToSpf: got=%v", spf)
	}
	fps := SpfToFps(spf)
	diffFps := float32(fps) - 30
	if diffFps < 0 {
		diffFps = -diffFps
	}
	if diffFps > 1e-5 {
		t.Errorf("SpfToFps: got=%v", fps)
	}
}

// TestClampFrame はClampの挙動を確認する。
func TestClampFrame(t *testing.T) {
	if got := ClampFrame(5, 0, 10); got != 5 {
		t.Errorf("ClampFrame center: got=%v", got)
	}
	if got := ClampFrame(-1, 0, 10); got != 0 {
		t.Errorf("ClampFrame min: got=%v", got)
	}
	if got := ClampFrame(11, 0, 10); got != 10 {
		t.Errorf("ClampFrame max: got=%v", got)
	}

	r := FrameRange{Start: 0, End: 10}
	if !IsFrameInRange(0, r) || !IsFrameInRange(10, r) {
		t.Errorf("IsFrameInRange inclusive failed")
	}
	if IsFrameInRange(11, r) {
		t.Errorf("IsFrameInRange out of range")
	}
}

// TestTimeZeroCases はゼロ除算回避を確認する。
func TestTimeZeroCases(t *testing.T) {
	if FramesToSeconds(10, 0) != 0 {
		t.Errorf("FramesToSeconds fps=0 should be 0")
	}
	if FpsToSpf(0) != 0 {
		t.Errorf("FpsToSpf fps=0 should be 0")
	}
	if SpfToFps(0) != 0 {
		t.Errorf("SpfToFps spf=0 should be 0")
	}
}
