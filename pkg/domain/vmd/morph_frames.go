package vmd

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type MorphFrames struct {
	data map[string]*MorphNameFrames
}

func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		data: make(map[string]*MorphNameFrames, 0),
	}
}

func (morphFrames *MorphFrames) Delete(morphName string) {
	delete(morphFrames.data, morphName)
}

func (morphFrames *MorphFrames) Contains(morphName string) bool {
	if _, ok := morphFrames.data[morphName]; ok {
		if morphFrames.data[morphName] != nil && morphFrames.data[morphName].Length() > 0 {
			return true
		}
	}

	return false
}

func (morphFrames *MorphFrames) Update(morphNameFrames *MorphNameFrames) {
	morphFrames.data[morphNameFrames.Name] = morphNameFrames
}

func (morphFrames *MorphFrames) Names() []string {
	names := make([]string, 0, len(morphFrames.data))
	for name := range morphFrames.data {
		names = append(names, name)
	}
	return names
}

func (morphFrames *MorphFrames) Get(morphName string) *MorphNameFrames {
	if !morphFrames.Contains(morphName) {
		morphFrames.Update(NewMorphNameFrames(morphName))
	}
	return morphFrames.data[morphName]
}

func (morphFrames *MorphFrames) MaxFrame() float32 {
	maxFno := float32(0)
	for _, mnfs := range morphFrames.data {
		fno := float32(mnfs.MaxFrame())
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (morphFrames *MorphFrames) MinFrame() float32 {
	minFno := float32(math.MaxFloat32)
	for _, mnfs := range morphFrames.data {
		fno := float32(mnfs.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (morphFrames *MorphFrames) Length() int {
	count := 0
	for _, fs := range morphFrames.data {
		count += fs.RegisteredIndexes.Length()
	}
	return count
}

func (morphFrames *MorphFrames) Clean() {
	for morphName := range morphFrames.data {
		if !morphFrames.Get(morphName).ContainsActive() {
			morphFrames.Delete(morphName)
		}
	}
}

func (morphFrames *MorphFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, morphFrames := range morphFrames.data {
		morphFrames.Indexes.ForEach(func(index float32) {
			indexes = append(indexes, int(index))
		})
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (morphFrames *MorphFrames) RegisteredIndexes() []int {
	indexes := make([]int, 0)
	for _, morphFrames := range morphFrames.data {
		morphFrames.RegisteredIndexes.ForEach(func(index float32) {
			indexes = append(indexes, int(index))
		})
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (morphFrames *MorphFrames) ForEach(f func(morphName string, morphNameFrames *MorphNameFrames)) {
	for morphName, morphNameFrames := range morphFrames.data {
		f(morphName, morphNameFrames)
	}
}
