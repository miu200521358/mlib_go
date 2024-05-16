package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type MorphNameFrames struct {
	*BaseFrames[*MorphFrame]
	Name string // ボーン名
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		BaseFrames: NewBaseFrames[*MorphFrame](NewMorphFrame),
		Name:       name,
	}
}

func (i *MorphNameFrames) NewFrame(index int) *MorphFrame {
	return NewMorphFrame(index)
}

func (fs *MorphNameFrames) AnimateVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			delta.Position.Add(offset.Position.MuledScalar(mf.Ratio))
		}
	}
}

func (fs *MorphNameFrames) AnimateAfterVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			delta.AfterPosition.Add(offset.Position.MuledScalar(mf.Ratio))
		}
	}
}

func (fs *MorphNameFrames) AnimateUv(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			uv := offset.Uv.MuledScalar(mf.Ratio).GetXY()
			delta.Uv.Add(uv)
		}
	}
}

func (fs *MorphNameFrames) AnimateUv1(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta()
			}
			uv := offset.Uv.MuledScalar(mf.Ratio)
			delta.Uv1.Add(uv.GetXY())
		}
	}
}

func (fs *MorphNameFrames) AnimateBone(
	frame int,
	model *pmx.PmxModel,
	deltas *BoneMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.BoneMorphOffset)
		if 0 < offset.BoneIndex && offset.BoneIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.BoneIndex]
			if delta == nil {
				delta = NewBoneMorphDelta()
			}
			delta.MorphPosition.Add(offset.Position.MuledScalar(mf.Ratio))
			delta.MorphLocalPosition.Add(offset.LocalPosition.MuledScalar(mf.Ratio))
			deltaRad := offset.Rotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphRotation.SetQuaternion(delta.MorphRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromRadians(deltaRad.GetX(), deltaRad.GetY(), deltaRad.GetZ())))
			deltaLocalRad := offset.LocalRotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphLocalRotation.SetQuaternion(delta.MorphLocalRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromRadians(deltaLocalRad.GetX(), deltaLocalRad.GetY(), deltaLocalRad.GetZ())))
			delta.MorphScale.Add(offset.Scale.MuledScalar(mf.Ratio))
			delta.MorphLocalScale.Add(offset.LocalScale.MuledScalar(mf.Ratio))
		}
	}
}

// AnimateMaterial 材質モーフの適用
func (fs *MorphNameFrames) AnimateMaterial(
	frame int,
	model *pmx.PmxModel,
	deltas *MaterialMorphDeltas,
) {
	mf := fs.Get(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	// 乗算→加算の順で処理
	for _, calcMode := range []pmx.MaterialMorphCalcMode{pmx.CALC_MODE_MULTIPLICATION, pmx.CALC_MODE_ADDITION} {
		for _, o := range morph.Offsets {
			offset := o.(*pmx.MaterialMorphOffset)
			if offset.CalcMode != calcMode {
				continue
			}
			if offset.MaterialIndex < 0 {
				// 全材質対象の場合
				for _, delta := range deltas.Data {
					if calcMode == pmx.CALC_MODE_MULTIPLICATION {
						delta.Mul(offset, mf.Ratio)
					} else {
						delta.Add(offset, mf.Ratio)
					}
				}
			} else if 0 < offset.MaterialIndex && offset.MaterialIndex <= len(deltas.Data) {
				// 特定材質のみの場合
				if calcMode == pmx.CALC_MODE_MULTIPLICATION {
					deltas.Data[offset.MaterialIndex].Mul(offset, mf.Ratio)
				} else {
					deltas.Data[offset.MaterialIndex].Add(offset, mf.Ratio)
				}
			}
		}
	}
}
