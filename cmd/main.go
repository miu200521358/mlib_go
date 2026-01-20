//go:build windows
// +build windows

// 指示: miu200521358
package main

import (
	"embed"
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/cmd/ui"
	"github.com/miu200521358/mlib_go/pkg/infra/app"
	infraconfig "github.com/miu200521358/mlib_go/pkg/infra/base/config"
	infraerr "github.com/miu200521358/mlib_go/pkg/infra/base/err"
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	infralogging "github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/viewer"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	sharedlogging "github.com/miu200521358/mlib_go/pkg/shared/base/logging"
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

	appConfig, err := infraconfig.LoadAppConfig(appFiles)
	if err != nil {
		infraerr.ShowFatalErrorDialog(nil, err)
		return
	}
	userConfig := infraconfig.NewUserConfigStore()
	if err := i18n.InitI18n(appI18nFiles, userConfig); err != nil {
		infraerr.ShowFatalErrorDialog(appConfig, err)
		return
	}
	logger := infralogging.NewLogger(i18n.Default())
	infralogging.SetDefaultLogger(logger)
	sharedlogging.SetDefaultLogger(logger)
	configStore := infraconfig.NewConfigStore(appConfig, userConfig)
	baseServices := &base.BaseServices{
		ConfigStore:   configStore,
		I18nService:   i18n.Default(),
		LoggerService: logger,
	}

	shared := state.NewSharedState(viewerCount)
	sharedState, ok := shared.(*state.SharedState)
	if !ok {
		infraerr.ShowFatalErrorDialog(appConfig, fmt.Errorf("共有状態の初期化に失敗しました"))
		return
	}

	widths, heights, positionXs, positionYs := app.GetCenterSizeAndWidth(appConfig, viewerCount)

	var controlWindow *controller.ControlWindow
	viewerManager := viewer.NewViewerManager(sharedState, baseServices)

	go func() {
		defer app.SafeExecute(appConfig, func() {
			widgets := &controller.MWidgets{
				Position: &walk.Point{X: positionXs[0], Y: positionYs[0]},
			}

			controlWindow, err = controller.NewControlWindow(
				sharedState,
				baseServices,
				ui.NewMenuItems(baseServices.I18n(), baseServices.Logger()),
				[]declarative.TabPage{ui.NewTabPage(widgets, baseServices)},
				widths[0], heights[0], positionXs[0], positionYs[0], viewerCount,
			)
			if err != nil {
				infraerr.ShowFatalErrorDialog(appConfig, err)
				return
			}

			widgets.SetWindow(controlWindow)
			widgets.OnLoaded()
			controlWindow.Run()
		})
	}()

	if err := glfw.Init(); err != nil {
		infraerr.ShowFatalErrorDialog(appConfig, fmt.Errorf("GLFWの初期化に失敗しました: %w", err))
		return
	}

	defer app.SafeExecute(appConfig, func() {
		for n := 0; n < viewerCount; n++ {
			idx := n + 1
			if err := viewerManager.AddWindow(
				fmt.Sprintf("Viewer%d", idx),
				widths[idx], heights[idx], positionXs[idx], positionYs[idx],
			); err != nil {
				infraerr.ShowFatalErrorDialog(appConfig, err)
				return
			}
		}

		viewerManager.InitOverlay()
		viewerManager.Run()
	})
}
