package vmd

import (
	"math"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type MorphFrames struct {
	data map[string]*MorphNameFrames
	lock sync.RWMutex // マップアクセス制御用
}

func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		data: make(map[string]*MorphNameFrames, 0),
	}
}

func (morphFrames *MorphFrames) Delete(morphName string) {
	morphFrames.lock.Lock()
	defer morphFrames.lock.Unlock()

	delete(morphFrames.data, morphName)
}

func (morphFrames *MorphFrames) Contains(morphName string) bool {
	morphFrames.lock.Lock()
	defer morphFrames.lock.Unlock()

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

// ContainsActive 有効なキーフレが存在するか
func (morphFrames *MorphFrames) ContainsActive(morphName string) bool {
	morphFrames.lock.RLock()
	defer morphFrames.lock.RUnlock()

	if _, ok := morphFrames.data[morphName]; !ok {
		return false
	}

	if morphFrames.data[morphName].Length() == 0 {
		return false
	}

	for mf := range morphFrames.data[morphName].Iterator() {
		if mf == nil {
			return false
		}

		if !mmath.NearEquals(mf.Ratio, 0.0, 1e-2) {
			return true
		}

		nextMf := morphFrames.data[morphName].Get(morphFrames.data[morphName].NextFrame(mf.Index()))

		if nextMf == nil {
			return false
		}

		if !mmath.NearEquals(nextMf.Ratio, 0.0, 1e-2) {
			return true
		}
	}

	return false
}

func (morphFrames *MorphFrames) Clean() {
	for morphName := range morphFrames.data {
		if !morphFrames.ContainsActive(morphName) {
			morphFrames.Delete(morphName)
		}
	}
}

func (morphFrames *MorphFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, morphFrames := range morphFrames.data {
		for bf := range morphFrames.Iterator() {
			indexes = append(indexes, int(bf.Index()))
		}
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (morphFrames *MorphFrames) RegisteredIndexes() []int {
	indexes := make([]int, 0)
	for _, morphFrames := range morphFrames.data {
		for index := range morphFrames.RegisteredIndexes.Iterator() {
			indexes = append(indexes, int(index))
		}
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}
