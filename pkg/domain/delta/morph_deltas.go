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

// VertexMorphDeltasにForEachメソッドを追加
func (vertexMorphDeltas *VertexMorphDeltas) ForEach(callback func(index int, value *VertexMorphDelta)) {
	for i, v := range vertexMorphDeltas.data {
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

func (boneMorphDeltas *BoneMorphDeltas) Iterator() <-chan struct {
	Index int
	Delta *BoneMorphDelta
} {
	ch := make(chan struct {
		Index int
		Delta *BoneMorphDelta
	})
	go func() {
		for i, b := range boneMorphDeltas.data {
			ch <- struct {
				Index int
				Delta *BoneMorphDelta
			}{Index: i, Delta: b}
		}
		close(ch)
	}()
	return ch
}

// ----------------------------

type MaterialMorphDeltas struct {
	data      map[int]*MaterialMorphDelta
	materials *pmx.Materials
}

func NewMaterialMorphDeltas(materials *pmx.Materials) *MaterialMorphDeltas {
	deltas := make(map[int]*MaterialMorphDelta)
	for m := range materials.Iterator() {
		deltas[m.Index] = NewMaterialMorphDelta(m.Value)
	}

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

func (materialMorphDeltas *MaterialMorphDeltas) Iterator() <-chan struct {
	Index int
	Value *MaterialMorphDelta
} {
	ch := make(chan struct {
		Index int
		Value *MaterialMorphDelta
	})
	go func() {
		for i, m := range materialMorphDeltas.data {
			ch <- struct {
				Index int
				Value *MaterialMorphDelta
			}{Index: i, Value: m}
		}
		close(ch)
	}()
	return ch
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
