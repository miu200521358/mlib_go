package mwidget

import (
	"embed"
	"fmt"
	"syscall"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type MWindow struct {
	*walk.MainWindow
	TabWidget            *MTabWidget  // タブウィジェット
	isHorizontal         bool         // 横並びであるか否か
	GlWindows            []*GlWindow  // 描画ウィンドウ
	frameDropAction      *walk.Action // フレームドロップON/OFF
	physicsAction        *walk.Action // 物理ON/OFF
	physicsResetAction   *walk.Action // 物理リセット
	boneDebugAction      *walk.Action // ボーンデバッグ表示
	rigidBodyDebugAction *walk.Action // 剛体デバッグ表示
	jointDebugAction     *walk.Action // ジョイントデバッグ表示
}

func NewMWindow(resourceFiles embed.FS, isHorizontal bool, width int, height int) (*MWindow, error) {
	appConfig := mutils.LoadAppConfig(resourceFiles)

	mainWindow := &MWindow{isHorizontal: isHorizontal, GlWindows: []*GlWindow{}}

	if err := (declarative.MainWindow{
		AssignTo: &mainWindow.MainWindow,
		Title:    fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion),
		Size:     getWindowSize(width, height),
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear, MarginsZero: true, SpacingZero: true},
		MenuItems: []declarative.MenuItem{
			declarative.Menu{
				Text: "&モデル描画",
				Items: []declarative.MenuItem{
					declarative.Action{
						Text:        "&フレームドロップON/OFF",
						Checkable:   true,
						OnTriggered: mainWindow.frameDropTriggered,
						AssignTo:    &mainWindow.frameDropAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        "&物理ON/OFF",
						Checkable:   true,
						OnTriggered: mainWindow.physicsTriggered,
						AssignTo:    &mainWindow.physicsAction,
					},
					declarative.Action{
						Text:        "&物理リセット",
						OnTriggered: mainWindow.physicsResetTriggered,
						AssignTo:    &mainWindow.physicsResetAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        "&ボーンデバッグ表示",
						Checkable:   true,
						OnTriggered: mainWindow.boneDebugViewTriggered,
						AssignTo:    &mainWindow.boneDebugAction,
					},
					declarative.Separator{},
					declarative.Action{
						Text:        "&剛体デバッグ表示",
						Checkable:   true,
						OnTriggered: mainWindow.rigidBodyDebugViewTriggered,
						AssignTo:    &mainWindow.rigidBodyDebugAction,
					},
					declarative.Action{
						Text:        "&ジョイントデバッグ表示",
						Checkable:   true,
						OnTriggered: mainWindow.jointDebugViewTriggered,
						AssignTo:    &mainWindow.jointDebugAction,
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
		if len(mainWindow.GlWindows) > 0 {
			for _, glWindow := range mainWindow.GlWindows {
				glWindow.SetShouldClose(true)
			}
		}
		walk.App().Exit(0)
	})

	iconImg, err := mutils.LoadIconFile(resourceFiles)
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
	CheckError(err, mainWindow, "背景色生成エラー")
	mainWindow.SetBackground(bg)

	return mainWindow, nil
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
