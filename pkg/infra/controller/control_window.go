//go:build windows
// +build windows

// 指示: miu200521358
package controller

import (
	"fmt"
	"io"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	sharedi18n "github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	sharedtime "github.com/miu200521358/mlib_go/pkg/shared/contracts/time"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// FpsLimit はFPS制限種別を表す。
type FpsLimit int

const (
	// FPS_LIMIT_30 は30fps制限。
	FPS_LIMIT_30 FpsLimit = 30
	// FPS_LIMIT_60 は60fps制限。
	FPS_LIMIT_60 FpsLimit = 60
	// FPS_LIMIT_UNLIMITED は無制限。
	FPS_LIMIT_UNLIMITED FpsLimit = -1
)

const (
	// langJa は日本語の言語コード。
	langJa sharedi18n.LangCode = "ja"
	// langEn は英語の言語コード。
	langEn sharedi18n.LangCode = "en"
	// langZh は中国語の言語コード。
	langZh sharedi18n.LangCode = "zh"
	// langKo は韓国語の言語コード。
	langKo sharedi18n.LangCode = "ko"
)

// ControlWindow はコントローラーウィンドウを表す。
type ControlWindow struct {
	*walk.MainWindow

	shared     *state.SharedState
	appConfig  *config.AppConfig
	userConfig config.IUserConfig
	translator sharedi18n.II18n
	logger     logging.ILogger

	tabWidget   *walk.TabWidget
	consoleView *ConsoleView
	progressBar *ProgressBar

	viewerCount int

	setEnabledInPlaying func(playing bool)
	onChangePlayingPre  func(playing bool)
	onChangePlayingPost func(playing bool)
	onEnabledInPlaying  func(playing bool)

	leftButtonPressed bool

	enabledFrameDropAction       *walk.Action
	enabledPhysicsAction         *walk.Action
	physicsResetAction           *walk.Action
	showNormalAction             *walk.Action
	showWireAction               *walk.Action
	showOverrideUpperAction      *walk.Action
	showOverrideLowerAction      *walk.Action
	showOverrideNoneAction       *walk.Action
	showSelectedVertexAction     *walk.Action
	showBoneAllAction            *walk.Action
	showBoneIkAction             *walk.Action
	showBoneEffectorAction       *walk.Action
	showBoneFixedAction          *walk.Action
	showBoneRotateAction         *walk.Action
	showBoneTranslateAction      *walk.Action
	showBoneVisibleAction        *walk.Action
	showRigidBodyFrontAction     *walk.Action
	showRigidBodyBackAction      *walk.Action
	showJointAction              *walk.Action
	showInfoAction               *walk.Action
	limitFps30Action             *walk.Action
	limitFps60Action             *walk.Action
	limitFpsUnLimitAction        *walk.Action
	cameraSyncAction             *walk.Action
	logLevelDebugAction          *walk.Action
	logLevelVerboseAction        *walk.Action
	logLevelIkVerboseAction      *walk.Action
	logLevelPhysicsVerboseAction *walk.Action
	logLevelViewerVerboseAction  *walk.Action
	linkWindowAction             *walk.Action

	verboseSinks map[logging.VerboseIndex]logging.IVerboseSink
}

// NewControlWindow はコントローラーウィンドウを生成する。
func NewControlWindow(shared *state.SharedState, baseServices base.IBaseServices,
	helpMenuItems []declarative.MenuItem, tabPages []declarative.TabPage,
	width, height, positionX, positionY, viewerCount int,
) (*ControlWindow, error) {
	var appConfig *config.AppConfig
	var userConfig config.IUserConfig
	var translator sharedi18n.II18n
	var logger logging.ILogger
	if baseServices != nil {
		if cfg := baseServices.Config(); cfg != nil {
			appConfig = cfg.AppConfig()
			userConfig = cfg.UserConfig()
		}
		translator = baseServices.I18n()
		logger = baseServices.Logger()
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}

	cw := &ControlWindow{
		shared:       shared,
		appConfig:    appConfig,
		userConfig:   userConfig,
		translator:   translator,
		logger:       logger,
		verboseSinks: map[logging.VerboseIndex]logging.IVerboseSink{},
		viewerCount:  viewerCount,
	}

	controllerItems := cw.buildControllerMenuItems()
	if len(helpMenuItems) > 0 {
		controllerItems = append(controllerItems, declarative.Separator{})
		controllerItems = append(controllerItems, helpMenuItems...)
	}

	var consoleContainer *walk.Composite

	if err := (declarative.MainWindow{
		AssignTo: &cw.MainWindow,
		Title:    cw.appTitle(),
		Size:     declarative.Size{Width: width, Height: height},
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear, MarginsZero: true, SpacingZero: true},
		Background: declarative.SolidColorBrush{
			Color: ColorWindowBackground,
		},
		MenuItems: []declarative.MenuItem{
			cw.buildViewerMenu(),
			declarative.Menu{Text: cw.t("&コントローラーウィンドウ"), Items: controllerItems},
			declarative.Menu{Text: cw.t("&言語"), Items: cw.buildLanguageMenu()},
		},
		Children: []declarative.Widget{
			declarative.TabWidget{AssignTo: &cw.tabWidget, Pages: tabPages},
			declarative.Composite{
				AssignTo:      &consoleContainer,
				Layout:        declarative.VBox{MarginsZero: true, SpacingZero: true},
				StretchFactor: 1,
			},
		},
		OnClosing: func(canceled *bool, reason walk.CloseReason) {
			cw.onClosing(canceled)
		},
		OnClickActivate: func() {
			cw.onActivate()
		},
		OnEnterSizeMove: func() {
			cw.onEnterSizeMove()
		},
		OnExitSizeMove: func() {
			cw.onExitSizeMove()
		},
		OnMinimize: func() {
			cw.onMinimize()
		},
		OnRestore: func() {
			cw.onRestore()
		},
	}).Create(); err != nil {
		return nil, NewControllerWindowInitFailed("コントローラーウィンドウの初期化に失敗しました", err)
	}

	cw.SetPos(positionX, positionY)
	cw.shared.SetControlWindowPosition(state.WindowPosition{X: positionX, Y: positionY})
	cw.shared.SetControlWindowHandle(state.WindowHandle(uintptr(cw.Handle())))
	cw.shared.SetControlWindowReady(true)
	cw.shared.SetControlWindowFocused(true)

	if consoleContainer != nil {
		if cv, err := NewConsoleView(consoleContainer, width/10, height/10); err != nil {
			return nil, err
		} else {
			cw.consoleView = cv
			cw.setConsoleSink(cv)
		}
		if pb, err := NewProgressBar(consoleContainer); err != nil {
			return nil, err
		} else {
			cw.progressBar = pb
		}
	}

	cw.shared.SetFrame(0)
	cw.shared.SetMaxFrame(1)
	cw.SetPhysicsEnabled(true)

	cw.applyUserConfig()
	return cw, nil
}

// Run はメインウィンドウを起動する。
func (cw *ControlWindow) Run() int {
	return cw.MainWindow.Run()
}

// Dispose はウィンドウを破棄する。
func (cw *ControlWindow) Dispose() {
	cw.Close()
	cw.closeVerboseSinks()
}

// infoLineLogger はタイトル付き区切り線ログのI/F。
type infoLineLogger interface {
	InfoLineTitle(title, msg string, params ...any)
}

// t は翻訳済み文言を返す。
func (cw *ControlWindow) t(key string) string {
	if cw == nil || cw.translator == nil || !cw.translator.IsReady() {
		return "●●" + key + "●●"
	}
	return cw.translator.T(key)
}

// loggerOrDefault は利用可能なロガーを返す。
func (cw *ControlWindow) loggerOrDefault() logging.ILogger {
	if cw == nil || cw.logger == nil {
		return logging.DefaultLogger()
	}
	return cw.logger
}

// setConsoleSink はメッセージ欄の出力先を設定する。
func (cw *ControlWindow) setConsoleSink(writer io.Writer) {
	logger := cw.loggerOrDefault()
	if consoleLogger, ok := logger.(logging.IConsoleLogger); ok {
		consoleLogger.SetConsoleSink(writer)
		return
	}
	logging.SetConsoleSink(writer)
}

// infoLineTitle は区切り線付きタイトルログを出力する。
func (cw *ControlWindow) infoLineTitle(title, msg string) {
	logger := cw.loggerOrDefault()
	if titled, ok := logger.(infoLineLogger); ok {
		titled.InfoLineTitle(title, msg)
		return
	}
	logger.Info("%s %s", title, msg)
}

// OnClose は閉じる処理を行う。
func (cw *ControlWindow) OnClose() {
	cw.shared.SetClosed(true)
}

// WindowSize はウィンドウサイズを返す。
func (cw *ControlWindow) WindowSize() (int, int) {
	size := cw.Size()
	return size.Width, size.Height
}

// ProgressBar は進捗バーを返す。
func (cw *ControlWindow) ProgressBar() *ProgressBar {
	return cw.progressBar
}

// Position はウィンドウ位置を返す。
func (cw *ControlWindow) Position() (int, int) {
	x, y := cw.GetPos()
	return x, y
}

// GetPos はウィンドウ位置を返す。
func (cw *ControlWindow) GetPos() (int, int) {
	return cw.X(), cw.Y()
}

// SetPosition はウィンドウ位置を設定する。
func (cw *ControlWindow) SetPosition(x, y int) {
	cw.SetPos(x, y)
}

// SetPos はウィンドウ位置を設定する。
func (cw *ControlWindow) SetPos(x, y int) {
	cw.SetX(x)
	cw.SetY(y)
}

// SetPlaying は再生状態を設定する。
func (cw *ControlWindow) SetPlaying(playing bool) {
	if playing {
		cw.shared.EnableFlag(state.STATE_FLAG_PLAYING)
	} else {
		cw.shared.DisableFlag(state.STATE_FLAG_PLAYING)
	}
}

// Playing は再生中か返す。
func (cw *ControlWindow) Playing() bool {
	return cw.shared.HasFlag(state.STATE_FLAG_PLAYING)
}

// SetFrame は現在フレームを設定する。
func (cw *ControlWindow) SetFrame(frame sharedtime.Frame) {
	cw.shared.SetFrame(frame)
}

// Frame は現在フレームを返す。
func (cw *ControlWindow) Frame() sharedtime.Frame {
	return cw.shared.Frame()
}

// SetMaxFrame は最大フレームを設定する。
func (cw *ControlWindow) SetMaxFrame(frame sharedtime.Frame) {
	cw.shared.SetMaxFrame(frame)
}

// MaxFrame は最大フレームを返す。
func (cw *ControlWindow) MaxFrame() sharedtime.Frame {
	return cw.shared.MaxFrame()
}

// SetEnabledInPlaying は再生中の有効化を反映する。
func (cw *ControlWindow) SetEnabledInPlaying(playing bool) {
	if cw.onEnabledInPlaying != nil {
		cw.onEnabledInPlaying(playing)
	}
	if cw.setEnabledInPlaying != nil {
		cw.setEnabledInPlaying(playing)
	}
}

// SetOnEnabledInPlaying は再生中の有効化コールバックを設定する。
func (cw *ControlWindow) SetOnEnabledInPlaying(callback func(playing bool)) {
	cw.onEnabledInPlaying = callback
}

// SetOnChangePlayingPre は再生前コールバックを設定する。
func (cw *ControlWindow) SetOnChangePlayingPre(callback func(playing bool)) {
	cw.onChangePlayingPre = callback
}

// SetOnChangePlayingPost は再生後コールバックを設定する。
func (cw *ControlWindow) SetOnChangePlayingPost(callback func(playing bool)) {
	cw.onChangePlayingPost = callback
}

// OnChangePlayingPre は再生前コールバックを実行する。
func (cw *ControlWindow) OnChangePlayingPre(playing bool) {
	if cw.onChangePlayingPre != nil {
		cw.onChangePlayingPre(playing)
	}
}

// OnChangePlayingPost は再生後コールバックを実行する。
func (cw *ControlWindow) OnChangePlayingPost(playing bool) {
	if cw.onChangePlayingPost != nil {
		cw.onChangePlayingPost(playing)
	}
}

// SetSaveDelta は差分保存の有効可否を設定する。
func (cw *ControlWindow) SetSaveDelta(windowIndex int, enabled bool) {
	cw.shared.SetDeltaSaveEnabled(windowIndex, enabled)
}

// SetSaveDeltaIndex は差分保存対象のインデックスを設定する。
func (cw *ControlWindow) SetSaveDeltaIndex(windowIndex int, index int) {
	cw.shared.SetDeltaSaveIndex(windowIndex, index)
}

// SetFpsLimit はFPS制限を設定する。
func (cw *ControlWindow) SetFpsLimit(limit FpsLimit) {
	interval := sharedtime.Seconds(-1)
	switch limit {
	case FPS_LIMIT_30:
		interval = sharedtime.FpsToSpf(sharedtime.Fps(30))
	case FPS_LIMIT_60:
		interval = sharedtime.FpsToSpf(sharedtime.Fps(60))
	case FPS_LIMIT_UNLIMITED:
		interval = sharedtime.Seconds(-1)
	}
	cw.shared.SetFrameInterval(interval)
	cw.shared.SetFpsLimitTriggered(true)
	cw.updateFpsMenu(limit)
	cw.saveUserInt(config.UserConfigKeyFpsLimit, int(limit))
}

// SetFrameDropEnabled はフレームドロップを設定する。
func (cw *ControlWindow) SetFrameDropEnabled(enabled bool) {
	cw.updateActionChecked(cw.enabledFrameDropAction, enabled)
	cw.setDisplayFlag(state.STATE_FLAG_FRAME_DROP, enabled)
	cw.saveUserBool(config.UserConfigKeyFrameDrop, enabled)
}

// SetWindowLinkageEnabled は画面移動連動を設定する。
func (cw *ControlWindow) SetWindowLinkageEnabled(enabled bool) {
	cw.updateActionChecked(cw.linkWindowAction, enabled)
	cw.setDisplayFlag(state.STATE_FLAG_WINDOW_LINKAGE, enabled)
	cw.saveUserBool(config.UserConfigKeyWindowLinkage, enabled)
}

// SetShowInfoEnabled は情報表示を設定する。
func (cw *ControlWindow) SetShowInfoEnabled(enabled bool) {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_INFO, enabled)
}

// SetDisplayFlag は表示フラグを設定する。
func (cw *ControlWindow) SetDisplayFlag(flag state.StateFlag, enabled bool) {
	switch flag {
	case state.STATE_FLAG_CAMERA_SYNC:
		if enabled {
			cw.setDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_UPPER, false)
			cw.setDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_LOWER, false)
			cw.updateDisplayAction(state.STATE_FLAG_SHOW_OVERRIDE_UPPER, false)
			cw.updateDisplayAction(state.STATE_FLAG_SHOW_OVERRIDE_LOWER, false)
		}
	case state.STATE_FLAG_SHOW_OVERRIDE_UPPER, state.STATE_FLAG_SHOW_OVERRIDE_LOWER:
		if enabled {
			cw.setDisplayFlag(state.STATE_FLAG_CAMERA_SYNC, false)
			cw.updateDisplayAction(state.STATE_FLAG_CAMERA_SYNC, false)
		}
	}
	cw.setDisplayFlag(flag, enabled)
	cw.updateDisplayAction(flag, enabled)
}

// RequestPhysicsReset は物理リセット種別を設定する。
func (cw *ControlWindow) RequestPhysicsReset(resetType state.PhysicsResetType) {
	cw.shared.SetPhysicsResetType(resetType)
}

// SetPhysicsEnabled は物理の有効可否を設定する。
func (cw *ControlWindow) SetPhysicsEnabled(enabled bool) {
	cw.updateActionChecked(cw.enabledPhysicsAction, enabled)
	cw.setDisplayFlag(state.STATE_FLAG_PHYSICS_ENABLED, enabled)
}

// SetModel はモデルを設定する。
func (cw *ControlWindow) SetModel(windowIndex, modelIndex int, modelData *model.PmxModel) {
	cw.shared.SetModel(windowIndex, modelIndex, modelData)
	cw.SetSelectedMaterialIndexes(windowIndex, modelIndex, allMaterialIndexes(modelData))
}

// Model はモデルを取得する。
func (cw *ControlWindow) Model(windowIndex, modelIndex int) *model.PmxModel {
	if raw := cw.shared.Model(windowIndex, modelIndex); raw != nil {
		if m, ok := raw.(*model.PmxModel); ok {
			return m
		}
	}
	return nil
}

// SetMotion はモーションを設定する。
func (cw *ControlWindow) SetMotion(windowIndex, modelIndex int, motionData *motion.VmdMotion) {
	cw.shared.SetMotion(windowIndex, modelIndex, motionData)
}

// Motion はモーションを取得する。
func (cw *ControlWindow) Motion(windowIndex, modelIndex int) *motion.VmdMotion {
	if raw := cw.shared.Motion(windowIndex, modelIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// SetPhysicsWorldMotion は物理ワールド用モーションを設定する。
func (cw *ControlWindow) SetPhysicsWorldMotion(windowIndex int, motionData *motion.VmdMotion) {
	cw.shared.SetPhysicsWorldMotion(windowIndex, motionData)
}

// PhysicsWorldMotion は物理ワールド用モーションを取得する。
func (cw *ControlWindow) PhysicsWorldMotion(windowIndex int) *motion.VmdMotion {
	if raw := cw.shared.PhysicsWorldMotion(windowIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// SetPhysicsModelMotion は物理モデル用モーションを設定する。
func (cw *ControlWindow) SetPhysicsModelMotion(windowIndex, modelIndex int, motionData *motion.VmdMotion) {
	cw.shared.SetPhysicsModelMotion(windowIndex, modelIndex, motionData)
}

// PhysicsModelMotion は物理モデル用モーションを取得する。
func (cw *ControlWindow) PhysicsModelMotion(windowIndex, modelIndex int) *motion.VmdMotion {
	if raw := cw.shared.PhysicsModelMotion(windowIndex, modelIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// SetWindMotion は風用モーションを設定する。
func (cw *ControlWindow) SetWindMotion(windowIndex int, motionData *motion.VmdMotion) {
	cw.shared.SetWindMotion(windowIndex, motionData)
}

// WindMotion は風用モーションを取得する。
func (cw *ControlWindow) WindMotion(windowIndex int) *motion.VmdMotion {
	if raw := cw.shared.WindMotion(windowIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// SetSelectedMaterialIndexes は選択材質を設定する。
func (cw *ControlWindow) SetSelectedMaterialIndexes(windowIndex, modelIndex int, indexes []int) {
	cw.shared.SetSelectedMaterialIndexes(windowIndex, modelIndex, indexes)
}

// SelectedMaterialIndexes は選択材質を取得する。
func (cw *ControlWindow) SelectedMaterialIndexes(windowIndex, modelIndex int) []int {
	return cw.shared.SelectedMaterialIndexes(windowIndex, modelIndex)
}

// SetDeltaMotion は差分モーションを設定する。
func (cw *ControlWindow) SetDeltaMotion(windowIndex, modelIndex, deltaIndex int, motionData *motion.VmdMotion) {
	cw.shared.SetDeltaMotion(windowIndex, modelIndex, deltaIndex, motionData)
}

// DeltaMotion は差分モーションを取得する。
func (cw *ControlWindow) DeltaMotion(windowIndex, modelIndex, deltaIndex int) *motion.VmdMotion {
	if raw := cw.shared.DeltaMotion(windowIndex, modelIndex, deltaIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// ClearDeltaMotion は差分モーションを削除する。
func (cw *ControlWindow) ClearDeltaMotion(windowIndex, modelIndex int) {
	cw.shared.ClearDeltaMotion(windowIndex, modelIndex)
}

// DeltaMotionCount は差分モーション数を返す。
func (cw *ControlWindow) DeltaMotionCount(windowIndex, modelIndex int) int {
	return cw.shared.DeltaMotionCount(windowIndex, modelIndex)
}

// appTitle はアプリタイトルを生成する。
func (cw *ControlWindow) appTitle() string {
	if cw.appConfig == nil {
		return ""
	}
	if cw.appConfig.AppName == "" {
		return cw.appConfig.Version
	}
	return fmt.Sprintf("%s %s", cw.appConfig.AppName, cw.appConfig.Version)
}

// onClosing はウィンドウ終了時の処理を行う。
func (cw *ControlWindow) onClosing(canceled *bool) {
	if cw.appConfig == nil || !cw.appConfig.IsCloseConfirmEnabled() {
		cw.shared.SetClosed(true)
		return
	}
	if cw.shared.IsClosed() {
		return
	}
	if result := walk.MsgBox(
		cw,
		cw.t("終了確認"),
		cw.t("終了確認メッセージ"),
		walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel,
	); result == walk.DlgCmdOK {
		cw.shared.SetClosed(true)
		return
	}
	*canceled = true
}

// onActivate はアクティブ化時の処理を行う。
func (cw *ControlWindow) onActivate() {
	cw.shared.SetFocusedWindowHandle(state.WindowHandle(uintptr(cw.Handle())))
	cw.shared.SetControlWindowFocused(true)
	cw.shared.SetAllViewerWindowsFocused(false)
}

// onEnterSizeMove は移動開始時の処理を行う。
func (cw *ControlWindow) onEnterSizeMove() {
	if !cw.shared.HasFlag(state.STATE_FLAG_WINDOW_LINKAGE) {
		return
	}
	x, y := cw.GetPos()
	cw.shared.SetControlWindowPosition(state.WindowPosition{X: x, Y: y})
}

// onExitSizeMove は移動終了時の処理を行う。
func (cw *ControlWindow) onExitSizeMove() {
	if !cw.shared.HasFlag(state.STATE_FLAG_WINDOW_LINKAGE) {
		return
	}
	x, y := cw.GetPos()
	prev := cw.shared.ControlWindowPosition()
	cw.shared.SetControlWindowPosition(state.WindowPosition{X: x, Y: y, DiffX: x - prev.X, DiffY: y - prev.Y})
	cw.shared.SetControlWindowMoving(true)
}

// onMinimize は最小化時の処理を行う。
func (cw *ControlWindow) onMinimize() {
	if !cw.shared.IsAllViewerWindowsReady() {
		return
	}
	for i := 0; i < cw.viewerCount; i++ {
		cw.shared.SyncMinimize(i)
	}
}

// onRestore は復元時の処理を行う。
func (cw *ControlWindow) onRestore() {
	if !cw.shared.IsAllViewerWindowsReady() {
		return
	}
	for i := 0; i < cw.viewerCount; i++ {
		cw.shared.SyncRestore(i)
	}
}

// buildViewerMenu はビューワーメニューを構築する。
func (cw *ControlWindow) buildViewerMenu() declarative.Menu {
	return declarative.Menu{
		Text: cw.t("&ビューワー"),
		Items: []declarative.MenuItem{
			declarative.Action{Text: cw.t("&フレームドロップON"), Checkable: true, OnTriggered: cw.TriggerEnabledFrameDrop, AssignTo: &cw.enabledFrameDropAction},
			declarative.Menu{Text: cw.t("&fps制限"), Items: []declarative.MenuItem{
				declarative.Action{Text: cw.t("&30fps制限"), Checkable: true, OnTriggered: cw.TriggerFps30Limit, AssignTo: &cw.limitFps30Action},
				declarative.Action{Text: cw.t("&60fps制限"), Checkable: true, OnTriggered: cw.TriggerFps60Limit, AssignTo: &cw.limitFps60Action},
				declarative.Action{Text: cw.t("&fps無制限"), Checkable: true, OnTriggered: cw.TriggerUnLimitFps, AssignTo: &cw.limitFpsUnLimitAction},
			}},
			declarative.Action{Text: cw.t("&情報表示"), Checkable: true, OnTriggered: cw.TriggerShowInfo, AssignTo: &cw.showInfoAction},
			declarative.Separator{},
			declarative.Action{Text: cw.t("&物理ON/OFF"), Checkable: true, OnTriggered: cw.TriggerEnabledPhysics, AssignTo: &cw.enabledPhysicsAction},
			declarative.Action{Text: cw.t("&物理リセット"), OnTriggered: cw.TriggerPhysicsReset, AssignTo: &cw.physicsResetAction},
			declarative.Separator{},
			declarative.Action{Text: cw.t("&法線表示"), Checkable: true, OnTriggered: cw.TriggerShowNormal, AssignTo: &cw.showNormalAction},
			declarative.Action{Text: cw.t("&ワイヤーフレーム表示"), Checkable: true, OnTriggered: cw.TriggerShowWire, AssignTo: &cw.showWireAction},
			declarative.Action{Text: cw.t("&カメラ同期"), Checkable: true, OnTriggered: cw.TriggerCameraSync, AssignTo: &cw.cameraSyncAction},
			declarative.Menu{Text: cw.t("&サブビューワーオーバーレイ"), Items: []declarative.MenuItem{
				declarative.Action{Text: cw.t("&上半身合わせ"), Checkable: true, OnTriggered: cw.TriggerShowOverrideUpper, AssignTo: &cw.showOverrideUpperAction},
				declarative.Action{Text: cw.t("&下半身合わせ"), Checkable: true, OnTriggered: cw.TriggerShowOverrideLower, AssignTo: &cw.showOverrideLowerAction},
				declarative.Action{Text: cw.t("&カメラ合わせなし"), Checkable: true, OnTriggered: cw.TriggerShowOverrideNone, AssignTo: &cw.showOverrideNoneAction},
				declarative.Action{Text: cw.t("&サブビューワーオーバーレイの使い方"), OnTriggered: cw.showOverrideHelp},
			}},
			declarative.Action{Text: cw.t("&選択頂点表示"), Checkable: true, OnTriggered: cw.TriggerShowSelectedVertex, AssignTo: &cw.showSelectedVertexAction},
			declarative.Menu{Text: cw.t("&ボーン表示"), Items: []declarative.MenuItem{
				declarative.Action{Text: cw.t("&全ボーン"), Checkable: true, OnTriggered: cw.TriggerShowBoneAll, AssignTo: &cw.showBoneAllAction},
				declarative.Action{Text: cw.t("&IKボーン"), Checkable: true, OnTriggered: cw.TriggerShowBoneIk, AssignTo: &cw.showBoneIkAction},
				declarative.Action{Text: cw.t("&付与親ボーン"), Checkable: true, OnTriggered: cw.TriggerShowBoneEffector, AssignTo: &cw.showBoneEffectorAction},
				declarative.Action{Text: cw.t("&軸制限ボーン"), Checkable: true, OnTriggered: cw.TriggerShowBoneFixed, AssignTo: &cw.showBoneFixedAction},
				declarative.Action{Text: cw.t("&回転ボーン"), Checkable: true, OnTriggered: cw.TriggerShowBoneRotate, AssignTo: &cw.showBoneRotateAction},
				declarative.Action{Text: cw.t("&移動ボーン"), Checkable: true, OnTriggered: cw.TriggerShowBoneTranslate, AssignTo: &cw.showBoneTranslateAction},
				declarative.Action{Text: cw.t("&表示ボーン"), Checkable: true, OnTriggered: cw.TriggerShowBoneVisible, AssignTo: &cw.showBoneVisibleAction},
				declarative.Action{Text: cw.t("&ボーン表示の使い方"), OnTriggered: cw.showBoneHelp},
			}},
			declarative.Menu{Text: cw.t("&剛体表示"), Items: []declarative.MenuItem{
				declarative.Action{Text: cw.t("&前面表示"), Checkable: true, OnTriggered: cw.TriggerShowRigidBodyFront, AssignTo: &cw.showRigidBodyFrontAction},
				declarative.Action{Text: cw.t("&埋め込み表示"), Checkable: true, OnTriggered: cw.TriggerShowRigidBodyBack, AssignTo: &cw.showRigidBodyBackAction},
				declarative.Action{Text: cw.t("&剛体表示の使い方"), OnTriggered: cw.showRigidBodyHelp},
			}},
			declarative.Action{Text: cw.t("&ジョイント表示"), Checkable: true, OnTriggered: cw.TriggerShowJoint, AssignTo: &cw.showJointAction},
			declarative.Separator{},
			declarative.Action{Text: cw.t("&ビューワーの使い方"), OnTriggered: cw.showViewerHelp},
		},
	}
}

// buildControllerMenuItems はコントローラーメニュー項目を構築する。
func (cw *ControlWindow) buildControllerMenuItems() []declarative.MenuItem {
	items := []declarative.MenuItem{
		declarative.Action{Text: cw.t("&使い方"), OnTriggered: cw.showControllerHelp},
		declarative.Separator{},
		declarative.Action{Text: cw.t("&画面移動連動"), Checkable: true, OnTriggered: cw.triggerWindowLinkage, AssignTo: &cw.linkWindowAction},
		declarative.Separator{},
		declarative.Action{Text: cw.t("&デバッグログ表示"), Checkable: true, OnTriggered: cw.triggerLogLevelDebug, AssignTo: &cw.logLevelDebugAction},
		declarative.Action{Text: cw.t("&ログ保存"), OnTriggered: cw.triggerSaveLog},
	}

	if cw.appConfig != nil && cw.appConfig.IsDev() {
		items = append(items,
			declarative.Action{Text: cw.t("&モーション冗長ログ表示"), Checkable: true, OnTriggered: cw.triggerLogLevelVerbose, AssignTo: &cw.logLevelVerboseAction},
			declarative.Action{Text: cw.t("&IK冗長ログ表示"), Checkable: true, OnTriggered: cw.triggerLogLevelIkVerbose, AssignTo: &cw.logLevelIkVerboseAction},
			declarative.Action{Text: cw.t("&物理冗長ログ表示"), Checkable: true, OnTriggered: cw.triggerLogLevelPhysicsVerbose, AssignTo: &cw.logLevelPhysicsVerboseAction},
			declarative.Action{Text: cw.t("&ビューワー冗長ログ表示"), Checkable: true, OnTriggered: cw.triggerLogLevelViewerVerbose, AssignTo: &cw.logLevelViewerVerboseAction},
		)
	}
	return items
}

// buildLanguageMenu は言語メニュー項目を構築する。
func (cw *ControlWindow) buildLanguageMenu() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{Text: "日本語", OnTriggered: func() { cw.onChangeLanguage(langJa) }},
		declarative.Action{Text: "English", OnTriggered: func() { cw.onChangeLanguage(langEn) }},
		declarative.Action{Text: "中文", OnTriggered: func() { cw.onChangeLanguage(langZh) }},
		declarative.Action{Text: "한국어", OnTriggered: func() { cw.onChangeLanguage(langKo) }},
	}
}

// onChangeLanguage は言語変更を行う。
func (cw *ControlWindow) onChangeLanguage(lang sharedi18n.LangCode) {
	if result := walk.MsgBox(cw, cw.t("言語変更"), cw.t("言語変更メッセージ"), walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel); result != walk.DlgCmdOK {
		return
	}
	if cw.translator == nil {
		return
	}
	if _, err := cw.translator.SetLang(lang); err != nil {
		cw.loggerOrDefault().Error("言語設定の保存に失敗しました: %s", err.Error())
		walk.MsgBox(cw, cw.t("言語変更"), err.Error(), walk.MsgBoxIconError)
		return
	}
	cw.shared.SetClosed(true)
}

// TriggerEnabledFrameDrop はフレームドロップを切り替える。
func (cw *ControlWindow) TriggerEnabledFrameDrop() {
	cw.SetFrameDropEnabled(cw.actionChecked(cw.enabledFrameDropAction))
}

// TriggerEnabledPhysics は物理の有効可否を切り替える。
func (cw *ControlWindow) TriggerEnabledPhysics() {
	cw.SetPhysicsEnabled(cw.actionChecked(cw.enabledPhysicsAction))
}

// TriggerPhysicsReset は物理リセットを要求する。
func (cw *ControlWindow) TriggerPhysicsReset() {
	cw.RequestPhysicsReset(state.PHYSICS_RESET_TYPE_START_FRAME)
}

// TriggerShowNormal は法線表示を切り替える。
func (cw *ControlWindow) TriggerShowNormal() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_NORMAL, cw.actionChecked(cw.showNormalAction))
}

// TriggerShowWire はワイヤーフレーム表示を切り替える。
func (cw *ControlWindow) TriggerShowWire() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_WIRE, cw.actionChecked(cw.showWireAction))
}

// TriggerShowOverrideUpper は上半身オーバーレイを切り替える。
func (cw *ControlWindow) TriggerShowOverrideUpper() {
	enabled := cw.actionChecked(cw.showOverrideUpperAction)
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_UPPER, enabled)
	if enabled {
		cw.SetDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_LOWER, false)
		cw.SetDisplayFlag(state.STATE_FLAG_CAMERA_SYNC, false)
	}
}

// TriggerShowOverrideLower は下半身オーバーレイを切り替える。
func (cw *ControlWindow) TriggerShowOverrideLower() {
	enabled := cw.actionChecked(cw.showOverrideLowerAction)
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_LOWER, enabled)
	if enabled {
		cw.SetDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_UPPER, false)
		cw.SetDisplayFlag(state.STATE_FLAG_CAMERA_SYNC, false)
	}
}

// TriggerShowOverrideNone はカメラ合わせなしを切り替える。
func (cw *ControlWindow) TriggerShowOverrideNone() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_NONE, cw.actionChecked(cw.showOverrideNoneAction))
}

// TriggerShowSelectedVertex は選択頂点表示を切り替える。
func (cw *ControlWindow) TriggerShowSelectedVertex() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX, cw.actionChecked(cw.showSelectedVertexAction))
}

// TriggerShowBoneAll は全ボーン表示を切り替える。
func (cw *ControlWindow) TriggerShowBoneAll() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_BONE_ALL, cw.actionChecked(cw.showBoneAllAction))
}

// TriggerShowBoneIk はIKボーン表示を切り替える。
func (cw *ControlWindow) TriggerShowBoneIk() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_BONE_IK, cw.actionChecked(cw.showBoneIkAction))
}

// TriggerShowBoneEffector は付与親ボーン表示を切り替える。
func (cw *ControlWindow) TriggerShowBoneEffector() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_BONE_EFFECTOR, cw.actionChecked(cw.showBoneEffectorAction))
}

// TriggerShowBoneFixed は軸制限ボーン表示を切り替える。
func (cw *ControlWindow) TriggerShowBoneFixed() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_BONE_FIXED, cw.actionChecked(cw.showBoneFixedAction))
}

// TriggerShowBoneRotate は回転ボーン表示を切り替える。
func (cw *ControlWindow) TriggerShowBoneRotate() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_BONE_ROTATE, cw.actionChecked(cw.showBoneRotateAction))
}

// TriggerShowBoneTranslate は移動ボーン表示を切り替える。
func (cw *ControlWindow) TriggerShowBoneTranslate() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_BONE_TRANSLATE, cw.actionChecked(cw.showBoneTranslateAction))
}

// TriggerShowBoneVisible は表示ボーン表示を切り替える。
func (cw *ControlWindow) TriggerShowBoneVisible() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_BONE_VISIBLE, cw.actionChecked(cw.showBoneVisibleAction))
}

// TriggerShowRigidBodyFront は剛体前面表示を切り替える。
func (cw *ControlWindow) TriggerShowRigidBodyFront() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_RIGID_BODY_FRONT, cw.actionChecked(cw.showRigidBodyFrontAction))
}

// TriggerShowRigidBodyBack は剛体埋め込み表示を切り替える。
func (cw *ControlWindow) TriggerShowRigidBodyBack() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_RIGID_BODY_BACK, cw.actionChecked(cw.showRigidBodyBackAction))
}

// TriggerShowJoint はジョイント表示を切り替える。
func (cw *ControlWindow) TriggerShowJoint() {
	cw.SetDisplayFlag(state.STATE_FLAG_SHOW_JOINT, cw.actionChecked(cw.showJointAction))
}

// TriggerShowInfo は情報表示を切り替える。
func (cw *ControlWindow) TriggerShowInfo() {
	cw.SetShowInfoEnabled(cw.actionChecked(cw.showInfoAction))
}

// TriggerFps30Limit は30fps制限を設定する。
func (cw *ControlWindow) TriggerFps30Limit() {
	cw.SetFpsLimit(FPS_LIMIT_30)
}

// TriggerFps60Limit は60fps制限を設定する。
func (cw *ControlWindow) TriggerFps60Limit() {
	cw.SetFpsLimit(FPS_LIMIT_60)
}

// TriggerUnLimitFps はFPS無制限を設定する。
func (cw *ControlWindow) TriggerUnLimitFps() {
	cw.SetFpsLimit(FPS_LIMIT_UNLIMITED)
}

// TriggerCameraSync はカメラ同期を切り替える。
func (cw *ControlWindow) TriggerCameraSync() {
	enabled := cw.actionChecked(cw.cameraSyncAction)
	cw.SetDisplayFlag(state.STATE_FLAG_CAMERA_SYNC, enabled)
	if enabled {
		cw.SetDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_UPPER, false)
		cw.SetDisplayFlag(state.STATE_FLAG_SHOW_OVERRIDE_LOWER, false)
	}
}

// triggerLogLevelDebug はデバッグログを切り替える。
func (cw *ControlWindow) triggerLogLevelDebug() {
	enabled := cw.actionChecked(cw.logLevelDebugAction)
	cw.resetVerboseActions()
	cw.updateActionChecked(cw.logLevelDebugAction, enabled)
	level := logging.LOG_LEVEL_INFO
	if enabled {
		level = logging.LOG_LEVEL_DEBUG
	}
	cw.loggerOrDefault().SetLevel(level)
}

// triggerLogLevelVerbose はモーション冗長ログを切り替える。
func (cw *ControlWindow) triggerLogLevelVerbose() {
	enabled := cw.actionChecked(cw.logLevelVerboseAction)
	cw.setVerbose(logging.VERBOSE_INDEX_MOTION, enabled, "motion")
	cw.updateActionChecked(cw.logLevelVerboseAction, enabled)
}

// triggerLogLevelIkVerbose はIK冗長ログを切り替える。
func (cw *ControlWindow) triggerLogLevelIkVerbose() {
	enabled := cw.actionChecked(cw.logLevelIkVerboseAction)
	cw.setVerbose(logging.VERBOSE_INDEX_IK, enabled, "ik")
	cw.updateActionChecked(cw.logLevelIkVerboseAction, enabled)
}

// triggerLogLevelPhysicsVerbose は物理冗長ログを切り替える。
func (cw *ControlWindow) triggerLogLevelPhysicsVerbose() {
	enabled := cw.actionChecked(cw.logLevelPhysicsVerboseAction)
	cw.setVerbose(logging.VERBOSE_INDEX_PHYSICS, enabled, "physics")
	cw.updateActionChecked(cw.logLevelPhysicsVerboseAction, enabled)
}

// triggerLogLevelViewerVerbose はビューワー冗長ログを切り替える。
func (cw *ControlWindow) triggerLogLevelViewerVerbose() {
	enabled := cw.actionChecked(cw.logLevelViewerVerboseAction)
	cw.setVerbose(logging.VERBOSE_INDEX_VIEWER, enabled, "viewer")
	cw.updateActionChecked(cw.logLevelViewerVerboseAction, enabled)
}

// triggerWindowLinkage は画面移動連動を切り替える。
func (cw *ControlWindow) triggerWindowLinkage() {
	cw.SetWindowLinkageEnabled(cw.actionChecked(cw.linkWindowAction))
}

// triggerSaveLog はログ保存を行う。
func (cw *ControlWindow) triggerSaveLog() {
	text := logging.ConsoleText()
	path, err := mfile.SaveConsoleSnapshot(cw.userConfig, "console", text)
	if err != nil {
		cw.loggerOrDefault().Error("ログ保存に失敗しました: %s", err.Error())
		walk.MsgBox(cw, cw.t("ログ保存"), err.Error(), walk.MsgBoxIconError)
		return
	}
	cw.loggerOrDefault().Info("ログを保存しました: %s", path)
}

// showControllerHelp はコントローラーの使い方を表示する。
func (cw *ControlWindow) showControllerHelp() {
	cw.infoLineTitle(cw.t("コントローラーウィンドウの使い方"), cw.t("コントローラーウィンドウの使い方メッセージ"))
}

// showOverrideHelp はオーバーレイの使い方を表示する。
func (cw *ControlWindow) showOverrideHelp() {
	cw.infoLineTitle(cw.t("&サブビューワーオーバーレイの使い方"), cw.t("サブビューワーオーバーレイの使い方メッセージ"))
}

// showBoneHelp はボーン表示の使い方を表示する。
func (cw *ControlWindow) showBoneHelp() {
	cw.infoLineTitle(cw.t("&ボーン表示の使い方"), cw.t("ボーン表示の使い方メッセージ"))
}

// showRigidBodyHelp は剛体表示の使い方を表示する。
func (cw *ControlWindow) showRigidBodyHelp() {
	cw.infoLineTitle(cw.t("&剛体表示の使い方"), cw.t("剛体表示の使い方メッセージ"))
}

// showViewerHelp はビューワーの使い方を表示する。
func (cw *ControlWindow) showViewerHelp() {
	cw.infoLineTitle(cw.t("&ビューワーの使い方"), cw.t("ビューワーの使い方メッセージ"))
}

// updateFpsMenu はFPS制限メニューの状態を更新する。
func (cw *ControlWindow) updateFpsMenu(limit FpsLimit) {
	cw.updateActionChecked(cw.limitFps30Action, limit == FPS_LIMIT_30)
	cw.updateActionChecked(cw.limitFps60Action, limit == FPS_LIMIT_60)
	cw.updateActionChecked(cw.limitFpsUnLimitAction, limit == FPS_LIMIT_UNLIMITED)
}

// applyUserConfig はユーザー設定を反映する。
func (cw *ControlWindow) applyUserConfig() {
	if cw.userConfig == nil {
		return
	}
	linkage, err := cw.userConfig.GetBool(config.UserConfigKeyWindowLinkage, true)
	if err == nil {
		cw.SetWindowLinkageEnabled(linkage)
	}
	fpsLimit, err := cw.userConfig.GetInt(config.UserConfigKeyFpsLimit, int(FPS_LIMIT_30))
	if err == nil {
		switch fpsLimit {
		case int(FPS_LIMIT_60):
			cw.SetFpsLimit(FPS_LIMIT_60)
		case int(FPS_LIMIT_UNLIMITED):
			cw.SetFpsLimit(FPS_LIMIT_UNLIMITED)
		default:
			cw.SetFpsLimit(FPS_LIMIT_30)
		}
	}
	frameDrop, err := cw.userConfig.GetBool(config.UserConfigKeyFrameDrop, true)
	if err == nil {
		cw.SetFrameDropEnabled(frameDrop)
	}
}

// saveUserBool はユーザー設定を保存する。
func (cw *ControlWindow) saveUserBool(key string, value bool) {
	if cw.userConfig == nil {
		return
	}
	_ = cw.userConfig.SetBool(key, value)
}

// saveUserInt はユーザー設定のint値を保存する。
func (cw *ControlWindow) saveUserInt(key string, value int) {
	if cw.userConfig == nil {
		return
	}
	_ = cw.userConfig.SetInt(key, value)
}

// setDisplayFlag は表示フラグを反映する。
func (cw *ControlWindow) setDisplayFlag(flag state.StateFlag, enabled bool) {
	if enabled {
		cw.shared.EnableFlag(flag)
	} else {
		cw.shared.DisableFlag(flag)
	}
}

// updateDisplayAction は表示メニューの状態を更新する。
func (cw *ControlWindow) updateDisplayAction(flag state.StateFlag, enabled bool) {
	switch flag {
	case state.STATE_FLAG_SHOW_NORMAL:
		cw.updateActionChecked(cw.showNormalAction, enabled)
	case state.STATE_FLAG_SHOW_WIRE:
		cw.updateActionChecked(cw.showWireAction, enabled)
	case state.STATE_FLAG_SHOW_OVERRIDE_UPPER:
		cw.updateActionChecked(cw.showOverrideUpperAction, enabled)
	case state.STATE_FLAG_SHOW_OVERRIDE_LOWER:
		cw.updateActionChecked(cw.showOverrideLowerAction, enabled)
	case state.STATE_FLAG_SHOW_OVERRIDE_NONE:
		cw.updateActionChecked(cw.showOverrideNoneAction, enabled)
	case state.STATE_FLAG_SHOW_SELECTED_VERTEX:
		cw.updateActionChecked(cw.showSelectedVertexAction, enabled)
	case state.STATE_FLAG_SHOW_BONE_ALL:
		cw.updateActionChecked(cw.showBoneAllAction, enabled)
	case state.STATE_FLAG_SHOW_BONE_IK:
		cw.updateActionChecked(cw.showBoneIkAction, enabled)
	case state.STATE_FLAG_SHOW_BONE_EFFECTOR:
		cw.updateActionChecked(cw.showBoneEffectorAction, enabled)
	case state.STATE_FLAG_SHOW_BONE_FIXED:
		cw.updateActionChecked(cw.showBoneFixedAction, enabled)
	case state.STATE_FLAG_SHOW_BONE_ROTATE:
		cw.updateActionChecked(cw.showBoneRotateAction, enabled)
	case state.STATE_FLAG_SHOW_BONE_TRANSLATE:
		cw.updateActionChecked(cw.showBoneTranslateAction, enabled)
	case state.STATE_FLAG_SHOW_BONE_VISIBLE:
		cw.updateActionChecked(cw.showBoneVisibleAction, enabled)
	case state.STATE_FLAG_SHOW_RIGID_BODY_FRONT:
		cw.updateActionChecked(cw.showRigidBodyFrontAction, enabled)
	case state.STATE_FLAG_SHOW_RIGID_BODY_BACK:
		cw.updateActionChecked(cw.showRigidBodyBackAction, enabled)
	case state.STATE_FLAG_SHOW_JOINT:
		cw.updateActionChecked(cw.showJointAction, enabled)
	case state.STATE_FLAG_SHOW_INFO:
		cw.updateActionChecked(cw.showInfoAction, enabled)
	case state.STATE_FLAG_CAMERA_SYNC:
		cw.updateActionChecked(cw.cameraSyncAction, enabled)
	case state.STATE_FLAG_FRAME_DROP:
		cw.updateActionChecked(cw.enabledFrameDropAction, enabled)
	case state.STATE_FLAG_WINDOW_LINKAGE:
		cw.updateActionChecked(cw.linkWindowAction, enabled)
	case state.STATE_FLAG_PHYSICS_ENABLED:
		cw.updateActionChecked(cw.enabledPhysicsAction, enabled)
	}
}

// actionChecked はアクションのチェック状態を返す。
func (cw *ControlWindow) actionChecked(action *walk.Action) bool {
	if action == nil {
		return false
	}
	return action.Checked()
}

// updateActionChecked はアクションのチェック状態を更新する。
func (cw *ControlWindow) updateActionChecked(action *walk.Action, enabled bool) {
	if action != nil {
		action.SetChecked(enabled)
	}
}

// setVerbose は冗長ログ設定を切り替える。
func (cw *ControlWindow) setVerbose(index logging.VerboseIndex, enabled bool, label string) {
	cw.resetVerboseActions()
	if enabled {
		cw.enableVerbose(index, label)
	}
}

// resetVerboseActions は冗長ログ設定をリセットする。
func (cw *ControlWindow) resetVerboseActions() {
	cw.updateActionChecked(cw.logLevelDebugAction, false)
	cw.updateActionChecked(cw.logLevelVerboseAction, false)
	cw.updateActionChecked(cw.logLevelIkVerboseAction, false)
	cw.updateActionChecked(cw.logLevelPhysicsVerboseAction, false)
	cw.updateActionChecked(cw.logLevelViewerVerboseAction, false)
	cw.loggerOrDefault().SetLevel(logging.LOG_LEVEL_INFO)
	cw.disableVerbose(logging.VERBOSE_INDEX_MOTION)
	cw.disableVerbose(logging.VERBOSE_INDEX_IK)
	cw.disableVerbose(logging.VERBOSE_INDEX_PHYSICS)
	cw.disableVerbose(logging.VERBOSE_INDEX_VIEWER)
}

// enableVerbose は冗長ログを有効化する。
func (cw *ControlWindow) enableVerbose(index logging.VerboseIndex, label string) {
	logger := cw.loggerOrDefault()
	logger.SetLevel(logging.LOG_LEVEL_VERBOSE)
	logger.EnableVerbose(index)
	if _, ok := cw.verboseSinks[index]; !ok {
		_, sink, err := mfile.OpenVerboseLogStream(cw.userConfig, label)
		if err != nil {
			logger.Error("冗長ログの開始に失敗しました: %s", err.Error())
			return
		}
		cw.verboseSinks[index] = sink
	}
	logger.AttachVerboseSink(index, cw.verboseSinks[index])
}

// disableVerbose は冗長ログを無効化する。
func (cw *ControlWindow) disableVerbose(index logging.VerboseIndex) {
	logger := cw.loggerOrDefault()
	logger.DisableVerbose(index)
	if sink, ok := cw.verboseSinks[index]; ok && sink != nil {
		_ = sink.Close()
		delete(cw.verboseSinks, index)
	}
}

// closeVerboseSinks は冗長ログの出力先を閉じる。
func (cw *ControlWindow) closeVerboseSinks() {
	for idx, sink := range cw.verboseSinks {
		if sink != nil {
			_ = sink.Close()
		}
		delete(cw.verboseSinks, idx)
	}
}

// allMaterialIndexes は全材質インデックスを返す。
func allMaterialIndexes(modelData *model.PmxModel) []int {
	if modelData == nil || modelData.Materials == nil {
		return []int{}
	}
	count := modelData.Materials.Len()
	indexes := make([]int, 0, count)
	for i := 0; i < count; i++ {
		indexes = append(indexes, i)
	}
	return indexes
}
