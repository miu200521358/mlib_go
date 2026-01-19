// 指示: miu200521358
package miter

import baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"

const (
	iterProcessFailedErrorID = "95301"
)

// newIterProcessFailed は並列処理失敗エラーを生成する。
func newIterProcessFailed(message string, cause error) error {
	return baseerr.NewCommonError(iterProcessFailedErrorID, baseerr.ErrorKindInternal, message, cause)
}
