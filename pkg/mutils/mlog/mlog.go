package mlog

import "log"

var level = 20

// SetLevel ログレベルの設定
func SetLevel(l int) {
	level = l
}

// Verbose 冗長ログ
func V(message string, param ...interface{}) {
	if level <= 0 {
		log.Printf(message, param...)
	}
}

// Debug デバッグログ
func D(message string, param ...interface{}) {
	if level <= 10 {
		log.Printf(message, param...)
	}
}

// Info 情報ログ
func I(message string, param ...interface{}) {
	if level <= 20 {
		log.Printf(message, param...)
	}
}

// Warn 警告ログ
func W(message string, param ...interface{}) {
	if level <= 30 {
		log.Printf(message, param...)
	}
}

// Error エラーログ
func E(message string, param ...interface{}) {
	if level <= 40 {
		log.Printf(message, param...)
	}
}

// Fatal 致命的エラーログ
func F(message string, param ...interface{}) {
	if level <= 50 {
		log.Fatalf(message, param...)
	}
}
