package mmath

import (
	"math"
	"sort"
)

func Search[T Number](a []T, x T) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func Sort[T Number](x []T) { sort.Sort(Slices[T](x)) }

// Slices implements Interface for a []T, sorting in increasing order,
// with not-a-number (NaN) values ordered before other values.
type Slices[T Number] []T

func (x Slices[T]) Len() int { return len(x) }

// Less reports whether x[i] should be ordered before x[j], as required by the sort Interface.
// Note that floating-point comparison by itself is not a transitive relation: it does not
// report a consistent ordering for not-a-number (NaN) values.
// This implementation of Less places NaN values before any others, by using:
//
//	x[i] < x[j] || (math.IsNaN(x[i]) && !math.IsNaN(x[j]))
func (x Slices[T]) Less(i, j int) bool { return x[i] < x[j] || (isNaN(x[i]) && !isNaN(x[j])) }
func (x Slices[T]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// isNaN is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN[T Number](f T) bool {
	return math.IsNaN(float64(f))
}

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x Slices[T]) Sort() { sort.Sort(x) }

// ------------------

// ArgSort 関数
func ArgSort[T Number](slice []T) []int {
	n := len(slice)
	indexes := make([]int, n)
	for i := range indexes {
		indexes[i] = i
	}

	// インデックスを slice の値に基づいてソート
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
