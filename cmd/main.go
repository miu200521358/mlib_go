//go:build windows
// +build windows

// 指示: miu200521358
package main

import (
	"embed"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/cmd/pkg/infra/controller/ui"
	"github.com/miu200521358/mlib_go/pkg/infra/app"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
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

// env はビルド時の -ldflags で埋め込む環境値。
var env string

// main はサンプルアプリを起動する。
func main() {
	app.Run(app.RunOptions{
		ViewerCount: 2,
		AppFiles:    appFiles,
		I18nFiles:   appI18nFiles,
		AdjustConfig: func(appConfig *config.AppConfig) {
			config.ApplyBuildEnv(appConfig, env)
		},
		BuildMenuItems: func(baseServices base.IBaseServices) []declarative.MenuItem {
			return ui.NewMenuItems(baseServices.I18n(), baseServices.Logger())
		},
		BuildTabPages: ui.NewTabPages,
	})
}
