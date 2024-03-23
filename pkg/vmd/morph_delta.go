package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"

)

type VertexMorphDelta struct {
	Position      *mmath.MVec3
	Uv            *mmath.MVec2
	Uv1           *mmath.MVec2
	AfterPosition *mmath.MVec3
}

func NewVertexMorphDelta() *VertexMorphDelta {
	return &VertexMorphDelta{
		Position:      mmath.NewMVec3(),
		Uv:            mmath.NewMVec2(),
		Uv1:           mmath.NewMVec2(),
		AfterPosition: mmath.NewMVec3(),
	}
}

func (md *VertexMorphDelta) GL() []float32 {
	p := md.Position.GL()
	ap := md.AfterPosition.GL()
	// UVは符号関係ないのでそのまま取得する
	return []float32{
		p[0], p[1], p[2],
		float32(md.Uv.GetX()), float32(md.Uv.GetY()),
		float32(md.Uv1.GetX()), float32(md.Uv1.GetY()),
		ap[0], ap[1], ap[2],
	}
}

type VertexMorphDeltas struct {
	Data []*VertexMorphDelta
}

func NewVertexMorphDeltas(vertexCount int) *VertexMorphDeltas {
	deltas := make([]*VertexMorphDelta, vertexCount)
	for i := 0; i < vertexCount; i++ {
		deltas[i] = NewVertexMorphDelta()
	}

	return &VertexMorphDeltas{
		Data: deltas,
	}
}

type BoneMorphDelta struct {
	*BoneFrame
}

func NewBoneMorphDelta() *BoneMorphDelta {
	return &BoneMorphDelta{
		BoneFrame: NewBoneFrame(0.0),
	}
}

type BoneMorphDeltas struct {
	Data []*BoneMorphDelta
}

func NewBoneMorphDeltas(boneCount int) *BoneMorphDeltas {
	deltas := make([]*BoneMorphDelta, boneCount)
	for i := 0; i < boneCount; i++ {
		deltas[i] = NewBoneMorphDelta()
	}

	return &BoneMorphDeltas{
		Data: deltas,
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
	md.TextureFactor.Add(m.TextureFactor.MuledScalar(ratio))
	md.SphereTextureFactor.Add(m.SphereTextureFactor.MuledScalar(ratio))
	md.ToonTextureFactor.Add(m.ToonTextureFactor.MuledScalar(ratio))
}

func (md *MaterialMorphDelta) Mul(m *pmx.MaterialMorphOffset, ratio float64) {
	md.Diffuse.Mul(m.Diffuse.MuledScalar(ratio))
	md.Specular.Mul(m.Specular.MuledScalar(ratio))
	md.Ambient.Mul(m.Ambient.MuledScalar(ratio))
	md.Edge.Mul(m.Edge.MuledScalar(ratio))
	md.EdgeSize += m.EdgeSize * ratio
	md.TextureFactor.Mul(m.TextureFactor.MuledScalar(ratio))
	md.SphereTextureFactor.Mul(m.SphereTextureFactor.MuledScalar(ratio))
	md.ToonTextureFactor.Mul(m.ToonTextureFactor.MuledScalar(ratio))
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

func NewMorphDeltas(vertexCount int, boneCount int, materials *pmx.Materials) *MorphDeltas {
	return &MorphDeltas{
		Vertices:  NewVertexMorphDeltas(vertexCount),
		Bones:     NewBoneMorphDeltas(boneCount),
		Materials: NewMaterialMorphDeltas(materials),
	}
}
