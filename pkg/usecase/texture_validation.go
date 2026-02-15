// 指示: miu200521358
package usecase

import (
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/usecase/messages"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

const (
	textureExistsValidationFailedErrorID = "13503"
	textureImageValidationFailedErrorID  = "13504"
)

// TextureValidationIssue はテクスチャ検証で無効だった要素を表す。
type TextureValidationIssue struct {
	Name string
	Path string
}

// TextureValidationResult はテクスチャ検証結果を表す。
type TextureValidationResult struct {
	Issues []TextureValidationIssue
	Errors []error
}

// ValidateModelTextures はモデルのテクスチャ有効性を検証する。
func ValidateModelTextures(modelData *model.PmxModel, validator io.ITextureValidator) *TextureValidationResult {
	result := &TextureValidationResult{}
	if modelData == nil || modelData.Textures == nil || validator == nil {
		return result
	}

	baseDir := filepath.Dir(modelData.Path())
	for _, texture := range modelData.Textures.Values() {
		if texture == nil {
			continue
		}
		name := texture.Name()
		if name == "" {
			texture.SetValid(false)
			result.Issues = append(result.Issues, TextureValidationIssue{Name: name})
			continue
		}
		texturePath := name
		if !filepath.IsAbs(texturePath) {
			texturePath = filepath.Join(baseDir, texturePath)
		}
		exists, err := validator.ExistsFile(texturePath)
		if err != nil {
			result.Errors = append(
				result.Errors,
				merr.NewCommonError(
					textureExistsValidationFailedErrorID,
					merr.ErrorKindValidate,
					messages.TextureExistsValidationFailed,
					err,
					texturePath,
				),
			)
			texture.SetValid(false)
			result.Issues = append(result.Issues, TextureValidationIssue{Name: name, Path: texturePath})
			continue
		}
		if !exists {
			texture.SetValid(false)
			result.Issues = append(result.Issues, TextureValidationIssue{Name: name, Path: texturePath})
			continue
		}
		if err := validator.ValidateImage(texturePath); err != nil {
			result.Errors = append(
				result.Errors,
				merr.NewCommonError(
					textureImageValidationFailedErrorID,
					merr.ErrorKindValidate,
					messages.TextureImageValidationFailed,
					err,
					texturePath,
				),
			)
			texture.SetValid(false)
			result.Issues = append(result.Issues, TextureValidationIssue{Name: name, Path: texturePath})
			continue
		}
		texture.SetValid(true)
	}
	return result
}
