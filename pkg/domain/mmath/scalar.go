package mmath

import (
	"math"
	"slices"
	"sort"
)

// 線形補間
func LerpFloat(v1, v2 float64, t float64) float64 {
	return v1 + ((v2 - v1) * t)
}

func Sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

func NearEquals(v float64, other float64, epsilon float64) bool {
	return math.Abs(v-other) <= epsilon
}

func ToRadian(degree float64) float64 {
	return degree * math.Pi / 180
}

func ToDegree(radian float64) float64 {
	return radian * 180 / math.Pi
}

// Clamp01 ベクトルの各要素をmin～maxの範囲内にクランプします
func ClampedFloat(v float64, min float64, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

// Clamp01 ベクトルの各要素をmin～maxの範囲内にクランプします
func ClampedFloat32(v float32, min float32, max float32) float32 {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

// ボーンから見た頂点ローカル位置を求める
// vertexPositions: グローバル頂点位置
// startBonePosition: 親ボーン位置
// endBonePosition: 子ボーン位置
func GetVertexLocalPositions(vertexPositions []*MVec3, startBonePosition *MVec3, endBonePosition *MVec3) []*MVec3 {
	vertexSize := len(vertexPositions)
	boneVector := endBonePosition.Sub(startBonePosition)
	boneDirection := boneVector.Normalized()

	localPositions := make([]*MVec3, vertexSize)
	for i := 0; i < vertexSize; i++ {
		vertexPosition := vertexPositions[i]
		subedVertexPosition := vertexPosition.Subed(startBonePosition)
		projection := subedVertexPosition.Project(boneDirection)
		localPosition := endBonePosition.Added(projection)
		localPositions[i] = localPosition
	}

	return localPositions
}

func ArgMin(distances []float64) int {
	minValue := math.MaxFloat64
	minIndex := -1
	for i, d := range distances {
		if d < minValue {
			minValue = d
			minIndex = i
		}
	}
	return minIndex
}

func ArgMax(distances []float64) int {
	maxValue := -math.MaxFloat64
	maxIndex := -1
	for i, d := range distances {
		if d > maxValue {
			maxValue = d
			maxIndex = i
		}
	}
	return maxIndex
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
