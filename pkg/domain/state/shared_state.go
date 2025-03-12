//go:build windows
// +build windows

package state

import (
	"sync/atomic"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type SharedState struct {
	flags                      uint32           // 32ビット分のフラグを格納
	frameValue                 atomic.Value     // 現在フレーム
	maxFrameValue              atomic.Value     // 最大フレーム
	frameIntervalValue         atomic.Value     // FPS制限
	controlWindowPosition      atomic.Value     // コントロールウィンドウの位置
	isActiveControlWindow      atomic.Bool      // コントロールウィンドウのアクティブ状態
	isActiveViewWindow         []atomic.Bool    // ビューウィンドウのアクティブ状態
	isInitializedControlWindow atomic.Bool      // コントロールウィンドウの初期化状態
	isInitializedViewWindow    []atomic.Bool    // ビューウィンドウの初期化状態
	focusControlWindow         atomic.Bool      // コントロールウィンドウのフォーカス状態
	focusViewWindow            atomic.Bool      // ビューウィンドウのフォーカス状態
	movedControlWindow         atomic.Bool      // コントロールウィンドウの移動状態
	models                     [][]atomic.Value // モデルデータ(ウィンドウ/モデルインデックス)
	motions                    [][]atomic.Value // モーションデータ(ウィンドウ/モデルインデックス)
}

// NewSharedState は2つのStateを注入して生成するコンストラクタ
func NewSharedState(viewerCount int) *SharedState {
	return &SharedState{
		flags:                   0,
		isActiveViewWindow:      make([]atomic.Bool, viewerCount),
		isInitializedViewWindow: make([]atomic.Bool, viewerCount),
		models:                  make([][]atomic.Value, viewerCount),
		motions:                 make([][]atomic.Value, viewerCount),
	}
}

type frameState struct {
	frame float32
}

type maxFrameState struct {
	maxFrame float32
}

type frameIntervalState struct {
	frameInterval float32
}

const (
	flagEnabledFrameDrop         = 1 << iota // フレームドロップON/OFF
	flagEnabledPhysics                       // 物理ON/OFF
	flagPhysicsReset                         // 物理リセット
	flagShowNormal                           // ボーンデバッグ表示
	flagShowWire                             // ワイヤーフレームデバッグ表示
	flagShowOverride                         // オーバーライドデバッグ表示
	flagShowSelectedVertex                   // 選択頂点デバッグ表示
	flagShowBoneAll                          // 全ボーンデバッグ表示
	flagShowBoneIk                           // IKボーンデバッグ表示
	flagShowBoneEffector                     // 付与親ボーンデバッグ表示
	flagShowBoneFixed                        // 軸制限ボーンデバッグ表示
	flagShowBoneRotate                       // 回転ボーンデバッグ表示
	flagShowBoneTranslate                    // 移動ボーンデバッグ表示
	flagShowBoneVisible                      // 表示ボーンデバッグ表示
	flagShowRigidBodyFront                   // 剛体デバッグ表示(前面)
	flagShowRigidBodyBack                    // 剛体デバッグ表示(埋め込み)
	flagShowJoint                            // ジョイントデバッグ表示
	flagShowInfo                             // 情報デバッグ表示
	flagCameraSync                           // カメラ同期
	flagPlaying                              // 再生中フラグ
	flagClosed                               // 描画ウィンドウクローズ
	flagWindowLinkage                        // ウィンドウリンクフラグ
	flagIsChangedEnableDropFrame             // フレームドロップON/OFF変更フラグ
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
	if len(ss.models[windowIndex]) <= modelIndex {
		for i := len(ss.models[windowIndex]); i <= modelIndex; i++ {
			ss.models[windowIndex] = append(ss.models[windowIndex], atomic.Value{})
		}
	}
	ss.models[windowIndex][modelIndex].Store(model)
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
	if len(ss.motions[windowIndex]) <= modelIndex {
		for i := len(ss.motions[windowIndex]); i <= modelIndex; i++ {
			ss.motions[windowIndex] = append(ss.motions[windowIndex], atomic.Value{})
		}
	}
	ss.motions[windowIndex][modelIndex].Store(motion)
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

// アトミックに取得
func (ss *SharedState) loadFlag() uint32 {
	return atomic.LoadUint32(&ss.flags)
}

// 特定ビットをON/OFFにする
func (ss *SharedState) setBit(bitMask uint32, enable bool) {
	for {
		oldVal := ss.loadFlag()
		newVal := oldVal
		if enable {
			newVal |= bitMask
		} else {
			newVal &= ^bitMask
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
	return ss.isBitSet(flagIsChangedEnableDropFrame)
}

func (ss *SharedState) SetChangedEnableDropFrame(changed bool) {
	ss.setBit(flagIsChangedEnableDropFrame, changed)
}

func (ss *SharedState) IsEnabledFrameDrop() bool {
	return ss.isBitSet(flagEnabledFrameDrop)
}

func (ss *SharedState) SetEnabledFrameDrop(enabled bool) {
	ss.setBit(flagEnabledFrameDrop, enabled)
}

func (ss *SharedState) IsEnabledPhysics() bool {
	return ss.isBitSet(flagEnabledPhysics)
}

func (ss *SharedState) SetEnabledPhysics(enabled bool) {
	ss.setBit(flagEnabledPhysics, enabled)
}

func (ss *SharedState) IsPhysicsReset() bool {
	return ss.isBitSet(flagPhysicsReset)
}

func (ss *SharedState) SetPhysicsReset(reset bool) {
	ss.setBit(flagPhysicsReset, reset)
}

func (ss *SharedState) IsShowNormal() bool {
	return ss.isBitSet(flagShowNormal)
}

func (ss *SharedState) SetShowNormal(show bool) {
	ss.setBit(flagShowNormal, show)
}

func (ss *SharedState) IsShowWire() bool {
	return ss.isBitSet(flagShowWire)
}

func (ss *SharedState) SetShowWire(show bool) {
	ss.setBit(flagShowWire, show)
}

func (ss *SharedState) IsShowOverride() bool {
	return ss.isBitSet(flagShowOverride)
}

func (ss *SharedState) SetShowOverride(show bool) {
	ss.setBit(flagShowOverride, show)
}

func (ss *SharedState) IsShowSelectedVertex() bool {
	return ss.isBitSet(flagShowSelectedVertex)
}

func (ss *SharedState) SetShowSelectedVertex(show bool) {
	ss.setBit(flagShowSelectedVertex, show)
}

func (ss *SharedState) IsShowBoneAll() bool {
	return ss.isBitSet(flagShowBoneAll)
}

func (ss *SharedState) SetShowBoneAll(show bool) {
	ss.setBit(flagShowBoneAll, show)
}

func (ss *SharedState) IsShowBoneIk() bool {
	return ss.isBitSet(flagShowBoneIk)
}

func (ss *SharedState) SetShowBoneIk(show bool) {
	ss.setBit(flagShowBoneIk, show)
}

func (ss *SharedState) IsShowBoneEffector() bool {
	return ss.isBitSet(flagShowBoneEffector)
}

func (ss *SharedState) SetShowBoneEffector(show bool) {
	ss.setBit(flagShowBoneEffector, show)
}

func (ss *SharedState) IsShowBoneFixed() bool {
	return ss.isBitSet(flagShowBoneFixed)
}

func (ss *SharedState) SetShowBoneFixed(show bool) {
	ss.setBit(flagShowBoneFixed, show)
}

func (ss *SharedState) IsShowBoneRotate() bool {
	return ss.isBitSet(flagShowBoneRotate)
}

func (ss *SharedState) SetShowBoneRotate(show bool) {
	ss.setBit(flagShowBoneRotate, show)
}

func (ss *SharedState) IsShowBoneTranslate() bool {
	return ss.isBitSet(flagShowBoneTranslate)
}

func (ss *SharedState) SetShowBoneTranslate(show bool) {
	ss.setBit(flagShowBoneTranslate, show)
}

func (ss *SharedState) IsShowBoneVisible() bool {
	return ss.isBitSet(flagShowBoneVisible)
}

func (ss *SharedState) SetShowBoneVisible(show bool) {
	ss.setBit(flagShowBoneVisible, show)
}

func (ss *SharedState) IsShowRigidBodyFront() bool {
	return ss.isBitSet(flagShowRigidBodyFront)
}

func (ss *SharedState) SetShowRigidBodyFront(show bool) {
	ss.setBit(flagShowRigidBodyFront, show)
}

func (ss *SharedState) IsShowRigidBodyBack() bool {
	return ss.isBitSet(flagShowRigidBodyBack)
}

func (ss *SharedState) SetShowRigidBodyBack(show bool) {
	ss.setBit(flagShowRigidBodyBack, show)
}

func (ss *SharedState) IsShowJoint() bool {
	return ss.isBitSet(flagShowJoint)
}

func (ss *SharedState) SetShowJoint(show bool) {
	ss.setBit(flagShowJoint, show)
}

func (ss *SharedState) IsShowInfo() bool {
	return ss.isBitSet(flagShowInfo)
}

func (ss *SharedState) SetShowInfo(show bool) {
	ss.setBit(flagShowInfo, show)
}

func (ss *SharedState) IsCameraSync() bool {
	return ss.isBitSet(flagCameraSync)
}

func (ss *SharedState) SetCameraSync(sync bool) {
	ss.setBit(flagCameraSync, sync)
}

func (ss *SharedState) Playing() bool {
	return ss.isBitSet(flagPlaying)
}

func (ss *SharedState) SetPlaying(p bool) {
	ss.setBit(flagPlaying, p)
}

func (ss *SharedState) IsClosed() bool {
	return ss.isBitSet(flagClosed)
}

func (ss *SharedState) SetClosed(closed bool) {
	ss.setBit(flagClosed, closed)
}

func (ss *SharedState) IsWindowLinkage() bool {
	return ss.isBitSet(flagWindowLinkage)
}

func (ss *SharedState) SetWindowLinkage(link bool) {
	ss.setBit(flagWindowLinkage, link)
}

func (ss *SharedState) Frame() float32 {
	return ss.frameValue.Load().(frameState).frame
}

func (ss *SharedState) SetFrame(frame float32) {
	ss.frameValue.Store(frameState{frame: frame})
}

func (ss *SharedState) MaxFrame() float32 {
	return ss.maxFrameValue.Load().(maxFrameState).maxFrame
}

func (ss *SharedState) UpdateMaxFrame(maxFrame float32) {
	if ss.MaxFrame() < maxFrame {
		ss.SetMaxFrame(maxFrame)
	}
}

func (ss *SharedState) SetMaxFrame(maxFrame float32) {
	ss.maxFrameValue.Store(maxFrameState{maxFrame: maxFrame})
}

func (ss *SharedState) FrameInterval() float32 {
	return ss.frameIntervalValue.Load().(frameIntervalState).frameInterval
}

func (ss *SharedState) SetFrameInterval(spf float32) {
	ss.frameIntervalValue.Store(frameIntervalState{frameInterval: spf})
}

func (ss *SharedState) ControlWindowPosition() (x, y, diffX, diffY int) {
	diff := ss.controlWindowPosition.Load().(mmath.MVec4)
	return int(diff.X), int(diff.Y), int(diff.Z), int(diff.W)
}

func (ss *SharedState) SetControlWindowPosition(x, y, diffX, diffY int) {
	ss.controlWindowPosition.Store(mmath.MVec4{X: float64(x), Y: float64(y), Z: float64(diffX), W: float64(diffY)})
}

func (ss *SharedState) IsActiveControlWindow() bool {
	return ss.isActiveControlWindow.Load()
}

func (ss *SharedState) SetActiveControlWindow(active bool) {
	ss.isActiveControlWindow.Store(active)
}

func (ss *SharedState) IsActivateViewWindow(windowIndex int) bool {
	return ss.isActiveViewWindow[windowIndex].Load()
}

func (ss *SharedState) SetActivateViewWindow(windowIndex int, active bool) {
	ss.isActiveViewWindow[windowIndex].Store(active)
}

func (ss *SharedState) IsInactiveALlViewWindows() bool {
	for i := range ss.isActiveViewWindow {
		if ss.isActiveViewWindow[i].Load() {
			return false
		}
	}
	return true
}

func (ss *SharedState) IsInactiveAllWindows() bool {
	return !ss.IsActiveControlWindow() && ss.IsInactiveALlViewWindows()
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

func (ss *SharedState) IsFocusViewWindow() bool {
	return ss.focusViewWindow.Load()
}

func (ss *SharedState) SetFocusViewWindow(focus bool) {
	ss.focusViewWindow.Store(focus)
}

func (ss *SharedState) IsMovedControlWindow() bool {
	return ss.movedControlWindow.Load()
}

func (ss *SharedState) SetMovedControlWindow(moving bool) {
	ss.movedControlWindow.Store(moving)
}
