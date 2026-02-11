//go:build windows
// +build windows

// 指示: miu200521358
package widget

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/mpresenter/messages"
	"time"

	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

const (
	audioVolumeDefault = 100
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
	audioPicker           *FilePicker
	audioPlayer           audio_api.IAudioPlayer
	audioPath             string
	userConfig            config.IUserConfig
	volumeEdit            *walk.NumberEdit
	volumeInitial         int
	updatingVolume        bool
}

// NewMotionPlayer はMotionPlayerを生成する。
func NewMotionPlayer(translator i18n.II18n) *MotionPlayer {
	player := new(MotionPlayer)
	player.translator = translator
	player.playingText = player.t(messages.MotionPlayerKey001)
	player.stoppedText = player.t(messages.ControlWindowKey046)
	player.volumeInitial = audioVolumeDefault
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
	if mp.audioPicker != nil {
		mp.audioPicker.SetWindow(window)
	}
}

// SetAudioPlayer は音声プレイヤーと設定を紐付ける。
func (mp *MotionPlayer) SetAudioPlayer(player audio_api.IAudioPlayer, userConfig config.IUserConfig) {
	mp.audioPlayer = player
	mp.userConfig = userConfig
	if mp.audioPlayer != nil {
		mp.audioPicker = NewAudioLoadFilePicker(
			userConfig,
			mp.translator,
			config.UserConfigKeyAudio,
			mp.t(messages.ControlWindowKey099),
			mp.t(messages.MotionPlayerKey002),
			mp.onAudioPathChanged,
		)
		if mp.window != nil {
			mp.audioPicker.SetWindow(mp.window)
		}
	}
	mp.applyVolumeFromConfig()
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
	playRow := declarative.Composite{
		Layout: declarative.HBox{},
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text:        mp.t(messages.ControlWindowKey046),
				ToolTipText: mp.t(messages.MotionPlayerKey003),
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
					frame := mtime.Frame(mp.frameEdit.Value())
					mp.applyFrameSeek(frame, "number")
					mp.frameSlider.ChangeValue(int(mp.frameEdit.Value()))
					mp.seekAudioByFrame(frame)
				},
				ToolTipText:            mp.t(messages.MotionPlayerKey004),
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
					frame := mtime.Frame(mp.frameSlider.Value())
					mp.applyFrameSeek(frame, "slider")
					mp.frameEdit.ChangeValue(float64(mp.frameSlider.Value()))
					mp.seekAudioByFrame(frame)
				},
				ToolTipText:   mp.t(messages.MotionPlayerKey005),
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
				ToolTipText:   mp.t(messages.MotionPlayerKey006),
				StretchFactor: 2,
			},
		},
	}

	children := []declarative.Widget{playRow}
	if mp.audioPicker != nil {
		audioChildren := declarative.Composite{
			Layout: declarative.HBox{
				MarginsZero: true,
				Alignment:   declarative.AlignHCenterVFar,
			},
			Children: []declarative.Widget{
				mp.volumeWidgets(),
				mp.audioPicker.Widgets(),
			},
		}
		children = append(children, audioChildren)
	}

	return declarative.Composite{
		Layout:   declarative.VBox{},
		Children: children,
	}
}

// volumeWidgets は音量ウィジェットの構成を返す。
func (mp *MotionPlayer) volumeWidgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.HBox{},
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text:        mp.t(messages.ControlWindowKey101),
				ToolTipText: mp.t(messages.MotionPlayerKey007),
			},
			declarative.NumberEdit{
				AssignTo:           &mp.volumeEdit,
				Decimals:           0,
				MinValue:           0,
				MaxValue:           100,
				Increment:          1,
				SpinButtonsVisible: true,
				MinSize:            declarative.Size{Width: 60, Height: 20},
				MaxSize:            declarative.Size{Width: 60, Height: 20},
				Value:              float64(mp.volumeInitial),
				OnValueChanged: func() {
					mp.handleVolumeChanged()
				},
				ToolTipText:            mp.t(messages.MotionPlayerKey007),
				StretchFactor:          3,
				ChangedBackgroundColor: walk.ColorWhite,
			},
		},
	}
}

// Reset は最大フレームを反映して再生UIを初期化する。
func (mp *MotionPlayer) Reset(maxFrame mtime.Frame) {
	mp.ChangeValue(0)
	mp.frameEdit.SetRange(0, float64(maxFrame))
	mp.frameSlider.SetRange(0, int(maxFrame))
	if mp.window != nil {
		mp.window.SetFrame(0)
		mp.window.SetMaxFrame(maxFrame)
	}
	mp.seekAudioByFrame(0)
}

// SetValue は表示値を設定する。
func (mp *MotionPlayer) SetValue(frame mtime.Frame) {
	mp.frameEdit.SetValue(float64(frame))
	mp.frameSlider.SetValue(int(frame))
}

// ChangeValue は表示値を変更する。
func (mp *MotionPlayer) ChangeValue(frame mtime.Frame) {
	mp.frameEdit.ChangeValue(float64(frame))
	mp.frameSlider.ChangeValue(int(frame))
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (mp *MotionPlayer) SetEnabledInPlaying(playing bool) {
	mp.frameEdit.SetEnabled(!playing)
	mp.frameSlider.SetEnabled(!playing)
	mp.playButton.SetEnabled(true)
	if mp.audioPicker != nil {
		mp.audioPicker.SetEnabledInPlaying(playing)
	}
	if mp.volumeEdit != nil {
		mp.volumeEdit.SetEnabled(true)
	}
}

// SetEnabled はウィジェットの有効状態を設定する。
func (mp *MotionPlayer) SetEnabled(enabled bool) {
	mp.frameEdit.SetEnabled(enabled)
	mp.frameSlider.SetEnabled(enabled)
	mp.playButton.SetEnabled(enabled)
	if mp.audioPicker != nil {
		mp.audioPicker.SetEnabled(enabled)
	}
	if mp.volumeEdit != nil {
		mp.volumeEdit.SetEnabled(enabled)
	}
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
		mp.startAudioPlayback(mp.currentFrame())
	} else {
		mp.pauseAudioPlayback()
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

	var prevFrame mtime.Frame
	for range ticker.C {
		currentFrame := mp.window.Frame()
		if currentFrame != prevFrame {
			mp.ChangeValue(currentFrame)
			mp.syncAudioOnLoop(currentFrame, prevFrame)
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

// applyFrameSeek は停止中のフレーム直接移動を適用して物理検証ログを出力する。
func (mp *MotionPlayer) applyFrameSeek(frame mtime.Frame, source string) {
	if mp == nil || mp.window == nil {
		return
	}
	prevFrame := mp.window.Frame()
	prevResetType := mp.window.PhysicsResetType()
	physicsEnabled := mp.window.PhysicsEnabled()
	mp.window.SetFrame(frame)
	nextResetType := mp.window.PhysicsResetType()
	mp.logPhysicsSeekVerbose(source, prevFrame, frame, prevResetType, nextResetType, physicsEnabled)
}

// logPhysicsSeekVerbose は停止中のフレーム直接移動に関する物理検証ログを出力する。
func (mp *MotionPlayer) logPhysicsSeekVerbose(
	source string,
	prevFrame mtime.Frame,
	nextFrame mtime.Frame,
	prevResetType state.PhysicsResetType,
	nextResetType state.PhysicsResetType,
	physicsEnabled bool,
) {
	logger := logging.DefaultLogger()
	if logger == nil || !logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
		return
	}
	frameDelta := nextFrame - prevFrame
	logger.Verbose(
		logging.VERBOSE_INDEX_PHYSICS,
		"物理検証シーク: source=%s prevFrame=%v nextFrame=%v frameDelta=%v physicsEnabled=%t resetBefore=%d resetAfter=%d",
		source,
		prevFrame,
		nextFrame,
		frameDelta,
		physicsEnabled,
		prevResetType,
		nextResetType,
	)
}

// currentFrame は現在フレームを取得する。
func (mp *MotionPlayer) currentFrame() mtime.Frame {
	if mp.window != nil {
		return mp.window.Frame()
	}
	if mp.frameEdit != nil {
		return mtime.Frame(mp.frameEdit.Value())
	}
	return 0
}

// startAudioPlayback は音声再生を開始する。
func (mp *MotionPlayer) startAudioPlayback(frame mtime.Frame) {
	if !mp.isAudioReady() {
		return
	}
	logger := logging.DefaultLogger()
	if err := mp.seekAudioByFrame(frame); err != nil {
		logger.Error("音声シークに失敗しました: %s", err.Error())
	}
	if err := mp.audioPlayer.Play(); err != nil {
		logger.Error("音声再生に失敗しました: %s", err.Error())
	}
}

// pauseAudioPlayback は音声再生を一時停止する。
func (mp *MotionPlayer) pauseAudioPlayback() {
	if !mp.isAudioReady() {
		return
	}
	logger := logging.DefaultLogger()
	if err := mp.audioPlayer.Pause(); err != nil {
		logger.Error("音声一時停止に失敗しました: %s", err.Error())
	}
}

// seekAudioByFrame はフレームに合わせて音声位置を調整する。
func (mp *MotionPlayer) seekAudioByFrame(frame mtime.Frame) error {
	if !mp.isAudioReady() {
		return nil
	}
	seconds := mtime.FramesToSeconds(frame, mtime.DefaultFps)
	return mp.audioPlayer.Seek(float64(seconds))
}

// syncAudioOnLoop はループ時の音声位置を補正する。
func (mp *MotionPlayer) syncAudioOnLoop(currentFrame, prevFrame mtime.Frame) {
	if !mp.isAudioReady() {
		return
	}
	if currentFrame < prevFrame {
		logger := logging.DefaultLogger()
		if err := mp.seekAudioByFrame(currentFrame); err != nil {
			logger.Error("音声シークに失敗しました: %s", err.Error())
		}
		if err := mp.audioPlayer.Play(); err != nil {
			logger.Error("音声再生に失敗しました: %s", err.Error())
		}
	}
}

// onAudioPathChanged は音声ファイル変更時の処理を行う。
func (mp *MotionPlayer) onAudioPathChanged(_ *controller.ControlWindow, _ io_common.IFileReader, path string) {
	if mp.audioPlayer == nil {
		return
	}
	if err := mp.audioPlayer.Load(path); err != nil {
		logger := logging.DefaultLogger()
		logger.Error("音楽ファイル読み込み失敗: %s", err.Error())
		controller.Beep()
		return
	}
	mp.audioPath = path
	if mp.window != nil && mp.window.Playing() {
		mp.startAudioPlayback(mp.window.Frame())
	}
}

// handleVolumeChanged は音量変更時の処理を行う。
func (mp *MotionPlayer) handleVolumeChanged() {
	if mp.volumeEdit == nil || mp.updatingVolume {
		return
	}
	volume := clampVolume(int(mp.volumeEdit.Value()))
	mp.updatingVolume = true
	mp.volumeEdit.SetValue(float64(volume))
	mp.updatingVolume = false
	if mp.audioPlayer != nil {
		if err := mp.audioPlayer.SetVolume(volume); err != nil {
			logger := logging.DefaultLogger()
			logger.Error("音量設定に失敗しました: %s", err.Error())
			controller.Beep()
		}
	}
	if mp.userConfig != nil {
		_ = mp.userConfig.SetInt(config.UserConfigKeyVolume, volume)
	}
}

// applyVolumeFromConfig はユーザー設定から音量を反映する。
func (mp *MotionPlayer) applyVolumeFromConfig() {
	volume := audioVolumeDefault
	if mp.userConfig != nil {
		if v, err := mp.userConfig.GetInt(config.UserConfigKeyVolume, audioVolumeDefault); err == nil {
			volume = v
		}
	}
	volume = clampVolume(volume)
	mp.volumeInitial = volume
	if mp.audioPlayer != nil {
		_ = mp.audioPlayer.SetVolume(volume)
	}
	if mp.volumeEdit != nil {
		mp.updatingVolume = true
		mp.volumeEdit.SetValue(float64(volume))
		mp.updatingVolume = false
	}
}

// isAudioReady は音声再生準備ができているか判定する。
func (mp *MotionPlayer) isAudioReady() bool {
	return mp.audioPlayer != nil && mp.audioPlayer.IsLoaded()
}

// clampVolume は音量値を0-100に丸める。
func clampVolume(volume int) int {
	if volume < 0 {
		return 0
	}
	if volume > 100 {
		return 100
	}
	return volume
}
