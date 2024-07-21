package core

import (
	"reflect"

	"github.com/jinzhu/copier"
)

type IIndexModel interface {
	IsValid() bool
	Copy() IIndexModel
	Index() int
	SetIndex(index int)
}

// INDEXを持つ基底クラス
type IndexModel struct {
	index int
}

func NewIndexModel(index int) *IndexModel {
	return &IndexModel{index: index}
}

func (iModel *IndexModel) Index() int {
	return iModel.index
}

func (iModel *IndexModel) SetIndex(index int) {
	iModel.index = index
}

func (iModel *IndexModel) IsValid() bool {
	return iModel != nil && iModel.Index() >= 0
}

func (iModel *IndexModel) Copy() IIndexModel {
	copied := IndexModel{index: iModel.index}
	copier.CopyWithOption(copied, iModel, copier.Option{DeepCopy: true})
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

func (iModels *IndexModels[T]) Get(index int) T {
	if index < 0 || index >= len(iModels.Data) {
		return iModels.nilFunc()
	}

	return iModels.Data[index]
}

func (iModels *IndexModels[T]) SetItem(index int, v T) {
	iModels.Data[index] = v
}

func (iModels *IndexModels[T]) Update(value T) {
	if value.Index() < 0 {
		panic("Index is not set")
	}
	iModels.Data[value.Index()] = value
	iModels.SetDirty(true)
}

func (iModels *IndexModels[T]) Append(value T) {
	if value.Index() < 0 {
		value.SetIndex(len(iModels.Data))
	}
	iModels.Data = append(iModels.Data, value)
	iModels.SetDirty(true)
}

func (iModels *IndexModels[T]) DeleteItem(index int) {
	iModels.Data[index] = iModels.nilFunc()
}

func (iModels *IndexModels[T]) Len() int {
	return len(iModels.Data)
}

func (iModels *IndexModels[T]) Contains(key int) bool {
	return iModels != nil && key >= 0 && key < len(iModels.Data) && !reflect.ValueOf(iModels.Data[key]).IsNil()
}

func (iModels *IndexModels[T]) IsDirty() bool {
	return iModels.isDirty
}

func (iModels *IndexModels[T]) SetDirty(dirty bool) {
	iModels.isDirty = dirty
}
