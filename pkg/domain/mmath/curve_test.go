// 指示: miu200521358
package mmath

import (
	"math"
	"testing"
)

func TestCurveBasics(t *testing.T) {
	c := NewCurve()
	if c.Start != (Vec2{20, 20}) || c.End != (Vec2{107, 107}) {
		t.Errorf("NewCurve")
	}
	c2 := &Curve{Start: Vec2{-1, -1}, End: Vec2{2, 2}}
	c2.Normalize(Vec2{0, 0}, Vec2{1, 1})
	if c2.Start.X != 20 || c2.End.X != 107 {
		t.Errorf("Normalize")
	}
	c3 := &Curve{Start: Vec2{0.5, 0.5}, End: Vec2{0.5, 0.5}}
	c3.Normalize(Vec2{0, 0}, Vec2{1, 1})
	if c3.Start != (Vec2{20, 20}) || c3.End != (Vec2{107, 107}) {
		t.Errorf("Normalize linear")
	}
}

func TestCurveEvaluate(t *testing.T) {
	c := NewCurve()
	if x, y, t0 := Evaluate(c, 0, 0, 10); x != 0 || y != 0 || t0 != 0 {
		t.Errorf("Evaluate early")
	}
	if x, y, t0 := Evaluate(c, 0, 10, 0); x != 0 || y != 0 || t0 != 0 {
		t.Errorf("Evaluate denom")
	}
	if x, y, t0 := Evaluate(c, 0, 2, 1); x != 1 || y != 1 || t0 != 1 {
		t.Errorf("Evaluate x>=1")
	}
	lin := &Curve{Start: Vec2{0, 0}, End: Vec2{127, 127}}
	if x, y, t0 := Evaluate(lin, 0, 25, 100); math.Abs(x-0.25) > 1e-6 || math.Abs(y-0.25) > 1e-6 || math.Abs(t0-0.25) > 1e-6 {
		t.Errorf("Evaluate linear")
	}
	non := &Curve{Start: Vec2{20, 107}, End: Vec2{107, 20}}
	if x, y, _ := Evaluate(non, 0, 25, 100); math.Abs(y-x) < 1e-3 {
		t.Errorf("Evaluate non-linear")
	}
}

func TestCurveNewton(t *testing.T) {
	if got := newton(0, 0, 0, 0, 1e-15, 1e-20); got != 0 {
		t.Errorf("newton")
	}
}

func TestCurveSplit(t *testing.T) {
	c := NewCurve()
	c1, c2 := SplitCurve(c, 0, 0, 10)
	if c1 == nil || c2 == nil {
		t.Errorf("SplitCurve early")
	}
	c3, c4 := SplitCurve(c, 0, 5, 10)
	if c3 == nil || c4 == nil {
		t.Errorf("SplitCurve")
	}
}

func TestCurveHelpers(t *testing.T) {
	if tryCurveNormalize(Vec2{0, 0}, Vec2{0, 0}, Vec2{0, 0}, Vec2{1, 1}, false) == nil {
		t.Errorf("tryCurveNormalize linear")
	}
	if tryCurveNormalize(Vec2{0, 0}, Vec2{0, 1}, Vec2{0, 2}, Vec2{0, 3}, false) == nil {
		t.Errorf("tryCurveNormalize diffX")
	}
	if tryCurveNormalize(Vec2{0, 0}, Vec2{1, 0}, Vec2{2, 0}, Vec2{3, 0}, false) == nil {
		t.Errorf("tryCurveNormalize diffY")
	}
	if tryCurveNormalize(Vec2{0, 0}, Vec2{2, 0}, Vec2{0, 2}, Vec2{1, 1}, false) != nil {
		t.Errorf("tryCurveNormalize out")
	}
	if tryCurveNormalize(Vec2{0, 0}, Vec2{0.2, 0.8}, Vec2{0.8, 0.2}, Vec2{1, 1}, true) == nil {
		t.Errorf("tryCurveNormalize decreasing")
	}
	if !isLinearInterpolation([]float64{1, 1, 1}, 1e-6) {
		t.Errorf("isLinearInterpolation same")
	}
	if isLinearInterpolation([]float64{0, 0.2, 1}, 1e-6) {
		t.Errorf("isLinearInterpolation diff")
	}
}

func TestCurveFit(t *testing.T) {
	if c, err := NewCurveFromValues([]float64{1, 2}, 1e-3); err != nil || c == nil {
		t.Errorf("NewCurveFromValues short")
	}
	if c, err := NewCurveFromValues([]float64{1, 1, 1}, 1e-3); err != nil || c == nil {
		t.Errorf("NewCurveFromValues flat")
	}
	if c, err := NewCurveFromValues([]float64{0, 1, 2, 3}, 1e-6); err != nil || c == nil {
		t.Errorf("NewCurveFromValues linear")
	}
	if c, err := NewCurveFromValues([]float64{0, 0.2, 0.8, 1}, 1e-6); err != nil || c == nil {
		t.Errorf("NewCurveFromValues curve")
	}
	if c, err := NewCurveFromValues([]float64{3, 2, 1, 0}, 1e-6); err != nil || c == nil {
		t.Errorf("NewCurveFromValues decreasing")
	}
	if _, err := NewCurveFromValues([]float64{math.NaN(), 0, 1}, 1e-6); err == nil {
		t.Errorf("NewCurveFromValues error")
	}
	xCoords := []float64{0, 0.5, 1}
	yCoords := []float64{0, 0.2, 1}
	if _, err := optimizePoints(xCoords, yCoords, Vec2{0, 0}, Vec2{0.3, 0.1}, Vec2{0.7, 0.9}, Vec2{1, 1}); err != nil {
		t.Errorf("optimizePoints")
	}
	if _, err := optimizePoints([]float64{math.NaN(), 0.5, 1}, []float64{0, 0.2, 1}, Vec2{0, 0}, Vec2{0.3, 0.1}, Vec2{0.7, 0.9}, Vec2{1, 1}); err == nil {
		t.Errorf("optimizePoints error")
	}
	if calculateError([]float64{0, 0.5, 1}, []float64{0, 0.9, 1}, Vec2{0.2, 0.2}, Vec2{0.8, 0.8}) == 0 {
		t.Errorf("calculateError")
	}
}

func TestCurveFitNormalizeError(t *testing.T) {
	orig := optimizePointsFunc
	optimizePointsFunc = func(xCoords, yCoords []float64, P0, P1, P2, P3 Vec2) (controlPoints, error) {
		return controlPoints{P0, Vec2{X: 2, Y: 0}, Vec2{X: 0, Y: 2}, P3}, nil
	}
	defer func() { optimizePointsFunc = orig }()

	if _, err := NewCurveFromValues([]float64{0, 1, 0, 1}, 1e-6); err == nil {
		t.Errorf("NewCurveFromValues normalize error")
	}
}
