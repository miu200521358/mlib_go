// 指示: miu200521358
package mstring

import "github.com/miu200521358/mlib_go/pkg/shared/base/merr"

const (
	intParseFailedErrorID = "15304"
)

// newIntParseFailed は数値パース失敗エラーを生成する。
func newIntParseFailed(message string, cause error) error {
	return merr.NewCommonError(intParseFailedErrorID, merr.ErrorKindValidate, message, cause)
}
