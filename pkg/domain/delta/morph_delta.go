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

func (vertexMorphDeltas *VertexMorphDeltas) Get(index int) *VertexMorphDelta {
	return vertexMorphDeltas.Data[index]
}

func (vertexMorphDeltas *VertexMorphDeltas) Update(v *VertexMorphDelta) {
	vertexMorphDeltas.Data[v.Index] = v
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
	BoneIndex               int
	FramePosition           *mmath.MVec3       // キーフレ位置の変動量
	FrameCancelablePosition *mmath.MVec3       // キャンセル位置の変動量
	FrameRotation           *mmath.MQuaternion // キーフレ回転の変動量
	FrameCancelableRotation *mmath.MQuaternion // キャンセル回転の変動量
	FrameScale              *mmath.MVec3       // キーフレスケールの変動量
	FrameCancelableScale    *mmath.MVec3       // キャンセルスケールの変動量
	FrameLocalMat           *mmath.MMat4       // キーフレのローカル変動行列
}

func NewBoneMorphDelta(boneIndex int) *BoneMorphDelta {
	return &BoneMorphDelta{
		BoneIndex: boneIndex,
	}
}

func (boneMorphDelta *BoneMorphDelta) Get(boneIndex int) *BoneMorphDelta {
	return boneMorphDelta
}

func (boneMorphDelta *BoneMorphDelta) FilledMorphPosition() *mmath.MVec3 {
	if boneMorphDelta.FramePosition == nil {
		boneMorphDelta.FramePosition = mmath.NewMVec3()
	}
	return boneMorphDelta.FramePosition
}

func (boneMorphDelta *BoneMorphDelta) FilledMorphCancelablePosition() *mmath.MVec3 {
	if boneMorphDelta.FrameCancelablePosition == nil {
		boneMorphDelta.FrameCancelablePosition = mmath.NewMVec3()
	}
	return boneMorphDelta.FrameCancelablePosition
}

func (boneMorph *BoneMorphDelta) FilledMorphRotation() *mmath.MQuaternion {
	if boneMorph.FrameRotation == nil {
		boneMorph.FrameRotation = mmath.NewMQuaternion()
	}
	return boneMorph.FrameRotation
}

func (boneMorph *BoneMorphDelta) FilledMorphCancelableRotation() *mmath.MQuaternion {
	if boneMorph.FrameCancelableRotation == nil {
		boneMorph.FrameCancelableRotation = mmath.NewMQuaternion()
	}
	return boneMorph.FrameCancelableRotation
}

func (boneMorphDelta *BoneMorphDelta) FilledMorphScale() *mmath.MVec3 {
	if boneMorphDelta.FrameScale == nil {
		boneMorphDelta.FrameScale = mmath.NewMVec3()
	}
	return boneMorphDelta.FrameScale
}

func (boneMorphDelta *BoneMorphDelta) FilledMorphCancelableScale() *mmath.MVec3 {
	if boneMorphDelta.FrameCancelableScale == nil {
		boneMorphDelta.FrameCancelableScale = mmath.NewMVec3()
	}
	return boneMorphDelta.FrameCancelableScale
}

func (boneMorphDelta *BoneMorphDelta) FilledMorphLocalMat() *mmath.MMat4 {
	if boneMorphDelta.FrameLocalMat == nil {
		boneMorphDelta.FrameLocalMat = mmath.NewMMat4()
	}
	return boneMorphDelta.FrameLocalMat
}

func (boneMorphDelta *BoneMorphDelta) Copy() *BoneMorphDelta {
	return &BoneMorphDelta{
		FramePosition:           boneMorphDelta.FilledMorphPosition().Copy(),
		FrameCancelablePosition: boneMorphDelta.FilledMorphCancelablePosition().Copy(),
		FrameRotation:           boneMorphDelta.FilledMorphRotation().Copy(),
		FrameCancelableRotation: boneMorphDelta.FilledMorphCancelableRotation().Copy(),
		FrameScale:              boneMorphDelta.FilledMorphScale().Copy(),
		FrameCancelableScale:    boneMorphDelta.FilledMorphCancelableScale().Copy(),
		FrameLocalMat:           boneMorphDelta.FilledMorphLocalMat().Copy(),
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
	}

	return mm
}

func (materialMorphDelta *MaterialMorphDelta) Add(m *pmx.MaterialMorphOffset, ratio float64) {
	materialMorphDelta.AddMaterial.Diffuse.X += m.Diffuse.X * ratio
	materialMorphDelta.AddMaterial.Diffuse.Y += m.Diffuse.Y * ratio
	materialMorphDelta.AddMaterial.Diffuse.Z += m.Diffuse.Z * ratio
	materialMorphDelta.AddMaterial.Diffuse.W += m.Diffuse.W * ratio
	materialMorphDelta.AddMaterial.Specular.X += m.Specular.X * ratio
	materialMorphDelta.AddMaterial.Specular.Y += m.Specular.Y * ratio
	materialMorphDelta.AddMaterial.Specular.Z += m.Specular.Z * ratio
	materialMorphDelta.AddMaterial.Specular.W += m.Specular.W * ratio
	materialMorphDelta.AddMaterial.Ambient.X += m.Ambient.X * ratio
	materialMorphDelta.AddMaterial.Ambient.Y += m.Ambient.Y * ratio
	materialMorphDelta.AddMaterial.Ambient.Z += m.Ambient.Z * ratio
	materialMorphDelta.AddMaterial.Edge.X += m.Edge.X * ratio
	materialMorphDelta.AddMaterial.Edge.Y += m.Edge.Y * ratio
	materialMorphDelta.AddMaterial.Edge.Z += m.Edge.Z * ratio
	materialMorphDelta.AddMaterial.Edge.W += m.Edge.W * ratio
	materialMorphDelta.AddMaterial.EdgeSize += m.EdgeSize * ratio
}

func (materialMorphDelta *MaterialMorphDelta) Mul(m *pmx.MaterialMorphOffset, ratio float64) {
	materialMorphDelta.MulMaterial.Diffuse.X *= mmath.LerpFloat(1, m.Diffuse.X, ratio)
	materialMorphDelta.MulMaterial.Diffuse.Y *= mmath.LerpFloat(1, m.Diffuse.Y, ratio)
	materialMorphDelta.MulMaterial.Diffuse.Z *= mmath.LerpFloat(1, m.Diffuse.Z, ratio)
	materialMorphDelta.MulMaterial.Diffuse.W *= mmath.LerpFloat(1, m.Diffuse.W, ratio)
	materialMorphDelta.MulMaterial.Specular.X *= mmath.LerpFloat(1, m.Specular.X, ratio)
	materialMorphDelta.MulMaterial.Specular.Y *= mmath.LerpFloat(1, m.Specular.Y, ratio)
	materialMorphDelta.MulMaterial.Specular.Z *= mmath.LerpFloat(1, m.Specular.Z, ratio)
	materialMorphDelta.MulMaterial.Specular.W *= mmath.LerpFloat(1, m.Specular.W, ratio)
	materialMorphDelta.MulMaterial.Ambient.X *= mmath.LerpFloat(1, m.Ambient.X, ratio)
	materialMorphDelta.MulMaterial.Ambient.Y *= mmath.LerpFloat(1, m.Ambient.Y, ratio)
	materialMorphDelta.MulMaterial.Ambient.Z *= mmath.LerpFloat(1, m.Ambient.Z, ratio)
	materialMorphDelta.MulMaterial.Edge.X *= mmath.LerpFloat(1, m.Edge.X, ratio)
	materialMorphDelta.MulMaterial.Edge.Y *= mmath.LerpFloat(1, m.Edge.Y, ratio)
	materialMorphDelta.MulMaterial.Edge.Z *= mmath.LerpFloat(1, m.Edge.Z, ratio)
	materialMorphDelta.MulMaterial.Edge.W *= mmath.LerpFloat(1, m.Edge.W, ratio)
	materialMorphDelta.MulMaterial.EdgeSize *= mmath.LerpFloat(1, m.EdgeSize, ratio)
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

func (materialMorphDeltas *MaterialMorphDeltas) Get(index int) *MaterialMorphDelta {
	if index < 0 || index >= len(materialMorphDeltas.Data) {
		return nil
	}

	return materialMorphDeltas.Data[index]
}

func (materialMorphDeltas *MaterialMorphDeltas) Update(m *MaterialMorphDelta) {
	materialMorphDeltas.Data[m.Index()] = m
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
