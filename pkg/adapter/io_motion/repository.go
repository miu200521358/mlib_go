// 指示: miu200521358
package io_motion

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion/vmd"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// VmdVpdRepository はVMD/VPDの切り替えを表す。
type VmdVpdRepository struct {
	vmdRepository *vmd.VmdRepository
	vpdRepository io_common.IFileReader
}

// NewVmdVpdRepository はVmdVpdRepositoryを生成する。
func NewVmdVpdRepository() *VmdVpdRepository {
	return &VmdVpdRepository{vmdRepository: vmd.NewVmdRepository()}
}

// SetVpdRepository はVPD読み取りリポジトリを設定する。
func (r *VmdVpdRepository) SetVpdRepository(repository io_common.IFileReader) {
	if r == nil {
		return
	}
	r.vpdRepository = repository
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *VmdVpdRepository) CanLoad(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".vmd", ".vpd":
		return true
	default:
		return false
	}
}

// InferName はパスから表示名を推定する。
func (r *VmdVpdRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load は拡張子に応じて読み込みを行う。
func (r *VmdVpdRepository) Load(path string) (hashable.IHashable, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".vpd" {
		if r.vpdRepository != nil {
			return r.vpdRepository.Load(path)
		}
		return nil, io_common.NewIoFormatNotSupported("VPD形式の読み込みは未実装です", nil)
	}
	return r.vmdRepository.Load(path)
}

// Save はVMDのみを保存する。
func (r *VmdVpdRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".vpd" {
		return io_common.NewIoEncodeFailed("VPD形式の保存は未実装です", nil)
	}
	return r.vmdRepository.Save(path, data, opts)
}
