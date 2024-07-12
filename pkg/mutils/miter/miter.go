package miter

import "sync"

var blockSize = 100

func GetBlockSize() int {
	return blockSize
}

// SetBlockSize 並列処理のブロックサイズ
func SetBlockSize(size int) {
	blockSize = size
}

// IterParallel は指定された全件数に対して、引数で指定された処理を並列または直列で実行する関数です。
// allCount: ループしたい全件数
// blockSize: 並列処理を行うブロックのサイズ
// processFunc: 実行したい処理を定義した関数
func IterParallel(allCount int, processFunc func(int)) {
	if blockSize <= 1 {
		// ブロックサイズが1以下の場合は直列処理
		for i := 0; i < allCount; i++ {
			processFunc(i)
		}
	} else {
		// ブロックサイズが1より大きい場合は並列処理
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
