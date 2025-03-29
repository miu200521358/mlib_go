package mmath

import (
	"math"
	"slices"
)

// Sum 合計値を返します
func Sum[T Number](values []T) T {
	sum := 0.0
	for _, v := range values {
		sum += float64(v)
	}
	return T(sum)
}

// Ratio 割合を返します
func Ratio[T Number](total T, values []T) []float64 {
	ratios := make([]float64, len(values))
	for i, v := range values {
		ratios[i] = float64(v) / float64(total)
	}
	return ratios
}

// Effective 有効数字を返します
func Effective[T Number](v T) T {
	if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
		return 0
	}
	return v
}

// Unique 重複を削除したスライスを返します
func Unique[T Number](values []T) []T {
	encountered := map[T]bool{}
	result := []T{}

	for _, v := range values {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

// Mean 平均値を返します
func Mean[T Number](values []T) float64 {
	sum := 0.0
	for _, v := range values {
		sum += float64(v)
	}
	return sum / float64(len(values))
}

// Median 中央値計算
func Median[T Number](values []T) T {
	sorted := make([]T, len(values))
	copy(sorted, values)

	Sort(sorted)
	middle := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[middle-1] + sorted[middle]) / 2
	} else {
		return sorted[middle]
	}
}

// Std 標準偏差
func Std[T Number](values []T) float64 {
	mean := T(Mean(values))
	variance := 0.0
	for _, num := range values {
		variance += math.Pow(float64(num-mean), 2)
	}
	return math.Sqrt(variance / float64(len(values)))
}

// Lerp 線形補間
func Lerp(v1, v2 float64, t float64) float64 {
	if t <= 0 {
		return v1
	} else if t >= 1 {
		return v2
	}
	return v1 + ((v2 - v1) * t)
}

// Sign 符号を返します
func Sign[T Number](v T) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

// NearEquals 2つの値がepsilon以内であるかどうかを返します
func NearEquals[T Number](v T, other T, epsilon float64) bool {
	return math.Abs(float64(v)-float64(other)) <= epsilon
}

// Utility functions to convert between degrees and radians
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func ThetaToRad(theta float64) float64 {
	return math.Asin(math.Max(-1.0, math.Min(1.0, theta)))
}

// Clamped ベクトルの各要素をmin～maxの範囲内にクランプします
func Clamped[T Number](v T, min T, max T) T {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

// Clamped01 ベクトルの各要素を0～1の範囲内にクランプします
func Clamped01[T Number](v T) T {
	if v < 0 {
		return 0
	} else if v > 1 {
		return 1
	}
	return v
}

// Contains slices.Contains の高速版
func Contains[T Number](s []T, v T) bool {
	if len(s) <= 1000 {
		return slices.Contains(s, v)
	}

	set := make(map[T]bool, len(s))
	for _, s := range s {
		set[s] = true
	}

	_, exists := set[v]
	return exists
}

// Max 最大値を返します
func Max[T Number](arr []T) T {
	if len(arr) == 0 {
		return 0
	}

	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

// Min 最小値を返します
func Min[T Number](arr []T) T {
	if len(arr) == 0 {
		return 0
	}

	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return min
}

// IntRanges 0からmax-1までの整数スライスを返します
func IntRanges(max int) []int {
	return IntRangesByStep(0, max, 1)
}

// IntRanges 0からmax-1までの整数スライスを返します
func IntRangesByStep(min, max, step int) []int {
	ranges := make([]int, max)
	for i := min; i < max; i += step {
		if i >= max {
			break
		}
		ranges[i] = i
	}
	return ranges
}

// 二次元配列の平均値計算(列基準)
func Mean2DVertical(nums [][]float64) []float64 {
	vertical := make([]float64, len(nums[0]))
	for _, num := range nums {
		for i, n := range num {
			vertical[i] += n
		}
	}
	for i, n := range vertical {
		vertical[i] = n / float64(len(nums))
	}
	return vertical
}

// 二次元配列の平均値計算(行基準)
func Mean2DHorizontal(nums [][]float64) []float64 {
	horizontal := make([]float64, len(nums))
	for i, num := range nums {
		horizontal[i] = Mean(num)
	}
	return horizontal
}

// ClampIfVerySmall ベクトルの各要素がとても小さい場合、ゼロを設定する
func ClampIfVerySmall[T Number](v T) T {
	epsilon := 1e-6
	if math.Abs(float64(v)) < epsilon {
		return 0
	}

	return v
}

// Round 四捨五入
func Round(v, threshold float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}

	vv := v * (1 / threshold)
	return math.Round(vv) * threshold
}

// IsAllSameValues すべての値が同じかどうかを返します
func IsAllSameValues(values []float64) bool {
	// すべて同じ値の場合、線形補間になる
	for n := range values {
		if values[0] != values[n] {
			// 厳密にイコールであることをチェックする
			return false
		}
	}
	return true
}

// IsAlmostAllSameValues ほぼすべての値が同じかどうかを返します
func IsAlmostAllSameValues(values []float64, threshold float64) bool {
	// すべて同じ値の場合、線形補間になる
	for n := range values {
		if !NearEquals(values[0], values[n], threshold) {
			// 大体イコールであることをチェックする
			return false
		}
	}
	return true
}

func DeepCopy[T Number](values []T) []T {
	copied := make([]T, len(values))
	copy(copied, values)
	return copied
}

func IsPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BoolToFlag(b bool) float64 {
	if b {
		return 1.0
	}
	return -1.0
}
