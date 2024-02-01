package mmath

import (
	"fmt"

)

type MRect struct {
	Min MVec2
	Max MVec2
}

func NewMRect() *MRect {
	return &MRect{Min: MVec2{0, 0}, Max: MVec2{0, 0}}
}

// String は MRect の文字列表現を返します。
func (v MRect) String() string {
	return fmt.Sprintf("[min=(%.5f, %.5f), max=(%.5f, %.5f)]", v.Min[0], v.Min[1], v.Max[0], v.Max[1])
}

// ContainsPoint returns if a point is contained within the rectangle.
func (rect *MRect) ContainsPoint(p *MVec2) bool {
	return p[0] >= rect.Min[0] && p[0] <= rect.Max[0] &&
		p[1] >= rect.Min[1] && p[1] <= rect.Max[1]
}

// Contains returns if other Rect is contained within the rectangle.
func (rect *MRect) Contains(other *MRect) bool {
	return rect.Min[0] <= other.Min[0] &&
		rect.Min[1] <= other.Min[1] &&
		rect.Max[0] >= other.Max[0] &&
		rect.Max[1] >= other.Max[1]
}

// Area calculates the area of the rectangle.
func (rect *MRect) Area() float64 {
	return (rect.Max[0] - rect.Min[0]) * (rect.Max[1] - rect.Min[1])
}
