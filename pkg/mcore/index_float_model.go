package mcore

import (
	"slices"

	"github.com/jinzhu/copier"
)

type IIndexFloatModel interface {
	IsValid() bool
	Copy() IIndexFloatModel
	GetIndex() float32
	SetIndex(index float32)
}

// INDEXを持つ基底クラス
type IndexFloatModel struct {
	Index float32
}

func (v *IndexFloatModel) GetIndex() float32 {
	return v.Index
}

func (v *IndexFloatModel) SetIndex(index float32) {
	v.Index = index
}

func (v *IndexFloatModel) IsValid() bool {
	return v.GetIndex() >= 0
}

func (v *IndexFloatModel) Copy() IIndexFloatModel {
	copied := IndexFloatModel{Index: v.Index}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IIndexFloatModels[T IIndexFloatModel] struct {
	Data    map[float32]T
	Indexes []float32
}

func NewIndexFloatModels[T IIndexFloatModel]() *IIndexFloatModels[T] {
	return &IIndexFloatModels[T]{
		Data:    make(map[float32]T, 0),
		Indexes: make([]float32, 0),
	}
}

func (c *IIndexFloatModels[T]) GetItem(index float32) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IIndexFloatModels[T]) SetItem(index float32, v T) {
	c.Data[index] = v
}

func (c *IIndexFloatModels[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(float32(len(c.Data)))
	}
	c.Data[value.GetIndex()] = value
	if !slices.Contains(c.Indexes, value.GetIndex()) {
		c.Indexes = append(c.Indexes, value.GetIndex())
	}
	c.SortIndexes()
}

func (c *IIndexFloatModels[T]) SortIndexes() {
	slices.Sort(c.Indexes)
}

func (c *IIndexFloatModels[T]) DeleteItem(index float32) {
	delete(c.Data, index)
}

func (c *IIndexFloatModels[T]) Len() int {
	return len(c.Data)
}

func (c *IIndexFloatModels[T]) Contains(key float32) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IIndexFloatModels[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IIndexFloatModels[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IIndexFloatModels[T]) LastIndex() float32 {
	maxIndex := float32(0.0)
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (c *IIndexFloatModels[T]) GetSortedData() []T {
	sortedData := make([]T, len(c.Indexes))
	for i, index := range c.Indexes {
		sortedData[i] = c.Data[index]
	}
	return sortedData
}
