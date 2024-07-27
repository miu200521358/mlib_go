//go:build windows
// +build windows

package state

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
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
	DeformBeforePhysics(appState IAppState, model *pmx.PmxModel) (*delta.VmdDeltas, *delta.RenderDeltas)
	DeformPhysics(physics mbt.IPhysics, appState IAppState)
	DeformAfterPhysics(physics mbt.IPhysics, appState IAppState)
	RenderModel() IRenderModel
	SetRenderModel(model IRenderModel)
	Render(shader mgl.IShader, appState IAppState, cursorRect *mmath.MRect, windowSize *mmath.MVec2)
	Load(model *pmx.PmxModel)
}

type IRenderModel interface {
	Delete()
	Render(shader mgl.IShader, appState IAppState, animationState IAnimationState,
		cursorRect *mmath.MRect, windowSize *mmath.MVec2)
	Hash() string
}
