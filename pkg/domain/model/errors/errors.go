// 指示: miu200521358
package errors

import (
	stdErrors "errors"
	"fmt"
)

// IndexOutOfRangeError はインデックス範囲外エラーを表す。
type IndexOutOfRangeError struct {
	Index  int
	Length int
}

// NewIndexOutOfRangeError は IndexOutOfRangeError を生成する。
func NewIndexOutOfRangeError(index, length int) *IndexOutOfRangeError {
	return &IndexOutOfRangeError{Index: index, Length: length}
}

// Error はエラーメッセージを返す。
func (e *IndexOutOfRangeError) Error() string {
	return fmt.Sprintf("index out of range: index=%d length=%d", e.Index, e.Length)
}

// IsIndexOutOfRangeError は err が IndexOutOfRangeError か判定する。
func IsIndexOutOfRangeError(err error) bool {
	var target *IndexOutOfRangeError
	return stdErrors.As(err, &target)
}

// NameNotFoundError は名前未検出エラーを表す。
type NameNotFoundError struct {
	Name string
}

// NewNameNotFoundError は NameNotFoundError を生成する。
func NewNameNotFoundError(name string) *NameNotFoundError {
	return &NameNotFoundError{Name: name}
}

// Error はエラーメッセージを返す。
func (e *NameNotFoundError) Error() string {
	return fmt.Sprintf("name not found: %s", e.Name)
}

// IsNameNotFoundError は err が NameNotFoundError か判定する。
func IsNameNotFoundError(err error) bool {
	var target *NameNotFoundError
	return stdErrors.As(err, &target)
}

// NameConflictError は名前衝突エラーを表す。
type NameConflictError struct {
	Name string
}

// NewNameConflictError は NameConflictError を生成する。
func NewNameConflictError(name string) *NameConflictError {
	return &NameConflictError{Name: name}
}

// Error はエラーメッセージを返す。
func (e *NameConflictError) Error() string {
	return fmt.Sprintf("name conflict: %s", e.Name)
}

// IsNameConflictError は err が NameConflictError か判定する。
func IsNameConflictError(err error) bool {
	var target *NameConflictError
	return stdErrors.As(err, &target)
}

// NameMismatchError は更新時の名前不一致エラーを表す。
type NameMismatchError struct {
	Index    int
	Expected string
	Actual   string
}

// NewNameMismatchError は NameMismatchError を生成する。
func NewNameMismatchError(index int, expected, actual string) *NameMismatchError {
	return &NameMismatchError{Index: index, Expected: expected, Actual: actual}
}

// Error はエラーメッセージを返す。
func (e *NameMismatchError) Error() string {
	return fmt.Sprintf("name mismatch: index=%d expected=%s actual=%s", e.Index, e.Expected, e.Actual)
}

// IsNameMismatchError は err が NameMismatchError か判定する。
func IsNameMismatchError(err error) bool {
	var target *NameMismatchError
	return stdErrors.As(err, &target)
}

// ModelCopyFailed はモデルコピー失敗を表す。
type ModelCopyFailed struct {
	Cause error
}

// NewModelCopyFailed は ModelCopyFailed を生成する。
func NewModelCopyFailed(cause error) *ModelCopyFailed {
	return &ModelCopyFailed{Cause: cause}
}

// Error はエラーメッセージを返す。
func (e *ModelCopyFailed) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("model copy failed: %v", e.Cause)
}

// Unwrap は元の原因エラーを返す。
func (e *ModelCopyFailed) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// IsModelCopyFailed は err が ModelCopyFailed か判定する。
func IsModelCopyFailed(err error) bool {
	var target *ModelCopyFailed
	return stdErrors.As(err, &target)
}
