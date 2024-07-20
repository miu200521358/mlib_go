package core

import (
	"reflect"
	"sort"

	"github.com/jinzhu/copier"
)

type IIndexNameModel interface {
	IsValid() bool
	Copy() IIndexNameModel
	Index() int
	SetIndex(index int)
	Name() string
	SetName(name string)
	EnglishName() string
	SetEnglishName(englishName string)
}

// INDEXを持つ基底クラス
type IndexNameModel struct {
	index       int
	name        string
	englishName string
}

func NewIndexNameModel(index int, name string, englishName string) *IndexNameModel {
	return &IndexNameModel{index: index, name: name, englishName: englishName}
}

func (v *IndexNameModel) Index() int {
	return v.index
}

func (v *IndexNameModel) SetIndex(index int) {
	v.index = index
}

func (v *IndexNameModel) Name() string {
	return v.name
}

func (v *IndexNameModel) SetName(name string) {
	v.name = name
}

func (v *IndexNameModel) EnglishName() string {
	return v.englishName
}

func (v *IndexNameModel) SetEnglishName(englishName string) {
	v.englishName = englishName
}

func (v *IndexNameModel) IsValid() bool {
	return v != nil && v.index >= 0
}

func (v *IndexNameModel) Copy() IIndexNameModel {
	copied := IndexNameModel{index: v.index, name: v.name, englishName: v.englishName}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IndexNameModels[T IIndexNameModel] struct {
	Data        []T
	NameIndexes map[string]int
	nilFunc     func() T
	isDirty     bool
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
	if value.Index() < 0 {
		panic("Index is not set")
	}
	c.Data[value.Index()] = value
	if _, ok := c.NameIndexes[value.Name()]; !ok {
		// 名前は先勝ち
		c.NameIndexes[value.Name()] = value.Index()
	}
	c.SetDirty(true)
}

func (c *IndexNameModels[T]) Append(value T) {
	if value.Index() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data = append(c.Data, value)
	if _, ok := c.NameIndexes[value.Name()]; !ok {
		// 名前は先勝ち
		c.NameIndexes[value.Name()] = value.Index()
	}
	c.SetDirty(true)
}

func (c *IndexNameModels[T]) Indexes() []int {
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
	for index := range c.Indexes() {
		names[i] = c.Data[index].Name()
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

func (c *IndexNameModels[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexNameModels[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexNameModels[T]) GetByName(name string) T {
	if index, ok := c.NameIndexes[name]; ok {
		return c.Data[index]
	}
	return c.nilFunc()
}

func (v *IndexNameModels[T]) Contains(index int) bool {
	return index >= 0 && index < len(v.Data) && !reflect.ValueOf(v.Data[index]).IsNil()
}

func (v *IndexNameModels[T]) ContainsByName(name string) bool {
	_, ok := v.NameIndexes[name]
	return ok
}

func (v *IndexNameModels[T]) IsDirty() bool {
	return v.isDirty
}

func (v *IndexNameModels[T]) SetDirty(dirty bool) {
	v.isDirty = dirty
}
