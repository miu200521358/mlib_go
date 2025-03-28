package widget

import (
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type MPushButton struct {
	*walk.PushButton
	window        *controller.ControlWindow
	label         string
	minSize       declarative.Size
	maxSize       declarative.Size
	stretchFactor int
	onClicked     func(cw *controller.ControlWindow)
	tooltip       string
}

func NewMPushButton() *MPushButton {
	return new(MPushButton)
}

func (b *MPushButton) SetTooltip(tooltip string) {
	b.tooltip = tooltip
}

func (b *MPushButton) SetLabel(label string) {
	b.label = label
}

func (b *MPushButton) SetMaxSize(maxSize declarative.Size) {
	b.maxSize = maxSize
}

func (b *MPushButton) SetMinSize(minSize declarative.Size) {
	b.minSize = minSize
}

func (b *MPushButton) SetStretchFactor(stretchFactor int) {
	b.stretchFactor = stretchFactor
}

func (b *MPushButton) SetOnClicked(onClicked func(cw *controller.ControlWindow)) {
	b.onClicked = onClicked
}

func (b *MPushButton) SetEnabledInPlaying(playing bool) {
	b.PushButton.SetEnabled(!playing)
}

func (b *MPushButton) SetWindow(window *controller.ControlWindow) {
	b.window = window
}

func (b *MPushButton) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.HBox{},
		Children: []declarative.Widget{
			declarative.PushButton{
				AssignTo:      &b.PushButton,
				Text:          b.label,
				MinSize:       b.minSize,
				MaxSize:       b.maxSize,
				ToolTipText:   b.tooltip,
				StretchFactor: b.stretchFactor,
				OnClicked: func() {
					if b.onClicked != nil {
						b.onClicked(b.window)
					}
				},
			},
		},
	}
}
