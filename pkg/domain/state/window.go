package state

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type IControlWindow interface {
	IAppState
	Run()
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
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

type IPlayer interface {
	Playing() bool
	SetPlaying(playing bool)
	PrevFrame() int
	SetPrevFrame(v int)
	Frame() float64
	SetFrame(v float64)
	MaxFrame() int
	SetMaxFrame(max int)
	UpdateMaxFrame(max int)
	SetRange(min, max int)
	SetEnabled(enabled bool)
}
