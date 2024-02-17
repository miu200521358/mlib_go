package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type MorphNameFrames struct {
	*mcore.IndexFloatModelCorrection[*MorphFrame]
	Name              string    // ボーン名
	RegisteredIndexes []float32 // 登録対象キーフレリスト
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*MorphFrame](),
		Name:                      name,
		RegisteredIndexes:         []float32{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (mnfs *MorphNameFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	if idx := mutils.SearchFloat32s(mnfs.Indexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = mnfs.Indexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(mnfs.Indexes, index); idx == len(mnfs.Indexes) {
		nextIndex = slices.Max(mnfs.Indexes)
	} else {
		nextIndex = mnfs.Indexes[idx]
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

	if prevIndex == nextIndex && slices.Contains(mnfs.Indexes, nextIndex) {
		nextMf := mnfs.Data[nextIndex]
		copied := &MorphFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Ratio:     nextMf.Ratio,
		}
		return copied
	}

	var prevMf, nextMf *MorphFrame
	if slices.Contains(mnfs.Indexes, prevIndex) {
		prevMf = mnfs.Data[prevIndex]
	} else {
		prevMf = NewMorphFrame(index)
	}
	if slices.Contains(mnfs.Indexes, nextIndex) {
		nextMf = mnfs.Data[nextIndex]
	} else {
		nextMf = NewMorphFrame(index)
	}

	mf := NewMorphFrame(index)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*float64(index-prevIndex)/float64(nextIndex-prevIndex)

	return mf
}

func (mnfs *MorphNameFrames) Append(value *MorphFrame) {
	if !slices.Contains(mnfs.Indexes, value.Index) {
		mnfs.Indexes = append(mnfs.Indexes, value.Index)
		mutils.SortFloat32s(mnfs.Indexes)
	}

	if value.Registered {
		if !slices.Contains(mnfs.RegisteredIndexes, value.Index) {
			mnfs.RegisteredIndexes = append(mnfs.RegisteredIndexes, value.Index)
			mutils.SortFloat32s(mnfs.RegisteredIndexes)
		}
	}

	mnfs.Data[value.Index] = value
}
