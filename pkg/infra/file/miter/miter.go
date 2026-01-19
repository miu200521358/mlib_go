// 指示: miu200521358
package miter

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// IterParallelByList はリストに対して並列/直列で処理を実行する。
func IterParallelByList[T any](allData []T, blockSize int, logBlockSize int,
	processFunc func(index int, data T) error, logFunc func(iterIndex, allCount int)) error {
	if len(allData) == 0 {
		return nil
	}
	if blockSize <= 0 {
		blockSize, _ = GetBlockSize(len(allData))
	}
	if blockSize >= len(allData) {
		if err := iterSerial(allData, processFunc, logBlockSize, logFunc); err != nil {
			return err
		}
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workers := runtime.NumCPU()
	errChan := make(chan error, workers)
	var wg sync.WaitGroup
	var mu sync.Mutex
	iterIndex := 0

	for startIndex := 0; startIndex < len(allData); startIndex += blockSize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					stackTrace := debug.Stack()
					errMsg := fmt.Sprintf("%v", r)
					if e, ok := r.(error); ok {
						errMsg = e.Error()
					}
					errChan <- newIterProcessFailed("panicを検知しました", fmt.Errorf("panic: %s\n%s", errMsg, stackTrace))
					cancel()
				}
			}()

			select {
			case <-ctx.Done():
				return
			default:
			}

			end := start + blockSize
			if end > len(allData) {
				end = len(allData)
			}
			for i := start; i < end; i++ {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if err := processFunc(i, allData[i]); err != nil {
					errChan <- newIterProcessFailed("処理中にエラーが発生しました", err)
					cancel()
					return
				}
				if logFunc != nil && logBlockSize > 0 {
					mu.Lock()
					if iterIndex%logBlockSize == 0 && iterIndex > 0 {
						logFunc(iterIndex, len(allData))
					}
					iterIndex++
					mu.Unlock()
				}
			}
		}(startIndex)
	}

	var firstErr error
	go func() {
		for err := range errChan {
			if firstErr == nil {
				firstErr = err
			}
			cancel()
		}
	}()

	wg.Wait()
	close(errChan)

	if firstErr != nil {
		return firstErr
	}
	return nil
}

// GetBlockSize はCPU数からブロックサイズを算出する。
func GetBlockSize(totalTasks int) (blockSize int, blockCount int) {
	blockCount = runtime.NumCPU()
	if blockCount <= 0 {
		blockCount = 1
	}
	if totalTasks <= 0 {
		return 1, blockCount
	}
	blockSize = (totalTasks + blockCount - 1) / blockCount
	if blockSize <= 0 {
		blockSize = 1
	}
	return blockSize, blockCount
}

func iterSerial[T any](allData []T, processFunc func(index int, data T) error, logBlockSize int, logFunc func(iterIndex, allCount int)) error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				errMsg := fmt.Sprintf("%v", r)
				if e, ok := r.(error); ok {
					errMsg = e.Error()
				}
				err = newIterProcessFailed("panicを検知しました", fmt.Errorf("panic: %s\n%s", errMsg, stackTrace))
			}
		}()
		for i := range allData {
			if processErr := processFunc(i, allData[i]); processErr != nil {
				err = newIterProcessFailed("処理中にエラーが発生しました", processErr)
				return
			}
			if logFunc != nil && logBlockSize > 0 {
				if i%logBlockSize == 0 && i > 0 {
					logFunc(i, len(allData))
				}
			}
		}
	}()
	return err
}
