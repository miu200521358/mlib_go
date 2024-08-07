//go:build windows
// +build windows

package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
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
	SetFrame(frame float32)
	UpdateSelectedVertexIndexes(indexes [][][]int)
	SetUpdateSelectedVertexIndexesFunc(f func([][][]int))
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
	GetViewerParameter() (float64, float64, *mmath.MVec2, *mmath.MVec3, *mmath.MVec3, *mmath.MVec3)
	UpdateViewerParameter(yaw, pitch float64, size *mmath.MVec2, cameraPos, cameraUp, lookAtCenter *mmath.MVec3)
}

type IPlayer interface {
	Playing() bool
	TriggerPlay(playing bool)
	PrevFrame() float32
	SetPrevFrame(v float32)
	Frame() float32
	SetFrame(v float32)
	MaxFrame() float32
	SetMaxFrame(max float32)
	UpdateMaxFrame(max float32)
	SetRange(min, max int)
	SetEnabled(enabled bool)
	SetEnabledOnlyButton(enabled bool)
	Enabled() bool
}
