package mutils

import (
	"math"
	"sort"

)

func BoolToInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
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
