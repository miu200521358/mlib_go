package controller

import (
	"fmt"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/window"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"

	"github.com/miu200521358/walk/pkg/declarative"
)

type ControlWindow struct {
	*walk.MainWindow
	tabWidget                   *widget.MTabWidget            // タブウィジェット
	motionPlayer                *widget.MotionPlayer          // モーションプレイヤー
	appConfig                   *mconfig.AppConfig            // アプリケーション設定
	appState                    window.IAppState              // UI状態
	enabledFrameDropAction      *walk.Action                  // フレームドロップON/OFF
	enabledPhysicsAction        *walk.Action                  // 物理ON/OFF
	physicsResetAction          *walk.Action                  // 物理リセット
	showNormalAction            *walk.Action                  // ボーンデバッグ表示
	showWireAction              *walk.Action                  // ワイヤーフレームデバッグ表示
	showSelectedVertexAction    *walk.Action                  // 選択頂点デバッグ表示
	showBoneAllAction           *walk.Action                  // 全ボーンデバッグ表示
	showBoneIkAction            *walk.Action                  // IKボーンデバッグ表示
	showBoneEffectorAction      *walk.Action                  // 付与親ボーンデバッグ表示
	showBoneFixedAction         *walk.Action                  // 軸制限ボーンデバッグ表示
	showBoneRotateAction        *walk.Action                  // 回転ボーンデバッグ表示
	showBoneTranslateAction     *walk.Action                  // 移動ボーンデバッグ表示
	showBoneVisibleAction       *walk.Action                  // 表示ボーンデバッグ表示
	showRigidBodyFrontAction    *walk.Action                  // 剛体デバッグ表示(前面)
	showRigidBodyBackAction     *walk.Action                  // 剛体デバッグ表示(埋め込み)
	showJointAction             *walk.Action                  // ジョイントデバッグ表示
	showInfoAction              *walk.Action                  // 情報デバッグ表示
	limitFps30Action            *walk.Action                  // 30FPS制限
	limitFps60Action            *walk.Action                  // 60FPS制限
	limitFpsUnLimitAction       *walk.Action                  // FPS無制限
	limitFpsDeformUnLimitAction *walk.Action                  // デフォームFPS無制限
	logLevelDebugAction         *walk.Action                  // デバッグメッセージ表示
	logLevelVerboseAction       *walk.Action                  // 冗長メッセージ表示
	logLevelIkVerboseAction     *walk.Action                  // IK冗長メッセージ表示
	animationStates             [][]*renderer.AnimationStates // アニメーションステート
}

func NewControlWindow(
	appConfig *mconfig.AppConfig,
	appState window.IAppState,
	helpMenuItemsFunc func() []declarative.MenuItem,
) *ControlWindow {
	controlWindow := &ControlWindow{
		appConfig: appConfig,
		appState:  appState,
	}

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

	if !appConfig.IsEnvProd() {
		// 開発時にだけ描画無制限モードを追加
		fpsLImitMenuItems = append(fpsLImitMenuItems,
			declarative.Action{
				Text:        "&デフォームfps無制限",
				Checkable:   true,
				OnTriggered: controlWindow.onTriggerUnLimitFpsDeform,
				AssignTo:    &controlWindow.limitFpsDeformUnLimitAction,
			})
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
						OnTriggered: func() { controlWindow.langTriggered("ja") },
					},
					declarative.Action{
						Text:        "English",
						OnTriggered: func() { controlWindow.langTriggered("en") },
					},
					declarative.Action{
						Text:        "中文",
						OnTriggered: func() { controlWindow.langTriggered("zh") },
					},
					declarative.Action{
						Text:        "한국어",
						OnTriggered: func() { controlWindow.langTriggered("ko") },
					},
				},
			},
		},
	}).Create(); err != nil {
		widget.RaiseError(err)
	}

	controlWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		controlWindow.appState.SetClosed(true)
	})

	icon, err := walk.NewIconFromImageForDPI(*appConfig.IconImage, 96)
	if err != nil {
		widget.RaiseError(err)
	}
	controlWindow.SetIcon(icon)

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DShadow)
	if err != nil {
		widget.RaiseError(err)
	}
	controlWindow.SetBackground(bg)

	return controlWindow
}

func (w *ControlWindow) Dispose() {
	w.Close()
}

func (w *ControlWindow) Close() {
	w.MainWindow.Close()
	w.appState.SetClosed(true)
}

func (w *ControlWindow) Run() {
	w.MainWindow.Run()
}

func (w *ControlWindow) Size() (int, int) {
	size := w.MainWindow.Size()
	return size.Width, size.Height
}

func (w *ControlWindow) SetPosition(x, y int) {
	w.MainWindow.SetX(x)
	w.MainWindow.SetY(y)
}

func (w *ControlWindow) AppState() window.IAppState {
	return w.appState
}

func (w *ControlWindow) GetMainWindow() *walk.MainWindow {
	return w.MainWindow
}

func (w *ControlWindow) InitTabWidget() {
	w.tabWidget = widget.NewMTabWidget(w.GetMainWindow())
}

func (w *ControlWindow) AddTabPage(tabPage *walk.TabPage) {
	err := w.tabWidget.Pages().Add(tabPage)
	if err != nil {
		widget.RaiseError(err)
	}
}

func (w *ControlWindow) SetPlayer(player *widget.MotionPlayer) {
	w.motionPlayer = player
}

func (w *ControlWindow) langTriggered(lang string) {
	mi18n.SetLang(lang)
	walk.MsgBox(
		w.MainWindow,
		mi18n.TWithLocale(lang, "LanguageChanged.Title"),
		mi18n.TWithLocale(lang, "LanguageChanged.Message"),
		walk.MsgBoxOK|walk.MsgBoxIconInformation,
	)
	w.appState.SetClosed(true)
}

func (w *ControlWindow) logLevelTriggered() {
	mlog.SetLevel(mlog.INFO)
	if w.logLevelDebugAction.Checked() {
		mlog.SetLevel(mlog.DEBUG)
	}
	if w.logLevelIkVerboseAction.Checked() {
		mlog.SetLevel(mlog.IK_VERBOSE)
	}
	if w.logLevelVerboseAction.Checked() {
		mlog.SetLevel(mlog.VERBOSE)
	}
}

func (w *ControlWindow) onTriggerEnabledFrameDrop() {
	w.appState.SetEnabledFrameDrop(w.enabledFrameDropAction.Checked())
}

func (w *ControlWindow) onTriggerEnabledPhysics() {
	w.appState.SetEnabledPhysics(w.enabledPhysicsAction.Checked())
}

func (w *ControlWindow) onTriggerPhysicsReset() {
	w.appState.SetPhysicsReset(true)
}

func (w *ControlWindow) onTriggerShowNormal() {
	w.appState.SetShowNormal(w.showNormalAction.Checked())
}

func (w *ControlWindow) onTriggerShowWire() {
	w.appState.SetShowWire(w.showWireAction.Checked())
}

func (w *ControlWindow) onTriggerShowSelectedVertex() {
	w.appState.SetShowSelectedVertex(w.showSelectedVertexAction.Checked())
}

func (w *ControlWindow) onTriggerShowBoneAll() {
	w.appState.SetShowBoneAll(true)
	w.appState.SetShowBoneIk(false)
	w.appState.SetShowBoneEffector(false)
	w.appState.SetShowBoneFixed(false)
	w.appState.SetShowBoneRotate(false)
	w.appState.SetShowBoneTranslate(false)
	w.appState.SetShowBoneVisible(false)

	w.showBoneAllAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneIk() {
	w.appState.SetShowBoneAll(false)
	w.appState.SetShowBoneIk(true)
	w.appState.SetShowBoneEffector(false)
	w.appState.SetShowBoneFixed(false)
	w.appState.SetShowBoneRotate(false)
	w.appState.SetShowBoneTranslate(false)
	w.appState.SetShowBoneVisible(false)

	w.showBoneIkAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneEffector() {
	w.appState.SetShowBoneAll(false)
	w.appState.SetShowBoneIk(false)
	w.appState.SetShowBoneEffector(true)
	w.appState.SetShowBoneFixed(false)
	w.appState.SetShowBoneRotate(false)
	w.appState.SetShowBoneTranslate(false)
	w.appState.SetShowBoneVisible(false)

	w.showBoneEffectorAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneFixed() {
	w.appState.SetShowBoneAll(false)
	w.appState.SetShowBoneIk(false)
	w.appState.SetShowBoneEffector(false)
	w.appState.SetShowBoneFixed(true)
	w.appState.SetShowBoneRotate(false)
	w.appState.SetShowBoneTranslate(false)
	w.appState.SetShowBoneVisible(false)

	w.showBoneFixedAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneRotate() {
	w.appState.SetShowBoneAll(false)
	w.appState.SetShowBoneIk(false)
	w.appState.SetShowBoneEffector(false)
	w.appState.SetShowBoneFixed(false)
	w.appState.SetShowBoneRotate(true)
	w.appState.SetShowBoneTranslate(false)
	w.appState.SetShowBoneVisible(false)

	w.showBoneRotateAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneTranslate() {
	w.appState.SetShowBoneAll(false)
	w.appState.SetShowBoneIk(false)
	w.appState.SetShowBoneEffector(false)
	w.appState.SetShowBoneFixed(false)
	w.appState.SetShowBoneRotate(false)
	w.appState.SetShowBoneTranslate(true)
	w.appState.SetShowBoneVisible(false)

	w.showBoneTranslateAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneVisible() {
	w.appState.SetShowBoneAll(false)
	w.appState.SetShowBoneIk(false)
	w.appState.SetShowBoneEffector(false)
	w.appState.SetShowBoneFixed(false)
	w.appState.SetShowBoneRotate(false)
	w.appState.SetShowBoneTranslate(false)
	w.appState.SetShowBoneVisible(true)

	w.showBoneVisibleAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowRigidBodyFront() {
	w.appState.SetShowRigidBodyFront(w.showRigidBodyFrontAction.Checked())
}

func (w *ControlWindow) onTriggerShowRigidBodyBack() {
	w.appState.SetShowRigidBodyBack(w.showRigidBodyBackAction.Checked())
}

func (w *ControlWindow) onTriggerShowJoint() {
	w.appState.SetShowJoint(w.showJointAction.Checked())
}

func (w *ControlWindow) onTriggerShowInfo() {
	w.appState.SetShowInfo(w.showInfoAction.Checked())
}

func (w *ControlWindow) onTriggerFps30Limit() {
	w.limitFps30Action.SetChecked(true)
	w.limitFps60Action.SetChecked(false)
	w.limitFpsUnLimitAction.SetChecked(false)
	w.limitFpsDeformUnLimitAction.SetChecked(false)
	w.appState.SetSpfLimit(1 / 30.0)
}

func (w *ControlWindow) onTriggerFps60Limit() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(true)
	w.limitFpsUnLimitAction.SetChecked(false)
	w.limitFpsDeformUnLimitAction.SetChecked(false)
	w.appState.SetSpfLimit(1 / 60.0)
}

func (w *ControlWindow) onTriggerUnLimitFps() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(false)
	w.limitFpsUnLimitAction.SetChecked(true)
	w.limitFpsDeformUnLimitAction.SetChecked(false)
	w.appState.SetSpfLimit(-1.0)
}

func (w *ControlWindow) onTriggerUnLimitFpsDeform() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(false)
	w.limitFpsUnLimitAction.SetChecked(false)
	w.limitFpsDeformUnLimitAction.SetChecked(true)
	w.appState.SetSpfLimit(-2.0)
}

func (w *ControlWindow) Frame() float64 {
	return w.motionPlayer.Frame()
}

func (w *ControlWindow) SetFrame(frame float64) {
	go func() {
		w.motionPlayer.SetFrame(frame)
	}()
}

func (w *ControlWindow) AddFrame(v float64) {
	go func() {
		w.motionPlayer.SetFrame(w.motionPlayer.Frame() + v)
	}()
}

func (w *ControlWindow) MaxFrame() int {
	return w.motionPlayer.MaxFrame()
}

func (w *ControlWindow) SetMaxFrame(maxFrame int) {
	go func() {
		w.motionPlayer.SetMaxFrame(maxFrame)
	}()
}

func (w *ControlWindow) PrevFrame() int {
	return w.motionPlayer.PrevFrame()
}

func (w *ControlWindow) SetPrevFrame(prevFrame int) {
	go func() {
		w.motionPlayer.SetPrevFrame(prevFrame)
	}()
}
