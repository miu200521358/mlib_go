//go:build windows
// +build windows

package render

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func newVertexMorphDeltasGl(mds *delta.VertexMorphDeltas) ([]int, [][]float32) {
	vertices := make([][]float32, 0)
	indexes := make([]int, 0)
	mds.ForEach(func(i int, v *delta.VertexMorphDelta) {
		vertices = append(vertices, newVertexMorphDeltaGl(v))
		indexes = append(indexes, i)
	})
	return indexes, vertices
}

func newVertexMorphDeltaGl(md *delta.VertexMorphDelta) []float32 {
	var p0, p1, p2 float32
	if md.Position != nil {
		p := mmath.NewGlVec3(md.Position)
		p0, p1, p2 = p[0], p[1], p[2]
	}
	var ap0, ap1, ap2 float32
	if md.AfterPosition != nil {
		ap := mmath.NewGlVec3(md.AfterPosition)
		ap0, ap1, ap2 = ap[0], ap[1], ap[2]
	}
	// UVは符号関係ないのでそのまま取得する
	var u0x, u0y, u1x, u1y float32
	if md.Uv != nil {
		u0x = float32(md.Uv.X)
		u0y = float32(md.Uv.Y)
	}
	if md.Uv1 != nil {
		u1x = float32(md.Uv1.X)
		u1y = float32(md.Uv1.Y)
	}
	return []float32{
		p0, p1, p2,
		u0x, u0y, 0, 0,
		u1x, u1y, 0, 0,
		ap0, ap1, ap2,
	}
}
