// 指示: miu200521358
package err

import "testing"

// TestTerminateError は基本情報を確認する。
func TestTerminateError(t *testing.T) {
	err := NewTerminateError("stop")
	if err.ErrorID() != TerminateErrorID {
		t.Errorf("ErrorID: got=%v", err.ErrorID())
	}
	if err.Error() == "" {
		t.Errorf("Error message empty")
	}
	if !IsTerminateError(err) {
		t.Errorf("IsTerminateError failed")
	}
}

// TestBaseErrorNil はnil受信の挙動を確認する。
func TestBaseErrorNil(t *testing.T) {
	var be *BaseError
	if be.Error() != "" {
		t.Errorf("nil Error should be empty")
	}
	if be.StackTrace() != "" {
		t.Errorf("nil StackTrace should be empty")
	}
	if IsTerminateError(be) {
		t.Errorf("IsTerminateError should be false for nil")
	}
}

// TestBaseErrorValues は値の取得を確認する。
func TestBaseErrorValues(t *testing.T) {
	be := &BaseError{
		msg:        "msg",
		stackTrace: "stack",
		ErrorKind:  "",
		ErrorID:    "id",
	}
	if be.Error() != "msg" {
		t.Errorf("Error: got=%v", be.Error())
	}
	if be.StackTrace() != "stack" {
		t.Errorf("StackTrace: got=%v", be.StackTrace())
	}
	if IsTerminateError(be) {
		t.Errorf("IsTerminateError should be false")
	}
}
