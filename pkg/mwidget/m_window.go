package mwidget

import (
	"embed"
	"fmt"
	"syscall"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutil"
)

const (
	MAIN_WINDOW_WIDTH  = 512
	MAIN_WINDOW_HEIGHT = 768
)

type MWindow struct {
	// 横並びであるか否か
	isHorizontal bool
	walk.MainWindow
	GlWindow *GlWindow
}

func NewMWindow(resourceFiles embed.FS, isHorizontal bool) (*MWindow, error) {
	appConfig := mutil.ReadAppConfig(resourceFiles)

	var mw *walk.MainWindow

	if err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion),
		Size:     declarative.Size{Width: MAIN_WINDOW_WIDTH, Height: MAIN_WINDOW_HEIGHT},
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear},
	}).Create(); err != nil {
		return nil, err
	}
	mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		walk.App().Exit(0)
	})
	iconImg, err := mutil.ReadIconFile(resourceFiles)
	if err != nil {
		return nil, err
	}
	icon, err := walk.NewIconFromImageForDPI(iconImg, 96)
	if err != nil {
		return nil, err
	}
	mw.SetIcon(icon)

	return &MWindow{MainWindow: *mw, isHorizontal: isHorizontal}, nil
}

func (mWindow *MWindow) Center() {
	// スクリーンの解像度を取得
	screenWidth := GetSystemMetrics(SM_CXSCREEN)
	screenHeight := GetSystemMetrics(SM_CYSCREEN)

	// ウィンドウのサイズを取得
	windowSize := mWindow.Size()

	glWindowSize := walk.Size{Width: 0, Height: 0}
	if mWindow.GlWindow != nil {
		glWindowSize = mWindow.GlWindow.Size()
	}

	// ウィンドウを中央に配置
	if mWindow.isHorizontal {
		centerX := (screenWidth - (windowSize.Width + glWindowSize.Width)) / 2
		centerY := (screenHeight - windowSize.Height) / 2

		mWindow.SetX(centerX + glWindowSize.Width)
		mWindow.SetY(centerY)

		if mWindow.GlWindow != nil {
			mWindow.GlWindow.SetPos(centerX, centerY)
		}
	} else {
		centerX := (screenWidth - windowSize.Width) / 2
		centerY := (screenHeight - (windowSize.Height + glWindowSize.Height)) / 2

		mWindow.SetX(centerX)
		mWindow.SetY(centerY + glWindowSize.Height)

		if mWindow.GlWindow != nil {
			mWindow.GlWindow.SetPos(centerX, centerY)
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
