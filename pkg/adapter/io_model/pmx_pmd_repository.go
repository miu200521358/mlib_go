// 指示: miu200521358
package io_model

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmd"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// PmxPmdRepository はPMX/PMD入出力のルーティングを表す。
type PmxPmdRepository struct {
	pmxRepository *pmx.PmxRepository
	pmdRepository *pmd.PmdRepository
}

// NewPmxPmdRepository はPmxPmdRepositoryを生成する。
func NewPmxPmdRepository() *PmxPmdRepository {
	return &PmxPmdRepository{
		pmxRepository: pmx.NewPmxRepository(),
		pmdRepository: pmd.NewPmdRepository(),
	}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *PmxPmdRepository) CanLoad(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pmx", ".pmd":
		return true
	default:
		return false
	}
}

// Load は拡張子に応じて読み込みを行う。
func (r *PmxPmdRepository) Load(path string) (hashable.IHashable, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pmx":
		return r.pmxRepository.Load(path)
	case ".pmd":
		return r.pmdRepository.Load(path)
	default:
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
}

// InferName はパスから表示名を推定する。
func (r *PmxPmdRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Save は拡張子に応じて保存を行う。
func (r *PmxPmdRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pmx":
		return r.pmxRepository.Save(path, data, opts)
	case ".pmd":
		return r.pmdRepository.Save(path, data, opts)
	default:
		return io_common.NewIoEncodeFailed("保存形式が未対応です", nil)
	}
}
