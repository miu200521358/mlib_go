package mcore

import (
	"github.com/jinzhu/copier"
	"github.com/petar/GoLLRB/llrb"

)

type IIndexFloatModel interface {
	IsValid() bool
	Copy() IIndexFloatModel
	GetIndex() float32
	SetIndex(index float32)
	Less(than llrb.Item) bool
}

// INDEXを持つ基底クラス
type IndexFloatModel struct {
	Index float32
}

func NewIndexFloatModel(index float32) *IndexFloatModel {
	return &IndexFloatModel{
		Index: index,
	}
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

func (v *IndexFloatModel) Less(than llrb.Item) bool {
	if than == nil {
		return false
	}
	return v.Index < than.(IIndexFloatModel).GetIndex()
}

type IIndexFloatModels[T IIndexFloatModel] interface {
	Has(index float32) bool
	Get(index float32) T
	Append(value T)
	Insert(value T)
	Delete(index float32)
	Len() int
	IsEmpty() bool
	IsNotEmpty() bool
	LastIndex() float32
}

// Tのリスト基底クラス
type IndexFloatModels[T IIndexFloatModel] struct {
	*llrb.LLRB
}

func NewIndexFloatModels[T IIndexFloatModel]() *IndexFloatModels[T] {
	return &IndexFloatModels[T]{
		LLRB: llrb.New(),
	}
}

func (c *IndexFloatModels[T]) Get(index float32) T {
	return c.LLRB.Get(NewIndexFloatModel(index)).(T)
}

func (c *IndexFloatModels[T]) Append(value T) {
	c.LLRB.ReplaceOrInsert(value)
}

func (c *IndexFloatModels[T]) Insert(value T) {
	c.LLRB.ReplaceOrInsert(value)
}

func (c *IndexFloatModels[T]) Has(index float32) bool {
	return c.LLRB.Has(NewIndexFloatModel(index))
}

func (c *IndexFloatModels[T]) IsEmpty() bool {
	return c.Len() == 0
}

func (c *IndexFloatModels[T]) IsNotEmpty() bool {
	return !c.IsEmpty()
}

func (c *IndexFloatModels[T]) LastIndex() float32 {
	return c.Max().(T).GetIndex()
}
