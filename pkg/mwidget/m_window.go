//go:build windows
// +build windows

package mwidget

import (
	"fmt"
	"image"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/sys/windows"

	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type MWindow struct {
	*walk.MainWindow
	TabWidget                 *MTabWidget  // タブウィジェット
	isHorizontal              bool         // 横並びであるか否か
	GlWindows                 []*GlWindow  // 描画ウィンドウ
	ConsoleView               *ConsoleView // コンソールビュー
	frameDropAction           *walk.Action // フレームドロップON/OFF
	physicsAction             *walk.Action // 物理ON/OFF
	physicsResetAction        *walk.Action // 物理リセット
	normalDebugAction         *walk.Action // ボーンデバッグ表示
	wireDebugAction           *walk.Action // ワイヤーフレームデバッグ表示
	selectedVertexDebugAction *walk.Action // 選択頂点デバッグ表示
	boneDebugAllAction        *walk.Action // 全ボーンデバッグ表示
	boneDebugIkAction         *walk.Action // IKボーンデバッグ表示
	boneDebugEffectorAction   *walk.Action // 付与親ボーンデバッグ表示
	boneDebugFixedAction      *walk.Action // 軸制限ボーンデバッグ表示
	boneDebugRotateAction     *walk.Action // 回転ボーンデバッグ表示
	boneDebugTranslateAction  *walk.Action // 移動ボーンデバッグ表示
	boneDebugVisibleAction    *walk.Action // 表示ボーンデバッグ表示
	rigidBodyFrontDebugAction *walk.Action // 剛体デバッグ表示(前面)
	rigidBodyBackDebugAction  *walk.Action // 剛体デバッグ表示(埋め込み)
	jointDebugAction          *walk.Action // ジョイントデバッグ表示
	infoDebugAction           *walk.Action // 情報デバッグ表示
	fps30LimitAction          *walk.Action // 30FPS制限
	fps60LimitAction          *walk.Action // 60FPS制限
	fpsUnLimitAction          *walk.Action // FPS無制限
	logLevelDebugAction       *walk.Action // デバッグメッセージ表示
	logLevelVerboseAction     *walk.Action // 冗長メッセージ表示
	logLevelIkVerboseAction   *walk.Action // IK冗長メッセージ表示
	fpsDeformUnLimitAction    *walk.Action // デフォームFPS無制限
}

func NewMWindow(
	width int,
	height int,
	funcHelpMenuItems func() []declarative.MenuItem,
	iconImg *image.Image,
	appConfig *mconfig.AppConfig,
	isHorizontal bool,
) (*MWindow, error) {
	mainWindow := &MWindow{
		isHorizontal: isHorizontal,
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
			OnTriggered: mainWindow.fps30LimitTriggered,
			AssignTo:    &mainWindow.fps30LimitAction,
		},
		declarative.Action{
			Text:        mi18n.T("&60fps制限"),
			Checkable:   true,
			OnTriggered: mainWindow.fps60LimitTriggered,
			AssignTo:    &mainWindow.fps60LimitAction,
		},
		declarative.Action{
			Text:        mi18n.T("&fps無制限"),
			Checkable:   true,
			OnTriggered: mainWindow.fpsUnLimitTriggered,
			AssignTo:    &mainWindow.fpsUnLimitAction,
		},
	}

	if !appConfig.IsEnvProd() {
		// 開発時にだけ描画無制限モードを追加
		fpsLImitMenuItems = append(fpsLImitMenuItems,
			declarative.Action{
				Text:        "&デフォームfps無制限",
				Checkable:   true,
				OnTriggered: mainWindow.fpsDeformUnLimitTriggered,
				AssignTo:    &mainWindow.fpsDeformUnLimitAction,
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
						OnTriggered: mainWindow.frameDropTriggered,
						AssignTo:    &mainWindow.frameDropAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&物理ON/OFF"),
						Checkable:   true,
						OnTriggered: mainWindow.physicsTriggered,
						AssignTo:    &mainWindow.physicsAction,
					},
					declarative.Action{
						Text:        mi18n.T("&物理リセット"),
						OnTriggered: mainWindow.physicsResetTriggered,
						AssignTo:    &mainWindow.physicsResetAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&法線表示"),
						Checkable:   true,
						OnTriggered: mainWindow.normalDebugViewTriggered,
						AssignTo:    &mainWindow.normalDebugAction,
					},
					declarative.Action{
						Text:        mi18n.T("&ワイヤーフレーム表示"),
						Checkable:   true,
						OnTriggered: mainWindow.wireDebugViewTriggered,
						AssignTo:    &mainWindow.wireDebugAction,
					},
					declarative.Action{
						Text:        mi18n.T("&選択頂点表示"),
						Checkable:   true,
						OnTriggered: mainWindow.selectedVertexDebugViewTriggered,
						AssignTo:    &mainWindow.selectedVertexDebugAction,
					},
					declarative.Separator{},
					declarative.Menu{
						Text: mi18n.T("&ボーン表示"),
						Items: []declarative.MenuItem{
							declarative.Action{
								Text:        mi18n.T("&全ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.boneDebugViewAllTriggered,
								AssignTo:    &mainWindow.boneDebugAllAction,
							},
							declarative.Separator{},
							declarative.Action{
								Text:        mi18n.T("&IKボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.boneDebugViewIndividualTriggered,
								AssignTo:    &mainWindow.boneDebugIkAction,
							},
							declarative.Action{
								Text:        mi18n.T("&付与親ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.boneDebugViewIndividualTriggered,
								AssignTo:    &mainWindow.boneDebugEffectorAction,
							},
							declarative.Action{
								Text:        mi18n.T("&軸制限ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.boneDebugViewIndividualTriggered,
								AssignTo:    &mainWindow.boneDebugFixedAction,
							},
							declarative.Action{
								Text:        mi18n.T("&回転ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.boneDebugViewIndividualTriggered,
								AssignTo:    &mainWindow.boneDebugRotateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&移動ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.boneDebugViewIndividualTriggered,
								AssignTo:    &mainWindow.boneDebugTranslateAction,
							},
							declarative.Action{
								Text:        mi18n.T("&表示ボーン"),
								Checkable:   true,
								OnTriggered: mainWindow.boneDebugViewIndividualTriggered,
								AssignTo:    &mainWindow.boneDebugVisibleAction,
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
								OnTriggered: mainWindow.rigidBodyDebugFrontViewTriggered,
								AssignTo:    &mainWindow.rigidBodyFrontDebugAction,
							},
							declarative.Action{
								Text:        mi18n.T("&埋め込み表示"),
								Checkable:   true,
								OnTriggered: mainWindow.rigidBodyDebugBackViewTriggered,
								AssignTo:    &mainWindow.rigidBodyBackDebugAction,
							},
						},
					},
					declarative.Action{
						Text:        mi18n.T("&ジョイント表示"),
						Checkable:   true,
						OnTriggered: mainWindow.jointDebugViewTriggered,
						AssignTo:    &mainWindow.jointDebugAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&情報表示"),
						Checkable:   true,
						OnTriggered: mainWindow.infoDebugViewTriggered,
						AssignTo:    &mainWindow.infoDebugAction,
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
				Items: funcHelpMenuItems(),
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
	mainWindow.physicsAction.SetChecked(true)
	// 最初はフレームドロップON
	mainWindow.frameDropAction.SetChecked(true)
	// 最初は30fps制限
	mainWindow.fps30LimitAction.SetChecked(true)

	mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if len(mainWindow.GlWindows) > 0 && !CheckOpenGLError() {
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
	mainWindow.TabWidget = NewMTabWidget(mainWindow)
	mainWindow.Children().Add(mainWindow.TabWidget)

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DShadow)
	CheckError(err, mainWindow, mi18n.T("背景色生成エラー"))
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
	w.wireDebugAction.SetChecked(checked)
	w.wireDebugViewTriggered()
}

func (w *MWindow) SetCheckSelectedVertexDebugView(checked bool) {
	w.selectedVertexDebugAction.SetChecked(checked)
	w.selectedVertexDebugViewTriggered()
}

func (w *MWindow) normalDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.VisibleNormal = w.normalDebugAction.Checked()
	}
}

func (w *MWindow) wireDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.VisibleWire = w.wireDebugAction.Checked()
	}
}

func (w *MWindow) selectedVertexDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.VisibleSelectedVertex = w.selectedVertexDebugAction.Checked()
	}
}

func (w *MWindow) boneDebugViewAllTriggered() {
	w.boneDebugIkAction.SetChecked(false)
	w.boneDebugEffectorAction.SetChecked(false)
	w.boneDebugFixedAction.SetChecked(false)
	w.boneDebugRotateAction.SetChecked(false)
	w.boneDebugTranslateAction.SetChecked(false)
	w.boneDebugVisibleAction.SetChecked(false)

	w.boneDebugViewTriggered()
}

func (w *MWindow) boneDebugViewIndividualTriggered() {
	w.boneDebugAllAction.SetChecked(false)

	w.boneDebugViewTriggered()
}

func (w *MWindow) boneDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		// 全ボーン表示
		glWindow.VisibleBones[pmx.BONE_FLAG_NONE] = w.boneDebugAllAction.Checked()
		// IKボーン表示
		glWindow.VisibleBones[pmx.BONE_FLAG_IS_IK] = w.boneDebugIkAction.Checked()
		// 付与親ボーン表示
		glWindow.VisibleBones[pmx.BONE_FLAG_IS_EXTERNAL_ROTATION] = w.boneDebugEffectorAction.Checked()
		glWindow.VisibleBones[pmx.BONE_FLAG_IS_EXTERNAL_TRANSLATION] = w.boneDebugEffectorAction.Checked()
		// 軸制限ボーン表示
		glWindow.VisibleBones[pmx.BONE_FLAG_HAS_FIXED_AXIS] = w.boneDebugFixedAction.Checked()
		// 回転ボーン表示
		glWindow.VisibleBones[pmx.BONE_FLAG_CAN_ROTATE] = w.boneDebugRotateAction.Checked()
		// 移動ボーン表示
		glWindow.VisibleBones[pmx.BONE_FLAG_CAN_TRANSLATE] = w.boneDebugTranslateAction.Checked()
		// 表示ボーン表示
		glWindow.VisibleBones[pmx.BONE_FLAG_IS_VISIBLE] = w.boneDebugVisibleAction.Checked()
	}
}

func (w *MWindow) rigidBodyDebugFrontViewTriggered() {
	w.rigidBodyBackDebugAction.SetChecked(false)
	for _, glWindow := range w.GlWindows {
		glWindow.Physics.VisibleRigidBody(w.rigidBodyFrontDebugAction.Checked())
		glWindow.Shader.IsDrawRigidBodyFront = true
	}
}

func (w *MWindow) rigidBodyDebugBackViewTriggered() {
	w.rigidBodyFrontDebugAction.SetChecked(false)
	for _, glWindow := range w.GlWindows {
		glWindow.Physics.VisibleRigidBody(w.rigidBodyBackDebugAction.Checked())
		glWindow.Shader.IsDrawRigidBodyFront = false
	}
}

func (w *MWindow) jointDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.Physics.VisibleJoint(w.jointDebugAction.Checked())
	}
}

func (w *MWindow) infoDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.isShowInfo = w.infoDebugAction.Checked()
	}
}

func (w *MWindow) fps30LimitTriggered() {
	w.fps30LimitAction.SetChecked(true)
	w.fps60LimitAction.SetChecked(false)
	w.fpsUnLimitAction.SetChecked(false)
	w.fpsDeformUnLimitAction.SetChecked(false)
	for _, glWindow := range w.GlWindows {
		glWindow.spfLimit = 1 / 30.0
	}
}

func (w *MWindow) fps60LimitTriggered() {
	w.fps30LimitAction.SetChecked(false)
	w.fps60LimitAction.SetChecked(true)
	w.fpsUnLimitAction.SetChecked(false)
	w.fpsDeformUnLimitAction.SetChecked(false)
	for _, glWindow := range w.GlWindows {
		glWindow.spfLimit = 1 / 60.0
	}
}

func (w *MWindow) fpsUnLimitTriggered() {
	w.fps30LimitAction.SetChecked(false)
	w.fps60LimitAction.SetChecked(false)
	w.fpsUnLimitAction.SetChecked(true)
	w.fpsDeformUnLimitAction.SetChecked(false)
	for _, glWindow := range w.GlWindows {
		glWindow.spfLimit = -1.0
	}
}

func (w *MWindow) fpsDeformUnLimitTriggered() {
	w.fps30LimitAction.SetChecked(false)
	w.fps60LimitAction.SetChecked(false)
	w.fpsUnLimitAction.SetChecked(false)
	w.fpsDeformUnLimitAction.SetChecked(true)
	for _, glWindow := range w.GlWindows {
		glWindow.spfLimit = -2.0
	}
}

func (w *MWindow) physicsTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.TriggerPhysicsEnabled(w.physicsAction.Checked())
	}
}

func (w *MWindow) physicsResetTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.TriggerPhysicsReset()
	}
}

func (w *MWindow) frameDropTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.EnableFrameDrop = w.frameDropAction.Checked()
	}
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

func (w *MWindow) GetMainGlWindow() *GlWindow {
	if len(w.GlWindows) > 0 {
		return w.GlWindows[0]
	}
	return nil
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
