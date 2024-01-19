package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/miu200521358/mlib_go/pkg/front/theme"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&theme.MTheme{})
	w := a.NewWindow("font")
	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			widget.NewLabel("こんにちは、ファイン"),
			widget.NewLabel("これは日本語のラベルです"),
			widget.NewButton("これはボタンです", func() {
				dialog.ShowInformation("確認", "これはダイアログです", w)
			}),
			layout.NewSpacer(),
		),
	)

	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}
