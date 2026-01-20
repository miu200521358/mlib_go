// 指示: miu200521358
package merrors

import (
	"errors"
	"testing"
)

func TestIndexOutOfRangeError(t *testing.T) {
	err := NewIndexOutOfRangeError(2, 5)
	if err.Index != 2 || err.Length != 5 {
		t.Fatalf("unexpected fields: %+v", err)
	}
	if err.Error() == "" {
		t.Fatalf("Error should not be empty")
	}
	if !IsIndexOutOfRangeError(err) {
		t.Fatalf("IsIndexOutOfRangeError should be true")
	}
}

func TestNameErrors(t *testing.T) {
	missing := NewNameNotFoundError("m")
	if missing.Error() == "" || !IsNameNotFoundError(missing) {
		t.Fatalf("NameNotFoundError check failed")
	}

	conflict := NewNameConflictError("c")
	if conflict.Error() == "" || !IsNameConflictError(conflict) {
		t.Fatalf("NameConflictError check failed")
	}

	mismatch := NewNameMismatchError(1, "a", "b")
	if mismatch.Error() == "" || !IsNameMismatchError(mismatch) {
		t.Fatalf("NameMismatchError check failed")
	}
}

func TestModelCopyFailed(t *testing.T) {
	cause := errors.New("boom")
	err := NewModelCopyFailed(cause)
	if err.Error() == "" {
		t.Fatalf("Error should not be empty")
	}
	if err.Unwrap() != cause {
		t.Fatalf("Unwrap should return cause")
	}
	if !IsModelCopyFailed(err) {
		t.Fatalf("IsModelCopyFailed should be true")
	}

	var nilErr *ModelCopyFailed
	if nilErr.Error() != "" {
		t.Fatalf("nil Error should be empty")
	}
	if nilErr.Unwrap() != nil {
		t.Fatalf("nil Unwrap should be nil")
	}
}
