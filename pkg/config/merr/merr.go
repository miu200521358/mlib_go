package merr

import (
	"bytes"
	"embed"
	"errors"
	"image/png"
	"os"
	"os/exec"
	"runtime"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

var NameNotFoundError = errors.New("name not found")

var ParentNotFoundError = errors.New("parent not found")

var TerminateError = errors.New("terminate error")

func dumpAllGoroutines() string {
	buf := make([]byte, 1<<20) // 1MB のバッファを確保
	n := runtime.Stack(buf, true)
	return string(bytes.ReplaceAll(buf[:n], []byte("\n"), []byte("\r\n")))
}

func showErrorDialog(appConfig *mconfig.AppConfig, err error, titleKey, msgKey, btnKey string,
	icon *walk.Icon, isAppClose bool) bool {
	errMsg := err.Error()
	stackTrace := dumpAllGoroutines()

	if !appConfig.IsSetEnv() {
		panic(err)
	}

	var errT *walk.TextEdit
	var mw *walk.MainWindow
	if _, err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    mi18n.T(titleKey),
		Icon:     appConfig.Icon,
		Size:     declarative.Size{Width: 500, Height: 400},
		MinSize:  declarative.Size{Width: 500, Height: 400},
		MaxSize:  declarative.Size{Width: 500, Height: 400},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{MarginsZero: true},
				Children: []declarative.Widget{
					declarative.ImageView{
						Image:   icon,
						MinSize: declarative.Size{Width: 48, Height: 48},
						MaxSize: declarative.Size{Width: 48, Height: 48},
					},
					declarative.TextLabel{
						Text: mi18n.T(msgKey, map[string]any{"AppName": appConfig.Name, "AppVersion": appConfig.Version}),
					},
				},
			},
			declarative.TextEdit{
				Text: string("```") +
					string(bytes.ReplaceAll([]byte(errMsg), []byte("\n"), []byte("\r\n"))) +
					string("\r\n------------\r\n") +
					stackTrace +
					string("```"),
				ReadOnly: true,
				AssignTo: &errT,
				VScroll:  true,
				HScroll:  true,
			},
			declarative.PushButton{
				Text:      mi18n.T("コミュニティ報告"),
				Alignment: declarative.AlignHFarVNear,
				OnClicked: func() {
					if err := walk.Clipboard().SetText(errT.Text()); err != nil {
						walk.MsgBox(nil, mi18n.T("クリップボードコピー失敗"),
							string(stackTrace), walk.MsgBoxIconError)
					}
					exec.Command("cmd", "/c", "start", "https://discord.gg/MW2Bn47aCN").Start()
				},
			},
			declarative.PushButton{
				Text: mi18n.T(btnKey),
				OnClicked: func() {
					if isAppClose {
						os.Exit(1)
					} else {
						mw.Close()
					}
				},
			},
		},
	}).Run(); err != nil {
		walk.MsgBox(nil, mi18n.T("エラーダイアログ起動失敗"), string(stackTrace), walk.MsgBoxIconError)
		return false
	}

	return true
}

//go:embed *.png
var images embed.FS

func ShowErrorDialog(appConfig *mconfig.AppConfig, err error) bool {
	icon, _ := readIconFromEmbedFS(images, "error_48dp_EF6C00_FILL1_wght400_GRAD0_opsz48.png")
	return showErrorDialog(appConfig, err, "通常エラーが発生しました", "通常エラーヘッダー", "エラーダイアログを閉じる",
		icon, false)
}

func ShowFatalErrorDialog(appConfig *mconfig.AppConfig, err error) bool {
	icon, _ := readIconFromEmbedFS(images, "dangerous_48dp_C62828_FILL1_wght400_GRAD0_opsz48.png")
	return showErrorDialog(appConfig, err, "予期せぬエラーが発生しました", "予期せぬエラーヘッダー", "アプリを終了",
		icon, true)
}

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
