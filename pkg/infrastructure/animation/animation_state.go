//go:build windows
// +build windows

package animation

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

type AnimationState struct {
	windowIndex  int                 // ウィンドウインデックス
	modelIndex   int                 // モデルインデックス
	frame        float64             // フレーム
	renderModel  *RenderModel        // 描画モデル
	model        *pmx.PmxModel       // モデル
	motion       *vmd.VmdMotion      // モーション
	vmdDeltas    *delta.VmdDeltas    // モーション変化量
	renderDeltas *delta.RenderDeltas // 描画変化量
}

func (a *AnimationState) WindowIndex() int {
	return a.windowIndex
}

func (a *AnimationState) SetWindowIndex(index int) {
	a.windowIndex = index
}

func (a *AnimationState) ModelIndex() int {
	return a.modelIndex
}

func (a *AnimationState) SetModelIndex(index int) {
	a.modelIndex = index
}

func (a *AnimationState) Frame() float64 {
	return a.frame
}

func (a *AnimationState) SetFrame(frame float64) {
	a.frame = frame
}

func (a *AnimationState) Model() *pmx.PmxModel {
	return a.model
}

func (a *AnimationState) SetModel(model *pmx.PmxModel) {
	a.model = model
}

func (a *AnimationState) RenderModel() state.IRenderModel {
	return a.renderModel
}

func (a *AnimationState) SetRenderModel(renderModel state.IRenderModel) {
	a.renderModel = renderModel.(*RenderModel)
}

func (a *AnimationState) Motion() *vmd.VmdMotion {
	return a.motion
}

func (a *AnimationState) SetMotion(motion *vmd.VmdMotion) {
	a.motion = motion
}

func (a *AnimationState) VmdDeltas() *delta.VmdDeltas {
	return a.vmdDeltas
}

func (a *AnimationState) SetVmdDeltas(deltas *delta.VmdDeltas) {
	a.vmdDeltas = deltas
}

func (a *AnimationState) RenderDeltas() *delta.RenderDeltas {
	return a.renderDeltas
}

func (a *AnimationState) SetRenderDeltas(deltas *delta.RenderDeltas) {
	a.renderDeltas = deltas
}

func (a *AnimationState) Load(model *pmx.PmxModel) {
	if a.renderModel == nil || a.renderModel.Hash() != model.Hash() {
		if a.renderModel != nil {
			a.renderModel.Delete()
		}
		a.renderModel = NewRenderModel(a.windowIndex, model)
		a.model = model
	}
}

func (animationState *AnimationState) Render(shader mgl.IShader, appState state.IAppState) {
	if animationState.renderModel != nil && animationState.model != nil {
		animationState.renderModel.Render(shader, appState, animationState)
	}
}

func NewAnimationState(windowIndex, modelIndex int) *AnimationState {
	return &AnimationState{
		windowIndex:  windowIndex,
		modelIndex:   modelIndex,
		frame:        -1,
		renderDeltas: delta.NewRenderDeltas(),
	}
}
