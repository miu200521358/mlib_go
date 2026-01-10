// 指示: miu200521358
package err

import "testing"

// TestTerminateError は基本情報を確認する。
func TestTerminateError(t *testing.T) {
	err := NewTerminateError("stop")
	if err.ErrorID != TerminateErrorID {
		t.Errorf("ErrorID: got=%v", err.ErrorID)
	}
	if err.Error() == "" {
		t.Errorf("Error message empty")
	}
	if !IsTerminateError(err) {
		t.Errorf("IsTerminateError failed")
	}
}
