package vmd

import (
	"math"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type BoneFrames struct {
	data map[string]*BoneNameFrames
	lock sync.RWMutex // マップアクセス制御用
}

func NewBoneFrames() *BoneFrames {
	return &BoneFrames{
		data: make(map[string]*BoneNameFrames, 0),
		lock: sync.RWMutex{},
	}
}

func (boneFrames *BoneFrames) Contains(boneName string) bool {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	if _, ok := boneFrames.data[boneName]; ok {
		if boneFrames.data[boneName] != nil && boneFrames.data[boneName].Length() > 0 {
			return true
		}
	}

	return false
}

func (boneFrames *BoneFrames) Append(boneNameFrames *BoneNameFrames) {
	boneFrames.lock.Lock()
	defer boneFrames.lock.Unlock()

	boneFrames.data[boneNameFrames.Name] = boneNameFrames
}

func (boneFrames *BoneFrames) Delete(boneName string) {
	boneFrames.lock.Lock()
	defer boneFrames.lock.Unlock()

	delete(boneFrames.data, boneName)
}

func (boneFrames *BoneFrames) Get(boneName string) *BoneNameFrames {
	if !boneFrames.Contains(boneName) {
		boneFrames.Append(NewBoneNameFrames(boneName))
	}

	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	return boneFrames.data[boneName]
}

func (boneFrames *BoneFrames) Names() []string {
	names := make([]string, 0, len(boneFrames.data))
	for name := range boneFrames.data {
		names = append(names, name)
	}
	return names
}

func (boneFrames *BoneFrames) Indexes() []int {
	indexes := make([]int, 0)
	for _, boneFrames := range boneFrames.data {
		for bf := range boneFrames.Iterator() {
			indexes = append(indexes, int(bf.Index()))
		}
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (boneFrames *BoneFrames) RegisteredIndexes() []int {
	indexes := make([]int, 0)
	for _, boneFrames := range boneFrames.data {
		for index := range boneFrames.RegisteredIndexes.Iterator() {
			indexes = append(indexes, int(index))
		}
	}
	mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

func (boneFrames *BoneFrames) Length() int {
	count := 0
	for _, boneFrames := range boneFrames.data {
		count += boneFrames.RegisteredIndexes.Len()
	}
	return count
}

func (boneFrames *BoneFrames) MaxFrame() float32 {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	maxFno := float32(0)
	for _, boneFrames := range boneFrames.data {
		fno := float32(boneFrames.MaxFrame())
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (boneFrames *BoneFrames) MinFrame() float32 {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	minFno := float32(math.MaxFloat32)
	for _, boneFrames := range boneFrames.data {
		fno := float32(boneFrames.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (boneFrames *BoneFrames) Clean() {
	for boneName, boneNameFrames := range boneFrames.data {
		if !boneNameFrames.ContainsActive() {
			boneFrames.Delete(boneName)
		}
	}
}

func (boneFrames *BoneFrames) Reduce() *BoneFrames {
	reduced := NewBoneFrames()
	var wg sync.WaitGroup
	for _, boneNameFrames := range boneFrames.data {
		wg.Add(1)
		go func(bnf *BoneNameFrames) {
			defer wg.Done()
			reduced.Append(bnf.Reduce())
		}(boneNameFrames)
	}
	wg.Wait()
	return reduced
}

func (boneFrames *BoneFrames) Iterator() <-chan *BoneNameFrames {
	ch := make(chan *BoneNameFrames)
	go func() {
		boneFrames.lock.RLock()
		defer boneFrames.lock.RUnlock()
		for _, boneNameFrames := range boneFrames.data {
			ch <- boneNameFrames
		}
		close(ch)
	}()
	return ch
}
