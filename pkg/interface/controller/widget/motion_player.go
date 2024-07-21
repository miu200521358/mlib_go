//go:build windows
// +build windows

package widget

import (
	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type MotionPlayer struct {
	walk.WidgetBase
	controlWindow app.IControlWindow // アプリ状態
	prevFrame     int                // 前回フレーム
	playing       bool               // 再生中かどうか
	frameEdit     *walk.NumberEdit   // フレーム番号入力欄
	frameSlider   *walk.Slider       // フレームスライダー
	playButton    *walk.PushButton   // 一時停止ボタン
}

const MotionPlayerClass = "MotionPlayer Class"

func NewMotionPlayer(
	parent walk.Container,
	controlWindow app.IControlWindow,
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

	bg, err := walk.NewSystemColorBrush(walk.SysColor3DFace)
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
			player.controlWindow.SetFrame(player.frameEdit.Value())
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
			player.controlWindow.SetFrame(float64(player.frameSlider.Value()))
		}
	})

	player.playButton, err = walk.NewPushButton(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	player.playButton.SetText(mi18n.T("再生"))
	player.playButton.Clicked().Attach(func() {
		player.TriggerPlay(!player.Playing())
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

func (player *MotionPlayer) PrevFrame() int {
	return player.prevFrame
}

func (player *MotionPlayer) SetPrevFrame(v int) {
	player.prevFrame = v
}

func (player *MotionPlayer) Frame() float64 {
	return player.frameEdit.Value()
}

func (player *MotionPlayer) SetFrame(v float64) {
	if player.playing && v > player.frameEdit.MaxValue() {
		v = 0
	}
	value := mmath.ClampedFloat(v, player.frameEdit.MinValue(), player.frameEdit.MaxValue())
	player.frameEdit.ChangeValue(value)
	player.frameSlider.ChangeValue(int(value))
}

func (player *MotionPlayer) MaxFrame() int {
	return player.frameSlider.MaxValue()
}

func (player *MotionPlayer) SetMaxFrame(max int) {
	player.frameEdit.SetRange(player.frameEdit.MinValue(), float64(max))
	player.frameSlider.SetRange(int(player.frameEdit.MinValue()), max)
}

func (player *MotionPlayer) UpdateMaxFrame(max int) {
	nowMax := player.frameSlider.MaxValue()
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

func (player *MotionPlayer) SetEnabledOnlyButton(enabled bool) {
	player.playButton.SetEnabled(enabled)
}

func (player *MotionPlayer) Playing() bool {
	return player.playing
}

func (player *MotionPlayer) TriggerPlay(playing bool) {
	player.playing = playing

	if playing {
		player.playButton.SetText(mi18n.T("一時停止"))
		player.controlWindow.SetEnabled(false)
		player.SetEnabled(false)
	} else {
		player.playButton.SetText(mi18n.T("再生"))
		player.controlWindow.SetEnabled(true)
		player.SetEnabled(true)
	}
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
