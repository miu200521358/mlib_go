package mmath

import (
	"fmt"
)

type MRect struct {
	Min *MVec2
	Max *MVec2
}

func NewMRect() MRect {
	return MRect{Min: &MVec2{0, 0}, Max: &MVec2{0, 0}}
}

// String は MRect の文字列表現を返します。
func (rect MRect) String() string {
	return fmt.Sprintf("[min=(%.7f, %.7f), max=(%.7f, %.7f)]", rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
}

// ContainsPoint returns if a point is contained within the rectangle.
func (rect *MRect) ContainsPoint(p *MVec2) bool {
	return p.X >= rect.Min.X && p.X <= rect.Max.X &&
		p.Y >= rect.Min.Y && p.Y <= rect.Max.Y
}

// Contains returns if other Rect is contained within the rectangle.
func (rect *MRect) Contains(other *MRect) bool {
	return rect.Min.X <= other.Min.X &&
		rect.Min.Y <= other.Min.Y &&
		rect.Max.X >= other.Max.X &&
		rect.Max.Y >= other.Max.Y
}

// Area calculates the area of the rectangle.
func (rect *MRect) Area() float64 {
	return (rect.Max.X - rect.Min.X) * (rect.Max.Y - rect.Min.Y)
}
