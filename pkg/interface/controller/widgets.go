package controller

import (
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type IMWidget interface {
	Widgets() declarative.Composite
	SetWindow(window *ControlWindow)
	EnabledInPlaying(enable bool)
}

type MWidgets struct {
	Widgets             []IMWidget
	onLoadedFunc        func()
	onChangePlayingFunc func(playing bool)
	window              *ControlWindow
	Position            *walk.Point
}

func (mw *MWidgets) SetWindow(window *ControlWindow) {
	for _, w := range mw.Widgets {
		w.SetWindow(window)
	}
}

func (mw *MWidgets) Window() *ControlWindow {
	return mw.window
}

func (mw *MWidgets) EnabledInPlaying(enable bool) {
	for _, w := range mw.Widgets {
		w.EnabledInPlaying(enable)
	}
	if mw.onChangePlayingFunc != nil {
		mw.onChangePlayingFunc(enable)
	}
}

func (mw *MWidgets) OnLoaded() {
	if mw.onLoadedFunc != nil {
		mw.onLoadedFunc()
	}
}

func (mw *MWidgets) SetOnLoaded(f func()) {
	mw.onLoadedFunc = f
}

func (mw *MWidgets) SetOnChangePlaying(f func(playing bool)) {
	mw.onChangePlayingFunc = f
}
