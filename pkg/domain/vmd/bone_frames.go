package vmd

import (
	"math"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

type BoneFrames struct {
	Data map[string]*BoneNameFrames
	lock sync.RWMutex // マップアクセス制御用
}

func NewBoneFrames() *BoneFrames {
	return &BoneFrames{
		Data: make(map[string]*BoneNameFrames, 0),
		lock: sync.RWMutex{},
	}
}

func (fs *BoneFrames) Contains(boneName string) bool {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	_, ok := fs.Data[boneName]
	return ok
}

func (fs *BoneFrames) Append(nfs *BoneNameFrames) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Data[nfs.Name] = nfs
}

func (fs *BoneFrames) Delete(boneName string) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	delete(fs.Data, boneName)
}

func (fs *BoneFrames) Get(boneName string) *BoneNameFrames {
	if !fs.Contains(boneName) {
		fs.Append(NewBoneNameFrames(boneName))
	}

	fs.lock.RLock()
	defer fs.lock.RUnlock()

	return fs.Data[boneName]
}

func (fs *BoneFrames) GetNames() []string {
	names := make([]string, 0, len(fs.Data))
	for name := range fs.Data {
		names = append(names, name)
	}
	return names
}

func (fs *BoneFrames) GetIndexes() []int {
	indexes := core.NewIntIndexes()
	for _, fs := range fs.Data {
		for _, f := range fs.List() {
			indexes.ReplaceOrInsert(f.Index)
		}
	}
	return indexes.List()
}

func (fs *BoneFrames) GetRegisteredIndexes() []int {
	indexes := core.NewIntIndexes()
	for _, fs := range fs.Data {
		for _, index := range fs.RegisteredIndexes.List() {
			indexes.ReplaceOrInsert(core.NewInt(index))
		}
	}
	return indexes.List()
}

func (fs *BoneFrames) Len() int {
	count := 0
	for _, fs := range fs.Data {
		count += fs.RegisteredIndexes.Len()
	}
	return count
}

func (fs *BoneFrames) GetMaxFrame() int {
	maxFno := int(0)
	for _, fs := range fs.Data {
		fno := fs.GetMaxFrame()
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (fs *BoneFrames) GetMinFrame() int {
	minFno := math.MaxInt
	for _, fs := range fs.Data {
		fno := fs.GetMinFrame()
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}
