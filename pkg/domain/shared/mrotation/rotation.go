// Package mrotation は回転表現の共通契約を定義する。
// 内部表現はクォータニオン（x,y,z,w）を採用し、Euler→Quat は XYZ 順序。
package mrotation

// QuaternionComponentOrder はクォータニオンの成分順を表す。
type QuaternionComponentOrder int

const (
	QuaternionComponentOrderXYZW QuaternionComponentOrder = iota
)

// EulerOrder はオイラー角の回転順序を表す。
type EulerOrder int

const (
	EulerOrderXYZ EulerOrder = iota
)

// 回転の前提:
// - IK の角度制限は条件で Z*X*Y / X*Y*Z / Y*Z*X を使い分ける。
// - ボーンのローカル軸は子ボーン方向から生成し、固定軸があれば拘束する。
// - VMD カメラ回転は Degrees 名だが内部はラジアンとして扱う。
// - Euler の正規化は明示せず、asin/atan2 の返値範囲に従う。

// ShortestPathInterpolation は最短経路補間を行うかの規約。
const ShortestPathInterpolation = true

// Mat4ToQuaternionNegateXYZ は行列→クォータニオン変換時に XYZ の符号反転を行う規約。
const Mat4ToQuaternionNegateXYZ = true

// RotateOnIOFlipSign は入出力で回転の符号反転を行うかの規約（行わない）。
const RotateOnIOFlipSign = false
