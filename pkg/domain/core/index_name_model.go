package core

import (
	"reflect"
	"sort"
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

func (iModels *IndexNameModels[T]) Get(index int) T {
	if index < 0 || index >= len(iModels.Data) {
		return iModels.nilFunc()
	}
	return iModels.Data[index]
}

func (iModels *IndexNameModels[T]) SetItem(index int, v T) {
	iModels.Data[index] = v
}

func (iModels *IndexNameModels[T]) Update(value T) {
	if value.Index() < 0 {
		panic("Index is not set")
	}
	iModels.Data[value.Index()] = value
	if _, ok := iModels.NameIndexes[value.Name()]; !ok {
		// 名前は先勝ち
		iModels.NameIndexes[value.Name()] = value.Index()
	}
	iModels.SetDirty(true)
}

func (iModels *IndexNameModels[T]) Append(value T) {
	if value.Index() < 0 {
		value.SetIndex(len(iModels.Data))
	}
	iModels.Data = append(iModels.Data, value)
	if _, ok := iModels.NameIndexes[value.Name()]; !ok {
		// 名前は先勝ち
		iModels.NameIndexes[value.Name()] = value.Index()
	}
	iModels.SetDirty(true)
}

func (iModels *IndexNameModels[T]) Indexes() []int {
	indexes := make([]int, len(iModels.NameIndexes))
	i := 0
	for _, index := range iModels.NameIndexes {
		indexes[i] = index
		i++
	}
	sort.Ints(indexes)
	return indexes
}

func (iModels *IndexNameModels[T]) GetNames() []string {
	names := make([]string, len(iModels.NameIndexes))
	i := 0
	for index := range iModels.Indexes() {
		names[i] = iModels.Data[index].Name()
		i++
	}
	return names
}

func (iModels *IndexNameModels[T]) DeleteItem(index int) {
	iModels.Data[index] = iModels.nilFunc()
}

func (iModels *IndexNameModels[T]) Len() int {
	return len(iModels.Data)
}

func (iModels *IndexNameModels[T]) IsEmpty() bool {
	return len(iModels.Data) == 0
}

func (iModels *IndexNameModels[T]) IsNotEmpty() bool {
	return len(iModels.Data) > 0
}

func (iModels *IndexNameModels[T]) GetByName(name string) T {
	if index, ok := iModels.NameIndexes[name]; ok {
		return iModels.Data[index]
	}
	return iModels.nilFunc()
}

func (iModels *IndexNameModels[T]) Contains(index int) bool {
	return index >= 0 && index < len(iModels.Data) && !reflect.ValueOf(iModels.Data[index]).IsNil()
}

func (iModels *IndexNameModels[T]) ContainsByName(name string) bool {
	_, ok := iModels.NameIndexes[name]
	return ok
}

func (iModels *IndexNameModels[T]) IsDirty() bool {
	return iModels.isDirty
}

func (iModels *IndexNameModels[T]) SetDirty(dirty bool) {
	iModels.isDirty = dirty
}
