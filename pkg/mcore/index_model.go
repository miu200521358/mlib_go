package mcore

import (
	"github.com/jinzhu/copier"
)

type IIndexIntModel interface {
	IsValid() bool
	Copy() IIndexIntModel
	GetIndex() Int
	SetIndex(index Int)
}

// INDEXを持つ基底クラス
type IndexModel struct {
	Index Int
}

func (v *IndexModel) GetIndex() Int {
	return v.Index
}

func (v *IndexModel) SetIndex(index Int) {
	v.Index = index
}

func (v *IndexModel) IsValid() bool {
	return v.GetIndex() >= 0
}

func (v *IndexModel) Copy() IIndexIntModel {
	copied := IndexModel{Index: v.Index}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IndexModelCorrection[T IIndexIntModel] struct {
	Data    map[Int]T
	Indexes *TreeIndexes[Int]
}

func NewIndexModelCorrection[T IIndexIntModel]() *IndexModelCorrection[T] {
	return &IndexModelCorrection[T]{
		Data:    make(map[Int]T, 0),
		Indexes: NewTreeIndexes[Int](),
	}
}

func (c *IndexModelCorrection[T]) GetItem(index Int) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[IndexModelCorrection] index out of range: index: " + string(rune(index)))
}

func (c *IndexModelCorrection[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(Int(len(c.Data)))
	}
	c.Data[value.GetIndex()] = value
	c.Indexes.Insert(value.GetIndex())
}

func (c *IndexModelCorrection[T]) DeleteItem(index Int) {
	delete(c.Data, index)
}

func (c *IndexModelCorrection[T]) Len() int {
	return len(c.Data)
}

func (c *IndexModelCorrection[T]) Contains(key Int) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexModelCorrection[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexModelCorrection[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexModelCorrection[T]) LastIndex() Int {
	if c.IsEmpty() {
		return -1
	}
	return c.Indexes.GetMax()
}

func (c *IndexModelCorrection[T]) GetSortedData() []T {
	sortedData := make([]T, len(c.Indexes.GetValues()))
	for i, index := range c.Indexes.GetValues() {
		sortedData[i] = c.Data[index]
	}
	return sortedData
}
