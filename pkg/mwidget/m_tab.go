package mwidget

import (
	"fmt"

	"github.com/miu200521358/walk/pkg/walk"
)

type MTabWidget struct {
	*walk.TabWidget
}

func NewMTabWidget(mWindow *MWindow) *MTabWidget {
	tabWidget, err := walk.NewTabWidget(mWindow)
	CheckError(err, mWindow, "タブウィジェット生成エラー")

	bg, err := walk.NewSystemColorBrush(walk.SysColorBackground)
	CheckError(err, mWindow, "背景色生成エラー")
	tabWidget.SetBackground(bg)

	return &MTabWidget{TabWidget: tabWidget}
}

type MTabPage struct {
	*walk.TabPage
}

func NewMTabPage(mWindow *MWindow, tabWidget *MTabWidget, title string) *MTabPage {
	tabPage, err := walk.NewTabPage()
	CheckError(err, mWindow, fmt.Sprintf("[%s]タブページ生成エラー", title))

	tabWidget.Pages().Add(tabPage)
	tabPage.SetTitle(title)

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DFace)
	CheckError(err, mWindow, "背景色生成エラー")
	tabPage.SetBackground(bg)

	return &MTabPage{tabPage}
}
