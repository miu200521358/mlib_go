package delta

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

type VertexMorphDeltas struct {
	data map[int]*VertexMorphDelta
}

func NewVertexMorphDeltas() *VertexMorphDeltas {
	return &VertexMorphDeltas{
		data: make(map[int]*VertexMorphDelta),
	}
}

func (vertexMorphDeltas *VertexMorphDeltas) Get(index int) *VertexMorphDelta {
	return vertexMorphDeltas.data[index]
}

func (vertexMorphDeltas *VertexMorphDeltas) Update(v *VertexMorphDelta) {
	vertexMorphDeltas.data[v.Index] = v
}

type WireVertexMorphDeltas struct {
	*VertexMorphDeltas
}

func NewWireVertexMorphDeltas() *WireVertexMorphDeltas {
	return &WireVertexMorphDeltas{
		VertexMorphDeltas: NewVertexMorphDeltas(),
	}
}

// ----------------------------

type BoneMorphDeltas struct {
	data []*BoneMorphDelta
}

func NewBoneMorphDeltas(bones *pmx.Bones) *BoneMorphDeltas {
	return &BoneMorphDeltas{
		data: make([]*BoneMorphDelta, bones.Length()),
	}
}

func (boneMorphDeltas *BoneMorphDeltas) Get(boneIndex int) *BoneMorphDelta {
	if boneIndex < 0 || boneIndex >= len(boneMorphDeltas.data) {
		return nil
	}

	return boneMorphDeltas.data[boneIndex]
}

func (boneMorphDeltas *BoneMorphDeltas) Update(b *BoneMorphDelta) {
	boneMorphDeltas.data[b.BoneIndex] = b
}

// ----------------------------

type MaterialMorphDeltas struct {
	data []*MaterialMorphDelta
}

func NewMaterialMorphDeltas(materials *pmx.Materials) *MaterialMorphDeltas {
	deltas := make([]*MaterialMorphDelta, materials.Length())
	for m := range materials.Iterator() {
		deltas[m.Index()] = NewMaterialMorphDelta(m)
	}

	return &MaterialMorphDeltas{
		data: deltas,
	}
}

func (materialMorphDeltas *MaterialMorphDeltas) Get(index int) *MaterialMorphDelta {
	if index < 0 || index >= len(materialMorphDeltas.data) {
		return nil
	}

	return materialMorphDeltas.data[index]
}

func (materialMorphDeltas *MaterialMorphDeltas) Update(m *MaterialMorphDelta) {
	materialMorphDeltas.data[m.Index()] = m
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
