// 指示: miu200521358
package mmath

import "github.com/petar/GoLLRB/llrb"

type ILlrbItem[T Number] interface {
	Less(than T) bool
}

type LlrbItem[T Number] struct {
	value T
}

func NewLlrbItem[T Number](v T) *LlrbItem[T] {
	return &LlrbItem[T]{value: v}
}

func (g LlrbItem[T]) Less(than llrb.Item) bool {
	other, ok := than.(*LlrbItem[T])
	if !ok {
		return false
	}
	return g.value < other.value
}

func (g LlrbItem[T]) Value() T {
	return g.value
}

type LlrbIndexes[T Number] struct {
	*llrb.LLRB
}

func NewLlrbIndexes[T Number]() *LlrbIndexes[T] {
	return &LlrbIndexes[T]{
		LLRB: llrb.New(),
	}
}

func (li *LlrbIndexes[T]) Prev(index T) T {
	lindex := NewLlrbItem(index)
	prevValue := li.Min()
	found := false

	li.DescendLessOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(*LlrbItem[T])
		if item.value != lindex.value {
			prevValue = item.value
			found = true
			return false
		}
		return true
	})

	if !found {
		return li.Min()
	}

	return prevValue
}

func (li *LlrbIndexes[T]) Next(index T) T {
	lindex := NewLlrbItem(index)
	nextValue := index
	found := false

	li.AscendGreaterOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(*LlrbItem[T])
		if item.value != lindex.value {
			nextValue = item.value
			found = true
			return false
		}
		return true
	})

	if !found {
		return index
	}

	return nextValue
}

func (li *LlrbIndexes[T]) Has(index T) bool {
	return li.LLRB.Has(NewLlrbItem(index))
}

func (li *LlrbIndexes[T]) Max() T {
	if li.LLRB.Len() == 0 {
		return 0
	}
	return li.LLRB.Max().(*LlrbItem[T]).value
}

func (li *LlrbIndexes[T]) Min() T {
	if li.LLRB.Len() == 0 {
		return 0
	}
	return li.LLRB.Min().(*LlrbItem[T]).value
}

func (li *LlrbIndexes[T]) Length() int {
	return li.LLRB.Len()
}

func (li *LlrbIndexes[T]) ForEach(callback func(item T) bool) {
	li.LLRB.AscendGreaterOrEqual(li.LLRB.Min(), func(item llrb.Item) bool {
		index := item.(*LlrbItem[T]).Value()
		return callback(index)
	})
}

