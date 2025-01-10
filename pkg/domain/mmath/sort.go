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

// ArgSort returns the indices that would sort s.
func ArgSort(s sort.Interface) []int {
	indices := make([]int, s.Len())
	sortInto(s, indices)
	return indices
}

// sortInto sorts s and populates the indices slice with the indices that would sort the input slice.
func sortInto(s sort.Interface, indices []int) {
	for i := 0; i < s.Len(); i++ {
		indices[i] = i
	}
	sort.Stable(argsorter{s: s, m: indices})
}

type argsorter struct {
	s sort.Interface
	m []int
}

func (a argsorter) Less(i, j int) bool { return a.s.Less(a.m[i], a.m[j]) }
func (a argsorter) Len() int           { return a.s.Len() }
func (a argsorter) Swap(i, j int)      { a.m[i], a.m[j] = a.m[j], a.m[i] }

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
