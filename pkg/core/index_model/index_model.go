package index_model

import (
	"slices"

)

type IndexModelInterface interface {
	IsValid() bool
	Copy() IndexModelInterface
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

// Copy
func (v *IndexModel) Copy() IndexModelInterface {
	copied := *v
	return &copied
}

// Tのリスト基底クラス
type IndexModelCorrection[T IndexModelInterface] struct {
	Data    map[int]*T
	Indexes []int
}

func NewIndexModelCorrection[T IndexModelInterface]() *IndexModelCorrection[T] {
	return &IndexModelCorrection[T]{
		Data:    make(map[int]*T),
		Indexes: make([]int, 0),
	}
}

func (c *IndexModelCorrection[T]) GetItem(index int) T {
	if index < 0 {
		// マイナス指定の場合、後ろからの順番に置き換える
		index = len(c.Data) + index
		return *c.Data[c.Indexes[index]]
	}
	if val, ok := c.Data[index]; ok {
		return *val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IndexModelCorrection[T]) SetItem(index int, v *T) {
	c.Data[index] = v
}

func (c *IndexModelCorrection[T]) Append(value T, isSort bool) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data[value.GetIndex()] = &value
	if isSort {
		c.SortIndexes()
	} else {
		c.Indexes = append(c.Indexes, value.GetIndex())
	}
}

func (c *IndexModelCorrection[T]) SortIndexes() {
	slices.Sort(c.Indexes)
}

func (c *IndexModelCorrection[T]) DeleteItem(index int) {
	delete(c.Data, index)
}

func (c *IndexModelCorrection[T]) Len() int {
	return len(c.Data)
}

func (c *IndexModelCorrection[T]) Contains(key int) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexModelCorrection[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexModelCorrection[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexModelCorrection[T]) LastIndex() int {
	maxIndex := 0
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}
