// 指示: miu200521358
package mfile

import "github.com/miu200521358/mlib_go/pkg/shared/base/merr"

const (
	fileNotFoundErrorID        = "15301"
	fileReadFailedErrorID      = "15302"
	imageDecodeFailedErrorID   = "15303"
	logStreamOpenFailedErrorID = "15305"
	consoleSnapshotSaveErrorID = "15306"
)

// newFileNotFound はファイル未検出エラーを生成する。
func newFileNotFound(message string, cause error, params ...any) error {
	return merr.NewCommonError(fileNotFoundErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// newFileReadFailed はファイル読込失敗エラーを生成する。
func newFileReadFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(fileReadFailedErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// newImageDecodeFailed は画像デコード失敗エラーを生成する。
func newImageDecodeFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(imageDecodeFailedErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// newLogStreamOpenFailed はログストリーム生成失敗エラーを生成する。
func newLogStreamOpenFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(logStreamOpenFailedErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// newConsoleSnapshotSaveFailed はコンソールスナップショット保存失敗エラーを生成する。
func newConsoleSnapshotSaveFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(consoleSnapshotSaveErrorID, merr.ErrorKindValidate, message, cause, params...)
}
