//go:build windows
// +build windows

package animation

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

type AnimationState struct {
	windowIndex              int                 // ウィンドウインデックス
	modelIndex               int                 // モデルインデックス
	frame                    float64             // フレーム
	renderModel              *render.RenderModel // 描画モデル
	model                    *pmx.PmxModel       // モデル
	motion                   *vmd.VmdMotion      // モーション
	invisibleMaterialIndexes map[int]struct{}    // 非表示材質インデックス
	selectedVertexIndexes    map[int]struct{}    // 選択頂点インデックス
	noSelectedVertexIndexes  map[int]struct{}    // 非選択頂点インデックス
	vmdDeltas                *delta.VmdDeltas    // モーション変化量
}

func (animationState *AnimationState) WindowIndex() int {
	return animationState.windowIndex
}

func (animationState *AnimationState) SetWindowIndex(index int) {
	animationState.windowIndex = index
}

func (animationState *AnimationState) ModelIndex() int {
	return animationState.modelIndex
}

func (animationState *AnimationState) SetModelIndex(index int) {
	animationState.modelIndex = index
}

func (animationState *AnimationState) Frame() float64 {
	return animationState.frame
}

func (animationState *AnimationState) SetFrame(frame float64) {
	animationState.frame = frame
}

func (animationState *AnimationState) Model() *pmx.PmxModel {
	return animationState.model
}

func (animationState *AnimationState) SetModel(model *pmx.PmxModel) {
	animationState.model = model
}

func (animationState *AnimationState) RenderModel() state.IRenderModel {
	return animationState.renderModel
}

func (animationState *AnimationState) SetRenderModel(renderModel state.IRenderModel) {
	animationState.renderModel = renderModel.(*render.RenderModel)
}

func (animationState *AnimationState) Motion() *vmd.VmdMotion {
	return animationState.motion
}

func (animationState *AnimationState) SetMotion(motion *vmd.VmdMotion) {
	animationState.motion = motion
}

func (animationState *AnimationState) VmdDeltas() *delta.VmdDeltas {
	return animationState.vmdDeltas
}

func (animationState *AnimationState) SetVmdDeltas(deltas *delta.VmdDeltas) {
	animationState.vmdDeltas = deltas
}

func (animationState *AnimationState) InvisibleMaterialIndexes() []int {
	if animationState.invisibleMaterialIndexes == nil {
		return nil
	}

	indexes := make([]int, 0, len(animationState.invisibleMaterialIndexes))
	for i := range animationState.invisibleMaterialIndexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (animationState *AnimationState) SetInvisibleMaterialIndexes(indexes []int) {
	if len(indexes) == 0 {
		animationState.invisibleMaterialIndexes = nil
		return
	}
	animationState.invisibleMaterialIndexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		animationState.invisibleMaterialIndexes[i] = struct{}{}
	}
}

func (animationState *AnimationState) ExistInvisibleMaterialIndex(index int) bool {
	_, ok := animationState.invisibleMaterialIndexes[index]
	return ok
}

func (animationState *AnimationState) SelectedVertexIndexes() []int {
	indexes := make([]int, 0, len(animationState.selectedVertexIndexes))
	for i := range animationState.selectedVertexIndexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (animationState *AnimationState) NoSelectedVertexIndexes() []int {
	indexes := make([]int, 0, len(animationState.noSelectedVertexIndexes))
	for i := range animationState.noSelectedVertexIndexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (animationState *AnimationState) SetSelectedVertexIndexes(indexes []int) {
	animationState.selectedVertexIndexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		animationState.selectedVertexIndexes[i] = struct{}{}
	}
}

func (animationState *AnimationState) SetNoSelectedVertexIndexes(indexes []int) {
	animationState.noSelectedVertexIndexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		animationState.noSelectedVertexIndexes[i] = struct{}{}
	}
}

func (animationState *AnimationState) ClearSelectedVertexIndexes() {
	animationState.UpdateNoSelectedVertexIndexes(animationState.SelectedVertexIndexes())
	animationState.selectedVertexIndexes = make(map[int]struct{})
}

func (animationState *AnimationState) UpdateSelectedVertexIndexes(indexes []int) {
	for _, index := range indexes {
		animationState.selectedVertexIndexes[index] = struct{}{}
	}
}

func (animationState *AnimationState) UpdateNoSelectedVertexIndexes(indexes []int) {
	if animationState.noSelectedVertexIndexes == nil {
		animationState.noSelectedVertexIndexes = make(map[int]struct{}, len(indexes))
	}
	for _, index := range indexes {
		animationState.noSelectedVertexIndexes[index] = struct{}{}
		delete(animationState.selectedVertexIndexes, index)
	}
}

func (animationState *AnimationState) Load(model *pmx.PmxModel) {
	if animationState.renderModel == nil || animationState.renderModel.Hash() != model.Hash() {
		if animationState.renderModel != nil {
			animationState.renderModel.Delete()
		}
		animationState.renderModel = render.NewRenderModel(animationState.windowIndex, model)
		animationState.model = model
	}
}

func (animationState *AnimationState) Render(
	shader mgl.IShader, appState state.IAppState, leftCursorPositions, leftCursorRemovePositions, leftCursorWorldHistoryPositions, leftCursorRemoveWorldHistoryPositions []*mgl32.Vec3,
) {
	if !appState.IsShowSelectedVertex() {
		animationState.ClearSelectedVertexIndexes()
	}
	if animationState.renderModel != nil && animationState.model != nil {
		animationState.renderModel.Render(shader, appState, animationState,
			leftCursorPositions, leftCursorRemovePositions,
			leftCursorWorldHistoryPositions, leftCursorRemoveWorldHistoryPositions)
	}
}

func NewAnimationState(windowIndex, modelIndex int) *AnimationState {
	return &AnimationState{
		windowIndex:              windowIndex,
		modelIndex:               modelIndex,
		frame:                    0,
		invisibleMaterialIndexes: nil,
		selectedVertexIndexes:    nil,
		noSelectedVertexIndexes:  nil,
	}
}
