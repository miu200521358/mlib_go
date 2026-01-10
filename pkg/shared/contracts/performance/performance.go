// 指示: miu200521358
package performance

// ApplyTiming は適用タイミングの方針を表す。
type ApplyTiming int

const (
	// APPLY_TIMING_USER_SELECT はユーザー選択を表す。
	APPLY_TIMING_USER_SELECT ApplyTiming = iota
)

// PerformancePolicy は性能方針を表す。
type PerformancePolicy struct {
	RealtimePreferred  bool
	AllowDeferredApply bool
	ApplyTiming        ApplyTiming
}

// MAX_BONE_FRAMES はボーンフレーム数の上限。
var MAX_BONE_FRAMES = 600000

// DEFAULT_PERFORMANCE_POLICY は既定の性能方針。
var DEFAULT_PERFORMANCE_POLICY = PerformancePolicy{
	RealtimePreferred:  true,
	AllowDeferredApply: true,
	ApplyTiming:        APPLY_TIMING_USER_SELECT,
}
