//go:build windows
// +build windows

package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

type IControlWindow interface {
	Run()
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
}

type IViewWindow interface {
	InitRenderModel(modelIndex int, model *pmx.PmxModel) state.IRenderModel
	Render(states []state.IAnimationState, nextState state.IAnimationState, timeStep float32)
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	TriggerClose(window *glfw.Window)
	GetWindow() *glfw.Window
	ResetPhysics(animationStates []state.IAnimationState)
	AppState() state.IAppState
}

type IPlayer interface {
	Playing() bool
	TriggerPlay(playing bool)
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
