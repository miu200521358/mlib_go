// 指示: miu200521358
package mstring

import baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"

const (
	intParseFailedErrorID = "15304"
)

// newIntParseFailed は数値パース失敗エラーを生成する。
func newIntParseFailed(message string, cause error) error {
	return baseerr.NewCommonError(intParseFailedErrorID, baseerr.ErrorKindValidate, message, cause)
}
