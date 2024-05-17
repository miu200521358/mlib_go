package mcore

import "github.com/petar/GoLLRB/llrb"

type Int int

func NewInt(v int) Int {
	return Int(v)
}

func (v Int) Less(than llrb.Item) bool {
	if than == nil {
		return false
	}
	return v < than.(Int)
}

type IntIndexes struct {
	*llrb.LLRB
}

func NewIntIndexes() *IntIndexes {
	return &IntIndexes{
		LLRB: llrb.New(),
	}
}

func (i IntIndexes) Prev(index int) int {
	prevIndex := Int(0)

	i.DescendLessOrEqual(Int(index), func(i llrb.Item) bool {
		prevIndex = i.(Int)
		return false
	})

	return int(prevIndex)
}

func (i IntIndexes) Next(index int) int {
	nextIndex := Int(index)

	i.AscendGreaterOrEqual(Int(index), func(i llrb.Item) bool {
		nextIndex = i.(Int)
		return false
	})

	return int(nextIndex)
}

func (i IntIndexes) Has(index int) bool {
	return i.LLRB.Has(Int(index))
}

func (i IntIndexes) Max() int {
	if i.LLRB.Len() == 0 {
		return 0
	}
	return int(i.LLRB.Max().(Int))
}

func (i IntIndexes) Min() int {
	if i.LLRB.Len() == 0 {
		return 0
	}
	return int(i.LLRB.Min().(Int))
}

func (i IntIndexes) Len() int {
	return i.LLRB.Len()
}

func (i IntIndexes) List() []int {
	list := make([]int, 0, i.LLRB.Len())
	i.LLRB.AscendGreaterOrEqual(i.LLRB.Min(), func(item llrb.Item) bool {
		if int(item.(Int)) >= 0 {
			list = append(list, int(item.(Int)))
		}
		return true
	})
	return list
}
