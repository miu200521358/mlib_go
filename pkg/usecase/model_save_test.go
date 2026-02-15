// 指示: miu200521358
package usecase

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
	"github.com/miu200521358/mlib_go/pkg/usecase/messages"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

// modelSaveTestWriter は保存呼び出しを記録するスタブ。
type modelSaveTestWriter struct {
	saveErr   error
	savedPath string
	savedData hashable.IHashable
	savedOpts io.SaveOptions
}

// Save は引数を記録し、設定エラーを返す。
func (w *modelSaveTestWriter) Save(path string, data hashable.IHashable, opts io.SaveOptions) error {
	if w == nil {
		return nil
	}
	w.savedPath = path
	w.savedData = data
	w.savedOpts = opts
	return w.saveErr
}

// modelSaveTestPathService は保存先判定を制御するスタブ。
type modelSaveTestPathService struct {
	canSave    bool
	outputPath string
}

// CanSave は保存可否を返す。
func (s *modelSaveTestPathService) CanSave(path string) bool {
	if s == nil {
		return false
	}
	return s.canSave
}

// CreateOutputPath は固定の出力パスを返す。
func (s *modelSaveTestPathService) CreateOutputPath(originalPath, label string) string {
	if s == nil {
		return ""
	}
	return s.outputPath
}

// SplitPath は入力パスを分解する。
func (s *modelSaveTestPathService) SplitPath(path string) (dir, name, ext string) {
	ext = filepath.Ext(path)
	base := filepath.Base(path)
	name = base[:len(base)-len(ext)]
	dir = filepath.Dir(path)
	return dir, name, ext
}

func TestSaveModelAsPmx(t *testing.T) {
	tempDir := t.TempDir()
	modelData := model.NewPmxModel()
	modelData.SetPath(filepath.Join(tempDir, "sample.x"))
	outputPath := filepath.Join(tempDir, "out.pmx")

	cases := []struct {
		name           string
		request        PmxSaveRequest
		wantErrID      string
		wantErrKind    merr.ErrorKind
		wantErrKey     string
		wantOutputPath string
		wantWriterCall bool
	}{
		{
			name: "モデル未設定エラー",
			request: PmxSaveRequest{
				ModelPath: "",
			},
			wantErrID:   modelNotLoadedErrorID,
			wantErrKind: merr.ErrorKindValidate,
			wantErrKey:  messages.SaveModelNotLoaded,
		},
		{
			name: "保存先判定サービス未設定エラー",
			request: PmxSaveRequest{
				ModelPath: modelData.Path(),
				ModelData: modelData,
				Writer:    &modelSaveTestWriter{},
			},
			wantErrID:   savePathServiceNotConfiguredErrorID,
			wantErrKind: merr.ErrorKindInternal,
			wantErrKey:  messages.SavePathServiceNotConfigured,
		},
		{
			name: "保存リポジトリ未設定エラー",
			request: PmxSaveRequest{
				ModelPath:   modelData.Path(),
				ModelData:   modelData,
				PathService: &modelSaveTestPathService{canSave: true, outputPath: outputPath},
			},
			wantErrID:   saveRepositoryNotConfiguredErrorID,
			wantErrKind: merr.ErrorKindInternal,
			wantErrKey:  messages.SaveRepositoryNotConfigured,
		},
		{
			name: "保存先パス不正エラー",
			request: PmxSaveRequest{
				ModelPath:   modelData.Path(),
				ModelData:   modelData,
				Writer:      &modelSaveTestWriter{},
				PathService: &modelSaveTestPathService{canSave: false, outputPath: outputPath},
			},
			wantErrID:   savePathInvalidErrorID,
			wantErrKind: merr.ErrorKindValidate,
			wantErrKey:  messages.SavePathInvalid,
		},
		{
			name: "保存先パス不正エラー（カスタムキー）",
			request: PmxSaveRequest{
				ModelPath:              modelData.Path(),
				ModelData:              modelData,
				Writer:                 &modelSaveTestWriter{},
				PathService:            &modelSaveTestPathService{canSave: false, outputPath: ""},
				InvalidSavePathMessage: "custom.invalid.path.key",
			},
			wantErrID:   savePathInvalidErrorID,
			wantErrKind: merr.ErrorKindValidate,
			wantErrKey:  "custom.invalid.path.key",
		},
		{
			name: "正常保存",
			request: PmxSaveRequest{
				ModelPath:   modelData.Path(),
				ModelData:   modelData,
				Writer:      &modelSaveTestWriter{},
				PathService: &modelSaveTestPathService{canSave: true, outputPath: outputPath},
				SaveOptions: io.SaveOptions{IncludeSystem: true},
			},
			wantOutputPath: outputPath,
			wantWriterCall: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SaveModelAsPmx(tc.request)
			if tc.wantErrID != "" {
				if err == nil {
					t.Fatalf("エラーが必要ですが nil です")
				}
				ce, ok := err.(*merr.CommonError)
				if !ok {
					t.Fatalf("CommonError ではありません: %T", err)
				}
				if ce.ErrorID() != tc.wantErrID {
					t.Fatalf("ErrorID が不正です: got=%s want=%s", ce.ErrorID(), tc.wantErrID)
				}
				if ce.ErrorKind() != tc.wantErrKind {
					t.Fatalf("ErrorKind が不正です: got=%s want=%s", ce.ErrorKind(), tc.wantErrKind)
				}
				if ce.MessageKey() != tc.wantErrKey {
					t.Fatalf("MessageKey が不正です: got=%s want=%s", ce.MessageKey(), tc.wantErrKey)
				}
				return
			}
			if err != nil {
				t.Fatalf("想定外エラーです: %v", err)
			}
			if result == nil {
				t.Fatalf("結果が nil です")
			}
			if result.OutputPath != tc.wantOutputPath {
				t.Fatalf("出力パスが不正です: got=%s want=%s", result.OutputPath, tc.wantOutputPath)
			}

			writer, ok := tc.request.Writer.(*modelSaveTestWriter)
			if tc.wantWriterCall && !ok {
				t.Fatalf("writer が検証できません")
			}
			if tc.wantWriterCall {
				if writer.savedPath != tc.wantOutputPath {
					t.Fatalf("writer保存パスが不正です: got=%s want=%s", writer.savedPath, tc.wantOutputPath)
				}
				if writer.savedData != tc.request.ModelData {
					t.Fatalf("writer保存データが不正です")
				}
				if writer.savedOpts != tc.request.SaveOptions {
					t.Fatalf("writer保存オプションが不正です: got=%+v want=%+v", writer.savedOpts, tc.request.SaveOptions)
				}
			}
		})
	}
}

func TestSaveModelAsPmxWriterError(t *testing.T) {
	saveErr := errors.New("save failed")
	modelData := model.NewPmxModel()
	modelData.SetPath(filepath.Join(t.TempDir(), "sample.x"))
	writer := &modelSaveTestWriter{saveErr: saveErr}
	request := PmxSaveRequest{
		ModelPath:   modelData.Path(),
		ModelData:   modelData,
		Writer:      writer,
		PathService: &modelSaveTestPathService{canSave: true, outputPath: filepath.Join(t.TempDir(), "out.pmx")},
	}

	_, err := SaveModelAsPmx(request)
	if !errors.Is(err, saveErr) {
		t.Fatalf("writerエラーが伝播していません: %v", err)
	}
}

func TestIsPmxConvertiblePath(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{path: "sample.x", want: true},
		{path: "sample.X", want: true},
		{path: "sample.pmd", want: true},
		{path: "sample.PMD", want: true},
		{path: "sample.pmx", want: false},
		{path: "", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := IsPmxConvertiblePath(tc.path)
			if got != tc.want {
				t.Fatalf("判定が不正です: path=%s got=%v want=%v", tc.path, got, tc.want)
			}
		})
	}
}
