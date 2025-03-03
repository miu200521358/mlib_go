//go:build windows
// +build windows

package viewer

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
)

type ViewerList struct {
	shared     *state.SharedState // SharedState への参照
	appConfig  *mconfig.AppConfig // アプリケーション設定
	viewerList []*ViewWindow
}

func NewViewerList(shared *state.SharedState, appConfig *mconfig.AppConfig) *ViewerList {
	return &ViewerList{
		shared:     shared,
		appConfig:  appConfig,
		viewerList: make([]*ViewWindow, 0),
	}
}

// Add は ViewerList に ViewerWindow を追加します。
func (vl *ViewerList) Add(title string, width, height, positionX, positionY int) error {
	var mainViewerWindow *glfw.Window
	if len(vl.viewerList) > 0 {
		mainViewerWindow = vl.viewerList[0].Window
	}

	viewWindow, err := newViewWindow(
		len(vl.viewerList),
		title,
		width,
		height,
		positionX,
		positionY,
		vl.appConfig.IconImage,
		vl.appConfig.IsEnvProd(),
		mainViewerWindow,
		vl,
	)

	if err != nil {
		return err
	}

	vl.viewerList = append(vl.viewerList, viewWindow)

	return nil
}

func (vl *ViewerList) Run() {

	for !vl.shared.IsClosed() {
		glfw.PollEvents()

		for _, viewWindow := range vl.viewerList {
			viewWindow.Render(vl.shared)
		}
	}

	for _, viewWindow := range vl.viewerList {
		viewWindow.Destroy()
	}
}
