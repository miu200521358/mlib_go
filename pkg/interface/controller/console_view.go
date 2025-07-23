//go:build windows
// +build windows

package controller

import (
	"bytes"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

// ConsoleView は、walk.WidgetBaseを埋め込んでカスタムウィジェットを作成します。
type ConsoleView struct {
	walk.WidgetBase
	Console   *walk.TextEdit
	logChan   chan string
	mutex     sync.Mutex
	lineCount int // 現在の行数をカウント
}

const (
	ConsoleViewClass = "ConsoleView Class"
	TEM_APPENDTEXT   = win.WM_USER + 6
	MaxLines         = 10000 // 最大行数
)

func NewConsoleView(parent walk.Container, minWidth int, minHeight int) (*ConsoleView, error) {
	lc := make(chan string, 10000) // バッファサイズを10000に増加
	cv := &ConsoleView{logChan: lc, lineCount: 0}

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
	cv.Console.SendMessage(win.EM_SETLIMITTEXT, 4294967295, 0)

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
	// 行数をカウント
	controlView.lineCount += strings.Count(value, "\r\n")

	textLength := controlView.textLength()
	controlView.setTextSelection(textLength, textLength)
	uv, err := syscall.UTF16PtrFromString(value)
	if err == nil {
		controlView.Console.SendMessage(win.EM_REPLACESEL, 0, uintptr(unsafe.Pointer(uv)))
	}

	// 行数制限のチェック（100行ごとに実行してパフォーマンスを保つ）
	if controlView.lineCount%100 == 0 && controlView.lineCount > MaxLines {
		controlView.trimLines()
	}
}
func (controlView *ConsoleView) PostAppendText(value string) {
	controlView.mutex.Lock()
	defer controlView.mutex.Unlock()

	// 非ブロッキング送信：チャンネルがフルの場合は古いメッセージを破棄
	select {
	case controlView.logChan <- value:
		// 正常に送信できた場合
		win.PostMessage(controlView.Handle(), TEM_APPENDTEXT, 0, 0)
	default:
		// チャンネルがフルの場合、古いメッセージを1つ破棄して新しいメッセージを追加
		select {
		case <-controlView.logChan:
			// 古いメッセージを破棄
		default:
		}
		// 新しいメッセージを送信
		select {
		case controlView.logChan <- value:
			win.PostMessage(controlView.Handle(), TEM_APPENDTEXT, 0, 0)
		default:
			// それでも送信できない場合は諦める
		}
	}
}

// trimLines は、行数が上限を超えた場合に古い行を削除します
func (controlView *ConsoleView) trimLines() {
	text := controlView.Console.Text()
	lines := strings.Split(text, "\r\n")

	if len(lines) > MaxLines {
		// 先頭から削除する行数を計算
		linesToRemove := len(lines) - MaxLines
		// 新しいテキストを作成
		newText := strings.Join(lines[linesToRemove:], "\r\n")

		// テキストを更新
		controlView.Console.SetText(newText)
		controlView.lineCount = MaxLines

		// カーソルを最後に移動
		textLength := controlView.textLength()
		controlView.setTextSelection(textLength, textLength)
	}
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

	return len(p), nil
}
