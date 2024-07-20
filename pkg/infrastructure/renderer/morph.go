//go:build windows
// +build windows

package renderer

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

func newVertexMorphDeltasGl(mds *delta.VertexMorphDeltas) ([]int, [][]float32) {
	vertices := make([][]float32, 0)
	indices := make([]int, 0)
	for i, md := range mds.Data {
		vertices = append(vertices, newVertexMorphDeltaGl(md))
		indices = append(indices, i)
	}
	return indices, vertices
}

func newVertexMorphDeltaGl(md *delta.VertexMorphDelta) []float32 {
	var p0, p1, p2 float32
	if md.Position != nil {
		p := mgl.NewGlVec3(md.Position)
		p0, p1, p2 = p[0], p[1], p[2]
	}
	var ap0, ap1, ap2 float32
	if md.AfterPosition != nil {
		ap := mgl.NewGlVec3(md.AfterPosition)
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

func newSelectedVertexMorphDeltasGL(
	mds *SelectedVertexMorphDeltas,
	model *pmx.PmxModel, selectedVertexIndexes, nextSelectedVertexIndexes []int,
) ([]int, [][]float32) {
	indices := make([]int, 0)
	vertices := make([][]float32, 0)
	for i := range len(model.Vertices.Data) {
		// 選択頂点
		var selectedVertexDelta []float32
		// 選択頂点
		if selectedVertexIndexes != nil && nextSelectedVertexIndexes != nil {
			if slices.Contains(selectedVertexIndexes, i) {
				// 選択されている頂点のUVXを＋にして（フラグをたてて）非表示にする
				selectedVertexDelta = []float32{
					0, 0, 0,
					1, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0,
				}
			}
			if slices.Contains(nextSelectedVertexIndexes, i) {
				// 選択されている頂点のUVXを0にして（フラグを落として）表示する
				selectedVertexDelta = []float32{
					0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0,
				}
			}
		} else if selectedVertexIndexes != nil && slices.Contains(selectedVertexIndexes, i) {
			// 選択されている頂点のUVXを0にして（フラグを落として）表示する
			selectedVertexDelta = []float32{
				0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0,
			}
		}

		if selectedVertexDelta != nil {
			vertices = append(vertices, selectedVertexDelta)
			indices = append(indices, i)
		}
	}
	return indices, vertices
}
