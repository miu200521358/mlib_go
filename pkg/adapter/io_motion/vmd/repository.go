// 指示: miu200521358
package vmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// VmdRepository はVMDバイナリ入出力を表す。
type VmdRepository struct{}

// NewVmdRepository はVmdRepositoryを生成する。
func NewVmdRepository() *VmdRepository {
	return &VmdRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *VmdRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".vmd")
}

// InferName はパスから表示名を推定する。
func (r *VmdRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はVMDバイナリを読み込む。
func (r *VmdRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoParseFailed("VMDファイルのオープンに失敗しました", err)
	}
	defer file.Close()

	motionData := motion.NewVmdMotion(path)
	reader := newVmdReader(file)
	if err := reader.Read(motionData); err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, io_common.NewIoParseFailed("VMDファイル情報の取得に失敗しました", err)
	}
	motionData.SetFileModTime(info.ModTime().UnixNano())
	motionData.UpdateHash()
	return motionData, nil
}

// Save はVMDバイナリを保存する。
func (r *VmdRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		return io_common.NewIoEncodeFailed("VMD保存対象が不正です", nil)
	}
	savePath := path
	if savePath == "" {
		savePath = motionData.Path()
	}
	if savePath == "" {
		return io_common.NewIoSaveFailed("保存先パスが空です", nil)
	}

	file, err := os.Create(savePath)
	if err != nil {
		return io_common.NewIoSaveFailed("VMDファイルの作成に失敗しました", err)
	}
	defer file.Close()

	writer := newVmdWriter(file)
	if err := writer.Write(motionData); err != nil {
		return err
	}
	motionData.SetPath(savePath)
	motionData.UpdateHash()
	return nil
}
