package mproc

import "runtime"

func SetMaxProcess(isFull bool) {
	cpuNum := runtime.NumCPU()
	if !isFull {
		// FULL稼働ではない時にはシステム上の1/4の論理プロセッサを使用させる
		cpuNum = max(1, int(runtime.NumCPU()/4))
	}

	runtime.GOMAXPROCS(cpuNum)
}
