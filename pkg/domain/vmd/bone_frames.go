package vmd

import (
	"math"
	"slices"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
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

func (boneFrames *BoneFrames) Copy() *BoneFrames {
	copied := NewBoneFrames()
	for _, boneNameFrames := range boneFrames.Data {
		copied.Append(boneNameFrames.Copy())
	}
	return copied
}

func (boneFrames *BoneFrames) Contains(boneName string) bool {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	if _, ok := boneFrames.Data[boneName]; ok {
		if boneFrames.Data[boneName] != nil && boneFrames.Data[boneName].Len() > 0 {
			return true
		}
	}

	return false
}

// ContainsActive 有効なキーフレが存在するか
func (boneFrames *BoneFrames) ContainsActive(boneName string) bool {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	if _, ok := boneFrames.Data[boneName]; !ok {
		return false
	}

	if boneFrames.Data[boneName].Len() == 0 {
		return false
	}

	for _, f := range boneFrames.Data[boneName].Indexes.List() {
		bf := boneFrames.Data[boneName].Get(f)
		if bf == nil {
			return false
		}

		if (bf.Position != nil && !bf.Position.NearEquals(mmath.MVec3Zero, 1e-2)) ||
			(bf.Rotation != nil && !bf.Rotation.NearEquals(mmath.MQuaternionIdent, 1e-2)) {
			return true
		}

		nextBf := boneFrames.Data[boneName].Get(boneFrames.Data[boneName].Indexes.Next(f))

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

	boneFrames.Data[boneNameFrames.Name] = boneNameFrames
}

func (boneFrames *BoneFrames) Delete(boneName string) {
	boneFrames.lock.Lock()
	defer boneFrames.lock.Unlock()

	delete(boneFrames.Data, boneName)
}

func (boneFrames *BoneFrames) Get(boneName string) *BoneNameFrames {
	if !boneFrames.Contains(boneName) {
		boneFrames.Append(NewBoneNameFrames(boneName))
	}

	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	return boneFrames.Data[boneName]
}

func (boneFrames *BoneFrames) Names() []string {
	names := make([]string, 0, len(boneFrames.Data))
	for name := range boneFrames.Data {
		names = append(names, name)
	}
	return names
}

func (boneFrames *BoneFrames) Indexes() []int {
	indexes := core.NewIntIndexes()
	for _, boneFrames := range boneFrames.Data {
		for _, f := range boneFrames.List() {
			indexes.ReplaceOrInsert(core.Int(f.index))
		}
	}
	return indexes.List()
}

func (boneFrames *BoneFrames) GetRegisteredIndexes() []int {
	indexes := core.NewIntIndexes()
	for _, boneFrames := range boneFrames.Data {
		for _, index := range boneFrames.RegisteredIndexes.List() {
			indexes.ReplaceOrInsert(core.NewInt(int(index)))
		}
	}
	return indexes.List()
}

func (boneFrames *BoneFrames) Len() int {
	count := 0
	for _, boneFrames := range boneFrames.Data {
		count += boneFrames.RegisteredIndexes.Len()
	}
	return count
}

func (boneFrames *BoneFrames) MaxFrame() float32 {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	maxFno := float32(0)
	for _, boneFrames := range boneFrames.Data {
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
	for _, boneFrames := range boneFrames.Data {
		fno := float32(boneFrames.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (boneFrames *BoneFrames) registeredFramesMap(boneNames []string) map[float32]struct{} {
	boneFrames.lock.RLock()
	defer boneFrames.lock.RUnlock()

	frames := make(map[float32]struct{}, 0)
	for boneName, boneFrames := range boneFrames.Data {
		if boneNames != nil && !slices.Contains(boneNames, boneName) {
			continue
		}
		for _, f := range boneFrames.Indexes.List() {
			frames[f] = struct{}{}
		}
	}
	return frames
}

func (boneFrames *BoneFrames) RegisteredFrames(boneNames []string) []int {
	bFrames := boneFrames.registeredFramesMap(boneNames)

	frames := make([]int, 0, len(bFrames))
	for f := range bFrames {
		frames = append(frames, int(f))
	}
	mmath.SortInts(frames)

	return frames
}

func (boneFrames *BoneFrames) Clean() {
	for boneName := range boneFrames.Data {
		if !boneFrames.ContainsActive(boneName) {
			boneFrames.Delete(boneName)
		}
	}
}
