// 指示: miu200521358
package mmath

import (
	"math"
	"slices"
	"sort"

	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

const (
	mathCalculateXFailedErrorID = "92103"
)

// newMathCalculateXFailed はCalculateX失敗エラーを生成する。
func newMathCalculateXFailed() error {
	return baseerr.NewCommonError(mathCalculateXFailedErrorID, baseerr.ErrorKindInternal, "CalculateXに失敗しました", nil)
}

// Sum は合計を返す。
func Sum[T Number](values []T) T {
	sum := 0.0
	for _, v := range values {
		sum += float64(v)
	}
	return T(sum)
}

// Ratio は比率配列を返す。
func Ratio[T Number](total T, values []T) []float64 {
	denom := float64(total)
	if denom == 0 || math.IsNaN(denom) || math.IsInf(denom, 0) {
		return make([]float64, len(values))
	}
	ratios := make([]float64, len(values))
	for i, v := range values {
		ratio := float64(v) / denom
		if math.IsNaN(ratio) || math.IsInf(ratio, 0) {
			ratio = 0
		}
		ratios[i] = ratio
	}
	return ratios
}

// Effective は無効値を0に補正する。
func Effective[T Number](v T) T {
	if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
		return 0
	}
	return v
}

// Unique は重複を除いた配列を返す。
func Unique[T comparable](values []T) []T {
	encountered := map[T]bool{}
	result := make([]T, 0, len(values))
	for _, v := range values {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Mean は平均を返す。
func Mean[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += float64(v)
	}
	if math.IsNaN(sum) || math.IsInf(sum, 0) {
		return 0
	}
	result := sum / float64(len(values))
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0
	}
	return result
}

// Median は中央値を返す。
func Median[T Number](values []T) T {
	if len(values) == 0 {
		var zero T
		return zero
	}
	sorted := make([]T, len(values))
	copy(sorted, values)
	Sort(sorted)
	middle := len(sorted) / 2
	if len(sorted)%2 == 0 {
		result := (sorted[middle-1] + sorted[middle]) / 2
		if isNaN(result) {
			var zero T
			return zero
		}
		return result
	}
	result := sorted[middle]
	if isNaN(result) {
		var zero T
		return zero
	}
	return result
}

// Std は標準偏差を返す。
func Std[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := T(Mean(values))
	variance := 0.0
	for _, num := range values {
		variance += math.Pow(float64(num-mean), 2)
	}
	result := math.Sqrt(variance / float64(len(values)))
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0
	}
	return result
}

// Lerp は線形補間する。
func Lerp(v1, v2, t float64) float64 {
	if t <= 0 {
		return v1
	}
	if t >= 1 {
		return v2
	}
	return v1 + ((v2 - v1) * t)
}

// Sign は符号を返す。
func Sign[T Number](v T) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

// NearEquals は誤差内で等しいか判定する。
func NearEquals[T Number](v, other T, epsilon float64) bool {
	return math.Abs(float64(v)-float64(other)) <= epsilon
}

// DegToRad は度からラジアンに変換する。
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// RadToDeg はラジアンから度に変換する。
func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

// ThetaToRad はシータ値をラジアンに変換する。
func ThetaToRad(theta float64) float64 {
	return math.Asin(math.Max(-1.0, math.Min(1.0, theta)))
}

// Clamped は範囲内に収めた値を返す。
func Clamped[T Number](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Clamped01 は0〜1に収めた値を返す。
func Clamped01[T Number](v T) T {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// Contains は含まれているか判定する。
func Contains[T comparable](values []T, v T) bool {
	if len(values) <= 1000 {
		return slices.Contains(values, v)
	}
	set := make(map[T]bool, len(values))
	for _, value := range values {
		set[value] = true
	}
	_, exists := set[v]
	return exists
}

// Max は最大値を返す。
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

// Min は最小値を返す。
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

// IntRanges は0からの範囲配列を返す。
func IntRanges(max int) []int {
	return IntRangesByStep(0, max, 1)
}

// IntRangesByStep は刻み幅付きの範囲配列を返す。
func IntRangesByStep(min, max, step int) []int {
	values := make([]int, 0, int(max/step)+step)
	for i := min; i <= max; i += step {
		values = append(values, i)
	}
	return values
}

// Mean2DVertical は縦方向の平均を返す。
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
	for i, n := range vertical {
		val := n / float64(len(nums))
		if math.IsNaN(val) || math.IsInf(val, 0) {
			val = 0
		}
		vertical[i] = val
	}
	return vertical
}

// Mean2DHorizontal は横方向の平均を返す。
func Mean2DHorizontal(nums [][]float64) []float64 {
	if len(nums) == 0 {
		return nil
	}
	horizontal := make([]float64, len(nums))
	for i, num := range nums {
		val := Mean(num)
		if math.IsNaN(val) || math.IsInf(val, 0) {
			val = 0
		}
		horizontal[i] = val
	}
	return horizontal
}

// ClampIfVerySmall は微小値を0に丸める。
func ClampIfVerySmall[T Number](v T) T {
	epsilon := 1e-6
	if math.Abs(float64(v)) < epsilon {
		return 0
	}
	return v
}

// Round は丸める。
func Round(v, threshold float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	vv := v * (1 / threshold)
	return math.Round(vv) * threshold
}

// IsAllSameValues は全要素が同じか判定する。
func IsAllSameValues(values []float64) bool {
	for n := range values {
		if values[0] != values[n] {
			return false
		}
	}
	return true
}

// IsAlmostAllSameValues は誤差内で同一か判定する。
func IsAlmostAllSameValues(values []float64, threshold float64) bool {
	for n := range values {
		if !NearEquals(values[0], values[n], threshold) {
			return false
		}
	}
	return true
}

// DeepCopy はスライスを深いコピーで複製する。
func DeepCopy[T any](values []T) []T {
	copied := make([]T, len(values))
	copy(copied, values)
	return copied
}

// IsPowerOfTwo は2の累乗か判定する。
func IsPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}

// BoolToInt はboolをintに変換する。
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolToFlag はboolをフラグ値に変換する。
func BoolToFlag(b bool) float64 {
	if b {
		return 1.0
	}
	return -1.0
}

// CalculateX は三平方からXを計算する。
func CalculateX(length, y, z float64) (float64, error) {
	squareTerm := length*length - y*y - z*z
	if squareTerm < 0 {
		return 0, newMathCalculateXFailed()
	}
	return math.Sqrt(squareTerm), nil
}

// Flatten は2次元スライスを平坦化する。
func Flatten[T any](slices2 [][]T) []T {
	if len(slices2) == 0 {
		return nil
	}
	total := 0
	for _, slice := range slices2 {
		total += len(slice)
	}
	flattened := make([]T, 0, total)
	for _, slice := range slices2 {
		flattened = append(flattened, slice...)
	}
	return flattened
}

// Sort は昇順に並び替える。
func Sort[T Number](values []T) {
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j] || (isNaN(values[i]) && !isNaN(values[j]))
	})
}

// isNaN はNaNか判定する。
func isNaN[T Number](v T) bool {
	return math.IsNaN(float64(v))
}
