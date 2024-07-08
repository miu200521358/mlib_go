package mcore

import (
	"reflect"
	"sort"

	"github.com/jinzhu/copier"
)

type IIndexNameModel interface {
	IsValid() bool
	Copy() IIndexNameModel
	GetIndex() int
	SetIndex(index int)
	GetName() string
	SetName(name string)
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
	return v != nil && v.GetIndex() >= 0
}

func (v *IndexNameModel) Copy() IIndexNameModel {
	copied := IndexNameModel{Index: v.Index, Name: v.Name, EnglishName: v.EnglishName}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IndexNameModels[T IIndexNameModel] struct {
	Data        []T
	NameIndexes map[string]int
	nilFunc     func() T
}

func NewIndexNameModels[T IIndexNameModel](count int, nilFunc func() T) *IndexNameModels[T] {
	return &IndexNameModels[T]{
		Data:        make([]T, count),
		NameIndexes: make(map[string]int, 0),
		nilFunc:     nilFunc,
	}
}

func (c *IndexNameModels[T]) Get(index int) T {
	if index < 0 || index >= len(c.Data) {
		return c.nilFunc()
	}
	return c.Data[index]
}

func (c *IndexNameModels[T]) SetItem(index int, v T) {
	c.Data[index] = v
}

func (c *IndexNameModels[T]) Update(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data[value.GetIndex()] = value
	if _, ok := c.NameIndexes[value.GetName()]; !ok {
		// 名前は先勝ち
		c.NameIndexes[value.GetName()] = value.GetIndex()
	}
}

func (c *IndexNameModels[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data = append(c.Data, value)
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
	c.Data[index] = c.nilFunc()
}

func (c *IndexNameModels[T]) Len() int {
	return len(c.Data)
}

func (c *IndexNameModels[T]) ContainsIndex(key int) bool {
	return c != nil && key >= 0 && key < len(c.Data) && !reflect.ValueOf(c.Data[key]).IsNil()
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
	return index >= 0 && index < len(v.Data) && v.Data[index].IsValid()
}

func (v *IndexNameModels[T]) ContainsByName(name string) bool {
	_, ok := v.NameIndexes[name]
	return ok
}
