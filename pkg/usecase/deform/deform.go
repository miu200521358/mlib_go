package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/miter"
)

func DeformModel(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	frame int,
) *pmx.PmxModel {
	vmdDeltas := delta.NewVmdDeltas(float32(frame), model.Bones, "", "")
	vmdDeltas.Morphs = DeformMorph(model, motion.MorphFrames, float32(frame), nil)
	vmdDeltas = deformBoneByPhysicsFlag(model, motion, vmdDeltas, false, float32(frame), nil, false)

	// 頂点にボーン変形を適用
	model.Vertices.ForEach(func(index int, vertex *pmx.Vertex) bool {
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

		return true
	})

	// ボーンの位置を更新
	model.Bones.ForEach(func(index int, bone *pmx.Bone) bool {
		if vmdDeltas.Bones.Get(index) != nil {
			bone.Position = vmdDeltas.Bones.Get(index).FilledGlobalPosition()
		}

		return true
	})

	return model
}

func DeformIks(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	frame float32,
	ikBones []*pmx.Bone,
	ikTargetBones []*pmx.Bone,
	ikGlobalPositions []*mmath.MVec3,
	boneNames []string,
	loopCount int, // IKのループ回数
	isRemoveTwist bool, // IKの捻りを除去するかどうか
	isForceDebug bool, // IKのデバッグを強制的に有効にするかどうか
) (*delta.VmdDeltas, []int) {
	if boneNames == nil {
		boneNames = make([]string, 0)
	}
	for _, ikBone := range ikBones {
		ikTargetBone, _ := model.Bones.Get(ikBone.Ik.BoneIndex)
		boneNames = append(boneNames, ikTargetBone.Name())
		for _, link := range ikBone.Ik.Links {
			linkBone, _ := model.Bones.Get(link.BoneIndex)
			boneNames = append(boneNames, linkBone.Name())
		}
	}

	deformBoneIndexes, deltas := newVmdDeltas(model, motion, deltas, frame, boneNames, false)

	for range loopCount {
		for i, ikBone := range ikBones {
			// IK変形リスト（IKのターゲットで代用して、ボーンの子孫にあたるボーンインデックス一覧）
			ikTargetDeformBoneIndexes := model.Bones.DeformBoneIndexes[ikTargetBones[i].Index()]

			// 変形リストを再帰的に更新 (IKの前に対象ボーンを先に最新化)
			// IK対象ボーンの子階層がまだ最新でない場合、先に更新する
			deltas.Bones = fillBoneDeform(
				model,
				motion,
				deltas,
				frame,
				ikTargetDeformBoneIndexes,
				false, // IK再帰呼び出ししない
				false,
			)

			// 親→子の順にグローバル行列を再更新
			UpdateGlobalMatrix(deltas.Bones, ikTargetDeformBoneIndexes)

			// IK適用前のグローバル行列を保存
			for _, idx := range ikTargetDeformBoneIndexes {
				linkD := deltas.Bones.Get(idx)
				if linkD != nil {
					linkD.GlobalIkOffMatrix = linkD.GlobalMatrix.Copy()
					deltas.Bones.Update(linkD)
				}
			}

			deformIk(model, motion, deltas, frame, false, ikBone, ikGlobalPositions[i],
				ikTargetDeformBoneIndexes, 0, isRemoveTwist, isForceDebug)
		}
	}

	UpdateGlobalMatrix(deltas.Bones, deformBoneIndexes)

	return deltas, deformBoneIndexes
}

// DeformBone 前回情報なしでボーンデフォーム処理を実行する
func DeformBone(
	model *pmx.PmxModel,
	morphMotion *vmd.VmdMotion,
	boneMotion *vmd.VmdMotion,
	isCalcIk bool,
	iFrame int,
	boneNames []string,
) *delta.VmdDeltas {
	frame := float32(iFrame)

	vmdDeltas := delta.NewVmdDeltas(frame, model.Bones, model.Hash(), morphMotion.Hash())
	vmdDeltas.Morphs = deformBoneMorph(model, morphMotion.MorphFrames, frame, nil)
	return deformBoneByPhysicsFlag(model, boneMotion, vmdDeltas, isCalcIk, frame, boneNames, false)
}

// DeformBone 前回情報ありでボーンデフォーム処理を実行する
func DeformBoneWithDeltas(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	isCalcIk bool,
	iFrame int,
	boneNames []string,
) *delta.VmdDeltas {
	frame := float32(iFrame)

	return deformBoneByPhysicsFlag(model, motion, deltas, isCalcIk, frame, boneNames, false)
}

// deformBoneByPhysicsFlag ボーンデフォーム処理を実行する
func deformBoneByPhysicsFlag(
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
	UpdateGlobalMatrix(deltas.Bones, deformBoneIndexes)

	return deltas
}

// DeformBeforePhysicsReset 物理前のボーンデフォーム処理を実行する
func DeformBeforePhysicsReset(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	frame float32,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return deltas
	}

	if deltas == nil {
		deltas = delta.NewVmdDeltas(frame, model.Bones, model.Hash(), motion.Hash())
	} else {
		deltas.SetFrame(frame)
		deltas.SetModelHash(model.Hash())
		deltas.SetMotionHash(motion.Hash())
		deltas.Bones = delta.NewBoneDeltas(model.Bones)
	}

	deltas.Morphs = DeformMorph(model, motion.MorphFrames, frame, nil)

	// ボーンデフォーム情報を埋める(物理前後全部埋める)
	deltas.Bones = fillBoneDeform(model, motion, deltas, frame, model.Bones.LayerSortedBoneIndexes[false], true, false)
	deltas.Bones = fillBoneDeform(model, motion, deltas, frame, model.Bones.LayerSortedBoneIndexes[true], true, true)

	// ボーンデフォーム情報を更新する
	UpdateGlobalMatrix(deltas.Bones, model.Bones.LayerSortedIndexes)

	return deltas
}

// DeformBeforePhysics 物理前のボーンデフォーム処理を実行する
func DeformBeforePhysics(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	frame float32,
) *delta.VmdDeltas {
	if model == nil || motion == nil {
		return deltas
	}

	if deltas != nil && deltas.Frame() == frame &&
		deltas.ModelHash() == model.Hash() && deltas.MotionHash() == motion.Hash() {
		return deltas
	}

	// 前とは条件が違う場合のみ再計算
	if deltas == nil {
		deltas = delta.NewVmdDeltas(frame, model.Bones, model.Hash(), motion.Hash())
	} else {
		deltas.SetFrame(frame)
		deltas.SetModelHash(model.Hash())
		deltas.SetMotionHash(motion.Hash())
		deltas.Bones = delta.NewBoneDeltas(model.Bones)
	}

	deltas.Morphs = DeformMorph(model, motion.MorphFrames, frame, nil)

	// ボーンデフォーム情報を埋める(物理前後全部埋める)
	deltas.Bones = fillBoneDeform(model, motion, deltas, frame, model.Bones.LayerSortedIndexes, true, false)

	// ボーンデフォーム情報を更新する
	UpdateGlobalMatrix(deltas.Bones, model.Bones.LayerSortedIndexes)

	return deltas
}

// DeformForPhysics 物理剛体位置を更新する
func DeformForPhysics(
	shared *state.SharedState,
	physics physics.IPhysics,
	model *pmx.PmxModel,
	deltas *delta.VmdDeltas,
) *delta.VmdDeltas {
	if model == nil {
		return deltas
	}

	// 物理剛体位置を更新
	if err := miter.IterParallelByList(mmath.IntRanges(model.RigidBodies.Length()-1), 100, 0, func(i int, rigidBodyIndex int) error {
		rigidBody, err := model.RigidBodies.Get(rigidBodyIndex)
		if err != nil {
			return err
		}

		// 現在のボーン変形情報を保持
		if rigidBody.Bone == nil || deltas.Bones.Get(rigidBody.Bone.Index()) == nil {
			return nil
		}

		if (shared.IsEnabledPhysics() && rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC) ||
			shared.IsPhysicsReset() {
			// 通常はボーン追従剛体・物理＋ボーン剛体だけ。物理リセット時は全部更新
			physics.UpdateTransform(model.Index(), rigidBody.Bone,
				deltas.Bones.Get(rigidBody.Bone.Index()).FilledGlobalMatrix(), rigidBody)
		}

		return nil
	}, nil); err != nil {
		return deltas
	}

	return deltas
}

// DeformAfterPhysics 物理後のボーンデフォーム処理を実行する
func DeformAfterPhysics(
	shared *state.SharedState,
	physics physics.IPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	vmdDeltas *delta.VmdDeltas,
	frame float32,
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
				bd := delta.NewBoneDeltaByGlobalMatrix(bone, frame,
					bonePhysicsGlobalMatrix, vmdDeltas.Bones.Get(bone.ParentIndex))
				vmdDeltas.Bones.Update(bd)
			}
		}
	}

	// ボーンデフォーム情報を埋める(物理後のみ埋める)
	vmdDeltas.Bones = fillBoneDeform(model, motion, vmdDeltas, frame,
		model.Bones.LayerSortedBoneIndexes[true], true, true)

	// ボーンデフォーム情報を更新する
	UpdateGlobalMatrix(vmdDeltas.Bones, model.Bones.LayerSortedBoneIndexes[true])

	return vmdDeltas
}
