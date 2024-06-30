package mcore

import (
	"sort"

	"github.com/jinzhu/copier"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type IIndexNameModel interface {
	IsValid() bool
	Copy() IIndexNameModel
	GetIndex() int
	SetIndex(index int)
	GetName() string
	SetName(name string)
	GetMapKey() mmath.MVec3
	GetMapValue() *mmath.MVec3
}

// INDEXを持つ基底クラス
type IndexNameModel struct {
	Index       int
	Name        string
	EnglishName string
}

func (v *IndexNameModel) GetIndex() int {
	return v.Index
}

func (v *IndexNameModel) SetIndex(index int) {
	v.Index = index
}

func (v *IndexNameModel) GetName() string {
	return v.Name
}

func (v *IndexNameModel) SetName(name string) {
	v.Name = name
}

func (v *IndexNameModel) IsValid() bool {
	return v.GetIndex() >= 0
}

func (v *IndexNameModel) GetMapKey() mmath.MVec3 {
	return *mmath.MVec3Zero
}

func (v *IndexNameModel) GetMapValue() *mmath.MVec3 {
	return nil
}

func (v *IndexNameModel) Copy() IIndexNameModel {
	copied := IndexNameModel{Index: v.Index, Name: v.Name, EnglishName: v.EnglishName}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IndexNameModels[T IIndexNameModel] struct {
	Data        map[int]T
	NameIndexes map[string]int
	IndexMap    map[mmath.MVec3]map[int]T
	nilFunc     func() T
}

func NewIndexNameModels[T IIndexNameModel](nilFunc func() T) *IndexNameModels[T] {
	return &IndexNameModels[T]{
		Data:        make(map[int]T, 0),
		NameIndexes: make(map[string]int, 0),
		IndexMap:    make(map[mmath.MVec3]map[int]T),
		nilFunc:     nilFunc,
	}
}

func (c *IndexNameModels[T]) SetupMapKeys() {
	c.IndexMap = make(map[mmath.MVec3]map[int]T)
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
				c.IndexMap[key] = make(map[int]T)
			}
			c.IndexMap[key][k] = v
		}
	}
}

func (c *IndexNameModels[T]) GetMapValues(v T) ([]int, []*mmath.MVec3) {
	if c.Data == nil {
		return nil, nil
	}
	key := v.GetMapKey()
	indexes := make([]int, 0)
	values := make([]*mmath.MVec3, 0)
	if mapIndexes, ok := c.IndexMap[key]; ok {
		for i, iv := range mapIndexes {
			indexes = append(indexes, i)
			values = append(values, iv.GetMapValue())
		}
		return indexes, values
	}
	return nil, nil
}

func (c *IndexNameModels[T]) Get(index int) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	return c.nilFunc()
}

func (c *IndexNameModels[T]) SetItem(index int, v T) {
	c.Data[index] = v
}

func (c *IndexNameModels[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data[value.GetIndex()] = value
	if _, ok := c.NameIndexes[value.GetName()]; !ok {
		// 名前は先勝ち
		c.NameIndexes[value.GetName()] = value.GetIndex()
	}
}

func (c *IndexNameModels[T]) GetIndexes() []int {
	indexes := make([]int, len(c.NameIndexes))
	i := 0
	for _, index := range c.NameIndexes {
		indexes[i] = index
		i++
	}
	sort.Ints(indexes)
	return indexes
}

func (c *IndexNameModels[T]) GetNames() []string {
	names := make([]string, len(c.NameIndexes))
	i := 0
	for index := range c.GetIndexes() {
		names[i] = c.Data[index].GetName()
		i++
	}
	return names
}

func (c *IndexNameModels[T]) DeleteItem(index int) {
	delete(c.Data, index)
}

func (c *IndexNameModels[T]) Len() int {
	return len(c.Data)
}

func (c *IndexNameModels[T]) ContainsIndex(key int) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexNameModels[T]) ContainsName(key string) bool {
	_, ok := c.NameIndexes[key]
	return ok
}

func (c *IndexNameModels[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexNameModels[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexNameModels[T]) LastIndex() int {
	maxIndex := 0
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (c *IndexNameModels[T]) GetByName(name string) T {
	if index, ok := c.NameIndexes[name]; ok {
		return c.Data[index]
	}
	return c.nilFunc()
}

func (v *IndexNameModels[T]) Contains(index int) bool {
	_, ok := v.Data[index]
	return ok
}

func (v *IndexNameModels[T]) ContainsByName(name string) bool {
	_, ok := v.NameIndexes[name]
	return ok
}
