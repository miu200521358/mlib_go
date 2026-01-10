// 指示: miu200521358
package i18n

import (
	"embed"
	"testing"
	"testing/fstest"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
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

	if action := i.SetLang(LANG_JA); action != i18n.LANG_CHANGE_NONE {
		t.Errorf("SetLang same: got=%v", action)
	}
	if action := i.SetLang(LANG_EN); action != i18n.LANG_CHANGE_RESTART_REQUIRED {
		t.Errorf("SetLang change: got=%v", action)
	}
	if len(cfg.setLang) == 0 || cfg.setLang[0] != "en" {
		t.Errorf("SetLang saved: got=%v", cfg.setLang)
	}
}

// TestDefaultI18nFallbacks は未初期化時の挙動を確認する。
func TestDefaultI18nFallbacks(t *testing.T) {
	prev := defaultI18n
	defaultI18n = nil
	t.Cleanup(func() { defaultI18n = prev })

	if CurrentLang() != i18n.DefaultLang {
		t.Errorf("CurrentLang default: got=%v", CurrentLang())
	}
	if T("missing") != "●●missing●●" {
		t.Errorf("T default: got=%v", T("missing"))
	}
	if TWithLang(LANG_EN, "missing") != "●●missing●●" {
		t.Errorf("TWithLang default: got=%v", TWithLang(LANG_EN, "missing"))
	}
	if SetLang(LANG_EN) != i18n.LANG_CHANGE_RESTART_REQUIRED {
		t.Errorf("SetLang default expected restart")
	}
}

// TestI18nNilReceiver はnil受信の分岐を確認する。
func TestI18nNilReceiver(t *testing.T) {
	var i *I18n
	if i.Lang() != i18n.DefaultLang {
		t.Errorf("Lang nil: got=%v", i.Lang())
	}
	if i.SetLang(LANG_EN) != i18n.LANG_CHANGE_RESTART_REQUIRED {
		t.Errorf("SetLang nil expected restart")
	}
	if i.IsReady() {
		t.Errorf("IsReady nil should be false")
	}
	if i.T("key") != "●●key●●" {
		t.Errorf("T nil: got=%v", i.T("key"))
	}
	if i.TWithLang("", "key") != "●●key●●" {
		t.Errorf("TWithLang nil: got=%v", i.TWithLang("", "key"))
	}
}

// TestLangDetection は言語検出と正規化を確認する。
func TestLangDetection(t *testing.T) {
	if detectLang(nil) != i18n.DefaultLang {
		t.Errorf("detectLang nil")
	}
	cfgEmpty := &stubUserConfig{lang: []string{}}
	if detectLang(cfgEmpty) != i18n.DefaultLang {
		t.Errorf("detectLang empty")
	}
	cfg := &stubUserConfig{lang: []string{"zh"}}
	if detectLang(cfg) != LANG_ZH {
		t.Errorf("detectLang zh: got=%v", detectLang(cfg))
	}
	cfg = &stubUserConfig{lang: []string{"xx"}}
	if detectLang(cfg) != i18n.DefaultLang {
		t.Errorf("detectLang invalid: got=%v", detectLang(cfg))
	}
	if normalizeLang("xx") != i18n.DefaultLang {
		t.Errorf("normalizeLang invalid")
	}
	if normalizeLang(LANG_KO) != LANG_KO {
		t.Errorf("normalizeLang valid")
	}
}

// TestInitI18nGlobals はInitI18nとグローバル関数を確認する。
func TestInitI18nGlobals(t *testing.T) {
	var emptyFS embed.FS
	cfg := &stubUserConfig{lang: []string{"en"}}

	prev := defaultI18n
	i := InitI18n(emptyFS, cfg)
	t.Cleanup(func() { defaultI18n = prev })

	if i == nil || CurrentLang() != LANG_EN {
		t.Errorf("InitI18n CurrentLang: got=%v", CurrentLang())
	}
	if SetLang(LANG_EN) != i18n.LANG_CHANGE_NONE {
		t.Errorf("SetLang same expected none")
	}
	if SetLang(LANG_JA) != i18n.LANG_CHANGE_RESTART_REQUIRED {
		t.Errorf("SetLang change expected restart")
	}
	if T("開く") == "●●開く●●" {
		t.Errorf("T should resolve key")
	}
	if TWithLang("", "開く") == "●●開く●●" {
		t.Errorf("TWithLang default should resolve key")
	}
}

// TestI18nReadyFlag はreadyフラグの分岐を確認する。
func TestI18nReadyFlag(t *testing.T) {
	i := &I18n{ready: false}
	if i.IsReady() {
		t.Errorf("IsReady should be false")
	}
	if i.T("key") != "●●key●●" {
		t.Errorf("T should return placeholder when not ready")
	}
}

// TestLoadMessagesAndLookup は読み込みと参照の分岐を確認する。
func TestLoadMessagesAndLookup(t *testing.T) {
	fsys := fstest.MapFS{
		"ok.json":   &fstest.MapFile{Data: []byte(`[{"id":"a","translation":"A"},{"id":"","translation":"skip"}]`)},
		"bad.json":  &fstest.MapFile{Data: []byte(`{invalid`)},
	}
	if got := loadMessages(fsys, "missing.json"); len(got) != 0 {
		t.Errorf("loadMessages missing: got=%v", got)
	}
	if got := loadMessages(fsys, "bad.json"); len(got) != 0 {
		t.Errorf("loadMessages bad json: got=%v", got)
	}
	ok := loadMessages(fsys, "ok.json")
	if ok["a"] != "A" || ok[""] != "" {
		t.Errorf("loadMessages ok: got=%v", ok)
	}

	merged := mergeMessages(map[string]string{"a": "A", "b": "B"}, map[string]string{"b": "BB"})
	if merged["b"] != "BB" {
		t.Errorf("mergeMessages override: got=%v", merged["b"])
	}

	msgs := map[LangCode]map[string]string{
		LANG_JA: {"k": "v"},
	}
	if lookupMessage(msgs, LANG_EN, "k") != "v" {
		t.Errorf("lookupMessage fallback failed")
	}
	if lookupMessage(msgs, LANG_JA, "missing") != "▼▼missing▼▼" {
		t.Errorf("lookupMessage missing failed")
	}
}
