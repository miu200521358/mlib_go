//go:build windows
// +build windows

package controller

import (
	"fmt"
	"time"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// ControlWindow は操作画面(コントローラウィンドウ)を管理する
type ControlWindow struct {
	*walk.MainWindow

	shared    *state.SharedState // SharedState への参照
	appConfig *mconfig.AppConfig // アプリケーション設定

	tabWidget        *walk.TabWidget    // タブウィジェット
	consoleView      *ConsoleView       // コンソールビュー
	enabledInPlaying func(enabled bool) // 再生中に無効化するウィジェット

	leftButtonPressed bool // 左ボタン押下フラグ

	// メニューアクション
	enabledFrameDropAction      *walk.Action // フレームドロップON/OFF
	enabledPhysicsAction        *walk.Action // 物理ON/OFF
	physicsResetAction          *walk.Action // 物理リセット
	showNormalAction            *walk.Action // ボーンデバッグ表示
	showWireAction              *walk.Action // ワイヤーフレームデバッグ表示
	showOverrideAction          *walk.Action // オーバーライドデバッグ表示
	showSelectedVertexAction    *walk.Action // 選択頂点デバッグ表示
	showBoneAllAction           *walk.Action // 全ボーンデバッグ表示
	showBoneIkAction            *walk.Action // IKボーンデバッグ表示
	showBoneEffectorAction      *walk.Action // 付与親ボーンデバッグ表示
	showBoneFixedAction         *walk.Action // 軸制限ボーンデバッグ表示
	showBoneRotateAction        *walk.Action // 回転ボーンデバッグ表示
	showBoneTranslateAction     *walk.Action // 移動ボーンデバッグ表示
	showBoneVisibleAction       *walk.Action // 表示ボーンデバッグ表示
	showRigidBodyFrontAction    *walk.Action // 剛体デバッグ表示(前面)
	showRigidBodyBackAction     *walk.Action // 剛体デバッグ表示(埋め込み)
	showJointAction             *walk.Action // ジョイントデバッグ表示
	showInfoAction              *walk.Action // 情報デバッグ表示
	limitFps30Action            *walk.Action // 30FPS制限
	limitFps60Action            *walk.Action // 60FPS制限
	limitFpsUnLimitAction       *walk.Action // FPS無制限
	cameraSyncAction            *walk.Action // カメラ同期
	logLevelDebugAction         *walk.Action // デバッグメッセージ表示
	logLevelVerboseAction       *walk.Action // 冗長メッセージ表示
	logLevelIkVerboseAction     *walk.Action // IK冗長メッセージ表示
	logLevelViewerVerboseAction *walk.Action // ビューワー冗長メッセージ表示
	linkWindowAction            *walk.Action // ウィンドウ同期
}

// Run はコントローラウィンドウを実行する
func NewControlWindow(
	shared *state.SharedState,
	appConfig *mconfig.AppConfig,
	helpMenuItems []declarative.MenuItem,
	tabPages []declarative.TabPage,
	enabledInPlaying func(enabled bool),
	width, height, positionX, positionY int,
) (*ControlWindow, error) {
	cw := &ControlWindow{
		shared:           shared,
		appConfig:        appConfig,
		enabledInPlaying: enabledInPlaying,
	}

	logMenuItems := []declarative.MenuItem{
		declarative.Action{
			Text: mi18n.T("&使い方"),
			OnTriggered: func() {
				mlog.ILT(mi18n.T("メイン画面の使い方"), "%s", mi18n.T("メイン画面の使い方メッセージ"))
			},
		},
		declarative.Separator{},
		declarative.Action{
			Text:        mi18n.T("&画面移動連結"),
			Checkable:   true,
			OnTriggered: cw.triggerWindowLinkage,
			AssignTo:    &cw.linkWindowAction,
		},
		declarative.Separator{},
		declarative.Action{
			Text:        mi18n.T("&デバッグログ表示"),
			Checkable:   true,
			OnTriggered: cw.triggerLogLevel,
			AssignTo:    &cw.logLevelDebugAction,
		},
	}

	if appConfig.IsEnvDev() {
		// 開発時のみ冗長ログ表示を追加
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
					declarative.Action{
						Text:        mi18n.T("&フレームドロップON"),
						Checkable:   true,
						OnTriggered: cw.TriggerEnabledFrameDrop,
						AssignTo:    &cw.enabledFrameDropAction,
					},
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
					declarative.Action{
						Text:        mi18n.T("&情報表示"),
						Checkable:   true,
						OnTriggered: cw.TriggerShowInfo,
						AssignTo:    &cw.showInfoAction,
					},
					declarative.Separator{},
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
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&法線表示"),
						Checkable:   true,
						OnTriggered: cw.TriggerShowNormal,
						AssignTo:    &cw.showNormalAction,
					},
					declarative.Action{
						Text:        mi18n.T("&ワイヤーフレーム表示"),
						Checkable:   true,
						OnTriggered: cw.TriggerShowWire,
						AssignTo:    &cw.showWireAction,
					},
					declarative.Separator{},
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
					declarative.Action{
						Text:        mi18n.T("&カメラ同期"),
						Checkable:   true,
						OnTriggered: cw.TriggerCameraSync,
						AssignTo:    &cw.cameraSyncAction,
					},
					declarative.Action{
						Text:        mi18n.T("&サブビューワーオーバーレイ"),
						Checkable:   true,
						OnTriggered: cw.TriggerShowOverride,
						AssignTo:    &cw.showOverrideAction,
					},
					declarative.Action{
						Text: mi18n.T("&サブビューワーオーバーレイの使い方"),
						OnTriggered: func() {
							mlog.ILT(mi18n.T("&サブビューワーオーバーレイの使い方"),
								mi18n.T("サブビューワーオーバーレイの使い方メッセージ"))
						},
					},
					declarative.Separator{},
					declarative.Menu{
						Text: mi18n.T("&ボーン表示"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&全ボーン"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowBoneAll,
								AssignTo:    &cw.showBoneAllAction,
							},
							declarative.Separator{},
							declarative.Action{
								Text:        mi18n.T("&IKボーン"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowBoneIk,
								AssignTo:    &cw.showBoneIkAction,
							},
							declarative.Action{
								Text:        mi18n.T("&付与親ボーン"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowBoneEffector,
								AssignTo:    &cw.showBoneEffectorAction,
							},
							declarative.Action{
								Text:        mi18n.T("&軸制限ボーン"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowBoneFixed,
								AssignTo:    &cw.showBoneFixedAction,
							},
							declarative.Action{
								Text:        mi18n.T("&回転ボーン"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowBoneRotate,
								AssignTo:    &cw.showBoneRotateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&移動ボーン"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowBoneTranslate,
								AssignTo:    &cw.showBoneTranslateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&表示ボーン"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowBoneVisible,
								AssignTo:    &cw.showBoneVisibleAction,
							},
						},
					},
					declarative.Menu{
						Text: mi18n.T("&剛体表示"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&前面表示"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowRigidBodyFront,
								AssignTo:    &cw.showRigidBodyFrontAction,
							},
							declarative.Action{
								Text:        mi18n.T("&埋め込み表示"),
								Checkable:   true,
								OnTriggered: cw.TriggerShowRigidBodyBack,
								AssignTo:    &cw.showRigidBodyBackAction,
							},
						},
					},
					declarative.Action{
						Text:        mi18n.T("&ジョイント表示"),
						Checkable:   true,
						OnTriggered: cw.TriggerShowJoint,
						AssignTo:    &cw.showJointAction,
					},
					declarative.Separator{},
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
			// // ユーザーがOKを選んだ場合、 viewerState の isClosed を true にする
			// os.WriteFile(fmt.Sprintf("log_%s.txt", time.Now().Format("20060102_150405")),
			// 	[]byte(cw.consoleView.Console.Text()), 0644)

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
					cw.shared.SetClosed(true)
				} else {
					// 閉じない場合はキャンセル
					*canceled = true
				}
			}
		},
		OnClickActivate: func() {
			if !cw.shared.IsInitializedAllWindows() {
				// 初期化が終わってない場合、スルー
				return
			}
			cw.shared.TriggerLinkedFocus(-1)
		},
		// OnDeactivate: func() {
		// 	if !cw.shared.IsInitializedAllWindows() {
		// 		// 初期化が終わってない場合、スルー
		// 		return
		// 	}

		// 	// コントローラウィンドウが非アクティブ状態
		// 	mlog.IS("(C.2) ControlWindow deactivate")
		// },
		OnEnterSizeMove: func() {
			// 移動サイズ変更開始
			if cw.shared.IsWindowLinkage() {
				x, y := cw.GetPosition()
				cw.shared.SetControlWindowPosition(x, y, 0, 0)
			}
		},
		OnExitSizeMove: func() {
			// 移動サイズ変更終了
			if cw.shared.IsWindowLinkage() {
				cw.shared.SetMovedControlWindow(true)
				x, y := cw.GetPosition()
				prevPosX, prevPosY, _, _ := cw.shared.ControlWindowPosition()
				diffX := x - prevPosX
				diffY := y - prevPosY
				cw.shared.SetControlWindowPosition(x, y, diffX, diffY)
			}
		},
		OnMinimize: func() {
			if !cw.shared.IsInitializedAllWindows() {
				// 初期化が終わってない場合、スルー
				return
			}

			cw.shared.SyncMinimize(-1)
		},
		OnRestore: func() {
			if !cw.shared.IsInitializedAllWindows() {
				// 初期化が終わってない場合、スルー
				return
			}

			cw.shared.SyncRestore(-1)
		},
	}).Create(); err != nil {
		return nil, err
	}

	// 初期設定
	cw.shared.SetFrame(0.0)                  // フレーム初期化
	cw.shared.SetMaxFrame(1.0)               // 最大フレーム初期化
	cw.enabledPhysicsAction.SetChecked(true) // 物理ON
	cw.TriggerEnabledPhysics()

	// ウィンドウ移動同期
	if mconfig.LoadUserConfigBool(mconfig.KeyWindowLinkage, true) {
		cw.linkWindowAction.SetChecked(true)
		cw.triggerWindowLinkage()
	}

	//FPS制限
	fpsLimit := mconfig.LoadUserConfigInt(mconfig.KeyFpsLimit, 60)
	switch fpsLimit {
	case 30:
		cw.TriggerFps30Limit()
	case 60:
		cw.TriggerFps60Limit()
	case -1:
		cw.TriggerUnLimitFps()
	}

	// フレームドロップ
	if mconfig.LoadUserConfigBool(mconfig.KeyFrameDrop, true) {
		cw.enabledFrameDropAction.SetChecked(true)
		cw.TriggerEnabledFrameDrop()
	}

	// コンソールを追加で作成
	if cv, err := NewConsoleView(cw, width/10, height/10); err != nil {
		return nil, err
	} else {
		cw.consoleView = cv
	}
	// // ログ出力先をコンソールビューに設定
	// log.SetOutput(cw.consoleView)

	cw.SetPosition(positionX, positionY)
	cw.shared.SetInitializedControlWindow(true)
	// コントローラウィンドウハンドルを保持
	cw.shared.SetControlWindowHandle(int32(cw.Handle()))
	// コントローラウィンドウ位置を保持
	cw.shared.SetControlWindowPosition(positionX, positionY, 0, 0)

	// フォーカスチェック
	cw.checkFocus()

	return cw, nil
}

func (cw *ControlWindow) checkFocus() {
	go func() {
		for {
			if cw.shared.IsFocusControlWindow() {
				cw.Synchronize(func() {
					cw.SetForegroundWindow()
				})
				cw.shared.KeepFocus()
				cw.shared.SetFocusControlWindow(false)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
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

func (cw *ControlWindow) GetPosition() (int, int) {
	return cw.X(), cw.Y()
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

func (cw *ControlWindow) triggerWindowLinkage() {
	cw.shared.SetWindowLinkage(cw.linkWindowAction.Checked())
	mconfig.SaveUserConfigBool(mconfig.KeyWindowLinkage, cw.linkWindowAction.Checked())
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

func (cw *ControlWindow) EnabledInPlaying(enabled bool) {
	cw.enabledInPlaying(enabled)
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

func (cw *ControlWindow) StoreSelectedMaterialIndexes(windowIndex int, modelIndex int, indexes []int) {
	cw.shared.StoreSelectedMaterialIndexes(windowIndex, modelIndex, indexes)
}

func (cw *ControlWindow) LoadSelectedMaterialIndexes(windowIndex int, modelIndex int) []int {
	return cw.shared.LoadSelectedMaterialIndexes(windowIndex, modelIndex)
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
	mconfig.SaveUserConfigInt(mconfig.KeyFpsLimit, 30)
}

func (cw *ControlWindow) TriggerFps60Limit() {
	cw.limitFps30Action.SetChecked(false)
	cw.limitFps60Action.SetChecked(true)
	cw.limitFpsUnLimitAction.SetChecked(false)
	cw.shared.SetFrameInterval(1.0 / 60.0)
	mconfig.SaveUserConfigInt(mconfig.KeyFpsLimit, 60)
}

func (cw *ControlWindow) TriggerUnLimitFps() {
	cw.limitFps30Action.SetChecked(false)
	cw.limitFps60Action.SetChecked(false)
	cw.limitFpsUnLimitAction.SetChecked(true)
	cw.shared.SetFrameInterval(-1)
	mconfig.SaveUserConfigInt(mconfig.KeyFpsLimit, -1)
}

func (cw *ControlWindow) TriggerEnabledFrameDrop() {
	cw.shared.SetEnabledFrameDrop(cw.enabledFrameDropAction.Checked())
	cw.shared.SetChangedEnableDropFrame(true)
	mconfig.SaveUserConfigBool(mconfig.KeyFrameDrop, cw.enabledFrameDropAction.Checked())
}

func (cw *ControlWindow) TriggerShowNormal() {
	cw.shared.SetShowNormal(cw.showNormalAction.Checked())
}

func (cw *ControlWindow) TriggerShowWire() {
	cw.shared.SetShowWire(cw.showWireAction.Checked())
}

func (cw *ControlWindow) TriggerShowOverride() {
	cw.shared.SetShowOverride(cw.showOverrideAction.Checked())
}

func (cw *ControlWindow) TriggerShowSelectedVertex() {
	cw.shared.SetShowSelectedVertex(cw.showSelectedVertexAction.Checked())
}

func (cw *ControlWindow) updateShowBoneFlag() {
	cw.shared.UpdateFlags(
		map[uint32]bool{
			state.FlagShowBoneAll:       cw.showBoneAllAction.Checked(),
			state.FlagShowBoneIk:        cw.showBoneIkAction.Checked(),
			state.FlagShowBoneEffector:  cw.showBoneEffectorAction.Checked(),
			state.FlagShowBoneFixed:     cw.showBoneFixedAction.Checked(),
			state.FlagShowBoneRotate:    cw.showBoneRotateAction.Checked(),
			state.FlagShowBoneTranslate: cw.showBoneTranslateAction.Checked(),
			state.FlagShowBoneVisible:   cw.showBoneVisibleAction.Checked(),
		},
	)
}

func (cw *ControlWindow) TriggerShowBoneAll() {
	if cw.showBoneAllAction.Checked() {
		cw.showBoneIkAction.SetChecked(false)
		cw.showBoneEffectorAction.SetChecked(false)
		cw.showBoneFixedAction.SetChecked(false)
		cw.showBoneRotateAction.SetChecked(false)
		cw.showBoneTranslateAction.SetChecked(false)
		cw.showBoneVisibleAction.SetChecked(false)
	}
	cw.updateShowBoneFlag()
}

func (cw *ControlWindow) TriggerShowBoneIk() {
	if cw.showBoneIkAction.Checked() {
		cw.showBoneAllAction.SetChecked(false)
	}
	cw.updateShowBoneFlag()
}

func (cw *ControlWindow) TriggerShowBoneEffector() {
	if cw.showBoneEffectorAction.Checked() {
		cw.showBoneAllAction.SetChecked(false)
	}
	cw.updateShowBoneFlag()
}

func (cw *ControlWindow) TriggerShowBoneFixed() {
	if cw.showBoneFixedAction.Checked() {
		cw.showBoneAllAction.SetChecked(false)
	}
	cw.updateShowBoneFlag()
}

func (cw *ControlWindow) TriggerShowBoneRotate() {
	if cw.showBoneRotateAction.Checked() {
		cw.showBoneAllAction.SetChecked(false)
	}
	cw.updateShowBoneFlag()
}

func (cw *ControlWindow) TriggerShowBoneTranslate() {
	if cw.showBoneTranslateAction.Checked() {
		cw.showBoneAllAction.SetChecked(false)
	}
	cw.updateShowBoneFlag()
}

func (cw *ControlWindow) TriggerShowBoneVisible() {
	if cw.showBoneVisibleAction.Checked() {
		cw.showBoneAllAction.SetChecked(false)
	}
	cw.updateShowBoneFlag()
}

func (cw *ControlWindow) TriggerShowRigidBodyFront() {
	cw.shared.SetShowRigidBodyFront(cw.showRigidBodyFrontAction.Checked())
}

func (cw *ControlWindow) TriggerShowRigidBodyBack() {
	cw.shared.SetShowRigidBodyBack(cw.showRigidBodyBackAction.Checked())
}

func (cw *ControlWindow) TriggerShowJoint() {
	cw.shared.SetShowJoint(cw.showJointAction.Checked())
}

func (cw *ControlWindow) TriggerShowInfo() {
	cw.shared.SetShowInfo(cw.showInfoAction.Checked())
}

func (cw *ControlWindow) TriggerCameraSync() {
	cw.shared.SetCameraSync(cw.cameraSyncAction.Checked())
}
