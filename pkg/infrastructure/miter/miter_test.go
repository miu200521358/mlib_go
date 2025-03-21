package miter

import (
	"sync"
	"testing"
)

func TestIterParallel(t *testing.T) {
	// Test case 1: Serial processing
	var serialCount int
	serialFunc := func(d, i int) error {
		serialCount++
		return nil
	}
	IterParallelByList([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 100, 1, serialFunc, nil)
	if serialCount != 10 {
		t.Errorf("Expected serialCount to be 10, got %d", serialCount)
	}

	// Test case 2: Parallel processing
	var parallelCount int
	parallelFunc := func(d, i int) error {
		parallelCount++
		return nil
	}
	IterParallelByList([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 100, 0, parallelFunc, nil)
	if parallelCount != 10 {
		t.Errorf("Expected parallelCount to be 10, got %d", parallelCount)
	}
}

func TestIterParallel_BlockSize(t *testing.T) {
	// Test case 1: Serial processing with block size 1
	var serialCount int
	serialFunc := func(d, i int) error {
		serialCount++
		return nil
	}
	IterParallelByList([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 1, 0, serialFunc, nil)
	if serialCount != 10 {
		t.Errorf("Expected serialCount to be 10, got %d", serialCount)
	}

	// Test case 2: Parallel processing with block size 2
	var parallelCount int
	parallelFunc := func(d, i int) error {
		parallelCount++
		return nil
	}
	IterParallelByList([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, 0, parallelFunc, nil)
	if parallelCount != 10 {
		t.Errorf("Expected parallelCount to be 10, got %d", parallelCount)
	}
}

func TestIterParallel_Concurrency(t *testing.T) {
	// Test case 1: Verify concurrency
	var count int
	concurrentFunc := func(d, i int) error {
		count++
		return nil
	}
	IterParallelByList([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, 0, concurrentFunc, nil)
	if count != 10 {
		t.Errorf("Expected count to be 10, got %d", count)
	}

	// Test case 2: Verify concurrent execution
	var wg sync.WaitGroup
	concurrentFunc2 := func(i int) {
		defer wg.Done()
		count++
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go concurrentFunc2(i)
	}
	wg.Wait()
	if count != 20 {
		t.Errorf("Expected count to be 20, got %d", count)
	}
}
