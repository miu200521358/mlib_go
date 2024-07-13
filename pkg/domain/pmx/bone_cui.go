//go:build linux
// +build linux

package pmx

import "github.com/miu200521358/mlib_go/pkg/domain/core"

// ボーンリスト
type Bones struct {
	*core.IndexNameModels[*Bone]
	IkTreeIndexes     map[int][]int
	LayerSortedBones  map[bool][]*Bone
	LayerSortedNames  map[bool]map[string]int
	DeformBoneIndexes map[int][]int
}

func NewBones(count int) *Bones {
	return &Bones{
		IndexNameModels:  core.NewIndexNameModels[*Bone](count, func() *Bone { return nil }),
		IkTreeIndexes:    make(map[int][]int),
		LayerSortedBones: make(map[bool][]*Bone),
		LayerSortedNames: make(map[bool]map[string]int),
	}
}
