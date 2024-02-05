package vmd

import (
	"slices"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type ShadowFrames struct {
	*mcore.IndexModelCorrection[*ShadowFrame]
	RegisteredIndexes []int // 登録対象キーフレリスト
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*ShadowFrame](),
		RegisteredIndexes:    []int{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (sfs *ShadowFrames) GetRangeIndexes(index int) (int, int) {

	prevIndex := 0
	nextIndex := index

	if idx := sort.SearchInts(sfs.Indexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = sfs.Indexes[idx-1]
	}

	if idx := sort.SearchInts(sfs.Indexes, index); idx == len(sfs.Indexes) {
		nextIndex = slices.Max(sfs.Indexes)
	} else {
		nextIndex = sfs.Indexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (sfs *ShadowFrames) GetItem(index int) *ShadowFrame {
	if index < 0 {
		// マイナス指定の場合、後ろからの順番に置き換える
		index = len(sfs.Data) + index
		return sfs.Data[sfs.Indexes[index]]
	}
	if val, ok := sfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := sfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && slices.Contains(sfs.Indexes, nextIndex) {
		nextSf := sfs.Data[nextIndex]
		copied := &ShadowFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Distance:  nextSf.Distance,
		}
		return copied
	}

	var prevSf, nextSf *ShadowFrame
	if slices.Contains(sfs.Indexes, prevIndex) {
		prevSf = sfs.Data[prevIndex]
	} else {
		prevSf = NewShadowFrame(index)
	}
	if slices.Contains(sfs.Indexes, nextIndex) {
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
	if !slices.Contains(sfs.Indexes, value.Index) {
		sfs.Indexes = append(sfs.Indexes, value.Index)
		sort.Ints(sfs.Indexes)
	}
	if value.Registered {
		if !slices.Contains(sfs.RegisteredIndexes, value.Index) {
			sfs.RegisteredIndexes = append(sfs.RegisteredIndexes, value.Index)
			sort.Ints(sfs.RegisteredIndexes)
		}
	}

	sfs.Data[value.Index] = value
}
