//go:build windows
// +build windows

package widget

import (
	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"

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

const MotionPlayerClass = "MotionPlayer Class"

func NewMotionPlayer(
	parent walk.Container,
	controlWindow state.IControlWindow,
) *MotionPlayer {
	player := new(MotionPlayer)
	player.controlWindow = controlWindow

	if err := walk.InitWidget(
		player,
		parent,
		MotionPlayerClass,
		win.WS_DISABLED,
		0); err != nil {
		RaiseError(err)
	}

	playerComposite, err := walk.NewComposite(parent)
	if err != nil {
		RaiseError(err)
	}
	layout := walk.NewHBoxLayout()
	playerComposite.SetLayout(layout)

	bg, err := walk.NewSystemColorBrush(walk.SysColorInactiveCaption)
	if err != nil {
		RaiseError(err)
	}
	playerComposite.SetBackground(bg)

	// 再生エリア
	titleLabel, err := walk.NewTextLabel(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	titleLabel.SetText(mi18n.T("再生"))
	titleLabel.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		mlog.IL(mi18n.T("再生ウィジェットの使い方メッセージ"))
	})

	// キーフレ番号
	player.frameEdit, err = walk.NewNumberEdit(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	player.frameEdit.SetDecimals(0)
	player.frameEdit.SetRange(0, 1)
	player.frameEdit.SetValue(0)
	player.frameEdit.SetIncrement(1)
	player.frameEdit.SetSpinButtonsVisible(true)
	player.frameEdit.ValueChanged().Attach(func() {
		if !player.Playing() {
			player.controlWindow.SetFrame(float32(player.frameEdit.Value()))
		}
	})

	// フレームスライダー
	player.frameSlider, err = walk.NewSlider(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	player.frameSlider.SetRange(0, 1)
	player.frameSlider.SetValue(0)
	player.frameSlider.ValueChanged().Attach(func() {
		if !player.Playing() {
			player.controlWindow.SetFrame(float32(player.frameSlider.Value()))
		}
	})

	player.playButton, err = walk.NewPushButton(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	player.playButton.SetText(mi18n.T("再生"))
	player.playButton.Clicked().Attach(func() {
		player.SetPlaying(!player.Playing())
	})

	// レイアウト
	layout.SetStretchFactor(player.frameEdit, 3)
	layout.SetStretchFactor(player.frameSlider, 20)
	layout.SetStretchFactor(player.playButton, 2)

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
	if player.playing && frame > float32(player.frameEdit.MaxValue()) {
		frame = 0
	}
	value := mmath.ClampedFloat(float64(frame), player.frameEdit.MinValue(), player.frameEdit.MaxValue())
	player.frameEdit.ChangeValue(value)
	player.frameSlider.ChangeValue(int(value))
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

	player.controlWindow.SetPlaying(playing)

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
