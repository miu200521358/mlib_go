package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

func DeformMorph(
	model *pmx.PmxModel,
	mfs *vmd.MorphFrames,
	frame int,
	morphNames []string,
) *delta.MorphDeltas {
	if morphNames == nil {
		// モーフの指定がなければ全モーフチェック
		morphNames = make([]string, 0)
		for _, morph := range model.Morphs.Data {
			morphNames = append(morphNames, morph.Name())
		}
	}

	mds := delta.NewMorphDeltas(model.Materials, model.Bones)
	for _, morphName := range morphNames {
		if !mfs.Contains(morphName) || !model.Morphs.ContainsByName(morphName) {
			continue
		}

		mf := mfs.Get(morphName).Get(float64(frame))
		if mf == nil {
			continue
		}

		morph := model.Morphs.GetByName(morphName)
		switch morph.MorphType {
		case pmx.MORPH_TYPE_VERTEX:
			mds.Vertices = deformVertex(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_AFTER_VERTEX:
			mds.Vertices = deformAfterVertex(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_UV:
			mds.Vertices = deformUv(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_EXTENDED_UV1:
			mds.Vertices = deformUv1(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_BONE:
			mds.Bones = deformBone(morphName, model, mds.Bones, mf.Ratio)
		case pmx.MORPH_TYPE_MATERIAL:
			mds.Materials = deformMaterial(morphName, model, mds.Materials, mf.Ratio)
		case pmx.MORPH_TYPE_GROUP:
			// グループモーフは細分化
			for _, offset := range morph.Offsets {
				groupOffset := offset.(*pmx.GroupMorphOffset)
				groupMorph := model.Morphs.Get(groupOffset.MorphIndex)
				if groupMorph == nil {
					continue
				}
				switch groupMorph.MorphType {
				case pmx.MORPH_TYPE_VERTEX:
					mds.Vertices = deformVertex(
						groupMorph.Name(), model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_AFTER_VERTEX:
					mds.Vertices = deformAfterVertex(
						groupMorph.Name(), model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_UV:
					mds.Vertices = deformUv(
						groupMorph.Name(), model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_EXTENDED_UV1:
					mds.Vertices = deformUv1(
						groupMorph.Name(), model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_BONE:
					mds.Bones = deformBone(
						groupMorph.Name(), model, mds.Bones, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_MATERIAL:
					mds.Materials = deformMaterial(
						groupMorph.Name(), model, mds.Materials, mf.Ratio*groupOffset.MorphFactor)
				}
			}
		}
	}

	return mds
}

func deformVertex(
	morphName string,
	model *pmx.PmxModel,
	deltas *delta.VertexMorphDeltas,
	ratio float64,
) *delta.VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex {
			d := deltas.Data[offset.VertexIndex]
			if d == nil {
				d = delta.NewVertexMorphDelta(offset.VertexIndex)
			}
			if offset.Position != nil {
				if d.Position == nil {
					d.Position = offset.Position.MuledScalar(ratio)
				} else if !offset.Position.IsZero() {
					d.Position.Add(offset.Position.MuledScalar(ratio))
				}
			}
			deltas.Data[offset.VertexIndex] = d
		}
	}

	return deltas
}

func deformAfterVertex(
	morphName string,
	model *pmx.PmxModel,
	deltas *delta.VertexMorphDeltas,
	ratio float64,
) *delta.VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex {
			d := deltas.Data[offset.VertexIndex]
			if d == nil {
				d = delta.NewVertexMorphDelta(offset.VertexIndex)
			}
			if d.AfterPosition == nil {
				d.AfterPosition = mmath.NewMVec3()
			}
			d.AfterPosition.Add(offset.Position.MuledScalar(ratio))
			deltas.Data[offset.VertexIndex] = d
		}
	}

	return deltas
}

func deformUv(
	morphName string,
	model *pmx.PmxModel,
	deltas *delta.VertexMorphDeltas,
	ratio float64,
) *delta.VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex {
			d := deltas.Data[offset.VertexIndex]
			if d == nil {
				d = delta.NewVertexMorphDelta(offset.VertexIndex)
			}
			if d.Uv == nil {
				d.Uv = mmath.NewMVec2()
			}
			uv := offset.Uv.MuledScalar(ratio).XY()
			d.Uv.Add(uv)
			deltas.Data[offset.VertexIndex] = d
		}
	}

	return deltas
}

func deformUv1(
	morphName string,
	model *pmx.PmxModel,
	deltas *delta.VertexMorphDeltas,
	ratio float64,
) *delta.VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex {
			d := deltas.Data[offset.VertexIndex]
			if d == nil {
				d = delta.NewVertexMorphDelta(offset.VertexIndex)
			}
			if d.Uv1 == nil {
				d.Uv1 = mmath.NewMVec2()
			}
			uv := offset.Uv.MuledScalar(ratio)
			d.Uv1.Add(uv.XY())
			deltas.Data[offset.VertexIndex] = d
		}
	}

	return deltas
}

func deformBone(
	morphName string,
	model *pmx.PmxModel,
	deltas *delta.BoneMorphDeltas,
	ratio float64,
) *delta.BoneMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.BoneMorphOffset)
		if 0 < offset.BoneIndex {
			d := deltas.Get(offset.BoneIndex)
			if d == nil {
				d = delta.NewBoneMorphDelta(offset.BoneIndex)
			}

			offsetPos := offset.Position.MuledScalar(ratio)
			offsetQuat := offset.Rotation.Quaternion().MuledScalar(ratio).Normalize()
			offsetScale := offset.Extend.Scale.MuledScalar(ratio)

			if d.FramePosition == nil {
				d.FramePosition = offsetPos.Copy()
			} else {
				d.FramePosition.Add(offsetPos)
			}

			if d.FrameRotation == nil {
				d.FrameRotation = offsetQuat.Copy()
			} else {
				d.FrameRotation = offsetQuat.Muled(d.FrameRotation)
			}

			if d.FrameScale == nil {
				d.FrameScale = offsetScale.Copy()
			} else {
				d.FrameScale.Add(offsetScale)
			}

			deltas.Update(d)
		}
	}

	return deltas
}

// DeformMaterial 材質モーフの適用
func deformMaterial(
	morphName string,
	model *pmx.PmxModel,
	deltas *delta.MaterialMorphDeltas,
	ratio float64,
) *delta.MaterialMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	// 乗算→加算の順で処理
	for _, calcMode := range []pmx.MaterialMorphCalcMode{pmx.CALC_MODE_MULTIPLICATION, pmx.CALC_MODE_ADDITION} {
		for _, o := range morph.Offsets {
			offset := o.(*pmx.MaterialMorphOffset)
			if offset.CalcMode != calcMode {
				continue
			}
			if offset.MaterialIndex < 0 {
				// 全材質対象の場合
				for m, d := range deltas.Data {
					if d == nil {
						d = delta.NewMaterialMorphDelta(model.Materials.Get(m))
					}
					if calcMode == pmx.CALC_MODE_MULTIPLICATION {
						d.Mul(offset, ratio)
					} else {
						d.Add(offset, ratio)
					}
					deltas.Data[m] = d
				}
			} else if 0 <= offset.MaterialIndex && offset.MaterialIndex <= len(deltas.Data) {
				// 特定材質のみの場合
				d := deltas.Data[offset.MaterialIndex]
				if d == nil {
					d = delta.NewMaterialMorphDelta(model.Materials.Get(offset.MaterialIndex))
				}
				if calcMode == pmx.CALC_MODE_MULTIPLICATION {
					d.Mul(offset, ratio)
				} else {
					d.Add(offset, ratio)
				}
				deltas.Data[offset.MaterialIndex] = d
			}
		}
	}

	return deltas
}
