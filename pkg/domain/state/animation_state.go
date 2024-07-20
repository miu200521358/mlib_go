//go:build windows
// +build windows

package state

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
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
	AnimatePhysics(physics *mbt.MPhysics, appState IAppState)
	AnimateAfterPhysics(physics *mbt.MPhysics, appState IAppState)
	RenderModel() IRenderModel
	SetRenderModel(model IRenderModel)
	Render(shader *mgl.MShader, appState IAppState)
}

type IRenderModel interface {
	Render(shader *mgl.MShader, appState IAppState, animationState IAnimationState)
	Hash() string
}
