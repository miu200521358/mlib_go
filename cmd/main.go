package main

import (
	"embed"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/miu200521358/mlib_go/cmd/resources"
	"github.com/miu200521358/mlib_go/pkg/front/mtheme"
	"github.com/miu200521358/mlib_go/pkg/front/widget/file_picker"
	"github.com/miu200521358/mlib_go/pkg/pmx/pmx_reader"
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
	entry := widget.NewEntry()
	p, _ := file_picker.NewFilePicker(
		&w,
		"PmxPath",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		map[string]string{"*.pmx": "Pmx Files (*.pmx)", "*.*": "All Files (*.*)"},
		20,
		&pmx_reader.PmxReader{},
		func(path string) {
			entry.SetText(path)
		})
	container := container.New(layout.NewVBoxLayout(),
		widget.NewButton("これはボタンです", func() {
			dialog.ShowInformation("確認", "これはダイアログです", w)
		}),
		entry,
		p.Container,
		layout.NewSpacer())
	w.SetContent(container)

	w.Resize(fyne.NewSize(1024, 768))
	w.CenterOnScreen()
	w.ShowAndRun()
}
