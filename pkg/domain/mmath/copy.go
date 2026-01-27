// 指示: miu200521358
package mmath

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/tiendc/go-deepcopy"
)

const mathCopyFailedErrorID = "92101"

// deepCopy は汎用の深いコピーを行う。
func deepCopy[T any](src T) (T, error) {
	var dst T
	if err := deepcopy.Copy(&dst, src); err != nil {
		return dst, merr.NewCommonError(mathCopyFailedErrorID, merr.ErrorKindInternal, "数学型のコピーに失敗", err)
	}
	return dst, nil
}
