// 指示: miu200521358
package err

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

// BaseError はカスタムエラーの基底。
type BaseError struct {
	msg        string
	params     []any
	stackTrace string
	ErrorKind  merr.ErrorKind
	ErrorID    string
}

// Error はerror文字列を返す。
func (e *BaseError) Error() string {
	if e == nil {
		return ""
	}
	if e.msg == "" {
		return ""
	}
	if len(e.params) == 0 {
		return e.msg
	}
	return fmt.Sprintf(e.msg, e.params...)
}

// StackTrace はスタックトレースを返す。
func (e *BaseError) StackTrace() string {
	if e == nil {
		return ""
	}
	return e.stackTrace
}


// captureStackTrace はスタックトレースを取得する。
func captureStackTrace() string {
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	return string(bytes.ReplaceAll(buf[:n], []byte("\n"), []byte("\r\n")))
}

// TerminateErrorID はユーザー中止のErrorID。
const TerminateErrorID = "85401"

// TerminateError はユーザー中止エラー。
type TerminateError struct {
	*BaseError
	Reason string
}

// ErrorID はエラーIDを返す。
func (e *TerminateError) ErrorID() string {
	if e == nil || e.BaseError == nil {
		return ""
	}
	return e.BaseError.ErrorID
}

// ErrorKind はエラー種別を返す。
func (e *TerminateError) ErrorKind() merr.ErrorKind {
	if e == nil || e.BaseError == nil {
		return ""
	}
	return e.BaseError.ErrorKind
}

// MessageKey はメッセージキーを返す。
func (e *TerminateError) MessageKey() string {
	if e == nil || e.BaseError == nil {
		return ""
	}
	return e.BaseError.msg
}

// MessageParams はメッセージパラメータを返す。
func (e *TerminateError) MessageParams() []any {
	if e == nil || e.BaseError == nil {
		return nil
	}
	return e.BaseError.params
}

// NewTerminateError はTerminateErrorを生成する。
func NewTerminateError(reason string) *TerminateError {
	return &TerminateError{
		BaseError: &BaseError{
			msg:        "ユーザー中止で終了処理へ移行",
			params:     nil,
			stackTrace: captureStackTrace(),
			ErrorKind:  merr.ErrorKindExternal,
			ErrorID:    TerminateErrorID,
		},
		Reason: reason,
	}
}

// IsTerminateError はTerminateErrorか判定する。
func IsTerminateError(err error) bool {
	var term *TerminateError
	return errors.As(err, &term)
}
