package widget

import (
	"strconv"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type TextSlider struct {
	window         *controller.ControlWindow       // メインウィンドウ
	title          string                          // タイトル
	tooltip        string                          // ツールチップ
	valueEdit      *walk.TextEdit                  // スライダー値入力欄
	slider         *walk.Slider                    // スライダー
	sliderColumns  int                             // スライダーの列数
	gridColumns    int                             // グリッドの列数
	sliderMin      float32                         // スライダーの最小値
	sliderMax      float32                         // スライダーの最大値
	initialValue   float32                         // スライダーの初期値
	onValueChanged func(*controller.ControlWindow) // パス変更時のコールバック
}

func NewTextSlider(title, tooltip string,
	sliderMin, sliderMax, initialValue float32,
	gridColumns, sliderColumns int,
	onValueChanged func(*controller.ControlWindow),
) *TextSlider {
	return &TextSlider{
		title:          title,
		tooltip:        tooltip,
		sliderMin:      sliderMin,
		sliderMax:      sliderMax,
		initialValue:   initialValue,
		onValueChanged: onValueChanged,
		sliderColumns:  sliderColumns,
		gridColumns:    gridColumns,
	}
}

func (ts *TextSlider) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.Grid{Columns: ts.gridColumns},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{MarginsZero: true},
				Children: []declarative.Widget{
					declarative.TextLabel{
						Text:        ts.title,
						ToolTipText: ts.tooltip,
						OnMouseDown: func(x, y int, button walk.MouseButton) { mlog.IL("%s", ts.tooltip) },
					},
				},
			},
			declarative.Composite{
				Layout: declarative.HBox{MarginsZero: true},
				Children: []declarative.Widget{
					declarative.TextEdit{
						AssignTo: &ts.valueEdit,
						OnTextChanged: func() {
							ts.slider.ChangeValue(int(ts.Value() * ts.sliderMax))
							ts.onValueChanged(ts.window)
						},
						MinSize: declarative.Size{Width: 40, Height: 5},
						MaxSize: declarative.Size{Width: 40, Height: 5},
						Text:    strconv.FormatFloat(float64(ts.initialValue), 'f', 2, 32),
					},
					declarative.Slider{
						AssignTo:    &ts.slider,
						ToolTipText: ts.tooltip,
						MinValue:    int(ts.sliderMin * ts.sliderMax),
						MaxValue:    int(ts.sliderMax * ts.sliderMax),
						Value:       int(ts.initialValue * ts.sliderMax),
						OnValueChanged: func() {
							v := float32(ts.slider.Value()) / ts.sliderMax
							ts.valueEdit.ChangeText(strconv.FormatFloat(float64(v), 'f', 2, 32))
							ts.onValueChanged(ts.window)
						},
					},
				},
				Column: ts.sliderColumns,
			},
		},
	}
}

func (ts *TextSlider) SetWindow(window *controller.ControlWindow) {
	ts.window = window
}

func (ts *TextSlider) SetEnabledInPlaying(playing bool) {
	ts.valueEdit.SetEnabled(!playing)
	ts.slider.SetEnabled(!playing)
}

func (ts *TextSlider) SetEnabled(enabled bool) {
	ts.valueEdit.SetEnabled(enabled)
	ts.slider.SetEnabled(enabled)
}

func (ts *TextSlider) Value() float32 {
	if ts.valueEdit.Text() == "" {
		return 0
	}
	v, err := strconv.ParseFloat(ts.valueEdit.Text(), 32)
	if err != nil {
		mlog.ET("数値変換失敗", err, "")
		return 0
	}
	return float32(v)
}

func (ts *TextSlider) SetValue(v float32) {
	ts.slider.SetValue(int(v * float32(ts.sliderMax)))
}
