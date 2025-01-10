package core

import (
	"errors"
	"reflect"
	"sort"
)

type IIndexNameModel interface {
	IsValid() bool
	Index() int
	SetIndex(index int)
	Name() string
	SetName(name string)
	EnglishName() string
	SetEnglishName(englishName string)
}

// Tのリスト基底クラス
type IndexNameModels[T IIndexNameModel] struct {
	values      []T
	nameIndexes map[string]int
}

func NewIndexNameModels[T IIndexNameModel](capacity int) *IndexNameModels[T] {
	return &IndexNameModels[T]{
		values:      make([]T, capacity),
		nameIndexes: make(map[string]int, 0),
	}
}

func (im *IndexNameModels[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(im.values) {
		return *new(T), errors.New("index out of range")
	}
	return im.values[index], nil
}

func (im *IndexNameModels[T]) Update(value T) error {
	if value.Index() < 0 {
		// Update は index指定必須
		return errors.New("invalid index")
	}
	im.values[value.Index()] = value
	if _, ok := im.nameIndexes[value.Name()]; !ok {
		// 名前は先勝ち
		im.nameIndexes[value.Name()] = value.Index()
	}
	return nil
}

func (im *IndexNameModels[T]) Append(value T) error {
	if value.Index() < 0 {
		value.SetIndex(len(im.values))
	}
	im.values = append(im.values, value)
	if _, ok := im.nameIndexes[value.Name()]; !ok {
		// 名前は先勝ち
		im.nameIndexes[value.Name()] = value.Index()
	}
	return nil
}

func (im *IndexNameModels[T]) Indexes() []int {
	indexes := make([]int, len(im.nameIndexes))
	i := 0
	for _, index := range im.nameIndexes {
		indexes[i] = index
		i++
	}
	sort.Ints(indexes)
	return indexes
}

func (im *IndexNameModels[T]) Names() []string {
	names := make([]string, len(im.nameIndexes))
	i := 0
	for index := range im.Length() {
		names[i] = im.values[index].Name()
		i++
	}
	return names
}

func (im *IndexNameModels[T]) Remove(index int) error {
	if index < 0 || index >= len(im.values) {
		return errors.New("index out of range") // インデックスが範囲外の場合にエラーを返します
	}
	name := im.values[index].Name()
	delete(im.nameIndexes, name)

	im.values = append(im.values[:index], im.values[index+1:]...)
	return nil
}

func (im *IndexNameModels[T]) Length() int {
	return len(im.values)
}

func (im *IndexNameModels[T]) IsEmpty() bool {
	return len(im.values) == 0
}

func (im *IndexNameModels[T]) IsNotEmpty() bool {
	return len(im.values) > 0
}

func (im *IndexNameModels[T]) Contains(index int) bool {
	return index >= 0 && index < len(im.values) && !reflect.ValueOf(im.values[index]).IsNil()
}

func (im *IndexNameModels[T]) GetByName(name string) (T, error) {
	if index, ok := im.nameIndexes[name]; ok {
		return im.values[index], nil
	}
	return *new(T), errors.New("name not found")
}

func (im *IndexNameModels[T]) ContainsByName(name string) bool {
	_, ok := im.nameIndexes[name]
	return ok
}

func (im *IndexNameModels[T]) RemoveByName(name string) error {

	if index, ok := im.nameIndexes[name]; ok {
		im.values = append(im.values[:index], im.values[index+1:]...)
		delete(im.nameIndexes, name)
		return nil
	}
	return errors.New("name not found")
}

// Iterator はコレクションのイテレータを提供します
func (im *IndexNameModels[T]) Iterator() <-chan T {
	ch := make(chan T)
	go func() {
		for _, value := range im.values {
			ch <- value
		}
		close(ch)
	}()
	return ch
}
