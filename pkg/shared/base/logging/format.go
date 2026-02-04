// 指示: miu200521358
package logging

import (
	"strings"
	"sync/atomic"
)

var debugStackEnabled uint32

// IsDebugEnabled はデバッグログが有効か判定する。
func IsDebugEnabled(logger ILogger) bool {
	if logger == nil {
		return false
	}
	level := logger.Level()
	if level == LOG_LEVEL_VERBOSE {
		level = LOG_LEVEL_DEBUG
	}
	return level <= LOG_LEVEL_DEBUG
}

// SetDebugStackEnabled はスタックトレースの出力可否を設定する。
func SetDebugStackEnabled(enabled bool) {
	if enabled {
		atomic.StoreUint32(&debugStackEnabled, 1)
		return
	}
	atomic.StoreUint32(&debugStackEnabled, 0)
}

// IsDebugStackEnabled はスタックトレースの出力が有効か判定する。
func IsDebugStackEnabled() bool {
	return atomic.LoadUint32(&debugStackEnabled) == 1
}

// TrimStackTrace はエラーメッセージ末尾のスタックトレースを除去する。
func TrimStackTrace(message string) string {
	if message == "" {
		return ""
	}
	if idx := strings.Index(message, "\r\n\r\nStack:\r\n"); idx >= 0 {
		return message[:idx]
	}
	if idx := strings.Index(message, "\n\nStack:\n"); idx >= 0 {
		return message[:idx]
	}
	if idx := strings.Index(message, "\n\nStack:\r\n"); idx >= 0 {
		return message[:idx]
	}
	return message
}

// FormatError はログ出力用のエラーメッセージを整形する。
func FormatError(err error, logger ILogger) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if !IsDebugStackEnabled() {
		message = TrimStackTrace(message)
	}
	return message
}
