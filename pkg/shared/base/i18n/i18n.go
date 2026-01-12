// 指示: miu200521358
package i18n

// LangCode は言語コード。
type LangCode string

// LangChangeAction は言語変更時のアクション種別。
type LangChangeAction int

const (
	// LANG_CHANGE_NONE は変更なしを表す。
	LANG_CHANGE_NONE LangChangeAction = iota
	// LANG_CHANGE_RESTART_REQUIRED は再起動要求を表す。
	LANG_CHANGE_RESTART_REQUIRED
)

// DefaultLang は既定の言語。
const DefaultLang LangCode = "ja"

// II18n は多言語変換のI/F。
type II18n interface {
	Lang() LangCode
	SetLang(lang LangCode) (LangChangeAction, error)
	IsReady() bool
	T(key string) string
	TWithLang(lang LangCode, key string) string
}
