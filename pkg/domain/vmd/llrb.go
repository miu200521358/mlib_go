package vmd

import (
	"github.com/petar/GoLLRB/llrb"
)

type LlrbFrameIndexes struct {
	*llrb.LLRB
}

func NewLlrbFrameIndexes() *LlrbFrameIndexes {
	return &LlrbFrameIndexes{
		LLRB: llrb.New(),
	}
}

func (li *LlrbFrameIndexes) Prev(index float32) float32 {
	lindex := NewFrame(index)
	var prevValue float32 = li.Min()
	found := false

	// indexより小さい最大の値を探す
	li.DescendLessOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(IBaseFrame)
		if item.Index() != lindex.Index() {
			prevValue = item.Index()
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

func (li *LlrbFrameIndexes) Next(index float32) float32 {
	lindex := NewFrame(index)
	var nextValue float32 = index
	found := false

	// indexより大きい最初の値を見つけたら即座に終了
	li.AscendGreaterOrEqual(lindex, func(i llrb.Item) bool {
		item := i.(IBaseFrame)
		if item.Index() != lindex.Index() {
			nextValue = item.Index()
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

func (li *LlrbFrameIndexes) Has(index float32) bool {
	return li.LLRB.Has(NewFrame(index))
}

func (li *LlrbFrameIndexes) Max() float32 {
	if li.LLRB.Len() == 0 {
		return 0
	}
	return li.LLRB.Max().(IBaseFrame).Index()
}

func (li *LlrbFrameIndexes) Min() float32 {
	if li.LLRB.Len() == 0 {
		return 0
	}
	return li.LLRB.Min().(IBaseFrame).Index()
}

func (li *LlrbFrameIndexes) Length() int {
	return li.LLRB.Len()
}

// ForEach はコレクションのイテレータを提供します
func (li *LlrbFrameIndexes) ForEach(callback func(item float32) bool) {
	li.LLRB.AscendGreaterOrEqual(li.LLRB.Min(), func(item llrb.Item) bool {
		index := item.(IBaseFrame).Index()
		return callback(index)
	})
}
