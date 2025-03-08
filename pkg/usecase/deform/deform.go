package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/miter"
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
	for v := range model.Vertices.Iterator() {
		vertex := v.Value
		mat := &mmath.MMat4{}
		for j := range vertex.Deform.Indexes() {
			boneIndex := vertex.Deform.Indexes()[j]
			weight := vertex.Deform.Weights()[j]
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
				vmdDeltas.Bones.Get(sdef.Indexes()[0]).GlobalPosition,
				vmdDeltas.Bones.Get(sdef.Indexes()[1]).GlobalPosition,
				vertex.Position,
			)

			// SDEF-R0: 0番目のボーンとSDEF-Cの中点
			sdef.SdefR0 = vmdDeltas.Bones.Get(sdef.Indexes()[0]).GlobalPosition.Added(sdef.SdefC).MuledScalar(0.5)

			// SDEF-R1: 1番目のボーンとSDEF-Cの中点
			sdef.SdefR1 = vmdDeltas.Bones.Get(sdef.Indexes()[1]).GlobalPosition.Added(sdef.SdefC).MuledScalar(0.5)
		}
	}

	// ボーンの位置を更新
	for b := range model.Bones.Iterator() {
		bone := b.Value
		if vmdDeltas.Bones.Get(b.Index) != nil {
			bone.Position = vmdDeltas.Bones.Get(b.Index).FilledGlobalPosition()
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
	ikTargetBone, _ := model.Bones.Get(ikBone.Ik.BoneIndex)
	boneNames = append(boneNames, ikTargetBone.Name())
	for _, link := range ikBone.Ik.Links {
		linkBone, _ := model.Bones.Get(link.BoneIndex)
		boneNames = append(boneNames, linkBone.Name())
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

// DeformBoneByPhysicsFlag ボーンデフォーム処理を実行する
func DeformBoneByPhysicsFlag(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	isCalcIk bool,
	frame float32,
	boneNames []string,
	isAfterPhysics bool,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return deltas
	}

	deformBoneIndexes, deltas := newVmdDeltas(model, motion, deltas, frame, boneNames, isAfterPhysics)

	// ボーンデフォーム情報を埋める
	deltas.Bones = fillBoneDeform(model, motion, deltas, frame, deformBoneIndexes, isCalcIk, isAfterPhysics)

	// ボーンデフォーム情報を更新する
	updateGlobalMatrix(deltas.Bones, deformBoneIndexes)

	return deltas
}

func Deform(
	shared *state.SharedState,
	physics *mbt.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	timeStep float32,
) *delta.VmdDeltas {
	// 物理前変形
	deltas = deformBeforePhysics(model, motion, deltas, shared.Frame())

	if shared.IsEnabledPhysics() || shared.IsPhysicsReset() {
		// 物理更新
		physics.StepSimulation(timeStep)
	}

	// 物理変形
	if err := deformPhysics(shared, physics, model, deltas); err != nil {
		return deltas
	}

	// 物理後変形
	deltas = deformAfterPhysics(shared, physics, model, motion, deltas)

	return deltas
}

// deformBeforePhysics 物理演算前のボーンデフォーム処理を実行する
func deformBeforePhysics(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	vmdDeltas *delta.VmdDeltas,
	frame float32,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return vmdDeltas
	}

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

func deformPhysics(
	shared *state.SharedState,
	physics *mbt.MPhysics,
	model *pmx.PmxModel,
	vmdDeltas *delta.VmdDeltas,
) error {
	// 物理剛体位置を更新
	if err := miter.IterParallelByCount(model.RigidBodies.Length(), 100, func(i int) {
		rigidBody, err := model.RigidBodies.Get(i)
		if err != nil {
			return
		}

		// 現在のボーン変形情報を保持
		rigidBodyBone := rigidBody.Bone
		if rigidBodyBone == nil {
			rigidBodyBone = rigidBody.JointedBone
		}

		if rigidBodyBone == nil || vmdDeltas.Bones.Get(rigidBodyBone.Index()) == nil {
			return
		}

		if (shared.IsEnabledPhysics() && rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC) ||
			shared.IsPhysicsReset() {
			// 通常はボーン追従剛体・物理＋ボーン剛体だけ。物理リセット時は全部更新
			physics.UpdateTransform(model.Index(), rigidBodyBone,
				vmdDeltas.Bones.Get(rigidBodyBone.Index()).FilledGlobalMatrix(), rigidBody)
		}
	}); err != nil {
		return err
	}

	return nil
}

func deformAfterPhysics(
	shared *state.SharedState,
	physics *mbt.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	vmdDeltas *delta.VmdDeltas,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return vmdDeltas
	}

	if shared.IsEnabledPhysics() && !shared.IsPhysicsReset() {
		// 物理剛体位置を更新
		for _, boneIndex := range model.Bones.LayerSortedIndexes {
			bone, err := model.Bones.Get(boneIndex)
			if err != nil || bone == nil || bone.RigidBody == nil || bone.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
				continue
			}
			bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(model.Index(), bone.RigidBody)
			if vmdDeltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
				bd := delta.NewBoneDeltaByGlobalMatrix(bone, shared.Frame(),
					bonePhysicsGlobalMatrix, vmdDeltas.Bones.Get(bone.ParentIndex))
				vmdDeltas.Bones.Update(bd)
			}
		}
	}

	// ボーンデフォーム情報を埋める(物理後埋める)
	vmdDeltas.Bones = fillBoneDeform(model, motion, vmdDeltas, shared.Frame(),
		model.Bones.LayerSortedBoneIndexes[true], true, true)

	// ボーンデフォーム情報を更新する
	updateGlobalMatrix(vmdDeltas.Bones, model.Bones.LayerSortedBoneIndexes[true])

	return vmdDeltas
}
