//go:build windows
// +build windows

package renderer

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/interface/core"
)

func Animate(
	physics mbt.IPhysics, animationStates []core.IAnimationState, appState core.IAppState, timeStep float32,
) {
	// モデルの描画モデルが未設定の場合は設定
	for i := range animationStates {
		if animationStates[i].Model() != nil && animationStates[i].RenderModel() == nil {
			animationStates[i].SetRenderModel(
				NewRenderModel(animationStates[i].WindowIndex(), animationStates[i].Model()))
			physics.AddModel(animationStates[i].ModelIndex(), animationStates[i].Model())
		}
	}

	// 物理デフォーム
	{
		var wg sync.WaitGroup

		for i := range animationStates {
			if animationStates[i].Model() == nil || animationStates[i].RenderModel() == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii].AnimatePhysics(physics, appState)
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
			if animationStates[i].RenderModel() == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii].AnimateAfterPhysics(physics, appState)
			}(i)
		}

		wg.Wait()
	}
}

func (animationState *AnimationState) AnimateBeforePhysics(
	appState core.IAppState, model *pmx.PmxModel,
) (*delta.VmdDeltas, *delta.RenderDeltas) {
	if animationState.motion == nil {
		animationState.motion = vmd.NewVmdMotion("")
	}

	frame := int(appState.Frame())

	vmdDeltas := delta.NewVmdDeltas(model.Materials, model.Bones)
	vmdDeltas.Morphs = deform.DeformMorph(
		model, animationState.motion.MorphFrames, frame, nil)
	vmdDeltas = deform.DeformBoneByPhysicsFlag(model,
		animationState.motion, vmdDeltas, true, frame, nil, false)

	renderDeltas := delta.NewRenderDeltas()
	renderDeltas.VertexMorphDeltaIndexes, renderDeltas.VertexMorphDeltas =
		newVertexMorphDeltasGl(vmdDeltas.Morphs.Vertices)

	renderDeltas.MeshDeltas = make([]*delta.MeshDelta, len(model.Materials.Data))
	for i, md := range vmdDeltas.Morphs.Materials.Data {
		renderDeltas.MeshDeltas[i] = delta.NewMeshDelta(md)
	}

	return vmdDeltas, renderDeltas
}

func (animationState *AnimationState) AnimatePhysics(physics mbt.IPhysics, appState core.IAppState) {
	if appState.IsEnabledPhysics() && animationState.model != nil && animationState.vmdDeltas != nil {
		for _, rigidBody := range animationState.model.RigidBodies.Data {
			// 現在のボーン変形情報を保持
			rigidBodyBone := rigidBody.Bone
			if rigidBodyBone == nil {
				rigidBodyBone = rigidBody.JointedBone
			}
			if rigidBodyBone == nil || animationState.vmdDeltas.Bones.Get(rigidBodyBone.Index()) == nil {
				continue
			}

			if rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC {
				// ボーン追従剛体・物理＋ボーン位置もしくは強制更新の場合のみ剛体位置更新
				physics.UpdateTransform(animationState.ModelIndex(), rigidBodyBone,
					animationState.vmdDeltas.Bones.Get(rigidBodyBone.Index()).FilledGlobalMatrix(), rigidBody)
			}
		}
	}
}

func (animationState *AnimationState) AnimateAfterPhysics(physics mbt.IPhysics, appState core.IAppState) {
	// 物理剛体位置を更新
	if (appState.IsEnabledPhysics() || appState.IsPhysicsReset()) &&
		animationState.model != nil && animationState.vmdDeltas != nil {
		for _, isAfterPhysics := range []bool{false, true} {
			for _, bone := range animationState.model.Bones.LayerSortedBones[isAfterPhysics] {
				if bone.Extend.RigidBody == nil || bone.Extend.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
					continue
				}
				bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(
					animationState.ModelIndex(), bone.Extend.RigidBody)
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
}
