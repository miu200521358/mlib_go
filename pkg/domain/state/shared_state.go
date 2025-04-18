//go:build windows
// +build windows

package state

import (
	"sync/atomic"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/win"
)

type SharedState struct {
	flags                      uint32           // 32ビット分のフラグを格納
	frameValue                 atomic.Value     // 現在フレーム
	maxFrameValue              atomic.Value     // 最大フレーム
	frameIntervalValue         atomic.Value     // FPS制限
	linkingFocus               atomic.Bool      // 連動フォーカス中かどうかのフラグ
	controlWindowPosition      atomic.Value     // コントロールウィンドウの位置
	controlWindowHandle        atomic.Int32     // コントロールウィンドウのハンドル
	viewerWindowHandles        []atomic.Int32   // ビューウィンドウのハンドル
	focusedWindowHandle        atomic.Int32     // フォーカス中のウィンドウのハンドル
	isInitializedControlWindow atomic.Bool      // コントロールウィンドウの初期化状態
	isInitializedViewWindow    []atomic.Bool    // ビューウィンドウの初期化状態
	focusControlWindow         atomic.Bool      // コントロールウィンドウのフォーカス状態
	focusViewWindow            []atomic.Bool    // ビューウィンドウのフォーカス状態
	isTriggeredFpsLimit        atomic.Bool      // FPS制限トリガー
	movedControlWindow         atomic.Bool      // コントロールウィンドウの移動状態
	isClosed                   atomic.Bool      // ウィンドウのクローズ状態
	models                     [][]atomic.Value // モデルデータ(ウィンドウ/モデルインデックス)
	motions                    [][]atomic.Value // モーションデータ(ウィンドウ/モデルインデックス)
	selectedMaterialIndexes    [][]atomic.Value // 選択中のマテリアルインデックス(ウィンドウ/モデルインデックス)
}

// NewSharedState は2つのStateを注入して生成するコンストラクタ
func NewSharedState(viewerCount int) *SharedState {
	shared := &SharedState{
		flags:                   0,
		viewerWindowHandles:     make([]atomic.Int32, viewerCount),
		isInitializedViewWindow: make([]atomic.Bool, viewerCount),
		focusViewWindow:         make([]atomic.Bool, viewerCount),
		models:                  make([][]atomic.Value, viewerCount),
		motions:                 make([][]atomic.Value, viewerCount),
		selectedMaterialIndexes: make([][]atomic.Value, viewerCount),
	}

	shared.SetFrame(0)
	shared.SetMaxFrame(0)
	shared.SetFrameInterval(-1)
	shared.SetControlWindowPosition(0, 0, 0, 0)
	shared.SetControlWindowHandle(0)
	shared.SetFocusedWindowHandle(0)
	shared.SetInitializedControlWindow(false)
	shared.SetFocusedWindowHandle(0)
	shared.SetFocusControlWindow(false)
	shared.SetMovedControlWindow(false)
	shared.SetClosed(false)

	return shared
}

const (
	FlagEnabledFrameDrop         = 1 << iota // フレームドロップON/OFF
	FlagEnabledPhysics                       // 物理ON/OFF
	FlagPhysicsReset                         // 物理リセット
	FlagShowNormal                           // ボーンデバッグ表示
	FlagShowWire                             // ワイヤーフレームデバッグ表示
	FlagShowOverrideUpper                    // オーバーライドデバッグ(上半身)表示
	FlagShowOverrideLower                    // オーバーライドデバッグ(下半身)表示
	FlagShowOverrideNone                     // オーバーライドデバッグ(カメラ合わせなし)表示
	FlagShowSelectedVertex                   // 選択頂点デバッグ表示
	FlagShowBoneAll                          // 全ボーンデバッグ表示
	FlagShowBoneIk                           // IKボーンデバッグ表示
	FlagShowBoneEffector                     // 付与親ボーンデバッグ表示
	FlagShowBoneFixed                        // 軸制限ボーンデバッグ表示
	FlagShowBoneRotate                       // 回転ボーンデバッグ表示
	FlagShowBoneTranslate                    // 移動ボーンデバッグ表示
	FlagShowBoneVisible                      // 表示ボーンデバッグ表示
	FlagShowRigidBodyFront                   // 剛体デバッグ表示(前面)
	FlagShowRigidBodyBack                    // 剛体デバッグ表示(埋め込み)
	FlagShowJoint                            // ジョイントデバッグ表示
	FlagShowInfo                             // 情報デバッグ表示
	FlagCameraSync                           // カメラ同期
	FlagPlaying                              // 再生中フラグ
	FlagWindowLinkage                        // ウィンドウリンクフラグ
	FlagIsChangedEnableDropFrame             // フレームドロップON/OFF変更フラグ
)

func (ss *SharedState) ModelCount(windowIndex int) int {
	if len(ss.models) <= windowIndex {
		return 0
	}
	return len(ss.models[windowIndex])
}

func (ss *SharedState) MotionCount(windowIndex int) int {
	if len(ss.motions) <= windowIndex {
		return 0
	}
	// モデルが読み込まれていたらモーションは必須
	return max(len(ss.motions[windowIndex]), len(ss.models[windowIndex]))
}

// StoreModel は指定されたウィンドウとモデルインデックスにモデルを格納
func (ss *SharedState) StoreModel(windowIndex, modelIndex int, model *pmx.PmxModel) {
	if len(ss.models) <= windowIndex {
		return
	}
	for modelIndex >= len(ss.models[windowIndex]) {
		ss.models[windowIndex] = append(ss.models[windowIndex], atomic.Value{})
		ss.selectedMaterialIndexes[windowIndex] = append(ss.selectedMaterialIndexes[windowIndex], atomic.Value{})
	}
	ss.models[windowIndex][modelIndex].Store(model)
	if model != nil {
		ss.selectedMaterialIndexes[windowIndex][modelIndex].Store(model.Materials.Indexes())
	} else {
		ss.selectedMaterialIndexes[windowIndex][modelIndex].Store([]int{})
	}
}

// LoadModel は指定されたウィンドウとモデルインデックスのモデルを取得
func (ss *SharedState) LoadModel(windowIndex, modelIndex int) *pmx.PmxModel {
	if len(ss.models) <= windowIndex {
		return nil
	}
	if len(ss.models[windowIndex]) <= modelIndex {
		return nil
	}
	return ss.models[windowIndex][modelIndex].Load().(*pmx.PmxModel)
}

// StoreMotion は指定されたウィンドウとモデルインデックスにモーションを格納
func (ss *SharedState) StoreMotion(windowIndex, modelIndex int, motion *vmd.VmdMotion) {
	if len(ss.motions) <= windowIndex {
		return
	}
	for modelIndex >= len(ss.motions[windowIndex]) {
		ss.motions[windowIndex] = append(ss.motions[windowIndex], atomic.Value{})
	}
	if motion != nil {
		ss.motions[windowIndex][modelIndex].Store(motion)
	} else {
		ss.motions[windowIndex][modelIndex].Store(vmd.NewVmdMotion(""))
	}
}

// LoadMotion は指定されたウィンドウとモデルインデックスのモーションを取得
func (ss *SharedState) LoadMotion(windowIndex, modelIndex int) *vmd.VmdMotion {
	if len(ss.motions) <= windowIndex {
		return nil
	}
	if len(ss.motions[windowIndex]) <= modelIndex {
		return vmd.NewVmdMotion("")
	}
	return ss.motions[windowIndex][modelIndex].Load().(*vmd.VmdMotion)
}

// LoadSelectedMaterialIndexes は選択中のマテリアルインデックスを取得
func (ss *SharedState) LoadSelectedMaterialIndexes(windowIndex, modelIndex int) []int {
	if len(ss.selectedMaterialIndexes) <= windowIndex {
		return nil
	}
	if len(ss.selectedMaterialIndexes[windowIndex]) <= modelIndex {
		return nil
	}
	return ss.selectedMaterialIndexes[windowIndex][modelIndex].Load().([]int)
}

// StoreSelectedMaterialIndexes は選択中のマテリアルインデックスを格納
func (ss *SharedState) StoreSelectedMaterialIndexes(windowIndex, modelIndex int, indexes []int) {
	if len(ss.selectedMaterialIndexes) <= windowIndex {
		return
	}
	if len(ss.selectedMaterialIndexes[windowIndex]) <= modelIndex {
		for i := len(ss.selectedMaterialIndexes[windowIndex]); i <= modelIndex; i++ {
			ss.selectedMaterialIndexes[windowIndex] = append(ss.selectedMaterialIndexes[windowIndex], atomic.Value{})
		}
	}
	ss.selectedMaterialIndexes[windowIndex][modelIndex].Store(indexes)
}

// アトミックに取得
func (ss *SharedState) loadFlag() uint32 {
	return atomic.LoadUint32(&ss.flags)
}

// 新しいフラグ値を計算する（ビット操作のみ）
func (ss *SharedState) setBit(currentFlag uint32, bitMask uint32, enable bool) uint32 {
	if enable {
		return currentFlag | bitMask
	}
	return currentFlag &^ bitMask
}

// 複数フラグ一括更新用の関数
func (ss *SharedState) UpdateFlags(changes map[uint32]bool) {
	for {
		oldVal := ss.loadFlag()
		newVal := oldVal

		// すべての変更を適用
		for bitMask, enable := range changes {
			newVal = ss.setBit(newVal, bitMask, enable)
		}

		if atomic.CompareAndSwapUint32(&ss.flags, oldVal, newVal) {
			return
		}
	}
}

// ビットが立っているかどうか
func (ss *SharedState) isBitSet(bitMask uint32) bool {
	return (ss.loadFlag() & bitMask) != 0
}

func (ss *SharedState) IsChangedEnableDropFrame() bool {
	return ss.isBitSet(FlagIsChangedEnableDropFrame)
}

func (ss *SharedState) SetChangedEnableDropFrame(changed bool) {
	ss.UpdateFlags(map[uint32]bool{FlagIsChangedEnableDropFrame: changed})
}

func (ss *SharedState) IsEnabledFrameDrop() bool {
	return ss.isBitSet(FlagEnabledFrameDrop)
}

func (ss *SharedState) SetEnabledFrameDrop(enabled bool) {
	ss.UpdateFlags(map[uint32]bool{FlagEnabledFrameDrop: enabled})
}

func (ss *SharedState) IsEnabledPhysics() bool {
	return ss.isBitSet(FlagEnabledPhysics)
}

func (ss *SharedState) SetEnabledPhysics(enabled bool) {
	ss.UpdateFlags(map[uint32]bool{FlagEnabledPhysics: enabled})
}

func (ss *SharedState) IsPhysicsReset() bool {
	return ss.isBitSet(FlagPhysicsReset)
}

func (ss *SharedState) SetPhysicsReset(reset bool) {
	ss.UpdateFlags(map[uint32]bool{FlagPhysicsReset: reset})
}

func (ss *SharedState) IsShowNormal() bool {
	return ss.isBitSet(FlagShowNormal)
}

func (ss *SharedState) SetShowNormal(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowNormal: show})
}

func (ss *SharedState) IsShowWire() bool {
	return ss.isBitSet(FlagShowWire)
}

func (ss *SharedState) SetShowWire(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowWire: show})
}

func (ss *SharedState) IsShowOverride() bool {
	return ss.IsShowOverrideUpper() || ss.IsShowOverrideLower() || ss.IsShowOverrideNone()
}

func (ss *SharedState) IsShowOverrideUpper() bool {
	return ss.isBitSet(FlagShowOverrideUpper)
}

func (ss *SharedState) SetShowOverrideUpper(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowOverrideUpper: show})
}

func (ss *SharedState) IsShowOverrideLower() bool {
	return ss.isBitSet(FlagShowOverrideLower)
}

func (ss *SharedState) SetShowOverrideLower(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowOverrideLower: show})
}

func (ss *SharedState) IsShowOverrideNone() bool {
	return ss.isBitSet(FlagShowOverrideNone)
}

func (ss *SharedState) SetShowOverrideNone(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowOverrideNone: show})
}

func (ss *SharedState) IsShowSelectedVertex() bool {
	return ss.isBitSet(FlagShowSelectedVertex)
}

func (ss *SharedState) SetShowSelectedVertex(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowSelectedVertex: show})
}

func (ss *SharedState) IsShowBoneAll() bool {
	return ss.isBitSet(FlagShowBoneAll)
}

func (ss *SharedState) SetShowBoneAll(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowBoneAll: show})
}

func (ss *SharedState) IsShowBoneIk() bool {
	return ss.isBitSet(FlagShowBoneIk)
}

func (ss *SharedState) SetShowBoneIk(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowBoneIk: show})
}

func (ss *SharedState) IsShowBoneEffector() bool {
	return ss.isBitSet(FlagShowBoneEffector)
}

func (ss *SharedState) SetShowBoneEffector(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowBoneEffector: show})
}

func (ss *SharedState) IsShowBoneFixed() bool {
	return ss.isBitSet(FlagShowBoneFixed)
}

func (ss *SharedState) SetShowBoneFixed(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowBoneFixed: show})
}

func (ss *SharedState) IsShowBoneRotate() bool {
	return ss.isBitSet(FlagShowBoneRotate)
}

func (ss *SharedState) SetShowBoneRotate(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowBoneRotate: show})
}

func (ss *SharedState) IsShowBoneTranslate() bool {
	return ss.isBitSet(FlagShowBoneTranslate)
}

func (ss *SharedState) SetShowBoneTranslate(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowBoneTranslate: show})
}

func (ss *SharedState) IsShowBoneVisible() bool {
	return ss.isBitSet(FlagShowBoneVisible)
}

func (ss *SharedState) SetShowBoneVisible(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowBoneVisible: show})
}

func (ss *SharedState) IsShowRigidBodyFront() bool {
	return ss.isBitSet(FlagShowRigidBodyFront)
}

func (ss *SharedState) SetShowRigidBodyFront(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowRigidBodyFront: show})
}

func (ss *SharedState) IsShowRigidBodyBack() bool {
	return ss.isBitSet(FlagShowRigidBodyBack)
}

func (ss *SharedState) SetShowRigidBodyBack(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowRigidBodyBack: show})
}

func (ss *SharedState) IsShowJoint() bool {
	return ss.isBitSet(FlagShowJoint)
}

func (ss *SharedState) SetShowJoint(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowJoint: show})
}

func (ss *SharedState) IsShowInfo() bool {
	return ss.isBitSet(FlagShowInfo)
}

func (ss *SharedState) SetShowInfo(show bool) {
	ss.UpdateFlags(map[uint32]bool{FlagShowInfo: show})
}

func (ss *SharedState) IsCameraSync() bool {
	return ss.isBitSet(FlagCameraSync)
}

func (ss *SharedState) SetCameraSync(sync bool) {
	ss.UpdateFlags(map[uint32]bool{FlagCameraSync: sync})
}

func (ss *SharedState) Playing() bool {
	return ss.isBitSet(FlagPlaying)
}

func (ss *SharedState) SetPlaying(p bool) {
	ss.UpdateFlags(map[uint32]bool{FlagPlaying: p})
}

func (ss *SharedState) IsWindowLinkage() bool {
	return ss.isBitSet(FlagWindowLinkage)
}

func (ss *SharedState) SetWindowLinkage(link bool) {
	ss.UpdateFlags(map[uint32]bool{FlagWindowLinkage: link})
}

func (ss *SharedState) Frame() float32 {
	return ss.frameValue.Load().(float32)
}

func (ss *SharedState) SetFrame(frame float32) {
	ss.frameValue.Store(frame)
}

func (ss *SharedState) MaxFrame() float32 {
	return ss.maxFrameValue.Load().(float32)
}

func (ss *SharedState) SetMaxFrame(maxFrame float32) {
	ss.maxFrameValue.Store(maxFrame)
}

func (ss *SharedState) FrameInterval() float32 {
	return ss.frameIntervalValue.Load().(float32)
}

func (ss *SharedState) SetFrameInterval(spf float32) {
	ss.frameIntervalValue.Store(spf)
}

func (ss *SharedState) ControlWindowPosition() (x, y, diffX, diffY int) {
	diff := ss.controlWindowPosition.Load().(mmath.MVec4)
	return int(diff.X), int(diff.Y), int(diff.Z), int(diff.W)
}

func (ss *SharedState) SetControlWindowPosition(x, y, diffX, diffY int) {
	ss.controlWindowPosition.Store(mmath.MVec4{X: float64(x), Y: float64(y), Z: float64(diffX), W: float64(diffY)})
}

func (ss *SharedState) ControlWindowHandle() int32 {
	return ss.controlWindowHandle.Load()
}

func (ss *SharedState) SetControlWindowHandle(handle int32) {
	ss.controlWindowHandle.Store(handle)
}

func (ss *SharedState) ViewerWindowHandle(windowIndex int) int32 {
	return ss.viewerWindowHandles[windowIndex].Load()
}

func (ss *SharedState) SetViewerWindowHandle(windowIndex int, handle int32) {
	ss.viewerWindowHandles[windowIndex].Store(handle)
}

func (ss *SharedState) IsHitWindowHandle(handle int32) bool {
	if ss.ControlWindowHandle() == handle {
		return true
	}

	for i := range ss.viewerWindowHandles {
		if ss.viewerWindowHandles[i].Load() == handle {
			return true
		}
	}
	return false
}

func (ss *SharedState) FocusedWindowHandle() int32 {
	return ss.focusedWindowHandle.Load()
}

func (ss *SharedState) SetFocusedWindowHandle(handle int32) {
	ss.focusedWindowHandle.Store(handle)
}

func (ss *SharedState) IsInitializedControlWindow() bool {
	return ss.isInitializedControlWindow.Load()
}

func (ss *SharedState) SetInitializedControlWindow(initialized bool) {
	ss.isInitializedControlWindow.Store(initialized)
}

func (ss *SharedState) IsInitializedViewWindow(windowIndex int) bool {
	return ss.isInitializedViewWindow[windowIndex].Load()
}

func (ss *SharedState) SetInitializedViewWindow(windowIndex int, initialized bool) {
	ss.isInitializedViewWindow[windowIndex].Store(initialized)
}

func (ss *SharedState) IsInitializedAllViewWindows() bool {
	for i := range ss.isInitializedViewWindow {
		if !ss.isInitializedViewWindow[i].Load() {
			return false
		}
	}
	return true
}

func (ss *SharedState) IsInitializedAllWindows() bool {
	return ss.IsInitializedControlWindow() && ss.IsInitializedAllViewWindows()
}

func (ss *SharedState) IsFocusControlWindow() bool {
	return ss.focusControlWindow.Load()
}

func (ss *SharedState) SetFocusControlWindow(focus bool) {
	ss.focusControlWindow.Store(focus)
}

func (ss *SharedState) IsTriggeredFpsLimit() bool {
	return ss.isTriggeredFpsLimit.Load()
}

func (ss *SharedState) SetTriggeredFpsLimit(triggered bool) {
	ss.isTriggeredFpsLimit.Store(triggered)
}

func (ss *SharedState) IsFocusViewWindow(windowIndex int) bool {
	return ss.focusViewWindow[windowIndex].Load()
}

func (ss *SharedState) SetFocusViewWindow(windowIndex int, focus bool) {
	ss.focusViewWindow[windowIndex].Store(focus)
}

func (ss *SharedState) SetFocusAllViewWindows(focus bool) {
	for i := range ss.focusViewWindow {
		ss.focusViewWindow[i].Store(focus)
	}
}

func (ss *SharedState) IsMovedControlWindow() bool {
	return ss.movedControlWindow.Load()
}

func (ss *SharedState) SetMovedControlWindow(moving bool) {
	ss.movedControlWindow.Store(moving)
}

func (ss *SharedState) IsClosed() bool {
	return ss.isClosed.Load()
}

func (ss *SharedState) SetClosed(closed bool) {
	ss.isClosed.Store(closed)
}

func (ss *SharedState) IsLinkingFocus() bool {
	return ss.linkingFocus.Load()
}

func (ss *SharedState) SetLinkingFocus(val bool) {
	ss.linkingFocus.Store(val)
}

// 任意ビューワーでフォーカスが発生した際に呼び出す共通関数
// viewerIndex: フォーカスが発生したビューワーのインデックス(-1: コントロールウィンドウ)
func (ss *SharedState) TriggerLinkedFocus(viewerIndex int) {
	if ss.IsLinkingFocus() {
		return
	}

	// // すでに連動処理中なら再発火を防止
	// if ss.linkingFocus.CompareAndSwap(false, true) {
	// 	// コントロールウィンドウの前面化要求はそのまま行い、
	// 	// 連動対象は発生元のViewer以外に限定する
	// 	if viewerIndex == -1 {
	// 		// コントローラーウィンドウをフォーカスして発火した場合、コントローラーウィンドウのハンドルを保持
	// 		ss.SetFocusedWindowHandle(ss.ControlWindowHandle())
	// 	} else {
	// 		// ビューアウィンドウをフォーカスして発火した場合、ビューアウィンドウのハンドルを保持
	// 		ss.SetFocusedWindowHandle(ss.ViewerWindowHandle(viewerIndex))
	// 	}

	// 	if viewerIndex >= 0 && win.IsWindowCenterObscured(win.HWND(ss.ControlWindowHandle())) {
	// 		// コントロールウィンドウが前面にない場合は前面化
	// 		ss.SetFocusControlWindow(true)
	// 	}
	// 	for i := range ss.focusViewWindow {
	// 		if win.IsWindowCenterObscured(win.HWND(ss.ViewerWindowHandle(i))) {
	// 			// Viewerが前面にない場合は前面化
	// 			ss.SetFocusViewWindow(i, i != viewerIndex)
	// 		}
	// 	}

	// 	// 連動中フラグを一定時間後に解除
	// 	go func() {
	// 		time.Sleep(300 * time.Millisecond)
	// 		ss.SetLinkingFocus(false)
	// 		// フォーカスを解除
	// 		ss.SetFocusedWindowHandle(0)
	// 	}()
	// }
}

func (ss *SharedState) KeepFocus() {
	// 連動処理中の場合のみ処理を行う
	if ss.IsLinkingFocus() {
		// 対象ウィンドウの前面化要求
		win.SetForegroundWindow(win.HWND(ss.FocusedWindowHandle()))
	}
}

// SyncMinimize は全ウィンドウを最小化する
// viewerIndex: フォーカスが発生したビューワーのインデックス(-1: コントロールウィンドウ)
func (ss *SharedState) SyncMinimize(viewerIndex int) {
	// 対象ウィンドウの最小化
	if viewerIndex >= 0 {
		win.ShowWindow(win.HWND(ss.ControlWindowHandle()), win.SW_MINIMIZE)
	}
	for i := range ss.viewerWindowHandles {
		if i != viewerIndex {
			win.ShowWindow(win.HWND(ss.ViewerWindowHandle(i)), win.SW_MINIMIZE)
		}
	}
}

// SyncRestore は全ウィンドウの最小化を解除する
// viewerIndex: フォーカスが発生したビューワーのインデックス(-1: コントロールウィンドウ)
func (ss *SharedState) SyncRestore(viewerIndex int) {
	// 対象ウィンドウの最小化を解除
	if viewerIndex >= 0 {
		win.ShowWindow(win.HWND(ss.ControlWindowHandle()), win.SW_RESTORE)
	}
	for i := range ss.viewerWindowHandles {
		if i != viewerIndex {
			win.ShowWindow(win.HWND(ss.ViewerWindowHandle(i)), win.SW_RESTORE)
		}
	}
}
