// 指示: miu200521358
package deform

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// ComputeMorphDeltas はモーフ差分を算出して返す。
func ComputeMorphDeltas(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	frame motion.Frame,
	morphNames []string,
) *delta.MorphDeltas {
	if modelData == nil {
		return delta.NewMorphDeltas(nil, nil, nil)
	}
	deltas := delta.NewMorphDeltas(modelData.Vertices, modelData.Materials, modelData.Bones)
	if motionData == nil || motionData.MorphFrames == nil {
		return deltas
	}

	names := morphNames
	if names == nil {
		values := modelData.Morphs.Values()
		names = make([]string, 0, len(values))
		for _, morph := range values {
			if morph == nil {
				continue
			}
			names = append(names, morph.Name())
		}
	}

	visited := make(map[int]struct{})
	for _, name := range names {
		if !motionData.MorphFrames.Has(name) {
			continue
		}
		mf := motionData.MorphFrames.Get(name).Get(frame)
		if mf == nil || math.Abs(mf.Ratio) < 1e-12 {
			continue
		}
		morph, err := modelData.Morphs.GetByName(name)
		if err != nil {
			continue
		}
		applyMorphDelta(deltas, modelData, morph, mf.Ratio, visited)
	}
	return deltas
}

// computeBoneMorphDeltas はボーンモーフ差分を算出して返す。
func computeBoneMorphDeltas(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	frame motion.Frame,
	morphNames []string,
) *delta.BoneMorphDeltas {
	boneCount := 0
	if modelData != nil && modelData.Bones != nil {
		boneCount = modelData.Bones.Len()
	}
	deltas := delta.NewBoneMorphDeltas(boneCount)
	if modelData == nil || motionData == nil || motionData.MorphFrames == nil {
		return deltas
	}

	names := morphNames
	if names == nil {
		values := modelData.Morphs.Values()
		names = make([]string, 0, len(values))
		for _, morph := range values {
			if morph == nil {
				continue
			}
			names = append(names, morph.Name())
		}
	}

	visited := make(map[int]struct{})
	for _, name := range names {
		if !motionData.MorphFrames.Has(name) {
			continue
		}
		mf := motionData.MorphFrames.Get(name).Get(frame)
		if mf == nil || math.Abs(mf.Ratio) < 1e-12 {
			continue
		}
		morph, err := modelData.Morphs.GetByName(name)
		if err != nil || morph == nil {
			continue
		}
		applyBoneMorphDelta(deltas, modelData, morph, mf.Ratio, visited)
	}
	return deltas
}

// ApplyMorphDeltas はモーフ差分をモデルへ適用する。
func ApplyMorphDeltas(modelData *model.PmxModel, deltas *delta.MorphDeltas) {
	if modelData == nil || deltas == nil {
		return
	}
	applyVertexMorphDeltas(modelData, deltas.Vertices())
	applyMaterialMorphDeltas(modelData, deltas.Materials())
}

// applyMorphDelta はモーフ種別ごとに差分を適用する。
func applyMorphDelta(
	deltas *delta.MorphDeltas,
	modelData *model.PmxModel,
	morph *model.Morph,
	ratio float64,
	visited map[int]struct{},
) {
	if deltas == nil || modelData == nil || morph == nil {
		return
	}
	if _, ok := visited[morph.Index()]; ok {
		return
	}
	visited[morph.Index()] = struct{}{}
	defer delete(visited, morph.Index())

	switch morph.MorphType {
	case model.MORPH_TYPE_VERTEX:
		applyVertexMorph(deltas.Vertices(), morph, ratio)
	case model.MORPH_TYPE_AFTER_VERTEX:
		applyAfterVertexMorph(deltas.Vertices(), morph, ratio)
	case model.MORPH_TYPE_UV:
		applyUvMorph(deltas.Vertices(), morph, ratio)
	case model.MORPH_TYPE_EXTENDED_UV1:
		applyUv1Morph(deltas.Vertices(), morph, ratio)
	case model.MORPH_TYPE_BONE:
		applyBoneMorph(deltas.Bones(), morph, ratio)
	case model.MORPH_TYPE_MATERIAL:
		applyMaterialMorph(deltas.Materials(), modelData, morph, ratio)
	case model.MORPH_TYPE_GROUP:
		applyGroupMorph(deltas, modelData, morph, ratio, visited)
	}
}

// applyBoneMorphDelta はボーン/グループモーフの差分を集計する。
func applyBoneMorphDelta(
	deltas *delta.BoneMorphDeltas,
	modelData *model.PmxModel,
	morph *model.Morph,
	ratio float64,
	visited map[int]struct{},
) {
	if deltas == nil || modelData == nil || morph == nil {
		return
	}
	if math.Abs(ratio) < 1e-12 {
		return
	}
	if _, ok := visited[morph.Index()]; ok {
		return
	}
	visited[morph.Index()] = struct{}{}
	defer delete(visited, morph.Index())

	switch morph.MorphType {
	case model.MORPH_TYPE_BONE:
		applyBoneMorph(deltas, morph, ratio)
	case model.MORPH_TYPE_GROUP:
		for _, raw := range morph.Offsets {
			offset, ok := raw.(*model.GroupMorphOffset)
			if !ok || offset.MorphIndex < 0 {
				continue
			}
			groupMorph, err := modelData.Morphs.Get(offset.MorphIndex)
			if err != nil || groupMorph == nil {
				continue
			}
			applyBoneMorphDelta(deltas, modelData, groupMorph, ratio*offset.MorphFactor, visited)
		}
	}
}

// applyVertexMorph は頂点モーフ差分を加算する。
func applyVertexMorph(deltas *delta.VertexMorphDeltas, morph *model.Morph, ratio float64) {
	if deltas == nil || morph == nil {
		return
	}
	for _, raw := range morph.Offsets {
		offset, ok := raw.(*model.VertexMorphOffset)
		if !ok || offset.VertexIndex < 0 {
			continue
		}
		d := deltas.Get(offset.VertexIndex)
		if d == nil {
			d = delta.NewVertexMorphDelta(offset.VertexIndex)
		}
		offsetPos := offset.Position.MuledScalar(ratio)
		if d.Position == nil {
			d.Position = &offsetPos
		} else {
			pos := d.Position.Added(offsetPos)
			d.Position = &pos
		}
		deltas.Update(d)
	}
}

// applyAfterVertexMorph はボーン変形後頂点モーフ差分を加算する。
func applyAfterVertexMorph(deltas *delta.VertexMorphDeltas, morph *model.Morph, ratio float64) {
	if deltas == nil || morph == nil {
		return
	}
	for _, raw := range morph.Offsets {
		offset, ok := raw.(*model.VertexMorphOffset)
		if !ok || offset.VertexIndex < 0 {
			continue
		}
		d := deltas.Get(offset.VertexIndex)
		if d == nil {
			d = delta.NewVertexMorphDelta(offset.VertexIndex)
		}
		offsetPos := offset.Position.MuledScalar(ratio)
		if d.AfterPosition == nil {
			d.AfterPosition = &offsetPos
		} else {
			pos := d.AfterPosition.Added(offsetPos)
			d.AfterPosition = &pos
		}
		deltas.Update(d)
	}
}

// applyUvMorph はUVモーフ差分を加算する。
func applyUvMorph(deltas *delta.VertexMorphDeltas, morph *model.Morph, ratio float64) {
	if deltas == nil || morph == nil {
		return
	}
	for _, raw := range morph.Offsets {
		offset, ok := raw.(*model.UvMorphOffset)
		if !ok || offset.VertexIndex < 0 {
			continue
		}
		d := deltas.Get(offset.VertexIndex)
		if d == nil {
			d = delta.NewVertexMorphDelta(offset.VertexIndex)
		}
		uv := offset.Uv.MuledScalar(ratio).XY()
		if d.Uv == nil {
			d.Uv = &uv
		} else {
			out := d.Uv.Added(uv)
			d.Uv = &out
		}
		deltas.Update(d)
	}
}

// applyUv1Morph は追加UV1モーフ差分を加算する。
func applyUv1Morph(deltas *delta.VertexMorphDeltas, morph *model.Morph, ratio float64) {
	if deltas == nil || morph == nil {
		return
	}
	for _, raw := range morph.Offsets {
		offset, ok := raw.(*model.UvMorphOffset)
		if !ok || offset.VertexIndex < 0 {
			continue
		}
		d := deltas.Get(offset.VertexIndex)
		if d == nil {
			d = delta.NewVertexMorphDelta(offset.VertexIndex)
		}
		uv := offset.Uv.MuledScalar(ratio).XY()
		if d.Uv1 == nil {
			d.Uv1 = &uv
		} else {
			out := d.Uv1.Added(uv)
			d.Uv1 = &out
		}
		deltas.Update(d)
	}
}

// applyBoneMorph はボーンモーフ差分を加算する。
func applyBoneMorph(deltas *delta.BoneMorphDeltas, morph *model.Morph, ratio float64) {
	if deltas == nil || morph == nil {
		return
	}
	for _, raw := range morph.Offsets {
		offset, ok := raw.(*model.BoneMorphOffset)
		if !ok || offset.BoneIndex < 0 {
			continue
		}
		d := deltas.Get(offset.BoneIndex)
		if d == nil {
			d = delta.NewBoneMorphDelta(offset.BoneIndex)
		}
		if !offset.Position.IsZero() {
			offsetPos := offset.Position.MuledScalar(ratio)
			if d.FramePosition == nil {
				d.FramePosition = &offsetPos
			} else {
				pos := d.FramePosition.Added(offsetPos)
				d.FramePosition = &pos
			}
		}
		if !offset.Rotation.IsIdent() {
			offsetQuat := offset.Rotation.MuledScalar(ratio).Normalized()
			if d.FrameRotation == nil {
				d.FrameRotation = &offsetQuat
			} else {
				rot := offsetQuat.Muled(*d.FrameRotation)
				d.FrameRotation = &rot
			}
		}
		deltas.Update(d)
	}
}

// applyMaterialMorph は材質モーフ差分を加算する。
func applyMaterialMorph(
	deltas *delta.MaterialMorphDeltas,
	modelData *model.PmxModel,
	morph *model.Morph,
	ratio float64,
) {
	if deltas == nil || modelData == nil || morph == nil {
		return
	}
	for _, calcMode := range []model.MaterialMorphCalcMode{model.CALC_MODE_MULTIPLICATION, model.CALC_MODE_ADDITION} {
		for _, raw := range morph.Offsets {
			offset, ok := raw.(*model.MaterialMorphOffset)
			if !ok || offset.CalcMode != calcMode {
				continue
			}
			if offset.MaterialIndex < 0 {
				deltas.ForEach(func(index int, data *delta.MaterialMorphDelta) bool {
					if data == nil {
						mat, err := modelData.Materials.Get(index)
						if err != nil || mat == nil {
							return true
						}
						data = delta.NewMaterialMorphDelta(mat)
					}
					applyMaterialCalcMode(data, offset, ratio, calcMode)
					deltas.Update(data)
					return true
				})
				continue
			}
			if offset.MaterialIndex >= deltas.Len() {
				continue
			}
			data := deltas.Get(offset.MaterialIndex)
			if data == nil {
				mat, err := modelData.Materials.Get(offset.MaterialIndex)
				if err != nil || mat == nil {
					continue
				}
				data = delta.NewMaterialMorphDelta(mat)
			}
			applyMaterialCalcMode(data, offset, ratio, calcMode)
			deltas.Update(data)
		}
	}
}

// applyMaterialCalcMode は材質モーフ計算モードを適用する。
func applyMaterialCalcMode(
	data *delta.MaterialMorphDelta,
	offset *model.MaterialMorphOffset,
	ratio float64,
	calcMode model.MaterialMorphCalcMode,
) {
	if calcMode == model.CALC_MODE_MULTIPLICATION {
		data.Mul(offset, ratio)
		return
	}
	data.Add(offset, ratio)
}

// applyGroupMorph はグループモーフを展開する。
func applyGroupMorph(
	deltas *delta.MorphDeltas,
	modelData *model.PmxModel,
	morph *model.Morph,
	ratio float64,
	visited map[int]struct{},
) {
	if deltas == nil || modelData == nil || morph == nil {
		return
	}
	for _, raw := range morph.Offsets {
		offset, ok := raw.(*model.GroupMorphOffset)
		if !ok || offset.MorphIndex < 0 {
			continue
		}
		groupMorph, err := modelData.Morphs.Get(offset.MorphIndex)
		if err != nil {
			continue
		}
		applyMorphDelta(deltas, modelData, groupMorph, ratio*offset.MorphFactor, visited)
	}
}

// applyVertexMorphDeltas は頂点モーフ差分をモデルへ適用する。
func applyVertexMorphDeltas(modelData *model.PmxModel, deltas *delta.VertexMorphDeltas) {
	if modelData == nil || deltas == nil || modelData.Vertices == nil {
		return
	}
	deltas.ForEach(func(index int, d *delta.VertexMorphDelta) bool {
		if d == nil || d.IsZero() {
			return true
		}
		vertex, err := modelData.Vertices.Get(index)
		if err != nil || vertex == nil {
			return true
		}
		if d.Position != nil {
			vertex.Position = vertex.Position.Added(*d.Position)
		}
		// AfterPosition はスキニング後に適用するためここでは反映しない。
		if d.Uv != nil {
			vertex.Uv = vertex.Uv.Added(*d.Uv)
		}
		if d.Uv1 != nil && len(vertex.ExtendedUvs) > 0 {
			uv1 := vertex.ExtendedUvs[0]
			uv1.X += d.Uv1.X
			uv1.Y += d.Uv1.Y
			vertex.ExtendedUvs[0] = uv1
		}
		return true
	})
}

// applyMaterialMorphDeltas は材質モーフ差分をモデルへ適用する。
func applyMaterialMorphDeltas(modelData *model.PmxModel, deltas *delta.MaterialMorphDeltas) {
	if modelData == nil || deltas == nil || modelData.Materials == nil {
		return
	}
	deltas.ForEach(func(index int, d *delta.MaterialMorphDelta) bool {
		if d == nil {
			return true
		}
		material, err := modelData.Materials.Get(index)
		if err != nil || material == nil {
			return true
		}
		material.Diffuse = material.Diffuse.Muled(d.MulMaterial.Diffuse).Added(d.AddMaterial.Diffuse)
		material.Specular = material.Specular.Muled(d.MulMaterial.Specular).Added(d.AddMaterial.Specular)
		material.Ambient = material.Ambient.Muled(d.MulMaterial.Ambient).Added(d.AddMaterial.Ambient)
		material.Edge = material.Edge.Muled(d.MulMaterial.Edge).Added(d.AddMaterial.Edge)
		material.EdgeSize = material.EdgeSize*d.MulMaterial.EdgeSize + d.AddMaterial.EdgeSize
		material.TextureFactor = material.TextureFactor.Muled(d.MulMaterial.TextureFactor).Added(d.AddMaterial.TextureFactor)
		material.SphereTextureFactor = material.SphereTextureFactor.Muled(d.MulMaterial.SphereTextureFactor).Added(d.AddMaterial.SphereTextureFactor)
		material.ToonTextureFactor = material.ToonTextureFactor.Muled(d.MulMaterial.ToonTextureFactor).Added(d.AddMaterial.ToonTextureFactor)
		return true
	})
}
