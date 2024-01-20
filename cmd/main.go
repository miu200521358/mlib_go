package main

import (
	"embed"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/miu200521358/mlib_go/cmd/resources"
	"github.com/miu200521358/mlib_go/pkg/front/mtheme"
	"github.com/miu200521358/mlib_go/pkg/front/widget/file_picker"
	"github.com/miu200521358/mlib_go/pkg/utils/config"
)

//go:embed resources/app_config.json
var appConfig embed.FS

func main() {
	appConfig := config.ReadAppConfig(appConfig)
	a := app.New()
	a.Settings().SetTheme(&mtheme.MTheme{})
	a.SetIcon(resources.AppIcon)
	w := a.NewWindow(fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion))
	pmx_output_file_picker, _ := file_picker.NewPmxSaveFilePicker(
		&w,
		"",
		"出力Pmxファイル",
		"出力Pmxファイルパスを入力もしくは選択してください",
		func(path string) {})
	pmx_file_picker, _ := file_picker.NewPmxReadFilePicker(
		&w,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		func(path string) {
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
			pmx_output_file_picker.PathEntry.SetText(outputPath)
		})

	// glWindow, err := gl_window.NewGlWindow(&a, "GL Window", resources.AppIcon, func() {}, func() {}, func(fyne.Position, []fyne.URI) {})
	// if err != nil {
	// 	panic(err)
	// }

	container := container.New(layout.NewVBoxLayout(),
		pmx_file_picker.Container,
		pmx_output_file_picker.Container,
		widget.NewButton("変換", func() {
			// glWindow.ShowAndRun()
		}),
		layout.NewSpacer())
	w.SetContent(container)

	w.Resize(fyne.NewSize(1024, 768))
	w.CenterOnScreen()
	w.ShowAndRun()
}
