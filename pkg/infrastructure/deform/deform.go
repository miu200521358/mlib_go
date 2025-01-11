package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/miter"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

func DeformModel(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	frame int,
) *pmx.PmxModel {
	vmdDeltas := delta.NewVmdDeltas(float32(frame), model.Bones, "", "")
	vmdDeltas.Morphs = DeformMorph(model, motion.MorphFrames, float32(frame), nil)
	vmdDeltas = DeformBoneByPhysicsFlag(model, motion, vmdDeltas, false, float32(frame), nil, false)

	// 頂点にボーン変形を適用
	for _, vertex := range model.Vertices.Data {
		mat := &mmath.MMat4{}
		for j := range vertex.Deform.AllIndexes() {
			boneIndex := vertex.Deform.AllIndexes()[j]
			weight := vertex.Deform.AllWeights()[j]
			mat.Add(vmdDeltas.Bones.Get(boneIndex).FilledLocalMatrix().MuledScalar(weight))
		}

		var morphDelta *delta.VertexMorphDelta
		if vmdDeltas.Morphs != nil && vmdDeltas.Morphs.Vertices != nil {
			morphDelta = vmdDeltas.Morphs.Vertices.Get(vertex.Index())
		}

		// 頂点変形
		if morphDelta == nil {
			vertex.Position = mat.MulVec3(vertex.Position)
		} else {
			vertex.Position = mat.MulVec3(vertex.Position.Added(morphDelta.Position))
		}

		// 法線変形
		vertex.Normal = mat.MulVec3(vertex.Normal).Normalized()

		// SDEFの場合、パラメーターを再計算
		switch sdef := vertex.Deform.(type) {
		case *pmx.Sdef:
			// SDEF-C: ボーンのベクトルと頂点の交点
			sdef.SdefC = mmath.IntersectLinePoint(
				vmdDeltas.Bones.Get(sdef.AllIndexes()[0]).GlobalPosition,
				vmdDeltas.Bones.Get(sdef.AllIndexes()[1]).GlobalPosition,
				vertex.Position,
			)

			// SDEF-R0: 0番目のボーンとSDEF-Cの中点
			sdef.SdefR0 = vmdDeltas.Bones.Get(sdef.AllIndexes()[0]).GlobalPosition.Added(sdef.SdefC).MuledScalar(0.5)

			// SDEF-R1: 1番目のボーンとSDEF-Cの中点
			sdef.SdefR1 = vmdDeltas.Bones.Get(sdef.AllIndexes()[1]).GlobalPosition.Added(sdef.SdefC).MuledScalar(0.5)
		}
	}

	// ボーンの位置を更新
	for i, bone := range model.Bones.Data {
		if vmdDeltas.Bones.Get(i) != nil {
			bone.Position = vmdDeltas.Bones.Get(i).FilledGlobalPosition()
		}
	}

	return model
}

func DeformIk(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	frame float32,
	ikBone *pmx.Bone,
	ikGlobalPosition *mmath.MVec3,
	boneNames []string,
) *delta.VmdDeltas {
	if boneNames == nil {
		boneNames = make([]string, 0)
	}
	boneNames = append(boneNames, model.Bones.Get(ikBone.Ik.BoneIndex).Name())
	for _, link := range ikBone.Ik.Links {
		boneNames = append(boneNames, model.Bones.Get(link.BoneIndex).Name())
	}

	deformBoneIndexes, deltas := newVmdDeltas(model, motion, deltas, frame, boneNames, false)

	deformIk(model, motion, deltas, frame, false, ikBone, ikGlobalPosition, deformBoneIndexes, 0)

	updateGlobalMatrix(deltas.Bones, deformBoneIndexes)

	return deltas
}

// DeformBone 前回情報なしでボーンデフォーム処理を実行する
func DeformBone(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	isCalcIk bool,
	frame int,
	boneNames []string,
) *delta.BoneDeltas {
	return DeformBoneByPhysicsFlag(model, motion, nil, isCalcIk, float32(frame), boneNames, false).Bones
}

func deformBeforePhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion, vmdDeltas *delta.VmdDeltas,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return vmdDeltas
	}

	frame := appState.Frame()

	if vmdDeltas == nil || vmdDeltas.Frame() != frame ||
		vmdDeltas.ModelHash() != model.Hash() || vmdDeltas.MotionHash() != motion.Hash() {
		deltas := delta.NewVmdDeltas(frame, model.Bones, model.Hash(), motion.Hash())
		deltas.Morphs = DeformMorph(model, motion.MorphFrames, frame, nil)

		// ボーンデフォーム情報を埋める(物理前後全部埋める)
		deltas.Bones = fillBoneDeform(model, motion, deltas, frame, model.Bones.LayerSortedIndexes, true, false)

		// ボーンデフォーム情報を更新する
		updateGlobalMatrix(deltas.Bones, model.Bones.LayerSortedIndexes)

		return deltas
	}

	return vmdDeltas
}

func DeformPhysicsByBone(
	appState state.IAppState, model *pmx.PmxModel, vmdDeltas *delta.VmdDeltas, physics *mbt.MPhysics,
) error {
	// 物理剛体位置を更新
	if err := miter.IterParallelByCount(model.RigidBodies.Len(), 100, func(i int) {
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
	}); err != nil {
		return err
	}

	return nil
}

func DeformBonePyPhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion,
	vmdDeltas *delta.VmdDeltas, physics *mbt.MPhysics,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return vmdDeltas
	}

	if appState.IsEnabledPhysics() && !appState.IsPhysicsReset() {
		// 物理剛体位置を更新
		for _, boneIndex := range model.Bones.LayerSortedIndexes {
			bone := model.Bones.Get(boneIndex)
			if bone.Extend.RigidBody == nil || bone.Extend.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
				continue
			}
			bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(model.Index(), bone.Extend.RigidBody)
			if vmdDeltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
				bd := delta.NewBoneDeltaByGlobalMatrix(bone, appState.Frame(),
					bonePhysicsGlobalMatrix, vmdDeltas.Bones.Get(bone.ParentIndex))
				vmdDeltas.Bones.Update(bd)
			}
		}
	}

	// ボーンデフォーム情報を埋める(物理後埋める)
	vmdDeltas.Bones = fillBoneDeform(model, motion, vmdDeltas, appState.Frame(),
		model.Bones.LayerSortedBoneIndexes[true], true, true)

	// ボーンデフォーム情報を更新する
	updateGlobalMatrix(vmdDeltas.Bones, model.Bones.LayerSortedBoneIndexes[true])

	return vmdDeltas
}

func DeformForReset(
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
		// 物理前
		vmdDeltas[i] = DeformBoneByPhysicsFlag(models[i], motions[i], vmdDeltas[i], true, appState.Frame(), nil, false)
		// 物理後
		vmdDeltas[i] = DeformBoneByPhysicsFlag(models[i], motions[i], vmdDeltas[i], true, appState.Frame(), nil, true)
	}

	return vmdDeltas
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
) ([]*delta.VmdDeltas, error) {
	// 物理デフォーム
	for i := range models {
		for i >= len(vmdDeltas) {
			vmdDeltas = append(vmdDeltas, nil)
		}
		if models[i] == nil || vmdDeltas[i] == nil {
			continue
		}
		if err := DeformPhysicsByBone(appState, models[i], vmdDeltas[i], physics); err != nil {
			return vmdDeltas, err
		}
	}

	if appState.IsEnabledPhysics() || appState.IsPhysicsReset() {
		// 物理更新
		physics.StepSimulation(timeStep)
	}

	for i := range models {
		if models[i] == nil || motions[i] == nil {
			continue
		}
		vmdDeltas[i] = DeformBonePyPhysics(appState, models[i], motions[i], vmdDeltas[i], physics)
	}

	return vmdDeltas, nil
}
