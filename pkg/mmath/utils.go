package mmath

import (
	"math"
	"slices"
	"sort"
)

func BoolToInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// slices.Contains の高速版
func Contains[S ~[]E, E comparable](s S, v E) bool {
	if len(s) <= 20 {
		return slices.Contains(s, v)
	}

	set := make(map[E]bool, len(s))
	for _, s := range s {
		set[s] = true
	}

	_, exists := set[v]
	return exists
}

func MaxInt(arr []int) int {
	max := math.MinInt64
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

func MaxFloat(arr []float64) float64 {
	max := math.SmallestNonzeroFloat64
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

// 中央値計算
func Median(nums []float64) float64 {
	sortedNums := make([]float64, len(nums))
	copy(sortedNums, nums)
	sort.Float64s(sortedNums)
	middle := len(sortedNums) / 2
	if len(sortedNums)%2 == 0 {
		return (sortedNums[middle-1] + sortedNums[middle]) / 2
	} else {
		return sortedNums[middle]
	}
}

// 標準偏差計算
func Std(nums []float64) float64 {
	mean := Mean(nums)
	variance := 0.0
	for _, num := range nums {
		variance += math.Pow(num-mean, 2)
	}
	return math.Sqrt(variance / float64(len(nums)))
}

// 平均値計算
func Mean(nums []float64) float64 {
	total := 0.0
	for _, num := range nums {
		total += num
	}
	return total / float64(len(nums))
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

// ------------------

func SearchFloat32s(a []float32, x float32) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func SortFloat32s(x []float32) { sort.Sort(Float32Slice(x)) }

// Float32Slice implements Interface for a []float32, sorting in increasing order,
// with not-a-number (NaN) values ordered before other values.
type Float32Slice []float32

func (x Float32Slice) Len() int { return len(x) }

// Less reports whether x[i] should be ordered before x[j], as required by the sort Interface.
// Note that floating-point comparison by itself is not a transitive relation: it does not
// report a consistent ordering for not-a-number (NaN) values.
// This implementation of Less places NaN values before any others, by using:
//
//	x[i] < x[j] || (math.IsNaN(x[i]) && !math.IsNaN(x[j]))
func (x Float32Slice) Less(i, j int) bool { return x[i] < x[j] || (isNaN32(x[i]) && !isNaN32(x[j])) }
func (x Float32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// isNaN32 is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN32(f float32) bool {
	return math.IsNaN(float64(f))
}

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x Float32Slice) Sort() { sort.Sort(x) }

// ------------------

func SearchFloat64s(a []float64, x float64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func SortFloat64s(x []float64) { sort.Sort(Float64Slice(x)) }

// Float64Slice implements Interface for a []float64, sorting in increasing order,
// with not-a-number (NaN) values ordered before other values.
type Float64Slice []float64

func (x Float64Slice) Len() int { return len(x) }

// Less reports whether x[i] should be ordered before x[j], as required by the sort Interface.
// Note that floating-point comparison by itself is not a transitive relation: it does not
// report a consistent ordering for not-a-number (NaN) values.
// This implementation of Less places NaN values before any others, by using:
//
//	x[i] < x[j] || (math.IsNaN(x[i]) && !math.IsNaN(x[j]))
func (x Float64Slice) Less(i, j int) bool { return x[i] < x[j] || (isNaN64(x[i]) && !isNaN64(x[j])) }
func (x Float64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// isNaN64 is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN64(f float64) bool {
	return math.IsNaN(f)
}

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x Float64Slice) Sort() { sort.Sort(x) }

// ------------------

func SearchInts(a []int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func SortInts(x []int) { sort.Sort(IntSlice(x)) }

// IntSlice implements Interface for a []int, sorting in increasing order,
// with not-a-number (NaN) values ordered before other values.
type IntSlice []int

func (x IntSlice) Len() int { return len(x) }

// Less reports whether x[i] should be ordered before x[j], as required by the sort Interface.
// Note that floating-point comparison by itself is not a transitive relation: it does not
// report a consistent ordering for not-a-number (NaN) values.
// This implementation of Less places NaN values before any others, by using:
//
//	x[i] < x[j] || (math.IsNaN(x[i]) && !math.IsNaN(x[j]))
func (x IntSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x IntSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x IntSlice) Sort() { sort.Sort(x) }
