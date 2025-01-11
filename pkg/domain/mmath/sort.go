// Package argsort implements a variant of the sort function that returns a slice of indices that would sort the array.
//
// The name comes from the popular Python numpy.Argsort function.
package mmath

import (
	"math"
	"reflect"
	"sort"
)

// SortInto sorts s and populates the indices slice with the indices that would sort the input slice.
func SortInto(s sort.Interface, indices []int) {
	for i := 0; i < s.Len(); i++ {
		indices[i] = i
	}
	sort.Stable(argsorter{s: s, m: indices})
}

// ArgSort returns the indices that would sort s.
func ArgSort(s sort.Interface) []int {
	indices := make([]int, s.Len())
	SortInto(s, indices)
	return indices
}

// SortSliceInto sorts a slice and populates the indices slice with the indices that would sort the input slice.
func SortSliceInto(slice interface{}, indices []int, less func(i, j int) bool) {
	SortInto(dyn{slice, less}, indices)
}

// SortSlice return the indices that would sort a slice given a comparison function.
func SortSlice(slice interface{}, less func(i, j int) bool) []int {
	v := reflect.ValueOf(slice)
	indices := make([]int, v.Len())
	SortSliceInto(slice, indices, less)
	return indices
}

type argsorter struct {
	s sort.Interface
	m []int
}

func (a argsorter) Less(i, j int) bool { return a.s.Less(a.m[i], a.m[j]) }
func (a argsorter) Len() int           { return a.s.Len() }
func (a argsorter) Swap(i, j int)      { a.m[i], a.m[j] = a.m[j], a.m[i] }

type dyn struct {
	slice interface{}
	less  func(i, j int) bool
}

func (d dyn) Less(i, j int) bool { return d.less(i, j) }
func (d dyn) Len() int           { return reflect.ValueOf(d.slice).Len() }
func (d dyn) Swap(i, j int)      { panic("unnecessary") }

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

// ------------------

func UniqueFloat32s(floats []float32) []float32 {
	encountered := map[float32]bool{}
	result := []float32{}

	for _, v := range floats {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func UniqueInts(ints []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for _, v := range ints {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
