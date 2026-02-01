// 指示: miu200521358
package audio_api

// IAudioPlayer は音声再生の共通I/Fを表す。
type IAudioPlayer interface {
	// Load は音声ファイルを読み込む。
	Load(path string) error
	// Close は音声ファイルを閉じる。
	Close() error
	// Play は再生を開始する。
	Play() error
	// Pause は再生を一時停止する。
	Pause() error
	// Stop は再生を停止する。
	Stop() error
	// Seek は再生位置を指定秒に移動する。
	Seek(seconds float64) error
	// SetVolume は音量を0-100で設定する。
	SetVolume(volume int) error
	// Volume は現在の音量を0-100で返す。
	Volume() int
	// IsLoaded は音声が読み込まれているか返す。
	IsLoaded() bool
	// IsPlaying は再生中か返す。
	IsPlaying() bool
	// Path は読み込み済みのパスを返す。
	Path() string
}
