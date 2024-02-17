package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"

)

type LightFrames struct {
	*mcore.IndexFloatModelCorrection[*LightFrame]
	RegisteredIndexes map[float32]float32 // 登録対象キーフレリスト
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*LightFrame](),
		RegisteredIndexes:         make(map[float32]float32, 0),
	}
}

func (c *LightFrames) ContainsRegistered(key float32) bool {
	_, ok := c.RegisteredIndexes[key]
	return ok
}

func (c *LightFrames) GetSortedRegisteredIndexes() []float32 {
	keys := make([]float32, 0, len(c.RegisteredIndexes))
	for key := range c.RegisteredIndexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// 指定したキーフレの前後のキーフレ番号を返す
func (lfs *LightFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	lightIndexes := lfs.GetSortedIndexes()

	if idx := mutils.SearchFloat32s(lightIndexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = lightIndexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(lightIndexes, index); idx == len(lightIndexes) {
		nextIndex = slices.Max(lightIndexes)
	} else {
		nextIndex = lightIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (lfs *LightFrames) GetItem(index float32) *LightFrame {
	if val, ok := lfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := lfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && lfs.Contains(nextIndex) {
		nextLf := lfs.Data[nextIndex]
		copied := &LightFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Position:  nextLf.Position.Copy(),
			Color:     nextLf.Color.Copy(),
		}
		return copied
	}

	var prevLf, nextLf *LightFrame
	if lfs.Contains(prevIndex) {
		prevLf = lfs.Data[prevIndex]
	} else {
		prevLf = NewLightFrame(index)
	}
	if lfs.Contains(nextIndex) {
		nextLf = lfs.Data[nextIndex]
	} else {
		nextLf = NewLightFrame(index)
	}

	lf := NewLightFrame(index)

	t := (float64(index) - float64(prevIndex)) / (float64(nextIndex) - float64(prevIndex))

	lf.Position = mmath.LerpVec3(prevLf.Position, nextLf.Position, t)
	lf.Color = mmath.LerpVec3(prevLf.Color, nextLf.Color, t)

	return lf
}

func (lfs *LightFrames) Append(value *LightFrame) {
	if !lfs.Contains(value.Index) {
		lfs.Indexes[value.Index] = value.Index
	}
	if value.Registered {
		if !lfs.ContainsRegistered(value.Index) {
			lfs.RegisteredIndexes[value.Index] = value.Index
		}
	}

	lfs.Data[value.Index] = value
}
