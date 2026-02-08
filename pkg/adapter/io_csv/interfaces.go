// 指示: miu200521358
package io_csv

import "github.com/miu200521358/mlib_go/pkg/adapter/io_common"

// ICsvRepository はCSV入出力の共通契約を表す。
type ICsvRepository interface {
	io_common.IFileRepository
}
