package mmath

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

func gradient(values []float64) []float64 {
	result := make([]float64, len(values))
	for i := 1; i < len(values)-1; i++ {
		result[i] = (values[i+1] - values[i-1]) / 2.0
	}
	result[0] = values[1] - values[0]
	result[len(values)-1] = values[len(values)-1] - values[len(values)-2]
	return result
}

func FindInflectionPoints(values []float64, tolerance float64) map[int]int {
	ysPrime := gradient(values)

	primePoints := make(map[int]int)
	prevInflectionPoint := 0
	for i, v := range ysPrime {
		if i > 0 && math.Abs(v) > tolerance && ysPrime[i-1]*v < 0 && i-prevInflectionPoint > 2 {
			// ゼロに近しい許容値範囲外で前回と符号が変わっている場合、変曲点と見なす
			primePoints[prevInflectionPoint] = i
			prevInflectionPoint = i
		}
	}

	nonMovingPoints := make(map[int]int)
	startIdx := -1
	for i, v := range ysPrime {
		if math.Abs(v) <= tolerance {
			// ゼロに近しい許容範囲内
			if startIdx < 0 {
				// 開始INDEXが未設定の場合、設定
				startIdx = i
			} else {
				continue
			}
		} else {
			// 許容範囲外になった場合
			if startIdx >= 0 && i-startIdx > 2 {
				// 開始地点と終了地点を記録
				nonMovingPoints[startIdx] = i
				startIdx = -1
			}
		}
	}

	if startIdx > 0 && startIdx < len(ysPrime)-2 {
		// 最後に停止があった場合、最後のキーフレを保持
		nonMovingPoints[startIdx] = len(ysPrime) - 1
	}

	return MergeInflectionPoints(values, []map[int]int{primePoints, nonMovingPoints})
}

func MergeInflectionPoints(values []float64, inflectionPointsList []map[int]int) map[int]int {
	inflectionAllIndexes := make([]int, 0)
	for _, iPoints := range inflectionPointsList {
		for i, j := range iPoints {
			inflectionAllIndexes = append(inflectionAllIndexes, i)
			inflectionAllIndexes = append(inflectionAllIndexes, j)
		}
	}

	mutils.SortInts(inflectionAllIndexes)

	inflectionPoints := make(map[int]int)
	prevIdx := 0
	for i, iIdx := range inflectionAllIndexes {
		if i == 0 {
			prevIdx = iIdx
			continue
		}
		if iIdx-prevIdx > 1 {
			inflectionPoints[prevIdx] = iIdx
			prevIdx = iIdx
		}
	}

	return inflectionPoints
}
