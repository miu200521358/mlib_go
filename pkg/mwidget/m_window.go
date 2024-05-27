//go:build windows
// +build windows

package mwidget

import (
	"embed"
	"fmt"
	"syscall"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type MWindow struct {
	*walk.MainWindow
	TabWidget               *MTabWidget  // タブウィジェット
	isHorizontal            bool         // 横並びであるか否か
	GlWindows               []*GlWindow  // 描画ウィンドウ
	ConsoleView             *ConsoleView // コンソールビュー
	frameDropAction         *walk.Action // フレームドロップON/OFF
	physicsAction           *walk.Action // 物理ON/OFF
	physicsResetAction      *walk.Action // 物理リセット
	normalDebugAction       *walk.Action // ボーンデバッグ表示
	boneDebugAction         *walk.Action // ボーンデバッグ表示
	rigidBodyDebugAction    *walk.Action // 剛体デバッグ表示
	jointDebugAction        *walk.Action // ジョイントデバッグ表示
	logLevelDebugAction     *walk.Action // デバッグメッセージ表示
	logLevelVerboseAction   *walk.Action // 冗長メッセージ表示
	logLevelIkVerboseAction *walk.Action // IK冗長メッセージ表示
}

func NewMWindow(
	resourceFiles embed.FS,
	appConfig *mconfig.AppConfig,
	isHorizontal bool,
	width int,
	height int,
	funcHelpMenuItems func() []declarative.MenuItem,
) (*MWindow, error) {
	mi18n.Initialize(resourceFiles)

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

	if appConfig.Env == "dev" {
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
						Text:        mi18n.T("&フレームドロップON/OFF"),
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
						Text:        mi18n.T("&法線デバッグ表示"),
						Checkable:   true,
						OnTriggered: mainWindow.normalDebugViewTriggered,
						AssignTo:    &mainWindow.normalDebugAction,
					},
					declarative.Action{
						Text:        mi18n.T("&ボーンデバッグ表示"),
						Checkable:   true,
						OnTriggered: mainWindow.boneDebugViewTriggered,
						AssignTo:    &mainWindow.boneDebugAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        mi18n.T("&剛体デバッグ表示"),
						Checkable:   true,
						OnTriggered: mainWindow.rigidBodyDebugViewTriggered,
						AssignTo:    &mainWindow.rigidBodyDebugAction,
					},
					declarative.Action{
						Text:        mi18n.T("&ジョイントデバッグ表示"),
						Checkable:   true,
						OnTriggered: mainWindow.jointDebugViewTriggered,
						AssignTo:    &mainWindow.jointDebugAction,
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

	mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if len(mainWindow.GlWindows) > 0 && !CheckOpenGLError() {
			for _, glWindow := range mainWindow.GlWindows {
				glWindow.SetShouldClose(true)
			}
		}
		walk.App().Exit(0)
	})

	iconImg, err := mconfig.LoadIconFile(resourceFiles)
	if err != nil {
		return nil, err
	}
	icon, err := walk.NewIconFromImageForDPI(iconImg, 96)
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

func (w *MWindow) normalDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.VisibleNormal = w.normalDebugAction.Checked()
	}
}

func (w *MWindow) boneDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.VisibleBone = w.boneDebugAction.Checked()
	}
}

func (w *MWindow) rigidBodyDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.Physics.VisibleRigidBody(w.rigidBodyDebugAction.Checked())
	}
}

func (w *MWindow) jointDebugViewTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.Physics.VisibleJoint(w.jointDebugAction.Checked())
	}
}

func (w *MWindow) physicsTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.EnablePhysics = w.physicsAction.Checked()
	}
}

func (w *MWindow) physicsResetTriggered() {
	for _, glWindow := range w.GlWindows {
		glWindow.ResetPhysics()
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

	glWindowSize := walk.Size{Width: 0, Height: 0}
	if w.GetMainGlWindow() != nil {
		glWindowSize = w.GetMainGlWindow().Size()
	}

	// ウィンドウを中央に配置
	if w.isHorizontal {
		centerX := (screenWidth - (windowSize.Width + glWindowSize.Width)) / 2
		centerY := (screenHeight - windowSize.Height) / 2

		w.SetX(centerX + glWindowSize.Width)
		w.SetY(centerY)

		if w.GetMainGlWindow() != nil {
			w.GetMainGlWindow().SetPos(centerX, centerY)
		}
	} else {
		centerX := (screenWidth - windowSize.Width) / 2
		centerY := (screenHeight - (windowSize.Height + glWindowSize.Height)) / 2

		w.SetX(centerX)
		w.SetY(centerY + glWindowSize.Height)

		if w.GetMainGlWindow() != nil {
			w.GetMainGlWindow().SetPos(centerX, centerY)
		}
	}
}

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

func GetSystemMetrics(nIndex int) int {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(nIndex))
	return int(ret)
}
