// 指示: miu200521358
package merr

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"

	"github.com/miu200521358/mlib_go/pkg/shared/errorregistry"
)

// BaseError はカスタムエラーの基底。
type BaseError struct {
	msg        string
	stackTrace string
	ErrorKind  errorregistry.ErrorKind
	ErrorID    string
}

// Error はerror文字列を返す。
func (e *BaseError) Error() string {
	if e == nil {
		return ""
	}
	return e.msg
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

// NewTerminateError はTerminateErrorを生成する。
func NewTerminateError(reason string) *TerminateError {
	return &TerminateError{
		BaseError: &BaseError{
			msg:        fmt.Sprintf("terminate error: %s", reason),
			stackTrace: captureStackTrace(),
			ErrorKind:  errorregistry.ErrorKindExternal,
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
