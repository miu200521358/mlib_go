package widget

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type ValueSliders struct {
	window  *controller.ControlWindow // メインウィンドウ
	sliders []*ValueSlider            // スライダーのリスト
}

func NewValueSliders() *ValueSliders {
	return &ValueSliders{
		sliders: make([]*ValueSlider, 0),
	}
}

func (vs *ValueSliders) AddSlider(slider *ValueSlider) {
	slider.parent = vs
	vs.sliders = append(vs.sliders, slider)
}

func (ts *ValueSliders) SetWindow(window *controller.ControlWindow) {
	ts.window = window
}

func (ts *ValueSliders) SetEnabledInPlaying(playing bool) {
	for _, slider := range ts.sliders {
		slider.valueEdit.SetEnabled(!playing)
		slider.slider.SetEnabled(!playing)
	}
}

func (ts *ValueSliders) SetEnabled(enabled bool) {
	for _, slider := range ts.sliders {
		slider.valueEdit.SetEnabled(enabled)
		slider.slider.SetEnabled(enabled)
	}
}

func (ts *ValueSliders) Widgets() declarative.Composite {
	sliderWidgets := make([]declarative.Widget, 0)
	for _, slider := range ts.sliders {
		sliderWidgets = append(sliderWidgets, slider.widgets()...)
	}

	return declarative.Composite{
		Layout:   declarative.VBox{MarginsZero: true, SpacingZero: true},
		Children: sliderWidgets,
	}
}

// -----------------------------

type ValueSlider struct {
	parent         *ValueSliders                                 // 親スライダー
	title          string                                        // タイトル
	tooltip        string                                        // ツールチップ
	valueEdit      *walk.NumberEdit                              // スライダー値入力欄
	slider         *walk.Slider                                  // スライダー
	labelColumns   int                                           // ラベルの列数
	gridColumns    int                                           // グリッドの列数
	sliderMin      float64                                       // スライダーの最小値
	sliderMax      float64                                       // スライダーの最大値
	initialValue   float64                                       // スライダーの初期値
	amplification  float64                                       // 増幅値
	decimals       int                                           // 小数点以下の桁数
	increment      float64                                       // スライダーの増分
	onValueChanged func(v float64, cw *controller.ControlWindow) // パス変更時のコールバック
}

func NewValueSlider(title, tooltip string,
	sliderMin, sliderMax, initialValue float64,
	decimals int, increment float64,
	gridColumns, labelColumns int,
	onValueChanged func(v float64, cw *controller.ControlWindow),
) *ValueSlider {
	// 範囲の差分を計算
	rangeDiff := sliderMax - sliderMin

	// 桁数を計算してamplificationを決定
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

func (slider *ValueSlider) Value() float64 {
	return slider.valueEdit.Value()
}

func (slider *ValueSlider) SetValue(v float64) {
	slider.valueEdit.SetValue(v)
	slider.onValueChanged(v, slider.parent.window)
}

func (slider *ValueSlider) widgets() []declarative.Widget {
	return []declarative.Widget{
		declarative.TextLabel{
			Text:          slider.title,
			ToolTipText:   slider.tooltip,
			OnMouseDown:   func(x, y int, button walk.MouseButton) { mlog.IL("%s", slider.tooltip) },
			StretchFactor: 2,
			Column:        0,
		},
		declarative.NumberEdit{
			AssignTo: &slider.valueEdit,
			OnValueChanged: func() {
				slider.slider.ChangeValue(int(slider.Value() * slider.amplification))
				slider.onValueChanged(slider.valueEdit.Value(), slider.parent.window)
			},
			MinSize:            declarative.Size{Width: 60, Height: 20},
			MaxSize:            declarative.Size{Width: 60, Height: 20},
			Value:              slider.initialValue,
			Decimals:           slider.decimals,  // 小数点以下の桁数
			Increment:          slider.increment, // 増分
			SpinButtonsVisible: true,             // スピンボタンを表示
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
				slider.onValueChanged(slider.valueEdit.Value(), slider.parent.window)
			},
		},
	}
}
