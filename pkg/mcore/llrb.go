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
