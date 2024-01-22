package console_view

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"

)

type ConsoleView struct {
	walk.WidgetBase
	logChan chan string
}

const (
	TEM_APPENDTEXT = win.WM_USER + 6
)

func NewConsoleView(parent walk.Container) (*ConsoleView, error) {
	lc := make(chan string, 1024)
	lv := &ConsoleView{logChan: lc}

	if err := walk.InitWidget(
		lv,
		parent,
		"EDIT",
		win.WS_TABSTOP|win.WS_VISIBLE|win.WS_VSCROLL|win.ES_MULTILINE|win.ES_WANTRETURN,
		win.WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}
	lv.setReadOnly(true)
	lv.SendMessage(win.EM_SETLIMITTEXT, 4294967295, 0)
	return lv, nil
}

func (*ConsoleView) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return walk.NewGreedyLayoutItem()
}

func (lv *ConsoleView) setTextSelection(start, end int) {
	lv.SendMessage(win.EM_SETSEL, uintptr(start), uintptr(end))
}

func (lv *ConsoleView) textLength() int {
	return int(lv.SendMessage(0x000E, uintptr(0), uintptr(0)))
}

func (lv *ConsoleView) AppendText(value string) {
	textLength := lv.textLength()
	lv.setTextSelection(textLength, textLength)
	s, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return
	}
	lv.SendMessage(win.EM_REPLACESEL, 0, uintptr(unsafe.Pointer(s)))
}

func (lv *ConsoleView) setReadOnly(readOnly bool) error {
	if lv.SendMessage(win.EM_SETREADONLY, uintptr(win.BoolToBOOL(readOnly)), 0) == 0 {
		return errors.New("fail to call EM_SETREADONLY")
	}

	return nil
}

func (lv *ConsoleView) PostAppendText(value string) {
	lv.logChan <- value
	win.PostMessage(lv.Handle(), TEM_APPENDTEXT, 0, 0)
}

func (lv *ConsoleView) Write(p []byte) (int, error) {
	lv.PostAppendText(string(p) + "\r\n")
	return len(p), nil
}

func (lv *ConsoleView) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_GETDLGCODE:
		if wParam == win.VK_RETURN {
			return win.DLGC_WANTALLKEYS
		}

		return win.DLGC_HASSETSEL | win.DLGC_WANTARROWS | win.DLGC_WANTCHARS
	case TEM_APPENDTEXT:
		select {
		case value := <-lv.logChan:
			lv.AppendText(value)
		default:
			return 0
		}
	}

	return lv.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
