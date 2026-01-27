// 指示: miu200521358
package motion

import "testing"

// TestBoneCurvesEvaluateMerge はボーン曲線の補間と合成を確認する。
func TestBoneCurvesEvaluateMerge(t *testing.T) {
	curves := NewBoneCurves()
	xy, yy, zy, ry := curves.Evaluate(0, 5, 10)
	if xy < 0.49 || xy > 0.51 || yy < 0.49 || yy > 0.51 || zy < 0.49 || zy > 0.51 || ry < 0.49 || ry > 0.51 {
		t.Fatalf("Evaluate linear: xy=%v yy=%v zy=%v ry=%v", xy, yy, zy, ry)
	}
	merged := curves.Merge(false)
	if len(merged) != 64 {
		t.Fatalf("Merge length: got=%d", len(merged))
	}
	if merged[2] != 99 || merged[3] != 15 {
		t.Fatalf("Merge physics flag: got=%v,%v", merged[2], merged[3])
	}
}

// TestBoneCurvesCopy はコピーの独立性を確認する。
func TestBoneCurvesCopy(t *testing.T) {
	curves := NewBoneCurves()
	copied, err := curves.Copy()
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}
	if len(copied.Values) != len(curves.Values) {
		t.Fatalf("Copy length mismatch")
	}
	copied.Values[0] = 99
	if curves.Values[0] == 99 {
		t.Fatalf("Copy shares Values")
	}
}

// TestCameraCurvesEvaluateMerge はカメラ曲線の補間と合成を確認する。
func TestCameraCurvesEvaluateMerge(t *testing.T) {
	curves := NewCameraCurves()
	xy, yy, zy, ry, dy, vy := curves.Evaluate(0, 5, 10)
	if xy < 0.49 || yy < 0.49 || zy < 0.49 || ry < 0.49 || dy < 0.49 || vy < 0.49 {
		t.Fatalf("Evaluate linear: xy=%v yy=%v", xy, yy)
	}
	merged := curves.Merge()
	if len(merged) != 24 {
		t.Fatalf("Merge length: got=%d", len(merged))
	}
}
