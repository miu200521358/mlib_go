//go:build windows
// +build windows

package mwidget

import (
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
)

type MTabWidget struct {
	*walk.TabWidget
}

func NewMTabWidget(mWindow *MWindow) *MTabWidget {
	tabWidget, err := walk.NewTabWidget(mWindow)
	CheckError(err, mWindow, mi18n.T("タブウィジェット生成エラー"))

	bg, err := walk.NewSystemColorBrush(walk.SysColorBackground)
	CheckError(err, mWindow, mi18n.T("背景色生成エラー"))
	tabWidget.SetBackground(bg)

	return &MTabWidget{TabWidget: tabWidget}
}

type MTabPage struct {
	*walk.TabPage
	mWindow *MWindow
}

func NewMTabPage(mWindow *MWindow, tabWidget *MTabWidget, title string) (*MTabPage, error) {
	tabPage, err := walk.NewTabPage()
	if err != nil {
		return nil, err
	}

	tabWidget.Pages().Add(tabPage)
	tabPage.SetTitle(title)

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DFace)
	if err != nil {
		return nil, err
	}
	tabPage.SetBackground(bg)

	return &MTabPage{TabPage: tabPage, mWindow: mWindow}, nil
}
