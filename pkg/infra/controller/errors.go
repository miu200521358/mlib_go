//go:build windows
// +build windows

// 指示: miu200521358
package controller

import (
	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
)

const (
	controllerWindowInitFailedErrorID = "95501"
	consoleViewInitFailedErrorID      = "95502"
	progressBarInitFailedErrorID      = "95503"
)

// NewControllerWindowInitFailed はコントローラー初期化失敗エラーを生成する。
func NewControllerWindowInitFailed(message string, cause error) error {
	return baseerr.NewCommonError(controllerWindowInitFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}

// NewConsoleViewInitFailed はコンソール初期化失敗エラーを生成する。
func NewConsoleViewInitFailed(message string, cause error) error {
	return baseerr.NewCommonError(consoleViewInitFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}

// NewProgressBarInitFailed は進捗バー初期化失敗エラーを生成する。
func NewProgressBarInitFailed(message string, cause error) error {
	return baseerr.NewCommonError(progressBarInitFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}
