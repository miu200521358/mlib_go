// 指示: miu200521358
package i18n

import (
	"embed"
	"encoding/json"
	"io/fs"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
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
func InitI18n(appFiles embed.FS, userConfig config.IUserConfig) error {
	instance, err := initI18nFS(appFiles, userConfig)
	if err != nil {
		return err
	}
	defaultMu.Lock()
	defaultI18n = instance
	defaultMu.Unlock()
	return nil
}

// initI18nFS はFSからi18nを初期化する。
func initI18nFS(appFiles fs.FS, userConfig config.IUserConfig) (*I18n, error) {
	lang, err := detectLang(userConfig)
	if err != nil {
		return nil, err
	}
	langs := []LangCode{LANG_JA, LANG_EN, LANG_ZH, LANG_KO}

	messages := make(map[LangCode]map[string]string, len(langs))
	for _, lc := range langs {
		common, err := loadMessages(commonI18nFiles, "i18n/common."+string(lc)+".json")
		if err != nil {
			return nil, err
		}
		app, err := loadMessages(appFiles, "i18n/app."+string(lc)+".json")
		if err != nil {
			return nil, err
		}
		merged := mergeMessages(common, app)
		messages[lc] = merged
	}

	return &I18n{
		lang:       lang,
		ready:      true,
		messages:   messages,
		userConfig: userConfig,
	}, nil
}

// SetLang は言語を保存する。
func SetLang(lang LangCode) (i18n.LangChangeAction, error) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultI18n == nil {
		return i18n.LANG_CHANGE_RESTART_REQUIRED, nil
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

// Default は既定のi18nインスタンスを返す。
func Default() i18n.II18n {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	return defaultI18n
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
func (i *I18n) SetLang(lang LangCode) (i18n.LangChangeAction, error) {
	if i == nil {
		return i18n.LANG_CHANGE_RESTART_REQUIRED, nil
	}
	lang = normalizeLang(lang)
	if lang == i.lang {
		return i18n.LANG_CHANGE_NONE, nil
	}
	if i.userConfig != nil {
		if err := i.userConfig.SetStringSlice(config.UserConfigKeyLang, []string{string(lang)}, 1); err != nil {
			return i18n.LANG_CHANGE_RESTART_REQUIRED, err
		}
	}
	return i18n.LANG_CHANGE_RESTART_REQUIRED, nil
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

// detectLang はユーザー設定から言語を取得する（読込失敗時は error を返す）。
func detectLang(userConfig config.IUserConfig) (LangCode, error) {
	if userConfig == nil {
		return i18n.DefaultLang, nil
	}
	values, err := userConfig.GetStringSlice(config.UserConfigKeyLang)
	if err != nil {
		return i18n.DefaultLang, err
	}
	if len(values) == 0 {
		return i18n.DefaultLang, nil
	}
	switch LangCode(values[0]) {
	case LANG_JA, LANG_EN, LANG_ZH, LANG_KO:
		return LangCode(values[0]), nil
	default:
		return i18n.DefaultLang, nil
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
func loadMessages(files fs.FS, path string) (map[string]string, error) {
	data, err := fs.ReadFile(files, path)
	if err != nil {
		return nil, baseerr.NewFsPackageError("i18nメッセージ読込に失敗しました: "+path, err)
	}
	var entries []messageEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, baseerr.NewJsonPackageError("i18nメッセージ解析に失敗しました: "+path, err)
	}
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.ID == "" {
			continue
		}
		out[e.ID] = e.Translation
	}
	return out, nil
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
