package controller

import (
	"fmt"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/window"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"

	"github.com/miu200521358/walk/pkg/declarative"
)

type ControlWindow struct {
	*walk.MainWindow
	appConfig                   *mconfig.AppConfig // アプリケーション設定
	uiState                     window.IUiState    // UI状態
	enabledFrameDropAction      *walk.Action       // フレームドロップON/OFF
	enabledPhysicsAction        *walk.Action       // 物理ON/OFF
	physicsResetAction          *walk.Action       // 物理リセット
	showNormalAction            *walk.Action       // ボーンデバッグ表示
	showWireAction              *walk.Action       // ワイヤーフレームデバッグ表示
	showSelectedVertexAction    *walk.Action       // 選択頂点デバッグ表示
	showBoneAllAction           *walk.Action       // 全ボーンデバッグ表示
	showBoneIkAction            *walk.Action       // IKボーンデバッグ表示
	showBoneEffectorAction      *walk.Action       // 付与親ボーンデバッグ表示
	showBoneFixedAction         *walk.Action       // 軸制限ボーンデバッグ表示
	showBoneRotateAction        *walk.Action       // 回転ボーンデバッグ表示
	showBoneTranslateAction     *walk.Action       // 移動ボーンデバッグ表示
	showBoneVisibleAction       *walk.Action       // 表示ボーンデバッグ表示
	showRigidBodyFrontAction    *walk.Action       // 剛体デバッグ表示(前面)
	showRigidBodyBackAction     *walk.Action       // 剛体デバッグ表示(埋め込み)
	showJointAction             *walk.Action       // ジョイントデバッグ表示
	showInfoAction              *walk.Action       // 情報デバッグ表示
	limitFps30Action            *walk.Action       // 30FPS制限
	limitFps60Action            *walk.Action       // 60FPS制限
	limitFpsUnLimitAction       *walk.Action       // FPS無制限
	limitFpsDeformUnLimitAction *walk.Action       // デフォームFPS無制限
	logLevelDebugAction         *walk.Action       // デバッグメッセージ表示
	logLevelVerboseAction       *walk.Action       // 冗長メッセージ表示
	logLevelIkVerboseAction     *walk.Action       // IK冗長メッセージ表示
}

func NewMWindow(
	appConfig *mconfig.AppConfig,
	uiState window.IUiState,
	helpMenuItemsFunc func() []declarative.MenuItem,
) *ControlWindow {
	mWindow := &ControlWindow{
		appConfig: appConfig,
		uiState:   uiState,
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
			OnTriggered: mWindow.logLevelTriggered,
			AssignTo:    &mWindow.logLevelDebugAction,
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
				OnTriggered: mWindow.logLevelTriggered,
				AssignTo:    &mWindow.logLevelVerboseAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&IK冗長ログ表示"),
				Checkable:   true,
				OnTriggered: mWindow.logLevelTriggered,
				AssignTo:    &mWindow.logLevelIkVerboseAction,
			})
	}

	fpsLImitMenuItems := []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("&30fps制限"),
			Checkable:   true,
			OnTriggered: mWindow.onTriggerFps30Limit,
			AssignTo:    &mWindow.limitFps30Action,
		},
		declarative.Action{
			Text:        mi18n.T("&60fps制限"),
			Checkable:   true,
			OnTriggered: mWindow.onTriggerFps60Limit,
			AssignTo:    &mWindow.limitFps60Action,
		},
		declarative.Action{
			Text:        mi18n.T("&fps無制限"),
			Checkable:   true,
			OnTriggered: mWindow.onTriggerUnLimitFps,
			AssignTo:    &mWindow.limitFpsUnLimitAction,
		},
	}

	if !appConfig.IsEnvProd() {
		// 開発時にだけ描画無制限モードを追加
		fpsLImitMenuItems = append(fpsLImitMenuItems,
			declarative.Action{
				Text:        "&デフォームfps無制限",
				Checkable:   true,
				OnTriggered: mWindow.onTriggerUnLimitFpsDeform,
				AssignTo:    &mWindow.limitFpsDeformUnLimitAction,
			})
	}

	if err := (declarative.MainWindow{
		AssignTo: &mWindow.MainWindow,
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
						OnTriggered: mWindow.onTriggerEnabledFrameDrop,
						AssignTo:    &mWindow.enabledFrameDropAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&物理ON/OFF"),
						Checkable:   true,
						OnTriggered: mWindow.onTriggerEnabledPhysics,
						AssignTo:    &mWindow.enabledPhysicsAction,
					},
					declarative.Action{
						Text:        mi18n.T("&物理リセット"),
						OnTriggered: mWindow.onTriggerPhysicsReset,
						AssignTo:    &mWindow.physicsResetAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&法線表示"),
						Checkable:   true,
						OnTriggered: mWindow.onTriggerShowNormal,
						AssignTo:    &mWindow.showNormalAction,
					},
					declarative.Action{
						Text:        mi18n.T("&ワイヤーフレーム表示"),
						Checkable:   true,
						OnTriggered: mWindow.onTriggerShowWire,
						AssignTo:    &mWindow.showWireAction,
					},
					declarative.Action{
						Text:        mi18n.T("&選択頂点表示"),
						Checkable:   true,
						OnTriggered: mWindow.onTriggerShowSelectedVertex,
						AssignTo:    &mWindow.showSelectedVertexAction,
					},
					declarative.Separator{},
					declarative.Menu{
						Text: mi18n.T("&ボーン表示"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&全ボーン"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowBoneAll,
								AssignTo:    &mWindow.showBoneAllAction,
							},
							declarative.Separator{},
							declarative.Action{
								Text:        mi18n.T("&IKボーン"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowBoneIk,
								AssignTo:    &mWindow.showBoneIkAction,
							},
							declarative.Action{
								Text:        mi18n.T("&付与親ボーン"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowBoneEffector,
								AssignTo:    &mWindow.showBoneEffectorAction,
							},
							declarative.Action{
								Text:        mi18n.T("&軸制限ボーン"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowBoneFixed,
								AssignTo:    &mWindow.showBoneFixedAction,
							},
							declarative.Action{
								Text:        mi18n.T("&回転ボーン"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowBoneRotate,
								AssignTo:    &mWindow.showBoneRotateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&移動ボーン"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowBoneTranslate,
								AssignTo:    &mWindow.showBoneTranslateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&表示ボーン"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowBoneVisible,
								AssignTo:    &mWindow.showBoneVisibleAction,
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
								OnTriggered: mWindow.onTriggerShowRigidBodyFront,
								AssignTo:    &mWindow.showRigidBodyFrontAction,
							},
							declarative.Action{
								Text:        mi18n.T("&埋め込み表示"),
								Checkable:   true,
								OnTriggered: mWindow.onTriggerShowRigidBodyBack,
								AssignTo:    &mWindow.showRigidBodyBackAction,
							},
						},
					},
					declarative.Action{
						Text:        mi18n.T("&ジョイント表示"),
						Checkable:   true,
						OnTriggered: mWindow.onTriggerShowJoint,
						AssignTo:    &mWindow.showJointAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&情報表示"),
						Checkable:   true,
						OnTriggered: mWindow.onTriggerShowInfo,
						AssignTo:    &mWindow.showInfoAction,
					},
					declarative.Menu{
						Text:  mi18n.T("&fps制限"),
						Items: fpsLImitMenuItems,
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
				Text:  mi18n.T("&メイン画面"),
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
						OnTriggered: func() { mWindow.langTriggered("ja") },
					},
					declarative.Action{
						Text:        "English",
						OnTriggered: func() { mWindow.langTriggered("en") },
					},
					declarative.Action{
						Text:        "中文",
						OnTriggered: func() { mWindow.langTriggered("zh") },
					},
					declarative.Action{
						Text:        "???",
						OnTriggered: func() { mWindow.langTriggered("ko") },
					},
				},
			},
		},
	}).Create(); err != nil {
		return nil
	}

	mWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		mWindow.uiState.SetClosed(true)
	})

	icon, err := walk.NewIconFromImageForDPI(*appConfig.IconImage, 96)
	if err != nil {
		return nil
	}
	mWindow.SetIcon(icon)

	// // タブウィジェット追加
	// mWindow.TabWidget = widget.NewMTabWidget(mWindow.MainWindow)
	// mWindow.Children().Add(mWindow.TabWidget)

	// bg, err := walk.NewSystemColorBrush(walk.SysColor3DShadow)
	// widget.CheckError(err, mWindow.MainWindow, mi18n.T("背景色生成エラー"))
	// mWindow.SetBackground(bg)

	return mWindow
}

func (w *ControlWindow) Dispose() {
	w.Close()
}

func (w *ControlWindow) Close() {
	w.MainWindow.Close()
	w.uiState.SetClosed(true)
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

func (w *ControlWindow) langTriggered(lang string) {
	mi18n.SetLang(lang)
	walk.MsgBox(
		w.MainWindow,
		mi18n.TWithLocale(lang, "LanguageChanged.Title"),
		mi18n.TWithLocale(lang, "LanguageChanged.Message"),
		walk.MsgBoxOK|walk.MsgBoxIconInformation,
	)
	w.uiState.SetClosed(true)
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
	w.uiState.SetEnabledFrameDrop(w.enabledFrameDropAction.Checked())
}

func (w *ControlWindow) onTriggerEnabledPhysics() {
	w.uiState.SetEnabledPhysics(w.enabledPhysicsAction.Checked())
}

func (w *ControlWindow) onTriggerPhysicsReset() {
	w.uiState.SetPhysicsReset(true)
}

func (w *ControlWindow) onTriggerShowNormal() {
	w.uiState.SetShowNormal(w.showNormalAction.Checked())
}

func (w *ControlWindow) onTriggerShowWire() {
	w.uiState.SetShowWire(w.showWireAction.Checked())
}

func (w *ControlWindow) onTriggerShowSelectedVertex() {
	w.uiState.SetShowSelectedVertex(w.showSelectedVertexAction.Checked())
}

func (w *ControlWindow) onTriggerShowBoneAll() {
	w.uiState.SetShowBoneAll(true)
	w.uiState.SetShowBoneIk(false)
	w.uiState.SetShowBoneEffector(false)
	w.uiState.SetShowBoneFixed(false)
	w.uiState.SetShowBoneRotate(false)
	w.uiState.SetShowBoneTranslate(false)
	w.uiState.SetShowBoneVisible(false)

	w.showBoneAllAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneIk() {
	w.uiState.SetShowBoneAll(false)
	w.uiState.SetShowBoneIk(true)
	w.uiState.SetShowBoneEffector(false)
	w.uiState.SetShowBoneFixed(false)
	w.uiState.SetShowBoneRotate(false)
	w.uiState.SetShowBoneTranslate(false)
	w.uiState.SetShowBoneVisible(false)

	w.showBoneIkAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneEffector() {
	w.uiState.SetShowBoneAll(false)
	w.uiState.SetShowBoneIk(false)
	w.uiState.SetShowBoneEffector(true)
	w.uiState.SetShowBoneFixed(false)
	w.uiState.SetShowBoneRotate(false)
	w.uiState.SetShowBoneTranslate(false)
	w.uiState.SetShowBoneVisible(false)

	w.showBoneEffectorAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneFixed() {
	w.uiState.SetShowBoneAll(false)
	w.uiState.SetShowBoneIk(false)
	w.uiState.SetShowBoneEffector(false)
	w.uiState.SetShowBoneFixed(true)
	w.uiState.SetShowBoneRotate(false)
	w.uiState.SetShowBoneTranslate(false)
	w.uiState.SetShowBoneVisible(false)

	w.showBoneFixedAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneRotate() {
	w.uiState.SetShowBoneAll(false)
	w.uiState.SetShowBoneIk(false)
	w.uiState.SetShowBoneEffector(false)
	w.uiState.SetShowBoneFixed(false)
	w.uiState.SetShowBoneRotate(true)
	w.uiState.SetShowBoneTranslate(false)
	w.uiState.SetShowBoneVisible(false)

	w.showBoneRotateAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneTranslate() {
	w.uiState.SetShowBoneAll(false)
	w.uiState.SetShowBoneIk(false)
	w.uiState.SetShowBoneEffector(false)
	w.uiState.SetShowBoneFixed(false)
	w.uiState.SetShowBoneRotate(false)
	w.uiState.SetShowBoneTranslate(true)
	w.uiState.SetShowBoneVisible(false)

	w.showBoneTranslateAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowBoneVisible() {
	w.uiState.SetShowBoneAll(false)
	w.uiState.SetShowBoneIk(false)
	w.uiState.SetShowBoneEffector(false)
	w.uiState.SetShowBoneFixed(false)
	w.uiState.SetShowBoneRotate(false)
	w.uiState.SetShowBoneTranslate(false)
	w.uiState.SetShowBoneVisible(true)

	w.showBoneVisibleAction.SetChecked(true)
}

func (w *ControlWindow) onTriggerShowRigidBodyFront() {
	w.uiState.SetShowRigidBodyFront(w.showRigidBodyFrontAction.Checked())
}

func (w *ControlWindow) onTriggerShowRigidBodyBack() {
	w.uiState.SetShowRigidBodyBack(w.showRigidBodyBackAction.Checked())
}

func (w *ControlWindow) onTriggerShowJoint() {
	w.uiState.SetShowJoint(w.showJointAction.Checked())
}

func (w *ControlWindow) onTriggerShowInfo() {
	w.uiState.SetShowInfo(w.showInfoAction.Checked())
}

func (w *ControlWindow) onTriggerFps30Limit() {
	w.limitFps30Action.SetChecked(true)
	w.limitFps60Action.SetChecked(false)
	w.limitFpsUnLimitAction.SetChecked(false)
	w.limitFpsDeformUnLimitAction.SetChecked(false)
	w.uiState.SetSpfLimit(1 / 30.0)
}

func (w *ControlWindow) onTriggerFps60Limit() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(true)
	w.limitFpsUnLimitAction.SetChecked(false)
	w.limitFpsDeformUnLimitAction.SetChecked(false)
	w.uiState.SetSpfLimit(1 / 60.0)
}

func (w *ControlWindow) onTriggerUnLimitFps() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(false)
	w.limitFpsUnLimitAction.SetChecked(true)
	w.limitFpsDeformUnLimitAction.SetChecked(false)
	w.uiState.SetSpfLimit(-1.0)
}

func (w *ControlWindow) onTriggerUnLimitFpsDeform() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(false)
	w.limitFpsUnLimitAction.SetChecked(false)
	w.limitFpsDeformUnLimitAction.SetChecked(true)
	w.uiState.SetSpfLimit(-2.0)
}
