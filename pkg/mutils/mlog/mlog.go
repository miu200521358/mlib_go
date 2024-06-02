package mlog

import (
	"log"
	"runtime"
)

var level = 20

const (
	VERBOSE    = 0
	IK_VERBOSE = 1
	VERBOSE2   = 2
	DEBUG      = 10
	INFO       = 20
	WARN       = 30
	ERROR      = 40
	FATAL      = 50
)

func init() {
	log.SetFlags(0)
}

func IsVerbose() bool {
	return level < DEBUG
}

func IsVerbose2() bool {
	return level == VERBOSE2
}

func IsIkVerbose() bool {
	return level == IK_VERBOSE
}

func IsDebug() bool {
	return level < INFO
}

// SetLevel ログレベルの設定
func SetLevel(l int) {
	level = l
}

// Verbose 冗長ログ
func V(message string, param ...interface{}) {
	if level < DEBUG {
		log.Printf(message, param...)
	}
}

// Verbose2 冗長ログ
func V2(message string, param ...interface{}) {
	if VERBOSE2 <= level && level < DEBUG {
		log.Printf(message, param...)
	}
}

// Debug デバッグログ
func D(message string, param ...interface{}) {
	if level < INFO {
		log.Printf(message, param...)
	}
}

// L ログの区切り線
func L() {
	log.Println("---------------------------------")
}

// Info 情報ログ
func I(message string, param ...interface{}) {
	log.Printf(message, param...)
}

// IL 情報ログ（区切り線付き）
func IL(message string, param ...interface{}) {
	L()
	I(message, param...)
}

// IT 情報ログ（タイトル付き）
func IT(title string, message string, param ...interface{}) {
	log.Printf("■■■■■ %s ■■■■■", title)
	I(message, param...)
}

// ILT 情報ログ（区切り線・タイトル付き）
func ILT(title string, message string, param ...interface{}) {
	L()
	IT(title, message, param...)
}

// Warn 警告ログ
func W(message string, param ...interface{}) {
	WT("WARN", message, param...)
}

// Warn 警告ログ
func WT(title string, message string, param ...interface{}) {
	log.Printf("~~~~~~~~~~ %s ~~~~~~~~~~", title)
	log.Printf(message, param...)
}

// Error エラーログ
func E(message string, param ...interface{}) {
	ET("ERROR", message, param...)
}

// Error エラーログ
func ET(title string, message string, param ...interface{}) {
	log.Printf("********** %s **********", title)
	log.Printf(message, param...)
}

// Fatal 致命的エラーログ
func F(message string, param ...interface{}) {
	FT("FATAL ERROR", message, param...)
}

// Error エラーログ
func FT(title string, message string, param ...interface{}) {
	log.Printf("!!!!!!!!!! %s !!!!!!!!!!", title)
	log.Printf(message, param...)
}

var prevMem uint64

func Memory(prefix string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mem := bToMb(m.Alloc)
	if prevMem != mem {
		log.Printf("[%s] Alloc = %v -> %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
			prefix, prevMem, mem, bToMb(m.TotalAlloc), m.NumGC, bToMb(m.Sys))
		prevMem = mem
	}
	// log.Printf("[%s] Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
	// 	prefix, bToMb(m.Alloc), bToMb(m.TotalAlloc), m.NumGC, bToMb(m.Sys))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
