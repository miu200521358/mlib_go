// 指示: miu200521358
package io_common

import (
	"path/filepath"

	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
)

const (
	ioFileNotFoundErrorID     = "14101"
	ioExtInvalidErrorID       = "14102"
	ioFormatNotSupportedErrorID = "14103"
	ioEncodingUnknownErrorID  = "14104"
	ioParseFailedErrorID      = "14105"
	ioEncodeFailedErrorID     = "14106"
	ioNameEncodeFailedErrorID = "14107"
	ioSaveFailedErrorID       = "14108"
)

// NewIoFileNotFound は入力ファイル不存在エラーを生成する。
func NewIoFileNotFound(path string, cause error) error {
	message := "入力ファイルが存在しません: " + filepath.Base(path)
	return baseerr.NewCommonError(ioFileNotFoundErrorID, baseerr.ErrorKindValidate, message, cause)
}

// NewIoExtInvalid は拡張子不正エラーを生成する。
func NewIoExtInvalid(path string, cause error) error {
	message := "拡張子が不正です: " + filepath.Base(path)
	return baseerr.NewCommonError(ioExtInvalidErrorID, baseerr.ErrorKindValidate, message, cause)
}

// NewIoFormatNotSupported は形式未対応エラーを生成する。
func NewIoFormatNotSupported(message string, cause error) error {
	return baseerr.NewCommonError(ioFormatNotSupportedErrorID, baseerr.ErrorKindValidate, message, cause)
}

// NewIoEncodingUnknown は未知エンコードエラーを生成する。
func NewIoEncodingUnknown(message string, cause error) error {
	return baseerr.NewCommonError(ioEncodingUnknownErrorID, baseerr.ErrorKindValidate, message, cause)
}

// NewIoParseFailed は解析失敗エラーを生成する。
func NewIoParseFailed(message string, cause error) error {
	return baseerr.NewCommonError(ioParseFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}

// NewIoEncodeFailed はエンコード失敗エラーを生成する。
func NewIoEncodeFailed(message string, cause error) error {
	return baseerr.NewCommonError(ioEncodeFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}

// NewIoNameEncodeFailed は名前エンコード失敗エラーを生成する。
func NewIoNameEncodeFailed(message string, cause error) error {
	return baseerr.NewCommonError(ioNameEncodeFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}

// NewIoSaveFailed は保存失敗エラーを生成する。
func NewIoSaveFailed(message string, cause error) error {
	return baseerr.NewCommonError(ioSaveFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}
