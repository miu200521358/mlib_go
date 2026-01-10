// 指示: miu200521358
package mi18n

import (
	"testing"
	"testing/fstest"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	basei18n "github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
)

type stubUserConfig struct {
	lang     []string
	setLang  []string
}

// Get は未使用のため空実装。
func (s *stubUserConfig) Get(key string) (any, bool) { return nil, false }

// Set は未使用のため空実装。
func (s *stubUserConfig) Set(key string, value any) error { return nil }

// GetStringSlice は言語設定を返す。
func (s *stubUserConfig) GetStringSlice(key string) []string {
	if key == config.UserConfigKeyLang {
		return s.lang
	}
	return nil
}

// SetStringSlice は言語設定の保存を記録する。
func (s *stubUserConfig) SetStringSlice(key string, values []string, limit int) error {
	if key == config.UserConfigKeyLang {
		s.setLang = values
	}
	return nil
}

// GetBool は未使用のため既定値を返す。
func (s *stubUserConfig) GetBool(key string, defaultValue bool) bool { return defaultValue }

// SetBool は未使用のため空実装。
func (s *stubUserConfig) SetBool(key string, value bool) error { return nil }

// GetInt は未使用のため既定値を返す。
func (s *stubUserConfig) GetInt(key string, defaultValue int) int { return defaultValue }

// SetInt は未使用のため空実装。
func (s *stubUserConfig) SetInt(key string, value int) error { return nil }

// GetAll は未使用のため空実装。
func (s *stubUserConfig) GetAll(key string) ([]string, map[string]any) { return nil, map[string]any{} }

// AppRootDir は未使用のため空実装。
func (s *stubUserConfig) AppRootDir() (string, error) { return "", nil }

// TestI18nTranslations は翻訳取得を確認する。
func TestI18nTranslations(t *testing.T) {
	fsys := fstest.MapFS{
		"i18n/app.ja.json": &fstest.MapFile{Data: []byte(`[{"id":"hello","translation":"こんにちは"}]`)},
		"i18n/app.en.json": &fstest.MapFile{Data: []byte(`[{"id":"hello","translation":"hello"}]`)},
		"i18n/app.zh.json": &fstest.MapFile{Data: []byte(`[]`)},
		"i18n/app.ko.json": &fstest.MapFile{Data: []byte(`[]`)},
	}
	cfg := &stubUserConfig{lang: []string{"ja"}}
	i := initI18nFS(fsys, cfg)

	if i.Lang() != LANG_JA {
		t.Errorf("Lang: got=%v", i.Lang())
	}
	if i.T("hello") != "こんにちは" {
		t.Errorf("T: got=%v", i.T("hello"))
	}
	if i.TWithLang(LANG_EN, "hello") != "hello" {
		t.Errorf("TWithLang: got=%v", i.TWithLang(LANG_EN, "hello"))
	}
	if i.T("missing") != "▼▼missing▼▼" {
		t.Errorf("missing key: got=%v", i.T("missing"))
	}
}

// TestSetLangAction は言語変更の戻り値を確認する。
func TestSetLangAction(t *testing.T) {
	fsys := fstest.MapFS{
		"i18n/app.ja.json": &fstest.MapFile{Data: []byte(`[]`)},
		"i18n/app.en.json": &fstest.MapFile{Data: []byte(`[]`)},
		"i18n/app.zh.json": &fstest.MapFile{Data: []byte(`[]`)},
		"i18n/app.ko.json": &fstest.MapFile{Data: []byte(`[]`)},
	}
	cfg := &stubUserConfig{lang: []string{"ja"}}
	i := initI18nFS(fsys, cfg)

	if action := i.SetLang(LANG_JA); action != basei18n.LANG_CHANGE_NONE {
		t.Errorf("SetLang same: got=%v", action)
	}
	if action := i.SetLang(LANG_EN); action != basei18n.LANG_CHANGE_RESTART_REQUIRED {
		t.Errorf("SetLang change: got=%v", action)
	}
	if len(cfg.setLang) == 0 || cfg.setLang[0] != "en" {
		t.Errorf("SetLang saved: got=%v", cfg.setLang)
	}
}
