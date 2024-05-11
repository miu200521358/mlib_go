package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type IkFrames struct {
	*mcore.IIndexFloatModels[*IkFrame]
	RegisteredIndexes []float32 // 登録対象キーフレリスト
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		IIndexFloatModels: mcore.NewIndexFloatModels[*IkFrame](),
		RegisteredIndexes: []float32{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (ifs *IkFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	if idx := mutils.SearchFloat32s(ifs.Indexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = ifs.Indexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(ifs.Indexes, index); idx == len(ifs.Indexes) {
		nextIndex = slices.Max(ifs.Indexes)
	} else {
		nextIndex = ifs.Indexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (ifs *IkFrames) GetItem(index float32) *IkFrame {
	if val, ok := ifs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := ifs.GetRangeIndexes(index)

	if prevIndex == nextIndex && slices.Contains(ifs.Indexes, nextIndex) {
		nextIf := ifs.Data[nextIndex]
		copied := nextIf.Copy()
		return copied.(*IkFrame)
	}

	var prevIf *IkFrame
	if slices.Contains(ifs.Indexes, prevIndex) {
		prevIf = ifs.Data[prevIndex]
	} else {
		prevIf = NewIkFrame(index)
	}

	nif := prevIf.Copy()
	return nif.(*IkFrame)
}

func (ifs *IkFrames) Append(value *IkFrame) {
	if !slices.Contains(ifs.Indexes, value.Index) {
		ifs.Indexes = append(ifs.Indexes, value.Index)
		mutils.SortFloat32s(ifs.Indexes)
	}
	if value.Registered {
		if !slices.Contains(ifs.RegisteredIndexes, value.Index) {
			ifs.RegisteredIndexes = append(ifs.RegisteredIndexes, value.Index)
			mutils.SortFloat32s(ifs.RegisteredIndexes)
		}
	}

	ifs.Data[value.Index] = value
}
