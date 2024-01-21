package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/miu200521358/mlib_go/cmd/resources"
	"github.com/miu200521358/mlib_go/pkg/front/mtheme"
	"github.com/miu200521358/mlib_go/pkg/front/widget/file_picker"
	"github.com/miu200521358/mlib_go/pkg/front/widget/gl_window"
	"github.com/miu200521358/mlib_go/pkg/pmx/pmx_model"
	"github.com/miu200521358/mlib_go/pkg/utils/config"
)

//go:embed resources/app_config.json
var appConfig embed.FS

func init() {
	runtime.LockOSThread()
}

func main() {
	appConfig := config.ReadAppConfig(appConfig)
	a := app.New()
	a.Settings().SetTheme(&mtheme.MTheme{})
	a.SetIcon(resources.AppIcon)
	w := a.NewWindow(fmt.Sprintf("%s %s", appConfig.AppName, appConfig.AppVersion))
	pmxOutputFilePicker, _ := file_picker.NewPmxSaveFilePicker(
		&w,
		"",
		"出力Pmxファイル",
		"出力Pmxファイルパスを入力もしくは選択してください",
		func(path string) {})
	pmxFilePicker, _ := file_picker.NewPmxReadFilePicker(
		&w,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		func(path string) {
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_out"+ext)
			pmxOutputFilePicker.PathEntry.SetText(outputPath)
		})
	console := widget.NewMultiLineEntry()

	glWindow, err := gl_window.NewGlWindow(
		&a,
		"GL Window",
		resources.AppIcon,
		func() {},
		func() {},
		func(fyne.Position, []fyne.URI) {},
	)
	if err != nil {
		panic(err)
	}

	container := container.New(layout.NewVBoxLayout(),
		pmxFilePicker.Container,
		pmxOutputFilePicker.Container,
		widget.NewButton("変換", func() {
			data, err := pmxFilePicker.GetData()
			if err != nil {
				panic(err)
			}
			model := data.(*pmx_model.PmxModel)
			glWindow.AddData(model)

			console.Append(fmt.Sprintf("モデル名: %s\n", model.Name))
			console.Append(fmt.Sprintf("頂点数: %d\n", len(model.Vertices.Indexes)))
			console.Append(fmt.Sprintf("面数: %d\n", len(model.Faces.Indexes)))
			console.Append(fmt.Sprintf("材質数: %d\n", len(model.Materials.Indexes)))
			console.Append(fmt.Sprintf("ボーン数: %d\n", len(model.Bones.Indexes)))
			console.Append(fmt.Sprintf("表情数: %d\n", len(model.Morphs.Indexes)))
		}),
		console,
		layout.NewSpacer())
	w.SetContent(container)

	w.Resize(fyne.NewSize(1024, 768))
	w.CenterOnScreen()
	w.ShowAndRun()
}
