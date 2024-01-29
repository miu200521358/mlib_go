package mcore

import "slices"

type IndexNameModelInterface interface {
	IsValid() bool
	Copy() IndexNameModelInterface
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
	return v.GetIndex() >= 0
}

// Copy
func (v *IndexNameModel) Copy() IndexNameModelInterface {
	copied := *v
	return &copied
}

// Tのリスト基底クラス
type IndexNameModelCorrection[T IndexNameModelInterface] struct {
	Data        map[int]T
	NameIndexes map[string]int
}

func NewIndexNameModelCorrection[T IndexNameModelInterface]() *IndexNameModelCorrection[T] {
	return &IndexNameModelCorrection[T]{
		Data:        make(map[int]T, 0),
		NameIndexes: make(map[string]int, 0),
	}
}

func (c *IndexNameModelCorrection[T]) GetItem(index int) T {
	if val, ok := c.Data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (c *IndexNameModelCorrection[T]) SetItem(index int, v T) {
	c.Data[index] = v
}

func (c *IndexNameModelCorrection[T]) Append(value T) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(c.Data))
	}
	c.Data[value.GetIndex()] = value
	c.NameIndexes[value.GetName()] = value.GetIndex()
}

func (c *IndexNameModelCorrection[T]) GetIndexes() []int {
	indexes := make([]int, 0, len(c.NameIndexes))
	for _, value := range c.NameIndexes {
		indexes = append(indexes, value)
	}
	slices.Sort(indexes)
	return indexes
}

func (c *IndexNameModelCorrection[T]) GetNames() []string {
	names := make([]string, 0, len(c.NameIndexes))
	for index := range c.GetIndexes() {
		names = append(names, c.Data[index].GetName())
	}
	return names
}

func (c *IndexNameModelCorrection[T]) DeleteItem(index int) {
	delete(c.Data, index)
}

func (c *IndexNameModelCorrection[T]) Len() int {
	return len(c.Data)
}

func (c *IndexNameModelCorrection[T]) ContainsIndex(key int) bool {
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

func (c *IndexNameModelCorrection[T]) LastIndex() int {
	maxIndex := 0
	for index := range c.Data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func (c *IndexNameModelCorrection[T]) GetSortedData() []T {
	sortedData := make([]T, 0, len(c.NameIndexes))
	for index := range c.GetIndexes() {
		sortedData = append(sortedData, c.Data[index])
	}
	return sortedData
}

func (c *IndexNameModelCorrection[T]) GetItemByName(name string) T {
	if index, ok := c.NameIndexes[name]; ok {
		return c.Data[index]
	}
	panic("[BaseIndexDictModel] name not found: name: " + name)
}

func (v *IndexNameModelCorrection[T]) Contains(index int) bool {
	_, ok := v.Data[index]
	return ok
}

func (v *IndexNameModelCorrection[T]) ContainsByName(name string) bool {
	_, ok := v.NameIndexes[name]
	return ok
}
