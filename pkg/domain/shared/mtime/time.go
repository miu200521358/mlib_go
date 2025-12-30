// Package mtime は時間/フレームの共通契約を定義する。
// 内部の時刻はフレーム番号（float32）で保持する。
package mtime

// Frame はフレーム番号の表現。
type Frame float32

// 補間/再生の前提:
// - bone/camera は曲線補間、morph は線形補間。
// - 範囲外フレームは直前/直後を返す（wrap なし）。
// - 再生ループは frame > MaxFrame で 0 に戻す（保存時は停止）。
// - frame=0 にキーが無い場合はゼロ/単位を適用する。
// - 複数モーションのブレンドは行わない。

// DefaultFPS は既定の再生FPS。
const DefaultFPS float32 = 30.0

// DefaultSecondsPerFrame は既定の秒/フレーム。
const DefaultSecondsPerFrame float32 = 1.0 / 30.0

// PhysicsFixedTimeStepNum は物理の固定タイムステップ分割数。
const PhysicsFixedTimeStepNum = 60

// PhysicsSecondsPerFrame は物理の固定タイムステップ（秒）。
const PhysicsSecondsPerFrame float32 = 1.0 / 60.0
