// 指示: miu200521358
package motion

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/contracts/performance"
)

// TestPhysicsFramesGetDefaults は物理フレームの既定値を確認する。
func TestPhysicsFramesGetDefaults(t *testing.T) {
	maxFrames := NewMaxSubStepsFrames()
	if maxFrames.Get(0).MaxSubSteps != performance.DefaultMaxSubSteps {
		t.Fatalf("MaxSubSteps default")
	}
	fixedFrames := NewFixedTimeStepFrames()
	if fixedFrames.Get(0).FixedTimeStepNum != 60 {
		t.Fatalf("FixedTimeStep default")
	}
	gravityFrames := NewGravityFrames()
	g := gravityFrames.Get(0)
	if g.Gravity == nil || !g.Gravity.NearEquals(vec3(0, -9.8, 0), 1e-8) {
		t.Fatalf("Gravity default")
	}
	resetFrames := NewPhysicsResetFrames()
	if resetFrames.Get(0).PhysicsResetType != PHYSICS_RESET_TYPE_NONE {
		t.Fatalf("PhysicsReset default")
	}
}

// TestPhysicsFramesGetNext は次フレーム優先を確認する。
func TestPhysicsFramesGetNext(t *testing.T) {
	frames := NewMaxSubStepsFrames()
	f10 := NewMaxSubStepsFrame(10)
	f10.MaxSubSteps = 5
	frames.Append(f10)
	if frames.Get(5).MaxSubSteps != 5 {
		t.Fatalf("Get should return next")
	}
}

// TestFixedTimeStep は固定タイムステップの計算を確認する。
func TestFixedTimeStep(t *testing.T) {
	frame := NewFixedTimeStepFrame(0)
	frame.FixedTimeStepNum = 0
	if frame.FixedTimeStep() != 1.0/60.0 {
		t.Fatalf("FixedTimeStep default")
	}
	frame.FixedTimeStepNum = 30
	if frame.FixedTimeStep() != 1.0/30.0 {
		t.Fatalf("FixedTimeStep calc")
	}
}

// TestWindFramesDefaults は風フレームの既定値と係数計算を確認する。
func TestWindFramesDefaults(t *testing.T) {
	enabledFrames := NewWindEnabledFrames()
	enabled := enabledFrames.Get(0)
	if enabled == nil || enabled.Enabled {
		t.Fatalf("WindEnabled default")
	}

	directionFrames := NewWindDirectionFrames()
	direction := directionFrames.Get(0)
	if direction == nil || direction.Direction == nil || !direction.Direction.NearEquals(vec3(0, 0, 0), 1e-8) {
		t.Fatalf("WindDirection default")
	}

	lift := NewWindLiftCoeffFrame(0)
	if lift.WindLiftCoeff() != 1.0/60.0 {
		t.Fatalf("WindLiftCoeff default")
	}
	lift.LiftCoeff = 0
	if lift.WindLiftCoeff() != 1 {
		t.Fatalf("WindLiftCoeff zero")
	}

	drag := NewWindDragCoeffFrame(0)
	if drag.WindDragCoeff() != 1.0/60.0 {
		t.Fatalf("WindDragCoeff default")
	}
	drag.DragCoeff = 0
	if drag.WindDragCoeff() != 1 {
		t.Fatalf("WindDragCoeff zero")
	}

	randomness := NewWindRandomnessFrame(0)
	if randomness.WindRandomness() != 0 {
		t.Fatalf("WindRandomness default")
	}
	randomness.Randomness = 10
	if randomness.WindRandomness() != 0.1 {
		t.Fatalf("WindRandomness calc")
	}

	speed := NewWindSpeedFrame(0)
	if speed.WindSpeed() != 1.0/60.0 {
		t.Fatalf("WindSpeed default")
	}

	freq := NewWindTurbulenceFreqHzFrame(0)
	if freq.WindTurbulenceFreqHz() != 0 {
		t.Fatalf("WindTurbulenceFreqHz default")
	}
	freq.TurbulenceFreqHz = 4
	if freq.WindTurbulenceFreqHz() != 0.25 {
		t.Fatalf("WindTurbulenceFreqHz calc")
	}
}

// TestWindFramesGetNext は風フレームの次フレーム優先を確認する。
func TestWindFramesGetNext(t *testing.T) {
	frames := NewWindDirectionFrames()
	f10 := NewWindDirectionFrame(10)
	f10.Direction = vec3Ptr(1, 0, 0)
	frames.Append(f10)
	got := frames.Get(5)
	if got == nil || got.Direction == nil || got.Direction.X != 1 {
		t.Fatalf("Get should return next")
	}
}
