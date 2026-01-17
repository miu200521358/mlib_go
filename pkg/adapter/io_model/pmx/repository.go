// 指示: miu200521358
package pmx

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

// PmxRepository はPMXバイナリ入出力を表す。
type PmxRepository struct{}

// NewPmxRepository はPmxRepositoryを生成する。
func NewPmxRepository() *PmxRepository {
	return &PmxRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *PmxRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".pmx")
}

// InferName はパスから表示名を推定する。
func (r *PmxRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はPMXバイナリを読み込む。
func (r *PmxRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoParseFailed("PMXファイルのオープンに失敗しました", err)
	}
	defer file.Close()

	modelData := model.NewPmxModel()
	reader := newPmxReader(file)
	if err := reader.Read(modelData); err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, io_common.NewIoParseFailed("PMXファイル情報の取得に失敗しました", err)
	}
	modelData.SetPath(path)
	modelData.SetFileModTime(info.ModTime().UnixNano())
	modelData.UpdateHash()
	return modelData, nil
}

// Save はPMXバイナリを保存する。
func (r *PmxRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		return io_common.NewIoEncodeFailed("PMX保存対象が不正です", nil)
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
		return io_common.NewIoSaveFailed("PMXファイルの作成に失敗しました", err)
	}
	defer file.Close()

	writer := newPmxWriter(file)
	if err := writer.Write(modelData, opts); err != nil {
		return err
	}
	modelData.SetPath(savePath)
	modelData.UpdateHash()
	return nil
}
