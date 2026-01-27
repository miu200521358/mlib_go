// 指示: miu200521358
package usecase

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	portio "github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

const (
	modelNotLoadedErrorID               = "95505"
	savePathInvalidErrorID              = "95506"
	saveRepositoryNotConfiguredErrorID  = "95507"
	savePathServiceNotConfiguredErrorID = "95508"
)

// PmxSaveRequest はPMX保存要求を表す。
type PmxSaveRequest struct {
	ModelPath              string
	ModelData              *model.PmxModel
	Writer                 portio.IFileWriter
	PathService            portio.IPathService
	SaveOptions            portio.SaveOptions
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

func buildPmxOutputPath(service portio.IPathService, path string) string {
	if service == nil || path == "" {
		return ""
	}
	dir, name, _ := service.SplitPath(path)
	if dir == "" || name == "" {
		return ""
	}
	return service.CreateOutputPath(filepath.Join(dir, name+".pmx"), "")
}

func newModelNotLoadedError(message string) error {
	if message == "" {
		message = "XまたはPMDファイルが読み込まれていません"
	}
	return merr.NewCommonError(modelNotLoadedErrorID, merr.ErrorKindValidate, message, nil)
}

func newSavePathInvalidError(message string) error {
	if message == "" {
		message = "保存先パスが不正です"
	}
	return merr.NewCommonError(savePathInvalidErrorID, merr.ErrorKindValidate, message, nil)
}

func newSaveRepositoryNotConfiguredError() error {
	return merr.NewCommonError(saveRepositoryNotConfiguredErrorID, merr.ErrorKindInternal, "保存リポジトリがありません", nil)
}

func newPathServiceNotConfiguredError() error {
	return merr.NewCommonError(savePathServiceNotConfiguredErrorID, merr.ErrorKindInternal, "保存先判定ができません", nil)
}
