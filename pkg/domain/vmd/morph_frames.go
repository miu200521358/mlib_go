package vmd

import (
	"math"
	"slices"
	"sync/atomic"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type MorphFrames struct {
	names  atomic.Value
	values []*MorphNameFrames
}

func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		names:  atomic.Value{},
		values: make([]*MorphNameFrames, 0),
	}
}

func (morphFrames *MorphFrames) getNames() []string {
	names := morphFrames.names.Load()
	if names == nil {
		return make([]string, 0)
	}
	return names.([]string)
}

func (morphFrames *MorphFrames) setNames(names []string) {
	morphFrames.names.Store(names)
}

func (morphFrames *MorphFrames) Contains(morphName string) bool {
	if slices.Contains(morphFrames.getNames(), morphName) {
		return true
	}

	return false
}

func (morphFrames *MorphFrames) Update(morphNameFrames *MorphNameFrames) {
	index := slices.Index(morphFrames.getNames(), morphNameFrames.Name)
	if index < 0 {
		morphFrames.setNames(append(morphFrames.getNames(), morphNameFrames.Name))
		morphFrames.values = append(morphFrames.values, morphNameFrames)
	} else {
		morphFrames.getNames()[index] = morphNameFrames.Name
		morphFrames.values[index] = morphNameFrames
	}
}

func (morphFrames *MorphFrames) Delete(morphName string) {
	index := slices.Index(morphFrames.getNames(), morphName)
	if index < 0 {
		return
	}
	morphFrames.setNames(append(morphFrames.getNames()[:index], morphFrames.getNames()[index+1:]...))
	morphFrames.values = append(morphFrames.values[:index], morphFrames.values[index+1:]...)
}

func (morphFrames *MorphFrames) Get(morphName string) *MorphNameFrames {
	index := slices.Index(morphFrames.getNames(), morphName)

	if index < 0 || index >= len(morphFrames.values) {
		morphNameFrames := NewMorphNameFrames(morphName)
		morphFrames.setNames(append(morphFrames.getNames(), morphNameFrames.Name))
		morphFrames.values = append(morphFrames.values, morphNameFrames)
		return morphNameFrames
	}

	return morphFrames.values[index]
}

func (morphFrames *MorphFrames) Names() []string {
	return morphFrames.getNames()
}

func (morphFrames *MorphFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, morphFrames := range morphFrames.values {
		morphFrames.values.ForEach(func(index float32) bool {
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
	for _, morphName := range morphFrames.getNames() {
		if !slices.Contains(names, morphName) {
			continue
		}
		morphFrames.Get(morphName).values.ForEach(func(index float32) bool {
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
		count += morphFrames.values.Len()
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
	for _, morphName := range morphFrames.getNames() {
		if !morphFrames.Get(morphName).ContainsActive() {
			morphFrames.Delete(morphName)
		}
	}
}

func (morphFrames *MorphFrames) ForEach(fn func(morphName string, morphNameFrames *MorphNameFrames)) {
	for _, morphName := range morphFrames.getNames() {
		fn(morphName, morphFrames.Get(morphName))
	}
}

func (morphFrames *MorphFrames) Copy() (*MorphFrames, error) {
	copied := new(MorphFrames)
	err := deepcopy.Copy(copied, morphFrames)
	return copied, err
}
