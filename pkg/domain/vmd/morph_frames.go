package vmd

import (
	"math"
)

type MorphFrames struct {
	Data map[string]*MorphNameFrames
}

func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		Data: make(map[string]*MorphNameFrames, 0),
	}
}

func (mfs *MorphFrames) Contains(morphName string) bool {
	_, ok := mfs.Data[morphName]
	return ok
}

func (mfs *MorphFrames) Update(fs *MorphNameFrames) {
	mfs.Data[fs.Name] = fs
}

func (mfs *MorphFrames) GetNames() []string {
	names := make([]string, 0, len(mfs.Data))
	for name := range mfs.Data {
		names = append(names, name)
	}
	return names
}

func (mfs *MorphFrames) Get(morphName string) *MorphNameFrames {
	if !mfs.Contains(morphName) {
		mfs.Update(NewMorphNameFrames(morphName))
	}
	return mfs.Data[morphName]
}

func (mfs *MorphFrames) GetMaxFrame() int {
	maxFno := int(0)
	for _, mnfs := range mfs.Data {
		fno := mnfs.GetMaxFrame()
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (mfs *MorphFrames) GetMinFrame() int {
	minFno := math.MaxInt
	for _, mnfs := range mfs.Data {
		fno := mnfs.GetMinFrame()
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (fs *MorphFrames) Len() int {
	count := 0
	for _, fs := range fs.Data {
		count += fs.RegisteredIndexes.Len()
	}
	return count
}
