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

func (md *BoneMorphDelta) FilledMorphPosition() *mmath.MVec3 {
	if md.FramePosition == nil {
		md.FramePosition = mmath.NewMVec3()
	}
	return md.FramePosition
}

func (md *BoneMorphDelta) FilledMorphRotation() *mmath.MQuaternion {
	if md.FrameRotation == nil {
		md.FrameRotation = mmath.NewMQuaternion()
	}
	return md.FrameRotation
}

func (md *BoneMorphDelta) FilledMorphScale() *mmath.MVec3 {
	if md.FrameScale == nil {
		md.FrameScale = mmath.NewMVec3()
	}
	return md.FrameScale
}

func (md *BoneMorphDelta) Copy() *BoneMorphDelta {
	return &BoneMorphDelta{
		FramePosition: md.FilledMorphPosition().Copy(),
		FrameRotation: md.FilledMorphRotation().Copy(),
		FrameScale:    md.FilledMorphScale().Copy(),
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

func (md *MaterialMorphDelta) Add(m *pmx.MaterialMorphOffset, ratio float64) {
	md.AddMaterial.Diffuse.Add(m.Diffuse)
	md.AddMaterial.Specular.Add(m.Specular)
	md.AddMaterial.Ambient.Add(m.Ambient)
	md.AddMaterial.Edge.Add(m.Edge)
	md.AddMaterial.EdgeSize += m.EdgeSize
	if m.Diffuse.X != 0 {
		md.AddRatios.Diffuse.X += ratio
	}
	if m.Diffuse.Y != 0 {
		md.AddRatios.Diffuse.Y += ratio
	}
	if m.Diffuse.Z != 0 {
		md.AddRatios.Diffuse.Z += ratio
	}
	if m.Diffuse.W != 0 {
		md.AddRatios.Diffuse.W += ratio
	}
	if m.Specular.X != 0 {
		md.AddRatios.Specular.X += ratio
	}
	if m.Specular.Y != 0 {
		md.AddRatios.Specular.Y += ratio
	}
	if m.Specular.Z != 0 {
		md.AddRatios.Specular.Z += ratio
	}
	if m.Specular.W != 0 {
		md.AddRatios.Specular.W += ratio
	}
	if m.Ambient.X != 0 {
		md.AddRatios.Ambient.X += ratio
	}
	if m.Ambient.Y != 0 {
		md.AddRatios.Ambient.Y += ratio
	}
	if m.Ambient.Z != 0 {
		md.AddRatios.Ambient.Z += ratio
	}
	if m.Edge.X != 0 {
		md.AddRatios.Edge.X += ratio
	}
	if m.Edge.Y != 0 {
		md.AddRatios.Edge.Y += ratio
	}
	if m.Edge.Z != 0 {
		md.AddRatios.Edge.Z += ratio
	}
	if m.Edge.W != 0 {
		md.AddRatios.Edge.W += ratio
	}
	if m.EdgeSize != 0 {
		md.AddRatios.EdgeSize += ratio
	}
}

func (md *MaterialMorphDelta) Mul(m *pmx.MaterialMorphOffset, ratio float64) {
	if m.Diffuse.X != 1 {
		md.MulMaterial.Diffuse.X *= md.Material.Diffuse.X - m.Diffuse.X
		md.MulRatios.Diffuse.X *= 1 - ratio
	}
	if m.Diffuse.Y != 1 {
		md.MulMaterial.Diffuse.Y *= md.Material.Diffuse.Y - m.Diffuse.Y
		md.MulRatios.Diffuse.Y *= 1 - ratio
	}
	if m.Diffuse.Z != 1 {
		md.MulMaterial.Diffuse.Z *= md.Material.Diffuse.Z - m.Diffuse.Z
		md.MulRatios.Diffuse.Z *= 1 - ratio
	}
	if m.Diffuse.W != 1 {
		md.MulMaterial.Diffuse.W *= md.Material.Diffuse.W - m.Diffuse.W
		md.MulRatios.Diffuse.W *= 1 - ratio
	}
	if m.Specular.X != 1 {
		md.MulMaterial.Specular.X *= md.Material.Specular.X - m.Specular.X
		md.MulRatios.Specular.X *= 1 - ratio
	}
	if m.Specular.Y != 1 {
		md.MulMaterial.Specular.Y *= md.Material.Specular.Y - m.Specular.Y
		md.MulRatios.Specular.Y *= 1 - ratio
	}
	if m.Specular.Z != 1 {
		md.MulMaterial.Specular.Z *= md.Material.Specular.Z - m.Specular.Z
		md.MulRatios.Specular.Z *= 1 - ratio
	}
	if m.Specular.W != 1 {
		md.MulMaterial.Specular.W *= md.Material.Specular.W - m.Specular.W
		md.MulRatios.Specular.W *= 1 - ratio
	}
	if m.Ambient.X != 1 {
		md.MulMaterial.Ambient.X *= md.Material.Ambient.X - m.Ambient.X
		md.MulRatios.Ambient.X *= 1 - ratio
	}
	if m.Ambient.Y != 1 {
		md.MulMaterial.Ambient.Y *= md.Material.Ambient.Y - m.Ambient.Y
		md.MulRatios.Ambient.Y *= 1 - ratio
	}
	if m.Ambient.Z != 1 {
		md.MulMaterial.Ambient.Z *= md.Material.Ambient.Z - m.Ambient.Z
		md.MulRatios.Ambient.Z *= 1 - ratio
	}
	if m.Edge.X != 1 {
		md.MulMaterial.Edge.X *= md.Material.Edge.X - m.Edge.X
		md.MulRatios.Edge.X *= 1 - ratio
	}
	if m.Edge.Y != 1 {
		md.MulMaterial.Edge.Y *= md.Material.Edge.Y - m.Edge.Y
		md.MulRatios.Edge.Y *= 1 - ratio
	}
	if m.Edge.Z != 1 {
		md.MulMaterial.Edge.Z *= md.Material.Edge.Z - m.Edge.Z
		md.MulRatios.Edge.Z *= 1 - ratio
	}
	if m.Edge.W != 1 {
		md.MulMaterial.Edge.W *= md.Material.Edge.W - m.Edge.W
		md.MulRatios.Edge.W *= 1 - ratio
	}
	if m.EdgeSize != 1 {
		md.MulMaterial.EdgeSize *= md.Material.EdgeSize - m.EdgeSize
		md.MulRatios.EdgeSize *= ratio
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
