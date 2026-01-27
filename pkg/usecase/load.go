// 指示: miu200521358
package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

// iOverrideBoneInserter は不足ボーン補完のI/F。
type iOverrideBoneInserter interface {
	InsertShortageOverrideBones() error
}

const repositoryNotConfiguredErrorID = "95504"

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep io_common.IFileReader, path string) (*model.PmxModel, error) {
	if path == "" {
		return nil, nil
	}
	if rep == nil {
		return nil, newRepositoryNotConfiguredError("モデル")
	}
	data, err := rep.Load(path)
	if err != nil {
		return nil, err
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		return nil, io_common.NewIoFormatNotSupported("モデル形式が不正です", nil)
	}
	if modelData.Bones != nil {
		if inserter, ok := any(modelData.Bones).(iOverrideBoneInserter); ok {
			// システム用ボーンの追加に失敗しても処理は継続する。
			inserter.InsertShortageOverrideBones()
		}
	}
	return modelData, nil
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep io_common.IFileReader, path string) (*motion.VmdMotion, error) {
	if path == "" {
		return nil, nil
	}
	if rep == nil {
		return nil, newRepositoryNotConfiguredError("モーション")
	}
	data, err := rep.Load(path)
	if err != nil {
		return nil, err
	}
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		return nil, io_common.NewIoFormatNotSupported("モーション形式が不正です", nil)
	}
	return motionData, nil
}

// newRepositoryNotConfiguredError は読み込みリポジトリ未設定エラーを生成する。
func newRepositoryNotConfiguredError(target string) error {
	return merr.NewCommonError(repositoryNotConfiguredErrorID, merr.ErrorKindInternal, "読み込みリポジトリがありません: %s", nil, target)
}
