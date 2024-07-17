//go:build windows
// +build windows

package widget

import (
	"bytes"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/sys/windows"
)

var (
	user32          = windows.NewLazySystemDLL("user32.dll")
	procMessageBeep = user32.NewProc("MessageBeep")
	MB_ICONASTERISK = 0x00000040
)

func Beep() {
	procMessageBeep.Call(uintptr(MB_ICONASTERISK))
}

func RaiseError(err error) {
	errMsg := err.Error()
	stackTrace := debug.Stack()

	var errT *walk.TextEdit
	if _, err := (declarative.MainWindow{
		Title:   mi18n.T("予期せぬエラーが発生しました"),
		Size:    declarative.Size{Width: 500, Height: 400},
		MinSize: declarative.Size{Width: 500, Height: 400},
		MaxSize: declarative.Size{Width: 500, Height: 400},
		Layout:  declarative.VBox{},
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text: mi18n.T("予期せぬエラーヘッダー"),
			},
			declarative.TextEdit{
				Text: string(bytes.ReplaceAll([]byte(errMsg), []byte("\n"), []byte("\r\n"))) +
					string("\r\n------------\r\n") +
					string(bytes.ReplaceAll(stackTrace, []byte("\n"), []byte("\r\n"))),
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
					exec.Command("cmd", "/c", "start", "https://com.nicovideo.jp/community/co5387214").Start()
				},
			},
			declarative.PushButton{
				Text: mi18n.T("アプリを終了"),
				OnClicked: func() {
					os.Exit(1)
				},
			},
		},
	}).Run(); err != nil {
		walk.MsgBox(nil, mi18n.T("エラーダイアログ起動失敗"), string(stackTrace), walk.MsgBoxIconError)
	}
}
