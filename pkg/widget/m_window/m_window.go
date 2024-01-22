package m_window

import (
	"embed"
	"fmt"
	"syscall"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/utils/config"

)

type MWindow struct {
	walk.MainWindow
}

func NewMWindow(appConfigFile embed.FS) (*MWindow, error) {
	appConfig := config.ReadAppConfig(appConfigFile)

	var mw *walk.MainWindow

	if err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion),
		Size:     declarative.Size{Width: 1024, Height: 768},
		Layout:   declarative.VBox{Alignment: declarative.AlignHNearVNear},
	}).Create(); err != nil {
		panic(err)
	}

	return &MWindow{MainWindow: *mw}, nil
}

func (mWindow *MWindow) Center() {
	// スクリーンの解像度を取得
	screenWidth := GetSystemMetrics(SM_CXSCREEN)
	screenHeight := GetSystemMetrics(SM_CYSCREEN)

	// ウィンドウのサイズを取得
	windowSize := mWindow.Size()

	// ウィンドウを中央に配置
	mWindow.SetX((screenWidth - windowSize.Width) / 2)
	mWindow.SetY((screenHeight - windowSize.Height) / 2)
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
