package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
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
	Data map[int]*VertexDelta
}

func NewVertexDeltas() *VertexDeltas {
	return &VertexDeltas{
		Data: make(map[int]*VertexDelta),
	}
}

func (vds *VertexDeltas) Get(boneIndex int) *VertexDelta {
	if _, ok := vds.Data[boneIndex]; ok {
		return vds.Data[boneIndex]
	}
	return nil
}

func (vds *VertexDeltas) GetNearestVertexIndexes(worldPos *mmath.MVec3) []int {
	vertexIndexes := make([]int, 0)
	distances := make([]float64, len(vds.Data))
	for i := range len(vds.Data) {
		vd := vds.Get(i)
		distances[i] = worldPos.Distance(vd.Position)
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
