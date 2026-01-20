// 指示: miu200521358
package miter

import "github.com/miu200521358/mlib_go/pkg/shared/base/merr"

const (
	iterProcessFailedErrorID = "95301"
)

// newIterProcessFailed は並列処理失敗エラーを生成する。
func newIterProcessFailed(message string, cause error) error {
	return merr.NewCommonError(iterProcessFailedErrorID, merr.ErrorKindInternal, message, cause)
}
