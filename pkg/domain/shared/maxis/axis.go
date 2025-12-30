// Package maxis は座標系/軸の共通契約を定義する。
// 内部処理は右手系（+X right / +Y up / forward -Z）を基準にし、
// I/O は MMD/DirectX の左手系で扱う。
// Viewer 表示時のみ X 反転と面のワインディング反転を行う。
package maxis

// Handedness は座標系の左右を表す。
type Handedness int

const (
	HandednessRight Handedness = iota
	HandednessLeft
)

// Axis は軸の種別を表す。
type Axis int

const (
	AxisX Axis = iota
	AxisY
	AxisZ
)

// Direction は軸方向の符号を表す。
type Direction int

const (
	DirectionPositive Direction = 1
	DirectionNegative Direction = -1
)

// AxisDirection は軸と向きをまとめた表現。
type AxisDirection struct {
	Axis      Axis
	Direction Direction
}

// CoordinateBasis は座標系の基準軸を定義する。
type CoordinateBasis struct {
	Handedness Handedness
	Up         AxisDirection
	Forward    AxisDirection
	Right      AxisDirection
}

// RotationHandedness は回転の左右規約を軸ごとに表す。
// I/O の回転は X が左手回り、Y/Z が右手回り。
type RotationHandedness struct {
	X Handedness
	Y Handedness
	Z Handedness
}

// INTERNAL_BASIS は内部処理の基準座標系。
var INTERNAL_BASIS = CoordinateBasis{
	Handedness: HandednessRight,
	Up:         AxisDirection{Axis: AxisY, Direction: DirectionPositive},
	Forward:    AxisDirection{Axis: AxisZ, Direction: DirectionNegative},
	Right:      AxisDirection{Axis: AxisX, Direction: DirectionPositive},
}

// IO_BASIS は入力/出力の基準座標系（MMD/DirectX）。
var IO_BASIS = CoordinateBasis{
	Handedness: HandednessLeft,
	Up:         AxisDirection{Axis: AxisY, Direction: DirectionPositive},
	Forward:    AxisDirection{Axis: AxisZ, Direction: DirectionNegative},
	Right:      AxisDirection{Axis: AxisX, Direction: DirectionPositive},
}

// IO_ROTATION_HANDEDNESS は入力/出力の回転規約。
var IO_ROTATION_HANDEDNESS = RotationHandedness{
	X: HandednessLeft,
	Y: HandednessRight,
	Z: HandednessRight,
}

// ViewerFlipX は Viewer 表示時に X を反転する規約。
const ViewerFlipX = true

// ViewerReverseWinding は Viewer 表示時に面の表裏を維持するためのワインディング反転。
const ViewerReverseWinding = true
