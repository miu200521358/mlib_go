package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 面データ
type Face struct {
	index         int    // 面INDEX
	VertexIndexes [3]int // 頂点INDEXリスト
}

type FaceGL struct {
	VertexIndexes [3]uint32
}

func NewFace() *Face {
	return &Face{
		index:         -1,
		VertexIndexes: [3]int{0, 0, 0},
	}
}

func (face *Face) Index() int {
	return face.index
}

func (face *Face) SetIndex(index int) {
	face.index = index
}

func (face *Face) IsValid() bool {
	return face != nil && face.Index() >= 0
}

func (face *Face) Copy() core.IIndexModel {
	return &Face{
		index:         face.index,
		VertexIndexes: [3]int{face.VertexIndexes[0], face.VertexIndexes[1], face.VertexIndexes[2]},
	}
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

func (faces *Faces) Copy() *Faces {
	copied := NewFaces(len(faces.Data))
	for i, face := range faces.Data {
		copied.SetItem(i, face.Copy().(*Face))
	}
	return copied
}
