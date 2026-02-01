// 指示: miu200521358
package io

// IPathService はパス操作と保存可否判定の契約を表す。
type IPathService interface {
	// CanSave は保存可能なパスか判定する。
	CanSave(path string) bool
	// CreateOutputPath は出力パスを生成する。
	CreateOutputPath(originalPath, label string) string
	// SplitPath はパスを dir/name/ext に分割する。
	SplitPath(path string) (dir, name, ext string)
}
