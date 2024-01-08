package math

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/mmath"
)

func TestMVec2_CalcByRatio(t *testing.T) {
	prev := mmath.MVec2{1.0, 2.0}
	next := mmath.MVec2{3.0, 4.0}
	x := 0.5
	y := 0.5

	expected := mmath.MVec2{2.0, 3.0}
	result := prev.CalcByRatio(next, x, y)

	if result != expected {
		t.Errorf("CalcByRatio() failed, expected %v, got %v", expected, result)
	}
}
