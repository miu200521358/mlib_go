package face

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"

)

// 面データ
type Face struct {
	*index_model.IndexModel
	// 頂点INDEXリスト
	VertexIndexes [3]int
}

func NewFace() *Face {
	return &Face{
		IndexModel:    &index_model.IndexModel{Index: -1},
		VertexIndexes: [3]int{0, 0, 0},
	}
}

// 面リスト
type Faces struct {
	*index_model.IndexModelCorrection[*Face]
}

func NewFaces() *Faces {
	return &Faces{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Face](),
	}
}
