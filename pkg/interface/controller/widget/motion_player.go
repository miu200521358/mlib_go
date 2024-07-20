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
	appState    state.IAppState  // アプリ状態
	prevFrame   int              // 前回フレーム
	playing     bool             // 再生中かどうか
	frameEdit   *walk.NumberEdit // フレーム番号入力欄
	frameSlider *walk.Slider     // フレームスライダー
	playButton  *walk.PushButton // 一時停止ボタン
}

const MotionPlayerClass = "MotionPlayer Class"

func NewMotionPlayer(
	parent walk.Container,
	appState state.IAppState,
) *MotionPlayer {
	mp := new(MotionPlayer)
	mp.appState = appState

	if err := walk.InitWidget(
		mp,
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
	mp.frameEdit, err = walk.NewNumberEdit(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	mp.frameEdit.SetDecimals(0)
	mp.frameEdit.SetRange(0, 1)
	mp.frameEdit.SetValue(0)
	mp.frameEdit.SetIncrement(1)
	mp.frameEdit.SetSpinButtonsVisible(true)
	mp.frameEdit.ValueChanged().Attach(func() {
		if !mp.Playing() {
			mp.appState.SetFrame(mp.frameEdit.Value())
		}
	})

	// フレームスライダー
	mp.frameSlider, err = walk.NewSlider(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	mp.frameSlider.SetRange(0, 1)
	mp.frameSlider.SetValue(0)
	mp.frameSlider.ValueChanged().Attach(func() {
		if !mp.Playing() {
			mp.appState.SetFrame(float64(mp.frameSlider.Value()))
		}
	})

	mp.playButton, err = walk.NewPushButton(playerComposite)
	if err != nil {
		RaiseError(err)
	}
	mp.playButton.SetText(mi18n.T("再生"))
	mp.playButton.Clicked().Attach(func() {
		mp.TriggerPlay(!mp.Playing())
	})

	// レイアウト
	layout.SetStretchFactor(mp.frameEdit, 3)
	layout.SetStretchFactor(mp.frameSlider, 20)
	layout.SetStretchFactor(mp.playButton, 2)

	return mp
}

func (mp *MotionPlayer) Dispose() {
	mp.WidgetBase.Dispose()
	mp.frameEdit.Dispose()
	mp.frameSlider.Dispose()
	mp.playButton.Dispose()
}

func (mp *MotionPlayer) PrevFrame() int {
	return mp.prevFrame
}

func (mp *MotionPlayer) SetPrevFrame(v int) {
	mp.prevFrame = v
}

func (mp *MotionPlayer) Frame() float64 {
	return mp.frameEdit.Value()
}

func (mp *MotionPlayer) SetFrame(v float64) {
	if mp.playing && v > mp.frameEdit.MaxValue() {
		v = 0
	}
	value := mmath.ClampedFloat(v, mp.frameEdit.MinValue(), mp.frameEdit.MaxValue())
	mp.frameEdit.ChangeValue(value)
	mp.frameSlider.ChangeValue(int(value))
}

func (mp *MotionPlayer) MaxFrame() int {
	return mp.frameSlider.MaxValue()
}

func (mp *MotionPlayer) SetMaxFrame(max int) {
	mp.frameEdit.SetRange(mp.frameEdit.MinValue(), float64(max))
	mp.frameSlider.SetRange(int(mp.frameEdit.MinValue()), max)
}

func (mp *MotionPlayer) UpdateMaxFrame(max int) {
	nowMax := mp.frameSlider.MaxValue()
	if nowMax < max {
		mp.SetMaxFrame(max)
	}
}

func (mp *MotionPlayer) SetRange(min, max int) {
	mp.frameEdit.SetRange(float64(min), float64(max))
	mp.frameSlider.SetRange(min, max)
}

func (mp *MotionPlayer) SetEnabled(enabled bool) {
	mp.frameEdit.SetEnabled(enabled)
	mp.frameSlider.SetEnabled(enabled)
	mp.playButton.SetEnabled(enabled)
}

func (mp *MotionPlayer) Playing() bool {
	return mp.playing
}

func (mp *MotionPlayer) TriggerPlay(playing bool) {
	mp.playing = playing

	if playing {
		mp.playButton.SetText(mi18n.T("一時停止"))
		mp.SetEnabled(false)
		mp.playButton.SetEnabled(true)
	} else {
		mp.playButton.SetText(mi18n.T("再生"))
		mp.SetEnabled(true)
	}
}

// --------------------------------

func (*MotionPlayer) CreateLayoutItem(ctx *walk.LayoutContext) walk.LayoutItem {
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
