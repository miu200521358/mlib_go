// 指示: miu200521358
package numeric

import "math"

// Scalar はCPU側のスカラー型。
type Scalar float64

// GpuScalar はGPU側のスカラー型。
type GpuScalar float32

// NormalizeZeroPolicy は正規化時のゼロ扱い方針。
type NormalizeZeroPolicy int

const (
	// NORMALIZE_ZERO_POLICY_KEEP は長さ0/1をそのまま返す。
	NORMALIZE_ZERO_POLICY_KEEP NormalizeZeroPolicy = iota
)

// QuaternionZeroPolicy はクォータニオンのゼロ扱い方針。
type QuaternionZeroPolicy int

const (
	// QUATERNION_ZERO_POLICY_UNIT はゼロ近傍を単位回転にする。
	QUATERNION_ZERO_POLICY_UNIT QuaternionZeroPolicy = iota
)

// NearEquals は絶対誤差のみで近似一致を判定する。
func NearEquals(a, b, eps Scalar) bool {
	return Scalar(math.Abs(float64(a-b))) <= eps
}

// Clamp は値を範囲内に丸める。
func Clamp(val, min, max Scalar) Scalar {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// IsFinite はNaN/Infを除外した有限値か判定する。
func IsFinite(val Scalar) bool {
	return !math.IsNaN(float64(val)) && !math.IsInf(float64(val), 0)
}
