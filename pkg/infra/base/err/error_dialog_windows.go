//go:build windows
// +build windows

// 指示: miu200521358
package err

import (
	"embed"
	"errors"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	sharedi18n "github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

// ShowErrorDialog は通常エラーのダイアログを表示する。
func ShowErrorDialog(appConfig *config.AppConfig, err error) bool {
	return showErrorDialog(appConfig, err, i18n.T("通常エラーが発生しました"), i18n.T("通常エラーヘッダー"), false)
}

// ShowFatalErrorDialog は致命エラーのダイアログを表示する。
func ShowFatalErrorDialog(appConfig *config.AppConfig, err error) bool {
	return showErrorDialog(appConfig, err, i18n.T("予期せぬエラーが発生しました"), i18n.T("予期せぬエラーヘッダー"), true)
}

// BuildErrorText はエラー本文を生成して返す。
func BuildErrorText(err error) string {
	return buildErrorText(i18n.Default(), err)
}

// showErrorDialog はエラーダイアログの表示を行う。
func showErrorDialog(appConfig *config.AppConfig, err error, title string, header string, terminate bool) bool {
	message := replaceAppInfo(header, appConfig)
	errText := buildErrorText(i18n.Default(), err)
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
	iconName := "error_48dp_EF6C00_FILL1_wght400_GRAD0_opsz48.png"
	if terminate {
		iconName = "dangerous_48dp_C62828_FILL1_wght400_GRAD0_opsz48.png"
	}
	errorIcon, _ := readIconFromEmbedFS(errorIcons, iconName)
	closeText := i18n.T("エラーダイアログを閉じる")
	if terminate {
		closeText = i18n.T("アプリを終了")
	}
	var mw *walk.MainWindow
	var errView *walk.TextEdit
	if _, dialogErr := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    title,
		Icon:     errorIcon,
		Size:     declarative.Size{Width: 680, Height: 520},
		MinSize:  declarative.Size{Width: 680, Height: 520},
		MaxSize:  declarative.Size{Width: 1200, Height: 900},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.ImageView{
						Image:   errorIcon,
						MinSize: declarative.Size{Width: 48, Height: 48},
						MaxSize: declarative.Size{Width: 48, Height: 48},
					},
					declarative.TextLabel{
						Text: replaceAppInfo(header, appConfig),
					},
				},
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
						Text: i18n.T("エラーをダウンロード"),
						OnClicked: func() {
							if errView == nil {
								return
							}
							fd := new(walk.FileDialog)
							fd.Title = i18n.T("エラーをダウンロード")
							fd.Filter = i18n.T("テキストファイル") + " (*.txt)|*.txt|" + i18n.T("すべてのファイル") + " (*.*)|*.*"
							jst := time.FixedZone("JST", 9*60*60)
							fd.FilePath = "mlib_error_" + time.Now().In(jst).Format("200601021504") + ".txt"
							fd.Flags |= win.OFN_OVERWRITEPROMPT
							ok, dlgErr := fd.ShowSave(mw)
							if dlgErr != nil {
								walk.MsgBox(mw, i18n.T("保存失敗"), dlgErr.Error(), walk.MsgBoxIconError)
								return
							}
							if !ok {
								return
							}
							path := fd.FilePath
							if filepath.Ext(path) == "" {
								path += ".txt"
							}
							if writeErr := os.WriteFile(path, []byte(errView.Text()), 0o644); writeErr != nil {
								walk.MsgBox(mw, i18n.T("保存失敗"), writeErr.Error(), walk.MsgBoxIconError)
								return
							}
						},
					},
					declarative.PushButton{
						Text: i18n.T("コミュニティで報告"),
						OnClicked: func() {
							exec.Command("cmd", "/c", "start", "https://discord.gg/MW2Bn47aCN").Start()
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

type (
	iErrorIDProvider interface {
		ErrorID() string
	}
	iErrorKindProvider interface {
		ErrorKind() merr.ErrorKind
	}
	iErrorMessageProvider interface {
		MessageKey() string
		MessageParams() []any
	}
)

// buildErrorText はエラーダイアログ向けの本文を生成する。
func buildErrorText(translator sharedi18n.II18n, err error) string {
	if err == nil {
		return ""
	}
	errID := extractErrorID(err)
	msg := formatErrorMessage(translator, err)
	if msg == "" {
		msg = err.Error()
	}
	if errID != "" {
		msg = fmt.Sprintf("%s: %s\n%s", translateKey(translator, "エラーID"), errID, msg)
	}
	if extractErrorKind(err) == merr.ErrorKindValidate {
		if remedy := formatErrorRemedy(translator, errID); remedy != "" {
			msg += "\n\n" + fmt.Sprintf("%s:\n%s", translateKey(translator, "対処方法"), remedy)
		}
	}
	return msg
}

// formatErrorMessage はエラーのメッセージキーを翻訳して返す。
func formatErrorMessage(translator sharedi18n.II18n, err error) string {
	summary := formatErrorSummary(translator, err)
	detail := formatErrorDetail(translator, err)
	if summary == "" {
		return detail
	}
	if detail == "" || detail == summary {
		return summary
	}
	return summary + "\n" + fmt.Sprintf("%s: %s", translateKey(translator, "詳細"), detail)
}

// formatErrorRemedy はエラー管理表から対処法メッセージを取得する。
func formatErrorRemedy(translator sharedi18n.II18n, errID string) string {
	rec, err := merr.FindRecord(errID)
	if err != nil || rec == nil || rec.Remedy == "" {
		return ""
	}
	return translateKey(translator, rec.Remedy)
}

// formatErrorSummary はエラー管理表のSummaryを翻訳して返す。
func formatErrorSummary(translator sharedi18n.II18n, err error) string {
	rec, err := merr.FindRecord(extractErrorID(err))
	if err != nil || rec == nil || rec.Summary == "" {
		return ""
	}
	return translateKey(translator, rec.Summary)
}

// formatErrorDetail はエラーのメッセージキーを翻訳して返す。
func formatErrorDetail(translator sharedi18n.II18n, err error) string {
	key, params := extractMessageKey(err)
	if key == "" {
		return ""
	}
	text := translateKey(translator, key)
	if len(params) > 0 {
		return fmt.Sprintf(text, params...)
	}
	return text
}

// translateKey は翻訳済みのキーを返し、欠落時はキー自身を返す。
func translateKey(translator sharedi18n.II18n, key string) string {
	if key == "" {
		return ""
	}
	if translator == nil || !translator.IsReady() {
		return key
	}
	out := translator.T(key)
	if isMissingTranslation(out, key) {
		return key
	}
	return out
}

// extractErrorID はエラーIDを取得する。
func extractErrorID(err error) string {
	if err == nil {
		return ""
	}
	var provider iErrorIDProvider
	if errors.As(err, &provider) {
		return provider.ErrorID()
	}
	return ""
}

// extractErrorKind はエラー種別を取得する。
func extractErrorKind(err error) merr.ErrorKind {
	if err == nil {
		return ""
	}
	var provider iErrorKindProvider
	if errors.As(err, &provider) {
		return provider.ErrorKind()
	}
	return ""
}

// extractMessageKey はメッセージキーとパラメータを取得する。
func extractMessageKey(err error) (string, []any) {
	if err == nil {
		return "", nil
	}
	var provider iErrorMessageProvider
	if errors.As(err, &provider) {
		return provider.MessageKey(), provider.MessageParams()
	}
	return "", nil
}

// isMissingTranslation は未定義キーの表示か判定する。
func isMissingTranslation(text string, key string) bool {
	if text == "●●"+key+"●●" {
		return true
	}
	if text == "▼▼"+key+"▼▼" {
		return true
	}
	return false
}

//go:embed *.png
var errorIcons embed.FS

// readIconFromEmbedFS は埋め込み画像からアイコンを生成する。
func readIconFromEmbedFS(f embed.FS, name string) (*walk.Icon, error) {
	file, err := f.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	icon, err := walk.NewIconFromImageForDPI(img, 96)
	if err != nil {
		return nil, err
	}

	return icon, nil
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
