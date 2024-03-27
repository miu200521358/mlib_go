package mcore

import (
	"github.com/jinzhu/copier"

)

type IIndexNameModel interface {
	IsValid() bool
	Copy() IIndexNameModel
	GetIndex() Int
	SetIndex(index Int)
	GetName() string
	SetName(name string)
}

// INDEXを持つ基底クラス
type IndexNameModel struct {
	Index       Int
	Name        string
	EnglishName string
}

func (v *IndexNameModel) GetIndex() Int {
	return v.Index
}

func (v *IndexNameModel) SetIndex(index Int) {
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

func (v *IndexNameModel) Copy() IIndexNameModel {
	copied := IndexNameModel{Index: v.Index, Name: v.Name, EnglishName: v.EnglishName}
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return &copied
}

// Tのリスト基底クラス
type IndexNameModelCorrection[T IIndexNameModel] struct {
	Data        map[Int]T
	NameIndexes map[string]Int
	Indexes     *TreeIndexes[Int]
}

func NewIndexNameModelCorrection[T IIndexNameModel]() *IndexNameModelCorrection[T] {
	return &IndexNameModelCorrection[T]{
		Data:        make(map[Int]T, 0),
		NameIndexes: make(map[string]Int, 0),
		Indexes:     NewTreeIndexes[Int](),
	}
}

func (c *IndexNameModelCorrection[T]) GetItem(index Int) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IndexNameModelCorrection[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(Int(len(c.Data)))
	}
	c.Data[value.GetIndex()] = value
	c.Indexes.Insert(value.GetIndex())
	if _, ok := c.NameIndexes[value.GetName()]; !ok {
		// 名前は先勝ち
		c.NameIndexes[value.GetName()] = value.GetIndex()
	}
}

func (c *IndexNameModelCorrection[T]) GetNames() []string {
	names := make([]string, 0, len(c.NameIndexes))
	for _, index := range c.Indexes.GetValues() {
		names = append(names, c.Data[index].GetName())
	}
	return names
}

func (c *IndexNameModelCorrection[T]) DeleteItem(index Int) {
	delete(c.Data, index)
}

func (c *IndexNameModelCorrection[T]) Len() int {
	return len(c.Data)
}

func (c *IndexNameModelCorrection[T]) ContainsIndex(key Int) bool {
	_, ok := c.Data[key]
	return ok
}

func (c *IndexNameModelCorrection[T]) ContainsName(key string) bool {
	_, ok := c.NameIndexes[key]
	return ok
}

func (c *IndexNameModelCorrection[T]) IsEmpty() bool {
	return len(c.Data) == 0
}

func (c *IndexNameModelCorrection[T]) IsNotEmpty() bool {
	return len(c.Data) > 0
}

func (c *IndexNameModelCorrection[T]) LastIndex() Int {
	if c.IsEmpty() {
		return -1
	}
	return c.Indexes.GetMax()
}

func (c *IndexNameModelCorrection[T]) GetSortedData() []T {
	sortedData := make([]T, len(c.Indexes.GetValues()))
	for i, index := range c.Indexes.GetValues() {
		sortedData[i] = c.Data[index]
	}
	return sortedData
}

func (c *IndexNameModelCorrection[T]) GetItemByName(name string) T {
	if index, ok := c.NameIndexes[name]; ok {
		return c.Data[index]
	}
	panic("[BaseIndexDictModel] name not found: name: " + name)
}

func (v *IndexNameModelCorrection[T]) Contains(index Int) bool {
	_, ok := v.Data[index]
	return ok
}

func (v *IndexNameModelCorrection[T]) ContainsByName(name string) bool {
	_, ok := v.NameIndexes[name]
	return ok
}
