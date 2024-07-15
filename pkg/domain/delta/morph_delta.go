package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
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

type WireVertexMorphDeltas struct {
	*VertexMorphDeltas
}

func NewWireVertexMorphDeltas() *WireVertexMorphDeltas {
	return &WireVertexMorphDeltas{
		VertexMorphDeltas: NewVertexMorphDeltas(),
	}
}

type BoneMorphDelta struct {
	BoneIndex int
	// *vmd.BoneFrame
	*MorphBoneDelta
}

func NewBoneMorphDelta(boneIndex int) *BoneMorphDelta {
	return &BoneMorphDelta{
		BoneIndex: boneIndex,
		// BoneFrame:       NewBoneFrame(boneIndex),
		MorphBoneDelta: NewMorphBoneDelta(),
	}
}

type BoneMorphDeltas struct {
	Data []*BoneMorphDelta
}

func NewBoneMorphDeltas(bones *pmx.Bones) *BoneMorphDeltas {
	return &BoneMorphDeltas{
		Data: make([]*BoneMorphDelta, bones.Len()),
	}
}

func (bts *BoneMorphDeltas) Get(boneIndex int) *BoneMorphDelta {
	if boneIndex < 0 || boneIndex >= len(bts.Data) {
		return nil
	}

	return bts.Data[boneIndex]
}

func (bts *BoneMorphDeltas) Update(b *BoneMorphDelta) {
	bts.Data[b.BoneIndex] = b
}

type MaterialMorphDelta struct {
	*pmx.Material
	AddMaterial *pmx.Material
	MulMaterial *pmx.Material
	AddRatios   *pmx.Material
	MulRatios   *pmx.Material
}

func NewMaterialMorphDelta(m *pmx.Material) *MaterialMorphDelta {
	mm := &MaterialMorphDelta{
		Material: m.Copy().(*pmx.Material),
		AddMaterial: &pmx.Material{
			Diffuse:             &mmath.MVec4{},
			Specular:            &mmath.MVec4{},
			Ambient:             &mmath.MVec3{},
			Edge:                &mmath.MVec4{},
			EdgeSize:            0,
			TextureFactor:       &mmath.MVec4{},
			SphereTextureFactor: &mmath.MVec4{},
			ToonTextureFactor:   &mmath.MVec4{},
		},
		MulMaterial: &pmx.Material{
			Diffuse:             &mmath.MVec4{1, 1, 1, 1},
			Specular:            &mmath.MVec4{1, 1, 1, 1},
			Ambient:             &mmath.MVec3{1, 1, 1},
			Edge:                &mmath.MVec4{1, 1, 1, 1},
			EdgeSize:            1,
			TextureFactor:       &mmath.MVec4{1, 1, 1, 1},
			SphereTextureFactor: &mmath.MVec4{1, 1, 1, 1},
			ToonTextureFactor:   &mmath.MVec4{1, 1, 1, 1},
		},
		AddRatios: &pmx.Material{
			Diffuse:             &mmath.MVec4{},
			Specular:            &mmath.MVec4{},
			Ambient:             &mmath.MVec3{},
			Edge:                &mmath.MVec4{},
			EdgeSize:            0,
			TextureFactor:       &mmath.MVec4{},
			SphereTextureFactor: &mmath.MVec4{},
			ToonTextureFactor:   &mmath.MVec4{},
		},
		MulRatios: &pmx.Material{
			Diffuse:             &mmath.MVec4{1, 1, 1, 1},
			Specular:            &mmath.MVec4{1, 1, 1, 1},
			Ambient:             &mmath.MVec3{1, 1, 1},
			Edge:                &mmath.MVec4{1, 1, 1, 1},
			EdgeSize:            1,
			TextureFactor:       &mmath.MVec4{1, 1, 1, 1},
			SphereTextureFactor: &mmath.MVec4{1, 1, 1, 1},
			ToonTextureFactor:   &mmath.MVec4{1, 1, 1, 1},
		},
	}

	if mm.Material.TextureFactor == nil {
		mm.Material.TextureFactor = &mmath.MVec4{1, 1, 1, 1}
	}
	if mm.Material.SphereTextureFactor == nil {
		mm.Material.SphereTextureFactor = &mmath.MVec4{1, 1, 1, 1}
	}
	if mm.Material.ToonTextureFactor == nil {
		mm.Material.ToonTextureFactor = &mmath.MVec4{1, 1, 1, 1}
	}

	return mm
}

func (md *MaterialMorphDelta) Add(m *pmx.MaterialMorphOffset, ratio float64) {
	md.AddMaterial.Diffuse.Add(m.Diffuse)
	md.AddMaterial.Specular.Add(m.Specular)
	md.AddMaterial.Ambient.Add(m.Ambient)
	md.AddMaterial.Edge.Add(m.Edge)
	md.AddMaterial.EdgeSize += m.EdgeSize
	md.AddMaterial.TextureFactor.Add(m.TextureFactor)
	md.AddMaterial.SphereTextureFactor.Add(m.SphereTextureFactor)
	md.AddMaterial.ToonTextureFactor.Add(m.ToonTextureFactor)
	if m.Diffuse.GetX() != 0 {
		md.AddRatios.Diffuse.AddX(ratio)
	}
	if m.Diffuse.GetY() != 0 {
		md.AddRatios.Diffuse.AddY(ratio)
	}
	if m.Diffuse.GetZ() != 0 {
		md.AddRatios.Diffuse.AddZ(ratio)
	}
	if m.Diffuse.GetW() != 0 {
		md.AddRatios.Diffuse.AddW(ratio)
	}
	if m.Specular.GetX() != 0 {
		md.AddRatios.Specular.AddX(ratio)
	}
	if m.Specular.GetY() != 0 {
		md.AddRatios.Specular.AddY(ratio)
	}
	if m.Specular.GetZ() != 0 {
		md.AddRatios.Specular.AddZ(ratio)
	}
	if m.Specular.GetW() != 0 {
		md.AddRatios.Specular.AddW(ratio)
	}
	if m.Ambient.GetX() != 0 {
		md.AddRatios.Ambient.AddX(ratio)
	}
	if m.Ambient.GetY() != 0 {
		md.AddRatios.Ambient.AddY(ratio)
	}
	if m.Ambient.GetZ() != 0 {
		md.AddRatios.Ambient.AddZ(ratio)
	}
	if m.Edge.GetX() != 0 {
		md.AddRatios.Edge.AddX(ratio)
	}
	if m.Edge.GetY() != 0 {
		md.AddRatios.Edge.AddY(ratio)
	}
	if m.Edge.GetZ() != 0 {
		md.AddRatios.Edge.AddZ(ratio)
	}
	if m.Edge.GetW() != 0 {
		md.AddRatios.Edge.AddW(ratio)
	}
	if m.EdgeSize != 0 {
		md.AddRatios.EdgeSize += ratio
	}
	if m.TextureFactor.GetX() != 0 {
		md.AddRatios.TextureFactor.AddX(ratio)
	}
	if m.TextureFactor.GetY() != 0 {
		md.AddRatios.TextureFactor.AddY(ratio)
	}
	if m.TextureFactor.GetZ() != 0 {
		md.AddRatios.TextureFactor.AddZ(ratio)
	}
	if m.TextureFactor.GetW() != 0 {
		md.AddRatios.TextureFactor.AddW(ratio)
	}
	if m.SphereTextureFactor.GetX() != 0 {
		md.AddRatios.SphereTextureFactor.AddX(ratio)
	}
	if m.SphereTextureFactor.GetY() != 0 {
		md.AddRatios.SphereTextureFactor.AddY(ratio)
	}
	if m.SphereTextureFactor.GetZ() != 0 {
		md.AddRatios.SphereTextureFactor.AddZ(ratio)
	}
	if m.SphereTextureFactor.GetW() != 0 {
		md.AddRatios.SphereTextureFactor.AddW(ratio)
	}
	if m.ToonTextureFactor.GetX() != 0 {
		md.AddRatios.ToonTextureFactor.AddX(ratio)
	}
	if m.ToonTextureFactor.GetY() != 0 {
		md.AddRatios.ToonTextureFactor.AddY(ratio)
	}
	if m.ToonTextureFactor.GetZ() != 0 {
		md.AddRatios.ToonTextureFactor.AddZ(ratio)
	}
	if m.ToonTextureFactor.GetW() != 0 {
		md.AddRatios.ToonTextureFactor.AddW(ratio)
	}
}

func (md *MaterialMorphDelta) Mul(m *pmx.MaterialMorphOffset, ratio float64) {
	if m.Diffuse.GetX() != 1 {
		md.MulMaterial.Diffuse.MulX(md.Material.Diffuse.GetX() - m.Diffuse.GetX())
		md.MulRatios.Diffuse.MulX(1 - ratio)
	}
	if m.Diffuse.GetY() != 1 {
		md.MulMaterial.Diffuse.MulY(md.Material.Diffuse.GetY() - m.Diffuse.GetY())
		md.MulRatios.Diffuse.MulY(1 - ratio)
	}
	if m.Diffuse.GetZ() != 1 {
		md.MulMaterial.Diffuse.MulZ(md.Material.Diffuse.GetZ() - m.Diffuse.GetZ())
		md.MulRatios.Diffuse.MulZ(1 - ratio)
	}
	if m.Diffuse.GetW() != 1 {
		md.MulMaterial.Diffuse.MulW(md.Material.Diffuse.GetW() - m.Diffuse.GetW())
		md.MulRatios.Diffuse.MulW(1 - ratio)
	}
	if m.Specular.GetX() != 1 {
		md.MulMaterial.Specular.MulX(md.Material.Specular.GetX() - m.Specular.GetX())
		md.MulRatios.Specular.MulX(1 - ratio)
	}
	if m.Specular.GetY() != 1 {
		md.MulMaterial.Specular.MulY(md.Material.Specular.GetY() - m.Specular.GetY())
		md.MulRatios.Specular.MulY(1 - ratio)
	}
	if m.Specular.GetZ() != 1 {
		md.MulMaterial.Specular.MulZ(md.Material.Specular.GetZ() - m.Specular.GetZ())
		md.MulRatios.Specular.MulZ(1 - ratio)
	}
	if m.Specular.GetW() != 1 {
		md.MulMaterial.Specular.MulW(md.Material.Specular.GetW() - m.Specular.GetW())
		md.MulRatios.Specular.MulW(1 - ratio)
	}
	if m.Ambient.GetX() != 1 {
		md.MulMaterial.Ambient.MulX(md.Material.Ambient.GetX() - m.Ambient.GetX())
		md.MulRatios.Ambient.MulX(1 - ratio)
	}
	if m.Ambient.GetY() != 1 {
		md.MulMaterial.Ambient.MulY(md.Material.Ambient.GetY() - m.Ambient.GetY())
		md.MulRatios.Ambient.MulY(1 - ratio)
	}
	if m.Ambient.GetZ() != 1 {
		md.MulMaterial.Ambient.MulZ(md.Material.Ambient.GetZ() - m.Ambient.GetZ())
		md.MulRatios.Ambient.MulZ(1 - ratio)
	}
	if m.Edge.GetX() != 1 {
		md.MulMaterial.Edge.MulX(md.Material.Edge.GetX() - m.Edge.GetX())
		md.MulRatios.Edge.MulX(1 - ratio)
	}
	if m.Edge.GetY() != 1 {
		md.MulMaterial.Edge.MulY(md.Material.Edge.GetY() - m.Edge.GetY())
		md.MulRatios.Edge.MulY(1 - ratio)
	}
	if m.Edge.GetZ() != 1 {
		md.MulMaterial.Edge.MulZ(md.Material.Edge.GetZ() - m.Edge.GetZ())
		md.MulRatios.Edge.MulZ(1 - ratio)
	}
	if m.Edge.GetW() != 1 {
		md.MulMaterial.Edge.MulW(md.Material.Edge.GetW() - m.Edge.GetW())
		md.MulRatios.Edge.MulW(1 - ratio)
	}
	if m.EdgeSize != 1 {
		md.MulMaterial.EdgeSize *= md.Material.EdgeSize - m.EdgeSize
		md.MulRatios.EdgeSize *= ratio
	}
	if m.TextureFactor.GetX() != 1 {
		md.MulMaterial.TextureFactor.MulX(md.Material.TextureFactor.GetX() - m.TextureFactor.GetX())
		md.MulRatios.TextureFactor.MulX(1 - ratio)
	}
	if m.TextureFactor.GetY() != 1 {
		md.MulMaterial.TextureFactor.MulY(md.Material.TextureFactor.GetY() - m.TextureFactor.GetY())
		md.MulRatios.TextureFactor.MulY(1 - ratio)
	}
	if m.TextureFactor.GetZ() != 1 {
		md.MulMaterial.TextureFactor.MulZ(md.Material.TextureFactor.GetZ() - m.TextureFactor.GetZ())
		md.MulRatios.TextureFactor.MulZ(1 - ratio)
	}
	if m.TextureFactor.GetW() != 1 {
		md.MulMaterial.TextureFactor.MulW(md.Material.TextureFactor.GetW() - m.TextureFactor.GetW())
		md.MulRatios.TextureFactor.MulW(1 - ratio)
	}
	if m.SphereTextureFactor.GetX() != 1 {
		md.MulMaterial.SphereTextureFactor.MulX(md.Material.SphereTextureFactor.GetX() - m.SphereTextureFactor.GetX())
		md.MulRatios.SphereTextureFactor.MulX(1 - ratio)
	}
	if m.SphereTextureFactor.GetY() != 1 {
		md.MulMaterial.SphereTextureFactor.MulY(md.Material.SphereTextureFactor.GetY() - m.SphereTextureFactor.GetY())
		md.MulRatios.SphereTextureFactor.MulY(1 - ratio)
	}
	if m.SphereTextureFactor.GetZ() != 1 {
		md.MulMaterial.SphereTextureFactor.MulZ(md.Material.SphereTextureFactor.GetZ() - m.SphereTextureFactor.GetZ())
		md.MulRatios.SphereTextureFactor.MulZ(1 - ratio)
	}
	if m.SphereTextureFactor.GetW() != 1 {
		md.MulMaterial.SphereTextureFactor.MulW(md.Material.SphereTextureFactor.GetW() - m.SphereTextureFactor.GetW())
		md.MulRatios.SphereTextureFactor.MulW(1 - ratio)
	}
	if m.ToonTextureFactor.GetX() != 1 {
		md.MulMaterial.ToonTextureFactor.MulX(md.Material.ToonTextureFactor.GetX() - m.ToonTextureFactor.GetX())
		md.MulRatios.ToonTextureFactor.MulX(1 - ratio)
	}
	if m.ToonTextureFactor.GetY() != 1 {
		md.MulMaterial.ToonTextureFactor.MulY(md.Material.ToonTextureFactor.GetY() - m.ToonTextureFactor.GetY())
		md.MulRatios.ToonTextureFactor.MulY(1 - ratio)
	}
	if m.ToonTextureFactor.GetZ() != 1 {
		md.MulMaterial.ToonTextureFactor.MulZ(md.Material.ToonTextureFactor.GetZ() - m.ToonTextureFactor.GetZ())
		md.MulRatios.ToonTextureFactor.MulZ(1 - ratio)
	}
	if m.ToonTextureFactor.GetW() != 1 {
		md.MulMaterial.ToonTextureFactor.MulW(md.Material.ToonTextureFactor.GetW() - m.ToonTextureFactor.GetW())
		md.MulRatios.ToonTextureFactor.MulW(1 - ratio)
	}
}

type MaterialMorphDeltas struct {
	Data []*MaterialMorphDelta
}

func NewMaterialMorphDeltas(materials *pmx.Materials) *MaterialMorphDeltas {
	deltas := make([]*MaterialMorphDelta, len(materials.Data))
	for i := range materials.Data {
		m := materials.Get(i)
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

func NewMorphDeltas(materials *pmx.Materials, bones *pmx.Bones) *MorphDeltas {
	return &MorphDeltas{
		Vertices:  NewVertexMorphDeltas(),
		Bones:     NewBoneMorphDeltas(bones),
		Materials: NewMaterialMorphDeltas(materials),
	}
}
