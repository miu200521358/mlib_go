package controller

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/sys/windows"
)

// ControlWindow は操作画面(コントローラウィンドウ)を管理する
type ControlWindow struct {
	mainWindow *walk.MainWindow

	// SharedState への参照
	shared *state.SharedState

	// UI要素 (メニューアクションなど)
	// enabledFrameDropAction      *walk.Action // フレームドロップON/OFF
	enabledPhysicsAction *walk.Action // 物理ON/OFF
	physicsResetAction   *walk.Action // 物理リセット
	// showNormalAction            *walk.Action // ボーンデバッグ表示
	// showWireAction              *walk.Action // ワイヤーフレームデバッグ表示
	// showOverrideAction          *walk.Action // オーバーライドデバッグ表示
	// showSelectedVertexAction    *walk.Action // 選択頂点デバッグ表示
	// showBoneAllAction           *walk.Action // 全ボーンデバッグ表示
	// showBoneIkAction            *walk.Action // IKボーンデバッグ表示
	// showBoneEffectorAction      *walk.Action // 付与親ボーンデバッグ表示
	// showBoneFixedAction         *walk.Action // 軸制限ボーンデバッグ表示
	// showBoneRotateAction        *walk.Action // 回転ボーンデバッグ表示
	// showBoneTranslateAction     *walk.Action // 移動ボーンデバッグ表示
	// showBoneVisibleAction       *walk.Action // 表示ボーンデバッグ表示
	// showRigidBodyFrontAction    *walk.Action // 剛体デバッグ表示(前面)
	// showRigidBodyBackAction     *walk.Action // 剛体デバッグ表示(埋め込み)
	// showJointAction             *walk.Action // ジョイントデバッグ表示
	// showInfoAction              *walk.Action // 情報デバッグ表示
	// limitFps30Action            *walk.Action // 30FPS制限
	// limitFps60Action            *walk.Action // 60FPS制限
	// limitFpsUnLimitAction       *walk.Action // FPS無制限
	// cameraSyncAction            *walk.Action // カメラ同期
	logLevelDebugAction         *walk.Action // デバッグメッセージ表示
	logLevelVerboseAction       *walk.Action // 冗長メッセージ表示
	logLevelIkVerboseAction     *walk.Action // IK冗長メッセージ表示
	logLevelViewerVerboseAction *walk.Action // ビューワー冗長メッセージ表示
}

// NewControlWindow はコントローラウィンドウを生成する
func NewControlWindow(
	shared *state.SharedState,
	appConfig *mconfig.AppConfig,
	helpMenuItemsFunc func() []declarative.MenuItem,
) (*ControlWindow, error) {
	controlWindow := &ControlWindow{
		shared: shared,
	}

	logMenuItems := []declarative.MenuItem{
		declarative.Action{
			Text: mi18n.T("&使い方"),
			OnTriggered: func() {
				mlog.ILT(mi18n.T("メイン画面の使い方"), "%s", mi18n.T("メイン画面の使い方メッセージ"))
			},
		},
	}

	if !appConfig.IsEnvProd() {
		// 開発時のみ冗長ログ表示を追加
		logMenuItems = append(logMenuItems,
			declarative.Separator{})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&デバッグログ表示"),
				Checkable:   true,
				OnTriggered: controlWindow.triggerLogLevel,
				AssignTo:    &controlWindow.logLevelDebugAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&冗長ログ表示"),
				Checkable:   true,
				OnTriggered: controlWindow.triggerLogLevel,
				AssignTo:    &controlWindow.logLevelVerboseAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&IK冗長ログ表示"),
				Checkable:   true,
				OnTriggered: controlWindow.triggerLogLevel,
				AssignTo:    &controlWindow.logLevelIkVerboseAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&ビューワー冗長ログ表示"),
				Checkable:   true,
				OnTriggered: controlWindow.triggerLogLevel,
				AssignTo:    &controlWindow.logLevelViewerVerboseAction,
			})
	}

	if err := (declarative.MainWindow{
		AssignTo: &controlWindow.mainWindow,
		Title:    fmt.Sprintf("%s %s", appConfig.Name, appConfig.Version),
		Size:     GetWindowSize(appConfig.ControlWindowSize.Width, appConfig.ControlWindowSize.Height),
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear, MarginsZero: true, SpacingZero: true},
		MenuItems: []declarative.MenuItem{
			declarative.Menu{
				Text: mi18n.T("&ビューワー"),
				Items: []declarative.MenuItem{
					// declarative.Action{
					// 	Text:        mi18n.T("&フレームドロップON"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerEnabledFrameDrop,
					// 	AssignTo:    &controlWindow.enabledFrameDropAction,
					// },
					// declarative.Menu{
					// 	Text: mi18n.T("&fps制限"),
					// 	Items: []declarative.MenuItem{
					// 		declarative.Action{
					// 			Text:        mi18n.T("&30fps制限"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerFps30Limit,
					// 			AssignTo:    &controlWindow.limitFps30Action,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&60fps制限"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerFps60Limit,
					// 			AssignTo:    &controlWindow.limitFps60Action,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&fps無制限"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerUnLimitFps,
					// 			AssignTo:    &controlWindow.limitFpsUnLimitAction,
					// 		},
					// 	},
					// },
					// declarative.Action{
					// 	Text:        mi18n.T("&情報表示"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerShowInfo,
					// 	AssignTo:    &controlWindow.showInfoAction,
					// },
					// declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&物理ON/OFF"),
						Checkable:   true,
						OnTriggered: controlWindow.TriggerEnabledPhysics,
						AssignTo:    &controlWindow.enabledPhysicsAction,
					},
					declarative.Action{
						Text:        mi18n.T("&物理リセット"),
						OnTriggered: controlWindow.TriggerPhysicsReset,
						AssignTo:    &controlWindow.physicsResetAction,
					},
					// declarative.Separator{},
					// declarative.Action{
					// 	Text:        mi18n.T("&法線表示"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerShowNormal,
					// 	AssignTo:    &controlWindow.showNormalAction,
					// },
					// declarative.Action{
					// 	Text:        mi18n.T("&ワイヤーフレーム表示"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerShowWire,
					// 	AssignTo:    &controlWindow.showWireAction,
					// },
					// declarative.Separator{},
					// declarative.Action{
					// 	Text:        mi18n.T("&頂点ライン選択"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerShowSelectedVertex,
					// 	AssignTo:    &controlWindow.showSelectedVertexAction,
					// },
					// declarative.Action{
					// 	Text: mi18n.T("&頂点ライン選択使い方"),
					// 	OnTriggered: func() {
					// 		mlog.ILT(mi18n.T("&頂点ライン選択使い方"), mi18n.T("頂点ライン選択使い方メッセージ"))
					// 	},
					// },
					// declarative.Separator{},
					// declarative.Action{
					// 	Text:        mi18n.T("&カメラ同期"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerCameraSync,
					// 	AssignTo:    &controlWindow.cameraSyncAction,
					// },
					// declarative.Action{
					// 	Text:        mi18n.T("&サブビューワーオーバーレイ"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerShowOverride,
					// 	AssignTo:    &controlWindow.showOverrideAction,
					// },
					// declarative.Action{
					// 	Text: mi18n.T("&サブビューワーオーバーレイの使い方"),
					// 	OnTriggered: func() {
					// 		mlog.ILT(mi18n.T("&サブビューワーオーバーレイの使い方"),
					// 			mi18n.T("サブビューワーオーバーレイの使い方メッセージ"))
					// 	},
					// },
					// declarative.Separator{},
					// declarative.Menu{
					// 	Text: mi18n.T("&ボーン表示"),
					// 	Items: []declarative.MenuItem{
					// 		declarative.Action{
					// 			Text:        mi18n.T("&全ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowBoneAll,
					// 			AssignTo:    &controlWindow.showBoneAllAction,
					// 		},
					// 		declarative.Separator{},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&IKボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowBoneIk,
					// 			AssignTo:    &controlWindow.showBoneIkAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&付与親ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowBoneEffector,
					// 			AssignTo:    &controlWindow.showBoneEffectorAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&軸制限ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowBoneFixed,
					// 			AssignTo:    &controlWindow.showBoneFixedAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&回転ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowBoneRotate,
					// 			AssignTo:    &controlWindow.showBoneRotateAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&移動ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowBoneTranslate,
					// 			AssignTo:    &controlWindow.showBoneTranslateAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&表示ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowBoneVisible,
					// 			AssignTo:    &controlWindow.showBoneVisibleAction,
					// 		},
					// 	},
					// },
					// declarative.Menu{
					// 	Text: mi18n.T("&剛体表示"),
					// 	Items: []declarative.MenuItem{
					// 		declarative.Action{
					// 			Text:        mi18n.T("&前面表示"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowRigidBodyFront,
					// 			AssignTo:    &controlWindow.showRigidBodyFrontAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&埋め込み表示"),
					// 			Checkable:   true,
					// 			OnTriggered: controlWindow.TriggerShowRigidBodyBack,
					// 			AssignTo:    &controlWindow.showRigidBodyBackAction,
					// 		},
					// 	},
					// },
					// declarative.Action{
					// 	Text:        mi18n.T("&ジョイント表示"),
					// 	Checkable:   true,
					// 	OnTriggered: controlWindow.TriggerShowJoint,
					// 	AssignTo:    &controlWindow.showJointAction,
					// },
					// declarative.Separator{},
					declarative.Action{
						Text: mi18n.T("&ビューワーの使い方"),
						OnTriggered: func() {
							mlog.ILT(mi18n.T("&ビューワーの使い方"), "%s", mi18n.T("ビューワーの使い方メッセージ"))
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
		return nil, err
	}

	controlWindow.mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		// controllerStateを読み取り
		controlWindow.shared.ReadViewerState(func(vs state.IViewerState) {
			// ビューワーがまだ閉じていない場合のみ、確認ダイアログを表示
			if !vs.IsClosed() {
				if result := walk.MsgBox(
					nil,
					mi18n.T("終了確認"),
					mi18n.T("終了確認メッセージ"),
					walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel,
				); result == walk.DlgCmdOK {
					// ユーザーがOKを選んだ場合、controllerState の isClosed を true にする
					controlWindow.shared.UpdateControllerState(func(cs state.IControllerState) {
						cs.SetClosed(true)
					})
				} else {
					// 閉じない場合はキャンセル
					*canceled = true
				}
			}
		})
	})

	controlWindow.mainWindow.SetIcon(appConfig.Icon)

	if bg, err := walk.NewSystemColorBrush(walk.SysColor3DShadow); err != nil {
		return nil, err
	} else {
		controlWindow.mainWindow.SetBackground(bg)
	}

	// 初期設定
	// controlWindow.limitFps30Action.SetChecked(true)       // 物理ON
	controlWindow.enabledPhysicsAction.SetChecked(true) // フレームドロップON
	// controlWindow.enabledFrameDropAction.SetChecked(true) // 30fps制限

	return controlWindow, nil
}

// OnClose はウィンドウを閉じるときの処理
func (cw *ControlWindow) OnClose() {
	// コントローラStateのisClosedをtrueにする
	cw.shared.UpdateControllerState(func(cs state.IControllerState) {
		cs.SetClosed(true)
	})
}

// Run はメインウィンドウを起動する
func (cw *ControlWindow) Run() {
	cw.mainWindow.Run()
}

func (controlWindow *ControlWindow) Dispose() {
	controlWindow.mainWindow.Close()
}

func (controlWindow *ControlWindow) WindowSize() (int, int) {
	size := controlWindow.mainWindow.Size()
	return size.Width, size.Height
}

func (controlWindow *ControlWindow) SetPosition(x, y int) {
	controlWindow.mainWindow.SetX(x)
	controlWindow.mainWindow.SetY(y)
}

func (controlWindow *ControlWindow) onChangeLanguage(lang string) {
	if result := walk.MsgBox(
		controlWindow.mainWindow,
		mi18n.TWithLocale(lang, "言語変更"),
		mi18n.TWithLocale(lang, "言語変更メッセージ"),
		walk.MsgBoxOKCancel|walk.MsgBoxIconInformation,
	); result == walk.DlgCmdOK {
		mi18n.SetLang(lang)
		controlWindow.shared.UpdateControllerState(func(cs state.IControllerState) {
			cs.SetClosed(true)
		})
	}
}

func (controlWindow *ControlWindow) triggerLogLevel() {
	mlog.SetLevel(mlog.INFO)
	if controlWindow.logLevelDebugAction.Checked() {
		mlog.SetLevel(mlog.DEBUG)
	}
	if controlWindow.logLevelViewerVerboseAction.Checked() {
		mlog.I("exe階層に「viewerPng」フォルダを作成し、画面描画中の連番pngを出力し続けます\n画面サイズ: 1920x1080、視野角: 40.0、カメラ位置: (0, 10, 45)、カメラ角度: (0, 0, 0) ")
		mlog.SetLevel(mlog.VIEWER_VERBOSE)
	}
	if controlWindow.logLevelIkVerboseAction.Checked() {
		mlog.SetLevel(mlog.IK_VERBOSE)
	}
	if controlWindow.logLevelVerboseAction.Checked() {
		mlog.SetLevel(mlog.VERBOSE)
	}
}

// ------- 以下、メニューから呼ばれるトリガーメソッド -------

func (controlWindow *ControlWindow) TriggerEnabledPhysics() {
	controlWindow.shared.UpdateControllerState(func(cs state.IControllerState) {
		cs.SetEnabledPhysics(controlWindow.enabledPhysicsAction.Checked())
	})
}

func (controlWindow *ControlWindow) TriggerPhysicsReset() {
	controlWindow.shared.UpdateControllerState(func(cs state.IControllerState) {
		cs.SetPhysicsReset(true)
	})
}

// ----------------------

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func getSystemMetrics(nIndex int) int {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(nIndex))
	return int(ret)
}

func GetWindowSize(width int, height int) declarative.Size {
	screenWidth := getSystemMetrics(SM_CXSCREEN)
	screenHeight := getSystemMetrics(SM_CYSCREEN)

	if width > screenWidth-50 {
		width = screenWidth - 50
	}
	if height > screenHeight-50 {
		height = screenHeight - 50
	}

	return declarative.Size{Width: width, Height: height}
}
