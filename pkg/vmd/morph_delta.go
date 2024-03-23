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
	uv := md.Uv.GL()
	uv1 := md.Uv1.GL()
	ap := md.AfterPosition.GL()
	return []float32{
		p[0], p[1], p[2],
		uv[0], uv[1],
		uv1[0], uv1[1],
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
	}

	return &VertexMorphDeltas{
		Data: deltas,
	}
}

type MorphDeltas struct {
	Vertices *VertexMorphDeltas
}

func NewMorphDeltas(vertexCount int) *MorphDeltas {
	return &MorphDeltas{
		Vertices: NewVertexMorphDeltas(vertexCount),
	}
}

func (mds *MorphDeltas) GetFrameNos() []float32 {
	frames := make([]float32, 0)
	for _, v := range mds.Vertices.Data {
		if !slices.Contains(frames, v.Frame) {
			frames = append(frames, v.Frame)
		}
	}
	return frames
}
