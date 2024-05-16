package mcore

import (
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
type IndexModels[T IIndexModel] struct {
	Data map[int]T
}

func NewIndexModels[T IIndexModel]() *IndexModels[T] {
	return &IndexModels[T]{
		Data: make(map[int]T, 0),
	}
}

func (c *IndexModels[T]) GetItem(index int) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IndexModels[T]) SetItem(index int, v T) {
	c.Data[index] = v
}

func (c *IndexModels[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data[value.GetIndex()] = value
}

func (c *IndexModels[T]) DeleteItem(index int) {
	delete(c.Data, index)
}

func (c *IndexModels[T]) Len() int {
	return len(c.Data)
}

func (c *IndexModels[T]) Contains(key int) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexModels[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexModels[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}
