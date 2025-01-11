//go:build windows
// +build windows

package widget

import (
	"github.com/miu200521358/walk/pkg/walk"
)

type MTabWidget struct {
	*walk.TabWidget
}

func NewMTabWidget(w *walk.MainWindow) *MTabWidget {
	tabWidget, err := walk.NewTabWidget(w)
	if err != nil {
		RaiseError(err)
	}

	bg, err := walk.NewSystemColorBrush(walk.SysColorBackground)
	if err != nil {
		RaiseError(err)
	}
	tabWidget.SetBackground(bg)

	return &MTabWidget{TabWidget: tabWidget}
}

func (tabWidget *MTabWidget) Enabled(enabled bool) {
	for i := range tabWidget.Pages().Len() {
		for j := range tabWidget.Pages().At(i).Children().Len() {
			tabWidget.Pages().At(i).Children().At(j).SetEnabled(enabled)
		}
	}
}

type MTabPage struct {
	*walk.TabPage
}

func NewMTabPage(title string) *MTabPage {
	tabPage, err := walk.NewTabPage()
	if err != nil {
		RaiseError(err)
	}
	tabPage.SetTitle(title)

	bg, err := walk.NewSystemColorBrush(walk.SysColorInactiveCaption)
	if err != nil {
		RaiseError(err)
	}
	tabPage.SetBackground(bg)

	return &MTabPage{TabPage: tabPage}
}

func (tabPage *MTabPage) SetEnabled(enabled bool) {
	for i := range tabPage.Children().Len() {
		tabPage.Children().At(i).SetEnabled(enabled)
	}
}
