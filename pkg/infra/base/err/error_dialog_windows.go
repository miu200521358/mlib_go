//go:build windows
// +build windows

// 指示: miu200521358
package err

import (
	"os"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/walk/pkg/walk"
)

// ShowErrorDialog は通常エラーのダイアログを表示する。
func ShowErrorDialog(appConfig *config.AppConfig, err error) bool {
	return showErrorDialog(appConfig, err, i18n.T("通常エラーが発生しました"), i18n.T("通常エラーヘッダー"), false)
}

// ShowFatalErrorDialog は致命エラーのダイアログを表示する。
func ShowFatalErrorDialog(appConfig *config.AppConfig, err error) bool {
	return showErrorDialog(appConfig, err, i18n.T("予期せぬエラーが発生しました"), i18n.T("予期せぬエラーヘッダー"), true)
}

// showErrorDialog はエラーダイアログの表示を行う。
func showErrorDialog(appConfig *config.AppConfig, err error, title string, header string, terminate bool) bool {
	message := replaceAppInfo(header, appConfig)
	if err != nil {
		message += "\n\n" + err.Error()
	}
	walk.MsgBox(nil, title, message, walk.MsgBoxIconError|walk.MsgBoxOK)
	if terminate {
		os.Exit(1)
	}
	return true
}

// replaceAppInfo はアプリ名/バージョンのプレースホルダを置換する。
func replaceAppInfo(message string, appConfig *config.AppConfig) string {
	if appConfig == nil {
		return message
	}
	name := appConfig.AppName
	version := appConfig.Version
	if name == "" {
		name = appConfig.Version
	}
	if version == "" {
		version = appConfig.AppName
	}
	out := strings.ReplaceAll(message, "{{.AppName}}", name)
	out = strings.ReplaceAll(out, "{{.AppVersion}}", version)
	return out
}
