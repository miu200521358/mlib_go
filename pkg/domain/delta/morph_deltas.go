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

func (vertexMorphDeltas *VertexMorphDeltas) Length() int {
	return len(vertexMorphDeltas.data)
}

func (vertexMorphDeltas *VertexMorphDeltas) Get(index int) *VertexMorphDelta {
	return vertexMorphDeltas.data[index]
}

func (vertexMorphDeltas *VertexMorphDeltas) Update(v *VertexMorphDelta) {
	vertexMorphDeltas.data[v.Index] = v
}

func (vertexMorphDeltas *VertexMorphDeltas) Iterator() <-chan struct {
	Index int
	Value *VertexMorphDelta
} {
	ch := make(chan struct {
		Index int
		Value *VertexMorphDelta
	})
	go func() {
		for i, v := range vertexMorphDeltas.data {
			ch <- struct {
				Index int
				Value *VertexMorphDelta
			}{Index: i, Value: v}
		}
		close(ch)
	}()
	return ch
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

func (boneMorphDeltas *BoneMorphDeltas) Length() int {
	return len(boneMorphDeltas.data)
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
	data []*MaterialMorphDelta
}

func NewMaterialMorphDeltas(materials *pmx.Materials) *MaterialMorphDeltas {
	deltas := make([]*MaterialMorphDelta, materials.Length())
	for m := range materials.Iterator() {
		deltas[m.Index] = NewMaterialMorphDelta(m.Value)
	}

	return &MaterialMorphDeltas{
		data: deltas,
	}
}

func (materialMorphDeltas *MaterialMorphDeltas) Length() int {
	return len(materialMorphDeltas.data)
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

func NewMorphDeltas(materials *pmx.Materials, bones *pmx.Bones) *MorphDeltas {
	return &MorphDeltas{
		Vertices:  NewVertexMorphDeltas(),
		Bones:     NewBoneMorphDeltas(bones),
		Materials: NewMaterialMorphDeltas(materials),
	}
}
