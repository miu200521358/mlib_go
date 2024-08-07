//go:build windows
// +build windows

package animation

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
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

func (animationState *AnimationState) DeformPhysics(physics mbt.IPhysics, appState state.IAppState) {
	if animationState.model == nil || (!appState.IsEnabledPhysics() && !appState.IsPhysicsReset()) {
		return
	}

	for _, rigidBody := range animationState.model.RigidBodies.Data {
		// 現在のボーン変形情報を保持
		rigidBodyBone := rigidBody.Bone
		if rigidBodyBone == nil {
			rigidBodyBone = rigidBody.JointedBone
		}

		if rigidBodyBone == nil || animationState.vmdDeltas.Bones.Get(rigidBodyBone.Index()) == nil {
			continue
		}

		if (appState.IsEnabledPhysics() && rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC) ||
			appState.IsPhysicsReset() {
			// 通常はボーン追従剛体・物理＋ボーン剛体だけ。物理リセット時は全部更新
			physics.UpdateTransform(animationState.ModelIndex(), rigidBodyBone,
				animationState.vmdDeltas.Bones.Get(rigidBodyBone.Index()).FilledGlobalMatrix(), rigidBody)
		}
	}
}

func (animationState *AnimationState) DeformAfterPhysics(physics mbt.IPhysics, appState state.IAppState) {
	if animationState.model != nil && appState.IsEnabledPhysics() && !appState.IsPhysicsReset() {
		// 物理剛体位置を更新
		for _, isAfterPhysics := range []bool{false, true} {
			for _, bone := range animationState.model.Bones.LayerSortedBones[isAfterPhysics] {
				if bone.Extend.RigidBody == nil || bone.Extend.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
					continue
				}
				bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(
					animationState.ModelIndex(), bone.Extend.RigidBody)
				if animationState.vmdDeltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
					bd := delta.NewBoneDeltaByGlobalMatrix(bone, appState.Frame(),
						bonePhysicsGlobalMatrix, animationState.vmdDeltas.Bones.Get(bone.ParentIndex))
					animationState.vmdDeltas.Bones.Update(bd)
				}
			}
		}
	}

	// 物理後のデフォーム情報
	animationState.vmdDeltas = deform.DeformBoneByPhysicsFlag(animationState.model,
		animationState.motion, animationState.vmdDeltas, true, appState.Frame(), nil, true)
}

func Deform(
	physics mbt.IPhysics, animationStates []state.IAnimationState, appState state.IAppState, timeStep float32,
) {
	// 物理デフォーム
	{
		var wg sync.WaitGroup

		for i := range animationStates {
			if animationStates[i] == nil || animationStates[i].Model() == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii].DeformPhysics(physics, appState)
			}(i)
		}

		wg.Wait()
	}

	if appState.IsPhysicsReset() {
		physics.UpdateFlags(true)
	}

	if appState.IsEnabledPhysics() || appState.IsPhysicsReset() {
		// 物理更新
		physics.StepSimulation(timeStep)
	}

	// 物理後デフォーム
	{
		var wg sync.WaitGroup

		for i := range animationStates {
			if animationStates[i] == nil || animationStates[i].Model() == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii].DeformAfterPhysics(physics, appState)
			}(i)
		}

		wg.Wait()
	}
}
