package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"

)

type IkFrames struct {
	*mcore.IndexFloatModelCorrection[*IkFrame]
	RegisteredIndexes *mcore.TreeIndexes[mcore.Float32] // 登録対象キーフレリスト
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*IkFrame](),
		RegisteredIndexes:         mcore.NewTreeIndexes[mcore.Float32](),
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (ifs *IkFrames) GetRangeIndexes(index mcore.Float32) (mcore.Float32, mcore.Float32) {
	if ifs.RegisteredIndexes.IsEmpty() {
		return 0.0, 0.0
	}

	nowIndexes := ifs.RegisteredIndexes.Search(index)
	if nowIndexes != nil {
		return index, index
	}

	var prevIndex, nextIndex mcore.Float32

	prevIndexes := ifs.RegisteredIndexes.SearchLeft(index)
	nextIndexes := ifs.RegisteredIndexes.SearchRight(index)

	if prevIndexes == nil {
		prevIndex = mcore.NewFloat32(0)
	} else {
		prevIndex = prevIndexes.Value
	}

	if nextIndexes == nil {
		nextIndex = index
	} else {
		nextIndex = nextIndexes.Value
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (ifs *IkFrames) GetItem(index mcore.Float32) *IkFrame {
	if val, ok := ifs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := ifs.GetRangeIndexes(index)

	if prevIndex == nextIndex && ifs.Indexes.Contains(nextIndex) {
		nextIf := ifs.Data[nextIndex]
		copied := nextIf.Copy()
		return copied.(*IkFrame)
	}

	var prevIf *IkFrame
	if ifs.Indexes.Contains(prevIndex) {
		prevIf = ifs.Data[prevIndex]
	} else {
		prevIf = NewIkFrame(index)
	}

	nif := prevIf.Copy()
	return nif.(*IkFrame)
}

func (ifs *IkFrames) Append(value *IkFrame) {
	if !ifs.Indexes.Contains(value.Index) {
		ifs.Indexes.Insert(value.Index)
	}
	if value.Registered {
		if !ifs.RegisteredIndexes.Contains(value.Index) {
			ifs.RegisteredIndexes.Insert(value.Index)
		}
	}

	ifs.Data[value.Index] = value
}
