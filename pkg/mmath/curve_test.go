package mmath

import (
	"math"
	"testing"
)

func TestEvaluate(tst *testing.T) {
	inter := &Curve{}
	inter.Start = &MVec2{20.0, 20.0}
	inter.End = &MVec2{107.0, 107.0}

	x, y, t := Evaluate(inter, 0, 50, 100)

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
	inter := &Curve{}
	inter.Start = &MVec2{10.0, 30.0}
	inter.End = &MVec2{100.0, 80.0}

	x, y, t := Evaluate(inter, 0, 2, 10)

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

func TestEvaluate3(tst *testing.T) {
	inter := &Curve{}
	inter.Start = &MVec2{46.0, 33.0}
	inter.End = &MVec2{80.0, 88.0}

	x, y, t := Evaluate(inter, 272, 278, 280)

	if x != 0.75 {
		tst.Errorf("Expected x to be 0.2, but got %f", x)
	}

	expectedY := 0.7628337691389842
	if math.Abs(y-expectedY) > 1e-10 {
		tst.Errorf("Expected y to be %.20f, but got %.20f", expectedY, y)
	}

	expectedT := 0.7614940998277816
	if math.Abs(t-expectedT) > 1e-10 {
		tst.Errorf("Expected t to be %.20f, but got %.20f", expectedT, t)
	}
}

func TestSplitCurve(t *testing.T) {
	curve := &Curve{}
	curve.Start = &MVec2{89.0, 2.0}
	curve.End = &MVec2{52.0, 106.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{50, 7}
	if !startCurve.Start.PracticallyEquals(&expectedStartStart, 1e-1) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, *startCurve.Start)
	}

	expectedStartEnd := MVec2{91, 52}
	if !startCurve.End.PracticallyEquals(&expectedStartEnd, 1e-1) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, *startCurve.End)
	}

	expectedEndStart := MVec2{71, 21}
	if !endCurve.Start.PracticallyEquals(&expectedEndStart, 1e-1) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, *endCurve.Start)
	}

	expectedEndEnd := MVec2{44, 108}
	if !endCurve.End.PracticallyEquals(&expectedEndEnd, 1e-1) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, *endCurve.End)
	}
}

func TestSplitCurve2(t *testing.T) {
	curve := &Curve{}
	curve.Start = &MVec2{89.0, 2.0}
	curve.End = &MVec2{52.0, 106.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{50, 7}
	if !startCurve.Start.PracticallyEquals(&expectedStartStart, 1e-1) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, *startCurve.Start)
	}

	expectedStartEnd := MVec2{91, 52}
	if !startCurve.End.PracticallyEquals(&expectedStartEnd, 1e-1) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, *startCurve.End)
	}

	expectedEndStart := MVec2{71, 21}
	if !endCurve.Start.PracticallyEquals(&expectedEndStart, 1e-1) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, *endCurve.Start)
	}

	expectedEndEnd := MVec2{44, 108}
	if !endCurve.End.PracticallyEquals(&expectedEndEnd, 1e-1) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, *endCurve.End)
	}
}

func TestSplitCurveLinear(t *testing.T) {
	curve := &Curve{}
	curve.Start = &MVec2{20.0, 20.0}
	curve.End = &MVec2{107.0, 107.0}

	startCurve, endCurve := SplitCurve(curve, 0, 50, 100)

	expectedStartStart := MVec2{20, 20}
	if !startCurve.Start.Equals(&expectedStartStart) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, *startCurve.Start)
	}

	expectedStartEnd := MVec2{107, 107}
	if !startCurve.End.Equals(&expectedStartEnd) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, *startCurve.End)
	}

	expectedEndStart := MVec2{20, 20}
	if !endCurve.Start.Equals(&expectedEndStart) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, *endCurve.Start)
	}

	expectedEndEnd := MVec2{107, 107}
	if !endCurve.End.Equals(&expectedEndEnd) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, *endCurve.End)
	}
}

func TestSplitCurveSamePoints(t *testing.T) {
	curve := &Curve{}
	curve.Start = &MVec2{10.0, 10.0}
	curve.End = &MVec2{10.0, 10.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{20, 20}
	if !startCurve.Start.Equals(&expectedStartStart) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, *startCurve.Start)
	}

	expectedStartEnd := MVec2{107, 107}
	if !startCurve.End.Equals(&expectedStartEnd) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, *startCurve.End)
	}

	expectedEndStart := MVec2{20, 20}
	if !endCurve.Start.Equals(&expectedEndStart) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, *endCurve.Start)
	}

	expectedEndEnd := MVec2{107, 107}
	if !endCurve.End.Equals(&expectedEndEnd) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, *endCurve.End)
	}
}

func TestSplitCurveOutOfRange(t *testing.T) {
	curve := &Curve{}
	curve.Start = &MVec2{25.0, 101.0}
	curve.End = &MVec2{127.0, 12.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{27, 65}
	if !startCurve.Start.Equals(&expectedStartStart) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, *startCurve.Start)
	}

	expectedStartEnd := MVec2{73, 103}
	if !startCurve.End.Equals(&expectedStartEnd) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, *startCurve.End)
	}

	expectedEndStart := MVec2{49, 44}
	if !endCurve.Start.Equals(&expectedEndStart) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, *endCurve.Start)
	}

	expectedEndEnd := MVec2{127, 0}
	if !endCurve.End.Equals(&expectedEndEnd) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, *endCurve.End)
	}
}
