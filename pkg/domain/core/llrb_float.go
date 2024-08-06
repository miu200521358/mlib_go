package core

import "github.com/petar/GoLLRB/llrb"

type Float float64

func NewFloat(v float64) Float {
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

func (vFloats FloatIndexes) Prev(index float64) float64 {
	prevIndex := Float(0)

	vFloats.DescendLessOrEqual(Float(index), func(i llrb.Item) bool {
		prevIndex = i.(Float)
		return false
	})

	return float64(prevIndex)
}

func (vFloats FloatIndexes) Next(index float64) float64 {
	nextIndex := Float(index)

	vFloats.AscendGreaterOrEqual(Float(index), func(i llrb.Item) bool {
		nextIndex = i.(Float)
		return false
	})

	return float64(nextIndex)
}

func (vFloats FloatIndexes) Has(index float64) bool {
	return vFloats.LLRB.Has(Float(index))
}

func (vFloats FloatIndexes) Max() float64 {
	if vFloats.LLRB.Len() == 0 {
		return 0
	}
	return float64(vFloats.LLRB.Max().(Float))
}

func (vFloats FloatIndexes) Min() float64 {
	if vFloats.LLRB.Len() == 0 {
		return 0
	}
	return float64(vFloats.LLRB.Min().(Float))
}

func (vFloats FloatIndexes) Len() int {
	return vFloats.LLRB.Len()
}

func (vFloats FloatIndexes) List() []float64 {
	values := make([]float64, 0, vFloats.LLRB.Len())
	vFloats.LLRB.AscendGreaterOrEqual(vFloats.LLRB.Min(), func(item llrb.Item) bool {
		if float64(item.(Float)) >= 0 {
			values = append(values, float64(item.(Float)))
		}
		return true
	})
	return values
}
