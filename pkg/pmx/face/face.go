package face

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
)

type T struct {
	index_model.T
	Index         int
	VertexIndexes [3]int
}

func NewFace(index, vertexIndex0, vertexIndex1, vertexIndex2 int) *T {
	return &T{
		Index:         index,
		VertexIndexes: [3]int{vertexIndex0, vertexIndex1, vertexIndex2},
	}
}

func (m *T) Copy() *T {
	copied := *m
	return &copied
}

// 面リスト
type C struct {
	index_model.C
	data    map[int]*T
	Indexes []int
}

func NewFaces() *C {
	return &C{
		data:    make(map[int]*T),
		Indexes: make([]int, 0),
	}
}
