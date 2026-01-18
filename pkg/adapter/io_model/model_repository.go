// 指示: miu200521358
package io_model

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// ModelRepository はモデル入出力のルーティングを表す。
type ModelRepository struct {
	pmxRepository *pmx.PmxRepository
	xRepository   io_common.IFileReader
}

// NewModelRepository はModelRepositoryを生成する。
func NewModelRepository() *ModelRepository {
	return &ModelRepository{
		pmxRepository: pmx.NewPmxRepository(),
	}
}

// SetXRepository はX読み取り用のリポジトリを設定する。
func (r *ModelRepository) SetXRepository(repository io_common.IFileReader) {
	if r == nil {
		return
	}
	r.xRepository = repository
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *ModelRepository) CanLoad(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pmx", ".x":
		return true
	default:
		return false
	}
}

// Load は拡張子に応じて読み込みを行う。
func (r *ModelRepository) Load(path string) (hashable.IHashable, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pmx":
		return r.pmxRepository.Load(path)
	case ".x":
		if r.xRepository != nil {
			return r.xRepository.Load(path)
		}
		return nil, io_common.NewIoFormatNotSupported("X形式の読み込みは未実装です", nil)
	default:
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
}

// InferName はパスから表示名を推定する。
func (r *ModelRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Save は拡張子に応じて保存を行う。
func (r *ModelRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pmx":
		return r.pmxRepository.Save(path, data, opts)
	case ".x":
		return io_common.NewIoEncodeFailed("X形式の保存は未実装です", nil)
	default:
		return io_common.NewIoEncodeFailed("保存形式が未対応です", nil)
	}
}
