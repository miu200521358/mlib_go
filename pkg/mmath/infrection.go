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
		if i > 0 && math.Abs(v) > tolerance && ysPrime[i-1]*v < 0 {
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
	inflectionIndexes := make([]int, 0)
	inflectionPoints := make(map[int]int)
	inflectionEndPoints := make(map[int]int)
	for i := range values {
		for _, iPoints := range inflectionPointsList {
			if _, ok := iPoints[i]; ok {
				maxInflectionIndex := mutils.MaxInt(inflectionIndexes)
				if maxInflectionIndex <= i {
					// 今回の開始INDEXが最大INDEXより大きい場合、そのまま追加
					inflectionPoints[i] = iPoints[i]
					inflectionEndPoints[iPoints[i]] = i

					inflectionIndexes = append(inflectionIndexes, i)
					inflectionIndexes = append(inflectionIndexes, iPoints[i])
				} else {
					// 最大INDEXより小さい場合、分割する
					inflectionMaxStart := inflectionPoints[maxInflectionIndex]
					inflectionMaxEnd := maxInflectionIndex

					vs := []int{inflectionMaxEnd, i, iPoints[i]}
					if i > 0 && inflectionMaxStart > 0 {
						vs = append(vs, inflectionMaxStart)
					} else {
						vs = append(vs, i)
					}
					mutils.SortInts(vs)

					if (vs[1]-vs[0] < 2 && vs[2]-vs[1] < 2 && vs[3]-vs[2] < 2) ||
						(vs[1]-vs[0] < 2 && vs[2]-vs[1] >= 2 && vs[3]-vs[2] < 2) {
						// 全部連続している場合、0-3を繋ぐ
						// 0,1...2,3 の場合、1,2を削除して0-3を繋ぐ
						delete(inflectionPoints, vs[1])
						delete(inflectionPoints, vs[2])
						delete(inflectionEndPoints, vs[1])
						delete(inflectionEndPoints, vs[2])

						inflectionPoints[vs[0]] = vs[3]
						inflectionEndPoints[vs[3]] = vs[0]
					} else if vs[1]-vs[0] < 2 && vs[2]-vs[1] < 2 && vs[3]-vs[2] >= 2 {
						// 0,1,2...3 の場合、1を削除して0-1,1-3を繋ぐ
						delete(inflectionPoints, vs[1])
						delete(inflectionEndPoints, vs[1])

						inflectionPoints[vs[0]] = vs[2]
						inflectionEndPoints[vs[2]] = vs[0]

						inflectionPoints[vs[2]] = vs[3]
						inflectionEndPoints[vs[3]] = vs[2]
					} else if vs[1]-vs[0] >= 2 && vs[2]-vs[1] < 2 && vs[3]-vs[2] < 2 {
						// 0...1,2,3 の場合、2を削除して0-2,2-3を繋ぐ
						delete(inflectionPoints, vs[2])
						delete(inflectionEndPoints, vs[2])

						inflectionPoints[vs[0]] = vs[1]
						inflectionEndPoints[vs[1]] = vs[0]

						inflectionPoints[vs[1]] = vs[3]
						inflectionEndPoints[vs[3]] = vs[1]
					} else {
						// 全部離れている場合、全部登録
						inflectionPoints[vs[0]] = vs[1]
						inflectionEndPoints[vs[1]] = vs[0]

						inflectionPoints[vs[1]] = vs[2]
						inflectionEndPoints[vs[2]] = vs[1]

						inflectionPoints[vs[2]] = vs[3]
						inflectionEndPoints[vs[3]] = vs[2]
					}

					inflectionIndexes = append(inflectionIndexes, i)
					inflectionIndexes = append(inflectionIndexes, iPoints[i])
				}
			}
		}
	}

	return inflectionPoints
}
