// 指示: miu200521358
package io_model

import "github.com/miu200521358/mlib_go/pkg/adapter/io_common"

// IModelRepository はモデル入出力の共通契約を表す。
type IModelRepository interface {
	io_common.IFileRepository
}
