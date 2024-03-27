package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type LightFrames struct {
	*mcore.IndexFloatModels[*LightFrame]
	RegisteredIndexes []float32 // 登録対象キーフレリスト
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		IndexFloatModels:  mcore.NewIndexFloatModelCorrection[*LightFrame](),
		RegisteredIndexes: []float32{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (lfs *LightFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	if idx := mutils.SearchFloat32s(lfs.Indexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = lfs.Indexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(lfs.Indexes, index); idx == len(lfs.Indexes) {
		nextIndex = slices.Max(lfs.Indexes)
	} else {
		nextIndex = lfs.Indexes[idx]
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

	lf.Position = mmath.LerpVec3(prevLf.Position, nextLf.Position, t)
	lf.Color = mmath.LerpVec3(prevLf.Color, nextLf.Color, t)

	return lf
}

func (lfs *LightFrames) Append(value *LightFrame) {
	if !slices.Contains(lfs.Indexes, value.Index) {
		lfs.Indexes = append(lfs.Indexes, value.Index)
		mutils.SortFloat32s(lfs.Indexes)
	}
	if value.Registered {
		if !slices.Contains(lfs.RegisteredIndexes, value.Index) {
			lfs.RegisteredIndexes = append(lfs.RegisteredIndexes, value.Index)
			mutils.SortFloat32s(lfs.RegisteredIndexes)
		}
	}

	lfs.Data[value.Index] = value
}
