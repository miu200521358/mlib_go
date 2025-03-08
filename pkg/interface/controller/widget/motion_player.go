package widget

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type MotionPlayer struct {
	walk.WidgetBase
	window        *controller.ControlWindow // メインウィンドウ
	frameEdit     *walk.NumberEdit          // フレーム番号入力欄
	frameSlider   *walk.Slider              // フレームスライダー
	playButton    *walk.PushButton          // 一時停止ボタン
	onTriggerPlay func(playing bool)        // 再生トリガー
}

func NewMotionPlayer() *MotionPlayer {
	player := new(MotionPlayer)
	return player
}

func (mp *MotionPlayer) SetWindow(window *controller.ControlWindow) {
	mp.window = window
}

func (mp *MotionPlayer) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.HBox{},
		Children: []declarative.Widget{
			// 再生エリア
			declarative.TextLabel{
				Text:        mi18n.T("再生"),
				ToolTipText: mi18n.T("再生ウィジェットの使い方メッセージ"),
				OnMouseDown: func(x, y int, button walk.MouseButton) {
					mlog.IL(mi18n.T("再生ウィジェットの使い方メッセージ"))
				},
			},
			// キーフレ番号
			declarative.NumberEdit{
				AssignTo:           &mp.frameEdit,
				Decimals:           0,
				MinValue:           0,
				MaxValue:           1,
				Increment:          1,
				SpinButtonsVisible: true,
				MinSize:            declarative.Size{Width: 60, Height: 20},
				MaxSize:            declarative.Size{Width: 60, Height: 20},
				OnValueChanged: func() {
					if !mp.window.Playing() {
						mp.window.SetFrame(float32(mp.frameEdit.Value()))
						mp.frameSlider.ChangeValue(int(mp.frameEdit.Value()))
					}
				},
				ToolTipText:   mi18n.T("再生キーフレ説明"),
				StretchFactor: 3,
			},
			// フレームスライダー
			declarative.Slider{
				AssignTo:    &mp.frameSlider,
				MinValue:    0,
				MaxValue:    1,
				Orientation: walk.Horizontal,
				OnValueChanged: func() {
					if !mp.window.Playing() {
						mp.window.SetFrame(float32(mp.frameSlider.Value()))
						mp.frameEdit.ChangeValue(float64(mp.frameSlider.Value()))
					}
				},
				ToolTipText:   mi18n.T("再生スライダー説明"),
				Value:         0,
				StretchFactor: 20,
			},
			// 再生ボタン
			declarative.PushButton{
				AssignTo: &mp.playButton,
				Text:     mi18n.T("再生"),
				MinSize:  declarative.Size{Width: 90, Height: 20},
				MaxSize:  declarative.Size{Width: 90, Height: 20},
				OnClicked: func() {
					mp.window.SetPlaying(!mp.window.Playing())
				},
				ToolTipText:   mi18n.T("再生ボタン説明"),
				StretchFactor: 2,
			},
		},
	}
}

func (mp *MotionPlayer) Reset(maxFrame float32) {
	mp.frameEdit.SetValue(0)
	mp.frameSlider.SetValue(0)
	mp.frameEdit.SetRange(0, float64(maxFrame))
	mp.frameSlider.SetRange(0, int(maxFrame))
	mp.window.SetFrame(0)
	mp.window.SetMaxFrame(maxFrame)
}
