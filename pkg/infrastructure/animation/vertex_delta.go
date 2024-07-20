//go:build windows
// +build windows

package animation

import (
	"math"
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type SelectedVertexMorphDeltas struct {
	*delta.VertexMorphDeltas
}

func NewSelectedVertexMorphDeltas() *SelectedVertexMorphDeltas {
	return &SelectedVertexMorphDeltas{
		VertexMorphDeltas: delta.NewVertexMorphDeltas(),
	}
}

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

func (vds *VertexDeltas) FindNearestVertexIndexes(frontPos *mmath.MVec3, visibleMaterialIndexes []int) [][]int {
	nearestVertexIndexes := make([][]int, 0)
	distances := make([]float64, len(vds.Data))
	for i := range len(vds.Data) {
		vd := vds.Get(i)
		vertex := vds.Vertices.Get(i)
		if visibleMaterialIndexes == nil {
			distances[i] = frontPos.Distance(vd.Position)
		} else {
			for _, materialIndex := range visibleMaterialIndexes {
				if slices.Contains(vertex.MaterialIndexes, materialIndex) {
					distances[i] = frontPos.Distance(vd.Position)
					break
				} else {
					// 非表示材質は最長距離をとりあえず入れておく
					distances[i] = math.MaxFloat64
				}
			}
		}
	}
	if len(distances) == 0 {
		return nearestVertexIndexes
	}
	sortedDistances := mmath.Float64Slice(distances)
	sortedIndexes := mmath.ArgSort(sortedDistances)
	nearestVertex := vds.Get(sortedIndexes[0])

	nearestVertexIndexes = append(nearestVertexIndexes, make([]int, 0))
	nearestVertexIndexes[0] = append(nearestVertexIndexes[0], sortedIndexes[0])

	for i := range sortedIndexes {
		vd := vds.Get(sortedIndexes[i])
		if !vd.Position.NearEquals(nearestVertex.Position, 0.01) {
			break
		}
		nearestVertexIndexes[0] = append(nearestVertexIndexes[0], sortedIndexes[i])
	}
	return nearestVertexIndexes
}

func (vds *VertexDeltas) FindVerticesInBox(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos,
	prevXnowYBackPos, nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3,
	visibleMaterialIndexes []int) [][]int {

	boxVertexIndexes := make([][]int, 0)
	vertexIndexMap := make(map[mmath.MVec3][]int, 0)

	// 境界ボックスを計算
	minPos, maxPos := mmath.CalculateBoundingBox(
		prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
		nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos,
	)

	for i := range len(vds.Data) {
		vd := vds.Get(i)
		vertex := vds.Vertices.Get(i)
		if visibleMaterialIndexes == nil {
			if vd.Position.IsPointInsideBox(minPos, maxPos) {
				// 小数点第二位で四捨五入
				posKey := mmath.MVec3{
					X: math.Round(vd.Position.X*100) / 100,
					Y: math.Round(vd.Position.Y*100) / 100,
					Z: math.Round(vd.Position.Z*100) / 100}
				if _, ok := vertexIndexMap[posKey]; !ok {
					vertexIndexMap[posKey] = make([]int, 0)
				}
				vertexIndexMap[posKey] = append(vertexIndexMap[posKey], vertex.Index())
			}
		} else {
			for _, materialIndex := range visibleMaterialIndexes {
				if slices.Contains(vertex.MaterialIndexes, materialIndex) {
					if vd.Position.IsPointInsideBox(minPos, maxPos) {
						// 小数点第二位で四捨五入
						posKey := mmath.MVec3{
							X: math.Round(vd.Position.X*100) / 100,
							Y: math.Round(vd.Position.Y*100) / 100,
							Z: math.Round(vd.Position.Z*100) / 100}
						if _, ok := vertexIndexMap[posKey]; !ok {
							vertexIndexMap[posKey] = make([]int, 0)
						}
						vertexIndexMap[posKey] = append(vertexIndexMap[posKey], vertex.Index())
					}
					break
				}
			}
		}
	}

	for _, vertexIndexes := range vertexIndexMap {
		if len(vertexIndexes) > 0 {
			boxVertexIndexes = append(boxVertexIndexes, vertexIndexes)
		}
	}

	return boxVertexIndexes
}
