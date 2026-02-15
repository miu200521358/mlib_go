// 指示: miu200521358
package usecase

import (
	"errors"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
	"github.com/miu200521358/mlib_go/pkg/usecase/messages"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

// loadTestFileReader は usecase.Load* テスト向けの読み込みスタブ。
type loadTestFileReader struct {
	data hashable.IHashable
	err  error
}

// CanLoad は読み込み可否を返す。
func (r *loadTestFileReader) CanLoad(path string) bool {
	return true
}

// Load は事前設定したデータまたはエラーを返す。
func (r *loadTestFileReader) Load(path string) (hashable.IHashable, error) {
	if r == nil {
		return nil, nil
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.data, nil
}

// InferName はテスト用の固定名を返す。
func (r *loadTestFileReader) InferName(path string) string {
	return "test"
}

// assertCommonError は共通エラーのID/種別/キーを検証する。
func assertCommonError(t *testing.T, err error, wantID string, wantKind merr.ErrorKind, wantMessageKey string) {
	t.Helper()
	ce, ok := err.(*merr.CommonError)
	if !ok {
		t.Fatalf("CommonError ではありません: %T", err)
	}
	if ce.ErrorID() != wantID {
		t.Fatalf("ErrorID が不正です: got=%s want=%s", ce.ErrorID(), wantID)
	}
	if ce.ErrorKind() != wantKind {
		t.Fatalf("ErrorKind が不正です: got=%s want=%s", ce.ErrorKind(), wantKind)
	}
	if ce.MessageKey() != wantMessageKey {
		t.Fatalf("MessageKey が不正です: got=%s want=%s", ce.MessageKey(), wantMessageKey)
	}
}

func TestLoadModelWithMeta(t *testing.T) {
	originalInserter := runInsertShortageOverrideBones
	t.Cleanup(func() {
		runInsertShortageOverrideBones = originalInserter
	})

	loadFailedErr := errors.New("load failed")
	overrideFailedErr := errors.New("insert failed")

	cases := []struct {
		name             string
		path             string
		reader           io.IFileReader
		overrideInserter func(iOverrideBoneInserter) error
		wantErr          error
		wantErrID        string
		wantErrKind      merr.ErrorKind
		wantErrKey       string
		wantModel        bool
		wantWarningLen   int
		wantWarningKey   string
	}{
		{
			name:           "空パスは空結果を返す",
			path:           "",
			reader:         &loadTestFileReader{},
			wantModel:      false,
			wantWarningLen: 0,
		},
		{
			name:        "リポジトリ未設定は内部エラー",
			path:        "model.pmx",
			reader:      nil,
			wantErrID:   repositoryNotConfiguredErrorID,
			wantErrKind: merr.ErrorKindInternal,
			wantErrKey:  messages.LoadModelRepositoryNotConfigured,
		},
		{
			name:    "Load失敗はそのまま返す",
			path:    "model.pmx",
			reader:  &loadTestFileReader{err: loadFailedErr},
			wantErr: loadFailedErr,
		},
		{
			name:        "型不一致は形式エラー",
			path:        "model.pmx",
			reader:      &loadTestFileReader{data: hashable.NewHashableBase("", "dummy")},
			wantErrID:   ioFormatNotSupportedErrorID,
			wantErrKind: merr.ErrorKindValidate,
			wantErrKey:  messages.LoadModelFormatNotSupported,
		},
		{
			name: "不足ボーン補完失敗はWarningで継続",
			path: "model.pmx",
			reader: &loadTestFileReader{
				data: model.NewPmxModel(),
			},
			overrideInserter: func(iOverrideBoneInserter) error {
				return overrideFailedErr
			},
			wantModel:      true,
			wantWarningLen: 1,
			wantWarningKey: messages.LoadModelOverrideBoneInsertWarning,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runInsertShortageOverrideBones = originalInserter
			if tc.overrideInserter != nil {
				runInsertShortageOverrideBones = tc.overrideInserter
			}

			result, err := LoadModelWithMeta(tc.reader, tc.path)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("想定外エラーです: got=%v want=%v", err, tc.wantErr)
				}
				return
			}
			if tc.wantErrID != "" {
				if err == nil {
					t.Fatalf("エラーが必要ですが nil です")
				}
				assertCommonError(t, err, tc.wantErrID, tc.wantErrKind, tc.wantErrKey)
				return
			}
			if err != nil {
				t.Fatalf("想定外エラーです: %v", err)
			}
			if result == nil {
				t.Fatalf("結果が nil です")
			}
			if (result.Model != nil) != tc.wantModel {
				t.Fatalf("モデル有無が不正です: got=%v want=%v", result.Model != nil, tc.wantModel)
			}
			if len(result.Warnings) != tc.wantWarningLen {
				t.Fatalf("Warning件数が不正です: got=%d want=%d", len(result.Warnings), tc.wantWarningLen)
			}
			if tc.wantWarningKey != "" {
				if result.Warnings[0].MessageKey != tc.wantWarningKey {
					t.Fatalf("Warningキーが不正です: got=%s want=%s", result.Warnings[0].MessageKey, tc.wantWarningKey)
				}
				if len(result.Warnings[0].MessageParams) == 0 {
					t.Fatalf("Warningパラメータが不足しています")
				}
			}
		})
	}
}

func TestLoadModel(t *testing.T) {
	originalInserter := runInsertShortageOverrideBones
	t.Cleanup(func() {
		runInsertShortageOverrideBones = originalInserter
	})

	cases := []struct {
		name             string
		path             string
		reader           io.IFileReader
		overrideInserter func(iOverrideBoneInserter) error
		wantNilModel     bool
		wantErrID        string
		wantErrKind      merr.ErrorKind
		wantErrKey       string
	}{
		{
			name:         "空パスはnilモデル",
			path:         "",
			reader:       &loadTestFileReader{},
			wantNilModel: true,
		},
		{
			name: "Warning発生でもモデル返却",
			path: "model.pmx",
			reader: &loadTestFileReader{
				data: model.NewPmxModel(),
			},
			overrideInserter: func(iOverrideBoneInserter) error {
				return errors.New("warning")
			},
			wantNilModel: false,
		},
		{
			name:        "リポジトリ未設定エラー",
			path:        "model.pmx",
			reader:      nil,
			wantErrID:   repositoryNotConfiguredErrorID,
			wantErrKind: merr.ErrorKindInternal,
			wantErrKey:  messages.LoadModelRepositoryNotConfigured,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runInsertShortageOverrideBones = originalInserter
			if tc.overrideInserter != nil {
				runInsertShortageOverrideBones = tc.overrideInserter
			}

			modelData, err := LoadModel(tc.reader, tc.path)
			if tc.wantErrID != "" {
				if err == nil {
					t.Fatalf("エラーが必要ですが nil です")
				}
				assertCommonError(t, err, tc.wantErrID, tc.wantErrKind, tc.wantErrKey)
				return
			}
			if err != nil {
				t.Fatalf("想定外エラーです: %v", err)
			}
			if (modelData == nil) != tc.wantNilModel {
				t.Fatalf("モデル有無が不正です: got=%v want=%v", modelData == nil, tc.wantNilModel)
			}
		})
	}
}

func TestLoadMotion(t *testing.T) {
	cases := []struct {
		name        string
		path        string
		reader      io.IFileReader
		wantErrID   string
		wantErrKind merr.ErrorKind
		wantErrKey  string
		wantMotion  bool
	}{
		{
			name:        "リポジトリ未設定エラー",
			path:        "motion.vmd",
			reader:      nil,
			wantErrID:   repositoryNotConfiguredErrorID,
			wantErrKind: merr.ErrorKindInternal,
			wantErrKey:  messages.LoadMotionRepositoryNotConfigured,
		},
		{
			name:        "型不一致は形式エラー",
			path:        "motion.vmd",
			reader:      &loadTestFileReader{data: hashable.NewHashableBase("", "dummy")},
			wantErrID:   ioFormatNotSupportedErrorID,
			wantErrKind: merr.ErrorKindValidate,
			wantErrKey:  messages.LoadMotionFormatNotSupported,
		},
		{
			name:       "正常系",
			path:       "motion.vmd",
			reader:     &loadTestFileReader{data: motion.NewVmdMotion("motion.vmd")},
			wantMotion: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			motionData, err := LoadMotion(tc.reader, tc.path)
			if tc.wantErrID != "" {
				if err == nil {
					t.Fatalf("エラーが必要ですが nil です")
				}
				assertCommonError(t, err, tc.wantErrID, tc.wantErrKind, tc.wantErrKey)
				return
			}
			if err != nil {
				t.Fatalf("想定外エラーです: %v", err)
			}
			if (motionData != nil) != tc.wantMotion {
				t.Fatalf("モーション有無が不正です: got=%v want=%v", motionData != nil, tc.wantMotion)
			}
		})
	}
}
