//go:build windows
// +build windows

package mwidget

import (
	"fmt"

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

func NewMTabPage(mWindow *MWindow, tabWidget *MTabWidget, title string) *MTabPage {
	tabPage, err := walk.NewTabPage()
	CheckError(err, mWindow, fmt.Sprintf("[%s]%s", title, mi18n.T("タブページ生成エラー")))

	tabWidget.Pages().Add(tabPage)
	tabPage.SetTitle(title)

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DFace)
	CheckError(err, mWindow, mi18n.T("背景色生成エラー"))
	tabPage.SetBackground(bg)

	return &MTabPage{TabPage: tabPage, mWindow: mWindow}
}

type MultiPageMTabPage struct {
	*MTabPage
	Header                      *walk.Composite
	Pages                       []*walk.Composite
	navToolBar                  *walk.ToolBar
	PageComposite               *walk.Composite
	actions                     []*walk.Action
	currentAction               *walk.Action
	currentPage                 *walk.Composite
	currentPageChangedPublisher walk.EventPublisher
}

func NewMultiPageMTabPage(
	mWindow *MWindow, tabWidget *MTabWidget, title string, isHorizontal bool,
) (*MultiPageMTabPage, error) {
	mptp := &MultiPageMTabPage{
		MTabPage: NewMTabPage(mWindow, tabWidget, title),
		Pages:    make([]*walk.Composite, 0),
		actions:  make([]*walk.Action, 0),
	}
	mptp.SetLayout(walk.NewVBoxLayout())

	var err error
	mptp.Header, err = walk.NewComposite(mptp)
	if err != nil {
		return nil, err
	}
	mptp.Header.SetLayout(walk.NewVBoxLayout())

	// スクロール
	scrollView, err := walk.NewScrollView(mptp)
	if err != nil {
		return nil, err
	}
	scrollView.SetScrollbars(isHorizontal, !isHorizontal)
	scrollView.SetLayout(walk.NewHBoxLayout())

	// ナビゲーション用ツールバー
	mptp.navToolBar, err = walk.NewToolBarWithOrientationAndButtonStyle(
		scrollView, walk.Horizontal, walk.ToolBarButtonTextOnly)
	if err != nil {
		return nil, err
	}

	// ページ配置コンポーネント
	mptp.PageComposite, err = walk.NewComposite(mptp)
	if err != nil {
		return nil, err
	}
	mptp.PageComposite.SetLayout(walk.NewVBoxLayout())

	return mptp, nil
}

func (mptp *MultiPageMTabPage) AddPage(page *walk.Composite) error {
	action, err := mptp.newPageAction()
	if err != nil {
		return err
	}

	mptp.Pages = append(mptp.Pages, page)
	mptp.actions = append(mptp.actions, action)

	if err := mptp.updateNavigationToolBar(); err != nil {
		return err
	}

	if len(mptp.actions) > 0 {
		if err := mptp.setCurrentAction(len(mptp.actions) - 1); err != nil {
			return err
		}
	}

	return nil
}

func (mptp *MultiPageMTabPage) CurrentPage() *walk.Composite {
	return mptp.currentPage
}

func (mptp *MultiPageMTabPage) CurrentPageTitle() string {
	if mptp.currentAction == nil {
		return ""
	}

	return mptp.currentAction.Text()
}

func (mptp *MultiPageMTabPage) CurrentPageChanged() *walk.Event {
	return mptp.currentPageChangedPublisher.Event()
}

func (mptp *MultiPageMTabPage) newPageAction() (*walk.Action, error) {
	action := walk.NewAction()
	action.SetCheckable(true)
	action.SetExclusive(true)
	action.SetText(fmt.Sprintf("No. %d", len(mptp.actions)+1))

	action.Triggered().Attach(func() {
		mptp.setCurrentAction(len(mptp.actions) - 1)
	})

	return action, nil
}

func (mptp *MultiPageMTabPage) setCurrentAction(index int) error {
	defer func() {
		if !mptp.PageComposite.IsDisposed() {
			mptp.PageComposite.RestoreState()
		}
	}()

	mptp.SetFocus()

	if prevPage := mptp.currentPage; prevPage != nil {
		mptp.PageComposite.SaveState()
		prevPage.SetVisible(false)
		prevPage.SetParent(nil)
		prevPage.Dispose()
	}

	for i := range len(mptp.actions) {
		mptp.actions[i].SetChecked(false)
	}
	mptp.actions[index].SetChecked(true)
	mptp.currentPage = mptp.Pages[index]
	mptp.currentAction = mptp.actions[index]
	mptp.currentPageChangedPublisher.Publish()

	return nil
}

func (mptp *MultiPageMTabPage) updateNavigationToolBar() error {
	mptp.navToolBar.SetSuspended(true)
	defer mptp.navToolBar.SetSuspended(false)

	actions := mptp.navToolBar.Actions()

	if err := actions.Clear(); err != nil {
		return err
	}

	for _, action := range mptp.actions {
		if err := actions.Add(action); err != nil {
			return err
		}
	}

	if mptp.currentAction != nil {
		if !actions.Contains(mptp.currentAction) {
			for i, action := range mptp.actions {
				if action != mptp.currentAction {
					if err := mptp.setCurrentAction(i); err != nil {
						return err
					}

					break
				}
			}
		}
	}

	return nil
}
