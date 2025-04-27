package app

import (
	"fmt"
	"runtime/debug"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/merr"
)

// SafeExecute は関数でpanicが発生した場合にダイアログを表示する
func SafeExecute(appConfig *mconfig.AppConfig, f func()) {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			var errMsg string
			if recoveredErr, ok := r.(error); ok {
				errMsg = recoveredErr.Error()
			} else {
				errMsg = fmt.Sprintf("%v", r)
			}
			err := fmt.Errorf("panic: %s\n%s", errMsg, stackTrace)
			merr.ShowFatalErrorDialog(appConfig, err)
		}
	}()

	f()
}
