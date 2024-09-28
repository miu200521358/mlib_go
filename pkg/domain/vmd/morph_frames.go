package vmd

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type MorphFrames struct {
	Data map[string]*MorphNameFrames
}

func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		Data: make(map[string]*MorphNameFrames, 0),
	}
}

func (morphFrames *MorphFrames) Copy() *MorphFrames {
	copied := NewMorphFrames()
	for _, morphNameFrames := range morphFrames.Data {
		copied.Update(morphNameFrames.Copy())
	}
	return copied
}

func (morphFrames *MorphFrames) Contains(morphName string) bool {
	if _, ok := morphFrames.Data[morphName]; ok {
		if morphFrames.Data[morphName] != nil && morphFrames.Data[morphName].Len() > 0 {
			return true
		}
	}

	return false
}

func (morphFrames *MorphFrames) Update(morphNameFrames *MorphNameFrames) {
	morphFrames.Data[morphNameFrames.Name] = morphNameFrames
}

func (morphFrames *MorphFrames) Names() []string {
	names := make([]string, 0, len(morphFrames.Data))
	for name := range morphFrames.Data {
		names = append(names, name)
	}
	return names
}

func (morphFrames *MorphFrames) Get(morphName string) *MorphNameFrames {
	if !morphFrames.Contains(morphName) {
		morphFrames.Update(NewMorphNameFrames(morphName))
	}
	return morphFrames.Data[morphName]
}

func (morphFrames *MorphFrames) MaxFrame() float32 {
	maxFno := float32(0)
	for _, mnfs := range morphFrames.Data {
		fno := float32(mnfs.MaxFrame())
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (morphFrames *MorphFrames) MinFrame() float32 {
	minFno := float32(math.MaxFloat32)
	for _, mnfs := range morphFrames.Data {
		fno := float32(mnfs.MinFrame())
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (morphFrames *MorphFrames) Len() int {
	count := 0
	for _, fs := range morphFrames.Data {
		count += fs.RegisteredIndexes.Len()
	}
	return count
}

func (morphFrames *MorphFrames) registeredFramesMap() map[float32]struct{} {
	frames := make(map[float32]struct{}, 0)
	for _, boneFrames := range morphFrames.Data {
		for _, f := range boneFrames.Indexes.List() {
			frames[f] = struct{}{}
		}
	}
	return frames
}

func (morphFrames *MorphFrames) RegisteredFrames() []int {
	mFrames := morphFrames.registeredFramesMap()

	frames := make([]int, 0, len(mFrames))
	for f := range mFrames {
		frames = append(frames, int(f))
	}
	mmath.SortInts(frames)

	return frames
}
