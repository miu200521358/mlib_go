package mcore

import "github.com/petar/GoLLRB/llrb"

type Float32 float32

func NewFloat32(v float32) Float32 {
	return Float32(v)
}

func (v Float32) Less(than llrb.Item) bool {
	if than == nil {
		return false
	}
	return v < than.(Float32)
}

type FloatIndexes struct {
	*llrb.LLRB
}

func NewFloatIndexes() *FloatIndexes {
	return &FloatIndexes{
		LLRB: llrb.New(),
	}
}

func (i FloatIndexes) Has(index float32) bool {
	return i.LLRB.Has(Float32(index))
}

func (i FloatIndexes) Max() float32 {
	if i.LLRB.Len() == 0 {
		return 0
	}
	return float32(i.LLRB.Max().(Float32))
}
