//go:build windows
// +build windows

package state

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type IAnimationState interface {
	WindowIndex() int
	SetWindowIndex(index int)
	ModelIndex() int
	SetModelIndex(index int)
	Model() *pmx.PmxModel
	SetModel(model *pmx.PmxModel)
	Motion() *vmd.VmdMotion
	SetMotion(motion *vmd.VmdMotion)
	VmdDeltas() *delta.VmdDeltas
	SetVmdDeltas(deltas *delta.VmdDeltas)
	RenderDeltas() *delta.RenderDeltas
	SetRenderDeltas(deltas *delta.RenderDeltas)
	Frame() float64
	SetFrame(frame float64)
	AnimateBeforePhysics(appState IAppState, model *pmx.PmxModel) (*delta.VmdDeltas, *delta.RenderDeltas)
	AnimatePhysics(physics IPhysics, appState IAppState)
	AnimateAfterPhysics(physics IPhysics, appState IAppState)
	RenderModel() IRenderModel
	SetRenderModel(model IRenderModel)
	Render(shader IShader, appState IAppState)
}

type IRenderModel interface {
	Render(shader IShader, appState IAppState, animationState IAnimationState)
	Hash() string
}
