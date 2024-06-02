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
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			delta.Position.Add(offset.Position.MuledScalar(mf.Ratio))
			deltas.Data[offset.VertexIndex] = delta
		}
	}
}

func (fs *MorphNameFrames) DeformAfterVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			delta.AfterPosition.Add(offset.Position.MuledScalar(mf.Ratio))
			deltas.Data[offset.VertexIndex] = delta
		}
	}
}

func (fs *MorphNameFrames) DeformUv(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			uv := offset.Uv.MuledScalar(mf.Ratio).GetXY()
			delta.Uv.Add(uv)
			deltas.Data[offset.VertexIndex] = delta
		}
	}
}

func (fs *MorphNameFrames) DeformUv1(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			uv := offset.Uv.MuledScalar(mf.Ratio)
			delta.Uv1.Add(uv.GetXY())
			deltas.Data[offset.VertexIndex] = delta
		}
	}
}

func (fs *MorphNameFrames) DeformBone(
	frame int,
	model *pmx.PmxModel,
	deltas *BoneMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.BoneMorphOffset)
		if 0 < offset.BoneIndex && offset.BoneIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.BoneIndex]
			if delta == nil {
				delta = NewBoneMorphDelta()
			}

			if delta.MorphPosition == nil {
				delta.MorphPosition = offset.Position.MuledScalar(mf.Ratio).Copy()
			} else {
				delta.MorphPosition.Add(offset.Position.MuledScalar(mf.Ratio))
			}

			if delta.MorphLocalPosition == nil {
				delta.MorphLocalPosition = offset.LocalPosition.MuledScalar(mf.Ratio).Copy()
			} else {
				delta.MorphLocalPosition.Add(offset.LocalPosition.MuledScalar(mf.Ratio))
			}

			if delta.MorphRotation == nil {
				delta.MorphRotation = offset.Rotation.GetQuaternion().MuledScalar(mf.Ratio)
			} else {
				delta.MorphRotation.Mul(offset.Rotation.GetQuaternion().MuledScalar(mf.Ratio))
			}

			if delta.MorphLocalRotation == nil {
				delta.MorphLocalRotation = offset.LocalRotation.GetQuaternion().MuledScalar(mf.Ratio)
			} else {
				delta.MorphLocalRotation.Mul(offset.LocalRotation.GetQuaternion().MuledScalar(mf.Ratio))
			}

			if delta.MorphScale == nil {
				delta.MorphScale = offset.Scale.MuledScalar(mf.Ratio).Copy()
			} else {
				delta.MorphScale.Add(offset.Scale.MuledScalar(mf.Ratio))
			}

			if delta.MorphLocalScale == nil {
				delta.MorphLocalScale = offset.LocalScale.MuledScalar(mf.Ratio).Copy()
			} else {
				delta.MorphLocalScale.Add(offset.LocalScale.MuledScalar(mf.Ratio))
			}

			deltas.Data[offset.BoneIndex] = delta
		}
	}
}

// DeformMaterial 材質モーフの適用
func (fs *MorphNameFrames) DeformMaterial(
	frame int,
	model *pmx.PmxModel,
	deltas *MaterialMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
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
}
