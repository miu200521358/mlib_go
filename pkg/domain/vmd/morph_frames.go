package vmd

import (
	"math"
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type MorphFrames struct {
	names  []string
	values []*MorphNameFrames
}

func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		names:  make([]string, 0),
		values: make([]*MorphNameFrames, 0),
	}
}

func (morphFrames *MorphFrames) Contains(morphName string) bool {
	if slices.Contains(morphFrames.names, morphName) {
		return true
	}

	return false
}

func (morphFrames *MorphFrames) Update(morphNameFrames *MorphNameFrames) {
	index := slices.Index(morphFrames.names, morphNameFrames.Name)
	if index < 0 {
		morphFrames.names = append(morphFrames.names, morphNameFrames.Name)
		morphFrames.values = append(morphFrames.values, morphNameFrames)
	} else {
		morphFrames.names[index] = morphNameFrames.Name
		morphFrames.values[index] = morphNameFrames
	}
}

func (morphFrames *MorphFrames) Delete(morphName string) {
	index := slices.Index(morphFrames.names, morphName)
	if index < 0 {
		return
	}
	morphFrames.names = append(morphFrames.names[:index], morphFrames.names[index+1:]...)
	morphFrames.values = append(morphFrames.values[:index], morphFrames.values[index+1:]...)
}

func (morphFrames *MorphFrames) Get(morphName string) *MorphNameFrames {
	index := slices.Index(morphFrames.names, morphName)

	if index < 0 {
		morphNameFrames := NewMorphNameFrames(morphName)
		morphFrames.names = append(morphFrames.names, morphNameFrames.Name)
		morphFrames.values = append(morphFrames.values, morphNameFrames)
		return morphNameFrames
	}

	return morphFrames.values[index]
}

func (morphFrames *MorphFrames) Names() []string {
	return morphFrames.names
}

func (morphFrames *MorphFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, morphFrames := range morphFrames.values {
		morphFrames.Indexes.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (morphFrames *MorphFrames) IndexesByNames(names []string) []int {
	indexes := make([]int, 0)
	for _, morphName := range morphFrames.names {
		if !slices.Contains(names, morphName) {
			continue
		}
		morphFrames.Get(morphName).Indexes.ForEach(func(index float32) bool {
			indexes = append(indexes, int(index))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (morphFrames *MorphFrames) Length() int {
	count := 0
	for _, morphFrames := range morphFrames.values {
		count += morphFrames.Indexes.Len()
	}
	return count
}

func (morphFrames *MorphFrames) MaxFrame() float32 {
	maxFno := float32(0)
	for _, morphFrames := range morphFrames.values {
		fno := float32(morphFrames.MaxFrame())
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (morphFrames *MorphFrames) MinFrame() float32 {
	minFno := float32(math.MaxFloat32)
	for _, morphFrames := range morphFrames.values {
		fno := float32(morphFrames.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (morphFrames *MorphFrames) Clean() {
	for _, morphName := range morphFrames.names {
		if !morphFrames.Get(morphName).ContainsActive() {
			morphFrames.Delete(morphName)
		}
	}
}

func (morphFrames *MorphFrames) ForEach(fn func(morphName string, morphNameFrames *MorphNameFrames)) {
	for _, morphName := range morphFrames.names {
		fn(morphName, morphFrames.Get(morphName))
	}
}
