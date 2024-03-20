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
	walk.MainWindow
	// 横並びであるか否か
	isHorizontal bool
	// 描画ウィンドウ
	GlWindows []*GlWindow
}

func NewMWindow(resourceFiles embed.FS, isHorizontal bool, width int, height int) (*MWindow, error) {
	appConfig := mutils.LoadAppConfig(resourceFiles)

	var mw *walk.MainWindow

	if err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion),
		Size:     declarative.Size{Width: width, Height: height},
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear},
	}).Create(); err != nil {
		return nil, err
	}

	mainWindow := &MWindow{MainWindow: *mw, isHorizontal: isHorizontal, GlWindows: []*GlWindow{}}
	mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if len(mainWindow.GlWindows) > 0 {
			for _, glWindow := range mainWindow.GlWindows {
				glWindow.SetShouldClose(true)
			}
		}
		// mw.Dispose()
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

	return mainWindow, nil
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
