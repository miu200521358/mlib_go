package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type ShadowFrames struct {
	*mcore.IndexFloatModelCorrection[*ShadowFrame]
	RegisteredIndexes *mcore.TreeIndexes[mcore.Float32] // 登録対象キーフレリスト
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*ShadowFrame](),
		RegisteredIndexes:         mcore.NewTreeIndexes[mcore.Float32](),
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (sfs *ShadowFrames) GetRangeIndexes(index mcore.Float32) (mcore.Float32, mcore.Float32) {
	if sfs.RegisteredIndexes.IsEmpty() {
		return 0.0, 0.0
	}

	nowIndexes := sfs.RegisteredIndexes.Search(index)
	if nowIndexes != nil {
		return index, index
	}

	var prevIndex, nextIndex mcore.Float32

	prevIndexes := sfs.RegisteredIndexes.SearchLeft(index)
	nextIndexes := sfs.RegisteredIndexes.SearchRight(index)

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
func (sfs *ShadowFrames) GetItem(index mcore.Float32) *ShadowFrame {
	if val, ok := sfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := sfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && sfs.Indexes.Contains(nextIndex) {
		nextSf := sfs.Data[nextIndex]
		copied := &ShadowFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Distance:  nextSf.Distance,
		}
		return copied
	}

	var prevSf, nextSf *ShadowFrame
	if sfs.Indexes.Contains(prevIndex) {
		prevSf = sfs.Data[prevIndex]
	} else {
		prevSf = NewShadowFrame(index)
	}
	if sfs.Indexes.Contains(nextIndex) {
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
	if !sfs.Indexes.Contains(value.Index) {
		sfs.Indexes.Insert(value.Index)
	}
	if value.Registered {
		if !sfs.RegisteredIndexes.Contains(value.Index) {
			sfs.RegisteredIndexes.Insert(value.Index)
		}
	}

	sfs.Data[value.Index] = value
}
