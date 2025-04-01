package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 頂点リスト
type Vertices struct {
	*core.IndexModels[*Vertex]
	vertexMap map[int][]*Vertex
}

func NewVertices(capacity int) *Vertices {
	return &Vertices{
		IndexModels: core.NewIndexModels[*Vertex](capacity),
	}
}

// GetMapByBoneIndex はボーンINDEXをキーとして、ウェイト閾値以上の頂点リストを取得します
func (vertices *Vertices) GetMapByBoneIndex(weightThreshold float64) map[int][]*Vertex {
	if vertices.vertexMap != nil {
		return vertices.vertexMap
	}

	vertices.vertexMap = make(map[int][]*Vertex)
	vertices.ForEach(func(index int, vertex *Vertex) bool {
		if vertex.Deform != nil {
			for n, boneIndex := range vertex.Deform.Indexes() {
				if _, ok := vertices.vertexMap[boneIndex]; !ok {
					vertices.vertexMap[boneIndex] = make([]*Vertex, 0)
				}
				if vertex.Deform.Weights()[n] > weightThreshold {
					vertices.vertexMap[boneIndex] = append(vertices.vertexMap[boneIndex], vertex)
				}
			}
		}
		return true
	})

	return vertices.vertexMap
}

func (vertices *Vertices) Append(value *Vertex) {
	vertices.IndexModels.Append(value)
	vertices.vertexMap = nil
}
