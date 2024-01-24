package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
)

// 面データ
type Face struct {
	*mcore.IndexModel
	// 頂点INDEXリスト
	VertexIndexes [3]int
}

type FaceGL struct {
	VertexIndexes [3]uint32
}

func NewFace() *Face {
	return &Face{
		IndexModel:    &mcore.IndexModel{Index: -1},
		VertexIndexes: [3]int{0, 0, 0},
	}
}

func (f *Face) GL() *FaceGL {
	return &FaceGL{
		VertexIndexes: [3]uint32{uint32(f.VertexIndexes[2]), uint32(f.VertexIndexes[1]), uint32(f.VertexIndexes[0])},
	}
}

// 面リスト
type Faces struct {
	*mcore.IndexModelCorrection[*Face]
}

func NewFaces() *Faces {
	return &Faces{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Face](),
	}
}
