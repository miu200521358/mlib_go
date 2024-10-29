//go:build windows
// +build windows

package widget

import (
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type MotionPlayer struct {
	walk.WidgetBase
	controlWindow state.IControlWindow // コントローラー画面
	playing       bool                 // 再生中かどうか
	frameEdit     *walk.NumberEdit     // フレーム番号入力欄
	frameSlider   *walk.Slider         // フレームスライダー
	playButton    *walk.PushButton     // 一時停止ボタン
	onTriggerPlay func(playing bool)   // 再生トリガー
}

func NewMotionPlayer(
	parent walk.Container,
	controlWindow state.IControlWindow,
) *MotionPlayer {
	player := new(MotionPlayer)
	player.controlWindow = controlWindow

	composite := &declarative.Composite{
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
				AssignTo:           &player.frameEdit,
				Decimals:           0,
				MinValue:           0,
				MaxValue:           1,
				Increment:          1,
				SpinButtonsVisible: true,
				MinSize:            declarative.Size{Width: 60, Height: 20},
				MaxSize:            declarative.Size{Width: 60, Height: 20},
				OnValueChanged: func() {
					if !player.Playing() {
						player.controlWindow.SetFrame(float32(player.frameEdit.Value()))
					}
				},
				ToolTipText:   mi18n.T("再生キーフレ説明"),
				StretchFactor: 3,
			},
			// フレームスライダー
			declarative.Slider{
				AssignTo:    &player.frameSlider,
				MinValue:    0,
				MaxValue:    1,
				Orientation: walk.Horizontal,
				OnValueChanged: func() {
					if !player.Playing() {
						player.controlWindow.SetFrame(float32(player.frameSlider.Value()))
					}
				},
				ToolTipText:   mi18n.T("再生スライダー説明"),
				Value:         0,
				StretchFactor: 20,
			},
			// 再生ボタン
			declarative.PushButton{
				AssignTo: &player.playButton,
				Text:     mi18n.T("再生"),
				MinSize:  declarative.Size{Width: 90, Height: 20},
				MaxSize:  declarative.Size{Width: 90, Height: 20},
				OnClicked: func() {
					player.SetPlaying(!player.Playing())
				},
				ToolTipText:   mi18n.T("再生ボタン説明"),
				StretchFactor: 2,
			},
		},
	}

	if err := composite.Create(declarative.NewBuilder(parent)); err != nil {
		RaiseError(err)
	}

	return player
}

func (player *MotionPlayer) Dispose() {
	player.WidgetBase.Dispose()
	player.frameEdit.Dispose()
	player.frameSlider.Dispose()
	player.playButton.Dispose()
}

func (player *MotionPlayer) Frame() float32 {
	return float32(player.frameEdit.Value())
}

func (player *MotionPlayer) SetFrame(frame float32) {
	value := mmath.ClampedFloat(float64(frame), player.frameEdit.MinValue(), player.frameEdit.MaxValue())
	player.frameEdit.ChangeValue(value)
	player.frameSlider.ChangeValue(int(value))
	player.controlWindow.SetFrameChannel(float32(value))
}

func (player *MotionPlayer) MaxFrame() float32 {
	return float32(player.frameSlider.MaxValue())
}

func (player *MotionPlayer) SetMaxFrame(max float32) {
	player.frameEdit.SetRange(player.frameEdit.MinValue(), float64(max))
	player.frameSlider.SetRange(int(player.frameEdit.MinValue()), int(max))
}

func (player *MotionPlayer) UpdateMaxFrame(max float32) {
	nowMax := float32(player.frameSlider.MaxValue())
	if nowMax < max {
		player.SetMaxFrame(max)
	}
}

func (player *MotionPlayer) SetRange(min, max int) {
	player.frameEdit.SetRange(float64(min), float64(max))
	player.frameSlider.SetRange(min, max)
}

func (player *MotionPlayer) Enabled() bool {
	return player.frameEdit.Enabled()
}

func (player *MotionPlayer) SetEnabled(enabled bool) {
	player.frameEdit.SetEnabled(enabled)
	player.frameSlider.SetEnabled(enabled)
}

func (player *MotionPlayer) SetEnabledPlayButton(enabled bool) {
	player.playButton.SetEnabled(enabled)
}

func (player *MotionPlayer) Playing() bool {
	return player.playing
}

func (player *MotionPlayer) SetPlaying(playing bool) {
	player.playing = playing

	if playing {
		player.playButton.SetText(mi18n.T("一時停止"))
		player.SetEnabled(false)
	} else {
		player.playButton.SetText(mi18n.T("再生"))
		player.SetEnabled(true)
	}

	player.controlWindow.SetPlayingChannel(playing)

	if player.onTriggerPlay != nil {
		player.onTriggerPlay(playing)
	}
}

func (player *MotionPlayer) SetOnTriggerPlay(f func(playing bool)) {
	player.onTriggerPlay = f
}

// --------------------------------

func (player *MotionPlayer) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
	return &motionPlayerLayoutItem{idealSize: walk.SizeFrom96DPI(walk.Size{Width: 50, Height: 50}, ctx.DPI())}
}

type motionPlayerLayoutItem struct {
	walk.LayoutItemBase
	idealSize walk.Size // in native pixels
}

func (li *motionPlayerLayoutItem) LayoutFlags() walk.LayoutFlags {
	return 0
}

func (li *motionPlayerLayoutItem) IdealSize() walk.Size {
	return li.idealSize
}
