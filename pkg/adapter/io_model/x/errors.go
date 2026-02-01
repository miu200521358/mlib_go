// 指示: miu200521358
package x

import "github.com/miu200521358/mlib_go/pkg/adapter/io_common"

// newParseFailed はX読み込みの解析失敗エラーを生成する。
func newParseFailed(message string, params ...any) error {
	return io_common.NewIoParseFailed(message, nil, params...)
}
