package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type VertexMorphDelta struct {
	Index         int
	Position      *mmath.MVec3
	Uv            *mmath.MVec2
	Uv1           *mmath.MVec2
	AfterPosition *mmath.MVec3
}

func NewVertexMorphDelta(index int) *VertexMorphDelta {
	return &VertexMorphDelta{
		Index:         index,
		Position:      nil,
		Uv:            nil,
		Uv1:           nil,
		AfterPosition: nil,
	}
}

type VertexMorphDeltas struct {
	Data map[int]*VertexMorphDelta
}

func NewVertexMorphDeltas() *VertexMorphDeltas {
	return &VertexMorphDeltas{
		Data: make(map[int]*VertexMorphDelta),
	}
}

type BoneMorphDelta struct {
	*BoneFrame
}

func NewBoneMorphDelta(index int) *BoneMorphDelta {
	return &BoneMorphDelta{
		BoneFrame: NewBoneFrame(index),
	}
}

type BoneMorphDeltas struct {
	Data map[int]*BoneMorphDelta
}

func NewBoneMorphDeltas() *BoneMorphDeltas {
	return &BoneMorphDeltas{
		Data: make(map[int]*BoneMorphDelta),
	}
}

type MaterialMorphDelta struct {
	*pmx.Material
}

func NewMaterialMorphDelta(m *pmx.Material) *MaterialMorphDelta {
	return &MaterialMorphDelta{
		Material: m.Copy().(*pmx.Material),
	}
}

func (md *MaterialMorphDelta) Add(m *pmx.MaterialMorphOffset, ratio float64) {
	md.Diffuse.Add(m.Diffuse.MuledScalar(ratio))
	md.Specular.Add(m.Specular.MuledScalar(ratio))
	md.Ambient.Add(m.Ambient.MuledScalar(ratio))
	md.Edge.Add(m.Edge.MuledScalar(ratio))
	md.EdgeSize += m.EdgeSize * ratio
	if md.TextureFactor == nil {
		md.TextureFactor = m.TextureFactor.MuledScalar(ratio)
	} else {
		md.TextureFactor.Add(m.TextureFactor.MuledScalar(ratio))
	}
	if md.SphereTextureFactor == nil {
		md.SphereTextureFactor = m.SphereTextureFactor.MuledScalar(ratio)
	} else {
		md.SphereTextureFactor.Add(m.SphereTextureFactor.MuledScalar(ratio))
	}
	if md.ToonTextureFactor == nil {
		md.ToonTextureFactor = m.ToonTextureFactor.MuledScalar(ratio)
	} else {
		md.ToonTextureFactor.Add(m.ToonTextureFactor.MuledScalar(ratio))
	}
}

func (md *MaterialMorphDelta) Mul(m *pmx.MaterialMorphOffset, ratio float64) {
	md.Diffuse.Mul(m.Diffuse.MuledScalar(ratio))
	md.Specular.Mul(m.Specular.MuledScalar(ratio))
	md.Ambient.Mul(m.Ambient.MuledScalar(ratio))
	md.Edge.Mul(m.Edge.MuledScalar(ratio))
	md.EdgeSize += m.EdgeSize * ratio
	if md.TextureFactor == nil {
		md.TextureFactor = m.TextureFactor.MuledScalar(ratio)
	} else {
		md.TextureFactor.Mul(m.TextureFactor.MuledScalar(ratio))
	}
	if md.SphereTextureFactor == nil {
		md.SphereTextureFactor = m.SphereTextureFactor.MuledScalar(ratio)
	} else {
		md.SphereTextureFactor.Mul(m.SphereTextureFactor.MuledScalar(ratio))
	}
	if md.ToonTextureFactor == nil {
		md.ToonTextureFactor = m.ToonTextureFactor.MuledScalar(ratio)
	} else {
		md.ToonTextureFactor.Mul(m.ToonTextureFactor.MuledScalar(ratio))
	}
}

type MaterialMorphDeltas struct {
	Data []*MaterialMorphDelta
}

func NewMaterialMorphDeltas(materials *pmx.Materials) *MaterialMorphDeltas {
	deltas := make([]*MaterialMorphDelta, len(materials.Data))
	for i, m := range materials.GetSortedData() {
		deltas[i] = NewMaterialMorphDelta(m)
	}

	return &MaterialMorphDeltas{
		Data: deltas,
	}
}

type MorphDeltas struct {
	Vertices  *VertexMorphDeltas
	Bones     *BoneMorphDeltas
	Materials *MaterialMorphDeltas
}

func NewMorphDeltas(materials *pmx.Materials) *MorphDeltas {
	return &MorphDeltas{
		Vertices:  NewVertexMorphDeltas(),
		Bones:     NewBoneMorphDeltas(),
		Materials: NewMaterialMorphDeltas(materials),
	}
}
