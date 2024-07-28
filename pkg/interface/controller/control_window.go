//go:build windows
// +build windows

package controller

import (
	"fmt"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"

	"github.com/miu200521358/walk/pkg/declarative"
)

type ControlWindow struct {
	*walk.MainWindow
	*controlState                               // 操作状態
	tabWidget                *widget.MTabWidget // タブウィジェット
	appConfig                *mconfig.AppConfig // アプリケーション設定
	spfLimit                 float64            // FPS制限
	enabledFrameDropAction   *walk.Action       // フレームドロップON/OFF
	enabledPhysicsAction     *walk.Action       // 物理ON/OFF
	physicsResetAction       *walk.Action       // 物理リセット
	showNormalAction         *walk.Action       // ボーンデバッグ表示
	showWireAction           *walk.Action       // ワイヤーフレームデバッグ表示
	showSelectedVertexAction *walk.Action       // 選択頂点デバッグ表示
	showBoneAllAction        *walk.Action       // 全ボーンデバッグ表示
	showBoneIkAction         *walk.Action       // IKボーンデバッグ表示
	showBoneEffectorAction   *walk.Action       // 付与親ボーンデバッグ表示
	showBoneFixedAction      *walk.Action       // 軸制限ボーンデバッグ表示
	showBoneRotateAction     *walk.Action       // 回転ボーンデバッグ表示
	showBoneTranslateAction  *walk.Action       // 移動ボーンデバッグ表示
	showBoneVisibleAction    *walk.Action       // 表示ボーンデバッグ表示
	showRigidBodyFrontAction *walk.Action       // 剛体デバッグ表示(前面)
	showRigidBodyBackAction  *walk.Action       // 剛体デバッグ表示(埋め込み)
	showJointAction          *walk.Action       // ジョイントデバッグ表示
	showInfoAction           *walk.Action       // 情報デバッグ表示
	limitFps30Action         *walk.Action       // 30FPS制限
	limitFps60Action         *walk.Action       // 60FPS制限
	limitFpsUnLimitAction    *walk.Action       // FPS無制限
	logLevelDebugAction      *walk.Action       // デバッグメッセージ表示
	logLevelVerboseAction    *walk.Action       // 冗長メッセージ表示
	logLevelIkVerboseAction  *walk.Action       // IK冗長メッセージ表示
}

func NewControlWindow(
	appConfig *mconfig.AppConfig,
	controlState *controlState,
	helpMenuItemsFunc func() []declarative.MenuItem,
) *ControlWindow {
	controlWindow := &ControlWindow{
		controlState: controlState,
		appConfig:    appConfig,
		spfLimit:     1 / 30.0,
	}
	controlState.SetControlWindow(controlWindow)

	logMenuItems := []declarative.MenuItem{
		declarative.Action{
			Text: mi18n.T("&使い方"),
			OnTriggered: func() {
				mlog.ILT(mi18n.T("メイン画面の使い方"), mi18n.T("メイン画面の使い方メッセージ"))
			},
		},
		declarative.Separator{},
		declarative.Action{
			Text:        mi18n.T("&デバッグログ表示"),
			Checkable:   true,
			OnTriggered: controlWindow.logLevelTriggered,
			AssignTo:    &controlWindow.logLevelDebugAction,
		},
	}

	if !appConfig.IsEnvProd() {
		// 開発時のみ冗長ログ表示を追加
		logMenuItems = append(logMenuItems,
			declarative.Separator{})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&冗長ログ表示"),
				Checkable:   true,
				OnTriggered: controlWindow.logLevelTriggered,
				AssignTo:    &controlWindow.logLevelVerboseAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&IK冗長ログ表示"),
				Checkable:   true,
				OnTriggered: controlWindow.logLevelTriggered,
				AssignTo:    &controlWindow.logLevelIkVerboseAction,
			})
	}

	fpsLImitMenuItems := []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("&30fps制限"),
			Checkable:   true,
			OnTriggered: controlWindow.onTriggerFps30Limit,
			AssignTo:    &controlWindow.limitFps30Action,
		},
		declarative.Action{
			Text:        mi18n.T("&60fps制限"),
			Checkable:   true,
			OnTriggered: controlWindow.onTriggerFps60Limit,
			AssignTo:    &controlWindow.limitFps60Action,
		},
		declarative.Action{
			Text:        mi18n.T("&fps無制限"),
			Checkable:   true,
			OnTriggered: controlWindow.onTriggerUnLimitFps,
			AssignTo:    &controlWindow.limitFpsUnLimitAction,
		},
	}

	if err := (declarative.MainWindow{
		AssignTo: &controlWindow.MainWindow,
		Title:    fmt.Sprintf("%s %s", appConfig.Name, appConfig.Version),
		Size:     app.GetWindowSize(appConfig.ControlWindowSize.Width, appConfig.ControlWindowSize.Height),
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear, MarginsZero: true, SpacingZero: true},
		MenuItems: []declarative.MenuItem{
			declarative.Menu{
				Text: mi18n.T("&ビューワー"),
				Items: []declarative.MenuItem{
					declarative.Action{
						Text:        mi18n.T("&フレームドロップON"),
						Checkable:   true,
						OnTriggered: controlWindow.onTriggerEnabledFrameDrop,
						AssignTo:    &controlWindow.enabledFrameDropAction,
					},
					declarative.Menu{
						Text:  mi18n.T("&fps制限"),
						Items: fpsLImitMenuItems,
					},
					declarative.Action{
						Text:        mi18n.T("&情報表示"),
						Checkable:   true,
						OnTriggered: controlWindow.onTriggerShowInfo,
						AssignTo:    &controlWindow.showInfoAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&物理ON/OFF"),
						Checkable:   true,
						OnTriggered: controlWindow.onTriggerEnabledPhysics,
						AssignTo:    &controlWindow.enabledPhysicsAction,
					},
					declarative.Action{
						Text:        mi18n.T("&物理リセット"),
						OnTriggered: controlWindow.onTriggerPhysicsReset,
						AssignTo:    &controlWindow.physicsResetAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&法線表示"),
						Checkable:   true,
						OnTriggered: controlWindow.onTriggerShowNormal,
						AssignTo:    &controlWindow.showNormalAction,
					},
					declarative.Action{
						Text:        mi18n.T("&ワイヤーフレーム表示"),
						Checkable:   true,
						OnTriggered: controlWindow.onTriggerShowWire,
						AssignTo:    &controlWindow.showWireAction,
					},
					declarative.Action{
						Text:        mi18n.T("&選択頂点表示"),
						Checkable:   true,
						OnTriggered: controlWindow.onTriggerShowSelectedVertex,
						AssignTo:    &controlWindow.showSelectedVertexAction,
					},
					declarative.Separator{},
					declarative.Menu{
						Text: mi18n.T("&ボーン表示"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&全ボーン"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowBoneAll,
								AssignTo:    &controlWindow.showBoneAllAction,
							},
							declarative.Separator{},
							declarative.Action{
								Text:        mi18n.T("&IKボーン"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowBoneIk,
								AssignTo:    &controlWindow.showBoneIkAction,
							},
							declarative.Action{
								Text:        mi18n.T("&付与親ボーン"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowBoneEffector,
								AssignTo:    &controlWindow.showBoneEffectorAction,
							},
							declarative.Action{
								Text:        mi18n.T("&軸制限ボーン"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowBoneFixed,
								AssignTo:    &controlWindow.showBoneFixedAction,
							},
							declarative.Action{
								Text:        mi18n.T("&回転ボーン"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowBoneRotate,
								AssignTo:    &controlWindow.showBoneRotateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&移動ボーン"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowBoneTranslate,
								AssignTo:    &controlWindow.showBoneTranslateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&表示ボーン"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowBoneVisible,
								AssignTo:    &controlWindow.showBoneVisibleAction,
							},
						},
					},
					declarative.Separator{},
					declarative.Menu{
						Text: mi18n.T("&剛体表示"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&前面表示"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowRigidBodyFront,
								AssignTo:    &controlWindow.showRigidBodyFrontAction,
							},
							declarative.Action{
								Text:        mi18n.T("&埋め込み表示"),
								Checkable:   true,
								OnTriggered: controlWindow.onTriggerShowRigidBodyBack,
								AssignTo:    &controlWindow.showRigidBodyBackAction,
							},
						},
					},
					declarative.Action{
						Text:        mi18n.T("&ジョイント表示"),
						Checkable:   true,
						OnTriggered: controlWindow.onTriggerShowJoint,
						AssignTo:    &controlWindow.showJointAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text: mi18n.T("&使い方"),
						OnTriggered: func() {
							mlog.ILT(mi18n.T("ビューワーの使い方"), mi18n.T("ビューワーの使い方メッセージ"))
						},
					},
				},
			},
			declarative.Menu{
				Text:  mi18n.T("&操作画面"),
				Items: logMenuItems,
			},
			declarative.Menu{
				Text:  mi18n.T("&使い方"),
				Items: helpMenuItemsFunc(),
			},
			declarative.Menu{
				Text: mi18n.T("&言語"),
				Items: []declarative.MenuItem{
					declarative.Action{
						Text:        "日本語",
						OnTriggered: func() { controlWindow.onChangeLanguage("ja") },
					},
					declarative.Action{
						Text:        "English",
						OnTriggered: func() { controlWindow.onChangeLanguage("en") },
					},
					declarative.Action{
						Text:        "中文",
						OnTriggered: func() { controlWindow.onChangeLanguage("zh") },
					},
					declarative.Action{
						Text:        "한국어",
						OnTriggered: func() { controlWindow.onChangeLanguage("ko") },
					},
				},
			},
		},
	}).Create(); err != nil {
		widget.RaiseError(err)
	}

	controlWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if !controlWindow.appState.IsClosed() {
			if result := walk.MsgBox(nil, mi18n.T("終了確認"), mi18n.T("終了確認メッセージ"),
				walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel); result == walk.DlgCmdOK {
				controlWindow.SetClosed(true)
			} else {
				// 閉じない場合はキャンセル
				*canceled = true
			}
		}
	})

	controlWindow.SetIcon(appConfig.Icon)

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DShadow)
	if err != nil {
		widget.RaiseError(err)
	}
	controlWindow.SetBackground(bg)

	// 初期設定
	controlWindow.limitFps30Action.SetChecked(true)       // 物理ON
	controlWindow.enabledPhysicsAction.SetChecked(true)   // フレームドロップON
	controlWindow.enabledFrameDropAction.SetChecked(true) // 30fps制限

	return controlWindow
}

func (controlWindow *ControlWindow) Dispose() {
	controlWindow.Close()
}

func (controlWindow *ControlWindow) Close() {
	controlWindow.MainWindow.Close()
}

func (controlWindow *ControlWindow) Run() {
	controlWindow.MainWindow.Run()
}

func (controlWindow *ControlWindow) Size() (int, int) {
	size := controlWindow.MainWindow.Size()
	return size.Width, size.Height
}

func (controlWindow *ControlWindow) SetPosition(x, y int) {
	controlWindow.MainWindow.SetX(x)
	controlWindow.MainWindow.SetY(y)
}

func (controlWindow *ControlWindow) GetMainWindow() *walk.MainWindow {
	return controlWindow.MainWindow
}

func (controlWindow *ControlWindow) InitTabWidget() {
	controlWindow.tabWidget = widget.NewMTabWidget(controlWindow.MainWindow)
}

func (controlWindow *ControlWindow) AddTabPage(tabPage *walk.TabPage) {
	err := controlWindow.tabWidget.Pages().Add(tabPage)
	if err != nil {
		widget.RaiseError(err)
	}
}

func (controlWindow *ControlWindow) SetPlayer(player app.IPlayer) {
	controlWindow.controlState.SetPlayer(player)
}

func (controlWindow *ControlWindow) onChangeLanguage(lang string) {
	if result := walk.MsgBox(
		controlWindow.MainWindow,
		mi18n.TWithLocale(lang, "言語変更"),
		mi18n.TWithLocale(lang, "言語変更メッセージ"),
		walk.MsgBoxOKCancel|walk.MsgBoxIconInformation,
	); result == walk.DlgCmdOK {
		mi18n.SetLang(lang)
		controlWindow.controlState.SetClosed(true)
	}
}

func (controlWindow *ControlWindow) logLevelTriggered() {
	mlog.SetLevel(mlog.INFO)
	if controlWindow.logLevelDebugAction.Checked() {
		mlog.SetLevel(mlog.DEBUG)
	}
	if controlWindow.logLevelIkVerboseAction.Checked() {
		mlog.SetLevel(mlog.IK_VERBOSE)
	}
	if controlWindow.logLevelVerboseAction.Checked() {
		mlog.SetLevel(mlog.VERBOSE)
	}
}

func (controlWindow *ControlWindow) onTriggerEnabledFrameDrop() {
	controlWindow.controlState.SetEnabledFrameDrop(controlWindow.enabledFrameDropAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerEnabledPhysics() {
	controlWindow.controlState.SetEnabledPhysics(controlWindow.enabledPhysicsAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerPhysicsReset() {
	controlWindow.controlState.SetPhysicsReset(true)
}

func (controlWindow *ControlWindow) onTriggerShowNormal() {
	controlWindow.controlState.SetShowNormal(controlWindow.showNormalAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowWire() {
	controlWindow.controlState.SetShowWire(controlWindow.showWireAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowSelectedVertex() {
	controlWindow.controlState.SetShowSelectedVertex(controlWindow.showSelectedVertexAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowBoneAll() {
	if controlWindow.showBoneAllAction.Checked() {
		controlWindow.showBoneIkAction.SetChecked(false)
		controlWindow.showBoneEffectorAction.SetChecked(false)
		controlWindow.showBoneFixedAction.SetChecked(false)
		controlWindow.showBoneRotateAction.SetChecked(false)
		controlWindow.showBoneTranslateAction.SetChecked(false)
		controlWindow.showBoneVisibleAction.SetChecked(false)

		controlWindow.controlState.SetShowBoneIk(false)
		controlWindow.controlState.SetShowBoneEffector(false)
		controlWindow.controlState.SetShowBoneFixed(false)
		controlWindow.controlState.SetShowBoneRotate(false)
		controlWindow.controlState.SetShowBoneTranslate(false)
		controlWindow.controlState.SetShowBoneVisible(false)
	}
	controlWindow.controlState.SetShowBoneAll(controlWindow.showBoneAllAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowBoneIk() {
	if controlWindow.showBoneIkAction.Checked() {
		controlWindow.showBoneAllAction.SetChecked(false)
		controlWindow.controlState.SetShowBoneAll(false)
	}
	controlWindow.controlState.SetShowBoneIk(controlWindow.showBoneIkAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowBoneEffector() {
	if controlWindow.showBoneEffectorAction.Checked() {
		controlWindow.showBoneAllAction.SetChecked(false)
		controlWindow.controlState.SetShowBoneAll(false)
	}
	controlWindow.controlState.SetShowBoneEffector(controlWindow.showBoneEffectorAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowBoneFixed() {
	if controlWindow.showBoneFixedAction.Checked() {
		controlWindow.showBoneAllAction.SetChecked(false)
		controlWindow.controlState.SetShowBoneAll(false)
	}
	controlWindow.controlState.SetShowBoneFixed(controlWindow.showBoneFixedAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowBoneRotate() {
	if controlWindow.showBoneRotateAction.Checked() {
		controlWindow.showBoneAllAction.SetChecked(false)
		controlWindow.controlState.SetShowBoneAll(false)
	}
	controlWindow.controlState.SetShowBoneRotate(controlWindow.showBoneRotateAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowBoneTranslate() {
	if controlWindow.showBoneTranslateAction.Checked() {
		controlWindow.showBoneAllAction.SetChecked(false)
		controlWindow.controlState.SetShowBoneAll(false)
	}
	controlWindow.controlState.SetShowBoneTranslate(controlWindow.showBoneTranslateAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowBoneVisible() {
	if controlWindow.showBoneVisibleAction.Checked() {
		controlWindow.showBoneAllAction.SetChecked(false)
		controlWindow.controlState.SetShowBoneAll(false)
	}
	controlWindow.controlState.SetShowBoneVisible(controlWindow.showBoneVisibleAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowRigidBodyFront() {
	controlWindow.controlState.SetShowRigidBodyFront(controlWindow.showRigidBodyFrontAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowRigidBodyBack() {
	controlWindow.controlState.SetShowRigidBodyBack(controlWindow.showRigidBodyBackAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowJoint() {
	controlWindow.controlState.SetShowJoint(controlWindow.showJointAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerShowInfo() {
	controlWindow.controlState.SetShowInfo(controlWindow.showInfoAction.Checked())
}

func (controlWindow *ControlWindow) onTriggerFps30Limit() {
	controlWindow.limitFps30Action.SetChecked(true)
	controlWindow.limitFps60Action.SetChecked(false)
	controlWindow.limitFpsUnLimitAction.SetChecked(false)
	controlWindow.SetSpfLimit(1 / 30.0)
	controlWindow.controlState.SetSpfLimit(controlWindow.SpfLimit())
}

func (controlWindow *ControlWindow) onTriggerFps60Limit() {
	controlWindow.limitFps30Action.SetChecked(false)
	controlWindow.limitFps60Action.SetChecked(true)
	controlWindow.limitFpsUnLimitAction.SetChecked(false)
	controlWindow.SetSpfLimit(1 / 60.0)
	controlWindow.controlState.SetSpfLimit(controlWindow.SpfLimit())
}

func (controlWindow *ControlWindow) onTriggerUnLimitFps() {
	controlWindow.limitFps30Action.SetChecked(false)
	controlWindow.limitFps60Action.SetChecked(false)
	controlWindow.limitFpsUnLimitAction.SetChecked(true)
	controlWindow.SetSpfLimit(-1.0)
	controlWindow.controlState.SetSpfLimit(controlWindow.SpfLimit())
}

func (controlWindow *ControlWindow) Frame() float64 {
	return controlWindow.controlState.motionPlayer.Frame()
}

func (controlWindow *ControlWindow) SetFrame(frame float64) {
	controlWindow.controlState.SetFrame(frame)
}

func (controlWindow *ControlWindow) AddFrame(v float64) {
	f := controlWindow.Frame() + v
	controlWindow.controlState.SetFrame(f)
}

func (controlWindow *ControlWindow) MaxFrame() int {
	return controlWindow.controlState.motionPlayer.MaxFrame()
}

func (controlWindow *ControlWindow) UpdateMaxFrame(maxFrame int) {
	if controlWindow.MaxFrame() < maxFrame {
		controlWindow.controlState.SetMaxFrame(maxFrame)
	}
}

func (controlWindow *ControlWindow) SetMaxFrame(maxFrame int) {
	controlWindow.controlState.SetMaxFrame(maxFrame)
}

func (controlWindow *ControlWindow) PrevFrame() int {
	return controlWindow.controlState.motionPlayer.PrevFrame()
}

func (controlWindow *ControlWindow) SetPrevFrame(prevFrame int) {
	controlWindow.controlState.SetPrevFrame(prevFrame)
}

func (controlWindow *ControlWindow) SetAnimationState(state state.IAnimationState) {
	controlWindow.controlState.SetAnimationState(state)
}

func (controlWindow *ControlWindow) IsEnabledFrameDrop() bool {
	return controlWindow.enabledFrameDropAction.Checked()
}

func (controlWindow *ControlWindow) SetEnabledFrameDrop(enabled bool) {
	controlWindow.enabledFrameDropAction.SetChecked(enabled)
}

func (controlWindow *ControlWindow) IsEnabledPhysics() bool {
	return controlWindow.enabledPhysicsAction.Checked()
}

func (controlWindow *ControlWindow) SetEnabledPhysics(enabled bool) {
	controlWindow.enabledPhysicsAction.SetChecked(enabled)
}

func (controlWindow *ControlWindow) IsPhysicsReset() bool {
	return controlWindow.physicsResetAction.Checked()
}

func (controlWindow *ControlWindow) SetPhysicsReset(reset bool) {
	controlWindow.physicsResetAction.SetChecked(reset)
}

func (controlWindow *ControlWindow) IsShowNormal() bool {
	return controlWindow.showNormalAction.Checked()
}

func (controlWindow *ControlWindow) SetShowNormal(show bool) {
	controlWindow.showNormalAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowWire() bool {
	return controlWindow.showWireAction.Checked()
}

func (controlWindow *ControlWindow) SetShowWire(show bool) {
	controlWindow.showWireAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowSelectedVertex() bool {
	return controlWindow.showSelectedVertexAction.Checked()
}

func (controlWindow *ControlWindow) SetShowSelectedVertex(show bool) {
	controlWindow.showSelectedVertexAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowBoneAll() bool {
	return controlWindow.showBoneAllAction.Checked()
}

func (controlWindow *ControlWindow) SetShowBoneAll(show bool) {
	controlWindow.showBoneAllAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowBoneIk() bool {
	return controlWindow.showBoneIkAction.Checked()
}

func (controlWindow *ControlWindow) SetShowBoneIk(show bool) {
	controlWindow.showBoneIkAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowBoneEffector() bool {
	return controlWindow.showBoneEffectorAction.Checked()
}

func (controlWindow *ControlWindow) SetShowBoneEffector(show bool) {
	controlWindow.showBoneEffectorAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowBoneFixed() bool {
	return controlWindow.showBoneFixedAction.Checked()
}

func (controlWindow *ControlWindow) SetShowBoneFixed(show bool) {
	controlWindow.showBoneFixedAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowBoneRotate() bool {
	return controlWindow.showBoneRotateAction.Checked()
}

func (controlWindow *ControlWindow) SetShowBoneRotate(show bool) {
	controlWindow.showBoneRotateAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowBoneTranslate() bool {
	return controlWindow.showBoneTranslateAction.Checked()
}

func (controlWindow *ControlWindow) SetShowBoneTranslate(show bool) {
	controlWindow.showBoneTranslateAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowBoneVisible() bool {
	return controlWindow.showBoneVisibleAction.Checked()
}

func (controlWindow *ControlWindow) SetShowBoneVisible(show bool) {
	controlWindow.showBoneVisibleAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowRigidBodyFront() bool {
	return controlWindow.showRigidBodyFrontAction.Checked()
}

func (controlWindow *ControlWindow) SetShowRigidBodyFront(show bool) {
	controlWindow.showRigidBodyFrontAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowRigidBodyBack() bool {
	return controlWindow.showRigidBodyBackAction.Checked()
}

func (controlWindow *ControlWindow) SetShowRigidBodyBack(show bool) {
	controlWindow.showRigidBodyBackAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowJoint() bool {
	return controlWindow.showJointAction.Checked()
}

func (controlWindow *ControlWindow) SetShowJoint(show bool) {
	controlWindow.showJointAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsShowInfo() bool {
	return controlWindow.showInfoAction.Checked()
}

func (controlWindow *ControlWindow) SetShowInfo(show bool) {
	controlWindow.showInfoAction.SetChecked(show)
}

func (controlWindow *ControlWindow) IsLimitFps30() bool {
	return controlWindow.limitFps30Action.Checked()
}

func (controlWindow *ControlWindow) SetLimitFps30(limit bool) {
	controlWindow.limitFps30Action.SetChecked(limit)
}

func (controlWindow *ControlWindow) IsLimitFps60() bool {
	return controlWindow.limitFps60Action.Checked()
}

func (controlWindow *ControlWindow) SetLimitFps60(limit bool) {
	controlWindow.limitFps60Action.SetChecked(limit)
}

func (controlWindow *ControlWindow) IsUnLimitFps() bool {
	return controlWindow.limitFpsUnLimitAction.Checked()
}

func (controlWindow *ControlWindow) SetUnLimitFps(limit bool) {
	controlWindow.limitFpsUnLimitAction.SetChecked(limit)
}

func (controlWindow *ControlWindow) IsLogLevelDebug() bool {
	return controlWindow.logLevelDebugAction.Checked()
}

func (controlWindow *ControlWindow) SetLogLevelDebug(log bool) {
	controlWindow.logLevelDebugAction.SetChecked(log)
}

func (controlWindow *ControlWindow) IsLogLevelVerbose() bool {
	return controlWindow.logLevelVerboseAction.Checked()
}

func (controlWindow *ControlWindow) SetLogLevelVerbose(log bool) {
	controlWindow.logLevelVerboseAction.SetChecked(log)
}

func (controlWindow *ControlWindow) IsLogLevelIkVerbose() bool {
	return controlWindow.logLevelIkVerboseAction.Checked()
}

func (controlWindow *ControlWindow) SetLogLevelIkVerbose(log bool) {
	controlWindow.logLevelIkVerboseAction.SetChecked(log)
}

func (controlWindow *ControlWindow) IsClosed() bool {
	return false
}

func (controlWindow *ControlWindow) SetClosed(closed bool) {
	controlWindow.controlState.SetClosed(closed)
}

func (controlWindow *ControlWindow) Playing() bool {
	return controlWindow.controlState.motionPlayer != nil && controlWindow.controlState.motionPlayer.Playing()
}

func (controlWindow *ControlWindow) TriggerPlay(p bool) {
	controlWindow.controlState.TriggerPlay(p)
}

func (controlWindow *ControlWindow) SpfLimit() float64 {
	return controlWindow.spfLimit
}

func (controlWindow *ControlWindow) SetSpfLimit(spf float64) {
	controlWindow.spfLimit = spf
}

func (controlWindow *ControlWindow) SetEnabled(enabled bool) {
	if controlWindow.tabWidget != nil {
		for i := range controlWindow.tabWidget.Pages().Len() {
			for j := range controlWindow.tabWidget.Pages().At(i).Children().Len() {
				controlWindow.tabWidget.Pages().At(i).Children().At(j).SetEnabled(enabled)
			}
		}
		// controlWindow.tabWidget.SetEnabled(enabled)
	}
	if controlWindow.controlState.motionPlayer != nil {
		controlWindow.controlState.motionPlayer.SetEnabled(enabled)
	}
}

func (controlWindow *ControlWindow) Enabled() bool {
	if controlWindow.tabWidget != nil {
		for i := range controlWindow.tabWidget.Pages().Len() {
			if controlWindow.tabWidget.Pages().At(i) != nil && !controlWindow.tabWidget.Pages().At(i).Enabled() {
				return false
			}
		}
	}

	return true
}
