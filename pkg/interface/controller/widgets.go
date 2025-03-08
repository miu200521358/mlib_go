package controller

import "github.com/miu200521358/walk/pkg/declarative"

type IMWidget interface {
	Widgets() declarative.Composite
	SetWindow(window *ControlWindow)
	EnabledInPlaying(enable bool)
}

type MWidgets struct {
	Widgets []IMWidget
}

func (mw *MWidgets) SetWindow(window *ControlWindow) {
	for _, w := range mw.Widgets {
		w.SetWindow(window)
	}
}

func (mw *MWidgets) EnabledInPlaying(enable bool) {
	for _, w := range mw.Widgets {
		w.EnabledInPlaying(enable)
	}
}
