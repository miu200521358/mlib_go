package mcore

import (
	"slices"

	"github.com/jinzhu/copier"
)

type IIndexModel interface {
	IsValid() bool
	Copy() IIndexModel
	GetIndex() int
	SetIndex(index int)
}

// INDEXを持つ基底クラス
type IndexModel struct {
	Index int
}

func (v *IndexModel) GetIndex() int {
	return v.Index
}

func (v *IndexModel) SetIndex(index int) {
	v.Index = index
}

func (v *IndexModel) IsValid() bool {
	return v.GetIndex() >= 0
}

func (v *IndexModel) Copy() IIndexModel {
	copied := IndexModel{Index: v.Index}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IIndexModels[T IIndexModel] struct {
	Data    map[int]T
	Indexes map[int]int
}

func NewIndexModelCorrection[T IIndexModel]() *IIndexModels[T] {
	return &IIndexModels[T]{
		Data:    make(map[int]T, 0),
		Indexes: make(map[int]int, 0),
	}
}

func (c *IIndexModels[T]) GetItem(index int) T {
	if index < 0 {
		// マイナス指定の場合、後ろからの順番に置き換える
		index = len(c.Data) + index
		return c.Data[c.Indexes[index]]
	}
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IIndexModels[T]) SetItem(index int, v T) {
	c.Data[index] = v
}

func (c *IIndexModels[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data[value.GetIndex()] = value
	if _, ok := c.Indexes[value.GetIndex()]; !ok {
		c.Indexes[value.GetIndex()] = value.GetIndex()
	}
}

func (c *IIndexModels[T]) GetSortedIndexes() []int {
	keys := make([]int, 0, len(c.Indexes))
	for key := range c.Indexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func (c *IIndexModels[T]) DeleteItem(index int) {
	delete(c.Data, index)
}

func (c *IIndexModels[T]) Len() int {
	return len(c.Data)
}

func (c *IIndexModels[T]) Contains(key int) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IIndexModels[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IIndexModels[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IIndexModels[T]) LastIndex() int {
	maxIndex := 0
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (c *IIndexModels[T]) GetSortedData() []T {
	sortedData := make([]T, len(c.Indexes))
	for i, index := range c.GetSortedIndexes() {
		sortedData[i] = c.Data[index]
	}
	return sortedData
}
