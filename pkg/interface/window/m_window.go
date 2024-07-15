//go:build windows
// +build windows

package window

import (
	"fmt"
	"image"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/sys/windows"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/interface/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type MWindow struct {
	*walk.MainWindow
	UiState                  *widget.UiState      // UI状態
	TabWidget                *widget.MTabWidget   // タブウィジェット
	isHorizontal             bool                 // 横並びであるか否か
	GlWindows                []*GlWindow          // 描画ウィンドウ
	ConsoleView              *widget.ConsoleView  // コンソールビュー
	MotionPlayer             *widget.MotionPlayer // モーションプレイヤー
	enabledFrameDropAction   *walk.Action         // フレームドロップON/OFF
	enabledPhysicsAction     *walk.Action         // 物理ON/OFF
	physicsResetAction       *walk.Action         // 物理リセット
	showNormalAction         *walk.Action         // ボーンデバッグ表示
	showWireAction           *walk.Action         // ワイヤーフレームデバッグ表示
	showSelectedVertexAction *walk.Action         // 選択頂点デバッグ表示
	showBoneAllAction        *walk.Action         // 全ボーンデバッグ表示
	showBoneIkAction         *walk.Action         // IKボーンデバッグ表示
	showBoneEffectorAction   *walk.Action         // 付与親ボーンデバッグ表示
	showBoneFixedAction      *walk.Action         // 軸制限ボーンデバッグ表示
	showBoneRotateAction     *walk.Action         // 回転ボーンデバッグ表示
	showBoneTranslateAction  *walk.Action         // 移動ボーンデバッグ表示
	showBoneVisibleAction    *walk.Action         // 表示ボーンデバッグ表示
	showRigidBodyFrontAction *walk.Action         // 剛体デバッグ表示(前面)
	showRigidBodyBackAction  *walk.Action         // 剛体デバッグ表示(埋め込み)
	showJointAction          *walk.Action         // ジョイントデバッグ表示
	showInfoAction           *walk.Action         // 情報デバッグ表示
	limitFps30Action         *walk.Action         // 30FPS制限
	limitFps60Action         *walk.Action         // 60FPS制限
	unLimitFpsAction         *walk.Action         // FPS無制限
	unLimitFpsDeformAction   *walk.Action         // デフォームFPS無制限
	logLevelDebugAction      *walk.Action         // デバッグメッセージ表示
	logLevelVerboseAction    *walk.Action         // 冗長メッセージ表示
	logLevelIkVerboseAction  *walk.Action         // IK冗長メッセージ表示
}

func NewMWindow(
	width int,
	height int,
	isHorizontal bool,
	helpMenuItemsFunc func() []declarative.MenuItem,
	iconImg *image.Image,
	appConfig *mconfig.AppConfig,
	uiState *widget.UiState,
) (*MWindow, error) {
	mainWindow := &MWindow{
		isHorizontal: isHorizontal,
		UiState:      uiState,
		GlWindows:    []*GlWindow{},
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
			OnTriggered: mainWindow.logLevelTriggered,
			AssignTo:    &mainWindow.logLevelDebugAction,
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
				OnTriggered: mainWindow.logLevelTriggered,
				AssignTo:    &mainWindow.logLevelVerboseAction,
			})
		logMenuItems = append(logMenuItems,
			declarative.Action{
				Text:        mi18n.T("&IK冗長ログ表示"),
				Checkable:   true,
				OnTriggered: mainWindow.logLevelTriggered,
				AssignTo:    &mainWindow.logLevelIkVerboseAction,
			})
	}

	fpsLImitMenuItems := []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("&30fps制限"),
			Checkable:   true,
			OnTriggered: mainWindow.onTriggerFps30Limit,
			AssignTo:    &mainWindow.limitFps30Action,
		},
		declarative.Action{
			Text:        mi18n.T("&60fps制限"),
			Checkable:   true,
			OnTriggered: mainWindow.onTriggerFps60Limit,
			AssignTo:    &mainWindow.limitFps60Action,
		},
		declarative.Action{
			Text:        mi18n.T("&fps無制限"),
			Checkable:   true,
			OnTriggered: mainWindow.onTriggerUnLimitFps,
			AssignTo:    &mainWindow.unLimitFpsAction,
		},
	}

	if !appConfig.IsEnvProd() {
		// 開発時にだけ描画無制限モードを追加
		fpsLImitMenuItems = append(fpsLImitMenuItems,
			declarative.Action{
				Text:        "&デフォームfps無制限",
				Checkable:   true,
				OnTriggered: mainWindow.onTriggerUnLimitFpsDeform,
				AssignTo:    &mainWindow.unLimitFpsDeformAction,
			})
	}

	if err := (declarative.MainWindow{
		AssignTo: &mainWindow.MainWindow,
		Title:    fmt.Sprintf("%s %s", appConfig.Name, appConfig.Version),
		Size:     getWindowSize(width, height),
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear, MarginsZero: true, SpacingZero: true},
		MenuItems: []declarative.MenuItem{
			declarative.Menu{
				Text: mi18n.T("&ビューワー"),
				Items: []declarative.MenuItem{
					declarative.Action{
						Text:        mi18n.T("&フレームドロップON"),
						Checkable:   true,
						OnTriggered: mainWindow.onTriggerFrameDrop,
						AssignTo:    &mainWindow.enabledFrameDropAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&物理ON/OFF"),
						Checkable:   true,
						OnTriggered: mainWindow.onTriggerPhysics,
						AssignTo:    &mainWindow.enabledPhysicsAction,
					},
					declarative.Action{
						Text:        mi18n.T("&物理リセット"),
						OnTriggered: mainWindow.onTriggerPhysicsReset,
						AssignTo:    &mainWindow.physicsResetAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&法線表示"),
						Checkable:   true,
						OnTriggered: mainWindow.onTriggerShowNormal,
						AssignTo:    &mainWindow.showNormalAction,
					},
					declarative.Action{
						Text:        mi18n.T("&ワイヤーフレーム表示"),
						Checkable:   true,
						OnTriggered: mainWindow.onTriggerShowWire,
						AssignTo:    &mainWindow.showWireAction,
					},
					declarative.Action{
						Text:        mi18n.T("&選択頂点表示"),
						Checkable:   true,
						OnTriggered: mainWindow.onTriggerShowSelectedVertex,
						AssignTo:    &mainWindow.showSelectedVertexAction,
					},
					declarative.Separator{},
					declarative.Menu{
						Text: mi18n.T("&ボーン表示"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&全ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowBoneAll,
								AssignTo:    &mainWindow.showBoneAllAction,
							},
							declarative.Separator{},
							declarative.Action{
								Text:        mi18n.T("&IKボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowBoneOne,
								AssignTo:    &mainWindow.showBoneIkAction,
							},
							declarative.Action{
								Text:        mi18n.T("&付与親ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowBoneOne,
								AssignTo:    &mainWindow.showBoneEffectorAction,
							},
							declarative.Action{
								Text:        mi18n.T("&軸制限ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowBoneOne,
								AssignTo:    &mainWindow.showBoneFixedAction,
							},
							declarative.Action{
								Text:        mi18n.T("&回転ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowBoneOne,
								AssignTo:    &mainWindow.showBoneRotateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&移動ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowBoneOne,
								AssignTo:    &mainWindow.showBoneTranslateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&表示ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowBoneOne,
								AssignTo:    &mainWindow.showBoneVisibleAction,
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
								OnTriggered: mainWindow.onTriggerShowRigidBodyFront,
								AssignTo:    &mainWindow.showRigidBodyFrontAction,
							},
							declarative.Action{
								Text:        mi18n.T("&埋め込み表示"),
								Checkable:   true,
								OnTriggered: mainWindow.onTriggerShowRigidBodyBack,
								AssignTo:    &mainWindow.showRigidBodyBackAction,
							},
						},
					},
					declarative.Action{
						Text:        mi18n.T("&ジョイント表示"),
						Checkable:   true,
						OnTriggered: mainWindow.onTriggerShowJoint,
						AssignTo:    &mainWindow.showJointAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&情報表示"),
						Checkable:   true,
						OnTriggered: mainWindow.onTriggerShowInfo,
						AssignTo:    &mainWindow.showInfoAction,
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
						OnTriggered: func() { mainWindow.langTriggered("ja") },
					},
					declarative.Action{
						Text:        "English",
						OnTriggered: func() { mainWindow.langTriggered("en") },
					},
					declarative.Action{
						Text:        "中文",
						OnTriggered: func() { mainWindow.langTriggered("zh") },
					},
					declarative.Action{
						Text:        "한국어",
						OnTriggered: func() { mainWindow.langTriggered("ko") },
					},
				},
			},
		},
	}).Create(); err != nil {
		return nil, err
	}

	// 最初は物理ON
	mainWindow.enabledPhysicsAction.SetChecked(true)
	mainWindow.onTriggerPhysics()
	// 最初はフレームドロップON
	mainWindow.enabledFrameDropAction.SetChecked(true)
	mainWindow.onTriggerFrameDrop()
	// 最初は30fps制限
	mainWindow.limitFps30Action.SetChecked(true)
	mainWindow.onTriggerFps30Limit()

	mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if len(mainWindow.GlWindows) > 0 && !mgl.CheckOpenGLError() {
			for _, glWindow := range mainWindow.GlWindows {
				glWindow.SetShouldClose(true)
			}
		}
		walk.App().Exit(0)
	})

	icon, err := walk.NewIconFromImageForDPI(*iconImg, 96)
	if err != nil {
		return nil, err
	}
	mainWindow.SetIcon(icon)

	// タブウィジェット追加
	mainWindow.TabWidget = widget.NewMTabWidget(mainWindow.MainWindow)
	mainWindow.Children().Add(mainWindow.TabWidget)

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DShadow)
	widget.CheckError(err, mainWindow.MainWindow, mi18n.T("背景色生成エラー"))
	mainWindow.SetBackground(bg)

	return mainWindow, nil
}

func (w *MWindow) logLevelTriggered() {
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

func (w *MWindow) langTriggered(lang string) {
	mi18n.SetLang(lang)
	walk.MsgBox(
		w.MainWindow,
		mi18n.TWithLocale(lang, "LanguageChanged.Title"),
		mi18n.TWithLocale(lang, "LanguageChanged.Message"),
		walk.MsgBoxOK|walk.MsgBoxIconInformation,
	)
	w.Close()
}

func (w *MWindow) SetCheckWireDebugView(checked bool) {
	w.showWireAction.SetChecked(checked)
	w.onTriggerShowWire()
}

func (w *MWindow) SetCheckSelectedVertexDebugView(checked bool) {
	w.showSelectedVertexAction.SetChecked(checked)
	w.onTriggerShowSelectedVertex()
}

func (w *MWindow) onTriggerShowNormal() {
	w.UiState.IsShowNormal = w.showNormalAction.Checked()
}

func (w *MWindow) onTriggerShowWire() {
	w.UiState.IsShowWire = w.showWireAction.Checked()
}

func (w *MWindow) onTriggerShowSelectedVertex() {
	w.UiState.IsShowSelectedVertex = w.showSelectedVertexAction.Checked()
}

func (w *MWindow) onTriggerShowBoneAll() {
	w.showBoneIkAction.SetChecked(false)
	w.showBoneEffectorAction.SetChecked(false)
	w.showBoneFixedAction.SetChecked(false)
	w.showBoneRotateAction.SetChecked(false)
	w.showBoneTranslateAction.SetChecked(false)
	w.showBoneVisibleAction.SetChecked(false)

	w.onTriggerShowBone()
}

func (w *MWindow) onTriggerShowBoneOne() {
	w.showBoneAllAction.SetChecked(false)

	w.onTriggerShowBone()
}

func (w *MWindow) onTriggerShowBone() {
	w.UiState.IsShowBones[pmx.BONE_FLAG_NONE] = w.showBoneAllAction.Checked()
	w.UiState.IsShowBones[pmx.BONE_FLAG_IS_IK] = w.showBoneIkAction.Checked()
	w.UiState.IsShowBones[pmx.BONE_FLAG_IS_EXTERNAL_ROTATION] = w.showBoneEffectorAction.Checked()
	w.UiState.IsShowBones[pmx.BONE_FLAG_IS_EXTERNAL_TRANSLATION] = w.showBoneEffectorAction.Checked()
	w.UiState.IsShowBones[pmx.BONE_FLAG_HAS_FIXED_AXIS] = w.showBoneFixedAction.Checked()
	w.UiState.IsShowBones[pmx.BONE_FLAG_CAN_ROTATE] = w.showBoneRotateAction.Checked()
	w.UiState.IsShowBones[pmx.BONE_FLAG_CAN_TRANSLATE] = w.showBoneTranslateAction.Checked()
	w.UiState.IsShowBones[pmx.BONE_FLAG_IS_VISIBLE] = w.showBoneVisibleAction.Checked()
}

func (w *MWindow) onTriggerShowRigidBodyFront() {
	w.showRigidBodyBackAction.SetChecked(false)
	w.UiState.IsShowRigidBodyFront = w.showRigidBodyFrontAction.Checked()
	w.UiState.IsShowRigidBodyBack = false
}

func (w *MWindow) onTriggerShowRigidBodyBack() {
	w.showRigidBodyFrontAction.SetChecked(false)
	w.UiState.IsShowRigidBodyBack = w.showRigidBodyBackAction.Checked()
	w.UiState.IsShowRigidBodyFront = false
}

func (w *MWindow) onTriggerShowJoint() {
	w.UiState.IsShowJoint = w.showJointAction.Checked()
}

func (w *MWindow) onTriggerShowInfo() {
	w.UiState.IsShowInfo = w.showInfoAction.Checked()
}

func (w *MWindow) onTriggerFps30Limit() {
	w.limitFps30Action.SetChecked(true)
	w.limitFps60Action.SetChecked(false)
	w.unLimitFpsAction.SetChecked(false)
	w.unLimitFpsDeformAction.SetChecked(false)

	w.UiState.IsLimitFps30 = true
	w.UiState.IsLimitFps60 = false
	w.UiState.IsUnLimitFps = false
	w.UiState.IsUnLimitFpsDeform = false

	w.UiState.SpfLimit = 1 / 30.0
}

func (w *MWindow) onTriggerFps60Limit() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(true)
	w.unLimitFpsAction.SetChecked(false)
	w.unLimitFpsDeformAction.SetChecked(false)

	w.UiState.IsLimitFps30 = false
	w.UiState.IsLimitFps60 = true
	w.UiState.IsUnLimitFps = false
	w.UiState.IsUnLimitFpsDeform = false

	w.UiState.SpfLimit = 1 / 60.0
}

func (w *MWindow) onTriggerUnLimitFps() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(false)
	w.unLimitFpsAction.SetChecked(true)
	w.unLimitFpsDeformAction.SetChecked(false)

	w.UiState.IsLimitFps30 = false
	w.UiState.IsLimitFps60 = false
	w.UiState.IsUnLimitFps = true
	w.UiState.IsUnLimitFpsDeform = false

	w.UiState.SpfLimit = -1.0
}

func (w *MWindow) onTriggerUnLimitFpsDeform() {
	w.limitFps30Action.SetChecked(false)
	w.limitFps60Action.SetChecked(false)
	w.unLimitFpsAction.SetChecked(false)
	w.unLimitFpsDeformAction.SetChecked(true)

	w.UiState.IsLimitFps30 = false
	w.UiState.IsLimitFps60 = false
	w.UiState.IsUnLimitFps = false
	w.UiState.IsUnLimitFpsDeform = true

	w.UiState.SpfLimit = -2.0
}

func (w *MWindow) onTriggerPhysics() {
	w.UiState.EnabledPhysics = w.enabledPhysicsAction.Checked()
}

func (w *MWindow) onTriggerPhysicsReset() {
	w.UiState.DoPhysicsReset = true
}

func (w *MWindow) onTriggerFrameDrop() {
	w.UiState.EnabledFrameDrop = w.enabledFrameDropAction.Checked()
}

func getWindowSize(width int, height int) declarative.Size {
	screenWidth := GetSystemMetrics(SM_CXSCREEN)
	screenHeight := GetSystemMetrics(SM_CYSCREEN)

	if width > screenWidth-50 {
		width = screenWidth - 50
	}
	if height > screenHeight-50 {
		height = screenHeight - 50
	}

	return declarative.Size{Width: width, Height: height}
}

func (w *MWindow) AddGlWindow(glWindow *GlWindow) {
	w.GlWindows = append(w.GlWindows, glWindow)
}

func (w *MWindow) Center() {
	// スクリーンの解像度を取得
	screenWidth := GetSystemMetrics(SM_CXSCREEN)
	screenHeight := GetSystemMetrics(SM_CYSCREEN)

	// ウィンドウのサイズを取得
	windowSize := w.Size()

	glWindowWidth := 0
	glWindowHeight := 0
	for _, glWindow := range w.GlWindows {
		glWindowWidth += glWindow.Size().Width
		glWindowHeight += glWindow.Size().Height
	}

	// ウィンドウを中央に配置
	if w.isHorizontal {
		centerX := (screenWidth - (windowSize.Width + glWindowWidth)) / 2
		centerY := (screenHeight - windowSize.Height) / 2

		centerX += glWindowWidth
		w.SetX(centerX)
		w.SetY(centerY)

		for _, glWindow := range w.GlWindows {
			centerX -= glWindow.Size().Width
			glWindow.SetPos(centerX, centerY)
		}
	} else {
		centerX := (screenWidth - windowSize.Width) / 2
		centerY := (screenHeight - (windowSize.Height + glWindowHeight)) / 2

		centerY += windowSize.Height
		w.SetX(centerX)
		w.SetY(centerY)

		for _, glWindow := range w.GlWindows {
			centerY -= glWindow.Size().Height
			glWindow.SetPos(centerX, centerY)
		}
	}
}

func (w *MWindow) Dispose() {
	for _, glWindow := range w.GlWindows {
		glWindow.TriggerClose(glWindow.Window)
	}
	w.MainWindow.Dispose()
	defer walk.App().Exit(0)
}

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	procMessageBeep      = user32.NewProc("MessageBeep")
	MB_ICONASTERISK      = 0x00000040
)

func (w *MWindow) Beep() {
	procMessageBeep.Call(uintptr(MB_ICONASTERISK))
}

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func GetSystemMetrics(nIndex int) int {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(nIndex))
	return int(ret)
}
