package window

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/walk"
)

type IControlWindow interface {
	Run()
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	AppState() IAppState
	GetMainWindow() *walk.MainWindow
	InitTabWidget()
	AddTabPage(tabPage *walk.TabPage)
}

type IViewWindow interface {
	Render()
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	TriggerClose(window *glfw.Window)
	GetWindow() *glfw.Window
	ResetPhysicsStart()
	AppState() IAppState
}
