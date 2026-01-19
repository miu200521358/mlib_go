//go:build windows
// +build windows

// 指示: miu200521358
package err

import (
	"os"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/walk/pkg/declarative"
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
	errText := ""
	if err != nil {
		errText = err.Error()
	}
	text := message
	if errText != "" {
		text += "\n\n" + errText
	}
	// ToolTip追加で失敗する環境があるため、エラーダイアログ生成中だけ抑止する。
	prevEnv, hasEnv := os.LookupEnv("Env")
	if setErr := os.Setenv("Env", "debug"); setErr == nil {
		defer func() {
			if hasEnv {
				_ = os.Setenv("Env", prevEnv)
			} else {
				_ = os.Unsetenv("Env")
			}
		}()
	}
	closeText := i18n.T("エラーダイアログを閉じる")
	if terminate {
		closeText = i18n.T("アプリを終了")
	}
	var mw *walk.MainWindow
	var errView *walk.TextEdit
	if _, dialogErr := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    title,
		Size:     declarative.Size{Width: 680, Height: 520},
		MinSize:  declarative.Size{Width: 680, Height: 520},
		MaxSize:  declarative.Size{Width: 1200, Height: 900},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text: replaceAppInfo(header, appConfig),
			},
			declarative.TextEdit{
				Text:     strings.ReplaceAll(text, "\n", "\r\n"),
				ReadOnly: true,
				AssignTo: &errView,
				VScroll:  true,
				HScroll:  true,
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						Text: i18n.T("コミュニティ報告"),
						OnClicked: func() {
							if errView == nil {
								return
							}
							if copyErr := walk.Clipboard().SetText(errView.Text()); copyErr != nil {
								walk.MsgBox(nil, i18n.T("クリップボードコピー失敗"), copyErr.Error(), walk.MsgBoxIconError)
							}
						},
					},
					declarative.HSpacer{},
					declarative.PushButton{
						Text: closeText,
						OnClicked: func() {
							if terminate {
								os.Exit(1)
							}
							if mw != nil {
								mw.Close()
							}
						},
					},
				},
			},
		},
	}).Run(); dialogErr != nil {
		walk.MsgBox(nil, i18n.T("エラーダイアログ起動失敗"), dialogErr.Error(), walk.MsgBoxIconError)
		if terminate {
			os.Exit(1)
		}
		return false
	}
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
