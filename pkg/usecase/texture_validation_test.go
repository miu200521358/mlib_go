// 指示: miu200521358
package usecase

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/usecase/messages"
)

// textureValidationTestValidator はテクスチャ検証の戻り値を制御するスタブ。
type textureValidationTestValidator struct {
	existsResults  map[string]bool
	existsErrors   map[string]error
	validateErrors map[string]error
}

// ExistsFile は事前設定した存在判定結果を返す。
func (v *textureValidationTestValidator) ExistsFile(path string) (bool, error) {
	if v == nil {
		return false, nil
	}
	if err, ok := v.existsErrors[path]; ok {
		return false, err
	}
	exists, ok := v.existsResults[path]
	if !ok {
		return false, nil
	}
	return exists, nil
}

// ValidateImage は事前設定した検証エラーを返す。
func (v *textureValidationTestValidator) ValidateImage(path string) error {
	if v == nil {
		return nil
	}
	if err, ok := v.validateErrors[path]; ok {
		return err
	}
	return nil
}

// newTextureValidationTestModel は単一テクスチャを持つモデルを生成する。
func newTextureValidationTestModel(t *testing.T, textureName string) *model.PmxModel {
	t.Helper()
	modelData := model.NewPmxModel()
	modelData.SetPath(filepath.Join(t.TempDir(), "model.pmx"))
	texture := model.NewTexture()
	texture.SetName(textureName)
	modelData.Textures.AppendRaw(texture)
	return modelData
}

func TestValidateModelTextures(t *testing.T) {
	existsFailedErr := errors.New("exists failed")
	validateFailedErr := errors.New("decode failed")

	cases := []struct {
		name             string
		textureName      string
		setupValidator   func(texturePath string) *textureValidationTestValidator
		wantIssueLen     int
		wantErrorLen     int
		wantErrorKind    merr.ErrorKind
		wantErrorID      string
		wantErrorKey     string
		wantIssuePath    bool
		wantTextureValid bool
	}{
		{
			name:        "空テクスチャ名はIssueのみ",
			textureName: "",
			setupValidator: func(texturePath string) *textureValidationTestValidator {
				return &textureValidationTestValidator{}
			},
			wantIssueLen:     1,
			wantErrorLen:     0,
			wantIssuePath:    false,
			wantTextureValid: false,
		},
		{
			name:        "存在確認エラーはCommonErrorで記録",
			textureName: "missing.png",
			setupValidator: func(texturePath string) *textureValidationTestValidator {
				return &textureValidationTestValidator{
					existsErrors: map[string]error{texturePath: existsFailedErr},
				}
			},
			wantIssueLen:     1,
			wantErrorLen:     1,
			wantErrorKind:    merr.ErrorKindValidate,
			wantErrorID:      textureExistsValidationFailedErrorID,
			wantErrorKey:     messages.TextureExistsValidationFailed,
			wantIssuePath:    true,
			wantTextureValid: false,
		},
		{
			name:        "存在しないテクスチャはIssueのみ",
			textureName: "missing.png",
			setupValidator: func(texturePath string) *textureValidationTestValidator {
				return &textureValidationTestValidator{
					existsResults: map[string]bool{texturePath: false},
				}
			},
			wantIssueLen:     1,
			wantErrorLen:     0,
			wantIssuePath:    true,
			wantTextureValid: false,
		},
		{
			name:        "画像検証エラーはCommonErrorで記録",
			textureName: "broken.png",
			setupValidator: func(texturePath string) *textureValidationTestValidator {
				return &textureValidationTestValidator{
					existsResults:  map[string]bool{texturePath: true},
					validateErrors: map[string]error{texturePath: validateFailedErr},
				}
			},
			wantIssueLen:     1,
			wantErrorLen:     1,
			wantErrorKind:    merr.ErrorKindValidate,
			wantErrorID:      textureImageValidationFailedErrorID,
			wantErrorKey:     messages.TextureImageValidationFailed,
			wantIssuePath:    true,
			wantTextureValid: false,
		},
		{
			name:        "正常テクスチャはValid=true",
			textureName: "ok.png",
			setupValidator: func(texturePath string) *textureValidationTestValidator {
				return &textureValidationTestValidator{
					existsResults: map[string]bool{texturePath: true},
				}
			},
			wantIssueLen:     0,
			wantErrorLen:     0,
			wantIssuePath:    false,
			wantTextureValid: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			modelData := newTextureValidationTestModel(t, tc.textureName)
			texturePath := ""
			if tc.textureName != "" {
				texturePath = filepath.Join(filepath.Dir(modelData.Path()), tc.textureName)
			}
			validator := tc.setupValidator(texturePath)

			result := ValidateModelTextures(modelData, validator)
			if result == nil {
				t.Fatalf("結果が nil です")
			}
			if len(result.Issues) != tc.wantIssueLen {
				t.Fatalf("Issue件数が不正です: got=%d want=%d", len(result.Issues), tc.wantIssueLen)
			}
			if len(result.Errors) != tc.wantErrorLen {
				t.Fatalf("Error件数が不正です: got=%d want=%d", len(result.Errors), tc.wantErrorLen)
			}
			if tc.wantErrorLen > 0 {
				ce, ok := result.Errors[0].(*merr.CommonError)
				if !ok {
					t.Fatalf("CommonError ではありません: %T", result.Errors[0])
				}
				if ce.ErrorKind() != tc.wantErrorKind {
					t.Fatalf("ErrorKind が不正です: got=%s want=%s", ce.ErrorKind(), tc.wantErrorKind)
				}
				if ce.MessageKey() != tc.wantErrorKey {
					t.Fatalf("MessageKey が不正です: got=%s want=%s", ce.MessageKey(), tc.wantErrorKey)
				}
				if ce.ErrorID() == "" {
					t.Fatalf("ErrorID が空です")
				}
				if ce.ErrorID() != tc.wantErrorID {
					t.Fatalf("ErrorID が不正です: got=%s want=%s", ce.ErrorID(), tc.wantErrorID)
				}
			}
			if tc.wantIssueLen > 0 && tc.wantIssuePath {
				if result.Issues[0].Path != texturePath {
					t.Fatalf("Issueパスが不正です: got=%s want=%s", result.Issues[0].Path, texturePath)
				}
			}
			texture := modelData.Textures.Values()[0]
			if texture.IsValid() != tc.wantTextureValid {
				t.Fatalf("テクスチャ有効フラグが不正です: got=%v want=%v", texture.IsValid(), tc.wantTextureValid)
			}
		})
	}
}

func TestValidateModelTexturesNilInput(t *testing.T) {
	result := ValidateModelTextures(nil, &textureValidationTestValidator{})
	if result == nil {
		t.Fatalf("結果が nil です")
	}
	if len(result.Issues) != 0 || len(result.Errors) != 0 {
		t.Fatalf("nil入力では空結果想定です: %+v", result)
	}
}
