package controller

import (
	"fmt"
	"time"

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

// 処理時間をフォーマットして出力する関数
func FormatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	// 時間がある場合は hh:mm:ss フォーマットで、無い場合は mm:ss フォーマット
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

type IMWidget interface {
	Widgets() declarative.Composite
	SetWindow(window *ControlWindow)
	SetEnabledInPlaying(playing bool)
}

type MWidgets struct {
	Widgets      []IMWidget
	onLoadedFunc func()
	window       *ControlWindow
	Position     *walk.Point
}

func (mw *MWidgets) SetWindow(window *ControlWindow) {
	mw.window = window
	for _, w := range mw.Widgets {
		w.SetWindow(window)
	}
}

func (mw *MWidgets) Window() *ControlWindow {
	return mw.window
}

func (mw *MWidgets) OnLoaded() {
	if mw.onLoadedFunc != nil {
		mw.onLoadedFunc()
	}
}

func (mw *MWidgets) SetOnLoaded(f func()) {
	mw.onLoadedFunc = f
}
