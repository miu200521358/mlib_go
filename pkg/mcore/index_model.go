package mcore

import (
	"github.com/jinzhu/copier"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type IIndexModel interface {
	IsValid() bool
	Copy() IIndexModel
	GetIndex() int
	SetIndex(index int)
	GetMapKey() mmath.MVec3
	GetMapValue() *mmath.MVec3
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

func (v *IndexModel) GetMapKey() mmath.MVec3 {
	return *mmath.MVec3Zero
}

func (v *IndexModel) GetMapValue() *mmath.MVec3 {
	return nil
}

func (v *IndexModel) Copy() IIndexModel {
	copied := IndexModel{Index: v.Index}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IndexModels[T IIndexModel] struct {
	Data     map[int]T
	nilFunc  func() T
	IndexMap map[mmath.MVec3][]int
}

func NewIndexModels[T IIndexModel](nilFunc func() T) *IndexModels[T] {
	return &IndexModels[T]{
		Data:     make(map[int]T, 0),
		nilFunc:  nilFunc,
		IndexMap: make(map[mmath.MVec3][]int, 0),
	}
}

func (c *IndexModels[T]) SetupMapKeys() {
	c.IndexMap = make(map[mmath.MVec3][]int)
	for k, v := range c.Data {
		baseKey := v.GetMapKey()
		// 前後のオフセット込みでマッピング
		for _, offset := range []*mmath.MVec3{
			{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1},
			{0, 0, 0}, {-1, 0, 0}, {0, -1, 0}, {0, 0, -1},
			{1, 1, 0}, {1, 0, 1}, {0, 1, 1}, {1, 1, 1},
			{-1, -1, 0}, {-1, 0, -1}, {0, -1, -1}, {-1, -1, -1},
			{1, -1, 0}, {1, 0, -1}, {0, 1, -1}, {1, -1, 1},
			{-1, 1, 0}, {-1, 0, 1}, {0, -1, 1}, {-1, 1, -1},
		} {
			key := *baseKey.Added(offset)
			if _, ok := c.IndexMap[key]; !ok {
				c.IndexMap[key] = make([]int, 0)
			}
			c.IndexMap[key] = append(c.IndexMap[key], k)
		}
	}
}

func (c *IndexModels[T]) GetMapValues(v T) ([]int, []*mmath.MVec3) {
	if c.Data == nil {
		return nil, nil
	}
	key := v.GetMapKey()
	values := make([]*mmath.MVec3, 0)
	if indexes, ok := c.IndexMap[key]; ok {
		for _, i := range indexes {
			values = append(values, c.Get(i).GetMapValue())
		}
		return indexes, values
	}
	return nil, nil
}

func (c *IndexModels[T]) Get(index int) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	return c.nilFunc()
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
