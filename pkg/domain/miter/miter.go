package miter

import "sync"

// IterParallel は指定された全件数に対して、引数で指定された処理を並列または直列で実行する関数です。
// allCount: ループしたい全件数
// blockSize: 並列処理を行うブロックのサイズ
// processFunc: 実行したい処理を定義した関数
func IterParallel(allCount int, blockSize int, processFunc func(int)) {
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
