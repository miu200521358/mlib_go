//go:build windows
// +build windows

// 指示: miu200521358
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

// Beep は通知音を鳴らす。
func Beep() {
	procMessageBeep.Call(uintptr(MB_ICONASTERISK))
}

// FormatDuration は処理時間を hh:mm:ss / mm:ss 形式で返す。
func FormatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

// IMWidget はコントローラーの共通ウィジェットI/F。
type IMWidget interface {
	Widgets() declarative.Composite
	SetWindow(window *ControlWindow)
	SetEnabledInPlaying(playing bool)
}

// MWidgets は複数ウィジェットの集合を表す。
type MWidgets struct {
	Widgets      []IMWidget
	onLoadedFunc func()
	window       *ControlWindow
	Position     *walk.Point
}

// SetWindow は配下ウィジェットにウィンドウを設定する。
func (mw *MWidgets) SetWindow(window *ControlWindow) {
	mw.window = window
	for _, w := range mw.Widgets {
		w.SetWindow(window)
	}
}

// Window は紐付けられたウィンドウを返す。
func (mw *MWidgets) Window() *ControlWindow {
	return mw.window
}

// OnLoaded はロード完了コールバックを実行する。
func (mw *MWidgets) OnLoaded() {
	if mw.onLoadedFunc != nil {
		mw.onLoadedFunc()
	}
}

// SetOnLoaded はロード完了コールバックを設定する。
func (mw *MWidgets) SetOnLoaded(f func()) {
	mw.onLoadedFunc = f
}
