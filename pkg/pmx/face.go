package pmx

import (
	"github.com/jinzhu/copier"
	"github.com/miu200521358/mlib_go/pkg/mcore"
)

// 面データ
type Face struct {
	*mcore.IndexModel
	VertexIndexes [3]int // 頂点INDEXリスト
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

func (f *Face) Copy() mcore.IIndexModel {
	copied := NewFace()
	copier.CopyWithOption(copied, f, copier.Option{DeepCopy: true})
	return copied
}

// 面リスト
type Faces struct {
	*mcore.IndexModels[*Face]
}

func NewFaces() *Faces {
	return &Faces{
		IndexModels: mcore.NewIndexModels[*Face](func() *Face { return nil }),
	}
}
