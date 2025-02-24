//go:build windows
// +build windows

package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
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

func ShowErrorDialog(isSetEnv bool, err error) {
	errMsg := err.Error()
	stackTrace := debug.Stack()

	if !isSetEnv {
		panic(err)
	}

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
				Text: string("```") +
					string(bytes.ReplaceAll([]byte(errMsg), []byte("\n"), []byte("\r\n"))) +
					string("\r\n------------\r\n") +
					string(bytes.ReplaceAll(stackTrace, []byte("\n"), []byte("\r\n"))) +
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

// SafeExecute は関数でpanicが発生した場合にダイアログを表示する
func SafeExecute(isSetEnv bool, f func()) {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			var errMsg string
			if recoveredErr, ok := r.(error); ok {
				errMsg = recoveredErr.Error()
			} else {
				errMsg = fmt.Sprintf("%v", r)
			}
			err := fmt.Errorf("panic: %s\n%s", errMsg, stackTrace)
			ShowErrorDialog(isSetEnv, err)
		}
	}()

	f()
}
