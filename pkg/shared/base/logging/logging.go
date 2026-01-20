// 指示: miu200521358
package logging

import "io"

// LogLevel はログレベルを表す。
type LogLevel int

const (
	// LOG_LEVEL_VERBOSE は冗長ログを表す。
	LOG_LEVEL_VERBOSE LogLevel = 0
	// LOG_LEVEL_DEBUG はデバッグログを表す。
	LOG_LEVEL_DEBUG LogLevel = 10
	// LOG_LEVEL_INFO は情報ログを表す。
	LOG_LEVEL_INFO LogLevel = 20
	// LOG_LEVEL_WARN は警告ログを表す。
	LOG_LEVEL_WARN LogLevel = 30
	// LOG_LEVEL_ERROR はエラーログを表す。
	LOG_LEVEL_ERROR LogLevel = 40
	// LOG_LEVEL_FATAL は致命ログを表す。
	LOG_LEVEL_FATAL LogLevel = 50
)

// VerboseIndex は冗長ログの対象種別。
type VerboseIndex int

const (
	// VERBOSE_INDEX_MOTION はモーション冗長ログ。
	VERBOSE_INDEX_MOTION VerboseIndex = 230
	// VERBOSE_INDEX_IK はIK冗長ログ。
	VERBOSE_INDEX_IK VerboseIndex = 320
	// VERBOSE_INDEX_PHYSICS は物理冗長ログ。
	VERBOSE_INDEX_PHYSICS VerboseIndex = 330
	// VERBOSE_INDEX_VIEWER はビューワー冗長ログ。
	VERBOSE_INDEX_VIEWER VerboseIndex = 540
)

// IMessageBuffer はメッセージ欄バッファI/F。
type IMessageBuffer interface {
	Text() string
	Lines() []string
	Clear()
}

// IVerboseSink は冗長ログ出力のI/F。
type IVerboseSink interface {
	WriteLine(text string)
	Close() error
}

// ILogger はログ出力I/F。
type ILogger interface {
	Level() LogLevel
	SetLevel(level LogLevel)
	EnableVerbose(idx VerboseIndex)
	DisableVerbose(idx VerboseIndex)
	IsVerboseEnabled(idx VerboseIndex) bool
	MessageBuffer() IMessageBuffer
	AttachVerboseSink(idx VerboseIndex, sink IVerboseSink)
	Verbose(idx VerboseIndex, msg string, params ...any)
	Debug(msg string, params ...any)
	Info(msg string, params ...any)
	Warn(msg string, params ...any)
	Error(msg string, params ...any)
	Fatal(msg string, params ...any)
}

// IConsoleLogger はコンソール出力先を扱えるロガーI/F。
type IConsoleLogger interface {
	ILogger
	SetConsoleSink(writer io.Writer)
	ConsoleText() string
}

// defaultLogger は共有の既定ロガー。
var defaultLogger ILogger = &noopLogger{}

// SetDefaultLogger は既定ロガーを設定する。
func SetDefaultLogger(logger ILogger) {
	if logger == nil {
		return
	}
	defaultLogger = logger
}

// DefaultLogger は既定ロガーを返す。
func DefaultLogger() ILogger {
	return defaultLogger
}

// SetConsoleSink は既定ロガーのコンソール出力先を設定する。
func SetConsoleSink(writer io.Writer) {
	logger, ok := defaultLogger.(IConsoleLogger)
	if !ok {
		return
	}
	logger.SetConsoleSink(writer)
}

// ConsoleText は既定ロガーのメッセージ欄全文を返す。
func ConsoleText() string {
	logger, ok := defaultLogger.(IConsoleLogger)
	if !ok {
		return ""
	}
	return logger.ConsoleText()
}

// noopLogger は既定ロガー未設定時のダミー実装。
type noopLogger struct{}

// Level はログレベルを返す。
func (l *noopLogger) Level() LogLevel {
	return LOG_LEVEL_INFO
}

// SetLevel はログレベルを設定する。
func (l *noopLogger) SetLevel(level LogLevel) {}

// EnableVerbose は冗長ログを有効化する。
func (l *noopLogger) EnableVerbose(idx VerboseIndex) {}

// DisableVerbose は冗長ログを無効化する。
func (l *noopLogger) DisableVerbose(idx VerboseIndex) {}

// IsVerboseEnabled は冗長ログの有効可否を返す。
func (l *noopLogger) IsVerboseEnabled(idx VerboseIndex) bool {
	return false
}

// MessageBuffer はメッセージ欄バッファを返す。
func (l *noopLogger) MessageBuffer() IMessageBuffer {
	return noopBuffer{}
}

// AttachVerboseSink は冗長ログの出力先を設定する。
func (l *noopLogger) AttachVerboseSink(idx VerboseIndex, sink IVerboseSink) {}

// Verbose は冗長ログを出力する。
func (l *noopLogger) Verbose(idx VerboseIndex, msg string, params ...any) {}

// Debug はデバッグログを出力する。
func (l *noopLogger) Debug(msg string, params ...any) {}

// Info は情報ログを出力する。
func (l *noopLogger) Info(msg string, params ...any) {}

// Warn は警告ログを出力する。
func (l *noopLogger) Warn(msg string, params ...any) {}

// Error はエラーログを出力する。
func (l *noopLogger) Error(msg string, params ...any) {}

// Fatal は致命ログを出力する。
func (l *noopLogger) Fatal(msg string, params ...any) {}

// noopBuffer は空のメッセージバッファ。
type noopBuffer struct{}

// Text は全文を返す。
func (b noopBuffer) Text() string {
	return ""
}

// Lines は行配列を返す。
func (b noopBuffer) Lines() []string {
	return []string{}
}

// Clear はバッファをクリアする。
func (b noopBuffer) Clear() {}
