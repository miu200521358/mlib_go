package mlog

import (
	"bytes"
	"log"
	"runtime"
	"time"
)

var level = 20.0

const (
	VERBOSE        = 0.0
	IK_VERBOSE     = 1.0
	VIEWER_VERBOSE = 1.1
	VERBOSE2       = 2.0
	VERBOSE9       = 9.0
	DEBUG          = 10.0
	INFO           = 20.0
	WARN           = 30.0
	ERROR          = 40.0
	FATAL          = 50.0
)

func init() {
	log.SetFlags(0)
}

func IsVerbose() bool {
	return level < DEBUG && level != IK_VERBOSE && level != VIEWER_VERBOSE
}

func IsVerbose2() bool {
	return level == VERBOSE2
}

func IsIkVerbose() bool {
	return level == IK_VERBOSE
}

func IsViewerVerbose() bool {
	return level == VIEWER_VERBOSE
}

func IsVerbose9() bool {
	return level < DEBUG
}

func IsDebug() bool {
	return level < INFO
}

// SetLevel ログレベルの設定
func SetLevel(l float64) {
	level = float64(l)
}

// Verbose 冗長ログ
func V(message string, param ...any) {
	if level < DEBUG {
		log.Printf(message, param...)
	}
}

// Verbose2 冗長ログ
func V2(message string, param ...any) {
	if VERBOSE2 <= level && level < DEBUG {
		log.Printf(message, param...)
	}
}

// Ik Verbose2 冗長ログ
func IV(message string, param ...any) {
	if level == IK_VERBOSE {
		log.Printf(message, param...)
	}
}

// Viewer Verbose 冗長ログ
func VV(message string, param ...any) {
	if level == VIEWER_VERBOSE {
		log.Printf(message, param...)
	}
}

// Debug デバッグログ
func D(message string, param ...any) {
	if level < INFO {
		log.Printf(message, param...)
	}
}

// L ログの区切り線
func L() {
	log.Println("---------------------------------")
}

// Info 情報ログ
func I(message string, param ...any) {
	log.Printf(message, param...)
}

// IL 情報ログ（区切り線付き）
func IL(message string, param ...any) {
	L()
	I(message, param...)
}

// IS 情報ログ（タイムスタンプ付き）
func IS(message string, param ...any) {
	I("[%s]"+message, append([]any{time.Now().Format("15:04:05.999999999")}, param...)...)
}

// IT 情報ログ（タイトル付き）
func IT(title string, message string, param ...any) {
	log.Printf("■■■■■ %s ■■■■■", title)
	I(message, param...)
}

// ILT 情報ログ（区切り線・タイトル付き）
func ILT(title string, message string, param ...any) {
	L()
	IT(title, message, param...)
}

// Warn 警告ログ
func W(message string, param ...any) {
	WT("WARN", message, param...)
}

// Warn 警告ログ
func WT(title string, message string, param ...any) {
	log.Printf("~~~~~~~~~~ %s ~~~~~~~~~~", title)
	log.Printf(message, param...)
}

// Error エラーログ
func E(message string, err error, param ...any) {
	ET("ERROR", err, message, param...)
}

// Error エラーログ
func ET(title string, err error, message string, param ...any) {
	log.Printf("********** %s **********", title)
	if err != nil {
		errMsg := err.Error()
		log.Printf("Error Message: %s", errMsg)
		log.Printf("Stack Trace:\n%s", dumpAllGoroutines())
	}
	if message != "" {
		log.Printf(message, param...)
	}
}

// Fatal 致命的エラーログ
func F(message string, err error, param ...any) {
	FT("FATAL ERROR", err, message, param...)
}

// Error エラーログ
func FT(title string, err error, message string, param ...any) {
	log.Printf("!!!!!!!!!! %s !!!!!!!!!!", title)
	if err != nil {
		errMsg := err.Error()
		log.Printf("Error Message: %s", errMsg)
		log.Printf("Stack Trace:\n%s", dumpAllGoroutines())
	}
	if message != "" {
		log.Printf(message, param...)
	}
}

func dumpAllGoroutines() string {
	buf := make([]byte, 1<<20) // 1MB のバッファを確保
	n := runtime.Stack(buf, true)
	return string(bytes.ReplaceAll(buf[:n], []byte("\n"), []byte("\r\n")))
}

var prevMem uint64

func Memory(prefix string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mem := bToMb(m.Alloc)
	if prevMem != mem {
		log.Printf("[%s] Alloc = %v -> %v MiB, HeapAlloc = %v MiB, HeapSys = %v MiB, HeapIdle = %v MiB, "+
			"HeapInuse = %v MiB, HeapReleased = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
			prefix, prevMem, mem, bToMb(m.HeapAlloc), bToMb(m.HeapSys), bToMb(m.HeapIdle),
			bToMb(m.HeapInuse), bToMb(m.HeapReleased), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
		prevMem = mem
	}
	// log.Printf("[%s] Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
	// 	prefix, bToMb(m.Alloc), bToMb(m.TotalAlloc), m.NumGC, bToMb(m.Sys))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
