//go:build windows
// +build windows

package mwidget

import (
	"embed"

	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type MotionPlayer struct {
	walk.WidgetBase
	mWindow     *MWindow         // メインウィンドウ
	FrameEdit   *walk.NumberEdit // フレーム番号入力欄
	FrameSlider *walk.Slider     // フレームスライダー
	PlayButton  *walk.PushButton // 一時停止ボタン
	playing     bool             // 再生中かどうか
	OnPlay      func(bool) error // 再生/一時停止時のコールバック
}

const MotionPlayerClass = "MotionPlayer Class"

func NewMotionPlayer(parent walk.Container, mWindow *MWindow, resourceFiles embed.FS) (*MotionPlayer, error) {
	mp := new(MotionPlayer)
	mp.mWindow = mWindow

	if err := walk.InitWidget(
		mp,
		parent,
		MotionPlayerClass,
		win.WS_DISABLED,
		0); err != nil {

		return nil, err
	}

	playerComposite, err := walk.NewComposite(parent)
	if err != nil {
		return nil, err
	}
	layout := walk.NewHBoxLayout()
	layout.SetMargins(MARGIN_SMALL)
	playerComposite.SetLayout(layout)

	// 再生エリア
	titleLabel, err := walk.NewTextLabel(playerComposite)
	if err != nil {
		return nil, err
	}
	titleLabel.SetText(mi18n.T("再生"))
	titleLabel.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		mlog.IL(mi18n.T("再生ウィジェットの使い方メッセージ"))
	})

	// キーフレ番号
	mp.FrameEdit, err = walk.NewNumberEdit(playerComposite)
	if err != nil {
		return nil, err
	}
	mp.FrameEdit.SetDecimals(0)
	mp.FrameEdit.SetRange(0, 1)
	mp.FrameEdit.SetValue(0)

	// フレームスライダー
	mp.FrameSlider, err = walk.NewSlider(playerComposite)
	if err != nil {
		return nil, err
	}
	mp.FrameSlider.SetRange(0, 1)
	mp.FrameSlider.SetValue(0)

	mp.PlayButton, err = walk.NewPushButton(playerComposite)
	if err != nil {
		return nil, err
	}
	mp.PlayButton.SetText(mi18n.T("再生"))
	mp.PlayButton.Clicked().Attach(func() {
		mp.playing = !mp.playing
		mp.Play(mp.playing)
	})

	mp.FrameEdit.ValueChanged().Attach(func() {
		if !mp.playing {
			mWindow.GetMainGlWindow().SetFrame(
				float32(mp.FrameEdit.Value()) / mWindow.GetMainGlWindow().Physics.MMDFps)
			mp.FrameSlider.SetValue(int(mp.FrameEdit.Value()))
		}
	})
	mp.FrameSlider.ValueChanged().Attach(func() {
		if !mp.playing {
			mWindow.GetMainGlWindow().SetFrame(
				float32(mp.FrameSlider.Value()) / mWindow.GetMainGlWindow().Physics.MMDFps)
			mp.FrameEdit.SetValue(float64(mp.FrameSlider.Value()))
		}
	})

	// レイアウト
	layout.SetStretchFactor(mp.FrameEdit, 2)
	layout.SetStretchFactor(mp.FrameSlider, 20)
	layout.SetStretchFactor(mp.PlayButton, 2)

	return mp, nil
}

func (mp *MotionPlayer) Playing() bool {
	return mp.playing
}

func (mp *MotionPlayer) Dispose() {
	mp.WidgetBase.Dispose()
	mp.FrameEdit.Dispose()
	mp.FrameSlider.Dispose()
	mp.PlayButton.Dispose()
}

func (mp *MotionPlayer) Play(playing bool) {
	for _, glWindow := range mp.mWindow.GlWindows {
		glWindow.Play(playing)
	}
	if playing {
		mp.PlayButton.SetText(mi18n.T("一時停止"))
		mp.SetEnabled(false)
	} else {
		mp.PlayButton.SetText(mi18n.T("再生"))
		mp.SetEnabled(true)
	}

	if mp.OnPlay != nil {
		err := mp.OnPlay(playing)
		if err != nil {
			mlog.ET("再生失敗", err.Error())
			mp.Play(false)
		}
	}
}

func (mp *MotionPlayer) SetRange(min, max float64) {
	mp.FrameEdit.SetRange(min, max)
	mp.FrameSlider.SetRange(int(min), int(max))
}

func (mp *MotionPlayer) SetValue(value float64) {
	mp.FrameEdit.SetValue(value)
	mp.FrameSlider.SetValue(int(value))
}

func (mp *MotionPlayer) SetEnabled(enabled bool) {
	mp.FrameEdit.SetEnabled(enabled)
	if !enabled {
		bg, err := walk.NewSystemColorBrush(walk.SysColor3DFace)
		if err != nil {
			return
		}
		mp.FrameEdit.SetBackground(bg)
	} else {
		bg, err := walk.NewSolidColorBrush(walk.RGB(255, 255, 255))
		if err != nil {
			return
		}
		mp.FrameEdit.SetBackground(bg)
	}
	mp.FrameSlider.SetEnabled(enabled)
}

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
