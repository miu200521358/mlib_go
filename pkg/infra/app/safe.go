//go:build windows
// +build windows

// 指示: miu200521358
package app

import (
	"fmt"
	"runtime/debug"

	"github.com/miu200521358/mlib_go/pkg/infra/base/err"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
)

// SafeExecute はpanicを捕捉して致命ダイアログを表示する。
func SafeExecute(appConfig *config.AppConfig, f func()) {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			errMsg := ""
			if recoveredErr, ok := r.(error); ok {
				errMsg = recoveredErr.Error()
			} else {
				errMsg = fmt.Sprintf("%v", r)
			}
			runErr := fmt.Errorf("panic: %s\n%s", errMsg, string(stackTrace))
			err.ShowFatalErrorDialog(appConfig, runErr)
		}
	}()

	f()
}
