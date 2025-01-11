package mmath

import (
	"github.com/petar/GoLLRB/llrb"
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

// Less メソッドを持つ ILlrbItem インターフェースを定義
type ILlrbItem[T Number] interface {
	Less(than T) bool
}

// ---------------------------------------------------------

// 汎用型の定義
type LlrbItem[T Number] struct {
	value T
}

func NewLlrbItem[T Number](v T) LlrbItem[T] {
	return LlrbItem[T]{value: v}
}

// Less メソッドを実装
func (g LlrbItem[T]) Less(than llrb.Item) bool {
	other, ok := than.(LlrbItem[T])
	if !ok {
		return false
	}
	return g.value < other.value
}

// ---------------------------------------------------------

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

	ary := NewLlrbIndexes[T]()
	li.DescendLessOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(LlrbItem[T])
		if item.value != lindex.value {
			ary.InsertNoReplace(item)
		}
		return true
	})

	if ary.Len() == 0 {
		return li.Min()
	}

	return ary.Max()
}

func (li *LlrbIndexes[T]) Next(index T) T {
	lindex := NewLlrbItem(index)

	ary := NewLlrbIndexes[T]()
	li.AscendGreaterOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(LlrbItem[T])
		if item.value != lindex.value {
			ary.InsertNoReplace(item)
		}
		return true
	})

	if ary.Len() == 0 {
		return index
	}

	return ary.Min()
}

func (li *LlrbIndexes[T]) Has(index T) bool {
	return li.LLRB.Has(NewLlrbItem(index))
}

func (li *LlrbIndexes[T]) Max() T {
	if li.LLRB.Len() == 0 {
		return 0
	}
	return li.LLRB.Max().(LlrbItem[T]).value
}

func (li *LlrbIndexes[T]) Min() T {
	if li.LLRB.Len() == 0 {
		return 0
	}
	return li.LLRB.Min().(LlrbItem[T]).value
}

func (li *LlrbIndexes[T]) Length() int {
	return li.LLRB.Len()
}

func (li *LlrbIndexes[T]) Iter(itemFunc func(item llrb.Item) bool) {
	li.LLRB.AscendGreaterOrEqual(li.LLRB.Min(), func(item llrb.Item) bool {
		return itemFunc(item)
	})
}
