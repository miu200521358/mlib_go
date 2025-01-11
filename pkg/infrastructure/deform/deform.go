package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
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
