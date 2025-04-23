//go:build windows
// +build windows

package controller

import (
	"bytes"
	"math"
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

func NewConsoleView(parent walk.Container, minWidth int, minHeight int) (*ConsoleView, error) {
	lc := make(chan string, 512)
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
	te.SetMinMaxSize(
		walk.Size{Width: minWidth, Height: minHeight},
		walk.Size{Width: minWidth * 100, Height: minHeight * 100})
	cv.Console = te
	cv.Console.SetReadOnly(true)
	cv.Console.SendMessage(win.EM_SETLIMITTEXT, math.MaxInt, 0)

	return cv, nil
}

func (controlView *ConsoleView) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return walk.NewGreedyLayoutItem()
}

func (controlView *ConsoleView) setTextSelection(start, end int) {
	controlView.Console.SendMessage(win.EM_SETSEL, uintptr(start), uintptr(end))
}

func (controlView *ConsoleView) textLength() int {
	return int(controlView.Console.SendMessage(0x000E, uintptr(0), uintptr(0)))
}

func (controlView *ConsoleView) AppendText(value string) {
	textLength := controlView.textLength()
	controlView.setTextSelection(textLength, textLength)
	uv, err := syscall.UTF16PtrFromString(value)
	if err == nil {
		controlView.Console.SendMessage(win.EM_REPLACESEL, 0, uintptr(unsafe.Pointer(uv)))
	}
}
func (controlView *ConsoleView) PostAppendText(value string) {
	controlView.mutex.Lock()
	defer controlView.mutex.Unlock()

	controlView.logChan <- value
	win.PostMessage(controlView.Handle(), TEM_APPENDTEXT, 0, 0)
}

func (controlView *ConsoleView) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_GETDLGCODE:
		if wParam == win.VK_RETURN {
			return win.DLGC_WANTALLKEYS
		}

		return win.DLGC_HASSETSEL | win.DLGC_WANTARROWS | win.DLGC_WANTCHARS
	case TEM_APPENDTEXT:
		select {
		case value := <-controlView.logChan:
			controlView.AppendText(value)
		default:
			return 0
		}
	}

	return controlView.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}

// Write はio.Writerインターフェースを実装し、logの出力をConsoleViewにリダイレクトします。
func (controlView *ConsoleView) Write(p []byte) (n int, err error) {

	// 改行が文字列内にある場合、コンソール内で改行が行われるよう置換する
	p = bytes.ReplaceAll(p, []byte("\n"), []byte("\r\n"))
	controlView.PostAppendText(string(p))

	return n, nil
}
