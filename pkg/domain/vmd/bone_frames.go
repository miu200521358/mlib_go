package vmd

import (
	"math"
	"slices"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type BoneFrames struct {
	names  []string
	values []*BoneNameFrames
}

func NewBoneFrames() *BoneFrames {
	return &BoneFrames{
		names:  make([]string, 0),
		values: make([]*BoneNameFrames, 0),
	}
}

func (boneFrames *BoneFrames) Contains(boneName string) bool {
	if slices.Contains(boneFrames.names, boneName) {
		return true
	}

	return false
}

func (boneFrames *BoneFrames) Update(boneNameFrames *BoneNameFrames) {
	index := slices.Index(boneFrames.names, boneNameFrames.Name)
	if index < 0 {
		boneFrames.names = append(boneFrames.names, boneNameFrames.Name)
		boneFrames.values = append(boneFrames.values, boneNameFrames)
	} else {
		boneFrames.names[index] = boneNameFrames.Name
		boneFrames.values[index] = boneNameFrames
	}
}

func (boneFrames *BoneFrames) Delete(boneName string) {
	index := slices.Index(boneFrames.names, boneName)
	if index < 0 {
		return
	}
	boneFrames.names = append(boneFrames.names[:index], boneFrames.names[index+1:]...)
	boneFrames.values = append(boneFrames.values[:index], boneFrames.values[index+1:]...)
}

func (boneFrames *BoneFrames) Get(boneName string) *BoneNameFrames {
	index := slices.Index(boneFrames.names, boneName)

	if index < 0 {
		boneNameFrames := NewBoneNameFrames(boneName)
		boneFrames.names = append(boneFrames.names, boneNameFrames.Name)
		boneFrames.values = append(boneFrames.values, boneNameFrames)
		return boneNameFrames
	}

	return boneFrames.values[index]
}

func (boneFrames *BoneFrames) Names() []string {
	return boneFrames.names
}

func (boneFrames *BoneFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, boneFrames := range boneFrames.values {
		boneFrames.Indexes.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (boneFrames *BoneFrames) IndexesByNames(names []string) []int {
	indexes := make([]int, 0)
	for _, boneName := range boneFrames.names {
		if !slices.Contains(names, boneName) {
			continue
		}
		boneFrames.Get(boneName).Indexes.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (boneFrames *BoneFrames) Length() int {
	count := 0
	for _, boneFrames := range boneFrames.values {
		count += boneFrames.Indexes.Len()
	}
	return count
}

func (boneFrames *BoneFrames) MaxFrame() float32 {
	maxFno := float32(0)
	for _, boneFrames := range boneFrames.values {
		fno := float32(boneFrames.MaxFrame())
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (boneFrames *BoneFrames) MinFrame() float32 {
	minFno := float32(math.MaxFloat32)
	for _, boneFrames := range boneFrames.values {
		fno := float32(boneFrames.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (boneFrames *BoneFrames) Clean() {
	for _, boneName := range boneFrames.names {
		if !boneFrames.Get(boneName).ContainsActive() {
			boneFrames.Delete(boneName)
		}
	}
}

func (boneFrames *BoneFrames) Reduce() *BoneFrames {
	reduced := NewBoneFrames()
	var wg sync.WaitGroup
	for _, boneNameFrames := range boneFrames.values {
		wg.Add(1)
		go func(bnf *BoneNameFrames) {
			defer wg.Done()
			reduced.Update(bnf.Reduce())
		}(boneNameFrames)
	}
	wg.Wait()
	return reduced
}

func (boneFrames *BoneFrames) ForEach(fn func(boneName string, boneNameFrames *BoneNameFrames)) {
	for _, boneName := range boneFrames.names {
		fn(boneName, boneFrames.Get(boneName))
	}
}
