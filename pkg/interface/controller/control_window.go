//go:build windows
// +build windows

package controller

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// ControlWindow は操作画面(コントローラウィンドウ)を管理する
type ControlWindow struct {
	*walk.MainWindow

	shared    *state.SharedState // SharedState への参照
	appConfig *mconfig.AppConfig // アプリケーション設定

	tabWidget *walk.TabWidget // タブウィジェット

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
	limitFps30Action      *walk.Action // 30FPS制限
	limitFps60Action      *walk.Action // 60FPS制限
	limitFpsUnLimitAction *walk.Action // FPS無制限
	// cameraSyncAction            *walk.Action // カメラ同期
	logLevelDebugAction         *walk.Action // デバッグメッセージ表示
	logLevelVerboseAction       *walk.Action // 冗長メッセージ表示
	logLevelIkVerboseAction     *walk.Action // IK冗長メッセージ表示
	logLevelViewerVerboseAction *walk.Action // ビューワー冗長メッセージ表示
}

// Run はコントローラウィンドウを実行する
func NewControlWindow(
	shared *state.SharedState,
	appConfig *mconfig.AppConfig,
	helpMenuItems []declarative.MenuItem,
	tabPages []declarative.TabPage,
	width, height, positionX, positionY int,
) (*ControlWindow, error) {
	cw := &ControlWindow{
		shared:    shared,
		appConfig: appConfig,
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
				OnTriggered: cw.triggerLogLevel,
				AssignTo:    &cw.logLevelDebugAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&冗長ログ表示"),
				Checkable:   true,
				OnTriggered: cw.triggerLogLevel,
				AssignTo:    &cw.logLevelVerboseAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&IK冗長ログ表示"),
				Checkable:   true,
				OnTriggered: cw.triggerLogLevel,
				AssignTo:    &cw.logLevelIkVerboseAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&ビューワー冗長ログ表示"),
				Checkable:   true,
				OnTriggered: cw.triggerLogLevel,
				AssignTo:    &cw.logLevelViewerVerboseAction,
			})
	}

	if err := (declarative.MainWindow{
		AssignTo: &cw.MainWindow,
		Title:    fmt.Sprintf("%s %s", appConfig.Name, appConfig.Version),
		Size:     declarative.Size{Width: width, Height: height},
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear, MarginsZero: true, SpacingZero: true},
		Background: declarative.SystemColorBrush{
			Color: walk.SysColor3DShadow,
		},
		Icon: appConfig.Icon,
		MenuItems: []declarative.MenuItem{
			declarative.Menu{
				Text: mi18n.T("&ビューワー"),
				Items: []declarative.MenuItem{
					// declarative.Action{
					// 	Text:        mi18n.T("&フレームドロップON"),
					// 	Checkable:   true,
					// 	OnTriggered: cw.TriggerEnabledFrameDrop,
					// 	AssignTo:    &cw.enabledFrameDropAction,
					// },
					declarative.Menu{
						Text: mi18n.T("&fps制限"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&30fps制限"),
								Checkable:   true,
								OnTriggered: cw.TriggerFps30Limit,
								AssignTo:    &cw.limitFps30Action,
							},
							declarative.Action{
								Text:        mi18n.T("&60fps制限"),
								Checkable:   true,
								OnTriggered: cw.TriggerFps60Limit,
								AssignTo:    &cw.limitFps60Action,
							},
							declarative.Action{
								Text:        mi18n.T("&fps無制限"),
								Checkable:   true,
								OnTriggered: cw.TriggerUnLimitFps,
								AssignTo:    &cw.limitFpsUnLimitAction,
							},
						},
					},
					// declarative.Action{
					// 	Text:        mi18n.T("&情報表示"),
					// 	Checkable:   true,
					// 	OnTriggered: cw.TriggerShowInfo,
					// 	AssignTo:    &cw.showInfoAction,
					// },
					// declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&物理ON/OFF"),
						Checkable:   true,
						OnTriggered: cw.TriggerEnabledPhysics,
						AssignTo:    &cw.enabledPhysicsAction,
					},
					declarative.Action{
						Text:        mi18n.T("&物理リセット"),
						OnTriggered: cw.TriggerPhysicsReset,
						AssignTo:    &cw.physicsResetAction,
					},
					// declarative.Separator{},
					// declarative.Action{
					// 	Text:        mi18n.T("&法線表示"),
					// 	Checkable:   true,
					// 	OnTriggered: cw.TriggerShowNormal,
					// 	AssignTo:    &cw.showNormalAction,
					// },
					// declarative.Action{
					// 	Text:        mi18n.T("&ワイヤーフレーム表示"),
					// 	Checkable:   true,
					// 	OnTriggered: cw.TriggerShowWire,
					// 	AssignTo:    &cw.showWireAction,
					// },
					// declarative.Separator{},
					// declarative.Action{
					// 	Text:        mi18n.T("&頂点ライン選択"),
					// 	Checkable:   true,
					// 	OnTriggered: cw.TriggerShowSelectedVertex,
					// 	AssignTo:    &cw.showSelectedVertexAction,
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
					// 	OnTriggered: cw.TriggerCameraSync,
					// 	AssignTo:    &cw.cameraSyncAction,
					// },
					// declarative.Action{
					// 	Text:        mi18n.T("&サブビューワーオーバーレイ"),
					// 	Checkable:   true,
					// 	OnTriggered: cw.TriggerShowOverride,
					// 	AssignTo:    &cw.showOverrideAction,
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
					// 			OnTriggered: cw.TriggerShowBoneAll,
					// 			AssignTo:    &cw.showBoneAllAction,
					// 		},
					// 		declarative.Separator{},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&IKボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowBoneIk,
					// 			AssignTo:    &cw.showBoneIkAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&付与親ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowBoneEffector,
					// 			AssignTo:    &cw.showBoneEffectorAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&軸制限ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowBoneFixed,
					// 			AssignTo:    &cw.showBoneFixedAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&回転ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowBoneRotate,
					// 			AssignTo:    &cw.showBoneRotateAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&移動ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowBoneTranslate,
					// 			AssignTo:    &cw.showBoneTranslateAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&表示ボーン"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowBoneVisible,
					// 			AssignTo:    &cw.showBoneVisibleAction,
					// 		},
					// 	},
					// },
					// declarative.Menu{
					// 	Text: mi18n.T("&剛体表示"),
					// 	Items: []declarative.MenuItem{
					// 		declarative.Action{
					// 			Text:        mi18n.T("&前面表示"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowRigidBodyFront,
					// 			AssignTo:    &cw.showRigidBodyFrontAction,
					// 		},
					// 		declarative.Action{
					// 			Text:        mi18n.T("&埋め込み表示"),
					// 			Checkable:   true,
					// 			OnTriggered: cw.TriggerShowRigidBodyBack,
					// 			AssignTo:    &cw.showRigidBodyBackAction,
					// 		},
					// 	},
					// },
					// declarative.Action{
					// 	Text:        mi18n.T("&ジョイント表示"),
					// 	Checkable:   true,
					// 	OnTriggered: cw.TriggerShowJoint,
					// 	AssignTo:    &cw.showJointAction,
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
				Items: helpMenuItems,
			},
			declarative.Menu{
				Text: mi18n.T("&言語"),
				Items: []declarative.MenuItem{
					declarative.Action{
						Text:        "日本語",
						OnTriggered: func() { cw.onChangeLanguage("ja") },
					},
					declarative.Action{
						Text:        "English",
						OnTriggered: func() { cw.onChangeLanguage("en") },
					},
					declarative.Action{
						Text:        "中文",
						OnTriggered: func() { cw.onChangeLanguage("zh") },
					},
					declarative.Action{
						Text:        "한국어",
						OnTriggered: func() { cw.onChangeLanguage("ko") },
					},
				},
			},
		},
		Children: []declarative.Widget{
			declarative.TabWidget{
				AssignTo: &cw.tabWidget,
				Pages:    tabPages,
			},
		},
		OnClosing: func(canceled *bool, reason walk.CloseReason) {
			// controllerStateを読み取り、ビューワーが閉じていない場合は確認ダイアログを表示
			if !cw.appConfig.IsCloseConfirm() {
				cw.shared.SetClosed(true)
				return
			}
			if !cw.shared.IsClosed() {
				if result := walk.MsgBox(
					nil,
					mi18n.T("終了確認"),
					mi18n.T("終了確認メッセージ"),
					walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel,
				); result == walk.DlgCmdOK {
					// ユーザーがOKを選んだ場合、 viewerState の isClosed を true にする
					cw.shared.SetClosed(true)
				} else {
					// 閉じない場合はキャンセル
					*canceled = true
				}
			}
		},
	}).Create(); err != nil {
		return nil, err
	}

	// 初期設定
	cw.shared.SetFrame(0.0) // フレーム初期化
	cw.TriggerFps30Limit()  // 30fps物理ON

	// cw.enabledPhysicsAction.SetChecked(true) // フレームドロップON
	// // cw.enabledFrameDropAction.SetChecked(true) // 30fps制限

	cw.SetPosition(positionX, positionY)

	return cw, nil
}

// OnClose はウィンドウを閉じるときの処理
func (cw *ControlWindow) OnClose() {
	// コントローラStateのisClosedをtrueにする
	cw.shared.SetClosed(true)
}

// Run はメインウィンドウを起動する
func (cw *ControlWindow) Run() {
	cw.MainWindow.Run()
}

func (cw *ControlWindow) Dispose() {
	cw.Close()
}

func (cw *ControlWindow) WindowSize() (int, int) {
	size := cw.Size()
	return size.Width, size.Height
}

func (cw *ControlWindow) SetPosition(x, y int) {
	cw.SetX(x)
	cw.SetY(y)
}

func (cw *ControlWindow) onChangeLanguage(lang string) {
	if result := walk.MsgBox(
		cw.MainWindow,
		mi18n.TWithLocale(lang, "言語変更"),
		mi18n.TWithLocale(lang, "言語変更メッセージ"),
		walk.MsgBoxOKCancel|walk.MsgBoxIconInformation,
	); result == walk.DlgCmdOK {
		mi18n.SetLang(lang)
		cw.shared.SetClosed(true)
	}
}

func (cw *ControlWindow) triggerLogLevel() {
	mlog.SetLevel(mlog.INFO)
	if cw.logLevelDebugAction.Checked() {
		mlog.SetLevel(mlog.DEBUG)
	}
	if cw.logLevelViewerVerboseAction.Checked() {
		mlog.I("exe階層に「viewerPng」フォルダを作成し、画面描画中の連番pngを出力し続けます\n画面サイズ: 1920x1080、視野角: 40.0、カメラ位置: (0, 10, 45)、カメラ角度: (0, 0, 0) ")
		mlog.SetLevel(mlog.VIEWER_VERBOSE)
	}
	if cw.logLevelIkVerboseAction.Checked() {
		mlog.SetLevel(mlog.IK_VERBOSE)
	}
	if cw.logLevelVerboseAction.Checked() {
		mlog.SetLevel(mlog.VERBOSE)
	}
}

// ------- 以下、再生状態の取得・設定メソッド -------

func (cw *ControlWindow) SetPlaying(playing bool) {
	cw.shared.SetPlaying(playing)
}

func (cw *ControlWindow) Playing() bool {
	return cw.shared.Playing()
}

func (cw *ControlWindow) SetFrame(frame float32) {
	cw.shared.SetFrame(frame)
}

func (cw *ControlWindow) Frame() float32 {
	return cw.shared.Frame()
}

func (cw *ControlWindow) SetMaxFrame(frame float32) {
	cw.shared.SetMaxFrame(frame)
}

func (cw *ControlWindow) MaxFrame() float32 {
	return cw.shared.MaxFrame()
}

// ------- 以下、モデルやモーションの格納・取得メソッド -------

func (cw *ControlWindow) StoreModel(windowIndex int, modelIndex int, model *pmx.PmxModel) {
	cw.shared.StoreModel(windowIndex, modelIndex, model)
}

func (cw *ControlWindow) LoadModel(windowIndex int, modelIndex int) *pmx.PmxModel {
	return cw.shared.LoadModel(windowIndex, modelIndex)
}

func (cw *ControlWindow) StoreMotion(windowIndex int, modelIndex int, motion *vmd.VmdMotion) {
	cw.shared.StoreMotion(windowIndex, modelIndex, motion)
}

func (cw *ControlWindow) LoadMotion(windowIndex int, modelIndex int) *vmd.VmdMotion {
	return cw.shared.LoadMotion(windowIndex, modelIndex)
}

// ------- 以下、メニューから呼ばれるトリガーメソッド -------

func (cw *ControlWindow) TriggerEnabledPhysics() {
	cw.shared.SetEnabledPhysics(cw.enabledPhysicsAction.Checked())
}

func (cw *ControlWindow) TriggerPhysicsReset() {
	cw.shared.SetPhysicsReset(true)
}

func (cw *ControlWindow) TriggerFps30Limit() {
	cw.limitFps30Action.SetChecked(true)
	cw.limitFps60Action.SetChecked(false)
	cw.limitFpsUnLimitAction.SetChecked(false)
	cw.shared.SetFrameInterval(1.0 / 30.0)
}

func (cw *ControlWindow) TriggerFps60Limit() {
	cw.limitFps30Action.SetChecked(false)
	cw.limitFps60Action.SetChecked(true)
	cw.limitFpsUnLimitAction.SetChecked(false)
	cw.shared.SetFrameInterval(1.0 / 60.0)
}

func (cw *ControlWindow) TriggerUnLimitFps() {
	cw.limitFps30Action.SetChecked(false)
	cw.limitFps60Action.SetChecked(false)
	cw.limitFpsUnLimitAction.SetChecked(true)
	cw.shared.SetFrameInterval(0)
}
