//go:build linux
// +build linux

package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
)

// ボーンリスト
type Bones struct {
	*mcore.IndexNameModels[*Bone]
	Vertices           map[int][]int
	IkTreeIndexes      map[int][]int
	LayerSortedIndexes map[int]string
	LayerSortedNames   map[string]int
}

func NewBones() *Bones {
	return &Bones{
		IndexNameModels:    mcore.NewIndexNameModels[*Bone](),
		Vertices:           make(map[int][]int),
		IkTreeIndexes:      make(map[int][]int),
		LayerSortedIndexes: make(map[int]string),
		LayerSortedNames:   make(map[string]int),
	}
}
