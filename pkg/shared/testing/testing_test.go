// 指示: miu200521358
package testing

import "testing"

// TestTestingConstants は共通テスト定数を確認する。
func TestTestingConstants(t *testing.T) {
	if !AllowAbsoluteTestPaths {
		t.Errorf("AllowAbsoluteTestPaths: got=%v", AllowAbsoluteTestPaths)
	}
	if GoldenResourcePolicy != "test_resources" {
		t.Errorf("GoldenResourcePolicy: got=%v", GoldenResourcePolicy)
	}
	if MmdReproTolerance != 0.03 {
		t.Errorf("MmdReproTolerance: got=%v", MmdReproTolerance)
	}
	if DEFAULT_EPSILON_RANGE.Min != 1e-10 || DEFAULT_EPSILON_RANGE.Max != 1e-5 {
		t.Errorf("DEFAULT_EPSILON_RANGE: got=%v", DEFAULT_EPSILON_RANGE)
	}
}
