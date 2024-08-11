package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/miter"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

func deformBeforePhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion, vmdDeltas *delta.VmdDeltas,
) *delta.VmdDeltas {
	frame := appState.Frame()

	if vmdDeltas == nil || vmdDeltas.Frame() != frame ||
		vmdDeltas.ModelHash() != model.Hash() || vmdDeltas.MotionHash() != motion.Hash() {
		vmdDeltas = delta.NewVmdDeltas(frame, model.Bones, model.Hash(), motion.Hash())
		vmdDeltas.Morphs = DeformMorph(model, motion.MorphFrames, frame, nil)
		vmdDeltas = DeformBoneByPhysicsFlag(model, motion, vmdDeltas, true, frame, nil, false)
	}

	return vmdDeltas
}

func DeformPhysicsByBone(
	appState state.IAppState, model *pmx.PmxModel, vmdDeltas *delta.VmdDeltas, physics *mbt.MPhysics,
) {
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

func DeformBonePyPhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion,
	vmdDeltas *delta.VmdDeltas, physics *mbt.MPhysics,
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
	physics *mbt.MPhysics, appState state.IAppState, timeStep float32,
	models []*pmx.PmxModel, motions []*vmd.VmdMotion, vmdDeltas []*delta.VmdDeltas,
) []*delta.VmdDeltas {
	// 物理前デフォーム
	for i := range models {
		if models[i] == nil || motions[i] == nil {
			continue
		}
		for i >= len(vmdDeltas) {
			vmdDeltas = append(vmdDeltas, nil)
		}
		vmdDeltas[i] = deformBeforePhysics(appState, models[i], motions[i], vmdDeltas[i])
	}

	return vmdDeltas
}

func DeformPhysics(
	physics *mbt.MPhysics, appState state.IAppState, timeStep float32,
	models []*pmx.PmxModel, motions []*vmd.VmdMotion, vmdDeltas []*delta.VmdDeltas,
) []*delta.VmdDeltas {
	// 物理デフォーム
	for i := range models {
		if models[i] == nil || vmdDeltas[i] == nil {
			continue
		}
		DeformPhysicsByBone(appState, models[i], vmdDeltas[i], physics)
	}

	if appState.IsEnabledPhysics() || appState.IsPhysicsReset() {
		// 物理更新
		physics.StepSimulation(timeStep)
	}

	for i := range models {
		if models[i] == nil || motions[i] == nil || vmdDeltas[i] == nil {
			continue
		}
		vmdDeltas[i] = DeformBonePyPhysics(appState, models[i], motions[i], vmdDeltas[i], physics)
	}

	return vmdDeltas
}
