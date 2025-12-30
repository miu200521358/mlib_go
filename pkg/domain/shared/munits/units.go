// Package munits は単位/スケールの共通契約を定義する。
// 長さは MMD 単位（ミクセル）をそのまま扱い、正規化は行わない。
// 角度は内部ラジアン、カメラFOVのみ度数を使用する。
package munits

// LengthUnit は長さ単位の種別。
type LengthUnit int

const (
	LengthUnitMMD LengthUnit = iota
)

// LengthScale は長さ単位のスケール倍率（MMD 単位そのまま）。
const LengthScale = 1.0

// NormalizeModelScale はモデルスケールの正規化有無。
const NormalizeModelScale = false

// AngleUnit は角度単位の種別。
type AngleUnit int

const (
	AngleUnitRadian AngleUnit = iota
	AngleUnitDegree
)

// InternalAngleUnit は内部処理の角度単位。
const InternalAngleUnit = AngleUnitRadian

// CameraFovAngleUnit はカメラFOVの角度単位。
const CameraFovAngleUnit = AngleUnitDegree

// GravityScaleForBullet は物理エンジン設定時の重力倍率（Y 軸に掛ける係数）。
const GravityScaleForBullet = 10.0

// UVFlipOnLoad はUVの上下反転を行うかの規約（行わない）。
const UVFlipOnLoad = false

// ConvertDegreesOnIO は入出力で角度の度数変換を行うかの規約（行わない）。
const ConvertDegreesOnIO = false
