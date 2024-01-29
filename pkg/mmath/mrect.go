package mmath

import (
	"fmt"

	"github.com/ungerik/go3d/float64/vec2"
)

type MRect vec2.Rect

func NewMRect() *MRect {
	return &MRect{Min: vec2.T{0, 0}, Max: vec2.T{0, 0}}
}

// String は MRect の文字列表現を返します。
func (v MRect) String() string {
	return fmt.Sprintf("[min=(%.5f, %.5f), max=(%.5f, %.5f)]", v.Min[0], v.Min[1], v.Max[0], v.Max[1])
}
