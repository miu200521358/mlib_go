//go:build windows
// +build windows

package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mgl"
)

// ボーンリスト
type Bones struct {
	*mcore.IIndexNameModels[*Bone]
	Vertices           map[int][]int
	IkTreeIndexes      map[int][]int
	LayerSortedIndexes map[int]string
	LayerSortedNames   map[string]int
	positionVao        *mgl.VAO
	positionIbo        *mgl.IBO
	positionIboCount   int32
	normalVao          *mgl.VAO
	normalIbo          *mgl.IBO
	normalIboCount     int32
}

func NewBones() *Bones {
	return &Bones{
		IIndexNameModels:   mcore.NewIndexNameModelCorrection[*Bone](),
		Vertices:           make(map[int][]int),
		IkTreeIndexes:      make(map[int][]int),
		LayerSortedIndexes: make(map[int]string),
		LayerSortedNames:   make(map[string]int),
		positionVao:        nil,
		positionIbo:        nil,
		positionIboCount:   0,
		normalVao:          nil,
		normalIbo:          nil,
		normalIboCount:     0,
	}
}
