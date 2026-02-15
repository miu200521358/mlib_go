// 指示: miu200521358
package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/usecase/messages"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

// iOverrideBoneInserter は不足ボーン補完のI/F。
type iOverrideBoneInserter interface {
	InsertShortageOverrideBones() error
}

const (
	repositoryNotConfiguredErrorID = "93501"
	ioFormatNotSupportedErrorID    = "14103"
)

// runInsertShortageOverrideBones は不足ボーン補完の実行関数です。
var runInsertShortageOverrideBones = func(inserter iOverrideBoneInserter) error {
	return inserter.InsertShortageOverrideBones()
}

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep io.IFileReader, path string) (*model.PmxModel, error) {
	result, err := LoadModelWithMeta(rep, path)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.Model, nil
}

// LoadModelWithMeta はモデル読み込み結果と継続警告を返す。
func LoadModelWithMeta(rep io.IFileReader, path string) (*ModelLoadResult, error) {
	result := &ModelLoadResult{}
	if path == "" {
		return result, nil
	}
	if rep == nil {
		return nil, newRepositoryNotConfiguredError(messages.LoadModelRepositoryNotConfigured)
	}
	data, err := rep.Load(path)
	if err != nil {
		return nil, err
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		return nil, newFormatNotSupportedError(messages.LoadModelFormatNotSupported)
	}
	if modelData.Bones != nil {
		if inserter, ok := any(modelData.Bones).(iOverrideBoneInserter); ok {
			if insertErr := runInsertShortageOverrideBones(inserter); insertErr != nil {
				result.Warnings = append(
					result.Warnings,
					newModelLoadWarning(messages.LoadModelOverrideBoneInsertWarning, insertErr.Error()),
				)
			}
		}
	}
	result.Model = modelData
	return result, nil
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep io.IFileReader, path string) (*motion.VmdMotion, error) {
	if path == "" {
		return nil, nil
	}
	if rep == nil {
		return nil, newRepositoryNotConfiguredError(messages.LoadMotionRepositoryNotConfigured)
	}
	data, err := rep.Load(path)
	if err != nil {
		return nil, err
	}
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		return nil, newFormatNotSupportedError(messages.LoadMotionFormatNotSupported)
	}
	return motionData, nil
}

// newRepositoryNotConfiguredError は読み込みリポジトリ未設定エラーを生成する。
func newRepositoryNotConfiguredError(messageKey string) error {
	return merr.NewCommonError(repositoryNotConfiguredErrorID, merr.ErrorKindInternal, messageKey, nil)
}

// newFormatNotSupportedError は形式未対応エラーを生成する。
func newFormatNotSupportedError(messageKey string) error {
	return merr.NewCommonError(ioFormatNotSupportedErrorID, merr.ErrorKindValidate, messageKey, nil)
}

// newModelLoadWarning はモデル読み込み継続時の警告情報を生成する。
func newModelLoadWarning(messageKey string, messageParams ...any) ModelLoadWarning {
	return ModelLoadWarning{
		MessageKey:    messageKey,
		MessageParams: messageParams,
	}
}
