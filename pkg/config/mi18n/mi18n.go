package mi18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
)

var bundle *i18n.Bundle
var localizer *i18n.Localizer

//go:embed i18n/*
var commonI18nFiles embed.FS

func Initialize(appI18nFiles embed.FS) {
	langs := mconfig.LoadUserConfig("lang")

	var lang string
	var langTag language.Tag
	if len(langs) == 0 {
		lang = "ja"
	} else {
		lang = langs[0]
	}

	switch lang {
	case "ja":
		langTag = language.Japanese
	case "en":
		langTag = language.English
	case "zh":
		langTag = language.Chinese
	case "ko":
		langTag = language.Korean
	default:
		langTag = language.Japanese
	}

	bundle = i18n.NewBundle(langTag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFileFS(commonI18nFiles, "i18n/common.ja.json")
	bundle.LoadMessageFileFS(commonI18nFiles, "i18n/common.en.json")
	bundle.LoadMessageFileFS(commonI18nFiles, "i18n/common.zh.json")
	bundle.LoadMessageFileFS(commonI18nFiles, "i18n/common.ko.json")

	bundle.LoadMessageFileFS(appI18nFiles, "i18n/app.ja.json")
	bundle.LoadMessageFileFS(appI18nFiles, "i18n/app.en.json")
	bundle.LoadMessageFileFS(appI18nFiles, "i18n/app.zh.json")
	bundle.LoadMessageFileFS(appI18nFiles, "i18n/app.ko.json")

	localizer = i18n.NewLocalizer(bundle, lang)
}

func SetLang(lang string) {
	mconfig.SaveUserConfig("lang", lang, 1)
}

// T メッセージIDを元にメッセージを取得する
func T(key string, params ...map[string]interface{}) string {
	return t(localizer, key, params...)
}

// TWithLocale メッセージIDを元に指定ロケールでメッセージを取得する
func TWithLocale(lang string, key string, params ...map[string]interface{}) string {
	return t(i18n.NewLocalizer(bundle, lang), key, params...)
}

func t(l *i18n.Localizer, key string, params ...map[string]interface{}) string {
	if l == nil {
		return fmt.Sprintf("●●%s●●", key)
	}

	if len(params) == 0 {
		if translated, err := l.Localize(&i18n.LocalizeConfig{MessageID: key}); err == nil {
			return translated
		} else {
			return fmt.Sprintf("●●%s●●", key)
		}
	}

	if translated, err := l.Localize(&i18n.LocalizeConfig{MessageID: key, TemplateData: params[0]}); err == nil {
		return translated
	} else {
		return fmt.Sprintf("●●%s●●", key)
	}
}
