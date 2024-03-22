package mutils

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18n struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

func NewI18n(resourceFiles embed.FS) *I18n {
	langs := LoadUserConfig("lang")

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

	bundle := i18n.NewBundle(langTag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFileFS(resourceFiles, fmt.Sprintf("resources/i18n/common.%s.json", lang))
	bundle.LoadMessageFileFS(resourceFiles, fmt.Sprintf("resources/i18n/app.%s.json", lang))

	localizer := i18n.NewLocalizer(bundle, lang)

	return &I18n{bundle: bundle, localizer: localizer}
}

func (i *I18n) SetLang(lang string) {
	SaveUserConfig("lang", lang, 1)
}

func (i *I18n) T(key string) string {
	return i.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key})
}

func (i *I18n) TP(key string, param map[string]interface{}) string {
	return i.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key, TemplateData: param})
}
