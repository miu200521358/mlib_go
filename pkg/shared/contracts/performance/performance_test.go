// 指示: miu200521358
package performance

import "testing"

// TestMaxBoneFrames は既定値とテスト用差し替えを確認する。
func TestMaxBoneFrames(t *testing.T) {
	if MAX_BONE_FRAMES != 600000 {
		t.Errorf("MAX_BONE_FRAMES default: got=%v", MAX_BONE_FRAMES)
	}

	old := MAX_BONE_FRAMES
	MAX_BONE_FRAMES = 20
	if MAX_BONE_FRAMES != 20 {
		t.Errorf("MAX_BONE_FRAMES override: got=%v", MAX_BONE_FRAMES)
	}
	MAX_BONE_FRAMES = old
}

// TestDefaultPerformancePolicy は既定の性能方針を確認する。
func TestDefaultPerformancePolicy(t *testing.T) {
	if !DEFAULT_PERFORMANCE_POLICY.RealtimePreferred {
		t.Errorf("RealtimePreferred: got=%v", DEFAULT_PERFORMANCE_POLICY.RealtimePreferred)
	}
	if !DEFAULT_PERFORMANCE_POLICY.AllowDeferredApply {
		t.Errorf("AllowDeferredApply: got=%v", DEFAULT_PERFORMANCE_POLICY.AllowDeferredApply)
	}
	if DEFAULT_PERFORMANCE_POLICY.ApplyTiming != APPLY_TIMING_USER_SELECT {
		t.Errorf("ApplyTiming: got=%v", DEFAULT_PERFORMANCE_POLICY.ApplyTiming)
	}
}

// TestDefaultMaxSubSteps は既定の最大サブステップ数を確認する。
func TestDefaultMaxSubSteps(t *testing.T) {
	if DefaultMaxSubSteps != 5 {
		t.Errorf("DefaultMaxSubSteps: got=%v", DefaultMaxSubSteps)
	}
}
