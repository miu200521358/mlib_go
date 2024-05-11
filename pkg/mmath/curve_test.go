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
	if !startCurve.Start.PracticallyEquals(&expectedStartStart, 1e-1) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{91, 52}
	if !startCurve.End.PracticallyEquals(&expectedStartEnd, 1e-1) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{71, 21}
	if !endCurve.Start.PracticallyEquals(&expectedEndStart, 1e-1) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{44, 108}
	if !endCurve.End.PracticallyEquals(&expectedEndEnd, 1e-1) {
		t.Errorf("Expected endCurve.End to be %v, but got %v", expectedEndEnd, endCurve.End)
	}
}

func TestSplitCurve2(t *testing.T) {
	curve := &Curve{}
	curve.Start = MVec2{89.0, 2.0}
	curve.End = MVec2{52.0, 106.0}

	startCurve, endCurve := SplitCurve(curve, 0, 2, 10)

	expectedStartStart := MVec2{50, 7}
	if !startCurve.Start.PracticallyEquals(&expectedStartStart, 1e-1) {
		t.Errorf("Expected startCurve.Start to be %v, but got %v", expectedStartStart, startCurve.Start)
	}

	expectedStartEnd := MVec2{91, 52}
	if !startCurve.End.PracticallyEquals(&expectedStartEnd, 1e-1) {
		t.Errorf("Expected startCurve.End to be %v, but got %v", expectedStartEnd, startCurve.End)
	}

	expectedEndStart := MVec2{71, 21}
	if !endCurve.Start.PracticallyEquals(&expectedEndStart, 1e-1) {
		t.Errorf("Expected endCurve.Start to be %v, but got %v", expectedEndStart, endCurve.Start)
	}

	expectedEndEnd := MVec2{44, 108}
	if !endCurve.End.PracticallyEquals(&expectedEndEnd, 1e-1) {
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

func TestNewCurveFromValues(t *testing.T) {
	// Test case 1: Empty values
	values := []float64{}
	expected := NewCurve()
	result := NewCurveFromValues(values)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 2: Single value
	values = []float64{1.0}
	expected = NewCurve()
	result = NewCurveFromValues(values)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 3: Two values
	values = []float64{1.0, 2.0}
	expected = NewCurve()
	result = NewCurveFromValues(values)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 4: Three values
	values = []float64{1.0, 2.0, 3.0}
	expected = NewCurve()
	result = NewCurveFromValues(values)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 5: Four values
	values = []float64{
		0.5979851484298706,
		0.521004855632782,
		0.37936708331108093,
		0.3069007694721222,
		0.30701854825019836,
		0.243390291929245,
		0.28262293338775635,
		0.2678329348564148,
		0.27642160654067993,
		0.276456743478775,
		0.3345196545124054,
		0.3592238426208496,
		0.34113433957099915,
		0.32507243752479553,
		0.32527124881744385,
		0.28771495819091797,
		0.16925211250782013,
		0.07562929391860962,
		-0.01844712160527706,
		-0.01857336051762104,
		-0.10410086810588837,
		-0.1692115068435669,
		-0.22358137369155884,
		-0.28167492151260376,
		-0.2816798686981201,
		-0.2936137318611145,
		-0.35715776681900024,
		-0.3935558795928955,
		-0.43389225006103516,
		-0.4339143633842468,
		-0.4917398691177368,
		-0.5314293503761292,
		-0.5342316031455994,
		-0.5640663504600525,
		-0.5642469525337219,
		-0.5541977882385254,
		-0.5907184481620789,
		-0.6129558682441711,
		-0.7112897038459778,
		-0.7110443711280823,
		-0.7198588848114014,
		-0.7699093818664551,
		-0.8582518100738525,
		-0.9080440998077393,
		-0.907961368560791,
		-0.9225611090660095,
		-1.009341835975647,
		-1.0600273609161377,
		-1.1095038652420044,
		-1.1092925071716309,
		-1.1037554740905762,
		-1.1115431785583496,
	}
	expected = &Curve{
		Start: MVec2{2, 37},
		End:   MVec2{125, 94},
	}
	result = NewCurveFromValues(values)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
