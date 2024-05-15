//go:build windows
// +build windows

package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mview"
)

// ボーンリスト
type Bones struct {
	*mcore.IndexNameModels[*Bone]
	Vertices         map[int][]int
	IkTreeIndexes    map[int][]int
	LayerSortedBones map[int]*Bone
	LayerSortedNames map[string]int
	positionVao      *mview.VAO
	positionIbo      *mview.IBO
	positionIboCount int32
	normalVao        *mview.VAO
	normalIbo        *mview.IBO
	normalIboCount   int32
}

func NewBones() *Bones {
	return &Bones{
		IndexNameModels:  mcore.NewIndexNameModels[*Bone](),
		Vertices:         make(map[int][]int),
		IkTreeIndexes:    make(map[int][]int),
		LayerSortedBones: make(map[int]*Bone),
		LayerSortedNames: make(map[string]int),
		positionVao:      nil,
		positionIbo:      nil,
		positionIboCount: 0,
		normalVao:        nil,
		normalIbo:        nil,
		normalIboCount:   0,
	}
}
