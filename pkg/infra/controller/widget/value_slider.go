//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// ValueSliders はスライダー群を表す。
type ValueSliders struct {
	window  *controller.ControlWindow
	sliders []*ValueSlider
}

// NewValueSliders はValueSlidersを生成する。
func NewValueSliders() *ValueSliders {
	return &ValueSliders{sliders: make([]*ValueSlider, 0)}
}

// AddSlider はスライダーを追加する。
func (vs *ValueSliders) AddSlider(slider *ValueSlider) {
	slider.parent = vs
	vs.sliders = append(vs.sliders, slider)
}

// SetWindow はウィンドウ参照を設定する。
func (vs *ValueSliders) SetWindow(window *controller.ControlWindow) {
	vs.window = window
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (vs *ValueSliders) SetEnabledInPlaying(playing bool) {
	for _, slider := range vs.sliders {
		slider.valueEdit.SetEnabled(!playing)
		slider.slider.SetEnabled(!playing)
	}
}

// SetEnabled はウィジェットの有効状態を設定する。
func (vs *ValueSliders) SetEnabled(enabled bool) {
	for _, slider := range vs.sliders {
		slider.valueEdit.SetEnabled(enabled)
		slider.slider.SetEnabled(enabled)
	}
}

// Widgets はUI構成を返す。
func (vs *ValueSliders) Widgets() declarative.Composite {
	sliderWidgets := make([]declarative.Widget, 0)
	for _, slider := range vs.sliders {
		sliderWidgets = append(sliderWidgets, slider.widgets()...)
	}
	return declarative.Composite{
		Layout:   declarative.VBox{MarginsZero: true, SpacingZero: true},
		Children: sliderWidgets,
	}
}

// ValueSlider は単一のスライダーウィジェットを表す。
type ValueSlider struct {
	parent         *ValueSliders
	title          string
	tooltip        string
	valueEdit      *walk.NumberEdit
	slider         *walk.Slider
	labelColumns   int
	gridColumns    int
	sliderMin      float64
	sliderMax      float64
	initialValue   float64
	amplification  float64
	decimals       int
	increment      float64
	onValueChanged func(v float64, cw *controller.ControlWindow)
}

// NewValueSlider はValueSliderを生成する。
func NewValueSlider(title, tooltip string,
	sliderMin, sliderMax, initialValue float64,
	decimals int, increment float64,
	gridColumns, labelColumns int,
	onValueChanged func(v float64, cw *controller.ControlWindow),
) *ValueSlider {
	rangeDiff := sliderMax - sliderMin
	digits := int(math.Log10(float64(rangeDiff))) + 1
	amplification := float64(math.Pow10(digits))

	return &ValueSlider{
		title:          title,
		tooltip:        tooltip,
		decimals:       decimals,
		increment:      increment,
		sliderMin:      sliderMin,
		sliderMax:      sliderMax,
		initialValue:   initialValue,
		amplification:  amplification,
		onValueChanged: onValueChanged,
		labelColumns:   labelColumns,
		gridColumns:    gridColumns,
	}
}

// Value は現在値を返す。
func (slider *ValueSlider) Value() float64 {
	return slider.valueEdit.Value()
}

// SetValue は値を設定してコールバックを呼ぶ。
func (slider *ValueSlider) SetValue(v float64) {
	slider.valueEdit.SetValue(v)
	if slider.onValueChanged != nil {
		slider.onValueChanged(v, slider.parent.window)
	}
}

// widgets はスライダーUIを返す。
func (slider *ValueSlider) widgets() []declarative.Widget {
	return []declarative.Widget{
		declarative.TextLabel{
			Text:        slider.title,
			ToolTipText: slider.tooltip,
			OnMouseDown: func(x, y int, button walk.MouseButton) {
				logging.DefaultLogger().Info("%s", slider.tooltip)
			},
			StretchFactor: 2,
			Column:        0,
		},
		declarative.NumberEdit{
			AssignTo: &slider.valueEdit,
			OnValueChanged: func() {
				slider.slider.ChangeValue(int(slider.Value() * slider.amplification))
				if slider.onValueChanged != nil {
					slider.onValueChanged(slider.valueEdit.Value(), slider.parent.window)
				}
			},
			MinSize:            declarative.Size{Width: 60, Height: 20},
			MaxSize:            declarative.Size{Width: 60, Height: 20},
			Value:              slider.initialValue,
			Decimals:           slider.decimals,
			Increment:          slider.increment,
			SpinButtonsVisible: true,
			StretchFactor:      2,
			Column:             1,
		},
		declarative.Slider{
			AssignTo:      &slider.slider,
			ToolTipText:   slider.tooltip,
			Orientation:   walk.Horizontal,
			StretchFactor: 20,
			Column:        2,
			MinValue:      int(slider.sliderMin * slider.amplification),
			MaxValue:      int(slider.sliderMax * slider.amplification),
			Value:         int(slider.initialValue * slider.amplification),
			Increment:     int(slider.increment * slider.amplification),
			OnValueChanged: func() {
				v := float64(slider.slider.Value()) / slider.amplification
				slider.valueEdit.ChangeValue(v)
			},
			OnMouseUp: func(x, y int, button walk.MouseButton) {
				v := float64(slider.slider.Value()) / slider.amplification
				slider.valueEdit.ChangeValue(v)
				if slider.onValueChanged != nil {
					slider.onValueChanged(slider.valueEdit.Value(), slider.parent.window)
				}
			},
		},
	}
}
