package face

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
)

type Face struct {
	index_model.IndexModel
	Index         int
	VertexIndexes [3]int
}

func NewFace(index, vertexIndex0, vertexIndex1, vertexIndex2 int) *Face {
	return &Face{
		Index:         index,
		VertexIndexes: [3]int{vertexIndex0, vertexIndex1, vertexIndex2},
	}
}

// 面リスト
type Faces struct {
	index_model.IndexModelCorrection[*Face]
}

func NewFaces() *Faces {
	return &Faces{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Face](),
	}
}
