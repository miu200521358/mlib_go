package mmath

import (
	"math"
	"slices"
	"sort"
)

// ----- 定数 -----

const (
	// Epsilon は浮動小数点比較のための小さな値です
	EPSILON = 1e-10

	// Pi は円周率です
	PI = math.Pi

	// Deg2Rad は度からラジアンへの変換係数です
	DEG_TO_RAD = PI / 180.0

	// Rad2Deg はラジアンから度への変換係数です
	RAD_TO_DEG = 180.0 / PI
)

// ----- 変換関数 -----

// DegToRad は度をラジアンに変換します
func DegToRad(deg float64) float64 {
	return deg * DEG_TO_RAD
}

// RadToDeg はラジアンを度に変換します
func RadToDeg(rad float64) float64 {
	return rad * RAD_TO_DEG
}

// ThetaToRad はtheta値をラジアンに変換します（-1から1の範囲にクランプ）
func ThetaToRad(theta float64) float64 {
	return math.Asin(math.Max(-1.0, math.Min(1.0, theta)))
}

// ----- 比較関数 -----

// NearEquals は2つの値がepsilon以内であるかどうかを返します
func NearEquals[T Number](v, other T, epsilon float64) bool {
	return math.Abs(float64(v)-float64(other)) <= epsilon
}

// Sign は値の符号を返します（負: -1, 非負: 1）
func Sign[T Number](v T) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

// ----- クランプ関数 -----

// Clamped は値をmin～maxの範囲内にクランプします
func Clamped[T Number](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Clamped01 は値を0～1の範囲内にクランプします
func Clamped01[T Number](v T) T {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// Truncate は値が非常に小さい場合にゼロにします
func Truncate[T Number](v T, epsilon float64) T {
	if math.Abs(float64(v)) < epsilon {
		return 0
	}
	return v
}

// ----- 補間関数 -----

// Lerp は線形補間を行います
func Lerp(v1, v2, t float64) float64 {
	if t <= 0 {
		return v1
	}
	if t >= 1 {
		return v2
	}
	return v1 + (v2-v1)*t
}

// ----- 集計関数 -----

// Sum は値の合計を返します
func Sum[T Number](values []T) T {
	var sum float64
	for _, v := range values {
		sum += float64(v)
	}
	return T(sum)
}

// Mean は平均値を返します
func Mean[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += float64(v)
	}
	return sum / float64(len(values))
}

// Median は中央値を返します
func Median[T Number](values []T) T {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]T, len(values))
	copy(sorted, values)
	Sort(sorted)
	middle := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[middle-1] + sorted[middle]) / 2
	}
	return sorted[middle]
}

// Std は標準偏差を返します
func Std[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := Mean(values)
	var variance float64
	for _, v := range values {
		diff := float64(v) - mean
		variance += diff * diff
	}
	return math.Sqrt(variance / float64(len(values)))
}

// Ratio は割合を返します
func Ratio[T Number](total T, values []T) []float64 {
	ratios := make([]float64, len(values))
	for i, v := range values {
		ratios[i] = float64(v) / float64(total)
	}
	return ratios
}

// ----- 最大・最小関数 -----

// Max はスライス内の最大値を返します
func Max[T Number](arr []T) T {
	if len(arr) == 0 {
		return 0
	}
	max := arr[0]
	for _, v := range arr[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Min はスライス内の最小値を返します
func Min[T Number](arr []T) T {
	if len(arr) == 0 {
		return 0
	}
	min := arr[0]
	for _, v := range arr[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// ----- 丸め関数 -----

// Round は指定した閾値で四捨五入します
func Round(v, threshold float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	vv := v * (1 / threshold)
	return math.Round(vv) * threshold
}

// Effective は有効な値を返します（NaN/Infの場合は0）
func Effective[T Number](v T) T {
	f := float64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return v
}

// ----- ソート関数 -----

// Sort はスライスを昇順にソートします
func Sort[T Number](arr []T) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
}

// ----- ユーティリティ関数 -----

// Unique は重複を削除したスライスを返します
func Unique[T Number](values []T) []T {
	encountered := make(map[T]bool)
	result := make([]T, 0, len(values))
	for _, v := range values {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Contains はスライスに値が含まれているかを返します
func Contains[T Number](s []T, v T) bool {
	if len(s) <= 1000 {
		return slices.Contains(s, v)
	}
	set := make(map[T]bool, len(s))
	for _, item := range s {
		set[item] = true
	}
	return set[v]
}

// DeepCopy はスライスのディープコピーを返します
func DeepCopy[T Number](values []T) []T {
	copied := make([]T, len(values))
	copy(copied, values)
	return copied
}

// IntRanges は0からmaxまでの整数スライスを返します
func IntRanges(max int) []int {
	return IntRangesByStep(0, max, 1)
}

// IntRangesByStep はminからmaxまでstep刻みの整数スライスを返します
func IntRangesByStep(min, max, step int) []int {
	values := make([]int, 0, (max-min)/step+1)
	for i := min; i <= max; i += step {
		values = append(values, i)
	}
	return values
}

// ----- チェック関数 -----

// IsAllSameValues はすべての値が同じかどうかを返します
func IsAllSameValues(values []float64) bool {
	if len(values) == 0 {
		return true
	}
	first := values[0]
	for _, v := range values[1:] {
		if v != first {
			return false
		}
	}
	return true
}

// IsAlmostAllSameValues はほぼすべての値が同じかどうかを返します
func IsAlmostAllSameValues(values []float64, threshold float64) bool {
	if len(values) == 0 {
		return true
	}
	first := values[0]
	for _, v := range values[1:] {
		if !NearEquals(first, v, threshold) {
			return false
		}
	}
	return true
}

// IsPowerOfTwo は値が2の累乗かどうかを返します
func IsPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}

// ----- 変換関数 -----

// BoolToInt はboolをintに変換します（true: 1, false: 0）
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolToFlag はboolをfloat64に変換します（true: 1.0, false: -1.0）
func BoolToFlag(b bool) float64 {
	if b {
		return 1.0
	}
	return -1.0
}

// ----- 二次元配列関数 -----

// Mean2DVertical は二次元配列の列ごとの平均値を計算します
func Mean2DVertical(nums [][]float64) []float64 {
	if len(nums) == 0 || len(nums[0]) == 0 {
		return nil
	}
	vertical := make([]float64, len(nums[0]))
	for _, num := range nums {
		for i, n := range num {
			vertical[i] += n
		}
	}
	for i := range vertical {
		vertical[i] /= float64(len(nums))
	}
	return vertical
}

// Mean2DHorizontal は二次元配列の行ごとの平均値を計算します
func Mean2DHorizontal(nums [][]float64) []float64 {
	horizontal := make([]float64, len(nums))
	for i, num := range nums {
		horizontal[i] = Mean(num)
	}
	return horizontal
}

// Flatten は二次元スライスを一次元にフラット化します
func Flatten[T any](slices [][]T) []T {
	if len(slices) == 0 {
		return nil
	}
	totalLen := 0
	for _, s := range slices {
		totalLen += len(s)
	}
	flattened := make([]T, 0, totalLen)
	for _, slice := range slices {
		flattened = append(flattened, slice...)
	}
	return flattened
}
