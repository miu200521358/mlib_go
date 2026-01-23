// 指示: miu200521358
package state

import (
	"sync"
	"sync/atomic"

	"github.com/miu200521358/mlib_go/pkg/adapter/physics_api"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// WindowHandle はウィンドウハンドル型。
type WindowHandle int32

// WindowPosition はウィンドウ位置を表す。
type WindowPosition struct {
	X     int
	Y     int
	DiffX int
	DiffY int
}

// StateFlag は状態フラグのビットを表す。
type StateFlag uint64

const (
	// STATE_FLAG_FRAME_DROP はフレームドロップ。
	STATE_FLAG_FRAME_DROP StateFlag = 1 << iota
	// STATE_FLAG_PHYSICS_ENABLED は物理ON/OFF。
	STATE_FLAG_PHYSICS_ENABLED
	// STATE_FLAG_SHOW_NORMAL は法線表示。
	STATE_FLAG_SHOW_NORMAL
	// STATE_FLAG_SHOW_WIRE はワイヤーフレーム表示。
	STATE_FLAG_SHOW_WIRE
	// STATE_FLAG_SHOW_OVERRIDE_UPPER は上半身合わせ。
	STATE_FLAG_SHOW_OVERRIDE_UPPER
	// STATE_FLAG_SHOW_OVERRIDE_LOWER は下半身合わせ。
	STATE_FLAG_SHOW_OVERRIDE_LOWER
	// STATE_FLAG_SHOW_OVERRIDE_NONE はカメラ合わせなし。
	STATE_FLAG_SHOW_OVERRIDE_NONE
	// STATE_FLAG_SHOW_SELECTED_VERTEX は選択頂点表示。
	STATE_FLAG_SHOW_SELECTED_VERTEX
	// STATE_FLAG_SHOW_BONE_ALL は全ボーン表示。
	STATE_FLAG_SHOW_BONE_ALL
	// STATE_FLAG_SHOW_BONE_IK はIKボーン表示。
	STATE_FLAG_SHOW_BONE_IK
	// STATE_FLAG_SHOW_BONE_EFFECTOR は付与親ボーン表示。
	STATE_FLAG_SHOW_BONE_EFFECTOR
	// STATE_FLAG_SHOW_BONE_FIXED は軸制限ボーン表示。
	STATE_FLAG_SHOW_BONE_FIXED
	// STATE_FLAG_SHOW_BONE_ROTATE は回転ボーン表示。
	STATE_FLAG_SHOW_BONE_ROTATE
	// STATE_FLAG_SHOW_BONE_TRANSLATE は移動ボーン表示。
	STATE_FLAG_SHOW_BONE_TRANSLATE
	// STATE_FLAG_SHOW_BONE_VISIBLE は表示ボーン表示。
	STATE_FLAG_SHOW_BONE_VISIBLE
	// STATE_FLAG_SHOW_RIGID_BODY_FRONT は剛体前面表示。
	STATE_FLAG_SHOW_RIGID_BODY_FRONT
	// STATE_FLAG_SHOW_RIGID_BODY_BACK は剛体埋め込み表示。
	STATE_FLAG_SHOW_RIGID_BODY_BACK
	// STATE_FLAG_SHOW_JOINT はジョイント表示。
	STATE_FLAG_SHOW_JOINT
	// STATE_FLAG_SHOW_INFO は情報表示。
	STATE_FLAG_SHOW_INFO
	// STATE_FLAG_CAMERA_SYNC はカメラ同期。
	STATE_FLAG_CAMERA_SYNC
	// STATE_FLAG_PLAYING は再生中フラグ。
	STATE_FLAG_PLAYING
	// STATE_FLAG_WINDOW_LINKAGE はウィンドウ連動。
	STATE_FLAG_WINDOW_LINKAGE
	// STATE_FLAG_CHANGED_ENABLE_FRAME_DROP はフレームドロップ変更通知。
	STATE_FLAG_CHANGED_ENABLE_FRAME_DROP
)

// StateFlagSet はフラグ集合。
type StateFlagSet uint64

// PlaybackState は再生状態。
type PlaybackState struct {
	Frame         mtime.Frame
	MaxFrame      mtime.Frame
	FrameInterval mtime.Seconds
	Playing       bool
}

// PhysicsResetType は物理リセット種別。
type PhysicsResetType int

const (
	// PHYSICS_RESET_TYPE_NONE はリセットなし。
	PHYSICS_RESET_TYPE_NONE PhysicsResetType = iota
	// PHYSICS_RESET_TYPE_CONTINUE_FRAME は現在フレーム継続。
	PHYSICS_RESET_TYPE_CONTINUE_FRAME
	// PHYSICS_RESET_TYPE_START_FRAME は開始フレームへ戻す。
	PHYSICS_RESET_TYPE_START_FRAME
	// PHYSICS_RESET_TYPE_START_FIT_FRAME は開始フレームへフィット。
	PHYSICS_RESET_TYPE_START_FIT_FRAME
)

// IStateModel は共有状態で扱うモデルI/F。
type IStateModel interface {
	hashable.IHashable
}

// IStateMotion は共有状態で扱うモーションI/F。
type IStateMotion interface {
	hashable.IHashable
}

// ISharedState は共有状態I/F。
type ISharedState interface {
	Flags() StateFlagSet
	SetFlags(flags StateFlagSet)
	EnableFlag(flag StateFlag)
	DisableFlag(flag StateFlag)
	HasFlag(flag StateFlag) bool
	IsAnyBoneVisible() bool
	IsShowOverride() bool
	Playback() PlaybackState
	SetPlayback(p PlaybackState)
	Frame() mtime.Frame
	SetFrame(frame mtime.Frame)
	MaxFrame() mtime.Frame
	SetMaxFrame(maxFrame mtime.Frame)
	FrameInterval() mtime.Seconds
	SetFrameInterval(spf mtime.Seconds)
	ControlWindowPosition() WindowPosition
	SetControlWindowPosition(pos WindowPosition)
	ControlWindowHandle() WindowHandle
	SetControlWindowHandle(handle WindowHandle)
	ViewerWindowHandle(viewerIndex int) WindowHandle
	SetViewerWindowHandle(viewerIndex int, handle WindowHandle)
	IsKnownWindowHandle(handle WindowHandle) bool
	FocusedWindowHandle() WindowHandle
	SetFocusedWindowHandle(handle WindowHandle)
	IsControlWindowReady() bool
	SetControlWindowReady(ready bool)
	IsViewerWindowReady(viewerIndex int) bool
	SetViewerWindowReady(viewerIndex int, ready bool)
	IsAllViewerWindowsReady() bool
	IsControlWindowFocused() bool
	SetControlWindowFocused(focus bool)
	IsViewerWindowFocused(viewerIndex int) bool
	SetViewerWindowFocused(viewerIndex int, focus bool)
	SetAllViewerWindowsFocused(focus bool)
	IsFpsLimitTriggered() bool
	SetFpsLimitTriggered(triggered bool)
	IsControlWindowMoving() bool
	SetControlWindowMoving(moving bool)
	IsClosed() bool
	SetClosed(closed bool)
	IsFocusLinkEnabled() bool
	SetFocusLinkEnabled(enabled bool)
	TriggerLinkedFocus(viewerIndex int)
	KeepFocus()
	SyncMinimize(viewerIndex int)
	SyncRestore(viewerIndex int)
	ModelCount(viewerIndex int) int
	MotionCount(viewerIndex int) int
	SetModel(viewerIndex, modelIndex int, model IStateModel)
	Model(viewerIndex, modelIndex int) IStateModel
	SetMotion(viewerIndex, modelIndex int, motion IStateMotion)
	Motion(viewerIndex, modelIndex int) IStateMotion
	SelectedMaterialIndexes(viewerIndex, modelIndex int) []int
	SetSelectedMaterialIndexes(viewerIndex, modelIndex int, indexes []int)
	IsDeltaSaveEnabled(viewerIndex int) bool
	SetDeltaSaveEnabled(viewerIndex int, enabled bool)
	DeltaSaveIndex(viewerIndex int) int
	SetDeltaSaveIndex(viewerIndex int, index int)
	SetDeltaMotion(viewerIndex, modelIndex, motionIndex int, motion IStateMotion)
	DeltaMotion(viewerIndex, modelIndex, deltaIndex int) IStateMotion
	ClearDeltaMotion(viewerIndex, modelIndex int)
	DeltaMotionCount(viewerIndex, modelIndex int) int
	SetPhysicsWorldMotion(viewerIndex int, motion IStateMotion)
	PhysicsWorldMotion(viewerIndex int) IStateMotion
	SetPhysicsModelMotion(viewerIndex, modelIndex int, motion IStateMotion)
	PhysicsModelMotion(viewerIndex, modelIndex int) IStateMotion
	SetWindMotion(viewerIndex int, motion IStateMotion)
	WindMotion(viewerIndex int) IStateMotion
	PhysicsResetType() PhysicsResetType
	SetPhysicsResetType(resetType PhysicsResetType)
}

type stateModelSlot struct {
	Model IStateModel
}

type stateMotionSlot struct {
	Motion IStateMotion
}

type stateIndexSlot struct {
	Indexes []int
}

// SharedState は共有状態の実装。
type SharedState struct {
	mu                    sync.Mutex
	flags                 atomic.Uint64
	frameValue            atomic.Value
	maxFrameValue         atomic.Value
	frameIntervalValue    atomic.Value
	controlWindowPosition atomic.Value
	controlWindowHandle   atomic.Int32
	viewerWindowHandles   []atomic.Int32
	focusedWindowHandle   atomic.Int32
	controlWindowReady    atomic.Bool
	viewerWindowReady     []atomic.Bool
	controlWindowFocused  atomic.Bool
	viewerWindowFocused   []atomic.Bool
	fpsLimitTriggered     atomic.Bool
	controlWindowMoving   atomic.Bool
	closed                atomic.Bool
	focusLinkEnabled      atomic.Bool
	linkingFocus          atomic.Bool
	models                [][]atomic.Value
	motions               [][]atomic.Value
	selectedIndexes       [][]atomic.Value
	deltaSaveEnabled      []atomic.Bool
	deltaSaveIndexes      []atomic.Int32
	deltaMotions          [][][]atomic.Value
	physicsWorldMotions   []atomic.Value
	physicsModelMotions   [][]atomic.Value
	windMotions           []atomic.Value
	physicsResetType      atomic.Int32
}

// NewSharedState は共有状態を生成する。
func NewSharedState(viewerCount int) ISharedState {
	ss := &SharedState{
		viewerWindowHandles: make([]atomic.Int32, viewerCount),
		viewerWindowReady:   make([]atomic.Bool, viewerCount),
		viewerWindowFocused: make([]atomic.Bool, viewerCount),
		models:              make([][]atomic.Value, viewerCount),
		motions:             make([][]atomic.Value, viewerCount),
		selectedIndexes:     make([][]atomic.Value, viewerCount),
		deltaSaveEnabled:    make([]atomic.Bool, viewerCount),
		deltaSaveIndexes:    make([]atomic.Int32, viewerCount),
		deltaMotions:        make([][][]atomic.Value, viewerCount),
		physicsWorldMotions: make([]atomic.Value, viewerCount),
		physicsModelMotions: make([][]atomic.Value, viewerCount),
		windMotions:         make([]atomic.Value, viewerCount),
	}

	ss.frameValue.Store(mtime.Frame(0))
	ss.maxFrameValue.Store(mtime.Frame(1))
	ss.frameIntervalValue.Store(mtime.Seconds(-1))
	ss.controlWindowPosition.Store(WindowPosition{})
	ss.focusedWindowHandle.Store(int32(0))
	ss.controlWindowHandle.Store(int32(0))
	ss.controlWindowReady.Store(false)
	ss.controlWindowFocused.Store(false)
	ss.fpsLimitTriggered.Store(false)
	ss.controlWindowMoving.Store(false)
	ss.closed.Store(false)
	ss.focusLinkEnabled.Store(true)
	ss.linkingFocus.Store(false)
	ss.physicsResetType.Store(int32(PHYSICS_RESET_TYPE_NONE))

	for i := 0; i < viewerCount; i++ {
		ss.physicsWorldMotions[i].Store(stateMotionSlot{Motion: newDefaultPhysicsWorldMotion()})
		ss.windMotions[i].Store(stateMotionSlot{Motion: newDefaultWindMotion()})
	}

	return ss
}

// Flags は現在のフラグ集合を返す。
func (ss *SharedState) Flags() StateFlagSet {
	return StateFlagSet(ss.flags.Load())
}

// SetFlags はフラグ集合を置換する。
func (ss *SharedState) SetFlags(flags StateFlagSet) {
	for {
		current := ss.flags.Load()
		if ss.flags.CompareAndSwap(current, uint64(flags)) {
			return
		}
	}
}

// EnableFlag はフラグを有効化する。
func (ss *SharedState) EnableFlag(flag StateFlag) {
	for {
		current := ss.flags.Load()
		next := current | uint64(flag)
		if ss.flags.CompareAndSwap(current, next) {
			return
		}
	}
}

// DisableFlag はフラグを無効化する。
func (ss *SharedState) DisableFlag(flag StateFlag) {
	for {
		current := ss.flags.Load()
		next := current &^ uint64(flag)
		if ss.flags.CompareAndSwap(current, next) {
			return
		}
	}
}

// HasFlag はフラグの有無を判定する。
func (ss *SharedState) HasFlag(flag StateFlag) bool {
	return ss.flags.Load()&uint64(flag) != 0
}

// IsAnyBoneVisible はボーン表示系フラグのいずれかが有効か判定する。
func (ss *SharedState) IsAnyBoneVisible() bool {
	return ss.HasFlag(STATE_FLAG_SHOW_BONE_ALL) ||
		ss.HasFlag(STATE_FLAG_SHOW_BONE_IK) ||
		ss.HasFlag(STATE_FLAG_SHOW_BONE_EFFECTOR) ||
		ss.HasFlag(STATE_FLAG_SHOW_BONE_FIXED) ||
		ss.HasFlag(STATE_FLAG_SHOW_BONE_ROTATE) ||
		ss.HasFlag(STATE_FLAG_SHOW_BONE_TRANSLATE) ||
		ss.HasFlag(STATE_FLAG_SHOW_BONE_VISIBLE)
}

// IsShowOverride はオーバーライド表示が有効か判定する。
func (ss *SharedState) IsShowOverride() bool {
	return ss.HasFlag(STATE_FLAG_SHOW_OVERRIDE_UPPER) ||
		ss.HasFlag(STATE_FLAG_SHOW_OVERRIDE_LOWER) ||
		ss.HasFlag(STATE_FLAG_SHOW_OVERRIDE_NONE)
}

// Playback は再生状態を返す。
func (ss *SharedState) Playback() PlaybackState {
	return PlaybackState{
		Frame:         ss.Frame(),
		MaxFrame:      ss.MaxFrame(),
		FrameInterval: ss.FrameInterval(),
		Playing:       ss.HasFlag(STATE_FLAG_PLAYING),
	}
}

// SetPlayback は再生状態を設定する。
func (ss *SharedState) SetPlayback(p PlaybackState) {
	ss.SetFrame(p.Frame)
	ss.SetMaxFrame(p.MaxFrame)
	ss.SetFrameInterval(p.FrameInterval)
	if p.Playing {
		ss.EnableFlag(STATE_FLAG_PLAYING)
	} else {
		ss.DisableFlag(STATE_FLAG_PLAYING)
	}
}

// Frame は現在フレームを返す。
func (ss *SharedState) Frame() mtime.Frame {
	return ss.frameValue.Load().(mtime.Frame)
}

// SetFrame は現在フレームを設定する。
func (ss *SharedState) SetFrame(frame mtime.Frame) {
	ss.frameValue.Store(frame)
}

// MaxFrame は最大フレームを返す。
func (ss *SharedState) MaxFrame() mtime.Frame {
	return ss.maxFrameValue.Load().(mtime.Frame)
}

// SetMaxFrame は最大フレームを設定する。
func (ss *SharedState) SetMaxFrame(maxFrame mtime.Frame) {
	ss.maxFrameValue.Store(maxFrame)
}

// FrameInterval はフレーム間隔を返す。
func (ss *SharedState) FrameInterval() mtime.Seconds {
	return ss.frameIntervalValue.Load().(mtime.Seconds)
}

// SetFrameInterval はフレーム間隔を設定する。
func (ss *SharedState) SetFrameInterval(spf mtime.Seconds) {
	ss.frameIntervalValue.Store(spf)
}

// ControlWindowPosition はコントロールウィンドウ位置を返す。
func (ss *SharedState) ControlWindowPosition() WindowPosition {
	return ss.controlWindowPosition.Load().(WindowPosition)
}

// SetControlWindowPosition はコントロールウィンドウ位置を設定する。
func (ss *SharedState) SetControlWindowPosition(pos WindowPosition) {
	ss.controlWindowPosition.Store(pos)
}

// ControlWindowHandle はコントロールウィンドウハンドルを返す。
func (ss *SharedState) ControlWindowHandle() WindowHandle {
	return WindowHandle(ss.controlWindowHandle.Load())
}

// SetControlWindowHandle はコントロールウィンドウハンドルを設定する。
func (ss *SharedState) SetControlWindowHandle(handle WindowHandle) {
	ss.controlWindowHandle.Store(int32(handle))
}

// ViewerWindowHandle はビューワーハンドルを返す。
func (ss *SharedState) ViewerWindowHandle(viewerIndex int) WindowHandle {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowHandles) {
		return 0
	}
	return WindowHandle(ss.viewerWindowHandles[viewerIndex].Load())
}

// SetViewerWindowHandle はビューワーハンドルを設定する。
func (ss *SharedState) SetViewerWindowHandle(viewerIndex int, handle WindowHandle) {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowHandles) {
		return
	}
	ss.viewerWindowHandles[viewerIndex].Store(int32(handle))
}

// IsKnownWindowHandle は既知ハンドルか判定する。
func (ss *SharedState) IsKnownWindowHandle(handle WindowHandle) bool {
	if ss.ControlWindowHandle() == handle {
		return true
	}
	for i := 0; i < len(ss.viewerWindowHandles); i++ {
		if ss.ViewerWindowHandle(i) == handle {
			return true
		}
	}
	return false
}

// FocusedWindowHandle はフォーカス中ハンドルを返す。
func (ss *SharedState) FocusedWindowHandle() WindowHandle {
	return WindowHandle(ss.focusedWindowHandle.Load())
}

// SetFocusedWindowHandle はフォーカス中ハンドルを設定する。
func (ss *SharedState) SetFocusedWindowHandle(handle WindowHandle) {
	ss.focusedWindowHandle.Store(int32(handle))
}

// IsControlWindowReady はコントロールウィンドウ準備完了か判定する。
func (ss *SharedState) IsControlWindowReady() bool {
	return ss.controlWindowReady.Load()
}

// SetControlWindowReady はコントロールウィンドウ準備完了を設定する。
func (ss *SharedState) SetControlWindowReady(ready bool) {
	ss.controlWindowReady.Store(ready)
}

// IsViewerWindowReady はビューワー準備完了か判定する。
func (ss *SharedState) IsViewerWindowReady(viewerIndex int) bool {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowReady) {
		return false
	}
	return ss.viewerWindowReady[viewerIndex].Load()
}

// SetViewerWindowReady はビューワー準備完了を設定する。
func (ss *SharedState) SetViewerWindowReady(viewerIndex int, ready bool) {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowReady) {
		return
	}
	ss.viewerWindowReady[viewerIndex].Store(ready)
}

// IsAllViewerWindowsReady は全ビューワー準備完了か判定する。
func (ss *SharedState) IsAllViewerWindowsReady() bool {
	for i := 0; i < len(ss.viewerWindowReady); i++ {
		if !ss.viewerWindowReady[i].Load() {
			return false
		}
	}
	return true
}

// IsControlWindowFocused はコントロールウィンドウがフォーカス中か判定する。
func (ss *SharedState) IsControlWindowFocused() bool {
	return ss.controlWindowFocused.Load()
}

// SetControlWindowFocused はコントロールウィンドウのフォーカス状態を設定する。
func (ss *SharedState) SetControlWindowFocused(focus bool) {
	ss.controlWindowFocused.Store(focus)
}

// IsViewerWindowFocused はビューワーがフォーカス中か判定する。
func (ss *SharedState) IsViewerWindowFocused(viewerIndex int) bool {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowFocused) {
		return false
	}
	return ss.viewerWindowFocused[viewerIndex].Load()
}

// SetViewerWindowFocused はビューワーのフォーカス状態を設定する。
func (ss *SharedState) SetViewerWindowFocused(viewerIndex int, focus bool) {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowFocused) {
		return
	}
	ss.viewerWindowFocused[viewerIndex].Store(focus)
}

// SetAllViewerWindowsFocused は全ビューワーのフォーカス状態を設定する。
func (ss *SharedState) SetAllViewerWindowsFocused(focus bool) {
	for i := 0; i < len(ss.viewerWindowFocused); i++ {
		ss.viewerWindowFocused[i].Store(focus)
	}
}

// IsFpsLimitTriggered はFPS制限トリガー状態を返す。
func (ss *SharedState) IsFpsLimitTriggered() bool {
	return ss.fpsLimitTriggered.Load()
}

// SetFpsLimitTriggered はFPS制限トリガー状態を設定する。
func (ss *SharedState) SetFpsLimitTriggered(triggered bool) {
	ss.fpsLimitTriggered.Store(triggered)
}

// IsControlWindowMoving はコントロールウィンドウ移動中か判定する。
func (ss *SharedState) IsControlWindowMoving() bool {
	return ss.controlWindowMoving.Load()
}

// SetControlWindowMoving はコントロールウィンドウ移動状態を設定する。
func (ss *SharedState) SetControlWindowMoving(moving bool) {
	ss.controlWindowMoving.Store(moving)
}

// IsClosed は閉じる状態か判定する。
func (ss *SharedState) IsClosed() bool {
	return ss.closed.Load()
}

// SetClosed は閉じる状態を設定する。
func (ss *SharedState) SetClosed(closed bool) {
	ss.closed.Store(closed)
}

// IsFocusLinkEnabled はフォーカス連動が有効か判定する。
func (ss *SharedState) IsFocusLinkEnabled() bool {
	return ss.focusLinkEnabled.Load()
}

// SetFocusLinkEnabled はフォーカス連動の有効化を設定する。
func (ss *SharedState) SetFocusLinkEnabled(enabled bool) {
	ss.focusLinkEnabled.Store(enabled)
}

// TriggerLinkedFocus は連動フォーカスを開始する。
func (ss *SharedState) TriggerLinkedFocus(viewerIndex int) {
	if !ss.IsFocusLinkEnabled() {
		return
	}
	ss.linkingFocus.Store(true)
	if viewerIndex >= 0 && viewerIndex < len(ss.viewerWindowHandles) {
		ss.SetFocusedWindowHandle(ss.ViewerWindowHandle(viewerIndex))
		ss.SetControlWindowFocused(false)
		ss.SetAllViewerWindowsFocused(false)
		ss.SetViewerWindowFocused(viewerIndex, true)
	}
}

// KeepFocus は連動フォーカス中の状態を維持する。
func (ss *SharedState) KeepFocus() {
	ss.linkingFocus.Store(false)
}

// SyncMinimize はウィンドウ最小化連動を反映する。
func (ss *SharedState) SyncMinimize(viewerIndex int) {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowFocused) {
		return
	}
	ss.viewerWindowFocused[viewerIndex].Store(false)
}

// SyncRestore はウィンドウ復元連動を反映する。
func (ss *SharedState) SyncRestore(viewerIndex int) {
	if viewerIndex < 0 || viewerIndex >= len(ss.viewerWindowFocused) {
		return
	}
	ss.viewerWindowFocused[viewerIndex].Store(true)
}

// ModelCount はモデル数を返す。
func (ss *SharedState) ModelCount(viewerIndex int) int {
	if viewerIndex < 0 || viewerIndex >= len(ss.models) {
		return 0
	}
	return len(ss.models[viewerIndex])
}

// MotionCount はモーション数を返す。
func (ss *SharedState) MotionCount(viewerIndex int) int {
	if viewerIndex < 0 || viewerIndex >= len(ss.motions) {
		return 0
	}
	if len(ss.motions[viewerIndex]) < len(ss.models[viewerIndex]) {
		return len(ss.models[viewerIndex])
	}
	return len(ss.motions[viewerIndex])
}

// SetModel はモデルを設定する。
func (ss *SharedState) SetModel(viewerIndex, modelIndex int, model IStateModel) {
	slot := ss.ensureModelSlot(viewerIndex, modelIndex)
	if slot == nil {
		return
	}
	slot.Store(stateModelSlot{Model: model})
	idxSlot := ss.ensureIndexSlot(viewerIndex, modelIndex)
	if idxSlot != nil {
		idxSlot.Store(stateIndexSlot{Indexes: []int{}})
	}
}

// Model はモデルを取得する。
func (ss *SharedState) Model(viewerIndex, modelIndex int) IStateModel {
	if viewerIndex < 0 || viewerIndex >= len(ss.models) {
		return nil
	}
	if modelIndex < 0 || modelIndex >= len(ss.models[viewerIndex]) {
		return nil
	}
	slot := ss.models[viewerIndex][modelIndex].Load().(stateModelSlot)
	return slot.Model
}

// SetMotion はモーションを設定する。
func (ss *SharedState) SetMotion(viewerIndex, modelIndex int, motion IStateMotion) {
	slot := ss.ensureMotionSlot(viewerIndex, modelIndex)
	if slot == nil {
		return
	}
	slot.Store(stateMotionSlot{Motion: motion})
}

// Motion はモーションを取得する。
func (ss *SharedState) Motion(viewerIndex, modelIndex int) IStateMotion {
	if viewerIndex < 0 || viewerIndex >= len(ss.motions) {
		return nil
	}
	if modelIndex < 0 || modelIndex >= len(ss.motions[viewerIndex]) {
		return nil
	}
	slot := ss.motions[viewerIndex][modelIndex].Load().(stateMotionSlot)
	return slot.Motion
}

// SelectedMaterialIndexes は選択材質インデックスを返す。
func (ss *SharedState) SelectedMaterialIndexes(viewerIndex, modelIndex int) []int {
	if viewerIndex < 0 || viewerIndex >= len(ss.selectedIndexes) {
		return nil
	}
	if modelIndex < 0 || modelIndex >= len(ss.selectedIndexes[viewerIndex]) {
		return nil
	}
	slot := ss.selectedIndexes[viewerIndex][modelIndex].Load().(stateIndexSlot)
	return cloneIntSlice(slot.Indexes)
}

// SetSelectedMaterialIndexes は選択材質インデックスを設定する。
func (ss *SharedState) SetSelectedMaterialIndexes(viewerIndex, modelIndex int, indexes []int) {
	slot := ss.ensureIndexSlot(viewerIndex, modelIndex)
	if slot == nil {
		return
	}
	slot.Store(stateIndexSlot{Indexes: cloneIntSlice(indexes)})
}

// IsDeltaSaveEnabled は差分保存が有効か判定する。
func (ss *SharedState) IsDeltaSaveEnabled(viewerIndex int) bool {
	if viewerIndex < 0 || viewerIndex >= len(ss.deltaSaveEnabled) {
		return false
	}
	return ss.deltaSaveEnabled[viewerIndex].Load()
}

// SetDeltaSaveEnabled は差分保存の有効可否を設定する。
func (ss *SharedState) SetDeltaSaveEnabled(viewerIndex int, enabled bool) {
	if viewerIndex < 0 || viewerIndex >= len(ss.deltaSaveEnabled) {
		return
	}
	ss.deltaSaveEnabled[viewerIndex].Store(enabled)
}

// DeltaSaveIndex は差分保存インデックスを返す。
func (ss *SharedState) DeltaSaveIndex(viewerIndex int) int {
	if viewerIndex < 0 || viewerIndex >= len(ss.deltaSaveIndexes) {
		return 0
	}
	return int(ss.deltaSaveIndexes[viewerIndex].Load())
}

// SetDeltaSaveIndex は差分保存インデックスを設定する。
func (ss *SharedState) SetDeltaSaveIndex(viewerIndex int, index int) {
	if viewerIndex < 0 || viewerIndex >= len(ss.deltaSaveIndexes) {
		return
	}
	ss.deltaSaveIndexes[viewerIndex].Store(int32(index))
}

// SetDeltaMotion は差分モーションを設定する。
func (ss *SharedState) SetDeltaMotion(viewerIndex, modelIndex, motionIndex int, motion IStateMotion) {
	slot := ss.ensureDeltaMotionSlot(viewerIndex, modelIndex, motionIndex)
	if slot == nil {
		return
	}
	slot.Store(stateMotionSlot{Motion: motion})
}

// DeltaMotion は差分モーションを取得する。
func (ss *SharedState) DeltaMotion(viewerIndex, modelIndex, deltaIndex int) IStateMotion {
	if viewerIndex < 0 || viewerIndex >= len(ss.deltaMotions) {
		return nil
	}
	if modelIndex < 0 || modelIndex >= len(ss.deltaMotions[viewerIndex]) {
		return nil
	}
	if deltaIndex < 0 || deltaIndex >= len(ss.deltaMotions[viewerIndex][modelIndex]) {
		return nil
	}
	slot := ss.deltaMotions[viewerIndex][modelIndex][deltaIndex].Load().(stateMotionSlot)
	return slot.Motion
}

// ClearDeltaMotion は差分モーションをクリアする。
func (ss *SharedState) ClearDeltaMotion(viewerIndex, modelIndex int) {
	if viewerIndex < 0 || viewerIndex >= len(ss.deltaMotions) {
		return
	}
	if modelIndex < 0 || modelIndex >= len(ss.deltaMotions[viewerIndex]) {
		return
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.deltaMotions[viewerIndex][modelIndex] = nil
}

// DeltaMotionCount は差分モーション数を返す。
func (ss *SharedState) DeltaMotionCount(viewerIndex, modelIndex int) int {
	if viewerIndex < 0 || viewerIndex >= len(ss.deltaMotions) {
		return 0
	}
	if modelIndex < 0 || modelIndex >= len(ss.deltaMotions[viewerIndex]) {
		return 0
	}
	return len(ss.deltaMotions[viewerIndex][modelIndex])
}

// SetPhysicsWorldMotion は物理ワールドモーションを設定する。
func (ss *SharedState) SetPhysicsWorldMotion(viewerIndex int, motion IStateMotion) {
	if viewerIndex < 0 || viewerIndex >= len(ss.physicsWorldMotions) {
		return
	}
	ss.physicsWorldMotions[viewerIndex].Store(stateMotionSlot{Motion: motion})
}

// PhysicsWorldMotion は物理ワールドモーションを取得する。
func (ss *SharedState) PhysicsWorldMotion(viewerIndex int) IStateMotion {
	if viewerIndex < 0 || viewerIndex >= len(ss.physicsWorldMotions) {
		return nil
	}
	slot := ss.physicsWorldMotions[viewerIndex].Load().(stateMotionSlot)
	return slot.Motion
}

// SetPhysicsModelMotion は物理モデルモーションを設定する。
func (ss *SharedState) SetPhysicsModelMotion(viewerIndex, modelIndex int, motion IStateMotion) {
	slot := ss.ensurePhysicsModelSlot(viewerIndex, modelIndex)
	if slot == nil {
		return
	}
	slot.Store(stateMotionSlot{Motion: motion})
}

// PhysicsModelMotion は物理モデルモーションを取得する。
func (ss *SharedState) PhysicsModelMotion(viewerIndex, modelIndex int) IStateMotion {
	if viewerIndex < 0 || viewerIndex >= len(ss.physicsModelMotions) {
		return nil
	}
	if modelIndex < 0 || modelIndex >= len(ss.physicsModelMotions[viewerIndex]) {
		return nil
	}
	slot := ss.physicsModelMotions[viewerIndex][modelIndex].Load().(stateMotionSlot)
	return slot.Motion
}

// SetWindMotion は風モーションを設定する。
func (ss *SharedState) SetWindMotion(viewerIndex int, motion IStateMotion) {
	if viewerIndex < 0 || viewerIndex >= len(ss.windMotions) {
		return
	}
	ss.windMotions[viewerIndex].Store(stateMotionSlot{Motion: motion})
}

// WindMotion は風モーションを取得する。
func (ss *SharedState) WindMotion(viewerIndex int) IStateMotion {
	if viewerIndex < 0 || viewerIndex >= len(ss.windMotions) {
		return nil
	}
	slot := ss.windMotions[viewerIndex].Load().(stateMotionSlot)
	return slot.Motion
}

// PhysicsResetType は物理リセット種別を返す。
func (ss *SharedState) PhysicsResetType() PhysicsResetType {
	return PhysicsResetType(ss.physicsResetType.Load())
}

// SetPhysicsResetType は物理リセット種別を設定する。
func (ss *SharedState) SetPhysicsResetType(resetType PhysicsResetType) {
	ss.physicsResetType.Store(int32(resetType))
}

// ensureModelSlot はモデルスロットを確保する。
func (ss *SharedState) ensureModelSlot(viewerIndex, modelIndex int) *atomic.Value {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.ensureViewerIndex(viewerIndex) {
		return nil
	}
	ss.models[viewerIndex] = ensureSlotSlice(ss.models[viewerIndex], modelIndex, stateModelSlot{})
	ss.motions[viewerIndex] = ensureSlotSlice(ss.motions[viewerIndex], modelIndex, stateMotionSlot{})
	ss.selectedIndexes[viewerIndex] = ensureSlotSlice(ss.selectedIndexes[viewerIndex], modelIndex, stateIndexSlot{Indexes: []int{}})
	return &ss.models[viewerIndex][modelIndex]
}

// ensureMotionSlot はモーションスロットを確保する。
func (ss *SharedState) ensureMotionSlot(viewerIndex, modelIndex int) *atomic.Value {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.ensureViewerIndex(viewerIndex) {
		return nil
	}
	ss.motions[viewerIndex] = ensureSlotSlice(ss.motions[viewerIndex], modelIndex, stateMotionSlot{})
	return &ss.motions[viewerIndex][modelIndex]
}

// ensureIndexSlot はインデックススロットを確保する。
func (ss *SharedState) ensureIndexSlot(viewerIndex, modelIndex int) *atomic.Value {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.ensureViewerIndex(viewerIndex) {
		return nil
	}
	ss.selectedIndexes[viewerIndex] = ensureSlotSlice(ss.selectedIndexes[viewerIndex], modelIndex, stateIndexSlot{Indexes: []int{}})
	return &ss.selectedIndexes[viewerIndex][modelIndex]
}

// ensureDeltaMotionSlot は差分モーションスロットを確保する。
func (ss *SharedState) ensureDeltaMotionSlot(viewerIndex, modelIndex, motionIndex int) *atomic.Value {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.ensureViewerIndex(viewerIndex) {
		return nil
	}
	if len(ss.deltaMotions[viewerIndex]) <= modelIndex {
		ss.deltaMotions[viewerIndex] = append(ss.deltaMotions[viewerIndex], make([][]atomic.Value, modelIndex-len(ss.deltaMotions[viewerIndex])+1)...)
	}
	ss.deltaMotions[viewerIndex][modelIndex] = ensureSlotSlice(ss.deltaMotions[viewerIndex][modelIndex], motionIndex, stateMotionSlot{})
	return &ss.deltaMotions[viewerIndex][modelIndex][motionIndex]
}

// ensurePhysicsModelSlot は物理モデルモーションスロットを確保する。
func (ss *SharedState) ensurePhysicsModelSlot(viewerIndex, modelIndex int) *atomic.Value {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.ensureViewerIndex(viewerIndex) {
		return nil
	}
	ss.physicsModelMotions[viewerIndex] = ensureSlotSlice(ss.physicsModelMotions[viewerIndex], modelIndex, stateMotionSlot{})
	return &ss.physicsModelMotions[viewerIndex][modelIndex]
}

// ensureViewerIndex はviewerIndexの配列初期化を行い、成否を返す。
func (ss *SharedState) ensureViewerIndex(viewerIndex int) bool {
	if viewerIndex < 0 || viewerIndex >= len(ss.models) {
		return false
	}
	if ss.models[viewerIndex] == nil {
		ss.models[viewerIndex] = []atomic.Value{}
	}
	if ss.motions[viewerIndex] == nil {
		ss.motions[viewerIndex] = []atomic.Value{}
	}
	if ss.selectedIndexes[viewerIndex] == nil {
		ss.selectedIndexes[viewerIndex] = []atomic.Value{}
	}
	if ss.deltaMotions[viewerIndex] == nil {
		ss.deltaMotions[viewerIndex] = [][]atomic.Value{}
	}
	if ss.physicsModelMotions[viewerIndex] == nil {
		ss.physicsModelMotions[viewerIndex] = []atomic.Value{}
	}
	return true
}

// ensureSlotSlice は必要長までスロットを初期化する。
func ensureSlotSlice[T any](slots []atomic.Value, index int, initial T) []atomic.Value {
	if index < 0 {
		return slots
	}
	if len(slots) <= index {
		for i := len(slots); i <= index; i++ {
			var v atomic.Value
			v.Store(initial)
			slots = append(slots, v)
		}
	}
	return slots
}

// cloneIntSlice はスライスを複製する。
func cloneIntSlice(src []int) []int {
	if src == nil {
		return nil
	}
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

type defaultMotion struct {
	*hashable.HashableBase
	Gravity       float32
	MaxSubSteps   int
	FixedTimeStep int
	WindEnabled   bool
	WindDirection [3]float32
}

// GetHashParts はハッシュ部品を返す。
func (m *defaultMotion) GetHashParts() string {
	return ""
}

// newDefaultPhysicsWorldMotion は物理ワールドの既定モーションを生成する。
func newDefaultPhysicsWorldMotion() IStateMotion {
	base := hashable.NewHashableBase("", "")
	m := &defaultMotion{
		HashableBase:  base,
		Gravity:       -9.8,
		MaxSubSteps:   physics_api.PhysicsDefaultMaxSubSteps,
		FixedTimeStep: 60,
		WindEnabled:   false,
		WindDirection: [3]float32{0, 0, 0},
	}
	m.SetHashPartsFunc(m.GetHashParts)
	return m
}

// newDefaultWindMotion は風の既定モーションを生成する。
func newDefaultWindMotion() IStateMotion {
	base := hashable.NewHashableBase("", "")
	m := &defaultMotion{
		HashableBase:  base,
		WindEnabled:   false,
		WindDirection: [3]float32{0, 0, 0},
	}
	m.SetHashPartsFunc(m.GetHashParts)
	return m
}
