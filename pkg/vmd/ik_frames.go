package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type IkFrames struct {
	*mcore.IndexModels[*IkFrame]
	RegisteredIndexes []int // 登録対象キーフレリスト
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		IndexModels:       mcore.NewIndexModels[*IkFrame](),
		RegisteredIndexes: []int{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (ifs *IkFrames) GetRangeIndexes(index int) (int, int) {

	prevIndex := 0
	nextIndex := index

	if idx := mutils.SearchInts(ifs.RegisteredIndexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = ifs.RegisteredIndexes[idx-1]
	}

	if idx := mutils.SearchInts(ifs.RegisteredIndexes, index); idx == len(ifs.RegisteredIndexes) {
		nextIndex = slices.Max(ifs.RegisteredIndexes)
	} else {
		nextIndex = ifs.RegisteredIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (ifs *IkFrames) GetItem(index int) *IkFrame {
	if val, ok := ifs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := ifs.GetRangeIndexes(index)

	if prevIndex == nextIndex && slices.Contains(ifs.RegisteredIndexes, nextIndex) {
		nextIf := ifs.Data[nextIndex]
		copied := nextIf.Copy()
		return copied.(*IkFrame)
	}

	var prevIf *IkFrame
	if slices.Contains(ifs.RegisteredIndexes, prevIndex) {
		prevIf = ifs.Data[prevIndex]
	} else {
		prevIf = NewIkFrame(index)
	}

	nif := prevIf.Copy()
	return nif.(*IkFrame)
}

func (ifs *IkFrames) Append(value *IkFrame) {
	if !slices.Contains(ifs.RegisteredIndexes, value.Index) {
		ifs.RegisteredIndexes = append(ifs.RegisteredIndexes, value.Index)
		mutils.SortInts(ifs.RegisteredIndexes)
	}
	if value.Registered {
		if !slices.Contains(ifs.RegisteredIndexes, value.Index) {
			ifs.RegisteredIndexes = append(ifs.RegisteredIndexes, value.Index)
			mutils.SortInts(ifs.RegisteredIndexes)
		}
	}

	ifs.Data[value.Index] = value
}
