package vmd

import (
	"math"
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type VertexDelta struct {
	Position *mmath.MVec3
}

func NewVertexDelta(pos *mmath.MVec3) *VertexDelta {
	return &VertexDelta{
		Position: pos,
	}
}

type VertexDeltas struct {
	Data     []*VertexDelta
	Vertices *pmx.Vertices
}

func NewVertexDeltas(vertices *pmx.Vertices) *VertexDeltas {
	return &VertexDeltas{
		Data:     make([]*VertexDelta, vertices.Len()),
		Vertices: vertices,
	}
}

func (vds *VertexDeltas) Get(vertexIndex int) *VertexDelta {
	if vertexIndex < 0 || vertexIndex >= len(vds.Data) {
		return nil
	}

	return vds.Data[vertexIndex]
}

func (vds *VertexDeltas) GetNearestVertexIndexes(worldPos *mmath.MVec3, visibleMaterialIndexes []int) []int {
	vertexIndexes := make([]int, 0)
	distances := make([]float64, len(vds.Data))
	for i := range len(vds.Data) {
		vd := vds.Get(i)
		vertex := vds.Vertices.Get(i)
		if visibleMaterialIndexes == nil {
			distances[i] = worldPos.Distance(vd.Position)
		} else {
			for _, materialIndex := range visibleMaterialIndexes {
				if slices.Contains(vertex.MaterialIndexes, materialIndex) {
					distances[i] = worldPos.Distance(vd.Position)
					break
				} else {
					// 非表示材質は最長距離をとりあえず入れておく
					distances[i] = math.MaxFloat64
				}
			}
		}
	}
	if len(distances) == 0 {
		return vertexIndexes
	}
	sortedDistances := mmath.Float64Slice(distances)
	sortedIndexes := mmath.ArgSort(sortedDistances)
	nearestVertex := vds.Get(sortedIndexes[0])
	for i := range sortedIndexes {
		vd := vds.Get(sortedIndexes[i])
		if len(vertexIndexes) > 0 {
			if !vd.Position.NearEquals(nearestVertex.Position, 0.01) {
				break
			}
		}
		vertexIndexes = append(vertexIndexes, sortedIndexes[i])
	}
	return vertexIndexes
}
