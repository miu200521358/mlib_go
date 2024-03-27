package mcore

import (
	"github.com/jinzhu/copier"
)

type IIndexFloatModel interface {
	IsValid() bool
	Copy() IIndexFloatModel
	GetIndex() Float32
	SetIndex(index Float32)
}

// INDEXを持つ基底クラス
type IndexFloatModel struct {
	Index Float32
}

func (v *IndexFloatModel) GetIndex() Float32 {
	return v.Index
}

func (v *IndexFloatModel) SetIndex(index Float32) {
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
type IndexFloatModelCorrection[T IIndexFloatModel] struct {
	Data    map[Float32]T
	Indexes *TreeIndexes[Float32]
}

func NewIndexFloatModelCorrection[T IIndexFloatModel]() *IndexFloatModelCorrection[T] {
	return &IndexFloatModelCorrection[T]{
		Data:    make(map[Float32]T, 0),
		Indexes: NewTreeIndexes[Float32](),
	}
}

func (c *IndexFloatModelCorrection[T]) GetItem(index Float32) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IndexFloatModelCorrection[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(Float32(len(c.Data)))
	}
	c.Data[value.GetIndex()] = value
	c.Indexes.Insert(value.GetIndex())
}

func (c *IndexFloatModelCorrection[T]) DeleteItem(index Float32) {
	delete(c.Data, index)
}

func (c *IndexFloatModelCorrection[T]) Len() int {
	return len(c.Data)
}

func (c *IndexFloatModelCorrection[T]) Contains(key Float32) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexFloatModelCorrection[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexFloatModelCorrection[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexFloatModelCorrection[T]) LastIndex() Float32 {
	maxIndex := Float32(0.0)
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (c *IndexFloatModelCorrection[T]) GetSortedData() []T {
	sortedData := make([]T, len(c.Indexes.GetValues()))
	for i, index := range c.Indexes.GetValues() {
		sortedData[i] = c.Data[index]
	}
	return sortedData
}
