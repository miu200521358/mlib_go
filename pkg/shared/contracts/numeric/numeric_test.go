// 指示: miu200521358
package numeric

import "testing"

// TestNumericTypes は型定義を確認する。
func TestNumericTypes(t *testing.T) {
	var _ Scalar = 1.0
	var _ GpuScalar = 1.0
}

// TestNearEquals は絶対誤差の近似一致を確認する。
func TestNearEquals(t *testing.T) {
	if !NearEquals(1.0, 1.0000001, 0.001) {
		t.Errorf("NearEquals: expected true")
	}
	if NearEquals(1.0, 1.1, 0.001) {
		t.Errorf("NearEquals: expected false")
	}
}

// TestClamp は範囲丸めを確認する。
func TestClamp(t *testing.T) {
	if Clamp(5, 0, 10) != 5 {
		t.Errorf("Clamp center failed")
	}
	if Clamp(-1, 0, 10) != 0 {
		t.Errorf("Clamp min failed")
	}
	if Clamp(11, 0, 10) != 10 {
		t.Errorf("Clamp max failed")
	}
}

// TestIsFinite は有限値判定を確認する。
func TestIsFinite(t *testing.T) {
	if !IsFinite(1.0) {
		t.Errorf("IsFinite: expected true")
	}
	zero := Scalar(0)
	if IsFinite(zero / zero) {
		t.Errorf("IsFinite: expected false for NaN")
	}
}
