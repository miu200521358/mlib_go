package core

import "github.com/petar/GoLLRB/llrb"

type Float float32

func NewFloat(v float32) Float {
	return Float(v)
}

func (vFloat Float) Less(than llrb.Item) bool {
	if than == nil {
		return false
	}
	return vFloat < than.(Float)
}

type FloatIndexes struct {
	*llrb.LLRB
}

func NewFloatIndexes() *FloatIndexes {
	return &FloatIndexes{
		LLRB: llrb.New(),
	}
}

func (vFloats FloatIndexes) Prev(index float32) float32 {
	lIndex := Float(index)

	ary := NewFloatIndexes()
	vFloats.DescendLessOrEqual(Float(index), func(i llrb.Item) bool {
		if i.(Float) != lIndex {
			ary.InsertNoReplace(i)
		}
		return true
	})

	if ary.Len() == 0 {
		return vFloats.Min()
	}

	return ary.Max()
}

func (vFloats FloatIndexes) Next(index float32) float32 {
	lIndex := Float(index)

	ary := NewFloatIndexes()
	vFloats.AscendGreaterOrEqual(Float(index), func(i llrb.Item) bool {
		if i.(Float) != lIndex {
			ary.InsertNoReplace(i)
		}
		return true
	})

	if ary.Len() == 0 {
		return index
	}

	return ary.Min()
}

func (vFloats FloatIndexes) Has(index float32) bool {
	return vFloats.LLRB.Has(Float(index))
}

func (vFloats FloatIndexes) Max() float32 {
	if vFloats.LLRB.Len() == 0 {
		return 0
	}
	return float32(vFloats.LLRB.Max().(Float))
}

func (vFloats FloatIndexes) Min() float32 {
	if vFloats.LLRB.Len() == 0 {
		return 0
	}
	return float32(vFloats.LLRB.Min().(Float))
}

func (vFloats FloatIndexes) Len() int {
	return vFloats.LLRB.Len()
}

func (vFloats FloatIndexes) List() []float32 {
	values := make([]float32, 0, vFloats.LLRB.Len())
	vFloats.LLRB.AscendGreaterOrEqual(vFloats.LLRB.Min(), func(item llrb.Item) bool {
		if float32(item.(Float)) >= 0 {
			values = append(values, float32(item.(Float)))
		}
		return true
	})
	return values
}
