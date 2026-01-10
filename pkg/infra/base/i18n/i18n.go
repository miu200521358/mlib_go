// 指示: miu200521358
package i18n

import (
	"embed"
	"encoding/json"
	"io/fs"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
)

//go:embed i18n/common.*.json
var commonI18nFiles embed.FS

// LangCode は言語コード。
type LangCode = i18n.LangCode

const (
	// LANG_JA は日本語。
	LANG_JA LangCode = "ja"
	// LANG_EN は英語。
	LANG_EN LangCode = "en"
	// LANG_ZH は中国語。
	LANG_ZH LangCode = "zh"
	// LANG_KO は韓国語。
	LANG_KO LangCode = "ko"
)

type messageEntry struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
}

// I18n はi18nの実装。
type I18n struct {
	lang       LangCode
	ready      bool
	messages   map[LangCode]map[string]string
	userConfig config.IUserConfig
}

var defaultI18n *I18n
var defaultMu sync.Mutex

// InitI18n はi18nを初期化する。
func InitI18n(appFiles embed.FS, userConfig config.IUserConfig) {
	instance := initI18nFS(appFiles, userConfig)
	defaultMu.Lock()
	defaultI18n = instance
	defaultMu.Unlock()
}

// initI18nFS はFSからi18nを初期化する。
func initI18nFS(appFiles fs.FS, userConfig config.IUserConfig) *I18n {
	lang := detectLang(userConfig)
	langs := []LangCode{LANG_JA, LANG_EN, LANG_ZH, LANG_KO}

	messages := make(map[LangCode]map[string]string, len(langs))
	for _, lc := range langs {
		common := loadMessages(commonI18nFiles, "i18n/common."+string(lc)+".json")
		app := loadMessages(appFiles, "cmd/i18n/app."+string(lc)+".json")
		merged := mergeMessages(common, app)
		messages[lc] = merged
	}

	return &I18n{
		lang:       lang,
		ready:      true,
		messages:   messages,
		userConfig: userConfig,
	}
}

// SetLang は言語を保存する。
func SetLang(lang LangCode) i18n.LangChangeAction {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultI18n == nil {
		return i18n.LANG_CHANGE_RESTART_REQUIRED
	}
	return defaultI18n.SetLang(lang)
}

// CurrentLang は現在言語を返す。
func CurrentLang() LangCode {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultI18n == nil {
		return i18n.DefaultLang
	}
	return defaultI18n.Lang()
}

// T はメッセージを取得する。
func T(key string) string {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultI18n == nil {
		return "●●" + key + "●●"
	}
	return defaultI18n.T(key)
}

// TWithLang は指定言語でメッセージを取得する。
func TWithLang(lang LangCode, key string) string {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultI18n == nil {
		return "●●" + key + "●●"
	}
	return defaultI18n.TWithLang(lang, key)
}

// Lang は現在言語を返す。
func (i *I18n) Lang() LangCode {
	if i == nil {
		return i18n.DefaultLang
	}
	return i.lang
}

// SetLang は言語を保存し、再起動が必要か返す。
func (i *I18n) SetLang(lang LangCode) i18n.LangChangeAction {
	if i == nil {
		return i18n.LANG_CHANGE_RESTART_REQUIRED
	}
	lang = normalizeLang(lang)
	if lang == i.lang {
		return i18n.LANG_CHANGE_NONE
	}
	if i.userConfig != nil {
		_ = i.userConfig.SetStringSlice(config.UserConfigKeyLang, []string{string(lang)}, 1)
	}
	return i18n.LANG_CHANGE_RESTART_REQUIRED
}

// IsReady は初期化済みか判定する。
func (i *I18n) IsReady() bool {
	if i == nil {
		return false
	}
	return i.ready
}

// T は現在言語でメッセージを取得する。
func (i *I18n) T(key string) string {
	if i == nil || !i.ready {
		return "●●" + key + "●●"
	}
	return lookupMessage(i.messages, i.lang, key)
}

// TWithLang は指定言語でメッセージを取得する。
func (i *I18n) TWithLang(lang LangCode, key string) string {
	if i == nil || !i.ready {
		return "●●" + key + "●●"
	}
	if lang == "" {
		lang = i18n.DefaultLang
	}
	return lookupMessage(i.messages, lang, key)
}

// detectLang はユーザー設定から言語を取得する。
func detectLang(userConfig config.IUserConfig) LangCode {
	if userConfig == nil {
		return i18n.DefaultLang
	}
	values := userConfig.GetStringSlice(config.UserConfigKeyLang)
	if len(values) == 0 {
		return i18n.DefaultLang
	}
	switch LangCode(values[0]) {
	case LANG_JA, LANG_EN, LANG_ZH, LANG_KO:
		return LangCode(values[0])
	default:
		return i18n.DefaultLang
	}
}

// normalizeLang は言語コードを既知の値へ正規化する。
func normalizeLang(lang LangCode) LangCode {
	switch lang {
	case LANG_JA, LANG_EN, LANG_ZH, LANG_KO:
		return lang
	default:
		return i18n.DefaultLang
	}
}

// loadMessages はJSONを読み込んでマップ化する。
func loadMessages(files fs.FS, path string) map[string]string {
	data, err := fs.ReadFile(files, path)
	if err != nil {
		return map[string]string{}
	}
	var entries []messageEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.ID == "" {
			continue
		}
		out[e.ID] = e.Translation
	}
	return out
}

// mergeMessages は共通とアプリを結合する。
func mergeMessages(common, app map[string]string) map[string]string {
	merged := make(map[string]string, len(common)+len(app))
	for k, v := range common {
		merged[k] = v
	}
	for k, v := range app {
		merged[k] = v
	}
	return merged
}

// lookupMessage は辞書からメッセージを引く。
func lookupMessage(messages map[LangCode]map[string]string, lang LangCode, key string) string {
	langMap, ok := messages[lang]
	if !ok {
		langMap = messages[i18n.DefaultLang]
	}
	value, ok := langMap[key]
	if !ok {
		return "▼▼" + key + "▼▼"
	}
	return value
}
