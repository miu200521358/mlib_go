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

func NewLlrbItem[T Number](v T) *LlrbItem[T] {
	return &LlrbItem[T]{value: v}
}

// Less メソッドを実装
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
	var prevValue T = li.Min()
	found := false

	// indexより小さい最大の値を探す
	li.DescendLessOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(*LlrbItem[T])
		if item.value != lindex.value {
			prevValue = item.value
			found = true
			return false // 走査を停止
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
	var nextValue T = index
	found := false

	// indexより大きい最初の値を見つけたら即座に終了
	li.AscendGreaterOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(*LlrbItem[T])
		if item.value != lindex.value {
			nextValue = item.value
			found = true
			return false // 走査を停止
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

// Iterator はコレクションのイテレータを提供します
func (li *LlrbIndexes[T]) Iterator() <-chan T {
	ch := make(chan T)
	go func() {
		li.LLRB.AscendGreaterOrEqual(li.LLRB.Min(), func(item llrb.Item) bool {
			ch <- item.(*LlrbItem[T]).Value()
			return true
		})
		close(ch)
	}()
	return ch
}
