// 指示: miu200521358
package testing

// EpsilonRange は許容誤差の範囲を表す。
type EpsilonRange struct {
	Min float64
	Max float64
}

// AllowAbsoluteTestPaths は絶対パス許容フラグ。
const AllowAbsoluteTestPaths = true

// GoldenResourcePolicy はゴールデン資産の配置方針。
const GoldenResourcePolicy = "test_resources"

// MmdReproTolerance はMMD再現の許容差。
const MmdReproTolerance = 0.03

// DEFAULT_EPSILON_RANGE は既定の許容誤差。
var DEFAULT_EPSILON_RANGE = EpsilonRange{Min: 1e-10, Max: 1e-5}
