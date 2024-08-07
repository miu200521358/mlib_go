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

func (morphFrames *MorphFrames) Contains(morphName string) bool {
	_, ok := morphFrames.Data[morphName]
	return ok
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
