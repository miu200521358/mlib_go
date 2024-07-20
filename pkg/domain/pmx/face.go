package pmx

import (
	"github.com/jinzhu/copier"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 面データ
type Face struct {
	*core.IndexModel
	VertexIndexes [3]int // 頂点INDEXリスト
}

type FaceGL struct {
	VertexIndexes [3]uint32
}

func NewFace() *Face {
	return &Face{
		IndexModel:    core.NewIndexModel(-1),
		VertexIndexes: [3]int{0, 0, 0},
	}
}

func (f *Face) Copy() core.IIndexModel {
	copied := NewFace()
	copier.CopyWithOption(copied, f, copier.Option{DeepCopy: true})
	return copied
}

// 面リスト
type Faces struct {
	*core.IndexModels[*Face]
}

func NewFaces(count int) *Faces {
	return &Faces{
		IndexModels: core.NewIndexModels[*Face](count, func() *Face { return nil }),
	}
}
