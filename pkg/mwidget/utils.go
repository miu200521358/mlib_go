package mwidget

import (
	"bytes"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"

)

func CheckError(err error, w *MWindow, message string) {
	if err != nil {
		if w != nil {
			walk.MsgBox(w.MainWindow, message, err.Error(), walk.MsgBoxIconError)
			w.Close()
		} else {
			walk.MsgBox(nil, message, err.Error(), walk.MsgBoxIconError)
		}
	}
}

var MARGIN_ZERO = walk.Margins{HNear: 0, VNear: 0, HFar: 0, VFar: 0}
var MARGIN_SMALL = walk.Margins{HNear: 3, VNear: 3, HFar: 3, VFar: 3}
var MARGIN_MEDIUM = walk.Margins{HNear: 6, VNear: 6, HFar: 6, VFar: 6}

// エラー監視
func RecoverFromPanic(mWindow *MWindow) {
	if r := recover(); r != nil {
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
					Text:     string(bytes.ReplaceAll(stackTrace, []byte("\n"), []byte("\r\n"))),
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
						if mWindow != nil {
							mWindow.Close()
						}
					},
				},
			},
		}).Run(); err != nil {
			walk.MsgBox(nil, mi18n.T("エラーダイアログ起動失敗"), string(stackTrace), walk.MsgBoxIconError)
		}

		if mWindow != nil {
			mWindow.Close()
		}
	}
}
