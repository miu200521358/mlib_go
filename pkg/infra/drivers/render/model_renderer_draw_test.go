//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestSelectionRectIntersectsTriangle(t *testing.T) {
	triangle := selectionTriangle{
		{x: 10, y: 10},
		{x: 90, y: 20},
		{x: 40, y: 90},
	}
	tests := []struct {
		name    string
		min     mmath.Vec2
		max     mmath.Vec2
		expects bool
	}{
		{name: "triangle vertex inside rectangle", min: mmath.Vec2{X: 0, Y: 0}, max: mmath.Vec2{X: 20, Y: 20}, expects: true},
		{name: "partial edge overlap", min: mmath.Vec2{X: 70, Y: 0}, max: mmath.Vec2{X: 110, Y: 30}, expects: true},
		{name: "rectangle inside triangle", min: mmath.Vec2{X: 35, Y: 30}, max: mmath.Vec2{X: 45, Y: 40}, expects: true},
		{name: "fully separate", min: mmath.Vec2{X: 120, Y: 120}, max: mmath.Vec2{X: 140, Y: 140}, expects: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := selectionRectIntersectsTriangle(test.min, test.max, triangle); got != test.expects {
				t.Fatalf("selectionRectIntersectsTriangle()=%v, want %v", got, test.expects)
			}
		})
	}
}
