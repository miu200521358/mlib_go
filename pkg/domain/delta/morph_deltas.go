// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
)

// MorphDeltas は頂点/ボーン/材質モーフ差分をまとめる。
type MorphDeltas struct {
	vertices  *VertexMorphDeltas
	materials *MaterialMorphDeltas
	bones     *BoneMorphDeltas
}

// NewMorphDeltas はMorphDeltasを生成する。
func NewMorphDeltas(
	vertices *collection.IndexedCollection[*model.Vertex],
	materials *collection.NamedCollection[*model.Material],
	bones *model.BoneCollection,
) *MorphDeltas {
	vertexCount := 0
	materialCount := 0
	boneCount := 0
	if vertices != nil {
		vertexCount = vertices.Len()
	}
	if materials != nil {
		materialCount = materials.Len()
	}
	if bones != nil {
		boneCount = bones.Len()
	}
	return &MorphDeltas{
		vertices:  NewVertexMorphDeltas(vertexCount),
		materials: NewMaterialMorphDeltas(materials, materialCount),
		bones:     NewBoneMorphDeltas(boneCount),
	}
}

// Vertices は頂点モーフ差分を返す。
func (m *MorphDeltas) Vertices() *VertexMorphDeltas {
	if m == nil {
		return nil
	}
	return m.vertices
}

// Materials は材質モーフ差分を返す。
func (m *MorphDeltas) Materials() *MaterialMorphDeltas {
	if m == nil {
		return nil
	}
	return m.materials
}

// Bones はボーンモーフ差分を返す。
func (m *MorphDeltas) Bones() *BoneMorphDeltas {
	if m == nil {
		return nil
	}
	return m.bones
}

// VertexMorphDeltas は頂点モーフ差分の集合を表す。
type VertexMorphDeltas struct {
	data       []*VertexMorphDelta
	hasNonZero bool
}

// Len は要素数を返す。
func (d *VertexMorphDeltas) Len() int {
	if d == nil {
		return 0
	}
	return len(d.data)
}

// NewVertexMorphDeltas はVertexMorphDeltasを生成する。
func NewVertexMorphDeltas(vertexCount int) *VertexMorphDeltas {
	return &VertexMorphDeltas{data: make([]*VertexMorphDelta, vertexCount)}
}

// Get はindexの差分を返す。
func (d *VertexMorphDeltas) Get(index int) *VertexMorphDelta {
	if d == nil || index < 0 || index >= len(d.data) {
		return nil
	}
	return d.data[index]
}

// Update は差分を更新する。
func (d *VertexMorphDeltas) Update(delta *VertexMorphDelta) {
	if d == nil || delta == nil {
		return
	}
	if !delta.IsZero() {
		d.hasNonZero = true
	}
	idx := delta.Index
	if idx < 0 || idx >= len(d.data) {
		return
	}
	d.data[idx] = delta
}

// HasNonZero は非ゼロ差分が存在するかを返す。
func (d *VertexMorphDeltas) HasNonZero() bool {
	if d == nil {
		return false
	}
	return d.hasNonZero
}

// ForEach は全要素を走査する。
func (d *VertexMorphDeltas) ForEach(fn func(index int, delta *VertexMorphDelta) bool) {
	if d == nil || fn == nil {
		return
	}
	for i, v := range d.data {
		if !fn(i, v) {
			return
		}
	}
}

// BoneMorphDeltas はボーンモーフ差分の集合を表す。
type BoneMorphDeltas struct {
	data []*BoneMorphDelta
}

// Len は要素数を返す。
func (d *BoneMorphDeltas) Len() int {
	if d == nil {
		return 0
	}
	return len(d.data)
}

// NewBoneMorphDeltas はBoneMorphDeltasを生成する。
func NewBoneMorphDeltas(boneCount int) *BoneMorphDeltas {
	return &BoneMorphDeltas{data: make([]*BoneMorphDelta, boneCount)}
}

// Get はindexの差分を返す。
func (d *BoneMorphDeltas) Get(index int) *BoneMorphDelta {
	if d == nil || index < 0 || index >= len(d.data) {
		return nil
	}
	return d.data[index]
}

// Update は差分を更新する。
func (d *BoneMorphDeltas) Update(delta *BoneMorphDelta) {
	if d == nil || delta == nil {
		return
	}
	idx := delta.BoneIndex
	if idx < 0 || idx >= len(d.data) {
		return
	}
	d.data[idx] = delta
}

// ForEach は全要素を走査する。
func (d *BoneMorphDeltas) ForEach(fn func(index int, delta *BoneMorphDelta) bool) {
	if d == nil || fn == nil {
		return
	}
	for i, v := range d.data {
		if !fn(i, v) {
			return
		}
	}
}

// MaterialMorphDeltas は材質モーフ差分の集合を表す。
type MaterialMorphDeltas struct {
	data []*MaterialMorphDelta
}

// Len は要素数を返す。
func (d *MaterialMorphDeltas) Len() int {
	if d == nil {
		return 0
	}
	return len(d.data)
}

// NewMaterialMorphDeltas はMaterialMorphDeltasを生成する。
func NewMaterialMorphDeltas(materials *collection.NamedCollection[*model.Material], materialCount int) *MaterialMorphDeltas {
	deltas := &MaterialMorphDeltas{data: make([]*MaterialMorphDelta, materialCount)}
	if materials == nil {
		return deltas
	}
	for i, material := range materials.Values() {
		if material == nil {
			continue
		}
		deltas.data[i] = NewMaterialMorphDelta(material)
	}
	return deltas
}

// Get はindexの差分を返す。
func (d *MaterialMorphDeltas) Get(index int) *MaterialMorphDelta {
	if d == nil || index < 0 || index >= len(d.data) {
		return nil
	}
	return d.data[index]
}

// Update は差分を更新する。
func (d *MaterialMorphDeltas) Update(delta *MaterialMorphDelta) {
	if d == nil || delta == nil {
		return
	}
	idx := delta.Material.Index()
	if idx < 0 || idx >= len(d.data) {
		return
	}
	d.data[idx] = delta
}

// ForEach は全要素を走査する。
func (d *MaterialMorphDeltas) ForEach(fn func(index int, delta *MaterialMorphDelta) bool) {
	if d == nil || fn == nil {
		return
	}
	for i, v := range d.data {
		if !fn(i, v) {
			return
		}
	}
}
