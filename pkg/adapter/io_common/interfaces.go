// 指示: miu200521358
package io_common

import "github.com/miu200521358/mlib_go/pkg/shared/hashable"

// IFileReader は入出力共通の読み込み契約を表す。
type IFileReader interface {
	// CanLoad は読み込み可能か判定する。
	CanLoad(path string) bool
	// Load はファイルから読み込み、IHashable を返す。
	Load(path string) (hashable.IHashable, error)
	// InferName はパスから表示名を推定する。
	InferName(path string) string
}

// IFileWriter は入出力共通の書き込み契約を表す。
type IFileWriter interface {
	// Save はIHashableを指定パスへ保存する。
	Save(path string, data hashable.IHashable, opts SaveOptions) error
}

// IFileRepository は読み書きの共通契約を表す。
type IFileRepository interface {
	IFileReader
	IFileWriter
}

// SaveOptions は保存時のオプションを表す。
type SaveOptions struct {
	// IncludeSystem はシステム要素を含めるか示す。
	IncludeSystem bool
}
