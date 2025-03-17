package delta

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

type VertexMorphDeltas struct {
	data     map[int]*VertexMorphDelta
	vertices *pmx.Vertices
}

func NewVertexMorphDeltas(vertices *pmx.Vertices) *VertexMorphDeltas {
	return &VertexMorphDeltas{
		data:     make(map[int]*VertexMorphDelta),
		vertices: vertices,
	}
}

func (vertexMorphDeltas *VertexMorphDeltas) Length() int {
	return len(vertexMorphDeltas.data)
}

func (vertexMorphDeltas *VertexMorphDeltas) Get(index int) *VertexMorphDelta {
	if v, ok := vertexMorphDeltas.data[index]; ok {
		return v
	}
	return nil
}

func (vertexMorphDeltas *VertexMorphDeltas) Update(v *VertexMorphDelta) {
	vertexMorphDeltas.data[v.Index] = v
}

// ForEach は全ての頂点モーフデルタをコールバック関数に渡します
func (vd *VertexMorphDeltas) ForEach(callback func(index int, value *VertexMorphDelta)) {
	for i, v := range vd.data {
		callback(i, v)
	}
}

type WireVertexMorphDeltas struct {
	*VertexMorphDeltas
}

func NewWireVertexMorphDeltas(vertices *pmx.Vertices) *WireVertexMorphDeltas {
	return &WireVertexMorphDeltas{
		VertexMorphDeltas: NewVertexMorphDeltas(vertices),
	}
}

// ----------------------------

type BoneMorphDeltas struct {
	data  map[int]*BoneMorphDelta
	bones *pmx.Bones
}

func NewBoneMorphDeltas(bones *pmx.Bones) *BoneMorphDeltas {
	return &BoneMorphDeltas{
		data:  make(map[int]*BoneMorphDelta),
		bones: bones,
	}
}

func (boneMorphDeltas *BoneMorphDeltas) Length() int {
	return len(boneMorphDeltas.data)
}

func (boneMorphDeltas *BoneMorphDeltas) Get(boneIndex int) *BoneMorphDelta {
	if v, ok := boneMorphDeltas.data[boneIndex]; ok {
		return v
	}
	return nil
}

func (boneMorphDeltas *BoneMorphDeltas) Update(b *BoneMorphDelta) {
	boneMorphDeltas.data[b.BoneIndex] = b
}

// ForEach は全てのボーンモーフデルタをコールバック関数に渡します
func (bd *BoneMorphDeltas) ForEach(callback func(index int, value *BoneMorphDelta)) {
	for i, v := range bd.data {
		callback(i, v)
	}
}

// ----------------------------

type MaterialMorphDeltas struct {
	data      map[int]*MaterialMorphDelta
	materials *pmx.Materials
}

func NewMaterialMorphDeltas(materials *pmx.Materials) *MaterialMorphDeltas {
	deltas := make(map[int]*MaterialMorphDelta)
	materials.ForEach(func(i int, m *pmx.Material) {
		deltas[i] = NewMaterialMorphDelta(m)
	})

	return &MaterialMorphDeltas{
		data:      deltas,
		materials: materials,
	}
}

func (materialMorphDeltas *MaterialMorphDeltas) Length() int {
	return len(materialMorphDeltas.data)
}

func (materialMorphDeltas *MaterialMorphDeltas) Get(index int) *MaterialMorphDelta {
	if v, ok := materialMorphDeltas.data[index]; ok {
		return v
	}
	return nil
}

func (materialMorphDeltas *MaterialMorphDeltas) Update(m *MaterialMorphDelta) {
	materialMorphDeltas.data[m.Index()] = m
}

// ForEach は全ての材質モーフデルタをコールバック関数に渡します
func (md *MaterialMorphDeltas) ForEach(callback func(index int, value *MaterialMorphDelta)) {
	for i, v := range md.data {
		callback(i, v)
	}
}

type MorphDeltas struct {
	Vertices  *VertexMorphDeltas
	Bones     *BoneMorphDeltas
	Materials *MaterialMorphDeltas
}

func NewMorphDeltas(vertices *pmx.Vertices, materials *pmx.Materials, bones *pmx.Bones) *MorphDeltas {
	return &MorphDeltas{
		Vertices:  NewVertexMorphDeltas(vertices),
		Bones:     NewBoneMorphDeltas(bones),
		Materials: NewMaterialMorphDeltas(materials),
	}
}
