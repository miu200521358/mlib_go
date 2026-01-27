//go:build windows
// +build windows

// 指示: miu200521358
package main

import (
	"embed"
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/cmd/ui"
	"github.com/miu200521358/mlib_go/pkg/infra/app"
	"github.com/miu200521358/mlib_go/pkg/infra/base/err"
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/maudio"
	"github.com/miu200521358/mlib_go/pkg/infra/viewer"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

// init はOSスレッド固定とコンソール登録を行う。
func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(controller.ConsoleViewClass)
	})
}

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

// main はサンプルアプリを起動する。
func main() {
	viewerCount := 2

	boot, initErr := app.Init(appFiles, appI18nFiles, nil)
	if initErr != nil {
		if boot != nil {
			err.ShowFatalErrorDialog(boot.AppConfig, initErr)
		} else {
			err.ShowFatalErrorDialog(nil, initErr)
		}
		return
	}
	appConfig := boot.AppConfig
	baseServices := boot.BaseServices
	iconImage := boot.IconImage
	appIcon := boot.AppIcon
	audioPlayer := maudio.NewAudioPlayer()

	sharedState := state.NewSharedState(viewerCount)

	widths, heights, positionXs, positionYs := app.GetCenterSizeAndWidth(appConfig, viewerCount)

	var (
		controlWindow    *controller.ControlWindow
		controlWindowErr error
	)
	viewerManager := viewer.NewViewerManager(sharedState, baseServices)
	if iconImage != nil {
		viewerManager.SetWindowIcon(iconImage)
	}

	go func() {
		defer app.SafeExecute(appConfig, func() {
			widgets := &controller.MWidgets{
				Position: &walk.Point{X: positionXs[0], Y: positionYs[0]},
			}

			controlWindow, controlWindowErr = controller.NewControlWindow(
				sharedState,
				baseServices,
				ui.NewMenuItems(baseServices.I18n(), baseServices.Logger()),
				ui.NewTabPages(widgets, baseServices, audioPlayer),
				widths[0], heights[0], positionXs[0], positionYs[0], viewerCount,
			)
			if controlWindowErr != nil {
				err.ShowFatalErrorDialog(appConfig, controlWindowErr)
				return
			}
			if appIcon != nil {
				controlWindow.SetIcon(appIcon)
			}

			widgets.SetWindow(controlWindow)
			widgets.OnLoaded()
			controlWindow.Run()
		})
	}()

	if glfwErr := glfw.Init(); glfwErr != nil {
		err.ShowFatalErrorDialog(appConfig, fmt.Errorf(i18n.T("GLFWの初期化に失敗しました: %w"), glfwErr))
		return
	}

	defer app.SafeExecute(appConfig, func() {
		for n := 0; n < viewerCount; n++ {
			idx := n + 1
			if addWindowErr := viewerManager.AddWindow(
				fmt.Sprintf("Viewer%d", idx),
				widths[idx], heights[idx], positionXs[idx], positionYs[idx],
			); addWindowErr != nil {
				err.ShowFatalErrorDialog(appConfig, addWindowErr)
				return
			}
		}

		viewerManager.InitOverlay()
		viewerManager.Run()
	})
}
