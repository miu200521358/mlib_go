//go:build windows
// +build windows

// 指示: miu200521358
package viewer

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/state"
	"github.com/miu200521358/win"
)

func TestHandleWindowFocusDisabledDiscardsPendingNotifications(t *testing.T) {
	shared := state.NewSharedState(1)
	shared.SetFocusLinkEnabled(false)
	shared.SetControlWindowFocused(true)
	shared.SetViewerWindowFocused(0, true)
	setWindowPosCalls := 0
	manager := &ViewerManager{
		shared: shared,
		setWindowPos: func(
			hwnd, insertAfter win.HWND,
			x, y, width, height int32,
			flags uint32,
		) bool {
			setWindowPosCalls++
			return true
		},
	}

	manager.handleWindowFocus()
	if setWindowPosCalls != 0 {
		t.Fatalf("無効中に SetWindowPos が呼ばれました: %d 回", setWindowPosCalls)
	}
	if shared.IsControlWindowFocused() || shared.IsViewerWindowFocused(0) {
		t.Fatalf("無効中のフォーカス通知が残っています")
	}

	shared.SetFocusLinkEnabled(true)
	manager.handleWindowFocus()
	if setWindowPosCalls != 0 {
		t.Fatalf("再有効化後に古い通知で SetWindowPos が呼ばれました: %d 回", setWindowPosCalls)
	}
}
