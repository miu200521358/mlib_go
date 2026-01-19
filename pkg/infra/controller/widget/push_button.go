//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// MPushButton はコールバック付きボタンを表す。
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

// NewMPushButton はMPushButtonを生成する。
func NewMPushButton() *MPushButton {
	return new(MPushButton)
}

// SetTooltip はツールチップを設定する。
func (b *MPushButton) SetTooltip(tooltip string) {
	b.tooltip = tooltip
}

// SetLabel は表示ラベルを設定する。
func (b *MPushButton) SetLabel(label string) {
	b.label = label
}

// SetMaxSize は最大サイズを設定する。
func (b *MPushButton) SetMaxSize(maxSize declarative.Size) {
	b.maxSize = maxSize
}

// SetMinSize は最小サイズを設定する。
func (b *MPushButton) SetMinSize(minSize declarative.Size) {
	b.minSize = minSize
}

// SetStretchFactor は伸長率を設定する。
func (b *MPushButton) SetStretchFactor(stretchFactor int) {
	b.stretchFactor = stretchFactor
}

// SetOnClicked はクリック時コールバックを設定する。
func (b *MPushButton) SetOnClicked(onClicked func(cw *controller.ControlWindow)) {
	b.onClicked = onClicked
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (b *MPushButton) SetEnabledInPlaying(playing bool) {
	b.PushButton.SetEnabled(!playing)
}

// SetWindow はウィンドウ参照を設定する。
func (b *MPushButton) SetWindow(window *controller.ControlWindow) {
	b.window = window
}

// Widgets はUI構成を返す。
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
