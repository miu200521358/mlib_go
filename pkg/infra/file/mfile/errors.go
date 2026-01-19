// 指示: miu200521358
package mfile

import baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"

const (
	fileNotFoundErrorID         = "15301"
	fileReadFailedErrorID       = "15302"
	imageDecodeFailedErrorID    = "15303"
	logStreamOpenFailedErrorID  = "15305"
	consoleSnapshotSaveErrorID  = "15306"
)

// newFileNotFound はファイル未検出エラーを生成する。
func newFileNotFound(message string, cause error) error {
	return baseerr.NewCommonError(fileNotFoundErrorID, baseerr.ErrorKindValidate, message, cause)
}

// newFileReadFailed はファイル読込失敗エラーを生成する。
func newFileReadFailed(message string, cause error) error {
	return baseerr.NewCommonError(fileReadFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}

// newImageDecodeFailed は画像デコード失敗エラーを生成する。
func newImageDecodeFailed(message string, cause error) error {
	return baseerr.NewCommonError(imageDecodeFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}

// newLogStreamOpenFailed はログストリーム生成失敗エラーを生成する。
func newLogStreamOpenFailed(message string, cause error) error {
	return baseerr.NewCommonError(logStreamOpenFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}

// newConsoleSnapshotSaveFailed はコンソールスナップショット保存失敗エラーを生成する。
func newConsoleSnapshotSaveFailed(message string, cause error) error {
	return baseerr.NewCommonError(consoleSnapshotSaveErrorID, baseerr.ErrorKindValidate, message, cause)
}
