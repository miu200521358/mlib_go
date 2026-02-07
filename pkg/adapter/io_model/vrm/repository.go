// 指示: miu200521358
package vrm

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// VrmRepository はVRM入力の読み込み契約を表す。
type VrmRepository struct{}

// NewVrmRepository はVrmRepositoryを生成する。
func NewVrmRepository() *VrmRepository {
	return &VrmRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *VrmRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".vrm")
}

// InferName はパスから表示名を推定する。
func (r *VrmRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はVRMを読み込む。
func (r *VrmRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
	return nil, io_common.NewIoFormatNotSupported("VRM形式の読み込みは未実装です", nil)
}
