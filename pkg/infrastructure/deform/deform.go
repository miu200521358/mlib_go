package deform

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/miter"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

func DeformBeforePhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion,
) *delta.VmdDeltas {
	frame := int(appState.Frame())

	vmdDeltas := delta.NewVmdDeltas(model.Materials, model.Bones)
	vmdDeltas.Morphs = DeformMorph(model, motion.MorphFrames, float32(frame), nil)
	vmdDeltas = DeformBoneByPhysicsFlag(model, motion, vmdDeltas, true, float32(frame), nil, false)

	return vmdDeltas
}

func DeformPhysics(appState state.IAppState, model *pmx.PmxModel, vmdDeltas *delta.VmdDeltas, physics mbt.IPhysics) {
	// 物理剛体位置を更新
	processFunc := func(i int) {
		rigidBody := model.RigidBodies.Get(i)

		// 現在のボーン変形情報を保持
		rigidBodyBone := rigidBody.Bone
		if rigidBodyBone == nil {
			rigidBodyBone = rigidBody.JointedBone
		}

		if rigidBodyBone == nil || vmdDeltas.Bones.Get(rigidBodyBone.Index()) == nil {
			return
		}

		if (appState.IsEnabledPhysics() && rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC) ||
			appState.IsPhysicsReset() {
			// 通常はボーン追従剛体・物理＋ボーン剛体だけ。物理リセット時は全部更新
			physics.UpdateTransform(model.Index(), rigidBodyBone,
				vmdDeltas.Bones.Get(rigidBodyBone.Index()).FilledGlobalMatrix(), rigidBody)
		}
	}

	// 100件ずつ処理
	miter.IterParallelByCount(model.RigidBodies.Len(), 100, processFunc)
}

func DeformAfterPhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion,
	vmdDeltas *delta.VmdDeltas, physics mbt.IPhysics,
) *delta.VmdDeltas {
	if model != nil && appState.IsEnabledPhysics() && !appState.IsPhysicsReset() {
		// 物理剛体位置を更新
		processFunc := func(i int) {
			bone := model.Bones.Get(i)
			if bone.Extend.RigidBody == nil || bone.Extend.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
				return
			}
			bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(model.Index(), bone.Extend.RigidBody)
			if vmdDeltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
				bd := delta.NewBoneDeltaByGlobalMatrix(bone, appState.Frame(),
					bonePhysicsGlobalMatrix, vmdDeltas.Bones.Get(bone.ParentIndex))
				vmdDeltas.Bones.Update(bd)
			}
		}

		// 100件ずつ処理
		miter.IterParallelByList(model.Bones.LayerSortedIndexes, 100, processFunc)
	}

	// 物理後のデフォーム情報
	return DeformBoneByPhysicsFlag(model, motion, vmdDeltas, true, appState.Frame(), nil, true)
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
				DeformPhysics(appState, animationStates[ii].Model(), animationStates[ii].VmdDeltas(), physics)
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
				animationStates[ii].SetVmdDeltas(
					DeformAfterPhysics(appState, animationStates[ii].Model(), animationStates[ii].Motion(),
						animationStates[ii].VmdDeltas(), physics))
			}(i)
		}

		wg.Wait()
	}
}
