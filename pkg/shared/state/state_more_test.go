// 指示: miu200521358
package state

import (
	"sync/atomic"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

type testMotion struct {
	*hashable.HashableBase
}

func newTestModel() IStateModel {
	return &testModel{HashableBase: hashable.NewHashableBase("model", "path")}
}

func newTestMotion() IStateMotion {
	return &testMotion{HashableBase: hashable.NewHashableBase("motion", "path")}
}

// TestSharedStateFlagsAndWindows はフラグとウィンドウ系を確認する。
func TestSharedStateFlagsAndWindows(t *testing.T) {
	ss := NewSharedState(2)
	_ = ss.Flags()
	ss.SetFlags(StateFlagSet(0))
	ss.EnableFlag(STATE_FLAG_SHOW_BONE_ALL)
	if !ss.IsAnyBoneVisible() {
		t.Errorf("IsAnyBoneVisible expected true")
	}
	ss.DisableFlag(STATE_FLAG_SHOW_BONE_ALL)
	if ss.IsAnyBoneVisible() {
		t.Errorf("IsAnyBoneVisible expected false")
	}

	ss.EnableFlag(STATE_FLAG_SHOW_OVERRIDE_UPPER)
	if !ss.IsShowOverride() {
		t.Errorf("IsShowOverride expected true")
	}
	ss.DisableFlag(STATE_FLAG_SHOW_OVERRIDE_UPPER)
	if ss.IsShowOverride() {
		t.Errorf("IsShowOverride expected false")
	}

	ss.SetPlayback(PlaybackState{Frame: 2, MaxFrame: 3, FrameInterval: 0.5, Playing: true})
	p := ss.Playback()
	if p.Frame != 2 || p.MaxFrame != 3 || !p.Playing {
		t.Errorf("Playback: got=%v", p)
	}
	ss.SetPlayback(PlaybackState{Frame: 1, MaxFrame: 2, FrameInterval: 1, Playing: false})
	if ss.HasFlag(STATE_FLAG_PLAYING) {
		t.Errorf("Playing flag should be false")
	}

	ss.SetControlWindowPosition(WindowPosition{X: 1, Y: 2})
	if ss.ControlWindowPosition().X != 1 {
		t.Errorf("ControlWindowPosition mismatch")
	}
	ss.SetControlWindowHandle(10)
	ss.SetViewerWindowHandle(0, 20)
	ss.SetViewerWindowHandle(1, 21)
	ss.SetViewerWindowHandle(2, 22)
	if ss.ViewerWindowHandle(-1) != 0 || ss.ViewerWindowHandle(2) != 0 {
		t.Errorf("ViewerWindowHandle out-of-range should be 0")
	}
	if !ss.IsKnownWindowHandle(10) || !ss.IsKnownWindowHandle(21) || ss.IsKnownWindowHandle(99) {
		t.Errorf("IsKnownWindowHandle mismatch")
	}

	ss.SetFocusedWindowHandle(21)
	if ss.FocusedWindowHandle() != 21 {
		t.Errorf("FocusedWindowHandle mismatch")
	}

	ss.SetControlWindowReady(true)
	if !ss.IsControlWindowReady() {
		t.Errorf("ControlWindowReady mismatch")
	}
	ss.SetViewerWindowReady(0, true)
	ss.SetViewerWindowReady(1, false)
	ss.SetViewerWindowReady(-1, true)
	if !ss.IsViewerWindowReady(0) {
		t.Errorf("IsViewerWindowReady mismatch")
	}
	if ss.IsAllViewerWindowsReady() {
		t.Errorf("IsAllViewerWindowsReady expected false")
	}
	ss.SetViewerWindowReady(1, true)
	if !ss.IsAllViewerWindowsReady() {
		t.Errorf("IsAllViewerWindowsReady expected true")
	}
	if ss.IsViewerWindowReady(-1) {
		t.Errorf("IsViewerWindowReady out-of-range should be false")
	}

	ss.SetControlWindowFocused(true)
	if !ss.IsControlWindowFocused() {
		t.Errorf("ControlWindowFocused mismatch")
	}
	ss.SetViewerWindowFocused(0, true)
	if !ss.IsViewerWindowFocused(0) {
		t.Errorf("ViewerWindowFocused mismatch")
	}
	ss.SetViewerWindowFocused(-1, true)
	if ss.IsViewerWindowFocused(-1) {
		t.Errorf("ViewerWindowFocused out-of-range should be false")
	}
	ss.SetAllViewerWindowsFocused(false)
	if ss.IsViewerWindowFocused(0) || ss.IsViewerWindowFocused(1) {
		t.Errorf("SetAllViewerWindowsFocused failed")
	}

	ss.SetFpsLimitTriggered(true)
	if !ss.IsFpsLimitTriggered() {
		t.Errorf("FpsLimitTriggered mismatch")
	}
	ss.SetControlWindowMoving(true)
	if !ss.IsControlWindowMoving() {
		t.Errorf("ControlWindowMoving mismatch")
	}
	ss.SetClosed(true)
	if !ss.IsClosed() {
		t.Errorf("Closed mismatch")
	}

	ss.SetFocusLinkEnabled(false)
	ss.TriggerLinkedFocus(0)
	ss.SetFocusLinkEnabled(true)
	ss.SetViewerWindowHandle(0, 30)
	ss.TriggerLinkedFocus(0)
	if ss.FocusedWindowHandle() != 30 || !ss.IsViewerWindowFocused(0) || ss.IsControlWindowFocused() {
		t.Errorf("TriggerLinkedFocus mismatch")
	}
	ss.TriggerLinkedFocus(99)
	ss.KeepFocus()

	ss.SyncMinimize(-1)
	ss.SyncRestore(-1)
	ss.SetViewerWindowFocused(0, true)
	ss.SyncMinimize(0)
	if ss.IsViewerWindowFocused(0) {
		t.Errorf("SyncMinimize failed")
	}
	ss.SyncRestore(0)
	if !ss.IsViewerWindowFocused(0) {
		t.Errorf("SyncRestore failed")
	}
}

// TestSharedStateModelsAndSelections はモデル/モーションを確認する。
func TestSharedStateModelsAndSelections(t *testing.T) {
	ss := NewSharedState(1)
	ss.SetModel(-1, 0, newTestModel())
	ss.SetMotion(-1, 0, newTestMotion())

	ss.SetModel(0, 0, newTestModel())
	ss.SetMotion(0, 0, newTestMotion())
	if ss.Model(0, 0) == nil || ss.Motion(0, 0) == nil {
		t.Errorf("Model/Motion should be set")
	}
	if ss.Model(-1, 0) != nil || ss.Motion(-1, 0) != nil {
		t.Errorf("Out-of-range viewer should be nil")
	}
	if ss.Model(0, 1) != nil || ss.Motion(0, 1) != nil {
		t.Errorf("Out-of-range should be nil")
	}
	if ss.ModelCount(0) != 1 || ss.ModelCount(-1) != 0 {
		t.Errorf("ModelCount mismatch")
	}

	ss.SetModel(0, 1, newTestModel())
	if ss.MotionCount(0) != 2 {
		t.Errorf("MotionCount should follow models")
	}
	if ss.MotionCount(-1) != 0 {
		t.Errorf("MotionCount invalid index should be 0")
	}

	if ss.SelectedMaterialIndexes(-1, 0) != nil {
		t.Errorf("SelectedMaterialIndexes out-of-range should be nil")
	}
	if ss.SelectedMaterialIndexes(0, 99) != nil {
		t.Errorf("SelectedMaterialIndexes model out-of-range should be nil")
	}
	ss.SetSelectedMaterialIndexes(-1, 0, []int{1})
	ss.SetSelectedMaterialIndexes(0, 0, []int{1, 2})
	idxs := ss.SelectedMaterialIndexes(0, 0)
	if len(idxs) != 2 || idxs[0] != 1 {
		t.Errorf("SelectedMaterialIndexes mismatch: %v", idxs)
	}
	idxs[0] = 99
	if ss.SelectedMaterialIndexes(0, 0)[0] == 99 {
		t.Errorf("SelectedMaterialIndexes should be cloned")
	}

	if ss.SelectedVertexIndexes(-1, 0) != nil {
		t.Errorf("SelectedVertexIndexes out-of-range should be nil")
	}
	if ss.SelectedVertexIndexes(0, 99) != nil {
		t.Errorf("SelectedVertexIndexes model out-of-range should be nil")
	}
	ss.SetSelectedVertexIndexes(-1, 0, []int{3})
	ss.SetSelectedVertexIndexes(0, 0, []int{3, 4})
	vertexIdxs := ss.SelectedVertexIndexes(0, 0)
	if len(vertexIdxs) != 2 || vertexIdxs[0] != 3 {
		t.Errorf("SelectedVertexIndexes mismatch: %v", vertexIdxs)
	}
	vertexIdxs[0] = 99
	if ss.SelectedVertexIndexes(0, 0)[0] == 99 {
		t.Errorf("SelectedVertexIndexes should be cloned")
	}

	impl := ss
	impl.models[0] = make([]atomic.Value, 2)
	impl.motions[0] = make([]atomic.Value, 1)
	if impl.MotionCount(0) != 2 {
		t.Errorf("MotionCount length branch failed")
	}
}

// TestSharedStateTrajectoryPolylines は軌跡の viewer 分離と複製境界を確認する。
func TestSharedStateTrajectoryPolylines(t *testing.T) {
	ss := NewSharedState(2)
	src := []TrajectoryPolyline{{
		Width: 2, GroundedWidth: 7, CurrentMarkerSize: 11,
		Points: []TrajectoryPoint{{
			Frame: 3, Position: [3]float32{1, 2, 3},
			Color: TrajectoryColor{R: 1, A: 1}, Grounded: true,
			GroundColor: TrajectoryColor{G: 1, A: 1}, Current: true,
		}},
	}}

	ss.SetTrajectoryPolylines(1, src)
	src[0].Width = 99
	src[0].Points[0].Position[0] = 99
	got := ss.TrajectoryPolylines(1)
	if len(got) != 1 || got[0].Width != 2 || got[0].Points[0].Position[0] != 1 {
		t.Fatalf("設定時の複製が保持されていません: %+v", got)
	}
	if ss.TrajectoryPolylines(0) != nil {
		t.Fatal("別 viewer に軌跡が漏れています")
	}

	got[0].Points[0].Color.R = 0
	if ss.TrajectoryPolylines(1)[0].Points[0].Color.R != 1 {
		t.Fatal("取得値の変更が共有状態へ漏れています")
	}
	_, versionBefore := ss.TrajectoryPolylinesWithVersion(1)
	ss.ClearTrajectoryPolylines(1)
	cleared, versionAfter := ss.TrajectoryPolylinesWithVersion(1)
	if cleared != nil || versionAfter <= versionBefore {
		t.Fatalf("clear/version = %v/%d, before=%d", cleared, versionAfter, versionBefore)
	}
	ss.SetTrajectoryPolylines(-1, src)
	if ss.TrajectoryPolylines(-1) != nil {
		t.Fatal("範囲外 viewer は空である必要があります")
	}
}

// TestSharedStateOperationPointHandles は操作点の viewer 分離、複製、消去を確認する。
func TestSharedStateOperationPointHandles(t *testing.T) {
	ss := NewSharedState(2)
	src := []OperationPointHandle{{ModelIndex: 0, BoneName: "首"}}
	ss.SetOperationPointHandles(1, src)
	_, handleVersion := ss.OperationPointHandlesWithVersion(1)
	src[0].BoneName = "変更済み"

	got := ss.OperationPointHandles(1)
	if len(got) != 1 || got[0].BoneName != "首" {
		t.Fatalf("設定時の複製が保持されていません: %+v", got)
	}
	if ss.OperationPointHandles(0) != nil {
		t.Fatal("別 viewer に操作点が漏れています")
	}
	got[0].BoneName = "取得後変更"
	if ss.OperationPointHandles(1)[0].BoneName != "首" {
		t.Fatal("取得値の変更が共有状態へ漏れています")
	}

	baseEvent := OperationPointDragEvent{ModelIndex: 0, BoneName: "首", HandleVersion: handleVersion}
	grabbed := baseEvent
	grabbed.Phase = OPERATION_POINT_DRAG_PHASE_GRABBED
	ss.PublishOperationPointDragEvent(1, grabbed)
	moved := baseEvent
	moved.Phase, moved.WorldPosition = OPERATION_POINT_DRAG_PHASE_MOVED, [3]float64{1, 2, 3}
	ss.PublishOperationPointDragEvent(1, moved)
	moved.WorldPosition = [3]float64{4, 5, 6}
	ss.PublishOperationPointDragEvent(1, moved)
	released := baseEvent
	released.Phase, released.WorldPosition = OPERATION_POINT_DRAG_PHASE_RELEASED, [3]float64{4, 5, 6}
	ss.PublishOperationPointDragEvent(1, released)
	events := ss.DrainOperationPointDragEvents(1)
	if len(events) != 3 || events[0].Phase != OPERATION_POINT_DRAG_PHASE_GRABBED ||
		events[1].Phase != OPERATION_POINT_DRAG_PHASE_MOVED || events[1].WorldPosition != [3]float64{4, 5, 6} ||
		events[2].Phase != OPERATION_POINT_DRAG_PHASE_RELEASED {
		t.Fatalf("通知キューの move 集約が不正です: %+v", events)
	}
	if ss.DrainOperationPointDragEvents(1) != nil {
		t.Fatal("取得済み通知が残っています")
	}
	if ss.DrainOperationPointDragEvents(0) != nil {
		t.Fatal("別 viewer に通知が漏れています")
	}

	ss.PublishOperationPointDragEvent(1, moved)
	_, versionBefore := ss.OperationPointHandlesWithVersion(1)
	ss.ClearOperationPointHandles(1)
	ss.PublishOperationPointDragEvent(1, released)
	cleared, versionAfter := ss.OperationPointHandlesWithVersion(1)
	remaining := ss.DrainOperationPointDragEvents(1)
	if cleared != nil || versionAfter <= versionBefore || remaining != nil {
		t.Fatalf("clear/version/events = %v/%d/%v, before=%d", cleared, versionAfter, remaining, versionBefore)
	}
	ss.SetOperationPointHandles(-1, src)
	if ss.OperationPointHandles(-1) != nil {
		t.Fatal("範囲外 viewer は空である必要があります")
	}
}

// TestSharedStateDeltaAndPhysics は差分/物理系を確認する。
func TestSharedStateDeltaAndPhysics(t *testing.T) {
	ss := NewSharedState(1)
	ss.SetDeltaSaveEnabled(-1, true)
	if ss.IsDeltaSaveEnabled(-1) {
		t.Errorf("IsDeltaSaveEnabled out-of-range should be false")
	}
	ss.SetDeltaSaveIndex(-1, 1)
	ss.SetDeltaSaveEnabled(0, true)
	if !ss.IsDeltaSaveEnabled(0) {
		t.Errorf("IsDeltaSaveEnabled mismatch")
	}
	ss.SetDeltaSaveIndex(0, 2)
	if ss.DeltaSaveIndex(0) != 2 || ss.DeltaSaveIndex(-1) != 0 {
		t.Errorf("DeltaSaveIndex mismatch")
	}

	ss.SetDeltaMotion(-1, 0, 0, newTestMotion())
	if ss.DeltaMotion(-1, 0, 0) != nil {
		t.Errorf("DeltaMotion out-of-range should be nil")
	}
	if ss.DeltaMotion(0, -1, 0) != nil {
		t.Errorf("DeltaMotion model out-of-range should be nil")
	}
	ss.SetDeltaMotion(0, 1, 1, newTestMotion())
	if ss.DeltaMotionCount(0, 1) != 2 {
		t.Errorf("DeltaMotionCount mismatch: %v", ss.DeltaMotionCount(0, 1))
	}
	if ss.DeltaMotion(0, 1, 99) != nil {
		t.Errorf("DeltaMotion out-of-range should be nil")
	}
	if ss.DeltaMotion(0, 1, 1) == nil {
		t.Errorf("DeltaMotion should be set")
	}
	ss.ClearDeltaMotion(0, 1)
	if ss.DeltaMotionCount(0, 1) != 0 {
		t.Errorf("ClearDeltaMotion failed")
	}
	ss.ClearDeltaMotion(0, -1)
	if ss.DeltaMotionCount(-1, 0) != 0 {
		t.Errorf("DeltaMotionCount out-of-range should be 0")
	}
	if ss.DeltaMotionCount(0, -1) != 0 {
		t.Errorf("DeltaMotionCount model out-of-range should be 0")
	}
	ss.ClearDeltaMotion(-1, 0)

	if ss.PhysicsWorldMotion(-1) != nil {
		t.Errorf("PhysicsWorldMotion out-of-range should be nil")
	}
	ss.SetPhysicsWorldMotion(-1, newTestMotion())
	ss.SetPhysicsWorldMotion(0, newTestMotion())
	if ss.PhysicsWorldMotion(0) == nil {
		t.Errorf("PhysicsWorldMotion should be set")
	}

	if ss.PhysicsModelMotion(-1, 0) != nil {
		t.Errorf("PhysicsModelMotion out-of-range should be nil")
	}
	if ss.PhysicsModelMotion(0, -1) != nil {
		t.Errorf("PhysicsModelMotion model out-of-range should be nil")
	}
	ss.SetPhysicsModelMotion(-1, 0, newTestMotion())
	ss.SetPhysicsModelMotion(0, 0, newTestMotion())
	if ss.PhysicsModelMotion(0, 0) == nil {
		t.Errorf("PhysicsModelMotion should be set")
	}

	if ss.WindMotion(-1) != nil {
		t.Errorf("WindMotion out-of-range should be nil")
	}
	ss.SetWindMotion(-1, newTestMotion())
	ss.SetWindMotion(0, newTestMotion())
	if ss.WindMotion(0) == nil {
		t.Errorf("WindMotion should be set")
	}

	if ss.CameraMotion(-1) != nil {
		t.Errorf("CameraMotion out-of-range should be nil")
	}
	ss.SetCameraMotion(-1, newTestMotion())
	ss.SetCameraMotion(0, newTestMotion())
	if ss.CameraMotion(0) == nil {
		t.Errorf("CameraMotion should be set")
	}

	if ss.PhysicsResetType() != PHYSICS_RESET_TYPE_NONE {
		t.Errorf("PhysicsResetType default mismatch")
	}
	ss.SetPhysicsResetType(PHYSICS_RESET_TYPE_START_FRAME)
	if ss.PhysicsResetType() != PHYSICS_RESET_TYPE_START_FRAME {
		t.Errorf("PhysicsResetType set mismatch")
	}
}

// TestSharedStateHelpers は補助関数を確認する。
func TestSharedStateHelpers(t *testing.T) {
	if got := ensureSlotSlice[struct{}](nil, -1, struct{}{}); got != nil {
		t.Errorf("ensureSlotSlice negative index should keep nil")
	}
	values := []atomic.Value{}
	values = ensureSlotSlice(values, 1, stateIndexSlot{Indexes: []int{}})
	if len(values) != 2 {
		t.Errorf("ensureSlotSlice size: got=%v", len(values))
	}

	if cloneIntSlice(nil) != nil {
		t.Errorf("cloneIntSlice nil should be nil")
	}
	src := []int{1, 2}
	dst := cloneIntSlice(src)
	dst[0] = 9
	if src[0] == 9 {
		t.Errorf("cloneIntSlice should copy")
	}

	phys := newDefaultPhysicsWorldMotion().(*defaultMotion)
	if phys.GetHashParts() != "" || phys.Gravity != -9.8 {
		t.Errorf("newDefaultPhysicsWorldMotion mismatch: %v", phys)
	}
	wind := newDefaultWindMotion().(*defaultMotion)
	if wind.GetHashParts() != "" || wind.WindEnabled {
		t.Errorf("newDefaultWindMotion mismatch: %v", wind)
	}

	ss := NewSharedState(1)
	ss.SetFrame(mtime.Frame(5))
	if ss.Frame() != 5 {
		t.Errorf("Frame mismatch")
	}
}
