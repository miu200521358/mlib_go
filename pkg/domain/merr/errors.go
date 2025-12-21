// Package merr はドメイン層で使用するカスタムエラー型を定義します。
package merr

import (
	"errors"
	"fmt"
)

// IndexOutOfRangeError はインデックスが範囲外の場合のエラーです。
type IndexOutOfRangeError struct {
	Index int
	Max   int
}

// NewIndexOutOfRangeError は新しいIndexOutOfRangeErrorを生成します。
func NewIndexOutOfRangeError(index, max int) *IndexOutOfRangeError {
	return &IndexOutOfRangeError{
		Index: index,
		Max:   max,
	}
}

// Error はerrorインターフェースを実装します。
func (e *IndexOutOfRangeError) Error() string {
	return fmt.Sprintf("index %d out of range [0, %d]", e.Index, e.Max)
}

// IsIndexOutOfRangeError はエラーがIndexOutOfRangeErrorかどうか判定します。
func IsIndexOutOfRangeError(err error) bool {
	var target *IndexOutOfRangeError
	return errors.As(err, &target)
}

// InvalidIndexError は無効なインデックスのエラーです。
type InvalidIndexError struct {
	Index int
}

// NewInvalidIndexError は新しいInvalidIndexErrorを生成します。
func NewInvalidIndexError(index int) *InvalidIndexError {
	return &InvalidIndexError{
		Index: index,
	}
}

// Error はerrorインターフェースを実装します。
func (e *InvalidIndexError) Error() string {
	return fmt.Sprintf("invalid index: %d", e.Index)
}

// IsInvalidIndexError はエラーがInvalidIndexErrorかどうか判定します。
func IsInvalidIndexError(err error) bool {
	var target *InvalidIndexError
	return errors.As(err, &target)
}

// NameNotFoundError は指定された名前が見つからない場合のエラーです。
type NameNotFoundError struct {
	Name string
}

// NewNameNotFoundError は新しいNameNotFoundErrorを生成します。
func NewNameNotFoundError(name string) *NameNotFoundError {
	return &NameNotFoundError{
		Name: name,
	}
}

// Error はerrorインターフェースを実装します。
func (e *NameNotFoundError) Error() string {
	return fmt.Sprintf("name not found: %s", e.Name)
}

// IsNameNotFoundError はエラーがNameNotFoundErrorかどうか判定します。
func IsNameNotFoundError(err error) bool {
	var target *NameNotFoundError
	return errors.As(err, &target)
}

// InvalidArgumentError は無効な引数のエラーです。
type InvalidArgumentError struct {
	Param   string
	Message string
}

// NewInvalidArgumentError は新しいInvalidArgumentErrorを生成します。
func NewInvalidArgumentError(param, message string) *InvalidArgumentError {
	return &InvalidArgumentError{
		Param:   param,
		Message: message,
	}
}

// Error はerrorインターフェースを実装します。
func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument '%s': %s", e.Param, e.Message)
}

// IsInvalidArgumentError はエラーがInvalidArgumentErrorかどうか判定します。
func IsInvalidArgumentError(err error) bool {
	var target *InvalidArgumentError
	return errors.As(err, &target)
}

// InvalidOperationError は無効な操作のエラーです。
type InvalidOperationError struct {
	Message string
}

// NewInvalidOperationError は新しいInvalidOperationErrorを生成します。
func NewInvalidOperationError(message string) *InvalidOperationError {
	return &InvalidOperationError{
		Message: message,
	}
}

// Error はerrorインターフェースを実装します。
func (e *InvalidOperationError) Error() string {
	return fmt.Sprintf("invalid operation: %s", e.Message)
}

// IsInvalidOperationError はエラーがInvalidOperationErrorかどうか判定します。
func IsInvalidOperationError(err error) bool {
	var target *InvalidOperationError
	return errors.As(err, &target)
}
