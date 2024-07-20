//go:build windows
// +build windows

package renderer

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type AnimationState struct {
	windowIndex              int              // ウィンドウインデックス
	modelIndex               int              // モデルインデックス
	frame                    float64          // フレーム
	renderModel              *RenderModel     // 描画モデル
	model                    *pmx.PmxModel    // モデル
	motion                   *vmd.VmdMotion   // モーション
	vmdDeltas                *delta.VmdDeltas // モーション変化量
	invisibleMaterialIndexes []int            // 非表示材質インデックス
	selectedVertexIndexes    []int            // 選択頂点インデックス
	vertexMorphDeltaIndexes  []int            // 頂点モーフインデックス
	vertexMorphDeltas        [][]float32      // 頂点モーフデルタ
	meshDeltas               []*MeshDelta     // メッシュデルタ
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

func (a *AnimationState) RenderModel() *RenderModel {
	return a.renderModel
}

func (a *AnimationState) SetRenderModel(renderModel *RenderModel) {
	a.renderModel = renderModel
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

func (a *AnimationState) Load() {
	a.renderModel = NewRenderModel(a.windowIndex, a.model)
}

func NewAnimationState(windowIndex, modelIndex int) *AnimationState {
	return &AnimationState{
		windowIndex:              windowIndex,
		modelIndex:               modelIndex,
		frame:                    -1,
		invisibleMaterialIndexes: make([]int, 0),
		selectedVertexIndexes:    make([]int, 0),
	}
}

func Animate(
	physics *mbt.MPhysics, animationStates []*AnimationState, appState state.IAppState, timeStep float32,
) []*AnimationState {
	// 物理前デフォーム
	{
		var wg sync.WaitGroup

		for i := range animationStates {
			if animationStates[i].model == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii] = animateBeforePhysics(physics, animationStates[ii], appState)
			}(i)
		}

		wg.Wait()
	}

	if appState.IsEnabledPhysics() || appState.IsPhysicsReset() {
		// 物理更新
		physics.StepSimulation(timeStep)
	}

	// 物理後デフォーム
	{
		var wg sync.WaitGroup

		for i := range animationStates {
			if animationStates[i].renderModel == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii] = animateAfterPhysics(physics, animationStates[ii], appState)
			}(i)
		}

		wg.Wait()
	}

	return animationStates
}

func animateBeforePhysics(
	physics *mbt.MPhysics, animationState *AnimationState, appState state.IAppState,
) *AnimationState {
	if animationState.motion == nil {
		animationState.motion = vmd.NewVmdMotion("")
	}

	deltas := delta.NewVmdDeltas(animationState.model.Materials, animationState.model.Bones)

	if animationState.vmdDeltas == nil || animationState.frame != appState.Frame() {
		frame := int(appState.Frame())

		deltas.Morphs = deform.DeformMorph(animationState.model, animationState.motion.MorphFrames, frame, nil)
		deltas = deform.DeformBoneByPhysicsFlag(animationState.model,
			animationState.motion, deltas, true, frame, nil, false)

		animationState.vertexMorphDeltaIndexes, animationState.vertexMorphDeltas =
			newVertexMorphDeltasGl(deltas.Morphs.Vertices)

		animationState.meshDeltas = make([]*MeshDelta, len(animationState.model.Materials.Data))
		for i, md := range deltas.Morphs.Materials.Data {
			animationState.meshDeltas[i] = newMeshDelta(md)
		}

		animationState.frame = appState.Frame()
	} else {
		deltas.Morphs = animationState.vmdDeltas.Morphs
		deltas.Bones = animationState.vmdDeltas.Bones
	}

	modelIndex := 0

	if appState.IsEnabledPhysics() {
		for _, rigidBody := range animationState.model.RigidBodies.Data {
			// 現在のボーン変形情報を保持
			rigidBodyBone := rigidBody.Bone
			if rigidBodyBone == nil {
				rigidBodyBone = rigidBody.JointedBone
			}
			if rigidBodyBone == nil || deltas.Bones.Get(rigidBodyBone.Index) == nil {
				continue
			}

			if rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC {
				// ボーン追従剛体・物理＋ボーン位置もしくは強制更新の場合のみ剛体位置更新
				boneTransform := bt.NewBtTransform()
				defer bt.DeleteBtTransform(boneTransform)
				mat := mgl.NewGlMat4(deltas.Bones.Get(rigidBodyBone.Index).FilledGlobalMatrix())
				boneTransform.SetFromOpenGLMatrix(&mat[0])

				physics.UpdateTransform(modelIndex, rigidBodyBone, boneTransform, rigidBody)
			}
		}
	}

	animationState.vmdDeltas = deltas

	return animationState
}

func animateAfterPhysics(
	physics *mbt.MPhysics, animationState *AnimationState, appState state.IAppState,
) *AnimationState {
	modelIndex := 0

	// 物理剛体位置を更新
	if appState.IsEnabledPhysics() || appState.IsPhysicsReset() {
		for _, isAfterPhysics := range []bool{false, true} {
			for _, bone := range animationState.model.Bones.LayerSortedBones[isAfterPhysics] {
				if bone.Extend.RigidBody == nil || bone.Extend.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
					continue
				}
				bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(modelIndex, bone.Extend.RigidBody)
				if animationState.vmdDeltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
					bd := delta.NewBoneDeltaByGlobalMatrix(bone, int(appState.Frame()),
						bonePhysicsGlobalMatrix, animationState.vmdDeltas.Bones.Get(bone.ParentIndex))
					animationState.vmdDeltas.Bones.Update(bd)
				}
			}
		}
	}

	// 物理後のデフォーム情報
	animationState.vmdDeltas = deform.DeformBoneByPhysicsFlag(animationState.model,
		animationState.motion, animationState.vmdDeltas, true, int(appState.Frame()), nil, true)

	// // 選択頂点モーフの設定は常に更新する
	// SelectedVertexIndexesDeltas, SelectedVertexGlDeltas := renderer.SelectedVertexMorphDeltasGL(
	// 	SelectedVertexDeltas, model, selectedVertexIndexes, nextSelectedVertexIndexes)

	return animationState
}
