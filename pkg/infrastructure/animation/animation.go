//go:build windows
// +build windows

package animation

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

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

func (animationState *AnimationState) DeformBeforePhysics(
	appState state.IAppState, model *pmx.PmxModel,
) (*delta.VmdDeltas, *delta.RenderDeltas) {
	frame := int(appState.Frame())

	vmdDeltas := delta.NewVmdDeltas(model.Materials, model.Bones)
	vmdDeltas.Morphs = deform.DeformMorph(model, animationState.motion.MorphFrames, frame, nil)
	vmdDeltas = deform.DeformBoneByPhysicsFlag(model, animationState.motion, vmdDeltas, true, frame, nil, false)

	renderDeltas := delta.NewRenderDeltas()
	renderDeltas.VertexMorphDeltaIndexes, renderDeltas.VertexMorphDeltas =
		newVertexMorphDeltasGl(vmdDeltas.Morphs.Vertices)

	renderDeltas.MeshDeltas = make([]*delta.MeshDelta, len(model.Materials.Data))
	for i, md := range vmdDeltas.Morphs.Materials.Data {
		renderDeltas.MeshDeltas[i] = delta.NewMeshDelta(md)
	}

	return vmdDeltas, renderDeltas
}

func (animationState *AnimationState) DeformPhysics(physics mbt.IPhysics, appState state.IAppState) {
	if animationState.model == nil || !appState.IsEnabledPhysics() {
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

		if rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC {
			// ボーン追従剛体・物理＋ボーン位置もしくは強制更新の場合のみ剛体位置更新
			physics.UpdateTransform(animationState.ModelIndex(), rigidBodyBone,
				animationState.vmdDeltas.Bones.Get(rigidBodyBone.Index()).FilledGlobalMatrix(), rigidBody)
		}
	}
}

func (animationState *AnimationState) DeformAfterPhysics(physics mbt.IPhysics, appState state.IAppState) {
	if animationState.model == nil || !appState.IsEnabledPhysics() {
		return
	}

	// 物理剛体位置を更新
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

	// 物理後のデフォーム情報
	animationState.vmdDeltas = deform.DeformBoneByPhysicsFlag(animationState.model,
		animationState.motion, animationState.vmdDeltas, true, int(appState.Frame()), nil, true)

	// // 選択頂点モーフの設定は常に更新する
	// SelectedVertexIndexesDeltas, SelectedVertexGlDeltas := animation.SelectedVertexMorphDeltasGL(
	// 	SelectedVertexDeltas, model, selectedVertexIndexes, nextSelectedVertexIndexes)
}