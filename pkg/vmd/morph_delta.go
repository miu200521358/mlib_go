package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type VertexMorphDelta struct {
	Frame         float32
	Index         int
	Position      *mmath.MVec3
	Uv            *mmath.MVec2
	Uv1           *mmath.MVec2
	AfterPosition *mmath.MVec3
}

func NewVertexMorphDelta() *VertexMorphDelta {
	return &VertexMorphDelta{
		Frame:         0.0,
		Index:         0,
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
		deltas[i].Index = i
	}

	return &VertexMorphDeltas{
		Data: deltas,
	}
}

type BoneMorphDelta struct {
	BoneFrame
	Frame float32
	Index int
}

func NewBoneMorphDelta() *BoneMorphDelta {
	return &BoneMorphDelta{
		BoneFrame: *NewBoneFrame(0.0),
		Frame:     0.0,
		Index:     0,
	}
}

type BoneMorphDeltas struct {
	Data []*BoneMorphDelta
}

func NewBoneMorphDeltas(boneCount int) *BoneMorphDeltas {
	deltas := make([]*BoneMorphDelta, boneCount)
	for i := 0; i < boneCount; i++ {
		deltas[i] = NewBoneMorphDelta()
		deltas[i].Index = i
	}

	return &BoneMorphDeltas{
		Data: deltas,
	}
}

type MorphDeltas struct {
	Vertices *VertexMorphDeltas
	Bones    *BoneMorphDeltas
}

func NewMorphDeltas(vertexCount int, boneCount int) *MorphDeltas {
	return &MorphDeltas{
		Vertices: NewVertexMorphDeltas(vertexCount),
		Bones:    NewBoneMorphDeltas(boneCount),
	}
}

func (mds *MorphDeltas) GetFrameNos() []float32 {
	frames := make([]float32, 0)
	for _, v := range mds.Vertices.Data {
		if !slices.Contains(frames, v.Frame) {
			frames = append(frames, v.Frame)
		}
	}
	for _, b := range mds.Bones.Data {
		if !slices.Contains(frames, b.Frame) {
			frames = append(frames, b.Frame)
		}
	}
	return frames
}
