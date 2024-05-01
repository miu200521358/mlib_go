//go:build windows
// +build windows

package mwidget

import (
	"bytes"
	"sync"
	"syscall"
	"unsafe"

	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

// ConsoleView は、walk.WidgetBaseを埋め込んでカスタムウィジェットを作成します。
type ConsoleView struct {
	walk.WidgetBase
	Console *walk.TextEdit
	logChan chan string
	mutex   sync.Mutex
}

const (
	ConsoleViewClass = "ConsoleView Class"
	TEM_APPENDTEXT   = win.WM_USER + 6
)

func NewConsoleView(parent walk.Container) (*ConsoleView, error) {
	lc := make(chan string, 1024)
	cv := &ConsoleView{logChan: lc}

	if err := walk.InitWidget(
		cv,
		parent,
		ConsoleViewClass,
		win.WS_DISABLED,
		0); err != nil {
		return nil, err
	}

	// テキストエディットを作成
	te, err := walk.NewTextEditWithStyle(parent, win.WS_VISIBLE|win.WS_VSCROLL|win.ES_MULTILINE)
	if err != nil {
		return nil, err
	}
	cv.Console = te
	cv.Console.SetReadOnly(true)
	cv.Console.SendMessage(win.EM_SETLIMITTEXT, 4294967295, 0)

	return cv, nil
}

func (cv *ConsoleView) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return walk.NewGreedyLayoutItem()
}

func (cv *ConsoleView) setTextSelection(start, end int) {
	cv.Console.SendMessage(win.EM_SETSEL, uintptr(start), uintptr(end))
}

func (cv *ConsoleView) textLength() int {
	return int(cv.Console.SendMessage(0x000E, uintptr(0), uintptr(0)))
}

func (cv *ConsoleView) AppendText(value string) {
	textLength := cv.textLength()
	cv.setTextSelection(textLength, textLength)
	cv.Console.SendMessage(win.EM_REPLACESEL, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))))
}
func (cv *ConsoleView) PostAppendText(value string) {
	cv.mutex.Lock()
	defer cv.mutex.Unlock()

	cv.logChan <- value
	win.PostMessage(cv.Handle(), TEM_APPENDTEXT, 0, 0)
}

func (cv *ConsoleView) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_GETDLGCODE:
		if wParam == win.VK_RETURN {
			return win.DLGC_WANTALLKEYS
		}

		return win.DLGC_HASSETSEL | win.DLGC_WANTARROWS | win.DLGC_WANTCHARS
	case TEM_APPENDTEXT:
		select {
		case value := <-cv.logChan:
			cv.AppendText(value)
		default:
			return 0
		}
	}

	return cv.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

// Write はio.Writerインターフェースを実装し、logの出力をConsoleViewにリダイレクトします。
func (cv *ConsoleView) Write(p []byte) (n int, err error) {

	// 改行が文字列内にある場合、コンソール内で改行が行われるよう置換する
	p = bytes.ReplaceAll(p, []byte("\n"), []byte("\r\n"))
	cv.PostAppendText(string(p))

	return n, nil
}
