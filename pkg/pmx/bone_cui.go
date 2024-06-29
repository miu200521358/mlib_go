//go:build linux
// +build linux

package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

// ボーンリスト
type Bones struct {
	*mcore.IndexNameModels[*Bone]
	Positions        []*mmath.MVec3
	IkTreeIndexes    map[int][]int
	LayerSortedBones map[bool]map[int]*Bone
	LayerSortedNames map[bool]map[string]int
}

func NewBones() *Bones {
	return &Bones{
		IndexNameModels:  mcore.NewIndexNameModels[*Bone](func() *Bone { return nil }),
		IkTreeIndexes:    make(map[int][]int),
		LayerSortedBones: make(map[bool]map[int]*Bone),
		LayerSortedNames: make(map[bool]map[string]int),
	}
}
