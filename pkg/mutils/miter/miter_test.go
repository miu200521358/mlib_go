package miter

import (
	"sync"
	"testing"
)

func TestIterParallel(t *testing.T) {
	// Test case 1: Serial processing
	var serialCount int
	serialFunc := func(i int) {
		serialCount++
	}
	IterParallel(10, serialFunc)
	if serialCount != 10 {
		t.Errorf("Expected serialCount to be 10, got %d", serialCount)
	}

	// Test case 2: Parallel processing
	var parallelCount int
	parallelFunc := func(i int) {
		parallelCount++
	}
	IterParallel(10, parallelFunc)
	if parallelCount != 10 {
		t.Errorf("Expected parallelCount to be 10, got %d", parallelCount)
	}
}

func TestIterParallel_BlockSize(t *testing.T) {
	// Test case 1: Serial processing with block size 1
	var serialCount int
	serialFunc := func(i int) {
		serialCount++
	}
	blockSize = 1
	IterParallel(10, serialFunc)
	if serialCount != 10 {
		t.Errorf("Expected serialCount to be 10, got %d", serialCount)
	}

	// Test case 2: Parallel processing with block size 2
	var parallelCount int
	parallelFunc := func(i int) {
		parallelCount++
	}
	blockSize = 2
	IterParallel(10, parallelFunc)
	if parallelCount != 10 {
		t.Errorf("Expected parallelCount to be 10, got %d", parallelCount)
	}
}

func TestIterParallel_Concurrency(t *testing.T) {
	// Test case 1: Verify concurrency
	var count int
	concurrentFunc := func(i int) {
		count++
	}
	blockSize = 2
	IterParallel(10, concurrentFunc)
	if count != 10 {
		t.Errorf("Expected count to be 10, got %d", count)
	}

	// Test case 2: Verify concurrent execution
	var wg sync.WaitGroup
	concurrentFunc2 := func(i int) {
		defer wg.Done()
		count++
	}
	blockSize = 2
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go concurrentFunc2(i)
	}
	wg.Wait()
	if count != 20 {
		t.Errorf("Expected count to be 20, got %d", count)
	}
}
