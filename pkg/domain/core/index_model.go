package core

import (
	"reflect"

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
	return v != nil && v.GetIndex() >= 0
}

func (v *IndexModel) Copy() IIndexModel {
	copied := IndexModel{Index: v.Index}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IndexModels[T IIndexModel] struct {
	Data    []T
	nilFunc func() T
	isDirty bool
}

func NewIndexModels[T IIndexModel](count int, nilFunc func() T) *IndexModels[T] {
	return &IndexModels[T]{
		Data:    make([]T, count),
		nilFunc: nilFunc,
	}
}

func (c *IndexModels[T]) Get(index int) T {
	if index < 0 || index >= len(c.Data) {
		return c.nilFunc()
	}

	return c.Data[index]
}

func (c *IndexModels[T]) SetItem(index int, v T) {
	c.Data[index] = v
}

func (c *IndexModels[T]) Update(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data[value.GetIndex()] = value
	c.SetDirty(true)
}

func (c *IndexModels[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data = append(c.Data, value)
	c.SetDirty(true)
}

func (c *IndexModels[T]) DeleteItem(index int) {
	c.Data[index] = c.nilFunc()
}

func (c *IndexModels[T]) Len() int {
	return len(c.Data)
}

func (c *IndexModels[T]) Contains(key int) bool {
	return c != nil && key >= 0 && key < len(c.Data) && !reflect.ValueOf(c.Data[key]).IsNil()
}

func (c *IndexModels[T]) IsDirty() bool {
	return c.isDirty
}

func (c *IndexModels[T]) SetDirty(dirty bool) {
	c.isDirty = dirty
}
