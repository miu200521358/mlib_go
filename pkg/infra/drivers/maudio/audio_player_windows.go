//go:build windows
// +build windows

// 指示: miu200521358
package maudio

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"golang.org/x/sys/windows"
)

const (
	mciVolumeMin = 0
	mciVolumeMax = 1000
)

var (
	winmm              = windows.NewLazySystemDLL("winmm.dll")
	procMciSendString  = winmm.NewProc("mciSendStringW")
	procMciGetErrorMsg = winmm.NewProc("mciGetErrorStringW")
	procWaveOutSetVol  = winmm.NewProc("waveOutSetVolume")
	aliasSequence      uint64
)

// AudioPlayer はMCIを用いた音声プレイヤーを表す。
type AudioPlayer struct {
	mu      sync.Mutex
	alias   string
	path    string
	loaded  bool
	playing bool
	volume  int
}

// NewAudioPlayer は音声プレイヤーを生成する。
func NewAudioPlayer() *AudioPlayer {
	player := &AudioPlayer{volume: 50}
	return player
}

// Load は音声ファイルを読み込む。
func (p *AudioPlayer) Load(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if strings.TrimSpace(path) == "" {
		return merr.NewCommonError(merr.OsPackageErrorID, merr.ErrorKindValidate, "音楽ファイルが指定されていません", nil)
	}

	baseName := filepath.Base(path)
	if p.loaded {
		if err := p.closeLocked(); err != nil {
			return err
		}
	}

	alias := nextAlias()
	p.logVerbose("音楽ロード開始: file=%s alias=%s", baseName, alias)
	cmd := buildOpenCommand(path, alias)
	if err := sendMciCommand(cmd); err != nil {
		p.logVerbose("音楽ロード失敗: file=%s alias=%s err=%s", baseName, alias, err.Error())
		return wrapMciError("音楽ファイルの読み込みに失敗しました: "+filepath.Base(path), err)
	}

	p.logVerbose("音楽設定: time_format=milliseconds file=%s alias=%s", baseName, alias)
	if err := sendMciCommand(fmt.Sprintf("set %s time format milliseconds", alias)); err != nil {
		_ = sendMciCommand(fmt.Sprintf("close %s", alias))
		p.logVerbose("音楽設定失敗: file=%s alias=%s err=%s", baseName, alias, err.Error())
		return wrapMciError("音楽ファイルの再生設定に失敗しました: "+filepath.Base(path), err)
	}

	p.alias = alias
	p.path = path
	p.loaded = true
	p.playing = false

	if err := p.applyVolumeLocked(); err != nil {
		p.logVerbose("音量設定失敗: file=%s alias=%s volume=%d err=%s", baseName, alias, p.volume, err.Error())
		return err
	}
	p.logVerbose("音楽ロード完了: file=%s alias=%s", baseName, alias)
	return nil
}

// Close は音声ファイルを閉じる。
func (p *AudioPlayer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.loaded {
		return nil
	}
	return p.closeLocked()
}

// Play は再生を開始する。
func (p *AudioPlayer) Play() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.loaded {
		return nil
	}
	p.logVerbose("音楽再生開始: file=%s alias=%s", filepath.Base(p.path), p.alias)
	if err := sendMciCommand(fmt.Sprintf("play %s", p.alias)); err != nil {
		p.logVerbose("音楽再生失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
		return wrapMciError("音楽ファイルの再生に失敗しました: "+filepath.Base(p.path), err)
	}
	p.playing = true
	return nil
}

// Pause は再生を一時停止する。
func (p *AudioPlayer) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.loaded {
		return nil
	}
	p.logVerbose("音楽一時停止: file=%s alias=%s", filepath.Base(p.path), p.alias)
	if err := sendMciCommand(fmt.Sprintf("pause %s", p.alias)); err != nil {
		p.logVerbose("音楽一時停止失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
		return wrapMciError("音楽ファイルの一時停止に失敗しました: "+filepath.Base(p.path), err)
	}
	p.playing = false
	return nil
}

// Stop は再生を停止する。
func (p *AudioPlayer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.loaded {
		return nil
	}
	p.logVerbose("音楽停止: file=%s alias=%s", filepath.Base(p.path), p.alias)
	if err := sendMciCommand(fmt.Sprintf("stop %s", p.alias)); err != nil {
		p.logVerbose("音楽停止失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
		return wrapMciError("音楽ファイルの停止に失敗しました: "+filepath.Base(p.path), err)
	}
	if err := sendMciCommand(fmt.Sprintf("seek %s to start", p.alias)); err != nil {
		p.logVerbose("音楽シーク失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
		return wrapMciError("音楽ファイルのシークに失敗しました: "+filepath.Base(p.path), err)
	}
	p.playing = false
	return nil
}

// Seek は再生位置を指定秒に移動する。
func (p *AudioPlayer) Seek(seconds float64) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.loaded {
		return nil
	}
	wasPlaying := p.playing
	if wasPlaying {
		_ = sendMciCommand(fmt.Sprintf("stop %s", p.alias))
	}
	ms := int(seconds * 1000.0)
	if ms < 0 {
		ms = 0
	}
	p.logVerbose("音楽シーク: file=%s alias=%s ms=%d", filepath.Base(p.path), p.alias, ms)
	if err := sendMciCommand(fmt.Sprintf("seek %s to %d", p.alias, ms)); err != nil {
		p.logVerbose("音楽シーク失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
		return wrapMciError("音楽ファイルのシークに失敗しました: "+filepath.Base(p.path), err)
	}
	if wasPlaying {
		if err := sendMciCommand(fmt.Sprintf("play %s", p.alias)); err != nil {
			p.logVerbose("音楽再生失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
			return wrapMciError("音楽ファイルの再生に失敗しました: "+filepath.Base(p.path), err)
		}
		p.playing = true
	}
	return nil
}

// SetVolume は音量を0-100で設定する。
func (p *AudioPlayer) SetVolume(volume int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.volume = clampVolume(volume)
	if !p.loaded {
		return nil
	}
	p.logVerbose("音量設定: file=%s alias=%s volume=%d", filepath.Base(p.path), p.alias, p.volume)
	return p.applyVolumeLocked()
}

// Volume は現在の音量を0-100で返す。
func (p *AudioPlayer) Volume() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.volume
}

// IsLoaded は音声が読み込まれているか返す。
func (p *AudioPlayer) IsLoaded() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.loaded
}

// IsPlaying は再生中か返す。
func (p *AudioPlayer) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.playing
}

// Path は読み込み済みのパスを返す。
func (p *AudioPlayer) Path() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.path
}

// closeLocked は内部状態をロックしたままクローズする。
func (p *AudioPlayer) closeLocked() error {
	if !p.loaded {
		return nil
	}
	p.logVerbose("音楽クローズ開始: file=%s alias=%s", filepath.Base(p.path), p.alias)
	if err := sendMciCommand(fmt.Sprintf("stop %s", p.alias)); err != nil {
		p.logVerbose("音楽停止失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
		return wrapMciError("音楽ファイルの停止に失敗しました: "+filepath.Base(p.path), err)
	}
	if err := sendMciCommand(fmt.Sprintf("close %s", p.alias)); err != nil {
		p.logVerbose("音楽クローズ失敗: file=%s alias=%s err=%s", filepath.Base(p.path), p.alias, err.Error())
		return wrapMciError("音楽ファイルのクローズに失敗しました: "+filepath.Base(p.path), err)
	}
	alias := p.alias
	p.loaded = false
	p.playing = false
	p.alias = ""
	p.path = ""
	p.logVerbose("音楽クローズ完了: alias=%s", alias)
	return nil
}

// applyVolumeLocked は内部状態をロックしたまま音量を反映する。
func (p *AudioPlayer) applyVolumeLocked() error {
	volume := mciVolumeMin + (clampVolume(p.volume) * 10)
	if volume > mciVolumeMax {
		volume = mciVolumeMax
	}
	if isWaveAudio(p.path) {
		waveVolume := buildWaveOutVolume(p.volume)
		p.logVerbose("音量反映: file=%s alias=%s volume_raw=%d volume_wave=0x%08x", filepath.Base(p.path), p.alias, p.volume, waveVolume)
		if err := setWaveOutVolume(p.volume); err != nil {
			p.logVerbose("音量反映失敗: file=%s alias=%s volume_raw=%d err=%s", filepath.Base(p.path), p.alias, p.volume, err.Error())
			return err
		}
		return nil
	}
	p.logVerbose("音量反映: file=%s alias=%s volume_raw=%d volume_mci=%d", filepath.Base(p.path), p.alias, p.volume, volume)
	if err := sendMciCommand(fmt.Sprintf("setaudio %s volume to %d", p.alias, volume)); err != nil {
		p.logVerbose("音量反映失敗: file=%s alias=%s volume=%d err=%s", filepath.Base(p.path), p.alias, volume, err.Error())
		return wrapMciError("音量設定に失敗しました", err)
	}
	return nil
}

// logVerbose は冗長ログを出力する。
func (p *AudioPlayer) logVerbose(format string, args ...any) {
	logger := logging.DefaultLogger()
	if logger != nil && logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER) {
		logger.Verbose(logging.VERBOSE_INDEX_VIEWER, format, args...)
	}
}

// nextAlias はMCIの別名を生成する。
func nextAlias() string {
	seq := atomic.AddUint64(&aliasSequence, 1)
	return fmt.Sprintf("mlibAudio%04d", seq)
}

// buildOpenCommand はopenコマンド文字列を生成する。
func buildOpenCommand(path string, alias string) string {
	safePath := quoteMciPath(path)
	if isWaveAudio(path) {
		return fmt.Sprintf("open %s type waveaudio alias %s", safePath, alias)
	}
	return fmt.Sprintf("open %s type mpegvideo alias %s", safePath, alias)
}

// quoteMciPath はMCI用にパスをクォートする。
func quoteMciPath(path string) string {
	escaped := strings.ReplaceAll(path, "\"", "\\\"")
	return "\"" + escaped + "\""
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

// isWaveAudio は拡張子が wav かどうかを返す。
func isWaveAudio(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".wav"
}

// setWaveOutVolume はwaveOutSetVolumeで音量を反映する。
func setWaveOutVolume(volume int) error {
	level := buildWaveOutVolume(volume)
	ret, _, _ := procWaveOutSetVol.Call(0, uintptr(level))
	if ret == 0 {
		return nil
	}
	return merr.NewOsPackageError(fmt.Sprintf("音量設定に失敗しました: code=%d", ret), nil)
}

// buildWaveOutVolume は0-100の音量をwaveOut用の左右同一値に変換する。
func buildWaveOutVolume(volume int) uint32 {
	const waveOutVolumeMax = 0xFFFF
	level := uint32(clampVolume(volume)) * waveOutVolumeMax / 100
	return (level << 16) | level
}

// sendMciCommand はMCIコマンドを送信する。
func sendMciCommand(command string) error {
	ptr, err := windows.UTF16PtrFromString(command)
	if err != nil {
		return err
	}
	ret, _, _ := procMciSendString.Call(uintptr(unsafe.Pointer(ptr)), 0, 0, 0)
	if ret == 0 {
		return nil
	}
	return mciError(uint32(ret))
}

// mciError はMCIエラーコードを共通エラーに変換する。
func mciError(code uint32) error {
	message := mciErrorText(code)
	if message == "" {
		message = "MCIエラーが発生しました"
	}
	return merr.NewOsPackageError(message, nil)
}

// mciErrorText はMCIエラー文字列を取得する。
func mciErrorText(code uint32) string {
	buf := make([]uint16, 256)
	ret, _, _ := procMciGetErrorMsg.Call(uintptr(code), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if ret == 0 {
		return ""
	}
	return windows.UTF16ToString(buf)
}

// wrapMciError はMCIエラーを共通エラーに包む。
func wrapMciError(message string, cause error) error {
	if cause == nil {
		return merr.NewOsPackageError(message, nil)
	}
	return merr.NewOsPackageError(message, cause)
}
