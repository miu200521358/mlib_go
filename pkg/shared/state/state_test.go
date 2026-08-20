// 指示: miu200521358
package state

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/contracts/performance"
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
	if ss.SelectedVertexDepthMode() != SELECTED_VERTEX_DEPTH_MODE_ALL {
		t.Errorf("SelectedVertexDepthMode default mismatch: got=%v", ss.SelectedVertexDepthMode())
	}
	if ss.SelectedFaceDepthMode() != SELECTED_FACE_DEPTH_MODE_ALL {
		t.Errorf("SelectedFaceDepthMode default mismatch: got=%v", ss.SelectedFaceDepthMode())
	}
	if ss.SelectedFaceMode() != SELECTED_FACE_MODE_LINE {
		t.Errorf("SelectedFaceMode default mismatch: got=%v", ss.SelectedFaceMode())
	}
	ss.SetSelectedVertexDepthMode(SELECTED_VERTEX_DEPTH_MODE_FRONT)
	if ss.SelectedFaceDepthMode() != SELECTED_FACE_DEPTH_MODE_ALL {
		t.Errorf("face depth mode changed with vertex mode: got=%v", ss.SelectedFaceDepthMode())
	}
	ss.SetSelectedFaceDepthMode(SELECTED_FACE_DEPTH_MODE_FRONT)
	if ss.SelectedVertexDepthMode() != SELECTED_VERTEX_DEPTH_MODE_FRONT {
		t.Errorf("vertex depth mode changed with face mode: got=%v", ss.SelectedVertexDepthMode())
	}
	ss.SetSelectedFaceMode(SELECTED_FACE_MODE_BOX)
	if ss.SelectedFaceMode() != SELECTED_FACE_MODE_BOX {
		t.Errorf("SelectedFaceMode setter mismatch: got=%v", ss.SelectedFaceMode())
	}
	ss.SetSelectedFaceMode(SELECTED_FACE_MODE_LINE)
	if ss.SelectedFaceMode() != SELECTED_FACE_MODE_LINE {
		t.Errorf("SelectedFaceMode roundtrip mismatch: got=%v", ss.SelectedFaceMode())
	}

	phys, ok := ss.PhysicsWorldMotion(0).(*defaultMotion)
	if !ok {
		t.Fatalf("PhysicsWorldMotion type mismatch")
	}
	if phys.Gravity != -9.8 || phys.MaxSubSteps != performance.DefaultMaxSubSteps || phys.FixedTimeStep != 60 {
		t.Errorf("Physics defaults: got=%v", phys)
	}
	wind, ok := ss.WindMotion(0).(*defaultMotion)
	if !ok {
		t.Fatalf("WindMotion type mismatch")
	}
	if wind.WindEnabled || wind.WindDirection != [3]float32{0, 0, 0} {
		t.Errorf("Wind defaults: got=%v", wind)
	}
	if ss.CameraMotion(0) != nil {
		t.Errorf("CameraMotion default should be nil")
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
	ss.EnableFlag(STATE_FLAG_SHOW_SELECTED_VERTEX)
	ss.EnableFlag(STATE_FLAG_SHOW_SELECTED_FACE)
	if ss.HasFlag(STATE_FLAG_SHOW_SELECTED_VERTEX) || !ss.HasFlag(STATE_FLAG_SHOW_SELECTED_FACE) {
		t.Errorf("face selection should be exclusive with vertex selection")
	}
	ss.EnableFlag(STATE_FLAG_SHOW_SELECTED_VERTEX)
	if !ss.HasFlag(STATE_FLAG_SHOW_SELECTED_VERTEX) || ss.HasFlag(STATE_FLAG_SHOW_SELECTED_FACE) {
		t.Errorf("vertex selection should be exclusive with face selection")
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
	faceIdxs := ss.SelectedFaceIndexes(0, 0)
	if faceIdxs == nil {
		t.Errorf("SelectedFaceIndexes should not be nil")
	}
	if len(faceIdxs) != 0 {
		t.Errorf("SelectedFaceIndexes length: got=%v", len(faceIdxs))
	}
	ss.SetSelectedFaceIndexes(0, 0, []int{7, 8})
	faceIdxs, version := ss.SelectedFaceIndexesWithVersion(0, 0)
	if len(faceIdxs) != 2 || faceIdxs[0] != 7 || version == 0 {
		t.Errorf("SelectedFaceIndexesWithVersion mismatch: indexes=%v version=%d", faceIdxs, version)
	}
	faceCopy := ss.SelectedFaceIndexes(0, 0)
	faceCopy[0] = 99
	if ss.SelectedFaceIndexes(0, 0)[0] == 99 {
		t.Errorf("SelectedFaceIndexes should be cloned")
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
