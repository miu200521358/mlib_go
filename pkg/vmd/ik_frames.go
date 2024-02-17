package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type IkFrames struct {
	*mcore.IndexFloatModelCorrection[*IkFrame]
	RegisteredIndexes map[float32]float32 // 登録対象キーフレリスト
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*IkFrame](),
		RegisteredIndexes:         make(map[float32]float32, 0),
	}
}

func (c *IkFrames) ContainsRegistered(key float32) bool {
	_, ok := c.RegisteredIndexes[key]
	return ok
}

func (c *IkFrames) GetSortedRegisteredIndexes() []float32 {
	keys := make([]float32, 0, len(c.RegisteredIndexes))
	for key := range c.RegisteredIndexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// 指定したキーフレの前後のキーフレ番号を返す
func (ifs *IkFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	ikIndexes := ifs.GetSortedIndexes()

	if idx := mutils.SearchFloat32s(ikIndexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = ikIndexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(ikIndexes, index); idx == len(ikIndexes) {
		nextIndex = slices.Max(ikIndexes)
	} else {
		nextIndex = ikIndexes[idx]
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

	if prevIndex == nextIndex && ifs.Contains(nextIndex) {
		nextIf := ifs.Data[nextIndex]
		copied := nextIf.Copy()
		return copied.(*IkFrame)
	}

	var prevIf *IkFrame
	if ifs.Contains(prevIndex) {
		prevIf = ifs.Data[prevIndex]
	} else {
		prevIf = NewIkFrame(index)
	}

	nif := prevIf.Copy()
	return nif.(*IkFrame)
}

func (ifs *IkFrames) Append(value *IkFrame) {
	if !ifs.Contains(value.Index) {
		ifs.Indexes[value.Index] = value.Index
	}
	if value.Registered {
		if !ifs.ContainsRegistered(value.Index) {
			ifs.RegisteredIndexes[value.Index] = value.Index
		}
	}

	ifs.Data[value.Index] = value
}
