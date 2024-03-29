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

type FloatIndexes struct {
	*llrb.LLRB
}

func NewFloatIndexes() *FloatIndexes {
	return &FloatIndexes{
		LLRB: llrb.New(),
	}
}

func (i FloatIndexes) Has(index float32) bool {
	return i.LLRB.Has(Float32(index))
}

func (i FloatIndexes) Max() float32 {
	if i.LLRB.Len() == 0 {
		return 0
	}
	return float32(i.LLRB.Max().(Float32))
}

// Tのリスト基底クラス
type IndexFloatModels[T IIndexFloatModel] struct {
	Data    map[float32]T
	Indexes *FloatIndexes
}

func NewIndexFloatModelCorrection[T IIndexFloatModel]() *IndexFloatModels[T] {
	return &IndexFloatModels[T]{
		Data:    make(map[float32]T, 0),
		Indexes: NewFloatIndexes(),
	}
}

func (c *IndexFloatModels[T]) GetItem(index float32) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IndexFloatModels[T]) SetItem(index float32, v T) {
	c.Data[index] = v
}

func (c *IndexFloatModels[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(float32(len(c.Data)))
	}
	c.Data[value.GetIndex()] = value
	c.Indexes.ReplaceOrInsert(Float32(value.GetIndex()))
}

func (c *IndexFloatModels[T]) Delete(index float32) {
	delete(c.Data, index)
	c.Indexes.Delete(Float32(index))
}

func (c *IndexFloatModels[T]) Len() int {
	return len(c.Data)
}

func (c *IndexFloatModels[T]) Contains(key float32) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexFloatModels[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexFloatModels[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexFloatModels[T]) LastIndex() float32 {
	maxIndex := float32(0.0)
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}
