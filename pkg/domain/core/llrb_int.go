package core

import "github.com/petar/GoLLRB/llrb"

type Int int

func NewInt(v int) Int {
	return Int(v)
}

func (vInt Int) Less(than llrb.Item) bool {
	if than == nil {
		return false
	}
	return vInt < than.(Int)
}

type IntIndexes struct {
	*llrb.LLRB
}

func NewIntIndexes() *IntIndexes {
	return &IntIndexes{
		LLRB: llrb.New(),
	}
}

func (vInts IntIndexes) Prev(index int) int {
	lIndex := Int(index)

	ary := NewIntIndexes()
	vInts.DescendLessOrEqual(Int(index), func(i llrb.Item) bool {
		if i.(Int) != lIndex {
			ary.InsertNoReplace(i)
		}
		return true
	})

	if ary.Len() == 0 {
		return vInts.Min()
	}

	return ary.Max()
}

func (vInts IntIndexes) Next(index int) int {
	lIndex := Int(index)

	ary := NewIntIndexes()
	vInts.AscendGreaterOrEqual(Int(index), func(i llrb.Item) bool {
		if i.(Int) != lIndex {
			ary.InsertNoReplace(i)
		}
		return true
	})

	if ary.Len() == 0 {
		return index
	}

	return ary.Min()
}

func (vInts IntIndexes) Has(index int) bool {
	return vInts.LLRB.Has(Int(index))
}

func (vInts IntIndexes) Max() int {
	if vInts.LLRB.Len() == 0 {
		return 0
	}
	return int(vInts.LLRB.Max().(Int))
}

func (vInts IntIndexes) Min() int {
	if vInts.LLRB.Len() == 0 {
		return 0
	}
	return int(vInts.LLRB.Min().(Int))
}

func (vInts IntIndexes) Len() int {
	return vInts.LLRB.Len()
}

func (vInts IntIndexes) List() []int {
	values := make([]int, 0, vInts.LLRB.Len())
	vInts.LLRB.AscendGreaterOrEqual(vInts.LLRB.Min(), func(item llrb.Item) bool {
		if int(item.(Int)) >= 0 {
			values = append(values, int(item.(Int)))
		}
		return true
	})
	return values
}
