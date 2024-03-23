package mi18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

var bundle *i18n.Bundle
var localizer *i18n.Localizer

func Initialize(resourceFiles embed.FS) {
	langs := mutils.LoadUserConfig("lang")

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
	bundle.LoadMessageFileFS(resourceFiles, fmt.Sprintf("resources/i18n/common.%s.json", lang))
	bundle.LoadMessageFileFS(resourceFiles, fmt.Sprintf("resources/i18n/app.%s.json", lang))

	localizer = i18n.NewLocalizer(bundle, lang)
}

func SetLang(lang string) {
	mutils.SaveUserConfig("lang", lang, 1)
}

func T(key string, param ...interface{}) string {
	if localizer == nil {
		return key
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key, TemplateData: param})
}

func TWithLocale(lang string, key string, param ...interface{}) string {
	if bundle == nil {
		return key
	}
	return i18n.NewLocalizer(bundle, lang).MustLocalize(&i18n.LocalizeConfig{MessageID: key, TemplateData: param})
}
