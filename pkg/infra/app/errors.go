// 指示: miu200521358
package app

import "github.com/miu200521358/mlib_go/pkg/shared/base/merr"

const (
	glfwInitFailedErrorID        = "95401"
	panicDetectedErrorID         = "95402"
	sharedStateInitFailedErrorID = "95403"
)

// NewGlfwInitFailed はGLFW初期化失敗エラーを生成する。
func NewGlfwInitFailed(cause error) error {
	return merr.NewCommonError(glfwInitFailedErrorID, merr.ErrorKindInternal, "GLFWの初期化に失敗", cause)
}

// NewPanicDetected はpanic検知エラーを生成する。
func NewPanicDetected(errMsg string, stackTrace string) error {
	return merr.NewCommonError(panicDetectedErrorID, merr.ErrorKindInternal, "panicを検知しました: %s\n%s", nil, errMsg, stackTrace)
}

// NewSharedStateInitFailed は共有状態初期化失敗エラーを生成する。
func NewSharedStateInitFailed() error {
	return merr.NewCommonError(sharedStateInitFailedErrorID, merr.ErrorKindInternal, "共有状態の初期化に失敗しました", nil)
}
