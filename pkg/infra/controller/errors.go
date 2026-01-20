//go:build windows
// +build windows

// 指示: miu200521358
package controller

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

const (
	controllerWindowInitFailedErrorID = "95501"
	consoleViewInitFailedErrorID      = "95502"
	progressBarInitFailedErrorID      = "95503"
)

// NewControllerWindowInitFailed はコントローラー初期化失敗エラーを生成する。
func NewControllerWindowInitFailed(message string, cause error) error {
	return merr.NewCommonError(controllerWindowInitFailedErrorID, merr.ErrorKindInternal, message, cause)
}

// NewConsoleViewInitFailed はコンソール初期化失敗エラーを生成する。
func NewConsoleViewInitFailed(message string, cause error) error {
	return merr.NewCommonError(consoleViewInitFailedErrorID, merr.ErrorKindInternal, message, cause)
}

// NewProgressBarInitFailed は進捗バー初期化失敗エラーを生成する。
func NewProgressBarInitFailed(message string, cause error) error {
	return merr.NewCommonError(progressBarInitFailedErrorID, merr.ErrorKindInternal, message, cause)
}
