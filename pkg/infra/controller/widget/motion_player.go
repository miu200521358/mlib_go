//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"time"

	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	sharedtime "github.com/miu200521358/mlib_go/pkg/shared/contracts/time"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// MotionPlayer は再生操作ウィジェットを表す。
type MotionPlayer struct {
	walk.WidgetBase
	window                *controller.ControlWindow
	frameEdit             *walk.NumberEdit
	frameSlider           *walk.Slider
	playButton            *walk.PushButton
	translator            i18n.II18n
	playingText           string
	stoppedText           string
	onEnabledInPlaying    func(playing bool)
	onChangePlayingPre    func(playing bool)
	onChangePlayingPost   func(playing bool)
	startPlayingResetType func() state.PhysicsResetType
}

// NewMotionPlayer はMotionPlayerを生成する。
func NewMotionPlayer(translator i18n.II18n) *MotionPlayer {
	player := new(MotionPlayer)
	player.translator = translator
	player.playingText = player.t("一時停止")
	player.stoppedText = player.t("再生")
	player.startPlayingResetType = func() state.PhysicsResetType {
		return state.PHYSICS_RESET_TYPE_START_FRAME
	}
	return player
}

// SetLabelTexts は再生/停止ラベルを設定する。
func (mp *MotionPlayer) SetLabelTexts(playingText, stoppedText string) {
	mp.playingText = playingText
	mp.stoppedText = stoppedText
}

// SetWindow はウィンドウ参照を設定する。
func (mp *MotionPlayer) SetWindow(window *controller.ControlWindow) {
	mp.window = window
}

// t は翻訳済み文言を返す。
func (mp *MotionPlayer) t(key string) string {
	if mp == nil || mp.translator == nil || !mp.translator.IsReady() {
		return "●●" + key + "●●"
	}
	return mp.translator.T(key)
}

// Widgets はUI構成を返す。
func (mp *MotionPlayer) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.HBox{},
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text:        mp.t("再生"),
				ToolTipText: mp.t("再生ウィジェットの使い方メッセージ"),
			},
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
					if mp.window == nil || mp.window.Playing() {
						return
					}
					mp.window.SetFrame(sharedtime.Frame(mp.frameEdit.Value()))
					mp.frameSlider.ChangeValue(int(mp.frameEdit.Value()))
				},
				ToolTipText:            mp.t("再生キーフレ説明"),
				StretchFactor:          3,
				ChangedBackgroundColor: walk.ColorWhite,
			},
			declarative.Slider{
				AssignTo:    &mp.frameSlider,
				MinValue:    0,
				MaxValue:    1,
				Orientation: walk.Horizontal,
				OnValueChanged: func() {
					if mp.window == nil || mp.window.Playing() {
						return
					}
					mp.window.SetFrame(sharedtime.Frame(mp.frameSlider.Value()))
					mp.frameEdit.ChangeValue(float64(mp.frameSlider.Value()))
				},
				ToolTipText:   mp.t("再生スライダー説明"),
				Value:         0,
				StretchFactor: 20,
			},
			declarative.PushButton{
				AssignTo: &mp.playButton,
				Text:     mp.stoppedText,
				MinSize:  declarative.Size{Width: 90, Height: 20},
				MaxSize:  declarative.Size{Width: 90, Height: 20},
				OnClicked: func() {
					playing := true
					if mp.window != nil {
						playing = !mp.window.Playing()
					}
					mp.SetPlaying(playing)
				},
				ToolTipText:   mp.t("再生ボタン説明"),
				StretchFactor: 2,
			},
		},
	}
}

// Reset は最大フレームを反映して再生UIを初期化する。
func (mp *MotionPlayer) Reset(maxFrame sharedtime.Frame) {
	mp.ChangeValue(0)
	mp.frameEdit.SetRange(0, float64(maxFrame))
	mp.frameSlider.SetRange(0, int(maxFrame))
	if mp.window != nil {
		mp.window.SetFrame(0)
		mp.window.SetMaxFrame(maxFrame)
	}
}

// SetValue は表示値を設定する。
func (mp *MotionPlayer) SetValue(frame sharedtime.Frame) {
	mp.frameEdit.SetValue(float64(frame))
	mp.frameSlider.SetValue(int(frame))
}

// ChangeValue は表示値を変更する。
func (mp *MotionPlayer) ChangeValue(frame sharedtime.Frame) {
	mp.frameEdit.ChangeValue(float64(frame))
	mp.frameSlider.ChangeValue(int(frame))
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (mp *MotionPlayer) SetEnabledInPlaying(playing bool) {
	mp.frameEdit.SetEnabled(!playing)
	mp.frameSlider.SetEnabled(!playing)
	mp.playButton.SetEnabled(true)
}

// SetEnabled はウィジェットの有効状態を設定する。
func (mp *MotionPlayer) SetEnabled(enabled bool) {
	mp.frameEdit.SetEnabled(enabled)
	mp.frameSlider.SetEnabled(enabled)
	mp.playButton.SetEnabled(enabled)
}

// SetPlaying は再生状態を更新する。
func (mp *MotionPlayer) SetPlaying(playing bool) {
	if mp.playButton != nil {
		if playing {
			mp.playButton.SetText(mp.playingText)
		} else {
			mp.playButton.SetText(mp.stoppedText)
		}
	}

	resetType := mp.GetStartPlayingResetType()
	if mp.window != nil {
		mp.window.RequestPhysicsReset(resetType)
		mp.window.SetEnabledInPlaying(playing)
		mp.window.OnChangePlayingPre(playing)
		mp.window.SetPlaying(playing)
		mp.window.OnChangePlayingPost(playing)
	}

	if playing {
		go mp.syncWhilePlaying()
	}
}

// SetOnEnabledInPlaying は再生中の有効化コールバックを設定する。
func (mp *MotionPlayer) SetOnEnabledInPlaying(f func(playing bool)) {
	mp.onEnabledInPlaying = f
}

// EnabledInPlaying は再生中の有効化コールバックを呼び出す。
func (mp *MotionPlayer) EnabledInPlaying(playing bool) {
	if mp.onEnabledInPlaying != nil {
		mp.onEnabledInPlaying(playing)
	}
}

// SetOnChangePlayingPre は再生前コールバックを設定する。
func (mp *MotionPlayer) SetOnChangePlayingPre(f func(playing bool)) {
	mp.onChangePlayingPre = f
}

// OnChangePlayingPre は再生前コールバックを実行する。
func (mp *MotionPlayer) OnChangePlayingPre(playing bool) {
	if mp.onChangePlayingPre != nil {
		mp.onChangePlayingPre(playing)
	}
}

// SetOnChangePlayingPost は再生後コールバックを設定する。
func (mp *MotionPlayer) SetOnChangePlayingPost(f func(playing bool)) {
	mp.onChangePlayingPost = f
}

// OnChangePlayingPost は再生後コールバックを実行する。
func (mp *MotionPlayer) OnChangePlayingPost(playing bool) {
	if mp.onChangePlayingPost != nil {
		mp.onChangePlayingPost(playing)
	}
}

// SetStartPlayingResetType は再生開始時リセット種別を設定する。
func (mp *MotionPlayer) SetStartPlayingResetType(f func() state.PhysicsResetType) {
	mp.startPlayingResetType = f
}

// GetStartPlayingResetType は再生開始時リセット種別を返す。
func (mp *MotionPlayer) GetStartPlayingResetType() state.PhysicsResetType {
	if mp.startPlayingResetType != nil {
		return mp.startPlayingResetType()
	}
	return state.PHYSICS_RESET_TYPE_START_FRAME
}

// syncWhilePlaying は再生中のUI同期を行う。
func (mp *MotionPlayer) syncWhilePlaying() {
	if mp.window == nil {
		return
	}
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	var prevFrame sharedtime.Frame
	for range ticker.C {
		currentFrame := mp.window.Frame()
		if currentFrame != prevFrame {
			mp.ChangeValue(currentFrame)
			prevFrame = currentFrame
		}
		if !mp.window.Playing() {
			if mp.playButton != nil && mp.playButton.Text() != mp.stoppedText {
				mp.SetPlaying(false)
			}
			break
		}
	}
}
