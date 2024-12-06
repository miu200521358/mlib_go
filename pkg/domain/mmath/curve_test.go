package mmath

import (
	"math"
	"reflect"
	"testing"
)

func TestEvaluate(tst *testing.T) {
	inter := &Curve{}
	inter.Start = MVec2{20.0, 20.0}
	inter.End = MVec2{107.0, 107.0}

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
	inter.Start = MVec2{10.0, 30.0}
	inter.End = MVec2{100.0, 80.0}

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

func TestSplitCurve(t *testing.T) {
	curve := &Curve{}
	curve.Start = MVec2{89.0, 2.0}
	curve.End = MVec2{52.0, 106.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{50, 7}
	if !startCurve.Start.NearEquals(&expectedStartStart, 1e-1) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{91, 52}
	if !startCurve.End.NearEquals(&expectedStartEnd, 1e-1) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{71, 21}
	if !endCurve.Start.NearEquals(&expectedEndStart, 1e-1) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{44, 108}
	if !endCurve.End.NearEquals(&expectedEndEnd, 1e-1) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, endCurve.End)
	}
}

func TestSplitCurve2(t *testing.T) {
	curve := &Curve{}
	curve.Start = MVec2{89.0, 2.0}
	curve.End = MVec2{52.0, 106.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{50, 7}
	if !startCurve.Start.NearEquals(&expectedStartStart, 1e-1) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{91, 52}
	if !startCurve.End.NearEquals(&expectedStartEnd, 1e-1) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{71, 21}
	if !endCurve.Start.NearEquals(&expectedEndStart, 1e-1) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{44, 108}
	if !endCurve.End.NearEquals(&expectedEndEnd, 1e-1) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, endCurve.End)
	}
}

func TestSplitCurveLinear(t *testing.T) {
	curve := &Curve{}
	curve.Start = MVec2{20.0, 20.0}
	curve.End = MVec2{107.0, 107.0}

	startCurve, endCurve := SplitCurve(curve, 0, 50, 100)

	expectedStartStart := MVec2{20, 20}
	if !startCurve.Start.Equals(&expectedStartStart) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{107, 107}
	if !startCurve.End.Equals(&expectedStartEnd) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{20, 20}
	if !endCurve.Start.Equals(&expectedEndStart) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{107, 107}
	if !endCurve.End.Equals(&expectedEndEnd) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, endCurve.End)
	}
}

func TestSplitCurveSamePoints(t *testing.T) {
	curve := &Curve{}
	curve.Start = MVec2{10.0, 10.0}
	curve.End = MVec2{10.0, 10.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{20, 20}
	if !startCurve.Start.Equals(&expectedStartStart) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{107, 107}
	if !startCurve.End.Equals(&expectedStartEnd) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{20, 20}
	if !endCurve.Start.Equals(&expectedEndStart) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{107, 107}
	if !endCurve.End.Equals(&expectedEndEnd) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, endCurve.End)
	}
}

func TestSplitCurveOutOfRange(t *testing.T) {
	curve := &Curve{}
	curve.Start = MVec2{25.0, 101.0}
	curve.End = MVec2{127.0, 12.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{27, 65}
	if !startCurve.Start.Equals(&expectedStartStart) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{73, 103}
	if !startCurve.End.Equals(&expectedStartEnd) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{49, 44}
	if !endCurve.Start.Equals(&expectedEndStart) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{127, 0}
	if !endCurve.End.Equals(&expectedEndEnd) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, endCurve.End)
	}
}

func TestSplitCurveNan(t *testing.T) {
	curve := &Curve{}
	curve.Start = MVec2{127.0, 0.0}
	curve.End = MVec2{0.0, 127.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{50, 0}
	if !startCurve.Start.Equals(&expectedStartStart) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{92, 45}
	if !startCurve.End.Equals(&expectedStartEnd) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{104, 17}
	if !endCurve.Start.Equals(&expectedEndStart) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{0, 127}
	if !endCurve.End.Equals(&expectedEndEnd) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, endCurve.End)
	}
}

func TestNewCurveFromValues(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected *Curve
	}{
		{
			name:     "Empty values",
			values:   []float64{},
			expected: NewCurve(),
		},
		{
			name:     "Single value",
			values:   []float64{1.0},
			expected: NewCurve(),
		},
		{
			name:     "Two values",
			values:   []float64{1.0, 2.0},
			expected: NewCurve(),
		},
		{
			name:     "Three values",
			values:   []float64{1.0, 2.0, 3.0},
			expected: NewCurve(),
		},
		{
			name: "Gimme センターX 500-517",
			values: []float64{
				0.6999982,
				0.7076900,
				0.7293687,
				0.7631898,
				0.8075706,
				0.8611180,
				0.9225729,
				0.9907619,
				1.0645531,
				1.1428096,
				1.2243360,
				1.3078052,
				1.3916532,
				1.4739038,
				1.5518538,
				1.6214294,
				1.6756551,
				1.6999979,
			},
			expected: &Curve{
				Start: MVec2{48, 0},
				End:   MVec2{103, 127},
			},
		},
		{
			name: "Gimme センターZ 500-517",
			values: []float64{
				0.7000000,
				0.6969233,
				0.6882519,
				0.6747234,
				0.6569711,
				0.6355521,
				0.6109701,
				0.5836945,
				0.5541780,
				0.5228754,
				0.4902649,
				0.4568771,
				0.4233379,
				0.3904377,
				0.3592577,
				0.3314274,
				0.3097372,
				0.3000000,
			},
			expected: &Curve{
				Start: MVec2{48, 0},
				End:   MVec2{103, 127},
			},
		},
	}

	for n, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewCurveFromValues(tt.values)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[%d: %s] Expected %v, but got %v", n, tt.name, tt.expected, result)
			}
		})
	}
}
