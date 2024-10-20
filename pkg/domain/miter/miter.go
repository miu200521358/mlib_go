package miter

import (
	"runtime"
	"sync"
)

// IterParallelByCount は指定された全件数に対して、引数で指定された処理を並列または直列で実行する関数です。
func IterParallelByCount(allCount int, blockSize int, processFunc func(int)) {
	if blockSize <= 1 || blockSize >= allCount {
		// ブロックサイズが1以下、もしくは全件数より大きい場合は直列処理
		for i := 0; i < allCount; i++ {
			processFunc(i)
		}
	} else {
		// ブロックサイズが全件数より小さい場合は並列処理
		var wg sync.WaitGroup
		for i := 0; i < allCount; i += blockSize {
			wg.Add(1)
			go func(startIndex int) {
				defer wg.Done()
				endIndex := startIndex + blockSize
				if endIndex > allCount {
					endIndex = allCount
				}
				for j := startIndex; j < endIndex; j++ {
					processFunc(j)
				}
			}(i)
		}
		wg.Wait() // すべてのgoroutineが終了するまで待機
	}
}

// IterParallelByList は指定された全リストに対して、引数で指定された処理を並列または直列で実行する関数です。
func IterParallelByList(allData []int, blockSize int, processFunc func(data, index int)) {
	if blockSize <= 1 || blockSize >= len(allData) {
		// ブロックサイズが1以下、もしくは全件数より大きい場合は直列処理
		for i := 0; i < len(allData); i++ {
			processFunc(allData[i], i)
		}
	} else {
		// ブロックサイズが全件数より小さい場合は並列処理
		var wg sync.WaitGroup
		for i := 0; i < len(allData); i += blockSize {
			wg.Add(1)
			go func(startIndex int) {
				defer wg.Done()
				endIndex := startIndex + blockSize
				if endIndex > len(allData) {
					endIndex = len(allData)
				}
				for j := startIndex; j < endIndex; j++ {
					processFunc(allData[j], j)
				}
			}(i)
		}
		wg.Wait() // すべてのgoroutineが終了するまで待機
	}
}

// CPUコア数を元に、ブロックサイズを計算
func GetBlockSize(totalTasks int) int {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// ブロックサイズを切り上げで計算
	return (totalTasks + numCPU - 1) / numCPU
}
