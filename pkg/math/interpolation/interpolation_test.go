package interpolation_test

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/math/interpolation"
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"

)

func TestEvaluate(tst *testing.T) {
	inter := &interpolation.T{}
	inter.Start = mvec2.T{20.0, 20.0}
	inter.End = mvec2.T{107.0, 107.0}

	x, y, t := interpolation.Evaluate(inter, 0, 50, 100)

	if x != 0.5 {
		tst.Errorf("Expected x to be 0.5, but got %f", x)
	}

	if y != 0.5 {
		tst.Errorf("Expected y to be 0.5, but got %f", y)
	}

	if t != 0.5 {
		tst.Errorf("Expected t to be 0.5, but got %f", t)
	}
}

func TestEvaluate2(tst *testing.T) {
	inter := &interpolation.T{}
	inter.Start = mvec2.T{10.0, 30.0}
	inter.End = mvec2.T{100.0, 80.0}

	x, y, t := interpolation.Evaluate(inter, 0, 2, 10)

	if x != 0.2 {
		tst.Errorf("Expected x to be 0.2, but got %f", x)
	}

	expectedY := 0.24085271757748078
	if math.Abs(y-expectedY) > 1e-10 {
		tst.Errorf("Expected y to be %.20f, but got %.20f", expectedY, y)
	}

	expectedT := 0.2900272452240925
	if math.Abs(t-expectedT) > 1e-10 {
		tst.Errorf("Expected t to be %.20f, but got %.20f", expectedT, t)
	}
}