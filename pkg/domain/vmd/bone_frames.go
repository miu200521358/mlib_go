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

// ContainsActive 有効なキーフレが存在するか
func (boneFrames *BoneFrames) ContainsActive(boneName string) bool {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	if _, ok := boneFrames.data[boneName]; !ok {
		return false
	}

	if boneFrames.data[boneName].Length() == 0 {
		return false
	}

	for bf := range boneFrames.data[boneName].Iterator() {
		if bf == nil {
			return false
		}

		if (bf.Position != nil && !bf.Position.NearEquals(mmath.MVec3Zero, 1e-2)) ||
			(bf.Rotation != nil && !bf.Rotation.NearEquals(mmath.MQuaternionIdent, 1e-2)) {
			return true
		}

		nextBf := boneFrames.data[boneName].Get(boneFrames.data[boneName].NextFrame(bf.Index()))

		if nextBf == nil {
			return false
		}

		if bf.Position != nil && nextBf.Position != nil && !bf.Position.NearEquals(nextBf.Position, 1e-2) {
			return true
		}

		if bf.Rotation != nil && nextBf.Rotation != nil && !bf.Rotation.NearEquals(nextBf.Rotation, 1e-2) {
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
	for boneName := range boneFrames.data {
		if !boneFrames.ContainsActive(boneName) {
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
