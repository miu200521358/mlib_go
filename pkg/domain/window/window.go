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
	ControlState() IAppState
}

type IViewWindow interface {
	Render(states []IAnimationState)
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	TriggerClose(window *glfw.Window)
	GetWindow() *glfw.Window
	ResetPhysicsStart()
	AppState() IAppState
}
