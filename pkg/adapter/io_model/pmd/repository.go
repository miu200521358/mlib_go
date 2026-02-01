// 指示: miu200521358
package pmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// PmdRepository はPMDバイナリ入出力を表す。
type PmdRepository struct{}

// NewPmdRepository はPmdRepositoryを生成する。
func NewPmdRepository() *PmdRepository {
	return &PmdRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *PmdRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".pmd")
}

// InferName はパスから表示名を推定する。
func (r *PmdRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はPMDバイナリを読み込む。
func (r *PmdRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoParseFailed("PMDファイルのオープンに失敗しました", err)
	}
	defer file.Close()

	modelData := model.NewPmxModel()
	reader := newPmdReader(file)
	if err := reader.Read(modelData); err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, io_common.NewIoParseFailed("PMDファイル情報の取得に失敗しました", err)
	}
	modelData.SetPath(path)
	modelData.SetFileModTime(info.ModTime().UnixNano())
	modelData.UpdateHash()
	return modelData, nil
}

// Save はPMDバイナリを保存する。
func (r *PmdRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		return io_common.NewIoEncodeFailed("PMD保存対象が不正です", nil)
	}
	savePath := path
	if savePath == "" {
		savePath = modelData.Path()
	}
	if savePath == "" {
		return io_common.NewIoSaveFailed("保存先パスが空です", nil)
	}

	file, err := os.Create(savePath)
	if err != nil {
		return io_common.NewIoSaveFailed("PMDファイルの作成に失敗しました", err)
	}
	defer file.Close()

	writer := newPmdWriter(file)
	if err := writer.Write(modelData, opts); err != nil {
		return err
	}
	modelData.SetPath(savePath)
	modelData.UpdateHash()
	return nil
}
