package miter

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/config/mproc"
)

// IterParallelByList は指定された全リストに対して、引数で指定された処理を並列または直列で実行する関数です。
func IterParallelByList[T any](allData []T, blockSize int, logBlockSize int,
	processFunc func(index int, data T) error, logFunc func(iterIndex, allCount int)) error {
	numCPU := runtime.NumCPU()

	mproc.SetMaxProcess(true)
	defer mproc.SetMaxProcess(false)

	if blockSize >= len(allData) {
		// ブロックサイズが全件数より大きい場合は直列処理
		// パニックをキャッチするためにdeferとrecoverを使用
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					stackTrace := debug.Stack()
					errMsg := fmt.Sprintf("%v", r)
					if e, ok := r.(error); ok {
						errMsg = e.Error()
					}
					err = fmt.Errorf("panic: %s\n%s", errMsg, stackTrace)
				}
			}()

			for i := range allData {
				if processErr := processFunc(i, allData[i]); processErr != nil {
					err = processErr
					return
				}
			}
		}()

		if err != nil {
			return err
		}
	} else {
		// ブロックサイズが全件数より小さい場合は並列処理
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		errChan := make(chan error, numCPU)
		var wg sync.WaitGroup
		var mu sync.Mutex
		iterIndex := 0

		for startIndex := 0; startIndex < len(allData); startIndex += blockSize {
			wg.Add(1)
			go func(startIndex int) {
				// 最初にパニックハンドリングのためのdeferを設置
				defer func() {
					if r := recover(); r != nil {
						stackTrace := debug.Stack()
						errMsg := fmt.Sprintf("%v", r)
						if e, ok := r.(error); ok {
							errMsg = e.Error()
						}
						errChan <- fmt.Errorf("panic: %s\n%s", errMsg, stackTrace)
						cancel() // 他のゴルーチンも停止させる
					}
					wg.Done()
				}()

				// コンテキストがキャンセルされていないか確認
				select {
				case <-ctx.Done():
					return
				default:
					// 処理続行
				}

				endIndex := startIndex + blockSize
				if endIndex > len(allData) {
					endIndex = len(allData)
				}

				for j := startIndex; j < endIndex; j++ {
					// 定期的にコンテキストをチェック
					select {
					case <-ctx.Done():
						return
					default:
						// 処理続行
					}

					if err := processFunc(j, allData[j]); err != nil {
						errChan <- err
						cancel() // エラー発生時は他の処理も中断
						return
					}

					if logFunc != nil && logBlockSize > 0 {
						mu.Lock()
						if iterIndex%logBlockSize == 0 {
							logFunc(iterIndex, len(allData))
						}
						iterIndex++
						mu.Unlock()
					}
				}
			}(startIndex)
		}

		// エラー収集用ゴルーチン
		var firstErr error
		go func() {
			for err := range errChan {
				if firstErr == nil {
					firstErr = err
				}
				cancel() // 何かエラーが発生したら即座に他の処理をキャンセル
			}
		}()

		wg.Wait()
		close(errChan)

		if firstErr != nil {
			return firstErr
		}
	}

	return nil
}

// CPUコア数を元に、ブロックサイズを計算
func GetBlockSize(totalTasks int) (blockSize int, blockCount int) {
	blockCount = runtime.NumCPU()

	// ブロックサイズを切り上げで計算
	blockSize = max(1, (totalTasks+blockCount-1)/blockCount)

	return blockSize, blockCount
}
