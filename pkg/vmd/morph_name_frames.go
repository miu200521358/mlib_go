package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type MorphNameFrames struct {
	*BaseFrames[*MorphFrame]
	Name string // ボーン名
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		BaseFrames: NewBaseFrames[*MorphFrame](NewMorphFrame, NullMorphFrame),
		Name:       name,
	}
}

func (i *MorphNameFrames) NewFrame(index int) *MorphFrame {
	return NewMorphFrame(index)
}

func (fs *MorphNameFrames) DeformVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) *VertexMorphDeltas {
	mf := fs.Get(frame)
	if mf == nil || mf.Ratio == 0.0 {
		return deltas
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			if offset.Position != nil && !offset.Position.IsZero() {
				if delta.Position == nil {
					delta.Position = offset.Position.MuledScalar(mf.Ratio)
				} else if !offset.Position.IsZero() {
					delta.Position.Add(offset.Position.MuledScalar(mf.Ratio))
				}
			}
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (fs *MorphNameFrames) DeformAfterVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) *VertexMorphDeltas {
	mf := fs.Get(frame)
	if mf == nil || mf.Ratio == 0.0 {
		return deltas
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			delta.AfterPosition.Add(offset.Position.MuledScalar(mf.Ratio))
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (fs *MorphNameFrames) DeformUv(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) *VertexMorphDeltas {
	mf := fs.Get(frame)
	if mf == nil || mf.Ratio == 0.0 {
		return deltas
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			uv := offset.Uv.MuledScalar(mf.Ratio).GetXY()
			delta.Uv.Add(uv)
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (fs *MorphNameFrames) DeformUv1(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) *VertexMorphDeltas {
	mf := fs.Get(frame)
	if mf == nil || mf.Ratio == 0.0 {
		return deltas
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			uv := offset.Uv.MuledScalar(mf.Ratio)
			delta.Uv1.Add(uv.GetXY())
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (fs *MorphNameFrames) DeformBone(
	frame int,
	model *pmx.PmxModel,
	deltas *BoneMorphDeltas,
) *BoneMorphDeltas {
	mf := fs.Get(frame)
	if mf == nil || mf.Ratio == 0.0 {
		return deltas
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.BoneMorphOffset)
		if 0 < offset.BoneIndex {
			delta := deltas.Get(offset.BoneIndex)
			if delta == nil {
				delta = NewBoneMorphDelta(offset.BoneIndex)
			}

			offsetPos := offset.Position.MuledScalar(mf.Ratio)
			offsetQuat := offset.Rotation.GetQuaternion().MuledScalar(mf.Ratio).Normalize()
			offsetScale := offset.Scale.MuledScalar(mf.Ratio)

			if delta.MorphFrameDelta.framePosition == nil {
				delta.MorphFrameDelta.framePosition = offsetPos
			} else {
				delta.MorphFrameDelta.framePosition.Add(offsetPos)
			}

			if delta.MorphFrameDelta.frameRotation == nil {
				delta.MorphFrameDelta.frameRotation = offsetQuat
			} else {
				delta.MorphFrameDelta.frameRotation = offsetQuat.Mul(delta.MorphFrameDelta.frameRotation)
			}

			if delta.MorphFrameDelta.frameScale == nil {
				delta.MorphFrameDelta.frameScale = offsetScale
			} else {
				delta.MorphFrameDelta.frameScale.Add(offsetScale)
			}

			deltas.Append(delta)
		}
	}

	return deltas
}

// DeformMaterial 材質モーフの適用
func (fs *MorphNameFrames) DeformMaterial(
	frame int,
	model *pmx.PmxModel,
	deltas *MaterialMorphDeltas,
) *MaterialMorphDeltas {
	mf := fs.Get(frame)
	if mf == nil || mf.Ratio == 0.0 {
		return deltas
	}

	morph := model.Morphs.GetByName(fs.Name)
	// 乗算→加算の順で処理
	for _, calcMode := range []pmx.MaterialMorphCalcMode{pmx.CALC_MODE_MULTIPLICATION, pmx.CALC_MODE_ADDITION} {
		for _, o := range morph.Offsets {
			offset := o.(*pmx.MaterialMorphOffset)
			if offset.CalcMode != calcMode {
				continue
			}
			if offset.MaterialIndex < 0 {
				// 全材質対象の場合
				for m, delta := range deltas.Data {
					if delta == nil {
						delta = NewMaterialMorphDelta(model.Materials.Get(m))
					}
					if calcMode == pmx.CALC_MODE_MULTIPLICATION {
						delta.Mul(offset, mf.Ratio)
					} else {
						delta.Add(offset, mf.Ratio)
					}
					deltas.Data[m] = delta
				}
			} else if 0 < offset.MaterialIndex && offset.MaterialIndex <= len(deltas.Data) {
				// 特定材質のみの場合
				delta := deltas.Data[offset.MaterialIndex]
				if delta == nil {
					delta = NewMaterialMorphDelta(model.Materials.Get(offset.MaterialIndex))
				}
				if calcMode == pmx.CALC_MODE_MULTIPLICATION {
					delta.Mul(offset, mf.Ratio)
				} else {
					delta.Add(offset, mf.Ratio)
				}
				deltas.Data[offset.MaterialIndex] = delta
			}
		}
	}

	return deltas
}
