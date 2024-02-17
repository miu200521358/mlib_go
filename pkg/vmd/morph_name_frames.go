package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mutils"

)

type MorphNameFrames struct {
	*mcore.IndexFloatModelCorrection[*MorphFrame]
	Name              string              // ボーン名
	RegisteredIndexes map[float32]float32 // 登録対象キーフレリスト
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*MorphFrame](),
		Name:                      name,
		RegisteredIndexes:         make(map[float32]float32, 0),
	}
}

func (c *MorphNameFrames) ContainsRegistered(key float32) bool {
	_, ok := c.RegisteredIndexes[key]
	return ok
}

func (c *MorphNameFrames) GetSortedRegisteredIndexes() []float32 {
	keys := make([]float32, 0, len(c.RegisteredIndexes))
	for key := range c.RegisteredIndexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// 指定したキーフレの前後のキーフレ番号を返す
func (mnfs *MorphNameFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	morphIndexes := mnfs.GetSortedIndexes()

	if idx := mutils.SearchFloat32s(morphIndexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = morphIndexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(morphIndexes, index); idx == len(morphIndexes) {
		nextIndex = slices.Max(morphIndexes)
	} else {
		nextIndex = morphIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (mnfs *MorphNameFrames) GetItem(index float32) *MorphFrame {
	if val, ok := mnfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := mnfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && mnfs.Contains(nextIndex) {
		nextMf := mnfs.Data[nextIndex]
		copied := &MorphFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Ratio:     nextMf.Ratio,
		}
		return copied
	}

	var prevMf, nextMf *MorphFrame
	if mnfs.Contains(prevIndex) {
		prevMf = mnfs.Data[prevIndex]
	} else {
		prevMf = NewMorphFrame(index)
	}
	if mnfs.Contains(nextIndex) {
		nextMf = mnfs.Data[nextIndex]
	} else {
		nextMf = NewMorphFrame(index)
	}

	mf := NewMorphFrame(index)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*float64(index-prevIndex)/float64(nextIndex-prevIndex)

	return mf
}

func (mnfs *MorphNameFrames) Append(value *MorphFrame) {
	if !mnfs.Contains(value.Index) {
		mnfs.Indexes[value.Index] = value.Index
	}

	if value.Registered {
		if !mnfs.ContainsRegistered(value.Index) {
			mnfs.RegisteredIndexes[value.Index] = value.Index
		}
	}

	mnfs.Data[value.Index] = value
}
