package widget

import (
	"time"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

type MotionPlayer struct {
	walk.WidgetBase
	window                *controller.ControlWindow   // メインウィンドウ
	frameEdit             *walk.NumberEdit            // フレーム番号入力欄
	frameSlider           *walk.Slider                // フレームスライダー
	playButton            *walk.PushButton            // 一時停止ボタン
	playingText           string                      // 再生中のテキスト
	stoppedText           string                      // 停止中のテキスト
	onEnabledInPlaying    func(playing bool)          // 再生中に操作可能なウィジェットを有効化する
	onChangePlayingPre    func(playing bool)          // 再生前に呼ばれるコールバック
	onChangePlayingPost   func(playing bool)          // 再生後に呼ばれるコールバック
	startPlayingResetType func() vmd.PhysicsResetType // 再生開始時に設定する物理リセットタイプを返す関数
}

func NewMotionPlayer() *MotionPlayer {
	player := new(MotionPlayer)
	player.playingText = mi18n.T("一時停止")
	player.stoppedText = mi18n.T("再生")
	return player
}

func (mp *MotionPlayer) SetLabelTexts(playingText, stoppedText string) {
	mp.playingText = playingText
	mp.stoppedText = stoppedText
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
				Text:     mp.stoppedText,
				MinSize:  declarative.Size{Width: 90, Height: 20},
				MaxSize:  declarative.Size{Width: 90, Height: 20},
				OnClicked: func() {
					playing := !mp.window.Playing()
					mp.SetPlaying(playing)
				},
				ToolTipText:   mi18n.T("再生ボタン説明"),
				StretchFactor: 2,
			},
		},
	}
}

func (mp *MotionPlayer) Reset(maxFrame float32) {
	mp.ChangeValue(0.0)
	mp.frameEdit.SetRange(0, float64(maxFrame))
	mp.frameSlider.SetRange(0, int(maxFrame))
	mp.window.SetFrame(0.0)
	mp.window.SetMaxFrame(maxFrame)
}

func (mp *MotionPlayer) SetValue(frame float32) {
	mp.frameEdit.SetValue(float64(frame))
	mp.frameSlider.SetValue(int(frame))
}

func (mp *MotionPlayer) ChangeValue(frame float32) {
	mp.frameEdit.ChangeValue(float64(frame))
	mp.frameSlider.ChangeValue(int(frame))
}

func (mp *MotionPlayer) SetEnabledInPlaying(playing bool) {
	mp.frameEdit.SetEnabled(!playing)
	mp.frameSlider.SetEnabled(!playing)
	// 再生ボタンは常に有効
	mp.playButton.SetEnabled(true)
}

func (mp *MotionPlayer) SetEnabled(enabled bool) {
	mp.frameEdit.SetEnabled(enabled)
	mp.frameSlider.SetEnabled(enabled)
	mp.playButton.SetEnabled(enabled)
}

func (mp *MotionPlayer) SetPlaying(playing bool) {
	if playing {
		mp.playButton.SetText(mp.playingText)
	} else {
		mp.playButton.SetText(mp.stoppedText)
	}

	// 再生前処理
	mp.window.StorePhysicsReset(mp.GetStartPlayingResetType())
	mp.EnabledInPlaying(playing)
	mp.OnChangePlayingPre(playing)

	// 再生
	mp.window.SetPlaying(playing)

	// 再生後処理
	mp.OnChangePlayingPost(playing)

	// 再生中のみ、Ticker で定期的にフレーム情報を監視・更新する
	if playing {
		go func() {
			ticker := time.NewTicker(time.Second / 60)
			defer ticker.Stop()

			var prevFrame float32
			for range ticker.C {
				// 再生が停止されたらループを抜ける
				if !mp.window.Playing() {
					break
				}
				currentFrame := mp.window.Frame()
				// 前回のフレームと異なる場合に更新する
				if currentFrame != prevFrame {
					mp.ChangeValue(currentFrame)
					prevFrame = currentFrame
				}
			}
		}()
	}
}

func (mp *MotionPlayer) SetOnEnabledInPlaying(f func(playing bool)) {
	mp.onEnabledInPlaying = f
}

func (mp *MotionPlayer) EnabledInPlaying(playing bool) {
	if mp.onEnabledInPlaying != nil {
		mp.onEnabledInPlaying(playing)
	}
}

func (mp *MotionPlayer) SetOnChangePlayingPre(f func(playing bool)) {
	mp.onChangePlayingPre = f
}

func (mp *MotionPlayer) OnChangePlayingPre(playing bool) {
	if mp.onChangePlayingPre != nil {
		mp.onChangePlayingPre(playing)
	}
}

func (mp *MotionPlayer) SetOnChangePlayingPost(f func(playing bool)) {
	mp.onChangePlayingPost = f
}

func (mp *MotionPlayer) OnChangePlayingPost(playing bool) {
	if mp.onChangePlayingPost != nil {
		mp.onChangePlayingPost(playing)
	}
}

func (mp *MotionPlayer) SetStartPlayingResetType(f func() vmd.PhysicsResetType) {
	mp.startPlayingResetType = f
}

func (mp *MotionPlayer) GetStartPlayingResetType() vmd.PhysicsResetType {
	if mp.startPlayingResetType != nil {
		return mp.startPlayingResetType()
	}
	return vmd.PHYSICS_RESET_TYPE_START_FRAME
}
