// 指示: miu200521358
package usecase

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/usecase/messages"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

const (
	modelNotLoadedErrorID               = "13501"
	savePathInvalidErrorID              = "13502"
	saveRepositoryNotConfiguredErrorID  = "93502"
	savePathServiceNotConfiguredErrorID = "93503"
)

// PmxSaveRequest はPMX保存要求を表す。
type PmxSaveRequest struct {
	ModelPath              string
	ModelData              *model.PmxModel
	Writer                 io.IFileWriter
	PathService            io.IPathService
	SaveOptions            io.SaveOptions
	MissingModelMessage    string
	InvalidSavePathMessage string
}

// PmxSaveResult はPMX保存結果を表す。
type PmxSaveResult struct {
	OutputPath string
}

// SaveModelAsPmx はX/PMDモデルをPMX形式で保存する。
func SaveModelAsPmx(request PmxSaveRequest) (*PmxSaveResult, error) {
	result := &PmxSaveResult{}
	modelPath := request.ModelPath
	if modelPath == "" && request.ModelData != nil {
		modelPath = request.ModelData.Path()
	}
	if modelPath == "" || request.ModelData == nil || !IsPmxConvertiblePath(modelPath) {
		return result, newModelNotLoadedError(request.MissingModelMessage)
	}
	if request.PathService == nil {
		return result, newPathServiceNotConfiguredError()
	}
	if request.Writer == nil {
		return result, newSaveRepositoryNotConfiguredError()
	}

	outputPath := buildPmxOutputPath(request.PathService, modelPath)
	if outputPath == "" || !request.PathService.CanSave(outputPath) {
		return result, newSavePathInvalidError(request.InvalidSavePathMessage)
	}
	if err := request.Writer.Save(outputPath, request.ModelData, request.SaveOptions); err != nil {
		return result, err
	}
	result.OutputPath = outputPath
	return result, nil
}

// IsPmxConvertiblePath はPMX保存対象のパスか判定する。
func IsPmxConvertiblePath(path string) bool {
	if path == "" {
		return false
	}
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, ".x") || strings.EqualFold(ext, ".pmd")
}

// buildPmxOutputPath は入力モデルパスからPMX保存先パスを生成する。
func buildPmxOutputPath(service io.IPathService, path string) string {
	if service == nil || path == "" {
		return ""
	}
	dir, name, _ := service.SplitPath(path)
	if dir == "" || name == "" {
		return ""
	}
	return service.CreateOutputPath(filepath.Join(dir, name+".pmx"), "")
}

// newModelNotLoadedError は保存対象モデル未設定エラーを生成する。
func newModelNotLoadedError(messageKey string) error {
	if messageKey == "" {
		messageKey = messages.SaveModelNotLoaded
	}
	return merr.NewCommonError(modelNotLoadedErrorID, merr.ErrorKindValidate, messageKey, nil)
}

// newSavePathInvalidError は保存先パス不正エラーを生成する。
func newSavePathInvalidError(messageKey string) error {
	if messageKey == "" {
		messageKey = messages.SavePathInvalid
	}
	return merr.NewCommonError(savePathInvalidErrorID, merr.ErrorKindValidate, messageKey, nil)
}

// newSaveRepositoryNotConfiguredError は保存リポジトリ未設定エラーを生成する。
func newSaveRepositoryNotConfiguredError() error {
	return merr.NewCommonError(
		saveRepositoryNotConfiguredErrorID,
		merr.ErrorKindInternal,
		messages.SaveRepositoryNotConfigured,
		nil,
	)
}

// newPathServiceNotConfiguredError は保存先判定サービス未設定エラーを生成する。
func newPathServiceNotConfiguredError() error {
	return merr.NewCommonError(
		savePathServiceNotConfiguredErrorID,
		merr.ErrorKindInternal,
		messages.SavePathServiceNotConfigured,
		nil,
	)
}
