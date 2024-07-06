package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type MorphFrame struct {
	*BaseFrame         // キーフレ
	Ratio      float64 // モーフの割合
}

func NewMorphFrame(index int) *MorphFrame {
	return &MorphFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Ratio:     0.0,
	}
}

func NullMorphFrame() *MorphFrame {
	return nil
}

func (mf *MorphFrame) Add(v *MorphFrame) {
	mf.Ratio += v.Ratio
}

func (mf *MorphFrame) Added(v *MorphFrame) *MorphFrame {
	copied := mf.Copy().(*MorphFrame)
	copied.Ratio += v.Ratio
	return copied
}

func (mf *MorphFrame) Copy() IBaseFrame {
	copied := NewMorphFrame(mf.GetIndex())
	copied.Ratio = mf.Ratio
	return copied
}

func (nextMf *MorphFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevMf := prevFrame.(*MorphFrame)

	prevIndex := prevMf.GetIndex()
	nextIndex := nextMf.GetIndex()

	mf := NewMorphFrame(index)

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*ry

	return mf
}

func (mf *MorphFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}

type MorphFrameDelta struct {
	framePosition *mmath.MVec3       // キーフレ位置の変動量
	frameRotation *mmath.MQuaternion // キーフレ回転の変動量
	frameScale    *mmath.MVec3       // キーフレスケールの変動量
}

func (md *MorphFrameDelta) FramePosition() *mmath.MVec3 {
	if md.framePosition == nil {
		md.framePosition = mmath.NewMVec3()
	}
	return md.framePosition
}

func (md *MorphFrameDelta) FrameRotation() *mmath.MQuaternion {
	if md.frameRotation == nil {
		md.frameRotation = mmath.NewMQuaternion()
	}
	return md.frameRotation
}

func (md *MorphFrameDelta) FrameScale() *mmath.MVec3 {
	if md.frameScale == nil {
		md.frameScale = mmath.NewMVec3()
	}
	return md.frameScale
}

func NewMorphFrameDelta() *MorphFrameDelta {
	return &MorphFrameDelta{}
}

func (md *MorphFrameDelta) Copy() *MorphFrameDelta {
	return &MorphFrameDelta{
		framePosition: md.FramePosition().Copy(),
		frameRotation: md.FrameRotation().Copy(),
		frameScale:    md.FrameScale().Copy(),
	}
}

func (mf *MorphFrame) DeformVertex(
	morphName string,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
	ratio float64,
) *VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			if offset.Position != nil {
				if delta.Position == nil {
					delta.Position = offset.Position.MuledScalar(ratio)
				} else if !offset.Position.IsZero() {
					delta.Position.Add(offset.Position.MuledScalar(ratio))
				}
			}
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (mf *MorphFrame) DeformAfterVertex(
	morphName string,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
	ratio float64,
) *VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			if delta.AfterPosition == nil {
				delta.AfterPosition = mmath.NewMVec3()
			}
			delta.AfterPosition.Add(offset.Position.MuledScalar(ratio))
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (mf *MorphFrame) DeformUv(
	morphName string,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
	ratio float64,
) *VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			if delta.Uv == nil {
				delta.Uv = mmath.NewMVec2()
			}
			uv := offset.Uv.MuledScalar(ratio).GetXY()
			delta.Uv.Add(uv)
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (mf *MorphFrame) DeformUv1(
	morphName string,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
	ratio float64,
) *VertexMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex {
			delta := deltas.Data[offset.VertexIndex]
			if delta == nil {
				delta = NewVertexMorphDelta(offset.VertexIndex)
			}
			if delta.Uv1 == nil {
				delta.Uv1 = mmath.NewMVec2()
			}
			uv := offset.Uv.MuledScalar(ratio)
			delta.Uv1.Add(uv.GetXY())
			deltas.Data[offset.VertexIndex] = delta
		}
	}

	return deltas
}

func (mf *MorphFrame) DeformBone(
	morphName string,
	model *pmx.PmxModel,
	deltas *BoneMorphDeltas,
	ratio float64,
) *BoneMorphDeltas {
	morph := model.Morphs.GetByName(morphName)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.BoneMorphOffset)
		if 0 < offset.BoneIndex {
			delta := deltas.Get(offset.BoneIndex)
			if delta == nil {
				delta = NewBoneMorphDelta(offset.BoneIndex)
			}

			offsetPos := offset.Position.MuledScalar(ratio)
			offsetQuat := offset.Rotation.GetQuaternion().MuledScalar(ratio).Normalize()
			offsetScale := offset.Scale.MuledScalar(ratio)

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
func (mf *MorphFrame) DeformMaterial(
	morphName string,
	model *pmx.PmxModel,
	deltas *MaterialMorphDeltas,
	ratio float64,
) *MaterialMorphDeltas {
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
				for m, delta := range deltas.Data {
					if delta == nil {
						delta = NewMaterialMorphDelta(model.Materials.Get(m))
					}
					if calcMode == pmx.CALC_MODE_MULTIPLICATION {
						delta.Mul(offset, ratio)
					} else {
						delta.Add(offset, ratio)
					}
					deltas.Data[m] = delta
				}
			} else if 0 <= offset.MaterialIndex && offset.MaterialIndex <= len(deltas.Data) {
				// 特定材質のみの場合
				delta := deltas.Data[offset.MaterialIndex]
				if delta == nil {
					delta = NewMaterialMorphDelta(model.Materials.Get(offset.MaterialIndex))
				}
				if calcMode == pmx.CALC_MODE_MULTIPLICATION {
					delta.Mul(offset, ratio)
				} else {
					delta.Add(offset, ratio)
				}
				deltas.Data[offset.MaterialIndex] = delta
			}
		}
	}

	return deltas
}
