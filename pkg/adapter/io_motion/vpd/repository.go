// 指示: miu200521358
package vpd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// VpdRepository はVPDテキスト入出力の読み取りを表す。
type VpdRepository struct{}

// NewVpdRepository はVpdRepositoryを生成する。
func NewVpdRepository() *VpdRepository {
	return &VpdRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *VpdRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".vpd")
}

// InferName はパスから表示名を推定する。
func (r *VpdRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はVPDテキストを読み込む。
func (r *VpdRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoParseFailed("VPDファイルのオープンに失敗しました", err)
	}
	defer file.Close()

	motionData := motion.NewVmdMotion(path)
	reader := newVpdReader(file)
	if err := reader.Read(motionData); err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, io_common.NewIoParseFailed("VPDファイル情報の取得に失敗しました", err)
	}
	motionData.SetFileModTime(info.ModTime().UnixNano())
	motionData.UpdateHash()
	return motionData, nil
}
