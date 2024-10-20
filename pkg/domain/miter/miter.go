package miter

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// IterParallelByCount は指定された全件数に対して、引数で指定された処理を並列または直列で実行する関数です。
func IterParallelByCount(allCount int, blockSize int, processFunc func(int)) error {
	if blockSize <= 1 || blockSize >= allCount {
		// ブロックサイズが1以下、もしくは全件数より大きい場合は直列処理
		for i := 0; i < allCount; i++ {
			processFunc(i)
		}
	} else {
		numCPU := runtime.NumCPU()
		errorChan := make(chan error, numCPU)

		// ブロックサイズが全件数より小さい場合は並列処理
		var wg sync.WaitGroup
		for i := 0; i < allCount; i += blockSize {
			wg.Add(1)
			go func(startIndex int) {
				defer wg.Done()
				defer func() {
					// recoverによるpanicキャッチ
					if r := recover(); r != nil {
						stackTrace := debug.Stack()

						var errMsg string
						// パニックの値がerror型である場合、エラーメッセージを取得
						if err, ok := r.(error); ok {
							errMsg = err.Error()
						} else {
							// それ以外の型の場合は、文字列に変換
							errMsg = fmt.Sprintf("%v", r)
						}

						errorChan <- fmt.Errorf("panic: %s\n%s", errMsg, stackTrace)
					}
				}()

				endIndex := startIndex + blockSize
				if endIndex > allCount {
					endIndex = allCount
				}
				for j := startIndex; j < endIndex; j++ {
					processFunc(j)
				}
			}(i)
		}

		// すべてのゴルーチンの完了を待つ
		wg.Wait()
		close(errorChan) // 全てのゴルーチンが終了したらチャネルを閉じる

		// チャネルからエラーを受け取る
		for err := range errorChan {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IterParallelByList は指定された全リストに対して、引数で指定された処理を並列または直列で実行する関数です。
func IterParallelByList(allData []int, blockSize int, processFunc func(data, index int)) error {
	if blockSize <= 1 || blockSize >= len(allData) {
		// ブロックサイズが1以下、もしくは全件数より大きい場合は直列処理
		for i := 0; i < len(allData); i++ {
			processFunc(allData[i], i)
		}
	} else {
		numCPU := runtime.NumCPU()
		errorChan := make(chan error, numCPU)

		// ブロックサイズが全件数より小さい場合は並列処理
		var wg sync.WaitGroup
		for i := 0; i < len(allData); i += blockSize {
			wg.Add(1)
			go func(startIndex int) {
				defer wg.Done()
				defer func() {
					// recoverによるpanicキャッチ
					if r := recover(); r != nil {
						stackTrace := debug.Stack()

						var errMsg string
						// パニックの値がerror型である場合、エラーメッセージを取得
						if err, ok := r.(error); ok {
							errMsg = err.Error()
						} else {
							// それ以外の型の場合は、文字列に変換
							errMsg = fmt.Sprintf("%v", r)
						}

						errorChan <- fmt.Errorf("panic: %s\n%s", errMsg, stackTrace)
					}
				}()

				endIndex := startIndex + blockSize
				if endIndex > len(allData) {
					endIndex = len(allData)
				}
				for j := startIndex; j < endIndex; j++ {
					processFunc(allData[j], j)
				}
			}(i)
		}

		// すべてのゴルーチンの完了を待つ
		wg.Wait()
		close(errorChan) // 全てのゴルーチンが終了したらチャネルを閉じる

		// チャネルからエラーを受け取る
		for err := range errorChan {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CPUコア数を元に、ブロックサイズを計算
func GetBlockSize(totalTasks int) int {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// ブロックサイズを切り上げで計算
	return (totalTasks + numCPU - 1) / numCPU
}
