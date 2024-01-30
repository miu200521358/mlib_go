package vmd

import (
	"slices"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type LightFrames struct {
	*mcore.IndexModelCorrection[*LightFrame]
	RegisteredIndexes []int // 登録対象キーフレリスト
}

func NewLightNameFrames(name string) *LightFrames {
	return &LightFrames{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*LightFrame](),
		RegisteredIndexes:    []int{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (lfs *LightFrames) GetRangeIndexes(index int) (int, int) {

	prevIndex := 0
	nextIndex := index

	if idx := sort.SearchInts(lfs.Indexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = lfs.Indexes[idx-1]
	}

	if idx := sort.SearchInts(lfs.Indexes, index); idx == len(lfs.Indexes) {
		nextIndex = slices.Max(lfs.Indexes)
	} else {
		nextIndex = lfs.Indexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (lfs *LightFrames) GetItem(index int) *LightFrame {
	if index < 0 {
		// マイナス指定の場合、後ろからの順番に置き換える
		index = len(lfs.Data) + index
		return lfs.Data[lfs.Indexes[index]]
	}
	if val, ok := lfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := lfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && slices.Contains(lfs.Indexes, nextIndex) {
		nextLf := lfs.Data[nextIndex]
		copied := &LightFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Position:  nextLf.Position.Copy(),
			Color:     nextLf.Color.Copy(),
		}
		return copied
	}

	var prevLf, nextLf *LightFrame
	if slices.Contains(lfs.Indexes, prevIndex) {
		prevLf = lfs.Data[prevIndex]
	} else {
		prevLf = NewLightFrame(index)
	}
	if slices.Contains(lfs.Indexes, nextIndex) {
		nextLf = lfs.Data[nextIndex]
	} else {
		nextLf = NewLightFrame(index)
	}

	lf := NewLightFrame(index)

	t := (float64(index) - float64(prevIndex)) / (float64(nextIndex) - float64(prevIndex))

	position := mmath.LerpVec3(prevLf.Position, nextLf.Position, t)
	lf.Position = &position

	color := mmath.LerpVec3(prevLf.Color, nextLf.Color, t)
	lf.Color = &color

	return lf
}

func (lfs *LightFrames) Append(value *LightFrame) {
	if !slices.Contains(lfs.Indexes, value.Index) {
		lfs.Indexes = append(lfs.Indexes, value.Index)
		sort.Ints(lfs.Indexes)
	}
	if value.Registered {
		if !slices.Contains(lfs.RegisteredIndexes, value.Index) {
			lfs.RegisteredIndexes = append(lfs.RegisteredIndexes, value.Index)
			sort.Ints(lfs.RegisteredIndexes)
		}
	}

	lfs.Data[value.Index] = value
}
