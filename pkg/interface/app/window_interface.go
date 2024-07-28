//go:build windows
// +build windows

package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

type IControlWindow interface {
	Run()
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	Enabled() bool
	SetEnabled(enabled bool)
	SetFrame(frame float64)
}

type IViewWindow interface {
	Animate(states []state.IAnimationState, nextStates []state.IAnimationState, timeStep float32) ([]state.IAnimationState, []state.IAnimationState)
	Dispose()
	Close()
	Size() (int, int)
	SetPosition(x, y int)
	TriggerClose(window *glfw.Window)
	GetWindow() *glfw.Window
	ResetPhysics(animationStates []state.IAnimationState)
	AppState() state.IAppState
	Title() string
	OverrideTextureId() uint32
	SetOverrideTextureId(id uint32)
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
	SetEnabledOnlyButton(enabled bool)
	Enabled() bool
}
