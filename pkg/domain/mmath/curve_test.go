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
		{
			name: "Gimme グルーブY 0-96",
			values: []float64{
				-0.4784528315067291,
				-0.47836776438987616,
				-0.4782826972730232,
				-0.47819763015617023,
				-0.47811256303931726,
				-0.4780274959224643,
				-0.47794242880561133,
				-0.47785736168875836,
				-0.4777722945719054,
				-0.4776872274550525,
				-0.4776021603381995,
				-0.47751709322134656,
				-0.4774320261044936,
				-0.4773469589876406,
				-0.47726189187078766,
				-0.4771768247539347,
				-0.4770917576370817,
				-0.47700669052022876,
				-0.4769216234033758,
				-0.47683655628652283,
				-0.47675148916966986,
				-0.4766664220528169,
				-0.47658135493596393,
				-0.47649628781911096,
				-0.476411220702258,
				-0.4763261535854051,
				-0.4762410864685521,
				-0.47615601935169916,
				-0.4760709522348462,
				-0.4759858851179932,
				-0.47590081800114026,
				-0.4758157508842873,
				-0.4757306837674343,
				-0.47564561665058136,
				-0.4755605495337284,
				-0.47547548241687543,
				-0.47539041530002246,
				-0.4753053481831695,
				-0.47522028106631653,
				-0.47513521394946356,
				-0.4750501468326106,
				-0.47496507971575763,
				-0.4748800125989047,
				-0.47479494548205176,
				-0.4747098783651988,
				-0.4746248112483458,
				-0.47453974413149286,
				-0.4744546770146399,
				-0.4743696098977869,
				-0.47428454278093396,
				-0.474199475664081,
				-0.47411440854722803,
				-0.47402934143037506,
				-0.4739442743135221,
				-0.47385920719666913,
				-0.47377414007981616,
				-0.4736890729629632,
				-0.47360400584611023,
				-0.4735189387292573,
				-0.47343387161240436,
				-0.4733488044955514,
				-0.4732637373786984,
				-0.47317867026184546,
				-0.4730936031449925,
				-0.4730085360281395,
				-0.47292346891128656,
				-0.4728384017944336,
				-0.4728314965963364,
				-0.4728245913982391,
				-0.4728176862001419,
				-0.47281078100204466,
				-0.47280387580394745,
				-0.47279697060585024,
				-0.472790065407753,
				-0.4727831602096558,
				-0.4727762550115585,
				-0.4727693498134613,
			},
			expected: NewCurve(),
		},
	}

	for n, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewCurveFromValues(tt.values, 1e-2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[%d: %s] Expected %v, but got %v", n, tt.name, tt.expected, result)
			}
		})
	}
}
