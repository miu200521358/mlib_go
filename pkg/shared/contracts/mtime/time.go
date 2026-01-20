// 指示: miu200521358
package mtime

// Frame はフレーム番号の型。
type Frame float32

// Seconds は秒の型。
type Seconds float32

// Fps はフレームレートの型。
type Fps float32

// FrameRange はフレーム範囲を表す。
type FrameRange struct {
	Start Frame
	End   Frame
}

// DefaultFps は既定のFPS。
const DefaultFps Fps = 30

// FramesToSeconds はフレーム数を秒に変換する。
func FramesToSeconds(frames Frame, fps Fps) Seconds {
	if fps == 0 {
		return 0
	}
	return Seconds(float32(frames) / float32(fps))
}

// SecondsToFrames は秒をフレーム数に変換する。
func SecondsToFrames(seconds Seconds, fps Fps) Frame {
	return Frame(float32(seconds) * float32(fps))
}

// FpsToSpf はFPSを秒/フレームに変換する。
func FpsToSpf(fps Fps) Seconds {
	if fps == 0 {
		return 0
	}
	return Seconds(1.0 / float32(fps))
}

// SpfToFps は秒/フレームをFPSに変換する。
func SpfToFps(spf Seconds) Fps {
	if spf == 0 {
		return 0
	}
	return Fps(1.0 / float32(spf))
}

// ClampFrame はフレームを範囲内に丸める。
func ClampFrame(frame, min, max Frame) Frame {
	if frame < min {
		return min
	}
	if frame > max {
		return max
	}
	return frame
}

// IsFrameInRange はフレームが範囲内か判定する。
func IsFrameInRange(frame Frame, r FrameRange) bool {
	return frame >= r.Start && frame <= r.End
}
