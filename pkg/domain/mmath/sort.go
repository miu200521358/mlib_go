// 指示: miu200521358
package mmath

import (
	"math"
	"sort"
)

func Search[T Number](a []T, x T) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func Sort[T Number](x []T) { sort.Sort(Slices[T](x)) }

type Slices[T Number] []T

func (x Slices[T]) Len() int { return len(x) }

func (x Slices[T]) Less(i, j int) bool { return x[i] < x[j] || (isNaN(x[i]) && !isNaN(x[j])) }
func (x Slices[T]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func isNaN[T Number](f T) bool {
	return math.IsNaN(float64(f))
}

func (x Slices[T]) Sort() { sort.Sort(x) }

func ArgSort[T Number](slice []T) []int {
	n := len(slice)
	indexes := make([]int, n)
	for i := range indexes {
		indexes[i] = i
	}

	sort.Slice(indexes, func(i, j int) bool {
		return slice[indexes[i]] < slice[indexes[j]]
	})

	return indexes
}

func ArgMin[T Number](values []T) int {
	if len(values) == 0 {
		return -1
	}

	minValue := values[0]
	minIndex := 0
	for i, d := range values {
		if d < minValue {
			minValue = d
			minIndex = i
		}
	}
	return minIndex
}

func ArgMax[T Number](values []T) int {
	if len(values) == 0 {
		return -1
	}

	maxValue := values[0]
	maxIndex := 0
	for i, v := range values {
		if v > maxValue {
			maxValue = v
			maxIndex = i
		}
	}

	return maxIndex
}

