// 指示: miu200521358
package io_common

import (
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
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
	message := "入力ファイルが存在しない: %s"
	return merr.NewCommonError(ioFileNotFoundErrorID, merr.ErrorKindValidate, message, cause, filepath.Base(path))
}

// NewIoExtInvalid は拡張子不正エラーを生成する。
func NewIoExtInvalid(path string, cause error) error {
	message := "拡張子が不正です: %s"
	return merr.NewCommonError(ioExtInvalidErrorID, merr.ErrorKindValidate, message, cause, filepath.Base(path))
}

// NewIoFormatNotSupported は形式未対応エラーを生成する。
func NewIoFormatNotSupported(message string, cause error, params ...any) error {
	return merr.NewCommonError(ioFormatNotSupportedErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// NewIoEncodingUnknown は未知エンコードエラーを生成する。
func NewIoEncodingUnknown(message string, cause error, params ...any) error {
	return merr.NewCommonError(ioEncodingUnknownErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// NewIoParseFailed は解析失敗エラーを生成する。
func NewIoParseFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(ioParseFailedErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// NewIoEncodeFailed はエンコード失敗エラーを生成する。
func NewIoEncodeFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(ioEncodeFailedErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// NewIoNameEncodeFailed は名前エンコード失敗エラーを生成する。
func NewIoNameEncodeFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(ioNameEncodeFailedErrorID, merr.ErrorKindValidate, message, cause, params...)
}

// NewIoSaveFailed は保存失敗エラーを生成する。
func NewIoSaveFailed(message string, cause error, params ...any) error {
	return merr.NewCommonError(ioSaveFailedErrorID, merr.ErrorKindValidate, message, cause, params...)
}
