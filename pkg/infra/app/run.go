//go:build windows
// +build windows

// 指示: miu200521358
package app

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/infra/base/err"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/maudio"
	"github.com/miu200521358/mlib_go/pkg/infra/viewer"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

// MenuItemsBuilder はメニュー項目生成処理を表す。
type MenuItemsBuilder func(baseServices base.IBaseServices) []declarative.MenuItem

// TabPagesBuilder はタブページ生成処理を表す。
type TabPagesBuilder func(widgets *controller.MWidgets, baseServices base.IBaseServices, audioPlayer audio_api.IAudioPlayer) []declarative.TabPage

// RunOptions は起動時のオプションを表す。
type RunOptions struct {
	ViewerCount    int
	AppFiles       embed.FS
	I18nFiles      embed.FS
	AdjustConfig   AppConfigAdjuster
	BuildMenuItems MenuItemsBuilder
	BuildTabPages  TabPagesBuilder
}

// Run は共通の起動フローを実行する。
func Run(options RunOptions) {
	viewerCount := options.ViewerCount
	if viewerCount <= 0 {
		viewerCount = 1
	}

	boot, initErr := Init(options.AppFiles, options.I18nFiles, options.AdjustConfig)
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
	if sharedState == nil {
		err.ShowFatalErrorDialog(appConfig, NewSharedStateInitFailed())
		return
	}

	widths, heights, positionXs, positionYs := GetCenterSizeAndWidth(appConfig, viewerCount)

	viewerManager := viewer.NewViewerManager(sharedState, baseServices)
	if iconImage != nil {
		viewerManager.SetWindowIcon(iconImage)
	}

	go func() {
		defer SafeExecute(appConfig, func() {
			widgets := &controller.MWidgets{
				Position: &walk.Point{X: positionXs[0], Y: positionYs[0]},
			}
			menuItems := []declarative.MenuItem(nil)
			if options.BuildMenuItems != nil {
				menuItems = options.BuildMenuItems(baseServices)
			}
			tabPages := []declarative.TabPage(nil)
			if options.BuildTabPages != nil {
				tabPages = options.BuildTabPages(widgets, baseServices, audioPlayer)
			}
			controlWindow, controlWindowErr := controller.NewControlWindow(
				sharedState,
				baseServices,
				menuItems,
				tabPages,
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
		err.ShowFatalErrorDialog(appConfig, NewGlfwInitFailed(glfwErr))
		return
	}

	defer SafeExecute(appConfig, func() {
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

// FindInitialPath は起動引数から対象拡張子のパスを取得する。
func FindInitialPath(args []string, exts ...string) string {
	if len(args) <= 1 {
		return ""
	}
	allowed := map[string]struct{}{}
	for _, ext := range exts {
		trimmed := strings.ToLower(strings.TrimSpace(ext))
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, ".") {
			trimmed = "." + trimmed
		}
		allowed[trimmed] = struct{}{}
	}
	if len(allowed) == 0 {
		return ""
	}
	for i := 1; i < len(args); i++ {
		path := strings.TrimSpace(args[i])
		path = strings.Trim(path, "\"")
		path = strings.Trim(path, "'")
		if path == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := allowed[ext]; ok {
			return path
		}
	}
	return ""
}
