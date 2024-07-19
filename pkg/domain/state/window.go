package state

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type IControlWindow interface {
	Run()
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	AppState() IAppState
}

type IViewWindow interface {
	Render(states []IAnimationState, timeStep float32)
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	TriggerClose(window *glfw.Window)
	GetWindow() *glfw.Window
	ResetPhysicsStart()
	AppState() IAppState
}
