package index_model

import "sort"

// Tのインタフェース
type TInterface interface {
	GetIndex() int
	SetIndex(index int)
	IsValid() bool
	Add(v *TInterface) *T
	Copy() *T
}

// INDEXを持つ基底クラス
type T struct {
	Index int
}

func (v *T) GetIndex() int {
	return v.Index
}

func (v *T) SetIndex(index int) {
	v.Index = index
}

func NewT(index int) *T {
	return &T{
		Index: index,
	}
}

func (b *T) IsValid() bool {
	return b.GetIndex() >= 0
}

func (b *T) Add(v *TInterface) *T {
	// Implement your logic here
	return nil
}

// Copy
func (b *T) Copy() *T {
	copied := *b
	return &copied
}

// Cのインタフェース
type CInterface interface {
	GetItem(index int) *T
	Range(start, stop, step int) []*T
	SetItem(index int, v TInterface)
	Append(value TInterface, isSort bool)
}

// Tのリスト基底クラス
type C struct {
	data    map[int]*T
	Indexes []int
}

func NewC() *C {
	return &C{
		data:    make(map[int]*T),
		Indexes: make([]int, 0),
	}
}

func (b *C) GetItem(index int) *T {
	if index < 0 {
		// マイナス指定の場合、後ろからの順番に置き換える
		index = len(b.data) + index
		return b.data[b.Indexes[index]]
	}
	if val, ok := b.data[index]; ok {
		return val
	}
	// なかったらエラー
	panic("[BaseIndexDictModel] index out of range: index: " + string(rune(index)))
}

func (b *C) Range(start, stop, step int) []*T {
	if stop < 0 {
		// マイナス指定の場合、後ろからの順番に置き換える
		stop = len(b.data) + stop + 1
	}
	result := make([]*T, 0)
	for i := start; i < stop; i += step {
		result = append(result, b.data[b.Indexes[i]])
	}
	return result
}

func (b *C) SetItem(index int, v TInterface) {
	b.data[index] = v.(*T)
}

func (b *C) Append(value TInterface, isSort bool) {
	if value.GetIndex() < 0 {
		value.SetIndex(len(b.data))
	}
	b.data[value.GetIndex()] = value.(*T)
	if isSort {
		b.SortIndexes()
	} else {
		b.Indexes = append(b.Indexes, value.GetIndex())
	}
}

func (b *C) SortIndexes() {
	sort.Ints(b.Indexes)
}

func (b *C) DeleteItem(index int) {
	delete(b.data, index)
}

func (b *C) Len() int {
	return len(b.data)
}

func (b *C) Iterator() []*T {
	result := make([]*T, 0)
	for _, index := range b.Indexes {
		result = append(result, b.data[index])
	}
	return result
}

func (b *C) Contains(key int) bool {
	_, ok := b.data[key]
	return ok
}

func (b *C) IsEmpty() bool {
	return len(b.data) > 0
}

func (b *C) LastIndex() int {
	maxIndex := 0
	for index := range b.data {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}
