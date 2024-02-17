package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type ShadowFrames struct {
	*mcore.IndexFloatModelCorrection[*ShadowFrame]
	RegisteredIndexes map[float32]float32 // 登録対象キーフレリスト
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*ShadowFrame](),
		RegisteredIndexes:         make(map[float32]float32, 0),
	}
}

func (c *ShadowFrames) ContainsRegistered(key float32) bool {
	_, ok := c.RegisteredIndexes[key]
	return ok
}

func (c *ShadowFrames) GetSortedRegisteredIndexes() []float32 {
	keys := make([]float32, 0, len(c.RegisteredIndexes))
	for key := range c.RegisteredIndexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// 指定したキーフレの前後のキーフレ番号を返す
func (sfs *ShadowFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	shadowIndexes := sfs.GetSortedIndexes()

	if idx := mutils.SearchFloat32s(shadowIndexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = shadowIndexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(shadowIndexes, index); idx == len(shadowIndexes) {
		nextIndex = slices.Max(shadowIndexes)
	} else {
		nextIndex = shadowIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (sfs *ShadowFrames) GetItem(index float32) *ShadowFrame {
	if val, ok := sfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := sfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && sfs.Contains(nextIndex) {
		nextSf := sfs.Data[nextIndex]
		copied := &ShadowFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Distance:  nextSf.Distance,
		}
		return copied
	}

	var prevSf, nextSf *ShadowFrame
	if sfs.Contains(prevIndex) {
		prevSf = sfs.Data[prevIndex]
	} else {
		prevSf = NewShadowFrame(index)
	}
	if sfs.Contains(nextIndex) {
		nextSf = sfs.Data[nextIndex]
	} else {
		nextSf = NewShadowFrame(index)
	}

	sf := NewShadowFrame(index)

	t := (float64(index) - float64(prevIndex)) / (float64(nextIndex) - float64(prevIndex))
	sf.Distance = mmath.LerpFloat(prevSf.Distance, nextSf.Distance, t)

	return sf
}

func (sfs *ShadowFrames) Append(value *ShadowFrame) {
	if !sfs.Contains(value.Index) {
		sfs.Indexes[value.Index] = value.Index
	}
	if value.Registered {
		if !sfs.ContainsRegistered(value.Index) {
			sfs.RegisteredIndexes[value.Index] = value.Index
		}
	}

	sfs.Data[value.Index] = value
}
