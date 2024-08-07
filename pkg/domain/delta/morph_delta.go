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

func (vertexMorphDeltas *VertexMorphDeltas) Get(index int) *VertexMorphDelta {
	return vertexMorphDeltas.Data[index]
}

func (vertexMorphDeltas *VertexMorphDeltas) Update(v *VertexMorphDelta) {
	vertexMorphDeltas.Data[v.Index] = v
}

type BoneMorphDelta struct {
	BoneIndex     int
	FramePosition *mmath.MVec3       // キーフレ位置の変動量
	FrameRotation *mmath.MQuaternion // キーフレ回転の変動量
	FrameScale    *mmath.MVec3       // キーフレスケールの変動量
}

func NewBoneMorphDelta(boneIndex int) *BoneMorphDelta {
	return &BoneMorphDelta{
		BoneIndex: boneIndex,
	}
}

func (boneMorphDelta *BoneMorphDelta) FilledMorphPosition() *mmath.MVec3 {
	if boneMorphDelta.FramePosition == nil {
		boneMorphDelta.FramePosition = mmath.NewMVec3()
	}
	return boneMorphDelta.FramePosition
}

func (boneMorph *BoneMorphDelta) FilledMorphRotation() *mmath.MQuaternion {
	if boneMorph.FrameRotation == nil {
		boneMorph.FrameRotation = mmath.NewMQuaternion()
	}
	return boneMorph.FrameRotation
}

func (boneMorphDelta *BoneMorphDelta) FilledMorphScale() *mmath.MVec3 {
	if boneMorphDelta.FrameScale == nil {
		boneMorphDelta.FrameScale = mmath.NewMVec3()
	}
	return boneMorphDelta.FrameScale
}

func (boneMorphDelta *BoneMorphDelta) Copy() *BoneMorphDelta {
	return &BoneMorphDelta{
		FramePosition: boneMorphDelta.FilledMorphPosition().Copy(),
		FrameRotation: boneMorphDelta.FilledMorphRotation().Copy(),
		FrameScale:    boneMorphDelta.FilledMorphScale().Copy(),
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

func (boneMorphDeltas *BoneMorphDeltas) Get(boneIndex int) *BoneMorphDelta {
	if boneIndex < 0 || boneIndex >= len(boneMorphDeltas.Data) {
		return nil
	}

	return boneMorphDeltas.Data[boneIndex]
}

func (boneMorphDeltas *BoneMorphDeltas) Update(b *BoneMorphDelta) {
	boneMorphDeltas.Data[b.BoneIndex] = b
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
			Diffuse:  &mmath.MVec4{},
			Specular: &mmath.MVec4{},
			Ambient:  &mmath.MVec3{},
			Edge:     &mmath.MVec4{},
			EdgeSize: 0,
		},
		MulMaterial: &pmx.Material{
			Diffuse:  &mmath.MVec4{X: 1, Y: 1, Z: 1, W: 1},
			Specular: &mmath.MVec4{X: 1, Y: 1, Z: 1, W: 1},
			Ambient:  &mmath.MVec3{X: 1, Y: 1, Z: 1},
			Edge:     &mmath.MVec4{X: 1, Y: 1, Z: 1, W: 1},
			EdgeSize: 1,
		},
		AddRatios: &pmx.Material{
			Diffuse:  &mmath.MVec4{},
			Specular: &mmath.MVec4{},
			Ambient:  &mmath.MVec3{},
			Edge:     &mmath.MVec4{},
			EdgeSize: 0,
		},
		MulRatios: &pmx.Material{
			Diffuse:  &mmath.MVec4{X: 1, Y: 1, Z: 1, W: 1},
			Specular: &mmath.MVec4{X: 1, Y: 1, Z: 1, W: 1},
			Ambient:  &mmath.MVec3{X: 1, Y: 1, Z: 1},
			Edge:     &mmath.MVec4{X: 1, Y: 1, Z: 1, W: 1},
			EdgeSize: 1,
		},
	}

	return mm
}

func (materialMorphDelta *MaterialMorphDelta) Add(m *pmx.MaterialMorphOffset, ratio float64) {
	materialMorphDelta.AddMaterial.Diffuse.Add(m.Diffuse)
	materialMorphDelta.AddMaterial.Specular.Add(m.Specular)
	materialMorphDelta.AddMaterial.Ambient.Add(m.Ambient)
	materialMorphDelta.AddMaterial.Edge.Add(m.Edge)
	materialMorphDelta.AddMaterial.EdgeSize += m.EdgeSize
	if m.Diffuse.X != 0 {
		materialMorphDelta.AddRatios.Diffuse.X += ratio
	}
	if m.Diffuse.Y != 0 {
		materialMorphDelta.AddRatios.Diffuse.Y += ratio
	}
	if m.Diffuse.Z != 0 {
		materialMorphDelta.AddRatios.Diffuse.Z += ratio
	}
	if m.Diffuse.W != 0 {
		materialMorphDelta.AddRatios.Diffuse.W += ratio
	}
	if m.Specular.X != 0 {
		materialMorphDelta.AddRatios.Specular.X += ratio
	}
	if m.Specular.Y != 0 {
		materialMorphDelta.AddRatios.Specular.Y += ratio
	}
	if m.Specular.Z != 0 {
		materialMorphDelta.AddRatios.Specular.Z += ratio
	}
	if m.Specular.W != 0 {
		materialMorphDelta.AddRatios.Specular.W += ratio
	}
	if m.Ambient.X != 0 {
		materialMorphDelta.AddRatios.Ambient.X += ratio
	}
	if m.Ambient.Y != 0 {
		materialMorphDelta.AddRatios.Ambient.Y += ratio
	}
	if m.Ambient.Z != 0 {
		materialMorphDelta.AddRatios.Ambient.Z += ratio
	}
	if m.Edge.X != 0 {
		materialMorphDelta.AddRatios.Edge.X += ratio
	}
	if m.Edge.Y != 0 {
		materialMorphDelta.AddRatios.Edge.Y += ratio
	}
	if m.Edge.Z != 0 {
		materialMorphDelta.AddRatios.Edge.Z += ratio
	}
	if m.Edge.W != 0 {
		materialMorphDelta.AddRatios.Edge.W += ratio
	}
	if m.EdgeSize != 0 {
		materialMorphDelta.AddRatios.EdgeSize += ratio
	}
}

func (materialMorphDelta *MaterialMorphDelta) Mul(m *pmx.MaterialMorphOffset, ratio float64) {
	if m.Diffuse.X != 1 {
		materialMorphDelta.MulMaterial.Diffuse.X *= materialMorphDelta.Material.Diffuse.X - m.Diffuse.X
		materialMorphDelta.MulRatios.Diffuse.X *= 1 - ratio
	}
	if m.Diffuse.Y != 1 {
		materialMorphDelta.MulMaterial.Diffuse.Y *= materialMorphDelta.Material.Diffuse.Y - m.Diffuse.Y
		materialMorphDelta.MulRatios.Diffuse.Y *= 1 - ratio
	}
	if m.Diffuse.Z != 1 {
		materialMorphDelta.MulMaterial.Diffuse.Z *= materialMorphDelta.Material.Diffuse.Z - m.Diffuse.Z
		materialMorphDelta.MulRatios.Diffuse.Z *= 1 - ratio
	}
	if m.Diffuse.W != 1 {
		materialMorphDelta.MulMaterial.Diffuse.W *= materialMorphDelta.Material.Diffuse.W - m.Diffuse.W
		materialMorphDelta.MulRatios.Diffuse.W *= 1 - ratio
	}
	if m.Specular.X != 1 {
		materialMorphDelta.MulMaterial.Specular.X *= materialMorphDelta.Material.Specular.X - m.Specular.X
		materialMorphDelta.MulRatios.Specular.X *= 1 - ratio
	}
	if m.Specular.Y != 1 {
		materialMorphDelta.MulMaterial.Specular.Y *= materialMorphDelta.Material.Specular.Y - m.Specular.Y
		materialMorphDelta.MulRatios.Specular.Y *= 1 - ratio
	}
	if m.Specular.Z != 1 {
		materialMorphDelta.MulMaterial.Specular.Z *= materialMorphDelta.Material.Specular.Z - m.Specular.Z
		materialMorphDelta.MulRatios.Specular.Z *= 1 - ratio
	}
	if m.Specular.W != 1 {
		materialMorphDelta.MulMaterial.Specular.W *= materialMorphDelta.Material.Specular.W - m.Specular.W
		materialMorphDelta.MulRatios.Specular.W *= 1 - ratio
	}
	if m.Ambient.X != 1 {
		materialMorphDelta.MulMaterial.Ambient.X *= materialMorphDelta.Material.Ambient.X - m.Ambient.X
		materialMorphDelta.MulRatios.Ambient.X *= 1 - ratio
	}
	if m.Ambient.Y != 1 {
		materialMorphDelta.MulMaterial.Ambient.Y *= materialMorphDelta.Material.Ambient.Y - m.Ambient.Y
		materialMorphDelta.MulRatios.Ambient.Y *= 1 - ratio
	}
	if m.Ambient.Z != 1 {
		materialMorphDelta.MulMaterial.Ambient.Z *= materialMorphDelta.Material.Ambient.Z - m.Ambient.Z
		materialMorphDelta.MulRatios.Ambient.Z *= 1 - ratio
	}
	if m.Edge.X != 1 {
		materialMorphDelta.MulMaterial.Edge.X *= materialMorphDelta.Material.Edge.X - m.Edge.X
		materialMorphDelta.MulRatios.Edge.X *= 1 - ratio
	}
	if m.Edge.Y != 1 {
		materialMorphDelta.MulMaterial.Edge.Y *= materialMorphDelta.Material.Edge.Y - m.Edge.Y
		materialMorphDelta.MulRatios.Edge.Y *= 1 - ratio
	}
	if m.Edge.Z != 1 {
		materialMorphDelta.MulMaterial.Edge.Z *= materialMorphDelta.Material.Edge.Z - m.Edge.Z
		materialMorphDelta.MulRatios.Edge.Z *= 1 - ratio
	}
	if m.Edge.W != 1 {
		materialMorphDelta.MulMaterial.Edge.W *= materialMorphDelta.Material.Edge.W - m.Edge.W
		materialMorphDelta.MulRatios.Edge.W *= 1 - ratio
	}
	if m.EdgeSize != 1 {
		materialMorphDelta.MulMaterial.EdgeSize *= materialMorphDelta.Material.EdgeSize - m.EdgeSize
		materialMorphDelta.MulRatios.EdgeSize *= ratio
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
