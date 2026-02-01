// 指示: miu200521358
package mmath

import (
	"math"
	"testing"
)

func TestScalarBasics(t *testing.T) {
	values := []float64{1, 2, 3, 4}
	if Sum(values) != 10 {
		t.Errorf("Sum")
	}
	ratios := Ratio(10.0, []float64{1, 2})
	if len(ratios) != 2 || ratios[0] != 0.1 || ratios[1] != 0.2 {
		t.Errorf("Ratio")
	}
	ratiosZero := Ratio(0.0, []float64{1, 2})
	if len(ratiosZero) != 2 || ratiosZero[0] != 0 || ratiosZero[1] != 0 {
		t.Errorf("Ratio zero")
	}
	if Effective(math.NaN()) != 0 {
		t.Errorf("Effective NaN")
	}
	if Effective(1.5) != 1.5 {
		t.Errorf("Effective")
	}
	uniq := Unique([]int{1, 1, 2, 3, 2})
	if len(uniq) != 3 {
		t.Errorf("Unique")
	}
	if Mean(values) != 2.5 {
		t.Errorf("Mean")
	}
	if Mean([]float64{}) != 0 {
		t.Errorf("Mean empty")
	}
	if Median([]int{3, 1, 2}) != 2 {
		t.Errorf("Median odd")
	}
	if Median([]int{1, 2, 3, 4}) != 2 {
		t.Errorf("Median even")
	}
	if Median([]int{}) != 0 {
		t.Errorf("Median empty")
	}
	if math.Abs(Std([]float64{1, 2})-0.5) < 1e-9 {
		// 問題なし
	} else {
		t.Errorf("Std")
	}
	if Std([]float64{}) != 0 {
		t.Errorf("Std empty")
	}
	if Lerp(1, 2, -1) != 1 || Lerp(1, 2, 2) != 2 || Lerp(1, 3, 0.5) != 2 {
		t.Errorf("Lerp")
	}
	if Sign(-1) != -1 || Sign(1) != 1 {
		t.Errorf("Sign")
	}
	if !NearEquals(1.0, 1.0000001, 1e-3) || NearEquals(1.0, 1.1, 1e-3) {
		t.Errorf("NearEquals")
	}
	if DegToRad(180) != math.Pi || RadToDeg(math.Pi) != 180 {
		t.Errorf("Deg/Rad")
	}
	if ThetaToRad(1) != math.Asin(1) {
		t.Errorf("ThetaToRad")
	}
	if Clamped(5, 0, 3) != 3 || Clamped(-1, 0, 3) != 0 || Clamped(2, 0, 3) != 2 {
		t.Errorf("Clamped")
	}
	if Clamped01(-1.0) != 0 || Clamped01(2.0) != 1 || Clamped01(0.5) != 0.5 {
		t.Errorf("Clamped01")
	}
	if !Contains([]int{1, 2, 3}, 2) || Contains([]int{1, 2, 3}, 4) {
		t.Errorf("Contains small")
	}
	big := make([]int, 1001)
	for i := range big {
		big[i] = i
	}
	if !Contains(big, 1000) || Contains(big, 2000) {
		t.Errorf("Contains big")
	}
	if Max([]int{}) != 0 || Min([]int{}) != 0 {
		t.Errorf("Max/Min empty")
	}
	if Max([]int{1, 3, 2}) != 3 || Min([]int{1, 3, 2}) != 1 {
		t.Errorf("Max/Min")
	}
	if Min([]int{3, 2, 1}) != 1 {
		t.Errorf("Min update")
	}
	if len(IntRanges(3)) != 4 {
		t.Errorf("IntRanges")
	}
	if len(IntRangesByStep(1, 3, 5)) != 1 {
		t.Errorf("IntRangesByStep")
	}
	if got := Mean2DVertical([][]float64{{1, 2}, {3, 4}}); got[0] != 2 || got[1] != 3 {
		t.Errorf("Mean2DVertical")
	}
	if got := Mean2DVertical([][]float64{}); len(got) != 0 {
		t.Errorf("Mean2DVertical empty")
	}
	if got := Mean2DVertical([][]float64{{}}); len(got) != 0 {
		t.Errorf("Mean2DVertical empty row")
	}
	if got := Mean2DHorizontal([][]float64{{1, 2}, {3, 4}}); got[0] != 1.5 || got[1] != 3.5 {
		t.Errorf("Mean2DHorizontal")
	}
	if got := Mean2DHorizontal([][]float64{}); len(got) != 0 {
		t.Errorf("Mean2DHorizontal empty")
	}
	if ClampIfVerySmall(1e-7) != 0 || ClampIfVerySmall(1e-3) == 0 {
		t.Errorf("ClampIfVerySmall")
	}
	if Round(math.NaN(), 0.1) != 0 || Round(math.Inf(1), 0.1) != 0 || math.Abs(Round(1.234, 0.1)-1.2) > 1e-9 {
		t.Errorf("Round")
	}
	if !IsAllSameValues([]float64{1, 1, 1}) || IsAllSameValues([]float64{1, 2, 1}) {
		t.Errorf("IsAllSameValues")
	}
	if !IsAlmostAllSameValues([]float64{1, 1.0001, 0.9999}, 1e-2) || IsAlmostAllSameValues([]float64{1, 2}, 1e-2) {
		t.Errorf("IsAlmostAllSameValues")
	}
	orig := []int{1, 2}
	copy := DeepCopy(orig)
	copy[0] = 9
	if orig[0] == 9 {
		t.Errorf("DeepCopy")
	}
	if IsPowerOfTwo(0) || !IsPowerOfTwo(8) || IsPowerOfTwo(6) {
		t.Errorf("IsPowerOfTwo")
	}
	if BoolToInt(true) != 1 || BoolToInt(false) != 0 {
		t.Errorf("BoolToInt")
	}
	if BoolToFlag(true) != 1 || BoolToFlag(false) != -1 {
		t.Errorf("BoolToFlag")
	}
	if _, err := CalculateX(1, 1, 1); err == nil {
		t.Errorf("CalculateX error")
	}
	if v, err := CalculateX(5, 3, 4); err != nil || v != 0 {
		t.Errorf("CalculateX")
	}
	if Flatten([][]int{}) != nil {
		t.Errorf("Flatten empty")
	}
	flat := Flatten([][]int{{1, 2}, {3}})
	if len(flat) != 3 || flat[2] != 3 {
		t.Errorf("Flatten")
	}
	vals := []float64{math.NaN(), 2, 1}
	Sort(vals)
	if !math.IsNaN(vals[0]) {
		t.Errorf("Sort NaN")
	}
}

func TestDeepCopyError(t *testing.T) {
	var v any
	if _, err := deepCopy(v); err == nil {
		t.Errorf("deepCopy error")
	}
}
