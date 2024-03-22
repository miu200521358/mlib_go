package mwidget

import (
	"bytes"
	"sync"

	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

// ConsoleView は、walk.WidgetBaseを埋め込んでカスタムウィジェットを作成します。
type ConsoleView struct {
	walk.WidgetBase
	Console   *walk.TextEdit
	logBuffer *bytes.Buffer
	mutex     sync.Mutex
}

const ConsoleViewClass = "ConsoleView Class"

func NewConsoleView(parent walk.Container) (*ConsoleView, error) {
	lv := &ConsoleView{logBuffer: new(bytes.Buffer)}

	if err := walk.InitWidget(
		lv,
		parent,
		ConsoleViewClass,
		win.WS_DISABLED,
		win.WS_EX_CLIENTEDGE); err != nil {
		return nil, err
	}

	// テキストエディットを作成
	te, err := walk.NewTextEditWithStyle(parent, win.WS_VISIBLE|win.WS_VSCROLL|win.ES_MULTILINE)
	if err != nil {
		return nil, err
	}
	lv.Console = te
	lv.Console.SetReadOnly(true)

	return lv, nil
}

func (*ConsoleView) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return walk.NewGreedyLayoutItem()
}

// Write はio.Writerインターフェースを実装し、logの出力をConsoleViewにリダイレクトします。
func (cv *ConsoleView) Write(p []byte) (n int, err error) {
	cv.mutex.Lock()
	defer cv.mutex.Unlock()

	// ログメッセージをバッファに書き込む
	n, err = cv.logBuffer.Write(p)
	if err != nil {
		return 0, err
	}

	cv.Console.AppendText(string(p))

	return n, nil
}
