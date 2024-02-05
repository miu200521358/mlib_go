package mcore

import (
	"slices"

)

type IndexFloatModelInterface interface {
	IsValid() bool
	Copy() IndexFloatModelInterface
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

// Copy
func (v *IndexFloatModel) Copy() IndexFloatModelInterface {
	copied := *v
	return &copied
}

// Tのリスト基底クラス
type IndexFloatModelCorrection[T IndexFloatModelInterface] struct {
	Data    map[float32]T
	Indexes []float32
}

func NewIndexFloatModelCorrection[T IndexFloatModelInterface]() *IndexFloatModelCorrection[T] {
	return &IndexFloatModelCorrection[T]{
		Data:    make(map[float32]T, 0),
		Indexes: make([]float32, 0),
	}
}

func (c *IndexFloatModelCorrection[T]) GetItem(index float32) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IndexFloatModelCorrection[T]) SetItem(index float32, v T) {
	c.Data[index] = v
}

func (c *IndexFloatModelCorrection[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(float32(len(c.Data)))
	}
	c.Data[value.GetIndex()] = value
	if !slices.Contains(c.Indexes, value.GetIndex()) {
		c.Indexes = append(c.Indexes, value.GetIndex())
	}
	c.SortIndexes()
}

func (c *IndexFloatModelCorrection[T]) SortIndexes() {
	slices.Sort(c.Indexes)
}

func (c *IndexFloatModelCorrection[T]) DeleteItem(index float32) {
	delete(c.Data, index)
}

func (c *IndexFloatModelCorrection[T]) Len() int {
	return len(c.Data)
}

func (c *IndexFloatModelCorrection[T]) Contains(key float32) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexFloatModelCorrection[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexFloatModelCorrection[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexFloatModelCorrection[T]) LastIndex() float32 {
	maxIndex := float32(0.0)
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (c *IndexFloatModelCorrection[T]) GetSortedData() []T {
	sortedData := make([]T, len(c.Indexes))
	for i, index := range c.Indexes {
		sortedData[i] = c.Data[index]
	}
	return sortedData
}
