// Package mnumeric は数値/行列の共通契約を定義する。
// 内部計算は float64 を基準にし、I/O は float32、GPU は float32 を使用する。
package mnumeric

// MatrixLayout は行列のレイアウト規約。
type MatrixLayout int

const (
	MatrixLayoutColumnMajor MatrixLayout = iota
)

// VectorMultiplyOrder はベクトルの掛け順規約。
type VectorMultiplyOrder int

const (
	VectorMultiplyOrderMatrixVector VectorMultiplyOrder = iota
)

// MatrixCompositionOrder は行列合成の順序規約。
type MatrixCompositionOrder int

const (
	MatrixCompositionParentLocal MatrixCompositionOrder = iota
)

// LocalTransformOrder はローカル変換の積の順序規約。
type LocalTransformOrder int

const (
	LocalTransformOrderLocalScalePositionRotation LocalTransformOrder = iota
)

// BoneTransformOrder はボーンのローカル変換合成順序。
type BoneTransformOrder int

const (
	BoneTransformOrderLocalScalePositionRotation BoneTransformOrder = iota
)

// BoneOffsetOrder はボーンのオフセット合成順序。
type BoneOffsetOrder int

const (
	BoneOffsetOrderRevertOffsetThenUnit BoneOffsetOrder = iota
)

// InternalFloatBits は内部計算の浮動小数点精度（bit）。
const InternalFloatBits = 64

// IOFloatBits は入出力の浮動小数点精度（bit）。
const IOFloatBits = 32

// GPUFloatBits はGPUの浮動小数点精度（bit）。
const GPUFloatBits = 32

// RecommendedEpsilonFine は高精度比較の目安。
const RecommendedEpsilonFine = 1e-10

// RecommendedEpsilonCoarse は粗い比較の目安。
const RecommendedEpsilonCoarse = 1e-6

// NormalizeZeroVectorKeepsValue は長さ 0/1 の正規化で元値を返す規約。
const NormalizeZeroVectorKeepsValue = true

// NormalizeZeroQuaternionToIdentity はゼロ近傍の正規化/逆行列で単位を返す規約。
const NormalizeZeroQuaternionToIdentity = true

// ReplaceInvalidNumbersOnWrite は NaN/Inf を書き込み時に置換する規約。
const ReplaceInvalidNumbersOnWrite = true

// AllowInvalidNumbersOnRead は読み込み時に NaN/Inf をそのまま取り込む規約。
const AllowInvalidNumbersOnRead = true

// DeterminismRequiresBitwiseMatch は bitwise 一致を要求するかの規約（要求しない）。
const DeterminismRequiresBitwiseMatch = false

// 行列の前提:
// - 直交性の検証は行わない。
// - IsIdent/NearEquals は近似比較。

// LocalMatrixUsesGlobalThenOffset はローカル行列が Global * Offset で算出される規約。
const LocalMatrixUsesGlobalThenOffset = true

// OffsetMatrixUsesNegativePosition は Offset が -Position を使う規約。
const OffsetMatrixUsesNegativePosition = true
