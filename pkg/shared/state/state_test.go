// 指示: miu200521358
package state

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/adapter/physics_api"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

type testModel struct {
	*hashable.HashableBase
}

// TestNewSharedStateDefaults は既定値を確認する。
func TestNewSharedStateDefaults(t *testing.T) {
	ss := NewSharedState(2)
	if ss.Frame() != 0 {
		t.Errorf("Frame: got=%v", ss.Frame())
	}
	if ss.MaxFrame() != 1 {
		t.Errorf("MaxFrame: got=%v", ss.MaxFrame())
	}
	if ss.FrameInterval() != -1 {
		t.Errorf("FrameInterval: got=%v", ss.FrameInterval())
	}
	if ss.IsClosed() {
		t.Errorf("IsClosed: expected false")
	}

	phys, ok := ss.PhysicsWorldMotion(0).(*defaultMotion)
	if !ok {
		t.Fatalf("PhysicsWorldMotion type mismatch")
	}
	if phys.Gravity != -9.8 || phys.MaxSubSteps != physics_api.PhysicsDefaultMaxSubSteps || phys.FixedTimeStep != 60 {
		t.Errorf("Physics defaults: got=%v", phys)
	}
	wind, ok := ss.WindMotion(0).(*defaultMotion)
	if !ok {
		t.Fatalf("WindMotion type mismatch")
	}
	if wind.WindEnabled || wind.WindDirection != [3]float32{0, 0, 0} {
		t.Errorf("Wind defaults: got=%v", wind)
	}
}

// TestFlagUpdate はフラグ更新を確認する。
func TestFlagUpdate(t *testing.T) {
	ss := NewSharedState(1)
	ss.EnableFlag(STATE_FLAG_FRAME_DROP)
	if !ss.HasFlag(STATE_FLAG_FRAME_DROP) {
		t.Errorf("EnableFlag failed")
	}
	ss.DisableFlag(STATE_FLAG_FRAME_DROP)
	if ss.HasFlag(STATE_FLAG_FRAME_DROP) {
		t.Errorf("DisableFlag failed")
	}
}

// TestSetModelInitializesSelected は材質選択初期化を確認する。
func TestSetModelInitializesSelected(t *testing.T) {
	ss := NewSharedState(1)
	model := &testModel{HashableBase: hashable.NewHashableBase("", "")}
	ss.SetModel(0, 0, model)
	idxs := ss.SelectedMaterialIndexes(0, 0)
	if idxs == nil {
		t.Errorf("SelectedMaterialIndexes should not be nil")
	}
	if len(idxs) != 0 {
		t.Errorf("SelectedMaterialIndexes length: got=%v", len(idxs))
	}
	vertexIdxs := ss.SelectedVertexIndexes(0, 0)
	if vertexIdxs == nil {
		t.Errorf("SelectedVertexIndexes should not be nil")
	}
	if len(vertexIdxs) != 0 {
		t.Errorf("SelectedVertexIndexes length: got=%v", len(vertexIdxs))
	}
}

// TestOutOfRangeAccess は範囲外アクセスを確認する。
func TestOutOfRangeAccess(t *testing.T) {
	ss := NewSharedState(1)
	if ss.Model(0, 99) != nil {
		t.Errorf("Model out-of-range should be nil")
	}
	if ss.Motion(0, 99) != nil {
		t.Errorf("Motion out-of-range should be nil")
	}
}
