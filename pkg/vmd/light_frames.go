package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type LightFrames struct {
	*mcore.IndexFloatModelCorrection[*LightFrame]
	RegisteredIndexes *mcore.TreeIndexes[mcore.Float32] // 登録対象キーフレリスト
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*LightFrame](),
		RegisteredIndexes:         mcore.NewTreeIndexes[mcore.Float32](),
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (lfs *LightFrames) GetRangeIndexes(index mcore.Float32) (mcore.Float32, mcore.Float32) {
	if lfs.RegisteredIndexes.IsEmpty() {
		return 0.0, 0.0
	}

	nowIndexes := lfs.RegisteredIndexes.Search(index)
	if nowIndexes != nil {
		return index, index
	}

	var prevIndex, nextIndex mcore.Float32

	prevIndexes := lfs.RegisteredIndexes.SearchLeft(index)
	nextIndexes := lfs.RegisteredIndexes.SearchRight(index)

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
func (lfs *LightFrames) GetItem(index mcore.Float32) *LightFrame {
	if val, ok := lfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := lfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && lfs.Indexes.Contains(nextIndex) {
		nextLf := lfs.Data[nextIndex]
		copied := &LightFrame{
			BaseFrame: NewVmdBaseFrame(index),
			Position:  nextLf.Position.Copy(),
			Color:     nextLf.Color.Copy(),
		}
		return copied
	}

	var prevLf, nextLf *LightFrame
	if lfs.Indexes.Contains(prevIndex) {
		prevLf = lfs.Data[prevIndex]
	} else {
		prevLf = NewLightFrame(index)
	}
	if lfs.Indexes.Contains(nextIndex) {
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
	if !lfs.Indexes.Contains(value.Index) {
		lfs.Indexes.Insert(value.Index)
	}
	if value.Registered {
		if !lfs.RegisteredIndexes.Contains(value.Index) {
			lfs.RegisteredIndexes.Insert(value.Index)
		}
	}

	lfs.Data[value.Index] = value
}
